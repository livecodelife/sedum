package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/evidence"
)

// A marker can name a record that has no record file, as the fixture markers of
// a test file do. `sedum evidence` reports such regions in text mode and skips
// them: every record with a file still gets its verdict, and the exit code does
// not count them.

// evOrphanTree has the grown file of evRecord and a second file whose marker
// names evOther, a record with no file.
func evOrphanTree(t *testing.T) *evTree {
	t.Helper()
	tree := newEvTree(t, evMarker(evRecord, "sha256:aa", "sha256:01"))
	tree.write("app/models/other.rb", evMarker(evOther, "sha256:bb", "sha256:02"))
	tree.record(evRecord, "implemented", 3)
	return tree
}

func TestEvidenceSkipsRegionsOfARecordWithNoFileAndNamesThemInTextMode(t *testing.T) {
	tree := evOrphanTree(t)

	out, err := tree.run()

	if err != nil {
		t.Fatalf("evidence: %v\n%s", err, out)
	}
	if !strings.Contains(out, "skipped "+evOther) || !strings.Contains(out, "app/models/other.rb") || !strings.Contains(out, "no record file") {
		t.Errorf("no line names the record and the file of the skipped regions:\n%s", out)
	}
	if !strings.Contains(out, evRecord+"  none  now=") || !strings.HasSuffix(out, "\n0 fresh, 0 stale, 1 none, 0 unknown, 0 no regions\n") {
		t.Errorf("the record with a file did not get its verdict:\n%s", out)
	}
}

func TestEvidenceNamesEachFileOfASkippedRecordOnceAndEachRecordOnOneLine(t *testing.T) {
	tree := evOrphanTree(t)
	tree.write("app/models/third.rb", evMarker(evOther, "sha256:bb", "sha256:02")+evMarker(evOther, "sha256:bb", "sha256:02"))

	out, err := tree.run()

	if err != nil {
		t.Fatalf("evidence: %v\n%s", err, out)
	}
	if n := strings.Count(out, "skipped "+evOther); n != 1 {
		t.Errorf("the record is named on %d lines, want 1:\n%s", n, out)
	}
	if n := strings.Count(out, "app/models/third.rb"); n != 1 {
		t.Errorf("a file is named %d times, want once:\n%s", n, out)
	}
}

func TestEvidenceRefreshAppendsOnlyForRecordsWithAFile(t *testing.T) {
	tree := evOrphanTree(t)

	out, err := tree.run("--refresh")

	if err != nil {
		t.Fatalf("evidence --refresh: %v\n%s", err, out)
	}
	if !strings.Contains(out, "appended "+evRecord+" version 3 ") || strings.Contains(out, "appended "+evOther) {
		t.Errorf("output = %q", out)
	}
	if tree.evidenceFile() == "" {
		t.Error("no evidence was appended for the record with a file")
	}
	if got := evidenceFileOf(tree, evOther); got != "" {
		t.Errorf("evidence was appended for a record with no file: %q", got)
	}
}

func TestEvidenceStillFailsForAStaleRecordBesideASkippedOne(t *testing.T) {
	tree := evOrphanTree(t)
	if _, err := tree.run("--refresh"); err != nil {
		t.Fatal(err)
	}
	tree.write("app/models/user.rb", evMarker(evRecord, "sha256:ab", "sha256:01"))

	_, err := tree.run()

	wantExit(t, err, 1, "")
}

func TestEvidenceJSONDocumentIsUnchangedBySkippedRegions(t *testing.T) {
	tree := evOrphanTree(t)

	out, err := tree.run("--json")

	if err != nil {
		t.Fatalf("evidence --json: %v\n%s", err, out)
	}
	if strings.Contains(out, "skipped") || strings.Contains(out, evOther) {
		t.Errorf("the JSON document mentions the skipped record:\n%s", out)
	}
	if problems := sameKeys("evidence", decodeJSON(t, out), example(t, "evidence.example.json")); len(problems) > 0 {
		t.Errorf("the document disagrees with conformance/formats/evidence.example.json:\n%s", strings.Join(problems, "\n"))
	}
}

func TestEvidenceRecordFlagLeavesOutTheSkippedRegionsAndStillRefusesTheirRecord(t *testing.T) {
	tree := evOrphanTree(t)

	out, err := tree.run("--record", evRecord)
	if err != nil || strings.Contains(out, "skipped") {
		t.Errorf("narrowed to a real record: output = %q err = %v", out, err)
	}
	_, err = tree.run("--record", evOther)
	wantExit(t, err, 2, "no record named "+evOther)
}

func evidenceFileOf(e *evTree, id string) string {
	data, _ := os.ReadFile(evidence.Path(e.out, id))
	return string(data)
}
