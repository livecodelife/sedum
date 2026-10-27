package gitcheck

import (
	"errors"
	"fmt"
	"io"
)

// ProvenanceError is an error about provenance that cannot run or gave bad
// output.
func ProvenanceError(message string) error {
	return errors.New("provenance: " + message)
}

// GitError is an error about git that cannot read a commit or a commit message.
func GitError(message string) error {
	return errors.New("git: " + message)
}

// SkewError is an error about provenance output that lacks a field the command
// needs.
func SkewError(message string) error {
	return errors.New("provenance: " + message)
}

// UsageError is an error about how the command was invoked.
func UsageError(message string) error {
	return errors.New("usage: " + message)
}

// FormatError returns the message written to stderr for err.
func FormatError(err error) string {
	return err.Error()
}

// ExitCode maps the number of unwaived findings and an error to the exit code:
// 2 for any error, otherwise 1 when a finding is left, otherwise 0.
func ExitCode(unwaived int, err error) int {
	switch {
	case err != nil:
		return 2
	case unwaived > 0:
		return 1
	}
	return 0
}

// WriteHuman writes each remaining finding on one line, then the count of
// skipped scope paths, then the count of waived findings.
func WriteHuman(w io.Writer, findings []Finding, waived, skipped int) error {
	for _, f := range findings {
		var err error
		if f.File != "" {
			_, err = fmt.Fprintf(w, "%s %s: %s\n", f.RuleID, f.File, f.Message)
		} else {
			_, err = fmt.Fprintf(w, "%s: %s\n", f.RuleID, f.Message)
		}
		if err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, "skipped %d scope paths that no generator package claims\n", skipped); err != nil {
		return err
	}
	_, err := fmt.Fprintf(w, "waived: %d\n", waived)
	return err
}

// WriteSARIF writes the merged SARIF document followed by a newline.
func WriteSARIF(w io.Writer, merged Merged, skipped int) error {
	data, err := merged.JSON(skipped)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	_, err = io.WriteString(w, "\n")
	return err
}
