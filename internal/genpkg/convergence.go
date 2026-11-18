package genpkg

import (
	"fmt"
	"sort"
	"strings"

	"github.com/livecodelife/sedum/internal/render"
)

// ConvergencePlaceholder is the text every seeded, not-yet-converged helper
// carries. A consumer (graft) recognises an unconverged body by it, so it is
// fixed by the core rather than chosen per package: a package that seeded a
// different string would be invisible to that check.
const ConvergencePlaceholder = "graft: not converged"

// RuleConvergenceInvalid is the finding for a malformed convergence block.
const RuleConvergenceInvalid = "convergence_invalid"

// Convergence is the optional `convergence:` block of sedum.yaml. It tells a
// consumer how this language writes a stub that delegates to a helper, and how
// a seeded file of placeholder helpers looks.
//
// Sedum's core knows no language. It renders these templates with the same
// engine as every other template and supplies, as opaque strings, everything
// the grammar (no conditionals, no loops) cannot compute: parameter lists are
// pre-joined, a result list is rendered by arity, and the func-or-method
// choice is made by picking the variant. Nothing here is parsed as code.
type Convergence struct {
	// SynthSuffix is inserted between a file's stem and extension to name the
	// file that holds its helpers: dir/<stem><suffix><ext>.
	SynthSuffix string `yaml:"synth_suffix"`

	// Helper names the helper a stub delegates to. Values: stem, name, Name
	// (name with its first letter upper-cased), dir, ext.
	Helper string `yaml:"helper"`

	// Param renders one parameter (values: name, type). Arg renders it as a
	// call argument (name). SpreadPrefix, when set, marks a type that spreads
	// (Go's "..."); such a parameter is passed with SpreadArg (name) instead.
	Param        string `yaml:"param"`
	Arg          string `yaml:"arg"`
	SpreadPrefix string `yaml:"spread_prefix"`
	SpreadArg    string `yaml:"spread_arg"`

	// Result renders a signature's results by arity. None has no values; One
	// has result; Many has results, pre-joined with ", " (the grammar's join
	// trims its separator, so the core joins).
	Result ConvergenceResult `yaml:"result"`

	// Call renders the delegating call (values: helper, args, synth_stem, and
	// for a method receiver_name and recv_args, the receiver then the args). Body wraps it by whether the stub returns
	// anything (value: call).
	Call ConvergenceVariants `yaml:"call"`
	Body ConvergenceBody     `yaml:"body"`

	// ReceiverParam renders the receiver as a leading helper parameter (values:
	// receiver_name, receiver_type); ReceiverArg renders it as a leading call
	// argument (value: receiver_name). They let a language whose helpers live in
	// another module take the receiver explicitly. SynthImport renders the module
	// specifier a stub file must import to reach its helpers (values: synth_stem,
	// stem, dir, ext); empty when the helpers share the stub's namespace.
	ReceiverParam string `yaml:"receiver_param"`
	ReceiverArg   string `yaml:"receiver_arg"`
	SynthImport   string `yaml:"synth_import"`

	// Declare renders a helper's signature without a body (values: helper,
	// params, result, and for a method receiver_name, receiver_type and
	// recv_params, the receiver parameter then the params).
	Declare ConvergenceVariants `yaml:"declare"`

	// Note is the generated line of a stub's doc comment (values: name,
	// helper). Doc comments are written with the package's comment prefix.
	Note string `yaml:"note"`

	// Stub renders the delegating function or method, doc comment included
	// (values: name, params, result, body, doc, helper, and receiver_name,
	// receiver_type for a method).
	Stub ConvergenceVariants `yaml:"stub"`

	// ImportAlias says the language binds an import to a name: each import then
	// carries an alias (value alias) derived from its path (last element after
	// "/" or ":", extension dropped, other characters "_"), which the addImport
	// kwargs and the Import template can use.
	ImportAlias bool `yaml:"import_alias"`

	// MethodAction is the catalog action that writes a method stub (default
	// setMethod, which also takes a receiver kwarg). A language that has no
	// method declarations outside a class sets setFunction and MethodName, the
	// name the stub is declared under (values: receiver_type, name).
	MethodAction string `yaml:"method_action"`
	MethodName   string `yaml:"method_name"`

	// Import renders one import line (values: path, alias). ImportUse, when set, is
	// the text an import must be used as to be kept in the seed (values: base,
	// the last path segment, and alias): an import no helper signature mentions would
	// not compile in some languages. Unset keeps every import.
	Import    string `yaml:"import"`
	ImportUse string `yaml:"import_use"`

	// SeedDecl renders one placeholder helper (values: declare, placeholder)
	// and must contain the placeholder. Seed renders the whole synth file
	// (values: package, imports, decls, file, stem, dir); imports and decls
	// are pre-joined. Runs of blank lines are collapsed to one, so a template
	// can leave an empty import list in place.
	SeedDecl string `yaml:"seed_decl"`
	Seed     string `yaml:"seed"`

	// Package names the language-level package of a file that does not carry
	// one. Entry is used for a file whose directory matches an EntryDirs
	// pattern (path.Match); Default is a template over dirname (the last
	// directory element). Both are optional; the seed may not use package.
	Package ConvergencePackage `yaml:"package"`
}

// ConvergenceResult is the result-list templates by arity.
type ConvergenceResult struct {
	None string `yaml:"none"`
	One  string `yaml:"one"`
	Many string `yaml:"many"`
}

// ConvergenceVariants holds a template for a free function and one for a
// method, because the grammar cannot choose between them.
type ConvergenceVariants struct {
	Func   string `yaml:"func"`
	Method string `yaml:"method"`
}

// ConvergenceBody wraps the call for a stub with and without a result.
type ConvergenceBody struct {
	Void  string `yaml:"void"`
	Value string `yaml:"value"`
}

// ConvergencePackage is how a file's package name is derived.
type ConvergencePackage struct {
	Entry     string   `yaml:"entry"`
	EntryDirs []string `yaml:"entry_dirs"`
	Default   string   `yaml:"default"`
}

type convTemplate struct {
	label  string
	src    string
	values []string // names the template may reference
}

// templates lists every template with the values it may reference.
func (c *Convergence) templates() []convTemplate {
	hv := []string{"stem", "name", "Name", "dir", "ext"}
	sig := []string{"helper", "params", "result"}
	recv := []string{"receiver_name", "receiver_type", "recv_params"}
	cat := func(parts ...[]string) []string {
		var out []string
		for _, p := range parts {
			out = append(out, p...)
		}
		return out
	}
	return []convTemplate{
		{"helper", c.Helper, hv},
		{"param", c.Param, []string{"name", "type"}},
		{"arg", c.Arg, []string{"name"}},
		{"spread_arg", c.SpreadArg, []string{"name"}},
		{"result.none", c.Result.None, nil},
		{"result.one", c.Result.One, []string{"result"}},
		{"result.many", c.Result.Many, []string{"results"}},
		{"receiver_param", c.ReceiverParam, []string{"receiver_name", "receiver_type"}},
		{"receiver_arg", c.ReceiverArg, []string{"receiver_name"}},
		{"synth_import", c.SynthImport, []string{"synth_stem", "stem", "dir", "ext"}},
		{"call.func", c.Call.Func, []string{"helper", "args", "synth_stem"}},
		{"call.method", c.Call.Method, []string{"helper", "args", "synth_stem", "receiver_name", "recv_args"}},
		{"body.void", c.Body.Void, []string{"call"}},
		{"body.value", c.Body.Value, []string{"call"}},
		{"declare.func", c.Declare.Func, sig},
		{"declare.method", c.Declare.Method, cat(sig, recv)},
		{"note", c.Note, []string{"name", "helper"}},
		{"stub.func", c.Stub.Func, []string{"name", "params", "result", "body", "doc", "helper"}},
		{"stub.method", c.Stub.Method, []string{"name", "params", "result", "body", "doc", "helper", "receiver_name", "receiver_type", "recv_params"}},
		{"import", c.Import, []string{"path", "alias"}},
		{"import_use", c.ImportUse, []string{"base", "alias"}},
		{"method_name", c.MethodName, []string{"receiver_type", "name"}},
		{"seed_decl", c.SeedDecl, []string{"declare", "placeholder"}},
		{"seed", c.Seed, []string{"package", "imports", "decls", "file", "stem", "dir", "ext"}},
		{"package.default", c.Package.Default, []string{"dirname"}},
	}
}

// optionalConvTemplates may be empty.
var optionalConvTemplates = map[string]bool{
	"spread_arg": true, "import_use": true, "package.default": true,
	"result.none": true, "synth_import": true, "method_name": true, "receiver_param": true, "receiver_arg": true,
}

// checkConvergence validates the block: every required template present,
// compiling against the package's own transforms, referencing only values the
// core supplies, and a seed placeholder that carries the marker string.
func checkConvergence(pkg *Package, r *reporter) {
	c := pkg.Convergence
	if c == nil {
		return
	}
	bad := func(format string, args ...any) {
		r.errorf(manifestFile, RuleConvergenceInvalid, "convergence: "+format, args...)
	}
	if pkg.Engine == nil {
		return // the transform table already failed; nothing can compile
	}
	if c.SynthSuffix == "" {
		bad("synth_suffix is required")
	}
	if c.SpreadPrefix != "" && c.SpreadArg == "" {
		bad("spread_prefix is set, so spread_arg is required")
	}
	for _, t := range c.templates() {
		if t.src == "" {
			if !optionalConvTemplates[t.label] {
				bad("%s is required", t.label)
			}
			continue
		}
		tmpl, err := render.Compile(pkg.Engine, t.src)
		if err != nil {
			bad("%s does not compile: %v", t.label, err)
			continue
		}
		allowed := map[string]bool{}
		for _, v := range t.values {
			allowed[v] = true
		}
		for _, v := range tmpl.Values() {
			if !allowed[v] {
				bad("%s references %q, which is not one of its values (%s)",
					t.label, v, strings.Join(sortedStrings(t.values), ", "))
			}
		}
	}
	if !strings.Contains(c.SeedDecl, ConvergencePlaceholder) && !strings.Contains(c.SeedDecl, "{{placeholder") {
		bad("seed_decl must contain %q (or reference {{placeholder}}), or a converged body cannot be told from a seeded one", ConvergencePlaceholder)
	}
	if !strings.Contains(c.Seed, "{{decls") {
		bad("seed must reference {{decls}}")
	}
	switch c.MethodAction {
	case "", "setMethod", "setFunction":
	default:
		bad("method_action %q must be setMethod or setFunction", c.MethodAction)
	}
	if c.MethodAction == "setFunction" && c.MethodName == "" {
		bad("method_action setFunction needs method_name, or methods of different receivers would collide")
	}
	for _, p := range c.Package.EntryDirs {
		if p == "" {
			bad("package.entry_dirs has an empty pattern")
		}
	}
}

func sortedStrings(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}

// ConvergenceInfo is what a consumer needs to lay out files before it renders
// anything.
type ConvergenceInfo struct {
	Package       string   `json:"package"`
	Extensions    []string `json:"extensions"`
	CommentPrefix string   `json:"comment_prefix"`
	SynthSuffix   string   `json:"synth_suffix"`
	Placeholder   string   `json:"placeholder"`
}

// ConvergenceInfo reports the package's convergence facts, or an error when it
// declares no convergence block.
func (p *Package) ConvergenceInfo() (ConvergenceInfo, error) {
	if p.Convergence == nil {
		return ConvergenceInfo{}, fmt.Errorf("package %s declares no convergence block", p.Name)
	}
	return ConvergenceInfo{
		Package:       p.Name,
		Extensions:    p.Extensions,
		CommentPrefix: p.CommentPrefix,
		SynthSuffix:   p.Convergence.SynthSuffix,
		Placeholder:   ConvergencePlaceholder,
	}, nil
}
