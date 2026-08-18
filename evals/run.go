package evals

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/calebcowen/sedum/internal/pipeline"
	"github.com/calebcowen/sedum/internal/recording"
	"github.com/calebcowen/sedum/internal/selection"
)

// Sample is one run of one case against one model.
type Sample struct {
	// Counts is how many invocations of each action the model selected.
	Counts map[string]int
	// Total is every invocation, including actions no expectation names.
	Total int
	// First is the action the model selected first.
	//
	// Recorded because it turned out to be the whole story once: across
	// nineteen runs of the todo-rails case, the dropped action appeared if and
	// only if it appeared first, which says the model was not weighing it and
	// declining - it was not arriving there at all. A per-sample field costs
	// nothing and is the kind of fact that is invisible in an aggregate.
	First string

	// Invalid marks a run where the model answered and Phase 5 rejected it.
	//
	// This is a measurement, not a failure, and the harness got it wrong the
	// first time by treating the two alike. At the default budget of zero
	// retries, an answer that would have been repaired by the retry loop shows
	// up as no answer at all - and how often a first call validates is one of
	// the more interesting things to know about a model, not noise to drop. At
	// a raised budget this means the answer was rejected every time it was
	// asked, which is a different and weaker claim; Measurement.Retries is what
	// says which was measured.
	Invalid bool

	// Detail is why a sample was invalid or failed, kept so a rate is never
	// reported without the ability to say what the misses looked like.
	Detail string

	// Rules are the rule slugs this sample was rejected under, one entry per
	// violation, in attempt order and with repeats kept.
	//
	// Detail is the rendered first line of the same thing and is unusable in
	// aggregate: the description A/B produced 7 rejections in one arm and 11 in
	// the other, and every one of them read "the model's output did not
	// validate within 1 attempt(s)". A slug can be counted, which is the
	// difference between knowing the model failed more often and knowing it
	// failed differently (prov-2026-2256e6fa).
	//
	// Repeats are kept because a model making one mistake three times and three
	// mistakes once are different problems, which is the same reason
	// Rejection.Attempts keeps every attempt whole (prov-2026-0811425c).
	Rules []string

	// Invocations is what the model bound, action by action, for a sample that
	// produced an answer.
	//
	// Counts are a projection of this and were until now the only thing kept.
	// Every failure that motivated kwarg descriptions was a correctly selected
	// action with a wrong argument (prov-2026-c5697387), so a count cannot see
	// the thing they were added for, while the arguments were in memory one
	// line before being discarded.
	//
	// Storing them makes an entry re-scorable: a rule invented later can be run
	// against samples drawn before it existed, instead of re-drawing two arms
	// at eighty minutes a pair every time a question is sharpened.
	Invocations []recording.Invocation

	// Calls, Rejected, Completeness and the token counts are what the sample
	// cost, summed over the case's records.
	//
	// Cost in calls rather than only in seconds is what makes two runs at
	// different retry budgets comparable, and Rejected is what makes first-call
	// validity survive a raised budget: a sample with none validated first try
	// whatever the budget allowed (prov-2026-0811425c). They are recorded for
	// invalid samples too, which is where the count is most interesting.
	Calls        int
	Rejected     int
	Completeness int

	// PromptTokens and CompletionTokens are what those calls cost, as the
	// server accounted for it. Zero means the endpoint reported no usage, which
	// is why a report prints them only when there are some: a zero averaged in
	// would read as a measurement (prov-2026-096a4d4b).
	PromptTokens     int
	CompletionTokens int

	// Err is set when the run could not be made at all - an endpoint that was
	// down, a package that would not load. Excluded from every rate, because an
	// unreachable server is not a model that chose badly.
	Err error

	// Behavior is what applying this sample's selection produced, or nil when
	// behaviour was not measured - because the run did not ask, because the
	// case declares no target, or because this sample never produced a
	// selection to apply.
	//
	// A pointer rather than a zero value, so that "not measured" and "measured
	// and failed" are different states. Folding them together would report a
	// run that skipped behaviour as one where nothing worked.
	Behavior *BehaviorRun

	// Fill is how much of the work this sample's own files declared it
	// accounted for. Zero-valued for a sample that never reached Phase 3, and
	// its Rate reports absent rather than zero when nothing was planted.
	Fill AnchorFill

	// Syntax is what the target's parser made of what this sample wrote, or
	// the zero value when the case declares no check command.
	Syntax SyntaxCheck

	// Idempotent is what applying this sample's selection a second time did to
	// the bytes the first application wrote.
	Idempotent Idempotency

	// Elapsed is how long this sample took end to end.
	//
	// Recorded because the harness could not previously answer "why is this
	// slow" about itself - the 75s-per-call figure had to be reconstructed from
	// a stale run log. It is also how thermal throttling becomes visible: on a
	// fanless machine the last samples of a long run are slower than the first,
	// and only per-sample timing shows that.
	Elapsed time.Duration
}

// Measurement is what a case and model produced over N samples.
type Measurement struct {
	Case  Case
	Model Model

	Samples []Sample

	// Wall is how long the whole measurement took. Compared against the sum of
	// the samples it is what concurrency actually bought, which is not
	// derivable from either number alone.
	Wall time.Duration
	// Concurrency is how many samples were in flight at once.
	Concurrency int
	// Retries is the validation budget each sample was given.
	//
	// It travels with the measurement because "valid" means a different thing
	// at 0 than at 2, and no reading of the table recovers which one it was
	// (prov-2026-b4555efc).
	Retries int

	// Resolution is the question the sample size was drawn for.
	//
	// It travels for the same reason the retry budget does, and one step
	// further: the sample count alone cannot say whether five was chosen for a
	// question five can answer or inherited from a default. A smoke rate and a
	// coarse rate can be the same fraction and are not the same claim, so the
	// word is carried rather than reconstructed from n (prov-2026-3039750e).
	Resolution Resolution
}

// Options is what a measurement is drawn with.
//
// A struct rather than three positional ints, because samples, concurrency and
// retries are all small integers and nothing about a call site would look wrong
// if two of them were swapped.
type Options struct {
	// Resolution is the question the run is asking, and it is required: a run
	// that does not state one is a sample size nobody chose (prov-2026-3039750e).
	Resolution Resolution

	// Samples is how many runs to draw per model. Zero means the resolution's
	// own count, which is the usual case - an explicit count is for
	// oversampling a question deliberately, and one below the resolution's is
	// refused rather than quietly relabelled.
	Samples int

	// Concurrency is how many samples may be in flight at once. Anything below
	// one is sequential.
	Concurrency int

	// Retries is how many times a rejected answer may be re-prompted, and it
	// is zero by default on purpose: the question the harness was built for is
	// what one call produces, and a retry loop folds Phase 5's recovery into a
	// number meant to describe a first answer.
	//
	// Raising it answers a different question - what a complete answer contains
	// once the model has got past validating - and it buys sample size for that
	// question when first-call validity is varying. The cost is that a retried
	// sample pays for every rejected answer, so timing at one budget is not
	// comparable with timing at another. The budget is recorded in the entry so
	// that stays checkable rather than remembered.
	Retries int

	// Behavior applies every valid sample's selection to a scaffolded
	// application and asserts against it.
	//
	// Off by default, and a run that leaves it off is the run it was before
	// this existed. It costs a scaffold, a dependency install, a database and a
	// boot per sample - a different order of cost from a model call - and it
	// answers the one question none of the other numbers can: whether what the
	// model chose produces something that works (prov-2026-83340ba0).
	//
	// Ignored by a case that declares no behaviour target, which is not an
	// error: the flag says "measure it where it can be measured".
	Behavior bool
}

// Run measures one case against one model, as many times as its resolution
// calls for.
//
// The count comes from Options.Resolution rather than from a number the caller
// picked, because the size of a sample is a property of the question and not of
// the harness: the same five samples that comfortably separate a 2.7x
// difference cannot tell 4/5 from 5/5 (prov-2026-3039750e).
//
// Every sample is an independent run through Phases 0 to 5 with nothing written:
// packages load, records resolve, files are created in memory, the model is
// asked, and its answer is validated. Phases 6 and 7 are exercised too, because
// a dry run injects into what Phase 3 would have written - which keeps this
// measuring the same path a real run takes rather than a shortcut through it.
//
// Retries default to zero, and Options.Retries says what raising them measures
// instead. The default cannot move: every entry already recorded was drawn at
// zero and means "validated on the first call", so a new default would re-point
// that word while the field name and the report line stayed put.
//
// Samples run concurrently up to concurrency, because drawing N independent
// samples is the canonical batched-inference workload and running them one at a
// time leaves the server idle between tokens. A concurrency of 1 or less is
// sequential.
//
// Results are written by index rather than appended, so the report reads the
// same way whatever order they finish in. A measurement that reordered between
// runs would make two of them harder to compare for no reason.
func Run(ctx context.Context, c Case, model Model, opts Options) (Measurement, error) {
	if c.Arm != "sedum" {
		return Measurement{}, fmt.Errorf("case %s has arm %q; only the sedum arm can be run today", c.ID, c.Arm)
	}
	if opts.Concurrency < 1 {
		opts.Concurrency = 1
	}

	// The resolution is settled before a single call is made, because both ways
	// of getting it wrong are only expensive afterwards: a run drawn too small
	// for its question publishes a null result it could not have seen, and a run
	// labelled for a question it was not drawn for is worse than one carrying no
	// label at all.
	if opts.Resolution == "" {
		return Measurement{}, fmt.Errorf("case %s: no resolution stated; a sample size follows from the question, so say smoke, coarse or fine", c.ID)
	}
	if _, err := ParseResolution(string(opts.Resolution)); err != nil {
		return Measurement{}, fmt.Errorf("case %s: %w", c.ID, err)
	}
	samples, err := opts.Resolution.plan(opts.Samples)
	if err != nil {
		return Measurement{}, fmt.Errorf("case %s: %w", c.ID, err)
	}

	m := Measurement{
		Case:        c,
		Model:       model,
		Concurrency: opts.Concurrency,
		Retries:     opts.Retries,
		Resolution:  opts.Resolution,
		Samples:     make([]Sample, samples),
	}

	started := time.Now()
	var wg sync.WaitGroup
	slots := make(chan struct{}, opts.Concurrency)
	for i := 0; i < samples; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			slots <- struct{}{}
			defer func() { <-slots }()
			m.Samples[i] = sample(ctx, c, model.ID, opts)
		}(i)
	}
	wg.Wait()
	m.Wall = time.Since(started)

	return m, nil
}

func sample(ctx context.Context, c Case, model string, opts Options) Sample {
	retries := opts.Retries
	started := time.Now()
	client, err := selection.NewOpenAI(model)
	if err != nil {
		return Sample{Err: err, Elapsed: time.Since(started)}
	}

	// A fresh output directory per sample, rather than the working directory.
	//
	// Phase 3 stats each authorized path to decide whether it already exists -
	// and a path it finds is marked as existing and has its markers verified,
	// which can halt the run. Against the working directory that behavior
	// depends on where the harness happened to be invoked from, which is a
	// silent confound of exactly the kind these measurements keep getting
	// bitten by. An empty directory makes every sample see the same world.
	//
	// The sample writes into it rather than running dry, because two of the
	// three signals need the files to exist (prov-2026-d61010a4). A dry run
	// reports what Phase 7 would have done, and inject.Result carries a path,
	// an action and a variant and no content, so nothing downstream can parse
	// what was never written. A sample now exercises the write path as well as
	// the decision path, which is closer to what a real run does rather than
	// further from it. Nothing is written outside this directory and it does
	// not survive the sample.
	output, err := os.MkdirTemp("", "sedum-eval-")
	if err != nil {
		return Sample{Err: err, Elapsed: time.Since(started)}
	}
	defer os.RemoveAll(output)

	result, err := pipeline.Run(ctx, pipeline.Config{
		Generators: c.Generators,
		Records:    c.Records,
		Output:     output,
		Only:       c.Only,
		DryRun:     false,
		Client:     client,
		Retries:    retries,
		Variables:  c.Variables,
	})
	if err != nil {
		// An answer Phase 5 refused and a server that was not there are
		// different measurements, and they are told apart by type rather than
		// by matching the text of a diagnostic (prov-2026-0811425c). A rejected
		// answer carries what it cost, which is the sample where the count is
		// most interesting - it is the expensive one.
		var rejected *selection.Rejection
		if errors.As(err, &rejected) {
			return Sample{
				Invalid:          true,
				Detail:           firstLine(err.Error()),
				Rules:            rulesOf(rejected),
				Calls:            rejected.Calls,
				Rejected:         rejected.Rejected,
				Completeness:     rejected.Completeness,
				PromptTokens:     rejected.PromptTokens,
				CompletionTokens: rejected.CompletionTokens,
				Elapsed:          time.Since(started),
			}
		}
		return Sample{Err: err, Detail: firstLine(err.Error()), Elapsed: time.Since(started)}
	}

	// Elapsed is set at the end rather than here, because behaviour runs after
	// this and a sample's cost has to include it.
	s := Sample{Counts: map[string]int{}}
	for _, sel := range result.Selections {
		s.Calls += sel.Calls
		s.Rejected += sel.Rejected
		s.Completeness += sel.Completeness
		s.PromptTokens += sel.PromptTokens
		s.CompletionTokens += sel.CompletionTokens
		for _, inv := range sel.Invocations {
			if s.First == "" {
				s.First = inv.Action
			}
			s.Counts[inv.Action]++
			s.Total++
		}
		// Kept whole beside the counts rather than instead of them. Counts are
		// what every stored entry and every expectation is written in; the
		// invocations are what a question about arguments has to be answered
		// from (prov-2026-2256e6fa).
		s.Invocations = append(s.Invocations, sel.Invocations...)
	}

	// Derived from the package and the record rather than from an expectation,
	// which is what lets it say something the action counts structurally
	// cannot: whether a selection accounted for the work its own created files
	// declared (prov-2026-d61010a4).
	s.Fill = fillOf(result)

	// The target's own parser over the bytes Phase 7 produced. Malformed
	// output, not wrong output - the distinction is the constraint, and the
	// report keeps it (prov-2026-d61010a4).
	s.Syntax = syntaxOf(ctx, c.Check, output, result)

	// The same invocations applied a second time, against the files the first
	// application left. A rerun that differs is an injection or marker defect,
	// asserted here against a real package rather than a fixture.
	s.Idempotent = idempotencyOf(output, result)

	// Behaviour last, and only for a sample that produced something to apply.
	//
	// A sample with no selection is unmeasured rather than failed - the same
	// rule that keeps a rejected answer out of the selection denominator. There
	// is nothing to apply, so there is nothing to say about it.
	//
	// The elapsed time is taken after this, so a behaviour run's cost is inside
	// the sample's rather than beside it. A behaviour sample is minutes where a
	// selection sample is seconds, and a timing table that hid that would be
	// describing a run nobody made.
	if opts.Behavior && c.Expect.Behavior != nil && len(s.Invocations) > 0 {
		run := RunBehavior(ctx, c.Expect.Behavior.Target, s.Invocations, c.Variables)
		s.Behavior = &run
	}
	s.Elapsed = time.Since(started)
	return s
}

// rulesOf is every rule slug a rejection cited, in attempt order, repeats kept.
//
// A slug rather than the rendered violation: the prose names an action and a
// path and is different on every sample, so it cannot be counted, while the
// slug is the field that exists to be.
func rulesOf(r *selection.Rejection) []string {
	var out []string
	for _, a := range r.Attempts {
		for _, v := range a.Violations {
			out = append(out, v.Rule)
		}
	}
	return out
}

// ActionResult is one action's observed behavior across a measurement's samples.
type ActionResult struct {
	Action string
	// Want is the invocation count a complete answer contains.
	Want int
	// Selected is how many samples included the action at all. This is the
	// number that mattered in practice: an action is usually present in full or
	// absent entirely, and a rate hides that where a mean would not.
	Selected int
	// Exact is how many samples selected exactly Want invocations.
	Exact int
	// Mean is the average invocation count across scored samples.
	Mean float64
	// FirstRate is how often this action opened the answer.
	FirstRate float64
}

// Totals is how the samples divided.
//
// Valid and Invalid are both answers the model gave, and their ratio is the
// first-call validity rate. Failed is everything that never got that far, and is
// outside every denominator.
type Totals struct {
	Valid   int
	Invalid int
	Failed  int
}

// Answered is how many samples the model responded to at all.
func (t Totals) Answered() int { return t.Valid + t.Invalid }

// Cost is what a measurement spent, in calls.
//
// Seconds cannot answer this: a per-sample time is calls multiplied by an
// unknown and varying per-call cost, which is why two arms at 5/5 could not be
// compared on cost before Phase 5 reported what it spent.
type Cost struct {
	// Calls is every model call the measurement made.
	Calls int
	// Completeness is how many of them were the completeness observation.
	Completeness int
	// FirstTry is how many answered samples took no rejection at all. This is
	// first-call validity, and it survives any retry budget - which is the
	// measurement a raised budget used to destroy (prov-2026-0811425c).
	FirstTry int

	// PromptTokens and CompletionTokens are what the calls cost. Their split
	// is the point: a case whose calls are slow because its catalog is long and
	// one whose calls are slow because it has more to emit are different
	// problems, and a total would report them alike (prov-2026-096a4d4b).
	PromptTokens     int
	CompletionTokens int
}

// PerSecond is the measurement's completion-token throughput over its wall
// clock.
//
// Completion tokens only, and that is the whole point. A prompt token is billed
// whether or not it was computed, and this server computes almost none of them:
// llama.cpp selects a slot by longest-common-prefix similarity, so a case whose
// prompt is identical across samples evaluates one token of it and reuses the
// rest. Counting billed prompt tokens as work overstated this by about four
// times, and made a case look faster the larger its prompt was
// (prov-2026-e323b805).
//
// Derived from what was counted, never from a model's advertised rate: the
// advertised number describes a single stream on the vendor's hardware, and
// what a run is bounded by is this server at this concurrency.
//
// Zero when the wall clock is zero or no endpoint reported usage, because a
// throughput of zero is a measurement and this is its absence.
func (c Cost) PerSecond(wall time.Duration) float64 {
	if wall <= 0 || c.CompletionTokens == 0 {
		return 0
	}
	return float64(c.CompletionTokens) / wall.Seconds()
}

// Spent sums what the samples cost. Failed samples contribute nothing: a call
// that never returned an answer is not a cost attributable to a selection.
func (m Measurement) Spent() Cost {
	var c Cost
	for _, s := range m.Samples {
		if s.Err != nil {
			continue
		}
		c.Calls += s.Calls
		c.Completeness += s.Completeness
		c.PromptTokens += s.PromptTokens
		c.CompletionTokens += s.CompletionTokens
		if s.Rejected == 0 {
			c.FirstTry++
		}
	}
	return c
}

// Tally divides the samples without scoring any of them.
func (m Measurement) Tally() Totals {
	var t Totals
	for _, s := range m.Samples {
		switch {
		case s.Err != nil:
			t.Failed++
		case s.Invalid:
			t.Invalid++
		default:
			t.Valid++
		}
	}
	return t
}

// Observed returns the mean invocation count of every action the valid samples
// actually selected, whether or not an expectation names it.
//
// This is how a new fixture's expectations get established. Writing them by
// reading the generator package would make the first run agree with them by
// construction, which is the one thing that would make every later number
// meaningless - so a case starts with none, a run reports what it saw, and the
// counts are set from an answer that was actually complete.
//
// It is also worth having on a case that does have expectations, because an
// action the model selects that nothing expects is invisible otherwise.
func (m Measurement) Observed() map[string]float64 {
	totals := map[string]int{}
	valid := 0
	for _, s := range m.Samples {
		if s.Err != nil || s.Invalid {
			continue
		}
		valid++
		for name, n := range s.Counts {
			totals[name] += n
		}
	}
	if valid == 0 {
		return nil
	}

	out := make(map[string]float64, len(totals))
	for name, n := range totals {
		out[name] = float64(n) / float64(valid)
	}
	return out
}

// Details returns the distinct reasons samples were invalid or failed, so that a
// reported rate can always say what its misses looked like.
func (m Measurement) Details() []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range m.Samples {
		if s.Detail == "" || seen[s.Detail] {
			continue
		}
		seen[s.Detail] = true
		out = append(out, s.Detail)
	}
	sort.Strings(out)
	return out
}

// Scored reduces a measurement to one row per expected action.
//
// Only valid samples are scored, because an invalid one has no invocation list
// to count. The invalid ones are not discarded - they are reported alongside as
// their own rate, since a model that answers completely half the time and
// invalidly the other half is a different problem from one that answers validly
// but short.
func (m Measurement) Scored() (rows []ActionResult, scored int, failed int) {
	totals := m.Tally()
	var usable []Sample
	for _, s := range m.Samples {
		if s.Err != nil || s.Invalid {
			continue
		}
		usable = append(usable, s)
	}
	if len(usable) == 0 {
		return nil, 0, totals.Failed + totals.Invalid
	}

	names := make([]string, 0, len(m.Case.Expect.Actions))
	for name := range m.Case.Expect.Actions {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		want := m.Case.Expect.Actions[name]
		row := ActionResult{Action: name, Want: want}

		total := 0
		for _, s := range usable {
			got := s.Counts[name]
			total += got
			if got > 0 {
				row.Selected++
			}
			if got == want {
				row.Exact++
			}
			if s.First == name {
				row.FirstRate++
			}
		}
		row.Mean = float64(total) / float64(len(usable))
		row.FirstRate /= float64(len(usable))
		rows = append(rows, row)
	}
	return rows, len(usable), totals.Failed + totals.Invalid
}

// firstLine trims a multi-line diagnostic to its heading. Phase 5's error
// carries every violation of every attempt, which is the right thing for a run
// and too much for a table.
func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
