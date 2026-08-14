package pipeline

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// The pipeline's own job is ordering: each phase's output is the next phase's
// only input, and a stop point halts at a boundary with everything before it
// complete and nothing after it started. The phases themselves are tested where
// they live.

func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
	return root
}

func generators(t *testing.T) string {
	t.Helper()
	return writeTree(t, map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
`,
		"rails/files/app/models/{name}.rb": "class {{name|constantize}}\n  # sedum:anchor:class_body\nend\n",
		"rails/actions/actions.yaml":       "actions: {}\n",
	})
}

func recordsDir(t *testing.T, records map[string]string) string {
	t.Helper()
	files := map[string]string{}
	for id, scope := range records {
		files[id+".yml"] = "id: " + id + "\nintent: |\n  Something.\naffected_scope:\n" + scope
	}
	return writeTree(t, files)
}

func config(t *testing.T) Config {
	t.Helper()
	return Config{
		Generators: generators(t),
		Records: recordsDir(t, map[string]string{
			"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		}),
		Output: t.TempDir(),
	}
}

// tree returns every file under root, slash-separated and sorted.
func tree(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		out = append(out, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	sort.Strings(out)
	return out
}

func TestRunCarriesEveryPhaseThrough(t *testing.T) {
	cfg := config(t)

	result, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Packages.Packages) != 1 {
		t.Errorf("Phase 0 loaded %d packages, want 1", len(result.Packages.Packages))
	}
	if len(result.Records.Records) != 1 {
		t.Errorf("Phase 1 ingested %d records, want 1", len(result.Records.Records))
	}
	if len(result.Resolutions) != 1 {
		t.Errorf("Phase 2 resolved %d paths, want 1", len(result.Resolutions))
	}
	if len(result.Files) != 1 {
		t.Errorf("Phase 3 produced %d files, want 1", len(result.Files))
	}
	if got := tree(t, cfg.Output); len(got) != 1 || got[0] != "app/models/user.rb" {
		t.Errorf("output tree = %v", got)
	}
}

// Stopping after resolution leaves everything Phase 2 decided available and
// nothing on disk, which is what makes the stop point worth having
// (prov-2026-5696ff65).
func TestStopAfterResolutionDecidesEverythingAndWritesNothing(t *testing.T) {
	cfg := config(t)
	cfg.StopAfterPhase = PhaseResolve

	result, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Resolutions) != 1 {
		t.Fatalf("Phase 2 resolved %d paths, want 1", len(result.Resolutions))
	}
	res := result.Resolutions[0]
	if res.Package == nil || res.Package.Name != "rails" {
		t.Errorf("resolution has no package")
	}
	if res.Template != "app/models/{name}.rb" || res.Captures["name"] != "user" {
		t.Errorf("resolution = template %q captures %v; the matched template and its captures must be available here",
			res.Template, res.Captures)
	}
	if result.Files != nil {
		t.Errorf("Phase 3 ran despite a stop at resolution")
	}
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Errorf("stopping after resolution wrote %v", got)
	}
}

func TestStopAfterFilesCreatesFilesAndStops(t *testing.T) {
	cfg := config(t)
	cfg.StopAfterPhase = PhaseCreate

	result, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("Phase 3 produced %d files, want 1", len(result.Files))
	}
	if got := tree(t, cfg.Output); len(got) != 1 {
		t.Errorf("output tree = %v, want the one authorized file", got)
	}
}

// Resuming after either stop point is an ordinary rerun: Phases 0 through 2 are
// pure and Phase 3 is create-if-absent, so nothing needs to be preserved
// between the stop and the resume.
func TestResumingAfterAStopIsAnOrdinaryRerun(t *testing.T) {
	cfg := config(t)
	cfg.StopAfterPhase = PhaseCreate
	if _, err := Run(cfg); err != nil {
		t.Fatalf("first run: %v", err)
	}
	before := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))

	cfg.StopAfterPhase = 0
	result, err := Run(cfg)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}

	if !result.Files[0].Existed {
		t.Errorf("the resumed run treated an existing file as new")
	}
	if after := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb")); after != before {
		t.Errorf("the resumed run rewrote the file")
	}
}

func TestDryRunWritesNothing(t *testing.T) {
	cfg := config(t)
	cfg.DryRun = true

	result, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("a dry run reported %d files, want 1", len(result.Files))
	}
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Errorf("a dry run wrote %v", got)
	}
}

func TestOnlyRestrictsTheRunToNamedRecords(t *testing.T) {
	cfg := config(t)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		"prov-2026-bbbbbbbb": "  - app/models/order.rb\n",
	})
	cfg.Only = []string{"prov-2026-bbbbbbbb"}

	if _, err := Run(cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}

	want := []string{"app/models/order.rb"}
	if got := tree(t, cfg.Output); len(got) != 1 || got[0] != want[0] {
		t.Errorf("output tree = %v, want %v", got, want)
	}
}

// A phase that fails stops the run there. Phase 0 rejecting a package must not
// leave Phase 3 creating files under a half-loaded generators directory.
func TestARejectedPackageHaltsBeforeAnythingIsWritten(t *testing.T) {
	cfg := config(t)
	cfg.Generators = writeTree(t, map[string]string{
		"rails/sedum.yaml":           "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n",
		"rails/actions/actions.yaml": "actions:\n  broken:\n    anchor: nowhere\n",
	})

	_, err := Run(cfg)
	if err == nil {
		t.Fatal("a package with no template for a declared action was accepted")
	}
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Errorf("a failed load still wrote %v", got)
	}
}

func TestAnUnresolvablePathHaltsBeforeAnythingIsWritten(t *testing.T) {
	cfg := config(t)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - app/models/user.rb\n  - app/views/index.erb\n",
	})

	_, err := Run(cfg)
	if err == nil {
		t.Fatal("a path whose extension no package claims was accepted")
	}
	if !strings.Contains(err.Error(), ".erb") {
		t.Errorf("error = %q, does not name the unclaimed extension", err)
	}
	// Phase 2 resolves the whole path set before Phase 3 creates anything, so
	// one unresolvable path means no file is created at all rather than some.
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Errorf("a failed resolution still wrote %v", got)
	}
}

// Warnings are collected from every phase and handed back together, since the
// command is what decides where they go.
func TestWarningsAreCollected(t *testing.T) {
	cfg := config(t)
	cfg.Lang = []string{"nosuchpackage"}

	result, err := Run(cfg)
	if err != nil {
		t.Fatalf("an unusable --lang refused the run: %v", err)
	}

	var found bool
	for _, w := range result.Warnings {
		if strings.Contains(w, "nosuchpackage") {
			found = true
		}
	}
	if !found {
		t.Errorf("warnings = %v, do not mention the unusable --lang", result.Warnings)
	}
}

func TestMissingDirectoriesAreReported(t *testing.T) {
	t.Run("generators", func(t *testing.T) {
		cfg := config(t)
		cfg.Generators = filepath.Join(t.TempDir(), "absent")
		if _, err := Run(cfg); err == nil || !strings.Contains(err.Error(), "absent") {
			t.Errorf("error = %v, does not name the missing generators directory", err)
		}
	})
	t.Run("records", func(t *testing.T) {
		cfg := config(t)
		cfg.Records = filepath.Join(t.TempDir(), "absent")
		if _, err := Run(cfg); err == nil || !strings.Contains(err.Error(), "absent") {
			t.Errorf("error = %v, does not name the missing records directory", err)
		}
	})
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
