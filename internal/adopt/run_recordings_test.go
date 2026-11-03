package adopt

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/livecodelife/sedum/internal/pipeline"
	"github.com/livecodelife/sedum/internal/recording"
)

// A real adopt leaves the same per-record recordings a grow or a replay does,
// with an entry for each region it wrote, so a later check finds the record
// current.

const adoptedImprint = `id: prov-2026-aaaaaaaa
type: imprint
status: open
intent: adopt py files
created_at: "2026-10-16"
affected_scope: ["a.py", "b.py"]
`

const twoUnits = "def f():\n    return 1\n\ndef g():\n    return 2\n"

func recordingsOf(cfg Config) string {
	return filepath.Join(cfg.Root, ".sedum", "recordings")
}

func TestARealAdoptWritesTheRecordRecordingWithItsRegions(t *testing.T) {
	cfg := newRunFixture(t, map[string]string{"a.py": twoUnits})

	if _, err := Run(cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}

	got, err := recording.ReadRecord(recordingsOf(cfg), "prov-2026-aaaaaaaa")
	if err != nil {
		t.Fatalf("the record has no recording: %v", err)
	}
	if len(got.Records) != 1 || len(got.Records[0].Regions) != 2 {
		t.Fatalf("recording = %+v, want one record with a region for each of the two units", got.Records)
	}
	for _, r := range got.Records[0].Regions {
		if r.Path != "a.py" || r.Action != "setThing" || r.KwargsHash == "" || r.Pkg == "" {
			t.Errorf("region = %+v, want an entry for a unit of a.py with its hashes", r)
		}
	}
}

func TestACheckAfterARealAdoptFindsTheRecordCurrent(t *testing.T) {
	cfg := newRunFixture(t, map[string]string{"a.py": twoUnits, "b.py": "def h():\n    return 3\n"})
	if err := os.WriteFile(filepath.Join(filepath.Dir(cfg.Records), "prov", "imp.yml"), []byte(adoptedImprint), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}

	result, err := pipeline.Run(context.Background(), pipeline.Config{
		Generators: cfg.Generators,
		Records:    cfg.Records,
		Output:     cfg.Root,
		Check:      true,
	})
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if len(result.Classifications) != 1 || result.Classifications[0].Mode != pipeline.ModeSkip {
		t.Errorf("classifications = %+v, want the one record current", result.Classifications)
	}
}

func TestADryRunWritesNoRecordRecording(t *testing.T) {
	cfg := newRunFixture(t, map[string]string{"a.py": twoUnits})
	cfg.DryRun = true

	if _, err := Run(cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if _, err := os.Stat(recordingsOf(cfg)); !os.IsNotExist(err) {
		t.Errorf("a dry run left %s (stat error %v)", recordingsOf(cfg), err)
	}
}

func TestARunThatAdoptsNothingWritesNoRecordRecording(t *testing.T) {
	cfg := newRunFixture(t, map[string]string{"a.py": "# a license header\ndef f():\n    return 1\n"})

	rep, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if rep.Totals.Adopted != 0 {
		t.Fatalf("totals = %+v, want nothing adopted", rep.Totals)
	}
	if _, err := os.Stat(recordingsOf(cfg)); !os.IsNotExist(err) {
		t.Errorf("a run that adopted nothing left %s (stat error %v)", recordingsOf(cfg), err)
	}
}

func TestARecordingKeepsTheEntriesOfAnEarlierAdoptOfAnotherFile(t *testing.T) {
	cfg := newRunFixture(t, map[string]string{"a.py": twoUnits, "b.py": "def h():\n    return 3\n"})
	first := cfg
	first.Paths = []string{"a.py"}
	if _, err := Run(first); err != nil {
		t.Fatalf("first Run: %v", err)
	}
	second := cfg
	second.Paths = []string{"b.py"}
	if _, err := Run(second); err != nil {
		t.Fatalf("second Run: %v", err)
	}

	got, err := recording.ReadRecord(recordingsOf(cfg), "prov-2026-aaaaaaaa")
	if err != nil {
		t.Fatalf("the record has no recording: %v", err)
	}
	paths := map[string]int{}
	for _, r := range got.Records[0].Regions {
		paths[r.Path]++
	}
	if paths["a.py"] != 2 || paths["b.py"] != 1 {
		t.Errorf("regions by file = %v, want 2 for a.py and 1 for b.py", paths)
	}
}
