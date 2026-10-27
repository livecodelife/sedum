package evidence

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// The store appends to and reads a record's evidence file, and Latest picks the
// latest entry for a version.

func entry(record string, version int, hash string) Entry {
	return Entry{Record: record, Version: version, CodeHash: hash, At: "2026-10-27T12:00:00Z", Runner: Runner}
}

func TestPathIsUnderSedumEvidence(t *testing.T) {
	if got := Path("/out", "prov-2026-aaaaaaaa"); got != filepath.Join("/out", ".sedum", "evidence", "prov-2026-aaaaaaaa.jsonl") {
		t.Errorf("Path = %q", got)
	}
}

func TestTimestampIsUTCRFC3339WithoutFractionalSeconds(t *testing.T) {
	zone := time.FixedZone("x", 2*3600)
	got := Timestamp(time.Date(2026, 10, 27, 14, 30, 5, 987654321, zone))
	if got != "2026-10-27T12:30:05Z" {
		t.Errorf("Timestamp = %q", got)
	}
}

func TestAppendWritesOneLineOfJSONInFieldOrderAndCreatesTheDirectory(t *testing.T) {
	out := t.TempDir()
	e := entry("prov-2026-aaaaaaaa", 3, "sha256:aa<&>")

	if err := Append(out, e); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(Path(out, e.Record))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"record":"prov-2026-aaaaaaaa","version":3,"code_hash":"sha256:aa<&>","at":"2026-10-27T12:00:00Z","runner":"sedum evidence"}` + "\n"
	if string(data) != want {
		t.Errorf("file = %q, want %q (no HTML escaping)", data, want)
	}
}

func TestAppendNeverRewritesWhatTheFileHoldsAndRepairsAMissingFinalNewline(t *testing.T) {
	out := t.TempDir()
	path := Path(out, "prov-2026-aaaaaaaa")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	old := `{"record":"prov-2026-aaaaaaaa","version":1,"code_hash":"sha256:01","at":"2026-10-01T00:00:00Z","runner":"x"}`
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Append(out, entry("prov-2026-aaaaaaaa", 2, "sha256:02")); err != nil {
		t.Fatal(err)
	}
	if err := Append(out, entry("prov-2026-aaaaaaaa", 3, "sha256:03")); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) != 3 || lines[0] != old || !strings.Contains(lines[1], `"version":2`) || !strings.Contains(lines[2], `"version":3`) {
		t.Errorf("lines = %q", lines)
	}
	if !strings.HasSuffix(string(data), "\n") {
		t.Error("the file does not end with a newline")
	}
}

func TestAppendRefusesBadEntriesAndWritesNothing(t *testing.T) {
	out := t.TempDir()
	sep := string(os.PathSeparator)
	cases := map[string]Entry{
		"empty record":       entry("", 1, "sha256:aa"),
		"a slash":            entry("a/b", 1, "sha256:aa"),
		"the separator":      entry("a"+sep+"b", 1, "sha256:aa"),
		"dot dot":            entry("..", 1, "sha256:aa"),
		"version zero":       entry("prov-2026-aaaaaaaa", 0, "sha256:aa"),
		"an empty code hash": entry("prov-2026-aaaaaaaa", 1, ""),
	}
	for name, e := range cases {
		if err := Append(out, e); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	if _, err := os.Stat(filepath.Join(out, ".sedum")); !os.IsNotExist(err) {
		t.Error("a refused entry created something")
	}
}

func TestAppendErrorNamesTheFileWhenItCannotWrite(t *testing.T) {
	out := t.TempDir()
	if err := os.MkdirAll(Path(out, "prov-2026-aaaaaaaa"), 0o755); err != nil { // a directory where the file goes
		t.Fatal(err)
	}

	err := Append(out, entry("prov-2026-aaaaaaaa", 1, "sha256:aa"))

	if err == nil || !strings.Contains(err.Error(), "prov-2026-aaaaaaaa.jsonl") {
		t.Errorf("err = %v, want one naming the file", err)
	}
}

func TestReadReturnsEntriesInOrderSkippingBlankLinesAndNothingWhenTheFileIsMissing(t *testing.T) {
	out := t.TempDir()
	if got, err := Read(out, "prov-2026-aaaaaaaa"); err != nil || got != nil {
		t.Fatalf("missing file: %v %v", got, err)
	}
	for _, e := range []Entry{entry("prov-2026-aaaaaaaa", 1, "sha256:01"), entry("prov-2026-aaaaaaaa", 2, "sha256:02")} {
		if err := Append(out, e); err != nil {
			t.Fatal(err)
		}
	}
	f, _ := os.OpenFile(Path(out, "prov-2026-aaaaaaaa"), os.O_APPEND|os.O_WRONLY, 0)
	f.WriteString("   \n\n")
	f.Close()

	got, err := Read(out, "prov-2026-aaaaaaaa")

	if err != nil || len(got) != 2 || got[0].Version != 1 || got[1].CodeHash != "sha256:02" {
		t.Errorf("Read = %+v, %v", got, err)
	}
}

func TestReadErrorsNameTheFileAndTheLine(t *testing.T) {
	good := `{"record":"prov-2026-aaaaaaaa","version":1,"code_hash":"sha256:01","at":"t","runner":"r"}`
	cases := map[string]string{
		"not JSON":              "nonsense",
		"not an object":         "[1]",
		"an empty record":       `{"record":"","version":1,"code_hash":"sha256:01"}`,
		"version below one":     `{"record":"prov-2026-aaaaaaaa","version":0,"code_hash":"sha256:01"}`,
		"an empty code hash":    `{"record":"prov-2026-aaaaaaaa","version":1,"code_hash":""}`,
		"another record's line": `{"record":"prov-2026-bbbbbbbb","version":1,"code_hash":"sha256:01"}`,
	}
	for name, bad := range cases {
		out := t.TempDir()
		path := Path(out, "prov-2026-aaaaaaaa")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(good+"\n"+bad+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		_, err := Read(out, "prov-2026-aaaaaaaa")

		if err == nil || !strings.Contains(err.Error(), "prov-2026-aaaaaaaa.jsonl") || !strings.Contains(err.Error(), ":2") {
			t.Errorf("%s: err = %v, want one naming the file and line 2", name, err)
		}
	}
}

func TestReadErrorsNameTheFileWhenItCannotBeOpened(t *testing.T) {
	out := t.TempDir()
	if err := os.MkdirAll(Path(out, "prov-2026-aaaaaaaa"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(out, "prov-2026-aaaaaaaa"); err == nil || !strings.Contains(err.Error(), "prov-2026-aaaaaaaa.jsonl") {
		t.Errorf("err = %v", err)
	}
}

func TestReadChangesNothing(t *testing.T) {
	out := t.TempDir()
	Read(out, "prov-2026-aaaaaaaa")
	if _, err := os.Stat(filepath.Join(out, ".sedum")); !os.IsNotExist(err) {
		t.Error("Read created a directory")
	}
}

func TestLatestReturnsTheLastEntryOfTheVersion(t *testing.T) {
	entries := []Entry{entry("r", 1, "sha256:01"), entry("r", 2, "sha256:02"), entry("r", 2, "sha256:03"), entry("r", 1, "sha256:04")}

	got, ok := Latest(entries, 2)

	if !ok || got.CodeHash != "sha256:03" {
		t.Errorf("Latest = %+v %v", got, ok)
	}
	if got, ok := Latest(entries, 9); ok || !reflect.DeepEqual(got, Entry{}) {
		t.Errorf("Latest of an absent version = %+v %v", got, ok)
	}
}
