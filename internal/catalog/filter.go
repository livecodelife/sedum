package catalog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/livecodelife/sedum/internal/render"
)

// Removal is one action Phase 4 will not offer a given record, and why.
//
// It is reported rather than silently dropped, because "why did the model
// never see this action" is exactly the question sedum actions --all exists to
// answer for the hidden case (prov-2026-...), and a filter that could not say
// why it removed something would be as opaque as the failure it replaces.
type Removal struct {
	Action  string
	Package string
	Reason  string
}

// FilterForRecord keeps only the actions a record's authorized paths make it
// possible to invoke.
//
// Nothing has bound a kwarg yet when this runs, so the question it asks is
// never "does this invocation's path match" - only "could any invocation of
// this action possibly land on one of these paths." An injects_into pattern
// is walked the same way render.Parse already walks it: its fixed, literal
// text has to appear in an authorized path, and every {{...}} expression -
// bare or piped alike - is treated as matching any run of characters, because
// no placeholder value can be ruled out before anything is bound.
//
// The match is conservative in one direction only. An action is removed only
// when its fixed text cannot appear in any authorized path for any
// placeholder value; a false removal (dropping something that might still be
// legal) is a worse mistake than a false keep (showing something that turns
// out not to apply), so ambiguity resolves toward keeping. A free-target
// pattern - the whole of injects_into is one placeholder, per
// prov-2026-14c832bf - has no fixed text to fail to match, so it is never
// removed here; its applicability is anchor-based already.
//
// A composite is removed only when at least one of its entries (one per
// child, in the order catalog.targets already collected them) can never
// match, because expand.Targets requires every child's rendered target to be
// authorized for the composite's invocation to validate at all - one
// structurally unreachable child dooms the whole selection regardless of the
// others.
//
// Validation is unaffected: it runs against the catalog this never sees,
// built the same way it is today. This only decides what Phase 4 shows.
func FilterForRecord(c Catalog, authorizedPaths []string) (kept Catalog, removed []Removal) {
	for _, action := range c.Actions {
		if reason := unreachable(action, authorizedPaths); reason != "" {
			removed = append(removed, Removal{Action: action.Name, Package: action.Package, Reason: reason})
			continue
		}
		kept.Actions = append(kept.Actions, action)
	}
	return kept, removed
}

// unreachable names the first injects_into entry that cannot resolve to any
// authorized path, or "" when every entry can.
func unreachable(action Action, authorizedPaths []string) string {
	for _, pattern := range action.InjectsInto {
		if matchesAny(patternRegexp(pattern), authorizedPaths) {
			continue
		}
		return fmt.Sprintf(
			"injects_into %s cannot resolve to any of this record's authorized paths (%s)",
			strconv.Quote(pattern), quotePaths(authorizedPaths))
	}
	return ""
}

// patternRegexp turns an injects_into pattern into a regular expression that
// matches whatever the pattern could possibly render.
//
// It walks the pattern the same way render.translate does: literal text
// between expressions is copied through, escaped so nothing in it is read as
// a metacharacter, and each {{...}} expression - transforms and all - becomes
// a wildcard rather than a rendering, because rendering needs bound kwargs
// this filter does not have.
func patternRegexp(pattern string) *regexp.Regexp {
	exprs, problems := render.Parse(pattern)
	if len(problems) > 0 {
		// injects_into is validated at package load (genpkg checkTemplates),
		// so a pattern that fails to parse here is a defect somewhere else,
		// not something this filter can prove is unreachable. Matching
		// everything is the conservative answer.
		return matchEverything
	}

	var b strings.Builder
	b.WriteString("^")
	rest := pattern
	for _, e := range exprs {
		idx := strings.Index(rest, e.Source)
		b.WriteString(regexp.QuoteMeta(rest[:idx]))
		b.WriteString("(?s:.*)")
		rest = rest[idx+len(e.Source):]
	}
	b.WriteString(regexp.QuoteMeta(rest))
	b.WriteString("$")

	re, err := regexp.Compile(b.String())
	if err != nil {
		// QuoteMeta guarantees the literal segments are safe, and (?s:.*) is
		// fixed and valid, so this cannot happen - kept as the same fallback
		// rather than a panic, for the same reason as above.
		return matchEverything
	}
	return re
}

// matchEverything is what an injects_into pattern this filter cannot analyze
// - which never happens against a package Sedum loaded - resolves to. It
// matches, never removes.
var matchEverything = regexp.MustCompile(`^(?s:.*)$`)

func matchesAny(re *regexp.Regexp, paths []string) bool {
	for _, p := range paths {
		if re.MatchString(p) {
			return true
		}
	}
	return false
}

func quotePaths(paths []string) string {
	if len(paths) == 0 {
		return "none"
	}
	quoted := make([]string, 0, len(paths))
	for _, p := range paths {
		quoted = append(quoted, strconv.Quote(p))
	}
	return strings.Join(quoted, ", ")
}
