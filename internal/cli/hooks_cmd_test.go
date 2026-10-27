package cli

import (
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// `sedum hooks install` installs the sedum git hooks into the hooks directory
// of the repository containing a directory and prints the result.

func cleanEnv(t *testing.T) []string {
	t.Helper()
	var env []string
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "GIT_") {
			env = append(env, kv)
		}
	}
	return env
}

func gitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmd := osexec.Command("git", "-C", dir, "init", "-q")
	cmd.Env = cleanEnv(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	return dir
}

func realPath(t *testing.T, p string) string {
	t.Helper()
	r, err := filepath.EvalSymlinks(p)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestHooksInstallWritesBothHooksAndPrintsTheResult(t *testing.T) {
	repo := gitRepo(t)

	out, err := exec(t, "hooks", "install", "--dir", repo)

	if err != nil {
		t.Fatalf("hooks install: %v\n%s", err, out)
	}
	hooksDir := filepath.Join(realPath(t, repo), ".git", "hooks")
	want := "hooks directory: " + hooksDir + "\n" +
		"commit-msg: installed\n" +
		"pre-push: installed\n" +
		"\n" +
		"note: a later provenance hooks install replaces these two hooks again, so run sedum hooks install after it\n" +
		"note: the hooks run $SEDUM else sedum on PATH, and $PROVENANCE else provenance on PATH, never a file in the repository\n"
	if realOut := strings.Replace(out, realPath(t, repo), realPath(t, repo), 1); realOut != want {
		t.Errorf("output =\n%s\nwant\n%s", realOut, want)
	}
	for _, name := range []string{"commit-msg", "pre-push"} {
		info, err := os.Stat(filepath.Join(hooksDir, name))
		if err != nil || info.Mode().Perm() != 0o755 {
			t.Errorf("%s: %v %v", name, info, err)
		}
	}
}

func TestHooksInstallKeepsAForeignHookAndReportsUnchangedOnASecondRun(t *testing.T) {
	repo := gitRepo(t)
	hooks := filepath.Join(repo, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hooks, "commit-msg"), []byte("#!/bin/sh\necho mine\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	first, err := exec(t, "hooks", "install", "--dir", repo)
	if err != nil || !strings.Contains(first, "commit-msg: replaced (previous hook kept as commit-msg.bak)\n") {
		t.Fatalf("first: %v\n%s", err, first)
	}
	second, err := exec(t, "hooks", "install", "--dir", repo)

	if err != nil || !strings.Contains(second, "commit-msg: unchanged\n") || !strings.Contains(second, "pre-push: unchanged\n") {
		t.Errorf("second: %v\n%s", err, second)
	}
}

func TestHooksInstallFromASubdirectoryWritesToTheTopLevelHooks(t *testing.T) {
	repo := gitRepo(t)
	sub := filepath.Join(repo, "a", "b")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	if _, err := exec(t, "hooks", "install", "--dir", sub); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(repo, ".git", "hooks", "commit-msg")); err != nil {
		t.Errorf("the hook is not in the top level hooks directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(sub, ".git")); !os.IsNotExist(err) {
		t.Error("a .git directory was created in the subdirectory")
	}
}

func TestHooksInstallIgnoresGitEnvironmentVariables(t *testing.T) {
	repo := gitRepo(t)
	t.Setenv("GIT_DIR", filepath.Join(t.TempDir(), "elsewhere"))

	if _, err := exec(t, "hooks", "install", "--dir", repo); err != nil {
		t.Errorf("a GIT_DIR in the environment broke the lookup: %v", err)
	}
}

func TestHooksInstallFailsOutsideARepositoryAndForAMissingDirectory(t *testing.T) {
	outside := t.TempDir()

	_, err := exec(t, "hooks", "install", "--dir", outside)
	wantErr(t, err, outside+" is not inside a git repository")

	missing := filepath.Join(t.TempDir(), "nope")
	_, err = exec(t, "hooks", "install", "--dir", missing)
	wantErr(t, err, `--dir "`+missing+`" is not an existing directory`)
}

func TestHooksIsAGroupWithInstallAsItsOnlySubcommand(t *testing.T) {
	cmd, _, err := NewRootCommand().Find([]string{"hooks"})
	if err != nil || cmd.Name() != "hooks" || len(cmd.Commands()) != 1 || cmd.Commands()[0].Name() != "install" {
		t.Fatalf("hooks command = %v, %v", cmd, err)
	}
	install := cmd.Commands()[0]
	if f := install.Flags().Lookup("dir"); f == nil || f.Shorthand != "" || f.DefValue != "." {
		t.Errorf("dir flag = %+v", f)
	}
	if _, err := exec(t, "hooks", "install", "extra"); err == nil {
		t.Error("a positional argument was accepted")
	}
}
