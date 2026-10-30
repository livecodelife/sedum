package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/forpackage.go","path":"fmt"},"kwargs":"sha256:86eacf2cdd4d0adb123e0fbd011683b529fa0099457de077e1c51c869965e849","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/forpackage.go","path":"strings"},"kwargs":"sha256:e47532b981cfdff965bbd15fd3de064bb5b717a4f340c6005a02fdecace4d304","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/forpackage.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:3da6e786065e55c80ad0cfce9584fc1070e8364371ec85611f0721af1a371384","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/forpackage.go","name":"ForPackage"},"kwargs":"sha256:55e148543e254abb8388b6d72ae7353818db4c9d7d7893f7cb6727d831371c54","body":"sha256:b073b22c6b49640fbaef05c70de6016ccc4f90ff1fbc984983ab976472a2e7cf"}
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
// sedum:setVar {"v":2,"tier":"owned","record":"prov-2026-1b97007c","gen":"sha256:53668ca0cb5e79371162faf44e46140660adad99ff2ee94d8e5c91f82882b8c8","pkg":"sha256:c7b9d9ac1655e8211b8b2ca8c49cf34b8d510d8d823d4f5eb60ec7e1018aa2eb","id":{"file":"internal/adopt/forpackage.go","name":"builtinAdopters"},"kwargs":"sha256:20cdbad7280a569f805758ba5a2b3a6a907aebff2912dc6418e68c8818e538e6","body":"sha256:74923371cca67f3552a3a977142d0b0159a36c2ee8cdf9d3a9d8302932294bd4"}
// builtinAdopters maps the name of a built-in adopter to its constructor.
var builtinAdopters = map[string]func() Adopter{"go": func() Adopter { return GoAdopter{} }}

// /sedum:setVar
