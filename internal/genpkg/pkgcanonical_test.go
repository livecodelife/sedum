package genpkg

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1e485f36","gen":"sha256:6d2b1ec4ae0cc1a45ab982a956905123aa16c8a8060c50f0a9e97c07f55d0b2f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/pkgcanonical_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e485f36","gen":"sha256:6d2b1ec4ae0cc1a45ab982a956905123aa16c8a8060c50f0a9e97c07f55d0b2f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/pkgcanonical_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e485f36","gen":"sha256:6d2b1ec4ae0cc1a45ab982a956905123aa16c8a8060c50f0a9e97c07f55d0b2f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/pkgcanonical_test.go","path":"github.com/livecodelife/provenance/record"}}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e485f36","gen":"sha256:6d2b1ec4ae0cc1a45ab982a956905123aa16c8a8060c50f0a9e97c07f55d0b2f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// TestPkgCanonicalBytes pins the bytes record.Canonical writes for every kind\n// of value a kwarg default can carry. The want string is what Sedum's own\n// encoder wrote before pkg/1 moved to the library's, captured while both\n// existed, so a difference here is a change to every pkg hash.\nfunc TestPkgCanonicalBytes(t *testing.T) {\n\tvalue := map[string]any{\n\t\t\"s\": \"plain\", \"e\": \"\", \"n\": nil, \"t\": true, \"f\": false,\n\t\t\"i\": 3, \"neg\": -7, \"i64\": int64(9007199254740993), \"u\": uint64(18446744073709551615),\n\t\t\"l\":      []any{\"a\", \"\", nil, true, 0, -1, []any{}},\n\t\t\"q\":      \"ctl\\x01 nl\\n uni\\u00e9\",\n\t\t\"m\":      map[string]any{\"z\": 1, \"a\": map[string]any{}},\n\t\t\"emptyl\": []any{},\n\t}\n\twant := `{\"e\":\"\",\"emptyl\":[],\"f\":false,\"i\":3,\"i64\":9007199254740993,\"l\":[\"a\",\"\",null,true,0,-1,[]],\"m\":{\"a\":{},\"z\":1},\"n\":null,\"neg\":-7,\"q\":\"ctl\\u0001 nl\\u000a uni` + \"\\u00e9\" + `\",\"s\":\"plain\",\"t\":true,\"u\":18446744073709551615}`\n\tif got := record.Canonical(value); got != want {\n\t\tt.Errorf(\"record.Canonical differs from the pkg/1 bytes\\n got: %s\\nwant: %s\", got, want)\n\t}\n}","file":"internal/genpkg/pkgcanonical_test.go","name":"TestPkgCanonicalBytes"}}
// TestPkgCanonicalBytes pins the bytes record.Canonical writes for every kind
// of value a kwarg default can carry. The want string is what Sedum's own
// encoder wrote before pkg/1 moved to the library's, captured while both
// existed, so a difference here is a change to every pkg hash.
func TestPkgCanonicalBytes(t *testing.T) {
	value := map[string]any{
		"s": "plain", "e": "", "n": nil, "t": true, "f": false,
		"i": 3, "neg": -7, "i64": int64(9007199254740993), "u": uint64(18446744073709551615),
		"l":      []any{"a", "", nil, true, 0, -1, []any{}},
		"q":      "ctl\x01 nl\n uni\u00e9",
		"m":      map[string]any{"z": 1, "a": map[string]any{}},
		"emptyl": []any{},
	}
	want := `{"e":"","emptyl":[],"f":false,"i":3,"i64":9007199254740993,"l":["a","",null,true,0,-1,[]],"m":{"a":{},"z":1},"n":null,"neg":-7,"q":"ctl\u0001 nl\u000a uni` + "\u00e9" + `","s":"plain","t":true,"u":18446744073709551615}`
	if got := record.Canonical(value); got != want {
		t.Errorf("record.Canonical differs from the pkg/1 bytes\n got: %s\nwant: %s", got, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e485f36","gen":"sha256:6d2b1ec4ae0cc1a45ab982a956905123aa16c8a8060c50f0a9e97c07f55d0b2f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// TestPkgCanonicalUnhashable checks that a default the library cannot encode\n// comes back as an error naming the package and action, never as a panic.\nfunc TestPkgCanonicalUnhashable(t *testing.T) {\n\tpkg := &Package{Name: \"p\"}\n\taction := &Action{Name: \"a\", Kwargs: map[string]Kwarg{\"x\": {Type: \"int\", Default: 1.5}}}\n\t_, err := pkgCanonical(pkg, action)\n\tif err == nil || !strings.Contains(err.Error(), \"package p: action a: cannot be hashed\") {\n\t\tt.Fatalf(\"pkgCanonical error = %v, want one naming package p, action a and cannot be hashed\", err)\n\t}\n}","file":"internal/genpkg/pkgcanonical_test.go","name":"TestPkgCanonicalUnhashable"}}
// TestPkgCanonicalUnhashable checks that a default the library cannot encode
// comes back as an error naming the package and action, never as a panic.
func TestPkgCanonicalUnhashable(t *testing.T) {
	pkg := &Package{Name: "p"}
	action := &Action{Name: "a", Kwargs: map[string]Kwarg{"x": {Type: "int", Default: 1.5}}}
	_, err := pkgCanonical(pkg, action)
	if err == nil || !strings.Contains(err.Error(), "package p: action a: cannot be hashed") {
		t.Fatalf("pkgCanonical error = %v, want one naming package p, action a and cannot be hashed", err)
	}
}

// /sedum:setFunction
