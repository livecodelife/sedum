package pipeline

import (
	"context"
	"strings"
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

// The duplicate-path check moved from ingestion to Phase 4's entry. The grow
// path's behaviour is unchanged, and that is the test - two records naming one
// file still halt with the same diagnostic. What changed is when: before the
// first model call rather than before records are read (prov-2026-dc227be7).

func TestTwoRecordsNamingOneFileStillHalt(t *testing.T) {
	cfg := config(t)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		"prov-2026-bbbbbbbb": "  - app/models/user.rb\n",
	})
	client := &stub{}
	cfg.Client = client

	_, err := Run(context.Background(), cfg)
	if err == nil {
		t.Fatal("a path named by two records did not halt the run")
	}
	for _, want := range []string{"app/models/user.rb", "prov-2026-aaaaaaaa", "prov-2026-bbbbbbbb"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the diagnostic does not name %q: %v", want, err)
		}
	}

	// The point of moving it rather than deleting it: nothing was paid for.
	if client.calls != 0 {
		t.Errorf("the run made %d model call(s) before halting on a duplicate path", client.calls)
	}
}

// Records are still ingested and resolved when they share a path. The rule is
// Phase 4's, so everything before Phase 4 has to be reachable - this is what a
// caller supplying records purely for scope validation depends on.
func TestADuplicatePathSurvivesEverythingBeforePhaseFour(t *testing.T) {
	cfg := config(t)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		"prov-2026-bbbbbbbb": "  - app/models/user.rb\n",
	})
	cfg.StopAfterPhase = PhaseCreate

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("a run stopping before Phase 4 halted on a rule that belongs to it: %v", err)
	}
	if len(result.Records.Records) != 2 {
		t.Errorf("want both records ingested, got %d", len(result.Records.Records))
	}
}
