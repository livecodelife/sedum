package pipeline

import (
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
	}

	for _, s := range result.Selections {
		rec := recording.Record{RecordID: s.RecordID}

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

		out.Records = append(out.Records, rec)
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
