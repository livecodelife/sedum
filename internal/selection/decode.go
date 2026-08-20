package selection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/livecodelife/sedum/internal/recording"
)

// Decoding the model's response.
//
// The wire shape is an object wrapping an array - {"invocations": [...]} -
// because an OpenAI-compatible server constrains output with a JSON mode that
// requires the document's root to be an object. A bare array is accepted too,
// since that is the shape the PRD describes and models frequently emit
// (prov-2026-abd43bb4).
//
// That is the whole of the tolerance. Unwrapping a code fence and accepting two
// envelope shapes are formatting concessions; nothing here repairs a malformed
// invocation, infers a missing field, or drops a key it does not recognise.
// Once decoded, an invocation from a bare array faces exactly the checks one
// from an envelope faces.

// envelope is the requested response shape. Unknown fields are rejected rather
// than ignored: a response carrying commentary beside the invocations has
// understood the task differently than the catalog describes it, and one cheap
// call says so.
type envelope struct {
	Invocations []json.RawMessage `json:"invocations"`
}

// wireInvocation is one selection as it arrives.
type wireInvocation struct {
	Action string         `json:"action"`
	Kwargs map[string]any `json:"kwargs"`
}

// decode turns a raw response into invocations, or into the violations that
// explain why it could not.
//
// A decode failure is one violation about the response as a whole, because
// there is nothing to attribute it to. A response that decodes but holds a
// malformed entry produces one violation per entry, so that a response with
// three bad entries re-prompts with three.
func decode(raw string) ([]recording.Invocation, []Violation) {
	body := unfence(raw)
	if body == "" {
		return nil, []Violation{{
			Rule:   "empty_response",
			Detail: "the response was empty; return a JSON object of the form {\"invocations\": [{\"action\": ..., \"kwargs\": {...}}]}",
		}}
	}

	entries, v := entriesOf(body)
	if v != nil {
		return nil, []Violation{*v}
	}

	var (
		out        []recording.Invocation
		violations []Violation
	)
	for i, entry := range entries {
		var wire wireInvocation
		if err := readJSON(entry, &wire); err != nil {
			violations = append(violations, Violation{
				Index: i + 1,
				Rule:  "invocation_shape",
				Detail: fmt.Sprintf(
					"invocation %d is not readable as {\"action\": ..., \"kwargs\": {...}}: %s; it was %s",
					i+1, err, truncate(string(entry))),
			})
			continue
		}
		if wire.Action == "" {
			violations = append(violations, Violation{
				Index:  i + 1,
				Rule:   "missing_action",
				Detail: fmt.Sprintf("invocation %d names no action", i+1),
			})
			continue
		}
		if wire.Kwargs == nil {
			wire.Kwargs = map[string]any{}
		}
		out = append(out, recording.Invocation{Action: wire.Action, Kwargs: wire.Kwargs})
	}

	return out, violations
}

// entriesOf finds the invocation array in either accepted envelope.
func entriesOf(body string) ([]json.RawMessage, *Violation) {
	if strings.HasPrefix(body, "[") {
		var entries []json.RawMessage
		if err := json.Unmarshal([]byte(body), &entries); err != nil {
			return nil, &Violation{
				Rule:   "response_shape",
				Detail: fmt.Sprintf("the response opens as a JSON array but does not parse as one: %s", err),
			}
		}
		return entries, nil
	}

	var env envelope
	if err := readJSON([]byte(body), &env); err != nil {
		return nil, &Violation{
			Rule: "response_shape",
			Detail: fmt.Sprintf(
				"the response is not a JSON object of the form {\"invocations\": [...]}: %s; it was %s",
				err, truncate(body)),
		}
	}
	// An envelope with no invocations key at all decodes to nil, which is a
	// different statement from an empty array. The first is a response that
	// did not answer; the second is a considered "nothing to do", and a record
	// whose intent needs no action is legitimate.
	if env.Invocations == nil {
		return nil, &Violation{
			Rule: "response_shape",
			Detail: fmt.Sprintf(
				"the response carries no \"invocations\" key; return {\"invocations\": [...]}, using an empty array if no action applies. It was %s",
				truncate(body)),
		}
	}
	return env.Invocations, nil
}

// readJSON decodes one complete JSON value, ignoring keys the shape does not
// name.
//
// It used to reject them. A response is an answer to a question rather than a
// schema, and an entry naming the right action with the right arguments was
// being discarded whole for carrying anything beside them - so a model that
// annotated its own choice scored nothing. Five of fifty samples on the 4B
// described row never reached selection for it, while the other forty-five
// produced 900 of 900 assertions (prov-2026-986ac4ca).
//
// Leniency does not hide a mistyped key. An entry sending "kwarg" rather than
// "kwargs" decodes to no arguments and is rejected under missing_kwarg, naming
// the arguments it lacks, which is more actionable than being told the entry was
// unreadable. What is ignored is never read: a key the format does not define is
// not a channel.
func readJSON(data []byte, into any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(into); err != nil {
		return err
	}
	// Trailing content after a complete value is a malformed response, not a
	// second one.
	if dec.More() {
		return fmt.Errorf("trailing content after the JSON value")
	}
	return nil
}

// unfence strips a surrounding markdown code fence.
//
// A fence is a formatting habit rather than a claim about content, and models
// asked for JSON emit one constantly. Stripping it is the only rewriting done
// to a response; everything inside is decoded as it arrived.
func unfence(raw string) string {
	body := strings.TrimSpace(raw)
	if !strings.HasPrefix(body, "```") {
		return body
	}

	body = strings.TrimPrefix(body, "```")
	// An opening fence may carry a language tag, which is the rest of that
	// first line and is not part of the document.
	if newline := strings.IndexByte(body, '\n'); newline >= 0 {
		if tag := strings.TrimSpace(body[:newline]); !strings.ContainsAny(tag, "{[") {
			body = body[newline+1:]
		}
	}
	if end := strings.LastIndex(body, "```"); end >= 0 {
		body = body[:end]
	}
	return strings.TrimSpace(body)
}

// truncate keeps a quoted response short enough to read in a diagnostic while
// leaving enough of it to recognise.
func truncate(s string) string {
	const limit = 300
	if len(s) <= limit {
		return strconv.Quote(s)
	}
	return strconv.Quote(s[:limit]) + " (truncated)"
}
