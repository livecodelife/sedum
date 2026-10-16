package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/provenance/record"
)

func TestDeprecate(t *testing.T) {
	root, cfg := completeRepo(t, "commit_on_status_change: true\n", "true\n")
	out := mustRun(t, "-c", cfg, "deprecate", "--record", chkID, "--reason", "Replaced by the v2 flow")
	if !strings.Contains(out, "✓ "+chkID+" marked as deprecated") {
		t.Fatalf("output:\n%s", out)
	}
	if rec := load(t, root, chkID); rec.Status != record.StatusDeprecated || rec.Version != 1 || rec.Reason != "" {
		t.Fatalf("status %s version %d reason %q", rec.Status, rec.Version, rec.Reason)
	}
	if msg := gitIn(t, root, "log", "-1", "--format=%s"); msg != "Deprecate provenance record ["+chkID+"@1]: Replaced by the v2 flow" {
		t.Fatalf("commit %q", msg)
	}
	if _, err := run(t, "-c", cfg, "deprecate", "--record", chkID, "--reason", "again"); err == nil {
		t.Fatal("deprecating twice passed")
	}
}

func TestDeprecateRefusalsAndRollback(t *testing.T) {
	root, cfg := completeRepo(t, "commit_on_status_change: true\n", "true\n")
	for name, args := range map[string][]string{
		"no record": {"deprecate", "--reason", "x"},
		"no reason": {"deprecate", "--record", chkID},
	} {
		if out, err := run(t, append([]string{"-c", cfg}, args...)...); err == nil {
			t.Errorf("%s: accepted\n%s", name, out)
		}
	}
	hooks := filepath.Join(root, ".git", "reject-hooks")
	writeFile(t, filepath.Join(hooks, "pre-commit"), "#!/bin/sh\nexit 1\n")
	if err := os.Chmod(filepath.Join(hooks, "pre-commit"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "config", "core.hooksPath", hooks)
	before := recordFile(t, root, chkID)
	if _, err := run(t, "-c", cfg, "deprecate", "--record", chkID, "--reason", "x"); err == nil {
		t.Fatal("deprecate passed a rejecting hook")
	}
	if recordFile(t, root, chkID) != before {
		t.Fatal("rejected deprecate changed the record")
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Fatalf("rollback left changes:\n%s", st)
	}
}
