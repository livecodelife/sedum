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

// The mirror of the check above: an action naming a marker nothing plants is a
// typo, and a marker nothing names is an injection point nothing can reach
// (prov-2026-a9e59197). Both are dead configuration, and both warn.
func TestMarkerNoActionTargetsWarns(t *testing.T) {
	// class_body_top is planted by the controller template and targeted by
	// addBeforeFilter. Dropping that action leaves the marker with nothing
	// that can fill it.
	files := mutated(map[string]*string{
		"rails/actions/actions.yaml": text(`actions:
  createControllerMethod:
    kwargs:
      controller: { type: string, required: true }
      name: { type: string, required: true }
      collection: { type: string, required: false }
    discriminator: name
    variants: [index, show]
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body
`),
		"rails/actions/addBeforeFilter.rb": nil,
	})

	_, findings := loadTree(t, files)

	f := findingFor(t, findings, RuleMarkerUnfilled)
	if f.Kind != KindWarning {
		t.Errorf("unfilled marker reported as %s, want warning", f.Kind)
	}
	if !strings.Contains(f.Message, "class_body_top") {
		t.Errorf("diagnostic does not name the marker: %s", f.Message)
	}
	if findings.HasErrors() {
		t.Error("an unfilled marker rejected the package; it is usable and the author may be mid-edit")
	}
	// --strict is the existing answer for a team that wants it enforced,
	// which is why the check ships without a suppression list.
	if !findings.Strict().HasErrors() {
		t.Error("--strict did not promote the unfilled marker to an error")
	}
}

// A region anchor names its markers through anchor_start and anchor_end. Not
// reading those would report both endpoints of every region as unfilled, which
// would train authors to ignore the warning.
func TestRegionEndpointsCountAsTargeted(t *testing.T) {
	files := mutated(map[string]*string{
		"rails/actions/actions.yaml": text(`actions:
  addBetween:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: region
    anchor_start: class_body_top
    anchor_end: class_body
`),
		"rails/actions/addBetween.rb":                      text("# between\n"),
		"rails/actions/addBeforeFilter.rb":                 nil,
		"rails/actions/createControllerMethod/index.rb":    nil,
		"rails/actions/createControllerMethod/show.rb":     nil,
		"rails/actions/createControllerMethod/_default.rb": nil,
	})

	_, findings := loadTree(t, files)

	for _, f := range findings {
		if f.Rule == RuleMarkerUnfilled {
			t.Errorf("a marker reached through a region anchor was reported unfilled: %s", f.Message)
		}
	}
}

// The valid package plants nothing it cannot fill, so the check is silent on
// it. This is what keeps the warning meaningful: it fires on a real gap rather
// than on every package that ships a _default template.
func TestFilledMarkersAreSilent(t *testing.T) {
	_, findings := loadTree(t, validPackage())

	for _, f := range findings {
		if f.Rule == RuleMarkerUnfilled {
			t.Errorf("a package whose markers are all targeted warned: %s", f.Message)
		}
	}
}

// A package declares the paths it does not write, so a record can name a file
// that is part of a change without Sedum halting on it (prov-2026-529954ab).
func TestUnmanagedDeclarationsAreValidated(t *testing.T) {
	withUnmanaged := func(entries string) map[string]string {
		return mutated(map[string]*string{
			"rails/sedum.yaml": text(`name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
  instantize: [plural, "prefix:@"]
unmanaged:
` + entries),
		})
	}

	// An unreadable pattern would match nothing at all, so a declaration
	// meant to keep Sedum out of a directory would keep it out of nothing.
	t.Run("unreadable pattern is rejected", func(t *testing.T) {
		_, findings := loadTree(t, withUnmanaged("  - \"config/[unclosed.rb\"\n"))

		f := findingFor(t, findings, RuleUnmanagedInvalid)
		if f.Kind != KindError {
			t.Errorf("unreadable unmanaged entry reported as %s, want error", f.Kind)
		}
	})

	t.Run("empty entry is rejected", func(t *testing.T) {
		_, findings := loadTree(t, withUnmanaged("  - \"\"\n"))

		if f := findingFor(t, findings, RuleUnmanagedInvalid); f.Kind != KindError {
			t.Errorf("empty unmanaged entry reported as %s, want error", f.Kind)
		}
	})

	// A package that ships a template for a path it also disowns has said
	// both that it knows how to write the file and that it does not.
	t.Run("templating a disowned path is rejected", func(t *testing.T) {
		_, findings := loadTree(t, withUnmanaged("  - \"app/models/**\"\n"))

		f := findingFor(t, findings, RuleUnmanagedContradiction)
		if f.Kind != KindError {
			t.Errorf("contradiction reported as %s, want error", f.Kind)
		}
		if !strings.Contains(f.Message, "app/models/{name}.rb") {
			t.Errorf("diagnostic does not name the template: %s", f.Message)
		}
	})

	// An action pointed at a file its own package says nothing writes.
	t.Run("action injecting into a disowned path is rejected", func(t *testing.T) {
		files := withUnmanaged("  - \"config/routes.rb\"\n")
		files["rails/actions/actions.yaml"] = `actions:
  addRoute:
    kwargs:
      name: { type: string, required: true }
    injects_into: "config/routes.rb"
    anchor: class_body
`
		delete(files, "rails/actions/addBeforeFilter.rb")
		delete(files, "rails/actions/createControllerMethod/index.rb")
		delete(files, "rails/actions/createControllerMethod/show.rb")
		delete(files, "rails/actions/createControllerMethod/_default.rb")
		files["rails/actions/addRoute.rb"] = "get \"/{{name}}\"\n"

		_, findings := loadTree(t, files)

		f := findingFor(t, findings, RuleUnmanagedContradiction)
		if !strings.Contains(f.Message, "addRoute") {
			t.Errorf("diagnostic does not name the action: %s", f.Message)
		}
	})

	// A declaration that contradicts nothing is silent, and the package is
	// usable with it.
	t.Run("a declaration naming nothing the package writes is accepted", func(t *testing.T) {
		set, findings := loadTree(t, withUnmanaged("  - Gemfile\n  - \"config/credentials/\"\n"))

		if findings.HasErrors() {
			t.Fatalf("a legal unmanaged declaration rejected the package: %v", findings)
		}
		pkg, _ := set.Lookup("rails")
		if len(pkg.Unmanaged) != 2 {
			t.Errorf("loaded %d unmanaged entries, want 2", len(pkg.Unmanaged))
		}

		// The union is what Phase 2 consults, since it runs before a path
		// has been resolved to any package.
		declarer, entry, ok := set.Unmanaged("Gemfile")
		if !ok || declarer != "rails" || entry != "Gemfile" {
			t.Errorf("Set.Unmanaged(Gemfile) = %q, %q, %v; want rails, Gemfile, true", declarer, entry, ok)
		}
		if _, _, ok := set.Unmanaged("app/models/todo.rb"); ok {
			t.Error("a path the package does write reported as unmanaged")
		}
	})
}
