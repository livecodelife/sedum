package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/adopt_test.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/adopt_test.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/adopt_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/adopt_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/adopt_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/adopt_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func newAdoptCLIFixture(t *testing.T) (gen, prov, proj string) {\n\tt.Helper()\n\troot := t.TempDir()\n\tgen = filepath.Join(root, \"gen\")\n\tprov = filepath.Join(root, \"prov\")\n\tproj = filepath.Join(root, \"proj\")\n\tfiles := map[string]string{\n\t\tfilepath.Join(gen, \"py\", \"sedum.yaml\"):              \"name: py\\nextensions: [\\\".py\\\"]\\ncomment_prefix: \\\"#\\\"\\nadopt:\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  unit_end: dedent\\n  actions:\\n    def: setThing\\n\",\n\t\tfilepath.Join(gen, \"py\", \"files\", \"{name}.py\"):      \"# sedum:anchor:decls\\n\",\n\t\tfilepath.Join(gen, \"py\", \"actions\", \"actions.yaml\"): \"actions:\\n  setThing:\\n    kwargs:\\n      file: { type: string, required: true }\\n      name: { type: string, required: true }\\n      code: { type: string, required: true }\\n    injects_into: \\\"{{file}}\\\"\\n    anchor: decls\\n\",\n\t\tfilepath.Join(gen, \"py\", \"actions\", \"setThing.py\"):  \"{{code}}\\n\",\n\t\tfilepath.Join(prov, \"imp.yml\"):                      \"id: prov-2026-aaaaaaaa\\ntype: imprint\\nstatus: open\\nintent: adopt py files\\ncreated_at: \\\"2026-10-16\\\"\\naffected_scope: [\\\"*.py\\\"]\\n\",\n\t\tfilepath.Join(proj, \"a.py\"):                         \"def f():\\n    return 1\\n\",\n\t\tfilepath.Join(proj, \"b.py\"):                         \"import os\\n\\ndef g():\\n    return 1\\n\",\n\t}\n\tfor path, content := range files {\n\t\tif err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t\tif err := os.WriteFile(path, []byte(content), 0o644); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\treturn gen, prov, proj\n}","file":"internal/cli/adopt_test.go","name":"newAdoptCLIFixture"}}
func newAdoptCLIFixture(t *testing.T) (gen, prov, proj string) {
	t.Helper()
	root := t.TempDir()
	gen = filepath.Join(root, "gen")
	prov = filepath.Join(root, "prov")
	proj = filepath.Join(root, "proj")
	files := map[string]string{
		filepath.Join(gen, "py", "sedum.yaml"):              "name: py\nextensions: [\".py\"]\ncomment_prefix: \"#\"\nadopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  unit_end: dedent\n  actions:\n    def: setThing\n",
		filepath.Join(gen, "py", "files", "{name}.py"):      "# sedum:anchor:decls\n",
		filepath.Join(gen, "py", "actions", "actions.yaml"): "actions:\n  setThing:\n    kwargs:\n      file: { type: string, required: true }\n      name: { type: string, required: true }\n      code: { type: string, required: true }\n    injects_into: \"{{file}}\"\n    anchor: decls\n",
		filepath.Join(gen, "py", "actions", "setThing.py"):  "{{code}}\n",
		filepath.Join(prov, "imp.yml"):                      "id: prov-2026-aaaaaaaa\ntype: imprint\nstatus: open\nintent: adopt py files\ncreated_at: \"2026-10-16\"\naffected_scope: [\"*.py\"]\n",
		filepath.Join(proj, "a.py"):                         "def f():\n    return 1\n",
		filepath.Join(proj, "b.py"):                         "import os\n\ndef g():\n    return 1\n",
	}
	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return gen, prov, proj
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAdoptCommandText(t *testing.T) {\n\tgen, prov, proj := newAdoptCLIFixture(t)\n\troot := NewRootCommand()\n\tvar buf bytes.Buffer\n\troot.SetOut(&buf)\n\troot.SetErr(&buf)\n\troot.SetArgs([]string{\"adopt\", \"--generators\", gen, \"--records\", prov, \"--output\", proj, \"--dry-run\", \".\"})\n\tif err := root.Execute(); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tout := buf.String()\n\tfor _, want := range []string{\"adopt: 2 file(s): 1 adopted, 1 skipped, 0 rejected\", \"other_residual\"} {\n\t\tif !strings.Contains(out, want) {\n\t\t\tt.Errorf(\"output missing %q:\\n%s\", want, out)\n\t\t}\n\t}\n\tgot, err := os.ReadFile(filepath.Join(proj, \"a.py\"))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif string(got) != \"def f():\\n    return 1\\n\" {\n\t\tt.Errorf(\"a.py changed on dry run: %q\", got)\n\t}\n}","file":"internal/cli/adopt_test.go","name":"TestAdoptCommandText"}}
func TestAdoptCommandText(t *testing.T) {
	gen, prov, proj := newAdoptCLIFixture(t)
	root := NewRootCommand()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"adopt", "--generators", gen, "--records", prov, "--output", proj, "--dry-run", "."})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"adopt: 2 file(s): 1 adopted, 1 skipped, 0 rejected", "other_residual"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
	got, err := os.ReadFile(filepath.Join(proj, "a.py"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "def f():\n    return 1\n" {
		t.Errorf("a.py changed on dry run: %q", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAdoptCommandJSON(t *testing.T) {\n\tgen, prov, proj := newAdoptCLIFixture(t)\n\troot := NewRootCommand()\n\tvar buf bytes.Buffer\n\troot.SetOut(&buf)\n\troot.SetErr(&buf)\n\troot.SetArgs([]string{\"adopt\", \"--generators\", gen, \"--records\", prov, \"--output\", proj, \"--dry-run\", \"--json\", \".\"})\n\tif err := root.Execute(); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tvar m map[string]any\n\tif err := json.Unmarshal(buf.Bytes(), &m); err != nil {\n\t\tt.Fatalf(\"output is not JSON: %v\\n%s\", err, buf.String())\n\t}\n\ttotals, ok := m[\"totals\"].(map[string]any)\n\tif !ok {\n\t\tt.Fatalf(\"missing totals in %v\", m)\n\t}\n\tif totals[\"adopted\"] != float64(1) {\n\t\tt.Errorf(\"adopted = %v, want 1\", totals[\"adopted\"])\n\t}\n}","file":"internal/cli/adopt_test.go","name":"TestAdoptCommandJSON"}}
func TestAdoptCommandJSON(t *testing.T) {
	gen, prov, proj := newAdoptCLIFixture(t)
	root := NewRootCommand()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"adopt", "--generators", gen, "--records", prov, "--output", proj, "--dry-run", "--json", "."})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, buf.String())
	}
	totals, ok := m["totals"].(map[string]any)
	if !ok {
		t.Fatalf("missing totals in %v", m)
	}
	if totals["adopted"] != float64(1) {
		t.Errorf("adopted = %v, want 1", totals["adopted"])
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAdoptCommandWritesTheFile(t *testing.T) {\n\tgen, prov, proj := newAdoptCLIFixture(t)\n\troot := NewRootCommand()\n\tvar buf bytes.Buffer\n\troot.SetOut(&buf)\n\troot.SetErr(&buf)\n\troot.SetArgs([]string{\"adopt\", \"--generators\", gen, \"--records\", prov, \"--output\", proj, \".\"})\n\tif err := root.Execute(); err != nil {\n\t\tt.Fatalf(\"%v\\n%s\", err, buf.String())\n\t}\n\ta, err := os.ReadFile(filepath.Join(proj, \"a.py\"))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif !strings.Contains(string(a), \"# sedum:setThing\") {\n\t\tt.Errorf(\"a.py not adopted:\\n%s\", a)\n\t}\n\tb, err := os.ReadFile(filepath.Join(proj, \"b.py\"))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif string(b) != \"import os\\n\\ndef g():\\n    return 1\\n\" {\n\t\tt.Errorf(\"b.py changed:\\n%s\", b)\n\t}\n}","file":"internal/cli/adopt_test.go","name":"TestAdoptCommandWritesTheFile"}}
func TestAdoptCommandWritesTheFile(t *testing.T) {
	gen, prov, proj := newAdoptCLIFixture(t)
	root := NewRootCommand()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"adopt", "--generators", gen, "--records", prov, "--output", proj, "."})
	if err := root.Execute(); err != nil {
		t.Fatalf("%v\n%s", err, buf.String())
	}
	a, err := os.ReadFile(filepath.Join(proj, "a.py"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(a), "# sedum:setThing") {
		t.Errorf("a.py not adopted:\n%s", a)
	}
	b, err := os.ReadFile(filepath.Join(proj, "b.py"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "import os\n\ndef g():\n    return 1\n" {
		t.Errorf("b.py changed:\n%s", b)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAdoptCommandRequiresGeneratorsAndRecords(t *testing.T) {\n\troot := NewRootCommand()\n\tvar buf bytes.Buffer\n\troot.SetOut(&buf)\n\troot.SetErr(&buf)\n\troot.SetArgs([]string{\"adopt\", \".\"})\n\terr := root.Execute()\n\tif err == nil {\n\t\tt.Fatal(\"expected an error without --generators and --records\")\n\t}\n\tif !strings.Contains(err.Error(), \"generators\") && !strings.Contains(err.Error(), \"records\") {\n\t\tt.Errorf(\"error does not mention the missing flag: %v\", err)\n\t}\n}","file":"internal/cli/adopt_test.go","name":"TestAdoptCommandRequiresGeneratorsAndRecords"}}
func TestAdoptCommandRequiresGeneratorsAndRecords(t *testing.T) {
	root := NewRootCommand()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"adopt", "."})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected an error without --generators and --records")
	}
	if !strings.Contains(err.Error(), "generators") && !strings.Contains(err.Error(), "records") {
		t.Errorf("error does not mention the missing flag: %v", err)
	}
}

// /sedum:setFunction
