package evals

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// The harness's own arithmetic is deterministic and is tested as such. Only the
// runner needs a model, and it is behind a build tag; none of that is what these
// assert. What they assert is the distinction the first version got wrong.

func measurement(samples ...Sample) Measurement {
	m := Measurement{
		Model:   "test-model",
		Samples: samples,
	}
	m.Case.ID = "fixture"
	m.Case.Arm = "sedum"
	m.Case.Expect.Actions = map[string]int{"addColumn": 2, "createEndpoint": 5}
	return m
}

func valid(counts map[string]int, first string) Sample {
	return Sample{Counts: counts, First: first}
}

// An answer Phase 5 rejected is a measurement, not a failure. The first version
// of this harness counted the two alike and reported four of six samples
// "failed" without saying why - which is how a model that answers invalidly half
// the time would have looked identical to an endpoint that was down.
func TestInvalidAndFailedAreDifferentOutcomes(t *testing.T) {
	m := measurement(
		valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"),
		Sample{Invalid: true, Detail: "did not validate"},
		Sample{Err: errors.New("connection refused"), Detail: "connection refused"},
	)

	got := m.Tally()
	if got.Valid != 1 || got.Invalid != 1 || got.Failed != 1 {
		t.Errorf("tally is %+v, want one of each", got)
	}

	// A failed run is outside every denominator; an invalid one is inside the
	// one that asks how often a first call is acceptable.
	if got.Answered() != 2 {
		t.Errorf("answered is %d, want 2 - a run that never reached the model did not answer", got.Answered())
	}
}

// Only valid samples have an invocation list to count, so an invalid one must
// not be scored as though the model selected nothing. Scoring it would report a
// rejected answer as a complete miss on every action.
func TestOnlyValidSamplesAreScored(t *testing.T) {
	m := measurement(
		valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"),
		valid(map[string]int{"createEndpoint": 5}, "createEndpoint"),
		Sample{Invalid: true, Detail: "did not validate"},
	)

	rows, scored, _ := m.Scored()
	if scored != 2 {
		t.Fatalf("scored %d samples, want 2", scored)
	}

	by := map[string]ActionResult{}
	for _, r := range rows {
		by[r.Action] = r
	}

	if got := by["addColumn"]; got.Selected != 1 || got.Exact != 1 {
		t.Errorf("addColumn selected %d/%d exact %d, want 1 and 1", got.Selected, scored, got.Exact)
	}
	if got := by["createEndpoint"]; got.Selected != 2 || got.Exact != 2 {
		t.Errorf("createEndpoint selected %d exact %d, want 2 and 2", got.Selected, got.Exact)
	}
}

// Selected and exact answer different questions. An action present but short is
// a different failure from one absent entirely, and a single rate would hide
// which of the two happened.
func TestSelectedAndExactAreDistinct(t *testing.T) {
	m := measurement(
		valid(map[string]int{"addColumn": 1}, "addColumn"), // present, short
		valid(map[string]int{"addColumn": 2}, "addColumn"), // complete
	)

	rows, scored, _ := m.Scored()
	if scored != 2 {
		t.Fatalf("scored %d, want 2", scored)
	}
	for _, r := range rows {
		if r.Action != "addColumn" {
			continue
		}
		if r.Selected != 2 {
			t.Errorf("selected %d, want 2 - both samples included it", r.Selected)
		}
		if r.Exact != 1 {
			t.Errorf("exact %d, want 1 - only one sample had both invocations", r.Exact)
		}
		if r.Mean != 1.5 {
			t.Errorf("mean %v, want 1.5", r.Mean)
		}
	}
}

// The positional fact that turned out to be the whole story on todo-rails: the
// dropped action appeared if and only if it appeared first. A per-sample field
// is what makes that visible, because an aggregate cannot show it.
func TestFirstRateIsTracked(t *testing.T) {
	m := measurement(
		valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"),
		valid(map[string]int{"createEndpoint": 5}, "createEndpoint"),
		valid(map[string]int{"createEndpoint": 5}, "createEndpoint"),
	)

	rows, _, _ := m.Scored()
	for _, r := range rows {
		if r.Action == "addColumn" && r.FirstRate != 1.0/3.0 {
			t.Errorf("addColumn first rate %v, want 1/3", r.FirstRate)
		}
	}
}

// A rate reported without its sample size cannot be re-run, and will be believed
// longer than it is true. That is the failure prov-2026-6d87dc11 had to be
// corrected for, so the header carries the numbers that let someone reproduce it.
func TestReportCarriesWhatMakesItReproducible(t *testing.T) {
	m := measurement(
		valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"),
		Sample{Invalid: true, Detail: "record X: the model's output did not validate"},
	)

	var buf bytes.Buffer
	Report(&buf, m)
	out := buf.String()

	for _, want := range []string{"test-model", "n=2", "valid first call: 1/2", "did not validate"} {
		if !strings.Contains(out, want) {
			t.Errorf("report omits %q:\n%s", want, out)
		}
	}
}

// Details is what keeps a rate honest: whatever the misses were, the report can
// say so rather than leaving a number with no account of itself.
func TestDetailsAreDeduplicated(t *testing.T) {
	m := measurement(
		Sample{Invalid: true, Detail: "same reason"},
		Sample{Invalid: true, Detail: "same reason"},
		Sample{Err: errors.New("x"), Detail: "different reason"},
	)

	got := m.Details()
	if len(got) != 2 {
		t.Errorf("details are %v, want two distinct reasons", got)
	}
}
