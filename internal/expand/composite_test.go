package expand

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/record"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/resolve"
)

// The claim M5 exists to make is that one invocation lands correctly in two
// files. Proving it against expansion's return value alone would leave the last
// step assumed, so these cases run the phases a real run would: packages load,
// a record is ingested, its paths resolve and are created from file templates,
// one composite invocation is expanded, and Phase 7 writes.
//
// The target is a language that does not exist, so nothing Sedum could
// plausibly know about it can be doing the work. Its two file kinds stand in
// for the motivating case - a declaration in one file and the definition that
// needs it in another, neither valid alone.

func compositeGenerators() map[string]string {
	return map[string]string{
		"cairn/sedum.yaml": `name: cairn
extensions: [".crn"]
comment_prefix: ";;"
transforms:
  unitname: [singular, pascal]
  slug: [plural, kebab]
  shout: [snake, upper]
`,
		"cairn/files/Units/{name}/Manifest.crn": "unit {{name|unitname}}\n\n  ;; sedum:anchor:steps\n\nend\n",
		"cairn/files/Shared/{name}.crn":         "shared {{name|unitname}}\n",

		"cairn/actions/actions.yaml": `actions:
  provisionStep:
    composes: [declareConstant, addStep]

  addStep:
    kwargs:
      unit: { type: string, required: true }
      step: { type: string, required: true }
    injects_into: "Units/{{unit|slug}}/Manifest.crn"
    anchor: steps
    exposed: false

  declareConstant:
    kwargs:
      unit: { type: string, required: true }
      name: { type: string, required: true }
    injects_into: "Shared/{{name|slug}}.crn"
    anchor: end_of_file
    exposed: false
`,
		"cairn/actions/addStep.crn":         "  step {{step|slug}}\n    for {{unit|unitname}}\n  end\n",
		"cairn/actions/declareConstant.crn": "const {{name|shout}} = \"{{unit|slug}}\"\n",
	}
}

const (
	manifestPath = "Units/users/Manifest.crn"
	sharedPath   = "Shared/handles.crn"
)

// run performs Phases 0 through 3, 6 and 7 over one record authorizing the
// given paths, and returns the output tree.
func run(t *testing.T, authorized []string, invocations []recording.Invocation) (string, error) {
	t.Helper()

	set, findings, err := genpkg.Load(writeTree(t, compositeGenerators()), genpkg.Options{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("fixture package does not load: %s", f)
		}
	}

	scope := "id: prov-2026-aaaaaaaa\nintent: |\n  Provision a step.\naffected_scope:\n"
	for _, path := range authorized {
		scope += "  - " + path + "\n"
	}
	records, _, err := record.Load(writeTree(t, map[string]string{"prov-2026-aaaaaaaa.yml": scope}), record.Options{})
	if err != nil {
		t.Fatalf("record.Load: %v", err)
	}

	resolutions, _, err := resolve.Paths(set, records, nil)
	if err != nil {
		t.Fatalf("resolve.Paths: %v", err)
	}

	output := t.TempDir()
	files, err := resolve.Create(resolutions, resolve.Options{Output: output})
	if err != nil {
		t.Fatalf("resolve.Create: %v", err)
	}

	resolved, err := Expand("prov-2026-aaaaaaaa", files, invocations, nil)
	if err != nil {
		return output, err
	}
	if _, err := inject.Apply(resolved, inject.Options{Output: output}); err != nil {
		return output, err
	}
	return output, nil
}

func readFile(t *testing.T, root, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func provisionStep() []recording.Invocation {
	return []recording.Invocation{{
		Action: "provisionStep",
		Kwargs: map[string]any{"unit": "user", "name": "handle", "step": "build"},
	}}
}

func TestCompositeInjectsIntoBothFiles(t *testing.T) {
	output, err := run(t, []string{manifestPath, sharedPath}, provisionStep())
	if err != nil {
		t.Fatalf("run: %v", err)
	}

	// The declaration goes in the shared file, at end_of_file, under its own
	// name rather than the composite's.
	shared := readFile(t, output, sharedPath)
	want := "shared Handle\n" +
		`;; sedum:declareConstant {"tier":"owned","record":"prov-2026-aaaaaaaa","kwargs":{"name":"handle","unit":"user"}}` + "\n" +
		"const HANDLE = \"users\"\n" +
		";; /sedum:declareConstant\n"
	if shared != want {
		t.Errorf("%s:\n%s\nwant:\n%s", sharedPath, shared, want)
	}

	// The step goes in the manifest, at the marker its file template planted,
	// indented to it.
	manifest := readFile(t, output, manifestPath)
	want = "unit User\n" +
		"\n" +
		"  ;; sedum:anchor:steps\n" +
		`  ;; sedum:addStep {"tier":"owned","record":"prov-2026-aaaaaaaa","kwargs":{"step":"build","unit":"user"}}` + "\n" +
		"  step builds\n" +
		"    for User\n" +
		"  end\n" +
		"  ;; /sedum:addStep\n" +
		"\n" +
		"end\n"
	if manifest != want {
		t.Errorf("%s:\n%s\nwant:\n%s", manifestPath, manifest, want)
	}
}

// Expansion is deterministic and every region an expanded child owns is
// replaced in place rather than appended beside, so a composite is as safe to
// rerun as the simple invocations it becomes.
func TestRerunningACompositeIsByteIdentical(t *testing.T) {
	output, err := run(t, []string{manifestPath, sharedPath}, provisionStep())
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	first := []string{readFile(t, output, manifestPath), readFile(t, output, sharedPath)}

	// A second run against the same tree: Phase 3 is create-if-absent, so
	// both files are the ones the first run left.
	set, _, err := genpkg.Load(writeTree(t, compositeGenerators()), genpkg.Options{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	pkg, _ := set.Lookup("cairn")
	files := []resolve.File{
		{Resolution: resolve.Resolution{RecordID: "prov-2026-aaaaaaaa", Path: manifestPath, Package: pkg}},
		{Resolution: resolve.Resolution{RecordID: "prov-2026-aaaaaaaa", Path: sharedPath, Package: pkg}},
	}
	resolved, err := Expand("prov-2026-aaaaaaaa", files, provisionStep(), nil)
	if err != nil {
		t.Fatalf("second Expand: %v", err)
	}
	if _, err := inject.Apply(resolved, inject.Options{Output: output}); err != nil {
		t.Fatalf("second Apply: %v", err)
	}

	second := []string{readFile(t, output, manifestPath), readFile(t, output, sharedPath)}
	for i, path := range []string{manifestPath, sharedPath} {
		if first[i] != second[i] {
			t.Errorf("%s changed on rerun:\n%s\nwas:\n%s", path, second[i], first[i])
		}
	}
}

// Nothing is created that a provenance record did not authorize, and expansion
// never conjures the companion file however obvious the convention. A record
// naming the manifest and forgetting the shared file halts, naming the file the
// author left out - and nothing is written, so the manifest the record did
// authorize is left as Phase 3 scaffolded it.
func TestOmittingACompositesOtherFileHalts(t *testing.T) {
	output, err := run(t, []string{manifestPath}, provisionStep())
	if err == nil {
		t.Fatal("a composite expanded against a record that authorized only one of its two files")
	}
	for _, want := range []string{"provisionStep", "declareConstant", sharedPath} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not name %q: %v", want, err)
		}
	}

	if _, statErr := os.Stat(filepath.Join(output, filepath.FromSlash(sharedPath))); statErr == nil {
		t.Error("the unauthorized file was created")
	}
	if got := readFile(t, output, manifestPath); strings.Contains(got, "sedum:addStep") {
		t.Errorf("the authorized file was injected into even though the composite failed:\n%s", got)
	}
}
