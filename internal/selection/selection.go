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
func Select(ctx context.Context, client Client, req Request, opts Options) ([]recording.Invocation, error) {
	packages := expand.Packages(req.Files)
	cat := catalog.Build(packages, catalog.Options{})

	log := opts.Log
	if log == nil {
		log = runlog.Discard()
	}

	prompt, err := Prompt(req, cat)
	if err != nil {
		return nil, err
	}
	messages := prompt

	// Two budgets, deliberately separate. A validation retry says the answer
	// was wrong; a completeness re-prompt says it may be incomplete. Exhausting
	// one must not consume the other, and they are reported differently
	// (prov-2026-6d87dc11).
	var rejected []attempt
	askedForCompleteness := false

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
		if err != nil {
			return nil, fmt.Errorf("record %s: model call failed on attempt %d: %w", req.RecordID, i+1, err)
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
					log.Info("selection leaves anchors unfilled", "record", req.RecordID,
						"attempt", i+1, "unfilled", anchorSummary(unfilled))
					messages = append(messages,
						Message{Role: RoleAssistant, Content: raw},
						Message{Role: RoleUser, Content: note(unfilled)},
					)
					continue
				}
			}
			return invocations, nil
		}

		for _, v := range violations {
			log.Info("model output rejected", "record", req.RecordID, "attempt", i+1,
				"rule", v.Rule, "detail", v.Detail)
		}
		rejected = append(rejected, attempt{number: i + 1, violations: violations})
		if len(rejected) > opts.Retries {
			return nil, exhausted(req.RecordID, opts.Retries, rejected)
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

// attempt is one rejected response, kept so that exhaustion reports every
// violation rather than only the last response's.
type attempt struct {
	number     int
	violations []Violation
}

// exhausted reports a run that never produced a valid response.
//
// Every attempt's violations are reported, not just the final one. A model
// making a different mistake each time and a model making the same mistake
// three times are different problems with different fixes - the first is a
// prompt that underdetermines the answer, the second is usually a catalog an
// author needs to look at - and only the accumulated record distinguishes them.
func exhausted(recordID string, retries int, attempts []attempt) error {
	var b strings.Builder
	fmt.Fprintf(&b, "record %s: the model's output did not validate within %d attempt(s)",
		recordID, retries+1)

	for _, a := range attempts {
		fmt.Fprintf(&b, "\n\n  attempt %d:", a.number)
		for _, v := range a.violations {
			fmt.Fprintf(&b, "\n    %s", v)
		}
	}
	return fmt.Errorf("%s", b.String())
}
