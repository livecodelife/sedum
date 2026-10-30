package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/result_test.go","path":"encoding/json"},"kwargs":"sha256:d912b68dd0f0058ebe7c9e566ba0f6013fcca82d80452f37c62370cb7aa17c03","body":"sha256:872d7c3e9dee608e79a7f8ef57930d24b5e30f2a57d891a52549179ace45946a"}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/result_test.go","path":"strings"},"kwargs":"sha256:2c29be2bd4168ae22e86b0532ce7863239cfe8210c1bf802f3cc43bd21d557e9","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/result_test.go","path":"testing"},"kwargs":"sha256:484142e46872e9b232382277275387d5b07c47b174de0e41a2b55fda4efa3f71","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/result_test.go","name":"testSamples"},"kwargs":"sha256:411ac32803165f079404a5e0d8cd977ae788b97c80d822237d1b887c1775085d","body":"sha256:98b732646c2d5253fe8ed92979d8816979cf42c81c4b07bff9097a47744c1cd4"}
func testSamples() []Sample {
	return []Sample{
		{Task: "b", Variant: VariantCode, N: 1, Compiled: true, Passed: true, Calls: 1, PromptTokens: 100, CompletionTokens: 10},
		{Task: "a", Variant: VariantCode, N: 1, Compiled: true, Passed: false, Calls: 2, Rejected: 1, Repairs: 1, PromptTokens: 200, CompletionTokens: 20},
		{Task: "a", Variant: VariantBehaviour, N: 1, Compiled: false, Calls: 3, Rejected: 3, PromptTokens: 50, CompletionTokens: 5},
		{Task: "a", Variant: VariantBehaviour, N: 2, Compiled: true, Passed: true, Calls: 1, PromptTokens: 50, CompletionTokens: 5},
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/result_test.go","name":"TestAggregate"},"kwargs":"sha256:5bb1875c62052bc7ffcb07bd0b6191172093a522f139c68a79be4a0c774505e1","body":"sha256:3848b061b8b27f5efa9602e914585a5c737e6a3fa684e6891f8f8f4d3894cc11"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/result_test.go","name":"TestReportTotalsJSONAndTable"},"kwargs":"sha256:ed613a812892adbb75ade1e31094bdd7946075765a41db9c38fd02d4d6e57156","body":"sha256:6e9690210c3ae5f8f04febb5e772adfc2a102c11893334dfdc432e272b3235f5"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/result_test.go","name":"TestTaskHash"},"kwargs":"sha256:dcc2a2f04ca76aa87700d83454251b8ab8e2e53f7aa04ceb99483184af954ef8","body":"sha256:21988fa88192ca0c98c933b29ae13fdfbc21e55db5d5bbe1b7a745476dd40452"}
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
