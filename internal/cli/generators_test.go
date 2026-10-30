package cli

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/generators_test.go","path":"os"},"kwargs":"sha256:a6908c64920cfb4dafcd6204026778564bf49436dea6d83ef3c40542398f5694","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/generators_test.go","path":"path/filepath"},"kwargs":"sha256:583c842c8dd20b286a5fb65686caee1c919389a9931b2ce3ad710159dbe2fe9b","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/generators_test.go","path":"reflect"},"kwargs":"sha256:e5b8e59b0d50a80ad6f97fc501eef84c5820f105a4d8ddcf7a1d9137151358b0","body":"sha256:893cf368d31e1d8fe04189a0aba8627e8e69b39eafe10bea46d39a275e986fa2"}
import "reflect"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/generators_test.go","path":"strings"},"kwargs":"sha256:fa9a866a442c9cb8f90cc4aef4a858bb99ba8401ce4b83e0ad972935f2c9dc9f","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/generators_test.go","path":"testing"},"kwargs":"sha256:6309aff393b38e3a0842d9e77b5628d9b661e8cc1fc03aedaf6b71b91266cecc","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/generators_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:a22cb7cd19ce26d82c070124991700301a6ce9c594148fad4b4a8114a82ad530","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"snapshotTree"},"kwargs":"sha256:f316e24050e69374c2c43f7cf39d8e4776fe11b93868f72a799ea3e802237d65","body":"sha256:60f999f9617fc6af2a5442e26fab6d1d2f3937bb2509b24f5af7064f28af8089"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallBuiltin"},"kwargs":"sha256:1d978f4f6c2b41a61685f543437ba567c34a4732a7a2b9744575b2ec4d2cb2fc","body":"sha256:9b786f41d19bc9771b8fcaf15b212b2a53304352284c27ad5bb869b4ee7f7d9b"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallTwiceIsIdentical"},"kwargs":"sha256:e55f42883c5180eab6362089522c54fc80fe50e06036d5efcc09bd96ff3575e5","body":"sha256:0635bd24c28ad6adf6b7f7f37051c134e7a35286b4deb0be72c36a45b22d69a4"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallRefusesDifferingFile"},"kwargs":"sha256:e7f9e2ebf41944dfd00d6b4da21cb4e69992911626f91a32c34c3064dd91c2f3","body":"sha256:f002068d2e4fe1df1cbd167fe32495dafcf8785462d652f8cd753b1ed9f04eae"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallForceOverwritesOnlyDifferingFiles"},"kwargs":"sha256:37b7221d48bcb96481985f70719dedcbbb3c607fba74e59e754b0fc57f93b15f","body":"sha256:3544f1cb90f1115e4d015743c893ce6d396cfd7689002ce0315004a58c9f8554"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallKeepsExistingGoMod"},"kwargs":"sha256:d2980fe58fbd54b43174a6de2823a99d130e859ea5cdb8e32547a1e107233b94","body":"sha256:cb6d6d93284dbcbf8d82c0cbe87033675b3e7653a8ef8dd09bcff57a9f3cbe2e"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"TestGeneratorsRejectsEscape"},"kwargs":"sha256:f2ae6537064f8dbcad3b3d163a40cebaec496333d723d990d472c21a2c58e284","body":"sha256:bc462e4b1319c7a162ff3660426546ee9ff8c4924ae4316b2f9dca07e4c5cc85"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"TestGeneratorsDryRunWritesNothing"},"kwargs":"sha256:3181b6d5110772461d0b7ac920dc4458939ca6318d6f84e8ba5f3bb6455590a2","body":"sha256:9288d8d1071545870bad8d56da33ff4f66acc19d6be8b420454705721bcdd33c"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"TestGeneratorsUnknownBuiltinHintsFrom"},"kwargs":"sha256:5cb075d6b579211007017dbabf8ff96288679e660fb6be8a7f52e0c50e8b9a90","body":"sha256:b9e011715a9d37ef09ca8e9c56302c44ec3c79bdab4365381b55cb91d18d072d"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallFromDirectory"},"kwargs":"sha256:6468b53452be42b01023c9c97261bb04d73c52fa7efe1f4cf817143eef536f40","body":"sha256:2ff9f0215ac50fa0aa6f0f66d2f8faef0b4cdbfc1a9729fb76806ca9bb234095"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/generators_test.go","name":"TestGeneratorsInstallFromMissingPackage"},"kwargs":"sha256:2cc4b34e4e243d02d8a9136ffbd942c6218f510a88ebdcedf8365f1ce1cb100e","body":"sha256:69f62e3124524f625b5f0dece674edd1d23c949fc07af76e2137d49a28d86fce"}
func TestGeneratorsInstallFromMissingPackage(t *testing.T) {
	src := t.TempDir()
	gen := filepath.Join(t.TempDir(), "generators")
	if _, err := exec(t, "generators", "install", "absent", "--dir", gen, "--from", src); err == nil {
		t.Fatal("expected error for missing package directory")
	}
}

// /sedum:setFunction
