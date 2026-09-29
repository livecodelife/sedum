package selection

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

// These test the state machine and parsing in local.go directly - the part
// worth testing without a process - the same way client.go's retry loop is
// tested without a network call. Spawning a real goinfer-serve against a
// real gguf is proof for the record's completion, not a unit test: nothing
// here depends on goinfer-serve or a model existing, only on a process and
// an HTTP server existing.

func TestStripThinking(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty think block", "<think>\n\n</think>\n\nOK, I can hear you.", "OK, I can hear you."},
		{"non-empty think block", "<think>reasoning here</think>\nOK.", "OK."},
		{"no think block", "OK, I can hear you.", "OK, I can hear you."},
		{"leading whitespace before think", "  \n<think></think>OK.", "OK."},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := stripThinking(c.in); got != c.want {
				t.Errorf("stripThinking(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func modelsHandler(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]string{{"id": id}},
		})
	}
}

func TestDiscoverModelID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(modelsHandler("Qwen3.5-2B.Q8_0"))
		defer srv.Close()

		id, err := discoverModelID(context.Background(), srv.URL)
		if err != nil {
			t.Fatalf("discoverModelID: %v", err)
		}
		if id != "Qwen3.5-2B.Q8_0" {
			t.Errorf("id = %q, want %q", id, "Qwen3.5-2B.Q8_0")
		}
	})

	t.Run("no model loaded", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]string{}})
		}))
		defer srv.Close()

		if _, err := discoverModelID(context.Background(), srv.URL); err == nil {
			t.Fatal("expected an error for an empty model list")
		}
	})

	t.Run("non-200", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer srv.Close()

		if _, err := discoverModelID(context.Background(), srv.URL); err == nil {
			t.Fatal("expected an error for a non-200 response")
		}
	})
}

func TestFreePort(t *testing.T) {
	port, err := freePort()
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	if port <= 0 || port > 65535 {
		t.Errorf("freePort() = %d, not a valid port", port)
	}
}

func TestWaitReady(t *testing.T) {
	t.Run("process exits before answering", func(t *testing.T) {
		cmd := exec.Command("sh", "-c", "echo boom >&2; exit 1")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Start(); err != nil {
			t.Fatalf("Start: %v", err)
		}

		_, err := waitReady(context.Background(), "http://127.0.0.1:1", cmd, &stderr)
		if err == nil {
			t.Fatal("expected an error when the process exits before answering")
		}
		if !strings.Contains(err.Error(), "boom") {
			t.Errorf("error %q does not carry the process's own stderr", err)
		}
	})

	t.Run("context cancelled", func(t *testing.T) {
		cmd := exec.Command("sleep", "5")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Start(); err != nil {
			t.Fatalf("Start: %v", err)
		}
		defer func() { _ = cmd.Process.Kill() }()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if _, err := waitReady(ctx, "http://127.0.0.1:1", cmd, &stderr); err == nil {
			t.Fatal("expected an error when ctx is already cancelled")
		}
	})

	t.Run("becomes ready", func(t *testing.T) {
		srv := httptest.NewServer(modelsHandler("test-model"))
		defer srv.Close()

		cmd := exec.Command("sleep", "5")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Start(); err != nil {
			t.Fatalf("Start: %v", err)
		}
		defer func() { _ = cmd.Process.Kill() }()

		id, err := waitReady(context.Background(), srv.URL, cmd, &stderr)
		if err != nil {
			t.Fatalf("waitReady: %v", err)
		}
		if id != "test-model" {
			t.Errorf("id = %q, want %q", id, "test-model")
		}
	})
}
