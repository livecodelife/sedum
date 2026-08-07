package filetmpl

// Overlap answers a question the ranking cannot: does any path exist that two
// patterns both match? A package loader holds a set of patterns and no path, so
// without this it would report app/{name}.rb and lib/{name}.rb as a tie merely
// because they are equally specific, and reject a package over two directories
// that can never collide (prov-2026-2ac1e4ca).
//
// It is decided exactly rather than estimated. Both failure directions are
// invisible: over-approximating rejects legitimate packages, and
// under-approximating admits a tie into generation, which is written to assume
// ties were already eliminated and would pick a template by iteration order.

// An atom is one step of a segment: a specific byte, or a wildcard that accepts
// any byte and may repeat. Expanding a segment into atoms turns it into an NFA
// whose states are positions in the list, which makes intersection a walk over
// the product of two such lists.
type atom struct {
	ch   byte
	wild bool
}

func atoms(seg segment) []atom {
	var out []atom
	for _, p := range seg {
		if p.kind != literal {
			out = append(out, atom{wild: true})
			continue
		}
		for i := 0; i < len(p.text); i++ {
			out = append(out, atom{ch: p.text[i]})
		}
	}
	return out
}

// overlap reports whether some path matches both patterns.
func overlap(a, b *pattern) bool {
	if len(a.segments) != len(b.segments) {
		return false
	}
	for i := range a.segments {
		if !segmentsOverlap(atoms(a.segments[i]), atoms(b.segments[i])) {
			return false
		}
	}
	return true
}

// segmentsOverlap decides whether the two atom lists accept a common string, by
// breadth-first search over the product state space. It is bounded by
// (len(a)+1) * (len(b)+1) states, which is small for anything a file template
// could plausibly be.
func segmentsOverlap(a, b []atom) bool {
	type state struct{ i, j int }

	start := state{}
	seen := map[state]bool{start: true}
	queue := []state{start}

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]

		if st.i == len(a) && st.j == len(b) {
			return true
		}
		// One side is exhausted while the other still demands at least one
		// more character: no common string extends this state.
		if st.i == len(a) || st.j == len(b) {
			continue
		}

		for _, c := range commonBytes(a[st.i], b[st.j]) {
			for _, i := range step(a, st.i, c) {
				for _, j := range step(b, st.j, c) {
					next := state{i, j}
					if !seen[next] {
						seen[next] = true
						queue = append(queue, next)
					}
				}
			}
		}
	}
	return false
}

// commonBytes returns the bytes both atoms can consume at this position. When
// both are wildcards any byte works and they all lead to the same successor
// states, so one representative is enough.
func commonBytes(x, y atom) []byte {
	switch {
	case x.wild && y.wild:
		return []byte{'a'}
	case x.wild:
		return []byte{y.ch}
	case y.wild:
		return []byte{x.ch}
	case x.ch == y.ch:
		return []byte{x.ch}
	default:
		return nil
	}
}

// step returns the states reachable from i by consuming c. A wildcard either
// keeps consuming or ends here; either way it has taken at least one byte,
// which is what makes a wildcard match non-empty.
func step(as []atom, i int, c byte) []int {
	if as[i].wild {
		return []int{i, i + 1}
	}
	if as[i].ch != c {
		return nil
	}
	return []int{i + 1}
}
