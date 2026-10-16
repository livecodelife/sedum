package genpkg

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:235dd3daaefdd4a98dd3fb2bab1e2988d47e6bc2aab63e74e24144bec61de0c5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/adopt_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:235dd3daaefdd4a98dd3fb2bab1e2988d47e6bc2aab63e74e24144bec61de0c5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/adopt_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:235dd3daaefdd4a98dd3fb2bab1e2988d47e6bc2aab63e74e24144bec61de0c5","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"const setThingActions = `  setThing:\n    kwargs:\n      file: { type: string, required: true }\n      name: { type: string, required: true }\n      code: { type: string, required: true }\n    injects_into: \"{{file}}\"\n    anchor: class_body\n`","file":"internal/genpkg/adopt_test.go","name":"setThingActions"}}
const setThingActions = `  setThing:
    kwargs:
      file: { type: string, required: true }
      name: { type: string, required: true }
      code: { type: string, required: true }
    injects_into: "{{file}}"
    anchor: class_body
`

// /sedum:setConst
// sedum:setFunction {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:235dd3daaefdd4a98dd3fb2bab1e2988d47e6bc2aab63e74e24144bec61de0c5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// adoptTree loads the rails fixture with a setThing action and the given adopt: block appended to sedum.yaml.\nfunc adoptTree(t *testing.T, block string) (*Set, Findings) {\n\tt.Helper()\n\tfiles := validPackage()\n\tfiles[\"rails/actions/actions.yaml\"] += setThingActions\n\tfiles[\"rails/actions/setThing.rb\"] = \"{{code}}\\n\"\n\tfiles[\"rails/sedum.yaml\"] += block\n\treturn loadTree(t, files)\n}","file":"internal/genpkg/adopt_test.go","name":"adoptTree"}}
// adoptTree loads the rails fixture with a setThing action and the given adopt: block appended to sedum.yaml.
func adoptTree(t *testing.T, block string) (*Set, Findings) {
	t.Helper()
	files := validPackage()
	files["rails/actions/actions.yaml"] += setThingActions
	files["rails/actions/setThing.rb"] = "{{code}}\n"
	files["rails/sedum.yaml"] += block
	return loadTree(t, files)
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:235dd3daaefdd4a98dd3fb2bab1e2988d47e6bc2aab63e74e24144bec61de0c5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAdoptRules(t *testing.T) {\n\tconst declarative = \"adopt:\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  unit_end: dedent\\n  actions:\\n    def: setThing\\n\"\n\tcases := []struct {\n\t\tname, block, want string\n\t}{\n\t\t{\"declarative\", declarative, \"\"},\n\t\t{\"command\", \"adopt:\\n  command: [\\\"go\\\", \\\"vet\\\"]\\n\", \"\"},\n\t\t{\"both set\", \"adopt:\\n  command: [\\\"go\\\"]\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  actions:\\n    def: setThing\\n\", \"exactly one\"},\n\t\t{\"neither set\", \"adopt:\\n  header: '^x'\\n  actions:\\n    def: setThing\\n\", \"exactly one\"},\n\t\t{\"empty command\", \"adopt:\\n  command: [\\\"\\\"]\\n\", \"first element is empty\"},\n\t\t{\"unit_start does not compile\", \"adopt:\\n  unit_start: '^(?P<kind>def'\\n  actions:\\n    def: setThing\\n\", \"unit_start does not compile\"},\n\t\t{\"unit_start lacks name\", \"adopt:\\n  unit_start: '^(?P<kind>def)'\\n  actions:\\n    def: setThing\\n\", \"lacks the named group\"},\n\t\t{\"unit_end does not compile\", \"adopt:\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  unit_end: '('\\n  actions:\\n    def: setThing\\n\", \"unit_end does not compile\"},\n\t\t{\"header does not compile\", \"adopt:\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  header: '('\\n  actions:\\n    def: setThing\\n\", \"header does not compile\"},\n\t\t{\"attach_above does not compile\", \"adopt:\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  attach_above: ['(']\\n  actions:\\n    def: setThing\\n\", \"does not compile\"},\n\t\t{\"empty actions\", \"adopt:\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  actions: {}\\n\", \"must not be empty\"},\n\t\t{\"unknown action\", \"adopt:\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  actions:\\n    def: noSuchAction\\n\", \"is not an action\"},\n\t\t{\"extra required kwarg\", \"adopt:\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  actions:\\n    def: createControllerMethod\\n\", \"outside file, name, code\"},\n\t\t{\"command with declarative field\", \"adopt:\\n  command: [\\\"go\\\"]\\n  header: '^x'\\n\", \"cannot be combined\"},\n\t}\n\tfor _, c := range cases {\n\t\tt.Run(c.name, func(t *testing.T) {\n\t\t\tset, findings := adoptTree(t, c.block)\n\t\t\tvar got []string\n\t\t\tfor _, f := range findings {\n\t\t\t\tif f.Rule == RuleAdoptInvalid {\n\t\t\t\t\tgot = append(got, f.Message)\n\t\t\t\t}\n\t\t\t}\n\t\t\tif c.want == \"\" {\n\t\t\t\tif len(got) != 0 {\n\t\t\t\t\tt.Fatalf(\"unexpected findings: %v\", got)\n\t\t\t\t}\n\t\t\t\tif set.Packages[0].Adopt == nil {\n\t\t\t\t\tt.Fatal(\"Adopt not loaded\")\n\t\t\t\t}\n\t\t\t\treturn\n\t\t\t}\n\t\t\tfor _, m := range got {\n\t\t\t\tif strings.Contains(m, c.want) {\n\t\t\t\t\treturn\n\t\t\t\t}\n\t\t\t}\n\t\t\tt.Fatalf(\"no %q finding under %s; got %v\", c.want, RuleAdoptInvalid, got)\n\t\t})\n\t}\n}","file":"internal/genpkg/adopt_test.go","name":"TestAdoptRules"}}
func TestAdoptRules(t *testing.T) {
	const declarative = "adopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  unit_end: dedent\n  actions:\n    def: setThing\n"
	cases := []struct {
		name, block, want string
	}{
		{"declarative", declarative, ""},
		{"command", "adopt:\n  command: [\"go\", \"vet\"]\n", ""},
		{"both set", "adopt:\n  command: [\"go\"]\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  actions:\n    def: setThing\n", "exactly one"},
		{"neither set", "adopt:\n  header: '^x'\n  actions:\n    def: setThing\n", "exactly one"},
		{"empty command", "adopt:\n  command: [\"\"]\n", "first element is empty"},
		{"unit_start does not compile", "adopt:\n  unit_start: '^(?P<kind>def'\n  actions:\n    def: setThing\n", "unit_start does not compile"},
		{"unit_start lacks name", "adopt:\n  unit_start: '^(?P<kind>def)'\n  actions:\n    def: setThing\n", "lacks the named group"},
		{"unit_end does not compile", "adopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  unit_end: '('\n  actions:\n    def: setThing\n", "unit_end does not compile"},
		{"header does not compile", "adopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  header: '('\n  actions:\n    def: setThing\n", "header does not compile"},
		{"attach_above does not compile", "adopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  attach_above: ['(']\n  actions:\n    def: setThing\n", "does not compile"},
		{"empty actions", "adopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  actions: {}\n", "must not be empty"},
		{"unknown action", "adopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  actions:\n    def: noSuchAction\n", "is not an action"},
		{"extra required kwarg", "adopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  actions:\n    def: createControllerMethod\n", "outside file, name, code"},
		{"command with declarative field", "adopt:\n  command: [\"go\"]\n  header: '^x'\n", "cannot be combined"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			set, findings := adoptTree(t, c.block)
			var got []string
			for _, f := range findings {
				if f.Rule == RuleAdoptInvalid {
					got = append(got, f.Message)
				}
			}
			if c.want == "" {
				if len(got) != 0 {
					t.Fatalf("unexpected findings: %v", got)
				}
				if set.Packages[0].Adopt == nil {
					t.Fatal("Adopt not loaded")
				}
				return
			}
			for _, m := range got {
				if strings.Contains(m, c.want) {
					return
				}
			}
			t.Fatalf("no %q finding under %s; got %v", c.want, RuleAdoptInvalid, got)
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:235dd3daaefdd4a98dd3fb2bab1e2988d47e6bc2aab63e74e24144bec61de0c5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// TestAdoptAbsent checks that a package without an adopt: key has a nil Adopt and no adopt findings.\nfunc TestAdoptAbsent(t *testing.T) {\n\tset, findings := loadTree(t, validPackage())\n\tfor _, f := range findings {\n\t\tif f.Rule == RuleAdoptInvalid {\n\t\t\tt.Fatalf(\"unexpected adopt finding: %s\", f.Message)\n\t\t}\n\t}\n\tif set.Packages[0].Adopt != nil {\n\t\tt.Fatal(\"Adopt set without an adopt: key\")\n\t}\n}","file":"internal/genpkg/adopt_test.go","name":"TestAdoptAbsent"}}
// TestAdoptAbsent checks that a package without an adopt: key has a nil Adopt and no adopt findings.
func TestAdoptAbsent(t *testing.T) {
	set, findings := loadTree(t, validPackage())
	for _, f := range findings {
		if f.Rule == RuleAdoptInvalid {
			t.Fatalf("unexpected adopt finding: %s", f.Message)
		}
	}
	if set.Packages[0].Adopt != nil {
		t.Fatal("Adopt set without an adopt: key")
	}
}

// /sedum:setFunction
