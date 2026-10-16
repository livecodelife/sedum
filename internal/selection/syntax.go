package selection

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/livecodelife/sedum/internal/expand"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/resolve"
	"github.com/livecodelife/sedum/internal/syntax"
)

// RuleSyntax is a selection whose files would not pass their package's
// built-in syntax checker (prov-2026-5403d162).
const RuleSyntax = "syntax"

// CheckSyntax applies a validated selection to its files in memory and runs
// every file whose package declares a syntax: key through that built-in
// checker. Each failure is a violation, so in grow it is fed back to the model
// within the retry budget, and in replay it is terminal - the same treatment
// as any other rule an answer breaks.
//
// A file that parses but is not in its formatter's form is reported with the
// formatted form, because "bind exactly this" is an instruction the model can
// follow, where "your code is not gofmt-clean" is not.
//
// It reads each file from output when it is on disk, and otherwise from what
// Phase 3 rendered for it. A selection that cannot be applied at all - an
// ownership conflict, a missing anchor - is not reported here; Phase 7 reports
// it with its own diagnostic.
func CheckSyntax(recordID string, files []resolve.File, invocations []recording.Invocation, variables map[string]string, output string) []Violation {
	checked := false
	unwritten := map[string]string{}
	for _, f := range files {
		if f.Package != nil && f.Package.Syntax != "" {
			checked = true
		}
		if _, err := os.Stat(filepath.Join(output, filepath.FromSlash(f.Path))); err != nil {
			unwritten[f.Path] = f.Rendered
		}
	}
	if !checked || len(invocations) == 0 {
		return nil
	}

	expanded, err := expand.Expand(recordID, files, invocations, variables)
	if err != nil {
		return nil
	}
	buffers, err := inject.Preview(expanded, inject.Options{Output: output, Unwritten: unwritten})
	if err != nil {
		return nil
	}

	packages := map[string]string{}
	for _, f := range files {
		if f.Package != nil {
			packages[f.Path] = f.Package.Syntax
		}
	}
	paths := make([]string, 0, len(buffers))
	for path := range buffers {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var out []Violation
	for _, path := range paths {
		name := packages[path]
		checker, ok := syntax.Lookup(name)
		if name == "" || !ok {
			continue
		}
		err := checker.Check(path, buffers[path])
		if err == nil {
			continue
		}
		detail := err.Error()
		var problem *syntax.Problem
		if errors.As(err, &problem) && problem.Formatted != "" {
			detail += ". Formatted, the whole file would read as follows; bind each code value exactly as it appears here:\n\n" + problem.Formatted
		}
		out = append(out, Violation{Rule: RuleSyntax, Detail: detail})
	}
	return out
}
