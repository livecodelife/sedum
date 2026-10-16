package genpkg

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/adopt_builtin_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/adopt_builtin_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAdoptBuiltin(t *testing.T) {\n\tcases := []struct {\n\t\tname, block, msg string\n\t}{\n\t\t{\"go\", \"adopt:\\n  builtin: go\\n\", \"\"},\n\t\t{\"unknown\", \"adopt:\\n  builtin: cobol\\n\", \"not a built-in adopter\"},\n\t\t{\"with command\", \"adopt:\\n  builtin: go\\n  command: [\\\"x\\\"]\\n\", \"cannot be combined\"},\n\t\t{\"with unit_start\", \"adopt:\\n  builtin: go\\n  unit_start: '^(?P<kind>def) (?P<name>\\\\w+)'\\n\", \"cannot be combined\"},\n\t\t{\"with header\", \"adopt:\\n  builtin: go\\n  header: '^//'\\n\", \"cannot be combined\"},\n\t\t{\"empty adopt\", \"adopt: {}\\n\", \"exactly one\"},\n\t}\n\tfor _, tc := range cases {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\tfiles := validPackage()\n\t\t\tfiles[\"rails/sedum.yaml\"] += tc.block\n\t\t\tset, findings := loadTree(t, files)\n\t\t\tvar got []string\n\t\t\tfor _, f := range findings {\n\t\t\t\tif f.Rule == RuleAdoptInvalid {\n\t\t\t\t\tgot = append(got, f.Message)\n\t\t\t\t}\n\t\t\t}\n\t\t\tif tc.msg == \"\" {\n\t\t\t\tif len(got) != 0 {\n\t\t\t\t\tt.Fatalf(\"unexpected adopt findings: %v\", got)\n\t\t\t\t}\n\t\t\t\tif set.Packages[0].Adopt.Builtin != \"go\" {\n\t\t\t\t\tt.Fatalf(\"Builtin = %q, want go\", set.Packages[0].Adopt.Builtin)\n\t\t\t\t}\n\t\t\t\treturn\n\t\t\t}\n\t\t\tfor _, m := range got {\n\t\t\t\tif strings.Contains(m, tc.msg) {\n\t\t\t\t\treturn\n\t\t\t\t}\n\t\t\t}\n\t\t\tt.Fatalf(\"no %s finding containing %q; got %v\", RuleAdoptInvalid, tc.msg, got)\n\t\t})\n\t}\n}","file":"internal/genpkg/adopt_builtin_test.go","name":"TestAdoptBuiltin"}}
func TestAdoptBuiltin(t *testing.T) {
	cases := []struct {
		name, block, msg string
	}{
		{"go", "adopt:\n  builtin: go\n", ""},
		{"unknown", "adopt:\n  builtin: cobol\n", "not a built-in adopter"},
		{"with command", "adopt:\n  builtin: go\n  command: [\"x\"]\n", "cannot be combined"},
		{"with unit_start", "adopt:\n  builtin: go\n  unit_start: '^(?P<kind>def) (?P<name>\\w+)'\n", "cannot be combined"},
		{"with header", "adopt:\n  builtin: go\n  header: '^//'\n", "cannot be combined"},
		{"empty adopt", "adopt: {}\n", "exactly one"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := validPackage()
			files["rails/sedum.yaml"] += tc.block
			set, findings := loadTree(t, files)
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
				if set.Packages[0].Adopt.Builtin != "go" {
					t.Fatalf("Builtin = %q, want go", set.Packages[0].Adopt.Builtin)
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
