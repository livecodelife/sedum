package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_more_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_more_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_more_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeDedentEnd(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"py\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"dedent\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tsrc := []byte(\"def f():\\n    return 1\\nx = 2\\n\")\n\tspans, err := d.Split(\"a.py\", src)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tvar got []byte\n\tunits := 0\n\tfor _, s := range spans {\n\t\tgot = append(got, s.Text...)\n\t\tif s.Kind == KindUnit {\n\t\t\tunits++\n\t\t}\n\t}\n\tif string(got) != string(src) {\n\t\tt.Fatalf(\"spans do not cover the file: %q\", got)\n\t}\n\tif units != 1 {\n\t\tt.Fatalf(\"got %d units, want 1\", units)\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeDedentEnd"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeRegexpEndWithMatch(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"py\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   `^end$`,\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tsrc := []byte(\"def f():\\n    pass\\nend\\nx = 1\\n\")\n\tspans, err := d.Split(\"a.py\", src)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tvar got []byte\n\tunits := 0\n\tfor _, s := range spans {\n\t\tgot = append(got, s.Text...)\n\t\tif s.Kind == KindUnit {\n\t\t\tunits++\n\t\t}\n\t}\n\tif string(got) != string(src) {\n\t\tt.Fatalf(\"spans do not cover the file: %q\", got)\n\t}\n\tif units != 1 {\n\t\tt.Fatalf(\"got %d units, want 1\", units)\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeRegexpEndWithMatch"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeRegexpEndWithoutMatchRunsToEOF(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"py\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   `^END$`,\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tsrc := []byte(\"def f():\\n    pass\\n    return 1\\n\")\n\tspans, err := d.Split(\"a.py\", src)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tvar got []byte\n\tunits := 0\n\tfor _, s := range spans {\n\t\tgot = append(got, s.Text...)\n\t\tif s.Kind == KindUnit {\n\t\t\tunits++\n\t\t\tif s.Text != string(src) {\n\t\t\t\tt.Fatalf(\"unit ends at %q, want end of file\", s.Text)\n\t\t\t}\n\t\t}\n\t}\n\tif string(got) != string(src) {\n\t\tt.Fatalf(\"spans do not cover the file: %q\", got)\n\t}\n\tif units != 1 {\n\t\tt.Fatalf(\"got %d units, want 1\", units)\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeRegexpEndWithoutMatchRunsToEOF"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeAttachAbove(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"py\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart:   `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:     \"dedent\",\n\t\t\tAttachAbove: []string{`^@`},\n\t\t\tActions:     map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tsrc := []byte(\"@dec\\ndef f():\\n    pass\\n\\n@dec\\n\\ndef g():\\n    pass\\n\")\n\tspans, err := d.Split(\"a.py\", src)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tvar got []byte\n\tvar units []string\n\tfor _, s := range spans {\n\t\tgot = append(got, s.Text...)\n\t\tif s.Kind == KindUnit {\n\t\t\tunits = append(units, s.Text)\n\t\t}\n\t}\n\tif string(got) != string(src) {\n\t\tt.Fatalf(\"spans do not cover the file: %q\", got)\n\t}\n\tif len(units) != 2 {\n\t\tt.Fatalf(\"got %d units, want 2\", len(units))\n\t}\n\tif units[0] != \"@dec\\ndef f():\\n    pass\\n\" {\n\t\tt.Fatalf(\"decorator not attached above: %q\", units[0])\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeAttachAbove"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeUnknownKindStaysResidual(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"py\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"dedent\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tsrc := []byte(\"class C:\\n    pass\\n\")\n\tspans, err := d.Split(\"a.py\", src)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfor _, s := range spans {\n\t\tif s.Kind == KindUnit {\n\t\t\tt.Fatalf(\"unknown kind became a unit: %q\", s.Text)\n\t\t}\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeUnknownKindStaysResidual"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeKwargsAreOnlyTheDeclaredOnes(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"py\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"dedent\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tspans, err := d.Split(\"a.py\", []byte(\"def f():\\n    pass\\n\"))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfor _, s := range spans {\n\t\tif s.Kind != KindUnit {\n\t\t\tcontinue\n\t\t}\n\t\tif _, ok := s.Kwargs[\"file\"]; ok {\n\t\t\tt.Fatalf(\"undeclared file kwarg was set: %v\", s.Kwargs)\n\t\t}\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeKwargsAreOnlyTheDeclaredOnes"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestNewDeclarativeRejectsCommandAndMissingAdopt(t *testing.T) {\n\tcases := []struct {\n\t\tname string\n\t\tpkg  *genpkg.Package\n\t}{\n\t\t{\"command\", &genpkg.Package{Name: \"x\", Adopt: &genpkg.Adopt{Command: []string{\"splitter\"}}}},\n\t\t{\"missing adopt\", &genpkg.Package{Name: \"x\"}},\n\t}\n\tfor _, c := range cases {\n\t\tt.Run(c.name, func(t *testing.T) {\n\t\t\tif _, err := NewDeclarative(c.pkg); err == nil {\n\t\t\t\tt.Fatal(\"expected an error\")\n\t\t\t}\n\t\t})\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestNewDeclarativeRejectsCommandAndMissingAdopt"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:2a225d388295fddeb0df23abd21cde5b9563d98ef8fd88eba9e9f176affc56be","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestNewDeclarativeRejectsEmptyUnitEnd(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"py\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def) (?P<name>\\w+)`,\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\t_, err := NewDeclarative(pkg)\n\tif err == nil || !strings.Contains(err.Error(), \"unit_end is required\") {\n\t\tt.Fatalf(\"got error %v, want one containing %q\", err, \"unit_end is required\")\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestNewDeclarativeRejectsEmptyUnitEnd"}}
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
