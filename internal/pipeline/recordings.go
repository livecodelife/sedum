package pipeline

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/livecodelife/sedum/internal/recording"
)

// WriteRecordings writes the per-record recordings a grow run produced into
// dir (prov-2026-5cc2f7f6).
//
// A record the run selected has its recording written whole, from Capture. A
// record that was re-rendered keeps its existing recording, with the entries of
// the regions just written replaced; when it has none, one is built from what
// Capture returns for it. A record that was skipped, edited, or has nothing to
// grow has no file written or changed.
func WriteRecordings(result *Result, dir string) error {
	if len(result.Classifications) == 0 {
		return nil
	}

	captured := Capture(result)
	byID := map[string]recording.Record{}
	for _, r := range captured.Records {
		byID[r.RecordID] = r
	}

	for _, c := range result.Classifications {
		var rec recording.Recording
		switch c.Mode {
		case ModeSelect, ModeReselect:
			r, ok := byID[c.RecordID]
			if !ok {
				continue
			}
			rec = recordingOf(captured, r)
		case ModeRerender:
			existing, err := recording.ReadRecord(dir, c.RecordID)
			switch {
			case err == nil && len(existing.Records) > 0:
				rec = existing
				rec.Records = []recording.Record{rerenderedRecord(result, existing.Records[0], c.RecordID)}
			case err == nil || errors.Is(err, fs.ErrNotExist):
				r, ok := byID[c.RecordID]
				if !ok {
					continue
				}
				rec = recordingOf(captured, r)
			default:
				return fmt.Errorf("reading recording for record %s: %w", c.RecordID, err)
			}
		default:
			continue
		}

		if err := recording.WriteRecords(dir, rec); err != nil {
			return fmt.Errorf("writing recording for record %s: %w", c.RecordID, err)
		}
	}
	return nil
}

// recordingOf is the recording of one record of a captured run.
func recordingOf(captured recording.Recording, r recording.Record) recording.Recording {
	captured.Records = []recording.Record{r}
	return captured
}

// rerenderedRecord is a re-rendered record's recording with the regions just
// written added, replacing the entries that have their keys.
func rerenderedRecord(result *Result, existing recording.Record, recordID string) recording.Record {
	written := regionsOf(result, recordID)
	rewritten := map[string]bool{}
	for _, r := range written {
		rewritten[recording.RegionKey(r.Path, r.Action, r.Variant, r.ID)] = true
	}

	var regions []recording.Region
	for _, r := range existing.Regions {
		if !rewritten[recording.RegionKey(r.Path, r.Action, r.Variant, r.ID)] {
			regions = append(regions, r)
		}
	}
	existing.Regions = append(regions, written...)
	return existing
}

// WriteReplayRecordings writes the per-record recordings a replay produced,
// keeping what was replayed (prov-2026-5cc2f7f6): each record of the replayed
// recording is written with its own files and phases, with the regions the
// replay wrote added to the entries its existing recording holds for the
// regions the replay did not write and, when the replay restamped it, its new
// gen.
func WriteReplayRecordings(replayed recording.Recording, result *Result, dir string) error {
	for _, r := range replayed.Records {
		one := r
		one.Regions = regionsOf(result, r.RecordID)

		// A region the replay did not write is still in the tree, so the
		// entry the record's recording already holds for it is kept.
		// A recording that cannot be read holds nothing to keep; the write
		// below replaces it.
		if existing, err := recording.ReadRecord(dir, r.RecordID); err == nil && len(existing.Records) > 0 {
			one.Regions = rerenderedRecord(result, existing.Records[0], r.RecordID).Regions
		}
		for _, s := range result.Restamps {
			if s.RecordID == r.RecordID {
				one.Gen = s.New
			}
		}

		rec := recording.Recording{
			SedumVersion: replayed.SedumVersion,
			Packages:     replayed.Packages,
			Variables:    replayed.Variables,
			Records:      []recording.Record{one},
		}
		if err := recording.WriteRecords(dir, rec); err != nil {
			return fmt.Errorf("writing recording for record %s: %w", r.RecordID, err)
		}
	}
	return nil
}
