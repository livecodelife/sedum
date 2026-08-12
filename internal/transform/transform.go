// Package transform is the transform engine: the closed set of built-in
// operations Sedum ships, the named pipelines a generator package composes from
// them, and the per-package exception tables that handle irregular mappings.
//
// The operation set is fixed and language-neutral. What a target's conventions
// mean is expressed by composing those operations into pipelines in a package's
// sedum.yaml - "constantize" means whatever a package says it means - so a team
// adding a language writes configuration rather than waiting for Sedum to grow
// an operation.
//
// Every operation is pure: its output depends on its input, its literal
// argument, and the package configuration the engine was built from, and on
// nothing else. Arguments are string literals only. Dynamic arguments would
// begin the construction of an expression language, which is a different
// product.
//
// This package is a leaf. It imports nothing else of Sedum's, so that the
// engine cannot come to depend on how packages are laid out or loaded.
package transform

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
	"text/template"
)

// Operations is the closed set of built-in operations, always available in
// every package. This is the single declaration of the vocabulary: the package
// loader checks references against it and the engine below implements exactly
// these names (prov-2026-73127e53).
var Operations = []string{
	"pascal", "camel", "snake", "kebab",
	"upper", "lower",
	"plural", "singular",
	"prefix", "suffix",
}

// parameterized operations take a string literal after a colon: prefix:@,
// suffix:_path.
var parameterized = map[string]bool{"prefix": true, "suffix": true}

// Parameterized reports whether an operation takes a literal argument.
func Parameterized(op string) bool { return parameterized[op] }

// exceptionOperations are the operations that consult a package's
// op_exceptions table.
//
// Only pascal and camel split their input into words, and a per-word
// substitution needs words to substitute into. Declaring a table for any other
// operation is an error rather than a table that quietly does nothing
// (prov-2026-6feb227b).
var exceptionOperations = []string{"pascal", "camel"}

// Ref is one transform reference as an author writes it: a name, and for a
// parameterized operation the literal argument after the colon.
type Ref struct {
	Name string
	Arg  string
	// HasArg distinguishes "prefix:" from "prefix" - an empty argument is a
	// mistake, and an absent one is a different mistake.
	HasArg bool
}

// ParseRef splits a reference at its first colon. The argument is opaque text:
// it is a string literal, so there is nothing in it to resolve.
func ParseRef(s string) Ref {
	s = strings.TrimSpace(s)
	if name, arg, ok := strings.Cut(s, ":"); ok {
		return Ref{Name: strings.TrimSpace(name), Arg: strings.TrimSpace(arg), HasArg: true}
	}
	return Ref{Name: s}
}

func (r Ref) String() string {
	if r.HasArg {
		return r.Name + ":" + r.Arg
	}
	return r.Name
}

// Value is a bound value that may carry explicit inflected forms.
//
// Since the model is already binding arguments, it may return the forms where
// they matter rather than leaving morphology to a rule table. A supplied form
// overrides the table for that value; the table is the default, not the
// authority.
type Value struct {
	Text     string
	Plural   string
	Singular string
}

// String makes a bare reference to a Value render as the text the model bound.
func (v Value) String() string { return v.Text }

// Config is one package's transform configuration, as declared in its
// sedum.yaml.
type Config struct {
	// Pipelines are named compositions of built-in operations:
	// instantize: [plural, "prefix:@"].
	Pipelines map[string][]string
	// Exceptions are per-operation irregular mappings, applied per word.
	Exceptions map[string]map[string]string
	// Language selects the inflection rule table. Empty means
	// DefaultLanguage. Tables are data files, not code paths.
	Language string
}

// Engine applies transforms under one package's conventions.
//
// Everything it needs is held on the engine, never in a process-global. Two
// packages loaded in one run must be able to disagree about how a token
// capitalizes, which is also why strcase.ConfigureAcronym is never called - it
// writes to a package-level map, so the last package to configure would win
// process-wide.
type Engine struct {
	pipelines  map[string][]Ref
	exceptions map[string]map[string]string
	inflector  *inflector
}

// New builds an engine from a package's declared configuration.
//
// Every problem found is reported rather than the first, because an author
// fixes a package in one pass. A returned error means the package is rejected:
// there is no engine that half-works.
func New(cfg Config) (*Engine, error) {
	e := &Engine{
		pipelines:  map[string][]Ref{},
		exceptions: map[string]map[string]string{},
	}

	inf, err := loadInflector(cfg.Language)
	if err != nil {
		return nil, err
	}
	e.inflector = inf

	var problems []error

	for _, name := range sortedKeys(cfg.Exceptions) {
		table := cfg.Exceptions[name]
		switch {
		case !slices.Contains(Operations, name):
			problems = append(problems, fmt.Errorf(
				"op_exceptions declares a table for %q, which is not a built-in operation (%s)",
				name, strings.Join(Operations, ", ")))
			continue
		case !slices.Contains(exceptionOperations, name):
			problems = append(problems, fmt.Errorf(
				"op_exceptions declares a table for %q, which does not consult one; only %s split their input into words, so only they can substitute one",
				name, strings.Join(exceptionOperations, " and ")))
			continue
		}

		// Keys are normalized because the lookup happens after the input
		// has been normalized, so URL and url name the same word.
		normalized := make(map[string]string, len(table))
		for word, replacement := range table {
			normalized[strings.ToLower(word)] = replacement
		}
		e.exceptions[name] = normalized
	}

	for _, name := range sortedKeys(cfg.Pipelines) {
		steps := cfg.Pipelines[name]
		if slices.Contains(Operations, name) {
			problems = append(problems, fmt.Errorf(
				"pipeline %q shadows the built-in operation of the same name; a pipeline that redefines an operation makes every template that uses it ambiguous", name))
			continue
		}
		if len(steps) == 0 {
			problems = append(problems, fmt.Errorf("pipeline %q declares no operations", name))
			continue
		}

		refs := make([]Ref, 0, len(steps))
		ok := true
		for _, step := range steps {
			ref := ParseRef(step)
			if err := checkOperation(ref); err != nil {
				// A step naming another pipeline is worth its own
				// message: it is a plausible thing to try, and the
				// reason it is refused is not obvious from "unknown
				// operation".
				if _, isPipeline := cfg.Pipelines[ref.Name]; isPipeline {
					err = fmt.Errorf("step %q names another pipeline; a pipeline composes built-in operations only", step)
				}
				problems = append(problems, fmt.Errorf("pipeline %q: %w", name, err))
				ok = false
				continue
			}
			refs = append(refs, ref)
		}
		if ok {
			e.pipelines[name] = refs
		}
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	return e, nil
}

// Pipelines returns the names of the pipelines this engine declares, sorted.
func (e *Engine) Pipelines() []string { return sortedKeys(e.pipelines) }

// Check reports whether a reference resolves and is well formed. It is what
// makes an undefined transform a load-time failure: a package whose templates
// reference a transform that does not exist is broken whether or not a given
// run happens to reach the template.
func (e *Engine) Check(r Ref) error {
	if _, ok := e.pipelines[r.Name]; ok {
		if r.HasArg {
			return fmt.Errorf("pipeline %q takes no argument, but %q supplies one; only %s do",
				r.Name, r.String(), strings.Join(parameterizedNames(), " and "))
		}
		return nil
	}
	if !slices.Contains(Operations, r.Name) {
		return fmt.Errorf("%q is neither a built-in operation (%s) nor a pipeline this package declares (%s)",
			r.Name, strings.Join(Operations, ", "), strings.Join(e.Pipelines(), ", "))
	}
	return checkOperation(r)
}

// checkOperation validates a reference already known to name a built-in.
func checkOperation(r Ref) error {
	if !slices.Contains(Operations, r.Name) {
		return fmt.Errorf("%q is not a built-in operation (%s)", r.Name, strings.Join(Operations, ", "))
	}
	if !Parameterized(r.Name) {
		if r.HasArg {
			return fmt.Errorf("operation %q takes no argument, but %q supplies one", r.Name, r.String())
		}
		return nil
	}
	switch {
	case !r.HasArg:
		return fmt.Errorf("operation %q takes a literal argument, written %s:X", r.Name, r.Name)
	case r.Arg == "":
		return fmt.Errorf("operation %q declares an empty argument", r.Name)
	case strings.Contains(r.Arg, "{{"):
		return fmt.Errorf("operation %q takes string literals only, but %q reads as a dynamic argument; supporting those begins an expression language",
			r.Name, r.String())
	}
	return nil
}

// Apply applies one reference - a built-in operation or one of the package's
// pipelines - to a value.
func (e *Engine) Apply(r Ref, v any) (string, error) {
	if err := e.Check(r); err != nil {
		return "", err
	}
	if steps, ok := e.pipelines[r.Name]; ok {
		return e.run(steps, v)
	}
	return e.operate(r, v)
}

// run applies a pipeline's steps in order.
//
// Only the first step sees the original value, which is what carries any
// model-supplied inflected forms. After one operation the text is no longer the
// value those forms were bound for.
func (e *Engine) run(steps []Ref, v any) (string, error) {
	current := v
	for _, step := range steps {
		out, err := e.operate(step, current)
		if err != nil {
			return "", err
		}
		current = out
	}
	return scalar(current)
}

// Funcs is the FuncMap a rendered template is executed against: every built-in
// operation and every pipeline this package declares, under the names an author
// writes in a template.
//
// A parameterized operation takes its argument first, because a Go template
// pipeline appends the piped value as the last argument.
func (e *Engine) Funcs() template.FuncMap {
	funcs := template.FuncMap{}

	for _, op := range Operations {
		if Parameterized(op) {
			funcs[op] = func(arg string, v any) (string, error) {
				return e.operate(Ref{Name: op, Arg: arg, HasArg: true}, v)
			}
			continue
		}
		funcs[op] = func(v any) (string, error) {
			return e.operate(Ref{Name: op}, v)
		}
	}
	for name, steps := range e.pipelines {
		funcs[name] = func(v any) (string, error) { return e.run(steps, v) }
	}
	return funcs
}

func parameterizedNames() []string {
	var out []string
	for _, op := range Operations {
		if Parameterized(op) {
			out = append(out, op)
		}
	}
	return out
}

// scalar extracts the text a transform operates on.
//
// A kwarg may be declared as a list, and a list has no case and no plural.
// Refusing it names the mistake, where formatting a Go slice into someone's
// source file would bury it.
func scalar(v any) (string, error) {
	switch t := v.(type) {
	case string:
		return t, nil
	case Value:
		return t.Text, nil
	case fmt.Stringer:
		return t.String(), nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return fmt.Sprint(v), nil
	case reflect.Slice, reflect.Array:
		return "", fmt.Errorf("a transform takes a single value, but this one is a list; case and inflection have no meaning for a list")
	case reflect.Invalid:
		return "", fmt.Errorf("a transform takes a single value, but this one is empty")
	default:
		return "", fmt.Errorf("a transform takes a single value, but this one is a %s", rv.Kind())
	}
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
