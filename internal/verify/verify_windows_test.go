//go:build windows

package verify

// NOTE: this file only compiles and runs on Windows. It CANNOT run on macOS or
// Linux, where it is excluded by the build constraint above; it is type-checked
// from other hosts only by TestCrossBuild (go vet with GOOS=windows). Run it on
// a Windows host with: go test ./internal/verify -count=1
//
// It judges prov-2026-b5c21813's Windows constraints: a timeout or cancellation
// terminates the command AND every process it started; the command goes through
// `sh -c <command>` when an sh is on PATH and through `cmd /C` otherwise.

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	helperEnv  = "SEDUM_VERIFY_HELPER" // "parent" or "grandchild"
	markDirEnv = "SEDUM_VERIFY_MARKDIR"
	exeEnv     = "SEDUM_VERIFY_TESTEXE"
	fakeShEnv  = "SEDUM_FAKE_SH_OUT"
)

// TestMain lets this test binary stand in for processes the tests need to
// spawn: a parent that starts a detached grandchild, the grandchild, and a
// fake sh.exe that records how it was invoked.
func TestMain(m *testing.M) {
	if strings.EqualFold(filepath.Base(os.Args[0]), "sh.exe") && os.Getenv(fakeShEnv) != "" {
		_ = os.WriteFile(os.Getenv(fakeShEnv), []byte(strings.Join(os.Args[1:], "\n")), 0o644)
		os.Exit(0)
	}
	switch os.Getenv(helperEnv) {
	case "parent":
		child := exec.Command(os.Args[0])
		child.Env = append(os.Environ(), helperEnv+"=grandchild")
		if err := child.Start(); err != nil {
			os.Exit(2)
		}
		time.Sleep(90 * time.Second)
		os.Exit(0)
	case "grandchild":
		dir := os.Getenv(markDirEnv)
		_ = os.WriteFile(filepath.Join(dir, "started"), []byte("x"), 0o644)
		time.Sleep(4 * time.Second)
		_ = os.WriteFile(filepath.Join(dir, "survived"), []byte("x"), 0o644)
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func runTreeKill(t *testing.T, ctx context.Context, timeout time.Duration) (markDir string, failures []Failure, err error) {
	t.Helper()
	exe, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	markDir = t.TempDir()
	t.Setenv(helperEnv, "parent")
	t.Setenv(markDirEnv, markDir)
	t.Setenv(exeEnv, exe)
	command := pick(`"$`+exeEnv+`"`, `"%`+exeEnv+`%"`)
	failures, err = Run(ctx, Config{Commands: []string{command}, Timeout: timeout}, tree(t, nil), nil, &bytes.Buffer{})
	return markDir, failures, err
}

func assertGrandchildKilled(t *testing.T, markDir string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(markDir, "started")); err != nil {
		t.Fatal("setup problem: the grandchild never started before the kill, so this run proves nothing")
	}
	time.Sleep(5 * time.Second) // longer than the grandchild's own 4s delay
	if _, err := os.Stat(filepath.Join(markDir, "survived")); err == nil {
		t.Error("a grandchild outlived the timeout/cancellation: only the shell was terminated, not the process tree")
	}
}

func TestTimeoutTerminatesAGrandchildToo(t *testing.T) {
	markDir, failures, err := runTreeKill(t, context.Background(), 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if len(failures) != 1 || !strings.HasPrefix(failures[0].Status, "timed out after") {
		t.Fatalf("failures = %+v, want one timeout", failures)
	}
	assertGrandchildKilled(t, markDir)
}

func TestCancellationTerminatesAGrandchildToo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	time.AfterFunc(2*time.Second, cancel)
	markDir, _, err := runTreeKill(t, ctx, 120*time.Second)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	assertGrandchildKilled(t, markDir)
}

// An sh found on PATH is the interpreter, invoked as `sh -c <command>`.
func TestShOnPathIsInvokedWithDashC(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "sh.exe"), data, 0o755); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "seen")
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv(fakeShEnv, out)

	command := `echo "a  b" && exit 7`
	failures, err := Run(context.Background(), Config{Commands: []string{command}, Timeout: 30 * time.Second}, tree(t, nil), nil, &bytes.Buffer{})
	if err != nil || len(failures) != 0 {
		t.Fatalf("failures = %+v, err = %v (the fake sh exits 0)", failures, err)
	}
	seen, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("the sh.exe on PATH was never run: %v", err)
	}
	if got, want := string(seen), "-c\n"+command; got != want {
		t.Errorf("sh was invoked with %q, want %q", got, want)
	}
}

// With no sh on PATH the command runs through cmd /C, text passed verbatim
// (quotes preserved, & as the command separator).
func TestWithoutShTheCommandRunsThroughCmd(t *testing.T) {
	system32 := filepath.Join(os.Getenv("SystemRoot"), "System32")
	t.Setenv("PATH", system32)
	if hasSh() {
		t.Skip("an sh is still resolvable with PATH reduced to System32")
	}
	cfg := Config{Commands: []string{`echo "quoted arg"& exit 4`}, Timeout: 30 * time.Second}
	failures, err := Run(context.Background(), cfg, tree(t, nil), nil, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if len(failures) != 1 || failures[0].Status != "exit status 4" {
		t.Fatalf("failures = %+v, want exit status 4", failures)
	}
	if !strings.Contains(failures[0].Output, `"quoted arg"`) {
		t.Errorf("Output = %q, want the quoted argument delivered to cmd verbatim", failures[0].Output)
	}
}

func TestWithoutShATimeoutStillTerminatesTheTree(t *testing.T) {
	t.Setenv("PATH", filepath.Join(os.Getenv("SystemRoot"), "System32"))
	if hasSh() {
		t.Skip("an sh is still resolvable with PATH reduced to System32")
	}
	markDir, failures, err := runTreeKill(t, context.Background(), 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if len(failures) != 1 || !strings.HasPrefix(failures[0].Status, "timed out after") {
		t.Fatalf("failures = %+v, want one timeout", failures)
	}
	assertGrandchildKilled(t, markDir)
}
