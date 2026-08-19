package recording

import (
	"encoding/json"
	"fmt"
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

// A recording is written indented and read back leniently. These are the two
// halves of "committed artifact": it has to diff well, and it has to survive
// being read by a Sedum other than the one that wrote it.

func TestWriteRoundTrips(t *testing.T) {
	want := Recording{
		SedumVersion: Version,
		Packages:     map[string]Package{"rails": {Extensions: []string{".rb", ".erb"}}},
		Records: []Record{{
			RecordID: "PR-014",
			Files: []File{{
				Path:     "app/controllers/users_controller.rb",
				Package:  "rails",
				Template: "app/controllers/{name}_controller.rb",
				Captures: map[string]string{"name": "users"},
			}},
			Phases: []Phase{{
				Name:        DefaultPhase,
				Invocations: []Invocation{{Action: "createControllerMethod", Kwargs: map[string]any{"name": "index"}}},
			}},
		}},
		Variables: map[string]string{"module": "example.com/app"},
	}

	var buf strings.Builder
	if err := Write(&buf, want); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := Read(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	if got.SedumVersion != want.SedumVersion ||
		got.Records[0].RecordID != want.Records[0].RecordID ||
		got.Records[0].Files[0].Path != want.Records[0].Files[0].Path ||
		got.Records[0].Files[0].Captures["name"] != "users" ||
		got.Records[0].Phases[0].Invocations[0].Action != "createControllerMethod" ||
		got.Variables["module"] != "example.com/app" {
		t.Errorf("round trip lost something:\ngot  %+v\nwant %+v", got, want)
	}
}

// One line of JSON diffs as one changed line whatever moved inside it, which
// would make "a change in model output is visible as a change in the file" true
// only in the letter.
func TestWriteIsDiffable(t *testing.T) {
	var buf strings.Builder
	if err := Write(&buf, Recording{
		SedumVersion: Version,
		Records:      []Record{{RecordID: "r", Phases: []Phase{{Name: DefaultPhase}}}},
	}); err != nil {
		t.Fatalf("write: %v", err)
	}

	if !strings.Contains(buf.String(), "\n  \"sedum_version\"") {
		t.Errorf("the recording is not indented:\n%s", buf.String())
	}
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Errorf("the recording has no trailing newline")
	}
}

// Equivalent runs must produce byte-identical bytes, which means map ordering
// cannot leak into the output.
func TestWriteIsStableAcrossRuns(t *testing.T) {
	build := func() string {
		var buf strings.Builder
		if err := Write(&buf, Recording{
			SedumVersion: Version,
			Packages: map[string]Package{
				"rails": {Extensions: []string{".rb"}},
				"chi":   {Extensions: []string{".go"}},
				"cairn": {Extensions: []string{".crn"}},
			},
			Variables: map[string]string{"b": "2", "a": "1", "c": "3"},
		}); err != nil {
			t.Fatalf("write: %v", err)
		}
		return buf.String()
	}

	first := build()
	for i := 0; i < 20; i++ {
		if got := build(); got != first {
			t.Fatalf("two writes of one recording differ:\n%s\n---\n%s", first, got)
		}
	}
}

// A field a later Sedum adds must be an addition rather than a migration across
// every repository holding a committed recording (prov-2026-88e36921).
func TestReadIgnoresKeysItDoesNotModel(t *testing.T) {
	got, err := Read(strings.NewReader(`{
	  "sedum_version": "9.9.9",
	  "something_a_later_version_added": {"nested": [1, 2]},
	  "records": [{"record_id": "r", "unknown_on_record": true,
	    "phases": [{"name": "default", "invocations": []}]}]
	}`))
	if err != nil {
		t.Fatalf("a recording from a later Sedum did not load: %v", err)
	}
	if got.SedumVersion != "9.9.9" || len(got.Records) != 1 {
		t.Errorf("the known fields did not survive: %+v", got)
	}
}

// Leniency stops at the door: bytes that are not JSON at all are still an
// error, and the diagnostic says what it was reading.
func TestReadRejectsWhatIsNotJSON(t *testing.T) {
	_, err := Read(strings.NewReader("phases: [default]"))
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "recording") {
		t.Errorf("the diagnostic does not name what was being read: %v", err)
	}
}

// A kwarg the caller wrote as 1000000 must not reach a template as 1e+06.
func TestReadPreservesNumbersInKwargs(t *testing.T) {
	got, err := Read(strings.NewReader(`{
	  "records": [{"record_id": "r", "phases": [{"name": "default",
	    "invocations": [{"action": "a", "kwargs": {"size": 1000000}}]}]}]
	}`))
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	if rendered := fmt.Sprintf("%v", got.Records[0].Phases[0].Invocations[0].Kwargs["size"]); rendered != "1000000" {
		t.Errorf("a JSON number decoded to %q", rendered)
	}
}
