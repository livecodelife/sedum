package evals

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// History renders one case's entries as a timeline, oldest first.
//
// Everything printed is already in the entries; nothing here runs a model,
// needs an endpoint, or recomputes a measurement. The file was built to be
// evidence a record can cite (prov-2026-eb283c56), and citing it meant
// re-deriving the summary by hand every time until this existed.
func History(out io.Writer, caseID string, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintf(out, "%s: no entries\n", caseID)
		return
	}

	fmt.Fprintf(out, "%s  (%d entries)\n", caseID, len(entries))

	// A header rather than labels repeated on every row: the columns are the
	// same for every entry, and a row carrying its own labels is unreadable at
	// this width.
	fmt.Fprintf(out, "  %-10s %-9s %2s %2s %2s  %-18s %-18s %-6s %-11s %-8s %s\n",
		"date", "commit", "n", "r", "c",
		"valid", "first call", "calls", "tok/call", "wall", "gen tok/s")

	var dirty bool
	var overlapped bool
	var absent bool

	for i, e := range entries {
		valid := wilson(e.Valid, e.Valid+e.Invalid)

		commit := e.Commit
		if commit == "" {
			commit = "-"
		}
		if !e.Clean {
			commit += "*"
			dirty = true
		}

		// A mark rather than a suppressed row: the numbers are what was
		// measured, and what the mark says is that this entry and the one
		// before it do not distinguish each other (prov-2026-0baaa119).
		row := " "
		if i > 0 {
			prev := entries[i-1]
			if valid.Overlaps(wilson(prev.Valid, prev.Valid+prev.Invalid)) {
				row = "~"
				overlapped = true
			}
		}

		fmt.Fprintf(out, "%s %-10s %-9s %2d %2d %2d  %-18s %-18s %-6s %-11s %-8s %s\n",
			row, e.At.Format("2006-01-02"), commit,
			e.Samples, e.Retries, e.Concurrency,
			valid, firstCall(e), callCost(e), tokenCost(e), wallOf(e), throughputOf(e))

		if !hasCounts(e) {
			absent = true
		}
	}

	// Legends are printed only when something in the table used them, so a
	// clean history stays short.
	if dirty {
		fmt.Fprintln(out, "  * tree was dirty: not re-runnable from that commit")
	}
	if overlapped {
		fmt.Fprintln(out, "  ~ interval overlaps the previous entry: these runs do not distinguish each other")
	}
	if absent {
		fmt.Fprintln(out, "  - not recorded: the entry predates that field")
	}
}

// hasCounts reports whether an entry carries per-call accounting at all.
//
// It is asked before reading Rejected, because an entry written before
// prov-2026-0811425c has a zero there that means "not recorded" rather than "no
// answer was rejected" - and those print differently for the reason the schema
// is additive.
func hasCounts(e Entry) bool {
	for _, r := range e.Runs {
		if r.Calls > 0 {
			return true
		}
	}
	return false
}

// firstCall recovers first-call validity from the rejection counts, which is
// what makes it readable at any retry budget (prov-2026-0811425c).
func firstCall(e Entry) string {
	if !hasCounts(e) {
		return "-"
	}
	answered, first := 0, 0
	for _, r := range e.Runs {
		if r.Outcome == "failed" {
			continue
		}
		answered++
		if r.Rejected == 0 {
			first++
		}
	}
	return wilson(first, answered).String()
}

func wallOf(e Entry) string {
	if e.WallMS <= 0 {
		return "-"
	}
	return round(time.Duration(e.WallMS) * time.Millisecond).String()
}

// throughputOf is the entry's completion-token rate over its wall clock.
//
// Completion only: a prompt token is billed whether or not the server computed
// it, and this one reuses cached prefixes so thoroughly that a prompt of two
// thousand tokens costs one (prov-2026-e323b805).
//
// It is the column a concurrency sweep is read from, and it is only a
// comparison between entries whose server was run the same way - which an entry
// cannot record (prov-2026-945fa0aa).
func throughputOf(e Entry) string {
	completion := 0
	for _, r := range e.Runs {
		completion += r.CompletionTokens
	}
	if completion == 0 || e.WallMS <= 0 {
		return "-"
	}
	return fmt.Sprintf("%.1f", float64(completion)/(float64(e.WallMS)/1000))
}

func callCost(e Entry) string {
	if !hasCounts(e) {
		return "-"
	}
	calls, samples := 0, 0
	for _, r := range e.Runs {
		if r.Outcome == "failed" {
			continue
		}
		calls += r.Calls
		samples++
	}
	if samples == 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f", float64(calls)/float64(samples))
}

func tokenCost(e Entry) string {
	prompt, completion, calls := 0, 0, 0
	for _, r := range e.Runs {
		prompt += r.PromptTokens
		completion += r.CompletionTokens
		calls += r.Calls
	}
	if calls == 0 || prompt+completion == 0 {
		return "-"
	}
	return fmt.Sprintf("%.0f+%.0f", float64(prompt)/float64(calls), float64(completion)/float64(calls))
}

// Cases returns the case ids a results directory holds, sorted.
func Cases(dir string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.jsonl"))
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, strings.TrimSuffix(filepath.Base(m), ".jsonl"))
	}
	sort.Strings(out)
	return out, nil
}
