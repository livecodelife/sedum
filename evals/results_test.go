package evals

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A result is an observation, not a decision: it is appended and never edited.
// These assert the properties that follow from that (prov-2026-eb283c56).

func entryFor(t *testing.T, samples ...Sample) Entry {
	t.Helper()
	m := measurement(samples...)
	m.Wall = 90 * time.Second
	return NewEntry(m, 0, "http://127.0.0.1:1234/v1")
}

// Appending must preserve what was already there. A run adds lines and rewrites
// nothing, which is what makes the history impossible to edit by accident.
func TestAppendPreservesHistory(t *testing.T) {
	dir := t.TempDir()

	for i := 0; i < 3; i++ {
		if err := Append(dir, entryFor(t, valid(map[string]int{"addColumn": 2}, "addColumn"))); err != nil {
			t.Fatalf("append %d: %v", i, err)
		}
	}

	got, err := ReadEntries(dir, "fixture")
	if err != nil {
		t.Fatalf("ReadEntries: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("read %d entries, want 3 - an append overwrote history", len(got))
	}
	for _, e := range got {
		if e.Schema != EntrySchema {
			t.Errorf("entry carries schema %d, want %d", e.Schema, EntrySchema)
		}
	}

	// One entry per line is what makes it append-only in the way that matters.
	raw, err := os.ReadFile(filepath.Join(dir, "fixture.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if n := len(splitLines(string(raw))); n != 3 {
		t.Errorf("file has %d lines for 3 entries; an entry spans lines", n)
	}
}

// Reading a case with no history is not an error. A first run has nothing to
// compare against and should say so by returning nothing.
func TestReadingAnAbsentHistoryIsNotAnError(t *testing.T) {
	got, err := ReadEntries(t.TempDir(), "never-run")
	if err != nil {
		t.Fatalf("ReadEntries on an absent file: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d entries from nothing", len(got))
	}
}

// The three outcomes survive storage, and an invalid or failed sample carries
// no counts - it has no invocation list to record.
func TestOutcomesSurviveStorage(t *testing.T) {
	e := entryFor(t,
		valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"),
		Sample{Invalid: true, Detail: "did not validate"},
		Sample{Err: os.ErrDeadlineExceeded, Detail: "timeout"},
	)

	if e.Valid != 1 || e.Invalid != 1 || e.Failed != 1 {
		t.Errorf("tally stored as valid=%d invalid=%d failed=%d, want one of each", e.Valid, e.Invalid, e.Failed)
	}
	want := []string{"valid", "invalid", "failed"}
	for i, w := range want {
		if e.Runs[i].Outcome != w {
			t.Errorf("run %d stored as %q, want %q", i, e.Runs[i].Outcome, w)
		}
	}
	if e.Runs[1].Counts != nil || e.Runs[2].Counts != nil {
		t.Error("a sample with no invocation list stored counts anyway")
	}
}

// The whole reason samples are kept individually: an aggregate cannot say
// whether any single answer was complete. The run that established the chi
// expectations reported means of 0.33 and 3.33 - numbers no answer can have.
func TestCompleteFindsTheSamplesAnAggregateCannot(t *testing.T) {
	e := entryFor(t,
		valid(map[string]int{"addColumn": 2, "createEndpoint": 5}, "addColumn"), // complete
		valid(map[string]int{"addColumn": 2, "createEndpoint": 4}, "addColumn"), // short
		valid(map[string]int{"addColumn": 1, "createEndpoint": 5}, "addColumn"), // short
	)

	// The mean of createEndpoint here is 4.67 - a value no sample had.
	complete := e.Complete(map[string]int{"addColumn": 2, "createEndpoint": 5})
	if len(complete) != 1 || complete[0] != 0 {
		t.Errorf("Complete found %v, want only sample 0 - the mean cannot identify it", complete)
	}
}

// An entry says what commit it ran against, and whether that commit means
// anything. A dirty tree makes the commit a lie, and the entry has to carry
// that rather than imply reproducibility it does not have.
func TestAnEntryCarriesItsCommitAndCleanliness(t *testing.T) {
	e := entryFor(t, valid(map[string]int{"addColumn": 2}, "addColumn"))

	if e.Commit == "" {
		t.Error("entry carries no commit; it cannot be re-run from anything")
	}
	if !e.Clean && len(e.Dirty) == 0 {
		t.Error("entry is marked dirty but names nothing; the caveat is unusable")
	}
	if e.Endpoint == "" {
		t.Error("entry carries no endpoint; a hosted row and a local row are not the same evidence")
	}
	if e.Concurrency == 0 || e.Samples == 0 {
		t.Error("entry omits the parameters a rate is conditioned on")
	}
}

func splitLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}
