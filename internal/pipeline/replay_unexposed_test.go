package pipeline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// hideAddField marks the generator's one action as hidden from the catalog,
// after a recording naming it was made.
func hideAddField(t *testing.T, generators string) {
	t.Helper()
	path := filepath.Join(generators, "rails", "actions", "actions.yaml")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(body, []byte("    exposed: false\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestReplayRejectsAHiddenActionByDefault(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)
	hideAddField(t, cfg.Generators)

	_, err := Replay(rec, replayConfig(t, cfg))
	if err == nil {
		t.Fatal("a recording naming a hidden action replayed")
	}
	if !strings.Contains(err.Error(), "unknown_action") {
		t.Errorf("the diagnostic is not the unknown-action rule: %v", err)
	}
}

func TestReplayAcceptsAHiddenActionWhenAllowed(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)
	hideAddField(t, cfg.Generators)

	rcfg := replayConfig(t, cfg)
	rcfg.AllowUnexposed = true
	if _, err := Replay(rec, rcfg); err != nil {
		t.Fatalf("a hidden action was rejected although allowed: %v", err)
	}
	for _, rel := range tree(t, rcfg.Output) {
		if strings.Contains(readFile(t, filepath.Join(rcfg.Output, filepath.FromSlash(rel))), "attribute :email") {
			return
		}
	}
	t.Error("the hidden action's injection is missing from the output")
}
