package pipeline

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/selection"
)

// The pipeline's own job is ordering: each phase's output is the next phase's
// only input, and a stop point halts at a boundary with everything before it
// complete and nothing after it started. The phases themselves are tested where
// they live.

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

func generators(t *testing.T) string {
	t.Helper()
	return writeTree(t, map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
`,
		"rails/files/app/models/{name}.rb": "class {{name|constantize}}\n  # sedum:anchor:class_body\nend\n",
		"rails/actions/actions.yaml": `actions:
  addField:
    kwargs:
      model: { type: string, required: true }
      field: { type: string, required: true }
    injects_into: "app/models/{{model|snake}}.rb"
    anchor: class_body
`,
		"rails/actions/addField.rb": "attribute :{{field}}\n",
	})
}

func recordsDir(t *testing.T, records map[string]string) string {
	t.Helper()
	files := map[string]string{}
	for id, scope := range records {
		files[id+".yml"] = "id: " + id + "\nintent: |\n  Something.\naffected_scope:\n" + scope
	}
	return writeTree(t, files)
}

// stub is the model. Phase 4 is the one phase that is not a pure function of
// its inputs, so the pipeline's own tests supply a canned answer: what is being
// tested here is that each phase's output reaches the next, not what a model
// would have chosen.
type stub struct {
	response string
	calls    int
	prompts  [][]selection.Message
}

func (s *stub) Complete(_ context.Context, messages []selection.Message) (selection.Completion, error) {
	s.calls++
	s.prompts = append(s.prompts, messages)
	if s.response == "" {
		return selection.Completion{Content: `{"invocations": []}`}, nil
	}
	return selection.Completion{Content: s.response}, nil
}

func config(t *testing.T) Config {
	t.Helper()
	return Config{
		Generators: generators(t),
		Records: recordsDir(t, map[string]string{
			"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		}),
		Output: t.TempDir(),
		Client: &stub{},
	}
}

// tree returns every file under root, slash-separated and sorted.
func tree(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		out = append(out, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	sort.Strings(out)
	return out
}

func TestRunCarriesEveryPhaseThrough(t *testing.T) {
	cfg := config(t)

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Packages.Packages) != 1 {
		t.Errorf("Phase 0 loaded %d packages, want 1", len(result.Packages.Packages))
	}
	if len(result.Records.Records) != 1 {
		t.Errorf("Phase 1 ingested %d records, want 1", len(result.Records.Records))
	}
	if len(result.Resolutions) != 1 {
		t.Errorf("Phase 2 resolved %d paths, want 1", len(result.Resolutions))
	}
	if len(result.Files) != 1 {
		t.Errorf("Phase 3 produced %d files, want 1", len(result.Files))
	}
	if got := tree(t, cfg.Output); len(got) != 1 || got[0] != "app/models/user.rb" {
		t.Errorf("output tree = %v", got)
	}
}

// Stopping after resolution leaves everything Phase 2 decided available and
// nothing on disk, which is what makes the stop point worth having
// (prov-2026-5696ff65).
func TestStopAfterResolutionDecidesEverythingAndWritesNothing(t *testing.T) {
	cfg := config(t)
	cfg.StopAfterPhase = PhaseResolve

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Resolutions) != 1 {
		t.Fatalf("Phase 2 resolved %d paths, want 1", len(result.Resolutions))
	}
	res := result.Resolutions[0]
	if res.Package == nil || res.Package.Name != "rails" {
		t.Errorf("resolution has no package")
	}
	if res.Template != "app/models/{name}.rb" || res.Captures["name"] != "user" {
		t.Errorf("resolution = template %q captures %v; the matched template and its captures must be available here",
			res.Template, res.Captures)
	}
	if result.Files != nil {
		t.Errorf("Phase 3 ran despite a stop at resolution")
	}
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Errorf("stopping after resolution wrote %v", got)
	}
}

func TestStopAfterFilesCreatesFilesAndStops(t *testing.T) {
	cfg := config(t)
	cfg.StopAfterPhase = PhaseCreate

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("Phase 3 produced %d files, want 1", len(result.Files))
	}
	if got := tree(t, cfg.Output); len(got) != 1 {
		t.Errorf("output tree = %v, want the one authorized file", got)
	}
}

// Resuming after either stop point is an ordinary rerun: Phases 0 through 2 are
// pure and Phase 3 is create-if-absent, so nothing needs to be preserved
// between the stop and the resume.
func TestResumingAfterAStopIsAnOrdinaryRerun(t *testing.T) {
	cfg := config(t)
	cfg.StopAfterPhase = PhaseCreate
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("first run: %v", err)
	}
	before := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))

	cfg.StopAfterPhase = 0
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}

	if !result.Files[0].Existed {
		t.Errorf("the resumed run treated an existing file as new")
	}
	if after := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb")); after != before {
		t.Errorf("the resumed run rewrote the file")
	}
}

func TestDryRunWritesNothing(t *testing.T) {
	cfg := config(t)
	cfg.DryRun = true

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("a dry run reported %d files, want 1", len(result.Files))
	}
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Errorf("a dry run wrote %v", got)
	}
}

func TestOnlyRestrictsTheRunToNamedRecords(t *testing.T) {
	cfg := config(t)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		"prov-2026-bbbbbbbb": "  - app/models/order.rb\n",
	})
	cfg.Only = []string{"prov-2026-bbbbbbbb"}

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}

	want := []string{"app/models/order.rb"}
	if got := tree(t, cfg.Output); len(got) != 1 || got[0] != want[0] {
		t.Errorf("output tree = %v, want %v", got, want)
	}
}

// A phase that fails stops the run there. Phase 0 rejecting a package must not
// leave Phase 3 creating files under a half-loaded generators directory.
func TestARejectedPackageHaltsBeforeAnythingIsWritten(t *testing.T) {
	cfg := config(t)
	cfg.Generators = writeTree(t, map[string]string{
		"rails/sedum.yaml":           "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n",
		"rails/actions/actions.yaml": "actions:\n  broken:\n    anchor: nowhere\n",
	})

	_, err := Run(context.Background(), cfg)
	if err == nil {
		t.Fatal("a package with no template for a declared action was accepted")
	}
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Errorf("a failed load still wrote %v", got)
	}
}

func TestAnUnresolvablePathHaltsBeforeAnythingIsWritten(t *testing.T) {
	cfg := config(t)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - app/models/user.rb\n  - app/views/index.erb\n",
	})

	_, err := Run(context.Background(), cfg)
	if err == nil {
		t.Fatal("a path whose extension no package claims was accepted")
	}
	if !strings.Contains(err.Error(), ".erb") {
		t.Errorf("error = %q, does not name the unclaimed extension", err)
	}
	// Phase 2 resolves the whole path set before Phase 3 creates anything, so
	// one unresolvable path means no file is created at all rather than some.
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Errorf("a failed resolution still wrote %v", got)
	}
}

// Warnings are collected from every phase and handed back together, since the
// command is what decides where they go.
func TestWarningsAreCollected(t *testing.T) {
	cfg := config(t)
	cfg.Lang = []string{"nosuchpackage"}

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("an unusable --lang refused the run: %v", err)
	}

	var found bool
	for _, w := range result.Warnings {
		if strings.Contains(w, "nosuchpackage") {
			found = true
		}
	}
	if !found {
		t.Errorf("warnings = %v, do not mention the unusable --lang", result.Warnings)
	}
}

func TestMissingDirectoriesAreReported(t *testing.T) {
	t.Run("generators", func(t *testing.T) {
		cfg := config(t)
		cfg.Generators = filepath.Join(t.TempDir(), "absent")
		if _, err := Run(context.Background(), cfg); err == nil || !strings.Contains(err.Error(), "absent") {
			t.Errorf("error = %v, does not name the missing generators directory", err)
		}
	})
	t.Run("records", func(t *testing.T) {
		cfg := config(t)
		cfg.Records = filepath.Join(t.TempDir(), "absent")
		if _, err := Run(context.Background(), cfg); err == nil || !strings.Contains(err.Error(), "absent") {
			t.Errorf("error = %v, does not name the missing records directory", err)
		}
	})
}

const selectsAField = `{"invocations":[{"action":"addField","kwargs":{"model":"user","field":"email"}}]}`

// The whole loop, with no hand-written invocation list: a record's intent goes
// to a model, what comes back is validated, expanded, and injected, and the
// file on disk carries the region with its ownership marker.
func TestARecordIsGeneratedFromItsIntentAlone(t *testing.T) {
	cfg := config(t)
	client := &stub{response: selectsAField}
	cfg.Client = client

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if client.calls != 1 {
		t.Errorf("the model was called %d times for one record, want 1", client.calls)
	}
	if len(result.Selections) != 1 || len(result.Selections[0].Invocations) != 1 {
		t.Fatalf("Phase 5 produced %+v", result.Selections)
	}
	if len(result.Injections) != 1 {
		t.Fatalf("Phase 7 applied %d regions, want 1", len(result.Injections))
	}

	body := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))
	for _, want := range []string{"attribute :email", "sedum:addField", "/sedum:addField", `"record":"prov-2026-aaaaaaaa"`} {
		if !strings.Contains(body, want) {
			t.Errorf("the generated file does not carry %q:\n%s", want, body)
		}
	}
}

// Rerunning replaces the region an action owns rather than appending beside it.
// The recording and the sidecar cache both exist to be unnecessary, and this is
// why: idempotency state lives in the file.
func TestRerunningReplacesRatherThanDuplicates(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: selectsAField}

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("first run: %v", err)
	}
	first := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("second run: %v", err)
	}
	second := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))

	if first != second {
		t.Errorf("a rerun changed the file:\n%s\n---\n%s", first, second)
	}
	if got := strings.Count(second, "attribute :email"); got != 1 {
		t.Errorf("the region appears %d times after two runs, want 1:\n%s", got, second)
	}
}

// A dry run decides everything and writes nothing, including into files it
// declined to create. Phase 7 takes what Phase 3 rendered for those
// (prov-2026-23653fdc).
func TestDryRunReportsInjectionsWithoutWriting(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: selectsAField}
	cfg.DryRun = true

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Injections) != 1 {
		t.Fatalf("a dry run reported %d injections, want 1", len(result.Injections))
	}
	if got := tree(t, cfg.Output); len(got) != 0 {
		t.Errorf("a dry run wrote %v", got)
	}
}

// A dry run against a tree where the file already exists reads it from disk, so
// what it reports is a replacement rather than a first injection. Preferring
// the rendered template would describe a run nobody asked for.
func TestDryRunReadsFilesThatAreAlreadyThere(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: selectsAField}
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("first run: %v", err)
	}
	before := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))

	cfg.DryRun = true
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("dry run: %v", err)
	}

	if len(result.Injections) != 1 || !result.Injections[0].Replaced {
		t.Errorf("a dry run over an existing region reported %+v, want a replacement", result.Injections)
	}
	if after := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb")); after != before {
		t.Errorf("a dry run modified the file")
	}
}

// An action resolving to a path no record authorized fails validation, before a
// retry is spent and before anything is injected.
func TestAnUnauthorizedTargetHaltsWithNothingInjected(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: `{"invocations":[{"action":"addField","kwargs":{"model":"order","field":"total"}}]}`}

	_, err := Run(context.Background(), cfg)
	if err == nil {
		t.Fatal("an action targeting an unauthorized path was accepted")
	}
	if !strings.Contains(err.Error(), "app/models/order.rb") {
		t.Errorf("the error does not name the unauthorized path: %v", err)
	}

	// The authorized file was still created by Phase 3, and carries no region.
	body := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))
	if strings.Contains(body, "sedum:addField") {
		t.Errorf("a halted run still injected:\n%s", body)
	}
}

// Stopping after validation leaves the files scaffolded and nothing injected.
// It is the point a person reads what the model decided and corrects it.
func TestStopAfterInvocationsLeavesFilesUninjected(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: selectsAField}
	cfg.StopAfterPhase = PhaseValidate

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(result.Selections) != 1 || len(result.Selections[0].Invocations) != 1 {
		t.Fatalf("the validated invocations are not available: %+v", result.Selections)
	}
	if result.Injections != nil {
		t.Errorf("stopping after validation still injected %+v", result.Injections)
	}

	body := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb"))
	if strings.Contains(body, "sedum:addField") {
		t.Errorf("the file was injected into after a stop at validation:\n%s", body)
	}
}

// A record whose every path is unmanaged has an empty catalog, so the only
// valid answer is an empty list. Paying for a call to be told so is a cost the
// run can see is pointless before it is incurred.
func TestARecordWithNothingToWriteCostsNoModelCall(t *testing.T) {
	cfg := config(t)
	cfg.Generators = writeTree(t, map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
unmanaged:
  - Gemfile
`,
		"rails/actions/actions.yaml": "actions: {}\n",
	})
	cfg.Records = recordsDir(t, map[string]string{"prov-2026-cccccccc": "  - Gemfile\n"})
	client := &stub{}
	cfg.Client = client

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if client.calls != 0 {
		t.Errorf("a record with no writable path cost %d model calls", client.calls)
	}
	if len(result.Unmanaged) != 1 {
		t.Errorf("the unmanaged path was not reported as the run's handoff: %+v", result.Unmanaged)
	}
}

// The prompt the model receives is built from the record it is being asked
// about and the packages that record's paths resolved to, and nothing else.
func TestEachRecordGetsItsOwnCall(t *testing.T) {
	cfg := config(t)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		"prov-2026-bbbbbbbb": "  - app/models/order.rb\n",
	})
	client := &stub{}
	cfg.Client = client

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}

	// Two per record. The stub selects nothing, so every anchor the model
	// template plants is left unfilled and each record earns one completeness
	// re-prompt before its answer stands (prov-2026-6d87dc11). What this test
	// is about is that the records do not share a call, which the isolation
	// check below is what actually asserts.
	if client.calls != 4 {
		t.Fatalf("two records cost %d model calls, want 4 - one call and one completeness re-prompt each", client.calls)
	}
	// The catalog names paths too, in its injects_into patterns, so the file
	// list is what is counted rather than the prompt as a whole.
	for i, prompt := range client.prompts {
		user := prompt[len(prompt)-1].Content
		files, _, _ := strings.Cut(user, "## Actions")
		if strings.Count(files, "app/models/") != 1 {
			t.Errorf("call %d saw more than its own record's files:\n%s", i+1, files)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// Config.FullCatalog reaches selection.Options, the same field it names. By
// default the prompt is filtered to what the record authorized; the flag
// turns that off, end to end from the command's configuration.
func TestFullCatalogConfigReachesThePrompt(t *testing.T) {
	generators := writeTree(t, map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
`,
		"rails/files/app/models/{name}.rb":      "class {{name}}\n  # sedum:anchor:class_body\nend\n",
		"rails/files/app/controllers/{name}.rb": "class {{name}}Controller\n  # sedum:anchor:class_body\nend\n",
		"rails/actions/actions.yaml": `actions:
  addField:
    kwargs:
      model: { type: string, required: true }
      field: { type: string, required: true }
    injects_into: "app/models/{{model}}.rb"
    anchor: class_body

  addRoute:
    kwargs:
      controller: { type: string, required: true }
      path: { type: string, required: true }
    injects_into: "app/controllers/{{controller}}.rb"
    anchor: class_body
`,
		"rails/actions/addField.rb": "attribute :{{field}}\n",
		"rails/actions/addRoute.rb": "route {{path}}\n",
	})
	records := recordsDir(t, map[string]string{
		"prov-2026-dddddddd": "  - app/models/user.rb\n",
	})

	filtered := &stub{}
	cfg := Config{Generators: generators, Records: records, Output: t.TempDir(), Client: filtered}
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}
	filteredPrompt := filtered.prompts[0][len(filtered.prompts[0])-1].Content
	if strings.Contains(filteredPrompt, "addRoute") {
		t.Errorf("addRoute reached the prompt though no authorized path can reach it:\n%s", filteredPrompt)
	}

	full := &stub{}
	cfg.Client = full
	cfg.FullCatalog = true
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("Run with --full-catalog: %v", err)
	}
	fullPrompt := full.prompts[0][len(full.prompts[0])-1].Content
	if !strings.Contains(fullPrompt, "addRoute") {
		t.Errorf("cfg.FullCatalog did not reach the prompt:\n%s", fullPrompt)
	}
}

// prov-2026-e4bab53b (--prompt) and prov-2026-8602b702 (--prior-recording).
//
// Assumed symbols, none of which exist yet. setConfigField reports a missing
// one as a clear failure rather than a compile error:
//
//	Config.Prompt string               // caller text, applied to every selected record
//	Config.Prior  *recording.Recording // already-read prior recording; filtered here to this run's records
//	Result.Prompt string               // what Capture records
//
// The command reads and validates the file (missing, non-JSON, no records, the
// 16384-byte text cap) before calling Run; the 8192-byte per-record bound needs
// the run's records, so it is asserted here.

func setConfigField(t *testing.T, cfg *Config, name string, value any) {
	t.Helper()
	f := reflect.ValueOf(cfg).Elem().FieldByName(name)
	if !f.IsValid() {
		t.Fatalf("pipeline.Config has no %s field yet", name)
	}
	v := reflect.ValueOf(value)
	if !v.Type().AssignableTo(f.Type()) {
		t.Fatalf("pipeline.Config.%s is %s, but these tests assume %s", name, f.Type(), v.Type())
	}
	f.Set(v)
}

func userPrompt(s *stub, call int) string {
	m := s.prompts[call]
	return m[len(m)-1].Content
}

// recordPrompt returns the Phase 4 prompt the model was given for the record
// whose file is the one named. The prompt is told apart from a completeness
// re-prompt (prov-2026-6d87dc11) by what it contains, never by its position in
// the sequence of calls, so these tests hold however many calls a record costs.
// Only the part before the actions catalog is searched for the file, since the
// catalog names paths of its own.
func recordPrompt(t *testing.T, s *stub, file string) string {
	t.Helper()
	for _, m := range s.prompts {
		user := m[len(m)-1].Content
		files, _, found := strings.Cut(user, "## Actions")
		if found && strings.Contains(files, file) {
			return user
		}
	}
	t.Fatalf("the model was never given a prompt for %s", file)
	return ""
}

// twoRecordConfig selects prov-2026-aaaaaaaa (user.rb) and prov-2026-bbbbbbbb
// (order.rb), whose files the addField fixture can both reach.
func twoRecordConfig(t *testing.T) (Config, *stub) {
	t.Helper()
	cfg := config(t)
	cfg.Records = recordsDir(t, map[string]string{
		"prov-2026-aaaaaaaa": "  - app/models/user.rb\n",
		"prov-2026-bbbbbbbb": "  - app/models/order.rb\n",
	})
	client := &stub{}
	cfg.Client = client
	return cfg, client
}

func priorInvocation(record, field string) recording.Record {
	return recording.Record{
		RecordID: record,
		Phases: []recording.Phase{{
			Name: recording.DefaultPhase,
			Invocations: []recording.Invocation{{
				Action: "addField",
				Kwargs: map[string]any{"model": "user", "field": field},
			}},
		}},
	}
}

func priorOf(records ...recording.Record) *recording.Recording {
	return &recording.Recording{SedumVersion: recording.Version, Packages: map[string]recording.Package{}, Records: records}
}

func TestPromptReachesEveryRecordTheRunSelects(t *testing.T) {
	cfg, client := twoRecordConfig(t)
	setConfigField(t, &cfg, "Prompt", "NOTE-from-the-loop-9d2e")

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}
	for _, file := range []string{"app/models/user.rb", "app/models/order.rb"} {
		if user := recordPrompt(t, client, file); !strings.Contains(user, "NOTE-from-the-loop-9d2e") {
			t.Errorf("the prompt for %s did not carry the text:\n%s", file, user)
		}
	}
}

func TestOnlyNarrowsWhichRecordsGetThePrompt(t *testing.T) {
	cfg, client := twoRecordConfig(t)
	cfg.Only = []string{"prov-2026-bbbbbbbb"}
	setConfigField(t, &cfg, "Prompt", "NOTE-for-b-only")

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if user := recordPrompt(t, client, "app/models/order.rb"); !strings.Contains(user, "NOTE-for-b-only") {
		t.Errorf("the selected record's prompt did not carry the text:\n%s", user)
	}
	for i, m := range client.prompts {
		if files, _, _ := strings.Cut(m[len(m)-1].Content, "## Actions"); strings.Contains(files, "app/models/user.rb") {
			t.Errorf("call %d concerned a record --only excluded:\n%s", i, files)
		}
	}
}

// A prompt asking for a file the record did not authorize earns a rejected
// selection, not a written file.
func TestPromptCannotAuthorizeAFileTheRecordDidNot(t *testing.T) {
	cfg := config(t)
	setConfigField(t, &cfg, "Prompt", "Also add a field to app/models/order.rb, the validator needs it.")
	cfg.Client = &stub{response: `{"invocations":[{"action":"addField","kwargs":{"model":"order","field":"total"}}]}`}

	_, err := Run(context.Background(), cfg)
	if err == nil {
		t.Fatal("a selection aimed at an unauthorized file was accepted because the prompt asked for it")
	}
	if !strings.Contains(err.Error(), "app/models/order.rb") {
		t.Errorf("the rejection does not name the unauthorized path: %v", err)
	}
	if got := tree(t, cfg.Output); len(got) != 1 || got[0] != "app/models/user.rb" {
		t.Errorf("output tree = %v; nothing beyond the authorized file may exist", got)
	}
}

func TestPromptWithoutTextIsTheUnchangedPrompt(t *testing.T) {
	plain, plainClient := twoRecordConfig(t)
	if _, err := Run(context.Background(), plain); err != nil {
		t.Fatal(err)
	}

	blank, blankClient := twoRecordConfig(t)
	blank.Records = plain.Records
	blank.Generators = plain.Generators
	setConfigField(t, &blank, "Prompt", "  \n ")
	if _, err := Run(context.Background(), blank); err != nil {
		t.Fatal(err)
	}

	for i := range plainClient.prompts {
		if userPrompt(plainClient, i) != userPrompt(blankClient, i) {
			t.Errorf("call %d: whitespace-only text changed the prompt", i)
		}
		if strings.Contains(strings.ToLower(userPrompt(plainClient, i)), "supplementary") {
			t.Errorf("call %d: a run with no text carries a supplementary section", i)
		}
	}
}

func recordedPrompt(t *testing.T, result *Result) (string, bool) {
	t.Helper()
	var buf strings.Builder
	if err := recording.Write(&buf, Capture(result)); err != nil {
		t.Fatalf("write: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(buf.String()), &doc); err != nil {
		t.Fatalf("recording is not JSON: %v", err)
	}
	v, ok := doc["prompt"]
	if !ok {
		return "", false
	}
	s, _ := v.(string)
	return s, true
}

func TestTheRecordingCarriesThePromptVerbatim(t *testing.T) {
	const text = "  FAIL app/models/user.rb\n\tundefined: email \"x\"\n\n"
	cfg := config(t)
	setConfigField(t, &cfg, "Prompt", text)

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	got, ok := recordedPrompt(t, result)
	if !ok || got != text {
		t.Errorf("recorded prompt = %q (present: %v), want %q verbatim", got, ok, text)
	}
}

func TestARunWithoutAPromptRecordsNone(t *testing.T) {
	result, err := Run(context.Background(), config(t))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got, ok := recordedPrompt(t, result); ok {
		t.Errorf("a run with no prompt recorded one: %q", got)
	}
}

func TestPriorSelectionsAreFilteredToTheRecordsInThisRun(t *testing.T) {
	cfg, client := twoRecordConfig(t)
	cfg.Only = []string{"prov-2026-aaaaaaaa"}
	setConfigField(t, &cfg, "Prior", priorOf(
		priorInvocation("prov-2026-aaaaaaaa", "kept-field-a1"),
		priorInvocation("prov-2026-bbbbbbbb", "other-record-field-b2"),
		priorInvocation("prov-2026-cccccccc", "absent-record-field-c3"),
	))

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}
	user := recordPrompt(t, client, "app/models/user.rb")
	if !strings.Contains(user, "kept-field-a1") {
		t.Errorf("this record's own prior selection is missing:\n%s", user)
	}
	for _, leaked := range []string{"other-record-field-b2", "absent-record-field-c3"} {
		if strings.Contains(user, leaked) {
			t.Errorf("another record's prior selection %q reached the prompt:\n%s", leaked, user)
		}
	}
}

// Without --only each record sees its own prior selections and nobody else's.
func TestEachRecordSeesOnlyItsOwnPriorSelections(t *testing.T) {
	cfg, client := twoRecordConfig(t)
	setConfigField(t, &cfg, "Prior", priorOf(
		priorInvocation("prov-2026-aaaaaaaa", "field-for-a"),
		priorInvocation("prov-2026-bbbbbbbb", "field-for-b"),
	))
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}
	for file, f := range map[string][2]string{
		"app/models/user.rb":  {"field-for-a", "field-for-b"},
		"app/models/order.rb": {"field-for-b", "field-for-a"},
	} {
		own, other := f[0], f[1]
		user := recordPrompt(t, client, file)
		if !strings.Contains(user, own) || strings.Contains(user, other) {
			t.Errorf("the prompt for %s should carry %q and not %q:\n%s", file, own, other, user)
		}
	}
}

// A prior recording holding none of this run's records adds nothing, and is
// not an error.
func TestAPriorRecordingForOtherRecordsAddsNothing(t *testing.T) {
	baseline := config(t)
	baselineClient := &stub{}
	baseline.Client = baselineClient
	if _, err := Run(context.Background(), baseline); err != nil {
		t.Fatal(err)
	}

	cfg := config(t)
	client := &stub{}
	cfg.Client = client
	cfg.Generators, cfg.Records = baseline.Generators, baseline.Records
	setConfigField(t, &cfg, "Prior", priorOf(priorInvocation("prov-2026-zzzzzzzz", "elsewhere")))

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("a prior recording with none of this run's records was an error: %v", err)
	}
	if userPrompt(client, 0) != userPrompt(baselineClient, 0) {
		t.Errorf("unrelated prior selections changed the prompt:\n%s", userPrompt(client, 0))
	}
}

func TestPriorSelectionsAndPromptShareOneSupplementarySection(t *testing.T) {
	cfg := config(t)
	client := &stub{}
	cfg.Client = client
	setConfigField(t, &cfg, "Prompt", "TEXT-5b7e")
	setConfigField(t, &cfg, "Prior", priorOf(priorInvocation("prov-2026-aaaaaaaa", "prior-field-4c4c")))

	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}
	user := userPrompt(client, 0)
	labelled := 0
	for _, line := range strings.Split(user, "\n") {
		if strings.HasPrefix(line, "#") && strings.Contains(strings.ToLower(line), "supplementary") {
			labelled++
		}
	}
	if labelled != 1 {
		t.Errorf("want one supplementary heading, got %d:\n%s", labelled, user)
	}
	for _, want := range []string{"TEXT-5b7e", "prior-field-4c4c"} {
		if !strings.Contains(user, want) {
			t.Errorf("prompt lacks %q:\n%s", want, user)
		}
	}
}

// The prior recording's own prompt is not fed back in, and is not carried into
// this run's recording as though this run had been given it.
func TestThePriorRecordingsOwnPromptIsNotFedBack(t *testing.T) {
	cfg := config(t)
	client := &stub{}
	cfg.Client = client
	prior := priorOf(priorInvocation("prov-2026-aaaaaaaa", "prior-field"))
	var buf strings.Builder
	if err := recording.Write(&buf, *prior); err != nil {
		t.Fatal(err)
	}
	withPrompt, err := recording.Read(strings.NewReader(
		strings.Replace(buf.String(), "{\n", "{\n  \"prompt\": \"STALE-PRIOR-PROMPT-e0e0\",\n", 1)))
	if err != nil {
		t.Fatal(err)
	}
	setConfigField(t, &cfg, "Prior", &withPrompt)

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if strings.Contains(userPrompt(client, 0), "STALE-PRIOR-PROMPT-e0e0") {
		t.Error("the prior recording's prompt was fed back into this run's prompt")
	}
	if got, ok := recordedPrompt(t, result); ok {
		t.Errorf("this run recorded a prompt it was never given: %q", got)
	}
}

// Context only. The model is still called and its answer is what is applied;
// the prior invocations are never executed or trusted.
func TestAPriorRecordingIsContextAndNeverReplayed(t *testing.T) {
	cfg := config(t)
	client := &stub{} // answers with no invocations
	cfg.Client = client
	setConfigField(t, &cfg, "Prior", priorOf(priorInvocation("prov-2026-aaaaaaaa", "must-not-be-written")))

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	// The prior is shown to the model as context, which is what makes the
	// absence of a replay below mean something.
	if user := recordPrompt(t, client, "app/models/user.rb"); !strings.Contains(user, "must-not-be-written") {
		t.Errorf("the prior selection never reached the model as context:\n%s", user)
	}
	if len(result.Selections) != 1 || len(result.Selections[0].Invocations) != 0 {
		t.Errorf("selections = %+v; the model chose nothing and the prior must not stand in for it", result.Selections)
	}
	if body := readFile(t, filepath.Join(cfg.Output, "app/models/user.rb")); strings.Contains(body, "must-not-be-written") {
		t.Errorf("a prior invocation was written:\n%s", body)
	}
	if len(result.Injections) != 0 {
		t.Errorf("%d regions injected from a run whose model chose nothing", len(result.Injections))
	}
}

// A prior selection that would be rejected today is still only context: the
// model's answer is validated as ever, and an unauthorized prior path does not
// become authorized by having been shown.
func TestPriorSelectionsCannotAuthorizeAFileTheRecordDidNot(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: `{"invocations":[{"action":"addField","kwargs":{"model":"order","field":"total"}}]}`}
	prior := priorInvocation("prov-2026-aaaaaaaa", "total")
	prior.Phases[0].Invocations[0].Kwargs["model"] = "order"
	setConfigField(t, &cfg, "Prior", priorOf(prior))

	if _, err := Run(context.Background(), cfg); err == nil {
		t.Fatal("a selection aimed at an unauthorized file was accepted because an earlier run had made it")
	}
	if got := tree(t, cfg.Output); len(got) != 1 {
		t.Errorf("output tree = %v", got)
	}
}

// 8192 bytes of serialized invocations per record, with an error: truncating
// would drop the tail, which is the part that says what went wrong. It halts
// before any model call has been paid for.
func TestPriorSelectionsOverEightKiBPerRecordAreAnError(t *testing.T) {
	cfg := config(t)
	client := &stub{}
	cfg.Client = client
	setConfigField(t, &cfg, "Prior", priorOf(priorInvocation("prov-2026-aaaaaaaa", strings.Repeat("y", 9000))))

	_, err := Run(context.Background(), cfg)
	if err == nil {
		t.Fatal("9000 bytes of prior selections for one record was accepted")
	}
	for _, want := range []string{"8192", "prov-2026-aaaaaaaa"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error does not name %q: %v", want, err)
		}
	}
	if client.calls != 0 {
		t.Errorf("the model was called %d times before the overflow was reported", client.calls)
	}
}

// The bound is per record, not per run: two records each under it are fine even
// though together they are over.
func TestTheEightKiBBoundIsPerRecord(t *testing.T) {
	cfg, client := twoRecordConfig(t)
	setConfigField(t, &cfg, "Prior", priorOf(
		priorInvocation("prov-2026-aaaaaaaa", strings.Repeat("a", 5000)),
		priorInvocation("prov-2026-bbbbbbbb", strings.Repeat("b", 5000)),
	))
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("two records each under 8192 bytes were rejected: %v", err)
	}
	for file, own := range map[string]string{
		"app/models/user.rb":  strings.Repeat("a", 5000),
		"app/models/order.rb": strings.Repeat("b", 5000),
	} {
		if user := recordPrompt(t, client, file); !strings.Contains(user, own) {
			t.Errorf("the prompt for %s lacks its own prior selections", file)
		}
	}
}
