package syntax

import (
	"errors"
	"go/format"
	"strings"
	"testing"
)

// The registry is closed and names its members (prov-2026-5403d162).
func TestTheRegistryIsClosed(t *testing.T) {
	if _, ok := Lookup("go"); !ok {
		t.Fatal("go is not registered")
	}
	if _, ok := Lookup("cobol"); ok {
		t.Fatal("an unregistered syntax was found")
	}
	if got := strings.Join(Names(), ","); got != "go" {
		t.Errorf("Names() = %s, want go", got)
	}
}

func TestGoAcceptsGofmtCleanSource(t *testing.T) {
	check, _ := Lookup("go")
	src := "package x\n\nimport \"os\"\n\nfunc A() { _ = os.Args }\n"
	if err := check.Check(File{Path: "x/a.go", Content: src}); err != nil {
		t.Errorf("clean source rejected: %v", err)
	}
}

func TestGoRejectsUnparsableSource(t *testing.T) {
	check, _ := Lookup("go")
	err := check.Check(File{Path: "x/a.go", Content: "package x\n\nfunc A( {\n"})
	var problem *Problem
	if !errors.As(err, &problem) || problem.Formatted != "" {
		t.Fatalf("err = %v, want a Problem with no formatted form", err)
	}
	if !strings.Contains(err.Error(), "x/a.go") || !strings.Contains(err.Error(), "parse") {
		t.Errorf("the diagnostic should name the file and say it does not parse: %v", err)
	}
}

// Source that parses but is not gofmt-clean is rejected with the formatted form,
// so the model can be shown exactly what to bind instead.
func TestGoRejectsUnformattedSourceWithItsFormattedForm(t *testing.T) {
	check, _ := Lookup("go")
	err := check.Check(File{Path: "x/a.go", Content: "package x\n\nfunc A()  {\n    return\n}\n"})
	var problem *Problem
	if !errors.As(err, &problem) {
		t.Fatalf("err = %v, want a Problem", err)
	}
	if want := "package x\n\nfunc A() {\n\treturn\n}\n"; problem.Formatted != want {
		t.Errorf("Formatted = %q, want %q", problem.Formatted, want)
	}
	if !strings.Contains(err.Error(), "gofmt") {
		t.Errorf("the diagnostic should say the source is not gofmt-clean: %v", err)
	}
}

// A package qualifier the file does not import is flagged, so Phase 5
// re-prompts instead of landing code that does not build (prov-2026-0ce03d1c).
func TestGoFlagsQualifiersTheFileDoesNotImport(t *testing.T) {
	check, _ := Lookup("go")
	src := "package cmd\n\nimport \"os\"\n\nfunc Run() error {\n\t_, err := record.Load(os.Args[1])\n\treturn err\n}\n"
	err := check.Check(File{Path: "cmd/x.go", Content: src})
	if err == nil || !strings.Contains(err.Error(), "record") || strings.Contains(err.Error(), "os,") {
		t.Fatalf("err = %v, want only the unimported qualifier record named", err)
	}
}

func TestGoAcceptsQualifiersItsImportsSupply(t *testing.T) {
	check, _ := Lookup("go")
	src := `package cmd

import (
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/livecodelife/provenance/record"
	"github.com/spf13/cobra/v2"
	"gopkg.in/yaml.v3"
	"github.com/x/go-foo"
	"github.com/x/bar-go"
)

type T struct{ n int }

func (t T) N() int { return t.n }

func Run(c *cobra.Command) error {
	local := T{}
	_ = local.N()
	_ = T.N
	_ = foo.X
	_ = bar.Y
	_, _ = git.PlainOpen(os.Args[0])
	_ = yaml.Marshal
	_, err := record.Load("x")
	return err
}
`
	// Formatted first: this case is about qualifiers, not layout.
	formatted, err := format.Source([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if err := check.Check(File{Path: "cmd/x.go", Content: string(formatted)}); err != nil {
		t.Errorf("imports that supply every qualifier were flagged: %v", err)
	}
}

// A package-level variable declared in another file of the package is not a
// missing import; flagging it would re-prompt a correct answer.
func TestGoAcceptsQualifiersDeclaredInSiblingFiles(t *testing.T) {
	check, _ := Lookup("go")
	src := "package main\n\nfunc init() {\n\trootCmd.AddCommand(nil)\n}\n"
	if err := check.Check(File{Path: "cmd/p/a.go", Content: src}); err == nil {
		t.Fatal("rootCmd with no declaration anywhere was not flagged")
	}
	sibling := "package main\n\nimport \"github.com/spf13/cobra\"\n\nvar rootCmd = &cobra.Command{}\n"
	if err := check.Check(File{Path: "cmd/p/a.go", Content: src, Siblings: []string{sibling}}); err != nil {
		t.Errorf("a sibling's package-level variable was flagged: %v", err)
	}
}

// A dot import makes every name a possible import, so the file is not checked.
func TestGoSkipsTheQualifierCheckUnderADotImport(t *testing.T) {
	check, _ := Lookup("go")
	src := "package x\n\nimport . \"strings\"\n\nvar _ = anything.X\n\nvar _ = ToUpper\n"
	if err := check.Check(File{Path: "x/a.go", Content: src}); err != nil {
		t.Errorf("a file with a dot import was checked: %v", err)
	}
}
