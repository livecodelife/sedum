package main

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// statusRepo: an opened blueprint at version 2 (via edit), an imprint
// implementing it at version 1 with a tag, and a deprecated draft.
func statusRepo(t *testing.T) (root, cfg string) {
	t.Helper()
	child, dep := "prov-2026-0000d011", "prov-2026-0000d012"
	root, cfg = openedRepo(t, draftID)
	mustRun(t, "-c", cfg, "edit", draftID, "--add-constraint", "More", "--reason", "Grew", "--fix")
	childRec := strings.Replace(v2Draft(child, "draft"), "type: blueprint\n", "type: imprint\nimplements: "+draftID+"\nimplements_version: 1\n", 1)
	childRec = strings.Replace(childRec, "tags: []", "tags:\n  - cli", 1)
	writeFile(t, filepath.Join(root, "provenance", child+".yml"), childRec)
	mustRun(t, "-c", cfg, "open", "--record", child)
	writeFile(t, filepath.Join(root, "provenance", dep+".yml"), v2Draft(dep, "deprecated"))
	return root, cfg
}

func TestStatusTable(t *testing.T) {
	_, cfg := statusRepo(t)
	out := mustRun(t, "-c", cfg, "status")
	for _, want := range []string{"PROVENANCE RECORDS", draftID + "@2", "prov-2026-0000d011@1", "prov-2026-0000d012@1", "open ⚠", "2 open records have no associated LineSpecs"} {
		if !strings.Contains(out, want) {
			t.Errorf("table lacks %q:\n%s", want, out)
		}
	}
	open := mustRun(t, "-c", cfg, "status", "--filter", "open")
	if strings.Contains(open, "prov-2026-0000d012") || !strings.Contains(open, draftID) {
		t.Fatalf("--filter open:\n%s", open)
	}
	tagged := mustRun(t, "-c", cfg, "status", "--filter", "tag:cli")
	if !strings.Contains(tagged, "prov-2026-0000d011") || strings.Contains(tagged, draftID+"@") {
		t.Fatalf("--filter tag:cli:\n%s", tagged)
	}
	if _, err := run(t, "-c", cfg, "status", "--filter", "superseded"); err == nil {
		t.Fatal("a v1-only status was accepted as a filter")
	}
}

func TestStatusDetail(t *testing.T) {
	_, cfg := statusRepo(t)
	out := mustRun(t, "-c", cfg, "status", "--record", draftID)
	for _, want := range []string{draftID + "@2", "Version:      2 (fix: Grew)", "Implemented by:", "prov-2026-0000d011  [imprint]", "More"} {
		if !strings.Contains(out, want) {
			t.Errorf("detail lacks %q:\n%s", want, out)
		}
	}
	child := mustRun(t, "-c", cfg, "status", "--record", "prov-2026-0000d011")
	if !strings.Contains(child, "Implements:   "+draftID+"@1  Draft "+draftID) {
		t.Fatalf("child detail:\n%s", child)
	}
	if _, err := run(t, "-c", cfg, "status", "--record", "prov-2026-0000ffff"); err == nil {
		t.Fatal("unknown record accepted")
	}
}

func TestStatusJSON(t *testing.T) {
	_, cfg := statusRepo(t)
	var all struct {
		Records []map[string]any `json:"records"`
	}
	if err := json.Unmarshal([]byte(mustRun(t, "-c", cfg, "status", "--format", "json")), &all); err != nil {
		t.Fatal(err)
	}
	if len(all.Records) != 3 {
		t.Fatalf("%d records", len(all.Records))
	}
	var one struct {
		Record map[string]any `json:"record"`
	}
	if err := json.Unmarshal([]byte(mustRun(t, "-c", cfg, "status", "--record", draftID, "--format", "json")), &one); err != nil {
		t.Fatal(err)
	}
	r := one.Record
	if r["id"] != draftID || r["version"] != float64(2) || r["reason"] != "Grew" || r["reason_kind"] != "fix" || r["status"] != "open" {
		t.Fatalf("record json %v", r)
	}
	if cons, _ := r["constraints"].([]any); len(cons) != 2 {
		t.Fatalf("constraints %v", r["constraints"])
	}
}
