package verify

// Compile-proof for prov-2026-b5c21813: the platform split must leave a tree
// that builds for Windows (no Unix-only symbol, cgo off) as well as Linux and
// macOS. It shells out to the host's go toolchain, so it needs `go` on PATH
// (skipped otherwise) and relies on the build cache to stay fast.

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// moduleRoot is the directory holding go.mod, found relative to this file.
func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate the test file")
	}
	dir := filepath.Dir(file)
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatalf("no go.mod above %s", file)
	return ""
}

func TestCrossBuild(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go is not on PATH")
	}
	root := moduleRoot(t)
	cases := []struct {
		name string
		goos string
		args []string
	}{
		{"build sedum for windows/amd64", "windows", []string{"build", "-o", os.DevNull, "./cmd/sedum"}},
		{"vet verify and its tests for windows/amd64", "windows", []string{"vet", "./internal/verify"}},
		{"build sedum for linux/amd64", "linux", []string{"build", "-o", os.DevNull, "./cmd/sedum"}},
		{"build sedum for darwin/amd64", "darwin", []string{"build", "-o", os.DevNull, "./cmd/sedum"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			cmd := exec.CommandContext(ctx, goBin, c.args...)
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "GOOS="+c.goos, "GOARCH=amd64", "CGO_ENABLED=0")
			var out bytes.Buffer
			cmd.Stdout, cmd.Stderr = &out, &out
			if err := cmd.Run(); err != nil {
				t.Fatalf("go %v with GOOS=%s GOARCH=amd64 CGO_ENABLED=0 failed: %v\n%s", c.args, c.goos, err, out.String())
			}
		})
	}
}
