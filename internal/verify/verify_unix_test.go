//go:build unix

package verify

// Unix semantics of verify commands (prov-2026-b5c21813): sh -c, a process
// group of its own, SIGKILL to the whole group on timeout or cancellation.
// The first test is the original process-group test, unchanged in meaning.

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// A timeout kills the command and everything it started, reports it as a
// failure, and still removes the scratch directory.
func TestRunKillsTheWholeProcessGroupOnTimeoutAndCleansUp(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "cwd")
	cfg := Config{Commands: []string{"pwd > " + marker + "; sleep 30 & sleep 30"}, Timeout: 300 * time.Millisecond}
	start := time.Now()
	failures, err := Run(context.Background(), cfg, tree(t, nil), nil, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Errorf("Run took %s: the background process held the run open", elapsed)
	}
	if len(failures) != 1 || !strings.HasPrefix(failures[0].Status, "timed out") {
		t.Fatalf("failures = %+v, want one timeout", failures)
	}
	data, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(strings.TrimSpace(string(data))); err == nil {
		t.Errorf("the scratch directory %s was not removed", strings.TrimSpace(string(data)))
	}
}

// processGone polls for pid to disappear (a killed orphan is reaped by init
// asynchronously, so allow a short grace period).
func processGone(pid int) bool {
	for i := 0; i < 50; i++ {
		if err := syscall.Kill(pid, 0); errors.Is(err, syscall.ESRCH) {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

func readPid(t *testing.T, path string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the command never recorded its background pid: %v", err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatalf("bad pid file %q: %v", data, err)
	}
	return pid
}

// No descendant of the shell survives a timeout, observed by pid rather than
// inferred from how long Run took.
func TestTimeoutLeavesNoDescendantAlive(t *testing.T) {
	pidFile := filepath.Join(t.TempDir(), "pid")
	cfg := Config{Commands: []string{"sleep 61 & echo $! > " + pidFile + "; sleep 62"}, Timeout: 700 * time.Millisecond}
	failures, err := Run(context.Background(), cfg, tree(t, nil), nil, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if len(failures) != 1 || !strings.HasPrefix(failures[0].Status, "timed out") {
		t.Fatalf("failures = %+v, want one timeout", failures)
	}
	pid := readPid(t, pidFile)
	if !processGone(pid) {
		_ = syscall.Kill(pid, syscall.SIGKILL)
		t.Errorf("background process %d outlived the timeout", pid)
	}
}

// Cancellation kills the same whole group.
func TestCancellationLeavesNoDescendantAlive(t *testing.T) {
	pidFile := filepath.Join(t.TempDir(), "pid")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	time.AfterFunc(700*time.Millisecond, cancel)
	cfg := Config{Commands: []string{"sleep 63 & echo $! > " + pidFile + "; sleep 64"}, Timeout: 120 * time.Second}
	_, err := Run(ctx, cfg, tree(t, nil), nil, &bytes.Buffer{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	pid := readPid(t, pidFile)
	if !processGone(pid) {
		_ = syscall.Kill(pid, syscall.SIGKILL)
		t.Errorf("background process %d outlived the cancellation", pid)
	}
}

// The command is run as `sh -c <command>`, looked up on PATH, with the scratch
// directory (not the output directory) as cwd, in a process group of its own.
func TestCommandRunsAsShDashCInTheScratchDirInItsOwnGroup(t *testing.T) {
	bin := t.TempDir()
	out := filepath.Join(t.TempDir(), "seen")
	script := "#!/bin/sh\n" +
		"printf '%s\\n' \"$@\" > \"$SEDUM_FAKE_SH_OUT\"\n" +
		"pwd -P >> \"$SEDUM_FAKE_SH_OUT\"\n" +
		"ps -o pgid= -p $$ >> \"$SEDUM_FAKE_SH_OUT\"\n" +
		"echo $$ >> \"$SEDUM_FAKE_SH_OUT\"\n"
	if err := os.WriteFile(filepath.Join(bin, "sh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("SEDUM_FAKE_SH_OUT", out)

	root := tree(t, nil)
	rootReal, _ := filepath.EvalSymlinks(root)
	command := `echo "a  b" && exit 7`
	failures, err := Run(context.Background(), Config{Commands: []string{command}, Timeout: 20 * time.Second}, root, nil, &bytes.Buffer{})
	if err != nil || len(failures) != 0 {
		t.Fatalf("failures = %+v, err = %v (the fake sh exits 0)", failures, err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("the fake sh on PATH was never run: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 5 || lines[0] != "-c" || lines[1] != command {
		t.Fatalf("sh was invoked with %q, want exactly [-c %q]", lines, command)
	}
	if cwd := lines[2]; cwd == rootReal {
		t.Errorf("cwd = %q, want the scratch directory, not the output directory %q", cwd, rootReal)
	}
	if pgid, pid := strings.TrimSpace(lines[3]), strings.TrimSpace(lines[4]); pgid != pid {
		t.Errorf("process group %s != pid %s: the command must lead its own process group", pgid, pid)
	}
}
