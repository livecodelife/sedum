// Package syntax is the closed registry of target-language checks
// (prov-2026-5403d162).
//
// Sedum's pipeline contains no language-specific knowledge. Where a target
// language has to be understood - generated Go must parse and be gofmt-clean,
// because model-written code must never land broken - the knowledge lives
// here and only here: named, opted into by a package's `syntax:` key, shipped
// with Sedum, and running no external program. The pipeline asks the registry
// by name and never branches on a language itself.
//
// The registry is closed because the alternative, a package-declared command,
// lets a generator package run an executable on every grow, and nothing from
// a copied or shared source may run automatically.
package syntax

import (
	"fmt"
	"go/format"
	"sort"
)

// Checker checks one whole file of its language.
type Checker interface {
	// Check returns nil when content is acceptable, and otherwise a *Problem
	// naming path.
	Check(path, content string) error
}

// Problem is a file a checker rejected.
type Problem struct {
	Path   string
	Reason string
	// Formatted is the file as the language's formatter would write it, when
	// the file parses but is not in that form. Empty when it does not parse.
	// It is what a re-prompt shows the model, so it can bind exactly that.
	Formatted string
}

func (p *Problem) Error() string {
	return fmt.Sprintf("%s: %s", p.Path, p.Reason)
}

var registry = map[string]Checker{
	"go": goChecker{},
}

// Lookup returns the checker registered under name.
func Lookup(name string) (Checker, bool) {
	c, ok := registry[name]
	return c, ok
}

// Names returns every registered syntax, sorted.
func Names() []string {
	out := make([]string, 0, len(registry))
	for name := range registry {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// goChecker accepts Go source that parses and is exactly what gofmt would
// write. It uses go/format, the library gofmt itself is built on.
type goChecker struct{}

func (goChecker) Check(path, content string) error {
	formatted, err := format.Source([]byte(content))
	if err != nil {
		return &Problem{Path: path, Reason: "does not parse as Go: " + err.Error()}
	}
	if string(formatted) != content {
		return &Problem{Path: path, Reason: "is not gofmt-clean", Formatted: string(formatted)}
	}
	return nil
}
