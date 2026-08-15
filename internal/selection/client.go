package selection

import (
	"context"
	"errors"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// The model client.
//
// One code path serves a local server and a hosted API, because the endpoint is
// a variable rather than a mode. That is the whole reason the client is
// OpenAI-compatible: the open-weight range worth evaluating is served by LM
// Studio, vLLM, and Ollama, and none of them should need a branch here.

// Environment variables the client reads. They are the names go-openai's
// ecosystem already uses, so a configuration that works for anything else works
// for Sedum (prov-2026-2f131ba6).
const (
	EnvBaseURL = "OPENAI_BASE_URL"
	EnvAPIKey  = "OPENAI_API_KEY"
)

// OpenAI is a Client backed by an OpenAI-compatible chat completions endpoint.
type OpenAI struct {
	client *openai.Client
	model  string
}

// NewOpenAI builds a client for the named model.
//
// Configuration is the model name plus the environment, and nothing else. No
// generator package and no provenance record may name a model, an endpoint, or
// a credential: both would make a committed artifact carry an operational
// dependency, and a recording replayed a year later would be asking for a model
// that no longer exists in order to run phases that never call one.
func NewOpenAI(model string) (*OpenAI, error) {
	if model == "" {
		return nil, errors.New("--model is required for a run that invokes a model; it names a model the configured endpoint serves")
	}

	baseURL := os.Getenv(EnvBaseURL)
	key := os.Getenv(EnvAPIKey)

	// A local server accepts any string as a key and most users leave it
	// unset, so an empty key is only fatal against the default hosted
	// endpoint - where it cannot succeed, and where failing now beats failing
	// after a round trip with an error the user has to map back to a variable
	// name.
	if key == "" && baseURL == "" {
		return nil, fmt.Errorf(
			"neither %s nor %s is set, so there is no endpoint to call and no credential for the default one",
			EnvBaseURL, EnvAPIKey)
	}

	cfg := openai.DefaultConfig(key)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	return &OpenAI{client: openai.NewClientWithConfig(cfg), model: model}, nil
}

// Complete sends the conversation and returns the model's message.
//
// Temperature is zero. Sampling is the one source of variation Sedum can
// actually turn down, and a run that produced different invocations from the
// same record and the same packages would make the recording's promise of
// reproducibility harder to reason about than it needs to be.
//
// Structured output, never tool calling. Most of the open-weight range worth
// evaluating has no tool support, and requiring it would exclude exactly the
// models this design exists to accommodate.
//
// No response_format is sent. Asking a server to constrain the envelope made
// the model stop selecting: grammar-constrained decoding takes the shortest
// legal completion at the array's first token, and an empty array is legal, so
// the same prompt that produced three correct invocations produced none
// (prov-2026-4bcabb2f). The contract lives in the prompt and in Phase 5, which
// is where it can be enforced specifically enough to re-prompt with.
// The usage block is carried through rather than dropped. It is the server's
// own accounting of what the call cost, and a harness that had to infer it from
// elapsed time would be reading tokens multiplied by an unknown token rate
// (prov-2026-096a4d4b). A server that fills none leaves the counts at zero,
// which callers keep distinct from a call that cost nothing.
func (o *OpenAI) Complete(ctx context.Context, messages []Message) (Completion, error) {
	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       o.model,
		Temperature: 0,
		Messages:    wire(messages),
	})
	if err != nil {
		return Completion{}, fmt.Errorf("model %s: %w", o.model, err)
	}
	if len(resp.Choices) == 0 {
		// Not a validation failure: the server returned successfully and said
		// nothing, which is a transport-level problem rather than a response
		// worth re-prompting.
		return Completion{}, fmt.Errorf("model %s returned no choices", o.model)
	}
	return Completion{
		Content:          resp.Choices[0].Message.Content,
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
	}, nil
}

func wire(messages []Message) []openai.ChatCompletionMessage {
	out := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, m := range messages {
		out = append(out, openai.ChatCompletionMessage{Role: m.Role, Content: m.Content})
	}
	return out
}
