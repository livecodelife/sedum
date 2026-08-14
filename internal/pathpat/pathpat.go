// Package pathpat is the path-pattern grammar Sedum reads wherever a
// declaration names a set of files rather than one file.
//
// It is the small vocabulary every scope notation converges on:
//
//	**  any number of segments, including none
//	*   any run of characters within one segment
//	?   one character within one segment
//	[…] one character from a set, as path.Match reads it
//
// Segment-by-segment matching over path.Match is what keeps * from crossing a
// separator, which is the whole difference between config/*.yml and
// config/**.yml.
//
// It lives on its own because two declarations now read it: a provenance
// record's affected_scope and forbidden_scope, and a generator package's
// unmanaged list. One grammar written once means an author who has learned to
// write a scope entry has learned to write an unmanaged entry, and there is no
// second implementation to drift.
//
// This is deliberately not the file-template pattern grammar. That one has
// captures because its whole job is binding values out of a path; this one
// answers yes or no and binds nothing.
package pathpat

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

// IsPattern reports whether an entry describes a set rather than naming one
// file. A trailing slash counts: lib/tasks/ describes a subtree, and no file can
// be called that.
func IsPattern(entry string) bool {
	return strings.HasSuffix(entry, "/") || strings.ContainsAny(entry, "*?[")
}

// Check rejects a pattern path.Match cannot read, such as one with an unclosed
// character class. Left alone it would match nothing at all, so an entry meant
// to describe a directory would silently describe nothing.
func Check(entry string) error {
	for _, seg := range Segments(entry) {
		if seg == "**" {
			continue
		}
		if _, err := path.Match(seg, "x"); err != nil {
			return fmt.Errorf("%q is not a pattern Sedum can read: %w", seg, err)
		}
	}
	return nil
}

// Match reports whether an entry covers a path. A literal entry matches only
// itself, which is what lets the same function serve entries of either kind.
func Match(entry, target string) bool {
	return matchSegments(Segments(entry), Segments(Normalize(target)))
}

// MatchAny reports whether any entry covers the path, returning the first that
// does so a diagnostic can name the declaration rather than the path.
func MatchAny(entries []string, target string) (string, bool) {
	for _, entry := range entries {
		if Match(entry, target) {
			return entry, true
		}
	}
	return "", false
}

// Normalize puts a pattern or a path into slash-separated, cleaned form so that
// a declaration authored on one platform reads identically on another.
func Normalize(p string) string {
	return path.Clean(filepath.ToSlash(p))
}

// Segments splits an entry into the parts matching compares. A trailing slash is
// read as a subtree, so lib/tasks/ and lib/tasks/** are the same entry.
func Segments(entry string) []string {
	entry = strings.TrimSpace(entry)
	subtree := strings.HasSuffix(entry, "/")

	out := strings.Split(strings.Trim(Normalize(entry), "/"), "/")
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
