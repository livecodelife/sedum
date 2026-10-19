package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result_test.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func testSamples() []Sample {\n\treturn []Sample{\n\t\t{Task: \"b\", Variant: VariantCode, N: 1, Compiled: true, Passed: true, Calls: 1, PromptTokens: 100, CompletionTokens: 10},\n\t\t{Task: \"a\", Variant: VariantCode, N: 1, Compiled: true, Passed: false, Calls: 2, Rejected: 1, Repairs: 1, PromptTokens: 200, CompletionTokens: 20},\n\t\t{Task: \"a\", Variant: VariantBehaviour, N: 1, Compiled: false, Calls: 3, Rejected: 3, PromptTokens: 50, CompletionTokens: 5},\n\t\t{Task: \"a\", Variant: VariantBehaviour, N: 2, Compiled: true, Passed: true, Calls: 1, PromptTokens: 50, CompletionTokens: 5},\n\t}\n}","file":"evals/design/result_test.go","name":"testSamples"}}
func testSamples() []Sample {
	return []Sample{
		{Task: "b", Variant: VariantCode, N: 1, Compiled: true, Passed: true, Calls: 1, PromptTokens: 100, CompletionTokens: 10},
		{Task: "a", Variant: VariantCode, N: 1, Compiled: true, Passed: false, Calls: 2, Rejected: 1, Repairs: 1, PromptTokens: 200, CompletionTokens: 20},
		{Task: "a", Variant: VariantBehaviour, N: 1, Compiled: false, Calls: 3, Rejected: 3, PromptTokens: 50, CompletionTokens: 5},
		{Task: "a", Variant: VariantBehaviour, N: 2, Compiled: true, Passed: true, Calls: 1, PromptTokens: 50, CompletionTokens: 5},
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAggregate(t *testing.T) {\n\tcells := Aggregate(testSamples())\n\tif len(cells) != 3 {\n\t\tt.Fatalf(\"got %d cells, want 3: %+v\", len(cells), cells)\n\t}\n\tif cells[0].Task != \"a\" || cells[0].Variant != VariantBehaviour || cells[1].Variant != VariantCode || cells[2].Task != \"b\" {\n\t\tt.Fatalf(\"cells out of order: %+v\", cells)\n\t}\n\tc := cells[0]\n\tif c.Samples != 2 || c.Compiled != 1 || c.Passed != 1 || c.Calls != 4 || c.Rejected != 3 || c.Tokens() != 110 {\n\t\tt.Fatalf(\"cell = %+v\", c)\n\t}\n}","file":"evals/design/result_test.go","name":"TestAggregate"}}
func TestAggregate(t *testing.T) {
	cells := Aggregate(testSamples())
	if len(cells) != 3 {
		t.Fatalf("got %d cells, want 3: %+v", len(cells), cells)
	}
	if cells[0].Task != "a" || cells[0].Variant != VariantBehaviour || cells[1].Variant != VariantCode || cells[2].Task != "b" {
		t.Fatalf("cells out of order: %+v", cells)
	}
	c := cells[0]
	if c.Samples != 2 || c.Compiled != 1 || c.Passed != 1 || c.Calls != 4 || c.Rejected != 3 || c.Tokens() != 110 {
		t.Fatalf("cell = %+v", c)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestReportTotalsJSONAndTable(t *testing.T) {\n\tr := NewReport(\"abc1234\", \"sedum\", \"claude-haiku-5-5\", 2, map[string]string{\"a\": \"h1\"}, testSamples())\n\tif r.Totals.Samples != 4 || r.Totals.Passed != 2 || r.Totals.Rejected != 4 || r.Totals.Repairs != 1 || r.Totals.Tokens() != 440 {\n\t\tt.Fatalf(\"totals = %+v\", r.Totals)\n\t}\n\tdata, err := r.JSON()\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif !strings.HasSuffix(string(data), \"\\n\") {\n\t\tt.Fatal(\"JSON must end with a newline\")\n\t}\n\tvar back Report\n\tif err := json.Unmarshal(data, &back); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif back.Commit != \"abc1234\" || back.K != 2 || len(back.Samples) != 4 || back.Totals.Samples != 4 {\n\t\tt.Fatalf(\"round trip = %+v\", back)\n\t}\n\ttable := r.Table()\n\tfor _, want := range []string{\"claude-haiku-5-5\", \"K=2\", \"direction, not a number\", \"behaviour\", \"signatures\", \"code\", \"1/2\", \"total\"} {\n\t\tif !strings.Contains(table, want) {\n\t\t\tt.Errorf(\"table lacks %q:\\n%s\", want, table)\n\t\t}\n\t}\n}","file":"evals/design/result_test.go","name":"TestReportTotalsJSONAndTable"}}
func TestReportTotalsJSONAndTable(t *testing.T) {
	r := NewReport("abc1234", "sedum", "claude-haiku-5-5", 2, map[string]string{"a": "h1"}, testSamples())
	if r.Totals.Samples != 4 || r.Totals.Passed != 2 || r.Totals.Rejected != 4 || r.Totals.Repairs != 1 || r.Totals.Tokens() != 440 {
		t.Fatalf("totals = %+v", r.Totals)
	}
	data, err := r.JSON()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(string(data), "\n") {
		t.Fatal("JSON must end with a newline")
	}
	var back Report
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatal(err)
	}
	if back.Commit != "abc1234" || back.K != 2 || len(back.Samples) != 4 || back.Totals.Samples != 4 {
		t.Fatalf("round trip = %+v", back)
	}
	table := r.Table()
	for _, want := range []string{"claude-haiku-5-5", "K=2", "direction, not a number", "behaviour", "signatures", "code", "1/2", "total"} {
		if !strings.Contains(table, want) {
			t.Errorf("table lacks %q:\n%s", want, table)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestTaskHash(t *testing.T) {\n\troot := t.TempDir()\n\ttask, err := LoadTask(writeTask(t, root, \"ok\", validTask))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\th1, err := TaskHash(task)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\th2, _ := TaskHash(task)\n\tif h1 != h2 || len(h1) != 12 {\n\t\tt.Fatalf(\"hash = %q then %q\", h1, h2)\n\t}\n\twriteTask(t, root, \"ok\", strings.Replace(validTask, \"code: c\", \"code: d\", 1))\n\tchanged, err := LoadTask(task.Dir)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif h3, _ := TaskHash(changed); h3 == h1 {\n\t\tt.Fatal(\"a changed task must change its hash\")\n\t}\n}","file":"evals/design/result_test.go","name":"TestTaskHash"}}
func TestTaskHash(t *testing.T) {
	root := t.TempDir()
	task, err := LoadTask(writeTask(t, root, "ok", validTask))
	if err != nil {
		t.Fatal(err)
	}
	h1, err := TaskHash(task)
	if err != nil {
		t.Fatal(err)
	}
	h2, _ := TaskHash(task)
	if h1 != h2 || len(h1) != 12 {
		t.Fatalf("hash = %q then %q", h1, h2)
	}
	writeTask(t, root, "ok", strings.Replace(validTask, "code: c", "code: d", 1))
	changed, err := LoadTask(task.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if h3, _ := TaskHash(changed); h3 == h1 {
		t.Fatal("a changed task must change its hash")
	}
}

// /sedum:setFunction
