package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativePythonModuleCoverage(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"python\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tHeader:    `\\A# [a-z]+\\n`,\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"dedent\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\n\t\t\t\t\"file\": {Type: \"string\", Required: true},\n\t\t\t\t\"name\": {Type: \"string\", Required: true},\n\t\t\t\t\"code\": {Type: \"literal\", Required: true},\n\t\t\t}},\n\t\t\t\"setType\": {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\n\t\t\t\t\"file\": {Type: \"string\", Required: true},\n\t\t\t\t\"name\": {Type: \"string\", Required: true},\n\t\t\t\t\"code\": {Type: \"literal\", Required: true},\n\t\t\t}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tsrc := []byte(\"# doc\\nimport os\\n\\n@dec\\ndef f(x):\\n    return x\\n\\nclass C:\\n    pass\\n\\n\")\n\tspans, err := d.Split(\"mod.py\", src)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(spans) == 0 || spans[0].Kind != KindHeader {\n\t\tt.Fatalf(\"want header span first, got %d spans\", len(spans))\n\t}\n\tunits := 0\n\tvar b strings.Builder\n\tfor _, s := range spans {\n\t\tif s.Kind == KindUnit {\n\t\t\tunits++\n\t\t}\n\t\tb.WriteString(s.Text)\n\t}\n\tif units != 2 {\n\t\tt.Fatalf(\"want 2 units, got %d\", units)\n\t}\n\tif b.String() != string(src) {\n\t\tt.Fatalf(\"spans do not concatenate to source\")\n\t}\n}","file":"internal/adopt/declarative_test.go","name":"TestDeclarativePythonModuleCoverage"}}
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
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{
				"file": {Type: "string", Required: true},
				"name": {Type: "string", Required: true},
				"code": {Type: "literal", Required: true},
			}},
			"setType": {Name: "setType", Kwargs: map[string]genpkg.Kwarg{
				"file": {Type: "string", Required: true},
				"name": {Type: "string", Required: true},
				"code": {Type: "literal", Required: true},
			}},
		},
	}
	d, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("# doc\nimport os\n\n@dec\ndef f(x):\n    return x\n\nclass C:\n    pass\n\n")
	spans, err := d.Split("mod.py", src)
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) == 0 || spans[0].Kind != KindHeader {
		t.Fatalf("want header span first, got %d spans", len(spans))
	}
	units := 0
	var b strings.Builder
	for _, s := range spans {
		if s.Kind == KindUnit {
			units++
		}
		b.WriteString(s.Text)
	}
	if units != 2 {
		t.Fatalf("want 2 units, got %d", units)
	}
	if b.String() != string(src) {
		t.Fatalf("spans do not concatenate to source")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeBlankEnd(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName: \"python\",\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"blank\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\n\t\t\t\t\"file\": {Type: \"string\", Required: true},\n\t\t\t\t\"name\": {Type: \"string\", Required: true},\n\t\t\t\t\"code\": {Type: \"literal\", Required: true},\n\t\t\t}},\n\t\t\t\"setType\": {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\n\t\t\t\t\"file\": {Type: \"string\", Required: true},\n\t\t\t\t\"name\": {Type: \"string\", Required: true},\n\t\t\t\t\"code\": {Type: \"literal\", Required: true},\n\t\t\t}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tsrc := []byte(\"def a():\\n    return 1\\n\\ndef b():\\n    return 2\\n\")\n\tspans, err := d.Split(\"mod.py\", src)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tunits := 0\n\tvar b strings.Builder\n\tfor _, s := range spans {\n\t\tif s.Kind == KindUnit {\n\t\t\tunits++\n\t\t}\n\t\tb.WriteString(s.Text)\n\t}\n\tif units != 2 {\n\t\tt.Fatalf(\"want 2 units, got %d\", units)\n\t}\n\tif b.String() != string(src) {\n\t\tt.Fatalf(\"spans do not concatenate to source\")\n\t}\n}","file":"internal/adopt/declarative_test.go","name":"TestDeclarativeBlankEnd"}}
func TestDeclarativeBlankEnd(t *testing.T) {
	pkg := &genpkg.Package{
		Name: "python",
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:   "blank",
			Actions:   map[string]string{"def": "setFunction", "class": "setType"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{
				"file": {Type: "string", Required: true},
				"name": {Type: "string", Required: true},
				"code": {Type: "literal", Required: true},
			}},
			"setType": {Name: "setType", Kwargs: map[string]genpkg.Kwarg{
				"file": {Type: "string", Required: true},
				"name": {Type: "string", Required: true},
				"code": {Type: "literal", Required: true},
			}},
		},
	}
	d, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("def a():\n    return 1\n\ndef b():\n    return 2\n")
	spans, err := d.Split("mod.py", src)
	if err != nil {
		t.Fatal(err)
	}
	units := 0
	var b strings.Builder
	for _, s := range spans {
		if s.Kind == KindUnit {
			units++
		}
		b.WriteString(s.Text)
	}
	if units != 2 {
		t.Fatalf("want 2 units, got %d", units)
	}
	if b.String() != string(src) {
		t.Fatalf("spans do not concatenate to source")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeEmptyInput(t *testing.T) {\n\td, err := NewDeclarative(&genpkg.Package{Adopt: &genpkg.Adopt{UnitStart: `^(?P<kind>def) (?P<name>\\w+)`, UnitEnd: \"dedent\"}})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tspans, err := d.Split(\"x.py\", nil)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif spans == nil || len(spans) != 0 {\n\t\tt.Fatalf(\"want empty non-nil slice, got %v\", spans)\n\t}\n}","file":"internal/adopt/declarative_test.go","name":"TestDeclarativeEmptyInput"}}
func TestDeclarativeEmptyInput(t *testing.T) {
	d, err := NewDeclarative(&genpkg.Package{Adopt: &genpkg.Adopt{UnitStart: `^(?P<kind>def) (?P<name>\w+)`, UnitEnd: "dedent"}})
	if err != nil {
		t.Fatal(err)
	}
	spans, err := d.Split("x.py", nil)
	if err != nil {
		t.Fatal(err)
	}
	if spans == nil || len(spans) != 0 {
		t.Fatalf("want empty non-nil slice, got %v", spans)
	}
}

// /sedum:setFunction
