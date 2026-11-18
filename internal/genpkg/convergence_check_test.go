package genpkg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// copyGo copies generators/go into a temp generators directory with its
// manifest edited by edit.
func copyGoWith(t *testing.T, edit func(string) string) string {
	t.Helper()
	src := filepath.Join("..", "..", "generators", "go")
	dst := filepath.Join(t.TempDir(), "go")
	err := filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		if info.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if rel == "sedum.yaml" {
			b = []byte(edit(string(b)))
		}
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Dir(dst)
}

func TestConvergenceBlockIsValidated(t *testing.T) {
	cases := map[string]func(string) string{
		"seed_decl without the placeholder": func(s string) string {
			return strings.Replace(s, `panic(\"{{placeholder}}\")`, `panic(\"todo\")`, 1)
		},
		"missing helper":    func(s string) string { return strings.Replace(s, "  helper: \"{{ident}}{{Name}}Body\"\n", "", 1) },
		"bad method action": func(s string) string { return strings.Replace(s, "  note:", "  method_action: setThing\n  note:", 1) },
	}
	cases["unknown value"] = func(s string) string {
		return strings.Replace(s, `arg: "{{name}}"`, `arg: "{{nope}}"`, 1)
	}
	cases["placeholder transformed"] = func(s string) string {
		return strings.Replace(s, `panic(\"{{placeholder}}\")`, `panic(\"{{placeholder|upper}}\")`, 1)
	}
	cases["unknown format arg"] = func(s string) string { return strings.Replace(s, "format: [gofmt]", `format: [gofmt, ""]`, 1) }
	for name, edit := range cases {
		t.Run(name, func(t *testing.T) {
			_, findings, err := Load(copyGoWith(t, edit), Options{Only: []string{"go"}})
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range findings {
				if f.Kind == KindError && f.Rule == RuleConvergenceInvalid {
					return
				}
			}
			t.Fatalf("no %s finding: %v", RuleConvergenceInvalid, findings)
		})
	}
	// Unedited, the package loads clean.
	_, findings, _ := Load(copyGoWith(t, func(s string) string { return s }), Options{Only: []string{"go"}})
	for _, f := range findings {
		if f.Kind == KindError {
			t.Fatalf("unedited package has an error: %v", f)
		}
	}
}
