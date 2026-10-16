package pipeline

import (
	"context"
	"encoding/json"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/recording"
)

// The go package grows one-off Go through literal-code actions, and every file
// it touches must parse and be gofmt-clean (prov-2026-5403d162).

// goGenerators copies the reference go package into a fresh generators
// directory, alone, so no other package contests .go unless a test adds one.
func goGenerators(t *testing.T, extra map[string]string) string {
	t.Helper()
	files := map[string]string{}
	root := filepath.Join("..", "..", "testdata", "generators", "go")
	err := filepath.WalkDir(root, func(full string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(filepath.Join("..", "..", "testdata", "generators"), full)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = readFile(t, full)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range extra {
		files[k] = v
	}
	return writeTree(t, files)
}

const storePath = "record/store.go"

func goInvocations() []recording.Invocation {
	set := func(action, name, code string) recording.Invocation {
		return recording.Invocation{Action: action, Kwargs: map[string]any{"file": storePath, "name": name, "code": code}}
	}
	return []recording.Invocation{
		{Action: "addImport", Kwargs: map[string]any{"file": storePath, "path": "os"}},
		{Action: "addImport", Kwargs: map[string]any{"file": storePath, "path": "strings"}},
		set("setConst", "Prefix", `const Prefix = "prov-"`),
		set("setVar", "Root", `var Root = os.Getenv("ROOT")`),
		set("setType", "Store", "type Store struct {\n\tDir string\n}"),
		set("setFunction", "New", "func New(dir string) *Store {\n\treturn &Store{Dir: strings.TrimSpace(dir)}\n}"),
		{Action: "setMethod", Kwargs: map[string]any{"file": storePath, "receiver": "Store", "name": "Path",
			"code": "func (s *Store) Path(id string) string {\n\treturn s.Dir + \"/\" + Prefix + id\n}"}},
	}
}

func goRecording(invocations []recording.Invocation) recording.Recording {
	return recording.Recording{
		SedumVersion: recording.Version,
		Packages:     map[string]recording.Package{"go": {Extensions: []string{".go"}}},
		Records: []recording.Record{{
			RecordID: "prov-2026-aaaaaaaa",
			Files: []recording.File{{
				Path: storePath, Package: "go", Template: "{pkg}/{name}.go",
				Captures: map[string]string{"pkg": "record", "name": "store"},
			}},
			Phases: []recording.Phase{{Name: recording.DefaultPhase, Invocations: invocations}},
		}},
	}
}

func mustBeGofmtClean(t *testing.T, path string) string {
	t.Helper()
	got := readFile(t, path)
	formatted, err := format.Source([]byte(got))
	if err != nil {
		t.Fatalf("%s does not parse: %v\n%s", path, err, got)
	}
	if string(formatted) != got {
		t.Errorf("%s is not gofmt-clean:\n got %q\nwant %q", path, got, formatted)
	}
	return got
}

func TestEveryGoActionProducesGofmtCleanGo(t *testing.T) {
	gens := goGenerators(t, nil)
	output := t.TempDir()
	if _, err := Replay(goRecording(goInvocations()), ReplayConfig{Generators: gens, Output: output}); err != nil {
		t.Fatalf("Replay: %v", err)
	}
	got := mustBeGofmtClean(t, filepath.Join(output, storePath))
	for _, want := range []string{"package record", `import "os"`, `import "strings"`, "type Store struct", "func (s *Store) Path"} {
		if !strings.Contains(got, want) {
			t.Errorf("the file lacks %q:\n%s", want, got)
		}
	}

	// Replaying again changes nothing.
	if _, err := Replay(goRecording(goInvocations()), ReplayConfig{Generators: gens, Output: output}); err != nil {
		t.Fatalf("second Replay: %v", err)
	}
	if again := readFile(t, filepath.Join(output, storePath)); again != got {
		t.Errorf("a second replay changed the file:\n%s", again)
	}
}

// Replay's validation is terminal, so Go that is not gofmt-clean fails the
// run and no region is written.
func TestUnformattedGoIsRejectedOnReplay(t *testing.T) {
	gens := goGenerators(t, nil)
	output := t.TempDir()
	bad := []recording.Invocation{{Action: "setFunction", Kwargs: map[string]any{
		"file": storePath, "name": "New", "code": "func New()  {\n    return\n}",
	}}}
	_, err := Replay(goRecording(bad), ReplayConfig{Generators: gens, Output: output})
	if err == nil || !strings.Contains(err.Error(), "gofmt") {
		t.Fatalf("err = %v, want a gofmt rejection", err)
	}
	if content, err := os.ReadFile(filepath.Join(output, storePath)); err == nil && strings.Contains(string(content), "func New") {
		t.Errorf("a rejected region was written:\n%s", content)
	}
}

func goSelection(code string) string {
	encoded, _ := json.Marshal(code)
	return `{"invocations": [{"action": "setFunction", "kwargs": {"file": "` + storePath + `", "name": "New", "code": ` + string(encoded) + `}}]}`
}

// In grow, a check failure is fed back to the model with the formatted form,
// within the retry budget; the corrected answer is what lands.
func TestGrowRepromptsWithTheGofmtForm(t *testing.T) {
	cfg := config(t)
	cfg.Generators = goGenerators(t, nil)
	cfg.Records = recordsDir(t, map[string]string{"prov-2026-aaaaaaaa": "  - " + storePath + "\n"})
	cfg.Retries = 2
	// No completeness question follows the corrected answer: the go package
	// declares its imports anchor optional (prov-2026-73d32e29).
	client := &stub{sequence: []string{
		goSelection("func New()  {\n    return\n}"),
		goSelection("func New() {\n\treturn\n}"),
	}}
	cfg.Client = client

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("grow: %v", err)
	}
	if client.calls != 2 {
		t.Fatalf("calls = %d, want 2: a syntax rejection and the corrected answer, with no question about the optional imports anchor", client.calls)
	}
	followUp := client.prompts[1][len(client.prompts[1])-1].Content
	if !strings.Contains(followUp, "func New() {\n\treturn\n}") {
		t.Errorf("the re-prompt does not show the gofmt form:\n%s", followUp)
	}
	mustBeGofmtClean(t, filepath.Join(cfg.Output, storePath))
}

// Two packages claim .go: the one whose file template matches a path most
// specifically wins it, and --lang still overrides.
func TestTheMostSpecificTemplateWinsAContestedExtension(t *testing.T) {
	cli := map[string]string{
		"cli/sedum.yaml":                    "name: cli\nextensions: [\".go\"]\ncomment_prefix: \"//\"\n",
		"cli/files/cmd/{bin}/{name}_cmd.go": "package main\n\n// sedum:anchor:command\n",
		"cli/actions/actions.yaml": `actions:
  describe:
    kwargs:
      bin: { type: string, required: true }
      name: { type: string, required: true }
    injects_into: "cmd/{{bin}}/{{name}}_cmd.go"
    anchor: command
`,
		"cli/actions/describe.go": "// {{name}}\n",
	}
	cfg := config(t)
	cfg.Generators = goGenerators(t, cli)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - cmd/tool/migrate_cmd.go\n  - cmd/tool/migrate_run.go\n",
	})
	cfg.StopAfterPhase = PhaseResolve

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	got := map[string]string{}
	for _, r := range result.Resolutions {
		got[r.Path] = r.Package.Name
	}
	if got["cmd/tool/migrate_cmd.go"] != "cli" || got["cmd/tool/migrate_run.go"] != "go" {
		t.Errorf("resolutions = %v, want migrate_cmd.go to cli and migrate_run.go to go", got)
	}

	cfg.Lang = []string{"go"}
	result, err = Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("resolve with --lang: %v", err)
	}
	for _, r := range result.Resolutions {
		if r.Package.Name != "go" {
			t.Errorf("--lang go did not override for %s (got %s)", r.Path, r.Package.Name)
		}
	}
}

// go.mod and go.sum belong to `go mod`: a record naming them is authorized but
// Sedum leaves them to it, at any depth (prov-2026-7a704b59).
func TestTheGoPackageLeavesModuleFilesToGoMod(t *testing.T) {
	cfg := config(t)
	cfg.Generators = goGenerators(t, nil)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - go.mod\n  - go.sum\n  - tools/go.mod\n  - " + storePath + "\n",
	})
	cfg.StopAfterPhase = PhaseResolve

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	unmanaged := map[string]bool{}
	for _, r := range result.Unmanaged {
		unmanaged[r.Path] = true
	}
	for _, p := range []string{"go.mod", "go.sum", "tools/go.mod"} {
		if !unmanaged[p] {
			t.Errorf("%s is not unmanaged; unmanaged = %v", p, unmanaged)
		}
	}
	if unmanaged[storePath] {
		t.Errorf("%s was left unmanaged", storePath)
	}
}

// An answer using a package it did not import is re-prompted in Phase 5, and
// the answer that adds the import is what lands (prov-2026-0ce03d1c).
func TestGrowRepromptsForAnUnimportedQualifier(t *testing.T) {
	cfg := config(t)
	cfg.Generators = goGenerators(t, nil)
	cfg.Records = recordsDir(t, map[string]string{"prov-2026-aaaaaaaa": "  - " + storePath + "\n"})
	cfg.Retries = 2
	code, _ := json.Marshal("func New() error {\n\treturn errors.New(\"no\")\n}")
	without := `{"invocations": [{"action": "setFunction", "kwargs": {"file": "` + storePath + `", "name": "New", "code": ` + string(code) + `}}]}`
	with := `{"invocations": [{"action": "addImport", "kwargs": {"file": "` + storePath + `", "path": "errors"}}, {"action": "setFunction", "kwargs": {"file": "` + storePath + `", "name": "New", "code": ` + string(code) + `}}]}`
	client := &stub{sequence: []string{without, with}}
	cfg.Client = client

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("grow: %v", err)
	}
	if client.calls != 2 {
		t.Fatalf("calls = %d, want 2", client.calls)
	}
	followUp := client.prompts[1][len(client.prompts[1])-1].Content
	if !strings.Contains(followUp, "errors") || !strings.Contains(followUp, "without importing") {
		t.Errorf("the re-prompt does not name the missing import:\n%s", followUp)
	}
	got := mustBeGofmtClean(t, filepath.Join(cfg.Output, storePath))
	if !strings.Contains(got, `import "errors"`) {
		t.Errorf("the corrected answer's import did not land:\n%s", got)
	}
}
