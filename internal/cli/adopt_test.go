package cli

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/adopt_test.go","path":"bytes"},"kwargs":"sha256:c484bffaee1ccf73b74798ea4035492b4b229ec30fffcf47ca095f9fe673e9a6","body":"sha256:468731f6dd513df0a66690d4649dc9e7244c96de227d8e6ead98cc5e9747c57a"}
import "bytes"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/adopt_test.go","path":"encoding/json"},"kwargs":"sha256:aeaf04b643e3517d9e9e8f002e35d6fd98ba672a0f1603f39e19f1690f207bf2","body":"sha256:872d7c3e9dee608e79a7f8ef57930d24b5e30f2a57d891a52549179ace45946a"}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/adopt_test.go","path":"os"},"kwargs":"sha256:da7523694e9649629f572984fce3e1a934c2e41a88c48bce7b82df729d444a96","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/adopt_test.go","path":"path/filepath"},"kwargs":"sha256:dec2c5bb4a2f92f57396e0cf77bb5f2b2410b646500cd387948687e0851c019d","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/adopt_test.go","path":"strings"},"kwargs":"sha256:da151e52dff07afbfb06385c9b216e4b6e659b795059acb30a23e322341cb970","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/adopt_test.go","path":"testing"},"kwargs":"sha256:69cc8d42d96abd9b25fbd3a3488ebc6b09f47b622bcfb4361096b12198234674","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/adopt_test.go","name":"newAdoptCLIFixture"},"kwargs":"sha256:de47892779f8030bc2a6b2434ff48ae23c03b6926070b215d75bf838010b9ed6","body":"sha256:c82997c200ee01ef6f783e9a319d000b9d3fe15b3f6f560fc324aa08a4f97309"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/adopt_test.go","name":"TestAdoptCommandText"},"kwargs":"sha256:83039ced10351fe28a7d57d4d1ac2293c18c57dfdab5b6245c511ca0e9bbf208","body":"sha256:ba7b4b9c13c7c3e1fb3581fb7e6b4f49b3174a51f14494dc524eda7ada79a845"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/adopt_test.go","name":"TestAdoptCommandJSON"},"kwargs":"sha256:fb89536ab4f2c71b9f97253c21e5303cd397d00befbb9efe9858fb58b4c4e57e","body":"sha256:1ec7d36f86984e6938f77f9895e587f3f244c0fb528d4324e3b4963eb2363d1e"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/adopt_test.go","name":"TestAdoptCommandWritesTheFile"},"kwargs":"sha256:1acdf8e41a49cdeabe6b7703ebabdb64c9a9fbb4f0ac88648c6004deaf375cd1","body":"sha256:af71e3dde8d98492247d57d4262f7e894dcdc13e55c9bd7ec24274b1dc8ddba9"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/adopt_test.go","name":"TestAdoptCommandRequiresGeneratorsAndRecords"},"kwargs":"sha256:7a51c50c2eda2287e732fce984c2199dd397aeac45ede0bb52e125cd13809101","body":"sha256:61109f37ae310ebbb51a5b8a920a25c79f45a5f6dbc9d6cd222a8e492699d022"}
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
