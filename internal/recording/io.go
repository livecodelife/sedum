package recording

import (
	"encoding/json"
	"fmt"
	"io"
)

// Write serializes a recording.
//
// The output is indented because a recording is a committed artifact and the
// property that matters is that two recordings of equivalent runs diff cleanly.
// One line of JSON diffs as one changed line whatever moved inside it, which
// would make "a change in model output is visible as a change in the file" true
// only in the letter.
//
// Map keys are emitted in sorted order by encoding/json, so package and capture
// ordering is stable without this package sorting anything itself.
func Write(w io.Writer, r Recording) error {
	payload, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing the recording: %w", err)
	}
	// A trailing newline, so the file is one a text editor and a diff both
	// treat as complete.
	payload = append(payload, '\n')
	if _, err := w.Write(payload); err != nil {
		return fmt.Errorf("writing the recording: %w", err)
	}
	return nil
}

// Read parses a recording.
//
// Decoding is lenient: a key this schema does not model is ignored rather than
// rejected. A recording is a committed artifact that outlives the Sedum that
// wrote it, and strict decoding would turn a field added by a later version
// into a migration across every repository holding one. What a recording asks
// Sedum to do is validated instead, terminally, once the packages it names are
// in hand (prov-2026-88e36921).
//
// Numbers inside kwargs are kept as json.Number for the same reason sedum
// render keeps them: a kwarg the caller wrote as 1000000 must not reach a
// template as 1e+06.
func Read(r io.Reader) (Recording, error) {
	dec := json.NewDecoder(r)
	dec.UseNumber()

	var out Recording
	if err := dec.Decode(&out); err != nil {
		return Recording{}, fmt.Errorf("the recording is not valid JSON: %w", err)
	}
	return out, nil
}
