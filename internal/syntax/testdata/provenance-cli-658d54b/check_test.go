package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const chkID = "prov-2026-0000c001"

// checkRepo is a git repo holding one open v2 record chkID (version 1
// committed, with its snapshot and manifest) whose affected_scope is a.go
// and spec.sh, with spec.sh as its associated spec. README.md is excluded.
func checkRepo(t *testing.T) (root, cfg string) {
	t.Helper()
	return checkRepoSpec(t, "true\n")
}

// checkRepoSpec is checkRepo with spec.sh holding spec when the record opens.
func checkRepoSpec(t *testing.T, spec string) (root, cfg string) {
	t.Helper()
	root = t.TempDir()
	cfg = filepath.Join(root, ".provenance.yml")
	writeFile(t, cfg, "dir: provenance\ncommit_tag_required: true\nexclude_paths: ['README.md']\n")
	rec := strings.Replace(v2Draft(chkID, "draft"), "affected_scope:\n  - a.go\n", "affected_scope:\n  - a.go\n  - spec.sh\n", 1)
	rec = strings.Replace(rec, "associated_specs: []\n", "associated_specs:\n  - path: spec.sh\n    run_command: sh\n", 1)
	writeFile(t, filepath.Join(root, "provenance", chkID+".yml"), rec)
	for f, body := range map[string]string{"a.go": "package a\n", "b.go": "package b\n", "spec.sh": spec, "README.md": "hi\n"} {
		writeFile(t, filepath.Join(root, f), body)
	}
	mustRun(t, "-c", cfg, "open", "--record", chkID)
	gitRepo(t, root)
	return root, cfg
}

// check runs check --staged with message and reports whether it passed.
func check(t *testing.T, cfg, message string) (string, bool) {
	t.Helper()
	msg := filepath.Join(filepath.Dir(cfg), ".git", "TEST_MSG")
	writeFile(t, msg, message+"\n")
	out, err := run(t, "-c", cfg, "check", "--staged", "--message-file", msg)
	return out, err == nil
}

func stage(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for f, body := range files {
		writeFile(t, filepath.Join(root, f), body)
		gitIn(t, root, "add", "--", f)
	}
}

func TestCheckTagsAndScope(t *testing.T) {
	root, cfg := checkRepo(t)
	stage(t, root, map[string]string{"a.go": "package a // changed\n"})
	for _, msg := range []string{"Change a [" + chkID + "]", "Change a [" + chkID + "@1]"} {
		if out, ok := check(t, cfg, msg); !ok || !strings.Contains(out, "✓ No forbidden scope violations in staged") {
			t.Errorf("%q refused:\n%s", msg, out)
		}
	}
	if out, ok := check(t, cfg, "Change a ["+chkID+"@2]"); ok {
		t.Errorf("a tag naming a version that does not exist passed:\n%s", out)
	}
	if out, ok := check(t, cfg, "No tag here"); ok || !strings.Contains(out, "Commit tag required") {
		t.Errorf("untagged commit passed:\n%s", out)
	}
	stage(t, root, map[string]string{"b.go": "package b // changed\n"})
	if out, ok := check(t, cfg, "Change b ["+chkID+"]"); ok || !strings.Contains(out, "b.go") {
		t.Errorf("out-of-scope b.go passed:\n%s", out)
	}
}

func TestCheckExcludedNeedsNoTag(t *testing.T) {
	root, cfg := checkRepo(t)
	stage(t, root, map[string]string{"README.md": "hello\n"})
	if out, ok := check(t, cfg, "Docs"); !ok {
		t.Fatalf("excluded-only commit refused:\n%s", out)
	}
}

func TestCheckVersionTags(t *testing.T) {
	root, cfg := checkRepo(t)
	mustRun(t, "-c", cfg, "edit", chkID, "--add-constraint", "Second", "--reason", "More")
	gitIn(t, root, "add", "-A")
	if out, ok := check(t, cfg, "Revise ["+chkID+"@2]"); !ok {
		t.Fatalf("the revise commit tagged @2 refused:\n%s", out)
	}
	if out, ok := check(t, cfg, "Revise ["+chkID+"@1]"); ok {
		t.Fatalf("tag @1 on a commit that stages version 2 passed:\n%s", out)
	}
}

func TestCheckSnapshotsAppendOnly(t *testing.T) {
	root, cfg := checkRepo(t)
	snaps, _ := filepath.Glob(filepath.Join(root, "provenance", ".versions", chkID, "1-*.yml"))
	if len(snaps) != 1 {
		t.Fatalf("snapshots %v", snaps)
	}
	if err := os.Chmod(snaps[0], 0o644); err != nil {
		t.Fatal(err)
	}
	rel, _ := filepath.Rel(root, snaps[0])
	stage(t, root, map[string]string{rel: readFile(t, snaps[0]) + "# tampered\n"})
	if out, ok := check(t, cfg, "Tamper ["+chkID+"]"); ok || !strings.Contains(out, "append-only") {
		t.Fatalf("a modified snapshot passed:\n%s", out)
	}
}

func TestCheckImplementedAndTransitions(t *testing.T) {
	root, cfg := checkRepo(t)
	path := filepath.Join("provenance", chkID+".yml")
	implemented := strings.Replace(readFile(t, filepath.Join(root, path)), "status: open", "status: implemented", 1)
	// The completion commit itself passes.
	stage(t, root, map[string]string{path: implemented})
	if out, ok := check(t, cfg, "Complete provenance record ["+chkID+"@1]"); !ok {
		t.Fatalf("completion commit refused:\n%s", out)
	}
	gitIn(t, root, "commit", "-q", "-m", "complete")
	// After that, tagging it is refused with revise advice.
	stage(t, root, map[string]string{"a.go": "package a // later\n"})
	if out, ok := check(t, cfg, "Later ["+chkID+"]"); ok || !strings.Contains(out, "provenance revise") {
		t.Fatalf("a commit on an implemented record passed:\n%s", out)
	}
	gitIn(t, root, "reset", "-q")
	// implemented -> open without a new version is not a transition.
	stage(t, root, map[string]string{path: strings.Replace(implemented, "status: implemented", "status: open", 1)})
	if out, ok := check(t, cfg, "Reopen ["+chkID+"]"); ok {
		t.Fatalf("implemented -> open without a bump passed:\n%s", out)
	}
}

// TestCheckReadsTheIndex: the working tree holds an unstaged version 2, the
// index still version 1; check judges the index.
func TestCheckReadsTheIndex(t *testing.T) {
	root, cfg := checkRepo(t)
	stage(t, root, map[string]string{"a.go": "package a // staged\n"})
	mustRun(t, "-c", cfg, "edit", chkID, "--add-constraint", "Unstaged", "--reason", "Not staged")
	if out, ok := check(t, cfg, "Change a ["+chkID+"@1]"); !ok {
		t.Fatalf("check read the working tree, not the index:\n%s", out)
	}
}

func TestCheckLintsStagedSpecs(t *testing.T) {
	root, cfg := checkRepo(t)
	// An unchanged record with a spec lints clean: no false PROV024.
	stage(t, root, map[string]string{"a.go": "package a // changed\n"})
	if out, ok := check(t, cfg, "Change a ["+chkID+"]"); !ok {
		t.Fatalf("unchanged record with a spec refused:\n%s", out)
	}
	// A staged spec edit without a version bump is PROV024.
	stage(t, root, map[string]string{"spec.sh": "true # changed\n"})
	if out, ok := check(t, cfg, "Change spec ["+chkID+"]"); ok || !strings.Contains(out, "PROV024") {
		t.Fatalf("spec edit without a bump passed:\n%s", out)
	}
}

func TestCheckNeedsStaged(t *testing.T) {
	_, cfg := checkRepo(t)
	if _, err := run(t, "-c", cfg, "check"); err == nil {
		t.Fatal("check without --staged accepted")
	}
}
