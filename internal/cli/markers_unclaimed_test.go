package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A record's scope lists every kind of file its change touches. markers
// migrate reads the files a package claims, skips the other paths, says how
// many it skipped, and exits 0.
func TestMarkersMigrateSkipsScopePathsNoGeneratorPackageClaims(t *testing.T) {
	f := newMigFixture(t, false)
	const rec = "prov-2026-aaaaaaaa"
	scope := "  - app/models/user.rb\n  - docs/NOTES.md\n  - scripts/run.sh\n  - Makefile\n"
	text := "id: " + rec + "\nintent: |\n  Add fields.\naffected_scope:\n" + scope
	if err := os.WriteFile(filepath.Join(f.recs, rec+".yml"), []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := f.run(t, "--dry-run")
	if err != nil {
		t.Fatalf("migrate --dry-run: %v\n%s", err, out)
	}
	if !strings.Contains(out, "would migrate app/models/user.rb  addField") {
		t.Errorf("the .rb file's regions were not listed:\n%s", out)
	}
	for _, unclaimed := range []string{"docs/NOTES.md", "scripts/run.sh", "Makefile"} {
		if strings.Contains(out, unclaimed) {
			t.Errorf("the unclaimed path %s is in the output:\n%s", unclaimed, out)
		}
	}
	if !strings.Contains(out, "skipped 3 scope paths no generator package claims\n") {
		t.Errorf("the skipped count is missing:\n%s", out)
	}

	// The line is the same when the migration is real.
	real, err := f.run(t)
	if err != nil {
		t.Fatalf("migrate: %v\n%s", err, real)
	}
	if !strings.Contains(real, "skipped 3 scope paths no generator package claims\n") {
		t.Errorf("the skipped count is missing from a real migration:\n%s", real)
	}

	clean := newMigFixture(t, false)
	if out, err := clean.run(t, "--dry-run"); err != nil || strings.Contains(out, "skipped") {
		t.Errorf("a scope with every path claimed reports a skip: %v\n%s", err, out)
	}
}
