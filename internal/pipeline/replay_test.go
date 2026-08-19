package pipeline

import (
	"testing"
)

// Replay enters at Phase 3 with resolution already decided, skips Phase 4
// entirely, and runs Phases 5 through 7 unchanged. These are the properties
// that make that possible; the replay behaviour itself is asserted below them
// as it lands.

// The ordering is what "enters at Phase 3 and skips Phase 4" is stated in
// terms of. A renumbering that left the constants inconsistent would make
// replay's stop points name the wrong boundaries, and the --stop-after table in
// the command layer maps onto these numbers with no translation.
func TestPhaseOrderingReplayDependsOn(t *testing.T) {
	ordered := []struct {
		name  string
		phase int
	}{
		{"load", PhaseLoad},
		{"ingest", PhaseIngest},
		{"resolve", PhaseResolve},
		{"create", PhaseCreate},
		{"select", PhaseSelect},
		{"validate", PhaseValidate},
		{"expand", PhaseExpand},
		{"inject", PhaseInject},
	}

	for i, p := range ordered {
		if p.phase != i {
			t.Errorf("%s is phase %d, want %d", p.name, p.phase, i)
		}
	}
}
