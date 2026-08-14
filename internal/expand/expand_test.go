package expand

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/recording"
	"github.com/calebcowen/sedum/internal/resolve"
)

// Phase 6 is the last deterministic step before anything is written. Its job is
// to turn an action name and a bag of kwargs into a file, a template, and
// rendered text - and to refuse, specifically, when any of the three cannot be
// decided.

func generators() map[string]string {
	return map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
  instantize: [plural, "prefix:@"]
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

  createMissingFile:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_helper.rb"
    anchor: class_body
`,
		"rails/actions/createControllerMethod/index.rb": "def index\n" +
			"  {{collection|instantize}} = {{collection|constantize}}.all\n" +
			"  render json: {{collection|instantize}}\n" +
			"end\n",
		"rails/actions/createControllerMethod/show.rb":     "def show\nend\n",
		"rails/actions/createControllerMethod/_default.rb": "def {{name|snake}}\nend\n",
		"rails/actions/createMissingFile.rb":               "# helper\n",

		// A second package, so that a record spanning packages draws its
		// catalog from the union of both.
		"cairn/sedum.yaml": `name: cairn
extensions: [".crn"]
comment_prefix: ";;"
transforms:
  slug: [plural, kebab]
  unitname: [singular, pascal]
  shout: [snake, upper]
`,
		"cairn/files/Units/{name}/Manifest.crn": "unit {{name|unitname}}\n  ;; sedum:anchor:steps\nend\n",
		"cairn/actions/actions.yaml": `actions:
  addStep:
    kwargs:
      unit: { type: string, required: true }
      step: { type: string, required: true }
    injects_into: "Units/{{unit|slug}}/Manifest.crn"
    anchor: steps
`,
		"cairn/actions/addStep.crn": "step {{step|shout}}\n",
	}
}

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

func loadSet(t *testing.T, files map[string]string) *genpkg.Set {
	t.Helper()

	set, findings, err := genpkg.Load(writeTree(t, files), genpkg.Options{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("fixture package does not load: %s", f)
		}
	}
	return set
}

// created builds the Phase 3 output a record would have produced: one file per
// path, resolved to the named package.
func created(t *testing.T, set *genpkg.Set, paths map[string]string) []resolve.File {
	t.Helper()

	var out []resolve.File
	for path, pkgName := range paths {
		pkg, ok := set.Lookup(pkgName)
		if !ok {
			t.Fatalf("package %q did not load", pkgName)
		}
		out = append(out, resolve.File{
			Resolution: resolve.Resolution{RecordID: "PR-014", Path: path, Package: pkg},
		})
	}
	return out
}

func TestExpandsSimpleInvocation(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{"app/controllers/users_controller.rb": "rails"})

	got, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "createControllerMethod",
		Kwargs: map[string]any{"controller": "users", "name": "index", "collection": "users"},
	}})
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expanded to %d invocations, want 1", len(got))
	}

	inv := got[0]
	if inv.Path != "app/controllers/users_controller.rb" {
		t.Errorf("path = %q, want the controller injects_into rendered to", inv.Path)
	}
	if inv.Variant != "index" {
		t.Errorf("variant = %q, want index", inv.Variant)
	}
	if inv.RecordID != "PR-014" {
		t.Errorf("record = %q, want PR-014", inv.RecordID)
	}

	// Transforms are applied: the package says what instantize and
	// constantize mean, and Sedum's core knows neither.
	want := "def index\n  @users = User.all\n  render json: @users\nend\n"
	if inv.Content != want {
		t.Errorf("rendered content:\n%s\nwant:\n%s", inv.Content, want)
	}
}

// A discriminator value with no dedicated template falls to _default rather
// than failing, and the variant recorded is the value that was bound - not
// "_default" - so the marker says what was asked for.
func TestFallsBackToDefaultVariant(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{"app/controllers/users_controller.rb": "rails"})

	got, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "createControllerMethod",
		Kwargs: map[string]any{"controller": "users", "name": "search"},
	}})
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	if got[0].Variant != "search" {
		t.Errorf("variant = %q, want the bound value rather than the template that served it", got[0].Variant)
	}
	if got[0].Content != "def search\nend\n" {
		t.Errorf("content = %q, want the _default template rendered", got[0].Content)
	}
}

// injects_into must resolve to a file the record created. Zero matches is the
// case a record omitting a companion file produces, and the diagnostic names
// the path that is missing alongside the ones that are not.
func TestInjectsIntoMustNameACreatedFile(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{"app/controllers/users_controller.rb": "rails"})

	_, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "createMissingFile",
		Kwargs: map[string]any{"controller": "users"},
	}})
	if err == nil {
		t.Fatal("an action injecting into a path no record authorized was expanded")
	}
	for _, want := range []string{"createMissingFile", "users_helper.rb", "users_controller.rb"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not name %q: %v", want, err)
		}
	}
}

// A record's catalog is the union of exposed actions across every package its
// paths resolved to, because resolution is per file rather than per run.
func TestCatalogSpansPackages(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{
		"app/controllers/users_controller.rb": "rails",
		"Units/users/Manifest.crn":            "cairn",
	})

	got, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "createControllerMethod",
		Kwargs: map[string]any{"controller": "users", "name": "show"},
	}, {
		Action: "addStep",
		Kwargs: map[string]any{"unit": "user", "step": "build"},
	}})
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expanded to %d invocations, want 2", len(got))
	}

	// Each resolved under its own package's conventions.
	if got[0].Package.Name != "rails" || got[1].Package.Name != "cairn" {
		t.Errorf("packages = %s, %s; want rails, cairn", got[0].Package.Name, got[1].Package.Name)
	}
	if got[1].Path != "Units/users/Manifest.crn" {
		t.Errorf("path = %q, want cairn's slug pipeline applied to the unit", got[1].Path)
	}
	if got[1].Content != "step BUILD\n" {
		t.Errorf("content = %q, want cairn's shout pipeline applied", got[1].Content)
	}
}

func TestUnknownActionIsAnError(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{"app/controllers/users_controller.rb": "rails"})

	_, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "createControllerMehtod",
		Kwargs: map[string]any{"controller": "users"},
	}})
	if err == nil {
		t.Fatal("an action no package declares was expanded")
	}
	if !strings.Contains(err.Error(), "createControllerMehtod") {
		t.Errorf("error does not name the action: %v", err)
	}
}

// Composite expansion is a later milestone. A composite reaching here is
// reported rather than treated as a simple action, because searching the
// filesystem for a composite's template is a bug rather than a fallback.
func TestCompositeIsReportedAsUnsupported(t *testing.T) {
	files := generators()
	files["rails/actions/actions.yaml"] += `
  createResource:
    composes: [createControllerMethod]
`
	set := loadSet(t, files)
	resolved := created(t, set, map[string]string{"app/controllers/users_controller.rb": "rails"})

	_, err := Expand("PR-014", resolved, []recording.Invocation{{
		Action: "createResource",
		Kwargs: map[string]any{"controller": "users", "name": "index"},
	}})
	if err == nil {
		t.Fatal("a composite was expanded as though it were a simple action")
	}
	if !strings.Contains(err.Error(), "composite") {
		t.Errorf("error does not say the action is a composite: %v", err)
	}
}

// Every problem is reported rather than the first, so a fixture invocation list
// with three mistakes is fixed in one pass.
func TestAllProblemsAreReported(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{"app/controllers/users_controller.rb": "rails"})

	_, err := Expand("PR-014", files, []recording.Invocation{
		{Action: "noSuchAction", Kwargs: map[string]any{}},
		{Action: "createMissingFile", Kwargs: map[string]any{"controller": "users"}},
	})
	if err == nil {
		t.Fatal("an invocation list with two mistakes was expanded")
	}
	for _, want := range []string{"noSuchAction", "createMissingFile"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not report %q: %v", want, err)
		}
	}
}

// A kwarg the template references but nothing bound fails here, with nothing
// written and the bound values named alongside the missing one.
func TestUnboundKwargIsAnError(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{"app/controllers/users_controller.rb": "rails"})

	_, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "createControllerMethod",
		Kwargs: map[string]any{"controller": "users", "name": "index"},
	}})
	if err == nil {
		t.Fatal("an invocation whose template references an unbound value was expanded")
	}
	if !strings.Contains(err.Error(), "collection") {
		t.Errorf("error does not name the value nothing bound: %v", err)
	}
}
