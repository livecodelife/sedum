package selection

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// prov-2026-91c54941: --response-schema sends response_format on the
// OpenAI-compatible request, opt-in and default off. This file holds the wire
// contract - the actual HTTP request bytes an endpoint receives - which is
// the level this record's own constraints are stated at ("the request Sedum
// sends is byte-for-byte identical", "sent as response_format ... on the
// OpenAI-compatible request", "fails with a clear error naming
// --response-schema"). None of it is testable through the stub Client
// selection_test.go's retry-loop tests use, because a stub never produces
// bytes on a wire; these tests run the real *OpenAI client against an
// httptest server standing in for the endpoint.
//
// CompleteWithSchema does not exist yet. It is this file's one assumption
// about the shape code-author gives *OpenAI's structured-output path - a new
// method beside the existing, untouched Complete, rather than a parameter
// added to Complete itself, specifically so the off path (used by every
// caller today) cannot change shape by construction. How Select and
// internal/cli/grow.go decide when to call it (a new Client-implementing
// interface, a type switch, an Options field threaded through) is
// code-author's decision and is not exercised here; only *OpenAI's own wire
// behavior is.
//
// Not covered here, and flagged rather than guessed at: that the compiled
// schema Select would pass in reaches run.log with the record id, and
// whether the bundled goinfer-serve backend (internal/selection/local.go)
// accepts response_format or reproduces prov-2026-4bcabb2f's collapse. The
// record itself says the second one has to be checked against a real
// goinfer-serve response, not assumed - this sandbox has no goinfer-serve
// binary or model to check it against.

// captureServer stands in for an OpenAI-compatible endpoint. It records every
// request body it receives and returns the next staged response in order.
type captureServer struct {
	*httptest.Server
	bodies [][]byte
}

func newCaptureServer(t *testing.T, responses ...func(w http.ResponseWriter)) *captureServer {
	t.Helper()
	cs := &captureServer{}
	i := 0
	cs.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("reading the captured request body: %v", err)
		}
		cs.bodies = append(cs.bodies, body)
		if i >= len(responses) {
			t.Fatalf("server received %d request(s), only %d response(s) staged", i+1, len(responses))
		}
		responses[i](w)
		i++
	}))
	return cs
}

func validCompletionResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"id":"x","object":"chat.completion","created":1,"model":"m",`+
		`"choices":[{"index":0,"message":{"role":"assistant","content":"{\"invocations\":[]}"},"finish_reason":"stop"}],`+
		`"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`)
}

// rejectingResponseFormat is what an endpoint that does not support
// structured output plausibly answers with - modelled on the shape an
// OpenAI-compatible 400 takes, not on any one server's exact wording, since
// the point under test is that Sedum surfaces it rather than retrying silent.
func rejectingResponseFormat(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, `{"error":{"message":"response_format is not supported by this endpoint","type":"invalid_request_error"}}`)
}

func newTestClient(t *testing.T, baseURL string) *OpenAI {
	t.Helper()
	t.Setenv("OPENAI_BASE_URL", baseURL)
	t.Setenv("OPENAI_API_KEY", "test-key")
	client, err := NewOpenAI("test-model")
	if err != nil {
		t.Fatalf("NewOpenAI: %v", err)
	}
	return client
}

func testMessages() []Message {
	return []Message{
		{Role: RoleSystem, Content: "system prompt"},
		{Role: RoleUser, Content: "user prompt"},
	}
}

// todaysRequestBody is what Complete sends for testMessages() against
// "test-model", captured from this repo's code as it stands before
// --response-schema exists (internal/selection's Complete is untouched by
// this record). It is the byte-for-byte baseline the off path - which is
// every call this record does not opt into - must keep reproducing.
const todaysRequestBody = `{"model":"test-model","messages":[{"role":"system","content":"system prompt"},{"role":"user","content":"user prompt"}],"max_tokens":16384,"stream":false}`

// Without --response-schema, Complete is called exactly as it is today, and
// today's request is unconditionally reproduced byte for byte. This is a
// regression pin on the method the flag must never touch, not a test of the
// flag itself - the flag's own off-behavior is that nothing here changes.
func TestCompleteWithoutSchemaMatchesTodaysRequestByteForByte(t *testing.T) {
	srv := newCaptureServer(t, validCompletionResponse)
	defer srv.Close()
	client := newTestClient(t, srv.URL)

	if _, err := client.Complete(context.Background(), testMessages()); err != nil {
		t.Fatalf("Complete: %v", err)
	}

	if got := string(srv.bodies[0]); got != todaysRequestBody {
		t.Errorf("Complete's request changed shape with no schema involved:\ngot:  %s\nwant: %s", got, todaysRequestBody)
	}
	if strings.Contains(string(srv.bodies[0]), "response_format") {
		t.Error("Complete's request carries response_format; it must be absent unless a schema is explicitly requested")
	}
}

// --response-schema sends response_format: {"type": "json_schema",
// "json_schema": {name, schema}} on the request, and it is purely additive:
// everything else about the request - model, messages, max_tokens - stays
// exactly what Complete already sends.
func TestCompleteWithSchemaSendsResponseFormat(t *testing.T) {
	srv := newCaptureServer(t, validCompletionResponse)
	defer srv.Close()
	client := newTestClient(t, srv.URL)

	schema := []byte(`{"type":"object","additionalProperties":false,"required":["invocations"],"properties":{"invocations":{"type":"array","items":{"oneOf":[]}}}}`)

	if _, err := client.CompleteWithSchema(context.Background(), testMessages(), "sedum_invocations", schema); err != nil {
		t.Fatalf("CompleteWithSchema: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(srv.bodies[0], &got); err != nil {
		t.Fatalf("request body is not valid JSON: %v\n%s", err, srv.bodies[0])
	}

	rf, ok := got["response_format"].(map[string]any)
	if !ok {
		t.Fatalf("request carries no response_format object:\n%s", srv.bodies[0])
	}
	if rf["type"] != "json_schema" {
		t.Errorf(`response_format.type = %v, want "json_schema"`, rf["type"])
	}
	js, ok := rf["json_schema"].(map[string]any)
	if !ok {
		t.Fatalf("response_format carries no json_schema object: %+v", rf)
	}
	if js["name"] != "sedum_invocations" {
		t.Errorf("json_schema.name = %v, want %q", js["name"], "sedum_invocations")
	}

	var wantSchema, gotSchema any
	if err := json.Unmarshal(schema, &wantSchema); err != nil {
		t.Fatalf("test's own schema fixture does not parse: %v", err)
	}
	gotSchemaJSON, err := json.Marshal(js["schema"])
	if err != nil {
		t.Fatalf("re-encoding json_schema.schema: %v", err)
	}
	if err := json.Unmarshal(gotSchemaJSON, &gotSchema); err != nil {
		t.Fatalf("json_schema.schema is not the schema this test passed in: %v", err)
	}

	delete(got, "response_format")
	var withoutSchema map[string]any
	if err := json.Unmarshal([]byte(todaysRequestBody), &withoutSchema); err != nil {
		t.Fatalf("todaysRequestBody does not parse: %v", err)
	}
	gotWithoutRF, _ := json.Marshal(got)
	wantWithoutRF, _ := json.Marshal(withoutSchema)
	if string(gotWithoutRF) != string(wantWithoutRF) {
		t.Errorf("adding a schema changed something other than response_format:\ngot:  %s\nwant: %s", gotWithoutRF, wantWithoutRF)
	}
}

// An endpoint that rejects response_format must fail clearly, naming
// --response-schema, and must never be silently retried unconstrained - a
// silent fallback would hide exactly the information a caller who opted in
// needs, and would mean --response-schema's presence or absence in the
// eventual request could not be trusted from the flag alone.
func TestCompleteWithSchemaRejectedByEndpointFailsClearly(t *testing.T) {
	srv := newCaptureServer(t, rejectingResponseFormat)
	defer srv.Close()
	client := newTestClient(t, srv.URL)

	schema := []byte(`{"type":"object"}`)
	_, err := client.CompleteWithSchema(context.Background(), testMessages(), "sedum_invocations", schema)
	if err == nil {
		t.Fatal("expected an error when the endpoint rejects response_format")
	}
	if !strings.Contains(err.Error(), "--response-schema") {
		t.Errorf("error does not name --response-schema, so a caller has no way to tell which flag caused it: %v", err)
	}
	if !strings.Contains(err.Error(), "response_format is not supported") {
		t.Errorf("error drops the endpoint's own explanation: %v", err)
	}

	if len(srv.bodies) != 1 {
		t.Errorf("server received %d request(s), want exactly 1 - a silent unconstrained retry must never happen", len(srv.bodies))
	}
	if strings.Contains(string(srv.bodies[0]), "response_format") == false {
		t.Fatal("test bug: the one request made did not carry response_format at all")
	}
}
