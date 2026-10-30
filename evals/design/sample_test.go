package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample_test.go","path":"context"},"kwargs":"sha256:290605a52bc1b57689fdb6fabc12d543d42f06ac03b1eed7ffd5893d14fcbfbb","body":"sha256:0f7443fb25d7ad6a970617d4d0df6123283e24801b94c1d84fbe431230f31e1e"}
import "context"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample_test.go","path":"errors"},"kwargs":"sha256:ef035e37947fff9b6b0508911fc8c3909a577617a4b7144f7ca9817ea81764d3","body":"sha256:2b718237ba0a56467c8fd494cc131e942ef410e6228c03b644d70476c71862dc"}
import "errors"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample_test.go","path":"os"},"kwargs":"sha256:afe985957d45cf832264d85bb77ccbb042d10f4421ccf76a80133c490a55d680","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample_test.go","path":"path/filepath"},"kwargs":"sha256:a64dc28a88589cac4016cd1612f471ff8e053e72bbcf57162c4abc094dc289bc","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample_test.go","path":"strings"},"kwargs":"sha256:7d679d3e9ba888bf8b0cb48b8a91dd3d2f3bae1e53e7d0138dbab58619ea21a2","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample_test.go","path":"testing"},"kwargs":"sha256:65e4a00eae38fb9c1342c528018cf6dfea0f9e4e52b16b3468e89bd639bb5566","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/sample_test.go","path":"github.com/livecodelife/sedum/internal/growlog"},"kwargs":"sha256:93ea44b485496bc02b23a3a673cafd152149dfaf86ef8236768e05e270cc4a1c","body":"sha256:b7968f24df7d76473f917ee9527158f424f930250a91e5407a0140aa4c56c91f"}
import "github.com/livecodelife/sedum/internal/growlog"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/sample_test.go","name":"referenceGrow"},"kwargs":"sha256:6e8f4d32dfb438dc6734df6480896764ea4da9965fd2df2936b48c60da3b20b2","body":"sha256:2f00e6d9675ecaff8f315916baeab397bb14a9de600b321d89014a78abdcf171"}
// referenceGrow writes the task's reference code, taken from its code intent,
// as the target file, and logs one grow the way sedum grow does.
func referenceGrow(t *testing.T, task Task) func(context.Context, string) error {
	return func(_ context.Context, dir string) error {
		parts := strings.SplitN(task.Code, ":\n\n", 2)
		if len(parts) != 2 {
			t.Fatalf("task %s: unexpected code intent", task.ID)
		}
		if err := os.WriteFile(filepath.Join(dir, task.Package, task.Target[0]), []byte(parts[1]), 0o644); err != nil {
			return err
		}
		return growlog.Append(filepath.Join(dir, ".sedum", "grow.jsonl"), growlog.Entry{Calls: 2, Rejected: 1, Repairs: 1, PromptTokens: 100, CompletionTokens: 10})
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/sample_test.go","name":"TestReferenceCodePassesEveryTask"},"kwargs":"sha256:a5268b8e6c0928a25a88770099fad20fafa95fe7d6352802869a82ddc65bf613","body":"sha256:612a287ac9384f35f4e000127881475170f754fe97cb9404c51e45d54dc1ebc9"}
func TestReferenceCodePassesEveryTask(t *testing.T) {
	tasks, err := LoadTasks("testdata/tasks")
	if err != nil {
		t.Fatal(err)
	}
	for _, task := range tasks {
		t.Run(task.ID, func(t *testing.T) {
			cfg := Config{Generators: "../../generators", Grow: referenceGrow(t, task)}
			s, err := RunSample(context.Background(), cfg, task, VariantCode, 1)
			if err != nil {
				t.Fatal(err)
			}
			if !s.Compiled || !s.Passed {
				t.Fatalf("the reference code must pass its own hidden tests: %+v", s)
			}
			if s.Calls != 2 || s.Rejected != 1 || s.Repairs != 1 || s.PromptTokens != 100 || s.CompletionTokens != 10 {
				t.Fatalf("the grow log was not read: %+v", s)
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/sample_test.go","name":"TestRunSampleOutcomes"},"kwargs":"sha256:63d6a8e598f90cfec1b7b4c25611ac8326c5afe588944576c9142dd12bc574f6","body":"sha256:c76e747021190c170b5673e41f684d8fb7ba65a037636b07a0bbca43d374e84c"}
func TestRunSampleOutcomes(t *testing.T) {
	task, err := LoadTask("testdata/tasks/since")
	if err != nil {
		t.Fatal(err)
	}
	write := func(content string) func(context.Context, string) error {
		return func(_ context.Context, dir string) error {
			return os.WriteFile(filepath.Join(dir, "since", "parse.go"), []byte(content), 0o644)
		}
	}
	cases := []struct {
		name           string
		grow           func(context.Context, string) error
		compiled, pass bool
		errPart        string
	}{
		{"does not compile", write("package since\n\nfunc broken( {\n"), false, false, ""},
		{"compiles but fails the tests", write("package since\n\nimport \"time\"\n\nfunc parseSince(s string, now time.Time) (time.Time, error) { return now, nil }\n"), true, false, ""},
		{"the grow fails", func(context.Context, string) error { return errors.New("model unreachable\nsecond line") }, false, false, "model unreachable"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := RunSample(context.Background(), Config{Generators: "../../generators", Grow: tc.grow}, task, VariantBehaviour, 3)
			if err != nil {
				t.Fatal(err)
			}
			if s.Compiled != tc.compiled || s.Passed != tc.pass || s.N != 3 || s.Variant != VariantBehaviour {
				t.Fatalf("sample = %+v", s)
			}
			if tc.errPart != "" && s.Error != tc.errPart {
				t.Fatalf("error = %q, want %q", s.Error, tc.errPart)
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:36e1a83c85979b86a1d983a8b7b8a9c818fb8f21b52e5dfffa2860d20230971c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/sample_test.go","name":"TestRunAndPlan"},"kwargs":"sha256:778ff71883adfc200c2448cc8e54fb1b173dca196815bca7a0280b702ef7cd2b","body":"sha256:a482f528e40e49cf9a2ec7b6f4bf6c949e50ec9ef5f86939ca9a2f2152d2816b"}
func TestRunAndPlan(t *testing.T) {
	task, err := LoadTask("testdata/tasks/since")
	if err != nil {
		t.Fatal(err)
	}
	cfg := Config{Generators: "../../generators", Grow: referenceGrow(t, task)}
	samples, err := Run(context.Background(), cfg, []Task{task}, []Variant{VariantCode, VariantSignatures}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) != 4 || samples[0].Variant != VariantCode || samples[3].Variant != VariantSignatures || samples[3].N != 2 {
		t.Fatalf("samples = %+v", samples)
	}
	if got := Plan([]Task{task}, Variants(), 5, 20000); got != "1 task(s) x 3 variant(s) x 5 sample(s) = 15 grow(s), at most ~300000 tokens in total\n" {
		t.Fatalf("plan = %q", got)
	}
}

// /sedum:setFunction
