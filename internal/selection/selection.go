// Package selection runs Phases 4 and 5: it asks a model which actions to
// invoke, and holds the answer to the catalog's declared schema.
//
// This is the only place in Sedum a model is consulted, and the only place
// anything is non-deterministic. Everything downstream consumes the validated
// invocation list and nothing else, which is what makes the rest of the run a
// pure function of that list and the generator packages.
//
// The model's job is bounded to selection: it picks actions from a closed
// vocabulary and binds their arguments. It does not write code, choose file
// paths, or decide structure. That boundary is what makes every failure mode
// here machine-checkable before anything runs, and it is why a wrong answer
// costs one model call rather than a compile, a service start, and a test run.
package selection

import (
	"context"
	"fmt"
	"strings"

	"github.com/calebcowen/sedum/internal/catalog"
	"github.com/calebcowen/sedum/internal/expand"
	"github.com/calebcowen/sedum/internal/recording"
	"github.com/calebcowen/sedum/internal/resolve"
	"github.com/calebcowen/sedum/internal/runlog"
)

// Roles in a model conversation.
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// Message is one turn. It is Sedum's own shape rather than the client
// library's, so that the retry loop and the prompt are testable without a
// server and without go-openai's types reaching the rest of the package.
type Message struct {
	Role    string
	Content string
}

// Client is the model.
//
// It is an interface for two reasons that matter equally: the retry loop is the
// part most worth testing and the part least worth a network call to test, and
// the structured-output contract must not quietly become a tool-calling one -
// a client that can only return a string cannot smuggle in tool calls.
type Client interface {
	Complete(ctx context.Context, messages []Message) (string, error)
}

// Request is one provenance record's Phase 4 input.
//
// Files is the run's Phase 3 output for this record. It decides two things: the
// packages the catalog is drawn from, and the paths an action is allowed to
// resolve to. Nothing else about the run is visible to the model.
type Request struct {
	RecordID    string
	Intent      string
	Constraints []string
	Files       []resolve.File
}

// Options bound and observe the loop.
type Options struct {
	// Retries is how many times a rejected response may be re-prompted. Zero
	// means the first response is the only one.
	Retries int

	// Log is the run log. A nil log discards.
	Log *runlog.Log
}

// Answer is what one record's selection produced, and what it cost.
//
// The counts are here rather than only in the run log because the loop is the
// only thing that knows them, and a caller reconstructing them from the outside
// has nothing to reconstruct from but a clock - which is calls multiplied by an
// unknown and varying per-call cost (prov-2026-0811425c).
type Answer struct {
	Invocations []recording.Invocation

	// Calls is every model call this record cost, including the rejected
	// answers and the completeness observation.
	Calls int

	// Rejected is how many answers Phase 5 refused.
	//
	// It is what makes first-call validity recoverable at any retry budget: a
	// record with no rejections validated on its first call, whatever the
	// budget would have allowed.
	Rejected int

	// Completeness is the observation call: 0 or 1.
	//
	// Counted apart from Rejected because it draws from its own budget
	// (prov-2026-6d87dc11) and because the answer that earned it was valid.
	// Folding the two together would report a complete answer as a rejected
	// one.
	Completeness int
}

// Rejection is an answer Phase 5 refused every time it asked.
//
// It is a type rather than a message so that a caller can tell it from a
// transport failure without matching on text. The two are different
// measurements - a model that chose badly and a server that was not there - and
// the harness classified them by matching "did not validate" until this existed.
type Rejection struct {
	RecordID string
	Retries  int

	// Attempts is every rejected answer with its violations, kept whole. A
	// model making a different mistake each time and one making the same
	// mistake three times are different problems with different fixes.
	Attempts []Attempt

	// Calls, Rejected and Completeness are Answer's counts for a record that
	// never produced one, so cost is reported for the samples that failed as
	// well as the samples that did not.
	Calls        int
	Rejected     int
	Completeness int
}

// Attempt is one rejected answer.
type Attempt struct {
	Number     int
	Violations []Violation
}

func (e *Rejection) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "record %s: the model's output did not validate within %d attempt(s)",
		e.RecordID, e.Retries+1)

	for _, a := range e.Attempts {
		fmt.Fprintf(&b, "\n\n  attempt %d:", a.Number)
		for _, v := range a.Violations {
			fmt.Fprintf(&b, "\n    %s", v)
		}
	}
	return b.String()
}

// Select runs one record's model call and validates the result.
//
// One call per record. The prompt carries the record's intent, its constraints,
// the paths created for it, and the catalog - the union of exposed actions
// across every package those paths resolved to.
//
// A rejected response is re-prompted with the specific violations appended,
// which is the whole reason Phase 5's checks are individually specific: the
// loop's value is that the model is told exactly what was wrong, and a generic
// "invalid output" would spend a call teaching it nothing.
//
// A transport failure is not re-prompted. Nothing about a connection error is
// the model's to correct, and retrying it would silently consume the budget
// reserved for the mistakes that are.
func Select(ctx context.Context, client Client, req Request, opts Options) (Answer, error) {
	packages := expand.Packages(req.Files)
	cat := catalog.Build(packages, catalog.Options{})

	log := opts.Log
	if log == nil {
		log = runlog.Discard()
	}

	prompt, err := Prompt(req, cat)
	if err != nil {
		return Answer{}, err
	}
	messages := prompt

	// Two budgets, deliberately separate. A validation retry says the answer
	// was wrong; a completeness re-prompt says it may be incomplete. Exhausting
	// one must not consume the other, and they are reported differently
	// (prov-2026-6d87dc11).
	var rejected []Attempt
	askedForCompleteness := false

	// Counted here rather than derived by the caller, because this loop is the
	// only thing that can tell a rejected answer from a completeness
	// observation after the fact.
	var answer Answer

	for i := 0; ; i++ {
		log.Info("invoking model", "record", req.RecordID, "attempt", i+1,
			"actions", len(cat.Actions), "files", len(req.Files))
		// The prompt is logged as well as the response. "Why did the model
		// pick nothing?" is the question a package author asks most, and it
		// is unanswerable from the response alone - the catalog it was shown
		// is usually where the answer is.
		log.Info("model prompt", "record", req.RecordID, "attempt", i+1,
			"prompt", messages[len(messages)-1].Content)

		raw, err := client.Complete(ctx, messages)
		answer.Calls++
		if err != nil {
			// A transport failure is not the model's mistake and is not a
			// Rejection. Its counts go with it: a call that never returned an
			// answer is not a cost attributable to a selection.
			return Answer{}, fmt.Errorf("record %s: model call failed on attempt %d: %w", req.RecordID, i+1, err)
		}
		log.Info("model responded", "record", req.RecordID, "attempt", i+1, "response", raw)

		invocations, violations := decode(raw)
		if len(violations) == 0 {
			violations = validate(cat, packages, req.Files, invocations)
		}
		if len(violations) == 0 {
			log.Info("model output validated", "record", req.RecordID,
				"attempt", i+1, "invocations", len(invocations))

			// The answer is valid either way. What remains is whether it is
			// complete, which is a question rather than a verdict: the model
			// keeps the judgment and simply stops making it blind. A response
			// leaving nothing unfilled is never re-prompted, so the case that
			// was already complete costs nothing.
			//
			// Nor is a difference the model could not have closed: an anchor
			// no action in the run targets is already excluded upstream, so a
			// package whose templates plant a region reserved for a later
			// record earns no call here at all (prov-2026-206fa618).
			if !askedForCompleteness {
				unfilled := expand.Unfilled(req.RecordID, req.Files, invocations)
				if len(unfilled) > 0 {
					askedForCompleteness = true
					answer.Completeness = 1
					log.Info("selection leaves anchors unfilled", "record", req.RecordID,
						"attempt", i+1, "unfilled", anchorSummary(unfilled))
					messages = append(messages,
						Message{Role: RoleAssistant, Content: raw},
						Message{Role: RoleUser, Content: note(unfilled)},
					)
					continue
				}
			}
			answer.Invocations = invocations
			return answer, nil
		}

		for _, v := range violations {
			log.Info("model output rejected", "record", req.RecordID, "attempt", i+1,
				"rule", v.Rule, "detail", v.Detail)
		}
		rejected = append(rejected, Attempt{Number: i + 1, Violations: violations})
		answer.Rejected++
		if len(rejected) > opts.Retries {
			return Answer{}, &Rejection{
				RecordID:     req.RecordID,
				Retries:      opts.Retries,
				Attempts:     rejected,
				Calls:        answer.Calls,
				Rejected:     answer.Rejected,
				Completeness: answer.Completeness,
			}
		}

		// The rejected response stays in the conversation. Re-prompting with
		// the violations but without what they refer to would ask the model to
		// correct something it can no longer see.
		messages = append(messages,
			Message{Role: RoleAssistant, Content: raw},
			Message{Role: RoleUser, Content: rejection(violations)},
		)
	}
}
