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
//
// CheckSyntax does not reconcile: the preview it builds leaves regions the
// selection does not touch as they are. Replay uses it for that reason; grow
// uses CheckSyntaxReconciled.
func CheckSyntax(recordID string, files []resolve.File, invocations []recording.Invocation, variables map[string]string, output string) []Violation {
	return checkSyntax(recordID, files, invocations, variables, output, false)
}

// CheckSyntaxReconciled is CheckSyntax with the preview reconciled the way
// Phase 7 of grow reconciles it: regions this record owns that the selection
// does not fill are cleared in the preview, so the checked content is what
// grow would write. Grow uses it; replay uses CheckSyntax.
func CheckSyntaxReconciled(recordID string, files []resolve.File, invocations []recording.Invocation, variables map[string]string, output string) []Violation {
	return checkSyntax(recordID, files, invocations, variables, output, true)
}

func checkSyntax(recordID string, files []resolve.File, invocations []recording.Invocation, variables map[string]string, output string, reconcile bool) []Violation {
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
	opts := inject.Options{Output: output, Unwritten: unwritten}
	if reconcile {
		var targets []inject.Target
		for _, f := range files {
			if f.Unmanaged || f.Package == nil {
				continue
			}
			targets = append(targets, inject.Target{Path: f.Path, Package: f.Package})
		}
		opts.Reconcile = []inject.Reconcile{{RecordID: recordID, Targets: targets}}
	}
	buffers, err := inject.Preview(expanded, opts)
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
		err := checker.Check(syntax.File{Path: path, Content: buffers[path], Siblings: syntax.Siblings(output, path, buffers)})
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
