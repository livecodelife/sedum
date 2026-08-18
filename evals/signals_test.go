package evals

import (
	"testing"

	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/recording"
	resolvepkg "github.com/calebcowen/sedum/internal/resolve"
)

// railsSet loads the package set the todo-rails cases run against, so the
// signals are exercised against a real package rather than a fixture written to
// agree with them.
func railsSet(t *testing.T) *genpkg.Set {
	t.Helper()
	set, findings, err := genpkg.Load("testdata/todo-rails/generators/defined", genpkg.Options{})
	if err != nil {
		t.Fatalf("loading the rails generators: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("the rails generators did not load: %s", f)
		}
	}
	return set
}

// planted builds the Phase 3 output for one file with the content it was
// created from, which is what plants the anchors an action fills.
func planted(t *testing.T, set *genpkg.Set, path, content string) resolvepkg.File {
	t.Helper()
	pkg, ok := set.Lookup("rails")
	if !ok {
		t.Fatal("the rails package did not load")
	}
	return resolvepkg.File{
		Resolution: resolvepkg.Resolution{RecordID: "PR-014", Path: path, Package: pkg},
		Rendered:   content,
	}
}

// The number comes from the package and the record, never from an expectation
// anybody wrote. That is the whole reason it sits beside the action counts:
// those can be fitted to the answer they score and this cannot
// (prov-2026-d61010a4).
func TestAnchorFillIsDerivedFromWhatTheTemplatesPlanted(t *testing.T) {
	set := railsSet(t)
	files := []resolvepkg.File{planted(t, set, "app/controllers/todos_controller.rb",
		"class TodosController\n  # sedum:anchor:actions\n\n  private\n  # sedum:anchor:private\nend\n")}

	empty := anchorFill("PR-014", files, nil, nil)
	if empty.Planted != 2 || empty.Filled != 0 {
		t.Fatalf("an empty selection filled %d of %d anchors, want 0 of 2", empty.Filled, empty.Planted)
	}
	if rate, ok := empty.Rate(); !ok || rate != 0 {
		t.Errorf("rate = %v, %v; want 0, true", rate, ok)
	}

	one := anchorFill("PR-014", files, []recording.Invocation{{
		Action: "addControllerAction",
		Kwargs: map[string]any{"resource": "todo", "verb": "index"},
	}}, nil)
	if one.Planted != 2 || one.Filled != 1 {
		t.Fatalf("one anchored action filled %d of %d, want 1 of 2", one.Filled, one.Planted)
	}
	if rate, _ := one.Rate(); rate != 0.5 {
		t.Errorf("rate = %v, want 0.5", rate)
	}
}

// A package that plants an anchor nothing targets has a ceiling below 100% that
// belongs to the package rather than to the answer. Counting it in the
// denominator would report that ceiling as the model falling short of one,
// which is precisely the attribution error this signal exists to prevent.
func TestAnAnchorNothingCanFillIsOutsideTheDenominator(t *testing.T) {
	set := railsSet(t)
	files := []resolvepkg.File{planted(t, set, "app/controllers/todos_controller.rb",
		"class TodosController\n  # sedum:anchor:actions\n  # sedum:anchor:extensions\nend\n")}

	fill := anchorFill("PR-014", files, []recording.Invocation{{
		Action: "addControllerAction",
		Kwargs: map[string]any{"resource": "todo", "verb": "index"},
	}}, nil)

	if fill.Planted != 1 {
		t.Fatalf("planted = %d, want 1: an anchor no action targets is not work the model omitted", fill.Planted)
	}
	if rate, ok := fill.Rate(); !ok || rate != 1 {
		t.Errorf("rate = %v, %v; want 1, true - the package's unfillable anchor depressed the model's number", rate, ok)
	}
}

// Nothing planted and nothing filled are not the same observation as a
// selection that filled none of what was there, and a bare 0.0 tells them
// apart from neither.
func TestAnchorFillIsAbsentRatherThanZeroWhenNothingWasPlanted(t *testing.T) {
	set := railsSet(t)
	files := []resolvepkg.File{planted(t, set, "app/models/todo.rb", "class Todo\nend\n")}

	fill := anchorFill("PR-014", files, nil, nil)
	if fill.Planted != 0 {
		t.Fatalf("planted = %d, want 0", fill.Planted)
	}
	if _, ok := fill.Rate(); ok {
		t.Error("a run that planted no anchors reported a fill rate")
	}
}
