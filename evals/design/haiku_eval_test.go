package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/haiku_eval_test.go","path":"context"}}
import "context"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/haiku_eval_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/haiku_eval_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/haiku_eval_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/haiku_eval_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// fakeWrapper writes an executable shell script and returns its path.\nfunc fakeWrapper(t *testing.T, dir, body string) string {\n\tt.Helper()\n\tpath := filepath.Join(dir, \"wrapper.sh\")\n\tif err := os.WriteFile(path, []byte(\"#!/bin/sh\\n\"+body), 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\treturn path\n}","file":"evals/design/haiku_eval_test.go","name":"fakeWrapper"}}
// fakeWrapper writes an executable shell script and returns its path.
func fakeWrapper(t *testing.T, dir, body string) string {
	t.Helper()
	path := filepath.Join(dir, "wrapper.sh")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestHaikuGrowRunsTheWrapperWithoutAModelFlag(t *testing.T) {\n\tdir := t.TempDir()\n\targs := filepath.Join(dir, \"args\")\n\twrapper := fakeWrapper(t, dir, \"echo \\\"$@\\\" > \"+args+\"\\necho \\\"SEDUM_BIN=$SEDUM_BIN\\\" >> \"+args+\"\\n\")\n\tif err := HaikuGrow(wrapper, \"/bin/sedum-x\")(context.Background(), dir); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tdata, err := os.ReadFile(args)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tgot := string(data)\n\tfor _, want := range []string{\"grow --generators ./generators --records ./provenance --prompt-file PROMPT.md --verify go vet ./...\", \"SEDUM_BIN=/bin/sedum-x\"} {\n\t\tif !strings.Contains(got, want) {\n\t\t\tt.Errorf(\"wrapper arguments lack %q:\\n%s\", want, got)\n\t\t}\n\t}\n\tif strings.Contains(got, \"--model\") {\n\t\tt.Errorf(\"no model flag may be passed:\\n%s\", got)\n\t}\n}","file":"evals/design/haiku_eval_test.go","name":"TestHaikuGrowRunsTheWrapperWithoutAModelFlag"}}
func TestHaikuGrowRunsTheWrapperWithoutAModelFlag(t *testing.T) {
	dir := t.TempDir()
	args := filepath.Join(dir, "args")
	wrapper := fakeWrapper(t, dir, "echo \"$@\" > "+args+"\necho \"SEDUM_BIN=$SEDUM_BIN\" >> "+args+"\n")
	if err := HaikuGrow(wrapper, "/bin/sedum-x")(context.Background(), dir); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(args)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	for _, want := range []string{"grow --generators ./generators --records ./provenance --prompt-file PROMPT.md --verify go vet ./...", "SEDUM_BIN=/bin/sedum-x"} {
		if !strings.Contains(got, want) {
			t.Errorf("wrapper arguments lack %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "--model") {
		t.Errorf("no model flag may be passed:\n%s", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestHaikuGrowReportsTheLastOutputLine(t *testing.T) {\n\tdir := t.TempDir()\n\twrapper := fakeWrapper(t, dir, \"echo first\\necho boom\\nexit 1\\n\")\n\terr := HaikuGrow(wrapper, \"/bin/sedum-x\")(context.Background(), dir)\n\tif err == nil || err.Error() != \"boom\" {\n\t\tt.Fatalf(\"err = %v, want boom\", err)\n\t}\n}","file":"evals/design/haiku_eval_test.go","name":"TestHaikuGrowReportsTheLastOutputLine"}}
func TestHaikuGrowReportsTheLastOutputLine(t *testing.T) {
	dir := t.TempDir()
	wrapper := fakeWrapper(t, dir, "echo first\necho boom\nexit 1\n")
	err := HaikuGrow(wrapper, "/bin/sedum-x")(context.Background(), dir)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("err = %v, want boom", err)
	}
}

// /sedum:setFunction
