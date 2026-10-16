package main

import (
	"path/filepath"
	"strings"
	"testing"
)

const specID = "prov-2026-0000d031"

// specRecord is an open v2 record whose associated_specs are given as YAML
// list items.
func specRepo(t *testing.T, runOnComplete bool, specs string) (root, cfg string) {
	t.Helper()
	rec := strings.Replace(v2Draft(specID, "draft"), "associated_specs: []\n", "associated_specs:\n"+specs, 1)
	root, cfg = v2Repo(t, false, map[string]string{specID: rec})
	if runOnComplete {
		writeFile(t, cfg, "dir: provenance\nrun_associated_specs_on_complete: true\n")
	}
	for _, f := range []string{"a.txt", "b.txt", "c.txt"} {
		writeFile(t, filepath.Join(root, f), f+"\n")
	}
	return root, cfg
}

func TestRunSpecsOffIsSilent(t *testing.T) {
	_, cfg := specRepo(t, false, "  - path: a.txt\n    run_command: \"false\"\n")
	if out := mustRun(t, "-c", cfg, "run-specs", "--record", specID); strings.TrimSpace(out) != "" {
		t.Fatalf("run-specs with the setting off printed:\n%s", out)
	}
}

func TestRunSpecsPassAndBatch(t *testing.T) {
	// Two specs sharing a command run as one process with both paths; a
	// {{path}} command and a different command run alone; an unknown type
	// is skipped. Commands run in the repo root.
	_, cfg := specRepo(t, true,
		"  - path: a.txt\n    run_command: cat\n"+
			"  - path: b.txt\n    run_command: cat\n"+
			"  - path: c.txt\n    run_command: \"test -f {{path}}\"\n"+
			"  - path: notes.md\n    type: config\n")
	out := mustRun(t, "-c", cfg, "run-specs", "--record", specID)
	for _, want := range []string{
		"Running 4 associated spec(s) for " + specID,
		"  · a.txt, b.txt\n    cat a.txt b.txt",
		"a.txt\nb.txt\n",
		"    test -f c.txt",
		"notes.md  (skipped — no type or run_command)",
		"All associated specs passed.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

func TestRunSpecsFailureLocalizes(t *testing.T) {
	_, cfg := specRepo(t, true,
		"  - path: a.txt\n    run_command: cat\n"+
			"  - path: missing.txt\n    run_command: cat\n")
	out, err := run(t, "-c", cfg, "run-specs", "--record", specID)
	if err == nil {
		t.Fatalf("a failing spec passed:\n%s", out)
	}
	for _, want := range []string{"✗ batch failed — localizing per spec:", "✓ a.txt passed", "✗ missing.txt failed"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

func TestRunSpecsWithoutHookGitEnv(t *testing.T) {
	_, cfg := specRepo(t, true, "  - path: a.txt\n    run_command: \"test -z \\\"$GIT_DIR$GIT_INDEX_FILE\\\" && test -f\"\n")
	t.Setenv("GIT_DIR", "/nonexistent/.git")
	t.Setenv("GIT_INDEX_FILE", "/nonexistent/.git/index")
	if out, err := run(t, "-c", cfg, "run-specs", "--record", specID); err != nil {
		t.Fatalf("spec saw the hook's GIT_* variables:\n%s", out)
	}
}

func TestRunSpecsRefusals(t *testing.T) {
	_, cfg := specRepo(t, true, "  - path: a.txt\n    run_command: cat\n")
	for name, args := range map[string][]string{
		"no record":      {"run-specs"},
		"unknown record": {"run-specs", "--record", "prov-2026-0000ffff"},
	} {
		if out, err := run(t, append([]string{"-c", cfg}, args...)...); err == nil {
			t.Errorf("%s: accepted\n%s", name, out)
		}
	}
}
