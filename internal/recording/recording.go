// Package recording is the schema of a recorded execution.
//
// Every phase after model invocation is deterministic, which means the model's
// contribution to a run can be captured once and replayed forever. This package
// is the shape of that capture.
//
// It is types only. There is no serialization here and no replay entry point:
// those are M7's work. The types exist now because M4 drives Phases 6 and 7
// from hand-written invocation lists, and a fixture shaped like a recording
// makes M7 largely serialization work rather than new execution machinery. The
// JSON tags are part of the shape rather than an anticipation of a writer -
// they are what the PRD's example document says the field names are, recorded
// beside the fields so the two cannot drift.
package recording

// DefaultPhase names the single phase every recording this implementation
// produces contains.
//
// Invocations are grouped under phases rather than listed flat because
// recordings are committed artifacts: a team's standard service scaffold lives
// in version control, and adding the grouping after those files exist means
// migrating them. Reserving the level costs one array nesting.
const DefaultPhase = "default"

// Recording is everything Sedum resolved and everything it was told to do.
//
// It carries no timestamps, run identifiers, or other volatile fields. Two
// recordings of equivalent runs diff cleanly, so that a change in model output
// is visible as a change in the file.
type Recording struct {
	SedumVersion string             `json:"sedum_version"`
	Packages     map[string]Package `json:"packages"`
	Records      []Record           `json:"records"`

	// Variables are the run's values for the project facts its packages
	// declare. A recording carries them because replay renders the same
	// invocations against the same packages, and a recording that rendered
	// different text depending on invisible run state would not be a recording
	// (prov-2026-6fc3d13d).
	Variables map[string]string `json:"variables,omitempty"`
}

// Package is a generator package as the run resolved it. Replay verifies that
// the package is still present and still claims these extensions, because a
// recording that resolved .rb to one package and is replayed against a
// directory where another claims it would generate under the wrong conventions.
type Package struct {
	Extensions []string `json:"extensions"`
}

// Record is one provenance record's contribution to a run.
type Record struct {
	RecordID string  `json:"record_id"`
	Files    []File  `json:"files"`
	Phases   []Phase `json:"phases"`
}

// File is one authorized path with the resolution decided for it: which package
// claimed it, which file template matched, and what that template's captures
// bound.
type File struct {
	Path     string            `json:"path"`
	Package  string            `json:"package"`
	Template string            `json:"template"`
	Captures map[string]string `json:"captures"`
}

// Phase is a named, ordered group of invocations. Replay executes phases in
// order.
type Phase struct {
	Name        string       `json:"name"`
	Invocations []Invocation `json:"invocations"`
}

// Invocation is one action selection with its bound arguments.
//
// Invocations are recorded pre-expansion: a composite is stored as the
// composite rather than as its children, because expansion is deterministic and
// re-running it costs nothing. This keeps a recording at the abstraction level
// an author edits in - changing a createMethod call means editing one entry,
// not keeping two injections in sync.
type Invocation struct {
	Action string         `json:"action"`
	Kwargs map[string]any `json:"kwargs"`
}
