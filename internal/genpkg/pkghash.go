package genpkg

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/pkghash.go","path":"crypto/sha256"}}
import "crypto/sha256"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/pkghash.go","path":"encoding/hex"}}
import "encoding/hex"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/pkghash.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/pkghash.go","path":"github.com/livecodelife/provenance/record"}}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// The pkg hash a marker carries: everything that decides what an action\n// renders for fixed kwargs, hashed at load (prov-2026-66bea9f3).\n//\n// A region whose pkg differs from its action's current one was rendered from an\n// older package, and is brought up to date by re-rendering it from the kwargs\n// already on its marker - no model call, because nothing the model chose has\n// changed. So the hash has to cover every input of rendering and should cover\n// nothing else:\n//\n//   - the resolved action declaration: name, kwargs with their types,\n//     required flags and defaults, discriminator, variants in declared order,\n//     injects_into, the anchor fields, and identity\n//   - the content of each of the action's templates, keyed by variant name,\n//     or by the empty string for a simple action\n//   - the package's transforms and op_exceptions, since any template may use\n//     any of them\n//\n// Descriptions are left out, of the action and of each kwarg: they reach the\n// model and never the renderer. The package-wide part errs the other way on\n// purpose. An unrelated transform edit costs a re-render, while leaving out an\n// input that does matter would make a stale region read as current.\n//\n// The inflection table ships with Sedum rather than with a package, so a Sedum\n// upgrade that changes an irregular plural is not seen here.\n//\n// The encoding is published in conformance/pkghash so a tool above Sedum can\n// compute the same value without linking it.\n\n// pkgFormat names this encoding inside the hashed object, so a later change to\n// what is covered changes every hash rather than colliding with this one.\nconst pkgFormat = \"sedum.pkg/1\"","file":"internal/genpkg/pkghash.go","name":"pkgFormat"}}
// The pkg hash a marker carries: everything that decides what an action
// renders for fixed kwargs, hashed at load (prov-2026-66bea9f3).
//
// A region whose pkg differs from its action's current one was rendered from an
// older package, and is brought up to date by re-rendering it from the kwargs
// already on its marker - no model call, because nothing the model chose has
// changed. So the hash has to cover every input of rendering and should cover
// nothing else:
//
//   - the resolved action declaration: name, kwargs with their types,
//     required flags and defaults, discriminator, variants in declared order,
//     injects_into, the anchor fields, and identity
//   - the content of each of the action's templates, keyed by variant name,
//     or by the empty string for a simple action
//   - the package's transforms and op_exceptions, since any template may use
//     any of them
//
// Descriptions are left out, of the action and of each kwarg: they reach the
// model and never the renderer. The package-wide part errs the other way on
// purpose. An unrelated transform edit costs a re-render, while leaving out an
// input that does matter would make a stale region read as current.
//
// The inflection table ships with Sedum rather than with a package, so a Sedum
// upgrade that changes an irregular plural is not seen here.
//
// The encoding is published in conformance/pkghash so a tool above Sedum can
// compute the same value without linking it.

// pkgFormat names this encoding inside the hashed object, so a later change to
// what is covered changes every hash rather than colliding with this one.
const pkgFormat = "sedum.pkg/1"

// /sedum:setConst
// sedum:setFunction {"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// stampPkgHashes sets PkgHash on every action in a loaded package.\nfunc stampPkgHashes(pkg *Package) error {\n\tfor _, action := range pkg.Actions {\n\t\tcanonical, err := pkgCanonical(pkg, action)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tsum := sha256.Sum256(canonical)\n\t\taction.PkgHash = \"sha256:\" + hex.EncodeToString(sum[:])\n\t}\n\treturn nil\n}","file":"internal/genpkg/pkghash.go","name":"stampPkgHashes"}}
// stampPkgHashes sets PkgHash on every action in a loaded package.
func stampPkgHashes(pkg *Package) error {
	for _, action := range pkg.Actions {
		canonical, err := pkgCanonical(pkg, action)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(canonical)
		action.PkgHash = "sha256:" + hex.EncodeToString(sum[:])
	}
	return nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// pkgCanonical is the exact bytes the pkg hash is taken over: the object below\n// encoded by the provenance library's record.Canonical, the same canonical form\n// gen_hash uses.\nfunc pkgCanonical(pkg *Package, action *Action) (canonical []byte, err error) {\n\t// record.Canonical panics on a value it cannot encode. Load validation\n\t// keeps kwarg defaults to string, int, bool and list, so this is\n\t// unreachable for a valid package; it is an error rather than a crash\n\t// for any other.\n\tdefer func() {\n\t\tif r := recover(); r != nil {\n\t\t\tcanonical, err = nil, fmt.Errorf(\"package %s: action %s: cannot be hashed: %v\", pkg.Name, action.Name, r)\n\t\t}\n\t}()\n\n\tkwargs := map[string]any{}\n\tfor name, k := range action.Kwargs {\n\t\tkwargs[name] = map[string]any{\n\t\t\t\"type\":     k.Type,\n\t\t\t\"required\": k.Required,\n\t\t\t\"default\":  k.Default,\n\t\t}\n\t}\n\n\ttemplates := map[string]any{}\n\tpaths := map[string]string{}\n\tif action.Template != \"\" {\n\t\tpaths[SoleTemplate] = action.Template\n\t}\n\tfor variant, path := range action.Templates {\n\t\tpaths[variant] = path\n\t}\n\tfor key, path := range paths {\n\t\tcontent, ok := pkg.ActionTemplate(path)\n\t\tif !ok {\n\t\t\treturn nil, fmt.Errorf(\"package %s: action %s: template %s was not loaded\", pkg.Name, action.Name, path)\n\t\t}\n\t\ttemplates[key] = content\n\t}\n\n\ttransforms := map[string]any{}\n\tfor name, ops := range pkg.Transforms {\n\t\ttransforms[name] = stringList(ops)\n\t}\n\texceptions := map[string]any{}\n\tfor op, table := range pkg.OpExceptions {\n\t\twords := map[string]any{}\n\t\tfor word, replacement := range table {\n\t\t\twords[word] = replacement\n\t\t}\n\t\texceptions[op] = words\n\t}\n\n\tobject := map[string]any{\n\t\t\"format\": pkgFormat,\n\t\t\"action\": map[string]any{\n\t\t\t\"name\":           action.Name,\n\t\t\t\"kwargs\":         kwargs,\n\t\t\t\"discriminator\":  action.Discriminator,\n\t\t\t\"variants\":       stringList(action.Variants),\n\t\t\t\"injects_into\":   action.InjectsInto,\n\t\t\t\"anchor\":         action.Anchor,\n\t\t\t\"anchor_start\":   action.AnchorStart,\n\t\t\t\"anchor_end\":     action.AnchorEnd,\n\t\t\t\"anchor_pattern\": action.AnchorPattern,\n\t\t\t\"identity\":       stringList(action.Identity),\n\t\t},\n\t\t\"templates\":     templates,\n\t\t\"transforms\":    transforms,\n\t\t\"op_exceptions\": exceptions,\n\t}\n\n\treturn []byte(record.Canonical(object)), nil\n}","file":"internal/genpkg/pkghash.go","name":"pkgCanonical"}}
// pkgCanonical is the exact bytes the pkg hash is taken over: the object below
// encoded by the provenance library's record.Canonical, the same canonical form
// gen_hash uses.
func pkgCanonical(pkg *Package, action *Action) (canonical []byte, err error) {
	// record.Canonical panics on a value it cannot encode. Load validation
	// keeps kwarg defaults to string, int, bool and list, so this is
	// unreachable for a valid package; it is an error rather than a crash
	// for any other.
	defer func() {
		if r := recover(); r != nil {
			canonical, err = nil, fmt.Errorf("package %s: action %s: cannot be hashed: %v", pkg.Name, action.Name, r)
		}
	}()

	kwargs := map[string]any{}
	for name, k := range action.Kwargs {
		kwargs[name] = map[string]any{
			"type":     k.Type,
			"required": k.Required,
			"default":  k.Default,
		}
	}

	templates := map[string]any{}
	paths := map[string]string{}
	if action.Template != "" {
		paths[SoleTemplate] = action.Template
	}
	for variant, path := range action.Templates {
		paths[variant] = path
	}
	for key, path := range paths {
		content, ok := pkg.ActionTemplate(path)
		if !ok {
			return nil, fmt.Errorf("package %s: action %s: template %s was not loaded", pkg.Name, action.Name, path)
		}
		templates[key] = content
	}

	transforms := map[string]any{}
	for name, ops := range pkg.Transforms {
		transforms[name] = stringList(ops)
	}
	exceptions := map[string]any{}
	for op, table := range pkg.OpExceptions {
		words := map[string]any{}
		for word, replacement := range table {
			words[word] = replacement
		}
		exceptions[op] = words
	}

	object := map[string]any{
		"format": pkgFormat,
		"action": map[string]any{
			"name":           action.Name,
			"kwargs":         kwargs,
			"discriminator":  action.Discriminator,
			"variants":       stringList(action.Variants),
			"injects_into":   action.InjectsInto,
			"anchor":         action.Anchor,
			"anchor_start":   action.AnchorStart,
			"anchor_end":     action.AnchorEnd,
			"anchor_pattern": action.AnchorPattern,
			"identity":       stringList(action.Identity),
		},
		"templates":     templates,
		"transforms":    transforms,
		"op_exceptions": exceptions,
	}

	return []byte(record.Canonical(object)), nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// stringList turns a possibly nil string list into a JSON list, so an absent list\n// and an empty one encode alike.\nfunc stringList(list []string) []any {\n\tout := make([]any, len(list))\n\tfor i, s := range list {\n\t\tout[i] = s\n\t}\n\treturn out\n}","file":"internal/genpkg/pkghash.go","name":"stringList"}}
// stringList turns a possibly nil string list into a JSON list, so an absent list
// and an empty one encode alike.
func stringList(list []string) []any {
	out := make([]any, len(list))
	for i, s := range list {
		out[i] = s
	}
	return out
}

// /sedum:setFunction
