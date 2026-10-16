package genpkg

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/gomod.go","path":"io/fs"}}
import "io/fs"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/gomod.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/gomod.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/gomod.go","path":"strings"}}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-806865fb","gen":"sha256:e631702fd73e9908ecef5a5b2b0ba1b871e3b9bb22f36b8b974964eb896a5156","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// checkGoModBoundary reports a warning when dir holds Go template files\n// anywhere below it but has no go.mod at its root.\nfunc checkGoModBoundary(dir string) Findings {\n\thasGo := false\n\terr := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tif d.Type().IsRegular() && strings.HasSuffix(d.Name(), \".go\") {\n\t\t\thasGo = true\n\t\t\treturn fs.SkipAll\n\t\t}\n\t\treturn nil\n\t})\n\tif err != nil || !hasGo {\n\t\treturn nil\n\t}\n\tif _, statErr := os.Stat(filepath.Join(dir, \"go.mod\")); statErr == nil {\n\t\treturn nil\n\t}\n\treturn Findings{{\n\t\tKind:    KindWarning,\n\t\tRule:    RuleGoModMissing,\n\t\tPackage: filepath.Base(dir),\n\t\tMessage: \"holds Go template files but has no go.mod, so go build/test ./... will try to compile them; run `sedum generators install` or add a go.mod with `module generators`\",\n\t}}\n}","file":"internal/genpkg/gomod.go","name":"checkGoModBoundary"}}
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
