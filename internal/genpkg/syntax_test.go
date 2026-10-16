package genpkg

import (
	"strings"
	"testing"
)

// A package opts into a built-in syntax checker by name; a name the closed
// registry does not hold rejects the package (prov-2026-5403d162).
func TestSyntaxNamesARegisteredChecker(t *testing.T) {
	manifest := validPackage()["rails/sedum.yaml"]

	set, findings := loadTree(t, mutated(map[string]*string{"rails/sedum.yaml": text(manifest + "syntax: go\n")}))
	if findings.HasErrors() {
		t.Fatalf("syntax: go was rejected: %v", findings)
	}
	if pkg, _ := set.Lookup("rails"); pkg.Syntax != "go" {
		t.Errorf("Syntax = %q, want go", pkg.Syntax)
	}

	_, findings = loadTree(t, mutated(map[string]*string{"rails/sedum.yaml": text(manifest + "syntax: cobol\n")}))
	f := findingFor(t, findings, RuleSyntaxUnknown)
	if !strings.Contains(f.Message, "cobol") || !strings.Contains(f.Message, "go") {
		t.Errorf("the diagnostic should name the unknown syntax and the registered ones: %s", f.Message)
	}
}
