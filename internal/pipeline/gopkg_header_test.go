package pipeline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/recording"
)

// setBuildConstraint is the one action that writes above the package clause.
// These tests run against generators/go itself, not the testdata reference
// copy, because the action lives only in the real package.

func realGoGenerators(t *testing.T) string {
	t.Helper()
	files := map[string]string{}
	base := filepath.Join("..", "..", "generators")
	err := filepath.WalkDir(filepath.Join(base, "go"), func(full string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(base, full)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = readFile(t, full)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return writeTree(t, files)
}

func constraintInvocations(constraint string) []recording.Invocation {
	return []recording.Invocation{
		{Action: "setBuildConstraint", Kwargs: map[string]any{"file": storePath, "constraint": constraint}},
		{Action: "setFunction", Kwargs: map[string]any{"file": storePath, "name": "hook", "code": "func hook() {}"}},
	}
}

func TestSetBuildConstraintPrecedesPackageAndIsGofmtClean(t *testing.T) {
	gens := realGoGenerators(t)
	output := t.TempDir()
	rec := goRecording(constraintInvocations("unix"))
	if _, err := Replay(rec, ReplayConfig{Generators: gens, Output: output}); err != nil {
		t.Fatalf("Replay: %v", err)
	}
	got := mustBeGofmtClean(t, filepath.Join(output, storePath))
	build := strings.Index(got, "//go:build unix\n\n")
	pkg := strings.Index(got, "package record")
	if build < 0 || pkg < 0 || build > pkg {
		t.Fatalf("the constraint must precede the package clause and be followed by a blank line:\n%s", got)
	}
	if strings.Count(got, "//go:build") != 1 {
		t.Errorf("want exactly one constraint:\n%s", got)
	}

	if _, err := Replay(rec, ReplayConfig{Generators: gens, Output: output}); err != nil {
		t.Fatalf("second Replay: %v", err)
	}
	if again := readFile(t, filepath.Join(output, storePath)); again != got {
		t.Errorf("a second replay changed the file:\n%s", again)
	}

	// A different constraint replaces the first rather than stacking.
	if _, err := Replay(goRecording(constraintInvocations("not_windows")), ReplayConfig{Generators: gens, Output: output}); err != nil {
		t.Fatalf("third Replay: %v", err)
	}
	revised := mustBeGofmtClean(t, filepath.Join(output, storePath))
	if strings.Count(revised, "//go:build") != 1 || !strings.Contains(revised, "//go:build !windows\n") {
		t.Errorf("the constraint was not replaced in place:\n%s", revised)
	}
}

// The constraint is a declared variant, so nothing else can be written above
// the package clause.
func TestSetBuildConstraintRejectsAnUndeclaredExpression(t *testing.T) {
	for _, bad := range []string{"!windows", "unix || darwin", "unix\npackage evil", ""} {
		gens := realGoGenerators(t)
		output := t.TempDir()
		_, err := Replay(goRecording(constraintInvocations(bad)), ReplayConfig{Generators: gens, Output: output})
		if err == nil {
			t.Errorf("constraint %q was accepted", bad)
		}
	}
}
