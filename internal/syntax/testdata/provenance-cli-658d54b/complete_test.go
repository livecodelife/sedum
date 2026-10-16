package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/provenance/record"
)

// completeRepo is checkRepo with extra config lines and the spec script
// body given.
func completeRepo(t *testing.T, extraConfig, spec string) (root, cfg string) {
	t.Helper()
	root, cfg = checkRepoSpec(t, spec)
	writeFile(t, cfg, readFile(t, cfg)+extraConfig)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "--allow-empty", "-m", "setup ["+chkID+"]")
	return root, cfg
}

func recordFile(t *testing.T, root, id string) string {
	return readFile(t, filepath.Join(root, "provenance", id+".yml"))
}

func TestCompleteSeals(t *testing.T) {
	root, cfg := completeRepo(t, "", "true\n")
	head := gitIn(t, root, "rev-parse", "HEAD")
	manifest := readFile(t, filepath.Join(root, manifestPath))
	out := mustRun(t, "-c", cfg, "complete", "--record", chkID)
	if !strings.Contains(out, "✓ "+chkID+" marked as implemented") {
		t.Fatalf("output:\n%s", out)
	}
	rec := load(t, root, chkID)
	if rec.Status != record.StatusImplemented || rec.SealedAtSHA != head || rec.Version != 1 {
		t.Fatalf("status %s sealed %q version %d", rec.Status, rec.SealedAtSHA, rec.Version)
	}
	if readFile(t, filepath.Join(root, manifestPath)) != manifest {
		t.Fatal("complete changed the version manifest")
	}
	if _, err := run(t, "-c", cfg, "complete", "--record", chkID); err == nil {
		t.Fatal("completing an implemented record again passed")
	}
}

func TestCompleteRefusals(t *testing.T) {
	draft := "prov-2026-0000c002"
	root, cfg := completeRepo(t, "", "true\n")
	writeFile(t, filepath.Join(root, "provenance", draft+".yml"), v2Draft(draft, "draft"))
	// A draft cannot be completed; a missing spec blocks unless --force.
	if _, err := run(t, "-c", cfg, "complete", "--record", draft); err == nil {
		t.Error("completing a draft passed")
	}
	if _, err := run(t, "-c", cfg, "complete"); err == nil {
		t.Error("complete without --record passed")
	}
	// A record whose spec never existed (its version 1 already hashed it as
	// missing) is refused unless --force. Deleting an existing spec would be
	// a content change instead, which lint blocks even with --force.
	nospec := "prov-2026-0000c004"
	rec := strings.Replace(v2Draft(nospec, "draft"), "associated_specs: []\n", "associated_specs:\n  - path: never.sh\n    run_command: sh\n", 1)
	writeFile(t, filepath.Join(root, "provenance", nospec+".yml"), rec)
	mustRun(t, "-c", cfg, "open", "--record", nospec)
	before := recordFile(t, root, nospec)
	if out, err := run(t, "-c", cfg, "complete", "--record", nospec); err == nil {
		t.Fatalf("a missing spec passed:\n%s", out)
	}
	if recordFile(t, root, nospec) != before {
		t.Fatal("refused complete changed the record")
	}
	mustRun(t, "-c", cfg, "complete", "--record", nospec, "--force")
	if load(t, root, nospec).Status != record.StatusImplemented {
		t.Fatal("--force did not complete")
	}
}

func TestCompleteRunsSpecsAndRollsBack(t *testing.T) {
	// The spec passes only once ok.flag exists, so its content (part of the
	// record's content_hash) never changes.
	root, cfg := completeRepo(t, "run_associated_specs_on_complete: true\n", "test -f ok.flag\n")
	before := recordFile(t, root, chkID)
	out, err := run(t, "-c", cfg, "complete", "--record", chkID)
	if err == nil {
		t.Fatalf("a failing spec passed:\n%s", out)
	}
	if recordFile(t, root, chkID) != before {
		t.Fatal("failed complete changed the record")
	}
	writeFile(t, filepath.Join(root, "ok.flag"), "")
	out = mustRun(t, "-c", cfg, "complete", "--record", chkID)
	if !strings.Contains(out, "Running 1 associated spec(s)") || !strings.Contains(out, "✓ passed") {
		t.Fatalf("specs not run:\n%s", out)
	}
}

func TestCompleteLintGate(t *testing.T) {
	root, cfg := completeRepo(t, "", "true\n")
	path := filepath.Join(root, "provenance", chkID+".yml")
	writeFile(t, path, strings.Replace(readFile(t, path), "title: Draft "+chkID+"\n", "title: \"\"\n", 1))
	before := recordFile(t, root, chkID)
	if out, err := run(t, "-c", cfg, "complete", "--record", chkID); err == nil {
		t.Fatalf("a record failing lint completed:\n%s", out)
	}
	if recordFile(t, root, chkID) != before {
		t.Fatal("lint-gated complete changed the record")
	}
}

func TestCompleteAutoCommits(t *testing.T) {
	root, cfg := completeRepo(t, "commit_on_status_change: true\n", "true\n")
	mustRun(t, "-c", cfg, "complete", "--record", chkID)
	if msg := gitIn(t, root, "log", "-1", "--format=%s"); msg != "Complete provenance record ["+chkID+"@1]" {
		t.Fatalf("commit %q", msg)
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Fatalf("left changes:\n%s", st)
	}
}

func TestCompleteRollsBackRejectedCommit(t *testing.T) {
	root, cfg := completeRepo(t, "commit_on_status_change: true\n", "true\n")
	hooks := filepath.Join(root, ".git", "reject-hooks")
	writeFile(t, filepath.Join(hooks, "pre-commit"), "#!/bin/sh\nexit 1\n")
	if err := os.Chmod(filepath.Join(hooks, "pre-commit"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "config", "core.hooksPath", hooks)
	before := recordFile(t, root, chkID)
	if out, err := run(t, "-c", cfg, "complete", "--record", chkID); err == nil {
		t.Fatalf("complete passed a rejecting hook:\n%s", out)
	}
	if recordFile(t, root, chkID) != before {
		t.Fatal("rejected complete changed the record")
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Fatalf("rollback left changes or staged files:\n%s", st)
	}
}

// overlapRepo: sealed record B governs a.go and has a spec that fails; A's
// tagged commit changed a.go.
func overlapRepo(t *testing.T, mode string) (root, cfg string) {
	t.Helper()
	sealed := "prov-2026-0000c003"
	root, cfg = completeRepo(t, "run_associated_specs_on_complete: true\noverlap_specs_on_complete: "+mode+"\n", "true\n")
	b := strings.Replace(v2Draft(sealed, "implemented"), `sealed_at_sha: ""`, "sealed_at_sha: "+gitIn(t, root, "rev-parse", "HEAD"), 1)
	b = strings.Replace(b, "associated_specs: []\n", "associated_specs:\n  - path: b-spec.sh\n    run_command: sh\n", 1)
	writeFile(t, filepath.Join(root, "provenance", sealed+".yml"), b)
	writeFile(t, filepath.Join(root, "b-spec.sh"), "exit 1\n")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "sealed b")
	writeFile(t, filepath.Join(root, "a.go"), "package a // touched\n")
	gitIn(t, root, "add", "a.go")
	gitIn(t, root, "commit", "-q", "-m", "touch a ["+chkID+"@1]")
	return root, cfg
}

func TestCompleteOverlapTeeth(t *testing.T) {
	root, cfg := overlapRepo(t, "block")
	if out, err := run(t, "-c", cfg, "complete", "--record", chkID); err == nil || !strings.Contains(out, "Verifying sealed record prov-2026-0000c003") {
		t.Fatalf("block mode let a broken sealed record through:\n%s", out)
	}
	if load(t, root, chkID).Status != record.StatusOpen {
		t.Fatal("blocked complete did not roll back")
	}

	_, cfg = overlapRepo(t, "warn")
	if out := mustRun(t, "-c", cfg, "complete", "--record", chkID); !strings.Contains(out, "warn mode") {
		t.Fatalf("warn mode did not report:\n%s", out)
	}

	_, cfg = overlapRepo(t, "off")
	if out := mustRun(t, "-c", cfg, "complete", "--record", chkID); strings.Contains(out, "prov-2026-0000c003") {
		t.Fatalf("off mode ran the overlap teeth:\n%s", out)
	}
}
