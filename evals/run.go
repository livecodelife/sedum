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
}

// Options is what a measurement is drawn with.
//
// A struct rather than three positional ints, because samples, concurrency and
// retries are all small integers and nothing about a call site would look wrong
// if two of them were swapped.
type Options struct {
	// Samples is how many runs to draw per model.
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
}

// Run measures one case against one model, samples times.
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

	m := Measurement{
		Case:        c,
		Model:       model,
		Concurrency: opts.Concurrency,
		Retries:     opts.Retries,
		Samples:     make([]Sample, opts.Samples),
	}

	started := time.Now()
	var wg sync.WaitGroup
	slots := make(chan struct{}, opts.Concurrency)
	for i := 0; i < opts.Samples; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			slots <- struct{}{}
			defer func() { <-slots }()
			m.Samples[i] = sample(ctx, c, model.ID, opts.Retries)
		}(i)
	}
	wg.Wait()
	m.Wall = time.Since(started)

	return m, nil
}

func sample(ctx context.Context, c Case, model string, retries int) Sample {
	started := time.Now()
	client, err := selection.NewOpenAI(model)
	if err != nil {
		return Sample{Err: err, Elapsed: time.Since(started)}
	}

	// A fresh output directory per sample, rather than the working directory.
	//
	// Nothing is written under a dry run, but Phase 3 still stats each
	// authorized path to decide whether it already exists - and a path it finds
	// is marked as existing and has its markers verified, which can halt the
	// run. Against the working directory that behavior depends on where the
	// harness happened to be invoked from, which is a silent confound of
	// exactly the kind these measurements keep getting bitten by. An empty
	// directory makes every sample see the same world.
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
		DryRun:     true,
		Client:     client,
		Retries:    retries,
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

	s := Sample{Counts: map[string]int{}, Elapsed: time.Since(started)}
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
	}
	return s
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

// PerSecond is the measurement's token throughput over its wall clock.
//
// Derived from what was counted, never from a model's advertised rate: the
// advertised number describes a single stream on the vendor's hardware, and
// what a run is bounded by is this server, at this concurrency, with whatever
// prefix caching it was configured for (prov-2026-945fa0aa).
//
// Zero when the wall clock is zero or no endpoint reported usage, because a
// throughput of zero is a measurement and this is its absence.
func (c Cost) PerSecond(wall time.Duration) float64 {
	total := c.PromptTokens + c.CompletionTokens
	if wall <= 0 || total == 0 {
		return 0
	}
	return float64(total) / wall.Seconds()
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
