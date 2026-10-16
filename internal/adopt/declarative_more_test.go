package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:0fbd41998cb5fb9c3f35cd2846f7ada93a6ff3c345794c055625bcdb02cebcff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_more_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:0fbd41998cb5fb9c3f35cd2846f7ada93a6ff3c345794c055625bcdb02cebcff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_more_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:0fbd41998cb5fb9c3f35cd2846f7ada93a6ff3c345794c055625bcdb02cebcff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative_more_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:0fbd41998cb5fb9c3f35cd2846f7ada93a6ff3c345794c055625bcdb02cebcff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeDedentEnd(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName:       \"py\",\n\t\tExtensions: []string{\".py\"},\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"dedent\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatalf(\"NewDeclarative: %v\", err)\n\t}\n\tsrc := []byte(\"def f():\\n    return 1\\n\\nx = 2\\n\")\n\tspans, err := d.Split(\"m.py\", src)\n\tif err != nil {\n\t\tt.Fatalf(\"Split: %v\", err)\n\t}\n\tvar units []Span\n\tfor _, s := range spans {\n\t\tif s.Kind == KindUnit {\n\t\t\tunits = append(units, s)\n\t\t}\n\t}\n\tif len(units) != 1 {\n\t\tt.Fatalf(\"got %d units, want 1\", len(units))\n\t}\n\tif units[0].Action != \"setFunction\" {\n\t\tt.Errorf(\"action = %q, want setFunction\", units[0].Action)\n\t}\n\tif !strings.HasPrefix(units[0].Text, \"def f():\") || strings.Contains(units[0].Text, \"x = 2\") {\n\t\tt.Errorf(\"unit text = %q, want the body only up to the dedent\", units[0].Text)\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeDedentEnd"}}
func TestDeclarativeDedentEnd(t *testing.T) {
	pkg := &genpkg.Package{
		Name:       "py",
		Extensions: []string{".py"},
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
		t.Fatalf("NewDeclarative: %v", err)
	}
	src := []byte("def f():\n    return 1\n\nx = 2\n")
	spans, err := d.Split("m.py", src)
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	var units []Span
	for _, s := range spans {
		if s.Kind == KindUnit {
			units = append(units, s)
		}
	}
	if len(units) != 1 {
		t.Fatalf("got %d units, want 1", len(units))
	}
	if units[0].Action != "setFunction" {
		t.Errorf("action = %q, want setFunction", units[0].Action)
	}
	if !strings.HasPrefix(units[0].Text, "def f():") || strings.Contains(units[0].Text, "x = 2") {
		t.Errorf("unit text = %q, want the body only up to the dedent", units[0].Text)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:0fbd41998cb5fb9c3f35cd2846f7ada93a6ff3c345794c055625bcdb02cebcff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeRegexpEndWithMatch(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName:       \"py\",\n\t\tExtensions: []string{\".py\"},\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   `^END$`,\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatalf(\"NewDeclarative: %v\", err)\n\t}\n\tspans, err := d.Split(\"m.py\", []byte(\"def f():\\n    x = 1\\nEND\\ntail\\n\"))\n\tif err != nil {\n\t\tt.Fatalf(\"Split: %v\", err)\n\t}\n\tvar units, residual []Span\n\tfor _, s := range spans {\n\t\tswitch s.Kind {\n\t\tcase KindUnit:\n\t\t\tunits = append(units, s)\n\t\tcase KindResidual:\n\t\t\tresidual = append(residual, s)\n\t\t}\n\t}\n\tif len(units) != 1 {\n\t\tt.Fatalf(\"got %d units, want 1\", len(units))\n\t}\n\tif units[0].Kwargs[\"name\"] != \"f\" {\n\t\tt.Errorf(\"name = %v, want f\", units[0].Kwargs[\"name\"])\n\t}\n\tif strings.Contains(units[0].Text, \"tail\") {\n\t\tt.Errorf(\"unit text %q runs past the end match\", units[0].Text)\n\t}\n\tfound := false\n\tfor _, s := range residual {\n\t\tif strings.Contains(s.Text, \"tail\") {\n\t\t\tfound = true\n\t\t}\n\t}\n\tif !found {\n\t\tt.Errorf(\"text after the end match should be residual\")\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeRegexpEndWithMatch"}}
func TestDeclarativeRegexpEndWithMatch(t *testing.T) {
	pkg := &genpkg.Package{
		Name:       "py",
		Extensions: []string{".py"},
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
		t.Fatalf("NewDeclarative: %v", err)
	}
	spans, err := d.Split("m.py", []byte("def f():\n    x = 1\nEND\ntail\n"))
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	var units, residual []Span
	for _, s := range spans {
		switch s.Kind {
		case KindUnit:
			units = append(units, s)
		case KindResidual:
			residual = append(residual, s)
		}
	}
	if len(units) != 1 {
		t.Fatalf("got %d units, want 1", len(units))
	}
	if units[0].Kwargs["name"] != "f" {
		t.Errorf("name = %v, want f", units[0].Kwargs["name"])
	}
	if strings.Contains(units[0].Text, "tail") {
		t.Errorf("unit text %q runs past the end match", units[0].Text)
	}
	found := false
	for _, s := range residual {
		if strings.Contains(s.Text, "tail") {
			found = true
		}
	}
	if !found {
		t.Errorf("text after the end match should be residual")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:0fbd41998cb5fb9c3f35cd2846f7ada93a6ff3c345794c055625bcdb02cebcff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeRegexpEndWithoutMatchRunsToEOF(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName:       \"py\",\n\t\tExtensions: []string{\".py\"},\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   `^NEVER$`,\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatalf(\"NewDeclarative: %v\", err)\n\t}\n\tsrc := []byte(\"def f():\\n    x = 1\\n    y = 2\\n\")\n\tspans, err := d.Split(\"m.py\", src)\n\tif err != nil {\n\t\tt.Fatalf(\"Split: %v\", err)\n\t}\n\tvar units []Span\n\tfor _, s := range spans {\n\t\tif s.Kind == KindUnit {\n\t\t\tunits = append(units, s)\n\t\t}\n\t}\n\tif len(units) != 1 {\n\t\tt.Fatalf(\"got %d units, want 1\", len(units))\n\t}\n\tif units[0].Text != string(src) {\n\t\tt.Errorf(\"unit text = %q, want the whole file %q\", units[0].Text, string(src))\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeRegexpEndWithoutMatchRunsToEOF"}}
func TestDeclarativeRegexpEndWithoutMatchRunsToEOF(t *testing.T) {
	pkg := &genpkg.Package{
		Name:       "py",
		Extensions: []string{".py"},
		Adopt: &genpkg.Adopt{
			UnitStart: `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:   `^NEVER$`,
			Actions:   map[string]string{"def": "setFunction", "class": "setType"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
			"setType":     {Name: "setType", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	d, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatalf("NewDeclarative: %v", err)
	}
	src := []byte("def f():\n    x = 1\n    y = 2\n")
	spans, err := d.Split("m.py", src)
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	var units []Span
	for _, s := range spans {
		if s.Kind == KindUnit {
			units = append(units, s)
		}
	}
	if len(units) != 1 {
		t.Fatalf("got %d units, want 1", len(units))
	}
	if units[0].Text != string(src) {
		t.Errorf("unit text = %q, want the whole file %q", units[0].Text, string(src))
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:0fbd41998cb5fb9c3f35cd2846f7ada93a6ff3c345794c055625bcdb02cebcff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeAttachAbove(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName:       \"py\",\n\t\tExtensions: []string{\".py\"},\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart:   `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:     \"dedent\",\n\t\t\tAttachAbove: []string{`^@\\w+`},\n\t\t\tActions:     map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatalf(\"NewDeclarative: %v\", err)\n\t}\n\tspans, err := d.Split(\"m.py\", []byte(\"@dec\\ndef f():\\n    x = 1\\n\"))\n\tif err != nil {\n\t\tt.Fatalf(\"Split: %v\", err)\n\t}\n\tvar units []Span\n\tfor _, s := range spans {\n\t\tif s.Kind == KindUnit {\n\t\t\tunits = append(units, s)\n\t\t}\n\t}\n\tif len(units) != 1 || !strings.HasPrefix(units[0].Text, \"@dec\\n\") {\n\t\tt.Fatalf(\"decorator line directly above should join the unit, got %+v\", units)\n\t}\n\tspans, err = d.Split(\"m.py\", []byte(\"@dec\\n\\ndef g():\\n    y = 1\\n\"))\n\tif err != nil {\n\t\tt.Fatalf(\"Split: %v\", err)\n\t}\n\tfor _, s := range spans {\n\t\tif s.Kind == KindUnit && strings.Contains(s.Text, \"@dec\") {\n\t\t\tt.Errorf(\"decorator separated by a blank line should not join the unit: %q\", s.Text)\n\t\t}\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeAttachAbove"}}
func TestDeclarativeAttachAbove(t *testing.T) {
	pkg := &genpkg.Package{
		Name:       "py",
		Extensions: []string{".py"},
		Adopt: &genpkg.Adopt{
			UnitStart:   `^(?P<kind>def|class) (?P<name>\w+)`,
			UnitEnd:     "dedent",
			AttachAbove: []string{`^@\w+`},
			Actions:     map[string]string{"def": "setFunction", "class": "setType"},
		},
		Actions: map[string]*genpkg.Action{
			"setFunction": {Name: "setFunction", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
			"setType":     {Name: "setType", Kwargs: map[string]genpkg.Kwarg{"file": {Type: "string", Required: true}, "name": {Type: "string", Required: true}, "code": {Type: "literal", Required: true}}},
		},
	}
	d, err := NewDeclarative(pkg)
	if err != nil {
		t.Fatalf("NewDeclarative: %v", err)
	}
	spans, err := d.Split("m.py", []byte("@dec\ndef f():\n    x = 1\n"))
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	var units []Span
	for _, s := range spans {
		if s.Kind == KindUnit {
			units = append(units, s)
		}
	}
	if len(units) != 1 || !strings.HasPrefix(units[0].Text, "@dec\n") {
		t.Fatalf("decorator line directly above should join the unit, got %+v", units)
	}
	spans, err = d.Split("m.py", []byte("@dec\n\ndef g():\n    y = 1\n"))
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	for _, s := range spans {
		if s.Kind == KindUnit && strings.Contains(s.Text, "@dec") {
			t.Errorf("decorator separated by a blank line should not join the unit: %q", s.Text)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:0fbd41998cb5fb9c3f35cd2846f7ada93a6ff3c345794c055625bcdb02cebcff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeUnknownKindStaysResidual(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName:       \"py\",\n\t\tExtensions: []string{\".py\"},\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"dedent\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\", \"class\": \"setType\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t\t\"setType\":     {Name: \"setType\", Kwargs: map[string]genpkg.Kwarg{\"file\": {Type: \"string\", Required: true}, \"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatalf(\"NewDeclarative: %v\", err)\n\t}\n\tsrc := []byte(\"enum E:\\n    A = 1\\n\")\n\tspans, err := d.Split(\"m.py\", src)\n\tif err != nil {\n\t\tt.Fatalf(\"Split: %v\", err)\n\t}\n\tvar text string\n\tfor _, s := range spans {\n\t\tif s.Kind == KindUnit {\n\t\t\tt.Errorf(\"unknown kind should not become a unit: %q\", s.Text)\n\t\t}\n\t\tif s.Kind == KindResidual {\n\t\t\ttext += s.Text\n\t\t}\n\t}\n\tif text != string(src) {\n\t\tt.Errorf(\"residual text = %q, want %q\", text, string(src))\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeUnknownKindStaysResidual"}}
func TestDeclarativeUnknownKindStaysResidual(t *testing.T) {
	pkg := &genpkg.Package{
		Name:       "py",
		Extensions: []string{".py"},
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
		t.Fatalf("NewDeclarative: %v", err)
	}
	src := []byte("enum E:\n    A = 1\n")
	spans, err := d.Split("m.py", src)
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	var text string
	for _, s := range spans {
		if s.Kind == KindUnit {
			t.Errorf("unknown kind should not become a unit: %q", s.Text)
		}
		if s.Kind == KindResidual {
			text += s.Text
		}
	}
	if text != string(src) {
		t.Errorf("residual text = %q, want %q", text, string(src))
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:0fbd41998cb5fb9c3f35cd2846f7ada93a6ff3c345794c055625bcdb02cebcff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestDeclarativeKwargsAreOnlyTheDeclaredOnes(t *testing.T) {\n\tpkg := &genpkg.Package{\n\t\tName:       \"py\",\n\t\tExtensions: []string{\".py\"},\n\t\tAdopt: &genpkg.Adopt{\n\t\t\tUnitStart: `^(?P<kind>def|class) (?P<name>\\w+)`,\n\t\t\tUnitEnd:   \"dedent\",\n\t\t\tActions:   map[string]string{\"def\": \"setFunction\"},\n\t\t},\n\t\tActions: map[string]*genpkg.Action{\n\t\t\t\"setFunction\": {Name: \"setFunction\", Kwargs: map[string]genpkg.Kwarg{\"name\": {Type: \"string\", Required: true}, \"code\": {Type: \"literal\", Required: true}}},\n\t\t},\n\t}\n\td, err := NewDeclarative(pkg)\n\tif err != nil {\n\t\tt.Fatalf(\"NewDeclarative: %v\", err)\n\t}\n\tspans, err := d.Split(\"m.py\", []byte(\"def f():\\n    x = 1\\n\"))\n\tif err != nil {\n\t\tt.Fatalf(\"Split: %v\", err)\n\t}\n\tvar units []Span\n\tfor _, s := range spans {\n\t\tif s.Kind == KindUnit {\n\t\t\tunits = append(units, s)\n\t\t}\n\t}\n\tif len(units) != 1 {\n\t\tt.Fatalf(\"got %d units, want 1\", len(units))\n\t}\n\tif _, ok := units[0].Kwargs[\"file\"]; ok {\n\t\tt.Errorf(\"kwargs should not carry file when the action does not declare it: %v\", units[0].Kwargs)\n\t}\n\tif units[0].Kwargs[\"name\"] != \"f\" {\n\t\tt.Errorf(\"name = %v, want f\", units[0].Kwargs[\"name\"])\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestDeclarativeKwargsAreOnlyTheDeclaredOnes"}}
func TestDeclarativeKwargsAreOnlyTheDeclaredOnes(t *testing.T) {
	pkg := &genpkg.Package{
		Name:       "py",
		Extensions: []string{".py"},
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
		t.Fatalf("NewDeclarative: %v", err)
	}
	spans, err := d.Split("m.py", []byte("def f():\n    x = 1\n"))
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	var units []Span
	for _, s := range spans {
		if s.Kind == KindUnit {
			units = append(units, s)
		}
	}
	if len(units) != 1 {
		t.Fatalf("got %d units, want 1", len(units))
	}
	if _, ok := units[0].Kwargs["file"]; ok {
		t.Errorf("kwargs should not carry file when the action does not declare it: %v", units[0].Kwargs)
	}
	if units[0].Kwargs["name"] != "f" {
		t.Errorf("name = %v, want f", units[0].Kwargs["name"])
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-88c83fa3","gen":"sha256:0fbd41998cb5fb9c3f35cd2846f7ada93a6ff3c345794c055625bcdb02cebcff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestNewDeclarativeRejectsCommandAndMissingAdopt(t *testing.T) {\n\tcases := []*genpkg.Package{\n\t\t{Name: \"x\", Adopt: &genpkg.Adopt{Command: []string{\"splitter\"}}},\n\t\t{Name: \"x\"},\n\t}\n\tfor i, pkg := range cases {\n\t\tif _, err := NewDeclarative(pkg); err == nil {\n\t\t\tt.Errorf(\"case %d: NewDeclarative returned no error, want one\", i)\n\t\t}\n\t}\n}\n","file":"internal/adopt/declarative_more_test.go","name":"TestNewDeclarativeRejectsCommandAndMissingAdopt"}}
func TestNewDeclarativeRejectsCommandAndMissingAdopt(t *testing.T) {
	cases := []*genpkg.Package{
		{Name: "x", Adopt: &genpkg.Adopt{Command: []string{"splitter"}}},
		{Name: "x"},
	}
	for i, pkg := range cases {
		if _, err := NewDeclarative(pkg); err == nil {
			t.Errorf("case %d: NewDeclarative returned no error, want one", i)
		}
	}
}

// /sedum:setFunction
