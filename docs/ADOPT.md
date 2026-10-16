# `sedum adopt`: bringing existing code under Sedum

Status: design proposal, no code. Written 2026-10-16 for provisio-mgr and Caleb; revised to be language-agnostic with Go as the first adopter.

## Problem

Sedum manages code it grew. A repository with little or no grown code has
nothing for `grow --check`, `--repair`, `replay --restamp-gen` or the
`revise`-an-existing-record workflow to act on, so every change opens a new
blueprint/imprint pair and records pile up (sedum 171, provenance 84, linespec
526). The way to reach a "critical mass" where most work is a `revise` of an
existing record is to make the existing code Sedum-managed, without a model and
without changing what the code does.

## Principle

Adoption is a **deterministic transcription**, not a generation. It reads Go
source, emits a recording whose invocations reproduce that source, and replays
it. The model is never involved, so adoption is cheap, repeatable and
reviewable as a diff. Making the code *better structured* is a separate step,
the lift passes below, and that step is checked against the transcription.

Adoption is language-agnostic in the core and language-specific in generator
packages, like everything else in Sedum. The core never parses a language: a
generator package declares an **adopter** that splits a file into units, and the
core does the rest. Go is the first adopter (below, *Adopters*). The rest of this
section and the Go mapping describe that first adopter.

Adoption for Go uses only what exists today: the `go` package's actions
(`setFunction`, `setMethod`, `setType`, `setVar`, `setConst`, `addImport`), the
file templates, recordings, `replay --execute`, `--stamp`/`--restamp-gen`, and
the marker format. New code is limited to a reader (Go source to invocations),
a driver, and a verifier.

## What "byte-for-byte" can and cannot mean

A managed file carries marker comments and the blank lines around them (see the
README's *Ownership and Idempotency*). So an adopted file can never equal the
original file; it equals the original **plus markers**. The invariant is:

> Strip every `sedum:` marker line from the adopted file and `gofmt` it. The
> result is byte-identical to `gofmt` of the original file.

Consequences, stated so nobody is surprised in review:

- Adoption rewrites each adopted file once. Its diff is large and mechanical:
  one marker pair per declaration. The verifier proves the diff is only markers.
- Files are compared after `gofmt`. A file that is not gofmt-clean fails the
  precondition (the `syntax: go` check already requires gofmt-clean output);
  `adopt` reports it and offers `--gofmt-first` as a separate prior commit.
- Blank lines the original had between declarations are normalised by the
  marker layout. gofmt keeps at most one, so this is invisible after gofmt.

## Pipeline (Go as the example)

```
sedum adopt [--records DIR] [--record ID] [--generators DIR] [--dry-run] <path>...
```

1. **Select files.** Expand the paths to Go files. Skip `unmanaged` paths and
   generated files (`// Code generated ... DO NOT EDIT.`); report both. Files
   already carrying Sedum markers are left alone; files that no `go` file
   template matches (a root-level `main.go` today) are reported as unsupported,
   not guessed at.
2. **Split.** Ask the file's package for its adopter (see *Adopters*). The
   adopter returns an ordered list of units: a header span, then top-level
   units with an identity, an action, kwargs and the exact source text, and
   a list of residual spans it could not claim. The core never looks inside a
   unit.
3. **Map to invocations**, in source order. For Go the adopter's mapping is the
   table below.
4. **Write a recording** with one record per target imprint, one `default`
   phase, invocations in source order. Invocation order within a record fixes
   region order at a shared anchor, so source order is preserved.
5. **Replay** it (`replay --execute`, scope-checked against the records). Markers
   get `gen` from the record and `pkg` from the package, so hashes are correct
   without a model and without `--stamp`.
6. **Verify.** The marker-stripped invariant above, `go build ./...` and
   `go vet` on the adopted packages, then the project's own tests. Any failure
   restores the tree (Phase 7 already writes all-or-nothing) and prints the
   first differing declaration.
7. **Report** counts: files, declarations by action, regions with literal code.

`--dry-run` stops after step 4 and writes only the recording, so a human can
read exactly what would be claimed.

### Declaration mapping

| Source | Invocation |
| --- | --- |
| `func F(...)` | `setFunction {file, name, code}` |
| `func (r *T) M(...)` | `setMethod {file, receiver: T, name, code}` |
| `type T ...`, single or grouped | `setType` per spec (see grouped, below) |
| `var` / `const` | `setVar` / `setConst` |
| one `import "x"` line | `addImport {file, path}` |
| `init()` | `setFunction` named `init`, numbered `init#N` (identity must be unique) |

`code` is the exact declaration text including its doc comment.

Constructs the current actions cannot express, each of which needs a decision
before implementation (see open questions): file headers (license comment,
`//go:build`, the package doc), floating comments between declarations,
`//go:embed` and `//go:generate` directives, cgo preambles, grouped `var (...)`
and `const (...)` blocks with iota, `import (...)` grouping, and blank
identifiers (`var _ I = (*T)(nil)`).

## Adopters

An **adopter** is the one language-specific piece of adoption. A generator
package declares it in `sedum.yaml`; the core calls it, validates what it
returns, and does everything else. Two forms:

**Declarative.** Regexes and delimiters, in the style of anchors, so no parser:

```yaml
adopt:
  header: '\A(#!.*\n)?(""".*?"""\n)?'      # kept as the file template's head
  unit_start: '^(def|class|async def)\s+(?P<name>\w+)'   # a unit starts here
  unit_end: dedent                          # next line at column 0, or EOF
  attach_above: ['^@', '^#']                 # decorators and comments belong to the unit
  action: {def: setFunction, class: setType}
```

**External command.** The package names an executable that reads a file path
and writes a recording fragment to stdout (the same JSON as a recording, which
is already a public input format: README *The recording is an input format*).
The core runs it, validates the fragment exactly as it validates a model's
response, and treats failures as terminal. The Go adopter is of this form,
built into the `go` package, because Go needs a real parse (grouped blocks,
raw strings containing `//`).

Both forms are subject to the same checks, whoever wrote them:

- **Coverage.** Header + unit texts + residual spans, concatenated in order, must
  equal the original file exactly. Anything the adopter cannot claim is a
  *residual* and is kept as a verbatim, unmanaged span, never dropped.
- **Round trip.** The marker-stripped invariant (below) is the contract, and it
  is checked by the core on the replayed result, not trusted from the adopter.
- **Unmanaged residue is reported**, per file, so a package that adopts 60% of
  a file says so.

The strictly decoded `sedum.yaml` means the `adopt` key must be added to the
loader first, the same ordering that `test_paths` followed.

### Questions 7 and 8, answered

*Q8, where it lives.* The core stays parserless. `sedum adopt` is a core
command, but it only orchestrates: select files, call the package's adopter,
validate coverage, replay, verify. Splitting source is the adopter's job, and
an adopter is generator-package content, exempt from provenance like the rest of
`generators/**`. Trade-off: the core gains a command and a schema key, and the
`adopt` command is a new public surface; in return nothing in the core knows a
language, and a team can write an adopter for its language without touching
Sedum. The alternative, a separate tool that writes a recording, is already
possible today and stays possible (it is the external-command form run by hand),
but then verification and hashes are each tool's own problem.

*Q7, how a package declares it.* A keyed `adopt:` entry in `sedum.yaml`, either
declarative or `command:`. Trade-off: declarative adopters are portable and
cheap to review but cover regular languages of top-level units only (see below);
command adopters can use real parsers but are code a user must trust and run, so
`adopt` should refuse a command adopter unless `--allow-command-adopters` is
given, and print the command first. Go ships as a built-in command adopter.

## Sanity check against non-Go shapes

Paper exercises; nothing here is implemented or measured.

### Python or TypeScript (decorators, exports)

A Python module: shebang and docstring header, imports, then `def`, `class` and
module-level assignments, decorators above defs, comments above decorators.

- **Expressible declaratively:** header (shebang + docstring), one unit per
  top-level `def`/`class` ended by dedent, decorators and comment runs attached
  above (`attach_above`), imports as one-line units.
- **Not expressible declaratively:** unit end inside multi-line strings that contain
  dedented text, nested decorators with arguments across lines, `if __name__ ==
  '__main__':` blocks (a unit with no name), conditional imports, and
  top-level statements with side effects, whose order matters. These become
  residual spans or force a command adopter.
- **Ordering:** `class` bodies are one unit; methods are not separate regions
  unless the package defines an action with an inner anchor. So the first
  adoption of a large class is one literal region, a known weakness (lifting
  methods needs class-body anchors the file template must plant, which an
  existing file does not have).

TypeScript adds `export` modifiers, `export default`, overloads (several
declarations with one name: identity must include a signature hash, or the
overload group is one unit), and `declare module` blocks. Brace matching ending
a unit is declarative only if the regex flavour supports balanced delimiters;
template literals containing braces break a simple counter. Expected result:
declarative for plain modules, command adopter (TypeScript compiler API) for the
rest.

### C++ (header + implementation)

A declaration in `foo.h` and its definition in `foo.cpp` are one logical thing
in two files.

- **Adoption handles each file independently.** The core would give them two
  units with different identities (`Foo::bar` declared, `Foo::bar` defined).
  The round trip passes, but the pair is not recognised as one thing, so a
  `revise` has to touch both regions by hand.
- **Cannot be expressed:** macros that generate declarations, include guards
  that wrap the whole file (a header span ending at EOF, so coverage works but
  `decls` anchors sit inside it), templates split across headers, and
  `#ifdef` regions that share one unit. An adopter can mark such spans residual.
- **What the design would need to add:** a *group id* shared by the units of two
  files (a kwarg, and a fact the marker already carries), so that a lift can
  join them. That belongs to the lift, not to adoption.

### What the declarative adopter cannot express, in summary

Units with no name; units ended by balanced delimiters in languages with
nested string syntax; order-dependent side-effect statements; one logical
unit across two files; macros/preprocessing; a unit nested in another
(methods in a class) without anchors the existing file lacks; and any language
whose top-level boundary depends on type or scope information.

## Language-free parts

The following is identical for every language and lives in the core:

- **Marker layout.** `comment_prefix`/`comment_suffix` per package, the `sedum:`
  namespace, the attribute JSON object, `/sedum:` closing marker, indentation
  copied from the anchor.
- **Region identity and ownership.** Action + variant + identity kwargs; `record`
  as owner; reconcile by identity; conflicts as errors; order at shared anchors.
- **Hashes.** `gen` is over the record (intent, constraints, scope), `pkg` over
  the action and its templates. Neither looks at the language. Adoption only
  needs a record to exist, then replay stamps both.
- **Recording and replay.** Format, phases, scope validation, `--restamp-gen`.
- **The verifier's shape:** strip markers, normalise, compare. Only the
  *normaliser* is language-specific (gofmt for Go; `black`, `prettier` or
  `clang-format` for others, declared by the package; none means byte-exact
  compare, which is stricter and often what a language without a canonical
  formatter needs).
- **Lift scaffolding.** Cluster by normal form, re-render, compare. The
  normalisation of identifiers and literals is the one language-specific part
  and can be the same declarative style (token classes) or the command adopter.

Language-specific, inside the package: the file templates, the actions, the
adopter, the syntax checker, the normaliser.

## Agent-assisted adopter authoring

Writing an adopter is a small, checkable task, which is the shape an agent does
well: the output is a short `adopt:` entry (or a small command), and the core
supplies an exact verifier. So Sedum would ship a Claude Code skill for it.

**`/sedum:adopter`** drafts and refines an adopter for one generator package:

1. Inspect the repo: which files the package's extensions claim, a sample of
   their top-level shapes, and the existing `sedum.yaml`.
2. Draft the `adopt:` entry (declarative first; a command only when the
   declarative form cannot express the shapes it found).
3. Run `sedum adopt --dry-run` and read the coverage report: files adopted,
   units found, and the **residual spans** with file and line.
4. Refine and repeat. Stop when coverage is acceptable or progress stalls
   (three rounds with no fewer residuals, or a declared round limit).
5. Stop with a report: the adopter, the coverage numbers, what stayed residual
   and why, and the commands the user would run next.

**`/sedum:generator`** (later, separate) would author a whole generator package
(file templates, actions, an adopter) from a repository's existing code. It
reuses the same loop and adds the lift passes' proposals as its output. It is
listed here only so the plugin's shape leaves room for it.

### What the skill may and may not do

May: read the repo; write the adopter entry in the package's `sedum.yaml` (or an
adopter command file beside it) and the dry-run recording; run `sedum adopt
--dry-run`, `sedum validate` and read-only git commands.

May not, without the user saying so in that session: replay or write any
managed file, create or edit provenance records, commit or push, or edit
anything outside the adopter entry and recording. A command adopter runs only
with `--allow-command-adopters`, which the skill passes only after showing the
user the command and getting a yes. It never uses the model to produce the
adopted code itself; the adopter only splits, so the invariant stays checkable
by the core, not by the agent's say-so. Model calls it triggers elsewhere go
through the herd's pinned wrapper like any other.

### How it ships

Assets (the skill files and a plugin manifest) live in the sedum repo, as the
provenance plugin's do in its repo's `assets/`, and are embedded in the binary
and written by `sedum install-plugin` / `sedum install-skills` into the user's
Claude Code config. Open point: whether the two plugins share one marketplace
entry or stay separate (they have different release cadences, which favours
separate).

### Order of work

1. **Adopt core**: `adopt:` schema key, the coverage check, `sedum adopt` with
   `--dry-run`, the marker-stripped verifier, the report format. The report must
   be machine-readable (JSON) as well as text, since the skill reads it.
2. **Go adopter**, built into the `go` package, as the first real adopter and
   the reference for the declarative dialect.
3. **Skill and plugin**, last: it needs a stable dry-run report and at least one
   adopter to be tested against. Evals reuse the existing `evals/` harness with
   cases of "repo X, language Y, expected coverage at least Z".

Nothing is built until PR #11 merges and the open questions it raises are
settled.

## Records adoption creates

Sedum does not write provenance records (TOOL_BOUNDARIES). Adoption therefore
needs the records to exist first, and its driver prints the commands rather than
running them.

- **One blueprint per adoption effort** ("Bring existing code under Sedum"),
  scoped to the paths being adopted.
- **One imprint per Go package** (directory), not per file. Its `affected_scope`
  is `<dir>/*.go`, minus tests unless `--tests` is passed, its intent is a
  generated one-paragraph inventory ("adopts 14 declarations in 3 files; no
  behaviour change"), and `constraints` state the invariant. Code never goes in
  the intent; it is in the recording and the markers.
- `gen_hash` covers intent, constraints and scope, so an imprint's `gen` is
  stable and small however large the package. This is the main difference from
  today's literal-code imprints, whose intents carry Go source inline.
- Order: create blueprint and imprint, `sedum adopt`, then `provenance complete`
  (the imprint's spec runs `sedum grow --check` on the package). Imprint scope
  must fit the blueprint's (PROV023), so the blueprint is scoped first.

## A repository with zero Sedum code

1. `provenance init` and `sedum generators install go` (the built-in go package
   and its `generators/go.mod`).
2. Add `.gitattributes`, `exclude_paths` for generators, and `unmanaged` entries
   (`go.mod`, `go.sum` already are).
3. Commit the baseline and make sure `gofmt -l .` is empty.
4. Create the blueprint and one imprint per package. A helper
   (`sedum adopt --plan`) lists packages and the commands.
5. Adopt package by package, leaves first, committing each. A failing package
   does not block the rest, and a partly adopted repo is a valid state:
   unadopted files are simply not Sedum's.

The only requirement on the repository is that its Go files match a `go` file
template. Root-level `.go` files do not today (README's file-template list); the
smallest fix is a `{name}.go` template, to be decided as part of this work.

## Ownership tier

Adopted code is still being hand-edited. If regions are `owned`, a hand edit is
overwritten by the next `grow` or `replay`. Options: write adopted regions
`seeded` (Sedum records but never rewrites them), or keep `owned` and make
`revise` the only edit path. The marker already carries `tier`, and `seeded`
exists on read today, so adopt can choose per run without a format change.
Recommendation: `owned` by default, since the point is the `revise` workflow, with
`--seeded` as a transition mode. See open question 3.

## Lift passes

Adoption leaves every region as a literal (`code` kwarg). A **lift** replaces a
group of literal invocations with higher-level ones, under a check that the
re-rendered region equals the literal region.

```
sedum adopt lift [--dry-run] --record ID [--generators DIR]
```

1. **Find shapes.** Normalise each literal region (replace identifiers and string
   and number literals with typed holes), cluster by normal form. A cluster of
   N >= 3 with identical normal form is a candidate: N getters, N cobra
   commands, N handler methods.
2. **Match or propose.** If an existing action's template renders every member
   byte-for-byte (after gofmt), the candidate is a *match*. Otherwise it is a
   *proposal*: a draft action and template written under `generators/` for a
   person to review (generators are exempt from provenance). A model may draft
   the template at author time through `sedum-haiku`; it never touches the
   adoption itself.
3. **Rewrite the recording.** Replace the N invocations with N invocations of the
   new action (kwargs are the holes), or one with a list kwarg where the action
   allows. The imprint is **revised** (new version, same record, rule B), then the
   recording is replayed with `--restamp-gen`.
4. **Verify** as in adopt, plus: each lifted region's re-render equals the
   literal one it replaced.

Because the verification is mechanical, a lift never needs a model call at
replay time, and a lift that does not reproduce the original is rejected, not
approximated. The "edit a template or recording and replay" path in COMMON.md
applies directly; the open gap is that changing a recording by hand needs a
sanctioned way to update hashes, which `replay --restamp-gen` already provides
for `gen`, and a `pkg` restamp on template edits is what `grow` re-render does.

## Cost and risks

- **Markers carry kwargs**, and for the generic go actions kwargs include the
  whole `code`. Adopted source therefore appears twice in each file, once in the
  marker JSON and once as the body. For a 5,000-line package this roughly doubles
  its size. It shrinks as lifts replace `code` with small kwargs.
- **Churn.** Adoption touches every adopted file once. Do it per package, in its
  own commits, before other changes land on them.
- **Merge conflicts** in adopted files become conflicts in marker JSON. Mitigated by
  keeping adoption commits isolated.
- **False confidence.** A hash says a region matches its record, not that the
  record is meaningful. Adoption imprints are explicitly inventory, not intent.
- **gofmt/goimports drift** between versions changes bytes; the verifier uses the
  same gofmt the `syntax: go` checker uses.

## Open questions

1. **Marker size.** Should markers for literal regions store `code` or a hash of it
   and keep the text only in the file? The text is recoverable from the body, so
   kwargs could elide it for literals, at the cost of replay then needing the
   file. This changes the replay contract (README *Replay semantics*).
2. **Headers and floating comments.** Add a `setHeader` action anchored at
   `start_of_file`, and an `anchor` for comment-only gaps, or adopt only files whose
   non-declaration content is just the package clause and imports, and report the rest?
3. **Tier.** `owned` or `seeded` by default for adopted regions (above).
4. **Grouped declarations.** Splitting `const ( A = iota; B )` or `var ( ... )`
   into regions changes meaning or style; treat the whole block as one `setConst` /
   `setVar` region. Is one region per block acceptable as the default?
5. **Import blocks.** The go package writes one `import "x"` per region because a
   marker inside a block breaks gofmt, so adoption flattens grouped imports and loses
   stdlib/third-party grouping. Is that acceptable, or does the go package need a
   block-aware import region?
6. **Tests.** `test_paths` is validated but unread. Adopt tests in the same
   imprint, in a second imprint, or not at all by default?
7. **Other languages.** Answered above (*Adopters*): a keyed `adopt:` entry,
   declarative or external command, Go built in. Still open: the declarative
   dialect (regex flavour, how a unit end is specified) and whether command
   adopters need a trust model beyond `--allow-command-adopters`.
8. **Where it lives.** Answered above: orchestration in the core, splitting in the
   package. Still open: whether the core should refuse a package with no
   adopter or fall back to treating the whole file as one residual.
11. **Formatter per language.** The verifier needs a normaliser; is it declared in
    the package (`adopt.normalize: black`), and is a missing formatter an error?
12. **Cross-file units** (C++ header/impl, Go test beside code). Group id on the
    marker, or leave to lift?
9. **Provenance side.** `provenance consolidate` (prov-lib's analysis) handles the
   records already created. Adoption should not duplicate it; it only prevents new
   ones.
10. **Generated and vendored code.** Skipped by default. Adopt generated code
    never; confirm.

## Suggested order of work

1. Decide questions 1, 2, 5, 7, 8 (they change formats or tool boundaries).
2. Blueprint: reader, mapping, recording writer, verifier, for plain top-level
   declarations only, with the marker-strip invariant as the spec.
3. Prove it on a package of provenance's `cmd/` (handwritten parts), then linespec;
   measure marker overhead.
4. Lift pass on one clear shape (the `*_test.go` table tests or CLI flag setters).
