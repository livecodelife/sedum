package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/catalog"
	"github.com/calebcowen/sedum/internal/genpkg"
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
	if !strings.Contains(plain, "2 action(s), 2 exposed") {
		t.Errorf("the default catalog should hold only the exposed actions:\n%s", plain)
	}

	all, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "chi", "--all")
	if err != nil {
		t.Fatalf("actions --all: %v\n%s", err, all)
	}
	if !strings.Contains(all, "addHandlerFunc") || !strings.Contains(all, "[unexposed]") {
		t.Errorf("--all should include hidden actions, marked:\n%s", all)
	}
	if !strings.Contains(all, "2 exposed to the model") {
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

// A package that is not there is a mistake worth naming, and the diagnostic
// says what the directory does declare so the author can see the typo.
func TestActionsNamesTheMissingPackage(t *testing.T) {
	out, err := exec(t, "actions", "--generators", fixtureGenerators(), "--package", "rails2")
	if err == nil {
		t.Fatalf("expected an error for an absent package, got:\n%s", out)
	}
	wantErr(t, err, "rails2")
}
