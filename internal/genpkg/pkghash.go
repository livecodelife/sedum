package genpkg

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

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

// pkgCanonical is the exact bytes the pkg hash is taken over.
func pkgCanonical(pkg *Package, action *Action) ([]byte, error) {
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

	var b strings.Builder
	if err := writeCanonical(&b, object); err != nil {
		return nil, fmt.Errorf("package %s: action %s: %w", pkg.Name, action.Name, err)
	}
	return []byte(b.String()), nil
}

// stringList turns a possibly nil string list into a JSON list, so an absent list
// and an empty one encode alike.
func stringList(list []string) []any {
	out := make([]any, len(list))
	for i, s := range list {
		out[i] = s
	}
	return out
}

// writeCanonical writes v as JSON with byte-sorted object keys, no
// insignificant whitespace, and only the escapes JSON requires - the same
// rules gen_hash's canonical form uses, so a reader learns one encoding.
func writeCanonical(b *strings.Builder, v any) error {
	switch v := v.(type) {
	case nil:
		b.WriteString("null")
	case string:
		writeCanonicalString(b, v)
	case bool, int, int64, uint64, float64:
		encoded, err := json.Marshal(v)
		if err != nil {
			return err
		}
		b.Write(encoded)
	case []any:
		b.WriteByte('[')
		for i, item := range v {
			if i > 0 {
				b.WriteByte(',')
			}
			if err := writeCanonical(b, item); err != nil {
				return err
			}
		}
		b.WriteByte(']')
	case map[string]any:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				b.WriteByte(',')
			}
			writeCanonicalString(b, k)
			b.WriteByte(':')
			if err := writeCanonical(b, v[k]); err != nil {
				return err
			}
		}
		b.WriteByte('}')
	default:
		return fmt.Errorf("a value of type %T cannot be hashed", v)
	}
	return nil
}

func writeCanonicalString(b *strings.Builder, s string) {
	const hexDigits = "0123456789abcdef"
	b.WriteByte('"')
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == '"':
			b.WriteString(`\"`)
		case r == '\\':
			b.WriteString(`\\`)
		case r < 0x20:
			b.WriteString(`\u00`)
			b.WriteByte(hexDigits[r>>4])
			b.WriteByte(hexDigits[r&0xf])
		default:
			b.WriteString(s[i : i+size])
		}
		i += size
	}
	b.WriteByte('"')
}
