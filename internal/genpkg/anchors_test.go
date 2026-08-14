package genpkg

import (
	"strings"
	"testing"
)

// Anchor defects Phase 0 can catch (prov-2026-42f5eedd).
//
// Phase 7's hard error on a missing anchor stays exactly as it is: it means the
// package and a real file disagree, which is the one thing load time cannot
// know. Everything below is checkable with nothing but the package, so leaving
// it to a run that has already written files is a choice rather than a
// necessity - and in the reserved-name case it is not even a failure, because
// the mistake is silent.

// An action named "anchor" writes markers that are indistinguishable from the
// anchor points a file template plants. Its regions are invisible to the
// matching that would replace them, so it injects a duplicate on every run and
// says nothing.
func TestReservedActionNameIsRejected(t *testing.T) {
	_, findings := loadTree(t, mutated(map[string]*string{
		"rails/actions/actions.yaml": text(`actions:
  anchor:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body
`),
		"rails/actions/addBeforeFilter.rb":                 nil,
		"rails/actions/createControllerMethod/index.rb":    nil,
		"rails/actions/createControllerMethod/show.rb":     nil,
		"rails/actions/createControllerMethod/_default.rb": nil,
		"rails/actions/anchor.rb":                          text("# whatever\n"),
	}))

	f := findingFor(t, findings, RuleActionNameReserved)
	if f.Kind != KindError {
		t.Errorf("reserved action name reported as %s, want error", f.Kind)
	}
	if !strings.Contains(f.Message, "anchor") {
		t.Errorf("diagnostic does not name the reserved name: %s", f.Message)
	}
	if !findings.HasErrors() {
		t.Error("package carrying an action named anchor was not rejected")
	}
}

// Every other action name is fine, including ones that merely contain the
// reserved word.
func TestReservedNameIsExactRatherThanSubstring(t *testing.T) {
	_, findings := loadTree(t, mutated(map[string]*string{
		"rails/actions/actions.yaml": text(`actions:
  anchorHelper:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body
`),
		"rails/actions/addBeforeFilter.rb":                 nil,
		"rails/actions/createControllerMethod/index.rb":    nil,
		"rails/actions/createControllerMethod/show.rb":     nil,
		"rails/actions/createControllerMethod/_default.rb": nil,
		"rails/actions/anchorHelper.rb":                    text("# helper\n"),
	}))

	for _, f := range findings {
		if f.Rule == RuleActionNameReserved {
			t.Errorf("anchorHelper was rejected as a reserved name: %s", f.Message)
		}
	}
}

// An expression that does not parse is a defect in the package. A package
// carrying one used to be declared wholly valid and then fail partway through a
// run that had already written files.
func TestUnparseableAnchorPatternIsRejectedAtLoad(t *testing.T) {
	_, findings := loadTree(t, mutated(map[string]*string{
		"rails/actions/actions.yaml": text(`actions:
  addBeforeFilter:
    kwargs:
      controller: { type: string, required: true }
      filter: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: after_match
    anchor_pattern: "class ("
`),
		"rails/actions/createControllerMethod/index.rb":    nil,
		"rails/actions/createControllerMethod/show.rb":     nil,
		"rails/actions/createControllerMethod/_default.rb": nil,
	}))

	f := findingFor(t, findings, RuleAnchorInvalid)
	if f.Kind != KindError {
		t.Errorf("unparseable anchor_pattern reported as %s, want error", f.Kind)
	}
	if !strings.Contains(f.Message, "class (") {
		t.Errorf("diagnostic does not quote the pattern: %s", f.Message)
	}
}

// ^ and $ match the bounds of the whole file unless (?m) is set, so a pattern
// meant to find a line finds nothing, and Phase 7 reports it as a fault in the
// file rather than in the pattern.
func TestPatternWithoutLineModeWarns(t *testing.T) {
	withPattern := func(pattern string) map[string]string {
		return mutated(map[string]*string{
			"rails/actions/actions.yaml": text(`actions:
  addBeforeFilter:
    kwargs:
      controller: { type: string, required: true }
      filter: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: after_match
    anchor_pattern: "` + pattern + `"
`),
			"rails/actions/createControllerMethod/index.rb":    nil,
			"rails/actions/createControllerMethod/show.rb":     nil,
			"rails/actions/createControllerMethod/_default.rb": nil,
		})
	}

	t.Run("warns and does not reject", func(t *testing.T) {
		_, findings := loadTree(t, withPattern("^class .*Controller$"))

		f := findingFor(t, findings, RuleAnchorPatternLineMode)
		if f.Kind != KindWarning {
			t.Errorf("line-mode finding reported as %s, want warning", f.Kind)
		}
		if !strings.Contains(f.Message, "(?m)") {
			t.Errorf("diagnostic does not name the fix: %s", f.Message)
		}
		if findings.HasErrors() {
			t.Error("a legal pattern rejected the package")
		}
	})

	// The expression is never rewritten on the author's behalf, so declaring
	// (?m) is what silences it.
	t.Run("silent once (?m) is declared", func(t *testing.T) {
		_, findings := loadTree(t, withPattern("(?m)^class .*Controller$"))

		for _, f := range findings {
			if f.Rule == RuleAnchorPatternLineMode {
				t.Errorf("a line-anchored pattern still warned: %s", f.Message)
			}
		}
	})

	// A pattern using neither anchor has nothing to warn about.
	t.Run("silent when the pattern is not anchored", func(t *testing.T) {
		_, findings := loadTree(t, withPattern("class .*Controller"))

		for _, f := range findings {
			if f.Rule == RuleAnchorPatternLineMode {
				t.Errorf("an unanchored pattern warned: %s", f.Message)
			}
		}
	})
}

// ^ and $ mean something other than an anchor when escaped or inside a
// character class, and warning about those would train authors to ignore the
// warning.
func TestLineAnchorDetection(t *testing.T) {
	cases := []struct {
		pattern string
		want    bool
	}{
		{pattern: `^class`, want: true},
		{pattern: `Controller$`, want: true},
		{pattern: `class .*Controller`, want: false},
		{pattern: `\$id = 1`, want: false},
		{pattern: `[$^]literal`, want: false},
		{pattern: `[^a-z]+`, want: false},
		{pattern: `[^a-z]+$`, want: true},
		{pattern: `\^escaped`, want: false},
	}

	for _, tc := range cases {
		if got := lineAnchored(tc.pattern); got != tc.want {
			t.Errorf("lineAnchored(%q) = %v, want %v", tc.pattern, got, tc.want)
		}
	}
}
