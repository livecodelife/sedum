package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/declarative_more_test.go","path":"strings"},"kwargs":"sha256:e1da9725de12f06f53a693d9685839e4cdc707ddde6711b1da17827d09f21e9e","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/declarative_more_test.go","path":"testing"},"kwargs":"sha256:a40946d19da7695ea4867b1576ca43eaef3856bf1066aab496d081f994c67aa8","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/declarative_more_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:d9123429a752253155d4d7026a876388fd0f3e01e6417f31cfe4186cc2fa8453","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeDedentEnd"},"kwargs":"sha256:50f5e829318aa6e41d8b023c2ffd274d8b7f5df73f1c1c4cf6c83569e6b6fe33","body":"sha256:4df76f93ba6fb6aaada9132f77070d08e65214ba9bfe4c8eae5b2ba29e0f059d"}
func TestDeclarativeDedentEnd(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "py",
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:   "dedent",
			Actions:   map[string]string{"def": "setFunction", "class": "setType"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
			"setType":     {Name: "setType", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	d, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("def f():\n    return 1\nx = 2\n")
	spans, err := d.Split("a.py", src)
	if err != nil {
		t.Fatal(err)
	}
	var got []byte
	units := 0
	for _, s := range spans {
		got = append(got, s.Text...)
		if s.Kind == KindUnit {
			units++
		}
	}
	if string(got) != string(src) {
		t.Fatalf("spans do not cover the file: %q", got)
	}
	if units != 1 {
		t.Fatalf("got %d units, want 1", units)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeRegexpEndWithMatch"},"kwargs":"sha256:c9ffa82656bc36676bad4a86ff4700a160b2458e4ce728760ccbf60046a3c119","body":"sha256:986e6db0cce8f68609c20b04ba8727c691d845310f6f51f7c8870d514683201b"}
func TestDeclarativeRegexpEndWithMatch(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "py",
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:   `^end$`,
			Actions:   map[string]string{"def": "setFunction", "class": "setType"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
			"setType":     {Name: "setType", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	d, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("def f():\n    pass\nend\nx = 1\n")
	spans, err := d.Split("a.py", src)
	if err != nil {
		t.Fatal(err)
	}
	var got []byte
	units := 0
	for _, s := range spans {
		got = append(got, s.Text...)
		if s.Kind == KindUnit {
			units++
		}
	}
	if string(got) != string(src) {
		t.Fatalf("spans do not cover the file: %q", got)
	}
	if units != 1 {
		t.Fatalf("got %d units, want 1", units)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeRegexpEndWithoutMatchRunsToEOF"},"kwargs":"sha256:707f27346eec4d47dc5560c9bc0ba7577fea5a3d538bf4485cc9c2788429ea35","body":"sha256:0b83ccd81f58ea03147f89f0dd1fe303a6a6ee4f5e6a3f537bb01a49d6bf8c2b"}
func TestDeclarativeRegexpEndWithoutMatchRunsToEOF(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "py",
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:   `^END$`,
			Actions:   map[string]string{"def": "setFunction", "class": "setType"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
			"setType":     {Name: "setType", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	d, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("def f():\n    pass\n    return 1\n")
	spans, err := d.Split("a.py", src)
	if err != nil {
		t.Fatal(err)
	}
	var got []byte
	units := 0
	for _, s := range spans {
		got = append(got, s.Text...)
		if s.Kind == KindUnit {
			units++
			if s.Text != string(src) {
				t.Fatalf("unit ends at %q, want end of file", s.Text)
			}
		}
	}
	if string(got) != string(src) {
		t.Fatalf("spans do not cover the file: %q", got)
	}
	if units != 1 {
		t.Fatalf("got %d units, want 1", units)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeAttachAbove"},"kwargs":"sha256:329b020310b0b49f43b0d0d03174dbc7c41a4b5719d57cc413a54441bd128eb0","body":"sha256:579bf3fa5830d2a41894445a5290c621c94329f66db9f77e97dba5fc912d6815"}
func TestDeclarativeAttachAbove(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "py",
		Adopt: &genpkg.Adopt{
			UnitStart:   `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:     "dedent",
			AttachAbove: []string{`^@`},
			Actions:     map[string]string{"def": "setFunction", "class": "setType"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
			"setType":     {Name: "setType", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	d, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("@dec\ndef f():\n    pass\n\n@dec\n\ndef g():\n    pass\n")
	spans, err := d.Split("a.py", src)
	if err != nil {
		t.Fatal(err)
	}
	var got []byte
	var units []string
	for _, s := range spans {
		got = append(got, s.Text...)
		if s.Kind == KindUnit {
			units = append(units, s.Text)
		}
	}
	if string(got) != string(src) {
		t.Fatalf("spans do not cover the file: %q", got)
	}
	if len(units) != 2 {
		t.Fatalf("got %d units, want 2", len(units))
	}
	if units[0] != "@dec\ndef f():\n    pass\n" {
		t.Fatalf("decorator not attached above: %q", units[0])
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeUnknownKindStaysResidual"},"kwargs":"sha256:d69fa7b2ac1b1f9247b3d84734f2f72b08d99dcd21a61d9a4655287e48da0357","body":"sha256:c16029631ad0118426ab9cee28ea5dac8b100a2b90c5e97259cbdbb6dda19105"}
func TestDeclarativeUnknownKindStaysResidual(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "py",
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:   "dedent",
			Actions:   map[string]string{"def": "setFunction"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	d, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("class C:\n    pass\n")
	spans, err := d.Split("a.py", src)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range spans {
		if s.Kind == KindUnit {
			t.Fatalf("unknown kind became a unit: %q", s.Text)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeKwargsAreOnlyTheDeclaredOnes"},"kwargs":"sha256:3e8961df076f6d7102ac058e8a41c3fe178480eb9b8254f8c1a5fc85d368c4d4","body":"sha256:b87351d2a0b88beb0e252a7274631322c245a9febcdcbc543bd081a51cb5ff76"}
func TestDeclarativeKwargsAreOnlyTheDeclaredOnes(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "py",
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:   "dedent",
			Actions:   map[string]string{"def": "setFunction"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{"name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	d, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	spans, err := d.Split("a.py", []byte("def f():\n    pass\n"))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range spans {
		if s.Kind != KindUnit {
			continue
		}
		if _, ok := s.Kwargs["file"]; ok {
			t.Fatalf("undeclared file kwarg was set: %v", s.Kwargs)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_more_test.go","name":"TestNewDeclarativeRejectsCommandAndMissingAdopt"},"kwargs":"sha256:a8748d80e8bdf30c275c6e423cc70b5c62d0641207e9af53e027964f0d42ebb9","body":"sha256:9ea80f96ce95d797e99fde12884295250f35f3c2364e26639bb9521dd0ca4f92"}
func TestNewDeclarativeRejectsCommandAndMissingAdopt(t *testing.T) {
	cases := []struct {
		name string
		pkg  *genpkg.Package
	}{
		{"command", &genpkg.Package{Name: "x", Adopt: &genpkg.Adopt{Command: []string{"splitter"}}}},
		{"missing adopt", &genpkg.Package{Name: "x"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := NewDeclarative(c.pkg); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_more_test.go","name":"TestNewDeclarativeRejectsEmptyUnitEnd"},"kwargs":"sha256:5a76164082cebd0a4c7f294235a87c85b649f1b1732e1257364d1f88151100ac","body":"sha256:4fafad7925a2ea1209ac5e21f5d807b1c706a204ab125590b2876617ed8f90b3"}
func TestNewDeclarativeRejectsEmptyUnitEnd(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "py",
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def) (?P<name>\w+)`,
			Actions:   map[string]string{"def": "setFunction"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	_, err := NewDeclarative(pkg)
	if err == nil || !strings.Contains(err.Error(), "unit_end is required") {
		t.Fatalf("got error %v, want one containing %q", err, "unit_end is required")
	}
}

// /sedum:setFunction
