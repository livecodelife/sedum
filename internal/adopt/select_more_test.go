package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_more_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_more_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_more_test.go","path":"reflect"}}
import "reflect"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_more_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_more_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select_more_test.go","path":"github.com/livecodelife/sedum/internal/record"}}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSelectWalksDirectories(t *testing.T) {\n\troot, set, _ := newSelectFixture(t)\n\tfor _, f := range []string{\"sub/b.py\", \"sub/a.py\", \"sub/.hidden/c.py\"} {\n\t\tp := filepath.Join(root, filepath.FromSlash(f))\n\t\tif err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t\tif err := os.WriteFile(p, []byte(\"x\"), 0o644); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\trecords := &record.Set{Records: []*record.Record{{ID: \"imp-1\", Type: \"imprint\", Patterns: []string{\"**/*.py\"}}}}\n\tcands, _, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{\"sub\", \"sub\"})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tvar got []string\n\tfor _, c := range cands {\n\t\tgot = append(got, c.Path)\n\t}\n\tif want := []string{\"sub/a.py\", \"sub/b.py\"}; !reflect.DeepEqual(got, want) {\n\t\tt.Fatalf(\"got %v, want %v\", got, want)\n\t}\n}\n","file":"internal/adopt/select_more_test.go","name":"TestSelectWalksDirectories"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSelectRejectsEscapingPath(t *testing.T) {\n\troot, set, _ := newSelectFixture(t)\n\trecords := &record.Set{Records: []*record.Record{{ID: \"imp-1\", Type: \"imprint\", Patterns: []string{\"**/*.py\"}}}}\n\t_, _, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{\"../x.py\"})\n\tif err == nil {\n\t\tt.Fatal(\"want error for path escaping root\")\n\t}\n}\n","file":"internal/adopt/select_more_test.go","name":"TestSelectRejectsEscapingPath"}}
func TestSelectRejectsEscapingPath(t *testing.T) {
	root, set, _ := newSelectFixture(t)
	records := &record.Set{Records: []*record.Record{{ID: "imp-1", Type: "imprint", Patterns: []string{"**/*.py"}}}}
	_, _, err := Select(SelectConfig{Root: root, Set: set, Records: records}, []string{"../x.py"})
	if err == nil {
		t.Fatal("want error for path escaping root")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSelectTestsFlag(t *testing.T) {\n\troot, set, _ := newSelectFixture(t)\n\tif err := os.WriteFile(filepath.Join(root, \"test_a.py\"), []byte(\"x\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\trecords := &record.Set{Records: []*record.Record{{ID: \"imp-1\", Type: \"imprint\", Patterns: []string{\"**/*.py\"}}}}\n\tcands, _, err := Select(SelectConfig{Root: root, Set: set, Records: records, Tests: true}, []string{\"test_a.py\"})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(cands) != 1 || cands[0].Path != \"test_a.py\" {\n\t\tt.Fatalf(\"got %+v, want one candidate test_a.py\", cands)\n\t}\n}\n","file":"internal/adopt/select_more_test.go","name":"TestSelectTestsFlag"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-b82a46e4","gen":"sha256:b27dee73481760e31a54dfc3618318671f30bb7d73a15296292a9efe33a7c99c","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSelectNoImprintAndAmbiguous(t *testing.T) {\n\troot, set, _ := newSelectFixture(t)\n\tif err := os.WriteFile(filepath.Join(root, \"a.py\"), []byte(\"x\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tambiguous := &record.Set{Records: []*record.Record{{ID: \"imp-1\", Type: \"imprint\", Patterns: []string{\"**/*.py\"}}, {ID: \"imp-2\", Type: \"imprint\", Patterns: []string{\"**/*.py\"}}}}\n\t_, skips, err := Select(SelectConfig{Root: root, Set: set, Records: ambiguous}, []string{\"a.py\"})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(skips) != 1 || skips[0].Reason != ReasonAmbiguousImprint {\n\t\tt.Fatalf(\"got skips %+v, want one ReasonAmbiguousImprint\", skips)\n\t}\n\tfor _, id := range []string{\"imp-1\", \"imp-2\"} {\n\t\tif !strings.Contains(skips[0].Detail, id) {\n\t\t\tt.Errorf(\"Detail %q does not mention %s\", skips[0].Detail, id)\n\t\t}\n\t}\n\tnone := &record.Set{}\n\t_, skips, err = Select(SelectConfig{Root: root, Set: set, Records: none}, []string{\"a.py\"})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(skips) != 1 || skips[0].Reason != ReasonNoImprint {\n\t\tt.Fatalf(\"got skips %+v, want one ReasonNoImprint\", skips)\n\t}\n}\n","file":"internal/adopt/select_more_test.go","name":"TestSelectNoImprintAndAmbiguous"}}
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
