package evals

import (
	"bytes"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/recording"
)

// addColumnBinding is the expectation the rails fixture carries, authored from
// the record: title is a string with no default, completed a boolean defaulting
// to false, both not null.
func addColumnBinding() ActionBinding {
	return ActionBinding{
		Because: "the record says the table carries title, a string that is not null and carries no default, and completed, a boolean that is not null and defaults to false",
		Key:     []string{"name"},
		Invocations: []map[string]any{
			{"name": "title", "type": "string", "nullable": false, "default": "nil"},
			{"name": "completed", "type": "boolean", "nullable": false, "default": "false"},
		},
	}
}

func bound(kwargs ...map[string]any) Sample {
	s := Sample{Counts: map[string]int{"addColumn": len(kwargs)}}
	for _, k := range kwargs {
		s.Invocations = append(s.Invocations, recording.Invocation{Action: "addColumn", Kwargs: k})
	}
	return s
}

func measurementWithBindings(samples ...Sample) Measurement {
	m := Measurement{Model: Model{ID: "test-model", Engine: "mlx"}, Samples: samples, Concurrency: 1}
	m.Case.ID = "fixture"
	m.Case.Arm = "sedum"
	m.Case.Expect.Actions = map[string]int{"addColumn": 2}
	m.Case.Expect.Bindings = map[string]ActionBinding{"addColumn": addColumnBinding()}
	return m
}

// The case this whole record was written from. Both samples select addColumn
// exactly twice - a perfect selection row - and bind default wrong on both
// columns, which no count can see (prov-2026-2b121b62).
func TestAPerfectSelectionCanCarryAWrongArgument(t *testing.T) {
	m := measurementWithBindings(
		// What the described package set actually produced, twice.
		bound(
			map[string]any{"name": "title", "type": "string", "nullable": false, "default": `"false"`},
			map[string]any{"name": "completed", "type": "boolean", "nullable": false, "default": `"false"`},
		),
		bound(
			map[string]any{"name": "title", "type": "string", "nullable": false, "default": `"false"`},
			map[string]any{"name": "completed", "type": "boolean", "nullable": false, "default": `"false"`},
		),
	)

	// Selection is perfect and says nothing.
	rows, scored, _ := m.Scored()
	if scored != 2 || rows[0].Exact != 2 {
		t.Fatalf("selection scored %d exact over %d samples, want 2 and 2", rows[0].Exact, scored)
	}

	by := map[string]KwargResult{}
	results := m.ScoredBindings()
	if len(results) != 1 {
		t.Fatalf("scored %d actions, want 1", len(results))
	}
	for _, k := range results[0].Kwargs {
		by[k.Kwarg] = k
	}

	// The kwarg the arms differ on, wrong on both columns of both samples.
	if got := by["default"]; got.Correct != 0 || got.Scored != 4 {
		t.Errorf("default bound correctly %d/%d, want 0/4", got.Correct, got.Scored)
	}
	// Per kwarg rather than per invocation: the rest are right, and a rate that
	// named the invocation would have reported four failures instead of one
	// broken argument.
	for _, kwarg := range []string{"name", "type", "nullable"} {
		if got := by[kwarg]; got.Correct != got.Scored || got.Scored == 0 {
			t.Errorf("%s bound correctly %d/%d, want all of them", kwarg, got.Correct, got.Scored)
		}
	}
	if results[0].Missing != 0 || results[0].Unexpected != 0 {
		t.Errorf("both columns were answered; got %d missing and %d unexpected",
			results[0].Missing, results[0].Unexpected)
	}
}

// The string "false" is not the boolean false, and the empty string is not nil.
// Stringifying either side before comparing would erase exactly the distinction
// the four samples that motivated this turn on.
func TestComparisonDoesNotStringify(t *testing.T) {
	for _, tc := range []struct {
		name       string
		want, got  any
		wantEquals bool
	}{
		{"a quoted ruby string is not the literal", "false", `"false"`, false},
		{"the boolean is not the string", false, "false", false},
		{"no default is not the empty string", "nil", "", false},
		{"the undescribed arm's boolean default is correct", "false", "false", true},
		{"a bool kwarg matches a bool", false, false, true},
		// The one normalisation: encoding/json decodes every number to
		// float64 and YAML decodes 5 to int, which is a difference between two
		// decoders rather than between two values.
		{"an int expectation matches a json number", 5, float64(5), true},
		{"a number is still not a string", 5, "5", false},
	} {
		if got := equalBinding(tc.want, tc.got); got != tc.wantEquals {
			t.Errorf("%s: equalBinding(%#v, %#v) = %v", tc.name, tc.want, tc.got, got)
		}
	}
}

// An invocation whose key kwarg is wrong is not a near-miss. Pairing it with the
// expectation it least resembles would report a fumbled argument where the model
// wrote about a column nobody asked for.
func TestAWrongKeyIsAMissAndAnUnexpectedInvocation(t *testing.T) {
	m := measurementWithBindings(bound(
		map[string]any{"name": "title", "type": "string", "nullable": false, "default": "nil"},
		// Not a column the record asks for. Its other kwargs happen to match
		// what completed wanted, which is exactly what best-match pairing would
		// have been fooled by.
		map[string]any{"name": "is_done", "type": "boolean", "nullable": false, "default": "false"},
	))

	got := m.ScoredBindings()[0]
	if got.Missing != 1 {
		t.Errorf("missing is %d, want 1 - completed was never answered", got.Missing)
	}
	if got.Unexpected != 1 {
		t.Errorf("unexpected is %d, want 1 - is_done is not a column the record asks for", got.Unexpected)
	}
	// Only title was paired, so only title's kwargs were scored. The unmatched
	// invocation contributes no correct bindings on the strength of resembling
	// one.
	for _, k := range got.Kwargs {
		if k.Scored != 1 {
			t.Errorf("%s scored %d comparisons, want 1 - only one invocation paired", k.Kwarg, k.Scored)
		}
	}
}

// Only kwargs the expectation names are scored, so a fixture can adopt this one
// kwarg at a time rather than having to be fully specified first.
func TestAnUnnamedKwargIsNotScored(t *testing.T) {
	m := measurementWithBindings(bound(
		map[string]any{"name": "title", "type": "string", "nullable": false, "default": "nil", "resource": "todo", "stamp": "20260814000000"},
		map[string]any{"name": "completed", "type": "boolean", "nullable": false, "default": "false", "resource": "todo", "stamp": "20260814000000"},
	))

	for _, k := range m.ScoredBindings()[0].Kwargs {
		if k.Kwarg == "resource" || k.Kwarg == "stamp" {
			t.Errorf("%s was scored, but the expectation does not name it", k.Kwarg)
		}
	}

	// And what is named scores clean, so this fixture is a passing one - the
	// scoring is not simply reporting everything wrong.
	for _, k := range m.ScoredBindings()[0].Kwargs {
		if k.Correct != k.Scored {
			t.Errorf("%s bound %d/%d on a correct answer", k.Kwarg, k.Correct, k.Scored)
		}
	}
}

// A rejected answer has no invocation list. Scoring it would report a rejection
// as a wrong argument, which is the same conflation the harness got wrong once
// already for selection.
func TestOnlyValidSamplesAreScoredForBinding(t *testing.T) {
	m := measurementWithBindings(
		bound(
			map[string]any{"name": "title", "type": "string", "nullable": false, "default": "nil"},
			map[string]any{"name": "completed", "type": "boolean", "nullable": false, "default": "false"},
		),
		Sample{Invalid: true, Detail: "did not validate", Rules: []string{"missing_kwarg"}},
	)

	got := m.ScoredBindings()[0]

	// The sharp assertion: a rejected sample scored as an answer would have
	// matched neither expected invocation and reported two more misses.
	if got.Missing != 0 {
		t.Errorf("the rejected sample was scored as %d missing invocation(s), want 0", got.Missing)
	}
	// Two comparisons per kwarg - one per paired invocation - from the one
	// valid sample, and none from the rejected one.
	for _, k := range got.Kwargs {
		if k.Scored != 2 {
			t.Errorf("%s scored %d comparisons, want 2 - the valid sample's two columns and nothing from the rejection", k.Kwarg, k.Scored)
		}
	}
}

// A case naming no bindings measures exactly what it measured before the field
// existed, and its report says nothing about binding.
func TestACaseWithNoBindingsIsUnchanged(t *testing.T) {
	m := measurementWithBindings(bound(map[string]any{"name": "title"}))
	m.Case.Expect.Bindings = nil

	if got := m.ScoredBindings(); got != nil {
		t.Errorf("a case with no bindings scored %v", got)
	}

	var buf bytes.Buffer
	Report(&buf, m)
	if strings.Contains(buf.String(), "kwarg") {
		t.Errorf("the report carries a binding table for a case that declares none:\n%s", buf.String())
	}
}

// Binding is reported beside selection, never folded into it. A combined rate
// would hide the perfect selection carrying a wrong argument.
func TestTheReportKeepsBindingBesideSelection(t *testing.T) {
	m := measurementWithBindings(bound(
		map[string]any{"name": "title", "type": "string", "nullable": false, "default": `"false"`},
		map[string]any{"name": "completed", "type": "boolean", "nullable": false, "default": `"false"`},
	))

	var buf bytes.Buffer
	Report(&buf, m)
	out := buf.String()

	// The selection table still reports a perfect row.
	if !strings.Contains(out, "action") || !strings.Contains(out, "exact") {
		t.Errorf("the selection table is gone:\n%s", out)
	}
	// And the binding table reports default at zero, in the same report.
	if !strings.Contains(out, "kwarg") || !strings.Contains(out, "default") {
		t.Errorf("the binding table is missing:\n%s", out)
	}
	if !strings.Contains(out, "0/2") {
		t.Errorf("default does not report as unbound:\n%s", out)
	}
}

// A binding expectation is a judgment about what is correct. Nothing mechanical
// can check that it is right, so what is enforced is that it says where it came
// from - a reader deciding whether to disagree should not have to find the
// provenance record.
func TestABindingExpectationStatesWhereItCameFrom(t *testing.T) {
	base := func() Case {
		return Case{
			ID: "x", Arm: "baseline", Models: []Model{{ID: "m", Engine: "test"}},
			Expect: Expectations{Bindings: map[string]ActionBinding{"addColumn": addColumnBinding()}},
		}
	}

	if err := base().validate("x.yaml"); err != nil {
		t.Fatalf("a complete binding was rejected: %v", err)
	}

	for _, tc := range []struct {
		name   string
		break_ func(*ActionBinding)
		want   string
	}{
		{"no because", func(b *ActionBinding) { b.Because = "  " }, "because"},
		{"no key", func(b *ActionBinding) { b.Key = nil }, "key"},
		{"no invocations", func(b *ActionBinding) { b.Invocations = nil }, "expects no invocations"},
		{"an invocation missing its key", func(b *ActionBinding) {
			b.Invocations[1] = map[string]any{"type": "boolean"}
		}, "key kwargs"},
		{"two invocations sharing a key", func(b *ActionBinding) {
			b.Invocations[1]["name"] = "title"
		}, "cannot pair"},
	} {
		c := base()
		b := c.Expect.Bindings["addColumn"]
		tc.break_(&b)
		c.Expect.Bindings["addColumn"] = b

		err := c.validate("x.yaml")
		if err == nil {
			t.Errorf("%s was accepted", tc.name)
			continue
		}
		if !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: error does not mention %q: %v", tc.name, tc.want, err)
		}
	}
}
