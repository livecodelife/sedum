package evals

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"time"
)

// Report writes a measurement in a form that can be pasted into a record.
//
// The header carries the case, the model and the sample size on one line,
// because a rate recorded without them cannot be re-run and will be believed
// longer than it is true - which is the failure prov-2026-6d87dc11 had to be
// corrected for.
func Report(out io.Writer, m Measurement) {
	rows, scored, _ := m.Scored()
	t := m.Tally()

	fmt.Fprintf(out, "%s  model=%s  arm=%s  tightness=%s  n=%d  retries=%d\n",
		m.Case.ID, m.Model.Label(), m.Case.Arm, m.Case.Tightness, len(m.Samples), m.Retries)
	fmt.Fprintln(out, "  "+timing(m))

	// Validity before anything else: it is the number every rate below is
	// conditioned on. At the default budget of no retries this is how often one
	// call produced something Phase 5 accepted. At a raised budget it is a
	// weaker claim about more calls, so the line says which it is rather than
	// calling itself "first call" either way (prov-2026-b4555efc).
	fmt.Fprintf(out, "  valid %s: %s", calls(m.Retries), wilson(t.Valid, t.Answered()))
	if t.Failed > 0 {
		fmt.Fprintf(out, "   (%d run(s) never reached the model and are excluded)", t.Failed)
	}
	fmt.Fprintln(out)

	// The stronger claim, still available at a raised budget: an answer with no
	// rejections validated first try whatever the budget allowed. Printed only
	// when the budget could have hidden it, since at zero retries the line
	// above already is this number (prov-2026-0811425c).
	spent := m.Spent()
	if m.Retries > 0 && t.Answered() > 0 {
		fmt.Fprintf(out, "  valid first call: %s\n", wilson(spent.FirstTry, t.Answered()))
	}
	if spent.Calls > 0 {
		fmt.Fprintf(out, "  cost: %d call(s) over %d sample(s), mean %.2f",
			spent.Calls, t.Answered(), float64(spent.Calls)/float64(t.Answered()))
		if spent.Completeness > 0 {
			fmt.Fprintf(out, "  (%d completeness observation(s))", spent.Completeness)
		}
		fmt.Fprintln(out)
	}

	// The prompt/completion split, printed only when the endpoint reported it.
	// A server that fills no usage block gets no line rather than a line of
	// zeroes, because zero tokens is a measurement and this is its absence
	// (prov-2026-096a4d4b). The split is the point: a long catalog and a long
	// answer are different problems with different responses.
	if spent.PromptTokens > 0 || spent.CompletionTokens > 0 {
		perCall := func(n int) float64 { return float64(n) / float64(spent.Calls) }
		fmt.Fprintf(out, "  tokens: %s prompt + %s completion, per call %.0f + %.0f\n",
			thousands(spent.PromptTokens), thousands(spent.CompletionTokens),
			perCall(spent.PromptTokens), perCall(spent.CompletionTokens))
	}

	for _, d := range m.Details() {
		fmt.Fprintf(out, "    %s\n", d)
	}

	if scored == 0 {
		fmt.Fprintln(out, "  no valid samples to score")
		return
	}

	// A case with no expectations is a new fixture, not a broken one. Report
	// what was selected so the counts can be set from an answer that was
	// actually complete, rather than written by reading the package and then
	// agreed with by construction.
	if len(m.Case.Expect.Actions) == 0 {
		fmt.Fprintln(out, "  no expectations declared - observed selections, for establishing them:")
		observed := m.Observed()
		names := make([]string, 0, len(observed))
		for name := range observed {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Fprintf(out, "    %-24s mean %.2f\n", name, observed[name])
		}
		return
	}

	fmt.Fprintf(out, "  %-24s %5s %18s %18s %7s %7s\n",
		"action", "want", "selected", "exact", "mean", "first")
	for _, r := range rows {
		fmt.Fprintf(out, "  %-24s %5d %18s %18s %7.2f %6.0f%%\n",
			r.Action, r.Want,
			wilson(r.Selected, scored), wilson(r.Exact, scored),
			r.Mean, r.FirstRate*100)
	}

	if m.Case.Expect.Behavior != nil {
		fmt.Fprintln(out, "  behavior: declared but not measured; applying and running the target is not implemented")
	}
}

// timing renders what the run cost, which the harness previously could not say
// about itself - the 75s-per-call figure behind prov-2026-6d87dc11 had to be
// reconstructed from a stale run log.
//
// Slowest is printed beside fastest rather than only a mean, because on a
// fanless machine the spread is the throttling signal: a run whose last samples
// are much slower than its first is thermally limited, and a mean hides that.
// Wall against the summed sample time is what concurrency actually bought.
func timing(m Measurement) string {
	if len(m.Samples) == 0 {
		return "no samples"
	}

	var sum, fastest, slowest time.Duration
	for i, s := range m.Samples {
		sum += s.Elapsed
		if i == 0 || s.Elapsed < fastest {
			fastest = s.Elapsed
		}
		if s.Elapsed > slowest {
			slowest = s.Elapsed
		}
	}
	mean := sum / time.Duration(len(m.Samples))

	out := fmt.Sprintf("wall %s  per-sample fastest %s / mean %s / slowest %s  concurrency %d",
		round(m.Wall), round(fastest), round(mean), round(slowest), m.Concurrency)
	if m.Concurrency > 1 && m.Wall > 0 {
		out += fmt.Sprintf("  (%.1fx over sequential)", float64(sum)/float64(m.Wall))
	}
	return out
}

func round(d time.Duration) time.Duration {
	if d >= time.Second {
		return d.Round(100 * time.Millisecond)
	}
	return d.Round(time.Millisecond)
}

// calls names the budget a validity rate was measured under.
//
// "first call" is the strongest form of the claim and the one every entry
// recorded so far carries. A raised budget makes it a different measurement,
// and a line that read the same either way would be the field-changing-meaning
// failure the entry format exists to avoid.
func calls(retries int) string {
	if retries <= 0 {
		return "first call"
	}
	return fmt.Sprintf("within %d calls", retries+1)
}

// Interval is a rate with the uncertainty its sample size gives it.
type Interval struct {
	Rate      float64
	Low       float64
	High      float64
	Successes int
	Samples   int
}

// String renders the fraction and the interval together.
//
// The fraction never disappears behind the interval: five samples is the fact
// that produced the width, and a reader who sees only the bounds has lost it
// (prov-2026-0baaa119).
func (i Interval) String() string {
	if i.Samples == 0 {
		return "-"
	}
	return fmt.Sprintf("%d/%d [%.2f,%.2f]", i.Successes, i.Samples, i.Low, i.High)
}

// Overlaps reports whether two intervals leave open that the rates are the same.
// It is not a test and decides nothing; it is the condition under which a
// reported difference is not distinguishable from sampling.
func (i Interval) Overlaps(o Interval) bool {
	return i.Low <= o.High && o.Low <= i.High
}

// wilson is the Wilson score interval at 95%.
//
// Chosen over the normal approximation because the approximation collapses to
// zero width at 0/n and n/n, and n/n is exactly the reading most likely to be
// over-believed: 5/5 is not evidence of a rate of 1.00, and an interval that
// said [1.00,1.00] would claim it was.
func wilson(successes, samples int) Interval {
	if samples == 0 {
		return Interval{}
	}

	const z = 1.96
	n := float64(samples)
	p := float64(successes) / n

	denom := 1 + z*z/n
	center := (p + z*z/(2*n)) / denom
	spread := z * math.Sqrt(p*(1-p)/n+z*z/(4*n*n)) / denom

	return Interval{
		Rate:      p,
		Low:       math.Max(0, center-spread),
		High:      math.Min(1, center+spread),
		Successes: successes,
		Samples:   samples,
	}
}

// thousands keeps a token count readable without losing the small ones: a
// four-figure count is exact, and anything larger is rounded to a tenth of a
// thousand.
func thousands(n int) string {
	if n < 10_000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%.1fk", float64(n)/1000)
}

func ratio(n, of int) string {
	if of == 0 {
		return "-"
	}
	return fmt.Sprintf("%d/%d", n, of)
}

// Compare writes two measurements side by side.
//
// This is the operation the manual measurement actually needed and had to
// improvise with a git worktree: a rate on its own is not evidence that a change
// helped, and the only way to find out is to run both against the same case.
func Compare(out io.Writer, a, b Measurement) {
	aRows, aScored, _ := a.Scored()
	bRows, bScored, _ := b.Scored()

	fmt.Fprintf(out, "%s  %s (n=%d)  vs  %s (n=%d)\n",
		a.Case.ID, label(a), aScored, label(b), bScored)

	if aScored == 0 || bScored == 0 {
		fmt.Fprintln(out, "  one side has no scoreable samples; nothing to compare")
		return
	}

	byName := map[string]ActionResult{}
	for _, r := range bRows {
		byName[r.Action] = r
	}

	fmt.Fprintf(out, "  %-24s %18s %18s %9s\n", "action", "selected A", "selected B", "delta")
	for _, ra := range aRows {
		rb, ok := byName[ra.Action]
		if !ok {
			continue
		}
		ia := wilson(ra.Selected, aScored)
		ib := wilson(rb.Selected, bScored)

		// An overlap is marked rather than the delta being suppressed. The
		// numbers are what was measured; what the mark says is that these
		// samples do not distinguish the two, which is the case this table
		// used to report as a difference (prov-2026-0baaa119).
		delta := fmt.Sprintf("%+8.0f%%", (ib.Rate-ia.Rate)*100)
		if ia.Overlaps(ib) {
			delta += " ~"
		}
		fmt.Fprintf(out, "  %-24s %18s %18s %s\n", ra.Action, ia, ib, delta)
	}

	fmt.Fprintln(out, "  ~ marks rows whose intervals overlap: these runs do not distinguish the arms")
}

func label(m Measurement) string {
	parts := []string{m.Model.Label()}
	if m.Case.Tightness != "" {
		parts = append(parts, m.Case.Tightness)
	}
	if m.Case.Arm != "" {
		parts = append(parts, m.Case.Arm)
	}
	return strings.Join(parts, "/")
}
