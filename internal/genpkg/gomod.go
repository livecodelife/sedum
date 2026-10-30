package genpkg

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/gomod.go","path":"io/fs"},"kwargs":"sha256:4b2c7d4c698fde82e5c5f06def41847ea55f94c3c5d09ce8a0e93dbd313f553b","body":"sha256:80f0597eb523a3ace402ac68fefdf64bf00074dda5ec249d04e8f1e525952547"}
import "io/fs"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/gomod.go","path":"os"},"kwargs":"sha256:e3625cb90235ab417a59119754690cd18124d825f7b858883a2dfe97d55b695f","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/gomod.go","path":"path/filepath"},"kwargs":"sha256:1c62db65b1c247f452d4eb5d0363447df5ede61117d7be6692b14dc07d0dc37c","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/gomod.go","path":"strings"},"kwargs":"sha256:dfe81107b051474524ab7255027ffb16f2bf78fab081987b28da739cfb2221e8","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/gomod.go","name":"checkGoModBoundary"},"kwargs":"sha256:505a0dd2ae71a653b758a93448dcccd50b109592c48f025502461947500ed8d5","body":"sha256:fd4da9f827bd81a8770936a03c8afb14ff0f49b865f367f9ab75943061bb0967"}
// checkGoModBoundary reports a warning when dir holds Go template files
// anywhere below it but has no go.mod at its root.
func checkGoModBoundary(dir string) Findings {
	hasGo := false
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() && strings.HasSuffix(d.Name(), ".go") {
			hasGo = true
			return fs.SkipAll
		}
		return nil
	})
	if err != nil || !hasGo {
		return nil
	}
	if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
		return nil
	}
	return Findings{{
		Kind:    KindWarning,
		Rule:    RuleGoModMissing,
		Package: filepath.Base(dir),
		Message: "holds Go template files but has no go.mod, so go build/test ./... will try to compile them; run `sedum generators install` or add a go.mod with `module generators`",
	}}
}

// /sedum:setFunction
