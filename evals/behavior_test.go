package evals

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/calebcowen/sedum/internal/recording"
)

// Behaviour is the only measurement here that runs somebody else's code, so
// what is tested is the seam rather than the applications: that a sample's own
// invocations reach the harness unchanged, that the three outcomes stay apart,
// and that a run which was never measured never reads as a run where nothing
// worked (prov-2026-83340ba0).
//
// Standing up Rails and Postgres inside `go test ./...` is not done. The
// harness is exercised end to end by running it, which is what evals/behavior's
// README documents.

func behaviorSample(b *BehaviorRun) Sample {
	return Sample{Counts: map[string]int{"addColumn": 1}, Behavior: b}
}

func caseWithTarget(m *Measurement, target string) {
	m.Case.ID = "fixture"
	m.Case.Arm = "sedum"
	m.Case.Expect.Actions = map[string]int{"addColumn": 1}
	m.Case.Expect.Behavior = &BehaviorExpectation{Target: target, PassRate: 1.0}
}

// The three outcomes are counted apart. A service that never booted and one
// that booted and answered wrongly have different fixes, and a single pass rate
// would report a broken generator package and a wrong one as the same number.
func TestBehaviorKeepsBrokenApartFromWrong(t *testing.T) {
	m := Measurement{Samples: []Sample{
		behaviorSample(&BehaviorRun{Outcome: "ok", Checks: 20, Passed: 20}),
		behaviorSample(&BehaviorRun{Outcome: "checks_failed", Checks: 20, Passed: 17,
			Failed: []string{"partial update keeps the title it never sent", "an empty update is a 400"}}),
		behaviorSample(&BehaviorRun{Outcome: "failed", FailedPhase: "build"}),
	}}
	caseWithTarget(&m, "todo-rails")

	got, measured := m.Behavior()
	if !measured {
		t.Fatal("three behaviour runs reported as not measured")
	}
	if got.Measured != 3 {
		t.Errorf("measured %d, want 3", got.Measured)
	}
	if got.Working != 1 || got.Disagreed != 1 || got.Broke != 1 {
		t.Errorf("working/disagreed/broke = %d/%d/%d, want 1/1/1",
			got.Working, got.Disagreed, got.Broke)
	}
	if got.Phases["build"] != 1 {
		t.Errorf("the failed phase was not counted: %v", got.Phases)
	}
	if got.Failures["an empty update is a 400"] != 1 {
		t.Errorf("failed assertions were not counted per assertion: %v", got.Failures)
	}
	if got.Rate() != 1.0/3.0 {
		t.Errorf("rate = %v, want one in three", got.Rate())
	}
}

// A sample the harness could not be run for at all is excluded rather than
// counted as a failure - the same rule that keeps an unreachable endpoint out
// of the selection denominator.
func TestAHarnessErrorIsNotABehaviorFailure(t *testing.T) {
	m := Measurement{Samples: []Sample{
		behaviorSample(&BehaviorRun{Outcome: "ok", Checks: 5, Passed: 5}),
		behaviorSample(&BehaviorRun{Err: os.ErrNotExist}),
	}}
	caseWithTarget(&m, "todo-rails")

	got, _ := m.Behavior()
	if got.Measured != 1 || got.Errored != 1 {
		t.Errorf("measured/errored = %d/%d, want 1/1", got.Measured, got.Errored)
	}
	if got.Rate() != 1 {
		t.Errorf("rate = %v; an errored run must not drag the rate down", got.Rate())
	}
}

// Not measured and measured-with-nothing-working are different states, and the
// report has to say which. A reader who cannot tell them apart reads the
// absence of a number as the absence of a problem.
func TestNotMeasuredIsNotZero(t *testing.T) {
	m := Measurement{Samples: []Sample{{Counts: map[string]int{"addColumn": 1}}}}
	caseWithTarget(&m, "todo-rails")

	if _, measured := m.Behavior(); measured {
		t.Error("a run with no behaviour samples reported as measured")
	}

	var out strings.Builder
	behavior(&out, m)
	if !strings.Contains(out.String(), "not measured") {
		t.Errorf("the report does not say behaviour went unmeasured:\n%s", out.String())
	}
	if strings.Contains(out.String(), "working:") {
		t.Errorf("the report printed a working rate for a run that measured nothing:\n%s", out.String())
	}
}

// The report names which assertions broke. A rate that says something is wrong
// without saying what is not the finding this harness exists to produce.
func TestTheReportNamesTheBrokenContracts(t *testing.T) {
	m := Measurement{Samples: []Sample{
		behaviorSample(&BehaviorRun{Outcome: "checks_failed", Checks: 19, Passed: 16,
			Failed: []string{"a todo with no title is rejected"}, Elapsed: 4 * time.Second}),
	}}
	caseWithTarget(&m, "todo-chi")

	var out strings.Builder
	behavior(&out, m)
	got := out.String()
	for _, want := range []string{"todo-chi", "disagreed", "a todo with no title is rejected", "16 of 19"} {
		if !strings.Contains(got, want) {
			t.Errorf("the behaviour table does not mention %q:\n%s", want, got)
		}
	}
}

// The answer file is the envelope Phase 4 decodes, carrying the sample's own
// bindings. A kwarg dropped or renamed on the way through would make the
// harness measure a selection nobody made.
func TestTheAnswerCarriesTheSamplesOwnBindings(t *testing.T) {
	path, err := writeAnswer(invocationsFixture())
	if err != nil {
		t.Fatalf("writeAnswer: %v", err)
	}
	defer os.Remove(path)

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the answer: %v", err)
	}

	var body struct {
		Invocations []struct {
			Action string         `json:"action"`
			Kwargs map[string]any `json:"kwargs"`
		} `json:"invocations"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("the answer is not the envelope Phase 4 decodes: %v\n%s", err, raw)
	}
	if len(body.Invocations) != 2 {
		t.Fatalf("wrote %d invocations, want 2", len(body.Invocations))
	}
	if body.Invocations[0].Action != "addColumn" {
		t.Errorf("action = %q, want addColumn", body.Invocations[0].Action)
	}
	if body.Invocations[0].Kwargs["name"] != "title" {
		t.Errorf("kwargs did not survive: %v", body.Invocations[0].Kwargs)
	}
	// The omission is the answer for a defaulted kwarg (prov-2026-f03916ba), so
	// it has to still be an omission by the time the harness replays it.
	if _, bound := body.Invocations[0].Kwargs["default"]; bound {
		t.Errorf("an unbound kwarg acquired a value on the way to the harness: %v",
			body.Invocations[0].Kwargs)
	}
}

// Nothing to apply is not a measurement. A caller is expected not to ask, and
// asking gets an error rather than a zero-check pass.
func TestAnEmptySelectionIsNotMeasured(t *testing.T) {
	got := RunBehavior(t.Context(), "todo-rails", nil)
	if got.Err == nil {
		t.Error("applying no invocations reported a result")
	}
	if got.Working() {
		t.Error("an unmeasurable sample reported as working")
	}
}

// invocationsFixture is one sample's worth of bindings, with the defaulted
// kwarg left unbound as a correct answer leaves it.
func invocationsFixture() []recording.Invocation {
	return []recording.Invocation{
		{Action: "addColumn", Kwargs: map[string]any{
			"resource": "todos", "stamp": "20260814000000",
			"name": "title", "type": "string", "nullable": false,
		}},
		{Action: "addColumn", Kwargs: map[string]any{
			"resource": "todos", "stamp": "20260814000000",
			"name": "completed", "type": "boolean", "nullable": false, "default": "false",
		}},
	}
}
