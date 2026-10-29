package resolve

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A check run classifies records and writes nothing, so a hand-written file its
// template cannot reconcile must not stop it. Phase 3 can be asked to skip such
// a file: it is left as it is and reported as unreconciled
// (prov-2026-4c49ca46).

const unreconcilable = "class User\n  def slug\n  end\nend\n"

func writeExisting(t *testing.T, out, rel, content string) string {
	t.Helper()
	full := filepath.Join(out, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return full
}

func TestSkipUnreconcilableLeavesTheFileAloneAndReportsIt(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()
	full := writeExisting(t, out, "app/models/user.rb", unreconcilable)

	resolutions, _, err := Paths(set, records(t, "app/models/user.rb"), nil)
	if err != nil {
		t.Fatalf("Paths: %v", err)
	}
	files, err := Create(resolutions, Options{Output: out, SkipUnreconcilable: true})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("got %d files, want the skipped one still reported", len(files))
	}
	if !files[0].Existed {
		t.Errorf("the skipped file was not reported as existing")
	}
	for _, want := range []string{"app/models/user.rb", "app/models/{name}.rb", "def slug", "no counterpart"} {
		if !strings.Contains(files[0].Unreconciled, want) {
			t.Errorf("Unreconciled = %q, want it to say %q", files[0].Unreconciled, want)
		}
	}
	got, err := os.ReadFile(full)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != unreconcilable {
		t.Errorf("a skipped file was modified:\n%s", got)
	}
}

func TestWithoutSkipUnreconcilableTheRunStillHalts(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()
	writeExisting(t, out, "app/models/user.rb", unreconcilable)

	resolutions, _, err := Paths(set, records(t, "app/models/user.rb"), nil)
	if err != nil {
		t.Fatalf("Paths: %v", err)
	}
	_, err = Create(resolutions, Options{Output: out})

	wantErr(t, err, "app/models/user.rb", "def slug", "no counterpart")
}

// Only a failure to reconcile is skipped: a file that cannot be read halts.
func TestSkipUnreconcilableDoesNotSkipAnUnreadableFile(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()
	full := writeExisting(t, out, "app/models/user.rb", "class User\nend\n")
	if err := os.Chmod(full, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := os.ReadFile(full); err == nil {
		t.Skip("the file is readable even with no permissions, so this user cannot make it unreadable")
	}

	resolutions, _, err := Paths(set, records(t, "app/models/user.rb"), nil)
	if err != nil {
		t.Fatalf("Paths: %v", err)
	}
	if _, err := Create(resolutions, Options{Output: out, SkipUnreconcilable: true}); err == nil {
		t.Error("an unreadable file was skipped as if it were unreconcilable")
	}
}

// A file that reconciles is reported as before when the request is made: it
// gains what its template adds and is not marked unreconciled.
func TestSkipUnreconcilableReportsAReconcilableFileAsBefore(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()
	full := writeExisting(t, out, "app/models/user.rb", "class User\nend\n")

	resolutions, _, err := Paths(set, records(t, "app/models/user.rb"), nil)
	if err != nil {
		t.Fatalf("Paths: %v", err)
	}
	files, err := Create(resolutions, Options{Output: out, SkipUnreconcilable: true})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if files[0].Unreconciled != "" {
		t.Errorf("Unreconciled = %q for a file that reconciles", files[0].Unreconciled)
	}
	got, err := os.ReadFile(full)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "sedum:anchor:class_body") {
		t.Errorf("the file did not gain the marker its template plants:\n%s", got)
	}
}

// One skipped file does not hide another file's problem or stop the others.
func TestSkipUnreconcilableKeepsTheOtherFilesOfTheRun(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()
	writeExisting(t, out, "app/models/user.rb", unreconcilable)

	resolutions, _, err := Paths(set, records(t, "app/models/user.rb", "app/models/post.rb"), nil)
	if err != nil {
		t.Fatalf("Paths: %v", err)
	}
	files, err := Create(resolutions, Options{Output: out, SkipUnreconcilable: true})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	unreconciled := 0
	for _, f := range files {
		if f.Unreconciled != "" {
			unreconciled++
		}
	}
	if len(files) != 2 || unreconciled != 1 {
		t.Errorf("got %d files with %d unreconciled, want 2 with 1", len(files), unreconciled)
	}
	if _, err := os.Stat(filepath.Join(out, "app/models/post.rb")); err != nil {
		t.Errorf("the other file of the run was not created: %v", err)
	}
}
