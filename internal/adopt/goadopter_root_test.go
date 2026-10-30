package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter_root_test.go","path":"testing"},"kwargs":"sha256:4afff9780c52ccfdb49bab3f8a5316726892a12786ddc591b1aa6d90c8ef9d65","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter_root_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:7ae503be1ecf4cfaa3431bc241ad26283e7a17b3246e69b361166def1ed50a05","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/goadopter_root_test.go","name":"TestGoAdopterRootPackages"},"kwargs":"sha256:71d42119a4432ad5b254e88cd3b815675b85e7cdc7ac63456a9108f319592609","body":"sha256:2572205c58346a0bfb0bb4e31a6ee632e07f9acf248e3448b07d86956f94e435"}
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
