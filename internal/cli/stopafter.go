package cli

import "strings"

// A stopPoint is a phase boundary --stop-after may name.
//
// The rules governing each one are data rather than a chain of conditionals so
// that M7 has a single place to correct if replay's phase set ever changes. See
// prov-2026-86a21f53 for why replay's coverage, not a fixed list of accepted
// values, decides legality under --execute.
type stopPoint struct {
	// name is the value the user passes to --stop-after.
	name string

	// afterPhase is the pipeline phase this boundary follows.
	afterPhase int

	// runByReplay reports whether --execute runs afterPhase at all. Replay
	// enters at Phase 3 with resolution already decided and skips Phase 4,
	// so Phases 1, 2, and 4 are never run.
	runByReplay bool

	// requiresRecord reports whether stopping here discards model output that
	// has already been paid for, making --record mandatory. Only meaningful
	// when a model call happens, so it is not consulted under --execute.
	requiresRecord bool
}

var stopPoints = []stopPoint{
	{name: "resolution", afterPhase: 2, runByReplay: false, requiresRecord: false},
	{name: "files", afterPhase: 3, runByReplay: true, requiresRecord: false},
	{name: "invocations", afterPhase: 5, runByReplay: true, requiresRecord: true},
	{name: "expansion", afterPhase: 6, runByReplay: true, requiresRecord: true},
}

func lookupStopPoint(name string) (stopPoint, bool) {
	for _, sp := range stopPoints {
		if sp.name == name {
			return sp, true
		}
	}
	return stopPoint{}, false
}

func stopPointNames() string {
	names := make([]string, 0, len(stopPoints))
	for _, sp := range stopPoints {
		names = append(names, sp.name)
	}
	return strings.Join(names, ", ")
}
