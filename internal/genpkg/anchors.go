package genpkg

import (
	"regexp"
	"slices"
	"strings"
)

// Transform reference checking.
//
// A reference is undefined unless it resolves to a pipeline the package
// declares or to a built-in operation. Loading needs the built-in names to
// decide that, and nothing else - it never applies a transform.
//
// BuiltinOperations is therefore names only, and it is M2's to claim: when the
// transform engine lands in internal/transform, this list moves there so that
// the vocabulary is declared once, next to the code that implements it
// (prov-2026-73127e53).

// BuiltinOperations is the closed set of operations Sedum ships. Pure
// string -> string, always available in every package.
var BuiltinOperations = []string{
	"pascal", "camel", "snake", "kebab",
	"upper", "lower",
	"plural", "singular",
	"prefix", "suffix",
}

// parameterizedOperations take a string literal argument after a colon:
// prefix:@, suffix:_path. Arguments are literals only - dynamic arguments
// would start the construction of an expression language - so there is nothing
// in the argument to resolve.
var parameterizedOperations = map[string]bool{"prefix": true, "suffix": true}

// expression matches a {{...}} span. Finding these and splitting on '|' is the
// whole extraction mechanism: enough to check that names resolve, and it
// commits the eventual renderer to nothing.
var expression = regexp.MustCompile(`\{\{([^{}]*)\}\}`)

// transformRefs returns the transform names referenced in s, in order of
// appearance and without duplicates. The first element of a piped expression
// is the value being transformed, not a transform.
func transformRefs(s string) []string {
	var (
		out  []string
		seen = map[string]bool{}
	)
	for _, m := range expression.FindAllStringSubmatch(s, -1) {
		parts := strings.Split(m[1], "|")
		for _, raw := range parts[1:] {
			name := strings.TrimSpace(raw)
			if name == "" || seen[name] {
				continue
			}
			seen[name] = true
			out = append(out, name)
		}
	}
	return out
}

// resolves reports whether a reference names a built-in operation or one of the
// package's declared pipelines.
func (p *Package) resolves(ref string) bool {
	// An operation argument is opaque: prefix:@ is the prefix operation.
	if name, _, ok := strings.Cut(ref, ":"); ok {
		return parameterizedOperations[name]
	}
	if slices.Contains(BuiltinOperations, ref) {
		// A bare prefix or suffix with no argument is not usable.
		return !parameterizedOperations[ref]
	}
	_, declared := p.Transforms[ref]
	return declared
}

// Anchor vocabulary. anchor is a scalar; these values are reserved and every
// other value names a marker planted by a file template (prov-2026-d1d61186).
const (
	AnchorStartOfFile = "start_of_file"
	AnchorEndOfFile   = "end_of_file"
	AnchorRegion      = "region"
	AnchorAfterMatch  = "after_match"
	AnchorBeforeMatch = "before_match"
	// AnchorMarker is reserved but not usable bare: an action anchored to
	// "marker" has not said which marker.
	AnchorMarker = "marker"
)

func reservedAnchor(anchor string) bool {
	switch anchor {
	case AnchorStartOfFile, AnchorEndOfFile, AnchorRegion, AnchorAfterMatch, AnchorBeforeMatch, AnchorMarker:
		return true
	}
	return false
}

// MarkerAnchor returns the marker name an action is anchored to, and whether it
// is anchored to a marker at all. A reserved keyword is not a marker name.
func (a *Action) MarkerAnchor() (string, bool) {
	if a.Anchor == "" || reservedAnchor(a.Anchor) {
		return "", false
	}
	return a.Anchor, true
}

// markerDecl is the text a file template plants to create an anchor point. The
// prefix is the package's declared comment_prefix, since #, // and -- all
// appear across targets.
func markerDecl(commentPrefix, name string) string {
	return commentPrefix + " sedum:anchor:" + name
}

// plantedMarkers returns every marker name planted across the package's file
// template contents.
func plantedMarkers(commentPrefix string, contents []string) map[string]bool {
	// The prefix is author-supplied, so it is escaped rather than
	// interpolated into a pattern.
	re := regexp.MustCompile(regexp.QuoteMeta(markerDecl(commentPrefix, "")) + `([A-Za-z0-9_.-]+)`)

	out := map[string]bool{}
	for _, c := range contents {
		for _, m := range re.FindAllStringSubmatch(c, -1) {
			out[m[1]] = true
		}
	}
	return out
}
