package evals

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math"
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
	m, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"}, Options{Samples: 6, Concurrency: 4})
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
		Options{Samples: 1, Concurrency: 1, Retries: 2})
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
	zero, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"}, Options{Samples: 1})
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
	if out := buf.String(); !strings.Contains(out, "tokens: 7000 prompt + 600 completion") {
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

// A baseline arm has no generator package and no action vocabulary, so it is
// declared in the matrix and refused at run time rather than silently skipped.
func TestTheBaselineArmIsRefusedRatherThanSkipped(t *testing.T) {
	c := Case{ID: "fixture", Arm: "baseline"}
	if _, err := Run(context.Background(), c, Model{ID: "none", Engine: "test"}, Options{Samples: 1, Concurrency: 1}); err == nil {
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
