package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/task_test.go","path":"os"},"kwargs":"sha256:a5f7d900ea70138d0d6bc5930335a6cb26a64f59080ff8b9c110279b175c7bd6","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/task_test.go","path":"path/filepath"},"kwargs":"sha256:7d6efb56ca6b1eeea70dce30282656262becd41c5fdef6856cd3eec46fc84f27","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/task_test.go","path":"strings"},"kwargs":"sha256:c6d90c359c38b3b8bf192684995dd574c13eff3e7fcfb9f6bc4a81b8071aa5c9","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/task_test.go","path":"testing"},"kwargs":"sha256:cf67251da33a66d83f5bd59f6f5e76ca1d0ffe68b3ce98aed15cc834fbb83ae4","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/task_test.go","name":"writeTask"},"kwargs":"sha256:4449cbdc6d27c942198c701ee303878bc4a0c6c6f7bbbb0dd9cfca74bc8f8ed0","body":"sha256:5bdf4fd743982cc83926167285231a7d91acafce7bb35a86e896117c241119bc"}
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
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"evals/design/task_test.go","name":"validTask"},"kwargs":"sha256:9d040c993c3fc2a659254718f5e81732d6299e8056fd8f232fb3ec306e99ec69","body":"sha256:ea78c3a1219f2195693c88d83e7049372118699461ab3dc93006c075d67eb3b6"}
const validTask = "id: ok\npackage: p\ntarget: [a.go]\ncontext: [ctx.go]\nhidden: [a_test.go]\nbehaviour: b\nsignatures: s\ncode: c\n"

// /sedum:setConst
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/task_test.go","name":"TestLoadTask"},"kwargs":"sha256:92b2ffa51ab0f2f28fde1d638a68c36e52652ed9c677529e5e4fdb734efb723a","body":"sha256:0dd9fe06ac763f4f4ada377d8977e88074922d9faf9bdd313f42a01eed1539cc"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/task_test.go","name":"TestLoadTaskRejects"},"kwargs":"sha256:bdb83e599dd49652f4c0a4aba96c6c7a77c17ccc3d5ee9855c6dfe62cf3b7601","body":"sha256:52a1494e9038d2a6abd7ac78551929e2567139df4278991f0bb239a667edca98"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/task_test.go","name":"TestLoadTasks"},"kwargs":"sha256:548fad2314f4aadbf9442e0885fb17186c1b37c9a2b45d68a8d711aa276653e7","body":"sha256:a7932abb841459e2145297ef558e54e774141c5528cc26d1b33094d66b1eb7c4"}
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
