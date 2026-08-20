package evals

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/livecodelife/sedum/internal/recording"
)

// Behaviour is the only measurement here that runs somebody else's code, so
// what is tested is the seam rather than the applications: that a sample's own
// invocations reach the harness unchanged, that the three outcomes stay apart,
// and that a run which was never measured never reads as a run where nothing
// worked (prov-2026-83340ba0).
//
// Standing up Rails and Postgres inside `go test ./...` is not done. The
// harness is exercised end to end by running it, which is what evals/behavior's
// README documents.

func behaviorSample(b *BehaviorRun) Sample {
	return Sample{Counts: map[string]int{"addColumn": 1}, Behavior: b}
}

func caseWithTarget(m *Measurement, target string) {
	m.Case.ID = "fixture"
	m.Case.Arm = "sedum"
	m.Case.Expect.Actions = map[string]int{"addColumn": 1}
	m.Case.Expect.Behavior = &BehaviorExpectation{Target: target, PassRate: 1.0}
}

// The three outcomes are counted apart. A service that never booted and one
// that booted and answered wrongly have different fixes, and a single pass rate
// would report a broken generator package and a wrong one as the same number.
func TestBehaviorKeepsBrokenApartFromWrong(t *testing.T) {
	m := Measurement{Samples: []Sample{
		behaviorSample(&BehaviorRun{Outcome: "ok", Checks: 20, Passed: 20}),
		behaviorSample(&BehaviorRun{Outcome: "checks_failed", Checks: 20, Passed: 17,
			Failed: []string{"partial update keeps the title it never sent", "an empty update is a 400"}}),
		behaviorSample(&BehaviorRun{Outcome: "failed", FailedPhase: "build",
			Detail: "db/todos.go:61:24: undefined: models"}),
	}}
	caseWithTarget(&m, "todo-rails")

	got, measured := m.Behavior()
	if !measured {
		t.Fatal("three behaviour runs reported as not measured")
	}
	if got.Measured != 3 {
		t.Errorf("measured %d, want 3", got.Measured)
	}
	if got.Working != 1 || got.Disagreed != 1 || got.Broke != 1 {
		t.Errorf("working/disagreed/broke = %d/%d/%d, want 1/1/1",
			got.Working, got.Disagreed, got.Broke)
	}
	if got.Phases["build"] != 1 {
		t.Errorf("the failed phase was not counted: %v", got.Phases)
	}
	if got.Failures["an empty update is a 400"] != 1 {
		t.Errorf("failed assertions were not counted per assertion: %v", got.Failures)
	}
	// The reason a phase died survives the project being deleted. Without it a
	// failure says "build" and the log that said why is gone
	// (prov-2026-93829987).
	if got.Details["db/todos.go:61:24: undefined: models"] != 1 {
		t.Errorf("the failing phase's reason was not carried: %v", got.Details)
	}
	if got.Rate() != 1.0/3.0 {
		t.Errorf("rate = %v, want one in three", got.Rate())
	}
}

// A sample the harness could not be run for at all is excluded rather than
// counted as a failure - the same rule that keeps an unreachable endpoint out
// of the selection denominator.
func TestAHarnessErrorIsNotABehaviorFailure(t *testing.T) {
	m := Measurement{Samples: []Sample{
		behaviorSample(&BehaviorRun{Outcome: "ok", Checks: 5, Passed: 5}),
		behaviorSample(&BehaviorRun{Err: os.ErrNotExist}),
	}}
	caseWithTarget(&m, "todo-rails")

	got, _ := m.Behavior()
	if got.Measured != 1 || got.Errored != 1 {
		t.Errorf("measured/errored = %d/%d, want 1/1", got.Measured, got.Errored)
	}
	if got.Rate() != 1 {
		t.Errorf("rate = %v; an errored run must not drag the rate down", got.Rate())
	}
}

// Not measured and measured-with-nothing-working are different states, and the
// report has to say which. A reader who cannot tell them apart reads the
// absence of a number as the absence of a problem.
func TestNotMeasuredIsNotZero(t *testing.T) {
	m := Measurement{Samples: []Sample{{Counts: map[string]int{"addColumn": 1}}}}
	caseWithTarget(&m, "todo-rails")

	if _, measured := m.Behavior(); measured {
		t.Error("a run with no behaviour samples reported as measured")
	}

	var out strings.Builder
	behavior(&out, m)
	if !strings.Contains(out.String(), "not measured") {
		t.Errorf("the report does not say behaviour went unmeasured:\n%s", out.String())
	}
	if strings.Contains(out.String(), "working:") {
		t.Errorf("the report printed a working rate for a run that measured nothing:\n%s", out.String())
	}
}

// The report names which assertions broke. A rate that says something is wrong
// without saying what is not the finding this harness exists to produce.
func TestTheReportNamesTheBrokenContracts(t *testing.T) {
	m := Measurement{Samples: []Sample{
		behaviorSample(&BehaviorRun{Outcome: "checks_failed", Checks: 19, Passed: 16,
			Failed: []string{"a todo with no title is rejected"}, Elapsed: 4 * time.Second}),
	}}
	caseWithTarget(&m, "todo-chi")

	var out strings.Builder
	behavior(&out, m)
	got := out.String()
	for _, want := range []string{"todo-chi", "disagreed", "a todo with no title is rejected", "16 of 19"} {
		if !strings.Contains(got, want) {
			t.Errorf("the behaviour table does not mention %q:\n%s", want, got)
		}
	}
}

// The answer file is the envelope Phase 4 decodes, carrying the sample's own
// bindings. A kwarg dropped or renamed on the way through would make the
// harness measure a selection nobody made.
func TestTheAnswerCarriesTheSamplesOwnBindings(t *testing.T) {
	path, err := writeAnswer(invocationsFixture())
	if err != nil {
		t.Fatalf("writeAnswer: %v", err)
	}
	defer os.Remove(path)

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the answer: %v", err)
	}

	var body struct {
		Invocations []struct {
			Action string         `json:"action"`
			Kwargs map[string]any `json:"kwargs"`
		} `json:"invocations"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("the answer is not the envelope Phase 4 decodes: %v\n%s", err, raw)
	}
	if len(body.Invocations) != 2 {
		t.Fatalf("wrote %d invocations, want 2", len(body.Invocations))
	}
	if body.Invocations[0].Action != "addColumn" {
		t.Errorf("action = %q, want addColumn", body.Invocations[0].Action)
	}
	if body.Invocations[0].Kwargs["name"] != "title" {
		t.Errorf("kwargs did not survive: %v", body.Invocations[0].Kwargs)
	}
	// The omission is the answer for a defaulted kwarg (prov-2026-f03916ba), so
	// it has to still be an omission by the time the harness replays it.
	if _, bound := body.Invocations[0].Kwargs["default"]; bound {
		t.Errorf("an unbound kwarg acquired a value on the way to the harness: %v",
			body.Invocations[0].Kwargs)
	}
}

// Nothing to apply is not a measurement. A caller is expected not to ask, and
// asking gets an error rather than a zero-check pass.
func TestAnEmptySelectionIsNotMeasured(t *testing.T) {
	got := RunBehavior(t.Context(), "todo-rails", nil, nil)
	if got.Err == nil {
		t.Error("applying no invocations reported a result")
	}
	if got.Working() {
		t.Error("an unmeasurable sample reported as working")
	}
}

// invocationsFixture is one sample's worth of bindings, with the defaulted
// kwarg left unbound as a correct answer leaves it.
func invocationsFixture() []recording.Invocation {
	return []recording.Invocation{
		{Action: "addColumn", Kwargs: map[string]any{
			"resource": "todos", "stamp": "20260814000000",
			"name": "title", "type": "string", "nullable": false,
		}},
		{Action: "addColumn", Kwargs: map[string]any{
			"resource": "todos", "stamp": "20260814000000",
			"name": "completed", "type": "boolean", "nullable": false, "default": "false",
		}},
	}
}

// A phase failure prints why, not only where, and identical reasons across
// samples collapse to one finding with a count - the same rule the failed
// assertions follow (prov-2026-93829987).
func TestTheReportSaysWhyAPhaseDied(t *testing.T) {
	broke := func() Sample {
		return behaviorSample(&BehaviorRun{Outcome: "failed", FailedPhase: "build",
			Detail: "# todo/db\ndb/todos.go:85:3: invalid character U+0024 '$'"})
	}
	m := Measurement{Samples: []Sample{broke(), broke()}}
	caseWithTarget(&m, "todo-chi")

	var out strings.Builder
	behavior(&out, m)
	got := out.String()

	for _, want := range []string{"broke: 2", "build x2", "2 of 2:", "invalid character U+0024"} {
		if !strings.Contains(got, want) {
			t.Errorf("the behaviour table does not mention %q:\n%s", want, got)
		}
	}
	// One reason repeated is one finding, so it is printed once.
	if strings.Count(got, "invalid character U+0024") != 1 {
		t.Errorf("an identical reason was printed per sample rather than counted:\n%s", got)
	}
}

// A phase that will not finish must fail the sample and not the run. The
// deadline signals the whole process group, because bash defers a signal while a
// foreground command is running and behave.sh runs every phase in the
// foreground (prov-2026-3957eed2).
func TestABehaviourRunThatPassesItsDeadlineFailsTheSample(t *testing.T) {
	t.Setenv("SEDUM_BEHAVIOR_DEADLINE", "3s")

	started := time.Now()
	run := runHarness(context.Background(), "todo-chi", nil, nil, started)
	elapsed := time.Since(started)

	if run.Outcome != "failed" {
		t.Fatalf("outcome %q with err %v, want failed", run.Outcome, run.Err)
	}
	if elapsed > 90*time.Second {
		t.Errorf("took %s; the deadline did not reach what the script was waiting on", elapsed)
	}
	if !strings.Contains(run.Detail, "deadline") {
		t.Errorf("detail does not say it was the deadline: %q", run.Detail)
	}
	t.Logf("failed after %s in phase %q", elapsed, run.FailedPhase)
}

// ---------------------------------------------------------------- attribution

// A `failed` sample says which phase died and what the compiler said. It did
// not say who wrote the line, and the ownership marker has said all along -
// reading five identical-looking `build` deaths as three distinct defects took
// a hand trace across every failed sample's invocations (prov-2026-27c10ac4).
//
// The lookup itself is behave.sh's, because attribution has to be computed
// while the generated project still exists. It is exercised here through the
// script's --attribute form, against a tree written by this test rather than by
// a run, so the marker cases are chosen rather than whatever a build happened
// to produce.

// attributeFixture writes a project and a phase log, and returns what the
// harness attributes.
func attributeFixture(t *testing.T, files map[string]string, log string) []Attribution {
	t.Helper()

	dir := t.TempDir()
	app := filepath.Join(dir, "app")
	for rel, body := range files {
		full := filepath.Join(app, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	logPath := filepath.Join(dir, "build.log")
	if err := os.WriteFile(logPath, []byte(log), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command("bash", harnessScript, "--attribute", logPath, app).CombinedOutput()
	if err != nil {
		t.Fatalf("attribute failed: %v\n%s", err, out)
	}
	var got []Attribution
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("attribution is not JSON: %v\n%s", err, out)
	}
	return got
}

const goRegions = `package db

type Grocery struct {
	ID       int
	Quantity int
}

func helper() error {
	return nil
}

// sedum:createQuery:insert {"tier":"owned","record":"PR-014","kwargs":{"table":"groceries"}}
func Insert(in Grocery) error {
	return use(in.quantity)
}
// /sedum:createQuery:insert
`

// The line the compiler named is looked up in the file it named, and the marker
// pair enclosing it says which action, which variant and which kwargs produced
// it. Nothing parses Go: the lookup is text over markers, which is what lets a
// target for a new stack inherit attribution without writing any.
func TestAFailedBuildNamesTheActionThatWroteTheLine(t *testing.T) {
	got := attributeFixture(t,
		map[string]string{"db/groceries.go": goRegions},
		"# example.com/app/db\ndb/groceries.go:14:14: in.quantity undefined "+
			"(type Grocery has no field or method quantity, but does have field Quantity)\n")

	if len(got) != 1 {
		t.Fatalf("attributed %d lines, want 1: %+v", len(got), got)
	}
	a := got[0]
	if a.Action != "createQuery" || a.Variant != "insert" {
		t.Errorf("attributed to %q, want createQuery:insert", a.Label())
	}
	if a.File != "db/groceries.go" || a.Line != 14 {
		t.Errorf("attributed %s:%d, want db/groceries.go:14", a.File, a.Line)
	}
	if a.Record != "PR-014" {
		t.Errorf("record %q, want PR-014", a.Record)
	}
	if a.Kwargs["table"] != "groceries" {
		t.Errorf("kwargs %v do not carry what the region was rendered from", a.Kwargs)
	}
}

// The compiler often names the file template's own text - the first of two
// declarations is the template's line, not the injected one - so a line no
// marker encloses is reported as unattributed rather than charged to whichever
// region happens to be nearest.
func TestALineNoMarkerEnclosesIsUnattributed(t *testing.T) {
	got := attributeFixture(t,
		map[string]string{"db/groceries.go": goRegions},
		"db/groceries.go:9:6: helper redeclared in this block\n")

	if len(got) != 1 {
		t.Fatalf("attributed %d lines, want 1: %+v", len(got), got)
	}
	if got[0].Attributed() {
		t.Errorf("a line outside every region was charged to %q", got[0].Label())
	}
}

// A build naming several files and lines is several findings. Collapsing to the
// first would report three distinct defects as one.
func TestSeveralNamedLinesYieldSeveralAttributions(t *testing.T) {
	got := attributeFixture(t, map[string]string{
		"db/groceries.go": goRegions,
		"db/todos.go": `package db

// sedum:addStructField {"tier":"owned","record":"PR-002","kwargs":{"field":"title"}}
type Todo struct {
	ID int
	ID string
}
// /sedum:addStructField
`,
	}, "db/groceries.go:14:14: in.quantity undefined\n"+
		"db/todos.go:6:2: other declaration of ID\n"+
		"db/groceries.go:14:14: in.quantity undefined\n")

	if len(got) != 2 {
		t.Fatalf("attributed %d lines, want 2 (the repeat is one finding): %+v", len(got), got)
	}
	labels := map[string]bool{}
	for _, a := range got {
		labels[a.Label()] = true
	}
	for _, want := range []string{"createQuery:insert", "addStructField"} {
		if !labels[want] {
			t.Errorf("no attribution to %s: %+v", want, got)
		}
	}
}

// The comment prefix is the package's and is never hardcoded, so the lookup
// matches the keyword rather than the prefix. Anchor declarations share the
// "sedum:" namespace and are not ownership markers.
func TestAttributionIsIndifferentToTheCommentPrefixAndSkipsAnchors(t *testing.T) {
	got := attributeFixture(t, map[string]string{
		"app/models/todo.rb": `class Todo
  # sedum:anchor:class_body
  # sedum:addValidation {"tier":"owned","record":"PR-002","kwargs":{"attr":"title"}}
  validates :title, presence: true
  # /sedum:addValidation
end
`,
	}, "app/models/todo.rb:4:in 'validates'\n")

	if len(got) != 1 {
		t.Fatalf("attributed %d lines, want 1: %+v", len(got), got)
	}
	if got[0].Action != "addValidation" {
		t.Errorf("attributed to %q, want addValidation; an anchor is not an ownership marker", got[0].Action)
	}
}

// A reference the project does not hold is not a line this run wrote. A
// host:port and a version string both look like file:line, and reporting them
// as unattributed would fill the finding with noise - unattributed is reserved
// for a line the project really holds.
func TestReferencesThatAreNotFilesInTheProjectAreDropped(t *testing.T) {
	got := attributeFixture(t,
		map[string]string{"db/groceries.go": goRegions},
		"go: downloading github.com/lib/pq v1.10.9\nlistening on 127.0.0.1:8080\n"+
			"/usr/local/go/src/net/http/server.go:3210:1: something\n")

	if len(got) != 0 {
		t.Errorf("attributed something outside the project: %+v", got)
	}
}

// Counted per action across samples, and each action once per sample however
// many of its lines the compiler named. Three samples dying in one action and
// three dying in three are different findings, and only this tells them apart.
func TestAttributionsAreCountedPerActionAcrossSamples(t *testing.T) {
	oneAction := func() Sample {
		return behaviorSample(&BehaviorRun{Outcome: "failed", FailedPhase: "build",
			Attribution: []Attribution{
				{File: "db/a.go", Line: 4, Action: "createQuery", Variant: "insert"},
				{File: "db/b.go", Line: 9, Action: "createQuery", Variant: "insert"},
			}})
	}
	m := Measurement{Samples: []Sample{
		oneAction(), oneAction(), oneAction(),
		behaviorSample(&BehaviorRun{Outcome: "failed", FailedPhase: "build",
			Attribution: []Attribution{{File: "db/c.go", Line: 2, Action: "addStructField"}}}),
		behaviorSample(&BehaviorRun{Outcome: "failed", FailedPhase: "build",
			Attribution: []Attribution{{File: "db/d.go", Line: 7}}}),
	}}
	caseWithTarget(&m, "todo-chi")

	got, _ := m.Behavior()
	if got.Actions["createQuery:insert"] != 3 {
		t.Errorf("createQuery:insert counted %d, want 3 samples (two lines in one sample is one sample)",
			got.Actions["createQuery:insert"])
	}
	if got.Actions["addStructField"] != 1 {
		t.Errorf("addStructField counted %d, want 1", got.Actions["addStructField"])
	}
	if got.Unattributed != 1 {
		t.Errorf("unattributed %d, want 1", got.Unattributed)
	}
	if got.AttributedSamples != 5 {
		t.Errorf("attributed samples %d, want 5", got.AttributedSamples)
	}

	var out strings.Builder
	behavior(&out, m)
	report := out.String()
	for _, want := range []string{"attributed to", "createQuery:insert", "addStructField", "(no enclosing region)"} {
		if !strings.Contains(report, want) {
			t.Errorf("the behaviour table does not mention %q:\n%s", want, report)
		}
	}
}

// An entry drawn before attribution existed carries none, and a reader must not
// take that for a build no action was responsible for. A dash is not a zero.
func TestAnEntryWithoutAttributionSaysNothingRatherThanNone(t *testing.T) {
	m := Measurement{Samples: []Sample{
		behaviorSample(&BehaviorRun{Outcome: "failed", FailedPhase: "build",
			Detail: "db/todos.go:61:24: undefined: models"}),
	}}
	caseWithTarget(&m, "todo-chi")

	got, _ := m.Behavior()
	if got.AttributedSamples != 0 || len(got.Actions) != 0 || got.Unattributed != 0 {
		t.Errorf("an entry carrying no attribution was tallied: %+v", got)
	}

	var out strings.Builder
	behavior(&out, m)
	if strings.Contains(out.String(), "attributed to") {
		t.Errorf("the report claimed an attribution the entry does not carry:\n%s", out.String())
	}
}
