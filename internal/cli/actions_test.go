package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/catalog"
	"github.com/livecodelife/sedum/internal/genpkg"
)

// sedum actions is the authoring feedback loop for exposure and catalog
// clarity. Its value rests entirely on printing what the model actually
// receives, so what is asserted here is that agreement - not the formatting.

func TestActionsPrintsTheExposedCatalog(t *testing.T) {
	out, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "rails")
	if err != nil {
		t.Fatalf("actions: %v\n%s", err, out)
	}

	for _, want := range []string{"createControllerMethod", "controller", "required", "4 exposed"} {
		if !strings.Contains(out, want) {
			t.Errorf("catalog does not mention %q:\n%s", want, out)
		}
	}
}

// The variant list is in the catalog so the model can see where the cliff is,
// and whether a fallback exists says whether stepping off it is survivable
// (prov-2026-21031113). An author reading this command should see both.
func TestActionsShowsVariantsAndTheirFallback(t *testing.T) {
	out, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "rails")
	if err != nil {
		t.Fatalf("actions: %v\n%s", err, out)
	}

	if !strings.Contains(out, "name selects a template: index, show, create") {
		t.Errorf("the discriminated action does not show its variants:\n%s", out)
	}
	if !strings.Contains(out, "_default") {
		t.Errorf("the package ships a _default template, but the catalog does not say so:\n%s", out)
	}
}

// An unexposed action is absent from the option set, which is the point of the
// tier. --all is the author's way to see that a hidden action exists at all,
// which is what they need when asking why it never gets picked.
func TestActionsHidesUnexposedUnlessAsked(t *testing.T) {
	plain, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "chi")
	if err != nil {
		t.Fatalf("actions: %v\n%s", err, plain)
	}
	// The composite names its children, so the hidden names appear on its
	// composes line. What must not appear is an entry of their own, which is
	// what would put them in the model's option set.
	if strings.Contains(plain, "\naddHandlerFunc") || strings.Contains(plain, "[unexposed]") {
		t.Errorf("an unexposed action reached the default catalog as its own entry:\n%s", plain)
	}
	// Three exposed: the composite, the pattern-targeted addImport, and the
	// free-target addImportTo the fixture carries so that both target forms
	// have a worked example (prov-2026-14c832bf).
	if !strings.Contains(plain, "3 action(s), 3 exposed") {
		t.Errorf("the default catalog should hold only the exposed actions:\n%s", plain)
	}

	all, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "chi", "--all")
	if err != nil {
		t.Fatalf("actions --all: %v\n%s", err, all)
	}
	if !strings.Contains(all, "addHandlerFunc") || !strings.Contains(all, "[unexposed]") {
		t.Errorf("--all should include hidden actions, marked:\n%s", all)
	}
	if !strings.Contains(all, "3 exposed to the model") {
		t.Errorf("--all should still say how many the model would see:\n%s", all)
	}
}

// The command and Phase 4 build their catalog through one code path. This is
// the assertion that keeps that true: --json emits exactly the bytes the
// catalog encoder produces, so a prompt embedding those bytes and this output
// cannot disagree.
func TestActionsJSONIsTheCatalogPayload(t *testing.T) {
	out, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "rails", "--json")
	if err != nil {
		t.Fatalf("actions --json: %v\n%s", err, out)
	}

	set, findings, err := genpkg.Load(fixtureGenerators(), genpkg.Options{Only: []string{"rails"}})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("fixture package does not load: %s", f)
		}
	}
	pkg, ok := set.Lookup("rails")
	if !ok {
		t.Fatal("the rails fixture did not load")
	}
	want, err := catalog.Build([]*genpkg.Package{pkg}, catalog.Options{}).JSON()
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}

	if strings.TrimSpace(out) != string(want) {
		t.Errorf("--json is not the catalog payload:\n got: %s\nwant: %s", out, want)
	}

	var decoded catalog.Catalog
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("--json does not decode: %v", err)
	}
}

// An unmanaged path resolved to no package, by design - the check runs before
// extension resolution so a path no extension can reach is still reportable. A
// run reporting what it created has to handle that as a different branch rather
// than a different line, and until a real record named a Gemfile nothing here
// did: the report dereferenced a package that was never there.
func TestGrowReportsUnmanagedPathsRatherThanCrashing(t *testing.T) {
	dir := t.TempDir()
	generators := writeFiles(t, map[string]string{
		"rails/sedum.yaml": "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n" +
			"unmanaged:\n  - Gemfile\n",
		"rails/files/app/models/{name}.rb": "class X\nend\n",
		"rails/actions/actions.yaml":       "actions: {}\n",
	})
	records := writeFiles(t, map[string]string{
		"r.yml": "id: prov-2026-dddddddd\nintent: |\n  Add a model and the gem it needs.\n" +
			"affected_scope:\n  - app/models/todo.rb\n  - Gemfile\n",
	})

	out, err := exec(t, "grow",
		"--generators", generators,
		"--records", records,
		"--output", dir,
		"--log", filepath.Join(t.TempDir(), "run.log"),
		"--stop-after", "files")
	if err != nil {
		t.Fatalf("grow: %v\n%s", err, out)
	}

	for _, want := range []string{"Gemfile", "left unmanaged", "for a person or another tool"} {
		if !strings.Contains(out, want) {
			t.Errorf("the report does not mention %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "created  Gemfile") {
		t.Errorf("an unmanaged path was reported as created:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(dir, "Gemfile")); err == nil {
		t.Error("an unmanaged path was created")
	}
}

func writeFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// A package that is not there is a mistake worth naming, and the diagnostic
// says what the directory does declare so the author can see the typo.
func TestActionsNamesTheMissingPackage(t *testing.T) {
	out, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "rails2")
	if err == nil {
		t.Fatalf("expected an error for an absent package, got:\n%s", out)
	}
	wantErr(t, err, "rails2")
}

// Everything an author writes to explain a kwarg was previously written in a
// YAML comment, which is exactly where the model cannot see it. A description
// puts it on the surface the model actually reads - and this command is where
// the author confirms that it landed (prov-2026-c5697387).
func TestActionsPrintsAuthoredDescriptions(t *testing.T) {
	out, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "rails")
	if err != nil {
		t.Fatalf("actions: %v\n%s", err, out)
	}

	if !strings.Contains(out, "Registers a before_action callback") {
		t.Errorf("the action's own description is missing:\n%s", out)
	}
	if !strings.Contains(out, "the template writes the leading colon") {
		t.Errorf("a kwarg description is missing:\n%s", out)
	}

	// Absent means absent. An action with no description gets no placeholder,
	// because an invented sentence reads with an authority it has not earned.
	if strings.Contains(out, "addFileHeader\n  <") || strings.Contains(out, "no description") {
		t.Errorf("an undescribed action was given placeholder text:\n%s", out)
	}
}

// The command's claim is that it is evidence of what the model receives. A
// field carried in --json and omitted from the human rendering would make the
// two disagree, which is the one thing it cannot afford (prov-2026-369544c1).
func TestActionsAndJSONAgreeOnWhatTemplatesNeed(t *testing.T) {
	human, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "rails")
	if err != nil {
		t.Fatalf("actions: %v\n%s", err, human)
	}
	raw, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "rails", "--json")
	if err != nil {
		t.Fatalf("actions --json: %v\n%s", err, raw)
	}

	var payload catalog.Catalog
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		t.Fatalf("payload does not decode: %v", err)
	}

	for _, a := range payload.Actions {
		for variant, needs := range a.VariantRequires {
			if len(needs) == 0 {
				continue
			}
			if !strings.Contains(human, variant) {
				t.Errorf("--json reports %s needs %v, but the printed catalog never names it:\n%s",
					variant, needs, human)
			}
		}
		if a.Description != "" && !strings.Contains(human, a.Description) {
			t.Errorf("--json carries a description the printed catalog drops: %q", a.Description)
		}
	}
}
