package genpkg

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/gomod_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/gomod_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/gomod_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// writeGenFile creates rel under dir with a placeholder body.\nfunc writeGenFile(t *testing.T, dir, rel string) {\n\tt.Helper()\n\tp := filepath.Join(dir, rel)\n\tif err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := os.WriteFile(p, []byte(\"package x\\n\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n}","file":"internal/genpkg/gomod_test.go","name":"writeGenFile"}}
// writeGenFile creates rel under dir with a placeholder body.
func writeGenFile(t *testing.T, dir, rel string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("package x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoModBoundaryWarnsWithoutGoMod(t *testing.T) {\n\tdir := t.TempDir()\n\twriteGenFile(t, dir, \"tpl/main.go\")\n\tgot := checkGoModBoundary(dir)\n\tif len(got) != 1 {\n\t\tt.Fatalf(\"got %d findings, want 1\", len(got))\n\t}\n\tif got[0].Rule != RuleGoModMissing || got[0].Kind != KindWarning {\n\t\tt.Fatalf(\"unexpected finding: %+v\", got[0])\n\t}\n\twant := \"holds Go template files but has no go.mod, so go build/test ./... will try to compile them; run `sedum generators install` or add a go.mod with `module generators`\"\n\tif got[0].Message != want {\n\t\tt.Fatalf(\"message = %q, want %q\", got[0].Message, want)\n\t}\n}","file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundaryWarnsWithoutGoMod"}}
func TestGoModBoundaryWarnsWithoutGoMod(t *testing.T) {
	dir := t.TempDir()
	writeGenFile(t, dir, "tpl/main.go")
	got := checkGoModBoundary(dir)
	if len(got) != 1 {
		t.Fatalf("got %d findings, want 1", len(got))
	}
	if got[0].Rule != RuleGoModMissing || got[0].Kind != KindWarning {
		t.Fatalf("unexpected finding: %+v", got[0])
	}
	want := "holds Go template files but has no go.mod, so go build/test ./... will try to compile them; run `sedum generators install` or add a go.mod with `module generators`"
	if got[0].Message != want {
		t.Fatalf("message = %q, want %q", got[0].Message, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoModBoundaryQuietWithGoMod(t *testing.T) {\n\tdir := t.TempDir()\n\twriteGenFile(t, dir, \"main.go\")\n\twriteGenFile(t, dir, \"go.mod\")\n\tif got := checkGoModBoundary(dir); got != nil {\n\t\tt.Fatalf(\"got %+v, want nil\", got)\n\t}\n}","file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundaryQuietWithGoMod"}}
func TestGoModBoundaryQuietWithGoMod(t *testing.T) {
	dir := t.TempDir()
	writeGenFile(t, dir, "main.go")
	writeGenFile(t, dir, "go.mod")
	if got := checkGoModBoundary(dir); got != nil {
		t.Fatalf("got %+v, want nil", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoModBoundaryIgnoresNonGoFiles(t *testing.T) {\n\tdir := t.TempDir()\n\twriteGenFile(t, dir, \"install.yaml\")\n\twriteGenFile(t, dir, \"tpl/render.rb\")\n\tif got := checkGoModBoundary(dir); got != nil {\n\t\tt.Fatalf(\"got %+v, want nil\", got)\n\t}\n}","file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundaryIgnoresNonGoFiles"}}
func TestGoModBoundaryIgnoresNonGoFiles(t *testing.T) {
	dir := t.TempDir()
	writeGenFile(t, dir, "install.yaml")
	writeGenFile(t, dir, "tpl/render.rb")
	if got := checkGoModBoundary(dir); got != nil {
		t.Fatalf("got %+v, want nil", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoModBoundaryFindsNestedGoFile(t *testing.T) {\n\tdir := t.TempDir()\n\twriteGenFile(t, dir, \"a/b/c/d/deep.go\")\n\tgot := checkGoModBoundary(dir)\n\tif len(got) != 1 || got[0].Rule != RuleGoModMissing {\n\t\tt.Fatalf(\"got %+v, want one %s finding\", got, RuleGoModMissing)\n\t}\n}","file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundaryFindsNestedGoFile"}}
func TestGoModBoundaryFindsNestedGoFile(t *testing.T) {
	dir := t.TempDir()
	writeGenFile(t, dir, "a/b/c/d/deep.go")
	got := checkGoModBoundary(dir)
	if len(got) != 1 || got[0].Rule != RuleGoModMissing {
		t.Fatalf("got %+v, want one %s finding", got, RuleGoModMissing)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoModBoundarySubdirGoModStillWarns(t *testing.T) {\n\tdir := t.TempDir()\n\twriteGenFile(t, dir, \"tpl/main.go\")\n\twriteGenFile(t, dir, \"sub/go.mod\")\n\tgot := checkGoModBoundary(dir)\n\tif len(got) != 1 || got[0].Rule != RuleGoModMissing {\n\t\tt.Fatalf(\"got %+v, want one %s finding\", got, RuleGoModMissing)\n\t}\n}","file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundarySubdirGoModStillWarns"}}
func TestGoModBoundarySubdirGoModStillWarns(t *testing.T) {
	dir := t.TempDir()
	writeGenFile(t, dir, "tpl/main.go")
	writeGenFile(t, dir, "sub/go.mod")
	got := checkGoModBoundary(dir)
	if len(got) != 1 || got[0].Rule != RuleGoModMissing {
		t.Fatalf("got %+v, want one %s finding", got, RuleGoModMissing)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoModBoundaryMissingDir(t *testing.T) {\n\tdir := filepath.Join(t.TempDir(), \"absent\")\n\tif got := checkGoModBoundary(dir); got != nil {\n\t\tt.Fatalf(\"got %+v, want nil\", got)\n\t}\n}","file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundaryMissingDir"}}
func TestGoModBoundaryMissingDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "absent")
	if got := checkGoModBoundary(dir); got != nil {
		t.Fatalf("got %+v, want nil", got)
	}
}

// /sedum:setFunction
