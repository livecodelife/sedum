package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// sedum render answers one question: given this action and these kwargs, which
// file does the invocation land in. It exists because TOOL_BOUNDARIES §4 tells
// a caller to answer that question and, before this command, gave it no way to
// (prov-2026-b5465dfa). A caller that cannot ask Sedum has to reimplement
// plural, so what is asserted here is that the answer comes from the package's
// own engine rather than from anything this command reproduces.

func TestRenderRendersAnActionsTarget(t *testing.T) {
	out, err := exec(t, "render",
		"--generators", fixtureGenerators(),
		"--package", "chi",
		"--action", "addImport",
		"--kwargs", `{"resource":"Grocery","path":"net/http"}`)
	if err != nil {
		t.Fatalf("render: %v\n%s", err, out)
	}

	if !strings.Contains(out, "internal/handlers/grocery.go") {
		t.Errorf("the rendered path is absent:\n%s", out)
	}
}

// The transform is the whole point. A caller reimplementing snake would agree
// on Grocery and diverge on the first word boundary, so the case that matters
// is the one a naive lowercase gets wrong.
func TestRenderAppliesThePackagesTransforms(t *testing.T) {
	out, err := exec(t, "render",
		"--generators", fixtureGenerators(),
		"--package", "chi",
		"--action", "addImport",
		"--kwargs", `{"resource":"ShoppingList","path":"net/http"}`)
	if err != nil {
		t.Fatalf("render: %v\n%s", err, out)
	}

	if !strings.Contains(out, "internal/handlers/shopping_list.go") {
		t.Errorf("snake was not applied by the package engine:\n%s", out)
	}
}

// --json is the form the written consumer calls. It carries the pattern beside
// the path because a caller comparing a candidate against a region needs to
// know which rule produced the answer, not only what the answer was.
func TestRenderJSONCarriesThePatternBesideThePath(t *testing.T) {
	out, err := exec(t, "render",
		"--generators", fixtureGenerators(),
		"--package", "chi",
		"--action", "addImport",
		"--kwargs", `{"resource":"Grocery","path":"net/http"}`,
		"--json")
	if err != nil {
		t.Fatalf("render: %v\n%s", err, out)
	}

	var got renderResult
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out)
	}

	if got.Package != "chi" || got.Action != "addImport" {
		t.Errorf("the result does not name what was asked: %+v", got)
	}
	if len(got.Targets) != 1 {
		t.Fatalf("want one target, got %d: %+v", len(got.Targets), got.Targets)
	}
	if got.Targets[0].Path != "internal/handlers/grocery.go" {
		t.Errorf("path = %q", got.Targets[0].Path)
	}
	if got.Targets[0].InjectsInto != "internal/handlers/{{resource|snake}}.go" {
		t.Errorf("injects_into = %q", got.Targets[0].InjectsInto)
	}
}

// A composite has no pattern of its own and takes its children's, which is what
// makes one selection visibly touch two files. The cairn fixture's two children
// render different paths, so this fails if only the first is reported or if the
// two are collapsed.
func TestRenderCompositeReportsEveryChildTarget(t *testing.T) {
	out, err := exec(t, "render",
		"--generators", fixtureGenerators(),
		"--package", "cairn",
		"--action", "provisionStep",
		"--kwargs", `{"unit":"Water Intake","name":"Daily Cap","step":"fill"}`,
		"--json")
	if err != nil {
		t.Fatalf("render: %v\n%s", err, out)
	}

	var got renderResult
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out)
	}

	if len(got.Targets) != 2 {
		t.Fatalf("want two targets, got %d: %+v", len(got.Targets), got.Targets)
	}
	// Execution order, which is the order the composite declares.
	if got.Targets[0].Action != "declareConstant" || got.Targets[1].Action != "addStep" {
		t.Errorf("targets are not in execution order: %+v", got.Targets)
	}
	// cairn's slug pipeline is [plural, kebab], so these paths are exactly the
	// divergence this command exists to prevent: a caller reimplementing the
	// transforms reaches "daily-cap" and reports a file no run would write.
	if got.Targets[0].Path != "Shared/daily-caps.crn" {
		t.Errorf("first path = %q", got.Targets[0].Path)
	}
	if got.Targets[1].Path != "Units/water-intakes/Manifest.crn" {
		t.Errorf("second path = %q", got.Targets[1].Path)
	}
}

// A kwarg the pattern references and nothing bound is the caller's most likely
// mistake, so the diagnostic has to be authorable against rather than merely
// debuggable.
func TestRenderNamesWhatNothingBound(t *testing.T) {
	out, err := exec(t, "render",
		"--generators", fixtureGenerators(),
		"--package", "chi",
		"--action", "addImport",
		"--kwargs", `{"path":"net/http"}`)
	wantErr(t, err, "resource")
	_ = out
}

func TestRenderRejectsAnUnknownAction(t *testing.T) {
	_, err := exec(t, "render",
		"--generators", fixtureGenerators(),
		"--package", "chi",
		"--action", "addNothing",
		"--kwargs", `{}`)
	wantErr(t, err, "addNothing", "chi")
}

func TestRenderRejectsAnUnknownPackage(t *testing.T) {
	_, err := exec(t, "render",
		"--generators", fixtureGenerators(),
		"--package", "nosuch",
		"--action", "addImport",
		"--kwargs", `{}`)
	wantErr(t, err, "nosuch")
}

// --kwargs takes a JSON object because that is what an invocation's kwargs are.
// A bare array or scalar is a caller error and is named as one.
func TestRenderRejectsKwargsThatAreNotAnObject(t *testing.T) {
	_, err := exec(t, "render",
		"--generators", fixtureGenerators(),
		"--package", "chi",
		"--action", "addImport",
		"--kwargs", `["resource"]`)
	wantErr(t, err, "--kwargs")
}

func TestRenderRejectsMalformedKwargs(t *testing.T) {
	_, err := exec(t, "render",
		"--generators", fixtureGenerators(),
		"--package", "chi",
		"--action", "addImport",
		"--kwargs", `{"resource":`)
	wantErr(t, err, "--kwargs")
}

// A JSON number is rendered as the caller wrote it. Decoding into float64 would
// render this one as 1e+06, and a path is not where that should be discovered.
func TestRenderPreservesNumbersAsWritten(t *testing.T) {
	out, err := exec(t, "render",
		"--generators", fixtureGenerators(),
		"--package", "chi",
		"--action", "addImport",
		"--kwargs", `{"resource":1000000,"path":"net/http"}`)
	if err != nil {
		t.Fatalf("render: %v\n%s", err, out)
	}
	if !strings.Contains(out, "internal/handlers/1000000.go") {
		t.Errorf("a JSON number did not survive rendering:\n%s", out)
	}
}
