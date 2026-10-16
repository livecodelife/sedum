package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"github.com/livecodelife/sedum/internal/record"}}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// newSelectFixture builds a generators tree with one python package, an\n// imprint record, and an empty root directory for files under test.\nfunc newSelectFixture(t *testing.T) (root string, set *genpkg.Set, records *record.Set) {\n\tt.Helper()\n\tgen := t.TempDir()\n\twrite := func(rel, body string) {\n\t\tt.Helper()\n\t\tfull := filepath.Join(gen, filepath.FromSlash(rel))\n\t\tif err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t\tif err := os.WriteFile(full, []byte(body), 0o644); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\twrite(\"py/sedum.yaml\", \"name: py\\nextensions: [\\\".py\\\"]\\ncomment_prefix: \\\"#\\\"\\ntest_paths: [\\\"test_*.py\\\"]\\nunmanaged: [\\\"setup.py\\\"]\\n\")\n\twrite(\"py/files/{name}.py\", \"# sedum:anchor:decls\\n\")\n\twrite(\"py/actions/actions.yaml\", \"actions:\\n  setThing:\\n    kwargs:\\n      file: { type: string, required: true }\\n      name: { type: string, required: true }\\n      code: { type: string, required: true }\\n    injects_into: \\\"{{file}}\\\"\\n    anchor: decls\\n\")\n\twrite(\"py/actions/setThing.py\", \"{{code}}\\n\")\n\tset, findings, err := genpkg.Load(gen, genpkg.Options{})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(findings) > 0 {\n\t\tt.Fatalf(\"generators: %+v\", findings)\n\t}\n\trecords = &record.Set{Records: []*record.Record{{ID: \"imp-1\", Type: \"imprint\", Patterns: []string{\"*.py\"}}}}\n\troot = t.TempDir()\n\treturn root, set, records\n}","file":"internal/adopt/select_test.go","name":"newSelectFixture"}}
// newSelectFixture builds a generators tree with one python package, an
// imprint record, and an empty root directory for files under test.
func newSelectFixture(t *testing.T) (root string, set *genpkg.Set, records *record.Set) {
	t.Helper()
	gen := t.TempDir()
	write := func(rel, body string) {
		t.Helper()
		full := filepath.Join(gen, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("py/sedum.yaml", "name: py\nextensions: [\".py\"]\ncomment_prefix: \"#\"\ntest_paths: [\"test_*.py\"]\nunmanaged: [\"setup.py\"]\n")
	write("py/files/{name}.py", "# sedum:anchor:decls\n")
	write("py/actions/actions.yaml", "actions:\n  setThing:\n    kwargs:\n      file: { type: string, required: true }\n      name: { type: string, required: true }\n      code: { type: string, required: true }\n    injects_into: \"{{file}}\"\n    anchor: decls\n")
	write("py/actions/setThing.py", "{{code}}\n")
	set, findings, err := genpkg.Load(gen, genpkg.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) > 0 {
		t.Fatalf("generators: %+v", findings)
	}
	records = &record.Set{Records: []*record.Record{{ID: "imp-1", Type: "imprint", Patterns: []string{"*.py"}}}}
	root = t.TempDir()
	return root, set, records
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSelectCandidate(t *testing.T) {\n\troot, set, records := newSelectFixture(t)\n\tbody := []byte(\"print('hi')\\n\")\n\tif err := os.WriteFile(filepath.Join(root, \"a.py\"), body, 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tcands, skips, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{\"a.py\"})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(skips) != 0 {\n\t\tt.Fatalf(\"unexpected skips: %+v\", skips)\n\t}\n\tif len(cands) != 1 {\n\t\tt.Fatalf(\"want 1 candidate, got %d\", len(cands))\n\t}\n\tc := cands[0]\n\tif c.Template != \"{name}.py\" {\n\t\tt.Errorf(\"Template = %q, want {name}.py\", c.Template)\n\t}\n\tif c.Record == nil || c.Record.ID != \"imp-1\" {\n\t\tt.Errorf(\"Record = %+v, want ID imp-1\", c.Record)\n\t}\n\tif c.Captures[\"name\"] != \"a\" {\n\t\tt.Errorf(\"Captures[name] = %q, want a\", c.Captures[\"name\"])\n\t}\n\tif !bytes.Equal(c.Src, body) {\n\t\tt.Errorf(\"Src = %q, want %q\", c.Src, body)\n\t}\n}","file":"internal/adopt/select_test.go","name":"TestSelectCandidate"}}
func TestSelectCandidate(t *testing.T) {
	root, set, records := newSelectFixture(t)
	body := []byte("print('hi')\n")
	if err := os.WriteFile(filepath.Join(root, "a.py"), body, 0o644); err != nil {
		t.Fatal(err)
	}
	cands, skips, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{"a.py"})
	if err != nil {
		t.Fatal(err)
	}
	if len(skips) != 0 {
		t.Fatalf("unexpected skips: %+v", skips)
	}
	if len(cands) != 1 {
		t.Fatalf("want 1 candidate, got %d", len(cands))
	}
	c := cands[0]
	if c.Template != "{name}.py" {
		t.Errorf("Template = %q, want {name}.py", c.Template)
	}
	if c.Record == nil || c.Record.ID != "imp-1" {
		t.Errorf("Record = %+v, want ID imp-1", c.Record)
	}
	if c.Captures["name"] != "a" {
		t.Errorf("Captures[name] = %q, want a", c.Captures["name"])
	}
	if !bytes.Equal(c.Src, body) {
		t.Errorf("Src = %q, want %q", c.Src, body)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSelectSkipReasons(t *testing.T) {\n\ttests := []struct {\n\t\tname string\n\t\tfile string\n\t\tbody string\n\t\twant Reason\n\t}{\n\t\t{\"unmanaged\", \"setup.py\", \"x\\n\", ReasonUnmanaged},\n\t\t{\"test\", \"test_a.py\", \"x\\n\", ReasonTest},\n\t\t{\"generated\", \"gen.py\", \"# DO NOT EDIT\\n\", ReasonGenerated},\n\t\t{\"adopted\", \"done.py\", \"# sedum:setThing {}\\n\", ReasonAlreadyAdopted},\n\t\t{\"no package\", \"notes.txt\", \"x\\n\", ReasonNoPackage},\n\t}\n\tfor _, tc := range tests {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\troot, set, records := newSelectFixture(t)\n\t\t\tif err := os.WriteFile(filepath.Join(root, tc.file), []byte(tc.body), 0o644); err != nil {\n\t\t\t\tt.Fatal(err)\n\t\t\t}\n\t\t\t_, skips, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{tc.file})\n\t\t\tif err != nil {\n\t\t\t\tt.Fatal(err)\n\t\t\t}\n\t\t\tif len(skips) != 1 || skips[0].Reason != tc.want {\n\t\t\t\tt.Fatalf(\"skips = %+v, want one skip with reason %q\", skips, tc.want)\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/adopt/select_test.go","name":"TestSelectSkipReasons"}}
func TestSelectSkipReasons(t *testing.T) {
	tests := []struct {
		name string
		file string
		body string
		want Reason
	}{
		{"unmanaged", "setup.py", "x\n", ReasonUnmanaged},
		{"test", "test_a.py", "x\n", ReasonTest},
		{"generated", "gen.py", "# DO NOT EDIT\n", ReasonGenerated},
		{"adopted", "done.py", "# sedum:setThing {}\n", ReasonAlreadyAdopted},
		{"no package", "notes.txt", "x\n", ReasonNoPackage},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root, set, records := newSelectFixture(t)
			if err := os.WriteFile(filepath.Join(root, tc.file), []byte(tc.body), 0o644); err != nil {
				t.Fatal(err)
			}
			_, skips, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{tc.file})
			if err != nil {
				t.Fatal(err)
			}
			if len(skips) != 1 || skips[0].Reason != tc.want {
				t.Fatalf("skips = %+v, want one skip with reason %q", skips, tc.want)
			}
		})
	}
}

// /sedum:setFunction
