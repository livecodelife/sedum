package evals

import (
	"context"
	"fmt"
	"sort"
	"strings"

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
	// first time by treating the two alike. Retries are zero here, so an answer
	// that would have been repaired by the retry loop shows up as no answer at
	// all - and how often a first call validates is one of the more interesting
	// things to know about a model, not noise to drop.
	Invalid bool

	// Detail is why a sample was invalid or failed, kept so a rate is never
	// reported without the ability to say what the misses looked like.
	Detail string

	// Err is set when the run could not be made at all - an endpoint that was
	// down, a package that would not load. Excluded from every rate, because an
	// unreachable server is not a model that chose badly.
	Err error
}

// Measurement is what a case and model produced over N samples.
type Measurement struct {
	Case    Case
	Model   string
	Samples []Sample
}

// Run measures one case against one model, samples times.
//
// Every sample is an independent run through Phases 0 to 5 with nothing written:
// packages load, records resolve, files are created in memory, the model is
// asked, and its answer is validated. Phases 6 and 7 are exercised too, because
// a dry run injects into what Phase 3 would have written - which keeps this
// measuring the same path a real run takes rather than a shortcut through it.
//
// Retries are zero on purpose. The question is what one call produces, and a
// retry loop would fold Phase 5's recovery into a number meant to describe the
// model's first answer.
func Run(ctx context.Context, c Case, model string, samples int) (Measurement, error) {
	if c.Arm != "sedum" {
		return Measurement{}, fmt.Errorf("case %s has arm %q; only the sedum arm can be run today", c.ID, c.Arm)
	}

	m := Measurement{Case: c, Model: model}
	for i := 0; i < samples; i++ {
		m.Samples = append(m.Samples, sample(ctx, c, model))
	}
	return m, nil
}

func sample(ctx context.Context, c Case, model string) Sample {
	client, err := selection.NewOpenAI(model)
	if err != nil {
		return Sample{Err: err}
	}

	result, err := pipeline.Run(ctx, pipeline.Config{
		Generators: c.Generators,
		Records:    c.Records,
		Output:     "",
		DryRun:     true,
		Client:     client,
		Retries:    0,
	})
	if err != nil {
		// Provisional classification. internal/selection exports no error type
		// for a rejected response, so this matches the message it builds. It is
		// fragile in the ordinary way a string match is, and the honest fix is
		// an exported sentinel there - worth doing before anything depends on
		// this number.
		if strings.Contains(err.Error(), "did not validate") {
			return Sample{Invalid: true, Detail: firstLine(err.Error())}
		}
		return Sample{Err: err, Detail: firstLine(err.Error())}
	}

	s := Sample{Counts: map[string]int{}}
	for _, sel := range result.Selections {
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
