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

	fmt.Fprintf(out, "%s  model=%s  arm=%s  tightness=%s  res=%s  n=%d  retries=%d\n",
		m.Case.ID, m.Model.Label(), m.Case.Arm, m.Case.Tightness,
		resolutionOf(m.Resolution), len(m.Samples), m.Retries)

	// A smoke run says what it is before it says any number, because the rate
	// below is the sort of thing that gets pasted into a record on its own. It
	// exists to prove the plumbing and is not a measurement, and the only place
	// that can be made unmissable is above the numbers (prov-2026-3039750e).
	if m.Resolution == Smoke {
		fmt.Fprintf(out, "  SMOKE: plumbing only at n=%d. Not a measurement, and not to be cited as one.\n",
			len(m.Samples))
	}

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
		fmt.Fprintf(out, "  tokens: %s prompt billed + %s completion, per call %.0f + %.0f\n",
			thousands(spent.PromptTokens), thousands(spent.CompletionTokens),
			perCall(spent.PromptTokens), perCall(spent.CompletionTokens))

		// Completion tokens only, because those are the ones the machine
		// produced. A server that reuses a cached prefix evaluates one token of
		// a two-thousand-token prompt, so counting billed prompt tokens as work
		// would make a case look faster the larger its prompt is
		// (prov-2026-e323b805).
		if tps := spent.PerSecond(m.Wall); tps > 0 {
			fmt.Fprintf(out, "  throughput: %.1f completion tok/s at concurrency %d\n",
				tps, m.Concurrency)
		}
	}

	for _, d := range m.Details() {
		fmt.Fprintf(out, "    %s\n", d)
	}

	if scored == 0 {
		fmt.Fprintln(out, "  no valid samples to score")
		return
	}

	// The baseline arm has no catalog, so "what was selected" is not a question
	// about it. What it produced is a set of paths against the ones the record
	// authorized, which is the completeness question that survives without a
	// vocabulary (prov-2026-a4dbe65c).
	if m.Case.WithoutPackage() {
		baseline(out, m)
		signals(out, m)
		behavior(out, m)
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
		// The derived rungs do not depend on an expectation, which is exactly
		// why they belong here: a fixture with no counts written yet still
		// reports what its own templates asked for and what its output parsed
		// as. Behaviour is on the same footing - it is scored against a
		// declared target rather than against the action counts.
		signals(out, m)
		behavior(out, m)
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

	bindings(out, m, scored)
	signals(out, m)
	behavior(out, m)
}

// behavior renders what applying the selections produced, as its own table.
//
// Never folded into selection. The case this exists for is a sample that
// selects perfectly and produces a service that violates its own record, and a
// combined rate would be the arithmetic that hides it - the same conflation
// bindings had to be split out of (prov-2026-2b121b62, prov-2026-83340ba0).
func behavior(out io.Writer, m Measurement) {
	if m.Case.Expect.Behavior == nil {
		return
	}

	t, measured := m.Behavior()
	if !measured {
		// Declared and not asked for. Said plainly, because a reader who cannot
		// tell "not measured" from "measured and nothing worked" will read the
		// absence of a number as the absence of a problem.
		fmt.Fprintf(out, "\n  behavior: target %q declared, not measured (run with -behavior)\n",
			m.Case.Expect.Behavior.Target)
		return
	}

	fmt.Fprintf(out, "\n  behavior  target=%s  %s\n",
		m.Case.Expect.Behavior.Target, round(t.Elapsed))

	if t.Measured == 0 {
		fmt.Fprintf(out, "    no sample produced a selection to apply")
		if t.Errored > 0 {
			fmt.Fprintf(out, "; %d harness error(s)", t.Errored)
		}
		fmt.Fprintln(out)
		return
	}

	fmt.Fprintf(out, "    working: %s of applied samples", wilson(t.Working, t.Measured))
	if want := m.Case.Expect.Behavior.PassRate; want > 0 {
		fmt.Fprintf(out, "   (expected %.0f%%)", want*100)
	}
	fmt.Fprintln(out)

	// The two ways of not working, kept apart. A service that never booted and
	// one that booted and answered wrongly are different findings.
	if t.Disagreed > 0 {
		fmt.Fprintf(out, "    disagreed: %d sample(s) built and booted and failed an assertion\n", t.Disagreed)
	}
	if t.Broke > 0 {
		fmt.Fprintf(out, "    broke: %d sample(s) never reached the assertions%s\n", t.Broke, byPhase(t.Phases))
		// Why, not just where. A phase name says a build failed; this says what
		// the compiler said, which is the difference between a finding and a
		// reason to re-run the harness by hand (prov-2026-93829987).
		for _, d := range byCount(t.Details) {
			fmt.Fprintf(out, "      %d of %d:\n", d.n, t.Broke)
			for _, line := range strings.Split(d.name, "\n") {
				fmt.Fprintf(out, "        %s\n", line)
			}
		}
	}
	if t.Errored > 0 {
		fmt.Fprintf(out, "    excluded: %d sample(s) the harness could not run at all\n", t.Errored)
	}
	if t.Checks > 0 {
		fmt.Fprintf(out, "    assertions: %d of %d held across %d applied sample(s)\n",
			t.Passed, t.Checks, t.Measured)
	}

	// Which contracts broke, most often first. A rate without this says
	// something is wrong and not what, which is the shape of finding this
	// harness was built to produce.
	if len(t.Failures) > 0 {
		fmt.Fprintln(out, "    failed:")
		for _, f := range byCount(t.Failures) {
			fmt.Fprintf(out, "      %3d/%-3d %s\n", f.n, t.Measured, f.name)
		}
	}
}

type counted struct {
	name string
	n    int
}

// byCount orders by frequency and then by name, so two runs of the same
// measurement print the same list.
func byCount(counts map[string]int) []counted {
	out := make([]counted, 0, len(counts))
	for name, n := range counts {
		out = append(out, counted{name: name, n: n})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].n != out[j].n {
			return out[i].n > out[j].n
		}
		return out[i].name < out[j].name
	})
	return out
}

func byPhase(phases map[string]int) string {
	if len(phases) == 0 {
		return ""
	}
	parts := make([]string, 0, len(phases))
	for _, p := range byCount(phases) {
		parts = append(parts, fmt.Sprintf("%s x%d", p.name, p.n))
	}
	return " (" + strings.Join(parts, ", ") + ")"
}

// bindings renders what the model bound, beside the selection table and never
// folded into it.
//
// A combined rate would hide the case this exists to surface: four smoke
// samples scored a perfect six-of-six on selection while every one of them
// bound `default` wrong (prov-2026-2b121b62). Selection and binding are
// different measurements and the table says so by being a different table.
func bindings(out io.Writer, m Measurement, scored int) {
	results := m.ScoredBindings()
	if len(results) == 0 {
		return
	}

	// The interval goes on samples and nowhere else. Two invocations inside one
	// answer are one observation seen twice, so an interval over invocations
	// would report confidence nobody bought - and the invocation ratio is
	// printed beside it without one, because a sample rate cannot say whether a
	// failing sample got one invocation wrong or all of them
	// (prov-2026-7cb96bf0).
	fmt.Fprintf(out, "  %-24s %-12s %18s %s\n", "action", "kwarg", "samples bound", "invocations")
	for _, r := range results {
		for _, k := range r.Kwargs {
			fmt.Fprintf(out, "  %-24s %-12s %18s %d/%d\n",
				r.Action, k.Kwarg, wilson(k.CleanSamples, k.Samples), k.Correct, k.Scored)
		}
		// An expected invocation nothing answered and an answered one nothing
		// expected are reported rather than dropped. A key bound wrong produces
		// one of each, which is the honest reading: the model did not fumble an
		// argument, it addressed something nobody asked for.
		if r.Missing > 0 || r.Unexpected > 0 {
			fmt.Fprintf(out, "  %-24s %d expected invocation(s) unanswered, %d answered invocation(s) nobody expected, over %d sample(s)\n",
				r.Action, r.Missing, r.Unexpected, scored)
		}
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

// resolutionOf names a measurement's resolution, and says so plainly when there
// is none rather than defaulting the word to "coarse". An entry drawn before
// resolutions existed was drawn at five samples for no stated reason, which is
// the fact this field was added to stop hiding.
func resolutionOf(r Resolution) string {
	if r == "" {
		return "unstated"
	}
	return string(r)
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

	fmt.Fprintf(out, "%s  %s (res=%s n=%d)  vs  %s (res=%s n=%d)\n",
		a.Case.ID, label(a), resolutionOf(a.Resolution), aScored,
		label(b), resolutionOf(b.Resolution), bScored)

	if aScored == 0 || bScored == 0 {
		fmt.Fprintln(out, "  one side has no scoreable samples; nothing to compare")
		return
	}

	// A comparison is the one operation a smoke run must never be read as, so
	// it is named here as well as on each arm's own report. Two plumbing checks
	// differ by whatever the model did that morning.
	if !a.Resolution.Cites() || !b.Resolution.Cites() {
		fmt.Fprintln(out, "  SMOKE: at least one arm proves plumbing only. This is not a comparison.")
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

	// Overlap is conservative on purpose, and saying so is what keeps it from
	// being read as a test. A two-proportion comparison separates two rates
	// earlier than their intervals stop overlapping, so a row marked ~ is "not
	// shown by these samples", never "shown to be the same" - and the sample
	// size overlap asks for is an upper bound on what the question needed
	// (prov-2026-3039750e).
	fmt.Fprintln(out, "    it is a conservative reading, not a test: ~ means not distinguished here, never no difference")
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

// signals prints the three derived rungs between "the right actions were
// selected" and "the application works".
//
// Each carries its own denominator, because the three count different things:
// anchors, files, and files again. A single combined number would be the
// arithmetic that hides which rung moved, which is the conflation the whole
// ladder exists to prevent (prov-2026-d61010a4).
func signals(out io.Writer, m Measurement) {
	t := m.Signals()
	if t.Applied == 0 {
		return
	}

	fmt.Fprintf(out, "\n  signals  %d sample(s) applied\n", t.Applied)

	// Anchor fill and idempotency need a package and an invocation list. The
	// baseline arm has neither, and "nothing fillable was planted" would read
	// as a finding about the answer rather than as the absence of a catalog -
	// so the two rungs are skipped rather than reported empty
	// (prov-2026-a4dbe65c). Syntactic validity survives, because a file is a
	// file whoever wrote it.
	packaged := m.Case.Arm != "baseline"

	// Derived from the package and the record rather than from an expectation,
	// which is what it is here to add. An anchor filled by the wrong action
	// counts as filled, so this is a completeness number and not a correctness
	// one.
	if _, ok := t.Fill.Rate(); ok && packaged {
		fmt.Fprintf(out, "    anchors filled: %s of fillable planted\n",
			wilson(t.Fill.Filled, t.Fill.Planted))
	} else if packaged {
		fmt.Fprintln(out, "    anchors filled: nothing fillable was planted")
	}

	// Syntactic validity, named as such on the line itself. Output that parses
	// and misbehaves passes this, and a reader who takes it for correctness has
	// been given the wrong rung (prov-2026-c5697387).
	switch _, ok := t.Syntax.Rate(); {
	case ok:
		fmt.Fprintf(out, "    parses: %s of file(s) checked (syntax only, not correctness)\n",
			wilson(t.Syntax.Parsed, t.Syntax.Checked))
	case len(t.Unavailable) > 0:
		// Not a rate of zero. The distinction is the whole reason the field
		// exists: an uninstalled interpreter would otherwise report as every
		// file being malformed.
		fmt.Fprintf(out, "    parses: not checked - %s is not installed here\n",
			strings.Join(t.Unavailable, ", "))
	case len(m.Case.Check) == 0:
		fmt.Fprintln(out, "    parses: the case declares no check command")
	default:
		fmt.Fprintln(out, "    parses: no written file matched a check command")
	}

	if _, ok := t.Idempotent.Rate(); ok && packaged {
		fmt.Fprintf(out, "    idempotent: %s of file(s) unchanged by a second application\n",
			wilson(t.Idempotent.Stable, t.Idempotent.Files))
	}

	// The rate says something is wrong; these say what. Both are needed, and
	// the count beside each name is how often it happened across the samples.
	for _, group := range []struct {
		label string
		lines []string
	}{
		{"anchors left unfilled", fillMissed(packaged, t)},
		{"malformed", t.Malformed},
		{"changed on reapply", t.Unstable},
		{"second application failed", t.Errs},
	} {
		if len(group.lines) == 0 {
			continue
		}
		fmt.Fprintf(out, "    %s:\n", group.label)
		counts := map[string]int{}
		for _, line := range group.lines {
			counts[line]++
		}
		for _, d := range byCount(counts) {
			fmt.Fprintf(out, "      %3d  %s\n", d.n, d.name)
		}
	}
}

// baseline renders what an arm with no package produced: which of the paths the
// record authorized the model wrote.
//
// No selection table, no bindings, no anchor fill and no idempotency. All four
// need a catalog or an invocation list, and printing zeroes for them would
// report a measurement that was not made (prov-2026-a4dbe65c).
func baseline(out io.Writer, m Measurement) {
	fmt.Fprintf(out, "  arm=baseline: no package, no vocabulary. Selection, binding, anchor fill\n")
	fmt.Fprintf(out, "    and idempotency are undefined here and are not reported.\n")

	// The comparison this arm exists for is only honest against a sedum entry
	// drawn the same way, and the retry budget is the axis they can differ on:
	// a baseline gets one call by construction, because re-prompting a build
	// error is a repair loop this repository does not build.
	fmt.Fprintf(out, "    one call per sample; comparable to a sedum entry at retries 0.\n")

	var wrote, want int
	missingBy := map[string]int{}
	unexpectedBy := map[string]int{}
	for _, s := range m.Samples {
		if s.Err != nil || s.Invalid {
			continue
		}
		wrote += len(s.Files)
		want += len(s.Files) + len(s.Missing)
		for _, p := range s.Missing {
			missingBy[p]++
		}
		for _, p := range s.Unexpected {
			unexpectedBy[p]++
		}
	}
	if want == 0 {
		return
	}

	fmt.Fprintf(out, "    authorized paths written: %s\n", wilson(wrote, want))
	if len(missingBy) > 0 {
		fmt.Fprintln(out, "    never written:")
		for _, d := range byCount(missingBy) {
			fmt.Fprintf(out, "      %3d  %s\n", d.n, d.name)
		}
	}
	// A path nothing authorized was not written to disk, and saying so matters:
	// code at the wrong path is a different finding from no code at all.
	if len(unexpectedBy) > 0 {
		fmt.Fprintln(out, "    written to a path nothing authorized (discarded):")
		for _, d := range byCount(unexpectedBy) {
			fmt.Fprintf(out, "      %3d  %s\n", d.n, d.name)
		}
	}
}

// fillMissed is the unfilled anchors, or nothing at all for an arm that has no
// anchors to leave unfilled.
func fillMissed(packaged bool, t SignalTally) []string {
	if !packaged {
		return nil
	}
	return t.Fill.Missed
}
