package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:42fdf246ebb6ab5bf7b3ee907f89f14c49b671a5bc996a27f4d0f9cc2397e626","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/span_header_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:42fdf246ebb6ab5bf7b3ee907f89f14c49b671a5bc996a27f4d0f9cc2397e626","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestValidHeaderKwargs(t *testing.T) {\n\tok := Span{Kind: KindHeader, Text: \"package x\\n\", Kwargs: map[string]any{\"pkg\": \"x\"}}\n\tif err := ok.Valid(); err != nil {\n\t\tt.Fatalf(\"header with kwargs: unexpected error: %v\", err)\n\t}\n\tbad := Span{Kind: KindHeader, Text: \"package x\\n\", Action: \"setFunction\"}\n\tif err := bad.Valid(); err == nil {\n\t\tt.Fatal(\"header with action: want error\")\n\t}\n\tres := Span{Kind: KindResidual, Text: \"x\", Kwargs: map[string]any{\"a\": 1}}\n\tif err := res.Valid(); err == nil {\n\t\tt.Fatal(\"residual with kwargs: want error\")\n\t}\n}","file":"internal/adopt/span_header_test.go","name":"TestValidHeaderKwargs"}}
func TestValidHeaderKwargs(t *testing.T) {
	ok := Span{Kind: KindHeader, Text: "package x\n", Kwargs: map[string]any{"pkg": "x"}}
	if err := ok.Valid(); err != nil {
		t.Fatalf("header with kwargs: unexpected error: %v", err)
	}
	bad := Span{Kind: KindHeader, Text: "package x\n", Action: "setFunction"}
	if err := bad.Valid(); err == nil {
		t.Fatal("header with action: want error")
	}
	res := Span{Kind: KindResidual, Text: "x", Kwargs: map[string]any{"a": 1}}
	if err := res.Valid(); err == nil {
		t.Fatal("residual with kwargs: want error")
	}
}

// /sedum:setFunction
