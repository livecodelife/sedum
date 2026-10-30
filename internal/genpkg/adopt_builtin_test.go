package genpkg

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/adopt_builtin_test.go","path":"strings"},"kwargs":"sha256:86e64543b39566766edd779e233ade105ca0d7bde90d3fe0f9926b34dce1bbbe","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/adopt_builtin_test.go","path":"testing"},"kwargs":"sha256:d3a254ec8e0edf2492c8e27c9b533ee0a577006ceecea6cdc670229172caf233","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/adopt_builtin_test.go","name":"TestAdoptBuiltin"},"kwargs":"sha256:1d04b778a137c7be004d68c3f641775399885a78d5db500251555795c6600730","body":"sha256:d2c4c25239a52456fb317614e449b07fe1fbd0d1f0ca742c498437035b245b59"}
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
