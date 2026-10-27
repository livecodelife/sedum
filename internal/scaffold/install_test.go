package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

// InstallHook writes one hook file into a hooks directory, keeping a different
// existing hook as a backup.

func TestInstallHookWritesANewHookExecutableAndCreatesTheDirectory(t *testing.T) {
	hooks := filepath.Join(t.TempDir(), "a", "hooks")

	result, err := InstallHook(hooks, "commit-msg", "#!/bin/sh\nexit 0\n")

	if err != nil || result != "installed" {
		t.Fatalf("result = %q, err = %v", result, err)
	}
	info, err := os.Stat(filepath.Join(hooks, "commit-msg"))
	if err != nil || info.Mode().Perm() != 0o755 {
		t.Errorf("hook mode = %v, err = %v", info.Mode().Perm(), err)
	}
}

func TestInstallHookLeavesAnIdenticalHookAloneButMakesItExecutable(t *testing.T) {
	hooks := t.TempDir()
	write(t, hooks, "pre-push", "same\n", 0o644)
	before, _ := os.Stat(filepath.Join(hooks, "pre-push"))

	result, err := InstallHook(hooks, "pre-push", "same\n")

	if err != nil || result != "unchanged" {
		t.Fatalf("result = %q, err = %v", result, err)
	}
	after, _ := os.Stat(filepath.Join(hooks, "pre-push"))
	if after.Mode().Perm() != 0o755 || !after.ModTime().Equal(before.ModTime()) {
		t.Errorf("mode = %v, modtime changed = %v", after.Mode().Perm(), !after.ModTime().Equal(before.ModTime()))
	}
}

func TestInstallHookKeepsADifferentHookAsABackupOnce(t *testing.T) {
	hooks := t.TempDir()
	write(t, hooks, "commit-msg", "mine\n", 0o755)

	first, err := InstallHook(hooks, "commit-msg", "sedum's\n")
	if err != nil || first != "replaced (previous hook kept as commit-msg.bak)" {
		t.Fatalf("first = %q, err = %v", first, err)
	}
	if read(t, hooks, "commit-msg.bak") != "mine\n" || read(t, hooks, "commit-msg") != "sedum's\n" {
		t.Error("the backup or the new hook is wrong")
	}

	write(t, hooks, "commit-msg", "a third hook\n", 0o644)
	second, err := InstallHook(hooks, "commit-msg", "sedum's\n")
	if err != nil || second != "replaced" {
		t.Fatalf("second = %q, err = %v", second, err)
	}
	if read(t, hooks, "commit-msg.bak") != "mine\n" {
		t.Error("an existing backup was touched")
	}
	info, _ := os.Stat(filepath.Join(hooks, "commit-msg"))
	if info.Mode().Perm() != 0o755 {
		t.Errorf("replaced hook mode = %v, want 0755", info.Mode().Perm())
	}
}

func TestInstallHookReturnsTheErrorWhenItCannotWrite(t *testing.T) {
	hooks := t.TempDir()
	// A directory where the hook file should be makes the read fail with an
	// error other than a missing file.
	if err := os.Mkdir(filepath.Join(hooks, "pre-push"), 0o755); err != nil {
		t.Fatal(err)
	}

	result, err := InstallHook(hooks, "pre-push", "x\n")

	if err == nil || result != "" {
		t.Errorf("result = %q, err = %v", result, err)
	}
}
