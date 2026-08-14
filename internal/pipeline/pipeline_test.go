package pipeline

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/selection"
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
		"rails/actions/actions.yaml": `actions:
  addField:
    kwargs:
      model: { type: string, required: true }
      field: { type: string, required: true }
    injects_into: "app/models/{{model|snake}}.rb"
    anchor: class_body
`,
		"rails/actions/addField.rb": "attribute :{{field}}\n",
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

// stub is the model. Phase 4 is the one phase that is not a pure function of
// its inputs, so the pipeline's own tests supply a canned answer: what is being
// tested here is that each phase's output reaches the next, not what a model
// would have chosen.
type stub struct {
	response string
	calls    int
	prompts  [][]selection.Message
}

func (s *stub) Complete(_ context.Context, messages []selection.Message) (string, error) {
	s.calls++
	s.prompts = append(s.prompts, messages)
	if s.response == "" {
		return `{"invocations": []}`, nil
	}
	return s.response, nil
}

func config(t *testing.T) Config {
	t.Helper()
	return Config{
		Generators: generators(t),
		Records: recordsDir(t, map[string]string{
			"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		}),
		Output: t.TempDir(),
		Client: &stub{},
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

	result, err := Run(context.Background(), cfg)
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

	result, err := Run(context.Background(), cfg)
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

	result, err := Run(context.Background(), cfg)
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
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("first run: %v", err)
	}
	before := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))

	cfg.StopAfterPhase = 0
	result, err := Run(context.Background(), cfg)
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

	result, err := Run(context.Background(), cfg)
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

	if _, err := Run(context.Background(), cfg); err != nil {
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

	_, err := Run(context.Background(), cfg)
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

	_, err := Run(context.Background(), cfg)
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

	result, err := Run(context.Background(), cfg)
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
		if _, err := Run(context.Background(), cfg); err == nil || !strings.Contains(err.Error(), "absent") {
			t.Errorf("error = %v, does not name the missing generators directory", err)
		}
	})
	t.Run("records", func(t *testing.T) {
		cfg := config(t)
		cfg.Records = filepath.Join(t.TempDir(), "absent")
		if _, err := Run(context.Background(), cfg); err == nil || !strings.Contains(err.Error(), "absent") {
			t.Errorf("error = %v, does not name the missing records directory", err)
		}
	})
}

const selectsAField = `{"invocations":[{"action":"addField","kwargs":{"model":"user","field":"email"}}]}`

// The whole loop, with no hand-written invocation list: a record's intent goes
// to a model, what comes back is validated, expanded, and injected, and the
// file on disk carries the region with its ownership marker.
func TestARecordIsGeneratedFromItsIntentAlone(t *testing.T) {
	cfg := config(t)
	client := &stub{response: selectsAField}
	cfg.Client = client

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if client.calls != 1 {
		t.Errorf("the model was called %d times for one record, want 1", client.calls)
	}
	if len(result.Selections) != 1 || len(result.Selections[0].Invocations) != 1 {
		t.Fatalf("Phase 5 produced %+v", result.Selections)
	}
	if len(result.Injections) != 1 {
		t.Fatalf("Phase 7 applied %d regions, want 1", len(result.Injections))
	}

	body := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))
	for _, want := range []string{"attribute :email", "sedum:addField", "/sedum:addField", `"record":"prov-2026-aaaaaaaa"`} {
		if !strings.Contains(body, want) {
			t.Errorf("the generated file does not carry %q:\n%s", want, body)
		}
	}
}

// Rerunning replaces the region an action owns rather than appending beside it.
// The recording and the sidecar cache both exist to be unnecessary, and this is
// why: idempotency state lives in the file.
func TestRerunningReplacesRatherThanDuplicates(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: selectsAField}

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("first run: %v", err)
	}
	first := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("second run: %v", err)
	}
	second := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))

	if first != second {
		t.Errorf("a rerun changed the file:\n%s\n---\n%s", first, second)
	}
	if got := strings.Count(second, "attribute :email"); got != 1 {
		t.Errorf("the region appears %d times after two runs, want 1:\n%s", got, second)
	}
}

// A dry run decides everything and writes nothing, including into files it
// declined to create. Phase 7 takes what Phase 3 rendered for those
// (prov-2026-23653fdc).
func TestDryRunReportsInjectionsWithoutWriting(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: selectsAField}
	cfg.DryRun = true

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Injections) != 1 {
		t.Fatalf("a dry run reported %d injections, want 1", len(result.Injections))
	}
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Errorf("a dry run wrote %v", got)
	}
}

// A dry run against a tree where the file already exists reads it from disk, so
// what it reports is a replacement rather than a first injection. Preferring
// the rendered template would describe a run nobody asked for.
func TestDryRunReadsFilesThatAreAlreadyThere(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: selectsAField}
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("first run: %v", err)
	}
	before := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))

	cfg.DryRun = true
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("dry run: %v", err)
	}

	if len(result.Injections) != 1 || !result.Injections[0].Replaced {
		t.Errorf("a dry run over an existing region reported %+v, want a replacement", result.Injections)
	}
	if after := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb")); after != before {
		t.Errorf("a dry run modified the file")
	}
}

// An action resolving to a path no record authorized fails validation, before a
// retry is spent and before anything is injected.
func TestAnUnauthorizedTargetHaltsWithNothingInjected(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: `{"invocations":[{"action":"addField","kwargs":{"model":"order","field":"total"}}]}`}

	_, err := Run(context.Background(), cfg)
	if err == nil {
		t.Fatal("an action targeting an unauthorized path was accepted")
	}
	if !strings.Contains(err.Error(), "app/models/order.rb") {
		t.Errorf("the error does not name the unauthorized path: %v", err)
	}

	// The authorized file was still created by Phase 3, and carries no region.
	body := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))
	if strings.Contains(body, "sedum:addField") {
		t.Errorf("a halted run still injected:\n%s", body)
	}
}

// Stopping after validation leaves the files scaffolded and nothing injected.
// It is the point a person reads what the model decided and corrects it.
func TestStopAfterInvocationsLeavesFilesUninjected(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: selectsAField}
	cfg.StopAfterPhase = PhaseValidate

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Selections) != 1 || len(result.Selections[0].Invocations) != 1 {
		t.Fatalf("the validated invocations are not available: %+v", result.Selections)
	}
	if result.Injections != nil {
		t.Errorf("stopping after validation still injected %+v", result.Injections)
	}

	body := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))
	if strings.Contains(body, "sedum:addField") {
		t.Errorf("the file was injected into after a stop at validation:\n%s", body)
	}
}

// A record whose every path is unmanaged has an empty catalog, so the only
// valid answer is an empty list. Paying for a call to be told so is a cost the
// run can see is pointless before it is incurred.
func TestARecordWithNothingToWriteCostsNoModelCall(t *testing.T) {
	cfg := config(t)
	cfg.Generators = writeTree(t, map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
unmanaged:
  - Gemfile
`,
		"rails/actions/actions.yaml": "actions: {}\n",
	})
	cfg.Records = recordsDir(t, map[string]string{"prov-2026-cccccccc": "  - Gemfile\n"})
	client := &stub{}
	cfg.Client = client

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if client.calls != 0 {
		t.Errorf("a record with no writable path cost %d model calls", client.calls)
	}
	if len(result.Unmanaged) != 1 {
		t.Errorf("the unmanaged path was not reported as the run's handoff: %+v", result.Unmanaged)
	}
}

// The prompt the model receives is built from the record it is being asked
// about and the packages that record's paths resolved to, and nothing else.
func TestEachRecordGetsItsOwnCall(t *testing.T) {
	cfg := config(t)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		"prov-2026-bbbbbbbb": "  - app/models/order.rb\n",
	})
	client := &stub{}
	cfg.Client = client

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if client.calls != 2 {
		t.Fatalf("two records cost %d model calls, want 2", client.calls)
	}
	// The catalog names paths too, in its injects_into patterns, so the file
	// list is what is counted rather than the prompt as a whole.
	for i, prompt := range client.prompts {
		user := prompt[len(prompt)-1].Content
		files, _, _ := strings.Cut(user, "## Actions")
		if strings.Count(files, "app/models/") != 1 {
			t.Errorf("call %d saw more than its own record's files:\n%s", i+1, files)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
