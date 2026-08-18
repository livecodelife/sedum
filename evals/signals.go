package evals

import (
	"github.com/calebcowen/sedum/internal/expand"
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
