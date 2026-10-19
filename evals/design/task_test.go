package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/task_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/task_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/task_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/task_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// writeTask writes a minimal valid task into root/id and returns its dir.\nfunc writeTask(t *testing.T, root, id, yaml string) string {\n\tt.Helper()\n\tdir := filepath.Join(root, id)\n\tfor _, d := range []string{\"context\", \"hidden\"} {\n\t\tif err := os.MkdirAll(filepath.Join(dir, d), 0o755); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\tfiles := map[string]string{\n\t\t\"task.yaml\":        yaml,\n\t\t\"hidden/a_test.go\": \"package p\\n\\nimport \\\"testing\\\"\\n\\nfunc TestA(t *testing.T) {}\\n\",\n\t\t\"context/ctx.go\":   \"package p\\n\",\n\t}\n\tfor name, content := range files {\n\t\tif err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\treturn dir\n}","file":"evals/design/task_test.go","name":"writeTask"}}
// writeTask writes a minimal valid task into root/id and returns its dir.
func writeTask(t *testing.T, root, id, yaml string) string {
	t.Helper()
	dir := filepath.Join(root, id)
	for _, d := range []string{"context", "hidden"} {
		if err := os.MkdirAll(filepath.Join(dir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	files := map[string]string{
		"task.yaml":        yaml,
		"hidden/a_test.go": "package p\n\nimport \"testing\"\n\nfunc TestA(t *testing.T) {}\n",
		"context/ctx.go":   "package p\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// /sedum:setFunction
// sedum:setConst {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"const validTask = \"id: ok\\npackage: p\\ntarget: [a.go]\\ncontext: [ctx.go]\\nhidden: [a_test.go]\\nbehaviour: b\\nsignatures: s\\ncode: c\\n\"","file":"evals/design/task_test.go","name":"validTask"}}
const validTask = "id: ok\npackage: p\ntarget: [a.go]\ncontext: [ctx.go]\nhidden: [a_test.go]\nbehaviour: b\nsignatures: s\ncode: c\n"

// /sedum:setConst
// sedum:setFunction {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestLoadTask(t *testing.T) {\n\troot := t.TempDir()\n\tdir := writeTask(t, root, \"ok\", validTask)\n\ttask, err := LoadTask(dir)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif task.ID != \"ok\" || task.Package != \"p\" || task.Dir != dir || len(task.Target) != 1 || task.Code != \"c\" {\n\t\tt.Fatalf(\"task = %+v\", task)\n\t}\n}","file":"evals/design/task_test.go","name":"TestLoadTask"}}
func TestLoadTask(t *testing.T) {
	root := t.TempDir()
	dir := writeTask(t, root, "ok", validTask)
	task, err := LoadTask(dir)
	if err != nil {
		t.Fatal(err)
	}
	if task.ID != "ok" || task.Package != "p" || task.Dir != dir || len(task.Target) != 1 || task.Code != "c" {
		t.Fatalf("task = %+v", task)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestLoadTaskRejects(t *testing.T) {\n\tcases := []struct {\n\t\tname, id, yaml, want string\n\t}{\n\t\t{\"id differs from directory\", \"ok\", strings.Replace(validTask, \"id: ok\", \"id: other\", 1), \"directory name\"},\n\t\t{\"no target\", \"ok\", strings.Replace(validTask, \"target: [a.go]\", \"target: []\", 1), \"required\"},\n\t\t{\"missing intent\", \"ok\", strings.Replace(validTask, \"code: c\\n\", \"\", 1), \"all required\"},\n\t\t{\"hidden is not a test file\", \"ok\", strings.Replace(validTask, \"hidden: [a_test.go]\", \"hidden: [ctx.go]\", 1), \"_test.go\"},\n\t\t{\"missing hidden file\", \"ok\", strings.Replace(validTask, \"a_test.go\", \"b_test.go\", 1), \"no such file\"},\n\t\t{\"target is also context\", \"ok\", strings.Replace(validTask, \"target: [a.go]\", \"target: [ctx.go]\", 1), \"also a context or hidden\"},\n\t}\n\tfor _, tc := range cases {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\tdir := writeTask(t, t.TempDir(), tc.id, tc.yaml)\n\t\t\t_, err := LoadTask(dir)\n\t\t\tif err == nil || !strings.Contains(err.Error(), tc.want) {\n\t\t\t\tt.Fatalf(\"err = %v, want it to contain %q\", err, tc.want)\n\t\t\t}\n\t\t})\n\t}\n}","file":"evals/design/task_test.go","name":"TestLoadTaskRejects"}}
func TestLoadTaskRejects(t *testing.T) {
	cases := []struct {
		name, id, yaml, want string
	}{
		{"id differs from directory", "ok", strings.Replace(validTask, "id: ok", "id: other", 1), "directory name"},
		{"no target", "ok", strings.Replace(validTask, "target: [a.go]", "target: []", 1), "required"},
		{"missing intent", "ok", strings.Replace(validTask, "code: c\n", "", 1), "all required"},
		{"hidden is not a test file", "ok", strings.Replace(validTask, "hidden: [a_test.go]", "hidden: [ctx.go]", 1), "_test.go"},
		{"missing hidden file", "ok", strings.Replace(validTask, "a_test.go", "b_test.go", 1), "no such file"},
		{"target is also context", "ok", strings.Replace(validTask, "target: [a.go]", "target: [ctx.go]", 1), "also a context or hidden"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := writeTask(t, t.TempDir(), tc.id, tc.yaml)
			_, err := LoadTask(dir)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want it to contain %q", err, tc.want)
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestLoadTasks(t *testing.T) {\n\ttasks, err := LoadTasks(\"testdata/tasks\")\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tvar ids []string\n\tfor _, task := range tasks {\n\t\tids = append(ids, task.ID)\n\t}\n\tif got := strings.Join(ids, \",\"); got != \"append,read,render,since,summarize\" {\n\t\tt.Fatalf(\"tasks = %s\", got)\n\t}\n\tif _, err := LoadTasks(filepath.Join(t.TempDir(), \"none\")); err == nil {\n\t\tt.Fatal(\"a missing tasks directory must be an error\")\n\t}\n}","file":"evals/design/task_test.go","name":"TestLoadTasks"}}
func TestLoadTasks(t *testing.T) {
	tasks, err := LoadTasks("testdata/tasks")
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, task := range tasks {
		ids = append(ids, task.ID)
	}
	if got := strings.Join(ids, ","); got != "append,read,render,since,summarize" {
		t.Fatalf("tasks = %s", got)
	}
	if _, err := LoadTasks(filepath.Join(t.TempDir(), "none")); err == nil {
		t.Fatal("a missing tasks directory must be an error")
	}
}

// /sedum:setFunction
