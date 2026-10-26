package pipeline

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/livecodelife/sedum/internal/expand"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/record"
	"github.com/livecodelife/sedum/internal/recording"
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
	// ModeEdited is a record whose regions are all current and at least one of
	// which no longer matches its body hash. Nothing is written for it
	// (prov-2026-9621daed).
	ModeEdited
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
	case ModeEdited:
		return "edited"
	}
	return "unknown mode"
}

// RegionRef names one region in a report. Reason says why a v2 region is
// stale, and is empty everywhere else.
type RegionRef struct {
	Path    string
	Action  string
	Variant string
	Reason  string
}

// The reasons a v2 region is stale against its recording.
const (
	reasonNoRecording = "no recording"
	reasonUnreadable  = "recording unreadable"
	reasonNoEntry     = "recording has no entry for this region"
	reasonNotMatching = "recording does not match the marker"
)

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
	// Edited are the current regions whose text no longer matches their body
	// hash. A record in ModeEdited keeps in Rerender the regions a re-render
	// would have rewritten.
	Edited []RegionRef
	// Unverifiable are the current owned-tier regions with no body hash to
	// compare: v1 markers, and v2 markers that carry none.
	Unverifiable []RegionRef
	// Foreign are the regions in the record's files whose marker has a version
	// this Sedum does not know. They decide nothing and are never rewritten.
	Foreign []RegionRef

	// found is every region in the record's files, whoever owns it: the
	// record's own feed its rebuilt prior selections, and other records' are
	// listed in its prompt as already present (prov-2026-cc73cf61).
	found []inject.OwnedRegion

	// rerendered are the invocations that bring the Rerender regions up to
	// date, built during classification so that a region which cannot be
	// re-rendered sends its record to re-select before any model is asked.
	rerendered []inject.Invocation
}

// Failing reports whether --check fails on this record: it is stale, needs
// re-rendering, or is edited. A record with no regions is unknown, which does
// not fail.
func (c Classification) Failing() bool {
	return (c.Mode == ModeReselect && !c.Forced) || c.Mode == ModeRerender || c.Mode == ModeEdited
}

// Refused reports whether a re-render was refused because of a hand edit: the
// record is edited and a region a re-render would have rewritten is among the
// edited ones.
func (c Classification) Refused() bool {
	if c.Mode != ModeEdited {
		return false
	}
	for _, e := range c.Edited {
		for _, r := range c.Rerender {
			if e.Path == r.Path && e.Action == r.Action && e.Variant == r.Variant {
				return true
			}
		}
	}
	return false
}

// needsModel reports whether the record goes to Phase 4.
func (c Classification) needsModel() bool {
	return c.Mode == ModeSelect || c.Mode == ModeReselect
}

func refOf(r inject.OwnedRegion) RegionRef {
	return RegionRef{Path: r.Path, Action: r.Region.Marker.Action, Variant: r.Region.Marker.Variant}
}

// recordingsDirOf is the directory a run reads and writes recordings in: the
// configured one, or .sedum/recordings under the output directory.
func recordingsDirOf(cfg Config) string {
	if cfg.Recordings != "" {
		return cfg.Recordings
	}
	return filepath.Join(cfg.Output, ".sedum", "recordings")
}

// entryFor finds the entry of a recording's regions that has the key of a
// marker's region.
func entryFor(regions []recording.Region, path string, m inject.Marker) (recording.Region, bool) {
	key := recording.RegionKey(path, m.Action, m.Variant, m.Kwargs)
	for _, e := range regions {
		if recording.RegionKey(e.Path, e.Action, e.Variant, e.ID) == key {
			return e, true
		}
	}
	return recording.Region{}, false
}

// recordingReason says why a v2 region is not vouched for by the regions of its
// record's recording, or the empty string when it is. readErr is what reading
// the recording returned.
func recordingReason(readErr error, regions []recording.Region, r inject.OwnedRegion) string {
	if readErr != nil {
		if errors.Is(readErr, fs.ErrNotExist) {
			return reasonNoRecording
		}
		return reasonUnreadable
	}
	m := r.Region.Marker
	entry, ok := entryFor(regions, r.Path, m)
	if !ok {
		return reasonNoEntry
	}
	if entry.KwargsHash != m.KwargsHash {
		return reasonNotMatching
	}
	hash, err := inject.KwargsHash(entry.Kwargs)
	if err != nil || hash != entry.KwargsHash {
		return reasonNotMatching
	}
	return ""
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
func classify(records *record.Set, files []resolve.File, owners map[string]inject.Owner, output string, unwritten map[string]string, reselect bool, skip []string, recordings string) ([]Classification, error) {
	skipped := skipSet(skip)
	var out []Classification
	for _, rec := range records.Records {
		mine := filesOf(files, rec.ID)
		c := Classification{RecordID: rec.ID}
		if !anyManaged(mine) {
			c.Mode = ModeNone
			out = append(out, c)
			continue
		}

		var targets []inject.Target
		for _, t := range targetsOf(mine) {
			if !within(skipped, filepath.Join(output, filepath.FromSlash(t.Path))) {
				targets = append(targets, t)
			}
		}
		found, err := inject.ScanRegions(targets, output, unwritten)
		if err != nil {
			return nil, err
		}
		c.found = found
		var owned []inject.OwnedRegion
		for _, r := range found {
			if r.Region.Marker.Foreign {
				c.Foreign = append(c.Foreign, refOf(r))
				continue
			}
			if inject.Owns(owners, rec.ID, r.Region.Marker.Record) {
				owned = append(owned, r)
			}
		}

		// The recording is read once, and only for a record with a v2 region
		// whose gen is current, since nothing else looks at it.
		var (
			read    bool
			readErr error
			entries []recording.Region
		)
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
			if m.Version == 2 {
				if !read {
					read = true
					var got recording.Recording
					got, readErr = recording.ReadRecord(recordings, rec.ID)
					if readErr == nil && len(got.Records) > 0 {
						entries = got.Records[0].Regions
					}
				}
				if reason := recordingReason(readErr, entries, r); reason != "" {
					ref := refOf(r)
					ref.Reason = reason
					c.Stale = append(c.Stale, ref)
					continue
				}
			}
			current = append(current, r)
		}
		for _, r := range current {
			switch r.BodyState() {
			case inject.BodyEdited:
				c.Edited = append(c.Edited, refOf(r))
			case inject.BodyUnverifiable:
				c.Unverifiable = append(c.Unverifiable, refOf(r))
			}
		}

		switch {
		case len(c.Stale) == 0 && len(current) == 0:
			c.Mode = ModeSelect
		case reselect:
			c.Mode, c.Forced = ModeReselect, len(c.Stale) == 0
		case len(c.Stale) > 0:
			c.Mode = ModeReselect
		default:
			c.Mode = rerenderOrSkip(&c, rec, mine, current, entries)
			if len(c.Edited) > 0 && (c.Mode == ModeSkip || c.Mode == ModeRerender) {
				c.Mode = ModeEdited
			}
		}
		out = append(out, c)
	}
	return out, nil
}

// rerenderOrSkip finishes classifying a record whose owned regions are all
// current against it: skip when every pkg matches, otherwise re-render the
// regions that do not - or re-select if one of them cannot be re-rendered.
// A v2 region is re-rendered from the kwargs of its entry in regions, the
// Regions of the record's recording; a v1 region from the kwargs on its own
// marker.
func rerenderOrSkip(c *Classification, rec *record.Record, files []resolve.File, current []inject.OwnedRegion, regions []recording.Region) Mode {
	for _, r := range current {
		m := r.Region.Marker
		action, ok := r.Package.Actions[m.Action]
		if ok && m.Pkg != "" && m.Pkg == action.PkgHash {
			continue
		}
		c.Rerender = append(c.Rerender, refOf(r))

		kwargs := m.Kwargs
		if m.Version == 2 {
			entry, found := entryFor(regions, r.Path, m)
			if !found {
				c.Rerender, c.rerendered = nil, nil
				return ModeReselect
			}
			kwargs = entry.Kwargs
		}
		inv, err := expand.Rerender(rec.ID, r.Package, files, m.Action, kwargs)
		if err != nil || inv.Path != r.Path {
			// The region can no longer be reproduced from what records it, so
			// only a fresh selection can say what it should be.
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
