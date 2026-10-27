package gitcheck

import (
	"bytes"
	"encoding/json"
)

// Log is the runs of one SARIF document as provenance printed them.
type Log struct {
	runs []map[string]json.RawMessage
}

// Result is one SARIF result as a JSON object, with its run index and its
// result index in the log it came from.
type Result struct {
	obj      map[string]json.RawMessage
	run, idx int
}

// ParseLog reads the bytes provenance printed as a SARIF document. Bytes that
// are not JSON are a provenance error. A document with no runs, or a result
// without a string ruleId or a message text, is skew.
func ParseLog(data []byte) (Log, error) {
	var doc struct {
		Runs []map[string]json.RawMessage `json:"runs"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return Log{}, ProvenanceError("output is not SARIF: " + err.Error())
	}
	if len(doc.Runs) == 0 {
		return Log{}, SkewError("SARIF has no runs")
	}
	log := Log{runs: doc.Runs}
	for _, r := range log.Results() {
		var rule string
		if err := json.Unmarshal(r.obj["ruleId"], &rule); err != nil {
			return Log{}, SkewError("a SARIF result has no string ruleId")
		}
		var msg struct {
			Text *string `json:"text"`
		}
		if err := json.Unmarshal(r.obj["message"], &msg); err != nil || msg.Text == nil {
			return Log{}, SkewError("a SARIF result has no message text")
		}
	}
	return log, nil
}

// Results returns the results of every run in order.
func (l Log) Results() []Result {
	var out []Result
	for run, r := range l.runs {
		var results []map[string]json.RawMessage
		_ = json.Unmarshal(r["results"], &results)
		for i, obj := range results {
			out = append(out, Result{obj: obj, run: run, idx: i})
		}
	}
	return out
}

// RunIndex is the index of the run the result came from.
func (r Result) RunIndex() int { return r.run }

// ResultIndex is the index of the result in its run.
func (r Result) ResultIndex() int { return r.idx }

// RuleID is the result's ruleId.
func (r Result) RuleID() string {
	var s string
	_ = json.Unmarshal(r.obj["ruleId"], &s)
	return s
}

// MessageText is the text of the result's message.
func (r Result) MessageText() string {
	var m struct {
		Text string `json:"text"`
	}
	_ = json.Unmarshal(r.obj["message"], &m)
	return m.Text
}

// Property is the string value of the named result property, or empty.
func (r Result) Property(name string) string {
	var props map[string]json.RawMessage
	if json.Unmarshal(r.obj["properties"], &props) != nil {
		return ""
	}
	var s string
	_ = json.Unmarshal(props[name], &s)
	return s
}

// FirstLocationURI is the artifact uri of the result's first location, or empty
// when it has none.
func (r Result) FirstLocationURI() string {
	var locs []struct {
		PhysicalLocation struct {
			ArtifactLocation struct {
				URI string `json:"uri"`
			} `json:"artifactLocation"`
		} `json:"physicalLocation"`
	}
	if json.Unmarshal(r.obj["locations"], &locs) != nil || len(locs) == 0 {
		return ""
	}
	return locs[0].PhysicalLocation.ArtifactLocation.URI
}

// Merged is the one SARIF run built from several provenance logs.
type Merged struct {
	rules     []map[string]json.RawMessage
	results   []map[string]json.RawMessage
	artifacts []map[string]json.RawMessage
}

// Merge combines the logs into one run. justifications holds, for each log, one
// string per result in Results order; a result with a non-empty justification
// gets one external suppression carrying it. Results are copied unchanged apart
// from ruleIndex, which is recomputed against the union of the logs' rules, and
// the suppression.
func Merge(logs []Log, justifications [][]string) Merged {
	var m Merged
	ruleAt := map[string]int{}
	artifactAt := map[string]int{}
	addRule := func(id string, rule map[string]json.RawMessage) int {
		if at, ok := ruleAt[id]; ok {
			return at
		}
		if rule == nil {
			idJSON, _ := json.Marshal(id)
			rule = map[string]json.RawMessage{"id": idJSON}
		}
		ruleAt[id] = len(m.rules)
		m.rules = append(m.rules, rule)
		return ruleAt[id]
	}
	for li, log := range logs {
		for _, run := range log.runs {
			var rules []map[string]json.RawMessage
			var tool struct {
				Driver struct {
					Rules []map[string]json.RawMessage `json:"rules"`
				} `json:"driver"`
			}
			if json.Unmarshal(run["tool"], &tool) == nil {
				rules = tool.Driver.Rules
			}
			for _, rule := range rules {
				var id string
				if json.Unmarshal(rule["id"], &id) == nil {
					addRule(id, rule)
				}
			}
			var artifacts []map[string]json.RawMessage
			_ = json.Unmarshal(run["artifacts"], &artifacts)
			for _, a := range artifacts {
				var loc struct {
					URI string `json:"uri"`
				}
				_ = json.Unmarshal(a["location"], &loc)
				if _, ok := artifactAt[loc.URI]; ok {
					continue
				}
				artifactAt[loc.URI] = len(m.artifacts)
				m.artifacts = append(m.artifacts, a)
			}
		}
		for i, r := range log.Results() {
			obj := make(map[string]json.RawMessage, len(r.obj)+1)
			for k, v := range r.obj {
				obj[k] = v
			}
			at := addRule(r.RuleID(), nil)
			obj["ruleIndex"], _ = json.Marshal(at)
			if li < len(justifications) && i < len(justifications[li]) && justifications[li][i] != "" {
				obj["suppressions"], _ = json.Marshal([]map[string]string{{"kind": "external", "justification": justifications[li][i]}})
			}
			m.results = append(m.results, obj)
		}
	}
	// Stable insertion sort by uri: the artifact lists are short.
	for i := 1; i < len(m.artifacts); i++ {
		for j := i; j > 0 && artifactURI(m.artifacts[j]) < artifactURI(m.artifacts[j-1]); j-- {
			m.artifacts[j], m.artifacts[j-1] = m.artifacts[j-1], m.artifacts[j]
		}
	}
	return m
}

func artifactURI(a map[string]json.RawMessage) string {
	var loc struct {
		URI string `json:"uri"`
	}
	_ = json.Unmarshal(a["location"], &loc)
	return loc.URI
}

// JSON serialises the merged SARIF document, with the count of skipped scope
// paths as the run property skippedScopePaths.
func (m Merged) JSON(skipped int) ([]byte, error) {
	nonNil := func(s []map[string]json.RawMessage) []map[string]json.RawMessage {
		if s == nil {
			return []map[string]json.RawMessage{}
		}
		return s
	}
	doc := map[string]any{
		"$schema": "https://json.schemastore.org/sarif-2.1.0.json",
		"version": "2.1.0",
		"runs": []map[string]any{{
			"tool": map[string]any{"driver": map[string]any{
				"name":  "sedum check",
				"rules": nonNil(m.rules),
			}},
			"results":    nonNil(m.results),
			"artifacts":  nonNil(m.artifacts),
			"properties": map[string]int{"skippedScopePaths": skipped},
		}},
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(doc); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}
