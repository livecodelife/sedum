package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_more_test.go","path":"os"},"kwargs":"sha256:e4c079363d3ee7d94500637353c348f07947efe816e448ff01a4f266c71c1d03","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_more_test.go","path":"path/filepath"},"kwargs":"sha256:fd6ed702268f2791a3cf915ad5561c15258b6a431670e82318e9dfb5e666b7c0","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_more_test.go","path":"reflect"},"kwargs":"sha256:1110a04f5fca374bcce9961d13353288e6fd66bdcd1d6da5410c8e3f96dd6674","body":"sha256:893cf368d31e1d8fe04189a0aba8627e8e69b39eafe10bea46d39a275e986fa2"}
import "reflect"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_more_test.go","path":"strings"},"kwargs":"sha256:ce0a85df238b9e0b41b856f989ef364da94abbe067f5152ab31890f756f5fc07","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_more_test.go","path":"testing"},"kwargs":"sha256:173d4c00c8fbed2f9cb0436d6c3e515afd022a5d4d8fe44db234bb249278ec0c","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/select_more_test.go","path":"github.com/livecodelife/sedum/internal/record"},"kwargs":"sha256:6d22e29dbff8e830eb0333753ef21b59f28afcf6b1dd1c95c2fde8a893f0ad0e","body":"sha256:527c0befadfcad873cb5f9da4ab560755bc3040424ba6241011c3abf87fcc5e0"}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/select_more_test.go","name":"TestSelectWalksDirectories"},"kwargs":"sha256:5df998e4d94adae00e7ac8183ba34415262f31c96738d10465a9aa8d87d93b61","body":"sha256:ccac508b767c1abb16319bcab4667ca350179784d6d1c172049daea3d65a9763"}
func TestSelectWalksDirectories(t *testing.T) {
	root, set, _ := newSelectFixture(t)
	for _, f := range []string{"sub/b.py", "sub/a.py", "sub/.hidden/c.py"} {
		p := filepath.Join(root, filepath.FromSlash(f))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	records := &record.Set{Records: []*record.Record{{ID: "imp-1", Type: "imprint", Patterns: []string{"**/*.py"}}}}
	cands, _, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{"sub", "sub"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, c := range cands {
		got = append(got, c.Path)
	}
	if want := []string{"sub/a.py", "sub/b.py"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/select_more_test.go","name":"TestSelectRejectsEscapingPath"},"kwargs":"sha256:00ca6944dc3e9c89587719feaaf2b3dc0afa3409fe56f5b967a7fe9c770dff3b","body":"sha256:4a7a51e75917ec091a64a83de466820f6dd8d0bb7b6ed60e14eef7a2703166f1"}
func TestSelectRejectsEscapingPath(t *testing.T) {
	root, set, _ := newSelectFixture(t)
	records := &record.Set{Records: []*record.Record{{ID: "imp-1", Type: "imprint", Patterns: []string{"**/*.py"}}}}
	_, _, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{"../x.py"})
	if err == nil {
		t.Fatal("want error for path escaping root")
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/select_more_test.go","name":"TestSelectTestsFlag"},"kwargs":"sha256:53c762d755ec97acd76405c2f31d9e9de37cfadd8b2f7d1da6b0bad7632903a7","body":"sha256:37145afd4d5a0a21aef47d796d83872a5f2e32a2b3e7153a716bff529a7f9091"}
func TestSelectTestsFlag(t *testing.T) {
	root, set, _ := newSelectFixture(t)
	if err := os.WriteFile(filepath.Join(root, "test_a.py"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	records := &record.Set{Records: []*record.Record{{ID: "imp-1", Type: "imprint", Patterns: []string{"**/*.py"}}}}
	cands, _, err := Select(SelectConfig{Root: root, Set: set, Records: records, Tests: true}, []string{"test_a.py"})
	if err != nil {
		t.Fatal(err)
	}
	if len(cands) != 1 || cands[0].Path != "test_a.py" {
		t.Fatalf("got %+v, want one candidate test_a.py", cands)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/select_more_test.go","name":"TestSelectNoImprintAndAmbiguous"},"kwargs":"sha256:fd6b3ff3311bcfa9992db79b14f6459e18623c3575c3d5c0de6c65e4e08cae46","body":"sha256:e6112bee259787ce5cc424372ba8434b63544a0ac9138e43eb70280f24a1533e"}
func TestSelectNoImprintAndAmbiguous(t *testing.T) {
	root, set, _ := newSelectFixture(t)
	if err := os.WriteFile(filepath.Join(root, "a.py"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	ambiguous := &record.Set{Records: []*record.Record{{ID: "imp-1", Type: "imprint", Patterns: []string{"**/*.py"}}, {ID: "imp-2", Type: "imprint", Patterns: []string{"**/*.py"}}}}
	_, skips, err := Select(SelectConfig{Root: root, Set: set, Records: ambiguous}, []string{"a.py"})
	if err != nil {
		t.Fatal(err)
	}
	if len(skips) != 1 || skips[0].Reason != ReasonAmbiguousImprint {
		t.Fatalf("got skips %+v, want one ReasonAmbiguousImprint", skips)
	}
	for _, id := range []string{"imp-1", "imp-2"} {
		if !strings.Contains(skips[0].Detail, id) {
			t.Errorf("Detail %q does not mention %s", skips[0].Detail, id)
		}
	}
	none := &record.Set{}
	_, skips, err = Select(SelectConfig{Root: root, Set: set, Records: none}, []string{"a.py"})
	if err != nil {
		t.Fatal(err)
	}
	if len(skips) != 1 || skips[0].Reason != ReasonNoImprint {
		t.Fatalf("got skips %+v, want one ReasonNoImprint", skips)
	}
}

// /sedum:setFunction
