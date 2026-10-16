package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestPrepare(t *testing.T) {\n\ttask, err := LoadTask(writeTask(t, t.TempDir(), \"ok\", validTask))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tdir := t.TempDir()\n\tif err := Prepare(task, VariantSignatures, \"../../generators\", dir); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tmod, err := os.ReadFile(filepath.Join(dir, \"go.mod\"))\n\tif err != nil || !strings.Contains(string(mod), \"module design/ok\") {\n\t\tt.Fatalf(\"go.mod = %q, %v\", mod, err)\n\t}\n\tfor _, f := range []string{\"p/ctx.go\", \"PROMPT.md\", \"generators/go/sedum.yaml\", \"provenance/design-ok-signatures.yml\"} {\n\t\tif _, err := os.Stat(filepath.Join(dir, f)); err != nil {\n\t\t\tt.Errorf(\"missing %s: %v\", f, err)\n\t\t}\n\t}\n\trec, err := os.ReadFile(filepath.Join(dir, \"provenance\", \"design-ok-signatures.yml\"))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfor _, want := range []string{\"status: open\", \"intent: s\", \"- p/a.go\"} {\n\t\tif !strings.Contains(string(rec), want) {\n\t\t\tt.Errorf(\"record lacks %q:\\n%s\", want, rec)\n\t\t}\n\t}\n\tprompt, err := os.ReadFile(filepath.Join(dir, \"PROMPT.md\"))\n\tif err != nil || !strings.Contains(string(prompt), \"p/a.go\") {\n\t\tt.Fatalf(\"PROMPT.md must list the target paths: %q, %v\", prompt, err)\n\t}\n\tif _, err := os.Stat(filepath.Join(dir, \"p\", \"a_test.go\")); err == nil {\n\t\tt.Fatal(\"the hidden test must not be written before grading\")\n\t}\n\tif err := Prepare(task, \"nope\", \"../../generators\", t.TempDir()); err == nil {\n\t\tt.Fatal(\"an unknown variant must be an error\")\n\t}\n}","file":"evals/design/prepare_test.go","name":"TestPrepare"}}
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
	for _, f := range []string{"p/ctx.go", "PROMPT.md", "generators/go/sedum.yaml", "provenance/design-ok-signatures.yml"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("missing %s: %v", f, err)
		}
	}
	rec, err := os.ReadFile(filepath.Join(dir, "provenance", "design-ok-signatures.yml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"status: open", "intent: s", "- p/a.go"} {
		if !strings.Contains(string(rec), want) {
			t.Errorf("record lacks %q:\n%s", want, rec)
		}
	}
	prompt, err := os.ReadFile(filepath.Join(dir, "PROMPT.md"))
	if err != nil || !strings.Contains(string(prompt), "p/a.go") {
		t.Fatalf("PROMPT.md must list the target paths: %q, %v", prompt, err)
	}
	if _, err := os.Stat(filepath.Join(dir, "p", "a_test.go")); err == nil {
		t.Fatal("the hidden test must not be written before grading")
	}
	if err := Prepare(task, "nope", "../../generators", t.TempDir()); err == nil {
		t.Fatal("an unknown variant must be an error")
	}
}

// /sedum:setFunction
