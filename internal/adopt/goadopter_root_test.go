package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter_root_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter_root_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoAdopterRootPackages(t *testing.T) {\n\tcases := []struct {\n\t\tpath, src string\n\t\twant      Kind\n\t}{\n\t\t{\"main.go\", \"package main\\n\\nfunc F() {}\\n\", KindHeader},\n\t\t{\"lib.go\", \"package lib\\n\\nfunc F() {}\\n\", KindResidual},\n\t\t{\"sub/lib.go\", \"package lib\\n\\nfunc F() {}\\n\", KindHeader},\n\t}\n\tfor _, tc := range cases {\n\t\tspans, err := GoAdopter{}.Split(tc.path, []byte(tc.src))\n\t\tif err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t\tif spans[0].Kind != tc.want {\n\t\t\tt.Errorf(\"%s: first span is %s, want %s\", tc.path, spans[0].Kind, tc.want)\n\t\t}\n\t}\n\tspans, _ := GoAdopter{}.Split(\"lib.go\", []byte(\"package lib\\n\\nfunc F() {}\\n\"))\n\tres := ResidualsOf([]byte(\"package lib\\n\\nfunc F() {}\\n\"), spans, genpkg.Comment{Prefix: \"//\"})\n\tif len(res) != 1 || res[0].Reason != ReasonOtherResidual || res[0].FirstLine != \"package lib\" {\n\t\tt.Errorf(\"residuals = %+v\", res)\n\t}\n}","file":"internal/adopt/goadopter_root_test.go","name":"TestGoAdopterRootPackages"}}
func TestGoAdopterRootPackages(t *testing.T) {
	cases := []struct {
		path, src string
		want      Kind
	}{
		{"main.go", "package main\n\nfunc F() {}\n", KindHeader},
		{"lib.go", "package lib\n\nfunc F() {}\n", KindResidual},
		{"sub/lib.go", "package lib\n\nfunc F() {}\n", KindHeader},
	}
	for _, tc := range cases {
		spans, err := GoAdopter{}.Split(tc.path, []byte(tc.src))
		if err != nil {
			t.Fatal(err)
		}
		if spans[0].Kind != tc.want {
			t.Errorf("%s: first span is %s, want %s", tc.path, spans[0].Kind, tc.want)
		}
	}
	spans, _ := GoAdopter{}.Split("lib.go", []byte("package lib\n\nfunc F() {}\n"))
	res := ResidualsOf([]byte("package lib\n\nfunc F() {}\n"), spans, genpkg.Comment{Prefix: "//"})
	if len(res) != 1 || res[0].Reason != ReasonOtherResidual || res[0].FirstLine != "package lib" {
		t.Errorf("residuals = %+v", res)
	}
}

// /sedum:setFunction
