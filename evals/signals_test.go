package evals

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/expand"
	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/inject"
	"github.com/calebcowen/sedum/internal/pipeline"
	"github.com/calebcowen/sedum/internal/recording"
	resolvepkg "github.com/calebcowen/sedum/internal/resolve"
)

// railsSet loads the package set the todo-rails cases run against, so the
// signals are exercised against a real package rather than a fixture written to
// agree with them.
func railsSet(t *testing.T) *genpkg.Set {
	t.Helper()
	set, findings, err := genpkg.Load("testdata/todo-rails/generators/defined", genpkg.Options{})
	if err != nil {
		t.Fatalf("loading the rails generators: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("the rails generators did not load: %s", f)
		}
	}
	return set
}

// planted builds the Phase 3 output for one file with the content it was
// created from, which is what plants the anchors an action fills.
func planted(t *testing.T, set *genpkg.Set, path, content string) resolvepkg.File {
	t.Helper()
	pkg, ok := set.Lookup("rails")
	if !ok {
		t.Fatal("the rails package did not load")
	}
	return resolvepkg.File{
		Resolution: resolvepkg.Resolution{RecordID: "PR-014", Path: path, Package: pkg},
		Rendered:   content,
	}
}

// The number comes from the package and the record, never from an expectation
// anybody wrote. That is the whole reason it sits beside the action counts:
// those can be fitted to the answer they score and this cannot
// (prov-2026-d61010a4).
func TestAnchorFillIsDerivedFromWhatTheTemplatesPlanted(t *testing.T) {
	set := railsSet(t)
	files := []resolvepkg.File{planted(t, set, "app/controllers/todos_controller.rb",
		"class TodosController\n  # sedum:anchor:actions\n\n  private\n  # sedum:anchor:private\nend\n")}

	empty := anchorFill("PR-014", files, nil, nil)
	if empty.Planted != 2 || empty.Filled != 0 {
		t.Fatalf("an empty selection filled %d of %d anchors, want 0 of 2", empty.Filled, empty.Planted)
	}
	if rate, ok := empty.Rate(); !ok || rate != 0 {
		t.Errorf("rate = %v, %v; want 0, true", rate, ok)
	}

	one := anchorFill("PR-014", files, []recording.Invocation{{
		Action: "addControllerAction",
		Kwargs: map[string]any{"resource": "todo", "verb": "index"},
	}}, nil)
	if one.Planted != 2 || one.Filled != 1 {
		t.Fatalf("one anchored action filled %d of %d, want 1 of 2", one.Filled, one.Planted)
	}
	// The count says a selection fell short; this says where. Without it the
	// stored number has to be replayed to be read at all.
	if len(one.Missed) != 1 || one.Missed[0] != "app/controllers/todos_controller.rb @private" {
		t.Errorf("missed = %v, want the unfilled anchor named", one.Missed)
	}
	if rate, _ := one.Rate(); rate != 0.5 {
		t.Errorf("rate = %v, want 0.5", rate)
	}
}

// A package that plants an anchor nothing targets has a ceiling below 100% that
// belongs to the package rather than to the answer. Counting it in the
// denominator would report that ceiling as the model falling short of one,
// which is precisely the attribution error this signal exists to prevent.
func TestAnAnchorNothingCanFillIsOutsideTheDenominator(t *testing.T) {
	set := railsSet(t)
	files := []resolvepkg.File{planted(t, set, "app/controllers/todos_controller.rb",
		"class TodosController\n  # sedum:anchor:actions\n  # sedum:anchor:extensions\nend\n")}

	fill := anchorFill("PR-014", files, []recording.Invocation{{
		Action: "addControllerAction",
		Kwargs: map[string]any{"resource": "todo", "verb": "index"},
	}}, nil)

	if fill.Planted != 1 {
		t.Fatalf("planted = %d, want 1: an anchor no action targets is not work the model omitted", fill.Planted)
	}
	if rate, ok := fill.Rate(); !ok || rate != 1 {
		t.Errorf("rate = %v, %v; want 1, true - the package's unfillable anchor depressed the model's number", rate, ok)
	}
}

// Nothing planted and nothing filled are not the same observation as a
// selection that filled none of what was there, and a bare 0.0 tells them
// apart from neither.
func TestAnchorFillIsAbsentRatherThanZeroWhenNothingWasPlanted(t *testing.T) {
	set := railsSet(t)
	files := []resolvepkg.File{planted(t, set, "app/models/todo.rb", "class Todo\nend\n")}

	fill := anchorFill("PR-014", files, nil, nil)
	if fill.Planted != 0 {
		t.Fatalf("planted = %d, want 0", fill.Planted)
	}
	if _, ok := fill.Rate(); ok {
		t.Error("a run that planted no anchors reported a fill rate")
	}
}

// wrote builds the run result and the on-disk output a syntax check reads,
// which is the pair the signal is defined over: paths Sedum created, and the
// bytes it left at them.
func wrote(t *testing.T, files map[string]string) (string, *pipeline.Result) {
	t.Helper()
	dir := t.TempDir()
	result := &pipeline.Result{}
	for path, content := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		result.Files = append(result.Files, resolvepkg.File{
			Resolution: resolvepkg.Resolution{RecordID: "PR-014", Path: path},
		})
	}
	sort.Slice(result.Files, func(i, j int) bool {
		return result.Files[i].Path < result.Files[j].Path
	})
	return dir, result
}

// The command's exit status is the whole of the verdict. The harness runs what
// the case named and reads the code; it does not parse the diagnostic, and it
// has no opinion about what the command was.
func TestTheSyntaxCheckIsTheCommandsExitStatus(t *testing.T) {
	dir, result := wrote(t, map[string]string{
		"app/models/todo.rb":  "class Todo\nend\n",
		"app/models/user.rb":  "class User\nend\n",
		"config/database.yml": "adapter: postgresql\n",
	})

	got := syntaxOf(context.Background(), Check{".rb": {"true"}}, dir, result)
	if got.Checked != 2 || got.Parsed != 2 {
		t.Fatalf("parsed %d of %d, want 2 of 2 - the .yml file has no entry and is not checked",
			got.Parsed, got.Checked)
	}
	if rate, ok := got.Rate(); !ok || rate != 1 {
		t.Errorf("rate = %v, %v; want 1, true", rate, ok)
	}

	rejected := syntaxOf(context.Background(), Check{".rb": {"false"}}, dir, result)
	if rejected.Parsed != 0 || rejected.Checked != 2 {
		t.Fatalf("parsed %d of %d, want 0 of 2", rejected.Parsed, rejected.Checked)
	}
	if len(rejected.Failures) != 2 {
		t.Errorf("failures = %v, want one per rejected file", rejected.Failures)
	}
	if !strings.HasPrefix(rejected.Failures[0], "app/models/todo.rb") {
		t.Errorf("a failure does not name its file first: %q", rejected.Failures[0])
	}
}

// A machine without the interpreter installed must not report every file as
// malformed. That number would be zero, alarming, and shaped exactly like a
// real finding about the model - the worst thing a measurement can be.
func TestAnAbsentInterpreterIsNotAFailingFile(t *testing.T) {
	dir, result := wrote(t, map[string]string{"app/models/todo.rb": "class Todo\nend\n"})

	got := syntaxOf(context.Background(), Check{".rb": {"sedum-no-such-interpreter"}}, dir, result)
	if got.Checked != 0 || got.Parsed != 0 || len(got.Failures) != 0 {
		t.Fatalf("an uninstalled command was counted: %+v", got)
	}
	if len(got.Unavailable) != 1 || got.Unavailable[0] != "sedum-no-such-interpreter" {
		t.Fatalf("unavailable = %v, want the command that could not be run", got.Unavailable)
	}
	if _, ok := got.Rate(); ok {
		t.Error("a rate was reported for files nothing checked")
	}
}

// A case that declares no command is measured on every other signal and simply
// reports no syntactic validity.
func TestACaseWithNoCheckReportsNoSyntaxRate(t *testing.T) {
	dir, result := wrote(t, map[string]string{"app/models/todo.rb": "class Todo\nend\n"})

	got := syntaxOf(context.Background(), nil, dir, result)
	if got.Checked != 0 {
		t.Fatalf("checked = %d, want 0", got.Checked)
	}
	if _, ok := got.Rate(); ok {
		t.Error("a case with no check command reported a syntax rate")
	}
}

// Against a real parser, end to end: the signal exists to catch output that
// does not parse, and the check that never fails on broken input is worthless.
func TestARealParserRejectsMalformedOutput(t *testing.T) {
	if _, err := exec.LookPath("gofmt"); err != nil {
		t.Skip("gofmt is not on PATH")
	}
	dir, result := wrote(t, map[string]string{
		"handlers/todo.go":   "package handlers\n\nfunc List() {}\n",
		"handlers/broken.go": "package handlers\n\nfunc List( {}\n",
	})

	got := syntaxOf(context.Background(), Check{".go": {"gofmt", "-e", "-l"}}, dir, result)
	if got.Checked != 2 {
		t.Fatalf("checked = %d, want 2", got.Checked)
	}
	if got.Parsed != 1 {
		t.Fatalf("parsed %d of 2, want 1 - the malformed file was accepted", got.Parsed)
	}
	if len(got.Failures) != 1 || !strings.HasPrefix(got.Failures[0], "handlers/broken.go") {
		t.Errorf("failures = %v, want the broken file named", got.Failures)
	}
}

// applied writes one file's template output to disk and injects into it once,
// which is the state the idempotency signal asks its question about: a file
// that already carries a region from an earlier application.
func applied(t *testing.T, dir string, file resolvepkg.File, invocations []recording.Invocation) *pipeline.Result {
	t.Helper()
	full := filepath.Join(dir, file.Path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(file.Rendered), 0o644); err != nil {
		t.Fatal(err)
	}
	return &pipeline.Result{
		Files: []resolvepkg.File{file},
		Selections: []pipeline.Selection{{
			RecordID:    file.RecordID,
			Files:       []resolvepkg.File{file},
			Invocations: invocations,
		}},
	}
}

// The property M4 asserts on hand-built fixtures, asserted here against a real
// package and a real action: applying the same invocations to the file they
// already wrote into changes nothing.
func TestASecondApplicationLeavesTheBytesAlone(t *testing.T) {
	set := railsSet(t)
	dir := t.TempDir()
	file := planted(t, set, "app/controllers/todos_controller.rb",
		"class TodosController\n  # sedum:anchor:actions\nend\n")
	invocations := []recording.Invocation{{
		Action: "addControllerAction",
		Kwargs: map[string]any{"resource": "todo", "verb": "index"},
	}}

	result := applied(t, dir, file, invocations)

	// The first application, which the signal treats as already having
	// happened - in a sample it is Phase 7.
	expanded, err := expand.Expand(file.RecordID, result.Files, invocations, nil)
	if err != nil {
		t.Fatalf("expanding: %v", err)
	}
	if _, err := inject.Apply(expanded, inject.Options{Output: dir}); err != nil {
		t.Fatalf("first application: %v", err)
	}

	got := idempotencyOf(dir, result)
	if got.Err != "" {
		t.Fatalf("the second application failed: %s", got.Err)
	}
	if got.Files != 1 || got.Stable != 1 {
		t.Fatalf("%d of %d files stable, want 1 of 1; differing: %v", got.Stable, got.Files, got.Differing)
	}
	if rate, ok := got.Rate(); !ok || rate != 1 {
		t.Errorf("rate = %v, %v; want 1, true", rate, ok)
	}
}

// A signal that cannot fail is not a signal. This is the shape of the defect it
// exists to catch: an application that does not converge, so the bytes after
// differ from the bytes before.
func TestAnApplicationThatChangesTheFileIsReported(t *testing.T) {
	set := railsSet(t)
	dir := t.TempDir()
	file := planted(t, set, "app/controllers/todos_controller.rb",
		"class TodosController\n  # sedum:anchor:actions\nend\n")

	// Written but never injected into, so the application the signal makes is
	// the first one and it necessarily changes the file.
	result := applied(t, dir, file, []recording.Invocation{{
		Action: "addControllerAction",
		Kwargs: map[string]any{"resource": "todo", "verb": "index"},
	}})

	got := idempotencyOf(dir, result)
	if got.Err != "" {
		t.Fatalf("applying: %s", got.Err)
	}
	if got.Stable != 0 || got.Files != 1 {
		t.Fatalf("%d of %d files stable, want 0 of 1", got.Stable, got.Files)
	}
	if len(got.Differing) != 1 || got.Differing[0] != "app/controllers/todos_controller.rb" {
		t.Errorf("differing = %v, want the file that changed", got.Differing)
	}
}

// signalSample is a sample that reached Phase 7, with the three signals set.
func signalSample(fill AnchorFill, syn SyntaxCheck, idem Idempotency) Sample {
	return Sample{
		Counts: map[string]int{"addControllerAction": 1}, Total: 1,
		Fill: fill, Syntax: syn, Idempotent: idem,
	}
}

// The three rungs count different things - anchors, files, files again - and
// each carries its own denominator. One combined number would hide which rung
// moved, which is the conflation the ladder exists to prevent.
func TestTheSignalsAreReportedWithTheirOwnDenominators(t *testing.T) {
	m := Measurement{
		Case:  Case{ID: "todo-rails-defined", Check: Check{".rb": {"ruby", "-c"}}},
		Model: Model{ID: "qwen", Engine: "mlx", Quant: "4bit"},
		Samples: []Sample{
			signalSample(AnchorFill{Planted: 4, Filled: 3},
				SyntaxCheck{Checked: 3, Parsed: 3},
				Idempotency{Files: 3, Stable: 3}),
			signalSample(AnchorFill{Planted: 4, Filled: 2, Missed: []string{"db/todos_test.go @tests", "main.go @imports"}},
				SyntaxCheck{Checked: 3, Parsed: 2, Failures: []string{"app/models/todo.rb: syntax error"}},
				Idempotency{Files: 3, Stable: 3}),
		},
	}

	var out strings.Builder
	Report(&out, m)
	got := out.String()

	for _, want := range []string{
		"signals  2 sample(s) applied",
		"anchors filled: 5/8",
		"parses: 5/6",
		"idempotent: 6/6",
		"app/models/todo.rb: syntax error",
		"anchors left unfilled",
		"db/todos_test.go @tests",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the report does not contain %q:\n%s", want, got)
		}
	}
}

// Syntactic validity is not correctness, and the line that prints it says so.
// Output that parses and raises at runtime is the failure this signal's own
// record was written from, and it passes here (prov-2026-c5697387).
func TestTheSyntaxLineDoesNotClaimCorrectness(t *testing.T) {
	m := Measurement{
		Case:    Case{ID: "todo-rails-defined", Check: Check{".rb": {"ruby", "-c"}}},
		Model:   Model{ID: "qwen", Engine: "mlx", Quant: "4bit"},
		Samples: []Sample{signalSample(AnchorFill{}, SyntaxCheck{Checked: 2, Parsed: 2}, Idempotency{})},
	}

	var out strings.Builder
	Report(&out, m)
	if !strings.Contains(out.String(), "syntax only, not correctness") {
		t.Errorf("the syntax line does not say what it is not:\n%s", out.String())
	}
}

// An uninstalled interpreter reports as not checked, never as a rate of zero.
// A zero would be alarming, wrong, and shaped exactly like a real finding.
func TestAnUninstalledInterpreterIsNotReportedAsAFailingRate(t *testing.T) {
	m := Measurement{
		Case:  Case{ID: "todo-rails-defined", Check: Check{".rb": {"ruby", "-c"}}},
		Model: Model{ID: "qwen", Engine: "mlx", Quant: "4bit"},
		Samples: []Sample{signalSample(AnchorFill{Planted: 2, Filled: 2},
			SyntaxCheck{Unavailable: []string{"ruby"}}, Idempotency{Files: 1, Stable: 1})},
	}

	var out strings.Builder
	Report(&out, m)
	got := out.String()

	if !strings.Contains(got, "parses: not checked - ruby is not installed here") {
		t.Errorf("an absent interpreter is not named as the reason:\n%s", got)
	}
	if strings.Contains(got, "parses: 0/") {
		t.Errorf("an absent interpreter was reported as a zero rate:\n%s", got)
	}
}

// A signal that lived only in a terminal is the failure the results file exists
// to prevent. Stored per sample, so a question sharpened after the run can be
// asked of the draws that produced the number.
func TestAnEntryCarriesEachSamplesSignals(t *testing.T) {
	m := Measurement{
		Case:  Case{ID: "todo-rails-defined"},
		Model: Model{ID: "qwen", Engine: "mlx", Quant: "4bit"},
		Samples: []Sample{
			signalSample(AnchorFill{Planted: 4, Filled: 3, Missed: []string{"db/todos_test.go @tests"}},
				SyntaxCheck{Checked: 3, Parsed: 2, Failures: []string{"app/models/todo.rb: syntax error"}},
				Idempotency{Files: 3, Stable: 2, Differing: []string{"config/routes.rb"}}),
			{Invalid: true, Detail: "rejected"},
		},
	}

	e := NewEntry(m, "http://localhost:1234/v1")
	if len(e.Runs) != 2 {
		t.Fatalf("runs = %d, want 2", len(e.Runs))
	}

	got := e.Runs[0]
	if got.Fill == nil || got.Fill.Filled != 3 || got.Fill.Planted != 4 {
		t.Errorf("fill = %+v, want 3 of 4", got.Fill)
	}
	if got.Fill == nil || len(got.Fill.Missed) != 1 || got.Fill.Missed[0] != "db/todos_test.go @tests" {
		t.Errorf("fill.Missed = %+v, want the unfilled anchor named", got.Fill)
	}
	if got.Syntax == nil || got.Syntax.Parsed != 2 || len(got.Syntax.Failures) != 1 {
		t.Errorf("syntax = %+v, want 2 parsed and the failing file named", got.Syntax)
	}
	if got.Idempotent == nil || got.Idempotent.Differing[0] != "config/routes.rb" {
		t.Errorf("idempotent = %+v, want the changed file named", got.Idempotent)
	}

	// A rejected answer never wrote a file. Storing its zeroes would put them
	// in a denominator they do not belong in.
	if e.Runs[1].Fill != nil || e.Runs[1].Syntax != nil || e.Runs[1].Idempotent != nil {
		t.Errorf("an invalid sample carries signals it could not have produced: %+v", e.Runs[1])
	}

	// Round trip, because the entry is read back and not only written
	// (prov-2026-c0f55691).
	encoded, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	var back Entry
	if err := json.Unmarshal(encoded, &back); err != nil {
		t.Fatal(err)
	}
	if back.Runs[0].Fill == nil || back.Runs[0].Fill.Filled != 3 {
		t.Errorf("the fill did not survive the round trip: %s", encoded)
	}
	if back.Runs[0].Fill == nil || len(back.Runs[0].Fill.Missed) != 1 {
		t.Errorf("the missed anchor did not survive the round trip: %s", encoded)
	}
	if back.Runs[1].Fill != nil {
		t.Error("an invalid sample gained a fill on the way back")
	}
}

// An entry written before the signals existed reads as not recorded, never as a
// sample that filled no anchors. The schema is additive.
func TestAnOlderEntryReportsNoSignalsRatherThanZeroes(t *testing.T) {
	const stored = `{"schema":1,"commit":"abc","clean":true,"case":"todo-rails-defined",
	  "samples":1,"valid":1,"runs":[{"outcome":"valid","calls":1}]}`

	var e Entry
	if err := json.Unmarshal([]byte(stored), &e); err != nil {
		t.Fatal(err)
	}
	if e.Runs[0].Fill != nil || e.Runs[0].Syntax != nil || e.Runs[0].Idempotent != nil {
		t.Errorf("an older entry reports signals it never carried: %+v", e.Runs[0])
	}
}
