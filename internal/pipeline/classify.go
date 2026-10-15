package pipeline

import (
	"github.com/livecodelife/sedum/internal/expand"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/record"
	"github.com/livecodelife/sedum/internal/resolve"
)

// Mode is what grow does with one record (prov-2026-9621daed).
type Mode int

const (
	// ModeNone is a record with no path a generator package writes. It is
	// skipped without a model call, as it always was, and owns nothing.
	ModeNone Mode = iota
	// ModeSelect is a record that owns no regions yet: a normal grow. --check
	// reports it as unknown, since there is nowhere its hash is kept.
	ModeSelect
	// ModeReselect is a record some of whose owned regions were written from
	// another version of it, or whose regenerate was forced.
	ModeReselect
	// ModeSkip is a record whose every owned region is current.
	ModeSkip
	// ModeRerender is a record whose regions are current against it but were
	// rendered from an older package. They are re-rendered from their markers
	// with no model call.
	ModeRerender
)

func (m Mode) String() string {
	switch m {
	case ModeNone:
		return "nothing to grow"
	case ModeSelect:
		return "select"
	case ModeReselect:
		return "re-select"
	case ModeSkip:
		return "skip"
	case ModeRerender:
		return "re-render"
	}
	return "unknown mode"
}

// RegionRef names one region in a report.
type RegionRef struct {
	Path    string
	Action  string
	Variant string
}

// Classification is the decision for one record and what it was made from.
type Classification struct {
	RecordID string
	Mode     Mode
	// Forced reports that the record was re-selected because the run asked
	// for it (--force, or --record), not because a hash moved.
	Forced bool
	// Stale are the owned regions whose gen differs from the record's.
	Stale []RegionRef
	// Rerender are the owned regions whose pkg differs from their action's.
	Rerender []RegionRef
	// StaleSeeded are seeded regions whose gen differs from the record's.
	// They decide nothing and are never rewritten; they are reported so a
	// person can look.
	StaleSeeded []RegionRef

	// rerendered are the invocations that bring the Rerender regions up to
	// date, built during classification so that a region which cannot be
	// re-rendered sends its record to re-select before any model is asked.
	rerendered []inject.Invocation
}

// Failing reports whether --check fails on this record: it is stale, or it
// needs re-rendering. A record with no regions is unknown, which does not fail.
func (c Classification) Failing() bool {
	return (c.Mode == ModeReselect && !c.Forced) || c.Mode == ModeRerender
}

// needsModel reports whether the record goes to Phase 4.
func (c Classification) needsModel() bool {
	return c.Mode == ModeSelect || c.Mode == ModeReselect
}

func refOf(r inject.OwnedRegion) RegionRef {
	return RegionRef{Path: r.Path, Action: r.Region.Marker.Action, Variant: r.Region.Marker.Variant}
}

// classify decides, from markers alone, what grow does with each record.
//
//   - no owned regions: select
//   - an owned region whose gen differs from the record's, or is unknown:
//     re-select
//   - reselect asked for (--force, --record): re-select
//   - every gen matches, every pkg matches: skip
//   - every gen matches, some pkg differs: re-render those regions, unless one
//     cannot be (its action is gone, or it would now render into another
//     file), in which case re-select
//
// Seeded regions take no part in the decision.
func classify(records *record.Set, files []resolve.File, owners map[string]inject.Owner, output string, unwritten map[string]string, reselect bool) ([]Classification, error) {
	var out []Classification
	for _, rec := range records.Records {
		mine := filesOf(files, rec.ID)
		c := Classification{RecordID: rec.ID}
		if !anyManaged(mine) {
			c.Mode = ModeNone
			out = append(out, c)
			continue
		}

		owned, err := inject.ScanOwned(rec.ID, targetsOf(mine), owners, output, unwritten)
		if err != nil {
			return nil, err
		}

		var current []inject.OwnedRegion
		for _, r := range owned {
			m := r.Region.Marker
			if m.Tier == inject.TierSeeded {
				if m.Gen != rec.Gen {
					c.StaleSeeded = append(c.StaleSeeded, refOf(r))
				}
				continue
			}
			if m.Gen == "" || m.Gen != rec.Gen {
				c.Stale = append(c.Stale, refOf(r))
				continue
			}
			current = append(current, r)
		}

		switch {
		case len(c.Stale) == 0 && len(current) == 0:
			c.Mode = ModeSelect
		case reselect:
			c.Mode, c.Forced = ModeReselect, len(c.Stale) == 0
		case len(c.Stale) > 0:
			c.Mode = ModeReselect
		default:
			c.Mode = rerenderOrSkip(&c, rec, mine, current)
		}
		out = append(out, c)
	}
	return out, nil
}

// rerenderOrSkip finishes classifying a record whose owned regions are all
// current against it: skip when every pkg matches, otherwise re-render the
// regions that do not - or re-select if one of them cannot be re-rendered.
func rerenderOrSkip(c *Classification, rec *record.Record, files []resolve.File, current []inject.OwnedRegion) Mode {
	for _, r := range current {
		m := r.Region.Marker
		action, ok := r.Package.Actions[m.Action]
		if ok && m.Pkg != "" && m.Pkg == action.PkgHash {
			continue
		}
		c.Rerender = append(c.Rerender, refOf(r))

		inv, err := expand.Rerender(rec.ID, r.Package, files, m.Action, m.Kwargs)
		if err != nil || inv.Path != r.Path {
			// The region can no longer be reproduced from its marker, so only
			// a fresh selection can say what it should be.
			c.Rerender, c.rerendered = nil, nil
			return ModeReselect
		}
		inv.Gen = rec.Gen
		inv.Tier = m.Tier
		c.rerendered = append(c.rerendered, inv)
	}
	if len(c.Rerender) == 0 {
		return ModeSkip
	}
	return ModeRerender
}

// targetsOf is the managed files of one record, as Phase 7 targets.
func targetsOf(files []resolve.File) []inject.Target {
	var out []inject.Target
	for _, f := range files {
		if f.Unmanaged || f.Package == nil {
			continue
		}
		out = append(out, inject.Target{Path: f.Path, Package: f.Package})
	}
	return out
}
