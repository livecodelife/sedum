package transform

import (
	"strings"
	"testing"
)

// The deliverable is a table of input/pipeline/expected triples. A transform
// engine that "works" is one where every built-in, every documented package
// pipeline, and every exception table produces the exact string an author
// predicted, so the cases below are written as those triples rather than as
// behavioral prose.

// railsConfig is the worked example from PRD.md: a package whose pipelines and
// exception table say what "constantize" means for its target.
func railsConfig() Config {
	return Config{
		Pipelines: map[string][]string{
			"constantize": {"singular", "pascal"},
			"instantize":  {"plural", "prefix:@"},
			"pathify":     {"plural", "suffix:_path"},
			"tablename":   {"plural", "snake"},
		},
		Exceptions: map[string]map[string]string{
			"pascal": {"url": "URL", "id": "ID"},
			"camel":  {"url": "URL", "id": "ID"},
		},
	}
}

func newEngine(t *testing.T, cfg Config) *Engine {
	t.Helper()
	e, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return e
}

// apply is the whole point of the package in one line: a reference as an author
// writes it, a value, and the string that comes out.
func apply(t *testing.T, e *Engine, ref, in string) string {
	t.Helper()
	out, err := e.Apply(ParseRef(ref), in)
	if err != nil {
		t.Fatalf("apply %s to %q: %v", ref, in, err)
	}
	return out
}

func TestBuiltinOperations(t *testing.T) {
	e := newEngine(t, Config{})

	cases := []struct{ ref, in, want string }{
		// Case conversion normalizes its input first, so the same word
		// arrives at the same answer however the author spelled it.
		{"pascal", "user", "User"},
		{"pascal", "users_controller", "UsersController"},
		{"pascal", "users-controller", "UsersController"},
		{"pascal", "usersController", "UsersController"},
		{"camel", "users_controller", "usersController"},
		{"camel", "UsersController", "usersController"},
		{"snake", "UsersController", "users_controller"},
		{"snake", "users-controller", "users_controller"},
		{"kebab", "UsersController", "users-controller"},

		// upper and lower do not split words: a screaming constant is the
		// composition [snake, upper], not an eleventh operation.
		{"upper", "user_url", "USER_URL"},
		{"lower", "USER_URL", "user_url"},

		{"plural", "user", "users"},
		{"singular", "users", "user"},

		{"prefix:@", "users", "@users"},
		{"suffix:_path", "users", "users_path"},
	}

	for _, c := range cases {
		if got := apply(t, e, c.ref, c.in); got != c.want {
			t.Errorf("%s(%q) = %q, want %q", c.ref, c.in, got, c.want)
		}
	}
}

func TestEveryOperationIsCovered(t *testing.T) {
	// A new built-in must not be able to ship without a case above, since an
	// operation nobody exercised is an operation nobody checked.
	e := newEngine(t, Config{})
	for _, op := range Operations {
		ref := op
		if Parameterized(op) {
			ref = op + ":x"
		}
		if _, err := e.Apply(ParseRef(ref), "user"); err != nil {
			t.Errorf("built-in %s is declared but not implemented: %v", op, err)
		}
	}
}

func TestDeclaredPipelines(t *testing.T) {
	e := newEngine(t, railsConfig())

	cases := []struct{ ref, in, want string }{
		{"constantize", "users", "User"},
		{"constantize", "user", "User"},
		{"instantize", "user", "@users"},
		{"instantize", "users", "@users"},
		{"pathify", "user", "users_path"},
		{"tablename", "UserProfile", "user_profiles"},
	}
	for _, c := range cases {
		if got := apply(t, e, c.ref, c.in); got != c.want {
			t.Errorf("%s(%q) = %q, want %q", c.ref, c.in, got, c.want)
		}
	}
}

func TestPipelinesAreLanguageAgnostic(t *testing.T) {
	// The same operation set, a different set of pipelines, a different
	// target's conventions. Nothing language-specific enters the engine.
	e := newEngine(t, Config{Pipelines: map[string][]string{
		"exported": {"pascal"},
		"receiver": {"singular", "camel"},
		"filename": {"snake", "suffix:.go"},
	}})

	cases := []struct{ ref, in, want string }{
		{"exported", "user_profile", "UserProfile"},
		{"receiver", "Users", "user"},
		{"filename", "UserProfile", "user_profile.go"},
	}
	for _, c := range cases {
		if got := apply(t, e, c.ref, c.in); got != c.want {
			t.Errorf("%s(%q) = %q, want %q", c.ref, c.in, got, c.want)
		}
	}
}

func TestExceptionsApplyPerWord(t *testing.T) {
	e := newEngine(t, railsConfig())

	cases := []struct{ ref, in, want string }{
		// The whole reason exceptions are per word: declaring url -> URL
		// once has to reach every token containing that word.
		{"pascal", "url", "URL"},
		{"pascal", "user_url", "UserURL"},
		{"pascal", "url_id", "URLID"},
		{"pascal", "userId", "UserID"},
		// camel skips the table at position zero: a leading acronym would
		// otherwise render uRL.
		{"camel", "url_id", "urlID"},
		{"camel", "user_url", "userURL"},
	}
	for _, c := range cases {
		if got := apply(t, e, c.ref, c.in); got != c.want {
			t.Errorf("%s(%q) = %q, want %q", c.ref, c.in, got, c.want)
		}
	}
}

func TestExceptionsAreScopedToTheirPackage(t *testing.T) {
	// Two packages loaded in one run must be able to disagree about how a
	// token capitalizes, which is why no transform state is process-global.
	with := newEngine(t, railsConfig())
	without := newEngine(t, Config{})

	if got := apply(t, with, "pascal", "user_url"); got != "UserURL" {
		t.Errorf("with exceptions: got %q", got)
	}
	if got := apply(t, without, "pascal", "user_url"); got != "UserUrl" {
		t.Errorf("without exceptions: got %q, want the unexceptional form", got)
	}
}

func TestInflectionRoundTrips(t *testing.T) {
	e := newEngine(t, Config{})

	// Regular, irregular, and uncountable nouns, each of which has to survive
	// a trip in both directions.
	pairs := []struct{ singular, plural string }{
		{"user", "users"},
		{"box", "boxes"},
		{"church", "churches"},
		{"city", "cities"},
		{"day", "days"},
		{"bus", "buses"},
		{"status", "statuses"},
		{"wolf", "wolves"},
		{"wife", "wives"},
		{"knife", "knives"},
		{"analysis", "analyses"},
		{"datum", "data"},
		{"person", "people"},
		{"child", "children"},
		{"man", "men"},
		{"foot", "feet"},
		{"mouse", "mice"},
		{"sheep", "sheep"},
		{"equipment", "equipment"},
	}
	for _, p := range pairs {
		if got := apply(t, e, "plural", p.singular); got != p.plural {
			t.Errorf("plural(%q) = %q, want %q", p.singular, got, p.plural)
		}
		if got := apply(t, e, "singular", p.plural); got != p.singular {
			t.Errorf("singular(%q) = %q, want %q", p.plural, got, p.singular)
		}
		// Already-inflected input is a fixed point in its own direction.
		if got := apply(t, e, "plural", p.plural); got != p.plural {
			t.Errorf("plural(%q) = %q, want it unchanged", p.plural, got)
		}
		if got := apply(t, e, "singular", p.singular); got != p.singular {
			t.Errorf("singular(%q) = %q, want it unchanged", p.singular, got)
		}
	}
}

func TestModelSuppliedFormsOverrideTheRuleTable(t *testing.T) {
	e := newEngine(t, railsConfig())

	// The rule table is the default, not the authority. A model that binds
	// explicit forms wins for that value.
	v := Value{Text: "octopus", Plural: "octopodes", Singular: "octopus"}

	got, err := e.Apply(ParseRef("plural"), v)
	if err != nil {
		t.Fatalf("plural: %v", err)
	}
	if got != "octopodes" {
		t.Errorf("plural(%v) = %q, want the model-supplied form", v, got)
	}

	// The override reaches the first step of a pipeline, which is the only
	// step that sees the value the model bound a form for.
	got, err = e.Apply(ParseRef("instantize"), v)
	if err != nil {
		t.Fatalf("instantize: %v", err)
	}
	if got != "@octopodes" {
		t.Errorf("instantize(%v) = %q, want @octopodes", v, got)
	}

	// A form the value does not supply falls back to the table.
	only := Value{Text: "user"}
	if got, err := e.Apply(ParseRef("plural"), only); err != nil || got != "users" {
		t.Errorf("plural(%v) = %q, %v; want the rule table's answer", only, got, err)
	}
}

func TestUnresolvableReferences(t *testing.T) {
	e := newEngine(t, railsConfig())

	cases := []struct{ ref, wants string }{
		{"nonesuch", "neither a built-in operation"},
		// A bare parameterized operation has not said what to prepend.
		{"prefix", "takes a literal argument"},
		// An argument on an operation that takes none is a typo, not a
		// value to ignore.
		{"pascal:x", "takes no argument"},
		{"prefix:", "empty argument"},
		// Dynamic arguments would begin an expression language.
		{"prefix:{{other}}", "string literals only"},
	}
	for _, c := range cases {
		err := e.Check(ParseRef(c.ref))
		if err == nil {
			t.Errorf("Check(%q) = nil, want an error", c.ref)
			continue
		}
		if !strings.Contains(err.Error(), c.wants) {
			t.Errorf("Check(%q) = %q, want it to mention %q", c.ref, err, c.wants)
		}
	}

	for _, ref := range []string{"pascal", "prefix:@", "constantize", "tablename"} {
		if err := e.Check(ParseRef(ref)); err != nil {
			t.Errorf("Check(%q) = %v, want it to resolve", ref, err)
		}
	}
}

func TestMalformedPipelinesAreRejectedAtLoad(t *testing.T) {
	cases := []struct {
		name  string
		cfg   Config
		wants string
	}{
		{
			name:  "step names no operation",
			cfg:   Config{Pipelines: map[string][]string{"tablename": {"plural", "bogus"}}},
			wants: "bogus",
		},
		{
			name:  "step takes a dynamic argument",
			cfg:   Config{Pipelines: map[string][]string{"instantize": {"prefix:{{sigil}}"}}},
			wants: "string literals only",
		},
		{
			name: "pipeline composes another pipeline",
			cfg: Config{Pipelines: map[string][]string{
				"constantize": {"singular", "pascal"},
				"modelname":   {"constantize"},
			}},
			wants: "built-in operations",
		},
		{
			name:  "pipeline shadows a built-in",
			cfg:   Config{Pipelines: map[string][]string{"pascal": {"upper"}}},
			wants: "shadows",
		},
		{
			name:  "pipeline has no steps",
			cfg:   Config{Pipelines: map[string][]string{"nothing": {}}},
			wants: "no operations",
		},
		{
			name:  "exception table for an operation that consults none",
			cfg:   Config{Exceptions: map[string]map[string]string{"snake": {"url": "URL"}}},
			wants: "consult",
		},
		{
			name:  "exception table for an unknown operation",
			cfg:   Config{Exceptions: map[string]map[string]string{"bogus": {"url": "URL"}}},
			wants: "bogus",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := New(c.cfg)
			if err == nil {
				t.Fatalf("New = nil error, want a rejection")
			}
			if !strings.Contains(err.Error(), c.wants) {
				t.Errorf("New = %q, want it to mention %q", err, c.wants)
			}
		})
	}
}

func TestNewReportsEveryProblem(t *testing.T) {
	// An author fixes a package in one pass, so a config with two broken
	// pipelines names both.
	_, err := New(Config{Pipelines: map[string][]string{
		"first":  {"bogus"},
		"second": {"alsobogus"},
	}})
	if err == nil {
		t.Fatal("New = nil error, want a rejection")
	}
	for _, want := range []string{"bogus", "alsobogus"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("New = %q, want it to mention %q", err, want)
		}
	}
}

func TestOperationArguments(t *testing.T) {
	e := newEngine(t, Config{})

	// The argument is opaque text, not a pattern or an expression.
	cases := []struct{ ref, in, want string }{
		{"prefix:@", "users", "@users"},
		{"prefix:m_", "users", "m_users"},
		{"suffix:.go", "user", "user.go"},
		{"suffix:_id", "user", "user_id"},
	}
	for _, c := range cases {
		if got := apply(t, e, c.ref, c.in); got != c.want {
			t.Errorf("%s(%q) = %q, want %q", c.ref, c.in, got, c.want)
		}
	}
}

func TestNonScalarValuesAreRefused(t *testing.T) {
	e := newEngine(t, Config{})

	// A list kwarg has no case and no plural. Failing loudly beats rendering
	// Go's own formatting of a slice into someone's source file.
	_, err := e.Apply(ParseRef("pascal"), []string{"a", "b"})
	if err == nil {
		t.Fatal("Apply over a list = nil error, want a refusal")
	}
	if !strings.Contains(err.Error(), "list") {
		t.Errorf("Apply over a list = %q, want it to name the kind of value", err)
	}
}
