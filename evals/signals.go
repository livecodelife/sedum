package evals

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/calebcowen/sedum/internal/expand"
	"github.com/calebcowen/sedum/internal/inject"
	"github.com/calebcowen/sedum/internal/pipeline"
	"github.com/calebcowen/sedum/internal/recording"
	resolvepkg "github.com/calebcowen/sedum/internal/resolve"
)

// AnchorFill is how much of the work a run's own files declared was accounted
// for by what the model selected.
//
// It is derived rather than declared, which is the whole of its value. A file
// template plants the markers an action's anchor targets, so a created file
// states on disk what work it expects, and the fraction of that work a
// selection covered comes from the package and the record rather than from
// anybody's expectation. The hand-authored action counts can be fitted to the
// answer they score; this cannot be (prov-2026-d61010a4).
//
// It does not replace those counts and is not replaced by them. An anchor
// filled by the wrong action counts as filled here and is caught there, and an
// expectation names actions this has no opinion about.
type AnchorFill struct {
	// Planted is the anchors this run created that something it loaded could
	// have written into.
	Planted int
	// Filled is the subset a selection accounted for.
	Filled int
}

// Rate is the fraction of fillable planted anchors a selection accounted for,
// and false when the run planted none.
//
// Absent rather than zero, for the same reason every other rate here is: a run
// with no anchors to fill and a run that filled none of them are different
// observations and a bare 0.0 tells them apart from neither.
func (a AnchorFill) Rate() (float64, bool) {
	if a.Planted == 0 {
		return 0, false
	}
	return float64(a.Filled) / float64(a.Planted), true
}

// Add accumulates another selection's anchors into this one.
func (a *AnchorFill) Add(b AnchorFill) {
	a.Planted += b.Planted
	a.Filled += b.Filled
}

// anchorFill counts one record's fillable planted anchors and how many its
// invocations accounted for.
//
// The denominator is the *fillable* planted anchors rather than every planted
// one, and the difference is not a rounding detail. expand.Unfilled already
// drops anchors nothing in the run targets - init.sql plants "extensions" and
// no action writes there - because they are not omissions the model could have
// acted on. Counting them in the denominator anyway would give the package a
// ceiling below 100% and report it as the model falling short of one, which is
// exactly the attribution error this signal exists to prevent.
//
// So the subtraction is done against the same set on both sides: Fillable is
// computed from the files passed here, which is the set Unfilled computes it
// from too.
func anchorFill(recordID string, files []resolvepkg.File, invocations []recording.Invocation, variables map[string]string) AnchorFill {
	fillable := expand.Fillable(files)

	var fill AnchorFill
	for _, a := range expand.Planted(files) {
		if fillable[a.Marker] {
			fill.Planted++
		}
	}
	fill.Filled = fill.Planted - len(expand.Unfilled(recordID, files, invocations, variables))
	return fill
}

// fillOf sums the anchor fill across every record a run selected for.
//
// Per record rather than per run, because Unfilled is answered per record: a
// record's files are its own, and an anchor one record planted is not work
// another record's selection failed to do.
func fillOf(result *pipeline.Result) AnchorFill {
	var fill AnchorFill
	for _, s := range result.Selections {
		fill.Add(anchorFill(s.RecordID, s.Files, s.Invocations, result.Variables))
	}
	return fill
}

// Check is the command a case runs over what Sedum wrote, keyed by the file
// extension it can read.
//
// Keyed by extension because a package writes more than one kind of file - the
// rails package claims .rb and .yml - and ruby -c over a YAML file would fail
// for a reason that has nothing to do with what the model chose. The harness
// matches the suffix and runs the command; it does not interpret either, and it
// has no table of languages to consult. A framework is added by writing a case,
// never by editing this package (prov-2026-d61010a4).
//
// An extension no entry covers is not checked and is not counted, which is why
// SyntaxCheck reports what it looked at rather than only what passed.
type Check map[string][]string

// SyntaxCheck is what the target's own parser said about the files a sample
// wrote.
//
// It is not correctness and must never be reported as it. It catches malformed
// output, not wrong output: params.permit(title, completed) parses as valid Ruby
// and raises NameError the first time it runs, so the failure that motivated
// this signal's sibling passes it. Anything read as a claim about behaviour is
// claiming the rung above (prov-2026-c5697387).
type SyntaxCheck struct {
	// Checked is the files an entry in the case's Check covered.
	Checked int
	// Parsed is the subset the command accepted.
	Parsed int
	// Failures is one line per file the command rejected, in path order.
	Failures []string
	// Unavailable names the commands that could not be run at all, because
	// the binary is not on this machine.
	//
	// Separated from a rejection deliberately. A laptop without ruby would
	// otherwise report every Ruby file as malformed and produce a zero rate
	// that looks like a catastrophic finding about the model, which is the
	// worst failure mode a measurement can have: wrong, alarming, and shaped
	// exactly like a real result.
	Unavailable []string
}

// Rate is the fraction of checked files the target's parser accepted, and false
// when nothing was checked - because the case declares no command, because no
// file matched one, or because the command is not installed here.
func (s SyntaxCheck) Rate() (float64, bool) {
	if s.Checked == 0 {
		return 0, false
	}
	return float64(s.Parsed) / float64(s.Checked), true
}

// syntaxOf runs the case's check commands over what the run wrote under output.
//
// The files are the ones Sedum created and wrote into, read from disk after
// Phase 7 rather than from anything in memory: the point is what a target's
// parser makes of the bytes that were actually produced.
func syntaxOf(ctx context.Context, check Check, output string, result *pipeline.Result) SyntaxCheck {
	var out SyntaxCheck
	if len(check) == 0 {
		return out
	}

	missing := map[string]bool{}
	for _, path := range writtenPaths(result) {
		cmd, ok := check[filepath.Ext(path)]
		if !ok || len(cmd) == 0 {
			continue
		}
		if _, err := exec.LookPath(cmd[0]); err != nil {
			if !missing[cmd[0]] {
				missing[cmd[0]] = true
				out.Unavailable = append(out.Unavailable, cmd[0])
			}
			continue
		}

		out.Checked++
		args := append(append([]string{}, cmd[1:]...), path)
		c := exec.CommandContext(ctx, cmd[0], args...)
		c.Dir = output
		if combined, err := c.CombinedOutput(); err != nil {
			out.Failures = append(out.Failures, path+": "+firstLine(strings.TrimSpace(string(combined))))
			continue
		}
		out.Parsed++
	}
	sort.Strings(out.Unavailable)
	return out
}

// writtenPaths is every path this run created or injected into, deduplicated
// and in a stable order.
//
// Unmanaged paths are excluded: Sedum did not write them, so their contents are
// not evidence about what the model chose. A path that was already on disk is
// excluded for the same reason unless something was injected into it, in which
// case the bytes are partly Sedum's and worth parsing.
func writtenPaths(result *pipeline.Result) []string {
	seen := map[string]bool{}
	for _, f := range result.Files {
		if f.Unmanaged || f.Existed {
			continue
		}
		seen[f.Path] = true
	}
	for _, i := range result.Injections {
		if i.Skipped {
			continue
		}
		seen[i.Path] = true
	}

	out := make([]string, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// Idempotency is what applying a sample's own selection a second time did to
// the bytes the first application left.
//
// A rerun that differs from its input is an injection or a marker defect: the
// region's identity did not match, or it matched and was rewritten differently.
// M4 asserts this property on hand-built fixtures; here it is asserted against
// real packages and model-selected invocations, which is the combination
// fixtures cannot cover - every marker a real template plants, filled with
// arguments nobody chose in advance (prov-2026-d61010a4).
//
// It is a regression signal and is expected to find nothing. That is the point
// of having it rather than an argument against it.
type Idempotency struct {
	// Files is the files compared before and after the second application.
	Files int
	// Stable is the subset whose bytes were identical.
	Stable int
	// Differing names the files a rerun changed, in path order.
	Differing []string
	// Err is set when the second application could not be made at all. It is
	// a finding rather than a harness failure: the same invocations that
	// applied cleanly once must apply cleanly again.
	Err error
}

// Rate is the fraction of files a rerun left alone, and false when nothing was
// compared.
func (i Idempotency) Rate() (float64, bool) {
	if i.Files == 0 {
		return 0, false
	}
	return float64(i.Stable) / float64(i.Files), true
}

// idempotencyOf applies the run's own invocations a second time and compares
// the bytes.
//
// Only Phases 6 and 7 are re-run. Phase 4 is a model call and re-running it
// would be sampling twice rather than applying once twice, which is a different
// question and an expensive way to ask it. Phase 3 is skipped because the files
// are already there, which is exactly the state the property is about: the
// second application lands on a file that already carries the first's regions.
func idempotencyOf(output string, result *pipeline.Result) Idempotency {
	var out Idempotency

	paths := writtenPaths(result)
	before := make(map[string][]byte, len(paths))
	for _, p := range paths {
		content, err := os.ReadFile(filepath.Join(output, p))
		if err != nil {
			// A path the run reported writing and that is not on disk is a
			// defect worth surfacing rather than a file to skip.
			out.Err = err
			return out
		}
		before[p] = content
	}
	if len(before) == 0 {
		return out
	}

	var resolved []inject.Invocation
	for _, s := range result.Selections {
		expanded, err := expand.Expand(s.RecordID, s.Files, s.Invocations, result.Variables)
		if err != nil {
			out.Err = fmt.Errorf("record %s: %w", s.RecordID, err)
			return out
		}
		resolved = append(resolved, expanded...)
	}
	if _, err := inject.Apply(resolved, inject.Options{Output: output}); err != nil {
		out.Err = err
		return out
	}

	for _, p := range paths {
		after, err := os.ReadFile(filepath.Join(output, p))
		if err != nil {
			out.Err = err
			return out
		}
		out.Files++
		if bytes.Equal(before[p], after) {
			out.Stable++
			continue
		}
		out.Differing = append(out.Differing, p)
	}
	return out
}
