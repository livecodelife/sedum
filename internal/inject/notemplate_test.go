package inject

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/inject/notemplate_test.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/inject/notemplate_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/inject/notemplate_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestExplainUnmatched(t *testing.T) {\n\tcause := errors.New(\"no marker\")\n\tgoTemplates := []string{\"{pkg}/{name}.go\", \"cmd/{bin}/{name}.go\"}\n\tcases := []struct {\n\t\tname string\n\t\tinv  Invocation\n\t\twant string\n\t}{\n\t\t{\n\t\t\tname: \"root path, no matching template\",\n\t\t\tinv:  Invocation{Package: &genpkg.Package{Name: \"go\", FileTemplates: goTemplates}, Action: &genpkg.Action{Name: \"setFunction\"}, Path: \"greet.go\"},\n\t\t\twant: \"action setFunction cannot be injected into greet.go: no file template in package go matches this path, so the file was created empty and has none of the markers a template plants (no marker); go's file templates are: {pkg}/{name}.go, cmd/{bin}/{name}.go\",\n\t\t},\n\t\t{\n\t\t\tname: \"matching path keeps today's message\",\n\t\t\tinv:  Invocation{Package: &genpkg.Package{Name: \"go\", FileTemplates: goTemplates}, Action: &genpkg.Action{Name: \"setFunction\"}, Path: \"internal/greet.go\"},\n\t\t\twant: \"action setFunction cannot be injected into internal/greet.go: no marker; the file is not shaped the way the action assumes\",\n\t\t},\n\t\t{\n\t\t\tname: \"_default package keeps today's message\",\n\t\t\tinv:  Invocation{Package: &genpkg.Package{Name: \"go\", FileTemplates: []string{\"_default.go\"}}, Action: &genpkg.Action{Name: \"setFunction\"}, Path: \"greet.go\"},\n\t\t\twant: \"action setFunction cannot be injected into greet.go: no marker; the file is not shaped the way the action assumes\",\n\t\t},\n\t\t{\n\t\t\tname: \"package with no templates\",\n\t\t\tinv:  Invocation{Package: &genpkg.Package{Name: \"bare\"}, Action: &genpkg.Action{Name: \"setFunction\"}, Path: \"greet.go\"},\n\t\t\twant: \"action setFunction cannot be injected into greet.go: no file template in package bare matches this path, so the file was created empty and has none of the markers a template plants (no marker); bare ships no file templates\",\n\t\t},\n\t}\n\tfor _, tc := range cases {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\terr := explainUnmatched(tc.inv, cause)\n\t\t\tif err == nil || err.Error() != tc.want {\n\t\t\t\tt.Errorf(\"got %v, want %q\", err, tc.want)\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/inject/notemplate_test.go","name":"TestExplainUnmatched"}}
func TestExplainUnmatched(t *testing.T) {
	cause := errors.New("no marker")
	goTemplates := []string{"{pkg}/{name}.go", "cmd/{bin}/{name}.go"}
	cases := []struct {
		name string
		inv  Invocation
		want string
	}{
		{
			name: "root path, no matching template",
			inv:  Invocation{Package: &genpkg.Package{Name: "go", FileTemplates: goTemplates}, Action: &genpkg.Action{Name: "setFunction"}, Path: "greet.go"},
			want: "action setFunction cannot be injected into greet.go: no file template in package go matches this path, so the file was created empty and has none of the markers a template plants (no marker); go's file templates are: {pkg}/{name}.go, cmd/{bin}/{name}.go",
		},
		{
			name: "matching path keeps today's message",
			inv:  Invocation{Package: &genpkg.Package{Name: "go", FileTemplates: goTemplates}, Action: &genpkg.Action{Name: "setFunction"}, Path: "internal/greet.go"},
			want: "action setFunction cannot be injected into internal/greet.go: no marker; the file is not shaped the way the action assumes",
		},
		{
			name: "_default package keeps today's message",
			inv:  Invocation{Package: &genpkg.Package{Name: "go", FileTemplates: []string{"_default.go"}}, Action: &genpkg.Action{Name: "setFunction"}, Path: "greet.go"},
			want: "action setFunction cannot be injected into greet.go: no marker; the file is not shaped the way the action assumes",
		},
		{
			name: "package with no templates",
			inv:  Invocation{Package: &genpkg.Package{Name: "bare"}, Action: &genpkg.Action{Name: "setFunction"}, Path: "greet.go"},
			want: "action setFunction cannot be injected into greet.go: no file template in package bare matches this path, so the file was created empty and has none of the markers a template plants (no marker); bare ships no file templates",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := explainUnmatched(tc.inv, cause)
			if err == nil || err.Error() != tc.want {
				t.Errorf("got %v, want %q", err, tc.want)
			}
		})
	}
}

// /sedum:setFunction
