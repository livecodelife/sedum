package selection

import (
	"os"
	"path/filepath"
	"testing"
)

// BundledLocalConfig resolves paths from os.Executable(), which a unit test
// cannot redirect - go test's own binary is what os.Executable() reports.
// What is tested here instead is isFile and the fixed filenames it looks
// for, by placing them next to that test binary directly, the same
// arrangement a release archive puts sedum's binary in.

func TestBundledLocalConfig(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}
	dir := filepath.Dir(exe)

	serverPath := filepath.Join(dir, bundledServerName)
	modelPath := filepath.Join(dir, bundledModelName)

	t.Run("neither present", func(t *testing.T) {
		removeIfExists(t, serverPath)
		removeIfExists(t, modelPath)

		if _, ok := BundledLocalConfig(); ok {
			t.Fatal("expected no bundled config with neither file present")
		}
	})

	t.Run("only one present", func(t *testing.T) {
		removeIfExists(t, serverPath)
		writeFile(t, modelPath)
		defer removeIfExists(t, modelPath)

		if _, ok := BundledLocalConfig(); ok {
			t.Fatal("expected no bundled config with only the model present")
		}
	})

	t.Run("both present", func(t *testing.T) {
		writeFile(t, serverPath)
		writeFile(t, modelPath)
		defer removeIfExists(t, serverPath)
		defer removeIfExists(t, modelPath)

		cfg, ok := BundledLocalConfig()
		if !ok {
			t.Fatal("expected a bundled config with both files present")
		}
		if cfg.ServerPath != serverPath {
			t.Errorf("ServerPath = %q, want %q", cfg.ServerPath, serverPath)
		}
		if cfg.ModelPath != modelPath {
			t.Errorf("ModelPath = %q, want %q", cfg.ModelPath, modelPath)
		}
	})
}

func writeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("stand-in for a test"), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func removeIfExists(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Fatalf("removing %s: %v", path, err)
	}
}
