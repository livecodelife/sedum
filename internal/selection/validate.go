package selection

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/calebcowen/sedum/internal/catalog"
	"github.com/calebcowen/sedum/internal/expand"
	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/recording"
	"github.com/calebcowen/sedum/internal/resolve"
)

// Phase 5: deterministic validation of the model's output.
//
// Every check produces an error specific enough to re-prompt with. That is not
// a quality goal, it is the mechanism: the retry loop's entire value is that
// the model is told exactly what was wrong, and a generic rejection would spend
// a call teaching it nothing.

// The rules Phase 5 can report.
//
// Named constants rather than literals at the point of use, because these
// outlived being an implementation detail: the eval harness stores them in an
// append-only results file (prov-2026-2256e6fa), so a slug written today is
// read by something a year from now with only the file to resolve it against.
// That makes the set a format, in the same sense the ownership marker is one
// (prov-2026-72775ae5), and a format defined by ten scattered literals drifts
// the first time one is reworded.
//
// Declared together so the vocabulary can be read in one place, and so a
// caller can test for a rule without matching on prose (prov-2026-9a554c93).
const (
	// RuleUnknownAction names an invocation of an action no package exposes.
	RuleUnknownAction = "unknown_action"
	// RuleAmbiguousAction names an action more than one package claims.
	RuleAmbiguousAction = "ambiguous_action"
	// RuleMissingKwarg names a required kwarg the model did not mention.
	RuleMissingKwarg = "missing_kwarg"
	// RuleEmptyKwarg names a required kwarg the model mentioned and left
	// empty. Distinct from RuleMissingKwarg on purpose: a kwarg nobody wrote
	// and one deliberately emptied are different mistakes, and a shared slug
	// would make them one row in every count drawn from a results file.
	RuleEmptyKwarg = "empty_kwarg"
	// RuleUnknownKwarg names a kwarg the action does not declare.
	RuleUnknownKwarg = "unknown_kwarg"
	// RuleKwargType names a kwarg bound to a value of the wrong type.
	RuleKwargType = "kwarg_type"
	// RuleMissingDerivedKwarg names a kwarg the selected template renders but
	// the invocation did not bind, whatever the schema calls optional.
	RuleMissingDerivedKwarg = "missing_derived_kwarg"
	// RuleVariant names a discriminator bound to a value with no template.
	RuleVariant = "variant"
	// RuleUnauthorizedPath names an invocation whose rendered target is not a
	// path the record authorizes.
	RuleUnauthorizedPath = "unauthorized_path"
)

// Violation is one Phase 5 failure.
//
// Rule is a short stable name for the run log, so that a run's failures can be
// counted by kind. Detail is what the model is told, and is written to be read
// by the thing that has to correct it.
type Violation struct {
	// Index is the 1-based position of the offending invocation, or zero for
	// a fault in the response as a whole.
	Index  int
	Action string
	Rule   string
	Detail string
}

func (v Violation) String() string {
	if v.Action == "" {
		return fmt.Sprintf("[%s] %s", v.Rule, v.Detail)
	}
	return fmt.Sprintf("[%s] %s: %s", v.Rule, v.Action, v.Detail)
}

// validate runs every check against every invocation.
//
// Two orderings are load-bearing. Within an invocation the checks are ordered,
// and the rendered-path check runs only after the schema checks pass: rendering
// injects_into against a missing or wrongly typed kwarg produces a diagnostic
// about a path when the fault is an argument, and re-prompting with that
// teaches the wrong lesson. Across invocations there is no ordering and no
// early exit - a response with three faults re-prompts with three, because one
// model call per fault is exactly the cost this design exists to avoid
// (prov-2026-9dcf2658).
func validate(cat catalog.Catalog, packages []*genpkg.Package, files []resolve.File, invocations []recording.Invocation) []Violation {
	var out []Violation
	for i, inv := range invocations {
		out = append(out, validateOne(cat, packages, files, inv, i+1)...)
	}
	return out
}

func validateOne(cat catalog.Catalog, packages []*genpkg.Package, files []resolve.File, inv recording.Invocation, index int) []Violation {
	entry, v := resolveEntry(cat, inv, index)
	if v != nil {
		return []Violation{*v}
	}

	// The schema checks are reported together: a response binding a wrong
	// type and omitting a required argument should hear about both.
	violations := checkKwargs(entry, inv, index)
	violations = append(violations, checkVariant(entry, inv, index)...)
	if len(violations) > 0 {
		return violations
	}

	return checkTargets(packages, files, inv, index)
}

// resolveEntry finds the catalog entry an invocation names.
//
// An unexposed action is not merely rejected here - it is absent from the
// catalog entirely, so naming one is indistinguishable from naming something
// that does not exist. That is the point of the tier: the invalidity is
// unrepresentable rather than caught.
func resolveEntry(cat catalog.Catalog, inv recording.Invocation, index int) (catalog.Action, *Violation) {
	matches := cat.Lookup(inv.Action)
	switch {
	case len(matches) == 1:
		return matches[0], nil

	case len(matches) == 0:
		return catalog.Action{}, &Violation{
			Index:  index,
			Action: inv.Action,
			Rule:   RuleUnknownAction,
			Detail: fmt.Sprintf(
				"no action named %q is available for this change; the catalog offers %s",
				inv.Action, quoteList(cat.Names())),
		}

	default:
		// The response names an action and its kwargs and never a package,
		// so nothing in it could say which was meant. It is reported rather
		// than resolved by an ordering rule nothing declared.
		var owners []string
		for _, m := range matches {
			owners = append(owners, m.Package)
		}
		return catalog.Action{}, &Violation{
			Index:  index,
			Action: inv.Action,
			Rule:   RuleAmbiguousAction,
			Detail: fmt.Sprintf(
				"more than one package this change spans declares %q (%s), so nothing says which is meant",
				inv.Action, strings.Join(owners, ", ")),
		}
	}
}

// checkKwargs holds the bound arguments to the action's declared schema: every
// required name present, no name the action does not declare, and every value
// of its declared type.
func checkKwargs(entry catalog.Action, inv recording.Invocation, index int) []Violation {
	var out []Violation

	var missing, empty []string
	for name, k := range entry.Kwargs {
		if !k.Required {
			continue
		}
		v, bound := inv.Kwargs[name]
		switch {
		case !bound:
			missing = append(missing, name)
		case isEmpty(v):
			empty = append(empty, name)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		out = append(out, Violation{
			Index:  index,
			Action: inv.Action,
			Rule:   RuleMissingKwarg,
			Detail: fmt.Sprintf("required %s not bound: %s",
				plural(len(missing), "kwarg is", "kwargs are"), quoteList(missing)),
		})
	}
	if len(empty) > 0 {
		sort.Strings(empty)
		out = append(out, Violation{
			Index:  index,
			Action: inv.Action,
			Rule:   RuleEmptyKwarg,
			Detail: fmt.Sprintf("required %s bound to nothing: %s; give each a value or say why the action applies without one",
				plural(len(empty), "kwarg is", "kwargs are"), quoteList(empty)),
		})
	}

	if v := checkDerived(entry, inv, index); v != nil {
		out = append(out, *v)
	}

	var unknown []string
	for name := range inv.Kwargs {
		if _, declared := entry.Kwargs[name]; !declared {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		out = append(out, Violation{
			Index:  index,
			Action: inv.Action,
			Rule:   RuleUnknownKwarg,
			Detail: fmt.Sprintf("%s not declared by this action: %s; it declares %s",
				plural(len(unknown), "kwarg is", "kwargs are"), quoteList(unknown), quoteList(declaredNames(entry))),
		})
	}

	for _, name := range sortedNames(inv.Kwargs) {
		k, declared := entry.Kwargs[name]
		if !declared {
			// Already reported as unknown. Type-checking a value against a
			// schema that does not exist would be a second violation for one
			// mistake.
			continue
		}
		if got, ok := typeOf(inv.Kwargs[name]); !ok || got != k.Type {
			out = append(out, Violation{
				Index:  index,
				Action: inv.Action,
				Rule:   RuleKwargType,
				Detail: fmt.Sprintf("kwarg %q is declared %s but was bound to %s",
					name, k.Type, describe(inv.Kwargs[name])),
			})
		}
	}

	return out
}

// checkDerived holds an invocation to what its selected template will actually
// render, which the kwarg schema alone cannot express.
//
// This is the half of the requirement that comes from the template rather than
// the declaration. A discriminated action shares one schema across every
// variant, so a kwarg that one variant needs and another forbids can only be
// declared optional - and then the catalog says optional, the model omits it,
// every Phase 5 check passes, and Phase 6 halts rendering a template with the
// retry loop already skipped because nothing was wrong with the selection.
//
// That is the shape prov-2026-9dcf2658 named for paths and this closes for
// values: a later phase must not reject what this one accepted, because the
// loop that could have fixed it never runs (prov-2026-369544c1).
//
// The violation is reported separately from missing_kwarg, and names the
// variant, because the two have different fixes. A declared requirement that is
// absent is a binding mistake; a derived one is usually the author's schema
// understating what a variant needs.
func checkDerived(entry catalog.Action, inv recording.Invocation, index int) *Violation {
	required := append([]string{}, entry.Requires...)

	variant := ""
	if entry.Discriminator != "" && len(entry.VariantRequires) > 0 {
		// Selection falls back exactly as template selection does, so a value
		// with no dedicated template inherits the fallback's requirements
		// rather than none. A discriminator that is absent or not a string is
		// already reported by checkVariant; adding a second violation for one
		// mistake would make the diagnostic harder to act on.
		if raw, bound := inv.Kwargs[entry.Discriminator]; bound {
			if value, ok := raw.(string); ok {
				variant = value
				if _, covered := entry.VariantRequires[value]; !covered {
					variant = genpkg.DefaultVariant
				}
				required = append(required, entry.VariantRequires[variant]...)
			}
		}
	}

	var missing []string
	seen := map[string]bool{}
	for _, name := range required {
		if seen[name] {
			continue
		}
		seen[name] = true
		if _, bound := inv.Kwargs[name]; !bound {
			missing = append(missing, name)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	sort.Strings(missing)

	where := "this action's template renders"
	if variant != "" {
		where = fmt.Sprintf("the %q template renders", variant)
	}
	return &Violation{
		Index:  index,
		Action: inv.Action,
		Rule:   RuleMissingDerivedKwarg,
		Detail: fmt.Sprintf("%s %s, so %s must be bound even where the schema declares %s optional",
			where, quoteList(missing),
			plural(len(missing), "it", "they"),
			plural(len(missing), "it", "them")),
	}
}

// checkVariant holds a discriminated action's selecting value to what the
// package can actually render.
//
// A value with no dedicated template is legal when the package ships a
// _default and an error when it does not. The catalog says which is the case,
// so a model that read it should never land here - and one that did is told
// which values are covered rather than that its answer was wrong.
func checkVariant(entry catalog.Action, inv recording.Invocation, index int) []Violation {
	if entry.Discriminator == "" {
		return nil
	}

	raw, bound := inv.Kwargs[entry.Discriminator]
	if !bound {
		// A required discriminator is already reported as missing. An
		// optional one that selects a template is a package-authoring
		// problem, and it is named as such rather than blamed on the model.
		if k, declared := entry.Kwargs[entry.Discriminator]; declared && k.Required {
			return nil
		}
		return []Violation{{
			Index:  index,
			Action: inv.Action,
			Rule:   RuleVariant,
			Detail: fmt.Sprintf(
				"kwarg %q selects this action's template but nothing bound it; bind it to one of %s",
				entry.Discriminator, quoteList(entry.Variants)),
		}}
	}

	value, ok := raw.(string)
	if !ok {
		// The type check has already reported this if the kwarg is declared a
		// string, which it must be to select a template file.
		return nil
	}

	for _, variant := range entry.Variants {
		if variant == value {
			return nil
		}
	}
	if entry.HasDefault {
		return nil
	}

	return []Violation{{
		Index:  index,
		Action: inv.Action,
		Rule:   RuleVariant,
		Detail: fmt.Sprintf(
			"%s %q has no template and this action has no fallback; it covers %s",
			entry.Discriminator, value, quoteList(entry.Variants)),
	}}
}

// checkTargets requires every file the invocation would inject into to be one
// this record authorized and this run created.
//
// It calls Phase 6's resolver rather than rendering injects_into a second time.
// Two implementations would eventually let this accept what expansion rejects,
// and that failure is the one the retry loop can never fix, because the loop is
// never entered (prov-2026-9dcf2658).
//
// This is where a record that named an implementation file but forgot its
// header lands - before a retry is spent, naming the file the author omitted.
func checkTargets(packages []*genpkg.Package, files []resolve.File, inv recording.Invocation, index int) []Violation {
	if _, err := expand.Targets(packages, files, inv); err != nil {
		return []Violation{{
			Index:  index,
			Action: inv.Action,
			Rule:   RuleUnauthorizedPath,
			Detail: err.Error(),
		}}
	}
	return nil
}

// typeOf maps a decoded JSON value onto the closed kwarg type set.
//
// JSON has one number type, so int is the constraint that the number is
// integral rather than that it arrived differently from a float. A value that
// is not one of the four is reported by its absence from this map: there is no
// fifth type an action can declare, so anything else is a violation.
func typeOf(value any) (string, bool) {
	switch v := value.(type) {
	case string:
		return "string", true
	case bool:
		return "bool", true
	case []any:
		return "list", true
	case float64:
		if v != math.Trunc(v) {
			return "float", false
		}
		return "int", true
	default:
		return "", false
	}
}

// describe names what a value actually was, for a diagnostic the model can act
// on. It quotes the value because "bound to a string" is less useful than
// "bound to \"3\"" when the mistake is a number sent as text.
func describe(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("the string %s", strconv.Quote(v))
	case bool:
		return fmt.Sprintf("the bool %v", v)
	case []any:
		return fmt.Sprintf("a list of %d", len(v))
	case float64:
		if v != math.Trunc(v) {
			return fmt.Sprintf("the fractional number %v, which is not an int", v)
		}
		return fmt.Sprintf("the number %v", v)
	case nil:
		return "null"
	case map[string]any:
		return "an object, which is not one of the four kwarg types"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// rejection is what a rejected response is re-prompted with.
//
// It says what was wrong and what to do, and nothing else. It does not restate
// the catalog or the intent: both are already in the conversation, and repeating
// them buries the correction in text the model has already read.
func rejection(violations []Violation) string {
	var b strings.Builder
	b.WriteString("That response was rejected. ")
	fmt.Fprintf(&b, "%s:\n", plural(len(violations), "One check failed", "These checks failed"))

	for _, v := range violations {
		if v.Index > 0 {
			fmt.Fprintf(&b, "\n  - invocation %d (%s): %s", v.Index, v.Action, v.Detail)
			continue
		}
		fmt.Fprintf(&b, "\n  - %s", v.Detail)
	}

	b.WriteString("\n\nReturn the corrected list in the same format. " +
		"Keep the invocations that were not named above exactly as they were.")
	return b.String()
}

func declaredNames(entry catalog.Action) []string {
	out := make([]string, 0, len(entry.Kwargs))
	for name := range entry.Kwargs {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func sortedNames(kwargs map[string]any) []string {
	out := make([]string, 0, len(kwargs))
	for name := range kwargs {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func quoteList(items []string) string {
	if len(items) == 0 {
		return "none"
	}
	quoted := make([]string, 0, len(items))
	for _, item := range items {
		quoted = append(quoted, strconv.Quote(item))
	}
	return strings.Join(quoted, ", ")
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// isEmpty reports whether a bound value carries nothing.
//
// A required kwarg the model mentioned and left empty is not bound. It renders:
// the rails standard's addColumn writes "default: {{default}}" verbatim, and an
// empty default produced "t.string :title, null: false, default:" - Ruby that
// does not parse, from an answer every other Phase 5 check passed
// (prov-2026-9a554c93).
//
// Only string and list are read. Zero is a number and false is a boolean, and
// both are values an author may legitimately want - neither is the model
// declining to answer, which is what this reads for. A string of spaces counts
// as empty for the same reason: it is a value nobody chose that renders as one.
func isEmpty(v any) bool {
	switch value := v.(type) {
	case string:
		return strings.TrimSpace(value) == ""
	case []any:
		return len(value) == 0
	}
	return false
}
