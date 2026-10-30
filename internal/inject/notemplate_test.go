package inject

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/inject/notemplate_test.go","path":"errors"},"kwargs":"sha256:52530d73df916250c061edec71c74486e507b3ae2096e90aa65379731b113806","body":"sha256:2b718237ba0a56467c8fd494cc131e942ef410e6228c03b644d70476c71862dc"}
import "errors"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/inject/notemplate_test.go","path":"testing"},"kwargs":"sha256:72ef73b79e420016fe9003beae4edca28f5d3cc2ff9d741d9f1fffbaea8bc3fe","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/inject/notemplate_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:fa3bb3a486e61ec1d757eb3e3e047d2805de1f65e8734522a177f56bffbe81ec","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/inject/notemplate_test.go","name":"TestExplainUnmatched"},"kwargs":"sha256:a8493d20a9f2438492a8e176c43670ee879a92aff8f100234963e0594fb928a7","body":"sha256:58ddb7998d1b13465a6d08ef11322db9801560532bff827e413b6838d4aeec4d"}
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
