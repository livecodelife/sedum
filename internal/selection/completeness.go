package selection

import (
	"fmt"
	"strings"

	"github.com/livecodelife/sedum/internal/expand"
)

// completeness is the one failure class Phase 5 is structurally blind to.
//
// Every other check reads what the model returned. None of them can read what it
// did not return, and a short list is valid output - so a run that selects
// thirteen of fourteen correct invocations passes every rule, first try, with no
// retry, because nothing about it was wrong.
//
// What closes it is that a created file states on disk what work it expects. A
// file template plants the markers an action's anchor targets, so planted minus
// filled is the set of regions this run made and nothing it chose will write
// into (prov-2026-6d87dc11).

// note renders the observation fed back to the model.
//
// It names the file and the anchor rather than guessing at an action, because
// which action fills a region is the model's judgment and this is deliberately
// not taking it. The question is asked as a question - a template may plant a
// region a given change does not need, and an instruction to fill it would be
// wrong exactly as often as the model's original omission was.
func note(unfilled []expand.Anchor) string {
	var b strings.Builder
	b.WriteString("Your selection is valid. Before it is accepted, one observation:\n\n")

	// The note names the file and says nothing about who created it: since
	// records share files, it may well have existed before the run
	// (prov-2026-53f4f82b).
	if len(unfilled) == 1 {
		fmt.Fprintf(&b, "%s carries an injection point named %q that nothing you selected writes into.\n",
			unfilled[0].Path, unfilled[0].Marker)
	} else {
		b.WriteString("These files carry injection points that nothing you selected writes into:\n")
		for _, a := range unfilled {
			fmt.Fprintf(&b, "  %s: %q\n", a.Path, a.Marker)
		}
	}

	b.WriteString("\nThat may be correct - a file can carry a region this change does not need, ")
	b.WriteString("or one a later change will fill. It may also mean an action was missed.\n\n")
	b.WriteString("Reply with the complete list of invocations you want applied, in the same format. ")
	b.WriteString("Include everything from your previous answer that should still be applied, ")
	b.WriteString("plus anything you now want to add. If your previous answer was already complete, repeat it unchanged.")

	return b.String()
}

// anchorSummary renders the unfilled set for the run log, in one line.
func anchorSummary(unfilled []expand.Anchor) string {
	parts := make([]string, 0, len(unfilled))
	for _, a := range unfilled {
		parts = append(parts, a.Path+":"+a.Marker)
	}
	return strings.Join(parts, ", ")
}

// withoutPresent drops the anchors a region owned by another record already
// fills (prov-2026-53f4f82b). Such an anchor is not work this record left
// undone, and the only answer that would fill it - selecting that record's
// region - is refused in Phase 7 as an ownership conflict. A region fills the
// anchor its action is declared against.
func withoutPresent(unfilled []expand.Anchor, req Request) []expand.Anchor {
	if len(req.Present) == 0 {
		return unfilled
	}
	filled := map[expand.Anchor]bool{}
	for _, p := range req.Present {
		for _, f := range req.Files {
			if f.Path != p.Path || f.Package == nil {
				continue
			}
			if action, ok := f.Package.Actions[p.Action]; ok && action.Anchor != "" {
				filled[expand.Anchor{Path: p.Path, Marker: action.Anchor}] = true
			}
		}
	}
	var out []expand.Anchor
	for _, a := range unfilled {
		if !filled[a] {
			out = append(out, a)
		}
	}
	return out
}
