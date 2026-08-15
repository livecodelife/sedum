package evals

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

// The harness's own arithmetic is deterministic and is tested as such. Only the
// runner needs a model, and it is behind a build tag; none of that is what these
// assert. What they assert is the distinction the first version got wrong.

func measurement(samples ...Sample) Measurement {
	m := Measurement{
		Model:       Model{ID: "test-model", Engine: "mlx", Quant: "4bit"},
		Samples:     samples,
		Concurrency: 1,
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

// Two runtimes over one checkpoint are two rows, not one. MLX 4-bit and a
// llama.cpp Q4_K_M build use different quantization schemes and do not produce
// identical output, so a rate measured on one must not read as a claim about
// the other.
func TestModelLabelDistinguishesEngineAndQuant(t *testing.T) {
	mlx := Model{ID: "qwen2.5-coder-14b", Engine: "mlx", Quant: "4bit"}
	gguf := Model{ID: "qwen2.5-coder-14b", Engine: "llama.cpp", Quant: "q4_k_m"}

	if mlx.Label() == gguf.Label() {
		t.Fatalf("one checkpoint under two engines produced one label: %q", mlx.Label())
	}
	if !strings.Contains(mlx.Label(), "mlx") || !strings.Contains(gguf.Label(), "q4_k_m") {
		t.Errorf("labels do not carry engine and quant: %q / %q", mlx.Label(), gguf.Label())
	}

	// A hosted model's weights are not ours to describe, so quant is optional.
	hosted := Model{ID: "qwen/qwen3.6-27b", Engine: "groq"}
	if hosted.Label() != "qwen/qwen3.6-27b/groq" {
		t.Errorf("label is %q, want the id and engine with no trailing separator", hosted.Label())
	}
}

// A case that omits the engine is the mistake this field exists to prevent, so
// it is rejected at load rather than silently merging two rows.
func TestACaseWithoutAnEngineIsRejected(t *testing.T) {
	c := Case{ID: "x", Arm: "baseline", Models: []Model{{ID: "some-model"}}}
	err := c.validate("x.yaml")
	if err == nil {
		t.Fatal("a model with no engine was accepted")
	}
	if !strings.Contains(err.Error(), "engine") {
		t.Errorf("error does not name the missing field: %v", err)
	}
}

// The report has to be able to say what the run cost. The harness could not do
// this before, which is why the 75s-per-call figure behind prov-2026-6d87dc11
// had to be reconstructed from a stale run log.
func TestReportCarriesTiming(t *testing.T) {
	m := measurement(
		Sample{Counts: map[string]int{"addColumn": 2}, First: "addColumn", Elapsed: 10 * time.Second},
		Sample{Counts: map[string]int{"addColumn": 2}, First: "addColumn", Elapsed: 30 * time.Second},
	)
	m.Wall = 32 * time.Second

	var buf bytes.Buffer
	Report(&buf, m)
	out := buf.String()

	// Fastest and slowest both appear, because the spread is the throttling
	// signal and a mean alone would hide it.
	for _, want := range []string{"wall 32s", "fastest 10s", "mean 20s", "slowest 30s"} {
		if !strings.Contains(out, want) {
			t.Errorf("timing line omits %q:\n%s", want, out)
		}
	}
}

// Concurrency must not reorder results. A measurement whose rows moved between
// runs would make two of them harder to compare for no reason.
func TestConcurrentSamplesKeepTheirOrder(t *testing.T) {
	c := Case{ID: "fixture", Arm: "sedum", Generators: "x", Records: "y"}
	c.Expect.Actions = map[string]int{"addColumn": 2}

	// Every sample fails identically here - there is no endpoint - which is
	// enough to prove the slice is filled by index rather than appended.
	m, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"}, 6, 4)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(m.Samples) != 6 {
		t.Fatalf("got %d samples, want 6 - concurrent writes lost one", len(m.Samples))
	}
	if m.Concurrency != 4 {
		t.Errorf("concurrency recorded as %d, want 4", m.Concurrency)
	}
	if m.Wall == 0 {
		t.Error("wall time was not recorded")
	}
}

// A baseline arm has no generator package and no action vocabulary, so it is
// declared in the matrix and refused at run time rather than silently skipped.
func TestTheBaselineArmIsRefusedRatherThanSkipped(t *testing.T) {
	c := Case{ID: "fixture", Arm: "baseline"}
	if _, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"}, 1, 1); err == nil {
		t.Fatal("the baseline arm ran; it is not implemented and must say so")
	}
}

// The fixtures are vendored so a number measured today can be re-run tomorrow.
// This asserts both load and resolve to real directories - a case pointing at a
// path that moved is the failure mode that made prov-2026-6d87dc11's original
// finding unreproducible.
func TestVendoredCasesLoad(t *testing.T) {
	cases, err := Load("cases", "testdata")
	if err != nil {
		t.Fatalf("loading cases: %v", err)
	}
	if len(cases) < 2 {
		t.Fatalf("loaded %d cases, want at least the rails and chi fixtures", len(cases))
	}

	seen := map[string]bool{}
	for _, c := range cases {
		seen[c.Language] = true
		for _, dir := range []string{c.Generators, c.Records} {
			if _, err := os.Stat(dir); err != nil {
				t.Errorf("case %s points at a directory that is not there: %v", c.ID, err)
			}
		}
	}

	// Two languages, which is what makes the framework axis more than a label.
	for _, lang := range []string{"ruby", "go"} {
		if !seen[lang] {
			t.Errorf("no vendored case covers %s", lang)
		}
	}
}

// A new fixture has no expectations, because inventing them from the package
// would make the first run agree by construction. The report says what it
// observed instead, which is how they get established.
func TestACaseWithNoExpectationsReportsWhatItObserved(t *testing.T) {
	m := measurement(
		valid(map[string]int{"createEndpoint": 5, "createQuery": 5}, "createEndpoint"),
		valid(map[string]int{"createEndpoint": 4, "createQuery": 5}, "createEndpoint"),
	)
	m.Case.Expect.Actions = nil

	var buf bytes.Buffer
	Report(&buf, m)
	out := buf.String()

	if !strings.Contains(out, "no expectations declared") {
		t.Errorf("report does not say the case declares none:\n%s", out)
	}
	// 5 and 4 across two samples.
	if !strings.Contains(out, "createEndpoint") || !strings.Contains(out, "4.50") {
		t.Errorf("report does not carry the observed mean:\n%s", out)
	}
}

// Observed counts every action selected, including ones no expectation names -
// otherwise an action the model reaches for that nobody anticipated is
// invisible.
func TestObservedIncludesUnexpectedActions(t *testing.T) {
	m := measurement(valid(map[string]int{"addColumn": 2, "somethingNobodyExpected": 1}, "addColumn"))

	got := m.Observed()
	if got["somethingNobodyExpected"] != 1 {
		t.Errorf("observed is %v; an unexpected action was dropped", got)
	}
}
