package genpkg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/transform"
)

// Load-time validation is where a generator package's mistakes are supposed to
// surface, so every rule in the milestone gets a case that breaks exactly that
// rule and asserts the specific diagnostic. A validator that reports "package
// invalid" is not the deliverable.
//
// The valid packages live on disk under testdata/generators, because they are
// also the worked example of what a generator package looks like. The malformed
// ones are built here as one-line mutations of a valid package, so that each
// case reads as the single thing it broke rather than as another near-identical
// directory tree.

// validPackage returns a complete, valid rails package keyed by path relative
// to the generators directory.
func validPackage() map[string]string {
	return map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
  instantize: [plural, "prefix:@"]
`,
		"rails/files/app/controllers/{name}_controller.rb": "class {{name|constantize}}Controller\n" +
			"  # sedum:anchor:class_body_top\n\n  # sedum:anchor:class_body\nend\n",
		"rails/files/app/models/{name}.rb": "class {{name|constantize}}\nend\n",
		// No anchor: a _default template is boilerplate for paths no other
		// template claims, and a marker here would be an injection point no
		// action in this package targets.
		"rails/files/_default.rb": "# generated\n",
		"rails/actions/actions.yaml": `actions:
  addBeforeFilter:
    kwargs:
      controller: { type: string, required: true }
      filter: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body_top
  createControllerMethod:
    kwargs:
      controller: { type: string, required: true }
      name: { type: string, required: true }
      collection: { type: string, required: false }
    discriminator: name
    variants: [index, show]
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body
`,
		"rails/actions/addBeforeFilter.rb": "before_action :{{filter}}\n",
		"rails/actions/createControllerMethod/index.rb": "def index\n" +
			"  {{collection|instantize}} = {{collection|constantize}}.all\nend\n",
		"rails/actions/createControllerMethod/show.rb":     "def show\nend\n",
		"rails/actions/createControllerMethod/_default.rb": "def {{name}}\nend\n",
	}
}

// writeTree materializes a generators directory and returns its path.
func writeTree(t *testing.T, files map[string]string) string {
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
	return root
}

// mutated applies edits to a copy of the valid package. A nil value deletes the
// entry, which is how "the file is missing" cases are expressed.
func mutated(edits map[string]*string) map[string]string {
	files := validPackage()
	for path, content := range edits {
		if content == nil {
			delete(files, path)
			continue
		}
		files[path] = *content
	}
	return files
}

func text(s string) *string { return &s }

func loadTree(t *testing.T, files map[string]string) (*Set, Findings) {
	t.Helper()
	set, findings, err := Load(writeTree(t, files), Options{})
	if err != nil {
		t.Fatalf("Load: unexpected I/O error: %v", err)
	}
	return set, findings
}

// findingFor returns the first finding carrying rule, or fails naming what was
// reported instead.
func findingFor(t *testing.T, findings Findings, rule string) Finding {
	t.Helper()
	for _, f := range findings {
		if f.Rule == rule {
			return f
		}
	}
	var got []string
	for _, f := range findings {
		got = append(got, f.Rule+": "+f.Message)
	}
	t.Fatalf("no finding with rule %q; got:\n  %s", rule, strings.Join(got, "\n  "))
	return Finding{}
}

func TestValidPackageLoadsClean(t *testing.T) {
	set, findings := loadTree(t, validPackage())

	for _, f := range findings {
		t.Errorf("valid package reported %s: %s", f.Rule, f.Message)
	}
	if len(set.Packages) != 1 {
		t.Fatalf("loaded %d packages, want 1", len(set.Packages))
	}

	pkg := set.Packages[0]
	if pkg.Name != "rails" {
		t.Errorf("package name = %q, want rails", pkg.Name)
	}
	if pkg.CommentPrefix != "#" {
		t.Errorf("comment prefix = %q, want #", pkg.CommentPrefix)
	}
	if len(pkg.Actions) != 2 {
		t.Errorf("loaded %d actions, want 2", len(pkg.Actions))
	}

	// exposed defaults to true: authoring an action is enough to make it
	// usable, and hiding is the deliberate act.
	for name, a := range pkg.Actions {
		if !a.Exposed {
			t.Errorf("action %s: Exposed = false, want the default of true", name)
		}
	}
}

// A loaded package carries the engine its templates were checked against, so
// that the vocabulary a package validates against and the one it later renders
// with cannot be two different things.
func TestLoadedPackageCarriesItsTransformEngine(t *testing.T) {
	set, _ := loadTree(t, validPackage())
	pkg := set.Packages[0]

	if pkg.Engine == nil {
		t.Fatal("valid package loaded with no transform engine")
	}
	got, err := pkg.Engine.Apply(transform.ParseRef("constantize"), "users")
	if err != nil {
		t.Fatalf("constantize: %v", err)
	}
	if got != "User" {
		t.Errorf("constantize(users) = %q, want User", got)
	}
}

// Action shape comes from the schema, never from what happens to be on disk.
func TestActionShapeResolution(t *testing.T) {
	set, _ := loadTree(t, validPackage())
	pkg := set.Packages[0]

	simple := pkg.Actions["addBeforeFilter"]
	if simple.Kind() != Simple {
		t.Errorf("addBeforeFilter kind = %v, want Simple", simple.Kind())
	}
	if filepath.Base(simple.Template) != "addBeforeFilter.rb" {
		t.Errorf("addBeforeFilter template = %q, want addBeforeFilter.rb", simple.Template)
	}
	if len(simple.Templates) != 0 {
		t.Errorf("addBeforeFilter resolved variant templates %v, want none", simple.Templates)
	}

	disc := pkg.Actions["createControllerMethod"]
	if disc.Kind() != Discriminated {
		t.Errorf("createControllerMethod kind = %v, want Discriminated", disc.Kind())
	}
	if disc.Template != "" {
		t.Errorf("createControllerMethod resolved a single template %q, want none", disc.Template)
	}
	for _, variant := range []string{"index", "show", DefaultVariant} {
		if disc.Templates[variant] == "" {
			t.Errorf("createControllerMethod has no template for variant %q", variant)
		}
	}
}

func TestExtensionMap(t *testing.T) {
	files := validPackage()
	// A second package claiming the same extension is legal at load. It
	// becomes an error only when a path with that extension appears and no
	// --lang flag disambiguates, which is a resolution concern, not a
	// loading one.
	files["sinatra/sedum.yaml"] = "name: sinatra\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n"
	files["sinatra/files/_default.rb"] = "# generated\n"
	files["sinatra/actions/actions.yaml"] = "actions: {}\n"

	set, findings := loadTree(t, files)
	for _, f := range findings {
		if f.Kind == KindError {
			t.Errorf("contested extension reported an error at load: %s: %s", f.Rule, f.Message)
		}
	}

	claimants := set.ForExtension(".rb")
	if len(claimants) != 2 {
		t.Fatalf("ForExtension(.rb) returned %d packages, want 2", len(claimants))
	}
	if got := set.ForExtension(".go"); len(got) != 0 {
		t.Errorf("ForExtension(.go) returned %d packages, want 0", len(got))
	}
}

func TestOnlyLoadsNamedPackages(t *testing.T) {
	files := validPackage()
	files["sinatra/sedum.yaml"] = "name: sinatra\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n"
	files["sinatra/actions/actions.yaml"] = "actions: {}\n"

	set, _, err := Load(writeTree(t, files), Options{Only: []string{"rails"}})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(set.Packages) != 1 || set.Packages[0].Name != "rails" {
		t.Fatalf("Only=[rails] loaded %d packages, want just rails", len(set.Packages))
	}
}

// Every rule the milestone declares an error for, one case each.
func TestErrorRules(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		rule  string
		// substrings the diagnostic must contain, so that a case fails when
		// the message stops naming the thing it concerns
		mentions []string
	}{
		{
			name:     "sedum.yaml missing",
			files:    mutated(map[string]*string{"rails/sedum.yaml": nil}),
			rule:     RuleManifestMissing,
			mentions: []string{"sedum.yaml"},
		},
		{
			name:     "sedum.yaml malformed",
			files:    mutated(map[string]*string{"rails/sedum.yaml": text("name: rails\nextensions: [.rb\n")}),
			rule:     RuleManifestMalformed,
			mentions: []string{"sedum.yaml"},
		},
		{
			name:     "sedum.yaml carries an unknown key",
			files:    mutated(map[string]*string{"rails/sedum.yaml": text("name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\nextnesions: [\".erb\"]\n")}),
			rule:     RuleManifestMalformed,
			mentions: []string{"sedum.yaml", "extnesions"},
		},
		{
			name:     "actions.yaml missing",
			files:    mutated(map[string]*string{"rails/actions/actions.yaml": nil}),
			rule:     RuleActionsMissing,
			mentions: []string{"actions.yaml"},
		},
		{
			name: "actions.yaml carries an unknown key",
			files: mutated(map[string]*string{"rails/actions/actions.yaml": text(
				"actions:\n  addBeforeFilter:\n    injekts_into: \"x.rb\"\n    anchor: class_body\n")}),
			rule:     RuleActionsMalformed,
			mentions: []string{"actions.yaml", "injekts_into"},
		},
		{
			name: "kwarg type outside the closed set",
			files: mutated(map[string]*string{"rails/actions/actions.yaml": text(
				"actions:\n  addBeforeFilter:\n    kwargs:\n      controller: { type: float, required: true }\n" +
					"    injects_into: \"app/controllers/{{controller}}.rb\"\n    anchor: class_body\n")}),
			rule:     RuleKwargTypeUnknown,
			mentions: []string{"addBeforeFilter", "controller", "float"},
		},
		{
			name:     "declared shape absent from disk",
			files:    mutated(map[string]*string{"rails/actions/addBeforeFilter.rb": nil}),
			rule:     RuleTemplateMissing,
			mentions: []string{"addBeforeFilter"},
		},
		{
			name: "declared shape present in the wrong form",
			files: mutated(map[string]*string{
				"rails/actions/addBeforeFilter.rb":       nil,
				"rails/actions/addBeforeFilter/index.rb": text("x\n"),
			}),
			rule:     RuleTemplateWrongForm,
			mentions: []string{"addBeforeFilter"},
		},
		{
			name: "discriminated action present as a file",
			files: mutated(map[string]*string{
				"rails/actions/createControllerMethod/index.rb":    nil,
				"rails/actions/createControllerMethod/show.rb":     nil,
				"rails/actions/createControllerMethod/_default.rb": nil,
				"rails/actions/createControllerMethod.rb":          text("def x\nend\n"),
			}),
			rule:     RuleTemplateWrongForm,
			mentions: []string{"createControllerMethod"},
		},
		{
			name:     "declared variant with no template file",
			files:    mutated(map[string]*string{"rails/actions/createControllerMethod/show.rb": nil}),
			rule:     RuleVariantTemplateMissing,
			mentions: []string{"createControllerMethod", "show"},
		},
		{
			name: "template references an undefined transform",
			files: mutated(map[string]*string{
				"rails/actions/addBeforeFilter.rb": text("before_action :{{filter|tablenaem}}\n")}),
			rule:     RuleTransformUndefined,
			mentions: []string{"tablenaem"},
		},
		{
			name: "path pattern references an undefined transform",
			files: mutated(map[string]*string{"rails/actions/actions.yaml": text(
				"actions:\n  addBeforeFilter:\n    kwargs:\n      controller: { type: string, required: true }\n" +
					"    injects_into: \"app/controllers/{{controller|snek}}.rb\"\n    anchor: class_body\n")}),
			rule:     RuleTransformUndefined,
			mentions: []string{"snek"},
		},
		{
			name: "file template references an undefined transform",
			files: mutated(map[string]*string{
				"rails/files/app/models/{name}.rb": text("class {{name|constantise}}\nend\n")}),
			rule:     RuleTransformUndefined,
			mentions: []string{"constantise"},
		},
		{
			// The pipelines themselves are checked, not only the
			// references to them. A pipeline whose step names no
			// operation would otherwise load clean and fail at render,
			// which is the failure this milestone exists to prevent.
			name: "pipeline step names no operation",
			files: mutated(map[string]*string{"rails/sedum.yaml": text(
				"name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n" +
					"transforms:\n  constantize: [singula, pascal]\n  instantize: [plural, \"prefix:@\"]\n")}),
			rule:     RuleTransformInvalid,
			mentions: []string{"constantize", "singula"},
		},
		{
			name: "pipeline step takes a dynamic argument",
			files: mutated(map[string]*string{"rails/sedum.yaml": text(
				"name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n" +
					"transforms:\n  constantize: [singular, pascal]\n  instantize: [plural, \"prefix:{{sigil}}\"]\n")}),
			rule:     RuleTransformInvalid,
			mentions: []string{"instantize", "string literals only"},
		},
		{
			name: "op_exceptions declares a table nothing consults",
			files: mutated(map[string]*string{"rails/sedum.yaml": text(
				"name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n" +
					"transforms:\n  constantize: [singular, pascal]\n  instantize: [plural, \"prefix:@\"]\n" +
					"op_exceptions:\n  snake:\n    url: URL\n")}),
			rule:     RuleTransformInvalid,
			mentions: []string{"snake"},
		},
		{
			// A Go template construct survives translation and would
			// work, so the grammar has to be enforced at load or it
			// stops being the grammar.
			name: "template uses a construct outside the grammar",
			files: mutated(map[string]*string{
				"rails/actions/addBeforeFilter.rb": text("{{if filter}}before_action :{{filter}}{{end}}\n")}),
			rule:     RuleTemplateSyntaxInvalid,
			mentions: []string{"addBeforeFilter", "if filter"},
		},
		{
			name: "file template leaves an expression unclosed",
			files: mutated(map[string]*string{
				"rails/files/app/models/{name}.rb": text("class {{name|constantize\nend\n")}),
			rule:     RuleTemplateSyntaxInvalid,
			mentions: []string{"never closed"},
		},
		{
			name: "composite composes another composite",
			files: mutated(map[string]*string{"rails/actions/actions.yaml": text(
				"actions:\n" +
					"  outer:\n    composes: [inner]\n" +
					"  inner:\n    composes: [addBeforeFilter]\n" +
					"  addBeforeFilter:\n    kwargs:\n      filter: { type: string, required: true }\n" +
					"    injects_into: \"a.rb\"\n    anchor: class_body\n")}),
			rule:     RuleCompositeNested,
			mentions: []string{"outer", "inner"},
		},
		{
			name: "composite names an action the package does not define",
			files: mutated(map[string]*string{"rails/actions/actions.yaml": text(
				"actions:\n" +
					"  outer:\n    composes: [addBeforeFilter, addSomethingElse]\n" +
					"  addBeforeFilter:\n    kwargs:\n      filter: { type: string, required: true }\n" +
					"    injects_into: \"a.rb\"\n    anchor: class_body\n")}),
			rule:     RuleCompositeUnknownChild,
			mentions: []string{"outer", "addSomethingElse"},
		},
		{
			name: "composite children declare the same kwarg with different types",
			files: mutated(map[string]*string{"rails/actions/actions.yaml": text(
				"actions:\n" +
					"  outer:\n    composes: [a, b]\n" +
					"  a:\n    kwargs:\n      only: { type: list, required: false }\n" +
					"    injects_into: \"a.rb\"\n    anchor: class_body\n    exposed: false\n" +
					"  b:\n    kwargs:\n      only: { type: string, required: false }\n" +
					"    injects_into: \"b.rb\"\n    anchor: class_body\n    exposed: false\n")}),
			rule:     RuleCompositeKwargConflict,
			mentions: []string{"outer", "only", "list", "string"},
		},
		{
			name: "two file templates tie under the specificity ranking",
			files: mutated(map[string]*string{
				"rails/files/app/{a}/user.rb": text("class User\nend\n"),
				"rails/files/app/{b}/user.rb": text("class User\nend\n"),
			}),
			rule:     RuleFileTemplateTie,
			mentions: []string{"{a}", "{b}"},
		},
		{
			name: "file template path is not a usable pattern",
			files: mutated(map[string]*string{
				"rails/files/app/{a}{b}.rb": text("x\n"),
			}),
			rule:     RuleFileTemplateInvalid,
			mentions: []string{"{a}{b}"},
		},
		{
			name:     "declared name disagrees with the directory name",
			files:    mutated(map[string]*string{"rails/sedum.yaml": text("name: railz\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n")}),
			rule:     RuleNameMismatch,
			mentions: []string{"railz", "rails"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set, findings := loadTree(t, tt.files)

			f := findingFor(t, findings, tt.rule)
			if f.Kind != KindError {
				t.Errorf("%s reported as a warning, want an error", tt.rule)
			}
			if f.Package == "" {
				t.Errorf("%s finding does not name a package", tt.rule)
			}
			for _, want := range tt.mentions {
				if !strings.Contains(f.Message, want) {
					t.Errorf("%s message = %q, does not mention %q", tt.rule, f.Message, want)
				}
			}

			// Packages are wholly valid or rejected. A package with one
			// broken action does not contribute its working actions.
			for _, p := range set.Packages {
				if p.Name == "rails" {
					t.Errorf("a package with a %s error was still loaded", tt.rule)
				}
			}
		})
	}
}

func TestWarningRules(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		rule     string
		mentions []string
	}{
		{
			name: "marker anchor appears in no file template",
			files: mutated(map[string]*string{"rails/actions/actions.yaml": text(
				"actions:\n  addBeforeFilter:\n    kwargs:\n      filter: { type: string, required: true }\n" +
					"    injects_into: \"a.rb\"\n    anchor: clas_body_top\n")}),
			rule:     RuleAnchorUnplanted,
			mentions: []string{"addBeforeFilter", "clas_body_top"},
		},
		{
			name: "unexposed action referenced by no composite is dead config",
			files: mutated(map[string]*string{"rails/actions/actions.yaml": text(
				"actions:\n  addBeforeFilter:\n    kwargs:\n      filter: { type: string, required: true }\n" +
					"    injects_into: \"a.rb\"\n    anchor: class_body\n    exposed: false\n")}),
			rule:     RuleActionDead,
			mentions: []string{"addBeforeFilter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set, findings := loadTree(t, tt.files)

			f := findingFor(t, findings, tt.rule)
			if f.Kind != KindWarning {
				t.Fatalf("%s reported as an error, want a warning", tt.rule)
			}
			for _, want := range tt.mentions {
				if !strings.Contains(f.Message, want) {
					t.Errorf("%s message = %q, does not mention %q", tt.rule, f.Message, want)
				}
			}

			// A warning does not reject the package.
			if len(set.Packages) != 1 {
				t.Errorf("a warning rejected the package; loaded %d packages, want 1", len(set.Packages))
			}
		})
	}
}

// A reserved anchor keyword is not a marker name, so it must not be reported as
// an unplanted marker.
func TestReservedAnchorsAreNotMarkers(t *testing.T) {
	files := mutated(map[string]*string{"rails/actions/actions.yaml": text(
		"actions:\n  addBeforeFilter:\n    kwargs:\n      filter: { type: string, required: true }\n" +
			"    injects_into: \"a.rb\"\n    anchor: end_of_file\n")})

	_, findings := loadTree(t, files)
	for _, f := range findings {
		if f.Rule == RuleAnchorUnplanted {
			t.Errorf("reserved anchor reported as an unplanted marker: %s", f.Message)
		}
	}
}

// An unexposed action that a composite does reference is not dead.
func TestComposedChildIsNotDeadConfig(t *testing.T) {
	files := mutated(map[string]*string{
		"rails/actions/actions.yaml": text(
			"actions:\n" +
				"  createMethod:\n    composes: [addBeforeFilter]\n" +
				"  addBeforeFilter:\n    kwargs:\n      filter: { type: string, required: true }\n" +
				"    injects_into: \"a.rb\"\n    anchor: class_body\n    exposed: false\n"),
		"rails/actions/createControllerMethod/index.rb":    nil,
		"rails/actions/createControllerMethod/show.rb":     nil,
		"rails/actions/createControllerMethod/_default.rb": nil,
	})

	_, findings := loadTree(t, files)
	for _, f := range findings {
		if f.Rule == RuleActionDead {
			t.Errorf("a composed child was reported as dead config: %s", f.Message)
		}
	}
}

// A composite's kwarg schema is the union of its children's: union of names,
// union of required flags. The model binds a shared kwarg once and both
// children receive it.
func TestCompositeKwargUnion(t *testing.T) {
	files := mutated(map[string]*string{
		"rails/actions/actions.yaml": text(
			"actions:\n" +
				"  createMethod:\n    composes: [addDecl, addDef]\n" +
				"  addDecl:\n    kwargs:\n      class: { type: string, required: true }\n" +
				"      name: { type: string, required: false }\n" +
				"    injects_into: \"a.rb\"\n    anchor: class_body\n    exposed: false\n" +
				"  addDef:\n    kwargs:\n      class: { type: string, required: true }\n" +
				"      name: { type: string, required: true }\n" +
				"      args: { type: list, required: false }\n" +
				"    injects_into: \"b.rb\"\n    anchor: class_body\n    exposed: false\n"),
		"rails/actions/addBeforeFilter.rb":                 nil,
		"rails/actions/addDecl.rb":                         text("decl\n"),
		"rails/actions/addDef.rb":                          text("def\n"),
		"rails/actions/createControllerMethod/index.rb":    nil,
		"rails/actions/createControllerMethod/show.rb":     nil,
		"rails/actions/createControllerMethod/_default.rb": nil,
	})

	set, findings := loadTree(t, files)
	for _, f := range findings {
		if f.Kind == KindError {
			t.Fatalf("valid composite rejected: %s: %s", f.Rule, f.Message)
		}
	}

	composite := set.Packages[0].Actions["createMethod"]
	if composite.Kind() != Composite {
		t.Fatalf("createMethod kind = %v, want Composite", composite.Kind())
	}
	// A composite has no template and triggers no filesystem lookup at all.
	if composite.Template != "" || len(composite.Templates) != 0 {
		t.Errorf("composite resolved a template: %q %v", composite.Template, composite.Templates)
	}

	want := map[string]Kwarg{
		"class": {Type: "string", Required: true},
		"name":  {Type: "string", Required: true}, // required by one child, so required
		"args":  {Type: "list", Required: false},
	}
	if len(composite.Kwargs) != len(want) {
		t.Fatalf("composite has %d kwargs %v, want %d", len(composite.Kwargs), composite.Kwargs, len(want))
	}
	for name, w := range want {
		if got := composite.Kwargs[name]; got != w {
			t.Errorf("composite kwarg %s = %+v, want %+v", name, got, w)
		}
	}
}

// Validation reports everything wrong with a package rather than halting on the
// first finding, so an author fixes it in one pass.
func TestAllFindingsAreReported(t *testing.T) {
	files := mutated(map[string]*string{
		"rails/sedum.yaml":                 text("name: railz\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n"),
		"rails/actions/addBeforeFilter.rb": nil,
		"rails/files/app/models/{name}.rb": text("class {{name|constantise}}\nend\n"),
	})

	_, findings := loadTree(t, files)
	for _, rule := range []string{RuleNameMismatch, RuleTemplateMissing, RuleTransformUndefined} {
		findingFor(t, findings, rule)
	}
}

func TestStrictPromotesWarnings(t *testing.T) {
	files := mutated(map[string]*string{"rails/actions/actions.yaml": text(
		"actions:\n  addBeforeFilter:\n    kwargs:\n      filter: { type: string, required: true }\n" +
			"    injects_into: \"a.rb\"\n    anchor: clas_body_top\n")})

	_, findings := loadTree(t, files)
	if findings.HasErrors() {
		t.Fatalf("warnings counted as errors without --strict: %v", findings)
	}
	if !findings.Strict().HasErrors() {
		t.Errorf("Strict() did not promote a warning to an error")
	}
}

// Loading runs with no records, no model, and no network. The check that
// matters is that nothing outside the generators directory is read, which a
// generators directory built in a temp dir with no records alongside it
// demonstrates.
func TestLoadNeedsOnlyTheGeneratorsDirectory(t *testing.T) {
	if _, findings := loadTree(t, validPackage()); len(findings) != 0 {
		t.Fatalf("valid package needed something outside its directory: %v", findings)
	}
}

func TestMissingGeneratorsDirectoryIsAnIOError(t *testing.T) {
	_, _, err := Load(filepath.Join(t.TempDir(), "nope"), Options{})
	if err == nil {
		t.Fatal("Load of a missing generators directory returned no error")
	}
	if !strings.Contains(err.Error(), "nope") {
		t.Errorf("error = %q, does not name the directory", err)
	}
}

// The on-disk fixtures are the worked example, so they must stay loadable.
func TestTestdataGeneratorsAreValid(t *testing.T) {
	set, findings, err := Load(filepath.Join("..", "..", "testdata", "generators"), Options{})
	if err != nil {
		t.Fatalf("Load(testdata/generators): %v", err)
	}
	for _, f := range findings {
		t.Errorf("testdata package %s reported %s: %s", f.Package, f.Rule, f.Message)
	}
	if len(set.Packages) < 2 {
		t.Errorf("testdata/generators holds %d packages, want at least 2 so that more than one stack is exercised", len(set.Packages))
	}
}

// A finding about a file template must name the file that contains the problem.
// Sorting the patterns while leaving the contents in walk order paired them by
// index across two different orderings, which agree until a file and a
// directory share a prefix: '.' sorts before '/', but a walk descends into the
// directory first (prov-2026-11af2675).
func TestFileTemplateFindingsNameTheFileTheyCameFrom(t *testing.T) {
	files := validPackage()
	files["rails/files/app.rb"] = "{{name|nosuchpipeline}}\n"
	files["rails/files/app/x.rb"] = "plain text with no expressions\n"

	_, findings := loadTree(t, files)

	var named []string
	for _, f := range findings {
		if strings.Contains(f.Message, "nosuchpipeline") {
			named = append(named, f.File)
		}
	}
	if len(named) == 0 {
		t.Fatalf("the undefined pipeline went unreported: %v", findings)
	}
	for _, file := range named {
		if file != "files/app.rb" {
			t.Errorf("undefined pipeline in files/app.rb reported against %s", file)
		}
	}
}

// Phase 3 renders the template a path matched, so it looks contents up by
// pattern rather than by position.
func TestFileTemplateContentsAreAddressableByPattern(t *testing.T) {
	set, findings := loadTree(t, validPackage())
	if len(findings) != 0 {
		t.Fatalf("valid package reported findings: %v", findings)
	}

	pkg, ok := set.Lookup("rails")
	if !ok {
		t.Fatal("rails did not load")
	}
	for _, pattern := range pkg.FileTemplates {
		if _, ok := pkg.FileTemplate(pattern); !ok {
			t.Errorf("file template %q has no contents", pattern)
		}
	}
	if _, ok := pkg.FileTemplate("files/does/not/exist.rb"); ok {
		t.Error("a pattern the package does not ship reported contents")
	}
}

// A value a template renders is a value that template needs, and Sedum reads
// that off the template rather than asking an author to restate it in a
// requires list that would drift the first time either was edited
// (prov-2026-369544c1).
func TestTemplateReferencesAreDerivedPerVariant(t *testing.T) {
	files := map[string]string{
		"rails/sedum.yaml": "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n",
		"rails/files/app/controllers/{name}_controller.rb": "class X\n  # sedum:anchor:body\nend\n",
		"rails/actions/actions.yaml": `actions:
  createControllerMethod:
    kwargs:
      controller: { type: string, required: true }
      name: { type: string, required: true }
      collection: { type: string, required: false }
    discriminator: name
    variants: [index, show]
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: body
`,
		"rails/actions/createControllerMethod/index.rb":    "def index\n  {{collection}}\nend\n",
		"rails/actions/createControllerMethod/show.rb":     "def show\nend\n",
		"rails/actions/createControllerMethod/_default.rb": "def {{name}}\nend\n",
	}

	set, findings := loadTree(t, files)
	for _, f := range findings {
		if f.Kind == KindError {
			t.Fatalf("package does not load: %s", f)
		}
	}
	pkg, _ := set.Lookup("rails")
	action := pkg.Actions["createControllerMethod"]

	// index needs collection; show needs nothing; _default needs the
	// discriminator it selects on. A shared schema cannot say any of this.
	for variant, want := range map[string][]string{
		"index":        {"collection"},
		"show":         nil,
		DefaultVariant: {"name"},
	} {
		got := action.Requires(variant)
		if len(got) != len(want) {
			t.Errorf("variant %s requires %v, want %v", variant, got, want)
			continue
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("variant %s requires %v, want %v", variant, got, want)
				break
			}
		}
	}

	// A value with no dedicated template renders the fallback, so it inherits
	// the fallback's requirements rather than none.
	if got := action.Requires("archive"); len(got) != 1 || got[0] != "name" {
		t.Errorf("an uncovered variant requires %v, want the fallback's [name]", got)
	}
}

// A template naming a value the action has no kwarg for is a typo, and finding
// it at load is finding it before a run pays for a model call.
func TestATemplateValueTheActionDoesNotDeclareIsALoadError(t *testing.T) {
	files := map[string]string{
		"rails/sedum.yaml":                 "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n",
		"rails/files/app/models/{name}.rb": "class X\n  # sedum:anchor:body\nend\n",
		"rails/actions/actions.yaml": `actions:
  createModelClass:
    kwargs:
      name: { type: string, required: true }
    injects_into: "app/models/{{name|snake}}.rb"
    anchor: body
`,
		"rails/actions/createModelClass.rb": "class {{name}} < {{parent}}\nend\n",
	}

	_, findings := loadTree(t, files)
	f := findingFor(t, findings, RuleTemplateValueUndeclared)
	if f.Kind != KindError {
		t.Errorf("an undeclared template value reported as %s, want error", f.Kind)
	}
	if !strings.Contains(f.Message, "parent") {
		t.Errorf("diagnostic does not name the value: %s", f.Message)
	}
}

// A single-template action that always renders a kwarg could have declared it
// required. The derivation covers it either way, so this says the declaration
// understates itself rather than failing the package - and a discriminated
// action never earns it, because optional is all its shared schema can say.
func TestAnUnderstatedRequirementWarnsAndOnlyForSingleTemplateActions(t *testing.T) {
	files := map[string]string{
		"rails/sedum.yaml":                 "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n",
		"rails/files/app/models/{name}.rb": "class X\n  # sedum:anchor:body\nend\n",
		"rails/actions/actions.yaml": `actions:
  createModelClass:
    kwargs:
      name: { type: string, required: true }
      table: { type: string, required: false }
    injects_into: "app/models/{{name|snake}}.rb"
    anchor: body
`,
		"rails/actions/createModelClass.rb": "class {{name}}\n  self.table_name = {{table}}\nend\n",
	}

	_, findings := loadTree(t, files)
	f := findingFor(t, findings, RuleKwargRequirementUnderstated)
	if f.Kind != KindWarning {
		t.Errorf("an understated requirement reported as %s, want warning", f.Kind)
	}
	if !strings.Contains(f.Message, "table") {
		t.Errorf("diagnostic does not name the kwarg: %s", f.Message)
	}

	// The same shape under a discriminator is not a defect and must not warn.
	discriminated := map[string]string{
		"rails/sedum.yaml":                 "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n",
		"rails/files/app/models/{name}.rb": "class X\n  # sedum:anchor:body\nend\n",
		"rails/actions/actions.yaml": `actions:
  createModelClass:
    kwargs:
      name: { type: string, required: true }
      table: { type: string, required: false }
    discriminator: name
    variants: [user]
    injects_into: "app/models/{{name|snake}}.rb"
    anchor: body
`,
		"rails/actions/createModelClass/user.rb": "class {{name}}\n  self.table_name = {{table}}\nend\n",
	}
	_, findings = loadTree(t, discriminated)
	for _, f := range findings {
		if f.Rule == RuleKwargRequirementUnderstated {
			t.Errorf("a discriminated action warned about a shared schema it cannot change: %s", f)
		}
	}
}

// A description is authored prose carried untouched. Nothing parses one,
// validates one, or derives a constraint from one - a description must never
// become a place where a rule lives, because a rule the model can read and
// Phase 5 cannot enforce is worse than no rule (prov-2026-c5697387).
func TestDescriptionsAreCarriedAndOptional(t *testing.T) {
	files := map[string]string{
		"rails/sedum.yaml":                 "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n",
		"rails/files/app/models/{name}.rb": "class X\n  # sedum:anchor:body\nend\n",
		"rails/actions/actions.yaml": `actions:
  createModelClass:
    description: Defines the model class body.
    kwargs:
      name:
        type: string
        required: true
        description: the model's singular name, bare
    injects_into: "app/models/{{name|snake}}.rb"
    anchor: body

  addConcern:
    kwargs:
      name: { type: string, required: true }
    injects_into: "app/models/{{name|snake}}.rb"
    anchor: body
`,
		"rails/actions/createModelClass.rb": "class {{name}}\nend\n",
		"rails/actions/addConcern.rb":       "include {{name}}\n",
	}

	set, findings := loadTree(t, files)
	for _, f := range findings {
		if f.Kind == KindError {
			t.Fatalf("a package declaring descriptions does not load: %s", f)
		}
	}
	pkg, _ := set.Lookup("rails")

	described := pkg.Actions["createModelClass"]
	if described.Description != "Defines the model class body." {
		t.Errorf("action description is %q, want it carried verbatim", described.Description)
	}
	if got := described.Kwargs["name"].Description; got != "the model's singular name, bare" {
		t.Errorf("kwarg description is %q, want it carried verbatim", got)
	}

	// Absent means absent. Nothing is synthesised to fill the gap.
	bare := pkg.Actions["addConcern"]
	if bare.Description != "" {
		t.Errorf("an undescribed action gained a description: %q", bare.Description)
	}
	if got := bare.Kwargs["name"].Description; got != "" {
		t.Errorf("an undescribed kwarg gained a description: %q", got)
	}
}

// A composite's schema is the union of its children's, so a description a
// child wrote has to survive the union or it is written where nothing reads
// it. Two children describing one kwarg differently is an authoring problem
// the composite cannot resolve, so the first in declaration order wins rather
// than the two being concatenated into a sentence neither author wrote.
func TestACompositeInheritsItsChildrensDescriptions(t *testing.T) {
	files := map[string]string{
		"chi/sedum.yaml":                        "name: chi\nextensions: [\".go\"]\ncomment_prefix: \"//\"\n",
		"chi/files/internal/handlers/{name}.go": "package handlers\n\n// sedum:anchor:handlers\n",
		"chi/actions/actions.yaml": `actions:
  createHandler:
    composes: [addHandlerFunc, addRoute]

  addHandlerFunc:
    kwargs:
      resource:
        type: string
        required: true
        description: the resource's plural name, lowercase
    injects_into: "internal/handlers/{{resource}}.go"
    anchor: handlers
    exposed: false

  addRoute:
    kwargs:
      resource:
        type: string
        required: true
        description: a second sentence about the same kwarg
    injects_into: "internal/handlers/{{resource}}.go"
    anchor: end_of_file
    exposed: false
`,
		"chi/actions/addHandlerFunc.go": "func {{resource}}() {}\n",
		"chi/actions/addRoute.go":       "// route {{resource}}\n",
	}

	set, findings := loadTree(t, files)
	for _, f := range findings {
		if f.Kind == KindError {
			t.Fatalf("package does not load: %s", f)
		}
	}
	pkg, _ := set.Lookup("chi")

	got := pkg.Actions["createHandler"].Kwargs["resource"].Description
	if got != "the resource's plural name, lowercase" {
		t.Errorf("composite kwarg description is %q, want the first child's", got)
	}
}
