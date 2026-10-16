package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage_test.go","path":"reflect"}}
import "reflect"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestForPackageDeclarative(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"py\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"dedent\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\"setFunction\": {Name: \"setFunction\"}},\n\t}\n\ta, err := ForPackage(pkg, false)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif _, ok := a.(*Declarative); !ok {\n\t\tt.Fatalf(\"got %T, want *Declarative\", a)\n\t}\n}","file":"internal/adopt/forpackage_test.go","name":"TestForPackageDeclarative"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestForPackageCommandNeedsPermission(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Adopt: &genpkg.Adopt{Command: []string{\"splitter\", \"-x\"}}}\n\t_, err := ForPackage(pkg, false)\n\twant := \"package py declares a command adopter (splitter -x); rerun with --allow-command-adopters to run it\"\n\tif err == nil || err.Error() != want {\n\t\tt.Fatalf(\"err = %v, want %q\", err, want)\n\t}\n}","file":"internal/adopt/forpackage_test.go","name":"TestForPackageCommandNeedsPermission"}}
func TestForPackageCommandNeedsPermission(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Adopt: &genpkg.Adopt{Command: []string{"splitter", "-x"}}}
	_, err := ForPackage(pkg, false)
	want := "package py declares a command adopter (splitter -x); rerun with --allow-command-adopters to run it"
	if err == nil || err.Error() != want {
		t.Fatalf("err = %v, want %q", err, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestForPackageCommandAllowed(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Adopt: &genpkg.Adopt{Command: []string{\"splitter\", \"-x\"}}}\n\ta, err := ForPackage(pkg, true)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tca, ok := a.(CommandAdopter)\n\tif !ok {\n\t\tt.Fatalf(\"got %T, want CommandAdopter\", a)\n\t}\n\tif !reflect.DeepEqual(ca.Argv, pkg.Adopt.Command) {\n\t\tt.Fatalf(\"Argv = %v, want %v\", ca.Argv, pkg.Adopt.Command)\n\t}\n}","file":"internal/adopt/forpackage_test.go","name":"TestForPackageCommandAllowed"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestForPackageNoAdopter(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\"}\n\t_, err := ForPackage(pkg, false)\n\twant := \"package py declares no adopter\"\n\tif err == nil || err.Error() != want {\n\t\tt.Fatalf(\"err = %v, want %q\", err, want)\n\t}\n}","file":"internal/adopt/forpackage_test.go","name":"TestForPackageNoAdopter"}}
func TestForPackageNoAdopter(t *testing.T) {
	pkg := &genpkg.Package{Name: "py"}
	_, err := ForPackage(pkg, false)
	want := "package py declares no adopter"
	if err == nil || err.Error() != want {
		t.Fatalf("err = %v, want %q", err, want)
	}
}

// /sedum:setFunction
