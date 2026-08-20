package resolve

import (
	"fmt"
	"strings"

	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/inject"
)

// Reconciling an existing file with its file template (prov-2026-4c49ca46).
//
// Phase 3 used to verify an existing file and halt when it lacked the markers
// its template plants. That refused three situations it could have handled: an
// empty file, a file another tool generated, and a file whose template has
// grown since it was written. What replaces it applies the template where every
// difference between the two is something the template adds, and halts
// otherwise.
//
// Nothing here reads a language. A line is structural unless it is blank or
// begins with the package's declared comment_prefix, and a marker is structural
// because it is the thing being planted. That is the only knowledge involved,
// and the package supplies it.

// reconciled is what comparing a file against its rendered template produced.
type reconciled struct {
	// Content is what the file must contain. Equal to the input when the
	// template adds nothing.
	Content string
	// Inserted counts the lines the template contributed, for the run log. A
	// team adopting Sedum into a repository has to be able to read what it did
	// to their files rather than discover it (prov-2026-4c49ca46).
	Inserted int
}

// reconcile applies to existing whatever its rendered template adds.
//
// The whole of the rule is that the file's structural lines must be a
// subsequence of the template's. That is exactly the insertions-only
// constraint: if every structural line of the file appears in the template, in
// order, then the template contains everything the file does and the difference
// is what it adds. If one does not, the file carries content the template does
// not account for, and reconciling would mean deciding what to do with it -
// which is a judgment about someone's code and not Sedum's to make.
//
// Greedy leftmost matching is used and is complete: where a subsequence exists
// at all, taking the earliest match for each line finds one.
func reconcile(commentPrefix, path, template, existing, rendered string) (reconciled, error) {
	file, trailing := splitLines(existing)
	tmpl, _ := splitLines(rendered)

	// The skeleton is the file with every marked region removed, so that what
	// is compared is the boilerplate rather than the boilerplate plus whatever
	// has been injected into it since.
	skeleton, origin, owner, err := decompose(commentPrefix, file, existing)
	if err != nil {
		return reconciled{}, fmt.Errorf("file %s: %w", path, err)
	}

	fileAt := structural(skeleton, commentPrefix)
	tmplAt := structural(tmpl, commentPrefix)

	matched, failed := align(pick(skeleton, fileAt), pick(tmpl, tmplAt))
	if failed >= 0 {
		line := origin[fileAt[failed]]
		return reconciled{}, fmt.Errorf(
			"file %s cannot be reconciled with its file template %s: line %d, %q, has no counterpart in the template. "+
				"Sedum applies a template to a file that already exists only when every difference is something the template adds; "+
				"this file carries content the template does not account for, and what to do with it is not Sedum's to decide",
			path, template, line+1, strings.TrimSpace(skeleton[fileAt[failed]]))
	}

	// Gaps are walked back to front so that an earlier edit does not move the
	// position a later one was computed against.
	out := append([]string{}, file...)
	for i := len(matched); i >= 0; i-- {
		block := templateRun(tmpl, tmplAt, matched, i)
		if len(block) == 0 {
			// The template adds nothing here, and nothing is ever removed, so
			// whatever the file has between these two lines is already right.
			continue
		}
		from, to := fillerRun(file, skeleton, origin, owner, fileAt, i)

		// Whitespace between two lines the file and the template share is the
		// same whitespace, so the template's arrangement of it governs. Where
		// the file has put anything else there - a comment its author wrote -
		// that is theirs, and what the template adds goes after it rather than
		// over it.
		if blankOnly(file[from:to]) {
			out = splice(out, from, to, block)
		} else {
			out = splice(out, to, to, block)
		}
	}

	content := joinLines(out, trailing || existing == "")
	if content == existing {
		return reconciled{Content: existing}, nil
	}
	return reconciled{Content: content, Inserted: len(out) - len(file)}, nil
}

// templateRun is what the template has between the two lines the file already
// shares with it: the structural lines it adds, and the blanks and comments
// they arrive with.
func templateRun(tmpl []string, tmplAt, matched []int, i int) []string {
	first := 0
	if i > 0 {
		first = matched[i-1] + 1
	}
	last := len(tmplAt)
	if i < len(matched) {
		last = matched[i]
	}
	if first >= last {
		return nil
	}

	lo := 0
	if first > 0 {
		lo = tmplAt[first-1] + 1
	}
	hi := len(tmpl)
	if last < len(tmplAt) {
		hi = tmplAt[last]
	}
	return tmpl[lo:hi]
}

// fillerRun is the same span in the file: everything between the two shared
// lines that is not one of them.
//
// It starts after whatever the preceding line owns, so an injected region is
// outside the span and can neither be replaced nor written into.
func fillerRun(file, skeleton []string, origin, owner, fileAt []int, i int) (from, to int) {
	from = 0
	if i > 0 {
		from = owner[fileAt[i-1]] + 1
	}
	to = len(file)
	if i < len(fileAt) {
		to = origin[fileAt[i]]
	}
	if to < from {
		to = from
	}
	return from, to
}

func blankOnly(lines []string) bool {
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			return false
		}
	}
	return true
}

func splice(dst []string, from, to int, block []string) []string {
	out := make([]string, 0, len(dst)-(to-from)+len(block))
	out = append(out, dst[:from]...)
	out = append(out, block...)
	return append(out, dst[to:]...)
}

// align matches each of a's lines to the earliest unused line of b that equals
// it, and reports the first that has none.
//
// Lines are compared whole, indentation included. That is what lets a bare
// closing line match the one that closes an outermost block rather than a
// nested one, and normalising whitespace first would place an insertion wrongly
// rather than fail.
func align(a, b []string) (matched []int, failed int) {
	matched = make([]int, 0, len(a))
	at := 0
	for i, line := range a {
		for at < len(b) && b[at] != line {
			at++
		}
		if at == len(b) {
			return nil, i
		}
		matched = append(matched, at)
		at++
	}
	return matched, -1
}

// decompose splits a file into the boilerplate and the injected regions, and
// records for each boilerplate line which file line it came from and which is
// the last it owns.
//
// A line owns the region that follows it, so an insertion computed against the
// boilerplate can be placed after the region rather than inside it.
func decompose(commentPrefix string, file []string, content string) (skeleton []string, origin, owner []int, err error) {
	regions, err := inject.FindRegions(commentPrefix, content)
	if err != nil {
		return nil, nil, nil, err
	}

	inRegion := make([]bool, len(file))
	for _, r := range regions {
		from, to := lineOf(content, r.Start), lineOf(content, r.End-1)
		for i := from; i <= to && i < len(file); i++ {
			inRegion[i] = true
		}
	}

	origin = make([]int, 0, len(file))
	owner = make([]int, 0, len(file))
	for i, line := range file {
		if inRegion[i] {
			if len(owner) > 0 {
				owner[len(owner)-1] = i
			}
			continue
		}
		skeleton = append(skeleton, line)
		origin = append(origin, i)
		owner = append(owner, i)
	}
	return skeleton, origin, owner, nil
}

// structural returns the indices of lines that carry structure: not blank, and
// not a plain comment. A marker is structural, because it is what is being
// planted and a file that has one differs from a file that does not.
func structural(lines []string, commentPrefix string) []int {
	var out []int
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, commentPrefix) && len(genpkg.MarkersIn(commentPrefix, line)) == 0 {
			continue
		}
		out = append(out, i)
	}
	return out
}

func pick(lines []string, at []int) []string {
	out := make([]string, 0, len(at))
	for _, i := range at {
		out = append(out, lines[i])
	}
	return out
}

// lineOf returns the index of the line a byte offset falls in.
func lineOf(content string, offset int) int {
	if offset < 0 {
		return 0
	}
	if offset > len(content) {
		offset = len(content)
	}
	return strings.Count(content[:offset], "\n")
}

// splitLines splits content into lines and reports whether it ended with a
// newline, so that rejoining reproduces the file rather than adding or dropping
// one.
func splitLines(content string) (lines []string, trailingNewline bool) {
	if content == "" {
		return nil, false
	}
	trailingNewline = strings.HasSuffix(content, "\n")
	if trailingNewline {
		content = content[:len(content)-1]
	}
	return strings.Split(content, "\n"), trailingNewline
}

func joinLines(lines []string, trailingNewline bool) string {
	out := strings.Join(lines, "\n")
	if trailingNewline {
		out += "\n"
	}
	return out
}
