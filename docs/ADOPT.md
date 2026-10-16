# `sedum adopt`: bringing existing code under Sedum

Status: design proposal, no code. Written 2026-10-16 for provisio-mgr and Caleb.

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

Adoption uses only what exists today: the `go` package's actions
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

## Pipeline

```
sedum adopt [--records DIR] [--record ID] [--generators DIR] [--dry-run] <path>...
```

1. **Select files.** Expand the paths to Go files. Skip `unmanaged` paths and
   generated files (`// Code generated ... DO NOT EDIT.`); report both. Files
   already carrying Sedum markers are left alone; files that no `go` file
   template matches (a root-level `main.go` today) are reported as unsupported,
   not guessed at.
2. **Read.** Parse each file with `go/parser` (the standard library; Sedum's
   `syntax: go` checker already links it, so this adds no dependency and does
   not conflict with the "no parsers" rule, which is about generator anchors,
   not adoption). Split into top-level declarations with their doc comments.
3. **Map to invocations**, in source order, per the table below.
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
7. **Other languages.** The reader is Go-only. Does a generator package declare an
   adopter (a keyed entry in `sedum.yaml`, which is strictly decoded), or is
   adoption per-language code in Sedum?
8. **Where it lives.** Sedum is a text-level tool with no language parsers by
   design (PRD). Does `adopt` belong in Sedum, or in a separate tool that writes a
   recording, which is already a public input format?
9. **Provenance side.** `provenance consolidate` (prov-lib's analysis) handles the
   records already created. Adoption should not duplicate it; it only prevents new
   ones.
10. **Generated and vendored code.** Skipped by default. Adopt generated code
    never; confirm.

## Suggested order of work

1. Decide questions 1, 2, 5, 8 (they change formats or tool boundaries).
2. Blueprint: reader, mapping, recording writer, verifier, for plain top-level
   declarations only, with the marker-strip invariant as the spec.
3. Prove it on a package of provenance's `cmd/` (handwritten parts), then linespec;
   measure marker overhead.
4. Lift pass on one clear shape (the `*_test.go` table tests or CLI flag setters).
