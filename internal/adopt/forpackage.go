package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/forpackage.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// ForPackage chooses the adopter form the package declares. A built-in adopter\n// needs no permission. A command adopter runs only when allowCommand is set.\nfunc ForPackage(pkg *genpkg.Package, allowCommand bool) (Adopter, error) {\n\tif pkg.Adopt == nil {\n\t\treturn nil, fmt.Errorf(\"package %s declares no adopter\", pkg.Name)\n\t}\n\tif name := pkg.Adopt.Builtin; name != \"\" {\n\t\tmk, ok := builtinAdopters[name]\n\t\tif !ok {\n\t\t\treturn nil, fmt.Errorf(\"package %s names the built-in adopter %q, which this sedum does not have\", pkg.Name, name)\n\t\t}\n\t\treturn mk(), nil\n\t}\n\tif len(pkg.Adopt.Command) > 0 {\n\t\tline := strings.Join(pkg.Adopt.Command, \" \")\n\t\tif !allowCommand {\n\t\t\treturn nil, fmt.Errorf(\"package %s declares a command adopter (%s); rerun with --allow-command-adopters to run it\", pkg.Name, line)\n\t\t}\n\t\treturn CommandAdopter{Argv: pkg.Adopt.Command}, nil\n\t}\n\treturn NewDeclarative(pkg)\n}","file":"internal/adopt/forpackage.go","name":"ForPackage"}}
// ForPackage chooses the adopter form the package declares. A built-in adopter
// needs no permission. A command adopter runs only when allowCommand is set.
func ForPackage(pkg *genpkg.Package, allowCommand bool) (Adopter, error) {
	if pkg.Adopt == nil {
		return nil, fmt.Errorf("package %s declares no adopter", pkg.Name)
	}
	if name := pkg.Adopt.Builtin; name != "" {
		mk, ok := builtinAdopters[name]
		if !ok {
			return nil, fmt.Errorf("package %s names the built-in adopter %q, which this sedum does not have", pkg.Name, name)
		}
		return mk(), nil
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
// sedum:setVar {"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:c7b9d9ac1655e8211b8b2ca8c49cf34b8d510d8d823d4f5eb60ec7e1018aa2eb","kwargs":{"code":"// builtinAdopters maps the name of a built-in adopter to its constructor.\nvar builtinAdopters = map[string]func() Adopter{\"go\": func() Adopter { return GoAdopter{} }}","file":"internal/adopt/forpackage.go","name":"builtinAdopters"}}
// builtinAdopters maps the name of a built-in adopter to its constructor.
var builtinAdopters = map[string]func() Adopter{"go": func() Adopter { return GoAdopter{} }}

// /sedum:setVar
