package main

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestLintClean(t *testing.T) {
	_, cfg := openedRepo(t, draftID)
	out := mustRun(t, "-c", cfg, "lint")
	// linespec's counts: total is passed records plus warning and error
	// issues; the record's one warning is constraints without specs.
	for _, want := range []string{"✓ Linting 2 provenance records (enforcement: none)", "✓ 1 passed  ⚠ 1 warnings  ✗ 0 errors"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

func TestLintErrorsFail(t *testing.T) {
	bad := "prov-2026-0000d021"
	root, cfg := openedRepo(t, draftID)
	writeFile(t, filepath.Join(root, "provenance", bad+".yml"), strings.Replace(v2Draft(bad, "draft"), "intent: Do the thing\n", "intent: \"\"\n", 1))
	out, err := run(t, "-c", cfg, "lint")
	if err == nil {
		t.Fatalf("lint passed with an error:\n%s", out)
	}
	if !strings.Contains(out, "✗ "+bad) || !strings.Contains(out, "(intent)") {
		t.Fatalf("output:\n%s", out)
	}
	// --record keeps only that record's issues and counts.
	if out := mustRun(t, "-c", cfg, "lint", "--record", draftID); strings.Contains(out, bad) || !strings.Contains(out, "✓ 0 passed  ⚠ 1 warnings  ✗ 0 errors") {
		t.Fatalf("--record %s:\n%s", draftID, out)
	}
	if _, err := run(t, "-c", cfg, "lint", "--record", bad); err == nil {
		t.Fatal("--record on the bad record passed")
	}
	// Under enforcement none... errors still fail; the override is reported.
	out, _ = run(t, "-c", cfg, "lint", "--enforcement", "warn")
	if !strings.Contains(out, "(enforcement: warn)") {
		t.Fatalf("--enforcement not applied:\n%s", out)
	}
}

func TestLintJSON(t *testing.T) {
	bad := "prov-2026-0000d021"
	root, cfg := openedRepo(t, draftID)
	writeFile(t, filepath.Join(root, "provenance", bad+".yml"), strings.Replace(v2Draft(bad, "draft"), "intent: Do the thing\n", "intent: \"\"\n", 1))
	out, _ := run(t, "-c", cfg, "lint", "--format", "json")
	var res struct {
		Enforcement string `json:"enforcement"`
		Total       int    `json:"total"`
		Errors      int    `json:"errors"`
		Warnings    int    `json:"warnings"`
		Issues      []struct {
			RecordID string `json:"record_id"`
			Field    string `json:"field"`
			Message  string `json:"message"`
			Severity string `json:"severity"`
			Rule     string `json:"rule"`
		} `json:"issues"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	found := false
	for _, is := range res.Issues {
		found = found || (is.RecordID == bad && is.Rule == "PROV002" && is.Field == "intent" && is.Severity == "error")
	}
	if res.Total != res.Errors+res.Warnings+1 || res.Errors < 1 || !found {
		t.Fatalf("json %+v", res)
	}
}

func TestLintChecksManifest(t *testing.T) {
	root, cfg := openedRepo(t, draftID)
	writeFile(t, filepath.Join(root, manifestPath), "")
	out, err := run(t, "-c", cfg, "lint", "--all")
	if err == nil || !strings.Contains(out, "PROV026") && !strings.Contains(strings.ToLower(out), "manifest") {
		t.Fatalf("an empty manifest passed lint:\n%s", out)
	}
}
