package builtin

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin_test.go","path":"bytes"},"kwargs":"sha256:af12bd3ff7cc333ab5b3a4ceb777cb41e81e076f4e27d905a526bd818558ac7a","body":"sha256:468731f6dd513df0a66690d4649dc9e7244c96de227d8e6ead98cc5e9747c57a"}
import "bytes"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin_test.go","path":"os"},"kwargs":"sha256:f112c7c9f6cd9c6fe684d527fa5f232e59a93deba96175c2efa13719a98ab194","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin_test.go","path":"path/filepath"},"kwargs":"sha256:1d36510eba9db1e6acd843fb7672fa4d23150f5f4179b8e45a8804fc955fd011","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin_test.go","path":"sort"},"kwargs":"sha256:52d5b862464688104055a07e1c1c7e0f07624d7c9794de8c6654a6132ad8a4fc","body":"sha256:f197172617ae21ef23b1855b09aafb14571020ae81bff877d63d3ff81f4bf63c"}
import "sort"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin_test.go","path":"strings"},"kwargs":"sha256:74f9e2f23aef3e795b7c02d1a44904e16b9f8b9229b211d45a47959722f9f4cd","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin_test.go","path":"testing"},"kwargs":"sha256:0ee2fd2edf9e17772d12f1ff422202bcb56ecec16d036dcddd5bad31a8b5710b","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/builtin/builtin_test.go","name":"TestNames"},"kwargs":"sha256:0ced5a0649047eb2c6f276cf90507d65124da5cf6fd197ab0b1ef5f4a38eb322","body":"sha256:1fcfb3510df84d70fb1ff06c3fe96cd9fea5fdfd80a3d51201f46ae1379d0cfd"}
func TestNames(t *testing.T) {
	got := Names()
	if len(got) != 1 || got[0] != "go" {
		t.Fatalf("Names() = %q, want [go]", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/builtin/builtin_test.go","name":"TestFilesGo"},"kwargs":"sha256:a701071028e4b2cd4fc12d59eacadcdcf18f7c9ba8cc5a50f5bf35cf331f1ead","body":"sha256:3778b0f4c9256edc0bfa6089a33a9e030424b353684391b4fdd60560ddd90ff9"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/builtin/builtin_test.go","name":"TestFilesUnknown"},"kwargs":"sha256:147582918d015af263c2a2eccd18858abf7917e76044a58e1ab04f2574bd6d01","body":"sha256:d52a4b6dd6b9ab662535b0123ab7380594e50dffeb2abc8888eefb6baa4d0907"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/builtin/builtin_test.go","name":"TestPathsSorted"},"kwargs":"sha256:0f48f4337c0a73af84cdb33c6a02cecc860b537a2068905343fa0756a747ab9d","body":"sha256:ed6b8cb8f73aac4a1a03822cac1e65d12b9e390fb2cc64560cfb32e2d3a41935"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/builtin/builtin_test.go","name":"TestGeneratorsDrift"},"kwargs":"sha256:9df199d3ab110e2512a6f11672a5d2feed9851279b615d7a30bea164a40f8840","body":"sha256:d13112a4ce2ee91f44042e0416a784e0c7e4bd7cc80873f11f14bb4550dc4453"}
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
