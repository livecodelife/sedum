package pipeline

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/recording"
)

// Marker v2: grow and replay write one recording per record
// (prov-2026-5cc2f7f6). The units under test here are regionsOf, Capture,
// WriteRecordings and WriteReplayRecordings.
//
// Hash literals below were derived with a throwaway script (sha256 of the
// canonical kwargs JSON and of the region's body text), not typed from memory:
//
//	kwargs {"field":"email","model":"user"}   sha256:6e868af5...d861d6
//	body   "attribute :email\n"               sha256:ca0675c2...8f08037
//	kwargs {"field":"phone","model":"user"}   sha256:bf8d45b3...752278
//	body   "attribute :phone\n"               sha256:4e8f4bee...b2223c
//	body   "attribute :email, :string\n"      sha256:33ee6f79...3531b1
const (
	hashKwargsEmail = "sha256:6e868af53f57da43edaf70a18475b41d08d3728d202e7cc8c4d0658371d861d6"
	hashBodyEmail   = "sha256:ca0675c2bdbe489ce05d122bdc0807d38a99e1bfbee86f085bccd10be8f08037"
	hashKwargsPhone = "sha256:bf8d45b33ceb7c171b5f9d13ee2ddc5de4aaf24e631b50afc453b73c1a752278"
	hashBodyPhone   = "sha256:4e8f4beeb511d18ddd89c3cf3cfcf5dc6229beb36e9d376520d7f9c732b2223c"
	hashBodyEmailT  = "sha256:33ee6f79a0412ab8199e3dd9d0850ed82fb464c3f2f84084558c768bc83531b1"
	hashKwargsTotal = "sha256:bb042f9fcb8027d8765fe52e57bf9ff52f5b8e56907f17cc96abcba950ee52ca"
	hashBodyTotal   = "sha256:f432211614be3569f02e06e2813d5484ebc932d3fe90e5e5579d0557a71323e4"
)

const v2RecordA = "prov-2026-aaaaaaaa"

// shared helpers for the marker v2 tests of this package.

func jsonOf(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal %v: %v", v, err)
	}
	return string(b)
}

// v2Dir is the directory a run reads and writes recordings in.
func v2Dir(cfg Config) string {
	if cfg.Recordings != "" {
		return cfg.Recordings
	}
	return filepath.Join(cfg.Output, ".sedum", "recordings")
}

// v2Grow runs a first grow that selects the email field and writes its
// recordings, as runGrow does after the report.
func v2Grow(t *testing.T) (Config, string) {
	t.Helper()
	cfg := config(t)
	cfg.Client = &stub{response: fields("email")}
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("first grow: %v", err)
	}
	dir := v2Dir(cfg)
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}
	return cfg, dir
}

// v2Run runs grow with the given model response and returns what it did.
func v2Run(t *testing.T, cfg Config, response string) (*Result, *stub) {
	t.Helper()
	client := &stub{response: response}
	cfg.Client = client
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("grow: %v", err)
	}
	return result, client
}

// v2Check classifies without a model and without writing.
func v2Check(t *testing.T, cfg Config) *Result {
	t.Helper()
	cfg.Check = true
	result, client := v2Run(t, cfg, fields("email"))
	if client.calls != 0 {
		t.Fatalf("--check made %d model call(s)", client.calls)
	}
	return result
}

func v2Class(t *testing.T, result *Result, recordID string) Classification {
	t.Helper()
	for _, c := range result.Classifications {
		if c.RecordID == recordID {
			return c
		}
	}
	t.Fatalf("no classification for %s in %+v", recordID, result.Classifications)
	return Classification{}
}

// v2Regions parses the regions of one generated Ruby file.
func v2Regions(t *testing.T, cfg Config, rel string) []inject.Region {
	t.Helper()
	regions, err := inject.FindRegions(genpkg.Comment{Prefix: "#"}, readFile(t, filepath.Join(cfg.Output, filepath.FromSlash(rel))))
	if err != nil {
		t.Fatalf("FindRegions %s: %v", rel, err)
	}
	return regions
}

func v2WriteFile(t *testing.T, cfg Config, rel, content string) {
	t.Helper()
	full := filepath.Join(cfg.Output, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// v2EditRecording reads a record's recording, lets edit change its record, and
// writes it back.
func v2EditRecording(t *testing.T, dir, id string, edit func(*recording.Record)) {
	t.Helper()
	rec, err := recording.ReadRecord(dir, id)
	if err != nil {
		t.Fatalf("ReadRecord %s: %v", id, err)
	}
	edit(&rec.Records[0])
	if err := recording.WriteRecords(dir, rec); err != nil {
		t.Fatalf("WriteRecords %s: %v", id, err)
	}
}

func v2ReadRecord(t *testing.T, dir, id string) recording.Record {
	t.Helper()
	rec, err := recording.ReadRecord(dir, id)
	if err != nil {
		t.Fatalf("ReadRecord %s: %v", id, err)
	}
	if len(rec.Records) != 1 {
		t.Fatalf("recording for %s holds %d records, want 1", id, len(rec.Records))
	}
	return rec.Records[0]
}

// v2Injection is what Phase 7 reports for one written v2 region.
func v2Injection(record, path, action, variant string, tier inject.Tier, id, kwargs map[string]any, kwargsHash, body string) inject.Result {
	return inject.Result{
		Path: path, Action: action, Variant: variant, Record: record,
		Marker: inject.Marker{
			Action: action, Variant: variant, Tier: tier, Record: record,
			Pkg: "sha256:pkg-" + action, Version: 2, Kwargs: id, KwargsHash: kwargsHash, Body: body,
		},
		Kwargs: kwargs,
	}
}

// regionsOf.

func TestRegionsOfListsTheRegionsOneRecordWrote(t *testing.T) {
	result := &Result{Injections: []inject.Result{
		v2Injection("r1", "a.rb", "addField", "", inject.TierOwned,
			map[string]any{"field": "email"}, map[string]any{"field": "email", "model": "user"}, "sha256:k1", "sha256:b1"),
		v2Injection("r2", "b.rb", "addField", "", inject.TierOwned,
			map[string]any{"field": "other"}, map[string]any{"field": "other"}, "sha256:k2", "sha256:b2"),
		v2Injection("r1", "c.rb", "addSetting", "strict", inject.TierSeeded,
			map[string]any{"name": "x"}, map[string]any{"name": "x", "mode": "strict"}, "sha256:k3", "sha256:b3"),
	}}

	got := regionsOf(result, "r1")
	want := []recording.Region{
		{
			Path: "a.rb", Action: "addField", Tier: "owned", Pkg: "sha256:pkg-addField",
			ID: map[string]any{"field": "email"}, Kwargs: map[string]any{"field": "email", "model": "user"},
			KwargsHash: "sha256:k1", Body: "sha256:b1",
		},
		{
			Path: "c.rb", Action: "addSetting", Variant: "strict", Tier: "seeded", Pkg: "sha256:pkg-addSetting",
			ID: map[string]any{"name": "x"}, Kwargs: map[string]any{"name": "x", "mode": "strict"},
			KwargsHash: "sha256:k3", Body: "sha256:b3",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("regionsOf =\n%+v\nwant\n%+v", got, want)
	}
}

func TestRegionsOfTakesIdFromTheMarkerAndKwargsFromTheInjection(t *testing.T) {
	inj := v2Injection("r1", "a.rb", "addField", "", inject.TierOwned,
		map[string]any{"id": "marker-side"}, map[string]any{"id": "marker-side", "all": "injection-side"}, "sha256:k", "sha256:b")

	got := regionsOf(&Result{Injections: []inject.Result{inj}}, "r1")
	if len(got) != 1 {
		t.Fatalf("regionsOf = %+v, want one entry", got)
	}
	if jsonOf(t, got[0].ID) != `{"id":"marker-side"}` {
		t.Errorf("ID = %v, want the marker's Kwargs", got[0].ID)
	}
	if jsonOf(t, got[0].Kwargs) != `{"all":"injection-side","id":"marker-side"}` {
		t.Errorf("Kwargs = %v, want the injection's Kwargs", got[0].Kwargs)
	}
}

func TestRegionsOfLeavesOutRegionsNothingWasWrittenFor(t *testing.T) {
	base := v2Injection("r1", "a.rb", "addField", "", inject.TierOwned,
		map[string]any{"f": "x"}, map[string]any{"f": "x"}, "sha256:k", "sha256:b")

	for name, mutate := range map[string]func(*inject.Result){
		"pruned":     func(r *inject.Result) { r.Pruned = true },
		"skipped":    func(r *inject.Result) { r.Skipped = true },
		"unselected": func(r *inject.Result) { r.Unselected = true },
		"held":       func(r *inject.Result) { r.Held = true },
		"orphan":     func(r *inject.Result) { r.Orphan = true },
		"v1 marker":  func(r *inject.Result) { r.Marker.Version = 0 },
		"other marker version": func(r *inject.Result) {
			r.Marker.Version = 3
		},
		"v1 marker version one": func(r *inject.Result) { r.Marker.Version = 1 },
	} {
		t.Run(name, func(t *testing.T) {
			inj := base
			mutate(&inj)
			if got := regionsOf(&Result{Injections: []inject.Result{inj}}, "r1"); got != nil {
				t.Errorf("regionsOf = %+v, want nil", got)
			}
		})
	}
}

func TestRegionsOfIncludesARegionRewrittenWithTheSameBytes(t *testing.T) {
	inj := v2Injection("r1", "a.rb", "addField", "", inject.TierOwned,
		map[string]any{"f": "x"}, map[string]any{"f": "x"}, "sha256:k", "sha256:b")
	inj.Replaced = true
	inj.Changed = false

	if got := regionsOf(&Result{Injections: []inject.Result{inj}}, "r1"); len(got) != 1 {
		t.Errorf("regionsOf = %+v, want the replaced, unchanged region", got)
	}
}

func TestRegionsOfKeepsTheOrderOfInjectionsAndReturnsNilForNone(t *testing.T) {
	var injections []inject.Result
	for _, p := range []string{"z.rb", "a.rb", "m.rb"} {
		injections = append(injections, v2Injection("r1", p, "addField", "", inject.TierOwned,
			map[string]any{"f": p}, map[string]any{"f": p}, "sha256:k"+p, "sha256:b"+p))
	}
	got := regionsOf(&Result{Injections: injections}, "r1")
	var paths []string
	for _, r := range got {
		paths = append(paths, r.Path)
	}
	if !reflect.DeepEqual(paths, []string{"z.rb", "a.rb", "m.rb"}) {
		t.Errorf("paths = %v, want the order of Injections", paths)
	}

	if got := regionsOf(&Result{Injections: injections}, "nobody"); got != nil {
		t.Errorf("a record with no regions: %+v, want nil", got)
	}
	if got := regionsOf(&Result{}, "r1"); got != nil {
		t.Errorf("a Result with no injections: %+v, want nil", got)
	}
}

// Capture.

func TestCaptureRegionsFillsEachRecordsRegionsWithoutSortingThem(t *testing.T) {
	result := &Result{
		Selections: []Selection{{RecordID: "r1"}, {RecordID: "r2"}, {RecordID: "r3"}},
		Injections: []inject.Result{
			v2Injection("r1", "z.rb", "addField", "", inject.TierOwned,
				map[string]any{"f": "1"}, map[string]any{"f": "1"}, "sha256:k1", "sha256:b1"),
			v2Injection("r2", "a.rb", "addField", "", inject.TierOwned,
				map[string]any{"f": "2"}, map[string]any{"f": "2"}, "sha256:k2", "sha256:b2"),
			v2Injection("r1", "a.rb", "addField", "", inject.TierOwned,
				map[string]any{"f": "3"}, map[string]any{"f": "3"}, "sha256:k3", "sha256:b3"),
		},
	}

	got := Capture(result)
	if len(got.Records) != 3 {
		t.Fatalf("Capture holds %d records, want 3", len(got.Records))
	}
	for _, rec := range got.Records {
		if want := regionsOf(result, rec.RecordID); !reflect.DeepEqual(rec.Regions, want) {
			t.Errorf("%s: Regions = %+v, want regionsOf's %+v", rec.RecordID, rec.Regions, want)
		}
	}

	// Not sorted: Capture keeps Injections' order.
	var r1 []string
	for _, r := range got.Records[0].Regions {
		r1 = append(r1, r.Path)
	}
	if !reflect.DeepEqual(r1, []string{"z.rb", "a.rb"}) {
		t.Errorf("r1 regions in order %v, want Injections' order [z.rb a.rb]", r1)
	}
	if got.Records[2].Regions != nil {
		t.Errorf("a record with no regions has Regions %+v, want nil", got.Records[2].Regions)
	}
}

func TestCaptureOfARealRunCarriesItsRegions(t *testing.T) {
	cfg := config(t)
	result, _ := v2Run(t, cfg, fields("email", "phone"))

	got := Capture(result)
	if len(got.Records) != 1 || len(got.Records[0].Regions) != 2 {
		t.Fatalf("Capture = %+v, want one record with two regions", got.Records)
	}
	byField := map[string]recording.Region{}
	for _, r := range got.Records[0].Regions {
		byField[r.ID["field"].(string)] = r
	}
	if r := byField["email"]; r.KwargsHash != hashKwargsEmail || r.Body != hashBodyEmail || r.Tier != "owned" {
		t.Errorf("email region = %+v", r)
	}
	if r := byField["phone"]; r.KwargsHash != hashKwargsPhone || r.Body != hashBodyPhone {
		t.Errorf("phone region = %+v", r)
	}
}

// WriteRecordings.

func TestWriteRecordingsOfAResultWithNoClassificationsWritesNothing(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "recordings")
	result := &Result{Selections: []Selection{{RecordID: "r1"}}}

	if err := WriteRecordings(result, dir); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("the directory was created for a run that classified nothing: %v", err)
	}
}

func TestWriteRecordingsWritesTheRecordingOfASelectedRecord(t *testing.T) {
	cfg, dir := v2Grow(t)

	rec := v2ReadRecord(t, dir, v2RecordA)
	if rec.RecordID != v2RecordA {
		t.Fatalf("RecordID = %q", rec.RecordID)
	}
	if len(rec.Phases) != 1 || rec.Phases[0].Name != recording.DefaultPhase || len(rec.Phases[0].Invocations) != 1 {
		t.Errorf("phases = %+v, want the selection in the default phase", rec.Phases)
	}
	if len(rec.Files) != 1 || rec.Files[0].Path != "app/models/user.rb" || rec.Files[0].Package != "rails" {
		t.Errorf("files = %+v", rec.Files)
	}
	if len(rec.Regions) != 1 {
		t.Fatalf("regions = %+v, want one", rec.Regions)
	}

	got := rec.Regions[0]
	marker := v2Regions(t, cfg, "app/models/user.rb")[0].Marker
	if got.Path != "app/models/user.rb" || got.Action != "addField" || got.Variant != "" || got.Tier != "owned" {
		t.Errorf("region = %+v", got)
	}
	if jsonOf(t, got.ID) != `{"field":"email","model":"user"}` || jsonOf(t, got.Kwargs) != `{"field":"email","model":"user"}` {
		t.Errorf("id = %v, kwargs = %v", got.ID, got.Kwargs)
	}
	if got.KwargsHash != hashKwargsEmail || got.Body != hashBodyEmail {
		t.Errorf("hashes = %q, %q; want the kwargs hash and body hash of the email region", got.KwargsHash, got.Body)
	}
	// The entry says what is on the marker.
	if got.Pkg != marker.Pkg || got.KwargsHash != marker.KwargsHash || got.Body != marker.Body {
		t.Errorf("entry %+v does not match the marker %+v", got, marker)
	}
	if rec.Gen == "" || rec.Gen != marker.Gen {
		t.Errorf("gen = %q, marker gen = %q", rec.Gen, marker.Gen)
	}
}

func TestWriteRecordingsRewritesAReselectedRecord(t *testing.T) {
	cfg, dir := v2Grow(t)
	editRecord(t, cfg, "Something else.")

	result, client := v2Run(t, cfg, fields("phone"))
	if client.calls != 1 || v2Class(t, result, v2RecordA).Mode != ModeReselect {
		t.Fatalf("calls = %d, classifications = %+v; want one re-select", client.calls, result.Classifications)
	}
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}

	rec := v2ReadRecord(t, dir, v2RecordA)
	if len(rec.Regions) != 1 || rec.Regions[0].KwargsHash != hashKwargsPhone || rec.Regions[0].Body != hashBodyPhone {
		t.Errorf("regions = %+v, want only the phone region the re-selection left", rec.Regions)
	}
	if rec.Gen != result.Records.Records[0].Gen {
		t.Errorf("gen = %q, want the record's current gen %q", rec.Gen, result.Records.Records[0].Gen)
	}
	if got := jsonOf(t, rec.Phases[0].Invocations[0].Kwargs); got != `{"field":"phone","model":"user"}` {
		t.Errorf("phase invocation kwargs = %s, want the new selection", got)
	}
}

func TestWriteRecordingsLeavesSkippedEditedAndNoneRecordsAlone(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "recordings")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Not valid JSON: a write would replace it and a read would fail on it, so
	// either way the test sees it.
	for _, id := range []string{"r-skip", "r-edited", "r-none"} {
		if err := os.WriteFile(filepath.Join(dir, id+".json"), []byte("keep"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	result := &Result{
		Classifications: []Classification{
			{RecordID: "r-skip", Mode: ModeSkip},
			{RecordID: "r-edited", Mode: ModeEdited},
			{RecordID: "r-none", Mode: ModeNone},
		},
		Selections: []Selection{{RecordID: "r-skip"}, {RecordID: "r-edited"}, {RecordID: "r-none"}},
	}

	if err := WriteRecordings(result, dir); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}
	for _, id := range []string{"r-skip", "r-edited", "r-none"} {
		if got := readFile(t, filepath.Join(dir, id+".json")); got != "keep" {
			t.Errorf("%s was changed to %q", id, got)
		}
	}
	if got := tree(t, dir); len(got) != 3 {
		t.Errorf("directory holds %v, want the three files only", got)
	}
}

func TestWriteRecordingsAfterASkippedGrowChangesNothing(t *testing.T) {
	cfg, dir := v2Grow(t)
	ageTree(t, cfg.Output)
	before := snap(t, cfg.Output)

	result, client := v2Run(t, cfg, fields("email"))
	if client.calls != 0 || v2Class(t, result, v2RecordA).Mode != ModeSkip {
		t.Fatalf("calls = %d, classifications = %+v; want a skip", client.calls, result.Classifications)
	}
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}
	sameSnapshot(t, before, snap(t, cfg.Output))
}

func TestWriteRecordingsWritesOnlyTheRecordsThatWereSelected(t *testing.T) {
	cfg, client := twoRecordConfig(t)
	client.sequence = []string{
		fields("email"),
		`{"invocations": [{"action": "addField", "kwargs": {"model": "order", "field": "total"}}]}`,
	}
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("first grow: %v", err)
	}
	dir := v2Dir(cfg)
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}
	if got := tree(t, dir); !reflect.DeepEqual(got, []string{"prov-2026-aaaaaaaa.json", "prov-2026-bbbbbbbb.json"}) {
		t.Fatalf("recordings = %v", got)
	}
	if rec := v2ReadRecord(t, dir, "prov-2026-bbbbbbbb"); len(rec.Regions) != 1 || rec.Regions[0].KwargsHash != hashKwargsTotal || rec.Regions[0].Body != hashBodyTotal {
		t.Errorf("order record regions = %+v", rec.Regions)
	}

	ageTree(t, cfg.Output)
	before := snap(t, cfg.Output)
	editRecord(t, cfg, "Something else.")
	again, _ := v2Run(t, cfg, fields("phone"))
	if v2Class(t, again, "prov-2026-aaaaaaaa").Mode != ModeReselect || v2Class(t, again, "prov-2026-bbbbbbbb").Mode != ModeSkip {
		t.Fatalf("classifications = %+v", again.Classifications)
	}
	if err := WriteRecordings(again, dir); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}

	after := snap(t, cfg.Output)
	b := ".sedum/recordings/prov-2026-bbbbbbbb.json"
	if before[b].content != after[b].content || !before[b].mtime.Equal(after[b].mtime) {
		t.Errorf("the skipped record's recording was written")
	}
	a := ".sedum/recordings/prov-2026-aaaaaaaa.json"
	if before[a].content == after[a].content {
		t.Errorf("the re-selected record's recording was not rewritten")
	}
}

// A re-rendered record keeps what its recording says besides the regions it
// just rewrote.
func TestWriteRecordingsOfAReRenderedRecordKeepsItsRecordingAndReplacesTheRewrittenEntries(t *testing.T) {
	cfg, dir := v2Grow(t)

	bogus := recording.Region{
		Path: "app/models/user.rb", Action: "addField", Tier: "owned", Pkg: "sha256:old",
		ID:         map[string]any{"field": "kept-entry", "model": "user"},
		Kwargs:     map[string]any{"field": "kept-entry", "model": "user"},
		KwargsHash: "sha256:kept", Body: "sha256:kept",
	}
	keptPhase := []recording.Invocation{{Action: "addField", Kwargs: map[string]any{"model": "user", "field": "kept-phase-field"}}}
	v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) {
		r.CreatedAt = "2026-02-03"
		r.Gen = "sha256:recorded-gen"
		r.Phases = []recording.Phase{{Name: recording.DefaultPhase, Invocations: keptPhase}}
		r.Regions = append(r.Regions, bogus)
	})
	before := v2ReadRecord(t, dir, v2RecordA)
	filesBefore, oldPkg := before.Files, before.Regions[0].Pkg

	editTemplate(t, cfg, "attribute :{{field}}, :string\n")
	result, client := v2Run(t, cfg, fields("email"))
	if client.calls != 0 || v2Class(t, result, v2RecordA).Mode != ModeRerender {
		t.Fatalf("calls = %d, classifications = %+v; want a re-render", client.calls, result.Classifications)
	}
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}

	rec := v2ReadRecord(t, dir, v2RecordA)
	if rec.Gen != "sha256:recorded-gen" || rec.CreatedAt != "2026-02-03" {
		t.Errorf("gen = %q, created_at = %q; the recording's own values must be kept", rec.Gen, rec.CreatedAt)
	}
	if !reflect.DeepEqual(rec.Files, filesBefore) {
		t.Errorf("files = %+v, want %+v", rec.Files, filesBefore)
	}
	if len(rec.Phases) != 1 || !reflect.DeepEqual(rec.Phases[0].Invocations, keptPhase) {
		t.Errorf("phases = %+v, want the recording's own", rec.Phases)
	}

	marker := v2Regions(t, cfg, "app/models/user.rb")[0].Marker
	byField := map[string]recording.Region{}
	for _, r := range rec.Regions {
		byField[r.ID["field"].(string)] = r
	}
	if len(rec.Regions) != 2 {
		t.Fatalf("regions = %+v, want the rewritten email entry and the kept one", rec.Regions)
	}
	if got := byField["kept-entry"]; !reflect.DeepEqual(got, bogus) {
		t.Errorf("the entry of a region not rewritten changed: %+v", got)
	}
	email := byField["email"]
	if email.KwargsHash != hashKwargsEmail || email.Body != hashBodyEmailT {
		t.Errorf("email entry hashes = %q, %q; want the new body %q", email.KwargsHash, email.Body, hashBodyEmailT)
	}
	if email.Pkg != marker.Pkg || email.Pkg == oldPkg {
		t.Errorf("email entry pkg = %q, marker pkg = %q, old pkg = %q; want the new marker's", email.Pkg, marker.Pkg, oldPkg)
	}
}

func TestWriteRecordingsOfAReRenderedV1RecordWithNoRecordingBuildsOneFromCapture(t *testing.T) {
	cfg, dir := v2Grow(t)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	v2ToV1(t, cfg)

	editTemplate(t, cfg, "attribute :{{field}}, :string\n")
	result, client := v2Run(t, cfg, fields("email"))
	if client.calls != 0 || v2Class(t, result, v2RecordA).Mode != ModeRerender {
		t.Fatalf("calls = %d, classifications = %+v; want a re-render of the v1 region", client.calls, result.Classifications)
	}
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}

	rec := v2ReadRecord(t, dir, v2RecordA)
	if len(rec.Files) != 1 || rec.Files[0].Path != "app/models/user.rb" || rec.Files[0].Package != "rails" {
		t.Errorf("files = %+v", rec.Files)
	}
	if len(rec.Phases) != 1 || rec.Phases[0].Name != recording.DefaultPhase || len(rec.Phases[0].Invocations) != 0 {
		t.Errorf("phases = %+v, want one empty default phase", rec.Phases)
	}
	if rec.Gen != result.Records.Records[0].Gen {
		t.Errorf("gen = %q, want the record's gen", rec.Gen)
	}
	if len(rec.Regions) != 1 {
		t.Fatalf("regions = %+v, want the entry of the v1 region that was re-rendered", rec.Regions)
	}
	got := rec.Regions[0]
	if jsonOf(t, got.Kwargs) != `{"field":"email","model":"user"}` || jsonOf(t, got.ID) != `{"field":"email","model":"user"}` ||
		got.KwargsHash != hashKwargsEmail || got.Body != hashBodyEmailT || got.Action != "addField" {
		t.Errorf("entry = %+v", got)
	}
}

func TestWriteRecordingsReportsAnUnreadableRecordingOfAReRenderedRecord(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "r1.json"), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	result := &Result{
		Classifications: []Classification{{RecordID: "r1", Mode: ModeRerender}},
		Selections:      []Selection{{RecordID: "r1"}},
	}

	err := WriteRecordings(result, dir)
	if err == nil {
		t.Fatal("WriteRecordings read an undecodable recording")
	}
	if !strings.HasPrefix(err.Error(), "reading recording for record r1: ") {
		t.Errorf("error = %q, want the prefix %q", err.Error(), "reading recording for record r1: ")
	}
}

func TestWriteRecordingsReportsAWriteFailure(t *testing.T) {
	dir := t.TempDir()
	// A directory where the file must go.
	if err := os.MkdirAll(filepath.Join(dir, "r1.json", "keep"), 0o755); err != nil {
		t.Fatal(err)
	}
	result := &Result{
		Classifications: []Classification{{RecordID: "r1", Mode: ModeSelect}},
		Selections:      []Selection{{RecordID: "r1"}},
	}

	err := WriteRecordings(result, dir)
	if err == nil {
		t.Fatal("WriteRecordings replaced a directory with a file")
	}
	if !strings.HasPrefix(err.Error(), "writing recording for record r1: ") {
		t.Errorf("error = %q, want the prefix %q", err.Error(), "writing recording for record r1: ")
	}
}

// v2ToV1 rewrites the opening marker of every region of user.rb as the v1
// marker Sedum wrote before marker v2, with the same gen and pkg, so a test can
// grow from v1 markers without a v1 writer. A v2 marker's id holds every kwarg
// of the addField fixture, whose kwargs are all required and none declared as
// identity, so it is also the v1 kwargs.
func v2ToV1(t *testing.T, cfg Config) {
	t.Helper()
	path := "app/models/user.rb"
	content := readFile(t, filepath.Join(cfg.Output, filepath.FromSlash(path)))
	regions := v2Regions(t, cfg, path)

	// Rewrite from the last region to the first so offsets stay valid.
	sort.Slice(regions, func(i, j int) bool { return regions[i].Start > regions[j].Start })
	for _, r := range regions {
		rest := content[r.Start:]
		indent := rest[:len(rest)-len(strings.TrimLeft(rest, " \t"))]
		end := r.Start + strings.Index(rest, "\n")
		line := indent + "# sedum:" + r.Marker.Label() + " " + v1Attrs(t, r.Marker)
		content = content[:r.Start] + line + content[end:]
	}
	v2WriteFile(t, cfg, path, content)
}

// v1Attrs encodes a marker's attributes as a v1 marker would: tier, record,
// gen, pkg, then every kwarg. For a v2 marker the kwargs are read from the
// action's identity kwargs, which for the addField fixture are all of them.
func v1Attrs(t *testing.T, m inject.Marker) string {
	t.Helper()
	attrs := struct {
		Tier   string         `json:"tier"`
		Record string         `json:"record,omitempty"`
		Gen    string         `json:"gen,omitempty"`
		Pkg    string         `json:"pkg,omitempty"`
		Kwargs map[string]any `json:"kwargs,omitempty"`
	}{string(m.Tier), m.Record, m.Gen, m.Pkg, m.Kwargs}
	return jsonOf(t, attrs)
}

func TestWriteReplayRecordingsKeepsWhatWasReplayedAndFillsRegionsAndGen(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "recordings")
	stale := recording.Region{Path: "old.rb", Action: "addField", ID: map[string]any{"f": "stale"}}
	replayed := recording.Recording{
		SedumVersion: "9.9.9",
		Packages: map[string]recording.Package{
			"rails": {Extensions: []string{".rb"}},
			"other": {Extensions: []string{".x"}},
		},
		Variables: map[string]string{"service": "orders"},
		Complete:  true,
		Prompt:    "not carried",
		Records: []recording.Record{
			{
				RecordID: "r1", Gen: "sha256:old-gen", CreatedAt: "2026-01-02",
				Files:  []recording.File{{Path: "a.rb", Package: "rails", Template: "{name}.rb"}},
				Phases: []recording.Phase{{Name: recording.DefaultPhase, Invocations: []recording.Invocation{{Action: "addField", Kwargs: map[string]any{"f": "1"}}}}},
			},
			{
				RecordID: "r2", Gen: "sha256:r2-gen",
				Files:   []recording.File{{Path: "b.rb", Package: "rails", Template: "{name}.rb"}},
				Phases:  []recording.Phase{{Name: recording.DefaultPhase}},
				Regions: []recording.Region{stale},
			},
		},
	}
	result := &Result{
		Restamps: []Restamp{{RecordID: "r1", Old: "sha256:old-gen", New: "sha256:restamped"}},
		Injections: []inject.Result{
			v2Injection("r1", "a.rb", "addField", "", inject.TierOwned,
				map[string]any{"f": "1"}, map[string]any{"f": "1"}, "sha256:k1", "sha256:b1"),
			v2Injection("other-record", "z.rb", "addField", "", inject.TierOwned,
				map[string]any{"f": "z"}, map[string]any{"f": "z"}, "sha256:kz", "sha256:bz"),
		},
	}

	if err := WriteReplayRecordings(replayed, result, dir); err != nil {
		t.Fatalf("WriteReplayRecordings: %v", err)
	}

	if got := tree(t, dir); !reflect.DeepEqual(got, []string{"r1.json", "r2.json"}) {
		t.Fatalf("recordings = %v, want one file per replayed record", got)
	}

	r1file, err := recording.ReadRecord(dir, "r1")
	if err != nil {
		t.Fatal(err)
	}
	if r1file.SedumVersion != "9.9.9" || !reflect.DeepEqual(r1file.Variables, map[string]string{"service": "orders"}) {
		t.Errorf("r1 file: version %q, variables %v", r1file.SedumVersion, r1file.Variables)
	}
	if want := map[string]recording.Package{"rails": {Extensions: []string{".rb"}}}; !reflect.DeepEqual(r1file.Packages, want) {
		t.Errorf("r1 file packages = %v, want %v", r1file.Packages, want)
	}
	r1 := r1file.Records[0]
	if r1.Gen != "sha256:restamped" {
		t.Errorf("r1 gen = %q, want the restamp's New", r1.Gen)
	}
	if r1.CreatedAt != "2026-01-02" || len(r1.Files) != 1 || len(r1.Phases) != 1 || len(r1.Phases[0].Invocations) != 1 {
		t.Errorf("r1 lost what was replayed: %+v", r1)
	}
	if len(r1.Regions) != 1 || r1.Regions[0].Path != "a.rb" || r1.Regions[0].KwargsHash != "sha256:k1" || r1.Regions[0].Body != "sha256:b1" {
		t.Errorf("r1 regions = %+v, want the regions of the replay's injections for r1 only", r1.Regions)
	}

	r2 := v2ReadRecord(t, dir, "r2")
	if r2.Gen != "sha256:r2-gen" {
		t.Errorf("r2 gen = %q, want the replayed gen (no restamp names it)", r2.Gen)
	}
	if len(r2.Regions) != 0 {
		t.Errorf("r2 regions = %+v, want regionsOf's (none)", r2.Regions)
	}
}

func TestWriteReplayRecordingsReportsAWriteFailure(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "r1.json", "keep"), 0o755); err != nil {
		t.Fatal(err)
	}
	replayed := recording.Recording{
		SedumVersion: recording.Version,
		Packages:     map[string]recording.Package{},
		Records:      []recording.Record{{RecordID: "r1"}},
	}

	err := WriteReplayRecordings(replayed, &Result{}, dir)
	if err == nil {
		t.Fatal("WriteReplayRecordings replaced a directory with a file")
	}
	if !strings.HasPrefix(err.Error(), "writing recording for record r1: ") {
		t.Errorf("error = %q, want the prefix %q", err.Error(), "writing recording for record r1: ")
	}
}

// Replay then recordings: what a replay writes, a later grow can classify
// against, so a replayed tree is current.
func TestAReplayedRecordingsFilesVouchForTheRegionsItWrote(t *testing.T) {
	cfg, dir := v2Grow(t)
	original := v2ReadRecord(t, dir, v2RecordA)

	replayOut := t.TempDir()
	replayDir := filepath.Join(replayOut, ".sedum", "recordings")
	full := recording.Recording{
		SedumVersion: recording.Version,
		Packages:     map[string]recording.Package{"rails": {Extensions: []string{".rb"}}},
		Records:      []recording.Record{original},
	}
	result, err := Replay(full, ReplayConfig{Generators: cfg.Generators, Output: replayOut, Records: cfg.Records})
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if err := WriteReplayRecordings(full, result, replayDir); err != nil {
		t.Fatalf("WriteReplayRecordings: %v", err)
	}

	written := v2ReadRecord(t, replayDir, v2RecordA)
	if len(written.Regions) != 1 || written.Regions[0].KwargsHash != hashKwargsEmail || written.Regions[0].Body != hashBodyEmail {
		t.Errorf("regions = %+v", written.Regions)
	}

	// A grow over the replayed tree, reading the replay's recordings, skips.
	replayCfg := cfg
	replayCfg.Output = replayOut
	check := v2Check(t, replayCfg)
	if got := v2Class(t, check, v2RecordA); got.Mode != ModeSkip {
		t.Errorf("a grow over a replayed tree: %+v, want skip", got)
	}
}
