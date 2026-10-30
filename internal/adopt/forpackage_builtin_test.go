package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/forpackage_builtin_test.go","path":"strings"},"kwargs":"sha256:f68b9f294598c9d6b11c7e5e8c19d63138a405b51f01b0c67f608aa69e70416a","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/forpackage_builtin_test.go","path":"testing"},"kwargs":"sha256:e77bbea07d9da9bb2ce96ce1fcb2a94e65348379f04a46aec5f7ae8ba6ce4303","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/forpackage_builtin_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:4878c6b78d2bb36b4221aec477c2c48e22cc62565afc086a8ca7a004201ac9cd","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/forpackage_builtin_test.go","name":"TestForPackageBuiltinGo"},"kwargs":"sha256:5f4e9255f485c017c7d8b900d38841becdd9446708d6954c61bdd83b9e59b246","body":"sha256:8a2f1761b640e386afb778d30174a99087261baa2f5f127ba0c0b61a72a692fb"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/forpackage_builtin_test.go","name":"TestForPackageBuiltinUnknown"},"kwargs":"sha256:f380c9598858b6c4526385b6dedb4f94ceda880ee421dc2273f58893b48069ef","body":"sha256:5a24c240c267319434c84a2068ca3fddc26b6d8c8e45c65246a97a3ae225c4ce"}
func TestForPackageBuiltinUnknown(t *testing.T) {
	pkg := &genpkg.Package{Name: "cobol", Adopt: &genpkg.Adopt{Builtin: "cobol"}}
	_, err := ForPackage(pkg, false)
	if err == nil || !strings.Contains(err.Error(), "does not have") {
		t.Fatalf("err = %v, want it to mention \"does not have\"", err)
	}
}

// /sedum:setFunction
