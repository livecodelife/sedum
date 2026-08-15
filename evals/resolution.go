package evals

import "fmt"

// Resolution is how finely a run is asking its question, and the sample size
// follows from it.
//
// The default of five samples was never chosen against a question - it was the
// number the first run happened to use, and every run since inherited it.
// prov-2026-0baaa119 made what it buys visible: a perfect 5/5 is consistent with
// a true rate of 0.57, while 30/30 narrows the same reading to [0.89, 1.00]. So
// the number is not a property of the harness. It is a property of the effect
// being looked for, and the questions this matrix asks differ by an order of
// magnitude in that.
//
// Deliberately not called a tier, because a case already has one: complexity is
// an ordinal tier of the application. Two axes sharing a word would read as one
// axis in a results file.
type Resolution string

const (
	// Smoke proves the plumbing and measures nothing. A new case loads, a new
	// package resolves, an endpoint answers.
	Smoke Resolution = "smoke"

	// Coarse separates differences that are enormous, which is most of what
	// widening the matrix asks: whether a 4b model selects usefully where a 14b
	// does, whether a framework's package works at all, whether one arm is
	// several times slower than another. The chi arm at 2.7 times the rails arm
	// is a coarse result and five samples established it comfortably.
	Coarse Resolution = "coarse"

	// Fine moves a rate that is already high - 4/5 to 5/5, or a selection rate
	// from 80% to 100%. These are the questions where an interval at five
	// samples is decoration.
	Fine Resolution = "fine"
)

// DefaultResolution is coarse, and that is now a decision rather than an
// inheritance: five samples is the size of a coarse question, and a run asking
// a finer one says so.
const DefaultResolution = Coarse

// Samples is the sample size the resolution implies.
//
// Fine is thirty, and thirty is an upper bound rather than a threshold. It comes
// from comparing intervals, which is deliberately conservative - two overlapping
// intervals are slow to claim a difference, and a proper two-proportion
// comparison separates earlier. What the conservatism does not change is the
// direction of the error at five: too few for a fine question by a wide margin
// either way.
//
// Smoke is two rather than one because a single sample exercises no
// concurrency, no aggregation over more than one outcome, and no per-sample
// state at all - the three things most likely to be broken in plumbing that has
// never run.
func (r Resolution) Samples() int {
	switch r {
	case Smoke:
		return 2
	case Coarse:
		return 5
	case Fine:
		return 30
	}
	return 0
}

// Cites reports whether a rate drawn at this resolution may be cited as a
// measurement. A smoke rate may not: it exists to prove the plumbing, and it
// says so on every line it appears on.
func (r Resolution) Cites() bool { return r != Smoke && r != "" }

// ParseResolution reads a resolution by name, naming what each is for rather
// than only listing them - the choice is the point, and a bare list of three
// words invites picking the middle one again.
func ParseResolution(s string) (Resolution, error) {
	switch r := Resolution(s); r {
	case Smoke, Coarse, Fine:
		return r, nil
	}
	return "", fmt.Errorf("resolution %q: want smoke (plumbing only, n=%d), coarse (large differences, n=%d) or fine (moving a rate that is already high, n=%d)",
		s, Smoke.Samples(), Coarse.Samples(), Fine.Samples())
}

// plan is the sample size a run will actually draw.
//
// An explicit count is honoured when it is at least what the resolution needs,
// because oversampling a coarse question is a cost decision and nobody is misled
// by it. Undersampling is refused: a run labelled fine and drawn at five would
// put a resolution in the results file that its own numbers cannot support,
// which is exactly the mislabelling this type exists to prevent. A null result
// from a run too small to have seen the effect is not a null result.
//
// Smoke is exempt, because there is no measurement to mislabel. One sample is a
// perfectly good plumbing check - the imprint's own taxonomy says n=1 or 2 - and
// a rate drawn there is uncitable whatever its size.
func (r Resolution) plan(explicit int) (int, error) {
	if explicit <= 0 {
		return r.Samples(), nil
	}
	if r.Cites() && explicit < r.Samples() {
		return 0, fmt.Errorf("%d sample(s) is below the %s resolution's %d: raise the count or state the resolution the question actually needs",
			explicit, r, r.Samples())
	}
	return explicit, nil
}
