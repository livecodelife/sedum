package recording

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

// A recording is a committed artifact. What is asserted here is the shape teams
// will have in version control - that it round-trips, that it carries nothing
// volatile, and that the schema decisions M7 was built on are properties of the
// type rather than of the writer that happens to fill it in.

// Two recordings of equivalent runs must diff cleanly, so no field may carry a
// timestamp, a run identifier, or anything else that changes between runs of
// the same inputs.
func TestRecordingCarriesNothingVolatile(t *testing.T) {
	payload, err := json.Marshal(Recording{
		SedumVersion: Version,
		Packages:     map[string]Package{"rails": {Extensions: []string{".rb"}}},
		Records: []Record{{
			RecordID: "prov-2026-00000000",
			Files:    []File{{Path: "app/models/order.rb", Package: "rails", Template: "app/models/{name}.rb"}},
			Phases:   []Phase{{Name: DefaultPhase}},
		}},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	for _, forbidden := range []string{"time", "timestamp", "run_id", "date", "seed", "duration"} {
		if strings.Contains(string(payload), forbidden) {
			t.Errorf("a recording carries %q, so two equivalent runs would not diff cleanly:\n%s", forbidden, payload)
		}
	}
}

// The phase level is reserved even though this implementation emits exactly
// one. Adding the grouping after recordings exist in version control is a
// migration across repositories rather than a schema change.
func TestPhasesAreAReservedLevel(t *testing.T) {
	if DefaultPhase != "default" {
		t.Errorf("DefaultPhase = %q; committed recordings name it %q", DefaultPhase, "default")
	}

	var got Recording
	if err := json.Unmarshal([]byte(`{
	  "sedum_version": "0.1.0",
	  "records": [{"record_id": "r", "phases": [{"name": "default", "invocations": [{"action": "a", "kwargs": {}}]}]}]
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Records) != 1 || len(got.Records[0].Phases) != 1 {
		t.Fatalf("phases did not survive the round trip: %+v", got)
	}
	if got.Records[0].Phases[0].Name != DefaultPhase {
		t.Errorf("phase name = %q", got.Records[0].Phases[0].Name)
	}
}

// A composite is stored as the composite. Expansion is deterministic and
// re-running it is free, so a recording stays at the abstraction level an
// author edits in.
func TestInvocationsAreHeldPreExpansion(t *testing.T) {
	var got Recording
	if err := json.Unmarshal([]byte(`{
	  "records": [{"record_id": "r", "phases": [{"name": "default", "invocations": [
	    {"action": "createHandler", "kwargs": {"resource": "Order", "verb": "get"}}
	  ]}]}]
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	inv := got.Records[0].Phases[0].Invocations
	if len(inv) != 1 || inv[0].Action != "createHandler" {
		t.Fatalf("the composite was not held as itself: %+v", inv)
	}
}

// The version a recording carries is the one the binary reports, and a caller
// pins a version floor against it.
func TestVersionIsABareSemver(t *testing.T) {
	if !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(Version) {
		t.Errorf("Version = %q, which is not a bare semver", Version)
	}
}
