package expand

import (
	"maps"
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
		"cairn/files/Shared/{name}.crn":         "shared {{name|unitname}}\nend\n",

		// One logical change spanning two files, which is what composites
		// exist for: the constant is declared in one file and the step
		// that names it goes in another.
		"cairn/actions/actions.yaml": `actions:
  provisionStep:
    composes: [declareConstant, addStep]

  addStep:
    kwargs:
      unit: { type: string, required: true }
      step: { type: string, required: true }
    injects_into: "Units/{{unit|slug}}/Manifest.crn"
    anchor: steps

  declareConstant:
    kwargs:
      unit: { type: string, required: true }
      name: { type: string, required: true }
    injects_into: "Shared/{{name|slug}}.crn"
    anchor: end_of_file
`,
		"cairn/actions/addStep.crn":         "step {{step|shout}}\n",
		"cairn/actions/declareConstant.crn": "const {{name|shout}} = \"{{unit|slug}}\"\n",
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

// The composite payoff: the caller binds the union once and two files receive
// correctly shaped injections. Children come back in declaration order, which
// is the only order there is - nothing reorders them and nothing resolves
// dependencies between them.
func TestCompositeExpandsToItsChildrenInDeclarationOrder(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{
		"Units/users/Manifest.crn": "cairn",
		"Shared/handles.crn":       "cairn",
	})

	got, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "provisionStep",
		Kwargs: map[string]any{"unit": "user", "name": "handle", "step": "build"},
	}})
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("one composite expanded to %d invocations, want one per child", len(got))
	}

	// declareConstant is declared first, so it is applied first.
	if got[0].Action.Name != "declareConstant" || got[1].Action.Name != "addStep" {
		t.Errorf("children ran as %s, %s; want declaration order declareConstant, addStep",
			got[0].Action.Name, got[1].Action.Name)
	}
	if got[0].Path != "Shared/handles.crn" || got[1].Path != "Units/users/Manifest.crn" {
		t.Errorf("children targeted %s, %s; want each child's own injects_into rendered",
			got[0].Path, got[1].Path)
	}
	if got[0].Content != "const HANDLE = \"users\"\n" {
		t.Errorf("declaration content = %q, want the constant rendered from name and unit", got[0].Content)
	}
	if got[1].Content != "step BUILD\n" {
		t.Errorf("step content = %q, want the step rendered from step", got[1].Content)
	}
}

// A region is owned by the action that rendered it, whether that action was
// invoked directly or reached through a composite (prov-2026-a0e37dae). The
// composite has no template, so naming it on a marker would point a reader at
// nothing.
func TestExpandedChildOwnsItsRegionUnderItsOwnName(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{
		"Units/users/Manifest.crn": "cairn",
		"Shared/handles.crn":       "cairn",
	})

	got, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "provisionStep",
		Kwargs: map[string]any{"unit": "user", "name": "handle", "step": "build"},
	}})
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	for _, inv := range got {
		if inv.Action.Name == "provisionStep" {
			t.Errorf("a child's region is owned by the composite; want the child that rendered it")
		}
		if inv.Variant != "" {
			t.Errorf("variant = %q; a child's position in a composite is not a variant", inv.Variant)
		}
	}
}

// Union kwargs are mapped onto the subset each child declares. A kwarg one
// child declares and the other does not reaches only the declaring child, so a
// region's recorded kwargs describe that region rather than the whole composite.
func TestCompositeChildrenReceiveOnlyWhatTheyDeclare(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{
		"Units/users/Manifest.crn": "cairn",
		"Shared/handles.crn":       "cairn",
	})

	got, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "provisionStep",
		Kwargs: map[string]any{"unit": "user", "name": "handle", "step": "build"},
	}})
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}

	want := []map[string]any{
		{"unit": "user", "name": "handle"}, // declareConstant
		{"unit": "user", "step": "build"},  // addStep
	}
	for i, inv := range got {
		if !maps.Equal(inv.Kwargs, want[i]) {
			t.Errorf("%s received %v, want %v", inv.Action.Name, inv.Kwargs, want[i])
		}
	}
}

// The composite's kwarg schema is the union of its children's, so a kwarg two
// children share is supplied once and passed to both.
func TestCompositeSchemaIsTheUnionOfItsChildren(t *testing.T) {
	set := loadSet(t, generators())
	pkg, ok := set.Lookup("cairn")
	if !ok {
		t.Fatal("fixture package cairn did not load")
	}

	want := map[string]genpkg.Kwarg{
		"unit": {Type: "string", Required: true},
		"name": {Type: "string", Required: true},
		"step": {Type: "string", Required: true},
	}
	if !maps.Equal(pkg.Actions["provisionStep"].Kwargs, want) {
		t.Errorf("composite schema = %v, want the union %v", pkg.Actions["provisionStep"].Kwargs, want)
	}
}

// A child targeting a file no record authorized fails loudly rather than
// creating it. This is the case a record that names the implementation and
// forgets the header lands in, and the diagnostic has to name the composite
// too: the child may be unexposed, so the author never invoked it by name.
func TestCompositeChildTargetingAnUnauthorizedPathIsAnError(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{"Units/users/Manifest.crn": "cairn"})

	_, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "provisionStep",
		Kwargs: map[string]any{"unit": "user", "name": "handle", "step": "build"},
	}})
	if err == nil {
		t.Fatal("a composite whose child targets an unauthorized path was expanded")
	}
	for _, want := range []string{"provisionStep", "declareConstant", "Shared/handles.crn"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not name %q: %v", want, err)
		}
	}
}

// Every child is attempted, so a composite with two mistakes in it reports two.
func TestCompositeReportsEveryChildsProblem(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, map[string]string{"app/controllers/users_controller.rb": "rails"})
	cairn, ok := set.Lookup("cairn")
	if !ok {
		t.Fatal("fixture package cairn did not load")
	}
	// Bring cairn into the record's catalog without authorizing either of
	// the paths its children target.
	files = append(files, resolve.File{
		Resolution: resolve.Resolution{RecordID: "PR-014", Path: "Units/other/Manifest.crn", Package: cairn},
	})

	_, err := Expand("PR-014", files, []recording.Invocation{{
		Action: "provisionStep",
		Kwargs: map[string]any{"unit": "user", "name": "handle", "step": "build"},
	}})
	if err == nil {
		t.Fatal("a composite with two unauthorized children was expanded")
	}
	for _, want := range []string{"declareConstant", "addStep"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not report %q: %v", want, err)
		}
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

// A path a package declares unmanaged is authorized and deliberately not
// written by Sedum, so an action targeting one is a mistake in the package
// rather than a record that forgot a file. The diagnostic says which
// (prov-2026-529954ab), because "declared unmanaged" and "never created" have
// different fixes.
func TestActionTargetingAnUnmanagedPathSaysSo(t *testing.T) {
	files := generators()
	files["rails/sedum.yaml"] += "unmanaged:\n  - \"app/controllers/*_helper.rb\"\n"
	set := loadSet(t, files)

	resolved := created(t, set, map[string]string{"app/controllers/users_controller.rb": "rails"})
	pkg, _ := set.Lookup("rails")
	resolved = append(resolved, resolve.File{Resolution: resolve.Resolution{
		RecordID:    "PR-014",
		Path:        "app/controllers/users_helper.rb",
		Unmanaged:   true,
		UnmanagedBy: pkg.Name,
		UnmanagedAs: "app/controllers/*_helper.rb",
	}})

	_, err := Expand("PR-014", resolved, []recording.Invocation{{
		Action: "createMissingFile",
		Kwargs: map[string]any{"controller": "users"},
	}})
	if err == nil {
		t.Fatal("an action injecting into an unmanaged path was expanded")
	}
	for _, want := range []string{"createMissingFile", "users_helper.rb", "unmanaged", "rails"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not name %q: %v", want, err)
		}
	}
}

// rendered builds the Phase 3 output with the template content the file was
// created from, which is what plants the anchors an action fills.
func rendered(t *testing.T, set *genpkg.Set, path, pkgName, content string) resolve.File {
	t.Helper()
	pkg, ok := set.Lookup(pkgName)
	if !ok {
		t.Fatalf("package %q did not load", pkgName)
	}
	return resolve.File{
		Resolution: resolve.Resolution{RecordID: "PR-014", Path: path, Package: pkg},
		Rendered:   content,
	}
}

// A created file states on disk what work it expects, because a file template
// planted the markers an action's anchor targets. Planted minus filled is the
// one failure class Phase 5 is otherwise blind to (prov-2026-6d87dc11).
func TestUnfilledAnchorsAreWhatTheRunMadeAndNothingFills(t *testing.T) {
	set := loadSet(t, generators())
	files := []resolve.File{rendered(t, set,
		"app/controllers/users_controller.rb", "rails",
		"class UsersController\n  # sedum:anchor:class_body\nend\n")}

	// Nothing selected: the anchor the template planted is unfilled.
	empty := Unfilled("PR-014", files, nil)
	if len(empty) != 1 || empty[0].Marker != "class_body" {
		t.Fatalf("an empty selection left %+v unfilled, want class_body", empty)
	}
	if empty[0].Path != "app/controllers/users_controller.rb" {
		t.Errorf("unfilled anchor names %q, want the file that planted it", empty[0].Path)
	}

	// An action anchored to it accounts for it.
	filled := Unfilled("PR-014", files, []recording.Invocation{{
		Action: "createControllerMethod",
		Kwargs: map[string]any{"controller": "users", "name": "index", "collection": "users"},
	}})
	if len(filled) != 0 {
		t.Errorf("class_body reported unfilled after an action anchored to it was selected: %+v", filled)
	}
}

// An unmanaged path is one Sedum did not write, so it has no claim about what
// the file contains and no business reporting an anchor missing from it.
func TestUnmanagedPathsPlantNothing(t *testing.T) {
	files := []resolve.File{{
		Resolution: resolve.Resolution{
			RecordID: "PR-014", Path: "Gemfile", Unmanaged: true, UnmanagedBy: "rails",
		},
		Rendered: "# sedum:anchor:gems\n",
	}}

	if got := Planted(files); len(got) != 0 {
		t.Errorf("an unmanaged path reported planted anchors: %+v", got)
	}
}

// A file whose template planted nothing leaves nothing unfilled. A migration or
// a plain config file may legitimately start blank, and a run that reported one
// as incomplete would make every such template a liability.
func TestAFileWithNoMarkersLeavesNothingUnfilled(t *testing.T) {
	set := loadSet(t, generators())
	files := []resolve.File{rendered(t, set,
		"app/controllers/plain_controller.rb", "rails", "class PlainController\nend\n")}

	if got := Unfilled("PR-014", files, nil); len(got) != 0 {
		t.Errorf("a file planting no markers reported %+v unfilled", got)
	}
}
