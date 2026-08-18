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

// A nil expectation says the kwarg must not be bound at all.
//
// It is the shape a package's own default creates: once addColumn's default
// carries nil for anything unbound (prov-2026-f03916ba), the correct answer for
// a column with no default is silence, and the expectation has to be able to
// say so. It costs no code - equalBinding falls through to want == got with
// both sides nil - which is exactly why it needs a test: a behaviour that works
// by falling through is one an unrelated edit can remove without noticing.
//
// The alternative considered was dropping the expectation. That scores nothing,
// so a model binding "" or "false" to a column that must carry no default would
// have gone unmeasured - which is the failure this expectation exists for.
func TestAnAbsentBindingIsExpectedWithNil(t *testing.T) {
	expectation := ActionBinding{
		Because: "the record says title carries no default, and the package supplies nil for a default nobody binds",
		Key:     []string{"name"},
		Invocations: []map[string]any{
			{"name": "title", "type": "string", "nullable": false, "default": nil},
		},
	}

	score := func(kwargs map[string]any) KwargResult {
		m := Measurement{Model: Model{ID: "test-model", Engine: "mlx"}, Samples: []Sample{bound(kwargs)}, Concurrency: 1}
		m.Case.ID = "fixture"
		m.Case.Arm = "sedum"
		m.Case.Expect.Actions = map[string]int{"addColumn": 1}
		m.Case.Expect.Bindings = map[string]ActionBinding{"addColumn": expectation}
		for _, k := range m.ScoredBindings()[0].Kwargs {
			if k.Kwarg == "default" {
				return k
			}
		}
		t.Fatal("default was not scored at all; a nil expectation must still be an observation")
		return KwargResult{}
	}

	omitted := score(map[string]any{"name": "title", "type": "string", "nullable": false})
	if omitted.Scored != 1 || omitted.Correct != 1 {
		t.Errorf("omitting the kwarg scored %d/%d, want 1/1", omitted.Correct, omitted.Scored)
	}

	// The two ways a model has actually got this wrong, both of which produce a
	// migration that parses and is not what the record asked for.
	for _, wrong := range []any{"", `""`, "nil", "false"} {
		got := score(map[string]any{"name": "title", "type": "string", "nullable": false, "default": wrong})
		if got.Correct != 0 {
			t.Errorf("binding %#v to a kwarg expected absent scored %d/%d, want 0 correct",
				wrong, got.Correct, got.Scored)
		}
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

// The rate is cited per sample. Two invocations inside one answer are one
// observation seen twice: a model that has decided on "false" writes it in both
// columns, and counting that as two failures buys confidence nobody paid for
// (prov-2026-7cb96bf0).
func TestABindingRateCountsSamplesNotInvocations(t *testing.T) {
	wrongBoth := bound(
		map[string]any{"name": "title", "type": "string", "nullable": false, "default": `"false"`},
		map[string]any{"name": "completed", "type": "boolean", "nullable": false, "default": `"false"`},
	)
	m := measurementWithBindings(wrongBoth, wrongBoth)

	var def KwargResult
	for _, k := range m.ScoredBindings()[0].Kwargs {
		if k.Kwarg == "default" {
			def = k
		}
	}

	// Two samples, each binding it wrong twice.
	if def.Samples != 2 || def.CleanSamples != 0 {
		t.Errorf("default bound in %d/%d samples, want 0/2", def.CleanSamples, def.Samples)
	}
	// The invocation tally survives, because 0/4 and a sample that got one of
	// two wrong are different findings.
	if def.Correct != 0 || def.Scored != 4 {
		t.Errorf("invocation tally is %d/%d, want 0/4", def.Correct, def.Scored)
	}

	// The point of the whole record: the interval is the one two samples give,
	// not the tighter one four correlated comparisons would have implied.
	got := wilson(def.CleanSamples, def.Samples)
	want := wilson(0, 2)
	if got.String() != want.String() {
		t.Errorf("interval is %s, want the two-sample %s", got, want)
	}
	if inflated := wilson(0, 4); got.High <= inflated.High {
		t.Errorf("the sample interval %s is no wider than the invocation interval %s; the clustering is still being counted as independence", got, inflated)
	}
}

// A kwarg is bound in a sample only when every invocation of it in that sample
// is correct. One right and one wrong is not half an observation.
func TestASampleIsCleanOnlyIfEveryInvocationIs(t *testing.T) {
	m := measurementWithBindings(bound(
		map[string]any{"name": "title", "type": "string", "nullable": false, "default": "nil"},
		map[string]any{"name": "completed", "type": "boolean", "nullable": false, "default": `"false"`},
	))

	for _, k := range m.ScoredBindings()[0].Kwargs {
		if k.Kwarg != "default" {
			continue
		}
		if k.CleanSamples != 0 || k.Samples != 1 {
			t.Errorf("default bound in %d/%d samples, want 0/1 - one of its two invocations was wrong", k.CleanSamples, k.Samples)
		}
		// And the invocation tally is what says only one of the two was wrong,
		// which is the reason it is kept.
		if k.Correct != 1 || k.Scored != 2 {
			t.Errorf("invocation tally is %d/%d, want 1/2", k.Correct, k.Scored)
		}
	}
}

// The report puts an interval on the sample rate and none on the invocation
// ratio. An interval on a clustered count is the original error with an extra
// column.
func TestOnlyTheSampleRateCarriesAnInterval(t *testing.T) {
	m := measurementWithBindings(bound(
		map[string]any{"name": "title", "type": "string", "nullable": false, "default": `"false"`},
		map[string]any{"name": "completed", "type": "boolean", "nullable": false, "default": `"false"`},
	))

	var buf bytes.Buffer
	Report(&buf, m)

	var line string
	for _, l := range strings.Split(buf.String(), "\n") {
		if strings.Contains(l, "default") {
			line = l
		}
	}
	if line == "" {
		t.Fatalf("no default row:\n%s", buf.String())
	}

	// One sample, bound wrong: 0/1 with its interval, then 0/4... no - 0/2
	// invocations, bare.
	if !strings.Contains(line, wilson(0, 1).String()) {
		t.Errorf("the row does not carry the sample interval:\n%s", line)
	}
	if !strings.HasSuffix(strings.TrimSpace(line), "0/2") {
		t.Errorf("the row does not end in a bare invocation ratio:\n%s", line)
	}
	// Exactly one bracketed interval on the row.
	if n := strings.Count(line, "["); n != 1 {
		t.Errorf("the row carries %d intervals, want 1 - only the sample rate gets one:\n%s", n, line)
	}
}

// Some distinctions the package erases. This standard's transforms are
// [plural, snake] and [singular, snake], so todo and todos render identically
// and scoring one wrong would measure the fixture author's preference rather
// than the model's correctness (prov-2026-66b0116d).
func TestAKwargMayExpectAnyOfSeveralLiterals(t *testing.T) {
	either := []any{"todo", "todos"}

	for _, got := range []any{"todo", "todos"} {
		if !equalBinding(either, got) {
			t.Errorf("%q did not match [todo todos]", got)
		}
	}
	// A list matches its members and nothing else - it is an enumeration, not a
	// pattern.
	for _, got := range []any{"Todo", "todoes", "", "to"} {
		if equalBinding(either, got) {
			t.Errorf("%q matched [todo todos]", got)
		}
	}
	// And the comparison stays type-sensitive inside one, so a list of strings
	// never matches a boolean the way a stringifying compare would.
	if equalBinding([]any{"false", "nil"}, false) {
		t.Error("a boolean matched a list of strings")
	}
	if !equalBinding([]any{5, 6}, float64(6)) {
		t.Error("a json number did not match an int member")
	}
}

// A key is looked up rather than compared, so a key kwarg that could be two
// things cannot be looked up as either.
func TestAKeyKwargMayNotExpectSeveralValues(t *testing.T) {
	b := addColumnBinding()
	b.Invocations[0]["name"] = []any{"title", "titles"}
	c := Case{
		ID: "x", Arm: "baseline", Models: []Model{{ID: "m", Engine: "test"}},
		Expect: Expectations{Bindings: map[string]ActionBinding{"addColumn": b}},
	}

	err := c.validate("x.yaml")
	if err == nil {
		t.Fatal("a key kwarg expecting several values was accepted")
	}
	if !strings.Contains(err.Error(), "key") {
		t.Errorf("the error does not say the fault is in the key: %v", err)
	}
}

// Two spellings of one answer. Omitting a kwarg that declares nil as its
// default and binding that literal produce the same rendered text, so an
// expectation that admitted only one would report a correct answer as a wrong
// argument - the one thing a binding expectation must never do
// (prov-2026-0d30e826).
//
// It is the shape prov-2026-66b0116d settled for a resource whose singular and
// plural render identically, and it is bounded the same way: the list admits
// the spellings the package makes identical, and nothing else.
func TestAnExpectationMayAdmitAbsentOrTheDeclaredDefault(t *testing.T) {
	expectation := ActionBinding{
		Because: "the package declares nil as this kwarg's default, so omitting it and binding nil render the same text",
		Key:     []string{"name"},
		Invocations: []map[string]any{
			{"name": "title", "type": "string", "default": []any{nil, "nil"}},
		},
	}

	score := func(kwargs map[string]any) KwargResult {
		m := Measurement{Model: Model{ID: "test-model", Engine: "mlx"}, Samples: []Sample{bound(kwargs)}, Concurrency: 1}
		m.Case.ID = "fixture"
		m.Case.Arm = "sedum"
		m.Case.Expect.Actions = map[string]int{"addColumn": 1}
		m.Case.Expect.Bindings = map[string]ActionBinding{"addColumn": expectation}
		for _, k := range m.ScoredBindings()[0].Kwargs {
			if k.Kwarg == "default" {
				return k
			}
		}
		t.Fatal("default was not scored")
		return KwargResult{}
	}

	for _, right := range []map[string]any{
		{"name": "title", "type": "string"},                   // omitted
		{"name": "title", "type": "string", "default": "nil"}, // written out
	} {
		got := score(right)
		if got.Correct != 1 {
			t.Errorf("%v scored %d/%d, want it counted correct", right, got.Correct, got.Scored)
		}
	}

	// Still strict on everything the package does not make identical.
	for _, wrong := range []any{"", `""`, "false", "0"} {
		got := score(map[string]any{"name": "title", "type": "string", "default": wrong})
		if got.Correct != 0 {
			t.Errorf("binding %#v scored %d correct, want 0", wrong, got.Correct)
		}
	}
}
