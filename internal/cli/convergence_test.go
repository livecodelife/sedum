package cli

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func runConvergence(t *testing.T, stdin string, args ...string) (string, error) {
	t.Helper()
	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetIn(strings.NewReader(stdin))
	cmd.SetArgs(append([]string{"convergence"}, args...))
	err := cmd.Execute()
	return out.String(), err
}

func TestConvergenceInfoAndRender(t *testing.T) {
	gens := filepath.Join("..", "..", "generators")
	out, err := runConvergence(t, "", "info", "--generators", gens, "--package", "typescript")
	if err != nil {
		t.Fatalf("info: %v\n%s", err, out)
	}
	var info map[string]any
	if err := json.Unmarshal([]byte(out), &info); err != nil {
		t.Fatal(err)
	}
	if info["placeholder"] != "graft: not converged" || info["synth_suffix"] != "_synth" || info["comment_prefix"] != "//" {
		t.Errorf("info = %v", info)
	}

	out, err = runConvergence(t, `{"files":[{"path":"a.go","stubs":[{"name":"f"}]}]}`,
		"render", "--generators", gens, "--package", "go")
	if err != nil {
		t.Fatalf("render: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"synth_path": "a_synth.go"`) || !strings.Contains(out, "aFBody") {
		t.Errorf("render output = %s", out)
	}

	if _, err := runConvergence(t, `{"files":[],"bogus":1}`, "render", "--generators", gens, "--package", "go"); err == nil {
		t.Error("an unknown input key should be refused")
	}
}
