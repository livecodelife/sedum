package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/provenance/record"
)

const draftID = "prov-2026-0000d001"

// v2Draft is a schema 2 draft record with no snapshot.
func v2Draft(id, status string) string {
	return "id: " + id + "\ntitle: Draft " + id + "\nstatus: " + status + "\nversion: 1\nreason: \"\"\ncreated_at: \"2026-10-15\"\nauthor: a@b\n" +
		"intent: Do the thing\nconstraints:\n  - Keep it small\naffected_scope:\n  - a.go\nforbidden_scope: []\ntype: blueprint\n" +
		"related: []\nsealed_at_sha: \"\"\nassociated_specs: []\nassociated_traces: []\nmonitors: []\ntags: []\n"
}

// v2Repo is a migrated repository holding records, committed when commit is
// set (commit_on_status_change then applies too).
func v2Repo(t *testing.T, commit bool, records map[string]string) (root, cfg string) {
	t.Helper()
	root = t.TempDir()
	cfg = filepath.Join(root, ".provenance.yml")
	body := "dir: provenance\n"
	if commit {
		body += "commit_on_status_change: true\n"
	}
	writeFile(t, cfg, body)
	for id, rec := range records {
		writeFile(t, filepath.Join(root, "provenance", id+".yml"), rec)
	}
	if commit {
		gitRepo(t, root)
	}
	return root, cfg
}

func TestOpenWritesVersionOne(t *testing.T) {
	root, cfg := v2Repo(t, false, map[string]string{draftID: v2Draft(draftID, "draft")})
	out := mustRun(t, "-c", cfg, "open", "--record", draftID)
	if !strings.Contains(out, "✓ "+draftID+" transitioned from draft → open") {
		t.Fatalf("output:\n%s", out)
	}
	store := &record.FilesStore{Root: root}
	rec, _, err := store.Load(draftID)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status != record.StatusOpen || rec.Version != 1 {
		t.Fatalf("status %s version %d", rec.Status, rec.Version)
	}
	snap, err := store.Snapshot(draftID, 1)
	if err != nil {
		t.Fatalf("no v1 snapshot: %v", err)
	}
	if snap.Author != "test@example.com" || snap.Reason != "" {
		t.Fatalf("snapshot author %q reason %q", snap.Author, snap.Reason)
	}
	manifest := readFile(t, filepath.Join(root, manifestPath))
	if manifest != `{"id":"`+draftID+`","version":1,"content_hash":"`+snap.ContentHash+`"}`+"\n" {
		t.Fatalf("manifest:\n%s", manifest)
	}
	// The record file itself changed only in its status line.
	if got := readFile(t, filepath.Join(root, "provenance", draftID+".yml")); got != v2Draft(draftID, "open") {
		t.Fatalf("record file:\n%s", got)
	}
}

func TestOpenRefusals(t *testing.T) {
	v1ID := "prov-2026-0000d002"
	openID := "prov-2026-0000d003"
	root, cfg := v2Repo(t, false, map[string]string{
		v1ID:   v1(v1ID, "blueprint", "draft", ""),
		openID: v2Draft(openID, "open"),
	})
	before := readFile(t, filepath.Join(root, "provenance", openID+".yml"))
	for name, args := range map[string][]string{
		"no --record": {"open"},
		"unknown id":  {"open", "--record", "prov-2026-0000ffff"},
		"v1 record":   {"open", "--record", v1ID},
		"not a draft": {"open", "--record", openID},
		"stray arg":   {"open", "--record", openID, "extra"},
	} {
		if out, err := run(t, append([]string{"-c", cfg}, args...)...); err == nil {
			t.Errorf("%s: accepted\n%s", name, out)
		}
	}
	if got := readFile(t, filepath.Join(root, "provenance", openID+".yml")); got != before {
		t.Fatalf("refused open changed the record:\n%s", got)
	}
	if _, err := os.Stat(filepath.Join(root, "provenance", ".versions")); !os.IsNotExist(err) {
		t.Fatalf("refused opens wrote snapshots (%v)", err)
	}
}

func TestOpenAutoCommits(t *testing.T) {
	root, cfg := v2Repo(t, true, map[string]string{draftID: v2Draft(draftID, "draft")})
	mustRun(t, "-c", cfg, "open", "--record", draftID)
	if msg := gitIn(t, root, "log", "-1", "--format=%s"); msg != "Open provenance record "+draftID+" ["+draftID+"@1]" {
		t.Fatalf("commit message %q", msg)
	}
	files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD")
	for _, want := range []string{"provenance/" + draftID + ".yml", "provenance/.versions/" + draftID + "/", manifestPath} {
		if !strings.Contains(files, want) {
			t.Errorf("commit lacks %s:\n%s", want, files)
		}
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Fatalf("left changes behind:\n%s", st)
	}
}

func TestOpenRollsBackRejectedCommit(t *testing.T) {
	root, cfg := v2Repo(t, true, map[string]string{draftID: v2Draft(draftID, "draft")})
	hooks := filepath.Join(root, ".git", "reject-hooks")
	writeFile(t, filepath.Join(hooks, "pre-commit"), "#!/bin/sh\necho rejected by test hook >&2\nexit 1\n")
	if err := os.Chmod(filepath.Join(hooks, "pre-commit"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "config", "core.hooksPath", hooks)
	if out, err := run(t, "-c", cfg, "open", "--record", draftID); err == nil {
		t.Fatalf("open succeeded past a rejecting hook\n%s", out)
	}
	if got := readFile(t, filepath.Join(root, "provenance", draftID+".yml")); got != v2Draft(draftID, "draft") {
		t.Fatalf("record not rolled back:\n%s", got)
	}
	for _, p := range []string{filepath.Join("provenance", ".versions", draftID), manifestPath} {
		if _, err := os.Stat(filepath.Join(root, p)); !os.IsNotExist(err) {
			t.Errorf("%s survived the rollback (%v)", p, err)
		}
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Fatalf("rollback left changes or staged files:\n%s", st)
	}
}
