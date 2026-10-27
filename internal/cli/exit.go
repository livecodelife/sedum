package cli

import "errors"

// ExitError is an error that carries the process exit code and the message that
// Execute prints for it. An empty Message prints nothing.
type ExitError struct {
	Code    int
	Message string
}

// Error returns the message.
func (e ExitError) Error() string { return e.Message }

// CodeOf reports whether err is, or wraps, an ExitError, and returns its code.
func CodeOf(err error) (int, bool) {
	var exit ExitError
	if errors.As(err, &exit) {
		return exit.Code, true
	}
	return 0, false
}
