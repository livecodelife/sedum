// Package render turns a template into text.
//
// Sedum's template syntax is deliberately tiny: {{name}} substitutes a bound
// value, and {{name|op|op:arg}} passes it through transforms first. That is the
// whole grammar. It is not Go template syntax, and it is not a subset of one -
// there are no conditionals, no loops, no arithmetic, and no field access.
//
// Rendering translates each expression into text/template source and executes
// it against a FuncMap built from a package's transforms. Reusing a stdlib
// engine is worth more than hand-rolling a substituter, but it comes with a
// hazard: text/template would happily accept an {{if}} or a {{range}} that
// somebody wrote by mistake or by habit. So the translation is one-way. This
// package parses the recognized syntax itself and rejects everything else,
// which is what keeps the grammar from drifting into Go's by accident.
//
// The same syntax and the same engine serve file templates, action templates,
// and path patterns, so a transform behaves identically wherever it is written.
package render

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/livecodelife/sedum/internal/transform"
)

const (
	openBrace  = "{{"
	closeBrace = "}}"
)

// valueName is the shape of a bound name. Restricting it is what makes
// {{action.name}} and {{.name}} errors rather than something Go's parser would
// find a meaning for.
var valueName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Expr is one {{...}} expression: the value it references and the transforms
// applied to it, in order.
type Expr struct {
	// Source is the expression as written, braces included, so a diagnostic
	// quotes what the author typed.
	Source     string
	Value      string
	Transforms []transform.Ref
}

// Parse reads the expressions in src.
//
// Every problem found is reported rather than the first, because this is what
// package loading uses to check templates, and a diagnostic that stops at the
// first mistake makes fixing a package iterative.
func Parse(src string) ([]Expr, []error) {
	var (
		exprs    []Expr
		problems []error
		rest     = src
	)

	for {
		start := strings.Index(rest, openBrace)
		if start < 0 {
			return exprs, problems
		}
		rest = rest[start+len(openBrace):]

		end := strings.Index(rest, closeBrace)
		if end < 0 {
			problems = append(problems, fmt.Errorf(
				"an expression opens with %s and is never closed with %s", openBrace, closeBrace))
			return exprs, problems
		}

		body := rest[:end]
		rest = rest[end+len(closeBrace):]

		expr, err := parseExpr(body)
		if err != nil {
			problems = append(problems, err)
			continue
		}
		exprs = append(exprs, expr)
	}
}

func parseExpr(body string) (Expr, error) {
	source := openBrace + body + closeBrace
	expr := Expr{Source: source}

	parts := strings.Split(body, "|")
	expr.Value = strings.TrimSpace(parts[0])

	if expr.Value == "" {
		return expr, fmt.Errorf("%s references no value", source)
	}
	if !valueName.MatchString(expr.Value) {
		return expr, fmt.Errorf(
			"%s: %q is not the name of a bound value; the recognized syntax is {{name}} and {{name|transform|transform:argument}}, and nothing else",
			source, expr.Value)
	}

	for _, part := range parts[1:] {
		raw := strings.TrimSpace(part)
		if raw == "" {
			return expr, fmt.Errorf("%s applies a transform with no name", source)
		}
		ref := transform.ParseRef(raw)
		if !valueName.MatchString(ref.Name) {
			return expr, fmt.Errorf("%s: %q is not the name of a transform", source, ref.Name)
		}
		expr.Transforms = append(expr.Transforms, ref)
	}
	return expr, nil
}

// Template is a compiled template, ready to render as often as needed.
type Template struct {
	values []string
	// bare names the values referenced with no transform at all. Those are
	// the ones Go's template engine would format itself, so they are the ones
	// Render has to screen (prov-2026-1cdeb03b).
	bare []string
	tmpl *template.Template
}

// Compile parses a template and resolves every transform it references against
// the engine.
//
// An undefined transform fails here, with no values in sight and nothing
// written. Half a run must not reach disk before a typo in a transform name
// surfaces.
func Compile(engine *transform.Engine, src string) (*Template, error) {
	exprs, problems := Parse(src)

	for _, expr := range exprs {
		for _, ref := range expr.Transforms {
			if err := engine.Check(ref); err != nil {
				problems = append(problems, fmt.Errorf("%s: %w", expr.Source, err))
			}
		}
	}
	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}

	translated, values := translate(src, exprs)
	tmpl, err := template.New("sedum").Funcs(engine.Funcs()).Parse(translated)
	if err != nil {
		return nil, err
	}
	return &Template{values: values, bare: bareValues(exprs), tmpl: tmpl}, nil
}

// Values returns the names this template references, sorted and deduplicated.
func (t *Template) Values() []string { return t.values }

// Render executes the template against a set of bound values.
func (t *Template) Render(values map[string]any) (string, error) {
	var missing []string
	for _, name := range t.values {
		if _, bound := values[name]; !bound {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		// Naming what is bound as well as what is missing turns a typo in
		// a template into a one-line diagnosis.
		return "", fmt.Errorf("template references %s, which nothing bound; the bound values are %s",
			quoteAll(missing), quoteAll(sortedKeys(values)))
	}

	// A list reaches a file only through join. Every other path refuses one —
	// except a bare reference, which text/template would format as [a b c] and
	// write into source with no error anywhere. So it is refused here, before
	// anything executes.
	for _, name := range t.bare {
		if isList(values[name]) {
			return "", fmt.Errorf("%q is bound to a list, and a list renders only through join; write {{%s|join:<separator>}}",
				name, name)
		}
	}

	var out strings.Builder
	if err := t.tmpl.Execute(&out, values); err != nil {
		return "", err
	}
	return out.String(), nil
}

// translate rewrites the source into text/template source.
//
// The literal text between expressions is copied verbatim, and it cannot
// contain an expression opener: parsing consumed every one of them, and an
// unclosed opener was already an error. So there is nothing left for Go's
// parser to find a meaning in.
func translate(src string, exprs []Expr) (string, []string) {
	var (
		out   strings.Builder
		seen  = map[string]bool{}
		names []string
		rest  = src
	)

	for _, expr := range exprs {
		start := strings.Index(rest, expr.Source)
		out.WriteString(rest[:start])
		out.WriteString(goExpr(expr))
		rest = rest[start+len(expr.Source):]

		if !seen[expr.Value] {
			seen[expr.Value] = true
			names = append(names, expr.Value)
		}
	}
	out.WriteString(rest)

	sort.Strings(names)
	return out.String(), names
}

// goExpr writes one expression as a text/template pipeline.
//
// The value is read with index rather than as a field, so that a name Go would
// not accept as a field is still a legal Sedum value name, and an argument is
// written as a quoted Go literal, which is what it is.
func goExpr(expr Expr) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s index . %s", openBrace, strconv.Quote(expr.Value))
	for _, ref := range expr.Transforms {
		b.WriteString(" | ")
		b.WriteString(ref.Name)
		if ref.HasArg {
			b.WriteString(" " + strconv.Quote(ref.Arg))
		}
	}
	b.WriteString(" " + closeBrace)
	return b.String()
}

func quoteAll(names []string) string {
	quoted := make([]string, 0, len(names))
	for _, n := range names {
		quoted = append(quoted, strconv.Quote(n))
	}
	if len(quoted) == 0 {
		return "none"
	}
	return strings.Join(quoted, ", ")
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// bareValues collects the names referenced with no transform applied.
//
// A name referenced both ways counts as bare: the bare occurrence is the one
// that would render Go's formatting, and it does not stop being a problem
// because the same value is joined properly somewhere else in the file.
func bareValues(exprs []Expr) []string {
	var (
		seen = map[string]bool{}
		out  []string
	)
	for _, expr := range exprs {
		if len(expr.Transforms) > 0 || seen[expr.Value] {
			continue
		}
		seen[expr.Value] = true
		out = append(out, expr.Value)
	}
	sort.Strings(out)
	return out
}

// isList reports whether a bound value is a list. A string is a slice of bytes
// to reflect and is not one of these.
func isList(v any) bool {
	switch v.(type) {
	case nil, string, transform.Value:
		return false
	}
	switch reflect.ValueOf(v).Kind() {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}
