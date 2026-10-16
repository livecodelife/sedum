package main

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestEditFields(t *testing.T) {
	root, cfg := openedRepo(t, draftID)
	path := filepath.Join(root, "provenance", draftID+".yml")
	writeFile(t, path, strings.Replace(readFile(t, path), "constraints:\n  - Keep it small\n", "constraints:\n  - Keep it small\n  - Second\n  - Third\n", 1))
	// The extra constraints are an unbumped edit; edit's own change makes the
	// version that covers both.
	out := mustRun(t, "-c", cfg, "edit", draftID,
		"--intent", "New intent",
		"--remove-constraint", "1", "--remove-constraint", "3", "--add-constraint", "Added, with a comma",
		"--remove-scope", "a.go", "--add-scope", "b.go", "--add-forbidden", "c.go",
		"--reason", "Edited")
	if !strings.Contains(out, draftID+" is now version 2") {
		t.Fatalf("output:\n%s", out)
	}
	rec := load(t, root, draftID)
	if rec.Intent != "New intent" || rec.Version != 2 || rec.Reason != "Edited" {
		t.Fatalf("record %+v", rec)
	}
	if !reflect.DeepEqual(rec.Constraints, []string{"Second", "Added, with a comma"}) {
		t.Fatalf("constraints %q", rec.Constraints)
	}
	if !reflect.DeepEqual(rec.AffectedScope, []string{"b.go"}) || !reflect.DeepEqual(rec.ForbiddenScope, []string{"c.go"}) {
		t.Fatalf("scope %v forbidden %v", rec.AffectedScope, rec.ForbiddenScope)
	}
	mustRun(t, "-c", cfg, "edit", "--record", draftID, "--remove-forbidden", "c.go", "--reason", "Drop it")
	if rec := load(t, root, draftID); len(rec.ForbiddenScope) != 0 || rec.Version != 3 {
		t.Fatalf("forbidden %v version %d", rec.ForbiddenScope, rec.Version)
	}
}

func TestEditTitleOnlyNoBump(t *testing.T) {
	root, cfg := openedRepo(t, draftID)
	mustRun(t, "-c", cfg, "edit", draftID, "--title", "Retitled")
	if rec := load(t, root, draftID); rec.Title != "Retitled" || rec.Version != 1 {
		t.Fatalf("title %q version %d", rec.Title, rec.Version)
	}
}

func TestEditRefusals(t *testing.T) {
	root, cfg := openedRepo(t, draftID)
	before := readFile(t, filepath.Join(root, "provenance", draftID+".yml"))
	for name, args := range map[string][]string{
		"nothing to edit":      {"edit", draftID, "--reason", "x"},
		"no reason":            {"edit", draftID, "--intent", "Changed"},
		"constraint 0":         {"edit", draftID, "--remove-constraint", "0", "--reason", "x"},
		"constraint too big":   {"edit", draftID, "--remove-constraint", "2", "--reason", "x"},
		"constraint not int":   {"edit", draftID, "--remove-constraint", "one", "--reason", "x"},
		"missing scope":        {"edit", draftID, "--remove-scope", "zzz.go", "--reason", "x"},
		"missing forbidden":    {"edit", draftID, "--remove-forbidden", "zzz.go", "--reason", "x"},
		"intent and bad scope": {"edit", draftID, "--intent", "Changed", "--remove-scope", "zzz.go", "--reason", "x"},
	} {
		if out, err := run(t, append([]string{"-c", cfg}, args...)...); err == nil {
			t.Errorf("%s: accepted\n%s", name, out)
		}
	}
	if got := readFile(t, filepath.Join(root, "provenance", draftID+".yml")); got != before || len(versions(t, root, draftID)) != 1 {
		t.Fatalf("refused edit wrote something")
	}
}
