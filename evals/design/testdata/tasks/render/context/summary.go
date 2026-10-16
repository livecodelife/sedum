package growlog

import (
	"sort"
	"time"
)

// Filter selects which entries Summarize counts. Zero fields match everything.
type Filter struct {
	Since  time.Time
	Record string
	Model  string
}

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

// LiteralShare is the fraction of invocations that carried model-written code,
// or 0 when there were none.
func (m ModelStats) LiteralShare() float64 {
	if m.Invocations == 0 {
		return 0
	}
	return float64(m.Literal) / float64(m.Invocations)
}

// FailureRate is the fraction of grows that failed, or 0 when there were none.
func (m ModelStats) FailureRate() float64 {
	if m.Grows == 0 {
		return 0
	}
	return float64(m.Failed) / float64(m.Grows)
}

// Summary is Summarize's result: per-model stats sorted by model name, and the
// total over all of them.
type Summary struct {
	Entries   int          `json:"entries"`
	Malformed int          `json:"malformed"`
	Models    []ModelStats `json:"models"`
	Total     ModelStats   `json:"total"`
}

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
