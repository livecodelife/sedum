package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample_test.go","path":"context"}}
import "context"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample_test.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/sample_test.go","path":"github.com/livecodelife/sedum/internal/growlog"}}
import "github.com/livecodelife/sedum/internal/growlog"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// referenceGrow writes the task's reference code, taken from its code intent,\n// as the target file, and logs one grow the way sedum grow does.\nfunc referenceGrow(t *testing.T, task Task) func(context.Context, string) error {\n\treturn func(_ context.Context, dir string) error {\n\t\tparts := strings.SplitN(task.Code, \":\\n\\n\", 2)\n\t\tif len(parts) != 2 {\n\t\t\tt.Fatalf(\"task %s: unexpected code intent\", task.ID)\n\t\t}\n\t\tif err := os.WriteFile(filepath.Join(dir, task.Target[0]), []byte(parts[1]), 0o644); err != nil {\n\t\t\treturn err\n\t\t}\n\t\treturn growlog.Append(filepath.Join(dir, \".sedum\", \"grow.jsonl\"), growlog.Entry{Calls: 2, Rejected: 1, Repairs: 1, PromptTokens: 100, CompletionTokens: 10})\n\t}\n}","file":"evals/design/sample_test.go","name":"referenceGrow"}}
// referenceGrow writes the task's reference code, taken from its code intent,
// as the target file, and logs one grow the way sedum grow does.
func referenceGrow(t *testing.T, task Task) func(context.Context, string) error {
	return func(_ context.Context, dir string) error {
		parts := strings.SplitN(task.Code, ":\n\n", 2)
		if len(parts) != 2 {
			t.Fatalf("task %s: unexpected code intent", task.ID)
		}
		if err := os.WriteFile(filepath.Join(dir, task.Target[0]), []byte(parts[1]), 0o644); err != nil {
			return err
		}
		return growlog.Append(filepath.Join(dir, ".sedum", "grow.jsonl"), growlog.Entry{Calls: 2, Rejected: 1, Repairs: 1, PromptTokens: 100, CompletionTokens: 10})
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestReferenceCodePassesEveryTask(t *testing.T) {\n\ttasks, err := LoadTasks(\"testdata/tasks\")\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfor _, task := range tasks {\n\t\tt.Run(task.ID, func(t *testing.T) {\n\t\t\tcfg := Config{Generators: \"../../generators\", Grow: referenceGrow(t, task)}\n\t\t\ts, err := RunSample(context.Background(), cfg, task, VariantCode, 1)\n\t\t\tif err != nil {\n\t\t\t\tt.Fatal(err)\n\t\t\t}\n\t\t\tif !s.Compiled || !s.Passed {\n\t\t\t\tt.Fatalf(\"the reference code must pass its own hidden tests: %+v\", s)\n\t\t\t}\n\t\t\tif s.Calls != 2 || s.Rejected != 1 || s.Repairs != 1 || s.PromptTokens != 100 || s.CompletionTokens != 10 {\n\t\t\t\tt.Fatalf(\"the grow log was not read: %+v\", s)\n\t\t\t}\n\t\t})\n\t}\n}","file":"evals/design/sample_test.go","name":"TestReferenceCodePassesEveryTask"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunSampleOutcomes(t *testing.T) {\n\ttask, err := LoadTask(\"testdata/tasks/since\")\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\twrite := func(content string) func(context.Context, string) error {\n\t\treturn func(_ context.Context, dir string) error {\n\t\t\treturn os.WriteFile(filepath.Join(dir, \"parse.go\"), []byte(content), 0o644)\n\t\t}\n\t}\n\tcases := []struct {\n\t\tname           string\n\t\tgrow           func(context.Context, string) error\n\t\tcompiled, pass bool\n\t\terrPart        string\n\t}{\n\t\t{\"does not compile\", write(\"package since\\n\\nfunc broken( {\\n\"), false, false, \"\"},\n\t\t{\"compiles but fails the tests\", write(\"package since\\n\\nimport \\\"time\\\"\\n\\nfunc parseSince(s string, now time.Time) (time.Time, error) { return now, nil }\\n\"), true, false, \"\"},\n\t\t{\"the grow fails\", func(context.Context, string) error { return errors.New(\"model unreachable\\nsecond line\") }, false, false, \"model unreachable\"},\n\t}\n\tfor _, tc := range cases {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\ts, err := RunSample(context.Background(), Config{Generators: \"../../generators\", Grow: tc.grow}, task, VariantBehaviour, 3)\n\t\t\tif err != nil {\n\t\t\t\tt.Fatal(err)\n\t\t\t}\n\t\t\tif s.Compiled != tc.compiled || s.Passed != tc.pass || s.N != 3 || s.Variant != VariantBehaviour {\n\t\t\t\tt.Fatalf(\"sample = %+v\", s)\n\t\t\t}\n\t\t\tif tc.errPart != \"\" && s.Error != tc.errPart {\n\t\t\t\tt.Fatalf(\"error = %q, want %q\", s.Error, tc.errPart)\n\t\t\t}\n\t\t})\n\t}\n}","file":"evals/design/sample_test.go","name":"TestRunSampleOutcomes"}}
func TestRunSampleOutcomes(t *testing.T) {
	task, err := LoadTask("testdata/tasks/since")
	if err != nil {
		t.Fatal(err)
	}
	write := func(content string) func(context.Context, string) error {
		return func(_ context.Context, dir string) error {
			return os.WriteFile(filepath.Join(dir, "parse.go"), []byte(content), 0o644)
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-e6976fa9","gen":"sha256:ebb1bd9ca8c21c30033cac7adfe5ee18f5ef9dd23e972d1502637991e68c5e50","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunAndPlan(t *testing.T) {\n\ttask, err := LoadTask(\"testdata/tasks/since\")\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tcfg := Config{Generators: \"../../generators\", Grow: referenceGrow(t, task)}\n\tsamples, err := Run(context.Background(), cfg, []Task{task}, []Variant{VariantCode, VariantSignatures}, 2)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(samples) != 4 || samples[0].Variant != VariantCode || samples[3].Variant != VariantSignatures || samples[3].N != 2 {\n\t\tt.Fatalf(\"samples = %+v\", samples)\n\t}\n\tif got := Plan([]Task{task}, Variants(), 5, 20000); got != \"1 task(s) x 3 variant(s) x 5 sample(s) = 15 grow(s), at most ~300000 tokens in total\\n\" {\n\t\tt.Fatalf(\"plan = %q\", got)\n\t}\n}","file":"evals/design/sample_test.go","name":"TestRunAndPlan"}}
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
