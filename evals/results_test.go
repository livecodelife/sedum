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
	return NewEntry(m, "http://127.0.0.1:1234/v1")
}

// The history is evidence a record can cite, and citing it meant re-deriving
// the summary by hand until this existed (prov-2026-c0f55691).
func TestHistoryReadsBackWhatWasRecorded(t *testing.T) {
	m := measurement(
		Sample{Counts: map[string]int{"addColumn": 2}, Calls: 1, PromptTokens: 2000, CompletionTokens: 400},
		Sample{Counts: map[string]int{"addColumn": 2}, Calls: 2, Rejected: 1, PromptTokens: 4000, CompletionTokens: 800},
	)
	m.Wall = 100 * time.Second
	m.Retries = 2
	m.Concurrency = 2
	e := NewEntry(m, "http://x")
	e.Commit = "abc1234"
	e.Clean = true

	var buf strings.Builder
	History(&buf, "fixture", []Entry{e})
	out := buf.String()

	for _, want := range []string{
		"fixture  (1 entries)",
		"abc1234",
		"2/2 [0.34,1.00]", // both samples answered acceptably
		"1/2 [0.09,0.91]", // one of them needed no retry
		"1.50",            // three calls over two samples
		"2000+400",        // 6000 prompt and 1200 completion over three calls
		"12.0",            // 1200 completion tokens over 100s; the prompt was billed, not computed
	} {
		if !strings.Contains(out, want) {
			t.Errorf("history omits %q:\n%s", want, out)
		}
	}

	// A clean tree earns no asterisk and no legend line about one.
	if strings.Contains(out, "not re-runnable") {
		t.Errorf("a clean entry was marked dirty:\n%s", out)
	}
}

// An entry written before a field existed says so rather than reading as a
// measurement of zero. Four such entries exist already, and the additive
// schema's whole promise is that they keep meaning what they meant.
func TestHistoryPrintsAnAbsentFieldAsAbsent(t *testing.T) {
	// An entry as it was written before call counts and tokens existed.
	old := Entry{
		Schema:  EntrySchema,
		Commit:  "0000000",
		Clean:   true,
		At:      time.Now().UTC(),
		Case:    "fixture",
		Samples: 5,
		Valid:   5,
		WallMS:  60_000,
		Runs: []SampleRun{
			{Outcome: "valid", MS: 12_000}, {Outcome: "valid", MS: 12_000},
			{Outcome: "valid", MS: 12_000}, {Outcome: "valid", MS: 12_000},
			{Outcome: "valid", MS: 12_000},
		},
	}

	var buf strings.Builder
	History(&buf, "fixture", []Entry{old})
	out := buf.String()

	if !strings.Contains(out, "predates that field") {
		t.Errorf("history does not explain its dashes:\n%s", out)
	}
	// The rate it does carry is still reported; only the missing fields dash.
	if !strings.Contains(out, "5/5 [0.57,1.00]") {
		t.Errorf("history dropped a rate the old entry did carry:\n%s", out)
	}
	if strings.Contains(out, "0.00") {
		t.Errorf("an unrecorded count printed as zero:\n%s", out)
	}
}

// Consecutive entries whose intervals overlap are marked, for the reason every
// rate got an interval: a column of fractions invites a reader to see a change
// that five samples cannot support (prov-2026-0baaa119).
func TestHistoryMarksEntriesThatDoNotDistinguishEachOther(t *testing.T) {
	entry := func(valid, invalid int) Entry {
		e := Entry{Schema: EntrySchema, Clean: true, At: time.Now().UTC(),
			Samples: valid + invalid, Valid: valid, Invalid: invalid, WallMS: 1000}
		for i := 0; i < valid; i++ {
			e.Runs = append(e.Runs, SampleRun{Outcome: "valid", Calls: 1})
		}
		for i := 0; i < invalid; i++ {
			e.Runs = append(e.Runs, SampleRun{Outcome: "invalid", Calls: 1, Rejected: 1})
		}
		return e
	}

	var overlapping strings.Builder
	History(&overlapping, "fixture", []Entry{entry(5, 0), entry(4, 1)})
	if !strings.Contains(overlapping.String(), "~") {
		t.Errorf("5/5 then 4/5 was not marked as indistinguishable:\n%s", overlapping.String())
	}

	var separated strings.Builder
	History(&separated, "fixture", []Entry{entry(5, 0), entry(0, 5)})
	if strings.Contains(separated.String(), "~") {
		t.Errorf("5/5 then 0/5 was marked as indistinguishable:\n%s", separated.String())
	}
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
