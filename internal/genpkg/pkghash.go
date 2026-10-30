package genpkg

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/pkghash.go","path":"crypto/sha256"},"kwargs":"sha256:50bc76c5ef0b7dab3f582c6649890d878e6c70e3218b56e4f2205d2a9a04fd27","body":"sha256:449ec1ed24bf227b1713290d01a3cedc4fc91ebc26ceac7be6fe2ad3eedd60a1"}
import "crypto/sha256"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/pkghash.go","path":"encoding/hex"},"kwargs":"sha256:f445bcad71f0cfbf5e10eac1b365c119fb07718a29f94cceec8eedd265582da0","body":"sha256:5b4bad0c96848696f832da102bacefeeafa1eec6d331d5241ac31b3f4b803432"}
import "encoding/hex"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/pkghash.go","path":"fmt"},"kwargs":"sha256:c05024ba196483203ac7e962c42bfa4e5b5359d5e497d1e7013e11ee3b556c27","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/pkghash.go","path":"github.com/livecodelife/provenance/record"},"kwargs":"sha256:25725f681597ed180e636790c4ab31d962cab94045a3a8b775ea047e6b0cd9cd","body":"sha256:acbf8628061291d11aa39e9ea64ecc91095a117057cfe6f1958ef2c437fac003"}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/genpkg/pkghash.go","name":"pkgFormat"},"kwargs":"sha256:b7146f9b2d99b8dbd5e1fafba1bc3918060e1ed181041b5b5331166162e792f8","body":"sha256:b0951460dc30cb5434ed441eb07a0c016cfa7e9f928fdb5e986a75488088788e"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/pkghash.go","name":"stampPkgHashes"},"kwargs":"sha256:32c779598fd60f2757156210d1e2cad7d9a9544a9bd53e5402e2b2ff805ba07e","body":"sha256:fdfc9c0ca0097c40577c5012e1c368de197f78a83db0e462c4f361d3dd8127c7"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/pkghash.go","name":"pkgCanonical"},"kwargs":"sha256:c8caeaf6763a1c9dcd09f2ce06d548a61d4cd443d7771cb4012244ef58bed8ec","body":"sha256:7257d99b64136c3a6588026ae5a29892509ad1be44549e8f2c5a175b54e9f4f4"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1a55bc63","gen":"sha256:50dbb2109211cfaffb8624aa1e7c541e9662b5cc018317ad8d168033916054e4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/pkghash.go","name":"stringList"},"kwargs":"sha256:650709db79226d4e1a667cf50dd8106dbf11f763174f041d3d7db3c74338e412","body":"sha256:ab218d772fbf12b8aa25e30d17ffc5ac619ab614162a20460a5b470377bf9d32"}
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
