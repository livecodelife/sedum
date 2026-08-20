package expand

import (
	"sort"

	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/resolve"
)

// Anchor is one marker in one file: the pair that identifies an injection site
// across a run.
type Anchor struct {
	Path   string
	Marker string
}

// Planted returns every marker anchor the run's files carry.
//
// A file template plants the markers an action's anchor targets, which is the
// whole basis of injection - and it means a created file states, on disk, what
// work it expects. Reading that back is not structure inference: these are
// markers Sedum wrote, which is the distinction the non-goal turns on.
//
// It lives here rather than in genpkg because the per-file mapping needs
// resolve.File and resolve imports genpkg. genpkg keeps MarkersIn, which is the
// primitive this is built from.
//
// Unmanaged paths are skipped. Sedum did not write them, so it has no claim
// about what they contain and no business reporting an anchor missing from one.
func Planted(files []resolve.File) []Anchor {
	var out []Anchor
	for _, f := range files {
		if f.Unmanaged || f.Package == nil || f.Rendered == "" {
			continue
		}
		for _, marker := range genpkg.MarkersIn(f.Package.CommentPrefix, f.Rendered) {
			out = append(out, Anchor{Path: f.Path, Marker: marker})
		}
	}
	sortAnchors(out)
	return out
}

// Filled returns the anchors a validated invocation list writes into.
//
// Expansion is what resolves a composite into children with a path and an action
// each, so a composite fills its children's anchors rather than one of its own.
// Only marker anchors participate: start_of_file, end_of_file and the match
// anchors name no region a file can be observed to be missing.
//
// An expansion error yields no anchors rather than an error of its own. This
// runs after validation has passed, so a failure here is a defect in a later
// phase and belongs to that phase's diagnostic - reporting it as a completeness
// problem would name the wrong thing.
func Filled(recordID string, files []resolve.File, invocations []recording.Invocation, variables map[string]string) []Anchor {
	expanded, err := Expand(recordID, files, invocations, variables)
	if err != nil {
		return nil
	}

	var out []Anchor
	for _, inv := range expanded {
		marker, ok := inv.Action.MarkerAnchor()
		if !ok {
			continue
		}
		out = append(out, Anchor{Path: inv.Path, Marker: marker})
	}
	sortAnchors(out)
	return out
}

// Fillable returns the markers something in this run could have injected at:
// every marker targeted by any action in any package the run's files resolved
// to.
//
// It is computed across the run's packages rather than per package, because
// that is the question being asked. genpkg's load-time check is package-scoped
// because it asks an authoring question - is this package internally coherent.
// This asks whether anything in the run can reach the anchor, and an action's
// injects_into renders to a path that resolves to a package by extension, so an
// action in one package may target a file another package owns. A per-package
// rule would drop anchors something else could have filled.
func Fillable(files []resolve.File) map[string]bool {
	out := map[string]bool{}
	for _, pkg := range Packages(files) {
		for _, action := range pkg.Actions {
			for _, marker := range action.TargetedMarkers() {
				out[marker] = true
			}
		}
	}
	return out
}

// Unfilled is Planted minus Filled, minus the anchors nothing in the run could
// have filled: the anchors this run created, something it loaded can write
// into, and nothing it selected does.
//
// It is not an error and must not become one. A template may plant a region for
// an action a later record will fill, or for one this change legitimately does
// not need, and a run that refused to proceed would make every such template a
// liability. What it is good for is saying so - once, to the model, and to
// anyone measuring how complete a selection was.
//
// An anchor no action targets is dropped rather than reported, because it is
// not an omission the model can act on: asking about it spends a call on a
// question with no available answer, on every run, forever
// (prov-2026-206fa618). The subtraction is silent. Phase 0 already names every
// such marker at load, by name and with the reason, and repeating it here would
// charge an authoring warning by the run.
func Unfilled(recordID string, files []resolve.File, invocations []recording.Invocation, variables map[string]string) []Anchor {
	filled := map[Anchor]bool{}
	for _, a := range Filled(recordID, files, invocations, variables) {
		filled[a] = true
	}
	fillable := Fillable(files)

	var out []Anchor
	for _, a := range Planted(files) {
		if filled[a] || !fillable[a.Marker] {
			continue
		}
		out = append(out, a)
	}
	return out
}

// sortAnchors keeps every list stable, so that a diagnostic built from one reads
// the same way twice and a test can compare without sorting first.
func sortAnchors(a []Anchor) {
	sort.Slice(a, func(i, j int) bool {
		if a[i].Path != a[j].Path {
			return a[i].Path < a[j].Path
		}
		return a[i].Marker < a[j].Marker
	})
}
