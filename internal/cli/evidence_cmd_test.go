package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/evidence"
)

// `sedum evidence` reads the regions and records, judges each record, refreshes
// on request, writes the report and returns the exit code.

const (
	evRecord = "prov-2026-aaaaaaaa"
	evOther  = "prov-2026-bbbbbbbb"
)

func evRecordText(id, status string, version int) string {
	return "id: " + id + "\ntitle: A record\nstatus: " + status + "\nversion: " + itoaEv(version) + "\ncreated_at: \"2026-10-27\"\nauthor: a@example.com\nintent: |\n  Add email.\nconstraints:\n  - Add an email field\naffected_scope:\n  - app/models/user.rb\n"
}

func itoaEv(n int) string { return string(rune('0' + n)) }

func evMarker(record, gen, pkg string) string {
	fields := `"tier":"owned","record":"` + record + `"`
	if gen != "" {
		fields += `,"gen":"` + gen + `"`
	}
	if pkg != "" {
		fields += `,"pkg":"` + pkg + `"`
	}
	return "class User\n  # sedum:addField {" + fields + `,"kwargs":{"model":"user","field":"email"}}` + "\n  attribute :email\n  # /sedum:addField\nend\n"
}

// evTree is an output directory with one grown file and a records directory.
type evTree struct {
	t       *testing.T
	out     string
	records string
}

func newEvTree(t *testing.T, marker string) *evTree {
	t.Helper()
	tree := &evTree{t: t, out: t.TempDir(), records: filepath.Join(t.TempDir(), "provenance")}
	tree.write("app/models/user.rb", marker)
	if err := os.MkdirAll(tree.records, 0o755); err != nil {
		t.Fatal(err)
	}
	return tree
}

func (e *evTree) write(rel, content string) {
	e.t.Helper()
	full := filepath.Join(e.out, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		e.t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		e.t.Fatal(err)
	}
}

func (e *evTree) record(id, status string, version int) {
	e.t.Helper()
	if err := os.WriteFile(filepath.Join(e.records, id+".yml"), []byte(evRecordText(id, status, version)), 0o644); err != nil {
		e.t.Fatal(err)
	}
}

func (e *evTree) run(extra ...string) (string, error) {
	e.t.Helper()
	args := []string{"evidence", "--generators", filepath.Join("..", "..", "testdata", "generators"), "--records", e.records, "--output", e.out}
	return exec(e.t, append(args, extra...)...)
}

func (e *evTree) evidenceFile() string {
	data, _ := os.ReadFile(evidence.Path(e.out, evRecord))
	return string(data)
}

func TestEvidenceNeedsGeneratorsThenRecords(t *testing.T) {
	_, err := exec(t, "evidence")
	wantExit(t, err, 2, "sedum evidence needs --generators")
	_, err = exec(t, "evidence", "--generators", "g")
	wantExit(t, err, 2, "sedum evidence needs --records")
}

func TestEvidenceTurnsAFlagErrorIntoExitTwo(t *testing.T) {
	_, err := exec(t, "evidence", "--nope")
	wantExit(t, err, 2, "")
}

func TestEvidenceReportsNoneWithoutAnEntryWritesNothingAndSucceeds(t *testing.T) {
	tree := newEvTree(t, evMarker(evRecord, "sha256:aa", "sha256:01"))
	tree.record(evRecord, "implemented", 3)

	out, err := tree.run()

	if err != nil {
		t.Fatalf("evidence: %v\n%s", err, out)
	}
	if !strings.Contains(out, evRecord+"  none  now=") || !strings.HasSuffix(out, "\n0 fresh, 0 stale, 1 none, 0 unknown, 0 no regions\n") {
		t.Errorf("output = %q", out)
	}
	if _, err := os.Stat(filepath.Join(tree.out, ".sedum", "evidence")); !os.IsNotExist(err) {
		t.Error("a report without --refresh wrote evidence")
	}
}

func TestEvidenceRefreshAppendsOnceThenReportsFreshAndStaleAfterACodeChange(t *testing.T) {
	tree := newEvTree(t, evMarker(evRecord, "sha256:aa", "sha256:01"))
	tree.record(evRecord, "implemented", 3)

	out, err := tree.run("--refresh")
	if err != nil {
		t.Fatalf("refresh: %v\n%s", err, out)
	}
	if !strings.HasPrefix(out, "appended "+evRecord+" version 3 sha256:") || !strings.Contains(out, evRecord+"  fresh  now=") {
		t.Errorf("refresh output = %q", out)
	}
	if lines := strings.Count(tree.evidenceFile(), "\n"); lines != 1 {
		t.Errorf("evidence file has %d lines, want 1", lines)
	}

	again, err := tree.run("--refresh")
	if err != nil || strings.Contains(again, "appended") || strings.Count(tree.evidenceFile(), "\n") != 1 {
		t.Errorf("a fresh record was appended to again: %q %v", again, err)
	}

	tree.write("app/models/user.rb", evMarker(evRecord, "sha256:aa", "sha256:99"))
	out, err = tree.run()
	exit, ok := err.(ExitError)
	if !ok || exit.Code != 1 || exit.Message != "" {
		t.Fatalf("err = %#v, want a silent exit 1 for a stale record", err)
	}
	if !strings.Contains(out, evRecord+"  stale  now=") {
		t.Errorf("the report is written before the exit: %q", out)
	}

	out, err = tree.run("--refresh")
	if err != nil || !strings.Contains(out, "appended "+evRecord+" version 3") || !strings.Contains(out, "  fresh  ") || strings.Count(tree.evidenceFile(), "\n") != 2 {
		t.Errorf("refresh after a change: %q %v", out, err)
	}
}

func TestEvidenceReportsUnknownAndNeverAppendsForIt(t *testing.T) {
	tree := newEvTree(t, evMarker(evRecord, "", ""))
	tree.record(evRecord, "implemented", 3)

	out, err := tree.run("--refresh")

	exit, ok := err.(ExitError)
	if !ok || exit.Code != 1 {
		t.Errorf("err = %#v", err)
	}
	if !strings.Contains(out, evRecord+"  unknown  ") || strings.Contains(out, "appended") {
		t.Errorf("output = %q", out)
	}
	if tree.evidenceFile() != "" {
		t.Error("evidence was appended for an unknown record")
	}
}

func TestEvidenceReportsNoRegionsForARecordWithoutMarkers(t *testing.T) {
	tree := newEvTree(t, "class User\nend\n")
	tree.record(evRecord, "implemented", 3)

	out, err := tree.run()

	if err != nil || !strings.Contains(out, evRecord+"  no regions  ") {
		t.Errorf("output = %q err = %v", out, err)
	}
}

func TestEvidenceIgnoresTheRegionsAndTheReportOfADeprecatedRecord(t *testing.T) {
	tree := newEvTree(t, evMarker(evOther, "sha256:aa", "sha256:01"))
	tree.record(evRecord, "implemented", 3)
	tree.record(evOther, "deprecated", 2)

	out, err := tree.run()

	if err != nil || strings.Contains(out, evOther) || !strings.Contains(out, evRecord+"  no regions  ") {
		t.Errorf("output = %q err = %v", out, err)
	}
	_, err = tree.run("--record", evOther)
	wantExit(t, err, 2, "no record named "+evOther)
}

func TestEvidenceRecordNarrowsTheReportAndAnUnknownNameFails(t *testing.T) {
	tree := newEvTree(t, evMarker(evRecord, "sha256:aa", "sha256:01"))
	tree.record(evRecord, "implemented", 3)
	tree.record(evOther, "implemented", 1)

	out, err := tree.run("--record", evRecord)
	if err != nil || strings.Contains(out, evOther) || !strings.Contains(out, evRecord) {
		t.Errorf("narrowed output = %q err = %v", out, err)
	}
	_, err = tree.run("--record", "prov-2026-cccccccc")
	wantExit(t, err, 2, "no record named prov-2026-cccccccc")
}

func TestEvidenceFailsOnARecordFileItCannotDecode(t *testing.T) {
	tree := newEvTree(t, "class User\nend\n")
	tree.record(evRecord, "implemented", 3)
	if err := os.WriteFile(filepath.Join(tree.records, evOther+".yml"), []byte("id: "+evOther+"\nstatus: ["), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := tree.run()

	if err == nil {
		t.Fatal("no error")
	}
	wantExit(t, err, 2, "")
}

func TestEvidenceFailsOnAnEvidenceFileItCannotRead(t *testing.T) {
	tree := newEvTree(t, evMarker(evRecord, "sha256:aa", "sha256:01"))
	tree.record(evRecord, "implemented", 3)
	tree.write(".sedum/evidence/"+evRecord+".jsonl", "garbage\n")

	_, err := tree.run()

	wantExit(t, err, 2, "")
	if !strings.Contains(err.Error(), evRecord+".jsonl") {
		t.Errorf("err = %v, want one naming the file", err)
	}
}

func TestEvidenceJSONMatchesItsDocumentedFormat(t *testing.T) {
	tree := newEvTree(t, evMarker(evRecord, "sha256:aa", "sha256:01"))
	tree.record(evRecord, "implemented", 3)
	if _, err := tree.run("--refresh"); err != nil {
		t.Fatal(err)
	}

	out, err := tree.run("--json")

	if err != nil {
		t.Fatalf("evidence --json: %v\n%s", err, out)
	}
	if problems := sameKeys("evidence", decodeJSON(t, out), example(t, "evidence.example.json")); len(problems) > 0 {
		t.Errorf("sedum evidence --json disagrees with conformance/formats/evidence.example.json:\n%s", strings.Join(problems, "\n"))
	}
}

func TestEvidenceCommandHasTheDocumentedFlags(t *testing.T) {
	cmd, _, err := NewRootCommand().Find([]string{"evidence"})
	if err != nil || cmd.Name() != "evidence" {
		t.Fatalf("evidence command = %v, %v", cmd, err)
	}
	defaults := map[string]string{"record": "", "output": defaultOutput, "generators": "", "records": "", "refresh": "false", "json": "false"}
	for name, def := range defaults {
		f := cmd.Flags().Lookup(name)
		if f == nil || f.DefValue != def {
			t.Errorf("flag --%s = %+v, want default %q", name, f, def)
		}
	}
	if len(cmd.Annotations) > 0 {
		t.Errorf("annotations = %v: --generators and --records must not be marked required to cobra", cmd.Annotations)
	}
	if _, err := exec(t, "evidence", "extra"); err == nil {
		t.Error("a positional argument was accepted")
	}
}
