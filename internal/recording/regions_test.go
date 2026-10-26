package recording

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// Marker v2: a record of a recording holds every region after expansion with
// the kwargs it was rendered from (prov-2026-e95f5ccd, prov-2026-446c050c).

func TestRegionTypeFieldsAndTypes(t *testing.T) {
	typ := reflect.TypeOf(Region{})
	if typ.Kind() != reflect.Struct {
		t.Fatalf("Region is %s, want a struct", typ.Kind())
	}

	want := map[string]reflect.Type{
		"Path":       reflect.TypeOf(""),
		"Action":     reflect.TypeOf(""),
		"Variant":    reflect.TypeOf(""),
		"Tier":       reflect.TypeOf(""),
		"Pkg":        reflect.TypeOf(""),
		"ID":         reflect.TypeOf(map[string]any{}),
		"Kwargs":     reflect.TypeOf(map[string]any{}),
		"KwargsHash": reflect.TypeOf(""),
		"Body":       reflect.TypeOf(""),
	}
	for name, wantType := range want {
		f, ok := typ.FieldByName(name)
		if !ok {
			t.Errorf("Region has no field %s", name)
			continue
		}
		if f.Type != wantType {
			t.Errorf("Region.%s is %s, want %s", name, f.Type, wantType)
		}
	}
	if typ.NumField() != len(want) {
		t.Errorf("Region has %d fields, want exactly the %d the record names", typ.NumField(), len(want))
	}
}

// None of the nine keys is omitted when empty: a region entry always says what
// it is, so a reader never has to guess whether a missing key means empty.
func TestRegionJSONKeysAreNeverOmitted(t *testing.T) {
	payload, err := json.Marshal(Region{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(payload, &keys); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	var got []string
	for k := range keys {
		got = append(got, k)
	}
	sort.Strings(got)
	want := []string{"action", "body", "id", "kwargs", "kwargs_hash", "path", "pkg", "tier", "variant"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("an empty Region marshals with keys %v, want %v\n%s", got, want, payload)
	}
}

func TestRegionJSONKeysCarryTheirFields(t *testing.T) {
	region := Region{
		Path:       "app/models/user.rb",
		Action:     "addSetting",
		Variant:    "strict",
		Tier:       "seeded",
		Pkg:        "sha256:pkg",
		ID:         map[string]any{"name": "x"},
		Kwargs:     map[string]any{"name": "x", "mode": "strict"},
		KwargsHash: "sha256:kw",
		Body:       "sha256:body",
	}
	payload, err := json.Marshal(region)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var keys map[string]any
	if err := json.Unmarshal(payload, &keys); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	for key, want := range map[string]string{
		"path":        "app/models/user.rb",
		"action":      "addSetting",
		"variant":     "strict",
		"tier":        "seeded",
		"pkg":         "sha256:pkg",
		"kwargs_hash": "sha256:kw",
		"body":        "sha256:body",
	} {
		if keys[key] != want {
			t.Errorf("key %q = %v, want %q\n%s", key, keys[key], want, payload)
		}
	}
	if !reflect.DeepEqual(keys["id"], map[string]any{"name": "x"}) {
		t.Errorf("key id = %v", keys["id"])
	}
	if !reflect.DeepEqual(keys["kwargs"], map[string]any{"name": "x", "mode": "strict"}) {
		t.Errorf("key kwargs = %v", keys["kwargs"])
	}
}

// A Region read back keeps its numbers as written, like every other kwargs map
// in a recording.
func TestRegionRoundTripsThroughAFile(t *testing.T) {
	in := Recording{
		SedumVersion: Version,
		Packages:     map[string]Package{},
		Records: []Record{{
			RecordID: "prov-2026-00000000",
			Files:    []File{},
			Phases:   []Phase{{Name: DefaultPhase}},
			Regions: []Region{{
				Path:       "a.rb",
				Action:     "addField",
				Tier:       "owned",
				Pkg:        "sha256:p",
				ID:         map[string]any{"limit": json.Number("1000000")},
				Kwargs:     map[string]any{"limit": json.Number("1000000"), "name": "n"},
				KwargsHash: "sha256:k",
				Body:       "sha256:b",
			}},
		}},
	}

	var buf bytes.Buffer
	if err := Write(&buf, in); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if !strings.Contains(buf.String(), `"limit": 1000000`) {
		t.Errorf("a large integer kwarg was not written as written:\n%s", buf.String())
	}

	out, err := Read(&buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if !reflect.DeepEqual(out.Records[0].Regions, in.Records[0].Regions) {
		t.Errorf("regions did not round trip:\n got %#v\nwant %#v", out.Records[0].Regions, in.Records[0].Regions)
	}
}

func TestRecordRegionsIsOmittedWhenEmpty(t *testing.T) {
	for name, rec := range map[string]Record{
		"nil":   {RecordID: "r"},
		"empty": {RecordID: "r", Regions: []Region{}},
	} {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(rec)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if strings.Contains(string(payload), "regions") {
				t.Errorf("a record with no regions carries the key:\n%s", payload)
			}
		})
	}
}

func TestRecordRegionsHasTheKeyRegions(t *testing.T) {
	f, ok := reflect.TypeOf(Record{}).FieldByName("Regions")
	if !ok {
		t.Fatal("Record has no Regions field")
	}
	if f.Type != reflect.TypeOf([]Region{}) {
		t.Errorf("Record.Regions is %s, want []Region", f.Type)
	}

	payload, err := json.Marshal(Record{RecordID: "r", Regions: []Region{{Path: "a.rb"}}})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(payload, &keys); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := keys["regions"]; !ok {
		t.Errorf("a record with a region has no regions key:\n%s", payload)
	}
}

// A recording written before regions existed reads exactly as before.
func TestARecordWithoutRegionsStillReads(t *testing.T) {
	out, err := Read(strings.NewReader(`{
	  "sedum_version": "0.7.0",
	  "packages": {},
	  "records": [{"record_id": "r", "files": [], "phases": [{"name": "default", "invocations": []}]}]
	}`))
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(out.Records) != 1 || len(out.Records[0].Regions) != 0 {
		t.Errorf("got %+v, want one record with no regions", out.Records)
	}
}
