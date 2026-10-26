package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/selection"
)

// A re-selected record is shown its earlier selections as context
// (prov-2026-cc73cf61, revised): those --prior-recording supplies, else the
// invocations in the phases of its own recording, else the selections rebuilt
// from its markers. Any of them keeps only exposed actions and the kwargs the
// action declares, and one over selection.MaxPriorBytes is left out.

// v2Reselect grows once, then returns a config whose record has changed, so the
// next grow re-selects it. The model is asked for prompts to inspect.
func v2Reselect(t *testing.T) (Config, string) {
	t.Helper()
	cfg, dir := v2Grow(t)
	editRecord(t, cfg, "Something else.")
	return cfg, dir
}

// v2PriorPrompt runs the re-selecting grow and returns the prompt the model was
// given for the record.
func v2PriorPrompt(t *testing.T, cfg Config) string {
	t.Helper()
	client := &stub{response: fields("phone")}
	cfg.Client = client
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("grow: %v", err)
	}
	if v2Class(t, result, v2RecordA).Mode != ModeReselect {
		t.Fatalf("classifications = %+v, want a re-select", result.Classifications)
	}
	if client.calls < 1 {
		t.Fatal("the model was never asked")
	}
	return userPrompt(client, 0)
}

func setPhases(invocations ...recording.Invocation) func(*recording.Record) {
	return func(r *recording.Record) {
		r.Phases = []recording.Phase{{Name: recording.DefaultPhase, Invocations: invocations}}
	}
}

func addFieldInvocation(field string) recording.Invocation {
	return recording.Invocation{Action: "addField", Kwargs: map[string]any{"model": "user", "field": field}}
}

func TestSelectAllPriorFromRecordingUsesThePhasesOfTheRecordsOwnRecording(t *testing.T) {
	cfg, dir := v2Reselect(t)
	v2EditRecording(t, dir, v2RecordA, setPhases(addFieldInvocation("phase-prior-4c4c")))

	prompt := v2PriorPrompt(t, cfg)
	if !strings.Contains(prompt, `"field":"phase-prior-4c4c"`) {
		t.Errorf("the recording's phases did not reach the prompt:\n%s", prompt)
	}
	if strings.Contains(prompt, `"field":"email"`) {
		t.Errorf("the markers' rebuilt selection was used beside the recording's:\n%s", prompt)
	}
}

func TestSelectAllPriorFromRecordingTakesEveryPhasesInvocations(t *testing.T) {
	cfg, dir := v2Reselect(t)
	v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) {
		r.Phases = []recording.Phase{
			{Name: "first", Invocations: []recording.Invocation{addFieldInvocation("first-phase-1a")}},
			{Name: "second", Invocations: []recording.Invocation{addFieldInvocation("second-phase-2b")}},
		}
	})

	prompt := v2PriorPrompt(t, cfg)
	for _, want := range []string{"first-phase-1a", "second-phase-2b"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("%s is missing from the prompt:\n%s", want, prompt)
		}
	}
	if strings.Index(prompt, "first-phase-1a") > strings.Index(prompt, "second-phase-2b") {
		t.Errorf("the phases' invocations are out of order:\n%s", prompt)
	}
}

func TestSelectAllPriorFromRecordingReadsTheRecordingsDirectoryOfTheRun(t *testing.T) {
	cfg := config(t)
	cfg.Recordings = filepath.Join(t.TempDir(), "elsewhere")
	first, _ := v2Run(t, cfg, fields("email"))
	if err := WriteRecordings(first, cfg.Recordings); err != nil {
		t.Fatal(err)
	}
	v2EditRecording(t, cfg.Recordings, v2RecordA, setPhases(addFieldInvocation("custom-dir-prior-7d")))
	editRecord(t, cfg, "Something else.")

	if prompt := v2PriorPrompt(t, cfg); !strings.Contains(prompt, "custom-dir-prior-7d") {
		t.Errorf("the recording in the run's recordings directory was not read:\n%s", prompt)
	}
}

func TestSelectAllPriorFromRecordingFallsBackToMarkersWhenThePhasesAreEmpty(t *testing.T) {
	for name, phases := range map[string][]recording.Phase{
		"no phases":            nil,
		"a phase with nothing": {{Name: recording.DefaultPhase}},
		"phases with nothing":  {{Name: "a"}, {Name: "b", Invocations: []recording.Invocation{}}},
	} {
		t.Run(name, func(t *testing.T) {
			cfg, dir := v2Reselect(t)
			v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) { r.Phases = phases })

			prompt := v2PriorPrompt(t, cfg)
			if !strings.Contains(prompt, `"field":"email"`) {
				t.Errorf("the prior selection was not rebuilt from the markers:\n%s", prompt)
			}
		})
	}
}

func TestSelectAllPriorFromRecordingFallsBackToMarkersWhenThereIsNoUsableRecording(t *testing.T) {
	for name, mutate := range map[string]func(t *testing.T, dir string){
		"missing": func(t *testing.T, dir string) {
			if err := os.RemoveAll(dir); err != nil {
				t.Fatal(err)
			}
		},
		"unreadable": func(t *testing.T, dir string) {
			if err := os.WriteFile(filepath.Join(dir, v2RecordA+".json"), []byte("not json"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
		"holding no such record": func(t *testing.T, dir string) {
			data := `{"sedum_version":"0.8.0","packages":{},"records":[{"record_id":"prov-2026-zzzzzzzz","files":[],"phases":[]}]}`
			if err := os.WriteFile(filepath.Join(dir, v2RecordA+".json"), []byte(data), 0o644); err != nil {
				t.Fatal(err)
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			cfg, dir := v2Reselect(t)
			mutate(t, dir)

			// The run does not fail on it: the recording is context, and
			// without it the markers still say what was selected.
			prompt := v2PriorPrompt(t, cfg)
			if !strings.Contains(prompt, `"field":"email"`) {
				t.Errorf("the prior selection was not rebuilt from the markers:\n%s", prompt)
			}
		})
	}
}

func TestSelectAllPriorFromRecordingAPriorRecordingStillWins(t *testing.T) {
	cfg, dir := v2Reselect(t)
	v2EditRecording(t, dir, v2RecordA, setPhases(addFieldInvocation("phase-prior-loses")))
	setConfigField(t, &cfg, "Prior", priorOf(priorInvocation(v2RecordA, "explicit-prior-9f2e")))

	prompt := v2PriorPrompt(t, cfg)
	if !strings.Contains(prompt, "explicit-prior-9f2e") {
		t.Errorf("--prior-recording did not reach the prompt:\n%s", prompt)
	}
	for _, loser := range []string{"phase-prior-loses", `"field":"email"`} {
		if strings.Contains(prompt, loser) {
			t.Errorf("%s reached the prompt beside an explicit prior recording:\n%s", loser, prompt)
		}
	}
}

func TestSelectAllPriorFromRecordingAPriorRecordingNamingOtherRecordsLeavesTheRecordingToSupply(t *testing.T) {
	cfg, dir := v2Reselect(t)
	v2EditRecording(t, dir, v2RecordA, setPhases(addFieldInvocation("phase-prior-5a5a")))
	setConfigField(t, &cfg, "Prior", priorOf(priorInvocation("prov-2026-zzzzzzzz", "elsewhere-3c3c")))

	prompt := v2PriorPrompt(t, cfg)
	if !strings.Contains(prompt, "phase-prior-5a5a") || strings.Contains(prompt, "elsewhere-3c3c") {
		t.Errorf("a prior recording that does not name the record displaced its own recording:\n%s", prompt)
	}
}

func TestSelectAllPriorFromRecordingKeepsOnlyExposedActionsAndDeclaredKwargs(t *testing.T) {
	cfg, dir := v2Reselect(t)
	v2EditRecording(t, dir, v2RecordA, setPhases(
		recording.Invocation{Action: "addField", Kwargs: map[string]any{
			"model": "user", "field": "declared-field-6e6e", "undeclared": "undeclared-value-1b1b",
		}},
		recording.Invocation{Action: "noSuchAction", Kwargs: map[string]any{"field": "unknown-action-2c2c"}},
	))

	prompt := v2PriorPrompt(t, cfg)
	if !strings.Contains(prompt, "declared-field-6e6e") {
		t.Errorf("the declared kwargs of an exposed action were dropped:\n%s", prompt)
	}
	for _, leaked := range []string{"undeclared-value-1b1b", "unknown-action-2c2c", "noSuchAction"} {
		if strings.Contains(prompt, leaked) {
			t.Errorf("%q reached the prompt:\n%s", leaked, prompt)
		}
	}
}

func TestSelectAllPriorFromRecordingKeepsOnlyExposedActionsFromAHiddenAction(t *testing.T) {
	cfg := config(t)
	cfg.Generators = writeTree(t, map[string]string{
		"rails/sedum.yaml":                 "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\ntransforms:\n  constantize: [singular, pascal]\n",
		"rails/files/app/models/{name}.rb": "class {{name|constantize}}\n  # sedum:anchor:class_body\nend\n",
		"rails/actions/actions.yaml": `actions:
  addField:
    kwargs:
      model: { type: string, required: true }
      field: { type: string, required: true }
    injects_into: "app/models/{{model|snake}}.rb"
    anchor: class_body
  addHidden:
    exposed: false
    kwargs:
      model: { type: string, required: true }
      field: { type: string, required: true }
    injects_into: "app/models/{{model|snake}}.rb"
    anchor: class_body
`,
		"rails/actions/addField.rb":  "attribute :{{field}}\n",
		"rails/actions/addHidden.rb": "hidden :{{field}}\n",
	})
	first, _ := v2Run(t, cfg, fields("email"))
	dir := v2Dir(cfg)
	if err := WriteRecordings(first, dir); err != nil {
		t.Fatal(err)
	}
	v2EditRecording(t, dir, v2RecordA, setPhases(
		addFieldInvocation("exposed-field-8a8a"),
		recording.Invocation{Action: "addHidden", Kwargs: map[string]any{"model": "user", "field": "hidden-field-9b9b"}},
	))
	editRecord(t, cfg, "Something else.")

	prompt := v2PriorPrompt(t, cfg)
	if !strings.Contains(prompt, "exposed-field-8a8a") {
		t.Errorf("the exposed action was left out:\n%s", prompt)
	}
	if strings.Contains(prompt, "hidden-field-9b9b") {
		t.Errorf("an unexposed action reached the prompt:\n%s", prompt)
	}
}

func TestSelectAllPriorFromRecordingOverTheBoundIsLeftOutRatherThanFailingTheRun(t *testing.T) {
	cfg, dir := v2Reselect(t)
	big := strings.Repeat("z", selection.MaxPriorBytes+1)
	v2EditRecording(t, dir, v2RecordA, setPhases(addFieldInvocation(big)))

	client := &stub{response: fields("phone")}
	cfg.Client = client
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("a prior selection over %d bytes failed the run: %v", selection.MaxPriorBytes, err)
	}
	if client.calls != 1 || v2Class(t, result, v2RecordA).Mode != ModeReselect {
		t.Fatalf("calls = %d, classifications = %+v", client.calls, result.Classifications)
	}
	if prompt := userPrompt(client, 0); strings.Contains(prompt, "zzzzzzzzzz") {
		t.Errorf("the oversize prior selection reached the prompt")
	}
}

func TestSelectAllPriorFromRecordingIsOnlyContextAndIsNeverApplied(t *testing.T) {
	cfg, dir := v2Reselect(t)
	v2EditRecording(t, dir, v2RecordA, setPhases(addFieldInvocation("never-applied-4d4d")))

	client := &stub{response: fields("phone")}
	cfg.Client = client
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("grow: %v", err)
	}
	got := userModel(t, cfg)
	if strings.Contains(got, "never-applied-4d4d") || !strings.Contains(got, "attribute :phone") {
		t.Errorf("the model's selection was not what was applied:\n%s", got)
	}
}

func TestSelectAllPriorFromRecordingOnlyAppliesToReselectedRecords(t *testing.T) {
	// A skipped record makes no call, so nothing reads its recording for a
	// prompt; a new record has no recording and no markers.
	cfg, dir := v2Grow(t)
	v2EditRecording(t, dir, v2RecordA, setPhases(addFieldInvocation("unused-prior-0e0e")))

	client := &stub{response: fields("email")}
	cfg.Client = client
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("grow: %v", err)
	}
	if client.calls != 0 {
		t.Errorf("a skipped record made %d model call(s)", client.calls)
	}

	fresh := config(t)
	freshClient := &stub{response: fields("email")}
	fresh.Client = freshClient
	if _, err := Run(context.Background(), fresh); err != nil {
		t.Fatalf("grow: %v", err)
	}
	if prompt := userPrompt(freshClient, 0); strings.Contains(prompt, "Selections made for this change on an earlier run") {
		t.Errorf("a new record's prompt carries prior selections:\n%s", prompt)
	}
}

// priorFromMarkers over a v2 region with no recording yields only the identity
// kwargs of that region that its action declares.
func TestSelectAllPriorFromRecordingPriorFromMarkersOverAV2RegionYieldsItsIdentityKwargs(t *testing.T) {
	cfg := config(t)
	cfg.Generators = v2TypedGenerators(t)
	response := `{"invocations": [{"action": "addTyped", "kwargs": {"model": "user", "field": "age", "kind": "integer"}}]}`
	result, _ := v2Run(t, cfg, response)
	dir := v2Dir(cfg)
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	env := v2Load(t, cfg)

	got := priorFromMarkers(v2RecordA, env.current(t, cfg), env.owners)
	if len(got) != 1 || got[0].Action != "addTyped" {
		t.Fatalf("priorFromMarkers = %+v, want the one addTyped invocation", got)
	}
	if kw := jsonOf(t, got[0].Kwargs); kw != `{"field":"age","model":"user"}` {
		t.Errorf("kwargs = %s, want only the identity kwargs the action declares (kind is on the recording, not the marker)", kw)
	}

	// With no recording, a re-selected record is shown just those.
	editRecord(t, cfg, "Something else.")
	client := &stub{response: response}
	cfg.Client = client
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("grow: %v", err)
	}
	if prompt := userPrompt(client, 0); !strings.Contains(prompt, `"field":"age"`) || strings.Contains(prompt, `"kind":"integer"`) {
		t.Errorf("prior rebuilt from a v2 marker:\n%s", prompt)
	}
}
