package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage_builtin_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage_builtin_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage_builtin_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestForPackageBuiltinGo(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"go\", Adopt: &genpkg.Adopt{Builtin: \"go\"}}\n\tad, err := ForPackage(pkg, false)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif _, ok := ad.(GoAdopter); !ok {\n\t\tt.Fatalf(\"got %T, want GoAdopter\", ad)\n\t}\n}","file":"internal/adopt/forpackage_builtin_test.go","name":"TestForPackageBuiltinGo"}}
func TestForPackageBuiltinGo(t *testing.T) {
	pkg := &genpkg.Package{Name: "go", Adopt: &genpkg.Adopt{Builtin: "go"}}
	ad, err := ForPackage(pkg, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := ad.(GoAdopter); !ok {
		t.Fatalf("got %T, want GoAdopter", ad)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestForPackageBuiltinUnknown(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"cobol\", Adopt: &genpkg.Adopt{Builtin: \"cobol\"}}\n\t_, err := ForPackage(pkg, false)\n\tif err == nil || !strings.Contains(err.Error(), \"does not have\") {\n\t\tt.Fatalf(\"err = %v, want it to mention \\\"does not have\\\"\", err)\n\t}\n}","file":"internal/adopt/forpackage_builtin_test.go","name":"TestForPackageBuiltinUnknown"}}
func TestForPackageBuiltinUnknown(t *testing.T) {
	pkg := &genpkg.Package{Name: "cobol", Adopt: &genpkg.Adopt{Builtin: "cobol"}}
	_, err := ForPackage(pkg, false)
	if err == nil || !strings.Contains(err.Error(), "does not have") {
		t.Fatalf("err = %v, want it to mention \"does not have\"", err)
	}
}

// /sedum:setFunction
