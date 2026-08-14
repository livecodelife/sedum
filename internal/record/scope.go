package record

import (
	"github.com/calebcowen/sedum/internal/pathpat"
)

// Scope entries are patterns or paths, and telling them apart is what decides
// whether an entry creates a file (prov-2026-e8671c88).
//
// The grammar itself lives in internal/pathpat, because a generator package's
// unmanaged list reads the same notation. An author who has learned to write a
// scope entry has learned to write that one, and there is one implementation
// rather than two that drift.

func isPattern(entry string) bool { return pathpat.IsPattern(entry) }

func checkPattern(entry string) error { return pathpat.Check(entry) }

// matchScope reports whether a scope entry covers a path. A literal entry
// matches only itself, which is what makes the same function serve
// forbidden_scope entries of either kind.
func matchScope(entry, target string) bool { return pathpat.Match(entry, target) }
