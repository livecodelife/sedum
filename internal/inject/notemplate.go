package inject

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/inject/notemplate.go","path":"fmt"},"kwargs":"sha256:0d2fffb47a00107b609381b9b2ec368eeefa46c01497052e8c3587d81de11e34","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/inject/notemplate.go","path":"path"},"kwargs":"sha256:63e8cc010f5dba2f30d56960990fba2d1503a94e8db32283a5ee48d5884e305d","body":"sha256:c8c46e2b4136d27afcfddd73e765db239e056b6cba04d989689431e735645f13"}
import "path"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/inject/notemplate.go","path":"slices"},"kwargs":"sha256:41bf52891d53360d07c4b0f88cf10fb6c53b7163067fd7339569cc9b76a3595e","body":"sha256:5467eadcc4cfa5d6fc93588581f935056ea78f2a557c0da73bba03418792008c"}
import "slices"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/inject/notemplate.go","path":"strings"},"kwargs":"sha256:f3bcb9586db8ac6def8d479edc3afcb72538c40eaaf7dd312be5de683632efd6","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/inject/notemplate.go","path":"github.com/livecodelife/sedum/internal/filetmpl"},"kwargs":"sha256:d2f937742a04bb2f7dea4290a72308a94ac03abae42eaa47248c66905973df70","body":"sha256:832740aee57ce657667cbf064bf0d989a5dbb35fc52190bc573f2f94f02ea7d6"}
import "github.com/livecodelife/sedum/internal/filetmpl"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-7a9b82ae","gen":"sha256:e4093a0ea395d6b79b75ea0951e0c8963a26a4dca5e5f416b1a1639809e1b354","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/inject/notemplate.go","name":"explainUnmatched"},"kwargs":"sha256:b261685f8b37f8fef0d3ec8aae99baabe1880c228461775270bb2c9a143af8c5","body":"sha256:ce8a287d9fc7d60cbb5536f92110388e2bb76c0a81d703d4c36e73fc60681ed2"}
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
