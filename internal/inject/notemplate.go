package inject

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/inject/notemplate.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/inject/notemplate.go","path":"path"}}
import "path"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/inject/notemplate.go","path":"slices"}}
import "slices"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/inject/notemplate.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/inject/notemplate.go","path":"github.com/livecodelife/sedum/internal/filetmpl"}}
import "github.com/livecodelife/sedum/internal/filetmpl"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// explainUnmatched explains why an action could not be placed in a file\n// because the file's template set has no template matching its path.\nfunc explainUnmatched(inv Invocation, cause error) error {\n\tbase := fmt.Errorf(\"action %s cannot be injected into %s: %w; the file is not shaped the way the action assumes\", inv.Action.Name, inv.Path, cause)\n\tif len(inv.Package.FileTemplates) == 0 {\n\t\treturn fmt.Errorf(\"action %s cannot be injected into %s: no file template in package %s matches this path, so the file was created empty and has none of the markers a template plants (%w); %s ships no file templates\", inv.Action.Name, inv.Path, inv.Package.Name, cause, inv.Package.Name)\n\t}\n\tset, err := filetmpl.NewSet(inv.Package.FileTemplates)\n\tif err != nil {\n\t\treturn base\n\t}\n\tif _, matched, err := set.Match(inv.Path); err != nil || matched {\n\t\treturn base\n\t}\n\tif slices.Contains(inv.Package.FileTemplates, \"_default\"+strings.ToLower(path.Ext(inv.Path))) {\n\t\treturn base\n\t}\n\treturn fmt.Errorf(\"action %s cannot be injected into %s: no file template in package %s matches this path, so the file was created empty and has none of the markers a template plants (%w); %s's file templates are: %s\", inv.Action.Name, inv.Path, inv.Package.Name, cause, inv.Package.Name, strings.Join(inv.Package.FileTemplates, \", \"))\n}","file":"internal/inject/notemplate.go","name":"explainUnmatched"}}
// explainUnmatched explains why an action could not be placed in a file
// because the file's template set has no template matching its path.
func explainUnmatched(inv Invocation, cause error) error {
	base := fmt.Errorf("action %s cannot be injected into %s: %w; the file is not shaped the way the action assumes", inv.Action.Name, inv.Path, cause)
	if len(inv.Package.FileTemplates) == 0 {
		return fmt.Errorf("action %s cannot be injected into %s: no file template in package %s matches this path, so the file was created empty and has none of the markers a template plants (%w); %s ships no file templates", inv.Action.Name, inv.Path, inv.Package.Name, cause, inv.Package.Name)
	}
	set, err := filetmpl.NewSet(inv.Package.FileTemplates)
	if err != nil {
		return base
	}
	if _, matched, err := set.Match(inv.Path); err != nil || matched {
		return base
	}
	if slices.Contains(inv.Package.FileTemplates, "_default"+strings.ToLower(path.Ext(inv.Path))) {
		return base
	}
	return fmt.Errorf("action %s cannot be injected into %s: no file template in package %s matches this path, so the file was created empty and has none of the markers a template plants (%w); %s's file templates are: %s", inv.Action.Name, inv.Path, inv.Package.Name, cause, inv.Package.Name, strings.Join(inv.Package.FileTemplates, ", "))
}

// /sedum:setFunction
