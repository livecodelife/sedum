package genpkg

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The go package's convergence block must reproduce, byte for byte, what graft
// wrote before the contract existed. The golden files are graft's own
// (internal/recording/testdata/golden at 9ba3040) and are not to be edited to
// fit a change here.
func TestGoConvergenceMatchesGraftGolden(t *testing.T) {
	dir := filepath.Join("testdata", "convergence", "go")
	read := func(name string) []byte {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	var raw struct {
		Files []json.RawMessage `json:"files"`
	}
	if err := json.Unmarshal(read("skeleton.json"), &raw); err != nil {
		t.Fatal(err)
	}
	var in ConvRenderInput
	for _, r := range raw.Files {
		var f ConvFile // decls are graft's business; unknown keys are ignored here
		if err := json.Unmarshal(r, &f); err != nil {
			t.Fatal(err)
		}
		in.Files = append(in.Files, f)
	}

	set, findings, err := Load(filepath.Join("..", "..", "generators"), Options{Only: []string{"go"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range findings {
		if f.Kind == KindError {
			t.Fatalf("go package: %s", f)
		}
	}
	pkg, ok := set.Lookup("go")
	if !ok {
		t.Fatal("no go package")
	}
	out, err := pkg.RenderConvergence(in)
	if err != nil {
		t.Fatal(err)
	}

	// Stub code: the golden recording's setFunction/setMethod invocations.
	var rec struct {
		Records []struct {
			Phases []struct {
				Invocations []struct {
					Action string            `json:"action"`
					Kwargs map[string]string `json:"kwargs"`
				} `json:"invocations"`
			} `json:"phases"`
		} `json:"records"`
	}
	if err := json.Unmarshal(read("recording.json"), &rec); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{} // file|receiver|name -> code
	for _, r := range rec.Records {
		for _, p := range r.Phases {
			for _, iv := range p.Invocations {
				if iv.Action == "setFunction" || iv.Action == "setMethod" {
					want[iv.Kwargs["file"]+"|"+iv.Kwargs["receiver"]+"|"+iv.Kwargs["name"]] = iv.Kwargs["code"]
				}
			}
		}
	}
	checked := 0
	for i, f := range out.Files {
		if len(f.Stubs) == 0 {
			continue
		}
		for j, s := range f.Stubs {
			recv := strings.TrimPrefix(in.Files[i].Stubs[j].ReceiverType, "*")
			code, ok := want[f.Path+"|"+recv+"|"+s.Name]
			if !ok {
				t.Fatalf("%s %s: golden has no such invocation", f.Path, s.Name)
			}
			if s.Kwargs["code"] != code || s.Kwargs["file"] != f.Path || s.Kwargs["name"] != s.Name {
				t.Errorf("%s %s kwargs differ: %v", f.Path, s.Name, s.Kwargs)
			}
			if wantAct := map[bool]string{true: "setMethod", false: "setFunction"}[recv != ""]; s.Action != wantAct || (recv != "") != (s.Kwargs["receiver"] == recv && recv != "") {
				t.Errorf("%s %s action %s kwargs %v", f.Path, s.Name, s.Action, s.Kwargs)
			}
			if s.Stub != code {
				t.Errorf("%s %s stub differs:\n got: %q\nwant: %q", f.Path, s.Name, s.Stub, code)
			}
			checked++
		}
		// Seed: golden file minus the two marker lines.
		g := string(read("seed-" + strings.TrimSuffix(f.SynthPath, ".go") + ".go"))
		lines := strings.SplitAfter(strings.TrimSuffix(g, "\n"), "\n")
		body := strings.Join(lines[1:len(lines)-1], "")
		if body != f.Seed {
			t.Errorf("%s seed differs:\n got: %q\nwant: %q", f.SynthPath, f.Seed, body)
		}
	}
	if checked == 0 {
		t.Fatal("nothing compared")
	}
}

func TestGoConvergenceFormatsAndPicksMain(t *testing.T) {
	pkg := loadConvPackage(t, "go")
	out, err := pkg.RenderConvergence(ConvRenderInput{Files: []ConvFile{
		{Path: "cmd/tool/sub/a.go", Imports: []string{"net"}, Stubs: []ConvStub{{
			Name: "f", Params: []ConvParam{{Name: "m", Type: "map[string] int"}, {Name: "c", Type: "net.Conn"}, {Name: "rest", Type: "...string"}}, Results: []string{"error"}}}},
		{Path: "x/my-mod.go", Stubs: []ConvStub{{Name: "g", Doc: " "}}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	a := out.Files[0]
	if !strings.HasPrefix(a.Seed, "package main\n") {
		t.Errorf("cmd/tool/sub is a main package: %q", a.Seed)
	}
	if !strings.Contains(a.Seed, "map[string]int") || strings.Contains(a.Seed, "map[string] int") {
		t.Errorf("seed is not gofmt-clean: %q", a.Seed)
	}
	if !strings.Contains(a.Stubs[0].Stub, "aFBody(m, c, rest...)") {
		t.Errorf("stub = %q", a.Stubs[0].Stub)
	}
	b := out.Files[1]
	if !strings.Contains(b.Stubs[0].Helper, "myModGBody") || !strings.HasPrefix(b.Stubs[0].Stub, "//\n") {
		t.Errorf("hyphenated stem / blank doc: %+v", b.Stubs[0])
	}
}
