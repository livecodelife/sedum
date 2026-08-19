package evals

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"maps"
	"math"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/recording"
	"github.com/calebcowen/sedum/internal/selection"
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
	m, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"},
		Options{Resolution: Smoke, Samples: 6, Concurrency: 4})
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

// The retry budget travels with the measurement and into the entry, because
// "valid" means a different thing at zero than at two and nothing in the table
// recovers which was measured (prov-2026-b4555efc).
func TestTheRetryBudgetIsCarriedNotAssumed(t *testing.T) {
	c := Case{ID: "fixture", Arm: "sedum", Generators: "x", Records: "y"}

	m, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"},
		Options{Resolution: Smoke, Samples: 1, Concurrency: 1, Retries: 2})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if m.Retries != 2 {
		t.Errorf("measurement recorded %d retries, want 2", m.Retries)
	}
	if got := NewEntry(m, "http://x").Retries; got != 2 {
		t.Errorf("entry recorded %d retries, want the budget the run was given", got)
	}

	// The default is zero and cannot drift: every entry already recorded means
	// "validated on the first call", and a new default would re-point that word
	// while the field name and the report line stayed put.
	zero, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"},
		Options{Resolution: Smoke, Samples: 1})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if zero.Retries != 0 {
		t.Errorf("an unset budget came out as %d, want 0", zero.Retries)
	}
}

// The validity line names the budget it was measured under. A run that spent
// retries did not measure first-call validity, and a line that read the same
// either way would be the entry format's one prohibited move - a field that
// changes meaning without changing name.
func TestTheValidityLineNamesItsBudget(t *testing.T) {
	m := measurement(valid(map[string]int{"addColumn": 2}, "addColumn"))

	var buf bytes.Buffer
	Report(&buf, m)
	if out := buf.String(); !strings.Contains(out, "valid first call:") || !strings.Contains(out, "retries=0") {
		t.Errorf("a zero-budget report does not say so:\n%s", out)
	}

	m.Retries = 2
	buf.Reset()
	Report(&buf, m)
	out := buf.String()
	if !strings.Contains(out, "valid within 3 calls:") {
		t.Errorf("a raised budget still claims first-call validity:\n%s", out)
	}
	if !strings.Contains(out, "retries=2") {
		t.Errorf("the header omits the budget:\n%s", out)
	}
}

// Cost is reported in calls, which is the comparison seconds could never
// support: a per-sample time is calls multiplied by an unknown and varying
// per-call cost (prov-2026-0811425c).
func TestCostIsCountedInCallsNotSeconds(t *testing.T) {
	m := measurement(
		Sample{Counts: map[string]int{"addColumn": 2}, Calls: 1},
		Sample{Counts: map[string]int{"addColumn": 2}, Calls: 3, Rejected: 1, Completeness: 1},
		Sample{Invalid: true, Calls: 3, Rejected: 3},
		// A sample that never reached the model contributes nothing: a call
		// that returned no answer is not a cost a selection can be charged for.
		Sample{Err: errors.New("connection refused"), Calls: 9},
	)

	spent := m.Spent()
	if spent.Calls != 7 {
		t.Errorf("spent %d calls, want 7 - the failed sample must not be charged", spent.Calls)
	}
	if spent.Completeness != 1 {
		t.Errorf("counted %d completeness observations, want 1", spent.Completeness)
	}
	// First-call validity, recovered from the rejection counts rather than from
	// the retry budget: one of the three answered samples took no rejection.
	if spent.FirstTry != 1 {
		t.Errorf("counted %d first-try answers, want 1", spent.FirstTry)
	}
}

// Throughput is derived from what was counted, never from a model's advertised
// rate: the advertised number describes one stream on the vendor's hardware,
// and what a run is bounded by is this server at this concurrency
// (prov-2026-945fa0aa).
func TestThroughputIsDerivedFromWhatWasCounted(t *testing.T) {
	m := measurement(
		Sample{Counts: map[string]int{"addColumn": 2}, Calls: 1, PromptTokens: 2000, CompletionTokens: 500},
		Sample{Counts: map[string]int{"addColumn": 2}, Calls: 1, PromptTokens: 3000, CompletionTokens: 500},
	)
	m.Wall = 10 * time.Second
	m.Concurrency = 2

	// Completion tokens only: 1000 over 10s. The 5000 prompt tokens are billed
	// but were almost certainly never computed, and counting them as work would
	// make a case look faster the larger its prompt is (prov-2026-e323b805).
	if got := m.Spent().PerSecond(m.Wall); math.Abs(got-100) > 0.01 {
		t.Errorf("throughput %.2f tok/s, want 100 - 1000 completion tokens over 10s", got)
	}

	var buf bytes.Buffer
	Report(&buf, m)
	out := buf.String()
	if !strings.Contains(out, "throughput: 100.0 completion tok/s at concurrency 2") {
		t.Errorf("report omits completion throughput at its concurrency:\n%s", out)
	}
	// The prompt half stays on the tokens line, labelled for what it is: a
	// hosted endpoint charges it and a rate limit counts it.
	if !strings.Contains(out, "prompt billed") {
		t.Errorf("the tokens line does not mark prompt tokens as billed:\n%s", out)
	}
}

// A server that reported no usage has no throughput to report, and zero would
// read as a measured rate rather than as its absence.
func TestThroughputIsAbsentRatherThanZero(t *testing.T) {
	m := measurement(Sample{Counts: map[string]int{"addColumn": 2}, Calls: 1})
	m.Wall = 10 * time.Second

	if got := m.Spent().PerSecond(m.Wall); got != 0 {
		t.Errorf("throughput %.2f from a run with no token counts", got)
	}

	var buf bytes.Buffer
	Report(&buf, m)
	if out := buf.String(); strings.Contains(out, "throughput:") {
		t.Errorf("a run with no usage reported a throughput line:\n%s", out)
	}
}

// Five samples do not distinguish 4/5 from 5/5, and a table that prints them in
// the same column invites a reader to think otherwise. The chi case has come in
// at 5/5, 2/5 and 4/5 across runs where the difference could not have been
// caused by the change between them (prov-2026-0baaa119).
func TestARateCarriesTheUncertaintyItsSampleSizeGivesIt(t *testing.T) {
	cases := []struct {
		successes, samples int
		low, high          float64
	}{
		// The reading most likely to be over-believed: 5/5 is not evidence of
		// a rate of 1.00, and the normal approximation would claim it was by
		// collapsing to zero width here.
		{5, 5, 0.57, 1.00},
		{4, 5, 0.38, 0.96},
		{2, 5, 0.12, 0.77},
		{0, 5, 0.00, 0.43},
	}

	for _, c := range cases {
		got := wilson(c.successes, c.samples)
		if math.Abs(got.Low-c.low) > 0.01 || math.Abs(got.High-c.high) > 0.01 {
			t.Errorf("wilson(%d,%d) = [%.2f,%.2f], want [%.2f,%.2f]",
				c.successes, c.samples, got.Low, got.High, c.low, c.high)
		}
	}

	// The fraction stays visible: the sample size is what produced the width.
	if got := wilson(4, 5).String(); !strings.HasPrefix(got, "4/5 ") {
		t.Errorf("an interval hid the fraction that produced it: %s", got)
	}
	if got := wilson(0, 0).String(); got != "-" {
		t.Errorf("an interval over no samples rendered as %q, want a dash", got)
	}
}

// Overlap is the condition under which a reported difference is not
// distinguishable from sampling. It decides nothing - the harness reports rates
// rather than verdicts - but a table that stays silent about it reports noise
// as a result.
func TestOverlappingIntervalsAreMarkedInAComparison(t *testing.T) {
	a := measurement(
		valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"),
		valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"),
	)
	// Same action, one sample short: 2/2 against 1/2, whose intervals overlap
	// heavily at this sample size.
	b := measurement(
		valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"),
		valid(map[string]int{"createEndpoint": 5}, "createEndpoint"),
	)

	var buf bytes.Buffer
	Compare(&buf, a, b)
	out := buf.String()

	if !strings.Contains(out, "~") {
		t.Errorf("a comparison of overlapping rates is not marked:\n%s", out)
	}
	if !strings.Contains(out, "do not distinguish") {
		t.Errorf("the legend does not say what the mark means:\n%s", out)
	}

	// Non-overlapping stays unmarked, or the mark would mean nothing.
	if wilson(5, 5).Overlaps(wilson(0, 5)) {
		t.Error("5/5 and 0/5 were reported as indistinguishable")
	}
}

// A call count says how many times the model was asked; tokens say what each
// asking cost. The split is what distinguishes a case that is slow because its
// catalog is long from one that is slow because it has more to emit
// (prov-2026-096a4d4b).
func TestTokensAreCountedBesideCalls(t *testing.T) {
	m := measurement(
		Sample{Counts: map[string]int{"addColumn": 2}, Calls: 1, PromptTokens: 2000, CompletionTokens: 200},
		Sample{Counts: map[string]int{"addColumn": 2}, Calls: 2, Rejected: 1, PromptTokens: 5000, CompletionTokens: 400},
	)

	spent := m.Spent()
	if spent.PromptTokens != 7000 || spent.CompletionTokens != 600 {
		t.Errorf("spent %d prompt and %d completion tokens, want 7000 and 600",
			spent.PromptTokens, spent.CompletionTokens)
	}

	var buf bytes.Buffer
	Report(&buf, m)
	if out := buf.String(); !strings.Contains(out, "tokens: 7000 prompt billed + 600 completion") {
		t.Errorf("report omits the token split:\n%s", out)
	}
}

// A server that fills no usage block reported nothing, which is not the same as
// calls that cost nothing. A line of zeroes would read as a measurement.
func TestAnEndpointReportingNoUsageGetsNoTokenLine(t *testing.T) {
	m := measurement(Sample{Counts: map[string]int{"addColumn": 2}, Calls: 1})

	var buf bytes.Buffer
	Report(&buf, m)
	if out := buf.String(); strings.Contains(out, "tokens:") {
		t.Errorf("a run with no usage reported a token line:\n%s", out)
	}
}

// The budget stops being a trade-off between two measurements: a run at two
// retries reports both what it validated within three calls and how often one
// call was enough.
func TestARaisedBudgetStillReportsFirstCallValidity(t *testing.T) {
	m := measurement(
		Sample{Counts: map[string]int{"addColumn": 2}, Calls: 1},
		Sample{Counts: map[string]int{"addColumn": 2}, Calls: 2, Rejected: 1},
	)
	m.Retries = 2

	var buf bytes.Buffer
	Report(&buf, m)
	out := buf.String()

	for _, want := range []string{"valid within 3 calls: 2/2", "valid first call: 1/2", "cost: 3 call(s)"} {
		if !strings.Contains(out, want) {
			t.Errorf("report omits %q:\n%s", want, out)
		}
	}
}

// The counts reach the file, additively. An entry written before them stays
// readable and keeps meaning what it meant, which is the one move the format
// forbids itself.
func TestAnEntryCarriesWhatEachSampleCost(t *testing.T) {
	m := measurement(Sample{Counts: map[string]int{"addColumn": 2}, Calls: 2, Rejected: 1, Completeness: 1})
	e := NewEntry(m, "http://x")

	if len(e.Runs) != 1 {
		t.Fatalf("got %d runs, want 1", len(e.Runs))
	}
	r := e.Runs[0]
	if r.Calls != 2 || r.Rejected != 1 || r.Completeness != 1 {
		t.Errorf("entry recorded calls=%d rejected=%d completeness=%d, want 2/1/1",
			r.Calls, r.Rejected, r.Completeness)
	}

	// An older entry has no counts, and reads as having none rather than as
	// having spent zero calls.
	var old Entry
	if err := json.Unmarshal([]byte(`{"schema":1,"case":"x","runs":[{"outcome":"valid","ms":10}]}`), &old); err != nil {
		t.Fatalf("an entry written before the counts no longer decodes: %v", err)
	}
	if old.Runs[0].Calls != 0 {
		t.Errorf("a pre-counts entry decoded calls as %d", old.Runs[0].Calls)
	}
}

// A baseline arm has no package, so selection, binding, anchor fill and
// idempotency are all undefined for it and behaviour is the only rung that can
// score one. A run that cannot collect the single number the arm exists to
// produce is refused before the first call rather than after paying for every
// sample (prov-2026-a4dbe65c).
func TestABaselineRunWithoutBehaviourIsRefused(t *testing.T) {
	c := Case{ID: "fixture", Arm: "baseline"}
	// A resolution is stated so the refusal that is asserted is the arm's, not
	// the one an unstated sample size would have produced first.
	_, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"},
		Options{Resolution: Smoke, Samples: 1, Concurrency: 1})
	if err == nil {
		t.Fatal("a baseline run with no behaviour produced a measurement of nothing")
	}
	if !strings.Contains(err.Error(), "behaviour") {
		t.Errorf("the refusal is %q, which does not say what is missing", err)
	}

	// Asked for and with nowhere to boot it is the same refusal from the other
	// side: there is no target to apply what the model wrote.
	_, err = Run(context.Background(), c, Model{ID: "none", Engine: "test"},
		Options{Resolution: Smoke, Samples: 1, Concurrency: 1, Behavior: true})
	if err == nil {
		t.Fatal("a baseline case with no behaviour target ran")
	}
	if !strings.Contains(err.Error(), "target") {
		t.Errorf("the refusal is %q, which does not name the missing target", err)
	}
}

// An arm nobody declared is a typo, and it names both real values rather than
// leaving the reader to find them.
func TestAnUnknownArmIsNamedWithTheRealOnes(t *testing.T) {
	c := Case{ID: "fixture", Arm: "sedumm"}
	_, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"},
		Options{Resolution: Smoke, Samples: 1, Concurrency: 1})
	if err == nil {
		t.Fatal("a misspelled arm ran")
	}
	for _, want := range []string{"sedumm", "sedum", "baseline"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the refusal %q does not name %q", err, want)
		}
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
		// Every case names a records directory, because the record is what is
		// being generated from either way. Only the sedum arm names a
		// generators one: a baseline has no package, and that absence is the
		// arm rather than an omission (prov-2026-a4dbe65c).
		dirs := []string{c.Records}
		if c.Arm == "sedum" {
			dirs = append(dirs, c.Generators)
		} else if c.Generators != "" {
			t.Errorf("case %s has arm %q and names a generators directory; the arm is the absence of one", c.ID, c.Arm)
		}
		for _, dir := range dirs {
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

// The two arms differ in exactly one thing: the stack. Their records ask for
// the same functionality constraint for constraint, which is what a language
// comparison needs and what the source projects' own records did not provide -
// those prescribed opposite solutions to deployment concerns no action could
// serve, so a delta across them measured the records rather than the stacks.
func TestTheLanguageArmsAreControlled(t *testing.T) {
	cases, err := Load("cases", "testdata")
	if err != nil {
		t.Fatalf("loading cases: %v", err)
	}

	by := map[string]Case{}
	for _, c := range cases {
		by[c.ID] = c
	}
	rails, chi := by["todo-rails-defined"], by["todo-chi-defined"]
	if rails.ID == "" || chi.ID == "" {
		t.Fatal("both language arms must exist")
	}

	if rails.Language == chi.Language {
		t.Error("the arms share a language; nothing is being varied")
	}
	for _, f := range []struct{ name, a, b string }{
		{"tightness", rails.Tightness, chi.Tightness},
		{"arm", rails.Arm, chi.Arm},
		{"application", rails.Application.Name, chi.Application.Name},
	} {
		if f.a != f.b {
			t.Errorf("the arms differ in %s (%s vs %s), which is a second variable", f.name, f.a, f.b)
		}
	}
	if rails.Application.Complexity != chi.Application.Complexity {
		t.Errorf("the arms are tier %d and tier %d; they ask for the same functionality, so rating them differently puts a package difference on the application axis",
			rails.Application.Complexity, chi.Application.Complexity)
	}
}

// The description arms differ in exactly one thing: which package set they name.
// The record is the same file, the models are the same rows, and the
// expectations are the same map - because an expectation is a property of the
// record, and two blocks that drifted apart would measure two questions and
// report one number (prov-2026-ac15ed2b).
func TestTheDescriptionArmsAreControlled(t *testing.T) {
	cases, err := Load("cases", "testdata")
	if err != nil {
		t.Fatalf("loading cases: %v", err)
	}

	by := map[string]Case{}
	for _, c := range cases {
		by[c.ID] = c
	}
	plain, described := by["todo-rails-defined"], by["todo-rails-described"]
	if plain.ID == "" || described.ID == "" {
		t.Fatal("both description arms must exist")
	}

	if plain.Generators == described.Generators {
		t.Error("the arms name one package set; nothing is being varied")
	}
	for _, f := range []struct{ name, a, b string }{
		{"records", plain.Records, described.Records},
		{"language", plain.Language, described.Language},
		{"framework", plain.Framework, described.Framework},
		{"tightness", plain.Tightness, described.Tightness},
		{"arm", plain.Arm, described.Arm},
		{"application", plain.Application.Name, described.Application.Name},
	} {
		if f.a != f.b {
			t.Errorf("the arms differ in %s (%s vs %s), which is a second variable", f.name, f.a, f.b)
		}
	}
	if plain.Application.Complexity != described.Application.Complexity {
		t.Errorf("the arms are tier %d and tier %d over one application",
			plain.Application.Complexity, described.Application.Complexity)
	}
	if !reflect.DeepEqual(plain.Models, described.Models) {
		t.Errorf("the arms name different models:\n  %v\n  %v", plain.Models, described.Models)
	}
	if !reflect.DeepEqual(plain.Expect.Actions, described.Expect.Actions) {
		t.Errorf("the arms expect different actions:\n  %v\n  %v", plain.Expect.Actions, described.Expect.Actions)
	}
	// Bindings are a property of the record too, and the arms answer one
	// record. Two expectations that drifted apart would grade the same answer
	// differently depending on which package produced it, which is the one
	// thing a controlled comparison cannot survive (prov-2026-2b121b62).
	if !reflect.DeepEqual(plain.Expect.Bindings, described.Expect.Bindings) {
		t.Errorf("the arms expect different bindings:\n  %+v\n  %+v", plain.Expect.Bindings, described.Expect.Bindings)
	}
	if len(plain.Expect.Bindings) == 0 {
		t.Error("neither arm expects any binding; the comparison measures selection only")
	}
	// A signal is a column of the matrix rather than a property of one arm, the
	// same argument the model list is held to. A check command present in one
	// arm and absent in the other would make the pair incomparable on
	// syntactic validity the moment anyone ran it unfiltered
	// (prov-2026-d61010a4).
	if !reflect.DeepEqual(plain.Check, described.Check) {
		t.Errorf("the arms run different syntax checks:\n  %v\n  %v", plain.Check, described.Check)
	}
	if len(plain.Check) == 0 {
		t.Error("neither arm checks what it wrote; the pair reports no syntactic validity")
	}
}

// The described set is the defined set with descriptions added, and this is what
// holds that true. A hand-written second package drifts, and a reworded kwarg or
// a dropped required flag would make the A/B measure something nobody chose -
// so the guard is a diff in the default suite rather than a rule about who may
// edit the directory (prov-2026-ac15ed2b).
//
// It compares loaded packages rather than file bytes, because "differs only by
// descriptions" is a statement about what reaches the catalog. Comments are the
// prose the model cannot see, and both sets keep theirs; a textual diff would
// fail on that and say nothing about the experiment.
func TestTheDescribedSetDiffersOnlyByDescriptions(t *testing.T) {
	plain := loadPackageSet(t, "testdata/todo-rails/generators/defined")
	described := loadPackageSet(t, "testdata/todo-rails/generators/described")

	if len(plain) != len(described) {
		t.Fatalf("the sets hold %d and %d packages", len(plain), len(described))
	}

	var descriptions int
	for name, a := range plain {
		b, ok := described[name]
		if !ok {
			t.Errorf("the described set has no package %s", name)
			continue
		}

		for _, f := range []struct {
			name string
			a, b any
		}{
			{"extensions", a.Extensions, b.Extensions},
			{"comment_prefix", a.CommentPrefix, b.CommentPrefix},
			{"transforms", a.Transforms, b.Transforms},
			{"op_exceptions", a.OpExceptions, b.OpExceptions},
			{"unmanaged", a.Unmanaged, b.Unmanaged},
			{"file templates", a.FileTemplates, b.FileTemplates},
		} {
			if !reflect.DeepEqual(f.a, f.b) {
				t.Errorf("%s: the sets declare different %s:\n  %v\n  %v", name, f.name, f.a, f.b)
			}
		}

		for _, pattern := range a.FileTemplates {
			x, _ := a.FileTemplate(pattern)
			y, ok := b.FileTemplate(pattern)
			if !ok {
				t.Errorf("%s: the described set has no file template %s", name, pattern)
				continue
			}
			if x != y {
				t.Errorf("%s: file template %s differs between the sets", name, pattern)
			}
		}

		if len(a.Actions) != len(b.Actions) {
			t.Errorf("%s: the sets declare %d and %d actions", name, len(a.Actions), len(b.Actions))
		}
		for action, x := range a.Actions {
			y, ok := b.Actions[action]
			if !ok {
				t.Errorf("%s: the described set has no action %s", name, action)
				continue
			}

			// Every field is compared, descriptions excluded, so an action
			// field added later is covered without this test being revisited.
			if !reflect.DeepEqual(undescribed(x), undescribed(y)) {
				t.Errorf("%s: action %s differs between the sets by more than a description:\n  %+v\n  %+v",
					name, action, undescribed(x), undescribed(y))
			}

			for _, path := range templatePaths(x) {
				p, _ := a.ActionTemplate(path)
				q, ok := b.ActionTemplate(path)
				if !ok {
					t.Errorf("%s: the described set has no template %s", name, path)
					continue
				}
				if p != q {
					t.Errorf("%s: template %s differs between the sets", name, path)
				}
			}

			if x.Description != "" {
				t.Errorf("%s: action %s carries a description in the undescribed arm", name, action)
			}
			if y.Description != "" {
				descriptions++
			}
			for kwarg, k := range x.Kwargs {
				if k.Description != "" {
					t.Errorf("%s: %s.%s carries a description in the undescribed arm", name, action, kwarg)
				}
			}
			for _, k := range y.Kwargs {
				if k.Description != "" {
					descriptions++
				}
			}
		}
	}

	// Two identical sets would satisfy every assertion above, and would be an
	// A/B that varies nothing while reading like one that does.
	if descriptions == 0 {
		t.Error("the described set carries no descriptions; the arms are identical and the comparison measures nothing")
	}
}

// loadPackageSet loads a generator package set, keyed by package name. A set
// that does not load clean is not an arm - it is a broken fixture that would
// hand every sample a different catalog.
func loadPackageSet(t *testing.T, dir string) map[string]*genpkg.Package {
	t.Helper()

	set, findings, err := genpkg.Load(dir, genpkg.Options{})
	if err != nil {
		t.Fatalf("loading %s: %v", dir, err)
	}
	for _, f := range findings {
		t.Errorf("%s: %s", dir, f)
	}

	out := map[string]*genpkg.Package{}
	for _, p := range set.Packages {
		out[p.Name] = p
	}
	if len(out) == 0 {
		t.Fatalf("%s holds no packages", dir)
	}
	return out
}

// undescribed is an action with every description blanked, which is the form the
// two arms must be identical in.
func undescribed(a *genpkg.Action) genpkg.Action {
	out := *a
	out.Description = ""
	out.Kwargs = maps.Clone(a.Kwargs)
	for name, k := range out.Kwargs {
		k.Description = ""
		out.Kwargs[name] = k
	}
	return out
}

// templatePaths is every template an action renders from, simple or
// discriminated. A composite has none of its own.
func templatePaths(a *genpkg.Action) []string {
	if a.Template != "" {
		return []string{a.Template}
	}
	return slices.Sorted(maps.Values(a.Templates))
}

// The sample size follows from the question. Five was never chosen against one -
// it was the first run's number and every run since inherited it - and what it
// buys is an interval a fine question cannot be asked inside of
// (prov-2026-3039750e).
func TestTheResolutionDeterminesTheSampleSize(t *testing.T) {
	for _, tc := range []struct {
		res  Resolution
		want int
	}{
		{Smoke, 2},
		{Coarse, 5},
		{Fine, 30},
	} {
		if got := tc.res.Samples(); got != tc.want {
			t.Errorf("%s draws %d samples, want %d", tc.res, got, tc.want)
		}
	}

	// Five remains the default, but as the coarse tier's own number rather than
	// as an inheritance: the questions widening the matrix asks are coarse ones.
	if DefaultResolution != Coarse || DefaultResolution.Samples() != 5 {
		t.Errorf("the default is %s at n=%d, want coarse at 5", DefaultResolution, DefaultResolution.Samples())
	}

	// A run that states no question is a sample size nobody chose, which is the
	// whole finding this record was written from.
	c := Case{ID: "fixture", Arm: "sedum", Generators: "x", Records: "y"}
	if _, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"}, Options{Samples: 5}); err == nil {
		t.Error("a run with no stated resolution was accepted")
	}
}

// A null result from a run too small to have seen the effect is not a null
// result. The refusal is at the top of Run rather than in the reading of the
// report, because by then the hours are spent.
func TestARunTooSmallForItsQuestionIsRefused(t *testing.T) {
	c := Case{ID: "fixture", Arm: "sedum", Generators: "x", Records: "y"}

	_, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"},
		Options{Resolution: Fine, Samples: 5})
	if err == nil {
		t.Fatal("a fine question drawn at five samples was accepted")
	}
	if !strings.Contains(err.Error(), "fine") {
		t.Errorf("the refusal is %q, which does not name the resolution it was drawn against", err)
	}

	// Oversampling is a cost decision and misleads nobody, so it is honoured
	// rather than clamped.
	over, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"},
		Options{Resolution: Coarse, Samples: 7})
	if err != nil {
		t.Fatalf("oversampling a coarse question was refused: %v", err)
	}
	if len(over.Samples) != 7 {
		t.Errorf("drew %d samples, want the 7 that were asked for", len(over.Samples))
	}

	// Smoke has no measurement to mislabel, so one sample is a perfectly good
	// plumbing check.
	if _, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"},
		Options{Resolution: Smoke, Samples: 1}); err != nil {
		t.Errorf("a single-sample smoke run was refused: %v", err)
	}
}

// A smoke rate is not a measurement and says so wherever it appears: on its own
// report, in a comparison, and in the history file long after whoever ran it has
// forgotten which it was.
func TestASmokeRateSaysItIsNotAMeasurement(t *testing.T) {
	m := measurement(valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"))
	m.Resolution = Smoke

	var buf bytes.Buffer
	Report(&buf, m)
	out := buf.String()
	if !strings.Contains(out, "res=smoke") {
		t.Errorf("the header does not carry the resolution:\n%s", out)
	}
	if !strings.Contains(out, "SMOKE:") || !strings.Contains(out, "Not a measurement") {
		t.Errorf("a smoke report reads as a measurement:\n%s", out)
	}

	// And a comparison is the reading it must survive, because two plumbing
	// checks differ by whatever the model did that morning.
	coarse := measurement(valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"))
	coarse.Resolution = Coarse
	buf.Reset()
	Compare(&buf, m, coarse)
	if out := buf.String(); !strings.Contains(out, "SMOKE:") {
		t.Errorf("a comparison against a smoke arm does not say so:\n%s", out)
	}

	// The header names the resolution even when it is coarse, so that a rate is
	// never read without the question it was drawn for.
	buf.Reset()
	Report(&buf, coarse)
	if out := buf.String(); !strings.Contains(out, "res=coarse") {
		t.Errorf("a coarse report omits its resolution:\n%s", out)
	}
}

// The resolution travels into the entry, because n alone cannot recover it: two
// samples drawn as a plumbing check and two survivors of a five-sample run print
// the same number and are not the same claim.
func TestTheResolutionIsRecordedAndReadBack(t *testing.T) {
	m := measurement(valid(map[string]int{"addColumn": 2}, "addColumn"))
	m.Resolution = Fine

	e := NewEntry(m, "http://x")
	if e.Resolution != "fine" {
		t.Errorf("entry recorded resolution %q, want fine", e.Resolution)
	}

	raw, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("marshalling: %v", err)
	}
	var back Entry
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshalling: %v", err)
	}
	if back.Resolution != "fine" {
		t.Errorf("resolution read back as %q", back.Resolution)
	}

	// The entries already in results/ have none, and are not retroactively
	// coarse: they were drawn at five because five was the default, and
	// stamping a decision on them now would be inventing one. The field is
	// additive, so an older entry keeps meaning what it meant.
	var old Entry
	if err := json.Unmarshal([]byte(`{"schema":1,"case":"x","samples":5,"valid":5}`), &old); err != nil {
		t.Fatalf("an entry written before resolutions no longer decodes: %v", err)
	}
	if old.Resolution != "" {
		t.Errorf("a pre-resolution entry decoded as %q, want no resolution at all", old.Resolution)
	}
	if resolutionOf(Resolution(old.Resolution)) != "unstated" {
		t.Errorf("a pre-resolution entry reports as %q, want unstated", resolutionOf(Resolution(old.Resolution)))
	}
}

// History marks a smoke entry and refuses to compare against it. A two-sample
// plumbing check overlaps everything, so leaving it in the comparison chain
// would print "these runs do not distinguish each other" about a run that was
// never a measurement.
func TestHistoryDoesNotCompareAgainstASmokeEntry(t *testing.T) {
	at := time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Case: "x", At: at, Commit: "aaaaaaa", Clean: true, Resolution: "coarse", Samples: 5, Invalid: 5},
		{Case: "x", At: at, Commit: "bbbbbbb", Clean: true, Resolution: "smoke", Samples: 2, Valid: 2},
		{Case: "x", At: at, Commit: "ccccccc", Clean: true, Resolution: "coarse", Samples: 5, Valid: 5},
	}

	var buf bytes.Buffer
	History(&buf, "x", entries)
	out := buf.String()

	if !strings.Contains(out, "res") || !strings.Contains(out, "smoke") {
		t.Errorf("the history does not carry resolutions:\n%s", out)
	}
	if !strings.Contains(out, "s smoke: plumbing only") {
		t.Errorf("the smoke row is not marked:\n%s", out)
	}

	// 0/5 and 5/5 are the two coarse entries, and their intervals do not
	// overlap - so no row may be marked ~. If the smoke entry were left in the
	// chain, the 5/5 row would have been compared against 2/2 [0.34,1.00] and
	// marked indistinguishable from it.
	if strings.Contains(out, "~ interval overlaps") {
		t.Errorf("a run was marked indistinguishable from a plumbing check:\n%s", out)
	}
}

// An entry that cannot be re-run from its commit, or that never stated the
// question it was asking, gets what a smoke entry gets: printed in full, marked,
// and left out of the comparison. Every entry in results/ was both until this
// held, and each one was being compared against the last (prov-2026-c5ad54ff).
func TestHistoryComparesOnlyEntriesThatCanBeStoodOn(t *testing.T) {
	at := time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Case: "x", At: at, Commit: "aaaaaaa", Clean: true, Resolution: "coarse", Samples: 5, Invalid: 5},
		{Case: "x", At: at, Commit: "bbbbbbb", Clean: false, Resolution: "coarse", Samples: 5, Valid: 5},
		{Case: "x", At: at, Commit: "ccccccc", Clean: true, Resolution: "", Samples: 5, Valid: 5},
		{Case: "x", At: at, Commit: "ddddddd", Clean: true, Resolution: "coarse", Samples: 5, Valid: 5},
	}

	var buf bytes.Buffer
	History(&buf, "x", entries)
	out := buf.String()

	for commit, want := range map[string]string{
		"aaaaaaa": " ",
		"bbbbbbb": "*",
		"ccccccc": "?",
		"ddddddd": " ",
	} {
		if got := markOf(out, commit); got != want {
			t.Errorf("the %s row is marked %q, want %q:\n%s", commit, got, want, out)
		}
	}

	for _, legend := range []string{
		"* tree was dirty",
		"? no resolution stated",
	} {
		if !strings.Contains(out, legend) {
			t.Errorf("the %q legend is missing:\n%s", legend, out)
		}
	}

	// The chain skips what it does not admit rather than resetting at it, so the
	// last row is compared against 0/5 [0.00,0.43] - the last entry that could
	// be stood on - and those do not overlap. Both entries between them are 5/5
	// [0.57,1.00], which does overlap, so a ~ here means one of them was
	// silently serving as the baseline.
	if strings.Contains(out, "~ interval overlaps") {
		t.Errorf("a run was compared against an entry that cannot be stood on:\n%s", out)
	}

	// The mark withholds the comparison, not the measurement. Both excluded rows
	// still carry their own rate and interval.
	if strings.Count(out, "5/5 [0.57,1.00]") < 3 {
		t.Errorf("an excluded entry lost its numbers:\n%s", out)
	}
}

// markOf returns the row-prefix character History printed for the entry at the
// given commit.
func markOf(out, commit string) string {
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, commit) {
			return line[:1]
		}
	}
	return ""
}

// Eleven rejections that all render the same line are eleven rejections nobody
// can compare. The slug is the field that can be counted, and it survives into
// the entry beside the prose rather than instead of it (prov-2026-2256e6fa).
func TestARejectedSampleRecordsWhichRulesRejectedIt(t *testing.T) {
	rejection := &selection.Rejection{
		RecordID: "eval-todo-rails",
		Retries:  1,
		Attempts: []selection.Attempt{
			{Number: 1, Violations: []selection.Violation{
				{Index: 1, Action: "addColumn", Rule: "missing_kwarg", Detail: "stamp"},
				{Index: 3, Action: "addRoute", Rule: "unknown_action", Detail: "no such action"},
			}},
			// A model making one mistake twice and two mistakes once are
			// different problems, so a repeat is kept rather than deduplicated.
			{Number: 2, Violations: []selection.Violation{
				{Index: 1, Action: "addColumn", Rule: "missing_kwarg", Detail: "stamp"},
			}},
		},
	}

	got := rulesOf(rejection)
	want := []string{"missing_kwarg", "unknown_action", "missing_kwarg"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("recorded %v, want %v in attempt order with repeats kept", got, want)
	}

	// The rendered prose is what every sample shares and is why the arms were
	// indistinguishable; it stays, because a rate still has to be able to say
	// what its misses looked like.
	e := NewEntry(Measurement{
		Case:    Case{ID: "x"},
		Model:   Model{ID: "m", Engine: "test"},
		Samples: []Sample{{Invalid: true, Detail: "did not validate", Rules: got}},
	}, "http://endpoint")

	if len(e.Runs) != 1 || e.Runs[0].Outcome != "invalid" {
		t.Fatalf("the sample was stored as %+v", e.Runs)
	}
	if !reflect.DeepEqual(e.Runs[0].Rules, want) {
		t.Errorf("the entry stored rules %v, want %v", e.Runs[0].Rules, want)
	}
}

// A slug says a sample was rejected under invocation_shape; the text says the
// entry was `{"action":...,"reason":...}`. The first is countable and the second
// is the only one that says what the answer looked like, so both are kept and
// neither is scored (prov-2026-986ac4ca).
func TestARejectedSampleKeepsTheTextThatSaysWhy(t *testing.T) {
	rejection := &selection.Rejection{
		RecordID: "eval-todo-rails",
		Retries:  1,
		Attempts: []selection.Attempt{
			{Number: 1, Violations: []selection.Violation{
				{Index: 1, Action: "addColumn", Rule: "missing_kwarg", Detail: `addColumn is missing "stamp"`},
				{Index: 3, Rule: "invocation_shape", Detail: `invocation 3 is not readable: it was {"act":"x"}`},
			}},
			{Number: 2, Violations: []selection.Violation{
				{Index: 1, Action: "addColumn", Rule: "missing_kwarg", Detail: `addColumn is missing "stamp"`},
			}},
		},
	}

	want := []string{
		`addColumn is missing "stamp"`,
		`invocation 3 is not readable: it was {"act":"x"}`,
		`addColumn is missing "stamp"`,
	}
	got := violationsOf(rejection)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("collected %v, want %v in attempt order with repeats kept", got, want)
	}

	// One entry per slug, in the same order, so the two lists read together.
	if rules := rulesOf(rejection); len(rules) != len(got) {
		t.Errorf("%d rules against %d violations; they are meant to be read side by side", len(rules), len(got))
	}

	e := NewEntry(Measurement{
		Case:    Case{ID: "x"},
		Model:   Model{ID: "m", Engine: "test"},
		Samples: []Sample{{Invalid: true, Detail: "did not validate", Rules: rulesOf(rejection), Violations: got}},
	}, "http://endpoint")

	if len(e.Runs) != 1 {
		t.Fatalf("the sample was stored as %+v", e.Runs)
	}
	if !reflect.DeepEqual(e.Runs[0].Violations, want) {
		t.Errorf("the entry stored violations %v, want %v", e.Runs[0].Violations, want)
	}
}

// An entry written before the field carries none, and that reads as "not
// recorded" rather than as "nothing was wrong".
func TestAnEntryWithoutViolationsReadsAsAbsent(t *testing.T) {
	e := NewEntry(Measurement{
		Case:    Case{ID: "x"},
		Model:   Model{ID: "m", Engine: "test"},
		Samples: []Sample{{Invalid: true, Detail: "did not validate", Rules: []string{"missing_kwarg"}}},
	}, "http://endpoint")

	if e.Runs[0].Violations != nil {
		t.Errorf("stored %v, want nil so the field is omitted entirely", e.Runs[0].Violations)
	}
}

// Counts are a projection of the invocations and were the only thing kept. A
// wrong argument on a correctly selected action - which is every failure kwarg
// descriptions were added for - is invisible in the projection and plain in the
// original (prov-2026-2256e6fa).
func TestAValidSampleRecordsWhatItBound(t *testing.T) {
	invocations := []recording.Invocation{
		{Action: "addColumn", Kwargs: map[string]any{"resource": "todo", "name": "title", "type": "string"}},
		{Action: "addColumn", Kwargs: map[string]any{"resource": "todo", "name": "completed", "type": "boolean"}},
	}

	e := NewEntry(Measurement{
		Case:    Case{ID: "x"},
		Model:   Model{ID: "m", Engine: "test"},
		Samples: []Sample{{Counts: map[string]int{"addColumn": 2}, First: "addColumn", Invocations: invocations}},
	}, "http://endpoint")

	if !reflect.DeepEqual(e.Runs[0].Invocations, invocations) {
		t.Fatalf("stored %+v, want the bound arguments", e.Runs[0].Invocations)
	}

	// Through JSON, because the value of storing this is that a rule invented
	// later can be run against a sample drawn before it existed.
	raw, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("marshalling: %v", err)
	}
	var back Entry
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshalling: %v", err)
	}
	if got := back.Runs[0].Invocations[1].Kwargs["name"]; got != "completed" {
		t.Errorf("a bound argument did not survive the round trip: %v", got)
	}

	// Counts are unchanged by this. Every expectation and every stored entry is
	// written in them, and the invocations are kept beside rather than instead.
	if back.Runs[0].Counts["addColumn"] != 2 {
		t.Errorf("the counts changed: %v", back.Runs[0].Counts)
	}
}

// The fields are additive, so an entry written before them decodes and keeps
// meaning what it meant - the same promise the retry and token fields were
// added under (prov-2026-eb283c56).
func TestAnEntryWrittenBeforeTheseFieldsStillDecodes(t *testing.T) {
	var old Entry
	raw := `{"schema":1,"case":"x","samples":5,"valid":5,"runs":[{"outcome":"valid","counts":{"addColumn":2}}]}`
	if err := json.Unmarshal([]byte(raw), &old); err != nil {
		t.Fatalf("an entry written before these fields no longer decodes: %v", err)
	}
	if old.Runs[0].Rules != nil || old.Runs[0].Invocations != nil {
		t.Errorf("an older entry invented data it does not carry: %+v", old.Runs[0])
	}
	if old.Runs[0].Counts["addColumn"] != 2 {
		t.Errorf("an older entry lost what it did carry: %v", old.Runs[0].Counts)
	}
}

// The shipped cases parse their declared sample cost, and the chi cases carry
// one because the harness's constant is a rails number they run three times
// (prov-2026-59ed14d5).
func TestTheChiCasesDeclareWhatTheirSamplesCost(t *testing.T) {
	cases, err := Load("cases", "testdata")
	if err != nil {
		t.Fatalf("loading cases: %v", err)
	}
	by := map[string]Case{}
	for _, c := range cases {
		by[c.ID] = c
	}

	for _, id := range []string{"todo-chi-defined", "grocery-chi-defined"} {
		c, ok := by[id]
		if !ok {
			t.Fatalf("%s did not load", id)
		}
		if c.PerSample < 2*time.Minute {
			t.Errorf("%s declares per_sample %s; the constant it would otherwise get is %s, and its samples run several times that",
				id, c.PerSample, 90*time.Second)
		}
	}

	// The rails cases are what the constant was measured on, so they declare
	// nothing and the default is the right answer for them.
	for _, id := range []string{"todo-rails-defined", "todo-rails-described"} {
		if c := by[id]; c.PerSample != 0 {
			t.Errorf("%s declares per_sample %s; the harness constant was calibrated on it", id, c.PerSample)
		}
	}
}

// The baseline arm is the column todo-rails-defined is read against, so the two
// declare the same model rows. A row present in one and absent from the other
// makes the pair incomparable the moment anyone runs it unfiltered, which is
// the rule the description A/B is already held to (prov-2026-de71f29b).
func TestTheBaselineAndSedumArmsDeclareTheSameRows(t *testing.T) {
	cases, err := Load("cases", "testdata")
	if err != nil {
		t.Fatalf("loading cases: %v", err)
	}
	by := map[string]Case{}
	for _, c := range cases {
		by[c.ID] = c
	}

	sedum, base := by["todo-rails-defined"], by["todo-rails-baseline"]
	if sedum.ID == "" || base.ID == "" {
		t.Fatal("both arms of the baseline comparison must exist")
	}
	if !reflect.DeepEqual(sedum.Models, base.Models) {
		t.Errorf("the arms name different models:\n  sedum:    %v\n  baseline: %v", sedum.Models, base.Models)
	}
	// One record, because the record is the baseline's whole prompt and the
	// sedum arm is measured on the same one.
	if sedum.Records != base.Records {
		t.Errorf("the arms are measured on different records: %s and %s", sedum.Records, base.Records)
	}
	if base.Arm != "baseline" || sedum.Arm != "sedum" {
		t.Errorf("arms are %q and %q", sedum.Arm, base.Arm)
	}
}
