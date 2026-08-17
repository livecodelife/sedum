// Package catalog builds the action catalog: the closed vocabulary the model
// selects from.
//
// One builder serves two callers. Phase 4 embeds a catalog in its prompt, and
// sedum actions prints one for a package author to read. They share this code
// path rather than each assembling a view of their own, because the command's
// whole value is being evidence of what the model receives - if the two could
// drift, the authoring feedback loop lies.
//
// A catalog contains only what the model may act on. An unexposed action is
// absent rather than present-and-rejected, which makes selecting one
// unrepresentable rather than merely invalid. The one exception is sedum
// actions --all, which includes them marked, because an author debugging why an
// action never gets picked needs to see that it is hidden.
package catalog

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/calebcowen/sedum/internal/genpkg"
)

// Kwarg is one argument's declared schema, restated here rather than reused
// from genpkg because this shape is what crosses the wire to the model. The
// two are identical today; the reason to keep them separate is that genpkg's
// shape answers to actions.yaml and this one answers to the prompt, and a
// change to either should not silently become a change to the other.
type Kwarg struct {
	Type     string `json:"type"`
	Required bool   `json:"required"`

	// Description is the author's sentence about what a correct value looks
	// like, passed through untouched.
	//
	// It is here because the catalog is the only thing the model reads, and
	// everything an author writes to explain a kwarg is otherwise written in a
	// YAML comment the model never sees. The type set cannot express "bare
	// names, comma separated" or "the rest of the line after the column name";
	// a sentence can, and getting it wrong is not a failure any Phase 5 check
	// can catch, because the wrong value is a valid string (prov-2026-c5697387).
	Description string `json:"description,omitempty"`

	// Default is what this kwarg resolves to when the model binds nothing.
	//
	// Carried to the model rather than kept internal, because omitting a kwarg
	// is only a choice if the consequence of omitting it is visible. A model
	// that could not see this would be guessing at whether silence is safe -
	// and an absent kwarg in a recording is only interpretable against the
	// declaration that says what absence resolved to (prov-2026-f03916ba).
	//
	// Absent from the catalog when the kwarg declares none, which is the
	// ordinary case.
	Default any `json:"default,omitempty"`
}

// Action is one entry as the model receives it.
//
// Package is carried even though the model does not return it. A record's
// paths may resolve to more than one package, and its catalog is their union,
// so an entry that did not say where it came from would leave the model no way
// to tell that a controller action belongs to the .rb path rather than the .ts
// one beside it.
type Action struct {
	Name    string `json:"action"`
	Package string `json:"package"`

	// Description is the author's sentence about what this action does. Absent
	// means the author wrote none; nothing is synthesised to fill it, because
	// an invented description reads with an authority it has not earned.
	Description string `json:"description,omitempty"`

	Kwargs map[string]Kwarg `json:"kwargs,omitempty"`

	// Discriminator and Variants describe a discriminated action's template
	// selection. Variants is included deliberately: without it there is an
	// invisible cliff, where a declared value gets a dedicated template and
	// anything else falls to _default with no way for the model to know it
	// fell off.
	Discriminator string   `json:"discriminator,omitempty"`
	Variants      []string `json:"variants,omitempty"`

	// HasDefault reports whether the package ships a _default template for
	// this action. It is the other half of what the variant list is for:
	// knowing which values are covered says where the cliff is, and knowing
	// whether a fallback exists says whether stepping off it is survivable
	// or a hard error (prov-2026-21031113).
	HasDefault bool `json:"has_default,omitempty"`

	// Composes names a composite's children, in the order they execute. A
	// composite's own kwarg schema is the union of theirs, so Kwargs above is
	// already everything the caller binds; this is here so the model can see
	// that one selection touches two files rather than guessing from the name.
	Composes []string `json:"composes,omitempty"`

	// InjectsInto is the action's target pattern, exactly as the package
	// author wrote it, one entry per file the action touches - so a composite
	// carries its children's, in execution order.
	//
	// Without it a model has a kwarg named controller and a file named
	// app/controllers/users_controller.rb and nothing connecting them, and it
	// binds the file's path to the kwarg. With it the reasoning is forward and
	// mechanical: match the pattern's literal segments against the authorized
	// file and bind what is left. The model still never chooses a path - the
	// path is the package's - it binds arguments so the package's own pattern
	// lands somewhere the record authorized (prov-2026-1bbb8e2e).
	InjectsInto []string `json:"injects_into,omitempty"`

	// Requires names values this action's template renders on every
	// invocation, beyond whatever the kwarg schema declares required. It is
	// empty for a discriminated action, whose requirements depend on which
	// template the discriminator selects.
	//
	// It is carried separately rather than folded into Kwargs[x].Required so
	// that the schema stays a faithful view of actions.yaml. A reader
	// comparing the two should see what the author wrote; the effective
	// requirement is the union, and stating it as a union is what lets a
	// diagnostic say which half a missing kwarg came from (prov-2026-369544c1).
	Requires []string `json:"requires,omitempty"`

	// VariantRequires names, per variant, the values that variant's template
	// renders. This is the information a shared kwarg schema structurally
	// cannot carry: a kwarg one variant needs and another forbids can only be
	// declared optional, so without this the catalog says "optional" and the
	// model believes it.
	//
	// DefaultVariant appears here when the package ships a fallback, because a
	// discriminator value with no dedicated template inherits the fallback's
	// requirements rather than none.
	VariantRequires map[string][]string `json:"variant_requires,omitempty"`

	// Exposed is nil for an exposed action and false for a hidden one, so
	// that it appears only where it says something. Only sedum actions --all
	// ever produces an entry carrying it; the catalog the model receives has
	// no unexposed entries to mark.
	Exposed *bool `json:"exposed,omitempty"`
}

// Catalog is the whole option set for one record or one package.
type Catalog struct {
	Actions []Action `json:"actions"`
}

// Options controls what a build includes.
type Options struct {
	// IncludeUnexposed adds hidden actions, marked. It is sedum actions --all
	// and nothing else: a catalog bound for a model never sets it, because an
	// action the model cannot legally select has no business in its option
	// set.
	IncludeUnexposed bool
}

// Build assembles the catalog for a set of packages.
//
// Entries are ordered by action name and then by package, so that the same
// packages produce the same bytes every time. That is not cosmetic: the prompt
// embeds this, and an option set that reordered between runs would make one
// model call's input differ from another's for no reason the run could explain.
func Build(packages []*genpkg.Package, opts Options) Catalog {
	var actions []Action
	for _, pkg := range packages {
		for _, action := range pkg.Actions {
			if !action.Exposed && !opts.IncludeUnexposed {
				continue
			}
			actions = append(actions, entry(pkg, action, opts))
		}
	}

	sort.Slice(actions, func(i, j int) bool {
		if actions[i].Name != actions[j].Name {
			return actions[i].Name < actions[j].Name
		}
		return actions[i].Package < actions[j].Package
	})
	return Catalog{Actions: actions}
}

func entry(pkg *genpkg.Package, action *genpkg.Action, opts Options) Action {
	out := Action{
		Name:          action.Name,
		Package:       pkg.Name,
		Description:   action.Description,
		Discriminator: action.Discriminator,
		Variants:      action.Variants,
		Composes:      action.Composes,
	}

	if len(action.Kwargs) > 0 {
		out.Kwargs = make(map[string]Kwarg, len(action.Kwargs))
		for name, k := range action.Kwargs {
			out.Kwargs[name] = Kwarg{
				Type: k.Type, Required: k.Required, Description: k.Description,
				Default: k.Default,
			}
		}
	}

	// A composite has no templates of its own, so it has no _default to
	// report. Its children's fallbacks are theirs, and the model does not
	// select a child.
	if action.Kind() == genpkg.Discriminated {
		_, out.HasDefault = action.Templates[genpkg.DefaultVariant]
	}

	out.InjectsInto = targets(pkg, action)
	out.Requires, out.VariantRequires = requirements(pkg, action)

	if opts.IncludeUnexposed && !action.Exposed {
		hidden := false
		out.Exposed = &hidden
	}
	return out
}

// targets collects the patterns an action's invocation will render.
//
// A composite has none of its own and takes its children's, in execution order,
// because that is what makes one selection visibly touch two files. A child a
// package does not define cannot reach here - load rejects the package - so the
// lookup is a filter rather than a check restated.
func targets(pkg *genpkg.Package, action *genpkg.Action) []string {
	if action.Kind() != genpkg.Composite {
		if action.InjectsInto == "" {
			return nil
		}
		return []string{action.InjectsInto}
	}

	var out []string
	for _, name := range action.Composes {
		child, ok := pkg.Actions[name]
		if !ok || child.InjectsInto == "" {
			continue
		}
		out = append(out, child.InjectsInto)
	}
	return out
}

// requirements collects what an action's templates render, split by whether the
// requirement is unconditional or depends on a selected variant.
//
// A composite has no template of its own and takes its children's, resolved per
// child against the template that child selects - so a simple child contributes
// to the unconditional set and a discriminated one to the per-variant map. A
// child a package does not define cannot reach here, because load rejects the
// package.
func requirements(pkg *genpkg.Package, action *genpkg.Action) ([]string, map[string][]string) {
	unconditional := map[string]bool{}
	byVariant := map[string]map[string]bool{}

	// A value the run supplies is not something an invocation has to bind, so
	// it is not a derived requirement. Without this the check that exists to
	// stop Phase 6 halting on an unbound value rejects every selection of an
	// action whose template references a variable - which is every action the
	// feature exists for (prov-2026-6fc3d13d).
	isVariable := func(name string) bool {
		_, ok := pkg.Variables[name]
		return ok
	}

	collect := func(a *genpkg.Action) {
		for variant, values := range a.TemplateRefs {
			if variant == genpkg.SoleTemplate {
				for _, v := range values {
					if !isVariable(v) {
						unconditional[v] = true
					}
				}
				continue
			}
			if byVariant[variant] == nil {
				byVariant[variant] = map[string]bool{}
			}
			for _, v := range values {
				if !isVariable(v) {
					byVariant[variant][v] = true
				}
			}
		}
	}

	if action.Kind() == genpkg.Composite {
		for _, name := range action.Composes {
			if child, ok := pkg.Actions[name]; ok {
				collect(child)
			}
		}
	} else {
		collect(action)
	}

	var out []string
	for name := range unconditional {
		out = append(out, name)
	}
	sort.Strings(out)

	var variants map[string][]string
	if len(byVariant) > 0 {
		variants = make(map[string][]string, len(byVariant))
		for variant, set := range byVariant {
			var names []string
			for name := range set {
				names = append(names, name)
			}
			sort.Strings(names)
			variants[variant] = names
		}
	}
	return out, variants
}

// Lookup returns every entry declaring a name.
//
// It returns a slice rather than one entry because a record spanning two
// packages that both declare a name has no way to say which is meant - the
// model returns an action and its kwargs, never a package. More than one match
// is that ambiguity, and reporting it is Phase 5's job rather than this one's.
func (c Catalog) Lookup(name string) []Action {
	var out []Action
	for _, a := range c.Actions {
		if a.Name == name {
			out = append(out, a)
		}
	}
	return out
}

// Names returns every action name in the catalog, in order and deduplicated.
// It is what a diagnostic lists when a model names something absent.
func (c Catalog) Names() []string {
	var out []string
	for _, a := range c.Actions {
		if len(out) > 0 && out[len(out)-1] == a.Name {
			continue
		}
		out = append(out, a.Name)
	}
	return out
}

// JSON renders the catalog as the model receives it.
//
// This is the only encoder. sedum actions --json writes exactly these bytes and
// the Phase 4 prompt embeds exactly these bytes, which is what makes the
// command evidence rather than a second implementation that happens to agree
// today.
//
// HTML escaping is off because a template pattern legitimately contains
// characters Go's encoder would otherwise escape, and a prompt carrying
// & where the package said & is harder for a model to read and harder for
// an author to recognise as their own configuration.
func (c Catalog) JSON() ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(c); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}
