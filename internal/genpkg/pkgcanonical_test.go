package genpkg

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e485f36","gen":"sha256:6d2b1ec4ae0cc1a45ab982a956905123aa16c8a8060c50f0a9e97c07f55d0b2f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/pkgcanonical_test.go","path":"strings"},"kwargs":"sha256:bf7d9c2728eab53921a66ea7a73f7f3b028ed5ca3d742b19c37ce6cd031e4564","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e485f36","gen":"sha256:6d2b1ec4ae0cc1a45ab982a956905123aa16c8a8060c50f0a9e97c07f55d0b2f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/pkgcanonical_test.go","path":"testing"},"kwargs":"sha256:0f4bd3e78329f9cf30a2e23799fb7db15bde98434e6b668a16f857b5a4d2187e","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e485f36","gen":"sha256:6d2b1ec4ae0cc1a45ab982a956905123aa16c8a8060c50f0a9e97c07f55d0b2f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/pkgcanonical_test.go","path":"github.com/livecodelife/provenance/record"},"kwargs":"sha256:68cda5b623156855c48e6951e79841bb8f60c3bbf23e36179713ecea1a84c40b","body":"sha256:acbf8628061291d11aa39e9ea64ecc91095a117057cfe6f1958ef2c437fac003"}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1e485f36","gen":"sha256:6d2b1ec4ae0cc1a45ab982a956905123aa16c8a8060c50f0a9e97c07f55d0b2f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/pkgcanonical_test.go","name":"TestPkgCanonicalBytes"},"kwargs":"sha256:8f6f5f770a58623726d6860df42b3392bdfded838391b9b11f0d9b4524064324","body":"sha256:14872667497bb21f81529e92331ea252fa13e0d32e4d31ed226d7f30564ea350"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1e485f36","gen":"sha256:6d2b1ec4ae0cc1a45ab982a956905123aa16c8a8060c50f0a9e97c07f55d0b2f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/pkgcanonical_test.go","name":"TestPkgCanonicalUnhashable"},"kwargs":"sha256:89676a7dd29c56d424060925c171b03dfe7f4f11bd1aa5462937aff0e1c145a0","body":"sha256:b71e14cd98dcb8d44efe8d9e8e6841452a87446551d70ef9a83280c392f8ffb4"}
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
