package gitcheck

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// The SARIF unit parses provenance logs, merges their results into one run and
// serialises the result.

const rulesAB = `[{"id":"SCOPE001","name":"Scope"},{"id":"STATE001","name":"State"}]`

func logText(rules string, results string, artifacts string) []byte {
	return []byte(`{"$schema":"x","version":"2.1.0","runs":[{"tool":{"driver":{"name":"provenance","rules":` + rules + `}},"results":` + results + `,"artifacts":` + artifacts + `}]}`)
}

func result(rule string, index int, msg, file, record string) string {
	loc := ""
	if file != "" {
		loc = `,"locations":[{"physicalLocation":{"artifactLocation":{"uri":"` + file + `"}}}]`
	}
	props := ""
	if record != "" {
		props = `,"properties":{"recordId":"` + record + `"}`
	}
	return `{"ruleId":"` + rule + `","ruleIndex":` + itoa(index) + `,"level":"error","message":{"text":"` + msg + `"}` + loc + props + `}`
}

func itoa(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}

func parse(t *testing.T, data []byte) Log {
	t.Helper()
	log, err := ParseLog(data)
	if err != nil {
		t.Fatalf("ParseLog: %v", err)
	}
	return log
}

func TestParseLogReadsResultsFromEveryRunInOrder(t *testing.T) {
	data := []byte(`{"version":"2.1.0","runs":[
		{"tool":{"driver":{"rules":[]}},"results":[` + result("SCOPE001", 0, "one", "a.go", "r1") + `,` + result("STATE001", 1, "two", "", "") + `]},
		{"tool":{"driver":{"rules":[]}},"results":[` + result("TAG001", 0, "three", "b.go", "r2") + `]}]}`)

	results := parse(t, data).Results()

	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	want := []struct {
		run, res int
		rule     string
		msg      string
		file     string
		record   string
	}{{0, 0, "SCOPE001", "one", "a.go", "r1"}, {0, 1, "STATE001", "two", "", ""}, {1, 0, "TAG001", "three", "b.go", "r2"}}
	for i, w := range want {
		r := results[i]
		if r.RunIndex() != w.run || r.ResultIndex() != w.res || r.RuleID() != w.rule || r.MessageText() != w.msg ||
			r.FirstLocationURI() != w.file || r.Property("recordId") != w.record {
			t.Errorf("result %d = run %d res %d %s %q %q %q", i, r.RunIndex(), r.ResultIndex(), r.RuleID(), r.MessageText(), r.FirstLocationURI(), r.Property("recordId"))
		}
	}
}

func TestParseLogFailsOnBytesThatAreNotJSONWithAProvenanceError(t *testing.T) {
	_, err := ParseLog([]byte("not json"))
	if err == nil || !strings.HasPrefix(err.Error(), "provenance:") {
		t.Errorf("err = %v, want a provenance: error", err)
	}
}

func TestParseLogReportsSkewWhenAFieldTheCommandNeedsIsMissing(t *testing.T) {
	for name, data := range map[string]string{
		"no runs":        `{"version":"2.1.0"}`,
		"empty runs":     `{"runs":[]}`,
		"no rule id":     `{"runs":[{"results":[{"message":{"text":"m"}}]}]}`,
		"rule id number": `{"runs":[{"results":[{"ruleId":1,"message":{"text":"m"}}]}]}`,
		"no message":     `{"runs":[{"results":[{"ruleId":"X"}]}]}`,
		"no text":        `{"runs":[{"results":[{"ruleId":"X","message":{}}]}]}`,
	} {
		_, err := ParseLog([]byte(data))
		if err == nil || !strings.HasPrefix(err.Error(), "provenance:") {
			t.Errorf("%s: err = %v, want a provenance: skew error", name, err)
		}
	}
}

func TestAMissingRecordIDIsNotSkew(t *testing.T) {
	log, err := ParseLog(logText(rulesAB, "["+result("SCOPE001", 0, "m", "", "")+"]", "[]"))
	if err != nil {
		t.Fatalf("ParseLog: %v", err)
	}
	if got := log.Results()[0].Property("recordId"); got != "" {
		t.Errorf("Property = %q, want empty", got)
	}
}

type doc struct {
	Version string `json:"version"`
	Runs    []struct {
		Tool struct {
			Driver struct {
				Name  string `json:"name"`
				Rules []struct {
					ID string `json:"id"`
				} `json:"rules"`
			} `json:"driver"`
		} `json:"tool"`
		Results []struct {
			RuleID       string `json:"ruleId"`
			RuleIndex    int    `json:"ruleIndex"`
			Message      struct{ Text string }
			Suppressions []struct {
				Kind          string `json:"kind"`
				Justification string `json:"justification"`
			} `json:"suppressions"`
		} `json:"results"`
		Artifacts []struct {
			Location struct{ URI string } `json:"location"`
		} `json:"artifacts"`
		Properties map[string]any `json:"properties"`
	} `json:"runs"`
}

func merged(t *testing.T, logs []Log, just [][]string, skipped int) doc {
	t.Helper()
	data, err := Merge(logs, just).JSON(skipped)
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}
	var d doc
	if err := json.Unmarshal(data, &d); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, data)
	}
	return d
}

func TestMergeWritesOneRunOfSedumCheckWithEveryResultInCommitOrder(t *testing.T) {
	first := parse(t, logText(rulesAB, "["+result("SCOPE001", 0, "first", "a.go", "r1")+","+result("STATE001", 1, "second", "", "")+"]", `[{"location":{"uri":"b.go"}},{"location":{"uri":"a.go"}}]`))
	second := parse(t, logText(`[{"id":"TAG001","name":"Tag"},{"id":"SCOPE001","name":"Scope"}]`, "["+result("TAG001", 0, "third", "", "")+","+result("SCOPE001", 1, "fourth", "a.go", "")+"]", `[{"location":{"uri":"a.go"}},{"location":{"uri":"c.go"}}]`))

	d := merged(t, []Log{first, second}, [][]string{{"", ""}, {"", ""}}, 3)

	if d.Version != "2.1.0" || len(d.Runs) != 1 {
		t.Fatalf("version %q, runs %d", d.Version, len(d.Runs))
	}
	run := d.Runs[0]
	if run.Tool.Driver.Name != "sedum check" {
		t.Errorf("driver name = %q", run.Tool.Driver.Name)
	}
	var texts []string
	for _, r := range run.Results {
		texts = append(texts, r.Message.Text)
	}
	if strings.Join(texts, ",") != "first,second,third,fourth" {
		t.Errorf("result order = %v", texts)
	}
	var rules []string
	for _, r := range run.Tool.Driver.Rules {
		rules = append(rules, r.ID)
	}
	if strings.Join(rules, ",") != "SCOPE001,STATE001,TAG001" {
		t.Errorf("rules = %v, want the union in order of first appearance", rules)
	}
	for _, r := range run.Results {
		if rules[r.RuleIndex] != r.RuleID {
			t.Errorf("result %q: ruleIndex %d points at %s, want %s", r.Message.Text, r.RuleIndex, rules[r.RuleIndex], r.RuleID)
		}
	}
	var arts []string
	for _, a := range run.Artifacts {
		arts = append(arts, a.Location.URI)
	}
	if strings.Join(arts, ",") != "a.go,b.go,c.go" {
		t.Errorf("artifacts = %v, want one per uri sorted", arts)
	}
	if got := run.Properties["skippedScopePaths"]; got != float64(3) {
		t.Errorf("skippedScopePaths = %v, want 3", got)
	}
}

func TestMergeSuppressesOnlyTheResultsWithAJustification(t *testing.T) {
	log := parse(t, logText(rulesAB, "["+result("STATE001", 1, "waived", "", "r1")+","+result("SCOPE001", 0, "kept", "", "")+"]", "[]"))

	d := merged(t, []Log{log}, [][]string{{"regeneration: abc1234", ""}}, 0)

	res := d.Runs[0].Results
	if len(res[0].Suppressions) != 1 || res[0].Suppressions[0].Kind != "external" || res[0].Suppressions[0].Justification != "regeneration: abc1234" {
		t.Errorf("waived result suppressions = %+v", res[0].Suppressions)
	}
	if res[1].Suppressions != nil {
		t.Errorf("kept result has suppressions: %+v", res[1].Suppressions)
	}
}

func TestMergeKeepsAResultsOtherFieldsUnchanged(t *testing.T) {
	log := parse(t, logText(rulesAB, `[{"ruleId":"SCOPE001","ruleIndex":0,"level":"error","message":{"text":"m"},"properties":{"recordId":"r1","field":"f"},"extra":{"a":[1,2]}}]`, "[]"))

	data, err := Merge([]Log{log}, [][]string{{""}}).JSON(0)
	if err != nil {
		t.Fatal(err)
	}

	var out struct {
		Runs []struct {
			Results []map[string]json.RawMessage `json:"results"`
		} `json:"runs"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	r := out.Runs[0].Results[0]
	compact := func(raw json.RawMessage) string {
		var buf bytes.Buffer
		if err := json.Compact(&buf, raw); err != nil {
			t.Fatal(err)
		}
		return buf.String()
	}
	if compact(r["extra"]) != `{"a":[1,2]}` || compact(r["level"]) != `"error"` || !strings.Contains(compact(r["properties"]), `"field":"f"`) {
		t.Errorf("result changed: extra=%s level=%s properties=%s", r["extra"], r["level"], r["properties"])
	}
}

func TestAnEmptyRangeGivesALogWithZeroResults(t *testing.T) {
	data, err := Merge(nil, nil).JSON(0)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"results": []`) {
		t.Errorf("output = %s, want an empty results array", data)
	}
	var d doc
	if err := json.Unmarshal(data, &d); err != nil || len(d.Runs) != 1 || d.Runs[0].Tool.Driver.Name != "sedum check" {
		t.Errorf("err=%v doc=%+v", err, d)
	}
}
