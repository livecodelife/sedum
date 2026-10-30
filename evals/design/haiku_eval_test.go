package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/haiku_eval_test.go","path":"context"},"kwargs":"sha256:4ba13ef2837961a6ae3a6ce421bfca21c26f3c9412c44baf7537c448ed1c095d","body":"sha256:0f7443fb25d7ad6a970617d4d0df6123283e24801b94c1d84fbe431230f31e1e"}
import "context"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/haiku_eval_test.go","path":"os"},"kwargs":"sha256:b5bcc27fe72027ec93c961b23f19bdbdf36075148565d6e499cca014d599324f","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/haiku_eval_test.go","path":"path/filepath"},"kwargs":"sha256:286992022d7bbf1e1b6d595c9e90afda7b0837371411b70033b774e5b7c93f7f","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/haiku_eval_test.go","path":"strings"},"kwargs":"sha256:b7dc2e3937c1d5fb55cb03128637e7af76e854bd10adede0677f71d0e589387b","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/haiku_eval_test.go","path":"testing"},"kwargs":"sha256:1f3a49827c9819406c33d8198d1075a3c7ede5dd52351b3b33f8b5aa002e3c8f","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/haiku_eval_test.go","name":"fakeWrapper"},"kwargs":"sha256:0533efcdd8af05dfcc0865fd271a9287c0878f2b83107d291bd27955e836c7e4","body":"sha256:81e2049bdccaf4b77fbd76a0694a2a42806778666958b87ffb508d2831a2f88d"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/haiku_eval_test.go","name":"TestHaikuGrowRunsTheWrapperWithoutAModelFlag"},"kwargs":"sha256:921695851306f223f983b6120c9bb18a18ad0995bba05107e5486afb4a40a8ac","body":"sha256:4dd923f1e38501056c6a7e09e8621f522a9e6882c382da66b20345f61757b4a5"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/haiku_eval_test.go","name":"TestHaikuGrowReportsTheLastOutputLine"},"kwargs":"sha256:13b4ca1f85f6d8e09cb1eae46e80ea278e500962e2fc2ef1317113d4ea890f20","body":"sha256:833b5920bf15cf6728541e7c1c7e21da44fabb4cfe737e5af347715af04462c2"}
func TestHaikuGrowReportsTheLastOutputLine(t *testing.T) {
	dir := t.TempDir()
	wrapper := fakeWrapper(t, dir, "echo first\necho boom\nexit 1\n")
	err := HaikuGrow(wrapper, "/bin/sedum-x")(context.Background(), dir)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("err = %v, want boom", err)
	}
}

// /sedum:setFunction
