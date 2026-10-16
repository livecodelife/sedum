package builtin

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin_test.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin_test.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestNames(t *testing.T) {\n\tgot := Names()\n\tif len(got) != 1 || got[0] != \"go\" {\n\t\tt.Fatalf(\"Names() = %q, want [go]\", got)\n\t}\n}","file":"internal/builtin/builtin_test.go","name":"TestNames"}}
func TestNames(t *testing.T) {
	got := Names()
	if len(got) != 1 || got[0] != "go" {
		t.Fatalf("Names() = %q, want [go]", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestFilesGo(t *testing.T) {\n\tfiles, err := Files(\"go\")\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfor _, want := range []string{\"sedum.yaml\", \"actions/actions.yaml\"} {\n\t\tif _, ok := files[want]; !ok {\n\t\t\tt.Errorf(\"Files(\\\"go\\\") missing %s\", want)\n\t\t}\n\t}\n\tfor k := range files {\n\t\tif strings.HasPrefix(k, \"/\") || strings.Contains(k, \"..\") {\n\t\t\tt.Errorf(\"Files(\\\"go\\\") returned unsafe path %q\", k)\n\t\t}\n\t}\n}","file":"internal/builtin/builtin_test.go","name":"TestFilesGo"}}
func TestFilesGo(t *testing.T) {
	files, err := Files("go")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"sedum.yaml", "actions/actions.yaml"} {
		if _, ok := files[want]; !ok {
			t.Errorf("Files(\"go\") missing %s", want)
		}
	}
	for k := range files {
		if strings.HasPrefix(k, "/") || strings.Contains(k, "..") {
			t.Errorf("Files(\"go\") returned unsafe path %q", k)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestFilesUnknown(t *testing.T) {\n\t_, err := Files(\"x\")\n\tif err == nil {\n\t\tt.Fatal(\"Files(\\\"x\\\") returned no error\")\n\t}\n\tif !strings.Contains(err.Error(), \"go\") {\n\t\tt.Errorf(\"error %q does not name the built-in package go\", err)\n\t}\n}","file":"internal/builtin/builtin_test.go","name":"TestFilesUnknown"}}
func TestFilesUnknown(t *testing.T) {
	_, err := Files("x")
	if err == nil {
		t.Fatal("Files(\"x\") returned no error")
	}
	if !strings.Contains(err.Error(), "go") {
		t.Errorf("error %q does not name the built-in package go", err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestPathsSorted(t *testing.T) {\n\tgot := Paths(map[string][]byte{\"b\": nil, \"a\": nil, \"c/d\": nil})\n\tif !sort.StringsAreSorted(got) {\n\t\tt.Errorf(\"Paths() = %q, not sorted\", got)\n\t}\n\twant := []string{\"a\", \"b\", \"c/d\"}\n\tif len(got) != len(want) {\n\t\tt.Fatalf(\"Paths() = %q, want %q\", got, want)\n\t}\n\tfor i := range want {\n\t\tif got[i] != want[i] {\n\t\t\tt.Fatalf(\"Paths() = %q, want %q\", got, want)\n\t\t}\n\t}\n}","file":"internal/builtin/builtin_test.go","name":"TestPathsSorted"}}
func TestPathsSorted(t *testing.T) {
	got := Paths(map[string][]byte{"b": nil, "a": nil, "c/d": nil})
	if !sort.StringsAreSorted(got) {
		t.Errorf("Paths() = %q, not sorted", got)
	}
	want := []string{"a", "b", "c/d"}
	if len(got) != len(want) {
		t.Fatalf("Paths() = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Paths() = %q, want %q", got, want)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGeneratorsDrift(t *testing.T) {\n\troot := filepath.Join(\"..\", \"..\", \"generators\", \"go\")\n\tif _, err := os.Stat(root); os.IsNotExist(err) {\n\t\tt.Skip(\"generators/go does not exist\")\n\t}\n\tfiles, err := Files(\"go\")\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tsrc := make(map[string][]byte)\n\terr = filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tif d.IsDir() {\n\t\t\treturn nil\n\t\t}\n\t\trel, err := filepath.Rel(root, p)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tdata, err := os.ReadFile(p)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tsrc[filepath.ToSlash(rel)] = data\n\t\treturn nil\n\t})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tfor k, want := range src {\n\t\tgot, ok := files[k]\n\t\tif !ok || !bytes.Equal(got, want) {\n\t\t\tt.Errorf(\"generators/go/%s differs from built-in package; run hack/sync-builtin.sh\", k)\n\t\t}\n\t}\n\tfor k := range files {\n\t\tif _, ok := src[k]; !ok {\n\t\t\tt.Errorf(\"built-in package has %s missing from generators/go; run hack/sync-builtin.sh\", k)\n\t\t}\n\t}\n}","file":"internal/builtin/builtin_test.go","name":"TestGeneratorsDrift"}}
func TestGeneratorsDrift(t *testing.T) {
	root := filepath.Join("..", "..", "generators", "go")
	if _, err := os.Stat(root); os.IsNotExist(err) {
		t.Skip("generators/go does not exist")
	}
	files, err := Files("go")
	if err != nil {
		t.Fatal(err)
	}
	src := make(map[string][]byte)
	err = filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		src[filepath.ToSlash(rel)] = data
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for k, want := range src {
		got, ok := files[k]
		if !ok || !bytes.Equal(got, want) {
			t.Errorf("generators/go/%s differs from built-in package; run hack/sync-builtin.sh", k)
		}
	}
	for k := range files {
		if _, ok := src[k]; !ok {
			t.Errorf("built-in package has %s missing from generators/go; run hack/sync-builtin.sh", k)
		}
	}
}

// /sedum:setFunction
