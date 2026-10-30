package pipeline

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup_test.go","path":"context"},"kwargs":"sha256:7c62edde09d02adfdf213dff4211c4d0f56a566ceb77b581d09d8df7aef35d5d","body":"sha256:0f7443fb25d7ad6a970617d4d0df6123283e24801b94c1d84fbe431230f31e1e"}
import "context"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup_test.go","path":"os"},"kwargs":"sha256:d7da66f80ba547c9c952754cd5537e8c3fdc74527f237ebfe02980aca2714283","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup_test.go","path":"path/filepath"},"kwargs":"sha256:443eb938aabfaed71a610738a1d22d74106aac37b74fe1248ef855dc861f8941","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup_test.go","path":"testing"},"kwargs":"sha256:e373363c4c0ae9e39bec9ac508eda4be4d397be03ce06c932d823becabdb975d","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup_test.go","path":"github.com/livecodelife/sedum/internal/resolve"},"kwargs":"sha256:f507ac83c5779a72d82c28ff7d3b5c2266402b9c495bef4cc9e66440c9b621dc","body":"sha256:e12df022982be7cfcdb6de1e42d5d1480f2b3d08431010742e4e2d9d5f13e140"}
import "github.com/livecodelife/sedum/internal/resolve"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/cleanup_test.go","name":"TestUndoRemovesWhatTheRunCreated"},"kwargs":"sha256:37dca5d24303c83fc67ece70957acfa7ef6dc837720fd13b1c50ce52c694493f","body":"sha256:fa8e029506b3f8bc901674925cb60f3838167a81aad96da69dfcddd285691fa4"}
func TestUndoRemovesWhatTheRunCreated(t *testing.T) {
	out := t.TempDir()
	if err := os.MkdirAll(filepath.Join(out, "app", "models"), 0o755); err != nil {
		t.Fatal(err)
	}
	c := newRunCleanup(out, []string{"app/models/user.rb", "lib/deep/x.rb"})
	if err := os.MkdirAll(filepath.Join(out, "lib", "deep"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"app/models/user.rb", "lib/deep/x.rb"} {
		if err := os.WriteFile(filepath.Join(out, filepath.FromSlash(p)), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	c.created([]resolve.File{
		{Resolution: resolve.Resolution{Path: "app/models/user.rb"}},
		{Resolution: resolve.Resolution{Path: "lib/deep/x.rb"}},
	})
	c.undo(nil)
	if _, err := os.Stat(filepath.Join(out, "lib")); !os.IsNotExist(err) {
		t.Fatal("lib created by the run should be removed")
	}
	if _, err := os.Stat(filepath.Join(out, "app", "models")); err != nil {
		t.Fatal("pre-existing app/models should be kept")
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/cleanup_test.go","name":"TestUndoNeverTouchesPreexistingFile"},"kwargs":"sha256:38492a6672acaadc89135da81dbbfc842d3b48d8358ead71d9f9bbf04cf85cad","body":"sha256:8893e82d5c62db1040b4dce985e87781f2a7e430d37b19ed438cfd97a58366ed"}
func TestUndoNeverTouchesPreexistingFile(t *testing.T) {
	out := t.TempDir()
	dir := filepath.Join(out, "app", "models")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, "user.rb")
	if err := os.WriteFile(file, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := newRunCleanup(out, []string{"app/models/user.rb"})
	c.created([]resolve.File{{Resolution: resolve.Resolution{Path: "app/models/user.rb"}, Existed: true}})
	c.undo(nil)
	got, err := os.ReadFile(file)
	if err != nil || string(got) != "old" {
		t.Fatalf("pre-existing file changed: %q, %v", got, err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/cleanup_test.go","name":"TestFailingRunLeavesOutputAsItWas"},"kwargs":"sha256:40321a035f912b053fb517c7e1a1ea15a03d976863bf9cf80139662dafcf7b40","body":"sha256:1624ab6ecfa5d3bbd0b2907c155de32ac2694a40d92bf07ff022f15d735e3233"}
func TestFailingRunLeavesOutputAsItWas(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: `{"invocations":[{"action":"addField","kwargs":{"model":"order","field":"total"}}]}`}
	_, err := Run(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected the run to fail")
	}
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Fatalf("failed run left files behind: %v", got)
	}
	for _, d := range []string{"app", filepath.Join("app", "models")} {
		if _, err := os.Stat(filepath.Join(cfg.Output, d)); !os.IsNotExist(err) {
			t.Fatalf("directory %s left behind", d)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/cleanup_test.go","name":"TestFailingRunKeepsPreexistingFile"},"kwargs":"sha256:48c98de1e19c88475ebfb3c6ef2e57fafdac8fc1a2a04573bda262f1c52047ad","body":"sha256:885d09789dcdb6f5ca19358ac7beff4e30908805d0e15cc3dd26a422394f8f3d"}
func TestFailingRunKeepsPreexistingFile(t *testing.T) {
	cfg := config(t)
	dir := filepath.Join(cfg.Output, "app", "models")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(dir, "other.rb")
	if err := os.WriteFile(other, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg.Client = &stub{response: `{"invocations":[{"action":"addField","kwargs":{"model":"order","field":"total"}}]}`}
	if _, err := Run(context.Background(), cfg); err == nil {
		t.Fatal("expected the run to fail")
	}
	got, err := os.ReadFile(other)
	if err != nil || string(got) != "keep" {
		t.Fatalf("pre-existing file changed: %q, %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(dir, "user.rb")); !os.IsNotExist(err) {
		t.Fatal("created user.rb should be gone")
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/cleanup_test.go","name":"TestSucceedingRunKeepsFileAndDryRunCreatesNothing"},"kwargs":"sha256:fe30ca0b9cfe55b859ab69cf2a854b39cc60bcd69840ab72392a18ed5c29e762","body":"sha256:5c819e772383a836cfa0a21b77cb64d0b34246ab89533d410ba807d5db85c45b"}
func TestSucceedingRunKeepsFileAndDryRunCreatesNothing(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{}
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(cfg.Output, "app", "models", "user.rb")); err != nil {
		t.Fatalf("successful run dropped its file: %v", err)
	}
	dry := config(t)
	dry.DryRun = true
	if _, err := Run(context.Background(), dry); err != nil {
		t.Fatal(err)
	}
	if got := tree(t, dry.Output); len(got) != 0 {
		t.Fatalf("dry run created files: %v", got)
	}
}

// /sedum:setFunction
