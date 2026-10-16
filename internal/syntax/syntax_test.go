package syntax

import (
	"errors"
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
	if err := check.Check("x/a.go", src); err != nil {
		t.Errorf("clean source rejected: %v", err)
	}
}

func TestGoRejectsUnparsableSource(t *testing.T) {
	check, _ := Lookup("go")
	err := check.Check("x/a.go", "package x\n\nfunc A( {\n")
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
	err := check.Check("x/a.go", "package x\n\nfunc A()  {\n    return\n}\n")
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
