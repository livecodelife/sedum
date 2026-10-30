package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/forpackage_test.go","path":"reflect"},"kwargs":"sha256:d1a91a2d847263f5c16ad40e6544b508e3dd0989c292b5e9fb091feabf91c745","body":"sha256:893cf368d31e1d8fe04189a0aba8627e8e69b39eafe10bea46d39a275e986fa2"}
import "reflect"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/forpackage_test.go","path":"testing"},"kwargs":"sha256:daa76f9a2123c7a9cf747704a3ded634482ff718edc6bd83f380ad9abdca75f5","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/forpackage_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:fd7ed75f4f803ba2b815822782772a2734bede7c0eea09e61d9591120e385b0d","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/forpackage_test.go","name":"TestForPackageDeclarative"},"kwargs":"sha256:641241ce29bc4cb3daae208c8ec4427a7ae061126e3d513545876ca8bbc3e9af","body":"sha256:a0fec31288d1a85f505d8e6d254aea2ff5de51beec8536ae62c09553ae6499e8"}
func TestForPackageDeclarative(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "py",
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def) (?P<name>\w+)`,
			UnitEnd:   "dedent",
			Actions:   map[string]string{"def": "setFunction"},
		},
		Actions: map[string]*genpkg.Action{"setFunction": {Name: "setFunction"}},
	}
	a, err := ForPackage(pkg, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := a.(*Declarative); !ok {
		t.Fatalf("got %T, want *Declarative", a)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/forpackage_test.go","name":"TestForPackageCommandNeedsPermission"},"kwargs":"sha256:efd192ba8f1121315a10898986ddeb9ac694244c3be039f3f58181b56b25a132","body":"sha256:574d66b69cf5d978f7d7ad75d194f74a8d24e5e421bfb491d793d8ecfcd3127f"}
func TestForPackageCommandNeedsPermission(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Adopt: &genpkg.Adopt{Command: []string{"splitter", "-x"}}}
	_, err := ForPackage(pkg, false)
	want := "package py declares a command adopter (splitter -x); rerun with --allow-command-adopters to run it"
	if err == nil || err.Error() != want {
		t.Fatalf("err = %v, want %q", err, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/forpackage_test.go","name":"TestForPackageCommandAllowed"},"kwargs":"sha256:0a97c72e524993f248a2f12c29fe6f58f45c9c694750d2f5c33a36e1e8bee408","body":"sha256:509e058210bb11fcaccb0289c9cba1b112a73ec7ba987546fcaab2cd1214730d"}
func TestForPackageCommandAllowed(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Adopt: &genpkg.Adopt{Command: []string{"splitter", "-x"}}}
	a, err := ForPackage(pkg, true)
	if err != nil {
		t.Fatal(err)
	}
	ca, ok := a.(CommandAdopter)
	if !ok {
		t.Fatalf("got %T, want CommandAdopter", a)
	}
	if !reflect.DeepEqual(ca.Argv, pkg.Adopt.Command) {
		t.Fatalf("Argv = %v, want %v", ca.Argv, pkg.Adopt.Command)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/forpackage_test.go","name":"TestForPackageNoAdopter"},"kwargs":"sha256:dc30caa5308e29bc8a73f4575400b53f82ef969066455ea03dcd6d004392002f","body":"sha256:f71c6e2ad0b0f02233f5328bfb6b4ac4588fbffa44134f2615963ce2f63f49f3"}
func TestForPackageNoAdopter(t *testing.T) {
	pkg := &genpkg.Package{Name: "py"}
	_, err := ForPackage(pkg, false)
	want := "package py declares no adopter"
	if err == nil || err.Error() != want {
		t.Fatalf("err = %v, want %q", err, want)
	}
}

// /sedum:setFunction
