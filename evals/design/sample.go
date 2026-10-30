package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample.go","path":"context"},"kwargs":"sha256:6b44c0f772b3b14d753c0fd3b2eef7df24cd5648739db6dd1d63fd0ba68e0a9a","body":"sha256:0f7443fb25d7ad6a970617d4d0df6123283e24801b94c1d84fbe431230f31e1e"}
import "context"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample.go","path":"fmt"},"kwargs":"sha256:e76a378fc58fd4f8ea2b40626e58dd1128832b3560226032b1e9dcacedede2ba","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample.go","path":"os"},"kwargs":"sha256:3f4a5772f45600b0c4ef86e83b8addec861f824af0a9cf95874711a637f4ea58","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample.go","path":"path/filepath"},"kwargs":"sha256:5dd8488a281433b4a7eeb81a4e17d482802f9367508b823276afc9bbd8f1245f","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample.go","path":"strings"},"kwargs":"sha256:cde1a60d15f7d66d5d29cef24446a6d2c41e3782cfa7b176c5a6017ba26404be","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample.go","path":"time"},"kwargs":"sha256:64312aacf136c9b6de9c9dcef889e16a44fd0c75f7d4492595bd13e41854659a","body":"sha256:e99128ebf4e36998f40b33adf2eef57f4d6c315f1a4d1bdd3b7fc09a0d7bc82d"}
import "time"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample.go","path":"github.com/livecodelife/sedum/internal/growlog"},"kwargs":"sha256:ad1b8e17381ac3184e80448d467600231b7af36e1e239f786372873c3366d7e3","body":"sha256:b7968f24df7d76473f917ee9527158f424f930250a91e5407a0140aa4c56c91f"}
import "github.com/livecodelife/sedum/internal/growlog"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"evals/design/sample.go","name":"Config"},"kwargs":"sha256:a276611c9a6cc6e3134fe4395a28fd8a34ea873b69106a7948192208c758cd9e","body":"sha256:549ebf326fae64fb33320b134a3475f68b8023a75f0a95bbf0abc6ea302211f5"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/sample.go","name":"logTotals"},"kwargs":"sha256:bb3bd48349c24b2d2fc39d8c82bc3c339c60022b79f61dc70bb0cb403d4e062a","body":"sha256:31fd91030c3188f1c75e95d2c1dbb62db597f53f0f75a39867013ce5e6e768ad"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/sample.go","name":"RunSample"},"kwargs":"sha256:8b9705b64d616a9df9ce616681030ebb28ad0806ac9bc31798f6195eeeddf882","body":"sha256:56086f2c0b0d43072ff4bf37abaae239d568a793c908214bf4fac31a693856fb"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/sample.go","name":"firstLine"},"kwargs":"sha256:42b8190f52380bbe23925a387ac53db0bfe971b2a064c9fd125bacb33b715dff","body":"sha256:8b6c79d10f4515d2611f3d3d677ef49defaab961bebdc2029b0cfab2bc1054ea"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/sample.go","name":"Run"},"kwargs":"sha256:70520cec922a4de2fa3143aa5a777f55b98d29ebc2150065c3c7f071fe8a7acd","body":"sha256:1217053f1eecd99e632212c361f7ff6dd1474c2e42b2ee4a85ca5f8284e1fedd"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/sample.go","name":"Plan"},"kwargs":"sha256:f5beafa5ce59cf9d596061b2c6b13288ddd7b6cd553b84d0530c9a0d98711884","body":"sha256:96a60e345aaf00552ba393aa22801f1639a5e43ecad19c4dcf74f67b3cd1b23f"}
// Plan describes a run without drawing anything: the cells, the number of
// grows, and a rough ceiling on tokens at maxPerGrow tokens a grow.
func Plan(tasks []Task, variants []Variant, k, maxPerGrow int) string {
	grows := len(tasks) * len(variants) * k
	return fmt.Sprintf("%d task(s) x %d variant(s) x %d sample(s) = %d grow(s), at most ~%d tokens in total\n", len(tasks), len(variants), k, grows, grows*maxPerGrow)
}

// /sedum:setFunction
