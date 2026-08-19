package evals

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

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
// intentSystem is baselineSystem without the two promises the intent arm cannot
// keep: it is given no constraints and no list of paths, so it has to decide for
// itself which files a change needs and what belongs in each.
//
// The fenced-file format is unchanged, because the arms differ in what the model
// is told about the change and not in how it answers.
const intentSystem = `You write the source files a change requires.

You will be given a change's intent. Decide which files the change needs, what
each of them is called, and what belongs in them.

Return one markdown code fence per file. The fence's info string is the file
path and the fence body is the complete contents of that file:

` + "```" + `path/to/file.ext
the entire contents of that file
` + "```" + `

Write no prose outside the fences.`

// intentPrompt is the record's intent and nothing else - no constraints, no
// list of files to write.
//
// The three arms are a ladder of what the model was given: a prompt, a record, a
// record and a catalog. The point of the first rung is what a sentence alone
// produces, so anything else in the prompt is a rung it is not standing on
// (prov-2026-672c6471).
func intentPrompt(rec *record.Record) []selection.Message {
	var b strings.Builder
	b.WriteString("# Intent\n\n")
	b.WriteString(strings.TrimSpace(rec.Intent))
	return []selection.Message{
		{Role: "system", Content: intentSystem},
		{Role: "user", Content: b.String()},
	}
}

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
// A nil allowed list accepts every path the answer carries. That is the intent
// arm, which is told no paths and so cannot be held to them: filtering against a
// list it never saw would score whether it guessed the ones this standard
// happens to use. It is therefore not a path-for-path comparison with the other
// arms, which is the point rather than a defect - it is measuring the absence of
// Sedum (prov-2026-672c6471).
func parseFencedFiles(raw string, allowed []string) (map[string]string, []string, error) {
	anyPath := allowed == nil
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
		case anyPath, ok[p]:
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
// The arm decides two things and nothing else: what the model is told, and
// whether the answer is held to the record's paths. Everything after - the
// scaffold, the build, the boot, the assertions - is the same question whoever
// wrote the code (prov-2026-672c6471).
func baselineAnswer(ctx context.Context, client selection.Client, rec *record.Record, arm string) (map[string]string, selection.Completion, []string, error) {
	prompt := baselinePrompt(rec)
	allowed := rec.Paths
	if arm == "intent" {
		prompt = intentPrompt(rec)
		// nil, not empty: told no paths, it cannot be held to them.
		allowed = nil
	}

	answer, err := client.Complete(ctx, prompt)
	if err != nil {
		return nil, answer, nil, err
	}
	files, unexpected, err := parseFencedFiles(answer.Content, allowed)
	return files, answer, unexpected, err
}

// baselineSample draws one baseline sample: one call for the files, then the
// behaviour harness over what it wrote.
//
// The outcome vocabulary is the sedum arm's, and it has to be, or the two
// cannot be read side by side. An answer that carried no authorized file is
// invalid - the model answered and what it said was unusable, which is what
// invalid means on the other arm too - and an endpoint that was not there is
// failed, outside every denominator.
func baselineSample(ctx context.Context, c Case, model string, opts Options) Sample {
	started := time.Now()

	client, err := selection.NewOpenAI(model)
	if err != nil {
		return Sample{Err: err, Elapsed: time.Since(started)}
	}

	rec, err := baselineRecord(c)
	if err != nil {
		return Sample{Err: err, Detail: firstLine(err.Error()), Elapsed: time.Since(started)}
	}

	files, answer, unexpected, err := baselineAnswer(ctx, client, rec, c.Arm)
	s := Sample{
		Counts:           map[string]int{},
		Calls:            1,
		PromptTokens:     answer.PromptTokens,
		CompletionTokens: answer.CompletionTokens,
	}
	if err != nil {
		// An unusable answer and an unreachable server are told apart the same
		// way the sedum arm tells them apart: one is a measurement and the
		// other is not (prov-2026-e6969eb3).
		if answer.Content == "" {
			s.Err = err
			s.Detail = firstLine(err.Error())
			s.Elapsed = time.Since(started)
			return s
		}
		s.Invalid = true
		s.Rejected = 1
		s.Detail = firstLine(err.Error())
		s.Elapsed = time.Since(started)
		return s
	}

	// Written rather than counted. There is no catalog, so what a baseline
	// produced is a set of paths, and Total carries how many of the authorized
	// ones it wrote so the report has a denominator that is not a selection.
	for p := range files {
		s.Counts[p]++
	}
	s.Total = len(files)
	s.Files = files
	s.Missing = missing(files, rec.Paths)
	s.Unexpected = unexpected

	// The one derived signal a baseline can carry: a file is a file, whoever
	// wrote it. Fill and idempotency need anchors and invocations, and a
	// baseline has neither.
	s.Syntax = syntaxOfContents(ctx, c.Check, files)

	run := RunBehaviorFiles(ctx, c.Expect.Behavior.Target, files, c.Variables)
	s.Behavior = &run
	s.Elapsed = time.Since(started)
	return s
}

// baselineRecord loads the single record a baseline case is measured on.
//
// One record, because the prompt is the record: two would mean two calls and a
// sample whose files came from separate answers, which is not what the sedum
// arm does either - it makes one call per record and this arm has no way to
// keep the halves apart.
func baselineRecord(c Case) (*record.Record, error) {
	set, _, err := record.Load(c.Records, record.Options{Only: c.Only})
	if err != nil {
		return nil, err
	}
	if len(set.Records) != 1 {
		return nil, fmt.Errorf("case %s resolves to %d records; a baseline case names exactly one, because the prompt is the record", c.ID, len(set.Records))
	}
	return set.Records[0], nil
}
