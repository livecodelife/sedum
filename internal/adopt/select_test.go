package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_test.go","path":"github.com/livecodelife/sedum/internal/record"}}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// newSelectFixture writes a small generators tree and returns a root directory,\n// a loaded generator set, and an imprint record set.\nfunc newSelectFixture(t *testing.T) (string, *genpkg.Set, *record.Set) {\n\tt.Helper()\n\tgen := t.TempDir()\n\twrite := func(rel, body string) {\n\t\tt.Helper()\n\t\tp := filepath.Join(gen, filepath.FromSlash(rel))\n\t\tif err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t\tif err := os.WriteFile(p, []byte(body), 0o644); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\twrite(\"py/sedum.yaml\", `name: py\nextensions: [\".py\"]\ncomment_prefix: \"#\"\ntest_paths: [\"test_*.py\"]\nunmanaged: [\"setup.py\"]\n`)\n\twrite(\"py/files/{name}.py\", \"# sedum:anchor:decls\\n\")\n\twrite(\"py/files/*/{name}.py\", \"# sedum:anchor:decls\\n\")\n\twrite(\"py/actions/actions.yaml\", `actions:\n  setThing:\n    kwargs:\n      file: { type: string, required: true }\n      name: { type: string, required: true }\n      code: { type: string, required: true }\n    injects_into: \"{{file}}\"\n    anchor: decls\n`)\n\twrite(\"py/actions/setThing.py\", \"{{code}}\\n\")\n\tset, findings, err := genpkg.Load(gen, genpkg.Options{})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(findings) > 0 {\n\t\tt.Fatalf(\"generator findings: %v\", findings)\n\t}\n\troot := t.TempDir()\n\trecords := &record.Set{Records: []*record.Record{{ID: \"imp-1\", Type: \"imprint\", Patterns: []string{\"*.py\"}}}}\n\treturn root, set, records\n}","file":"internal/adopt/select_test.go","name":"newSelectFixture"}}
// newSelectFixture writes a small generators tree and returns a root directory,
// a loaded generator set, and an imprint record set.
func newSelectFixture(t *testing.T) (string, *genpkg.Set, *record.Set) {
	t.Helper()
	gen := t.TempDir()
	write := func(rel, body string) {
		t.Helper()
		p := filepath.Join(gen, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("py/sedum.yaml", `name: py
extensions: [".py"]
comment_prefix: "#"
test_paths: ["test_*.py"]
unmanaged: ["setup.py"]
`)
	write("py/files/{name}.py", "# sedum:anchor:decls\n")
	write("py/files/*/{name}.py", "# sedum:anchor:decls\n")
	write("py/actions/actions.yaml", `actions:
  setThing:
    kwargs:
      file: { type: string, required: true }
      name: { type: string, required: true }
      code: { type: string, required: true }
    injects_into: "{{file}}"
    anchor: decls
`)
	write("py/actions/setThing.py", "{{code}}\n")
	set, findings, err := genpkg.Load(gen, genpkg.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) > 0 {
		t.Fatalf("generator findings: %v", findings)
	}
	root := t.TempDir()
	records := &record.Set{Records: []*record.Record{{ID: "imp-1", Type: "imprint", Patterns: []string{"*.py"}}}}
	return root, set, records
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSelectCandidate(t *testing.T) {\n\troot, set, records := newSelectFixture(t)\n\tsrc := []byte(`x = 1\n`)\n\tif err := os.WriteFile(filepath.Join(root, \"a.py\"), src, 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tcands, skips, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{\"a.py\"})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(skips) != 0 {\n\t\tt.Fatalf(\"skips = %+v, want none\", skips)\n\t}\n\tif len(cands) != 1 {\n\t\tt.Fatalf(\"got %d candidates, want 1\", len(cands))\n\t}\n\tc := cands[0]\n\tif c.Template != \"{name}.py\" {\n\t\tt.Errorf(\"Template = %q, want %q\", c.Template, \"{name}.py\")\n\t}\n\tif c.Record == nil || c.Record.ID != \"imp-1\" {\n\t\tt.Errorf(\"Record = %+v, want ID imp-1\", c.Record)\n\t}\n\tif c.Captures[\"name\"] != \"a\" {\n\t\tt.Errorf(\"Captures[name] = %q, want %q\", c.Captures[\"name\"], \"a\")\n\t}\n\tif !bytes.Equal(c.Src, src) {\n\t\tt.Errorf(\"Src = %q, want %q\", c.Src, src)\n\t}\n}","file":"internal/adopt/select_test.go","name":"TestSelectCandidate"}}
func TestSelectCandidate(t *testing.T) {
	root, set, records := newSelectFixture(t)
	src := []byte(`x = 1
`)
	if err := os.WriteFile(filepath.Join(root, "a.py"), src, 0o644); err != nil {
		t.Fatal(err)
	}
	cands, skips, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{"a.py"})
	if err != nil {
		t.Fatal(err)
	}
	if len(skips) != 0 {
		t.Fatalf("skips = %+v, want none", skips)
	}
	if len(cands) != 1 {
		t.Fatalf("got %d candidates, want 1", len(cands))
	}
	c := cands[0]
	if c.Template != "{name}.py" {
		t.Errorf("Template = %q, want %q", c.Template, "{name}.py")
	}
	if c.Record == nil || c.Record.ID != "imp-1" {
		t.Errorf("Record = %+v, want ID imp-1", c.Record)
	}
	if c.Captures["name"] != "a" {
		t.Errorf("Captures[name] = %q, want %q", c.Captures["name"], "a")
	}
	if !bytes.Equal(c.Src, src) {
		t.Errorf("Src = %q, want %q", c.Src, src)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSelectSkipReasons(t *testing.T) {\n\tcases := []struct {\n\t\tname string\n\t\tfile string\n\t\tbody string\n\t\twant Reason\n\t}{\n\t\t{\"unmanaged\", \"setup.py\", \"x = 1\\n\", ReasonUnmanaged},\n\t\t{\"test\", \"test_a.py\", \"x = 1\\n\", ReasonTest},\n\t\t{\"generated\", \"gen.py\", \"# DO NOT EDIT\\nx = 1\\n\", ReasonGenerated},\n\t\t{\"already adopted\", \"done.py\", \"# sedum:setThing {}\\n\", ReasonAlreadyAdopted},\n\t\t{\"no package\", \"notes.txt\", \"hello\\n\", ReasonNoPackage},\n\t}\n\tfor _, tc := range cases {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\troot, set, records := newSelectFixture(t)\n\t\t\tif err := os.WriteFile(filepath.Join(root, tc.file), []byte(tc.body), 0o644); err != nil {\n\t\t\t\tt.Fatal(err)\n\t\t\t}\n\t\t\tcands, skips, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{tc.file})\n\t\t\tif err != nil {\n\t\t\t\tt.Fatal(err)\n\t\t\t}\n\t\t\tif len(cands) != 0 || len(skips) != 1 {\n\t\t\t\tt.Fatalf(\"got %d candidates and %d skips, want 0 and 1\", len(cands), len(skips))\n\t\t\t}\n\t\t\tif skips[0].Reason != tc.want {\n\t\t\t\tt.Errorf(\"Reason = %q, want %q\", skips[0].Reason, tc.want)\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/adopt/select_test.go","name":"TestSelectSkipReasons"}}
func TestSelectSkipReasons(t *testing.T) {
	cases := []struct {
		name string
		file string
		body string
		want Reason
	}{
		{"unmanaged", "setup.py", "x = 1\n", ReasonUnmanaged},
		{"test", "test_a.py", "x = 1\n", ReasonTest},
		{"generated", "gen.py", "# DO NOT EDIT\nx = 1\n", ReasonGenerated},
		{"already adopted", "done.py", "# sedum:setThing {}\n", ReasonAlreadyAdopted},
		{"no package", "notes.txt", "hello\n", ReasonNoPackage},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, set, records := newSelectFixture(t)
			if err := os.WriteFile(filepath.Join(root, tc.file), []byte(tc.body), 0o644); err != nil {
				t.Fatal(err)
			}
			cands, skips, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{tc.file})
			if err != nil {
				t.Fatal(err)
			}
			if len(cands) != 0 || len(skips) != 1 {
				t.Fatalf("got %d candidates and %d skips, want 0 and 1", len(cands), len(skips))
			}
			if skips[0].Reason != tc.want {
				t.Errorf("Reason = %q, want %q", skips[0].Reason, tc.want)
			}
		})
	}
}

// /sedum:setFunction
