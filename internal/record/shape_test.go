package record

import (
	"strconv"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/genhash"
)

// Sedum reads a record's shape into its gen, so a marker's gen is the same
// gen_hash every other tool computes for the record.
func TestGenCoversTheRecordsShape(t *testing.T) {
	// shape is a schema 2 key.
	v2 := strings.Replace(validRecord, "status: implemented\n", "status: implemented\nversion: 2\nreason: test\n", 1)
	without := loadOne(t, v2)
	with := loadOne(t, v2+"shape:\n  - path: a.go\n    package: a\n")
	if with.Gen == without.Gen {
		t.Fatalf("a record with a shape has the shape-less gen %s", with.Gen)
	}
	want := genhash.Of(genhash.Input{Intent: with.Intent, Constraints: with.Constraints,
		AffectedScope: []string{"app/controllers/users_controller.rb"}, ForbiddenScope: []string{},
		Shape: []any{map[string]any{"path": "a.go", "package": "a"}}})
	if with.Gen != want {
		t.Errorf("gen %s, want %s", with.Gen, want)
	}
}

// shape is a key sedum reads now, so a malformed one is the library's error,
// reported at its line in the original file.
func TestMalformedShapeIsReported(t *testing.T) {
	_, _, err := Load(writeRecords(t, map[string]string{"r.yml": strings.Replace(validRecord, "status: implemented\n", "status: implemented\nversion: 2\nreason: test\n", 1) + "shape: nope\n"}), Options{})
	if err == nil || !strings.Contains(err.Error(), "shape") {
		t.Fatalf("a malformed shape loaded: %v", err)
	}
}

// A shape the library rejects as a whole carries the line of the shape key, as
// every other load error does: a shape on a schema 1 record, and an entry with
// an unknown key.
func TestRejectedShapeIsReportedAtItsLine(t *testing.T) {
	v1 := validRecord + "shape:\n  - path: a.go\n"
	v2 := strings.Replace(validRecord, "status: implemented\n", "status: implemented\nversion: 2\nreason: test\n", 1) + "shape:\n  - path: a.go\n    bogus: x\n"
	for name, c := range map[string]struct{ rec, want string }{
		"schema 1":    {v1, "schema 1"},
		"unknown key": {v2, "bogus"},
	} {
		_, _, err := Load(writeRecords(t, map[string]string{"r.yml": c.rec}), Options{})
		line := strings.Count(c.rec[:strings.Index(c.rec, "shape:")], "\n") + 1
		if err == nil || !strings.Contains(err.Error(), "line "+strconv.Itoa(line)+": ") || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: want a line %d error mentioning %q, got %v", name, line, c.want, err)
		}
	}
}
