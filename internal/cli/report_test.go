package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// The two JSON reports are a surface: conformance/formats documents them for
// consumers that cannot link Sedum, and this test holds real output to those
// documents key for key, in both directions (prov-2026-bf01cd6b).

const railsRecording = `{"sedum_version":"0.8.0","packages":{"rails":{"extensions":[".rb"]}},
 "records":[{"record_id":"prov-2026-aaaaaaaa","files":[{"path":"app/controllers/users_controller.rb","package":"rails","template":"app/controllers/{name}_controller.rb","captures":{"name":"users"}}],
 "phases":[{"name":"default","invocations":[{"action":"addBeforeFilter","kwargs":{"controller":"users","filter":"authenticate"}}]}]}]}`

func example(t *testing.T, name string) any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "conformance", "formats", name))
	if err != nil {
		t.Fatal(err)
	}
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("%s is not JSON: %v", name, err)
	}
	return v
}

// sameKeys reports every object key present in one document and not the other,
// descending through objects and through the first element of arrays both
// sides populate.
func sameKeys(path string, got, want any) []string {
	var out []string
	switch w := want.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			return []string{path + ": not an object"}
		}
		var keys []string
		for k := range w {
			keys = append(keys, k)
		}
		for k := range g {
			if _, ok := w[k]; !ok {
				out = append(out, path+"."+k+": emitted but not documented")
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			gv, ok := g[k]
			if !ok {
				out = append(out, path+"."+k+": documented but not emitted")
				continue
			}
			if k == "identity" {
				continue // an action's own kwargs: free-form by design
			}
			out = append(out, sameKeys(path+"."+k, gv, w[k])...)
		}
	case []any:
		g, ok := got.([]any)
		if !ok {
			return []string{path + ": not an array"}
		}
		if len(w) > 0 && len(g) > 0 {
			if _, isObj := w[0].(map[string]any); isObj {
				out = append(out, sameKeys(path+"[0]", g[0], w[0])...)
			}
		}
	}
	return out
}

func decodeJSON(t *testing.T, out string) any {
	t.Helper()
	var v any
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("output is not one JSON document: %v\n%s", err, out)
	}
	return v
}

func TestJSONReportsMatchTheirDocumentedFormats(t *testing.T) {
	dir := t.TempDir()
	gens := filepath.Join("..", "..", "testdata", "generators")
	output := filepath.Join(dir, "out")
	rec := filepath.Join(dir, "rec.json")
	if err := os.WriteFile(rec, []byte(railsRecording), 0o644); err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(dir, "run.log")

	out, err := exec(t, "grow", "--generators", gens, "--output", output, "--execute", rec, "--log", log, "--json")
	if err != nil {
		t.Fatalf("grow --json: %v\n%s", err, out)
	}
	grow := decodeJSON(t, out)
	if problems := sameKeys("grow", grow, example(t, "grow.example.json")); len(problems) > 0 {
		// orphans and stamped are empty here, so only their presence is held.
		t.Errorf("grow --json differs from conformance/formats:\n  %s", strings.Join(problems, "\n  "))
	}
	records := grow.(map[string]any)["records"].([]any)
	if len(records) != 1 || records[0].(map[string]any)["changed"] != true {
		t.Errorf("records = %v, want the one record, changed", records)
	}

	out, err = exec(t, "regions", "--generators", gens, "--output", output, "--json")
	if err != nil {
		t.Fatalf("regions --json: %v\n%s", err, out)
	}
	if problems := sameKeys("regions", decodeJSON(t, out), example(t, "regions.example.json")); len(problems) > 0 {
		t.Errorf("regions --json differs from conformance/formats:\n  %s", strings.Join(problems, "\n  "))
	}

	// Replaying the same recording again rewrites the region identically,
	// which is not a change.
	out, err = exec(t, "grow", "--generators", gens, "--output", output, "--execute", rec, "--log", log, "--json")
	if err != nil {
		t.Fatalf("second grow --json: %v", err)
	}
	again := decodeJSON(t, out).(map[string]any)["records"].([]any)[0].(map[string]any)
	if again["changed"] != false {
		t.Errorf("an identical rewrite was reported as a change: %v", again)
	}
	regions := again["regions"].([]any)
	if len(regions) != 1 || regions[0].(map[string]any)["change"] != "unchanged" {
		t.Errorf("regions = %v, want one unchanged", regions)
	}
}
