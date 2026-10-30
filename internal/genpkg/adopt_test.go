package genpkg

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/adopt_test.go","path":"strings"},"kwargs":"sha256:da6fbe7c9e1225de2337876caaaa69c305da7eb2ace4417caddd947638347563","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/adopt_test.go","path":"testing"},"kwargs":"sha256:fff1330a19420c4344c1e08f296972f1c2af851729605884d515152ba8722285","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/genpkg/adopt_test.go","name":"adoptSetThingYAML"},"kwargs":"sha256:56b7df76e829fc4ba6f06ade9934c489c680031cea2b9e4924309fdf192459ce","body":"sha256:42da831442e90e8cba8a4dcd86ec679633dd0c634f8eb148baac166437a32a34"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/adopt_test.go","name":"TestAdopt"},"kwargs":"sha256:3993d0f81720e3911d50e7ec721493504bea91ccf58a41396b8c2b535b96f594","body":"sha256:572a669aa020d2734ed2de1227b37af7bf64864f9589ffe9160c6019e2c3e504"}
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
