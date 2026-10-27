package cli

import (
	"errors"
	"testing"
)

// Execute runs the root command and returns the process exit code, honouring
// the code of an ExitError.

func TestExitCodeForReturnsTheCodeAndTheLineToPrint(t *testing.T) {
	cases := []struct {
		name string
		err  error
		code int
		line string
	}{
		{"nil", nil, 0, ""},
		{"an exit error with a message", ExitError{Code: 2, Message: "provenance: gone"}, 2, "sedum: provenance: gone"},
		{"an exit error with no message", ExitError{Code: 1}, 1, ""},
		{"a plain error", errors.New("bad flag"), 1, "sedum: bad flag"},
	}
	for _, c := range cases {
		if code, line := exitCodeFor(c.err); code != c.code || line != c.line {
			t.Errorf("%s: code %d line %q, want %d %q", c.name, code, line, c.code, c.line)
		}
	}
}
