package cli

import (
	"errors"
	"fmt"
	"testing"
)

// The exit mechanism lets a command choose the process exit code and message
// through an error that Execute honours.

func TestExitErrorIsAnErrorCarryingItsMessageAsAValueAndAsAPointer(t *testing.T) {
	var err error = ExitError{Code: 2, Message: "boom"}
	if err.Error() != "boom" {
		t.Errorf("Error() = %q", err.Error())
	}
	var ptr error = &ExitError{Code: 1}
	if ptr.Error() != "" {
		t.Errorf("pointer Error() = %q", ptr.Error())
	}
}

func TestCodeOfFindsAnExitErrorDirectlyOrWrapped(t *testing.T) {
	for name, err := range map[string]error{
		"direct":  ExitError{Code: 2, Message: "m"},
		"wrapped": fmt.Errorf("outer: %w", ExitError{Code: 2, Message: "m"}),
	} {
		if code, ok := CodeOf(err); !ok || code != 2 {
			t.Errorf("%s: CodeOf = %d, %v", name, code, ok)
		}
	}
	if _, ok := CodeOf(errors.New("plain")); ok {
		t.Error("a plain error was taken for an ExitError")
	}
	if _, ok := CodeOf(nil); ok {
		t.Error("nil was taken for an ExitError")
	}
}
