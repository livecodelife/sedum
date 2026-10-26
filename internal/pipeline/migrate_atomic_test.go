package pipeline

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// A file a person or a tool has locked (read-only) is still rewritten by a
// migration, which replaces the file by renaming a new one over it, and the
// file keeps its mode afterwards.
func TestMigrateRewritesAReadOnlyFileAndLeavesItReadOnly(t *testing.T) {
	f := newMigFixture(t)
	gen, pkg := f.current(t)
	f.writeUser(t, f.email(gen, pkg))
	full := filepath.Join(f.cfg.Output, filepath.FromSlash(migPath))
	if err := os.Chmod(full, 0o444); err != nil {
		t.Fatal(err)
	}

	report := f.migrate(t)

	if !reflect.DeepEqual(report.Files, []string{migPath}) {
		t.Fatalf("Files = %v, want the locked file rewritten", report.Files)
	}
	if got, want := readFile(t, full), migFile(f.emailV2(gen, pkg)); got != want {
		t.Errorf("locked file =\n%s\nwant\n%s", got, want)
	}
	info, err := os.Stat(full)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o444 {
		t.Errorf("mode = %v, want the lock kept (0444)", info.Mode().Perm())
	}
}

// A migration that cannot write one of its files changes nothing: no file is
// rewritten, no recording is written and no temporary file is left behind.
func TestMigrateThatCannotWriteOneFileChangesNothing(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("a directory cannot be made unwritable for root")
	}
	f := newMigFixture(t)
	gen, pkg := f.current(t)
	f.writeUser(t, f.email(gen, pkg))

	// A second file whose directory cannot take a new file, so a replacement
	// cannot be staged beside it. It sorts after app/, so a migration that
	// wrote as it went would already have rewritten the first file.
	lockedDir := filepath.Join(f.cfg.Output, "lib", "locked")
	other := migFile(migRegion("  ", "addField", migV1Attrs(v2RecordA, gen, pkg, migPhoneKwargs), "attribute :phone\n"))
	if err := os.MkdirAll(lockedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	otherPath := filepath.Join(lockedDir, "thing.rb")
	if err := os.WriteFile(otherPath, []byte(other), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(lockedDir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(lockedDir, 0o755) })

	before := treeContents(t, f.cfg.Output)

	if _, err := Migrate(f.mc); err == nil {
		t.Fatal("Migrate succeeded though a file cannot be replaced")
	}

	if after := treeContents(t, f.cfg.Output); !reflect.DeepEqual(before, after) {
		for path := range after {
			if before[path] != after[path] {
				t.Errorf("%s changed or appeared after the failed migration", path)
			}
		}
		for path := range before {
			if _, ok := after[path]; !ok {
				t.Errorf("%s disappeared", path)
			}
		}
	}
	for path := range treeContents(t, f.cfg.Output) {
		if strings.Contains(path, ".sedum") {
			t.Errorf("a recording was written: %s", path)
		}
	}
}
