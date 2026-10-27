package gitcheck

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// The output unit writes the SARIF or human output, maps results to exit codes
// and builds the three kinds of error.

func TestWriteHumanPrintsFindingsThenTheSkippedAndWaivedLines(t *testing.T) {
	var out bytes.Buffer
	findings := []Finding{
		{RuleID: "SCOPE001", Message: "forbids changes to a.go", File: "a.go"},
		{RuleID: "TAG001", Message: "tag required"},
	}

	if err := WriteHuman(&out, findings, 2, 4); err != nil {
		t.Fatal(err)
	}

	want := "SCOPE001 a.go: forbids changes to a.go\n" +
		"TAG001: tag required\n" +
		"skipped 4 scope paths that no generator package claims\n" +
		"waived: 2\n"
	if out.String() != want {
		t.Errorf("output =\n%s\nwant\n%s", out.String(), want)
	}
}

func TestWriteHumanAlwaysPrintsBothLinesEvenWithNothingToReport(t *testing.T) {
	var out bytes.Buffer
	if err := WriteHuman(&out, nil, 0, 0); err != nil {
		t.Fatal(err)
	}
	if out.String() != "skipped 0 scope paths that no generator package claims\nwaived: 0\n" {
		t.Errorf("output = %q", out.String())
	}
}

func TestWriteSARIFWritesTheMergedDocument(t *testing.T) {
	var out bytes.Buffer
	if err := WriteSARIF(&out, Merge(nil, nil), 5); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"skippedScopePaths": 5`) || !strings.HasSuffix(out.String(), "\n") {
		t.Errorf("output = %s", out.String())
	}
}

func TestExitCode(t *testing.T) {
	cases := []struct {
		unwaived int
		err      error
		want     int
	}{{0, nil, 0}, {1, nil, 1}, {7, nil, 1}, {0, errors.New("x"), 2}, {3, errors.New("x"), 2}}
	for _, c := range cases {
		if got := ExitCode(c.unwaived, c.err); got != c.want {
			t.Errorf("ExitCode(%d, %v) = %d, want %d", c.unwaived, c.err, got, c.want)
		}
	}
}

func TestTheErrorKindsStartWithTheirPrefix(t *testing.T) {
	cases := map[string]error{
		"provenance: a": ProvenanceError("a"),
		"git: b":        GitError("b"),
		"provenance: c": SkewError("c"),
		"usage: d":      UsageError("d"),
	}
	for want, err := range cases {
		if err.Error() != want || FormatError(err) != want {
			t.Errorf("error = %q, FormatError = %q, want %q", err.Error(), FormatError(err), want)
		}
	}
}
