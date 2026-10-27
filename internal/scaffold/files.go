// Package scaffold prepares a directory for sedum: the context file, the
// .gitignore line, the provenance config's exclude list and the git hooks.
package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The status an Item reports.
const (
	StatusCreated   = "created"
	StatusUpdated   = "updated"
	StatusUnchanged = "unchanged"
	StatusSkipped   = "skipped"
)

// Item is the result of one step of Init.
type Item struct {
	Name   string
	Status string
	Detail string
}

// ContextTemplate is the text of a new context file.
const ContextTemplate = "# Context for sedum grow\n" +
	"\n" +
	"Pass with: sedum grow --generators <generators dir> --records provenance --prompt-file .sedum/context.md --only <record id>\n" +
	"Advisory only: behaviour belongs in the record, and editing this file regenerates nothing.\n" +
	"Write here what the model needs but the records should not carry.\n" +
	"\n" +
	"## Conventions\n" +
	"- (package conventions, output rules, error handling)\n" +
	"\n" +
	"## Helpers that exist (call them; do not redefine)\n" +
	"- (signatures the generated code may call)\n"

// Init prepares dir, which must be an existing directory: the context file, the
// .gitignore line for the run log and the config's exclude entries. It returns
// the three items in that order, or the error of the first step that fails.
func Init(dir string) ([]Item, error) {
	steps := []func(string) (Item, error){ensureContext, ensureGitignore, ensureExcludes}
	names := []string{".sedum/context.md", ".gitignore", ".provenance.yml"}
	var items []Item
	for i, step := range steps {
		item, err := step(dir)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", names[i], err)
		}
		items = append(items, item)
	}
	return items, nil
}

func ensureContext(dir string) (Item, error) {
	item := Item{Name: ".sedum/context.md"}
	path := filepath.Join(dir, ".sedum", "context.md")
	if _, err := os.Stat(path); err == nil {
		item.Status = StatusUnchanged
		return item, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return item, err
	}
	if err := os.WriteFile(path, []byte(ContextTemplate), 0o644); err != nil {
		return item, err
	}
	item.Status = StatusCreated
	return item, nil
}

const ignoreLine = ".sedum/run.log"

func ensureGitignore(dir string) (Item, error) {
	item := Item{Name: ".gitignore"}
	path := filepath.Join(dir, ".gitignore")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		item.Status = StatusCreated
		return item, os.WriteFile(path, []byte(ignoreLine+"\n"), 0o644)
	}
	if err != nil {
		return item, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == ignoreLine {
			item.Status = StatusUnchanged
			return item, nil
		}
	}
	text := string(data)
	if text != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	info, err := os.Stat(path)
	if err != nil {
		return item, err
	}
	item.Status = StatusUpdated
	return item, os.WriteFile(path, []byte(text+ignoreLine+"\n"), info.Mode().Perm())
}

var wantedExcludes = []string{"generators/**", ".sedum/**"}

const excludesError = "exclude_paths is not a one-line list or a block list; add generators/** and .sedum/** to it by hand"

func ensureExcludes(dir string) (Item, error) {
	item := Item{Name: ".provenance.yml"}
	path := filepath.Join(dir, ".provenance.yml")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		item.Status = StatusSkipped
		item.Detail = "no .provenance.yml in this directory"
		return item, nil
	}
	if err != nil {
		return item, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return item, err
	}
	updated, changed, err := addExcludes(string(data))
	if err != nil {
		return item, err
	}
	if !changed {
		item.Status = StatusUnchanged
		return item, nil
	}
	item.Status = StatusUpdated
	return item, os.WriteFile(path, []byte(updated), info.Mode().Perm())
}

// addExcludes adds the wanted entries missing from the exclude_paths list of a
// config text, which is a one-line flow list or a block list, or absent.
func addExcludes(text string) (string, bool, error) {
	lines := strings.Split(text, "\n")
	key := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "exclude_paths:") {
			key = i
			break
		}
	}
	if key < 0 {
		var missing []string
		for _, w := range wantedExcludes {
			missing = append(missing, "'"+w+"'")
		}
		if text != "" && !strings.HasSuffix(text, "\n") {
			text += "\n"
		}
		return text + "exclude_paths: [" + strings.Join(missing, ", ") + "]\n", true, nil
	}
	rest := strings.TrimSpace(strings.TrimPrefix(lines[key], "exclude_paths:"))
	switch {
	case strings.HasPrefix(rest, "["):
		return addToFlowList(lines, key)
	case rest == "" || strings.HasPrefix(rest, "#"):
		return addToBlockList(lines, key)
	}
	return "", false, fmt.Errorf("%s", excludesError)
}

// commentStart returns the index of a "#" preceded by a blank that starts a
// trailing comment, outside quotes, or the length of the line.
func commentStart(line string) int {
	var quote byte
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '\'' || c == '"':
			quote = c
		case c == '#' && i > 0 && (line[i-1] == ' ' || line[i-1] == '\t'):
			return i
		}
	}
	return len(line)
}

// flowEntries splits the inside of a flow list at commas that lie outside quotes,
// brackets and slash-delimited regular expressions.
func flowEntries(inner string) []string {
	var entries []string
	var quote byte
	depth, start := 0, 0
	inRegex := false
	for i := 0; i < len(inner); i++ {
		c := inner[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case inRegex:
			if c == '/' {
				inRegex = false
			}
		case c == '\'' || c == '"':
			quote = c
		case c == '/' && strings.TrimSpace(inner[start:i]) == "":
			inRegex = true
		case c == '[':
			depth++
		case c == ']':
			depth--
		case c == ',' && depth == 0:
			entries = append(entries, inner[start:i])
			start = i + 1
		}
	}
	return append(entries, inner[start:])
}

func unquote(entry string) string {
	entry = strings.TrimSpace(entry)
	if len(entry) >= 2 && (entry[0] == '\'' || entry[0] == '"') && entry[len(entry)-1] == entry[0] {
		entry = entry[1 : len(entry)-1]
	}
	return strings.TrimSpace(entry)
}

func missingExcludes(present []string) []string {
	have := map[string]bool{}
	for _, p := range present {
		have[p] = true
	}
	var missing []string
	for _, w := range wantedExcludes {
		if !have[w] {
			missing = append(missing, w)
		}
	}
	return missing
}

func addToFlowList(lines []string, key int) (string, bool, error) {
	line := lines[key]
	body := line[:commentStart(line)]
	open := strings.Index(body, "[")
	closing := strings.LastIndex(body, "]")
	if closing < open {
		return "", false, fmt.Errorf("%s", excludesError)
	}
	inner := body[open+1 : closing]
	var present []string
	empty := true
	for _, e := range flowEntries(inner) {
		if strings.TrimSpace(e) != "" {
			empty = false
			present = append(present, unquote(e))
		}
	}
	missing := missingExcludes(present)
	if len(missing) == 0 {
		return strings.Join(lines, "\n"), false, nil
	}
	var add strings.Builder
	for i, m := range missing {
		if i > 0 || !empty {
			add.WriteString(", ")
		}
		add.WriteString("'" + m + "'")
	}
	lines[key] = line[:closing] + add.String() + line[closing:]
	return strings.Join(lines, "\n"), true, nil
}

func addToBlockList(lines []string, key int) (string, bool, error) {
	last := key
	indent := ""
	haveIndent := false
	var present []string
	for i := key + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "- ") {
			if !haveIndent {
				indent = lines[i][:len(lines[i])-len(strings.TrimLeft(lines[i], " \t"))]
				haveIndent = true
			}
			present = append(present, unquote(strings.TrimPrefix(trimmed, "- ")))
			last = i
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		break
	}
	if !haveIndent {
		indent = "  "
	}
	missing := missingExcludes(present)
	if len(missing) == 0 {
		return strings.Join(lines, "\n"), false, nil
	}
	var added []string
	for _, m := range missing {
		added = append(added, indent+"- '"+m+"'")
	}
	out := append(append(append([]string{}, lines[:last+1]...), added...), lines[last+1:]...)
	text := strings.Join(out, "\n")
	if last == len(lines)-1 && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	return text, true, nil
}
