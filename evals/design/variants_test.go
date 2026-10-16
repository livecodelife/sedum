package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/variants_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/variants_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/variants_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/variants_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestIntentVariants(t *testing.T) {\n\ttask := Task{Behaviour: \"b\", Signatures: \"s\", Code: \"c\"}\n\tfor v, want := range map[Variant]string{VariantBehaviour: \"b\", VariantSignatures: \"s\", VariantCode: \"c\"} {\n\t\tgot, err := task.Intent(v)\n\t\tif err != nil || got != want {\n\t\t\tt.Errorf(\"Intent(%s) = %q, %v, want %q\", v, got, err, want)\n\t\t}\n\t}\n\tif _, err := task.Intent(\"nope\"); err == nil {\n\t\tt.Fatal(\"an unknown variant must be an error\")\n\t}\n\tif got := Variants(); len(got) != 3 || got[0] != VariantBehaviour || got[2] != VariantCode {\n\t\tt.Fatalf(\"Variants() = %v\", got)\n\t}\n}","file":"evals/design/variants_test.go","name":"TestIntentVariants"}}
func TestIntentVariants(t *testing.T) {
	task := Task{Behaviour: "b", Signatures: "s", Code: "c"}
	for v, want := range map[Variant]string{VariantBehaviour: "b", VariantSignatures: "s", VariantCode: "c"} {
		got, err := task.Intent(v)
		if err != nil || got != want {
			t.Errorf("Intent(%s) = %q, %v, want %q", v, got, err, want)
		}
	}
	if _, err := task.Intent("nope"); err == nil {
		t.Fatal("an unknown variant must be an error")
	}
	if got := Variants(); len(got) != 3 || got[0] != VariantBehaviour || got[2] != VariantCode {
		t.Fatalf("Variants() = %v", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestIsolatedRealTasks(t *testing.T) {\n\ttasks, err := LoadTasks(\"testdata/tasks\")\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfor _, task := range tasks {\n\t\tif err := task.Isolated(); err != nil {\n\t\t\tt.Errorf(\"%s: %v\", task.ID, err)\n\t\t}\n\t}\n}","file":"evals/design/variants_test.go","name":"TestIsolatedRealTasks"}}
func TestIsolatedRealTasks(t *testing.T) {
	tasks, err := LoadTasks("testdata/tasks")
	if err != nil {
		t.Fatal(err)
	}
	for _, task := range tasks {
		if err := task.Isolated(); err != nil {
			t.Errorf("%s: %v", task.ID, err)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestIsolatedCatchesLeaks(t *testing.T) {\n\troot := t.TempDir()\n\tdir := writeTask(t, root, \"ok\", validTask)\n\ttask, err := LoadTask(dir)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := task.Isolated(); err != nil {\n\t\tt.Fatalf(\"a clean task: %v\", err)\n\t}\n\tnamed := task\n\tnamed.Signatures = \"see a_test.go\"\n\tif err := named.Isolated(); err == nil || !strings.Contains(err.Error(), \"a_test.go\") {\n\t\tt.Fatalf(\"a named hidden file: %v\", err)\n\t}\n\tleaked := task\n\tleaked.Behaviour = \"copy this\\nfunc TestA(t *testing.T) {}\\n\"\n\tif err := leaked.Isolated(); err == nil || !strings.Contains(err.Error(), \"func TestA\") {\n\t\tt.Fatalf(\"a leaked test line: %v\", err)\n\t}\n\tif err := os.Remove(filepath.Join(dir, \"hidden\", \"a_test.go\")); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := task.Isolated(); err == nil {\n\t\tt.Fatal(\"an unreadable hidden file must be an error\")\n\t}\n}","file":"evals/design/variants_test.go","name":"TestIsolatedCatchesLeaks"}}
func TestIsolatedCatchesLeaks(t *testing.T) {
	root := t.TempDir()
	dir := writeTask(t, root, "ok", validTask)
	task, err := LoadTask(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := task.Isolated(); err != nil {
		t.Fatalf("a clean task: %v", err)
	}
	named := task
	named.Signatures = "see a_test.go"
	if err := named.Isolated(); err == nil || !strings.Contains(err.Error(), "a_test.go") {
		t.Fatalf("a named hidden file: %v", err)
	}
	leaked := task
	leaked.Behaviour = "copy this\nfunc TestA(t *testing.T) {}\n"
	if err := leaked.Isolated(); err == nil || !strings.Contains(err.Error(), "func TestA") {
		t.Fatalf("a leaked test line: %v", err)
	}
	if err := os.Remove(filepath.Join(dir, "hidden", "a_test.go")); err != nil {
		t.Fatal(err)
	}
	if err := task.Isolated(); err == nil {
		t.Fatal("an unreadable hidden file must be an error")
	}
}

// /sedum:setFunction
