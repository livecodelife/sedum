package selection

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// A local model client.
//
// goinfer-serve is spawned as a subprocess Sedum owns for exactly this
// client's lifetime - on a scratch loopback port, closed with Close - and
// then addressed with the same OpenAI-compatible wire format NewOpenAI
// already speaks, because that is the shape goinfer-serve answers in.
//
// goinfer over llama-server: llama-server's own chat template renders this
// kind of fine-tuned model's full extended-thinking block by default, and
// suppressing it needs a request field go-openai's client has no slot for -
// risking a completion that is all reasoning and no content once the model
// outruns the token budget. goinfer's template never produces that failure,
// at the cost of exposing no thinking toggle at all: its answers carry a
// leading <think></think> whether or not there was anything to think about,
// which Complete strips below rather than passing on to the retry loop.

// LocalConfig configures a goinfer-serve-backed Client.
type LocalConfig struct {
	// ModelPath is the gguf file for goinfer-serve to load. Required.
	ModelPath string

	// ServerPath is the goinfer-serve binary. Empty resolves it from PATH -
	// this record neither downloads nor bundles it (prov-2026-a1ccdd65).
	ServerPath string

	// Backend is passed through as goinfer-serve's own -backend flag (cpu,
	// metal, cuda, webgpu). Empty leaves it at goinfer's own default.
	Backend string
}

// Local is a Client backed by a goinfer-serve subprocess.
type Local struct {
	cmd     *exec.Cmd
	client  *openai.Client
	modelID string
}

// How long NewLocal waits for goinfer-serve to finish loading the model and
// answer its own /v1/models before giving up. Generous relative to every
// load time actually measured (under 2s for a 2B model): a user's first run
// against a much larger checkpoint should get a real error about the model,
// not one that reads as "this never works."
const readyTimeout = 60 * time.Second

const readyPollInterval = 100 * time.Millisecond

// NewLocal starts goinfer-serve against cfg.ModelPath on a scratch loopback
// port and returns a Client that talks to it. The subprocess belongs to the
// returned Local; call Close when the run is done with it.
func NewLocal(ctx context.Context, cfg LocalConfig) (*Local, error) {
	if cfg.ModelPath == "" {
		return nil, errors.New("--local-model is required to start a local model server")
	}

	serverPath := cfg.ServerPath
	if serverPath == "" {
		resolved, err := exec.LookPath("goinfer-serve")
		if err != nil {
			return nil, fmt.Errorf("goinfer-serve not found on PATH; set --local-model-server to its location: %w", err)
		}
		serverPath = resolved
	}

	port, err := freePort()
	if err != nil {
		return nil, fmt.Errorf("could not find a free port for the local model server: %w", err)
	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	args := []string{"-model", cfg.ModelPath, "-addr", addr}
	if cfg.Backend != "" {
		args = append(args, "-backend", cfg.Backend)
	}
	cmd := exec.Command(serverPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not start %s: %w", serverPath, err)
	}

	base := "http://" + addr
	modelID, err := waitReady(ctx, base, cmd, &stderr)
	if err != nil {
		_ = cmd.Process.Kill()
		return nil, err
	}

	clientCfg := openai.DefaultConfig("local")
	clientCfg.BaseURL = base + "/v1"

	return &Local{
		cmd:     cmd,
		client:  openai.NewClientWithConfig(clientCfg),
		modelID: modelID,
	}, nil
}

// Complete sends the conversation to the local server and strips its
// leading <think>...</think> before returning the content, for the reason
// this file's own top comment gives.
func (l *Local) Complete(ctx context.Context, messages []Message) (Completion, error) {
	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       l.modelID,
		Temperature: 0,
		MaxTokens:   maxCompletionTokens(),
		Messages:    wire(messages),
	})
	if err != nil {
		return Completion{}, fmt.Errorf("local model %s: %w", l.modelID, err)
	}
	if len(resp.Choices) == 0 {
		return Completion{}, fmt.Errorf("local model %s returned no choices", l.modelID)
	}
	return Completion{
		Content:          stripThinking(resp.Choices[0].Message.Content),
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
	}, nil
}

// Close stops the goinfer-serve subprocess this client started.
//
// It only signals the process; it does not wait for the exit to be reaped.
// Sedum's own process is short-lived, so an un-waited child lives no longer
// than the CLI invocation that started it, and Wait is already spoken for -
// waitReady's own goroutine owns the one legal call to it.
func (l *Local) Close() error {
	if l.cmd == nil || l.cmd.Process == nil {
		return nil
	}
	return l.cmd.Process.Kill()
}

// waitReady polls base until goinfer-serve answers /v1/models, and returns
// the id of the model it loaded - discovering readiness and the served
// model's id in the one request, since a server that answers this route at
// all has finished loading.
//
// It also watches cmd for an early exit, so a server that fails fast (a bad
// flag, a model it cannot load) is reported with its own stderr rather than
// as an opaque timeout once readyTimeout elapses.
func waitReady(ctx context.Context, base string, cmd *exec.Cmd, stderr *bytes.Buffer) (string, error) {
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()

	deadline := time.Now().Add(readyTimeout)
	ticker := time.NewTicker(readyPollInterval)
	defer ticker.Stop()

	for {
		select {
		case err := <-exited:
			return "", fmt.Errorf("goinfer-serve exited before it became ready: %w\n%s", err, stderr.String())
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
			if id, err := discoverModelID(ctx, base); err == nil {
				return id, nil
			}
			if time.Now().After(deadline) {
				return "", fmt.Errorf("goinfer-serve did not become ready within %s\n%s", readyTimeout, stderr.String())
			}
		}
	}
}

// discoverModelID asks a running server what it loaded.
func discoverModelID(ctx context.Context, base string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/v1/models", nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var parsed struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	if len(parsed.Data) == 0 {
		return "", errors.New("server answered /v1/models with no model loaded")
	}
	return parsed.Data[0].ID, nil
}

// freePort finds an unused loopback port by opening and immediately closing
// a listener on port 0. The small window before goinfer-serve binds it is
// the same one every tool taking this approach (including Go's own
// httptest) accepts.
func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// thinkBlock matches a leading, possibly-empty <think>...</think> - the
// preamble goinfer-serve always emits for this model, per this file's own
// top comment.
var thinkBlock = regexp.MustCompile(`(?s)^\s*<think>.*?</think>\s*`)

func stripThinking(content string) string {
	return thinkBlock.ReplaceAllString(content, "")
}
