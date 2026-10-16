package growlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// Version is the schema version written in every line of the grow log.
const Version = 1

// maxError bounds the error text kept in one line.
const maxError = 300

// Regions counts what one grow did to a record's regions.
type Regions struct {
	Injected  int `json:"injected"`
	Replaced  int `json:"replaced"`
	Pruned    int `json:"pruned"`
	Unchanged int `json:"unchanged"`
}

// Entry is one line of the grow log: what one grow did for one record. It
// holds counts and timings only, never prompt or response text, kwargs, file
// contents or keys.
type Entry struct {
	V                int     `json:"v"`
	TS               string  `json:"ts"`
	Sedum            string  `json:"sedum"`
	Record           string  `json:"record"`
	Mode             string  `json:"mode"`
	Model            string  `json:"model"`
	Calls            int     `json:"calls"`
	Rejected         int     `json:"rejected"`
	Completeness     int     `json:"completeness"`
	Repairs          int     `json:"repairs"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	Invocations      int     `json:"invocations"`
	Literal          int     `json:"literal"`
	Declarative      int     `json:"declarative"`
	Regions          Regions `json:"regions"`
	Files            int     `json:"files"`
	MS               int64   `json:"ms"`
	Outcome          string  `json:"outcome"`
	Error            string  `json:"error,omitempty"`
}

// clipError keeps the first line of s, trimmed, at most maxError bytes and
// never cutting a UTF-8 sequence.
func clipError(s string) string {
	s, _, _ = strings.Cut(s, "\n")
	s = strings.TrimSpace(s)
	if len(s) <= maxError {
		return s
	}
	s = s[:maxError]
	for !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

// Append adds e as one line to the log at path, creating the file and its
// directory. V is set to Version when zero and Error is clipped. The line is
// written with a single write to an O_APPEND file, so concurrent writers never
// interleave their lines.
func Append(path string, e Entry) error {
	if e.V == 0 {
		e.V = Version
	}
	e.Error = clipError(e.Error)
	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("encode grow log entry: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create grow log directory for %s: %w", path, err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open grow log %s: %w", path, err)
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		f.Close()
		return fmt.Errorf("write grow log %s: %w", path, err)
	}
	return f.Close()
}
