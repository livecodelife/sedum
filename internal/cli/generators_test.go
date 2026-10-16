package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators_test.go","path":"reflect"}}
import "reflect"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func snapshotTree(t *testing.T, root string) map[string]string {\n\tt.Helper()\n\tout := map[string]string{}\n\terr := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tif info.IsDir() {\n\t\t\treturn nil\n\t\t}\n\t\tdata, err := os.ReadFile(p)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\trel, err := filepath.Rel(root, p)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tout[rel] = string(data)\n\t\treturn nil\n\t})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\treturn out\n}","file":"internal/cli/generators_test.go","name":"snapshotTree"}}
func snapshotTree(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		out[rel] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsInstallBuiltin(t *testing.T) {\n\tgen := filepath.Join(t.TempDir(), \"generators\")\n\tif out, err := exec(t, \"generators\", \"install\", \"go\", \"--dir\", gen); err != nil {\n\t\tt.Fatalf(\"install: %v\\n%s\", err, out)\n\t}\n\tfor _, p := range []string{filepath.Join(gen, \"go\", \"sedum.yaml\"), filepath.Join(gen, \"go.mod\")} {\n\t\tif _, err := os.Stat(p); err != nil {\n\t\t\tt.Errorf(\"missing %s: %v\", p, err)\n\t\t}\n\t}\n\tif _, _, err := genpkg.Load(gen, genpkg.Options{}); err != nil {\n\t\tt.Fatalf(\"load installed set: %v\", err)\n\t}\n}","file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallBuiltin"}}
func TestGeneratorsInstallBuiltin(t *testing.T) {
	gen := filepath.Join(t.TempDir(), "generators")
	if out, err := exec(t, "generators", "install", "go", "--dir", gen); err != nil {
		t.Fatalf("install: %v\n%s", err, out)
	}
	for _, p := range []string{filepath.Join(gen, "go", "sedum.yaml"), filepath.Join(gen, "go.mod")} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("missing %s: %v", p, err)
		}
	}
	if _, _, err := genpkg.Load(gen, genpkg.Options{}); err != nil {
		t.Fatalf("load installed set: %v", err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsInstallTwiceIsIdentical(t *testing.T) {\n\tgen := filepath.Join(t.TempDir(), \"generators\")\n\tif out, err := exec(t, \"generators\", \"install\", \"go\", \"--dir\", gen); err != nil {\n\t\tt.Fatalf(\"first install: %v\\n%s\", err, out)\n\t}\n\tbefore := snapshotTree(t, gen)\n\tout, err := exec(t, \"generators\", \"install\", \"go\", \"--dir\", gen)\n\tif err != nil {\n\t\tt.Fatalf(\"second install: %v\\n%s\", err, out)\n\t}\n\tif strings.Contains(out, \"create \") {\n\t\tt.Errorf(\"second install created files:\\n%s\", out)\n\t}\n\tif !reflect.DeepEqual(before, snapshotTree(t, gen)) {\n\t\tt.Error(\"second install changed the tree\")\n\t}\n}","file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallTwiceIsIdentical"}}
func TestGeneratorsInstallTwiceIsIdentical(t *testing.T) {
	gen := filepath.Join(t.TempDir(), "generators")
	if out, err := exec(t, "generators", "install", "go", "--dir", gen); err != nil {
		t.Fatalf("first install: %v\n%s", err, out)
	}
	before := snapshotTree(t, gen)
	out, err := exec(t, "generators", "install", "go", "--dir", gen)
	if err != nil {
		t.Fatalf("second install: %v\n%s", err, out)
	}
	if strings.Contains(out, "create ") {
		t.Errorf("second install created files:\n%s", out)
	}
	if !reflect.DeepEqual(before, snapshotTree(t, gen)) {
		t.Error("second install changed the tree")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsInstallRefusesDifferingFile(t *testing.T) {\n\tgen := filepath.Join(t.TempDir(), \"generators\")\n\ttarget := filepath.Join(gen, \"go\", \"sedum.yaml\")\n\tif err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := os.WriteFile(target, []byte(\"changed\\n\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tbefore := snapshotTree(t, gen)\n\tif _, err := exec(t, \"generators\", \"install\", \"go\", \"--dir\", gen); err == nil {\n\t\tt.Fatal(\"expected refusal\")\n\t}\n\tif !reflect.DeepEqual(before, snapshotTree(t, gen)) {\n\t\tt.Error(\"refused install still wrote files\")\n\t}\n\tif _, err := os.Stat(filepath.Join(gen, \"go.mod\")); err == nil {\n\t\tt.Error(\"go.mod was created despite refusal\")\n\t}\n}","file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallRefusesDifferingFile"}}
func TestGeneratorsInstallRefusesDifferingFile(t *testing.T) {
	gen := filepath.Join(t.TempDir(), "generators")
	target := filepath.Join(gen, "go", "sedum.yaml")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := snapshotTree(t, gen)
	if _, err := exec(t, "generators", "install", "go", "--dir", gen); err == nil {
		t.Fatal("expected refusal")
	}
	if !reflect.DeepEqual(before, snapshotTree(t, gen)) {
		t.Error("refused install still wrote files")
	}
	if _, err := os.Stat(filepath.Join(gen, "go.mod")); err == nil {
		t.Error("go.mod was created despite refusal")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsInstallForceOverwritesOnlyDifferingFiles(t *testing.T) {\n\tgen := filepath.Join(t.TempDir(), \"generators\")\n\tif out, err := exec(t, \"generators\", \"install\", \"go\", \"--dir\", gen); err != nil {\n\t\tt.Fatalf(\"first install: %v\\n%s\", err, out)\n\t}\n\ttarget := filepath.Join(gen, \"go\", \"sedum.yaml\")\n\textra := filepath.Join(gen, \"go\", \"extra.txt\")\n\tif err := os.WriteFile(target, []byte(\"changed\\n\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := os.WriteFile(extra, []byte(\"keep\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif out, err := exec(t, \"generators\", \"install\", \"go\", \"--dir\", gen, \"--force\"); err != nil {\n\t\tt.Fatalf(\"force install: %v\\n%s\", err, out)\n\t}\n\tif data, _ := os.ReadFile(extra); string(data) != \"keep\" {\n\t\tt.Error(\"extra file was touched\")\n\t}\n}","file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallForceOverwritesOnlyDifferingFiles"}}
func TestGeneratorsInstallForceOverwritesOnlyDifferingFiles(t *testing.T) {
	gen := filepath.Join(t.TempDir(), "generators")
	if out, err := exec(t, "generators", "install", "go", "--dir", gen); err != nil {
		t.Fatalf("first install: %v\n%s", err, out)
	}
	target := filepath.Join(gen, "go", "sedum.yaml")
	extra := filepath.Join(gen, "go", "extra.txt")
	if err := os.WriteFile(target, []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(extra, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := exec(t, "generators", "install", "go", "--dir", gen, "--force"); err != nil {
		t.Fatalf("force install: %v\n%s", err, out)
	}
	if data, _ := os.ReadFile(extra); string(data) != "keep" {
		t.Error("extra file was touched")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsInstallKeepsExistingGoMod(t *testing.T) {\n\tgen := filepath.Join(t.TempDir(), \"generators\")\n\tif err := os.MkdirAll(gen, 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tgoMod := filepath.Join(gen, \"go.mod\")\n\tif err := os.WriteFile(goMod, []byte(\"module custom\\n\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif out, err := exec(t, \"generators\", \"install\", \"go\", \"--dir\", gen); err != nil {\n\t\tt.Fatalf(\"install: %v\\n%s\", err, out)\n\t}\n\tif data, _ := os.ReadFile(goMod); string(data) != \"module custom\\n\" {\n\t\tt.Error(\"existing go.mod was overwritten\")\n\t}\n}","file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallKeepsExistingGoMod"}}
func TestGeneratorsInstallKeepsExistingGoMod(t *testing.T) {
	gen := filepath.Join(t.TempDir(), "generators")
	if err := os.MkdirAll(gen, 0o755); err != nil {
		t.Fatal(err)
	}
	goMod := filepath.Join(gen, "go.mod")
	if err := os.WriteFile(goMod, []byte("module custom\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := exec(t, "generators", "install", "go", "--dir", gen); err != nil {
		t.Fatalf("install: %v\n%s", err, out)
	}
	if data, _ := os.ReadFile(goMod); string(data) != "module custom\n" {
		t.Error("existing go.mod was overwritten")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsRejectsEscape(t *testing.T) {\n\tparent := t.TempDir()\n\twork := filepath.Join(parent, \"work\")\n\tif err := os.MkdirAll(work, 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := os.WriteFile(filepath.Join(parent, \"canary\"), []byte(\"c\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tbefore := snapshotTree(t, parent)\n\tfor _, name := range []string{\"../evil\", \"a/b\", \"..\", \".\", \"\"} {\n\t\tif _, err := exec(t, \"generators\", \"install\", name, \"--dir\", filepath.Join(work, \"generators\")); err == nil {\n\t\t\tt.Errorf(\"name %q was accepted\", name)\n\t\t}\n\t}\n\tif !reflect.DeepEqual(before, snapshotTree(t, parent)) {\n\t\tt.Error(\"escape attempt changed the parent tree\")\n\t}\n}","file":"internal/cli/generators_test.go","name":"TestGeneratorsRejectsEscape"}}
func TestGeneratorsRejectsEscape(t *testing.T) {
	parent := t.TempDir()
	work := filepath.Join(parent, "work")
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(parent, "canary"), []byte("c"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := snapshotTree(t, parent)
	for _, name := range []string{"../evil", "a/b", "..", ".", ""} {
		if _, err := exec(t, "generators", "install", name, "--dir", filepath.Join(work, "generators")); err == nil {
			t.Errorf("name %q was accepted", name)
		}
	}
	if !reflect.DeepEqual(before, snapshotTree(t, parent)) {
		t.Error("escape attempt changed the parent tree")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsDryRunWritesNothing(t *testing.T) {\n\tgen := filepath.Join(t.TempDir(), \"generators\")\n\tif err := os.MkdirAll(gen, 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tbefore := snapshotTree(t, gen)\n\tout, err := exec(t, \"generators\", \"install\", \"go\", \"--dir\", gen, \"--dry-run\")\n\tif err != nil {\n\t\tt.Fatalf(\"dry run: %v\\n%s\", err, out)\n\t}\n\tif !strings.Contains(out, \"create \") || !strings.Contains(out, \"dry run: nothing written\") {\n\t\tt.Errorf(\"plan not printed:\\n%s\", out)\n\t}\n\tif !reflect.DeepEqual(before, snapshotTree(t, gen)) {\n\t\tt.Error(\"dry run changed the tree\")\n\t}\n}","file":"internal/cli/generators_test.go","name":"TestGeneratorsDryRunWritesNothing"}}
func TestGeneratorsDryRunWritesNothing(t *testing.T) {
	gen := filepath.Join(t.TempDir(), "generators")
	if err := os.MkdirAll(gen, 0o755); err != nil {
		t.Fatal(err)
	}
	before := snapshotTree(t, gen)
	out, err := exec(t, "generators", "install", "go", "--dir", gen, "--dry-run")
	if err != nil {
		t.Fatalf("dry run: %v\n%s", err, out)
	}
	if !strings.Contains(out, "create ") || !strings.Contains(out, "dry run: nothing written") {
		t.Errorf("plan not printed:\n%s", out)
	}
	if !reflect.DeepEqual(before, snapshotTree(t, gen)) {
		t.Error("dry run changed the tree")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsUnknownBuiltinHintsFrom(t *testing.T) {\n\tgen := filepath.Join(t.TempDir(), \"generators\")\n\tout, err := exec(t, \"generators\", \"install\", \"nosuchpkg\", \"--dir\", gen)\n\tif err == nil {\n\t\tt.Fatalf(\"expected error, got:\\n%s\", out)\n\t}\n\tif !strings.Contains(err.Error()+out, \"--from\") {\n\t\tt.Errorf(\"error lacks --from hint: %v\", err)\n\t}\n}","file":"internal/cli/generators_test.go","name":"TestGeneratorsUnknownBuiltinHintsFrom"}}
func TestGeneratorsUnknownBuiltinHintsFrom(t *testing.T) {
	gen := filepath.Join(t.TempDir(), "generators")
	out, err := exec(t, "generators", "install", "nosuchpkg", "--dir", gen)
	if err == nil {
		t.Fatalf("expected error, got:\n%s", out)
	}
	if !strings.Contains(err.Error()+out, "--from") {
		t.Errorf("error lacks --from hint: %v", err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsInstallFromDirectory(t *testing.T) {\n\tsrc := t.TempDir()\n\tnested := filepath.Join(src, \"mypkg\", \"templates\", \"deep\")\n\tif err := os.MkdirAll(nested, 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := os.WriteFile(filepath.Join(src, \"mypkg\", \"sedum.yaml\"), []byte(\"name: mypkg\\n\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := os.WriteFile(filepath.Join(nested, \"t.tmpl\"), []byte(\"x\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tgen := filepath.Join(t.TempDir(), \"generators\")\n\tif out, err := exec(t, \"generators\", \"install\", \"mypkg\", \"--dir\", gen, \"--from\", src); err != nil {\n\t\tt.Fatalf(\"install: %v\\n%s\", err, out)\n\t}\n\tif _, err := os.Stat(filepath.Join(gen, \"mypkg\", \"templates\", \"deep\", \"t.tmpl\")); err != nil {\n\t\tt.Errorf(\"nested file not copied: %v\", err)\n\t}\n}","file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallFromDirectory"}}
func TestGeneratorsInstallFromDirectory(t *testing.T) {
	src := t.TempDir()
	nested := filepath.Join(src, "mypkg", "templates", "deep")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "mypkg", "sedum.yaml"), []byte("name: mypkg\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "t.tmpl"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	gen := filepath.Join(t.TempDir(), "generators")
	if out, err := exec(t, "generators", "install", "mypkg", "--dir", gen, "--from", src); err != nil {
		t.Fatalf("install: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(gen, "mypkg", "templates", "deep", "t.tmpl")); err != nil {
		t.Errorf("nested file not copied: %v", err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsInstallFromMissingPackage(t *testing.T) {\n\tsrc := t.TempDir()\n\tgen := filepath.Join(t.TempDir(), \"generators\")\n\tif _, err := exec(t, \"generators\", \"install\", \"absent\", \"--dir\", gen, \"--from\", src); err == nil {\n\t\tt.Fatal(\"expected error for missing package directory\")\n\t}\n}","file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallFromMissingPackage"}}
func TestGeneratorsInstallFromMissingPackage(t *testing.T) {
	src := t.TempDir()
	gen := filepath.Join(t.TempDir(), "generators")
	if _, err := exec(t, "generators", "install", "absent", "--dir", gen, "--from", src); err == nil {
		t.Fatal("expected error for missing package directory")
	}
}

// /sedum:setFunction
