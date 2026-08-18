package evals

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/calebcowen/sedum/internal/record"
	"github.com/calebcowen/sedum/internal/selection"
)

// baselineSystem is what the model is told it is doing.
//
// It states the response shape and nothing about any stack. Every sentence
// naming a framework, a version or a convention would be a package rendered as
// prose, and the package is the variable under test (prov-2026-a4dbe65c).
const baselineSystem = `You write the source files a change requires.

You will be given a change's intent, its constraints, and the exact list of file
paths the change is authorized to write. Write every one of those files.

Return one markdown code fence per file. The fence's info string is the file
path, exactly as it was given to you, and the fence body is the complete
contents of that file:

` + "```" + `path/to/file.ext
the entire contents of that file
` + "```" + `

Write no prose outside the fences. Write no file that was not asked for, and
leave none of the listed files out.`

// baselinePrompt is the record, and only the record.
//
// Intent, constraints and the literal paths affected_scope names - which is
// exactly what Phase 4 is driven by, minus the catalog. A baseline handed the
// framework and its version would answer "does a catalog beat a good prompt"
// rather than "does Sedum beat not using Sedum", and the second is the question
// the arm exists for.
func baselinePrompt(rec *record.Record) []selection.Message {
	var b strings.Builder

	b.WriteString("# Intent\n\n")
	b.WriteString(strings.TrimSpace(rec.Intent))

	if len(rec.Constraints) > 0 {
		b.WriteString("\n\n# Constraints\n\n")
		for _, c := range rec.Constraints {
			fmt.Fprintf(&b, "- %s\n", strings.TrimSpace(c))
		}
	}

	b.WriteString("\n# Files to write\n\n")
	for _, p := range rec.Paths {
		fmt.Fprintf(&b, "- %s\n", p)
	}

	return []selection.Message{
		{Role: "system", Content: baselineSystem},
		{Role: "user", Content: b.String()},
	}
}

// parseFencedFiles reads the model's answer into a path-keyed map.
//
// A fence per file, its info string the path. JSON was the obvious shape and is
// the wrong one: it would ask the model to escape whole source files into
// string literals, and a baseline that failed on an unescaped heredoc would be
// measuring escaping rather than whether the model can write the application.
//
// Only the paths the record authorized are kept. A model that invents a file is
// not granted one - the sedum arm cannot either, since Phase 3 creates only
// what affected_scope names - and a path outside the list is reported rather
// than silently dropped, because a model writing the right code at the wrong
// path is a different finding from one writing nothing.
func parseFencedFiles(raw string, allowed []string) (map[string]string, []string, error) {
	ok := map[string]bool{}
	for _, p := range allowed {
		ok[p] = true
	}

	files := map[string]string{}
	var unexpected []string

	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	for i := 0; i < len(lines); i++ {
		info, isOpen := fenceInfo(lines[i])
		if !isOpen || info == "" {
			continue
		}

		var body []string
		i++
		for ; i < len(lines); i++ {
			if _, closing := fenceInfo(lines[i]); closing {
				break
			}
			body = append(body, lines[i])
		}

		p := path.Clean(strings.TrimSpace(info))
		content := strings.Join(body, "\n")
		if len(body) > 0 {
			content += "\n"
		}
		switch {
		case ok[p]:
			files[p] = content
		default:
			unexpected = append(unexpected, p)
		}
	}

	if len(files) == 0 {
		return nil, unexpected, fmt.Errorf("the answer carried no file for any authorized path")
	}
	sort.Strings(unexpected)
	return files, unexpected, nil
}

// fenceInfo reports whether a line opens or closes a fence, and what its info
// string is. A closing fence has none, which is how the two are told apart.
func fenceInfo(line string) (string, bool) {
	t := strings.TrimSpace(line)
	if !strings.HasPrefix(t, "```") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(t, "```")), true
}

// missing is the authorized paths the answer left out, in the record's order.
//
// Reported rather than fatal. A baseline that wrote four of five files is a
// measurement - it will fail to build, which is the outcome - and refusing it
// would turn an observation about the model into a harness error.
func missing(files map[string]string, allowed []string) []string {
	var out []string
	for _, p := range allowed {
		if _, ok := files[p]; !ok {
			out = append(out, p)
		}
	}
	return out
}

// baselineAnswer makes the one call a baseline sample is allowed.
//
// One call, because a baseline has no Phase 5 to reject anything and the only
// check it could fail is the build - and re-prompting with a compiler error is
// the repair loop TOOL_BOUNDARIES assigns to another tool. It is comparable to
// the sedum arm at retries zero, which is what almost every stored entry was
// drawn at (prov-2026-a4dbe65c).
func baselineAnswer(ctx context.Context, client selection.Client, rec *record.Record) (map[string]string, selection.Completion, []string, error) {
	answer, err := client.Complete(ctx, baselinePrompt(rec))
	if err != nil {
		return nil, answer, nil, err
	}
	files, unexpected, err := parseFencedFiles(answer.Content, rec.Paths)
	return files, answer, unexpected, err
}
