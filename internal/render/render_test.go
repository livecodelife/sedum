package render

import (
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/transform"
)

// Rendering is the only place a template's text becomes a file's text, so the
// cases here are about two things: that the recognized syntax produces exactly
// what an author predicted, and that everything outside it is refused rather
// than quietly handed to Go's template engine.

func engine(t *testing.T) *transform.Engine {
	t.Helper()
	e, err := transform.New(transform.Config{
		Pipelines: map[string][]string{
			"constantize": {"singular", "pascal"},
			"instantize":  {"plural", "prefix:@"},
		},
		Exceptions: map[string]map[string]string{
			"pascal": {"url": "URL"},
		},
	})
	if err != nil {
		t.Fatalf("transform.New: %v", err)
	}
	return e
}

func render(t *testing.T, src string, values map[string]any) string {
	t.Helper()
	tmpl, err := Compile(engine(t), src)
	if err != nil {
		t.Fatalf("Compile(%q): %v", src, err)
	}
	out, err := tmpl.Render(values)
	if err != nil {
		t.Fatalf("Render(%q): %v", src, err)
	}
	return out
}

func TestRecognizedSyntax(t *testing.T) {
	values := map[string]any{
		"controller": "users",
		"collection": "users",
		"name":       "index",
	}

	cases := []struct{ src, want string }{
		// Text outside an expression is copied through untouched.
		{"class Foo\nend\n", "class Foo\nend\n"},
		{"def {{name}}\n", "def index\n"},
		{"{{collection|constantize}}", "User"},
		{"{{collection|instantize}}", "@users"},
		// Chained operations apply left to right.
		{"{{collection|singular|pascal|suffix:Serializer}}", "UserSerializer"},
		// Whitespace inside the braces is incidental.
		{"{{ collection | constantize }}", "User"},
		// Several expressions in one line, and repetition of one value.
		{
			"  {{collection|instantize}} = {{collection|constantize}}.all",
			"  @users = User.all",
		},
		// A path pattern is the same syntax over the same engine.
		{
			"app/controllers/{{controller|snake}}_controller.rb",
			"app/controllers/users_controller.rb",
		},
	}
	for _, c := range cases {
		if got := render(t, c.src, values); got != c.want {
			t.Errorf("render(%q) = %q, want %q", c.src, got, c.want)
		}
	}
}

func TestValuesCarryingInflectedForms(t *testing.T) {
	// A model may bind explicit forms; a bare reference still prints the text
	// it bound, and a transform over it sees the forms.
	values := map[string]any{
		"collection": transform.Value{Text: "person", Plural: "people", Singular: "person"},
	}
	if got := render(t, "{{collection}}", values); got != "person" {
		t.Errorf("bare reference = %q, want person", got)
	}
	if got := render(t, "{{collection|instantize}}", values); got != "@people" {
		t.Errorf("instantize = %q, want @people", got)
	}
}

func TestUndefinedTransformFailsAtCompile(t *testing.T) {
	// Half a run must not be written before a typo in a transform name
	// surfaces, so this is a compile-time failure with no values in sight.
	_, err := Compile(engine(t), "{{collection|constantise}}")
	if err == nil {
		t.Fatal("Compile = nil error, want a rejection")
	}
	if !strings.Contains(err.Error(), "constantise") {
		t.Errorf("Compile = %q, want it to name the transform", err)
	}
}

func TestUnboundValueIsAnError(t *testing.T) {
	tmpl, err := Compile(engine(t), "def {{name}}\n")
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	_, err = tmpl.Render(map[string]any{"controller": "users"})
	if err == nil {
		t.Fatal("Render = nil error, want a rejection")
	}
	if !strings.Contains(err.Error(), "name") {
		t.Errorf("Render = %q, want it to name the unbound value", err)
	}
}

func TestSyntaxOutsideTheGrammarIsRejected(t *testing.T) {
	// Go template constructs that happen to be expressible after translation
	// are refused here, so the recognized syntax cannot drift by accident.
	cases := []struct {
		name string
		src  string
	}{
		{"conditional", `{{if .exposed}}x{{end}}`},
		{"range", `{{range .args}}{{.}}{{end}}`},
		{"nested field access", `{{action.name}}`},
		{"go field reference", `{{.name}}`},
		{"function call", `{{printf "%s" name}}`},
		{"variable assignment", `{{$x := name}}`},
		{"empty expression", `{{}}`},
		{"empty transform", `{{name|}}`},
		{"comment", `{{/* nothing */}}`},
		{"unclosed expression", `class {{name`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := Compile(engine(t), c.src); err == nil {
				t.Fatalf("Compile(%q) = nil error, want a rejection", c.src)
			}
		})
	}
}

func TestParseReportsEveryProblem(t *testing.T) {
	// Load-time diagnostics that stop at the first mistake make fixing a
	// package iterative, so parsing collects them all.
	_, errs := Parse("{{first.thing}} and {{second.thing}}")
	if len(errs) != 2 {
		t.Fatalf("Parse returned %d errors, want 2: %v", len(errs), errs)
	}
}

func TestParseExtractsReferences(t *testing.T) {
	// The load-time transform check reads these, so the values and the
	// transform references have to come back separated and in order.
	exprs, errs := Parse("{{collection|plural|prefix:@}} {{name}}")
	if len(errs) != 0 {
		t.Fatalf("Parse: %v", errs)
	}
	if len(exprs) != 2 {
		t.Fatalf("Parse returned %d expressions, want 2", len(exprs))
	}

	first := exprs[0]
	if first.Value != "collection" {
		t.Errorf("first value = %q, want collection", first.Value)
	}
	if len(first.Transforms) != 2 {
		t.Fatalf("first transforms = %v, want 2", first.Transforms)
	}
	if first.Transforms[0].Name != "plural" {
		t.Errorf("first transform = %q, want plural", first.Transforms[0].Name)
	}
	if got := first.Transforms[1]; got.Name != "prefix" || got.Arg != "@" {
		t.Errorf("second transform = %+v, want prefix with argument @", got)
	}

	if exprs[1].Value != "name" || len(exprs[1].Transforms) != 0 {
		t.Errorf("second expression = %+v, want a bare value", exprs[1])
	}
}

func TestRenderIsDeterministic(t *testing.T) {
	// Everything after the model's response is a pure function of its output
	// and the package, which is what makes a rerun reproduce a run.
	src := "{{collection|instantize}} = {{collection|constantize}}.find({{name|snake}})"
	values := map[string]any{"collection": "users", "name": "userId"}

	first := render(t, src, values)
	for range 5 {
		if got := render(t, src, values); got != first {
			t.Fatalf("render is not deterministic: %q then %q", first, got)
		}
	}
	if want := "@users = User.find(user_id)"; first != want {
		t.Errorf("render = %q, want %q", first, want)
	}
}

// A list reaches a template only through join. Before this was so, a
// transform-less reference translated straight to index . "name" and Go's
// template engine formatted the slice, writing [a b c] into a source file with
// no error anywhere (prov-2026-1cdeb03b).
func TestABareListReferenceIsRefused(t *testing.T) {
	tmpl, err := Compile(engine(t), "const words = {{words}};")
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	out, err := tmpl.Render(map[string]any{"words": []any{"a", "b", "c"}})
	if err == nil {
		t.Fatalf("rendering a bare list = %q with nil error, want a refusal", out)
	}
	if !strings.Contains(err.Error(), "list") {
		t.Errorf("rendering a bare list = %q, want it to name the kind of value", err)
	}
	// The Go formatting is the thing this exists to prevent, so its absence is
	// asserted rather than inferred from the error.
	if strings.Contains(out, "[a b c]") {
		t.Errorf("rendering a bare list emitted Go's slice syntax: %q", out)
	}
}

// The same refusal reaches a transform that is not join, so no path renders a
// list by accident.
func TestATransformedListIsRefused(t *testing.T) {
	tmpl, err := Compile(engine(t), "const words = {{words|pascal}};")
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	if out, err := tmpl.Render(map[string]any{"words": []any{"a", "b"}}); err == nil {
		t.Fatalf("rendering a transformed list = %q with nil error, want a refusal", out)
	}
}

// join in a template is the whole point: it is what a package declaring a list
// kwarg writes. An array literal is the shape that reaches a target language,
// which is why the separator carries the quoting.
func TestJoinRendersAListInATemplate(t *testing.T) {
	got := render(t, "const words = ['{{words|join:','}}'];", map[string]any{
		"words": []any{"simply", "easily", "just"},
	})
	if want := "const words = ['simply','easily','just'];"; got != want {
		t.Errorf("render = %q, want %q", got, want)
	}
}

// A pipe cannot be a separator, because the pipe is what separates transforms.
// The grammar splits the expression body before an argument is read, so
// {{words|join:|}} reads as a second transform with no name.
//
// This is a real limit and it is asserted rather than left to be discovered:
// a package wanting alternation joins on something else and lets the target
// language build the alternation, which is also where escaping belongs.
func TestAPipeCannotBeAJoinSeparator(t *testing.T) {
	if _, err := Compile(engine(t), "{{words|join:|}}"); err == nil {
		t.Fatal("compiling a pipe separator = nil error, want the grammar to refuse it")
	}
}

// A joined list composes with the operations that follow it, because after the
// join it is text like any other.
func TestAJoinedListComposesWithLaterTransforms(t *testing.T) {
	got := render(t, "{{words|join:_|upper}}", map[string]any{
		"words": []any{"user", "url"},
	})
	if want := "USER_URL"; got != want {
		t.Errorf("render = %q, want %q", got, want)
	}
}
