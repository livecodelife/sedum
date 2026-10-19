package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample.go","path":"context"}}
import "context"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample.go","path":"time"}}
import "time"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample.go","path":"github.com/livecodelife/sedum/internal/growlog"}}
import "github.com/livecodelife/sedum/internal/growlog"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Config says how samples are drawn.\ntype Config struct {\n\t// Generators is the directory of generator packages copied into each\n\t// sample's project.\n\tGenerators string\n\t// Grow runs sedum grow in the prepared project dir. The real one calls the\n\t// herd's Haiku wrapper; tests supply a fake.\n\tGrow func(ctx context.Context, dir string) error\n\t// Keep leaves each sample's project on disk.\n\tKeep bool\n}","file":"evals/design/sample.go","name":"Config"}}
// Config says how samples are drawn.
type Config struct {
	// Generators is the directory of generator packages copied into each
	// sample's project.
	Generators string
	// Grow runs sedum grow in the prepared project dir. The real one calls the
	// herd's Haiku wrapper; tests supply a fake.
	Grow func(ctx context.Context, dir string) error
	// Keep leaves each sample's project on disk.
	Keep bool
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// logTotals sums the grow log a sample's grow wrote in dir.\nfunc logTotals(dir string) (calls, rejected, repairs, prompt, completion int) {\n\tentries, _, err := growlog.Read(filepath.Join(dir, \".sedum\", \"grow.jsonl\"))\n\tif err != nil {\n\t\treturn 0, 0, 0, 0, 0\n\t}\n\tfor _, e := range entries {\n\t\tcalls += e.Calls\n\t\trejected += e.Rejected\n\t\trepairs += e.Repairs\n\t\tprompt += e.PromptTokens\n\t\tcompletion += e.CompletionTokens\n\t}\n\treturn calls, rejected, repairs, prompt, completion\n}","file":"evals/design/sample.go","name":"logTotals"}}
// logTotals sums the grow log a sample's grow wrote in dir.
func logTotals(dir string) (calls, rejected, repairs, prompt, completion int) {
	entries, _, err := growlog.Read(filepath.Join(dir, ".sedum", "grow.jsonl"))
	if err != nil {
		return 0, 0, 0, 0, 0
	}
	for _, e := range entries {
		calls += e.Calls
		rejected += e.Rejected
		repairs += e.Repairs
		prompt += e.PromptTokens
		completion += e.CompletionTokens
	}
	return calls, rejected, repairs, prompt, completion
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// RunSample draws one sample: prepare a project, grow it, grade it. A failure\n// of the grow itself is recorded in the sample, not returned; an error is\n// returned only when the harness could not do its own work.\nfunc RunSample(ctx context.Context, cfg Config, t Task, v Variant, n int) (Sample, error) {\n\ts := Sample{Task: t.ID, Variant: v, N: n}\n\tdir, err := os.MkdirTemp(\"\", \"design-\"+t.ID+\"-\")\n\tif err != nil {\n\t\treturn s, err\n\t}\n\tif !cfg.Keep {\n\t\tdefer os.RemoveAll(dir)\n\t}\n\tif err := Prepare(t, v, cfg.Generators, dir); err != nil {\n\t\treturn s, err\n\t}\n\tstart := time.Now()\n\tgrowErr := cfg.Grow(ctx, dir)\n\ts.MS = time.Since(start).Milliseconds()\n\ts.Calls, s.Rejected, s.Repairs, s.PromptTokens, s.CompletionTokens = logTotals(dir)\n\tif growErr != nil {\n\t\ts.Error = firstLine(growErr.Error())\n\t\treturn s, nil\n\t}\n\tcompiled, passed, output, err := Grade(t, dir)\n\tif err != nil {\n\t\treturn s, err\n\t}\n\ts.Compiled, s.Passed = compiled, passed\n\tif !passed {\n\t\ts.Error = firstLine(output)\n\t}\n\treturn s, nil\n}","file":"evals/design/sample.go","name":"RunSample"}}
// RunSample draws one sample: prepare a project, grow it, grade it. A failure
// of the grow itself is recorded in the sample, not returned; an error is
// returned only when the harness could not do its own work.
func RunSample(ctx context.Context, cfg Config, t Task, v Variant, n int) (Sample, error) {
	s := Sample{Task: t.ID, Variant: v, N: n}
	dir, err := os.MkdirTemp("", "design-"+t.ID+"-")
	if err != nil {
		return s, err
	}
	if !cfg.Keep {
		defer os.RemoveAll(dir)
	}
	if err := Prepare(t, v, cfg.Generators, dir); err != nil {
		return s, err
	}
	start := time.Now()
	growErr := cfg.Grow(ctx, dir)
	s.MS = time.Since(start).Milliseconds()
	s.Calls, s.Rejected, s.Repairs, s.PromptTokens, s.CompletionTokens = logTotals(dir)
	if growErr != nil {
		s.Error = firstLine(growErr.Error())
		return s, nil
	}
	compiled, passed, output, err := Grade(t, dir)
	if err != nil {
		return s, err
	}
	s.Compiled, s.Passed = compiled, passed
	if !passed {
		s.Error = firstLine(output)
	}
	return s, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// firstLine returns the first non-empty line of s, at most 200 bytes.\nfunc firstLine(s string) string {\n\tfor _, l := range strings.Split(s, \"\\n\") {\n\t\tif l = strings.TrimSpace(l); l != \"\" {\n\t\t\tif len(l) > 200 {\n\t\t\t\tl = l[:200]\n\t\t\t}\n\t\t\treturn l\n\t\t}\n\t}\n\treturn \"\"\n}","file":"evals/design/sample.go","name":"firstLine"}}
// firstLine returns the first non-empty line of s, at most 200 bytes.
func firstLine(s string) string {
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			if len(l) > 200 {
				l = l[:200]
			}
			return l
		}
	}
	return ""
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Run draws k samples of every task and variant, in order, and returns them.\n// Every task is checked for isolation first.\nfunc Run(ctx context.Context, cfg Config, tasks []Task, variants []Variant, k int) ([]Sample, error) {\n\tfor _, t := range tasks {\n\t\tif err := t.Isolated(); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t}\n\tvar samples []Sample\n\tfor _, t := range tasks {\n\t\tfor _, v := range variants {\n\t\t\tfor n := 1; n <= k; n++ {\n\t\t\t\ts, err := RunSample(ctx, cfg, t, v, n)\n\t\t\t\tif err != nil {\n\t\t\t\t\treturn samples, fmt.Errorf(\"%s/%s/%d: %w\", t.ID, v, n, err)\n\t\t\t\t}\n\t\t\t\tsamples = append(samples, s)\n\t\t\t}\n\t\t}\n\t}\n\treturn samples, nil\n}","file":"evals/design/sample.go","name":"Run"}}
// Run draws k samples of every task and variant, in order, and returns them.
// Every task is checked for isolation first.
func Run(ctx context.Context, cfg Config, tasks []Task, variants []Variant, k int) ([]Sample, error) {
	for _, t := range tasks {
		if err := t.Isolated(); err != nil {
			return nil, err
		}
	}
	var samples []Sample
	for _, t := range tasks {
		for _, v := range variants {
			for n := 1; n <= k; n++ {
				s, err := RunSample(ctx, cfg, t, v, n)
				if err != nil {
					return samples, fmt.Errorf("%s/%s/%d: %w", t.ID, v, n, err)
				}
				samples = append(samples, s)
			}
		}
	}
	return samples, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Plan describes a run without drawing anything: the cells, the number of\n// grows, and a rough ceiling on tokens at maxPerGrow tokens a grow.\nfunc Plan(tasks []Task, variants []Variant, k, maxPerGrow int) string {\n\tgrows := len(tasks) * len(variants) * k\n\treturn fmt.Sprintf(\"%d task(s) x %d variant(s) x %d sample(s) = %d grow(s), at most ~%d tokens in total\\n\", len(tasks), len(variants), k, grows, grows*maxPerGrow)\n}","file":"evals/design/sample.go","name":"Plan"}}
// Plan describes a run without drawing anything: the cells, the number of
// grows, and a rough ceiling on tokens at maxPerGrow tokens a grow.
func Plan(tasks []Task, variants []Variant, k, maxPerGrow int) string {
	grows := len(tasks) * len(variants) * k
	return fmt.Sprintf("%d task(s) x %d variant(s) x %d sample(s) = %d grow(s), at most ~%d tokens in total\n", len(tasks), len(variants), k, grows, grows*maxPerGrow)
}

// /sedum:setFunction
