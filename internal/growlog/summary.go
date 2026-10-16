package growlog

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/summary.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/summary.go","path":"time"}}
import "time"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Filter selects which entries Summarize counts. Zero fields match everything.\ntype Filter struct {\n\tSince  time.Time\n\tRecord string\n\tModel  string\n}","file":"internal/growlog/summary.go","name":"Filter"}}
// Filter selects which entries Summarize counts. Zero fields match everything.
type Filter struct {
	Since  time.Time
	Record string
	Model  string
}

// /sedum:setType
// sedum:setMethod {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// match reports whether e passes the filter. An entry whose ts does not parse\n// fails a non-zero Since.\nfunc (f Filter) match(e Entry) bool {\n\tif f.Record != \"\" && e.Record != f.Record {\n\t\treturn false\n\t}\n\tif f.Model != \"\" && e.Model != f.Model {\n\t\treturn false\n\t}\n\tif !f.Since.IsZero() {\n\t\tts, err := time.Parse(time.RFC3339, e.TS)\n\t\tif err != nil || ts.Before(f.Since) {\n\t\t\treturn false\n\t\t}\n\t}\n\treturn true\n}","file":"internal/growlog/summary.go","name":"match","receiver":"Filter"}}
// match reports whether e passes the filter. An entry whose ts does not parse
// fails a non-zero Since.
func (f Filter) match(e Entry) bool {
	if f.Record != "" && e.Record != f.Record {
		return false
	}
	if f.Model != "" && e.Model != f.Model {
		return false
	}
	if !f.Since.IsZero() {
		ts, err := time.Parse(time.RFC3339, e.TS)
		if err != nil || ts.Before(f.Since) {
			return false
		}
	}
	return true
}

// /sedum:setMethod
// sedum:setType {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// ModelStats is the totals for one model. Retries are the answers the model\n// had refused (Entry.Rejected); Repairs are verification send-backs.\ntype ModelStats struct {\n\tModel            string `json:\"model\"`\n\tGrows            int    `json:\"grows\"`\n\tFailed           int    `json:\"failed\"`\n\tCalls            int    `json:\"calls\"`\n\tRetries          int    `json:\"retries\"`\n\tRepairs          int    `json:\"repairs\"`\n\tCompleteness     int    `json:\"completeness\"`\n\tPromptTokens     int    `json:\"prompt_tokens\"`\n\tCompletionTokens int    `json:\"completion_tokens\"`\n\tInvocations      int    `json:\"invocations\"`\n\tLiteral          int    `json:\"literal\"`\n}","file":"internal/growlog/summary.go","name":"ModelStats"}}
// ModelStats is the totals for one model. Retries are the answers the model
// had refused (Entry.Rejected); Repairs are verification send-backs.
type ModelStats struct {
	Model            string `json:"model"`
	Grows            int    `json:"grows"`
	Failed           int    `json:"failed"`
	Calls            int    `json:"calls"`
	Retries          int    `json:"retries"`
	Repairs          int    `json:"repairs"`
	Completeness     int    `json:"completeness"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	Invocations      int    `json:"invocations"`
	Literal          int    `json:"literal"`
}

// /sedum:setType
// sedum:setMethod {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// LiteralShare is the fraction of invocations that carried model-written code,\n// or 0 when there were none.\nfunc (m ModelStats) LiteralShare() float64 {\n\tif m.Invocations == 0 {\n\t\treturn 0\n\t}\n\treturn float64(m.Literal) / float64(m.Invocations)\n}","file":"internal/growlog/summary.go","name":"LiteralShare","receiver":"ModelStats"}}
// LiteralShare is the fraction of invocations that carried model-written code,
// or 0 when there were none.
func (m ModelStats) LiteralShare() float64 {
	if m.Invocations == 0 {
		return 0
	}
	return float64(m.Literal) / float64(m.Invocations)
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// FailureRate is the fraction of grows that failed, or 0 when there were none.\nfunc (m ModelStats) FailureRate() float64 {\n\tif m.Grows == 0 {\n\t\treturn 0\n\t}\n\treturn float64(m.Failed) / float64(m.Grows)\n}","file":"internal/growlog/summary.go","name":"FailureRate","receiver":"ModelStats"}}
// FailureRate is the fraction of grows that failed, or 0 when there were none.
func (m ModelStats) FailureRate() float64 {
	if m.Grows == 0 {
		return 0
	}
	return float64(m.Failed) / float64(m.Grows)
}

// /sedum:setMethod
// sedum:setType {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Summary is Summarize's result: per-model stats sorted by model name, and the\n// total over all of them.\ntype Summary struct {\n\tEntries   int          `json:\"entries\"`\n\tMalformed int          `json:\"malformed\"`\n\tModels    []ModelStats `json:\"models\"`\n\tTotal     ModelStats   `json:\"total\"`\n}","file":"internal/growlog/summary.go","name":"Summary"}}
// Summary is Summarize's result: per-model stats sorted by model name, and the
// total over all of them.
type Summary struct {
	Entries   int          `json:"entries"`
	Malformed int          `json:"malformed"`
	Models    []ModelStats `json:"models"`
	Total     ModelStats   `json:"total"`
}

// /sedum:setType
// sedum:setMethod {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// add counts e into m. A skip line is not a grow.\nfunc (m *ModelStats) add(e Entry) {\n\tif e.Mode != \"skip\" {\n\t\tm.Grows++\n\t\tif e.Outcome == \"failed\" {\n\t\t\tm.Failed++\n\t\t}\n\t}\n\tm.Calls += e.Calls\n\tm.Retries += e.Rejected\n\tm.Repairs += e.Repairs\n\tm.Completeness += e.Completeness\n\tm.PromptTokens += e.PromptTokens\n\tm.CompletionTokens += e.CompletionTokens\n\tm.Invocations += e.Invocations\n\tm.Literal += e.Literal\n}","file":"internal/growlog/summary.go","name":"add","receiver":"ModelStats"}}
// add counts e into m. A skip line is not a grow.
func (m *ModelStats) add(e Entry) {
	if e.Mode != "skip" {
		m.Grows++
		if e.Outcome == "failed" {
			m.Failed++
		}
	}
	m.Calls += e.Calls
	m.Retries += e.Rejected
	m.Repairs += e.Repairs
	m.Completeness += e.Completeness
	m.PromptTokens += e.PromptTokens
	m.CompletionTokens += e.CompletionTokens
	m.Invocations += e.Invocations
	m.Literal += e.Literal
}

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Summarize totals the entries that pass f, per model and overall.\nfunc Summarize(entries []Entry, f Filter) Summary {\n\tby := map[string]*ModelStats{}\n\tvar s Summary\n\tfor _, e := range entries {\n\t\tif !f.match(e) {\n\t\t\tcontinue\n\t\t}\n\t\ts.Entries++\n\t\tm := by[e.Model]\n\t\tif m == nil {\n\t\t\tm = &ModelStats{Model: e.Model}\n\t\t\tby[e.Model] = m\n\t\t}\n\t\tm.add(e)\n\t\ts.Total.add(e)\n\t}\n\tfor _, m := range by {\n\t\ts.Models = append(s.Models, *m)\n\t}\n\tsort.Slice(s.Models, func(i, j int) bool { return s.Models[i].Model < s.Models[j].Model })\n\tif s.Models == nil {\n\t\ts.Models = []ModelStats{}\n\t}\n\treturn s\n}","file":"internal/growlog/summary.go","name":"Summarize"}}
// Summarize totals the entries that pass f, per model and overall.
func Summarize(entries []Entry, f Filter) Summary {
	by := map[string]*ModelStats{}
	var s Summary
	for _, e := range entries {
		if !f.match(e) {
			continue
		}
		s.Entries++
		m := by[e.Model]
		if m == nil {
			m = &ModelStats{Model: e.Model}
			by[e.Model] = m
		}
		m.add(e)
		s.Total.add(e)
	}
	for _, m := range by {
		s.Models = append(s.Models, *m)
	}
	sort.Slice(s.Models, func(i, j int) bool { return s.Models[i].Model < s.Models[j].Model })
	if s.Models == nil {
		s.Models = []ModelStats{}
	}
	return s
}

// /sedum:setFunction
