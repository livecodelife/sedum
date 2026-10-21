package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// newRunFixture writes a generator tree, a records directory, and the given\n// project files under a temporary directory, and returns a Config for Run.\nfunc newRunFixture(t *testing.T, files map[string]string) (cfg Config) {\n\tt.Helper()\n\tbase := t.TempDir()\n\twrite := func(rel, body string) {\n\t\tp := filepath.Join(base, rel)\n\t\tif err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t\tif err := os.WriteFile(p, []byte(body), 0o644); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\twrite(\"gen/py/sedum.yaml\", `name: py\nextensions: [\".py\"]\ncomment_prefix: \"#\"\nadopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  unit_end: dedent\n  actions:\n    def: setThing\n`)\n\twrite(\"gen/py/files/{name}.py\", \"# sedum:anchor:decls\\n\")\n\twrite(\"gen/py/actions/actions.yaml\", `actions:\n  setThing:\n    kwargs:\n      file: { type: string, required: true }\n      name: { type: string, required: true }\n      code: { type: string, required: true }\n    injects_into: \"{{file}}\"\n    anchor: decls\n`)\n\twrite(\"gen/py/actions/setThing.py\", \"{{code}}\\n\")\n\twrite(\"prov/imp.yml\", `id: prov-2026-aaaaaaaa\ntype: imprint\nstatus: open\nintent: adopt py files\ncreated_at: \"2026-10-16\"\naffected_scope: [\"*.py\"]\n`)\n\tif err := os.MkdirAll(filepath.Join(base, \"proj\"), 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfor name, body := range files {\n\t\twrite(filepath.Join(\"proj\", name), body)\n\t}\n\treturn Config{\n\t\tRoot:       filepath.Join(base, \"proj\"),\n\t\tGenerators: filepath.Join(base, \"gen\"),\n\t\tRecords:    filepath.Join(base, \"prov\"),\n\t\tPaths:      []string{\".\"},\n\t}\n}","file":"internal/adopt/run_test.go","name":"newRunFixture"}}
// newRunFixture writes a generator tree, a records directory, and the given
// project files under a temporary directory, and returns a Config for Run.
func newRunFixture(t *testing.T, files map[string]string) (cfg Config) {
	t.Helper()
	base := t.TempDir()
	write := func(rel, body string) {
		p := filepath.Join(base, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("gen/py/sedum.yaml", `name: py
extensions: [".py"]
comment_prefix: "#"
adopt:
  unit_start: '^(?P<kind>def) (?P<name>\w+)'
  unit_end: dedent
  actions:
    def: setThing
`)
	write("gen/py/files/{name}.py", "# sedum:anchor:decls\n")
	write("gen/py/actions/actions.yaml", `actions:
  setThing:
    kwargs:
      file: { type: string, required: true }
      name: { type: string, required: true }
      code: { type: string, required: true }
    injects_into: "{{file}}"
    anchor: decls
`)
	write("gen/py/actions/setThing.py", "{{code}}\n")
	write("prov/imp.yml", `id: prov-2026-aaaaaaaa
type: imprint
status: open
intent: adopt py files
created_at: "2026-10-16"
affected_scope: ["*.py"]
`)
	if err := os.MkdirAll(filepath.Join(base, "proj"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		write(filepath.Join("proj", name), body)
	}
	return Config{
		Root:       filepath.Join(base, "proj"),
		Generators: filepath.Join(base, "gen"),
		Records:    filepath.Join(base, "prov"),
		Paths:      []string{"."},
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// TestRunAdopts checks that Run adopts a two-unit Python file, rewrites it\n// with the replayed unit, and writes the recording.\nfunc TestRunAdopts(t *testing.T) {\n\tconst orig = \"def f():\\n    return 1\\n\\ndef g():\\n    return 2\\n\"\n\tcfg := newRunFixture(t, map[string]string{\"a.py\": orig})\n\trep, err := Run(cfg)\n\tif err != nil {\n\t\tt.Fatalf(\"Run: %v\", err)\n\t}\n\tif rep.Totals.Adopted != 1 || rep.Totals.Units != 2 {\n\t\tt.Fatalf(\"totals = %+v, want 1 adopted and 2 units\", rep.Totals)\n\t}\n\tgot, err := os.ReadFile(filepath.Join(cfg.Root, \"a.py\"))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif !strings.Contains(string(got), \"# sedum:setThing\") {\n\t\tt.Fatalf(\"file not rewritten:\\n%s\", got)\n\t}\n\tif err := Verify(genpkg.Comment{Prefix: \"#\"}, orig, string(got)); err != nil {\n\t\tt.Fatalf(\"Verify: %v\", err)\n\t}\n\tif _, err := os.Stat(filepath.Join(cfg.Root, \".sedum\", \"adopt\", \"adopt.json\")); err != nil {\n\t\tt.Fatalf(\"recording missing: %v\", err)\n\t}\n}","file":"internal/adopt/run_test.go","name":"TestRunAdopts"}}
// TestRunAdopts checks that Run adopts a two-unit Python file, rewrites it
// with the replayed unit, and writes the recording.
func TestRunAdopts(t *testing.T) {
	const orig = "def f():\n    return 1\n\ndef g():\n    return 2\n"
	cfg := newRunFixture(t, map[string]string{"a.py": orig})
	rep, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if rep.Totals.Adopted != 1 || rep.Totals.Units != 2 {
		t.Fatalf("totals = %+v, want 1 adopted and 2 units", rep.Totals)
	}
	got, err := os.ReadFile(filepath.Join(cfg.Root, "a.py"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "# sedum:setThing") {
		t.Fatalf("file not rewritten:\n%s", got)
	}
	if err := Verify(genpkg.Comment{Prefix: "#"}, orig, string(got)); err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if _, err := os.Stat(filepath.Join(cfg.Root, ".sedum", "adopt", "adopt.json")); err != nil {
		t.Fatalf("recording missing: %v", err)
	}
}

// /sedum:setFunction
