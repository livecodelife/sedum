package record

import (
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
