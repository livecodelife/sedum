package catalog

import (
	"encoding/json"
	"sort"
)

// ResponseSchema compiles this catalog into a JSON Schema for
// {"invocations": [...]}, for --response-schema (prov-2026-91c54941): an
// endpoint that supports structured output can be asked to constrain the
// model's completion to it, on top of - never instead of - the prompt and
// Phase 5.
//
// It is compiled from exactly what a Catalog already carries - action names,
// kwargs, their declared types, required flags, discriminators, and variants
// - and nothing else, so that a caller compiling ResponseSchema over
// FilterForRecord's output gets a schema scoped to what survived the filter,
// the same way the prompt already is. Description is deliberately absent:
// it exists to be read by the model, not enforced against it, and the wire
// contract a validator holds the model to is the type, not the prose
// (prov-2026-c5697387).
//
// Every property this method can emit - type, const, properties, required,
// additionalProperties, items, oneOf - is chosen so the array itself carries
// no minItems or other non-emptiness constraint: prov-2026-4bcabb2f found
// that exact family of constraint collapses selection, and {"invocations":
// []} must stay legal here exactly as it is without --response-schema.
//
// Deterministic by construction rather than by care taken at each call site:
// every collection below is either already ordered (c.Actions, by Build's own
// sort) or sorted before being written, and Go's encoding/json sorts map
// string keys when marshaling, so the same catalog produces byte-identical
// output on every call.
func (c Catalog) ResponseSchema() ([]byte, error) {
	schema := map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"invocations"},
		"properties": map[string]any{
			"invocations": map[string]any{
				"type": "array",
				"items": map[string]any{
					"oneOf": actionBranches(c),
				},
			},
		},
	}
	return json.Marshal(schema)
}

// actionBranches is one oneOf branch per action, and per declared variant for
// a discriminated one - in c.Actions' own order, which Build already makes
// deterministic (name, then package), with each discriminated action's
// variants in the order the package author declared them.
func actionBranches(c Catalog) []any {
	branches := make([]any, 0, len(c.Actions))
	for _, action := range c.Actions {
		variants := action.Variants
		if len(variants) == 0 {
			variants = []string{""}
		}
		for _, variant := range variants {
			branches = append(branches, actionBranch(action, variant))
		}
	}
	return branches
}

// actionBranch is the schema for one {action, kwargs} invocation naming this
// action - pinned by action as a const, exactly as required for oneOf to stay
// exclusive between branches that would otherwise both accept a value neither
// declares (additionalProperties: false forbids an undeclared property, not
// one another branch declares and this one's instance simply omits).
func actionBranch(action Action, variant string) map[string]any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"action", "kwargs"},
		"properties": map[string]any{
			"action": map[string]any{"const": action.Name},
			"kwargs": kwargsSchema(action, variant),
		},
	}
}

// kwargsSchema is one branch's kwargs object: every declared kwarg mapped to
// its wire type, the discriminator (when this branch is a declared variant's)
// pinned by const instead of typed, and required computed from the one shared
// source both this compiler and Phase 5's derivedRequirements read from -
// Action.RequiredForVariant - unioned with the kwargs the schema itself
// declares required.
func kwargsSchema(action Action, variant string) map[string]any {
	props := make(map[string]any, len(action.Kwargs))
	for name, k := range action.Kwargs {
		if variant != "" && name == action.Discriminator {
			props[name] = map[string]any{"const": variant}
			continue
		}
		props[name] = kwargTypeSchema(k)
	}

	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             kwargsRequired(action, variant),
		"properties":           props,
	}
}

// kwargsRequired is the property set a branch must bind: every kwarg the
// schema itself declares required, unioned with what RequiredForVariant says
// this action's (and, for a discriminated action, this variant's) template
// actually renders - the same union Phase 5's derivedRequirements computes
// from the same fields, so the two cannot drift apart into independently
// maintained copies of the rule (prov-2026-91c54941).
func kwargsRequired(action Action, variant string) []string {
	set := make(map[string]bool, len(action.Kwargs))
	for name, k := range action.Kwargs {
		if k.Required {
			set[name] = true
		}
	}
	for _, name := range action.RequiredForVariant(variant) {
		set[name] = true
	}

	out := make([]string, 0, len(set))
	for name := range set {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// kwargTypeSchema maps a kwarg's declared type onto its JSON Schema shape.
//
// A kwarg is never modeled as nullable, whatever the schema calls optional:
// Phase 5 rejects a null value for an absent optional kwarg
// (prov-2026-9a554c93/prov-2026-9a491128), so a schema that allowed null
// would legalize an output shape the validator throws out anyway, at the
// cost of a wasted retry instead of a clean rejection here.
//
// literal carries on the wire as a JSON string - genpkg.TypeSatisfiedBy
// draws the same distinction for the same reason - and list maps to an array
// of strings, per the record's own stated shape; nothing here inspects what a
// kwarg's bound value would mean in any target language.
func kwargTypeSchema(k Kwarg) map[string]any {
	switch k.Type {
	case "list":
		return map[string]any{"type": "array", "items": map[string]any{"type": "string"}}
	case "int":
		return map[string]any{"type": "integer"}
	case "bool":
		return map[string]any{"type": "boolean"}
	default:
		// "string" and "literal" both carry a JSON string on the wire.
		return map[string]any{"type": "string"}
	}
}
