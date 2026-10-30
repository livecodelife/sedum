package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/declarative_test.go","path":"strings"},"kwargs":"sha256:68ce1e4c3f248b586e22c83c0115e2f7de17930b98d8cfeb2e260818680126e7","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/declarative_test.go","path":"testing"},"kwargs":"sha256:f18c87b97652bb929b8d1f82dd89bb5eb2d8a229e2a5c0f5a481328bfa63748e","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/declarative_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:1559f01245b3a471733dbc3563fe8418ad717052166d9dfa17cc3639f0f11c6e","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_test.go","name":"TestDeclarativePythonModuleCoverage"},"kwargs":"sha256:0c726e049d8547aa67e96151c091f979a9d86228b581855334e242d22faf524e","body":"sha256:db75e6ef187a80fca4a3e21102d024d689574c449673066f2db57edf59af2e20"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_test.go","name":"TestDeclarativeBlankEnd"},"kwargs":"sha256:b2af3150e2bc7e7a5a205689be43886526147baf49ec171e5b55ba29340e5d45","body":"sha256:c638b8e8ea8e3ccdf5858ad12ab494c8c79c1aec2603f3e91ae75ad6d211f1d9"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative_test.go","name":"TestDeclarativeEmptyInput"},"kwargs":"sha256:a23bc82367a075ca1403078a9e4741b28160b95c76858d9f8dfa140081f28140","body":"sha256:fba7f17b51fbc4f63c4e51a45a8a3a74a3e045f2c5780ec7249b3b25b06805b8"}
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
