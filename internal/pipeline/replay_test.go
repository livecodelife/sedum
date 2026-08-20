package pipeline

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/recording"
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

// Replay's whole claim is that the model's contribution can be captured once
// and replayed forever. The tests below are that claim: same recording, same
// output, no model.

// recorded runs a grow and returns the recording it produced, with the output
// tree it wrote.
func recorded(t *testing.T, cfg Config) (recording.Recording, string) {
	t.Helper()
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	return Capture(result), cfg.Output
}

// withField is a stub answer selecting one action, so a run has something to
// record beyond an empty list.
const withField = `{"invocations": [{"action": "addField", "kwargs": {"model": "user", "field": "email"}}]}`

func replayConfig(t *testing.T, cfg Config) ReplayConfig {
	t.Helper()
	return ReplayConfig{
		Generators: cfg.Generators,
		Output:     t.TempDir(),
	}
}

// Same recording, same output, and no sampling involved.
func TestReplayReproducesTheRunItRecorded(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, grown := recorded(t, cfg)

	rcfg := replayConfig(t, cfg)
	if _, err := Replay(rec, rcfg); err != nil {
		t.Fatalf("Replay: %v", err)
	}

	if got, want := tree(t, rcfg.Output), tree(t, grown); !equalStrings(got, want) {
		t.Fatalf("replay wrote a different tree:\n got %v\nwant %v", got, want)
	}
	for _, rel := range tree(t, grown) {
		if a, b := readFile(t, filepath.Join(grown, filepath.FromSlash(rel))), readFile(t, filepath.Join(rcfg.Output, filepath.FromSlash(rel))); a != b {
			t.Errorf("%s differs after replay:\n grow   %q\n replay %q", rel, a, b)
		}
	}
}

// The point of a recording is that replaying it costs no model call at all. A
// replay that reached a client would be a replay that could sample differently.
func TestReplayInvokesNoModel(t *testing.T) {
	cfg := config(t)
	client := &stub{response: withField}
	cfg.Client = client
	rec, _ := recorded(t, cfg)

	callsAfterGrow := client.calls
	if _, err := Replay(rec, replayConfig(t, cfg)); err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if client.calls != callsAfterGrow {
		t.Errorf("replay made %d model call(s)", client.calls-callsAfterGrow)
	}
}

// Stopping at invocations and resuming with a recording matches an
// uninterrupted run - which is what makes --stop-after a resume point rather
// than a way to throw away a call already paid for.
func TestStoppingAtInvocationsAndResumingMatchesAWholeRun(t *testing.T) {
	whole := config(t)
	whole.Client = &stub{response: withField}
	_, wholeTree := recorded(t, whole)

	stopped := config(t)
	stopped.Client = &stub{response: withField}
	stopped.Generators = whole.Generators
	stopped.Records = whole.Records
	stopped.StopAfterPhase = PhaseValidate
	partial, err := Run(context.Background(), stopped)
	if err != nil {
		t.Fatalf("stopped run: %v", err)
	}

	resumed := ReplayConfig{Generators: whole.Generators, Output: t.TempDir()}
	if _, err := Replay(Capture(partial), resumed); err != nil {
		t.Fatalf("resume: %v", err)
	}

	for _, rel := range tree(t, wholeTree) {
		if a, b := readFile(t, filepath.Join(wholeTree, filepath.FromSlash(rel))), readFile(t, filepath.Join(resumed.Output, filepath.FromSlash(rel))); a != b {
			t.Errorf("%s differs between an uninterrupted run and a resumed one:\n whole   %q\n resumed %q", rel, a, b)
		}
	}
}

// A hand-edited recording naming an action that does not exist fails the check
// a model response would fail, and nothing is written.
func TestAnInvalidActionHaltsWithNoPartialWrites(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)

	rec.Records[0].Phases[0].Invocations[0].Action = "addFeeld"

	rcfg := replayConfig(t, cfg)
	_, err := Replay(rec, rcfg)
	if err == nil {
		t.Fatal("a recording naming a nonexistent action did not halt")
	}
	var invalid *InvalidRecording
	if !errors.As(err, &invalid) {
		t.Fatalf("want an InvalidRecording a caller can read as data, got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "addFeeld") {
		t.Errorf("the diagnostic does not name the action: %v", err)
	}

	// Phase 3 creates files and Phase 5 is what rejected, so the file exists
	// and the injection must not have happened.
	for _, rel := range tree(t, rcfg.Output) {
		if strings.Contains(readFile(t, filepath.Join(rcfg.Output, filepath.FromSlash(rel))), "attribute :email") {
			t.Errorf("%s carries an injection from a recording that did not validate", rel)
		}
	}
}

// A required kwarg left unbound is the other half of "validation is identical".
func TestAMissingRequiredKwargHalts(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)

	delete(rec.Records[0].Phases[0].Invocations[0].Kwargs, "field")

	if _, err := Replay(rec, replayConfig(t, cfg)); err == nil {
		t.Fatal("a recording omitting a required kwarg did not halt")
	} else if !strings.Contains(err.Error(), "field") {
		t.Errorf("the diagnostic does not name the kwarg: %v", err)
	}
}

// A package whose extensions changed since the recording halts, because
// replaying it would generate under conventions the recording never chose.
func TestAPackageThatNoLongerClaimsItsExtensionsHalts(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)

	rec.Packages["rails"] = recording.Package{Extensions: []string{".rb", ".rake"}}

	_, err := Replay(rec, replayConfig(t, cfg))
	if err == nil {
		t.Fatal("a package that no longer claims a recorded extension did not halt")
	}
	for _, want := range []string{"rails", ".rake"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the diagnostic does not name %q: %v", want, err)
		}
	}
}

func TestAPackageMissingFromTheDirectoryHalts(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)

	rec.Packages["gone"] = recording.Package{Extensions: []string{".zz"}}

	_, err := Replay(rec, replayConfig(t, cfg))
	if err == nil {
		t.Fatal("a package absent from the generators directory did not halt")
	}
	if !strings.Contains(err.Error(), "gone") {
		t.Errorf("the diagnostic does not name the package: %v", err)
	}
}

// Records are optional. Supplied, the paths are checked; omitted, the recording
// executes as written - a hand-edited generic scaffold has no records.
func TestScopeIsValidatedOnlyWhenRecordsAreSupplied(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)

	// A path no supplied record authorizes. The invocation moves with it, so
	// the recording stays internally consistent and the only thing wrong with
	// it is the scope question - otherwise Phase 5's target check would reject
	// it first and this would not be testing scope at all.
	rec.Records[0].Files[0].Path = "app/models/unauthorized.rb"
	rec.Records[0].Phases[0].Invocations[0].Kwargs["model"] = "unauthorized"

	without := replayConfig(t, cfg)
	if _, err := Replay(rec, without); err != nil {
		t.Fatalf("replay without records should execute as written: %v", err)
	}

	with := replayConfig(t, cfg)
	with.Records = cfg.Records
	_, err := Replay(rec, with)
	if err == nil {
		t.Fatal("an unauthorized path passed scope validation")
	}
	if !strings.Contains(err.Error(), "app/models/unauthorized.rb") {
		t.Errorf("the diagnostic does not name the path: %v", err)
	}
}

// The combination a caller driving Sedum uses: concurrent single-invocation
// recordings with scope validation on. Two records refining regions in one file
// is what the marker's record attribute exists to support, and the
// duplicate-path rule must not reach it (prov-2026-dc227be7).
func TestReplaySupportsTwoRecordsOverOneFileWithScopeValidation(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)

	shared := recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		"prov-2026-bbbbbbbb": "  - app/models/user.rb\n",
	})

	rcfg := replayConfig(t, cfg)
	rcfg.Records = shared
	if _, err := Replay(rec, rcfg); err != nil {
		t.Fatalf("replay refused two records naming one file: %v", err)
	}
}

// A recording carries the run's variables, because a recording that rendered
// different text depending on invisible run state would not be one.
func TestCaptureCarriesTheRunsVariables(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	rec := Capture(result)
	if rec.SedumVersion != recording.Version {
		t.Errorf("sedum_version = %q, want %q", rec.SedumVersion, recording.Version)
	}
	if len(rec.Records) != 1 || rec.Records[0].RecordID != "prov-2026-aaaaaaaa" {
		t.Fatalf("unexpected records: %+v", rec.Records)
	}
	if len(rec.Records[0].Phases) != 1 || rec.Records[0].Phases[0].Name != recording.DefaultPhase {
		t.Errorf("phases are not the reserved single default: %+v", rec.Records[0].Phases)
	}
	if pkg, ok := rec.Packages["rails"]; !ok || len(pkg.Extensions) == 0 {
		t.Errorf("the recording does not carry the package it resolved to: %+v", rec.Packages)
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
