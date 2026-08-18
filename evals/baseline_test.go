package evals

import (
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/record"
)

func todoRecord() *record.Record {
	return &record.Record{
		ID:          "PR-014",
		Intent:      "A todo list the service can create, read, update and delete.",
		Constraints: []string{"Only title and completed are accepted from a body."},
		Paths:       []string{"app/models/todo.rb", "config/routes.rb"},
	}
}

// The record and only the record. A baseline handed the framework and its
// version would answer whether a catalog beats a good prompt rather than
// whether Sedum beats not using Sedum (prov-2026-a4dbe65c).
func TestTheBaselinePromptCarriesTheRecordAndNothingElse(t *testing.T) {
	msgs := baselinePrompt(todoRecord())
	if len(msgs) != 2 || msgs[0].Role != "system" || msgs[1].Role != "user" {
		t.Fatalf("prompt is %d message(s), want a system and a user turn", len(msgs))
	}

	user := msgs[1].Content
	for _, want := range []string{
		"A todo list the service can create",
		"Only title and completed are accepted",
		"app/models/todo.rb",
		"config/routes.rb",
	} {
		if !strings.Contains(user, want) {
			t.Errorf("the prompt does not carry %q:\n%s", want, user)
		}
	}

	// The whole of the arm's meaning is in what is absent. A stack named
	// anywhere in either turn is a generator package written as prose.
	whole := msgs[0].Content + "\n" + user
	for _, leak := range []string{"Rails", "rails", "Ruby", "ActiveRecord", "PostgreSQL", "Gemfile", "7.2"} {
		if strings.Contains(whole, leak) {
			t.Errorf("the baseline prompt names %q; the package is the variable under test:\n%s", leak, whole)
		}
	}
}

// A fence per file, its info string the path. JSON would ask the model to
// escape whole source files, and a baseline failing on an unescaped heredoc
// would be measuring escaping.
func TestFencedFilesAreReadBackByPath(t *testing.T) {
	raw := "Here you go:\n\n" +
		"```app/models/todo.rb\n" +
		"class Todo < ApplicationRecord\nend\n" +
		"```\n\n" +
		"```config/routes.rb\n" +
		"Rails.application.routes.draw do\n  resources :todos\nend\n" +
		"```\n"

	files, unexpected, err := parseFencedFiles(raw, []string{"app/models/todo.rb", "config/routes.rb"})
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}
	if len(unexpected) != 0 {
		t.Errorf("unexpected = %v, want none", unexpected)
	}
	if got := files["app/models/todo.rb"]; got != "class Todo < ApplicationRecord\nend\n" {
		t.Errorf("model file is %q", got)
	}
	if !strings.Contains(files["config/routes.rb"], "resources :todos") {
		t.Errorf("routes file is %q", files["config/routes.rb"])
	}
}

// A model that invents a path is not granted one: the sedum arm cannot either,
// because Phase 3 creates only what affected_scope names. It is reported rather
// than dropped, since code at the wrong path is a different finding from none.
func TestAFileTheRecordDidNotAuthorizeIsReportedNotWritten(t *testing.T) {
	raw := "```app/models/todo.rb\nclass Todo\nend\n```\n" +
		"```config/database.yml\nadapter: postgresql\n```\n"

	files, unexpected, err := parseFencedFiles(raw, []string{"app/models/todo.rb"})
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}
	if _, ok := files["config/database.yml"]; ok {
		t.Error("an unauthorized path was written")
	}
	if len(unexpected) != 1 || unexpected[0] != "config/database.yml" {
		t.Errorf("unexpected = %v, want the unauthorized path named", unexpected)
	}
}

// An answer that wrote some of the files is a measurement - it will fail to
// build, which is the outcome - so the omission is reported and not fatal.
func TestAPartialAnswerNamesWhatItLeftOut(t *testing.T) {
	raw := "```app/models/todo.rb\nclass Todo\nend\n```\n"
	allowed := []string{"app/models/todo.rb", "config/routes.rb", "db/migrate/1_create.rb"}

	files, _, err := parseFencedFiles(raw, allowed)
	if err != nil {
		t.Fatalf("a partial answer was refused: %v", err)
	}
	got := missing(files, allowed)
	if len(got) != 2 || got[0] != "config/routes.rb" || got[1] != "db/migrate/1_create.rb" {
		t.Errorf("missing = %v, want the two unwritten paths in the record's order", got)
	}
}

// An answer with no file for any authorized path is an error rather than an
// empty success, the same way a selection that parsed to nothing is.
func TestAnAnswerWithNoAuthorizedFileIsAnError(t *testing.T) {
	if _, _, err := parseFencedFiles("I cannot help with that.", []string{"app/models/todo.rb"}); err == nil {
		t.Fatal("an answer carrying no file was accepted")
	}
}

// A fence carrying a language tag rather than a path is not a file. Models
// reach for ```ruby unprompted, and reading that as a path named "ruby" would
// invent a file nobody authorized.
func TestALanguageTaggedFenceIsNotAPath(t *testing.T) {
	raw := "```ruby\nclass Todo\nend\n```\n"
	_, unexpected, err := parseFencedFiles(raw, []string{"app/models/todo.rb"})
	if err == nil {
		t.Fatal("a language-tagged fence was accepted as a file")
	}
	if len(unexpected) != 1 || unexpected[0] != "ruby" {
		t.Errorf("unexpected = %v, want the tag reported so the shape is visible", unexpected)
	}
}

// A baseline report prints no number the arm cannot have. Selection, binding,
// anchor fill and idempotency all need a catalog or an invocation list, and
// zeroes for them would report a measurement nobody made.
func TestABaselineReportPrintsNothingItCannotHave(t *testing.T) {
	m := Measurement{
		Case: Case{ID: "todo-rails-baseline", Arm: "baseline",
			Expect: Expectations{Behavior: &BehaviorExpectation{Target: "todo-rails"}}},
		Model: Model{ID: "qwen", Engine: "llama.cpp", Quant: "q4_k_m"},
		Samples: []Sample{
			{
				Counts: map[string]int{"app/models/todo.rb": 1}, Total: 1, Calls: 1,
				Wrote:   []string{"app/models/todo.rb", "config/routes.rb"},
				Missing: []string{"db/migrate/1_create_todos.rb"},
				Behavior: &BehaviorRun{Outcome: "broke", FailedPhase: "build",
					Detail: "undefined method"},
			},
		},
	}

	var out strings.Builder
	Report(&out, m)
	got := out.String()

	for _, want := range []string{
		"arm=baseline",
		"undefined here and are not reported",
		"comparable to a sedum entry at retries 0",
		"authorized paths written: 2/3",
		"db/migrate/1_create_todos.rb",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the report does not contain %q:\n%s", want, got)
		}
	}
	// The selection table's own headers. Their presence would mean an arm with
	// no catalog was scored against one.
	for _, leak := range []string{"selected", "exact", "anchors filled", "idempotent"} {
		if strings.Contains(got, leak) {
			t.Errorf("the baseline report prints %q, which needs a package:\n%s", leak, got)
		}
	}
}

// The two arms are stored in one schema and told apart by the field that says
// which is which, so neither is a special case a reader has to know about.
func TestABaselineEntryCarriesItsArm(t *testing.T) {
	m := Measurement{
		Case:  Case{ID: "todo-rails-baseline", Arm: "baseline", Language: "ruby", Framework: "rails"},
		Model: Model{ID: "qwen", Engine: "llama.cpp", Quant: "q4_k_m"},
		Samples: []Sample{{
			Counts: map[string]int{"app/models/todo.rb": 1}, Total: 1, Calls: 1,
			Behavior: &BehaviorRun{Outcome: "working", Checks: 20, Passed: 20},
		}},
	}

	e := NewEntry(m, "http://localhost:1234/v1")
	if e.Arm != "baseline" {
		t.Errorf("entry arm is %q, want baseline", e.Arm)
	}
	if len(e.Runs) != 1 || e.Runs[0].Behavior == nil || e.Runs[0].Behavior.Outcome != "working" {
		t.Errorf("the behaviour outcome did not reach the entry: %+v", e.Runs)
	}
}
