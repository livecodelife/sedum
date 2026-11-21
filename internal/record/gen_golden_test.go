package record

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var updateGolden = flag.Bool("update-gen-golden", false, "rewrite testdata/gen_golden.json from the current code")

// The golden file holds the gens sedum 0.13.0 computed, so a record written
// since (such as the one for 0.14.0 itself) is not in it: the test checks every
// golden record still loads with the same gen, and a new record has no 0.13.0
// value to pin. A record revised since 0.13.0 has its entry recomputed by the
// 0.13.0 code over the new content (prov-2026-5f1e4f0a, whose constraint was
// added in this change), so the entry still comes from the old implementation.
//
// goldenDirs are sedum's own record directories. The golden file holds each
// record's gen as sedum 0.13.0 computed it, before gen_hash covered shape. A
// record with no shape must keep that exact hash: provenance leaves the shape
// out of the hash input when it is empty, and this test proves sedum's markers
// do not move for such records.
var goldenDirs = []string{
	"../../provenance",
	"../../evals/testdata/grocery-api/records",
	"../../evals/testdata/todo-api/records",
	"../../evals/testdata/todo-rails/records",
	"../../evals/lifecycle/fixtures/step1/records",
	"../../evals/lifecycle/fixtures/step3/records",
	"../../evals/lifecycle/fixtures/step11/records",
}

func currentGens(t *testing.T) map[string]string {
	t.Helper()
	got := map[string]string{}
	for _, d := range goldenDirs {
		set, _, err := Load(d, Options{})
		if err != nil {
			t.Fatalf("load %s: %v", d, err)
		}
		for _, r := range set.Known {
			got[filepath.ToSlash(filepath.Join(strings.TrimPrefix(filepath.ToSlash(d), "../../"), r.ID))] = r.Gen
		}
	}
	return got
}

func TestGenOfRecordsWithoutShapeIsUnchanged(t *testing.T) {
	got := currentGens(t)
	const golden = "testdata/gen_golden.json"
	if *updateGolden {
		b, _ := json.MarshalIndent(got, "", "  ")
		if err := os.WriteFile(golden, append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	raw, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	var want map[string]string
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	if len(want) < 100 {
		t.Fatalf("golden holds %d records; expected sedum's whole corpus", len(want))
	}
	for k, w := range want {
		if g, ok := got[k]; !ok {
			t.Errorf("%s: record no longer loads", k)
		} else if g != w {
			t.Errorf("%s: gen %s, golden (v0.13.0) %s", k, g, w)
		}
	}
}
