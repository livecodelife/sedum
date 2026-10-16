package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// hooksRepo is checkRepo with git's hooks path pointed at its own
// .git/hooks (TestMain's global config points every repo at an empty dir).
func hooksRepo(t *testing.T) (root, cfg, hooks string) {
	t.Helper()
	root, cfg = checkRepo(t)
	hooks = filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "config", "core.hooksPath", hooks)
	return root, cfg, hooks
}

func TestHooksInstallBacksUp(t *testing.T) {
	root, cfg, hooks := hooksRepo(t)
	writeFile(t, filepath.Join(hooks, "pre-commit"), "#!/bin/sh\n# linespec hook\n")
	out := mustRun(t, "-c", cfg, "hooks", "install")
	if !strings.Contains(out, "Installed git hooks to "+hooks) {
		t.Fatalf("output:\n%s", out)
	}
	for _, name := range []string{"pre-commit", "commit-msg"} {
		st, err := os.Stat(filepath.Join(hooks, name))
		if err != nil || st.Mode()&0o111 == 0 {
			t.Fatalf("%s not installed executable (%v)", name, err)
		}
		if body := readFile(t, filepath.Join(hooks, name)); !strings.Contains(body, "provenance") || strings.Contains(body, "linespec hook") {
			t.Fatalf("%s:\n%s", name, body)
		}
	}
	if got := readFile(t, filepath.Join(hooks, "pre-commit.bak")); got != "#!/bin/sh\n# linespec hook\n" {
		t.Fatalf("backup:\n%s", got)
	}
	if _, err := os.Stat(filepath.Join(hooks, "commit-msg.bak")); !os.IsNotExist(err) {
		t.Fatalf("a backup was made of a hook that did not exist (%v)", err)
	}
	// Reinstalling keeps the original backup.
	mustRun(t, "-c", cfg, "hooks", "install")
	if got := readFile(t, filepath.Join(hooks, "pre-commit.bak")); got != "#!/bin/sh\n# linespec hook\n" {
		t.Fatalf("reinstall replaced the backup:\n%s", got)
	}
	_ = root
}

// TestHooksGateCommits commits through the installed hooks with the CLI
// built from this package.
func TestHooksGateCommits(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "provenance")
	build := exec.Command("go", "build", "-o", bin, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	t.Setenv("PROVENANCE", bin)
	root, cfg, _ := hooksRepo(t)
	mustRun(t, "-c", cfg, "hooks", "install")

	commit := func(msg string, files map[string]string) error {
		for f, body := range files {
			writeFile(t, filepath.Join(root, f), body)
			gitIn(t, root, "add", "--", f)
		}
		c := exec.Command("git", "-C", root, "commit", "-q", "-m", msg)
		c.Env = append(os.Environ(), "GIT_DIR="+filepath.Join(root, ".git"), "GIT_WORK_TREE="+root)
		out, err := c.CombinedOutput()
		if err != nil {
			t.Logf("commit %q refused:\n%s", msg, out)
		}
		return err
	}
	if err := commit("Change a ["+chkID+"]", map[string]string{"a.go": "package a // 1\n"}); err != nil {
		t.Fatal("a tagged in-scope commit was refused")
	}
	if err := commit("Untagged", map[string]string{"a.go": "package a // 2\n"}); err == nil {
		t.Fatal("an untagged commit passed the hooks")
	}
	gitIn(t, root, "reset", "-q", "--hard")
	if err := commit("Change b ["+chkID+"]", map[string]string{"b.go": "package b // 1\n"}); err == nil {
		t.Fatal("an out-of-scope commit passed the hooks")
	}
	gitIn(t, root, "reset", "-q", "--hard")
	// A record that fails lint is blocked by pre-commit.
	path := filepath.Join("provenance", chkID+".yml")
	broken := strings.Replace(readFile(t, filepath.Join(root, path)), "intent: Do the thing\n", "intent: \"\"\n", 1)
	if err := commit("Break it ["+chkID+"]", map[string]string{path: broken}); err == nil {
		t.Fatal("a record failing lint passed pre-commit")
	}
}
