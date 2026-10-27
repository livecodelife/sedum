package gitcheck

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// The runner resolves the provenance binary, builds its arguments and
// environment, runs it and reports how it ended.

func fakeBinary(t *testing.T, dir, name, script string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestResolveBinaryPrefersTheFlagThenTheEnvironmentThenPATH(t *testing.T) {
	top := t.TempDir()
	outside := t.TempDir()
	flag := fakeBinary(t, outside, "flag-prov", "exit 0")
	env := fakeBinary(t, outside, "env-prov", "exit 0")
	onPath := fakeBinary(t, filepath.Join(outside, "bin"), "provenance", "exit 0")
	t.Setenv("PATH", filepath.Dir(onPath)+string(os.PathListSeparator)+os.Getenv("PATH"))

	if got, err := ResolveBinary(flag, env, top); err != nil || got != flag {
		t.Errorf("flag given: %q, %v", got, err)
	}
	if got, err := ResolveBinary("", env, top); err != nil || got != env {
		t.Errorf("environment: %q, %v", got, err)
	}
	if got, err := ResolveBinary("", "", top); err != nil || got != onPath {
		t.Errorf("PATH: %q, %v", got, err)
	}
}

func TestResolveBinaryFailsWithAProvenanceErrorWhenNothingIsFound(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if _, err := ResolveBinary("", "", t.TempDir()); err == nil || !strings.HasPrefix(err.Error(), "provenance:") {
		t.Errorf("err = %v, want a provenance: error", err)
	}
}

func TestResolveBinaryRefusesABinaryUnderTheGitTopLevelEvenThroughASymlink(t *testing.T) {
	top := t.TempDir()
	inside := fakeBinary(t, filepath.Join(top, "tools"), "provenance", "exit 0")
	outside := t.TempDir()
	link := filepath.Join(outside, "provenance")
	if err := os.Symlink(inside, link); err != nil {
		t.Fatal(err)
	}

	for name, bin := range map[string]string{"direct": inside, "through a symlink": link} {
		if _, err := ResolveBinary(bin, "", top); err == nil || !strings.HasPrefix(err.Error(), "provenance:") {
			t.Errorf("%s: err = %v, want a refusal", name, err)
		}
	}
}

func TestCheckArgs(t *testing.T) {
	cases := []struct {
		name            string
		sha             string
		staged          bool
		message, config string
		want            []string
	}{
		{"a commit", "abc123", false, "", "", []string{"check", "--commit", "abc123", "--format", "sarif"}},
		{"staged", "", true, "", "", []string{"check", "--staged", "--format", "sarif"}},
		{"staged with a message file and a config", "", true, "/m/MSG", "/c/.provenance.yml", []string{"check", "--staged", "--message-file", "/m/MSG", "--config", "/c/.provenance.yml", "--format", "sarif"}},
		{"a commit with a config", "abc123", false, "/m/ignored", "/c/.provenance.yml", []string{"check", "--commit", "abc123", "--config", "/c/.provenance.yml", "--format", "sarif"}},
	}
	for _, c := range cases {
		if got := CheckArgs(c.sha, c.staged, c.message, c.config); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: CheckArgs = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestEnvironmentClearsGitVariablesExceptTheIndexVariablesWhenStaged(t *testing.T) {
	current := []string{"PATH=/bin", "HOME=/h", "GIT_DIR=/r/.git", "GIT_WORK_TREE=/r", "GIT_INDEX_FILE=/r/.git/index.lock", "GIT_AUTHOR_NAME=x", "GIT_CONFIG_GLOBAL=/g"}

	staged := Environment(true, current)
	ranged := Environment(false, current)

	if want := []string{"PATH=/bin", "HOME=/h", "GIT_DIR=/r/.git", "GIT_WORK_TREE=/r", "GIT_INDEX_FILE=/r/.git/index.lock"}; !reflect.DeepEqual(staged, want) {
		t.Errorf("staged = %v, want %v", staged, want)
	}
	if want := []string{"PATH=/bin", "HOME=/h"}; !reflect.DeepEqual(ranged, want) {
		t.Errorf("range = %v, want %v", ranged, want)
	}
}

func TestRunProvenanceReturnsTheOutputWhenItExitsZeroOrOneAndPrintsSomething(t *testing.T) {
	dir := t.TempDir()
	zero := fakeBinary(t, dir, "zero", `echo '{"runs":[]}'; exit 0`)
	one := fakeBinary(t, dir, "one", `echo '{"runs":[]}'; exit 1`)

	for name, bin := range map[string]string{"exit 0": zero, "exit 1": one} {
		out, err := RunProvenance(bin, []string{"check"}, dir, os.Environ())
		if err != nil || strings.TrimSpace(string(out)) != `{"runs":[]}` {
			t.Errorf("%s: out = %q, err = %v", name, out, err)
		}
	}
}

func TestRunProvenanceFailsOnEmptyOutputOnAStatusAboveOneAndWhenItCannotStart(t *testing.T) {
	dir := t.TempDir()
	cases := map[string]string{
		"empty output":  fakeBinary(t, dir, "empty", "exit 0"),
		"status 2":      fakeBinary(t, dir, "two", `echo '{}'; exit 2`),
		"cannot start":  filepath.Join(dir, "missing"),
		"empty and one": fakeBinary(t, dir, "emptyone", "exit 1"),
	}
	for name, bin := range cases {
		if _, err := RunProvenance(bin, []string{"check"}, dir, os.Environ()); err == nil || !strings.HasPrefix(err.Error(), "provenance:") {
			t.Errorf("%s: err = %v, want a provenance: error", name, err)
		}
	}
}

func TestRunProvenanceRunsInTheWorkingDirectoryWithTheGivenEnvironment(t *testing.T) {
	dir := t.TempDir()
	bin := fakeBinary(t, dir, "show", `echo "$PWD|$MARKER|$*"`)

	out, err := RunProvenance(bin, []string{"check", "--staged"}, dir, []string{"MARKER=m1", "PATH=/usr/bin:/bin"})

	if err != nil {
		t.Fatal(err)
	}
	real, _ := filepath.EvalSymlinks(dir)
	got, _ := filepath.EvalSymlinks(strings.SplitN(strings.TrimSpace(string(out)), "|", 2)[0])
	if got != real || !strings.Contains(string(out), "|m1|check --staged") {
		t.Errorf("out = %q, want the working directory %s, MARKER and the arguments", out, real)
	}
}
