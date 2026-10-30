package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/prepare_test.go","path":"os"},"kwargs":"sha256:74747564df59b8b51bd2b0f3fe37d0558ee11e0bf4f3c1a9542db52ae0ba0d43","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/prepare_test.go","path":"path/filepath"},"kwargs":"sha256:1112ecaa2fb6f9145283026b1a7a2b84a5720162bc30b2a4a66951c43b6d8c50","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/prepare_test.go","path":"strings"},"kwargs":"sha256:56e02d812818ebf79e8a13e4b83450bf6aa750faac0c1d002095897db1bb666e","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/prepare_test.go","path":"testing"},"kwargs":"sha256:0953def59dfd115a372ed689cc83cb7114947c71b478307536e85c70dd87f536","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/prepare_test.go","name":"TestPrepare"},"kwargs":"sha256:874280e2ae8c7cd4dfff4ef876b4609e30aece8b1ead455decc039833dcf62e9","body":"sha256:1d7e85b0ceceae828f19d004eaea999f8df887817db9b180e735f096cbb833c6"}
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
