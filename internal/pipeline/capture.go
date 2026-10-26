package pipeline

import (
	"github.com/livecodelife/sedum/internal/record"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/resolve"
)

// Capture turns what a run resolved and decided into a recording.
//
// It lives here rather than in the recording package because the import
// direction is fixed: pipeline knows about recording, and recording is a schema
// package that knows about nothing. A Capture in the other package would need
// the Result type and invert that.
//
// Everything it reads is already on the Result. Nothing is recomputed, so a
// recording cannot describe a run that did not happen.
func Capture(result *Result) recording.Recording {
	out := recording.Recording{
		SedumVersion: recording.Version,
		Packages:     map[string]recording.Package{},
		Variables:    result.Variables,
		Prompt:       result.Prompt,
	}

	for _, s := range result.Selections {
		rec := recording.Record{
			RecordID:  s.RecordID,
			Gen:       genOf(result.Records, s.RecordID),
			CreatedAt: createdAtOf(result.Records, s.RecordID),
		}

		for _, f := range s.Files {
			// An unmanaged path has no package by design, and replay does not
			// create it. Recording it as a file would say Sedum resolved
			// something it deliberately did not, and replay would have no
			// package name to verify it against.
			if f.Unmanaged || f.Package == nil {
				continue
			}
			rec.Files = append(rec.Files, recording.File{
				Path:     f.Path,
				Package:  f.Package.Name,
				Template: f.Template,
				Captures: capturesOf(f),
			})
			out.Packages[f.Package.Name] = recording.Package{Extensions: f.Package.Extensions}
		}

		// The phase level is reserved rather than used: this implementation
		// emits exactly one, named default. A record contributing no
		// invocations still carries the phase, so that every record in a
		// recording has the same shape.
		rec.Phases = []recording.Phase{{
			Name:        recording.DefaultPhase,
			Invocations: s.Invocations,
		}}

		rec.Regions = regionsOf(result, s.RecordID)

		out.Records = append(out.Records, rec)
	}

	return out
}

// regionsOf lists, as recording regions, the regions a run wrote for one
// record (prov-2026-5cc2f7f6). Only a region written with a v2 marker has
// something to record: a pruned, skipped, unselected, held or orphan region
// was not written, and a v1 marker carries no hashes. A region rewritten with
// the same bytes is included, because its entry still has to say what is on
// the marker. The entries keep the order of the injections; sorting is the
// recording writer's.
func regionsOf(result *Result, recordID string) []recording.Region {
	var out []recording.Region
	for _, in := range result.Injections {
		if in.Record != recordID {
			continue
		}
		if in.Pruned || in.Skipped || in.Unselected || in.Held || in.Orphan {
			continue
		}
		m := in.Marker
		if m.Version != 2 {
			continue
		}
		out = append(out, recording.Region{
			Path:       in.Path,
			Action:     in.Action,
			Variant:    in.Variant,
			Tier:       string(m.Tier),
			Pkg:        m.Pkg,
			ID:         m.Kwargs,
			Kwargs:     in.Kwargs,
			KwargsHash: m.KwargsHash,
			Body:       m.Body,
		})
	}
	return out
}

// capturesOf returns a file's captures, or nil when it bound none.
//
// nil rather than an empty map so the field is omitted: a capture-free file
// writing "captures": {} into every entry is noise in an artifact whose value
// is that it diffs cleanly.
func capturesOf(f resolve.File) map[string]string {
	if len(f.Captures) == 0 {
		return nil
	}
	return f.Captures
}

// createdAtOf returns the creation date of the record with the given id, or
// nothing when the run has no such record.
func createdAtOf(records *record.Set, recordID string) string {
	if records == nil {
		return ""
	}
	rec, ok := records.Lookup(recordID)
	if !ok {
		return ""
	}
	return rec.CreatedAt
}
