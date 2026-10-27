package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// A record's scope names every kind of file its change touches, and only a
// file a generator package claims can carry markers. A check run that is asked
// to tolerate the others skips them and lists each distinct path once; any
// other run still refuses the first path no package claims (prov-2026-9621daed).

const skipRecordA = "prov-2026-aaaaaaaa"
const skipRecordB = "prov-2026-bbbbbbbb"

func skipConfig(t *testing.T, records map[string]string) Config {
	t.Helper()
	cfg := config(t)
	cfg.Records = recordsDir(t, records)
	cfg.Check = true
	cfg.Client = &stub{}
	return cfg
}

func TestCheckWithSkipUnclaimedScopeSkipsAndListsUnclaimedPaths(t *testing.T) {
	cfg := skipConfig(t, map[string]string{
		skipRecordA: "  - app/models/user.rb\n  - docs/NOTES.md\n  - Makefile\n",
	})
	cfg.SkipUnclaimedScope = true

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if want := []string{"Makefile", "docs/NOTES.md"}; !reflect.DeepEqual(result.SkippedScopePaths, want) {
		t.Errorf("SkippedScopePaths = %v, want %v", result.SkippedScopePaths, want)
	}
	if got := modeOf(t, result, skipRecordA); got != ModeSelect {
		t.Errorf("the record with the claimed path was not classified: mode = %v", got)
	}
}

func TestSkippedScopePathsListsEachDistinctPathOnce(t *testing.T) {
	cfg := skipConfig(t, map[string]string{
		skipRecordA: "  - app/models/user.rb\n  - docs/NOTES.md\n  - scripts/run.sh\n",
		skipRecordB: "  - docs/NOTES.md\n",
	})
	cfg.SkipUnclaimedScope = true

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if want := []string{"docs/NOTES.md", "scripts/run.sh"}; !reflect.DeepEqual(result.SkippedScopePaths, want) {
		t.Errorf("SkippedScopePaths = %v, want each distinct path once: %v", result.SkippedScopePaths, want)
	}
}

func TestSkippedScopePathsIsEmptyWhenEveryPathIsClaimed(t *testing.T) {
	cfg := skipConfig(t, map[string]string{skipRecordA: "  - app/models/user.rb\n"})
	cfg.SkipUnclaimedScope = true

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if len(result.SkippedScopePaths) != 0 {
		t.Errorf("SkippedScopePaths = %v, want none", result.SkippedScopePaths)
	}
}

// The flag is off by default, so a run that does not ask keeps today's
// behaviour: a scope path no package claims stops the run.
func TestARunWithoutSkipUnclaimedScopeRefusesAnUnclaimedPath(t *testing.T) {
	cfg := skipConfig(t, map[string]string{
		skipRecordA: "  - app/models/user.rb\n  - docs/NOTES.md\n",
	})

	if _, err := Run(context.Background(), cfg); err == nil {
		t.Error("check without SkipUnclaimedScope accepted a path no package claims")
	}
}

// The tolerance belongs to check mode only: a run that writes still refuses.
func TestSkipUnclaimedScopeDoesNotApplyOutsideCheck(t *testing.T) {
	cfg := skipConfig(t, map[string]string{
		skipRecordA: "  - app/models/user.rb\n  - docs/NOTES.md\n",
	})
	cfg.Check = false
	cfg.SkipUnclaimedScope = true

	if _, err := Run(context.Background(), cfg); err == nil {
		t.Error("a run that is not a check accepted a path no package claims")
	}
}

// A path a package declares unmanaged is claimed: it is handed off, not skipped.
func TestADeclaredUnmanagedPathIsNotSkipped(t *testing.T) {
	cfg := skipConfig(t, map[string]string{
		skipRecordA: "  - app/models/user.rb\n  - docs/NOTES.md\n",
	})
	manifest := filepath.Join(cfg.Generators, "rails", "sedum.yaml")
	data, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifest, append(data, []byte("unmanaged:\n  - \"docs/*.md\"\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg.SkipUnclaimedScope = true

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if len(result.SkippedScopePaths) != 0 {
		t.Errorf("SkippedScopePaths = %v, want none: the path is declared unmanaged", result.SkippedScopePaths)
	}
	if len(result.Unmanaged) != 1 || result.Unmanaged[0].Path != "docs/NOTES.md" {
		t.Errorf("Unmanaged = %+v, want docs/NOTES.md handed off", result.Unmanaged)
	}
}
