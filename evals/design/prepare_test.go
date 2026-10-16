package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestPrepare(t *testing.T) {\n\ttask, err := LoadTask(writeTask(t, t.TempDir(), \"ok\", validTask))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tdir := t.TempDir()\n\tif err := Prepare(task, VariantSignatures, \"../../generators\", dir); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tmod, err := os.ReadFile(filepath.Join(dir, \"go.mod\"))\n\tif err != nil || !strings.Contains(string(mod), \"module design/ok\") {\n\t\tt.Fatalf(\"go.mod = %q, %v\", mod, err)\n\t}\n\tfor _, f := range []string{\"ctx.go\", \"PROMPT.md\", \"generators/go/sedum.yaml\", \"provenance/design-ok-signatures.yml\"} {\n\t\tif _, err := os.Stat(filepath.Join(dir, f)); err != nil {\n\t\t\tt.Errorf(\"missing %s: %v\", f, err)\n\t\t}\n\t}\n\trec, err := os.ReadFile(filepath.Join(dir, \"provenance\", \"design-ok-signatures.yml\"))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfor _, want := range []string{\"status: open\", \"intent: s\", \"- a.go\"} {\n\t\tif !strings.Contains(string(rec), want) {\n\t\t\tt.Errorf(\"record lacks %q:\\n%s\", want, rec)\n\t\t}\n\t}\n\tif _, err := os.Stat(filepath.Join(dir, \"a_test.go\")); err == nil {\n\t\tt.Fatal(\"the hidden test must not be written before grading\")\n\t}\n\tif err := Prepare(task, \"nope\", \"../../generators\", t.TempDir()); err == nil {\n\t\tt.Fatal(\"an unknown variant must be an error\")\n\t}\n}","file":"evals/design/prepare_test.go","name":"TestPrepare"}}
func TestPrepare(t *testing.T) {
	task, err := LoadTask(writeTask(t, t.TempDir(), "ok", validTask))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := Prepare(task, VariantSignatures, "../../generators", dir); err != nil {
		t.Fatal(err)
	}
	mod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil || !strings.Contains(string(mod), "module design/ok") {
		t.Fatalf("go.mod = %q, %v", mod, err)
	}
	for _, f := range []string{"ctx.go", "PROMPT.md", "generators/go/sedum.yaml", "provenance/design-ok-signatures.yml"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("missing %s: %v", f, err)
		}
	}
	rec, err := os.ReadFile(filepath.Join(dir, "provenance", "design-ok-signatures.yml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"status: open", "intent: s", "- a.go"} {
		if !strings.Contains(string(rec), want) {
			t.Errorf("record lacks %q:\n%s", want, rec)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "a_test.go")); err == nil {
		t.Fatal("the hidden test must not be written before grading")
	}
	if err := Prepare(task, "nope", "../../generators", t.TempDir()); err == nil {
		t.Fatal("an unknown variant must be an error")
	}
}

// /sedum:setFunction
