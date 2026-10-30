package genpkg

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/gomod_test.go","path":"os"},"kwargs":"sha256:e15650787027cac7c83179c3190e01bdb61dbbe639cfbb3af2a318cc0b32670f","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/gomod_test.go","path":"path/filepath"},"kwargs":"sha256:9ba986960e0ad133a20a1f0d87ac85567a315380888e8719fd04333fbbd5613f","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/gomod_test.go","path":"testing"},"kwargs":"sha256:d8c6088e080819df91039e9a195d3b974e6b945c6688380cbc44f19227a5e4d0","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/gomod_test.go","name":"writeGenFile"},"kwargs":"sha256:c2920c8c2f002c8028f096c69c3a63338251617425ce6317875c05b86b9811f3","body":"sha256:ae1fc01954f8763e502a649276801e0ee9604264b07f21b0ed0b1bcec36bdab3"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundaryWarnsWithoutGoMod"},"kwargs":"sha256:e719d6aace355e156ec6b133200fed1e9c57f9591bb374a0223f6c1a5015e01d","body":"sha256:9a06f06882d11e6a0f2014bdf6ca5913d7cc92ad30f6392677e339fd0d77bd57"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundaryQuietWithGoMod"},"kwargs":"sha256:8332d69c0f7aa51fe50102007173398801bf3fb5931641751a7ae4a40b625db4","body":"sha256:56b9dbc53a22f4e725c2c56fd66c3c39bbeddd6ff651aa3fd45165dc0bd2b7e6"}
func TestGoModBoundaryQuietWithGoMod(t *testing.T) {
	dir := t.TempDir()
	writeGenFile(t, dir, "main.go")
	writeGenFile(t, dir, "go.mod")
	if got := checkGoModBoundary(dir); got != nil {
		t.Fatalf("got %+v, want nil", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundaryIgnoresNonGoFiles"},"kwargs":"sha256:93c5a06382e2d17d5fc6df79b8982b97ee4145eca56efc26c90477f2c732df9e","body":"sha256:e280ca0205ce41ab5b695c19e0eea2398a3ef6243ed522c91b595609b17e4bdc"}
func TestGoModBoundaryIgnoresNonGoFiles(t *testing.T) {
	dir := t.TempDir()
	writeGenFile(t, dir, "install.yaml")
	writeGenFile(t, dir, "tpl/render.rb")
	if got := checkGoModBoundary(dir); got != nil {
		t.Fatalf("got %+v, want nil", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundaryFindsNestedGoFile"},"kwargs":"sha256:101925158599089e00696e8d5d7fcae7268e6d37b8bd7bc7c4eee4f58fcc171b","body":"sha256:92cf160222a30852b7c02b515fabbad4cb8d6a871f9a17f80b0f04bc689ce8d8"}
func TestGoModBoundaryFindsNestedGoFile(t *testing.T) {
	dir := t.TempDir()
	writeGenFile(t, dir, "a/b/c/d/deep.go")
	got := checkGoModBoundary(dir)
	if len(got) != 1 || got[0].Rule != RuleGoModMissing {
		t.Fatalf("got %+v, want one %s finding", got, RuleGoModMissing)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundarySubdirGoModStillWarns"},"kwargs":"sha256:cbb205988fc8b96295364cd2e5e71b13c5b906a866d2b5353a919fac86221817","body":"sha256:ce3543ca56d9f7f7109866551b760ef4b485508820d56e56f919ad91c1f8b3c8"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/gomod_test.go","name":"TestGoModBoundaryMissingDir"},"kwargs":"sha256:3ec44218c517f3c5f2e92f106d95341a6bdd8b0cc349fb213d469740814cc8a3","body":"sha256:f3c5d9bbd16d043ac388160a4c74ee91a8666c03bb25b1376e8fe7fbd349b792"}
func TestGoModBoundaryMissingDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "absent")
	if got := checkGoModBoundary(dir); got != nil {
		t.Fatalf("got %+v, want nil", got)
	}
}

// /sedum:setFunction
