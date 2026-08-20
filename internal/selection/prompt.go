package selection

import (
	"fmt"
	"sort"
	"strings"

	"github.com/calebcowen/sedum/internal/catalog"
	"github.com/calebcowen/sedum/internal/expand"
	"github.com/calebcowen/sedum/internal/resolve"
)

// Phase 4's prompt.
//
// Nothing in this file names a language, a framework, or a file type, and that
// is a hard constraint rather than a stylistic one. Every target-specific word
// the model sees comes from the record's intent or from the catalog the
// package author wrote. A prompt that explained what a controller is would put
// target knowledge in Sedum's core through the back door, and the claim that a
// team switching stacks writes a package rather than waiting on Sedum would
// stop being true.

// systemPrompt states the bounded job.
//
// It describes the task in terms of the artifacts - a catalog, kwargs, files -
// and never in terms of what the code being generated is. The model is not
// asked to write code, name files, or decide structure, because none of those
// are things a later phase would accept from it.
const systemPrompt = `You choose code-generation actions from a fixed catalog.

You will be given a change to make, the files that already exist for it, and a
catalog of actions. Your only job is to choose which actions carry out the
change and to bind their arguments. You are not writing code: each action
already has a template, and the arguments you bind are what fill it in.

Rules:
- Use only actions named in the catalog. There are no others, and an action not
  listed does not exist.
- Bind every kwarg marked required. Do not bind a kwarg the action does not
  declare.
- Each value must match its declared type: string, int, bool, list, or literal.
  A literal is written as a JSON string, but its contents are source code that
  is emitted into the file exactly as you write it - so bind the code itself,
  spelled the way the target language spells it, and never a description of it.
  Bind that string and nothing around it: not an object naming the type, and
  not a sentence about the value. Whatever you write appears in the file
  unchanged, punctuation and all.
  There is always something correct to write; a literal is never left empty.
- Where an action lists variants, prefer one of them when the change maps
  cleanly onto it. A value outside the list is allowed only when the action
  reports a fallback, and it renders a stub rather than a full implementation.
- You do not choose which file an action writes to. Each action lists its
  injects_into patterns: the file it writes to, with {{...}} standing for a
  value derived from one of your arguments. Bind your arguments so that each
  pattern lands on one of the files listed below.

  Read a pattern forwards. Match its fixed text against a listed file and bind
  the argument to the part left over - not to the file's path, and not to the
  whole filename. If no binding can make an action's pattern reach a listed
  file, that action does not apply to this change.
- An action with no injects_into pattern takes its target as an ordinary
  argument, and its anchor says which injection point it needs. Bind that
  argument to a file listed below whose injection points include that anchor.
  An action whose anchor a file does not carry does not apply to that file, no
  matter what else fits.
- An action may be selected more than once with different arguments.

Respond with JSON and nothing else - no prose, no explanation, no fields beyond
the two shown:

{"invocations": [{"action": "<name>", "kwargs": {"<name>": <value>}}]}

If no action in the catalog applies to the change, return {"invocations": []}.`

// Prompt builds the conversation for one record.
//
// The record's intent and constraints are passed through verbatim. They are the
// author's words about their own change, and paraphrasing them would be Sedum
// deciding what the change means - which is the one judgment this whole design
// keeps out of the tool.
func Prompt(req Request, cat catalog.Catalog) ([]Message, error) {
	payload, err := catalog.ForPrompt(cat).JSON()
	if err != nil {
		return nil, fmt.Errorf("record %s: building the action catalog: %w", req.RecordID, err)
	}

	var b strings.Builder

	b.WriteString("## The change\n\n")
	b.WriteString(strings.TrimSpace(req.Intent))
	b.WriteString("\n")

	if len(req.Constraints) > 0 {
		b.WriteString("\n## Constraints on it\n\n")
		for _, c := range req.Constraints {
			fmt.Fprintf(&b, "- %s\n", strings.TrimSpace(c))
		}
	}

	writable, unmanaged := partition(req.Files)

	b.WriteString("\n## Files\n\n")
	if len(writable) == 0 {
		b.WriteString("No file in this change can be written to, so no action applies.\n")
	} else {
		b.WriteString("These files exist and are the only ones the actions below may reach:\n\n")
		for _, f := range writable {
			fmt.Fprintf(&b, "- %s\n", f)
		}
	}

	// Unmanaged paths are named rather than hidden. They are part of the
	// change and the intent may well refer to them, so a model that could not
	// see them would keep trying to reach them with an action. Saying they are
	// somebody else's is cheaper than rejecting the attempt.
	if len(unmanaged) > 0 {
		b.WriteString("\nThese are part of the change but are not written by this tool. " +
			"No action can reach them, and nothing you return should try to:\n\n")
		for _, f := range unmanaged {
			fmt.Fprintf(&b, "- %s\n", f)
		}
	}

	b.WriteString("\n## Actions\n\n")
	b.Write(payload)
	b.WriteString("\n")

	return []Message{
		{Role: RoleSystem, Content: systemPrompt},
		{Role: RoleUser, Content: b.String()},
	}, nil
}

// partition splits the run's files into the ones an action may reach and the
// ones a package declared it does not write. Both are sorted, so that one
// record's prompt is the same text every time it is built.
func partition(files []resolve.File) (writable, unmanaged []string) {
	// The injection points a file carries, read from what Phase 3 rendered.
	// An action whose target is an argument rather than a convention has no
	// pattern to match against a filename, so the anchor is what tells the
	// model which of these files it can be aimed at - and the model has no
	// way to know which files carry which anchors unless the list says so
	// (prov-2026-14c832bf).
	//
	// It is the same reading the completeness pass and Phase 5's applicability
	// check are built on, so the prompt cannot describe a file differently
	// from the rule that will judge what the model does with it.
	anchors := map[string][]string{}
	for _, a := range expand.Planted(files) {
		anchors[a.Path] = append(anchors[a.Path], a.Marker)
	}

	for _, f := range files {
		if f.Unmanaged {
			unmanaged = append(unmanaged, fmt.Sprintf("%s (declared unmanaged by %s)", f.Path, f.UnmanagedBy))
			continue
		}
		line := fmt.Sprintf("%s (package %s", f.Path, f.Package.Name)
		if planted := anchors[f.Path]; len(planted) > 0 {
			line += "; injection points: " + strings.Join(planted, ", ")
		}
		writable = append(writable, line+")")
	}
	sort.Strings(writable)
	sort.Strings(unmanaged)
	return writable, unmanaged
}
