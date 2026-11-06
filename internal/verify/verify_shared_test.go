package verify

// Platform-neutral behaviour of Run/runCommand (prov-2026-b5c21813, "On every
// platform" constraints). Commands are chosen per host: sh syntax when an sh is
// on PATH (every Unix, and Windows with Git for Windows), cmd.exe syntax
// otherwise. Nothing here names a hook, so it judges any implementation.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// hasSh reports whether verify commands will be interpreted by sh on this host.
func hasSh() bool {
	_, err := exec.LookPath("sh")
	return err == nil
}

// pick returns the sh spelling when sh is available, else the cmd.exe one.
func pick(sh, cmd string) string {
	if hasSh() {
		return sh
	}
	return cmd
}

// sleepForever is a command that outlives every test deadline.
func sleepForever(prefix string) string {
	return pick(prefix+"; sleep 60", prefix+"& ping -n 60 127.0.0.1 >nul")
}

// redirectTemp points every platform's temp-dir lookup at a fresh directory
// and returns it, so a test can assert Run left nothing behind.
func redirectTemp(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("TMPDIR", dir)
	t.Setenv("TMP", dir)
	t.Setenv("TEMP", dir)
	return dir
}

func assertEmptyDir(t *testing.T, dir, when string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("scratch directory not removed %s: %v left in the temp dir", when, names)
	}
}

func TestRunExitZeroPasses(t *testing.T) {
	tmp := redirectTemp(t)
	cfg := Config{Commands: []string{"exit 0"}, Timeout: 20 * time.Second}
	failures, err := Run(context.Background(), cfg, tree(t, nil), nil, &bytes.Buffer{})
	if err != nil || len(failures) != 0 {
		t.Fatalf("failures = %+v, err = %v, want a pass", failures, err)
	}
	assertEmptyDir(t, tmp, "after a pass")
}

func TestRunNonZeroExitIsAFailureWithStatusAndOutput(t *testing.T) {
	tmp := redirectTemp(t)
	cfg := Config{
		Commands: []string{pick("echo boom; exit 3", "echo boom& exit 3"), "exit 0"},
		Timeout:  20 * time.Second,
	}
	failures, err := Run(context.Background(), cfg, tree(t, nil), nil, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if len(failures) != 1 {
		t.Fatalf("failures = %+v, want exactly the first command", failures)
	}
	if failures[0].Status != "exit status 3" {
		t.Errorf("Status = %q, want %q", failures[0].Status, "exit status 3")
	}
	if !strings.Contains(failures[0].Output, "boom") {
		t.Errorf("Output = %q, want it to hold the command's output", failures[0].Output)
	}
	assertEmptyDir(t, tmp, "after a failure")
}

func TestRunTimeoutReportsStatusAndOutputTailAndCleansUp(t *testing.T) {
	tmp := redirectTemp(t)
	timeout := 1500 * time.Millisecond
	cfg := Config{Commands: []string{sleepForever("echo before-timeout")}, Timeout: timeout}
	start := time.Now()
	failures, err := Run(context.Background(), cfg, tree(t, nil), nil, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > 30*time.Second {
		t.Errorf("Run took %s: the timeout did not stop the command", elapsed)
	}
	if len(failures) != 1 {
		t.Fatalf("failures = %+v, want one timeout", failures)
	}
	if want := fmt.Sprintf("timed out after %s", timeout); failures[0].Status != want {
		t.Errorf("Status = %q, want %q", failures[0].Status, want)
	}
	if !strings.Contains(failures[0].Output, "before-timeout") || len(failures[0].Output) > MaxOutput {
		t.Errorf("Output = %q, want the captured output (at most %d bytes)", failures[0].Output, MaxOutput)
	}
	assertEmptyDir(t, tmp, "on the timeout path")
}

func TestRunCancelledContextReturnsItsErrorAndCleansUp(t *testing.T) {
	tmp := redirectTemp(t)
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(700*time.Millisecond, cancel)
	defer cancel()
	cfg := Config{Commands: []string{sleepForever("echo waiting")}, Timeout: 120 * time.Second}
	start := time.Now()
	failures, err := Run(ctx, cfg, tree(t, nil), nil, &bytes.Buffer{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if len(failures) != 0 {
		t.Errorf("failures = %+v, want none: a cancellation is an error, not a failed check", failures)
	}
	if elapsed := time.Since(start); elapsed > 30*time.Second {
		t.Errorf("Run took %s: cancellation did not stop the command", elapsed)
	}
	assertEmptyDir(t, tmp, "on the cancellation path")
}

func TestRunAlreadyCancelledContextRunsNothingToCompletion(t *testing.T) {
	redirectTemp(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cfg := Config{Commands: []string{sleepForever("echo hi")}, Timeout: 120 * time.Second}
	start := time.Now()
	_, err := Run(ctx, cfg, tree(t, nil), nil, &bytes.Buffer{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if elapsed := time.Since(start); elapsed > 30*time.Second {
		t.Errorf("Run took %s with a context that was already cancelled", elapsed)
	}
}
