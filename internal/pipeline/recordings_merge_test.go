package pipeline

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/recording"
)

// A replay writes the regions it wrote and keeps the entries the record's
// recording already holds for regions it did not write, so a region that is in
// the tree and was not part of the replay still has its recording entry.
func TestWriteReplayRecordingsKeepsTheEntriesOfRegionsTheReplayDidNotWrite(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "recordings")
	kept := recording.Region{Path: "a.rb", Action: "addField", ID: map[string]any{"f": "kept"}, KwargsHash: "sha256:kept", Body: "sha256:keptbody"}
	replaced := recording.Region{Path: "a.rb", Action: "addField", ID: map[string]any{"f": "1"}, KwargsHash: "sha256:old", Body: "sha256:oldbody"}
	first := recording.Recording{
		SedumVersion: "9.9.9",
		Packages:     map[string]recording.Package{"rails": {Extensions: []string{".rb"}}},
		Records: []recording.Record{{
			RecordID: "r1", Gen: "sha256:g",
			Files:   []recording.File{{Path: "a.rb", Package: "rails", Template: "{name}.rb"}},
			Phases:  []recording.Phase{{Name: recording.DefaultPhase}},
			Regions: []recording.Region{kept, replaced},
		}},
	}
	if err := recording.WriteRecords(dir, first); err != nil {
		t.Fatal(err)
	}

	replayed := first
	replayed.Records = []recording.Record{first.Records[0]}
	replayed.Records[0].Regions = nil
	result := &Result{Injections: []inject.Result{
		v2Injection("r1", "a.rb", "addField", "", inject.TierOwned,
			map[string]any{"f": "1"}, map[string]any{"f": "1"}, "sha256:new", "sha256:newbody"),
		v2Injection("r1", "a.rb", "addField", "", inject.TierOwned,
			map[string]any{"f": "2"}, map[string]any{"f": "2"}, "sha256:two", "sha256:twobody"),
	}}
	if err := WriteReplayRecordings(replayed, result, dir); err != nil {
		t.Fatalf("WriteReplayRecordings: %v", err)
	}

	got, err := recording.ReadRecord(dir, "r1")
	if err != nil {
		t.Fatal(err)
	}
	var keys []string
	hashes := map[string]string{}
	for _, r := range got.Records[0].Regions {
		k := recording.RegionKey(r.Path, r.Action, r.Variant, r.ID)
		keys = append(keys, k)
		hashes[k] = r.KwargsHash
	}
	want := []string{
		recording.RegionKey("a.rb", "addField", "", map[string]any{"f": "1"}),
		recording.RegionKey("a.rb", "addField", "", map[string]any{"f": "2"}),
		recording.RegionKey("a.rb", "addField", "", map[string]any{"f": "kept"}),
	}
	if !reflect.DeepEqual(keys, want) {
		t.Fatalf("regions = %v, want the two written and the one kept, in region-key order %v", keys, want)
	}
	if hashes[want[0]] != "sha256:new" {
		t.Errorf("a region the replay wrote keeps its old entry: %q", hashes[want[0]])
	}
	if hashes[want[2]] != "sha256:kept" {
		t.Errorf("the entry of a region the replay did not write was dropped or changed: %q", hashes[want[2]])
	}
}
