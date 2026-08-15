package evals

import (
	"fmt"
	"io"
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

	fmt.Fprintf(out, "%s  model=%s  arm=%s  tightness=%s  n=%d\n",
		m.Case.ID, m.Model.Label(), m.Case.Arm, m.Case.Tightness, len(m.Samples))
	fmt.Fprintln(out, "  "+timing(m))

	// First-call validity before anything else. Retries are zero here, so this
	// is how often one call produced something Phase 5 accepted - the number
	// every rate below is conditioned on.
	fmt.Fprintf(out, "  valid first call: %s", ratio(t.Valid, t.Answered()))
	if t.Failed > 0 {
		fmt.Fprintf(out, "   (%d run(s) never reached the model and are excluded)", t.Failed)
	}
	fmt.Fprintln(out)

	for _, d := range m.Details() {
		fmt.Fprintf(out, "    %s\n", d)
	}

	if scored == 0 {
		fmt.Fprintln(out, "  no valid samples to score")
		return
	}

	fmt.Fprintf(out, "  %-24s %5s %9s %8s %7s %7s\n",
		"action", "want", "selected", "exact", "mean", "first")
	for _, r := range rows {
		fmt.Fprintf(out, "  %-24s %5d %8s %7s %7.2f %6.0f%%\n",
			r.Action, r.Want,
			ratio(r.Selected, scored), ratio(r.Exact, scored),
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

	fmt.Fprintf(out, "  %-24s %11s %11s %9s\n", "action", "selected A", "selected B", "delta")
	for _, ra := range aRows {
		rb, ok := byName[ra.Action]
		if !ok {
			continue
		}
		rateA := float64(ra.Selected) / float64(aScored)
		rateB := float64(rb.Selected) / float64(bScored)
		fmt.Fprintf(out, "  %-24s %10s %10s %+8.0f%%\n",
			ra.Action, ratio(ra.Selected, aScored), ratio(rb.Selected, bScored),
			(rateB-rateA)*100)
	}

	fmt.Fprintln(out, "  a delta smaller than the sampling noise at this n is not a result")
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
