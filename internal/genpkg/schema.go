package genpkg

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/calebcowen/sedum/internal/transform"
)

// The YAML shapes of a generator package's two declaration files, and the
// resolved types loading produces from them.
//
// Decoding is strict: an unrecognized key is an error naming the key and its
// file, not a silently ignored field. A typo in a key is otherwise invisible
// until generation, at which point it looks like a missing feature rather than
// a misspelling.

// manifest is sedum.yaml.
type manifest struct {
	Name          string                       `yaml:"name"`
	Extensions    []string                     `yaml:"extensions"`
	CommentPrefix string                       `yaml:"comment_prefix"`
	Transforms    map[string][]string          `yaml:"transforms"`
	OpExceptions  map[string]map[string]string `yaml:"op_exceptions"`
	Unmanaged     []string                     `yaml:"unmanaged"`
	Variables     map[string]Variable          `yaml:"variables"`
}

// Variable is a project fact this package's templates reference and cannot
// know: a Go module path, a Java group id, a .NET root namespace.
//
// The package declares the name because which facts a standard needs is the
// standard's business. The run supplies the value because whose project it is
// is the run's. That split is what keeps a package reusable across every repo
// that follows its conventions (prov-2026-6fc3d13d).
//
// Sedum never interprets one. It substitutes text under a name and does not
// read the name, parse the value, or attach a meaning to either - which is what
// keeps a module path from becoming something the core knows about Go.
type Variable struct {
	// Description is what a correct value looks like, addressed to the person
	// running the command. That is the opposite of a kwarg description, which
	// is addressed to the model (prov-2026-c5697387): a variable never reaches
	// the catalog, so nothing but a human ever reads this.
	Description string `yaml:"description"`

	// Default is the value used when the run supplies none. Absent means the
	// run must supply one, and a run that does not is halted before Phase 3.
	//
	// A pointer because a variable's value is a string, so an absent default and
	// a declared empty one are the same string and only the pointer tells them
	// apart. Kwarg.Default gets this for free by being `any`.
	Default *string `yaml:"default"`
}

// HasDefault reports whether this variable declares a value to fall back to.
func (v Variable) HasDefault() bool { return v.Default != nil }

// Value is the default this variable declares, or the empty string when it
// declares none. Callers check HasDefault first.
func (v Variable) Value() string {
	if v.Default == nil {
		return ""
	}
	return *v.Default
}

// actionsFile is actions/actions.yaml.
type actionsFile struct {
	Actions map[string]*actionDecl `yaml:"actions"`
}

// actionDecl is one entry under actions:. Every field an action will ever carry
// is modeled here, including the ones only later milestones read, because
// strict decoding rejects any key the schema omits (prov-2026-d1d61186).
//
// Exposed is a pointer so that "absent" is distinguishable from "false".
// exposed defaults to true: authoring an action is enough to make it usable,
// and hiding is the deliberate act.
type actionDecl struct {
	Kwargs        map[string]Kwarg `yaml:"kwargs"`
	Description   string           `yaml:"description"`
	Discriminator string           `yaml:"discriminator"`
	Variants      []string         `yaml:"variants"`
	InjectsInto   string           `yaml:"injects_into"`
	Anchor        string           `yaml:"anchor"`
	AnchorStart   string           `yaml:"anchor_start"`
	AnchorEnd     string           `yaml:"anchor_end"`
	AnchorPattern string           `yaml:"anchor_pattern"`
	Composes      []string         `yaml:"composes"`
	Exposed       *bool            `yaml:"exposed"`
}

// Kwarg is one entry in an action's argument schema - the shape the model is
// held to when it binds arguments.
type Kwarg struct {
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`

	// Description is the author's sentence about what a correct value looks
	// like. It exists because the type set is deliberately tiny - string, int,
	// bool, list - so it cannot express "a comma-separated list of bare names"
	// and should not try, while an author can say it in a sentence.
	//
	// Nothing interprets it. Sedum does not parse a description, validate one,
	// derive a constraint from one, or check that a bound value agrees with
	// one. A description must never become a place where a rule lives, because
	// a rule the model can read and Phase 5 cannot enforce is worse than no
	// rule at all (prov-2026-c5697387).
	Description string `yaml:"description"`

	// Default is the value this kwarg takes when nothing binds it.
	//
	// It exists for arguments whose correct answer is a property of the package
	// rather than of the change: a Rails column that carries no default is
	// written `default: nil`, which is how Rails spells absence and not a
	// decision about any particular column. Requiring the model to write it
	// requires the model to know Ruby, which it measurably does not
	// (prov-2026-f03916ba).
	//
	// Only meaningful on a kwarg that is not required - a required kwarg's
	// default could never be reached, so declaring both is a load error rather
	// than a precedence rule.
	//
	// nil means no default was declared. An empty default is still a value -
	// `default: ""` decodes to an interface holding the empty string, which is
	// not nil - so the two are distinguishable without a second field.
	Default any `yaml:"default"`
}

// HasDefault reports whether this kwarg declares a value to use when nothing
// binds it.
//
// A YAML `default: null` reads as no default, which is the right reading: there
// is no null kwarg type, so a null default could never satisfy the type check
// it would have to pass.
func (k Kwarg) HasDefault() bool { return k.Default != nil }

// KwargTypes is the closed set an action's kwarg may declare. It is
// deliberately small: sufficient for argument binding and nothing more.
//
// literal is a string the template emits verbatim as source in the target
// language - a Ruby default, a SQL set clause, a Go scan target. It is carried
// on the wire as a JSON string and rendered exactly like one; nothing in Sedum
// parses it, quotes it, or knows what language it is written in, and adding a
// type that did would put a per-language rendering table in the core the PRD
// keeps language-free.
//
// It exists because "string" was a true statement that told the model the wrong
// thing. A kwarg whose value is prose and one whose value is code are bound
// differently, and declaring both "string" left the difference sayable only in
// a description - which is where it lived, in the one package that wrote one.
var KwargTypes = []string{"string", "int", "bool", "list", "literal"}

func validKwargType(t string) bool { return slices.Contains(KwargTypes, t) }

// TypeSatisfiedBy reports whether a value whose concrete type is got may bind a
// kwarg declared as declared.
//
// Every declared type but one is its own carried type. A literal is source code
// and neither JSON nor YAML can carry code, so it arrives as a string - which
// makes this the one place a declared type and a carried type differ.
//
// Defined here, beside the type set, rather than at either of the two places
// that ask. Phase 5 asks about a JSON value from the model and load asks about
// a YAML value from a package author; if each spelled the rule itself, a third
// type with the same property would be added to one and not the other.
func TypeSatisfiedBy(declared, got string) bool {
	if declared == "literal" {
		return got == "string"
	}
	return declared == got
}

// ActionKind is how an action is realized. It is read from the schema and never
// inferred from the filesystem.
type ActionKind int

const (
	// Simple has one template file named for the action.
	Simple ActionKind = iota
	// Discriminated has a directory named for the action holding one
	// template per declared variant plus an optional _default.
	Discriminated
	// Composite has no template at all and triggers no filesystem lookup.
	Composite
)

func (k ActionKind) String() string {
	switch k {
	case Discriminated:
		return "discriminated"
	case Composite:
		return "composite"
	default:
		return "simple"
	}
}

// DefaultVariant names the fallback template in a discriminated action's
// directory. It is optional; a discriminator value that is not a declared
// variant falls to it when present.
const DefaultVariant = "_default"

// Action is a loaded, resolved action definition.
type Action struct {
	Name string

	// Kwargs is the declared schema for a simple action, and the union of
	// its children's schemas for a composite - union of names, union of
	// required flags.
	Kwargs map[string]Kwarg

	// Description is the author's sentence about what this action does,
	// carried to the catalog untouched. Empty means the author wrote none,
	// and nothing is invented to fill it - a synthesised description is worse
	// than none, because it reads with the same authority and carries none.
	Description string

	Discriminator string
	Variants      []string
	InjectsInto   string
	Anchor        string
	AnchorStart   string
	AnchorEnd     string
	AnchorPattern string
	Composes      []string
	Exposed       bool

	// Template is the path to a simple action's template, relative to the
	// package directory. Empty for the other kinds.
	Template string
	// Templates maps variant name to template path for a discriminated
	// action, including DefaultVariant when present. Empty for the others.
	Templates map[string]string

	// TemplateRefs are the values each of this action's templates renders,
	// sorted and deduplicated, keyed by variant for a discriminated action
	// and by SoleTemplate for a simple one. A composite has none of its own;
	// its children carry theirs.
	//
	// A value a template renders is a value that template needs. That holds
	// exactly rather than approximately, because the grammar is {{name}} and
	// {{name|op}} with no conditionals, loops, or field access - so there is
	// no shape of template for which "referenced" and "required" differ. If
	// the grammar ever gains one, this stops being a derivation and becomes a
	// guess (prov-2026-369544c1).
	TemplateRefs map[string][]string

	kind ActionKind
}

// SoleTemplate keys the single entry in a simple action's TemplateRefs. A
// simple action has one template and no variant to name it by, and giving it a
// key rather than a field of its own lets one lookup serve both kinds.
const SoleTemplate = ""

func (a *Action) Kind() ActionKind { return a.kind }

// Requires returns the values the template selected by variant renders.
//
// The variant is SoleTemplate for a simple action and a discriminator value
// for a discriminated one, falling back to DefaultVariant exactly as template
// selection does - so a value with no dedicated template inherits the
// fallback's requirements rather than none.
func (a *Action) Requires(variant string) []string {
	if a.TemplateRefs == nil {
		return nil
	}
	if refs, ok := a.TemplateRefs[variant]; ok {
		return refs
	}
	return a.TemplateRefs[DefaultVariant]
}

// Package is a loaded generator package: a team's conventions for one target
// stack.
type Package struct {
	Name string
	// Dir is the package directory as it was found on disk.
	Dir           string
	Extensions    []string
	CommentPrefix string
	// Transforms are the named pipelines the package declares over the
	// built-in operation set.
	Transforms   map[string][]string
	OpExceptions map[string]map[string]string

	// Unmanaged are the paths this package declares Sedum will not create,
	// render, or inject into, in the grammar scope entries already use. It
	// authorizes nothing: a record still decides what may be touched, and
	// this says only that Sedum is not what touches it. The usual reason is
	// a file a person or another tool owns - a Gemfile, a lockfile, a key.
	Unmanaged []string

	// Variables are the project facts this package's templates reference and
	// cannot know. The run supplies their values (prov-2026-6fc3d13d).
	Variables map[string]Variable
	// Engine applies this package's transforms. It is built from the two
	// fields above at load, so that a pipeline which cannot be built rejects
	// the package instead of failing mid-render (prov-2026-4675cebe). It is
	// nil only for a package that was rejected.
	Engine  *transform.Engine
	Actions map[string]*Action

	// FileTemplates are the path patterns under files/, slash-normalized and
	// relative to files/. Each template's own path is its pattern. They are
	// not part of the action catalog and are never loaded as actions.
	FileTemplates []string

	// fileContents holds what each of those templates contains, keyed by
	// pattern. Loading has already read every one of them to check the
	// transforms they reference, and Phase 3 renders them; re-reading at
	// generation time would mean a package could change shape between the
	// point it was declared valid and the point it was used.
	fileContents map[string]string

	// actionContents holds what each action template contains, keyed by its
	// path relative to the package directory - the same key Action.Template
	// and Action.Templates carry. It is kept for the same reason
	// fileContents is: Phase 6 renders these, loading has already read them,
	// and re-reading would reopen the window between validation and use.
	actionContents map[string]string
}

// FileTemplate returns the contents of the file template with the given
// pattern.
func (p *Package) FileTemplate(pattern string) (string, bool) {
	content, ok := p.fileContents[pattern]
	return content, ok
}

// ActionTemplate returns the contents of the action template at the given path,
// which is what Action.Template holds for a simple action and what an entry of
// Action.Templates holds for a discriminated one.
func (p *Package) ActionTemplate(path string) (string, bool) {
	content, ok := p.actionContents[path]
	return content, ok
}

// Exposed returns the actions a catalog may show, in no particular order.
// Only exposed actions ever reach the model.
func (p *Package) Exposed() []*Action {
	var out []*Action
	for _, a := range p.Actions {
		if a.Exposed {
			out = append(out, a)
		}
	}
	return out
}

// ResolveVariables binds every variable the packages declare, from the values a
// run supplied and the defaults they declare.
//
// Two failures, both before Phase 3 touches disk. A value for a variable no
// package declares is an error rather than an ignored key: a run cannot invent a
// variable, so a name nothing declares is a typo, and a typo silently ignored
// renders the template with nothing bound. A declared variable with neither a
// value nor a default is an error naming the package and printing the author's
// description, which is what the description is for (prov-2026-6fc3d13d).
func (s *Set) ResolveVariables(supplied map[string]string) (map[string]string, error) {
	declared := map[string]Variable{}
	owner := map[string]string{}
	for _, pkg := range s.Packages {
		for name, v := range pkg.Variables {
			declared[name] = v
			owner[name] = pkg.Name
		}
	}

	var unknown []string
	for name := range supplied {
		if _, ok := declared[name]; !ok {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		known := sortedKeys(declared)
		if len(known) == 0 {
			return nil, fmt.Errorf("no loaded package declares any variable, so there is nothing for %s to bind",
				quoteAll(unknown))
		}
		return nil, fmt.Errorf("no loaded package declares %s; the declared variables are %s",
			quoteAll(unknown), quoteAll(known))
	}

	out := make(map[string]string, len(declared))
	var missing []string
	for _, name := range sortedKeys(declared) {
		if value, ok := supplied[name]; ok {
			out[name] = value
			continue
		}
		if declared[name].HasDefault() {
			out[name] = declared[name].Value()
			continue
		}
		missing = append(missing, fmt.Sprintf("  %s (%s): %s", name, owner[name], described(declared[name])))
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("%s not set and %s no default:\n%s",
			plural(len(missing), "variable is", "variables are"),
			plural(len(missing), "declares", "declare"),
			strings.Join(missing, "\n"))
	}
	return out, nil
}

func described(v Variable) string {
	if strings.TrimSpace(v.Description) == "" {
		return "the package describes it no further"
	}
	return strings.TrimSpace(v.Description)
}

func quoteAll(names []string) string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = strconv.Quote(n)
	}
	return strings.Join(out, ", ")
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
