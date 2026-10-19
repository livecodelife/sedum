package growlog

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/summary_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/summary_test.go","path":"time"}}
import "time"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func summaryEntries() []Entry {\n\treturn []Entry{\n\t\t{TS: \"2026-10-01T00:00:00Z\", Record: \"r1\", Mode: \"select\", Model: \"haiku\", Outcome: \"ok\", Calls: 2, Rejected: 1, Repairs: 1, PromptTokens: 100, CompletionTokens: 20, Invocations: 4, Literal: 3},\n\t\t{TS: \"2026-10-02T00:00:00Z\", Record: \"r2\", Mode: \"select\", Model: \"haiku\", Outcome: \"failed\", Calls: 3, Rejected: 3, PromptTokens: 50, CompletionTokens: 5, Invocations: 0},\n\t\t{TS: \"2026-10-03T00:00:00Z\", Record: \"r1\", Mode: \"select\", Model: \"stub\", Outcome: \"ok\", Calls: 1, Invocations: 2},\n\t\t{TS: \"2026-10-04T00:00:00Z\", Record: \"r3\", Mode: \"skip\", Model: \"\", Outcome: \"ok\"},\n\t}\n}","file":"internal/growlog/summary_test.go","name":"summaryEntries"}}
func summaryEntries() []Entry {
	return []Entry{
		{TS: "2026-10-01T00:00:00Z", Record: "r1", Mode: "select", Model: "haiku", Outcome: "ok", Calls: 2, Rejected: 1, Repairs: 1, PromptTokens: 100, CompletionTokens: 20, Invocations: 4, Literal: 3},
		{TS: "2026-10-02T00:00:00Z", Record: "r2", Mode: "select", Model: "haiku", Outcome: "failed", Calls: 3, Rejected: 3, PromptTokens: 50, CompletionTokens: 5, Invocations: 0},
		{TS: "2026-10-03T00:00:00Z", Record: "r1", Mode: "select", Model: "stub", Outcome: "ok", Calls: 1, Invocations: 2},
		{TS: "2026-10-04T00:00:00Z", Record: "r3", Mode: "skip", Model: "", Outcome: "ok"},
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSummarizePerModel(t *testing.T) {\n\ts := Summarize(summaryEntries(), Filter{})\n\tif s.Entries != 4 || len(s.Models) != 3 {\n\t\tt.Fatalf(\"entries = %d, models = %d, want 4 and 3\", s.Entries, len(s.Models))\n\t}\n\tif s.Models[0].Model != \"\" || s.Models[1].Model != \"haiku\" || s.Models[2].Model != \"stub\" {\n\t\tt.Fatalf(\"models not sorted by name: %+v\", s.Models)\n\t}\n\th := s.Models[1]\n\tif h.Grows != 2 || h.Failed != 1 || h.Calls != 5 || h.Retries != 4 || h.Repairs != 1 || h.PromptTokens != 150 || h.CompletionTokens != 25 || h.Invocations != 4 || h.Literal != 3 {\n\t\tt.Fatalf(\"haiku = %+v\", h)\n\t}\n\tif h.FailureRate() != 0.5 || h.LiteralShare() != 0.75 {\n\t\tt.Fatalf(\"rates = %v, %v\", h.FailureRate(), h.LiteralShare())\n\t}\n\tif s.Models[0].Grows != 0 {\n\t\tt.Fatalf(\"a skip line is not a grow: %+v\", s.Models[0])\n\t}\n\tif s.Total.Grows != 3 || s.Total.Calls != 6 || s.Total.Retries != 4 {\n\t\tt.Fatalf(\"total = %+v\", s.Total)\n\t}\n}","file":"internal/growlog/summary_test.go","name":"TestSummarizePerModel"}}
func TestSummarizePerModel(t *testing.T) {
	s := Summarize(summaryEntries(), Filter{})
	if s.Entries != 4 || len(s.Models) != 3 {
		t.Fatalf("entries = %d, models = %d, want 4 and 3", s.Entries, len(s.Models))
	}
	if s.Models[0].Model != "" || s.Models[1].Model != "haiku" || s.Models[2].Model != "stub" {
		t.Fatalf("models not sorted by name: %+v", s.Models)
	}
	h := s.Models[1]
	if h.Grows != 2 || h.Failed != 1 || h.Calls != 5 || h.Retries != 4 || h.Repairs != 1 || h.PromptTokens != 150 || h.CompletionTokens != 25 || h.Invocations != 4 || h.Literal != 3 {
		t.Fatalf("haiku = %+v", h)
	}
	if h.FailureRate() != 0.5 || h.LiteralShare() != 0.75 {
		t.Fatalf("rates = %v, %v", h.FailureRate(), h.LiteralShare())
	}
	if s.Models[0].Grows != 0 {
		t.Fatalf("a skip line is not a grow: %+v", s.Models[0])
	}
	if s.Total.Grows != 3 || s.Total.Calls != 6 || s.Total.Retries != 4 {
		t.Fatalf("total = %+v", s.Total)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSummarizeFilters(t *testing.T) {\n\tsince := time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)\n\tif s := Summarize(summaryEntries(), Filter{Since: since}); s.Entries != 2 {\n\t\tt.Fatalf(\"since: entries = %d, want 2\", s.Entries)\n\t}\n\tif s := Summarize(summaryEntries(), Filter{Model: \"haiku\"}); s.Entries != 2 || s.Total.Calls != 5 {\n\t\tt.Fatalf(\"model: %+v\", s)\n\t}\n\tif s := Summarize(summaryEntries(), Filter{Record: \"r1\"}); s.Entries != 2 {\n\t\tt.Fatalf(\"record: entries = %d, want 2\", s.Entries)\n\t}\n\tbad := []Entry{{TS: \"yesterday\", Model: \"m\", Outcome: \"ok\"}}\n\tif s := Summarize(bad, Filter{Since: since}); s.Entries != 0 {\n\t\tt.Fatalf(\"an unparsable ts must fail Since: %+v\", s)\n\t}\n}","file":"internal/growlog/summary_test.go","name":"TestSummarizeFilters"}}
func TestSummarizeFilters(t *testing.T) {
	since := time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)
	if s := Summarize(summaryEntries(), Filter{Since: since}); s.Entries != 2 {
		t.Fatalf("since: entries = %d, want 2", s.Entries)
	}
	if s := Summarize(summaryEntries(), Filter{Model: "haiku"}); s.Entries != 2 || s.Total.Calls != 5 {
		t.Fatalf("model: %+v", s)
	}
	if s := Summarize(summaryEntries(), Filter{Record: "r1"}); s.Entries != 2 {
		t.Fatalf("record: entries = %d, want 2", s.Entries)
	}
	bad := []Entry{{TS: "yesterday", Model: "m", Outcome: "ok"}}
	if s := Summarize(bad, Filter{Since: since}); s.Entries != 0 {
		t.Fatalf("an unparsable ts must fail Since: %+v", s)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSummarizeEmpty(t *testing.T) {\n\ts := Summarize(nil, Filter{})\n\tif s.Models == nil || len(s.Models) != 0 || s.Entries != 0 {\n\t\tt.Fatalf(\"empty summary = %+v\", s)\n\t}\n}","file":"internal/growlog/summary_test.go","name":"TestSummarizeEmpty"}}
func TestSummarizeEmpty(t *testing.T) {
	s := Summarize(nil, Filter{})
	if s.Models == nil || len(s.Models) != 0 || s.Entries != 0 {
		t.Fatalf("empty summary = %+v", s)
	}
}

// /sedum:setFunction
