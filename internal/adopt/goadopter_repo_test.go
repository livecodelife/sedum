package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter_repo_test.go","path":"io/fs"},"kwargs":"sha256:af02218fce4e38107974bf87b9786f434b97d8f8f9107010082894a0ded3cfb3","body":"sha256:80f0597eb523a3ace402ac68fefdf64bf00074dda5ec249d04e8f1e525952547"}
import "io/fs"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter_repo_test.go","path":"os"},"kwargs":"sha256:c8bddac512c868a7b8a69167caacd6f46a3a657cab04b5332f420e8a97cd09a2","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter_repo_test.go","path":"path/filepath"},"kwargs":"sha256:83d6d2a1062f2335164784122d3b704e264e878896857f8654b3471e73d25d4f","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter_repo_test.go","path":"strings"},"kwargs":"sha256:e41d082adb278d087072adb61c26e33a0af176afbe1b2999481076b4196f85e7","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter_repo_test.go","path":"testing"},"kwargs":"sha256:5bd39a669e74fc52e2fefb8f29c44e683314411a881508aec6fe3be68d473679","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/goadopter_repo_test.go","name":"TestGoAdopterCoversRepoFiles"},"kwargs":"sha256:1176a8063fe290c142dd3ca20a0fdc3d41944e94761f231b37bd5de455ea652b","body":"sha256:d46af44184b4450fb4a9486152963901fbc7754540e5934305887f6b53a30ad6"}
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
