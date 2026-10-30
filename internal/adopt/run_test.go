package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/run_test.go","path":"os"},"kwargs":"sha256:ea6c0fc6c1de82545cc84254ad3bac24bf98728496cf80f2bcb1fa876ed8e86a","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/run_test.go","path":"path/filepath"},"kwargs":"sha256:337dc19833fae26f2d75f010ce2c403c67184bb38bb14e3aca558600cdee08fe","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/run_test.go","path":"strings"},"kwargs":"sha256:ff8be3df73e0993fd9ab63faa0b4301f04faf2c2c6f3f816bd4a8cc80fc4f530","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/run_test.go","path":"testing"},"kwargs":"sha256:14317bd18548753433b07d236422dedefc2b33512ae0beda7ab19a64830cef63","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/run_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:a3f58f0cca1d44b40489811ebf94bc54dd562c7e462c797916b7332ff38d6954","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/run_test.go","name":"newRunFixture"},"kwargs":"sha256:35278cf89b9e739f7ff6beda321663db62c20566aaabd7d2e387f1a8d41f8358","body":"sha256:83376a6bcb023d5f6d73cd73ea8ecc745c2ee0da6c4c4fdf40ed9f646510a439"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-520ee992","gen":"sha256:abce6ad3af5346e3865652129852e6d56fe2af3bcf443cff1e4e15ca97b66586","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/run_test.go","name":"TestRunAdopts"},"kwargs":"sha256:4ac2bfb5958f2f83331e328bad69d775eef2db782047968d302a766f004c497c","body":"sha256:39e319eea5529f284ccf357fa39c9f5ff686d1d7d8958081a51366ffdad563ac"}
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
