package record

import (
	"fmt"
	"path"
	"strings"
)

// Scope entries are patterns or paths, and telling them apart is what decides
// whether an entry creates a file (prov-2026-e8671c88).
//
// The pattern vocabulary is the small one every scope notation converges on:
//
//	**  any number of segments, including none
//	*   any run of characters within one segment
//	?   one character within one segment
//	[…] one character from a set, as path.Match reads it
//
// Segment-by-segment matching over path.Match is what keeps * from crossing a
// separator, which is the whole difference between config/*.yml and
// config/**.yml.

// isPattern reports whether a scope entry authorizes a set rather than naming a
// file. A trailing slash counts: lib/tasks/ describes a subtree, and no file can
// be called that.
func isPattern(entry string) bool {
	return strings.HasSuffix(entry, "/") || strings.ContainsAny(entry, "*?[")
}

// checkPattern rejects a pattern path.Match cannot read, such as one with an
// unclosed character class. Left alone it would match nothing at all, so a
// forbidden_scope entry meant to protect a directory would silently protect
// nothing.
func checkPattern(entry string) error {
	for _, seg := range segments(entry) {
		if seg == "**" {
			continue
		}
		if _, err := path.Match(seg, "x"); err != nil {
			return fmt.Errorf("%q is not a pattern Sedum can read: %w", seg, err)
		}
	}
	return nil
}

// matchScope reports whether a scope entry covers a path. A literal entry
// matches only itself, which is what makes the same function serve
// forbidden_scope entries of either kind.
func matchScope(entry, target string) bool {
	return matchSegments(segments(entry), segments(normalize(target)))
}

// segments splits a scope entry into the parts matching compares. A trailing
// slash is read as a subtree, so lib/tasks/ and lib/tasks/** are the same
// entry.
func segments(entry string) []string {
	entry = strings.TrimSpace(entry)
	subtree := strings.HasSuffix(entry, "/")

	out := strings.Split(strings.Trim(normalize(entry), "/"), "/")
	if subtree {
		out = append(out, "**")
	}
	return out
}

func matchSegments(pattern, target []string) bool {
	for len(pattern) > 0 {
		if pattern[0] == "**" {
			// ** consumes any number of segments including none, so every
			// split has to be tried before the rest of the pattern can be
			// called unmatched.
			if len(pattern) == 1 {
				return true
			}
			for i := 0; i <= len(target); i++ {
				if matchSegments(pattern[1:], target[i:]) {
					return true
				}
			}
			return false
		}
		if len(target) == 0 {
			return false
		}
		ok, err := path.Match(pattern[0], target[0])
		if err != nil || !ok {
			return false
		}
		pattern, target = pattern[1:], target[1:]
	}
	return len(target) == 0
}
