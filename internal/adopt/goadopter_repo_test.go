package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter_repo_test.go","path":"io/fs"}}
import "io/fs"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter_repo_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter_repo_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter_repo_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter_repo_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoAdopterCoversRepoFiles(t *testing.T) {\n\tn := 0\n\terr := filepath.WalkDir(\"../../internal\", func(p string, d fs.DirEntry, err error) error {\n\t\tif err != nil || d.IsDir() || !strings.HasSuffix(p, \".go\") {\n\t\t\treturn err\n\t\t}\n\t\tsrc, err := os.ReadFile(p)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tif strings.Contains(string(src), \"sedum:\") {\n\t\t\treturn nil\n\t\t}\n\t\tspans, err := GoAdopter{}.Split(p, src)\n\t\tif err != nil {\n\t\t\treturn nil\n\t\t}\n\t\tif err := CheckSpans(src, spans); err != nil {\n\t\t\tt.Errorf(\"%s: %v\", p, err)\n\t\t}\n\t\tn++\n\t\treturn nil\n\t})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif n == 0 {\n\t\tt.Fatal(\"no Go files were split\")\n\t}\n}","file":"internal/adopt/goadopter_repo_test.go","name":"TestGoAdopterCoversRepoFiles"}}
func TestGoAdopterCoversRepoFiles(t *testing.T) {
	n := 0
	err := filepath.WalkDir("../../internal", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".go") {
			return err
		}
		src, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if strings.Contains(string(src), "sedum:") {
			return nil
		}
		spans, err := GoAdopter{}.Split(p, src)
		if err != nil {
			return nil
		}
		if err := CheckSpans(src, spans); err != nil {
			t.Errorf("%s: %v", p, err)
		}
		n++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("no Go files were split")
	}
}

// /sedum:setFunction
