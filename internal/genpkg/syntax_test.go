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

// optional_anchors names anchors completeness never asks about; a name no
// template of the package plants rejects it (prov-2026-73d32e29).
func TestOptionalAnchorsMustBePlanted(t *testing.T) {
	manifest := validPackage()["rails/sedum.yaml"]

	set, findings := loadTree(t, mutated(map[string]*string{"rails/sedum.yaml": text(manifest + "optional_anchors: [class_body_top]\n")}))
	if findings.HasErrors() {
		t.Fatalf("a planted optional anchor was rejected: %v", findings)
	}
	if pkg, _ := set.Lookup("rails"); len(pkg.OptionalAnchors) != 1 || pkg.OptionalAnchors[0] != "class_body_top" {
		t.Errorf("OptionalAnchors = %v", pkg.OptionalAnchors)
	}

	_, findings = loadTree(t, mutated(map[string]*string{"rails/sedum.yaml": text(manifest + "optional_anchors: [nowhere]\n")}))
	f := findingFor(t, findings, RuleOptionalAnchorUnplanted)
	if !strings.Contains(f.Message, "nowhere") {
		t.Errorf("the diagnostic should name the anchor: %s", f.Message)
	}
}
