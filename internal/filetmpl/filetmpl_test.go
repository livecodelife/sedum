package filetmpl

import (
	"errors"
	"reflect"
	"sort"
	"testing"
)

// The ranking rules are contract, not convenience: PRD.md's "Match specificity"
// section is reproduced here as a table so that a change to the comparison
// fails a test rather than silently changing which boilerplate a file receives.
//
// Every case in this file is a set of pattern strings and a path. Nothing here
// touches the filesystem: collecting the pattern set belongs to whoever owns
// files/ on disk, not to this package.

func newSet(t *testing.T, patterns ...string) *Set {
	t.Helper()
	s, err := NewSet(patterns)
	if err != nil {
		t.Fatalf("NewSet(%q): %v", patterns, err)
	}
	return s
}

func TestMatchBindsCaptures(t *testing.T) {
	// The Rails tree from PRD.md, plus a Go one, so that no case depends on a
	// single language's shape.
	rails := []string{
		"app/controllers/{name}_controller.rb",
		"app/models/{name}.rb",
		"config/initializers/{name}.rb",
		"_default.rb",
	}

	tests := []struct {
		name     string
		patterns []string
		path     string
		want     string
		captures map[string]string
	}{
		{
			name:     "controller wins over model on the leftmost differing segment",
			patterns: rails,
			path:     "app/controllers/users_controller.rb",
			want:     "app/controllers/{name}_controller.rb",
			captures: map[string]string{"name": "users"},
		},
		{
			name:     "model matches under its own directory",
			patterns: rails,
			path:     "app/models/user.rb",
			want:     "app/models/{name}.rb",
			captures: map[string]string{"name": "user"},
		},
		{
			name:     "single-segment pattern matches only a single-segment path",
			patterns: rails,
			path:     "_default.rb",
			want:     "_default.rb",
			captures: map[string]string{},
		},
		{
			name:     "capture is sub-segment, not whole segment",
			patterns: []string{"internal/{dir}/{pkg}_test.go"},
			path:     "internal/filetmpl/filetmpl_test.go",
			want:     "internal/{dir}/{pkg}_test.go",
			captures: map[string]string{"dir": "filetmpl", "pkg": "filetmpl"},
		},
		{
			name:     "multiple captures in one segment",
			patterns: []string{"src/{area}_{kind}.go"},
			path:     "src/user_handler.go",
			want:     "src/{area}_{kind}.go",
			captures: map[string]string{"area": "user", "kind": "handler"},
		},
		{
			name:     "literals bind leftmost, so a capture takes the shortest text that still completes the match",
			patterns: []string{"{name}_handler.go"},
			path:     "user_handler.go",
			want:     "{name}_handler.go",
			captures: map[string]string{"name": "user"},
		},
		{
			name:     "a capture stretches when the leftmost literal placement cannot complete the match",
			patterns: []string{"{name}_handler.go"},
			path:     "user_handler_handler.go",
			want:     "{name}_handler.go",
			captures: map[string]string{"name": "user_handler"},
		},
		{
			name:     "glob matches but binds nothing",
			patterns: []string{"vendor/*/LICENSE"},
			path:     "vendor/cobra/LICENSE",
			want:     "vendor/*/LICENSE",
			captures: map[string]string{},
		},
		{
			name:     "a wholly literal pattern needs no captures",
			patterns: []string{"Makefile", "{name}"},
			path:     "Makefile",
			want:     "Makefile",
			captures: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := newSet(t, tt.patterns...).Match(tt.path)
			if err != nil {
				t.Fatalf("Match(%q): unexpected error: %v", tt.path, err)
			}
			if !ok {
				t.Fatalf("Match(%q): no match, want %q", tt.path, tt.want)
			}
			if got.Pattern != tt.want {
				t.Errorf("Match(%q) selected %q, want %q", tt.path, got.Pattern, tt.want)
			}
			if !reflect.DeepEqual(got.Captures, tt.captures) {
				t.Errorf("Match(%q) bound %v, want %v", tt.path, got.Captures, tt.captures)
			}
		})
	}
}

// Each case is a pair differing only in the rule under test, so a regression
// names the rule it broke rather than "some pattern won".
func TestSpecificityRules(t *testing.T) {
	tests := []struct {
		rule   string
		winner string
		loser  string
		path   string
	}{
		{
			rule:   "literal beats capture",
			winner: "app/controllers/users_controller.rb",
			loser:  "app/{dir}/users_controller.rb",
			path:   "app/controllers/users_controller.rb",
		},
		{
			rule:   "literal beats capture, sub-segment",
			winner: "app/users_{kind}.rb",
			loser:  "app/{name}_{kind}.rb",
			path:   "app/users_controller.rb",
		},
		{
			rule:   "capture beats glob",
			winner: "app/{name}/index.rb",
			loser:  "app/*/index.rb",
			path:   "app/models/index.rb",
		},
		{
			rule:   "capture beats glob, sub-segment",
			winner: "app/{name}_controller.rb",
			loser:  "app/*_controller.rb",
			path:   "app/users_controller.rb",
		},
		{
			rule:   "longer literal prefix beats shorter",
			winner: "app/user_{kind}.rb",
			loser:  "app/us{rest}.rb",
			path:   "app/user_controller.rb",
		},
		{
			rule:   "longer trailing literal beats shorter at the same position",
			winner: "{name}_controller.rb",
			loser:  "{name}.rb",
			path:   "users_controller.rb",
		},
		{
			rule:   "leftmost segment decides, even when the loser is more specific later",
			winner: "app/{a}/{b}.rb",
			loser:  "{a}/models/user.rb",
			path:   "app/models/user.rb",
		},
		{
			rule:   "an extra literal piece beats running out of pieces",
			winner: "{name}_x.rb",
			loser:  "{name}.rb",
			path:   "user_x.rb",
		},
		{
			rule:   "running out of pieces beats a trailing glob",
			winner: "src/{name}_x",
			loser:  "src/{name}_x*",
			path:   "src/a_x_x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.rule, func(t *testing.T) {
			// Both orders: ranking must not depend on the order patterns
			// were supplied in.
			for _, patterns := range [][]string{{tt.winner, tt.loser}, {tt.loser, tt.winner}} {
				got, ok, err := newSet(t, patterns...).Match(tt.path)
				if err != nil {
					t.Fatalf("Match(%q) with %q: unexpected error: %v", tt.path, patterns, err)
				}
				if !ok {
					t.Fatalf("Match(%q) with %q: no match", tt.path, patterns)
				}
				if got.Pattern != tt.winner {
					t.Errorf("Match(%q) with %q selected %q, want %q", tt.path, patterns, got.Pattern, tt.winner)
				}
			}

			// Both patterns must genuinely match the path on their own,
			// or the case proves nothing about ranking.
			for _, p := range []string{tt.winner, tt.loser} {
				if _, ok, err := newSet(t, p).Match(tt.path); err != nil || !ok {
					t.Fatalf("case is not about ranking: %q does not match %q on its own (ok=%v err=%v)", p, tt.path, ok, err)
				}
			}
		})
	}
}

func TestOrderIndependence(t *testing.T) {
	patterns := []string{
		"app/controllers/{name}_controller.rb",
		"app/{dir}/{name}_controller.rb",
		"app/*/{name}_controller.rb",
		"{a}/{b}/{c}.rb",
	}
	const path = "app/controllers/users_controller.rb"
	const want = "app/controllers/{name}_controller.rb"

	for _, perm := range permutations(patterns) {
		got, ok, err := newSet(t, perm...).Match(path)
		if err != nil {
			t.Fatalf("Match with order %q: unexpected error: %v", perm, err)
		}
		if !ok || got.Pattern != want {
			t.Fatalf("Match with order %q selected %q (ok=%v), want %q", perm, got.Pattern, ok, want)
		}
	}
}

func TestTieIsReportedNotResolved(t *testing.T) {
	patterns := []string{"app/{dir}/user.rb", "app/{name}/user.rb"}
	const path = "app/models/user.rb"

	for _, order := range [][]string{patterns, {patterns[1], patterns[0]}} {
		_, ok, err := newSet(t, order...).Match(path)
		if ok {
			t.Fatalf("Match(%q) with %q resolved a tie instead of reporting it", path, order)
		}
		var tie *TieError
		if !errors.As(err, &tie) {
			t.Fatalf("Match(%q) with %q returned %v, want a *TieError", path, order, err)
		}
		if tie.Path != path {
			t.Errorf("TieError.Path = %q, want %q", tie.Path, path)
		}
		// Named in a stable order so the diagnostic does not depend on
		// which order the caller walked files/ in.
		if tie.A != "app/{dir}/user.rb" || tie.B != "app/{name}/user.rb" {
			t.Errorf("TieError names %q and %q, want both tied patterns in sorted order", tie.A, tie.B)
		}
		for _, name := range []string{tie.A, tie.B} {
			if !containsString(tie.Error(), name) {
				t.Errorf("TieError.Error() = %q, does not name %q", tie.Error(), name)
			}
		}
	}
}

func TestNoMatchIsNotAnError(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		path     string
	}{
		{"no pattern of that shape", []string{"app/models/{name}.rb"}, "lib/tasks/build.rake"},
		{"segment count differs", []string{"app/models/{name}.rb"}, "app/models/admin/user.rb"},
		{"literal does not line up", []string{"app/{name}_controller.rb"}, "app/user_model.rb"},
		{"capture will not match empty text", []string{"app/{name}_controller.rb"}, "app/_controller.rb"},
		{"glob will not match empty text", []string{"app/*_controller.rb"}, "app/_controller.rb"},
		{"empty pattern set", nil, "app/models/user.rb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := newSet(t, tt.patterns...).Match(tt.path)
			if err != nil {
				t.Fatalf("Match(%q): returned an error for a non-match: %v", tt.path, err)
			}
			if ok {
				t.Fatalf("Match(%q): selected %q, want no match", tt.path, got.Pattern)
			}
		})
	}
}

func TestSeparatorsAreNormalized(t *testing.T) {
	// A package authored on Windows must match identically elsewhere, and a
	// caller that hands over a host-separator path must get the same answer.
	for _, pattern := range []string{`app\controllers\{name}_controller.rb`, "./app/controllers/{name}_controller.rb"} {
		for _, path := range []string{`app\controllers\users_controller.rb`, "app/controllers/users_controller.rb"} {
			got, ok, err := newSet(t, pattern).Match(path)
			if err != nil || !ok {
				t.Fatalf("Match(%q) against %q: ok=%v err=%v", path, pattern, ok, err)
			}
			if got.Captures["name"] != "users" {
				t.Errorf("Match(%q) against %q bound name=%q, want users", path, pattern, got.Captures["name"])
			}
		}
	}
}

func TestPatternErrors(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		names   []string // substrings the diagnostic must contain
	}{
		{"unclosed capture", "app/{name.rb", []string{"app/{name.rb", "{"}},
		{"unmatched close", "app/name}.rb", []string{"app/name}.rb", "}"}},
		{"capture with no name", "app/{}.rb", []string{"app/{}.rb"}},
		{"duplicate capture name", "app/{name}/{name}.rb", []string{"name"}},
		{"adjacent captures are ambiguous", "app/{a}{b}.rb", []string{"app/{a}{b}.rb"}},
		{"capture after glob is ambiguous", "app/*{b}.rb", []string{"app/*{b}.rb"}},
		{"adjacent globs are ambiguous", "app/**.rb", []string{"app/**.rb"}},
		{"empty pattern", "", nil},
		{"pattern that normalizes to nothing", ".", nil},
		{"bare separator", "/", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSet([]string{tt.pattern})
			if err == nil {
				t.Fatalf("NewSet(%q): accepted a malformed pattern", tt.pattern)
			}
			for _, want := range tt.names {
				if !containsString(err.Error(), want) {
					t.Errorf("NewSet(%q) error = %q, does not name %q", tt.pattern, err, want)
				}
			}
		})
	}
}

// Conflicts is the load-time half: it answers "do two of these templates tie?"
// without a path, which is what a package loader has at the time it must
// reject the package.
func TestConflicts(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		want     []Conflict
	}{
		{
			name:     "the PRD's Rails tree is conflict-free",
			patterns: []string{"app/controllers/{name}_controller.rb", "app/models/{name}.rb", "config/initializers/{name}.rb", "_default.rb"},
		},
		{
			name:     "two captures over the same shape tie",
			patterns: []string{"app/{dir}/user.rb", "app/{name}/user.rb"},
			want:     []Conflict{{A: "app/{dir}/user.rb", B: "app/{name}/user.rb"}},
		},
		{
			name:     "equal specificity over disjoint literals is not a tie",
			patterns: []string{"app/{name}.rb", "lib/{name}.rb"},
		},
		{
			name:     "differing specificity is not a tie even though both can match",
			patterns: []string{"app/{name}_controller.rb", "app/*_controller.rb"},
		},
		{
			name:     "differing segment counts cannot tie",
			patterns: []string{"{name}.rb", "app/{name}.rb"},
		},
		{
			name:     "overlapping sub-segment captures tie",
			patterns: []string{"src/{a}_{b}.go", "src/{x}_{y}.go"},
			want:     []Conflict{{A: "src/{a}_{b}.go", B: "src/{x}_{y}.go"}},
		},
		{
			name:     "sub-segment literals that cannot line up are not a tie",
			patterns: []string{"src/{a}_handler.go", "src/{a}_model.go"},
		},
		{
			name:     "a tie is reported once per pair, not once per ordering",
			patterns: []string{"{a}/x.rb", "{b}/x.rb", "{c}/x.rb"},
			want: []Conflict{
				{A: "{a}/x.rb", B: "{b}/x.rb"},
				{A: "{a}/x.rb", B: "{c}/x.rb"},
				{A: "{b}/x.rb", B: "{c}/x.rb"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSet(t, tt.patterns...).Conflicts()
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Conflicts() = %v, want %v", got, tt.want)
			}

			// A reported conflict must be a real one: some path both
			// patterns match and neither wins.
			for _, c := range got {
				if compareRaw(t, c.A, c.B) != 0 {
					t.Errorf("Conflicts() reported %v, but the two do not tie under the ranking", c)
				}
			}
		})
	}
}

// Conflicts must not depend on the order patterns were supplied in either,
// since a caller walking files/ has no control over directory iteration order.
func TestConflictsOrderIndependent(t *testing.T) {
	patterns := []string{"{a}/x.rb", "app/{n}.rb", "{b}/x.rb"}
	want := []Conflict{{A: "{a}/x.rb", B: "{b}/x.rb"}}

	for _, perm := range permutations(patterns) {
		if got := newSet(t, perm...).Conflicts(); !reflect.DeepEqual(got, want) {
			t.Errorf("Conflicts() with order %q = %v, want %v", perm, got, want)
		}
	}
}

func compareRaw(t *testing.T, a, b string) int {
	t.Helper()
	pa, err := parse(a)
	if err != nil {
		t.Fatalf("parse(%q): %v", a, err)
	}
	pb, err := parse(b)
	if err != nil {
		t.Fatalf("parse(%q): %v", b, err)
	}
	return compare(pa, pb)
}

func containsString(haystack, needle string) bool {
	return len(needle) == 0 || len(haystack) >= len(needle) && indexOf(haystack, needle) >= 0
}

func indexOf(haystack, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

func permutations(in []string) [][]string {
	if len(in) <= 1 {
		return [][]string{append([]string(nil), in...)}
	}
	var out [][]string
	for i := range in {
		rest := make([]string, 0, len(in)-1)
		rest = append(rest, in[:i]...)
		rest = append(rest, in[i+1:]...)
		for _, p := range permutations(rest) {
			out = append(out, append([]string{in[i]}, p...))
		}
	}
	return out
}

func TestConflictOrderingIsStable(t *testing.T) {
	got := newSet(t, "{b}/x.rb", "{a}/x.rb").Conflicts()
	want := []Conflict{{A: "{a}/x.rb", B: "{b}/x.rb"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Conflicts() = %v, want %v", got, want)
	}
	if !sort.SliceIsSorted(got, func(i, j int) bool {
		return got[i].A < got[j].A || (got[i].A == got[j].A && got[i].B < got[j].B)
	}) {
		t.Errorf("Conflicts() = %v, want sorted output", got)
	}
}
