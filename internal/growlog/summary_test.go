package growlog

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/summary_test.go","path":"testing"},"kwargs":"sha256:92424b18e0042fa081727c39d9d100501f00279efc309089b231db7513582f92","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/summary_test.go","path":"time"},"kwargs":"sha256:ce9252af7e9ef628a642eda046c9ebd4406830ce59bf0e4dcd1eb0d0a79c2f04","body":"sha256:e99128ebf4e36998f40b33adf2eef57f4d6c315f1a4d1bdd3b7fc09a0d7bc82d"}
import "time"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/summary_test.go","name":"summaryEntries"},"kwargs":"sha256:e6d50f188033f2b9016e5b6d13d90e1141dea7f4f221191989af0a5f857727f7","body":"sha256:ccfbf1a7a07ccc6df8e25e960407a90e12b06f1fe4e4cb3a154d6c92657f74cf"}
func summaryEntries() []Entry {
	return []Entry{
		{TS: "2026-10-01T00:00:00Z", Record: "r1", Mode: "select", Model: "haiku", Outcome: "ok", Calls: 2, Rejected: 1, Repairs: 1, PromptTokens: 100, CompletionTokens: 20, Invocations: 4, Literal: 3},
		{TS: "2026-10-02T00:00:00Z", Record: "r2", Mode: "select", Model: "haiku", Outcome: "failed", Calls: 3, Rejected: 3, PromptTokens: 50, CompletionTokens: 5, Invocations: 0},
		{TS: "2026-10-03T00:00:00Z", Record: "r1", Mode: "select", Model: "stub", Outcome: "ok", Calls: 1, Invocations: 2},
		{TS: "2026-10-04T00:00:00Z", Record: "r3", Mode: "skip", Model: "", Outcome: "ok"},
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/summary_test.go","name":"TestSummarizePerModel"},"kwargs":"sha256:67bfcb2539818724ff603353fe64c6ca9ede0530688c4de1d5edffc5b37327ed","body":"sha256:cee0bcc367bf374c8f0468242762cdd9d6c7d3d11667344b6ec320564b6ad51f"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/summary_test.go","name":"TestSummarizeFilters"},"kwargs":"sha256:fbd948455429e888e3d9ab7b73822c794b8877d67b55c569bce9348cb0de8138","body":"sha256:55de38b2f9f431abfdb05b4833b73a57bf7f51abb7b77035ac2cd29f47c87af1"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/summary_test.go","name":"TestSummarizeEmpty"},"kwargs":"sha256:b04e3dad79e16f27c55257d76b5bcb3810bfe3f1e188b24cce609974c3a3d3f5","body":"sha256:0fb021c8ec1a5acfb57aaa1c8e25b3b73c76c1b7f77ab2872bc548f26f4f0e09"}
func TestSummarizeEmpty(t *testing.T) {
	s := Summarize(nil, Filter{})
	if s.Models == nil || len(s.Models) != 0 || s.Entries != 0 {
		t.Fatalf("empty summary = %+v", s)
	}
}

// /sedum:setFunction
