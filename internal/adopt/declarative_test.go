package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativePythonModuleCoverage(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"python\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tHeader:    `\\A# [a-z]+\\n`,\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"dedent\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\tad, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tsrc := []byte(\"# doc\\nimport os\\n\\n@dec\\ndef f(x):\\n    return x\\n\\nclass C:\\n    pass\\n\\n\")\n\tspans, err := ad.Split(\"m.py\", src)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tvar got strings.Builder\n\theaders, units := 0, 0\n\tfor _, s := range spans {\n\t\tgot.WriteString(s.Text)\n\t\tswitch s.Kind {\n\t\tcase KindHeader:\n\t\t\theaders++\n\t\tcase KindUnit:\n\t\t\tunits++\n\t\t}\n\t}\n\tif got.String() != string(src) {\n\t\tt.Fatalf(\"spans do not concatenate to source\")\n\t}\n\tif headers != 1 || units != 2 {\n\t\tt.Fatalf(\"got %d headers and %d units, want 1 and 2\", headers, units)\n\t}\n}","file":"internal/adopt/declarative_test.go","name":"TestDeclarativePythonModuleCoverage"}}
func TestDeclarativePythonModuleCoverage(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "python",
		Adopt: &genpkg.Adopt{
			Header:    `\A# [a-z]+\n`,
			UnitStart: `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:   "dedent",
			Actions:   map[string]string{"def": "setFunction", "class": "setType"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
			"setType":     {Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	ad, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("# doc\nimport os\n\n@dec\ndef f(x):\n    return x\n\nclass C:\n    pass\n\n")
	spans, err := ad.Split("m.py", src)
	if err != nil {
		t.Fatal(err)
	}
	var got strings.Builder
	headers, units := 0, 0
	for _, s := range spans {
		got.WriteString(s.Text)
		switch s.Kind {
		case KindHeader:
			headers++
		case KindUnit:
			units++
		}
	}
	if got.String() != string(src) {
		t.Fatalf("spans do not concatenate to source")
	}
	if headers != 1 || units != 2 {
		t.Fatalf("got %d headers and %d units, want 1 and 2", headers, units)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeBlankEnd(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"python\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"blank\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\tad, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tsrc := []byte(\"def a():\\n    return 1\\n\\ndef b():\\n    return 2\\n\")\n\tspans, err := ad.Split(\"m.py\", src)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tvar got strings.Builder\n\tunits := 0\n\tfor _, s := range spans {\n\t\tgot.WriteString(s.Text)\n\t\tif s.Kind == KindUnit {\n\t\t\tunits++\n\t\t}\n\t}\n\tif got.String() != string(src) {\n\t\tt.Fatalf(\"spans do not concatenate to source\")\n\t}\n\tif units != 2 {\n\t\tt.Fatalf(\"got %d units, want 2\", units)\n\t}\n}","file":"internal/adopt/declarative_test.go","name":"TestDeclarativeBlankEnd"}}
func TestDeclarativeBlankEnd(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "python",
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:   "blank",
			Actions:   map[string]string{"def": "setFunction", "class": "setType"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
			"setType":     {Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	ad, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("def a():\n    return 1\n\ndef b():\n    return 2\n")
	spans, err := ad.Split("m.py", src)
	if err != nil {
		t.Fatal(err)
	}
	var got strings.Builder
	units := 0
	for _, s := range spans {
		got.WriteString(s.Text)
		if s.Kind == KindUnit {
			units++
		}
	}
	if got.String() != string(src) {
		t.Fatalf("spans do not concatenate to source")
	}
	if units != 2 {
		t.Fatalf("got %d units, want 2", units)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeEmptyInput(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"python\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"blank\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}}},\n\t\t},\n\t}\n\tad, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tspans, err := ad.Split(\"m.py\", []byte{})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif spans == nil || len(spans) != 0 {\n\t\tt.Fatalf(\"got %v, want empty non-nil slice\", spans)\n\t}\n}","file":"internal/adopt/declarative_test.go","name":"TestDeclarativeEmptyInput"}}
func TestDeclarativeEmptyInput(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "python",
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:   "blank",
			Actions:   map[string]string{"def": "setFunction"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}}},
		},
	}
	ad, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	spans, err := ad.Split("m.py", []byte{})
	if err != nil {
		t.Fatal(err)
	}
	if spans == nil || len(spans) != 0 {
		t.Fatalf("got %v, want empty non-nil slice", spans)
	}
}

// /sedum:setFunction
