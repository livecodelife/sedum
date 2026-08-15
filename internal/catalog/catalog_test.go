package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/genpkg"
)

// The catalog is the model's entire option set, and it is also what sedum
// actions prints. These tests hold both halves of that: what an entry says
// about an action, and that one build is the only source of either view.

func generators() map[string]string {
	return map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
`,
		"rails/files/app/controllers/{name}_controller.rb": "class {{name|constantize}}Controller\n" +
			"  # sedum:anchor:class_body\nend\n",
		"rails/actions/actions.yaml": `actions:
  createControllerMethod:
    kwargs:
      controller: { type: string, required: true }
      name: { type: string, required: true }
      collection: { type: string, required: false }
    discriminator: name
    variants: [index, show]
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body

  addBeforeFilter:
    kwargs:
      controller: { type: string, required: true }
      filter: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body
`,
		"rails/actions/createControllerMethod/index.rb":    "def index\nend\n",
		"rails/actions/createControllerMethod/show.rb":     "def show\nend\n",
		"rails/actions/createControllerMethod/_default.rb": "def {{name|snake}}\nend\n",
		"rails/actions/addBeforeFilter.rb":                 "before_action :{{filter}}\n",

		// A second package with a composite whose children are hidden. It is
		// the case exposure exists for: the composite is the only way in.
		"chi/sedum.yaml": `name: chi
extensions: [".go"]
comment_prefix: "//"
`,
		"chi/files/internal/handlers/{name}.go": "package handlers\n\n// sedum:anchor:handlers\n",
		"chi/actions/actions.yaml": `actions:
  createHandler:
    composes: [addHandlerFunc, addRoute]

  addHandlerFunc:
    kwargs:
      resource: { type: string, required: true }
      verb: { type: string, required: true }
    injects_into: "internal/handlers/{{resource|snake}}.go"
    anchor: handlers
    exposed: false

  addRoute:
    kwargs:
      resource: { type: string, required: true }
      verb: { type: string, required: false }
    injects_into: "internal/handlers/{{resource|snake}}.go"
    anchor: end_of_file
    exposed: false
`,
		"chi/actions/addHandlerFunc.go": "func {{resource|pascal}}{{verb|pascal}}() {}\n",
		"chi/actions/addRoute.go":       "// route {{resource}}\n",
	}
}

func loadPackages(t *testing.T, files map[string]string, names ...string) []*genpkg.Package {
	t.Helper()

	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}

	set, findings, err := genpkg.Load(root, genpkg.Options{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("fixture package does not load: %s", f)
		}
	}

	out := make([]*genpkg.Package, 0, len(names))
	for _, name := range names {
		pkg, ok := set.Lookup(name)
		if !ok {
			t.Fatalf("package %q did not load", name)
		}
		out = append(out, pkg)
	}
	return out
}

func find(t *testing.T, c Catalog, name string) Action {
	t.Helper()

	matches := c.Lookup(name)
	if len(matches) != 1 {
		t.Fatalf("looking up %q found %d entries, want 1; catalog holds %v", name, len(matches), c.Names())
	}
	return matches[0]
}

// An unexposed action is absent from the option set rather than present and
// rejected. That is the whole point of the tier: selecting one is not merely
// invalid, it is unrepresentable.
func TestUnexposedActionsAreAbsent(t *testing.T) {
	packages := loadPackages(t, generators(), "chi")

	c := Build(packages, Options{})

	if got, want := c.Names(), []string{"createHandler"}; !slices.Equal(got, want) {
		t.Fatalf("catalog holds %v, want %v", got, want)
	}
}

// sedum actions --all is the exception, because an author debugging why an
// action never gets picked needs to see that it is hidden. Hidden is what gets
// marked; exposed carries nothing, so the mark says something wherever it
// appears.
func TestAllIncludesUnexposedActionsMarked(t *testing.T) {
	packages := loadPackages(t, generators(), "chi")

	c := Build(packages, Options{IncludeUnexposed: true})

	want := []string{"addHandlerFunc", "addRoute", "createHandler"}
	if got := c.Names(); !slices.Equal(got, want) {
		t.Fatalf("catalog holds %v, want %v", got, want)
	}

	hidden := find(t, c, "addRoute")
	if hidden.Exposed == nil || *hidden.Exposed {
		t.Errorf("addRoute is unexposed but its entry does not say so: %+v", hidden.Exposed)
	}
	if exposed := find(t, c, "createHandler"); exposed.Exposed != nil {
		t.Errorf("createHandler is exposed, so its entry should carry no exposure mark; got %v", *exposed.Exposed)
	}
}

// The variant list exists so the model can see where the cliff is: a declared
// value gets a dedicated template and anything else falls to _default. Knowing
// the fallback exists is the other half of that, and it is why has_default is
// carried rather than left for the model to assume either way.
func TestDiscriminatedActionCarriesItsVariantsAndFallback(t *testing.T) {
	packages := loadPackages(t, generators(), "rails")

	c := Build(packages, Options{})

	entry := find(t, c, "createControllerMethod")
	if entry.Discriminator != "name" {
		t.Errorf("discriminator is %q, want %q", entry.Discriminator, "name")
	}
	if want := []string{"index", "show"}; !slices.Equal(entry.Variants, want) {
		t.Errorf("variants are %v, want %v", entry.Variants, want)
	}
	if !entry.HasDefault {
		t.Error("the package ships createControllerMethod/_default.rb, but the entry does not say a fallback exists")
	}

	// An action with no discriminator has neither, and says so by omission
	// rather than by an empty list the model has to interpret.
	plain := find(t, c, "addBeforeFilter")
	if plain.Discriminator != "" || plain.Variants != nil || plain.HasDefault {
		t.Errorf("addBeforeFilter has no variants, but its entry claims %+v", plain)
	}
}

// A composite's schema is the union of its children's, so the caller binds a
// shared kwarg once. The entry says which children it will reach, because one
// selection touching two files is not something to infer from a name.
func TestCompositeEntryShowsItsChildrenAndUnionSchema(t *testing.T) {
	packages := loadPackages(t, generators(), "chi")

	entry := find(t, Build(packages, Options{}), "createHandler")

	if want := []string{"addHandlerFunc", "addRoute"}; !slices.Equal(entry.Composes, want) {
		t.Errorf("composes is %v, want %v", entry.Composes, want)
	}
	// resource is required by both children; verb is required by one and
	// optional in the other, and the union takes the stricter flag.
	if got, ok := entry.Kwargs["verb"]; !ok || !got.Required {
		t.Errorf("verb is required by a child, so the composite requires it; got %+v (present=%v)", got, ok)
	}
	if got, ok := entry.Kwargs["resource"]; !ok || got.Type != "string" {
		t.Errorf("resource should be a required string; got %+v (present=%v)", got, ok)
	}
}

// Without the target pattern a model has a kwarg named controller, a file named
// app/controllers/users_controller.rb, and nothing connecting the two - so it
// binds the path to the kwarg and every retry repeats the mistake. The pattern
// is what makes the binding derivable forwards (prov-2026-1bbb8e2e).
func TestEntryCarriesItsTargetPattern(t *testing.T) {
	packages := loadPackages(t, generators(), "rails", "chi")
	c := Build(packages, Options{})

	entry := find(t, c, "createControllerMethod")
	want := []string{"app/controllers/{{controller|snake}}_controller.rb"}
	if !slices.Equal(entry.InjectsInto, want) {
		t.Errorf("injects_into is %v, want %v", entry.InjectsInto, want)
	}

	// A composite has no pattern of its own and takes its children's, in
	// execution order, because one selection touching two files is the fact
	// the entry has to convey.
	composite := find(t, c, "createHandler")
	if len(composite.InjectsInto) != 2 {
		t.Fatalf("the composite carries %v, want one pattern per child", composite.InjectsInto)
	}
	for _, pattern := range composite.InjectsInto {
		if !strings.Contains(pattern, "internal/handlers/") {
			t.Errorf("a child's pattern is missing from the composite entry: %v", composite.InjectsInto)
		}
	}

	// The pattern is passed through as the author wrote it. Interpreting or
	// expanding it here would put target knowledge in Sedum's core.
	payload, err := c.JSON()
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if !strings.Contains(string(payload), "{{controller|snake}}") {
		t.Errorf("the transform pipe did not survive into the payload:\n%s", payload)
	}
}

// A record's paths may resolve to more than one package, and its catalog is
// their union. Entries are ordered so that the same packages produce the same
// bytes every time - the prompt embeds this, and an option set that reordered
// between runs would make one call's input differ from another's for no reason
// the run could explain.
func TestCatalogSpansPackagesInStableOrder(t *testing.T) {
	packages := loadPackages(t, generators(), "rails", "chi")

	first := Build(packages, Options{})
	want := []string{"addBeforeFilter", "createControllerMethod", "createHandler"}
	if got := first.Names(); !slices.Equal(got, want) {
		t.Fatalf("catalog holds %v, want %v", got, want)
	}

	// Built from the packages in the other order, it is the same catalog.
	second := Build([]*genpkg.Package{packages[1], packages[0]}, Options{})

	a, err := first.JSON()
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}
	b, err := second.JSON()
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if string(a) != string(b) {
		t.Errorf("the same packages in a different order produced different catalogs:\n%s\n---\n%s", a, b)
	}
}

// Two packages a record spans may both declare a name. The model returns an
// action and its kwargs and never a package, so nothing in its response could
// say which was meant. Lookup reports the ambiguity rather than resolving it by
// an ordering rule nothing declared.
func TestLookupReportsAnAmbiguousName(t *testing.T) {
	packages := loadPackages(t, generators(), "rails")
	clashing := &genpkg.Package{
		Name:    "sinatra",
		Actions: map[string]*genpkg.Action{"addBeforeFilter": {Name: "addBeforeFilter", Exposed: true}},
	}

	c := Build(append(packages, clashing), Options{})

	matches := c.Lookup("addBeforeFilter")
	if len(matches) != 2 {
		t.Fatalf("two packages declare addBeforeFilter, but Lookup found %d entries", len(matches))
	}
	if matches[0].Package != "rails" || matches[1].Package != "sinatra" {
		t.Errorf("ambiguous entries are ordered %q, %q; want rails, sinatra",
			matches[0].Package, matches[1].Package)
	}
	// Names deduplicates, because a diagnostic listing the option set should
	// not print one name twice.
	if got := c.Names(); !slices.Equal(got, []string{"addBeforeFilter", "createControllerMethod"}) {
		t.Errorf("names are %v, want the option set once each", got)
	}
}

// The JSON is what the prompt embeds, so it has to be readable as the author's
// own configuration. Go's default encoder would escape the characters a path
// pattern legitimately contains.
func TestJSONDoesNotEscapeConfigurationBackToTheAuthor(t *testing.T) {
	packages := loadPackages(t, generators(), "rails")

	out, err := Build(packages, Options{}).JSON()
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}

	if strings.Contains(string(out), `&`) {
		t.Errorf("catalog JSON escaped a character the package wrote plainly:\n%s", out)
	}
	if strings.HasSuffix(string(out), "\n") {
		t.Error("catalog JSON carries a trailing newline, which the caller has to strip to embed it")
	}

	// It is JSON, and it round-trips: sedum actions --json is a payload
	// something other than a person may read.
	var back Catalog
	if err := json.Unmarshal(out, &back); err != nil {
		t.Fatalf("catalog JSON does not decode: %v\n%s", err, out)
	}
	if !slices.Equal(back.Names(), Build(packages, Options{}).Names()) {
		t.Errorf("decoded catalog holds %v", back.Names())
	}
}

// The catalog carries what a shared kwarg schema structurally cannot: which
// values each variant's template renders. Without it the entry says "optional"
// and every reader believes it - the model on its first call, and a caller
// submitting a recording through replay's terminal validation, which has no
// retry loop at all (prov-2026-369544c1).
func TestVariantRequirementsReachBothConsumers(t *testing.T) {
	// A local override rather than a change to the shared fixture: the
	// behavior under test is a template that renders an optional kwarg, and
	// every other test here wants the fixture it already has.
	files := generators()
	files["rails/actions/createControllerMethod/index.rb"] = "def index\n  {{collection}}\nend\n"

	c := Build(loadPackages(t, files, "rails"), Options{})
	entry := find(t, c, "createControllerMethod")

	if got := entry.VariantRequires["index"]; len(got) != 1 || got[0] != "collection" {
		t.Errorf("index requires %v, want [collection]", got)
	}
	if got := entry.VariantRequires["show"]; len(got) != 0 {
		t.Errorf("show requires %v, want nothing", got)
	}

	// The schema stays a faithful view of actions.yaml. The effective
	// requirement is the union, and folding the derivation into required
	// would leave a reader unable to see what the author actually wrote.
	if entry.Kwargs["collection"].Required {
		t.Error("the derivation rewrote the declared schema; it must be carried alongside it")
	}

	// One encoder serves the prompt and sedum actions --json, so what an
	// author inspects is what the model receives.
	raw, err := c.JSON()
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if !strings.Contains(string(raw), "variant_requires") {
		t.Errorf("the encoded catalog omits variant_requires:\n%s", raw)
	}
}

// The catalog is the only thing the model reads, so a description that does not
// reach it is a sentence written where nobody looks. Both consumers get it
// through the one encoder (prov-2026-c5697387).
func TestDescriptionsReachTheEncodedCatalog(t *testing.T) {
	files := generators()
	files["rails/actions/actions.yaml"] = `actions:
  addBeforeFilter:
    description: Registers a before_action callback.
    kwargs:
      controller:
        type: string
        required: true
        description: the resource name, without the Controller suffix
      filter:
        type: string
        required: true
        description: the method name, bare - the template writes the leading colon
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body
`
	delete(files, "rails/actions/createControllerMethod/index.rb")
	delete(files, "rails/actions/createControllerMethod/show.rb")
	delete(files, "rails/actions/createControllerMethod/_default.rb")

	c := Build(loadPackages(t, files, "rails"), Options{})
	entry := find(t, c, "addBeforeFilter")

	if entry.Description != "Registers a before_action callback." {
		t.Errorf("action description is %q, want it carried verbatim", entry.Description)
	}
	if got := entry.Kwargs["filter"].Description; !strings.Contains(got, "leading colon") {
		t.Errorf("kwarg description is %q, want the author's sentence", got)
	}

	raw, err := c.JSON()
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if !strings.Contains(string(raw), "leading colon") {
		t.Errorf("the encoded catalog drops kwarg descriptions:\n%s", raw)
	}

	// An undescribed kwarg carries no key at all, rather than an empty string
	// the model would have to read as meaningful.
	if strings.Contains(string(raw), `"description": ""`) {
		t.Errorf("an absent description was encoded as empty rather than omitted:\n%s", raw)
	}
}
