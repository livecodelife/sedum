package record

import (
	"regexp"
	"strings"

	"github.com/livecodelife/sedum/internal/pathpat"
)

// Scope entries are patterns or paths, and telling them apart is what decides
// whether an entry creates a file (prov-2026-e8671c88).
//
// The grammar itself lives in internal/pathpat, because a generator package's
// unmanaged list reads the same notation. An author who has learned to write a
// scope entry has learned to write that one, and there is one implementation
// rather than two that drift.

// isPattern reports whether an affected_scope entry describes a set rather than
// naming a file. A route entry is the exception: a Next.js route folder written
// with brackets names the file it spells (prov-2026-e8671c88).
func isPattern(entry string) bool { return pathpat.IsPattern(entry) && !isRouteEntry(entry) }

// routeBracket is a whole path segment written as [name], [...name] or
// [[...name]], where name is one or more letters, digits, '-' or '_'.
var routeBracket = regexp.MustCompile(`^(?:\[[A-Za-z0-9_-]+\]|\[\.\.\.[A-Za-z0-9_-]+\]|\[\[\.\.\.[A-Za-z0-9_-]+\]\])$`)

// isRouteEntry reports whether an entry has no '*', '?' or trailing slash, has
// at least one segment containing '[', and every segment containing '[' is
// exactly a route bracket form. Any other bracket stays a character class.
func isRouteEntry(entry string) bool {
	if strings.ContainsAny(entry, "*?") || strings.HasSuffix(entry, "/") {
		return false
	}
	bracketed := false
	for _, seg := range strings.Split(entry, "/") {
		if !strings.Contains(seg, "[") {
			continue
		}
		if !routeBracket.MatchString(seg) {
			return false
		}
		bracketed = true
	}
	return bracketed
}

func checkPattern(entry string) error { return pathpat.Check(entry) }

// matchScope reports whether a scope entry covers a path. A literal entry
// matches only itself, which is what makes the same function serve
// forbidden_scope entries of either kind.
func matchScope(entry, target string) bool { return pathpat.Match(entry, target) }
