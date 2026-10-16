package genpkg

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:1562000e469bcf67c6bb531f82219649cc8a8a6cea0860c74eab819190aed59a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/adopt_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:1562000e469bcf67c6bb531f82219649cc8a8a6cea0860c74eab819190aed59a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/adopt_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:1562000e469bcf67c6bb531f82219649cc8a8a6cea0860c74eab819190aed59a","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// adoptSetThingYAML is appended to rails/actions/actions.yaml: a declarative\n// action whose required kwargs are exactly file, name and code.\nconst adoptSetThingYAML = `  setThing:\n    kwargs:\n      file: { type: string, required: true }\n      name: { type: string, required: true }\n      code: { type: string, required: true }\n    injects_into: \"{{file}}\"\n    anchor: class_body\n`","file":"internal/genpkg/adopt_test.go","name":"adoptSetThingYAML"}}
// adoptSetThingYAML is appended to rails/actions/actions.yaml: a declarative
// action whose required kwargs are exactly file, name and code.
const adoptSetThingYAML = `  setThing:
    kwargs:
      file: { type: string, required: true }
      name: { type: string, required: true }
      code: { type: string, required: true }
    injects_into: "{{file}}"
    anchor: class_body
`

// /sedum:setConst
// sedum:setFunction {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:1562000e469bcf67c6bb531f82219649cc8a8a6cea0860c74eab819190aed59a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAdopt(t *testing.T) {\n\tvalid := \"adopt:\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  unit_end: dedent\\n  actions:\\n    def: setThing\\n\"\n\tcases := []struct {\n\t\tname, block, msg string\n\t}{\n\t\t{\"declarative\", valid, \"\"},\n\t\t{\"command\", \"adopt:\\n  command: [\\\"go\\\"]\\n\", \"\"},\n\t\t{\"none\", \"\", \"\"},\n\t\t{\"unit_end missing\", strings.Replace(valid, \"  unit_end: dedent\\n\", \"\", 1), \"unit_end\"},\n\t\t{\"unit_end bad\", strings.Replace(valid, \"unit_end: dedent\", \"unit_end: '('\", 1), \"unit_end\"},\n\t\t{\"empty command element\", \"adopt:\\n  command: [\\\"\\\"]\\n\", \"empty element\"},\n\t\t{\"command mixed\", \"adopt:\\n  command: [\\\"go\\\"]\\n  header: '^//'\\n\", \"cannot be combined\"},\n\t\t{\"bad unit_start\", \"adopt:\\n  unit_start: '(unclosed'\\n  unit_end: dedent\\n  actions:\\n    def: setThing\\n\", \"does not compile\"},\n\t\t{\"unit_start lacks name\", \"adopt:\\n  unit_start: '^(?P<kind>def)'\\n  unit_end: dedent\\n  actions:\\n    def: setThing\\n\", \"named group\"},\n\t\t{\"unknown action\", strings.Replace(valid, \"setThing\", \"nosuch\", 1), \"does not define\"},\n\t\t{\"action kwarg not supplied\", strings.Replace(valid, \"setThing\", \"createControllerMethod\", 1), \"requires kwarg\"},\n\t\t{\"empty actions\", \"adopt:\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n  unit_end: dedent\\n\", \"actions must\"},\n\t\t{\"header bad\", valid + \"  header: '('\\n\", \"header\"},\n\t\t{\"attach_above bad\", valid + \"  attach_above: ['(']\\n\", \"attach_above\"},\n\t}\n\tfor _, tc := range cases {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\tfiles := validPackage()\n\t\t\tfiles[\"rails/actions/actions.yaml\"] += adoptSetThingYAML\n\t\t\tfiles[\"rails/actions/setThing.rb\"] = \"{{code}}\\n\"\n\t\t\tfiles[\"rails/sedum.yaml\"] += tc.block\n\t\t\tset, findings := loadTree(t, files)\n\t\t\tif tc.block == \"\" && set.Packages[0].Adopt != nil {\n\t\t\t\tt.Fatalf(\"Adopt = %+v, want nil\", set.Packages[0].Adopt)\n\t\t\t}\n\t\t\tvar got []string\n\t\t\tfor _, f := range findings {\n\t\t\t\tif f.Rule == RuleAdoptInvalid {\n\t\t\t\t\tgot = append(got, f.Message)\n\t\t\t\t}\n\t\t\t}\n\t\t\tif tc.msg == \"\" {\n\t\t\t\tif len(got) != 0 {\n\t\t\t\t\tt.Fatalf(\"unexpected adopt findings: %v\", got)\n\t\t\t\t}\n\t\t\t\treturn\n\t\t\t}\n\t\t\tfor _, m := range got {\n\t\t\t\tif strings.Contains(m, tc.msg) {\n\t\t\t\t\treturn\n\t\t\t\t}\n\t\t\t}\n\t\t\tt.Fatalf(\"no %s finding containing %q; got %v\", RuleAdoptInvalid, tc.msg, got)\n\t\t})\n\t}\n}","file":"internal/genpkg/adopt_test.go","name":"TestAdopt"}}
func TestAdopt(t *testing.T) {
	valid := "adopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  unit_end: dedent\n  actions:\n    def: setThing\n"
	cases := []struct {
		name, block, msg string
	}{
		{"declarative", valid, ""},
		{"command", "adopt:\n  command: [\"go\"]\n", ""},
		{"none", "", ""},
		{"unit_end missing", strings.Replace(valid, "  unit_end: dedent\n", "", 1), "unit_end"},
		{"unit_end bad", strings.Replace(valid, "unit_end: dedent", "unit_end: '('", 1), "unit_end"},
		{"empty command element", "adopt:\n  command: [\"\"]\n", "empty element"},
		{"command mixed", "adopt:\n  command: [\"go\"]\n  header: '^//'\n", "cannot be combined"},
		{"bad unit_start", "adopt:\n  unit_start: '(unclosed'\n  unit_end: dedent\n  actions:\n    def: setThing\n", "does not compile"},
		{"unit_start lacks name", "adopt:\n  unit_start: '^(?P<kind>def)'\n  unit_end: dedent\n  actions:\n    def: setThing\n", "named group"},
		{"unknown action", strings.Replace(valid, "setThing", "nosuch", 1), "does not define"},
		{"action kwarg not supplied", strings.Replace(valid, "setThing", "createControllerMethod", 1), "requires kwarg"},
		{"empty actions", "adopt:\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n  unit_end: dedent\n", "actions must"},
		{"header bad", valid + "  header: '('\n", "header"},
		{"attach_above bad", valid + "  attach_above: ['(']\n", "attach_above"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := validPackage()
			files["rails/actions/actions.yaml"] += adoptSetThingYAML
			files["rails/actions/setThing.rb"] = "{{code}}\n"
			files["rails/sedum.yaml"] += tc.block
			set, findings := loadTree(t, files)
			if tc.block == "" && set.Packages[0].Adopt != nil {
				t.Fatalf("Adopt = %+v, want nil", set.Packages[0].Adopt)
			}
			var got []string
			for _, f := range findings {
				if f.Rule == RuleAdoptInvalid {
					got = append(got, f.Message)
				}
			}
			if tc.msg == "" {
				if len(got) != 0 {
					t.Fatalf("unexpected adopt findings: %v", got)
				}
				return
			}
			for _, m := range got {
				if strings.Contains(m, tc.msg) {
					return
				}
			}
			t.Fatalf("no %s finding containing %q; got %v", RuleAdoptInvalid, tc.msg, got)
		})
	}
}

// /sedum:setFunction
