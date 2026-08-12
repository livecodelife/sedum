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
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/calebcowen/sedum/internal/transform"
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
	tmpl   *template.Template
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
	return &Template{values: values, tmpl: tmpl}, nil
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
