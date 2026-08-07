// Package filetmpl matches a path against a set of file-template path patterns
// and ranks the matches by specificity.
//
// A generator package's files/ directory is a literal mirror of the target
// project's structure, and each template's own path is its pattern:
//
//	app/controllers/{name}_controller.rb
//	app/models/{name}.rb
//	config/initializers/{name}.rb
//	_default.rb
//
// Matching app/controllers/users_controller.rb against that set selects the
// first and binds name=users. Captures are sub-segment, so {name}_controller.rb
// is literal text around a capture rather than a segment that happens to be
// one.
//
// The package is pure. It takes pattern strings and a path and returns an
// answer; it never reads a directory, stats a file, or writes anything.
// Walking files/ to collect the pattern set belongs to whoever owns the
// filesystem, which keeps the ranking rules testable as a table of strings.
//
// It also has no knowledge of _default templates, empty-file fallback, or
// logging. A path that matches nothing is a reported outcome, not an error, and
// the caller decides what to do about it.
package filetmpl

import (
	"fmt"
	"sort"
)

// Result is the template a path resolved to and the values its captures bound.
type Result struct {
	// Pattern is the winning pattern exactly as it was supplied, so a caller
	// can map it back to the template file it came from.
	Pattern string
	// Captures is empty, never nil, for a pattern with no captures.
	Captures map[string]string
}

// Conflict is a pair of patterns that tie: they are equally specific and some
// path matches both. Names are ordered so a diagnostic reads the same on every
// run.
type Conflict struct {
	A, B string
}

// TieError reports that a path matched two patterns equally well. Resolving it
// silently would make which boilerplate a file receives depend on the order the
// caller happened to walk files/ in, so it is reported instead.
type TieError struct {
	Path string
	A, B string
}

func (e *TieError) Error() string {
	return fmt.Sprintf("path %s matches file templates %s and %s equally; neither is more specific, so which one wins would depend on directory order", e.Path, e.A, e.B)
}

// Set is a parsed collection of file-template path patterns. Parsing once and
// matching many paths against it is the expected use, since a run resolves
// every authorized path against the same package.
type Set struct {
	patterns []*pattern
}

// NewSet parses every pattern, rejecting the whole set if any one of them is
// malformed. A pattern set is wholly valid or unusable: silently dropping the
// broken one would mean a typo in a template's path removes that template from
// consideration without saying so.
func NewSet(patterns []string) (*Set, error) {
	s := &Set{patterns: make([]*pattern, 0, len(patterns))}
	for _, raw := range patterns {
		p, err := parse(raw)
		if err != nil {
			return nil, err
		}
		s.patterns = append(s.patterns, p)
	}
	return s, nil
}

// Match selects the most specific pattern matching path and binds its captures.
//
// Three outcomes: a match, a tie reported as a *TieError, or no match, which is
// reported as ok=false and is not an error. Not every path needs a template.
//
// The winner does not depend on the order patterns were supplied in.
func (s *Set) Match(target string) (Result, bool, error) {
	var (
		best     *pattern
		captures map[string]string
		tied     *pattern
	)

	for _, p := range s.patterns {
		caps, ok := p.match(target)
		if !ok {
			continue
		}
		if best == nil {
			best, captures, tied = p, caps, nil
			continue
		}
		switch c := compare(p, best); {
		case c < 0:
			// A strictly better pattern clears any tie behind it.
			best, captures, tied = p, caps, nil
		case c == 0:
			tied = p
		}
	}

	switch {
	case best == nil:
		return Result{}, false, nil
	case tied != nil:
		a, b := order(best.raw, tied.raw)
		return Result{}, false, &TieError{Path: target, A: a, B: b}
	default:
		return Result{Pattern: best.raw, Captures: captures}, true, nil
	}
}

// Conflicts reports every pair of patterns that tie, without reference to any
// path. It is the check a package loader runs over files/ before any path
// exists: a tie means some file's boilerplate would be chosen by directory
// order, which is a defect in the package rather than in the run.
//
// Equal specificity alone is not a tie. app/{name}.rb and lib/{name}.rb are
// equally specific and can never both match, so they are not reported.
//
// Pairs are reported once, in a stable order, regardless of input order.
func (s *Set) Conflicts() []Conflict {
	var out []Conflict
	for i := range s.patterns {
		for j := i + 1; j < len(s.patterns); j++ {
			a, b := s.patterns[i], s.patterns[j]
			if compare(a, b) != 0 || !overlap(a, b) {
				continue
			}
			x, y := order(a.raw, b.raw)
			out = append(out, Conflict{A: x, B: y})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].A != out[j].A {
			return out[i].A < out[j].A
		}
		return out[i].B < out[j].B
	})
	return out
}

func order(a, b string) (string, string) {
	if b < a {
		return b, a
	}
	return a, b
}
