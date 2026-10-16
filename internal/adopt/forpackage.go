package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:2894782c60458ebbc1f2f84fee366c5078e0aba56f2009f295c1f09131a45d47","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func ForPackage(pkg *genpkg.Package, allowCommand bool) (Adopter, error) {\n\tif pkg.Adopt == nil {\n\t\treturn nil, fmt.Errorf(\"package %s declares no adopter\", pkg.Name)\n\t}\n\tif len(pkg.Adopt.Command) > 0 {\n\t\tline := strings.Join(pkg.Adopt.Command, \" \")\n\t\tif !allowCommand {\n\t\t\treturn nil, fmt.Errorf(\"package %s declares a command adopter (%s); rerun with --allow-command-adopters to run it\", pkg.Name, line)\n\t\t}\n\t\treturn CommandAdopter{Argv: pkg.Adopt.Command}, nil\n\t}\n\treturn NewDeclarative(pkg)\n}","file":"internal/adopt/forpackage.go","name":"ForPackage"}}
func ForPackage(pkg *genpkg.Package, allowCommand bool) (Adopter, error) {
	if pkg.Adopt == nil {
		return nil, fmt.Errorf("package %s declares no adopter", pkg.Name)
	}
	if len(pkg.Adopt.Command) > 0 {
		line := strings.Join(pkg.Adopt.Command, " ")
		if !allowCommand {
			return nil, fmt.Errorf("package %s declares a command adopter (%s); rerun with --allow-command-adopters to run it", pkg.Name, line)
		}
		return CommandAdopter{Argv: pkg.Adopt.Command}, nil
	}
	return NewDeclarative(pkg)
}

// /sedum:setFunction
