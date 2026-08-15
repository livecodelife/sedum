package genpkg

import (
	"slices"

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
}

// KwargTypes is the closed set an action's kwarg may declare. It is
// deliberately small: sufficient for argument binding and nothing more.
var KwargTypes = []string{"string", "int", "bool", "list"}

func validKwargType(t string) bool { return slices.Contains(KwargTypes, t) }

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
	Kwargs        map[string]Kwarg
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
