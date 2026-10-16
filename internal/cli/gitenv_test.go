package cli

import (
	"os"
	"strings"
	"testing"
)

// TestMain isolates git before any test runs (COMMON.md, "Git in tests and
// scripts"; prov-2026-a500e28a). linespec runs specs from inside a git hook,
// which exports GIT_DIR, GIT_INDEX_FILE and the like; a git command a test
// runs - here or inside a library it calls - would otherwise act on whatever
// repository the hook was running for.
func TestMain(m *testing.M) {
	for _, kv := range os.Environ() {
		if name, _, _ := strings.Cut(kv, "="); strings.HasPrefix(name, "GIT_") {
			os.Unsetenv(name)
		}
	}
	home, err := os.MkdirTemp("", "sedum-test-home-")
	if err != nil {
		panic(err)
	}
	os.Setenv("HOME", home)
	os.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	os.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	code := m.Run()
	os.RemoveAll(home)
	os.Exit(code)
}

// The isolation above holds for every test in the package.
func TestGitIsIsolatedFromTheEnvironment(t *testing.T) {
	for _, kv := range os.Environ() {
		name, value, _ := strings.Cut(kv, "=")
		switch {
		case name == "GIT_CONFIG_NOSYSTEM" || name == "GIT_CONFIG_GLOBAL":
		case strings.HasPrefix(name, "GIT_"):
			t.Errorf("%s=%s survived isolation", name, value)
		}
	}
	if os.Getenv("GIT_CONFIG_GLOBAL") != os.DevNull {
		t.Error("global git config is not turned off")
	}
}
