package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/run_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// newRunFixture writes a generators tree, a records directory, and the given\n// project files under t.TempDir(), and returns the matching Config.\nfunc newRunFixture(t *testing.T, files map[string]string) (cfg Config) {\n\tt.Helper()\n\tbase := t.TempDir()\n\tgen := filepath.Join(base, \"gen\")\n\tprov := filepath.Join(base, \"prov\")\n\tproj := filepath.Join(base, \"proj\")\n\twrite := func(p, s string) {\n\t\tif err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t\tif err := os.WriteFile(p, []byte(s), 0o644); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\twrite(filepath.Join(gen, \"py\", \"sedum.yaml\"), `name: py\nextensions: [\".py\"]\ncomment_prefix: \"#\"\nadopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  unit_end: dedent\n  actions:\n    def: setThing\n`)\n\twrite(filepath.Join(gen, \"py\", \"files\", \"{name}.py\"), \"# sedum:anchor:decls\\n\")\n\twrite(filepath.Join(gen, \"py\", \"actions\", \"actions.yaml\"), `actions:\n  setThing:\n    kwargs:\n      file: { type: string, required: true }\n      name: { type: string, required: true }\n      code: { type: string, required: true }\n    injects_into: \"{{file}}\"\n    anchor: decls\n`)\n\twrite(filepath.Join(gen, \"py\", \"actions\", \"setThing.py\"), \"{{code}}\\n\")\n\twrite(filepath.Join(prov, \"imp.yml\"), `id: prov-2026-aaaaaaaa\ntype: imprint\nstatus: open\nintent: adopt py files\ncreated_at: \"2026-10-16\"\naffected_scope: [\"*.py\"]\n`)\n\tfor name, body := range files {\n\t\twrite(filepath.Join(proj, name), body)\n\t}\n\treturn Config{Root: proj, Generators: gen, Records: prov, Paths: []string{\".\"}}\n}","file":"internal/adopt/run_test.go","name":"newRunFixture"}}
// newRunFixture writes a generators tree, a records directory, and the given
// project files under t.TempDir(), and returns the matching Config.
func newRunFixture(t *testing.T, files map[string]string) (cfg Config) {
	t.Helper()
	base := t.TempDir()
	gen := filepath.Join(base, "gen")
	prov := filepath.Join(base, "prov")
	proj := filepath.Join(base, "proj")
	write := func(p, s string) {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(gen, "py", "sedum.yaml"), `name: py
extensions: [".py"]
comment_prefix: "#"
adopt:
  unit_start: '^(?P<kind>def) (?P<name>\w+)'
  unit_end: dedent
  actions:
    def: setThing
`)
	write(filepath.Join(gen, "py", "files", "{name}.py"), "# sedum:anchor:decls\n")
	write(filepath.Join(gen, "py", "actions", "actions.yaml"), `actions:
  setThing:
    kwargs:
      file: { type: string, required: true }
      name: { type: string, required: true }
      code: { type: string, required: true }
    injects_into: "{{file}}"
    anchor: decls
`)
	write(filepath.Join(gen, "py", "actions", "setThing.py"), "{{code}}\n")
	write(filepath.Join(prov, "imp.yml"), `id: prov-2026-aaaaaaaa
type: imprint
status: open
intent: adopt py files
created_at: "2026-10-16"
affected_scope: ["*.py"]
`)
	for name, body := range files {
		write(filepath.Join(proj, name), body)
	}
	return Config{Root: proj, Generators: gen, Records: prov, Paths: []string{"."}}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:47d87e16dfdbc12882d2504c04389cd1316fe520c4b77f08f20e484f53308841","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunAdopts(t *testing.T) {\n\torig := \"def f():\\n    return 1\\n\\ndef g():\\n    return 2\\n\"\n\tcfg := newRunFixture(t, map[string]string{\"a.py\": orig})\n\treport, err := Run(cfg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif report.Totals.Adopted != 1 {\n\t\tt.Errorf(\"adopted = %d, want 1\", report.Totals.Adopted)\n\t}\n\tif report.Totals.Units != 2 {\n\t\tt.Errorf(\"units = %d, want 2\", report.Totals.Units)\n\t}\n\tgot, err := os.ReadFile(filepath.Join(cfg.Root, \"a.py\"))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif !strings.Contains(string(got), \"# sedum:setThing\") {\n\t\tt.Errorf(\"file not rewritten: %q\", string(got))\n\t}\n\tif err := Verify(genpkg.Comment{Prefix: \"#\"}, orig, string(got)); err != nil {\n\t\tt.Errorf(\"verify: %v\", err)\n\t}\n\tif _, err := os.Stat(filepath.Join(cfg.Root, \".sedum\", \"adopt\", \"adopt.json\")); err != nil {\n\t\tt.Errorf(\"recording missing: %v\", err)\n\t}\n}","file":"internal/adopt/run_test.go","name":"TestRunAdopts"}}
func TestRunAdopts(t *testing.T) {
	orig := "def f():\n    return 1\n\ndef g():\n    return 2\n"
	cfg := newRunFixture(t, map[string]string{"a.py": orig})
	report, err := Run(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if report.Totals.Adopted != 1 {
		t.Errorf("adopted = %d, want 1", report.Totals.Adopted)
	}
	if report.Totals.Units != 2 {
		t.Errorf("units = %d, want 2", report.Totals.Units)
	}
	got, err := os.ReadFile(filepath.Join(cfg.Root, "a.py"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "# sedum:setThing") {
		t.Errorf("file not rewritten: %q", string(got))
	}
	if err := Verify(genpkg.Comment{Prefix: "#"}, orig, string(got)); err != nil {
		t.Errorf("verify: %v", err)
	}
	if _, err := os.Stat(filepath.Join(cfg.Root, ".sedum", "adopt", "adopt.json")); err != nil {
		t.Errorf("recording missing: %v", err)
	}
}

// /sedum:setFunction
