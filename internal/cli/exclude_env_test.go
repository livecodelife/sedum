package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// SEDUM_EXCLUDE adds directories to the ones the marker walk skips, for
// `sedum regions` and `sedum grow` (prov-2026-fc05416c). It is an environment
// variable and not a flag because the flag surfaces are exact contract tests.

const excludeEnvMarker = "# sedum:addField {\"tier\":\"owned\",\"record\":\"prov-2026-aaaaaaaa\"}\n# /sedum:addField\n"

func markedRubyFile(t *testing.T, output, rel string) {
	t.Helper()
	path := filepath.Join(output, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(excludeEnvMarker), 0o644); err != nil {
		t.Fatal(err)
	}
}

func regionsOutput(t *testing.T, output string) string {
	t.Helper()
	gens := filepath.Join("..", "..", "testdata", "generators")
	out, err := exec(t, "regions", "--generators", gens, "--output", output)
	if err != nil {
		t.Fatalf("regions: %v\n%s", err, out)
	}
	return out
}

func TestRegionsSkipsVendoredAndSedumExcludeDirectories(t *testing.T) {
	output := t.TempDir()
	markedRubyFile(t, output, "app/keep.rb")
	markedRubyFile(t, output, "node_modules/pkg/vendored.rb")
	markedRubyFile(t, output, "scratch/a.rb")
	markedRubyFile(t, output, "docs/gen/b.rb")

	t.Setenv("SEDUM_EXCLUDE", "")
	out := regionsOutput(t, output)
	if !strings.Contains(out, "app/keep.rb") || strings.Contains(out, "node_modules") {
		t.Errorf("without SEDUM_EXCLUDE, vendored files must be skipped and the rest read:\n%s", out)
	}
	if !strings.Contains(out, "scratch/a.rb") || !strings.Contains(out, "docs/gen/b.rb") {
		t.Errorf("without SEDUM_EXCLUDE, ordinary directories are read:\n%s", out)
	}

	// Comma-separated entries, blank items from stray commas ignored, and
	// whitespace around an entry trimmed.
	t.Setenv("SEDUM_EXCLUDE", " scratch ,,docs/gen/ ,")
	out = regionsOutput(t, output)
	if strings.Contains(out, "scratch") || strings.Contains(out, "docs/gen") || !strings.Contains(out, "app/keep.rb") {
		t.Errorf("SEDUM_EXCLUDE entries were not skipped:\n%s", out)
	}
}

// The variable adds to the config file's exclude list and never replaces it.
func TestSedumExcludeAddsToTheConfigFilesList(t *testing.T) {
	output := t.TempDir()
	markedRubyFile(t, output, "app/keep.rb")
	markedRubyFile(t, output, "generated/a.rb")
	markedRubyFile(t, output, "scratch/b.rb")
	if err := os.MkdirAll(filepath.Join(output, ".sedum"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, ".sedum", "config.yml"), []byte("exclude: [generated]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("SEDUM_EXCLUDE", "scratch")
	out := regionsOutput(t, output)
	if strings.Contains(out, "generated") || strings.Contains(out, "scratch") || !strings.Contains(out, "app/keep.rb") {
		t.Errorf("config and environment entries must both apply:\n%s", out)
	}
}

func TestSedumExcludeRejectsAnUnusableEntry(t *testing.T) {
	gens := filepath.Join("..", "..", "testdata", "generators")
	for _, entry := range []string{"/etc", "../x", "a/../b"} {
		t.Setenv("SEDUM_EXCLUDE", entry)
		output := t.TempDir()

		_, err := exec(t, "regions", "--generators", gens, "--output", output)
		wantErr(t, err, "SEDUM_EXCLUDE", entry)

		records := t.TempDir()
		_, err = exec(t, "grow", "--generators", gens, "--records", records, "--output", output, "--check",
			"--log", filepath.Join(t.TempDir(), "run.log"))
		wantErr(t, err, "SEDUM_EXCLUDE", entry)
	}
}

// grow reads SEDUM_EXCLUDE too: an orphan in an excluded directory is not seen.
func TestGrowCheckHonoursSedumExclude(t *testing.T) {
	gens := filepath.Join("..", "..", "testdata", "generators")
	records := t.TempDir()
	rec := "id: prov-2026-aaaaaaaa\nintent: |\n  Something.\naffected_scope:\n  - lib/**\n"
	if err := os.WriteFile(filepath.Join(records, "prov-2026-aaaaaaaa.yml"), []byte(rec), 0o644); err != nil {
		t.Fatal(err)
	}
	output := t.TempDir()
	markedRubyFile(t, output, "lib/old.rb")
	gone := strings.Replace(excludeEnvMarker, "prov-2026-aaaaaaaa", "prov-2026-dddddddd", 1)
	if err := os.WriteFile(filepath.Join(output, "lib", "old.rb"), []byte(gone), 0o644); err != nil {
		t.Fatal(err)
	}
	check := func() string {
		out, _ := exec(t, "grow", "--generators", gens, "--records", records, "--output", output, "--check",
			"--log", filepath.Join(t.TempDir(), "run.log"))
		return out
	}

	t.Setenv("SEDUM_EXCLUDE", "")
	if out := check(); !strings.Contains(out, "orphan region(s)") {
		t.Fatalf("control: the orphan in lib/ should be reported:\n%s", out)
	}
	t.Setenv("SEDUM_EXCLUDE", "lib")
	if out := check(); strings.Contains(out, "orphan region(s)") {
		t.Errorf("SEDUM_EXCLUDE=lib: the orphan should not be seen:\n%s", out)
	}
}
