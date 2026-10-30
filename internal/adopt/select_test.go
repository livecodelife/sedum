package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_test.go","path":"bytes"},"kwargs":"sha256:652a793ce7d4f991af1a91c575a202e58ababc75495f239eb960cf44bd322d8e","body":"sha256:468731f6dd513df0a66690d4649dc9e7244c96de227d8e6ead98cc5e9747c57a"}
import "bytes"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_test.go","path":"os"},"kwargs":"sha256:e0f323099f9c0258b9260d62aafabf836324d447712e9501a80e3d75f5480ee2","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_test.go","path":"path/filepath"},"kwargs":"sha256:956a9cb7c3e161abfd18ac85cb35087e542b8e81bfa257ad3da321dd4e107402","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_test.go","path":"testing"},"kwargs":"sha256:299224db639fcd3a10e4394c8de20f8bc422cff3ed77ed6b412493e3967ac677","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:ac80b9ca621871b015bf4c15c7baf2d59287bc5642dce23ca5c76716b28fd5fd","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_test.go","path":"github.com/livecodelife/sedum/internal/record"},"kwargs":"sha256:156ff325f0c5587fb5eefa5d2adbfda77570cb6d1b5051be2130759bb2d99bac","body":"sha256:527c0befadfcad873cb5f9da4ab560755bc3040424ba6241011c3abf87fcc5e0"}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/select_test.go","name":"newSelectFixture"},"kwargs":"sha256:af9eaa4daf44152653a18db0c6caf08b565d68645069513af732bbe1b5b8d868","body":"sha256:d5953683445357da2ba1b5f83a4cacf00c030ce7dcf965e1910fe1047d5708ed"}
// newSelectFixture writes a small generators tree and returns a root directory,
// a loaded generator set, and an imprint record set.
func newSelectFixture(t *testing.T) (string, *genpkg.Set, *record.Set) {
	t.Helper()
	gen := t.TempDir()
	write := func(rel, body string) {
		t.Helper()
		p := filepath.Join(gen, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("py/sedum.yaml", `name: py
extensions: [".py"]
comment_prefix: "#"
test_paths: ["test_*.py"]
unmanaged: ["setup.py"]
`)
	write("py/files/{name}.py", "# sedum:anchor:decls\n")
	write("py/files/*/{name}.py", "# sedum:anchor:decls\n")
	write("py/actions/actions.yaml", `actions:
  setThing:
    kwargs:
      file: { type: string, required: true }
      name: { type: string, required: true }
      code: { type: string, required: true }
    injects_into: "{{file}}"
    anchor: decls
`)
	write("py/actions/setThing.py", "{{code}}\n")
	set, findings, err := genpkg.Load(gen, genpkg.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) > 0 {
		t.Fatalf("generator findings: %v", findings)
	}
	root := t.TempDir()
	records := &record.Set{Records: []*record.Record{{ID: "imp-1", Type: "imprint", Patterns: []string{"*.py"}}}}
	return root, set, records
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/select_test.go","name":"TestSelectCandidate"},"kwargs":"sha256:25146d3467681863ba2bf7ffad3fe2af9dd03a176fb8caac0218e369e91b51e6","body":"sha256:7c47d0a8d82b08e870ff939cdc6623c131a8c2c154b4882421b1949a365595b2"}
func TestSelectCandidate(t *testing.T) {
	root, set, records := newSelectFixture(t)
	src := []byte(`x = 1
`)
	if err := os.WriteFile(filepath.Join(root, "a.py"), src, 0o644); err != nil {
		t.Fatal(err)
	}
	cands, skips, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{"a.py"})
	if err != nil {
		t.Fatal(err)
	}
	if len(skips) != 0 {
		t.Fatalf("skips = %+v, want none", skips)
	}
	if len(cands) != 1 {
		t.Fatalf("got %d candidates, want 1", len(cands))
	}
	c := cands[0]
	if c.Template != "{name}.py" {
		t.Errorf("Template = %q, want %q", c.Template, "{name}.py")
	}
	if c.Record == nil || c.Record.ID != "imp-1" {
		t.Errorf("Record = %+v, want ID imp-1", c.Record)
	}
	if c.Captures["name"] != "a" {
		t.Errorf("Captures[name] = %q, want %q", c.Captures["name"], "a")
	}
	if !bytes.Equal(c.Src, src) {
		t.Errorf("Src = %q, want %q", c.Src, src)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-37216199","gen":"sha256:49fad05c79349b33a272d47e3cbeab53d67715030ab9b097f5d518bfa6da09db","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/select_test.go","name":"TestSelectSkipReasons"},"kwargs":"sha256:50659557bdd211e110cef6081f0dd6f7f52d5b7d1c2d2397021a05f8cd6fed60","body":"sha256:4957120871563e4bc1f4c3d1d7dfe83877a67eb626c101e03b1af2b003cb19d3"}
func TestSelectSkipReasons(t *testing.T) {
	cases := []struct {
		name string
		file string
		body string
		want Reason
	}{
		{"unmanaged", "setup.py", "x = 1\n", ReasonUnmanaged},
		{"test", "test_a.py", "x = 1\n", ReasonTest},
		{"generated", "gen.py", "# DO NOT EDIT\nx = 1\n", ReasonGenerated},
		{"already adopted", "done.py", "# sedum:setThing {}\n", ReasonAlreadyAdopted},
		{"no package", "notes.txt", "hello\n", ReasonNoPackage},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, set, records := newSelectFixture(t)
			if err := os.WriteFile(filepath.Join(root, tc.file), []byte(tc.body), 0o644); err != nil {
				t.Fatal(err)
			}
			cands, skips, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{tc.file})
			if err != nil {
				t.Fatal(err)
			}
			if len(cands) != 0 || len(skips) != 1 {
				t.Fatalf("got %d candidates and %d skips, want 0 and 1", len(cands), len(skips))
			}
			if skips[0].Reason != tc.want {
				t.Errorf("Reason = %q, want %q", skips[0].Reason, tc.want)
			}
		})
	}
}

// /sedum:setFunction
