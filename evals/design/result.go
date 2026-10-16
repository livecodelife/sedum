package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result.go","path":"crypto/sha256"}}
import "crypto/sha256"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result.go","path":"encoding/hex"}}
import "encoding/hex"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/result.go","path":"text/tabwriter"}}
import "text/tabwriter"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Sample is one grow of one task with one intent variant.\ntype Sample struct {\n\tTask             string  `json:\"task\"`\n\tVariant          Variant `json:\"variant\"`\n\tN                int     `json:\"n\"`\n\tCompiled         bool    `json:\"compiled\"`\n\tPassed           bool    `json:\"passed\"`\n\tCalls            int     `json:\"calls\"`\n\tRejected         int     `json:\"rejected\"`\n\tRepairs          int     `json:\"repairs\"`\n\tPromptTokens     int     `json:\"prompt_tokens\"`\n\tCompletionTokens int     `json:\"completion_tokens\"`\n\tMS               int64   `json:\"ms\"`\n\tError            string  `json:\"error,omitempty\"`\n}","file":"evals/design/result.go","name":"Sample"}}
// Sample is one grow of one task with one intent variant.
type Sample struct {
	Task             string  `json:"task"`
	Variant          Variant `json:"variant"`
	N                int     `json:"n"`
	Compiled         bool    `json:"compiled"`
	Passed           bool    `json:"passed"`
	Calls            int     `json:"calls"`
	Rejected         int     `json:"rejected"`
	Repairs          int     `json:"repairs"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	MS               int64   `json:"ms"`
	Error            string  `json:"error,omitempty"`
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Cell is the totals of the samples of one task and variant.\ntype Cell struct {\n\tTask             string  `json:\"task\"`\n\tVariant          Variant `json:\"variant\"`\n\tSamples          int     `json:\"samples\"`\n\tCompiled         int     `json:\"compiled\"`\n\tPassed           int     `json:\"passed\"`\n\tCalls            int     `json:\"calls\"`\n\tRejected         int     `json:\"rejected\"`\n\tRepairs          int     `json:\"repairs\"`\n\tPromptTokens     int     `json:\"prompt_tokens\"`\n\tCompletionTokens int     `json:\"completion_tokens\"`\n}","file":"evals/design/result.go","name":"Cell"}}
// Cell is the totals of the samples of one task and variant.
type Cell struct {
	Task             string  `json:"task"`
	Variant          Variant `json:"variant"`
	Samples          int     `json:"samples"`
	Compiled         int     `json:"compiled"`
	Passed           int     `json:"passed"`
	Calls            int     `json:"calls"`
	Rejected         int     `json:"rejected"`
	Repairs          int     `json:"repairs"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
}

// /sedum:setType
// sedum:setMethod {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"func (c *Cell) add(s Sample) {\n\tc.Samples++\n\tif s.Compiled {\n\t\tc.Compiled++\n\t}\n\tif s.Passed {\n\t\tc.Passed++\n\t}\n\tc.Calls += s.Calls\n\tc.Rejected += s.Rejected\n\tc.Repairs += s.Repairs\n\tc.PromptTokens += s.PromptTokens\n\tc.CompletionTokens += s.CompletionTokens\n}","file":"evals/design/result.go","name":"add","receiver":"Cell"}}
func (c *Cell) add(s Sample) {
	c.Samples++
	if s.Compiled {
		c.Compiled++
	}
	if s.Passed {
		c.Passed++
	}
	c.Calls += s.Calls
	c.Rejected += s.Rejected
	c.Repairs += s.Repairs
	c.PromptTokens += s.PromptTokens
	c.CompletionTokens += s.CompletionTokens
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Tokens is the cell's prompt and completion tokens together.\nfunc (c Cell) Tokens() int { return c.PromptTokens + c.CompletionTokens }","file":"evals/design/result.go","name":"Tokens","receiver":"Cell"}}
// Tokens is the cell's prompt and completion tokens together.
func (c Cell) Tokens() int { return c.PromptTokens + c.CompletionTokens }

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Aggregate totals samples per task and variant, sorted by task and then by\n// variant order.\nfunc Aggregate(samples []Sample) []Cell {\n\tby := map[string]*Cell{}\n\tfor _, s := range samples {\n\t\tkey := s.Task + \"\\x00\" + string(s.Variant)\n\t\tc := by[key]\n\t\tif c == nil {\n\t\t\tc = &Cell{Task: s.Task, Variant: s.Variant}\n\t\t\tby[key] = c\n\t\t}\n\t\tc.add(s)\n\t}\n\torder := map[Variant]int{}\n\tfor i, v := range Variants() {\n\t\torder[v] = i\n\t}\n\tcells := make([]Cell, 0, len(by))\n\tfor _, c := range by {\n\t\tcells = append(cells, *c)\n\t}\n\tsort.Slice(cells, func(i, j int) bool {\n\t\tif cells[i].Task != cells[j].Task {\n\t\t\treturn cells[i].Task < cells[j].Task\n\t\t}\n\t\treturn order[cells[i].Variant] < order[cells[j].Variant]\n\t})\n\treturn cells\n}","file":"evals/design/result.go","name":"Aggregate"}}
// Aggregate totals samples per task and variant, sorted by task and then by
// variant order.
func Aggregate(samples []Sample) []Cell {
	by := map[string]*Cell{}
	for _, s := range samples {
		key := s.Task + "\x00" + string(s.Variant)
		c := by[key]
		if c == nil {
			c = &Cell{Task: s.Task, Variant: s.Variant}
			by[key] = c
		}
		c.add(s)
	}
	order := map[Variant]int{}
	for i, v := range Variants() {
		order[v] = i
	}
	cells := make([]Cell, 0, len(by))
	for _, c := range by {
		cells = append(cells, *c)
	}
	sort.Slice(cells, func(i, j int) bool {
		if cells[i].Task != cells[j].Task {
			return cells[i].Task < cells[j].Task
		}
		return order[cells[i].Variant] < order[cells[j].Variant]
	})
	return cells
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// TaskHash fingerprints a task: its task.yaml and every context and hidden\n// file, so a result says which version of the task it was drawn against.\nfunc TaskHash(t Task) (string, error) {\n\tsum := sha256.New()\n\tpaths := []string{filepath.Join(t.Dir, \"task.yaml\")}\n\tfor _, f := range t.Context {\n\t\tpaths = append(paths, filepath.Join(t.Dir, \"context\", f))\n\t}\n\tfor _, f := range t.Hidden {\n\t\tpaths = append(paths, filepath.Join(t.Dir, \"hidden\", f))\n\t}\n\tfor _, p := range paths {\n\t\tdata, err := os.ReadFile(p)\n\t\tif err != nil {\n\t\t\treturn \"\", err\n\t\t}\n\t\tfmt.Fprintf(sum, \"%s\\x00%d\\x00\", filepath.Base(p), len(data))\n\t\tsum.Write(data)\n\t}\n\treturn hex.EncodeToString(sum.Sum(nil))[:12], nil\n}","file":"evals/design/result.go","name":"TaskHash"}}
// TaskHash fingerprints a task: its task.yaml and every context and hidden
// file, so a result says which version of the task it was drawn against.
func TaskHash(t Task) (string, error) {
	sum := sha256.New()
	paths := []string{filepath.Join(t.Dir, "task.yaml")}
	for _, f := range t.Context {
		paths = append(paths, filepath.Join(t.Dir, "context", f))
	}
	for _, f := range t.Hidden {
		paths = append(paths, filepath.Join(t.Dir, "hidden", f))
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(sum, "%s\x00%d\x00", filepath.Base(p), len(data))
		sum.Write(data)
	}
	return hex.EncodeToString(sum.Sum(nil))[:12], nil
}

// /sedum:setFunction
// sedum:setType {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Report is a whole run: where it was drawn and every sample.\ntype Report struct {\n\tVersion    int               `json:\"version\"`\n\tCommit     string            `json:\"commit\"`\n\tSedum      string            `json:\"sedum\"`\n\tModel      string            `json:\"model\"`\n\tK          int               `json:\"k\"`\n\tTaskHashes map[string]string `json:\"task_hashes\"`\n\tSamples    []Sample          `json:\"samples\"`\n\tCells      []Cell            `json:\"cells\"`\n\tTotals     Cell              `json:\"totals\"`\n}","file":"evals/design/result.go","name":"Report"}}
// Report is a whole run: where it was drawn and every sample.
type Report struct {
	Version    int               `json:"version"`
	Commit     string            `json:"commit"`
	Sedum      string            `json:"sedum"`
	Model      string            `json:"model"`
	K          int               `json:"k"`
	TaskHashes map[string]string `json:"task_hashes"`
	Samples    []Sample          `json:"samples"`
	Cells      []Cell            `json:"cells"`
	Totals     Cell              `json:"totals"`
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// NewReport builds the report for samples, computing the cells and the totals.\nfunc NewReport(commit, sedum, model string, k int, hashes map[string]string, samples []Sample) Report {\n\tr := Report{Version: 1, Commit: commit, Sedum: sedum, Model: model, K: k, TaskHashes: hashes, Samples: samples, Cells: Aggregate(samples)}\n\tfor _, s := range samples {\n\t\tr.Totals.add(s)\n\t}\n\treturn r\n}","file":"evals/design/result.go","name":"NewReport"}}
// NewReport builds the report for samples, computing the cells and the totals.
func NewReport(commit, sedum, model string, k int, hashes map[string]string, samples []Sample) Report {
	r := Report{Version: 1, Commit: commit, Sedum: sedum, Model: model, K: k, TaskHashes: hashes, Samples: samples, Cells: Aggregate(samples)}
	for _, s := range samples {
		r.Totals.add(s)
	}
	return r
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// JSON renders the report as indented JSON with a trailing newline.\nfunc (r Report) JSON() ([]byte, error) {\n\tb, err := json.MarshalIndent(r, \"\", \"  \")\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn append(b, '\\n'), nil\n}","file":"evals/design/result.go","name":"JSON","receiver":"Report"}}
// JSON renders the report as indented JSON with a trailing newline.
func (r Report) JSON() ([]byte, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-8b91f1f9","gen":"sha256:12e888fde70a07acdfaa513cd82e9f69ed9081c8aab93d33b064dfa10f760e54","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Table renders the report: passes per task for each variant, then per-variant\n// totals with retries, repairs and tokens. A rate over a handful of samples is\n// a direction, not a number, and the header says so.\nfunc (r Report) Table() string {\n\tvar b strings.Builder\n\tfmt.Fprintf(&b, \"design eval: model %s, K=%d samples per cell (a rate over so few samples is a direction, not a number)\\n\\n\", r.Model, r.K)\n\ttasks := map[string]bool{}\n\tfor _, c := range r.Cells {\n\t\ttasks[c.Task] = true\n\t}\n\tnames := make([]string, 0, len(tasks))\n\tfor t := range tasks {\n\t\tnames = append(names, t)\n\t}\n\tsort.Strings(names)\n\tcell := func(task string, v Variant) Cell {\n\t\tfor _, c := range r.Cells {\n\t\t\tif c.Task == task && c.Variant == v {\n\t\t\t\treturn c\n\t\t\t}\n\t\t}\n\t\treturn Cell{}\n\t}\n\tvar buf bytes.Buffer\n\tw := tabwriter.NewWriter(&buf, 0, 8, 2, ' ', 0)\n\tfmt.Fprintf(w, \"passed\\t%s\\n\", strings.Join(names, \"\\t\"))\n\tfor _, v := range Variants() {\n\t\trow := []string{string(v)}\n\t\tfor _, n := range names {\n\t\t\tc := cell(n, v)\n\t\t\trow = append(row, fmt.Sprintf(\"%d/%d\", c.Passed, c.Samples))\n\t\t}\n\t\tfmt.Fprintln(w, strings.Join(row, \"\\t\"))\n\t}\n\tw.Flush()\n\tb.Write(buf.Bytes())\n\tb.WriteString(\"\\n\")\n\tbuf.Reset()\n\tw = tabwriter.NewWriter(&buf, 0, 8, 2, ' ', 0)\n\tfmt.Fprintln(w, \"variant\\tsamples\\tcompiled\\tpassed\\tretries\\trepairs\\ttokens\")\n\tfor _, v := range Variants() {\n\t\tvar t Cell\n\t\tfor _, c := range r.Cells {\n\t\t\tif c.Variant == v {\n\t\t\t\tt.Samples += c.Samples\n\t\t\t\tt.Compiled += c.Compiled\n\t\t\t\tt.Passed += c.Passed\n\t\t\t\tt.Rejected += c.Rejected\n\t\t\t\tt.Repairs += c.Repairs\n\t\t\t\tt.PromptTokens += c.PromptTokens\n\t\t\t\tt.CompletionTokens += c.CompletionTokens\n\t\t\t}\n\t\t}\n\t\tfmt.Fprintf(w, \"%s\\t%d\\t%d\\t%d\\t%d\\t%d\\t%d\\n\", v, t.Samples, t.Compiled, t.Passed, t.Rejected, t.Repairs, t.Tokens())\n\t}\n\tfmt.Fprintf(w, \"total\\t%d\\t%d\\t%d\\t%d\\t%d\\t%d\\n\", r.Totals.Samples, r.Totals.Compiled, r.Totals.Passed, r.Totals.Rejected, r.Totals.Repairs, r.Totals.Tokens())\n\tw.Flush()\n\tb.Write(buf.Bytes())\n\treturn b.String()\n}","file":"evals/design/result.go","name":"Table","receiver":"Report"}}
// Table renders the report: passes per task for each variant, then per-variant
// totals with retries, repairs and tokens. A rate over a handful of samples is
// a direction, not a number, and the header says so.
func (r Report) Table() string {
	var b strings.Builder
	fmt.Fprintf(&b, "design eval: model %s, K=%d samples per cell (a rate over so few samples is a direction, not a number)\n\n", r.Model, r.K)
	tasks := map[string]bool{}
	for _, c := range r.Cells {
		tasks[c.Task] = true
	}
	names := make([]string, 0, len(tasks))
	for t := range tasks {
		names = append(names, t)
	}
	sort.Strings(names)
	cell := func(task string, v Variant) Cell {
		for _, c := range r.Cells {
			if c.Task == task && c.Variant == v {
				return c
			}
		}
		return Cell{}
	}
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "passed\t%s\n", strings.Join(names, "\t"))
	for _, v := range Variants() {
		row := []string{string(v)}
		for _, n := range names {
			c := cell(n, v)
			row = append(row, fmt.Sprintf("%d/%d", c.Passed, c.Samples))
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
	b.Write(buf.Bytes())
	b.WriteString("\n")
	buf.Reset()
	w = tabwriter.NewWriter(&buf, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "variant\tsamples\tcompiled\tpassed\tretries\trepairs\ttokens")
	for _, v := range Variants() {
		var t Cell
		for _, c := range r.Cells {
			if c.Variant == v {
				t.Samples += c.Samples
				t.Compiled += c.Compiled
				t.Passed += c.Passed
				t.Rejected += c.Rejected
				t.Repairs += c.Repairs
				t.PromptTokens += c.PromptTokens
				t.CompletionTokens += c.CompletionTokens
			}
		}
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%d\t%d\n", v, t.Samples, t.Compiled, t.Passed, t.Rejected, t.Repairs, t.Tokens())
	}
	fmt.Fprintf(w, "total\t%d\t%d\t%d\t%d\t%d\t%d\n", r.Totals.Samples, r.Totals.Compiled, r.Totals.Passed, r.Totals.Rejected, r.Totals.Repairs, r.Totals.Tokens())
	w.Flush()
	b.Write(buf.Bytes())
	return b.String()
}

// /sedum:setMethod
