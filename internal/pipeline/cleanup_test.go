package pipeline

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup_test.go","path":"context"}}
import "context"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/cleanup_test.go","path":"github.com/livecodelife/sedum/internal/resolve"}}
import "github.com/livecodelife/sedum/internal/resolve"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestUndoRemovesWhatTheRunCreated(t *testing.T) {\n\tout := t.TempDir()\n\tif err := os.MkdirAll(filepath.Join(out, \"app\", \"models\"), 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tc := newRunCleanup(out, []string{\"app/models/user.rb\", \"lib/deep/x.rb\"})\n\tif err := os.MkdirAll(filepath.Join(out, \"lib\", \"deep\"), 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfor _, p := range []string{\"app/models/user.rb\", \"lib/deep/x.rb\"} {\n\t\tif err := os.WriteFile(filepath.Join(out, filepath.FromSlash(p)), []byte(\"x\"), 0o644); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\tc.created([]resolve.File{\n\t\t{Resolution: resolve.Resolution{Path: \"app/models/user.rb\"}},\n\t\t{Resolution: resolve.Resolution{Path: \"lib/deep/x.rb\"}},\n\t})\n\tc.undo(nil)\n\tif _, err := os.Stat(filepath.Join(out, \"lib\")); !os.IsNotExist(err) {\n\t\tt.Fatal(\"lib created by the run should be removed\")\n\t}\n\tif _, err := os.Stat(filepath.Join(out, \"app\", \"models\")); err != nil {\n\t\tt.Fatal(\"pre-existing app/models should be kept\")\n\t}\n}","file":"internal/pipeline/cleanup_test.go","name":"TestUndoRemovesWhatTheRunCreated"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestUndoNeverTouchesPreexistingFile(t *testing.T) {\n\tout := t.TempDir()\n\tdir := filepath.Join(out, \"app\", \"models\")\n\tif err := os.MkdirAll(dir, 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfile := filepath.Join(dir, \"user.rb\")\n\tif err := os.WriteFile(file, []byte(\"old\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tc := newRunCleanup(out, []string{\"app/models/user.rb\"})\n\tc.created([]resolve.File{{Resolution: resolve.Resolution{Path: \"app/models/user.rb\"}, Existed: true}})\n\tc.undo(nil)\n\tgot, err := os.ReadFile(file)\n\tif err != nil || string(got) != \"old\" {\n\t\tt.Fatalf(\"pre-existing file changed: %q, %v\", got, err)\n\t}\n}","file":"internal/pipeline/cleanup_test.go","name":"TestUndoNeverTouchesPreexistingFile"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestFailingRunLeavesOutputAsItWas(t *testing.T) {\n\tcfg := config(t)\n\tcfg.Client = &stub{response: `{\"invocations\":[{\"action\":\"addField\",\"kwargs\":{\"model\":\"order\",\"field\":\"total\"}}]}`}\n\t_, err := Run(context.Background(), cfg)\n\tif err == nil {\n\t\tt.Fatal(\"expected the run to fail\")\n\t}\n\tif got := tree(t, cfg.Output); len(got) != 0 {\n\t\tt.Fatalf(\"failed run left files behind: %v\", got)\n\t}\n\tfor _, d := range []string{\"app\", filepath.Join(\"app\", \"models\")} {\n\t\tif _, err := os.Stat(filepath.Join(cfg.Output, d)); !os.IsNotExist(err) {\n\t\t\tt.Fatalf(\"directory %s left behind\", d)\n\t\t}\n\t}\n}","file":"internal/pipeline/cleanup_test.go","name":"TestFailingRunLeavesOutputAsItWas"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestFailingRunKeepsPreexistingFile(t *testing.T) {\n\tcfg := config(t)\n\tdir := filepath.Join(cfg.Output, \"app\", \"models\")\n\tif err := os.MkdirAll(dir, 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tother := filepath.Join(dir, \"other.rb\")\n\tif err := os.WriteFile(other, []byte(\"keep\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tcfg.Client = &stub{response: `{\"invocations\":[{\"action\":\"addField\",\"kwargs\":{\"model\":\"order\",\"field\":\"total\"}}]}`}\n\tif _, err := Run(context.Background(), cfg); err == nil {\n\t\tt.Fatal(\"expected the run to fail\")\n\t}\n\tgot, err := os.ReadFile(other)\n\tif err != nil || string(got) != \"keep\" {\n\t\tt.Fatalf(\"pre-existing file changed: %q, %v\", got, err)\n\t}\n\tif _, err := os.Stat(filepath.Join(dir, \"user.rb\")); !os.IsNotExist(err) {\n\t\tt.Fatal(\"created user.rb should be gone\")\n\t}\n}","file":"internal/pipeline/cleanup_test.go","name":"TestFailingRunKeepsPreexistingFile"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSucceedingRunKeepsFileAndDryRunCreatesNothing(t *testing.T) {\n\tcfg := config(t)\n\tcfg.Client = &stub{}\n\tif _, err := Run(context.Background(), cfg); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif _, err := os.Stat(filepath.Join(cfg.Output, \"app\", \"models\", \"user.rb\")); err != nil {\n\t\tt.Fatalf(\"successful run dropped its file: %v\", err)\n\t}\n\tdry := config(t)\n\tdry.DryRun = true\n\tif _, err := Run(context.Background(), dry); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif got := tree(t, dry.Output); len(got) != 0 {\n\t\tt.Fatalf(\"dry run created files: %v\", got)\n\t}\n}","file":"internal/pipeline/cleanup_test.go","name":"TestSucceedingRunKeepsFileAndDryRunCreatesNothing"}}
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
