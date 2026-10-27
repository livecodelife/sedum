package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// `sedum init` prepares a directory for sedum grow by calling the scaffold step
// and printing what it did.

func TestInitCreatesTheFilesAndPrintsOneLinePerItem(t *testing.T) {
	dir := t.TempDir()

	out, err := exec(t, "init", "--dir", dir)

	if err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	want := "created    .sedum/context.md\n" +
		"created    .gitignore\n" +
		"skipped    .provenance.yml (no .provenance.yml in this directory)\n"
	if out != want {
		t.Errorf("output =\n%q\nwant\n%q", out, want)
	}
	for _, rel := range []string{".sedum/context.md", ".gitignore"} {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(rel))); err != nil {
			t.Errorf("%s not created: %v", rel, err)
		}
	}
}

func TestInitUpdatesAConfigAndASecondRunChangesNothing(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".provenance.yml"), []byte("exclude_paths: ['a']\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	first, err := exec(t, "init", "--dir", dir)
	if err != nil || !strings.Contains(first, "updated    .provenance.yml\n") {
		t.Fatalf("first: %v\n%s", err, first)
	}
	second, err := exec(t, "init", "--dir", dir)

	want := "unchanged  .sedum/context.md\nunchanged  .gitignore\nunchanged  .provenance.yml\n"
	if err != nil || second != want {
		t.Errorf("second run = %q, %v; want %q", second, err, want)
	}
}

func TestInitFailsAndWritesNothingWhenTheDirectoryIsNotAnExistingDirectory(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")

	out, err := exec(t, "init", "--dir", missing)

	wantErr(t, err, `--dir "`+missing+`" is not an existing directory`)
	if strings.Contains(out, "created") {
		t.Errorf("output = %q, want nothing from the scaffold step", out)
	}
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := exec(t, "init", "--dir", file); err == nil {
		t.Error("a file was accepted as the directory")
	}
}

func TestInitReturnsTheScaffoldErrorUnchangedAndPrintsNothing(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".provenance.yml"), []byte("exclude_paths: nope\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := exec(t, "init", "--dir", dir)

	wantErr(t, err, ".provenance.yml: exclude_paths is not a one-line list or a block list")
	if strings.Contains(out, "created") {
		t.Errorf("output = %q, want no item lines", out)
	}
}

func TestInitUsesTheCurrentDirectoryByDefaultAndTakesNoArguments(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	if _, err := exec(t, "init"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".gitignore")); err != nil {
		t.Errorf("the current directory was not prepared: %v", err)
	}
	if _, err := exec(t, "init", "extra"); err == nil {
		t.Error("a positional argument was accepted")
	}
	cmd, _, err := NewRootCommand().Find([]string{"init"})
	if err != nil || cmd.Name() != "init" || cmd.Flags().Lookup("dir") == nil || cmd.Flags().Lookup("dir").Shorthand != "" || cmd.Flags().Lookup("dir").DefValue != "." {
		t.Errorf("init command = %v, %v", cmd, err)
	}
}
