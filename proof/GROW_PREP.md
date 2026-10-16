# grow-prep: hash-aware grow

## Summary

`sedum grow` becomes safe to run on every record, every time. Each region's marker now records the version of the record that wrote it (`gen`) and of the package that rendered it (`pkg`). From those alone, with no manifest or lockfile, grow decides per record:

- **skip** when nothing changed: no model call, no file written;
- **re-render** when only the package changed: regions are rebuilt from their own markers, still no model call;
- **re-select** when the record changed: the model is asked again, and the record's regions are reconciled (replaced, inserted, and the ones it no longer selects pruned).

A region belongs to the record named on its marker; another record selecting it is an error naming both. `--check` reports staleness for CI, `--stamp` lets a migrated repository adopt existing code without re-selecting it, and `--prune-orphans` removes the regions of deleted or deprecated records. `sedum regions --json` and `grow --json` give the provenance tools exactly which code belongs to which record version and which records changed.

Alongside it: `gen_hash` and record decoding now come from the shared provenance library (`github.com/livecodelife/provenance`), with identical hashes on all 1,100 real records and 938 migrated ones; a `quote` transform; and a generic `go` generator package whose output must parse, be gofmt-clean and import what it uses, checked by a closed, built-in syntax registry.

Every output below was captured by [`proof/capture.sh`](capture.sh) from the binary built at the commit named in each file's header, against a local stub model ([`proof/stub_model.py`](stub_model.py)) that logs every call, so call counts are exact. `proof/capture.sh --check` re-runs and re-asserts all of it without writing; the parity section also needs `PROVENANCE_REAL_REPOS` pointing at a workspace with the four repositories.

The lifecycle testbed (`evals/lifecycle/`, from migrate-dry) is the runnable end-to-end demo: v1 records, migration, and every grow mode in one command with a pass/fail table.

**Size:** 96 commits on top of `main`, 109 files changed, 21 new provenance records (listed at the end).

---

## What changed, with proof

### Markers carry `gen` and `pkg` — `prov-2026-66bea9f3`

Every region records the `gen_hash` of the record that wrote it and the hash of the action and templates it was rendered from. Both are omitted when absent, so markers written before them are unchanged; absent reads as stale, never current. The marker conformance corpus gained the cases, and `conformance/pkghash/` publishes the `pkg` encoding.

[`out/01-first-grow.txt`](out/01-first-grow.txt), [`out/01-markers.txt`](out/01-markers.txt):

```
# sedum:addBeforeFilter {"tier":"owned","record":"prov-2026-aaaaaaaa","gen":"sha256:3d30…","pkg":"sha256:5949…","kwargs":{"controller":"users","filter":"authenticate"}}
before_action :authenticate
# /sedum:addBeforeFilter
```

The first grow made 4 model calls: one selection per record, and one completeness question per record about the controller's `class_body` anchor, which these records leave empty.

### A second grow does nothing — `prov-2026-9621daed`

[`out/02-second-grow.txt`](out/02-second-grow.txt): both records `(current; no model call)`, **0 model calls**, and every file's bytes and modification time unchanged. This is Phase 2's done-check.

### A template change re-renders without the model — `prov-2026-9621daed`

[`out/03-template-change.txt`](out/03-template-change.txt): after editing the action's template, both records report `(package changed; regions re-rendered from their markers, no model call)`, **0 model calls**, and [`out/03-markers.txt`](out/03-markers.txt) shows the new template's output in place.

### `--check` for CI — `prov-2026-9621daed`

[`out/04-check-stale.txt`](out/04-check-stale.txt): after one record's intent changes, `stale prov-2026-aaaaaaaa`, `current prov-2026-bbbbbbbb`, exit 1, nothing written. [`out/04-regrow.txt`](out/04-regrow.txt): the next grow spends its 2 calls on the changed record only.

### Region ownership — `prov-2026-375f5d4a`, `prov-2026-4df89974`

A region belongs to the record its marker names (or one of that record's aliases). [`out/05-ownership.txt`](out/05-ownership.txt):

```
sedum: file app/controllers/users_controller.rb: record prov-2026-cccccccc selected addBeforeFilter, which record prov-2026-aaaaaaaa owns (identity {"controller":"users","filter":"authenticate"}); revise prov-2026-aaaaaaaa, or change prov-2026-cccccccc's selection
```

Grow reconciles a re-selected record inside Phase 7's all-or-nothing write: replace, insert, and prune its own regions it no longer selects; an empty selection over owned regions needs `--force`. Records sharing a file grow in one run (the duplicate-path check is gone), and new regions at a shared anchor are ordered by owner `created_at`, then id, so regenerating from scratch is byte-identical. This replaces `prov-2026-36c8a99c`'s last-writer rule.

### `--stamp` — `prov-2026-cc73cf61`

For repositories whose markers predate the hashes. [`out/06-check-unstamped.txt`](out/06-check-unstamped.txt) fails; [`out/06-stamp.txt`](out/06-stamp.txt) writes the two missing hashes with **0 model calls**; [`out/06-check-stamped.txt`](out/06-check-stamped.txt) then reports both records `current`, exit 0. Only missing values are written, so stamping never hides a real change.

### Deprecated records and `--prune-orphans` — `prov-2026-cc73cf61`, `prov-2026-fc05416c`

[`out/07-check-orphan.txt`](out/07-check-orphan.txt): after record b is deprecated, its region is an `orphan` and `--check` fails. [`out/07-prune-orphans.txt`](out/07-prune-orphans.txt): `pruned orphan … (record prov-2026-bbbbbbbb is gone or deprecated)`; [`out/07-markers.txt`](out/07-markers.txt) no longer has it. Orphans are pruned only in files a current record authorizes, and never under the generators or records directories.

### `sedum regions --json` and `grow --json` — `prov-2026-bf01cd6b`

[`out/08-regions-json.txt`](out/08-regions-json.txt): the `sedum.regions/1` document the provenance library's evidence code hash reads. [`out/09-grow-json.txt`](out/09-grow-json.txt): `sedum.grow/1`, where `"changed": true` marks exactly the records whose regions' bytes moved, so provenance runs only their specs. Both formats are documented with examples in `conformance/formats/`, and a test holds real output to them key for key.

### The `go` package and the syntax registry — `prov-2026-5403d162`, `prov-2026-0ce03d1c`, `prov-2026-73d32e29`, `prov-2026-7a704b59`

`testdata/generators/go` grows one-off Go through `setFunction`, `setMethod`, `setType`, `setVar`, `setConst` and `addImport`. A package opting into `syntax: go` has every touched file checked by Sedum's built-in Go checker (stdlib only, no external program): it must parse, be gofmt-clean, and use no package qualifier nothing in reach supplies. In grow a failure re-prompts the model with the exact fix; in replay it is terminal; Phase 7 checks again before writing.

[`out/10-go-rejected.txt`](out/10-go-rejected.txt): unformatted code is refused with the whole file in its gofmt form. [`out/10-go-accepted.txt`](out/10-go-accepted.txt), [`out/10-go-file.txt`](out/10-go-file.txt): the clean version lands, and `gofmt -l` reports nothing.

When two packages claim `.go`, the one whose file template matches a path most specifically wins it, so a command package keeps its `_cmd.go` files and `go` takes the rest.

### `quote` — `prov-2026-f63c340a`

`{{usage|quote}}` renders any value as a valid Go/JavaScript/TypeScript/Java/C# string literal, so a quote or backslash a model binds cannot break the build.

### gen_hash and records from the shared library — `prov-2026-dac49f78`, `prov-2026-5f1e4f0a`

[`out/11-parity.txt`](out/11-parity.txt): for every record in linespec, sedum, scion and linespec-provenance-ui, sedum's `gen_hash` equals the library's — 521 + 142 + 253 + 184 = **1,100** v1 records as they are, and 425 + 124 + 218 + 171 = **938** after the library migrates a copy to v2. Records are decoded by the library over the keys Sedum reads, so a malformed field Sedum does not read never fails a run.

### Boundaries and docs — `prov-2026-a88771dd`

`TOOL_BOUNDARIES.md` §9 records what moved into grow (per-record skip, re-render, re-select and reconciliation) and what stays above Sedum (repair loops, escalation, the two-beat lifecycle). `PRD.md` and `README.md` describe every flag and format above.

### Build

`go.mod` requires `github.com/livecodelife/provenance` by pseudo-version (`v0.0.0-20261016004733-69bef6526acd`); there is no `replace`, so CI and release build from this checkout alone (`prov-2026-3c35801c`).

---

## With a real model

Every capture above uses a stub model, so the call counts are exact. The same build has also grown real code with a real model: the `provenance` CLI's commands, grown by Claude Haiku 5.5 through this branch's sedum (via the herd's `sedum-haiku` wrapper), from one imprint per command, with the `cobra` and `go` generator packages. The gate for each is a hand-written `<cmd>_test.go`; all are green. Files are at [`provenance@658d54b/cmd/provenance`](https://github.com/livecodelife/provenance/tree/658d54b/cmd/provenance) (branch `cli`). Figures from prov-cli.

Calls are the first grow plus any regrows after a record revision; the reason for each revision is listed.

| Command | Record | Model calls | Record revision | Grown files |
|---|---|---|---|---|
| `provenance migrate` | prov-2026-76a125e6 | 2 | wiring only; its logic is a hand-written bridge until its v2 revise | [`main.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/main.go), [`migrate_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/migrate_cmd.go) |
| `provenance open` | prov-2026-6db88b75 | 2 |  | [`open_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/open_cmd.go), [`open_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/open_run.go) |
| `provenance create` | prov-2026-44a07768 | 2 + 1 | a missing import | [`create_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/create_cmd.go), [`create_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/create_run.go) |
| `provenance revise` | prov-2026-5f69977b | 1 + 1 | printed to os.Stdout instead of the command's writer | [`revise_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/revise_cmd.go), [`revise_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/revise_run.go) |
| `provenance edit` | prov-2026-57b8a95b | 1 |  | [`edit_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/edit_cmd.go), [`edit_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/edit_run.go) |
| `provenance status` | prov-2026-12b1ca97 | 2 + 2 | grouped field types misread | [`status_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/status_cmd.go), [`status_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/status_run.go) |
| `provenance lint` | prov-2026-a03e6bc3 | 2 |  | [`lint_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/lint_cmd.go), [`lint_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/lint_run.go) |
| `provenance run-specs` | prov-2026-c31c9c59 | 1 + 2 | an unused variable | [`run_specs_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/run_specs_cmd.go), [`run_specs_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/run_specs_run.go) |
| `provenance check --staged` | prov-2026-063f1395 | 1 + 1 | an "except" list misread | [`check_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/check_cmd.go), [`check_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/check_run.go) |
| `provenance hooks install` | prov-2026-d1d789de | 1 | hook scripts byte-identical to linespec's | [`hooks_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/hooks_cmd.go), [`hooks_install_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/hooks_install_cmd.go), [`hooks_install_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/hooks_install_run.go) |
| `provenance complete` | prov-2026-41fb548c | 1 + 1 | message wording | [`complete_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/complete_cmd.go), [`complete_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/complete_run.go) |
| `provenance deprecate` | prov-2026-35439261 | 1 |  | [`deprecate_cmd.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/deprecate_cmd.go), [`deprecate_run.go`](https://github.com/livecodelife/provenance/blob/658d54b/cmd/provenance/deprecate_run.go) |

**25 calls for 12 commands**, 8 of them after a record revision. 4 more regrew `open`, `create`, `revise` and `edit` when the cobra package's `setRun` began passing the command's writer: a package change, picked up by regrowing. A regrow of unchanged records makes **0** calls. **Every fix was a record revision, never a code edit.** Two of those fixes, the missing import and the unused variable, are what the Go checker now catches in the same call, or will once the queued unused-variable check lands. `common_run.go` and `migrate_run.go` are the CLI's two remaining hand-written bridge files, to be regrown once its records are on v2.

---

## Upgrading from v0.8.0

What a repository generated by sedum v0.8.0 sees on its first run of this version.

- **Old markers read as stale.** A v0.8.0 marker carries no `gen` or `pkg`, and an unknown hash is stale, never current. So the first `grow` re-selects every record that owns regions: one full model pass. To skip that, run `sedum grow --stamp` once. It writes the current hashes onto markers missing them, with no model call, and asserts the existing code is current. The lifecycle testbed shows both paths with v0.8.0's own binary: step `2′` (one full re-select) and step `4b` (`--stamp`, then 0 calls) in `evals/lifecycle/`.
- **Markers are byte-identical until rewritten.** Reading a v0.8.0 marker and writing it back unchanged produces the same bytes, because absent `gen` and `pkg` are omitted, not written empty (conformance case `round-trip/absent-gen-and-pkg-stay-absent`). Only a region that is actually regenerated or stamped gains the attributes.
- **A marker's record now means its owner.** In v0.8.0 it named the last record to write the region, and any later record could take the region over. Now a different record selecting it is an error naming both. A repository where one record's region was later rewritten by another - typically a v1 supersession - must migrate its records first (`provenance migrate`), which folds such chains into one record whose aliases own the old markers. The testbed's `step3grow` shows the refusal before migrating, and `step5` the clean grow after.
- **Recordings stay compatible.** Recordings gain optional `gen` and `created_at` per record. A v0.8.0 recording, which has neither, still replays: replay takes `gen` from `--records` when given, and otherwise writes none, which reads as stale. Older sedum versions read new recordings, because recording decoding is lenient. `grow --record` now re-selects every record, since a recording captures selections.
- **Behaviour changes worth knowing.** Records sharing a file now grow in one run (the duplicate-path error is gone); a deprecated record is no longer grown; `grow` reports orphan regions; and the completeness follow-up no longer asks about anchors another record already fills.

## Known limits

- **Changed kwarg defaults and `--var` values need `--force`.** A marker records every kwarg its region was rendered from, defaults and variables included, and re-rendering replays them verbatim. Neither `gen` nor `pkg` covers them, so a changed default or `--var` reaches existing regions only with `--force`. The `--check` and `--force` help says so (`prov-2026-9621daed`).
- **Re-selection can churn kwargs.** When a record changes, the model selects again and may bind kwargs differently on regions whose meaning did not change. Rebuilt prior selections and the "already present" list in the prompt reduce this (`prov-2026-cc73cf61`); `--dry-run` shows the diff first. They do not eliminate it.
- **A record with no regions costs a call on every grow.** There is nowhere to keep its hash, so it is selected each time and `--check` reports it as unknown rather than stale. This is a decided trade-off (DESIGN.md), narrow in practice: records that only create files, or that sedum cannot inject anything for, already make no call.
- **The Go checker is a parser, not a type checker.** It catches code that does not parse, is not gofmt-clean, or uses a package it does not import. It infers an import's name from its path, and errs toward silence when it cannot. It does not catch unused variables or unused imports yet; that AST-level check is queued.
- **One syntax checker.** The registry holds only `go`. Other languages get a checker by adding one to Sedum; packages cannot supply their own, by design.
- **Formatter-friendly markers are the template author's job.** A template that ends in a blank line keeps one before its closing marker, which is what gofmt wants after a top-level declaration (`prov-2026-ab24df5a`). Sedum does not know which targets need it.

---

## Records on this branch

| Record | Type | Title |
|---|---|---|
| prov-2026-66bea9f3 | blueprint | Markers carry gen and pkg hashes as first-class attributes |
| prov-2026-b9298192 | imprint | Recordings carry each record's gen so replay reproduces grow |
| prov-2026-dac49f78 | blueprint | gen_hash comes from the shared provenance record library |
| prov-2026-375f5d4a | blueprint | Regions belong to the record named on their marker |
| prov-2026-43b7ec39 | imprint | Recordings carry each record's created_at so replay orders shared anchors like grow |
| prov-2026-4df89974 | imprint | A region's record is its owner, not its last writer |
| prov-2026-9621daed | blueprint | Grow skips, re-renders or re-selects each record by its hashes |
| prov-2026-cc73cf61 | blueprint | Orphans, deprecated records, stamping, and prompt context for grow |
| prov-2026-53f4f82b | imprint | Completeness ignores anchors other records already fill |
| prov-2026-14fa88c6 | imprint | The chi fixture plants its imports anchor inside the import block |
| prov-2026-ab24df5a | imprint | A template ending in a blank line keeps one before its closing marker |
| prov-2026-f63c340a | blueprint | A quote transform renders a value as a string literal |
| prov-2026-5403d162 | blueprint | A generic go generator package grows one-off Go code through literal-code actions |
| prov-2026-fc05416c | imprint | Grow's marker scans skip the generators and records directories |
| prov-2026-73d32e29 | imprint | A package can declare anchors completeness never asks about |
| prov-2026-0ce03d1c | imprint | The go checker flags package qualifiers the file does not import |
| prov-2026-a88771dd | blueprint | Per-record regeneration moves into grow: TOOL_BOUNDARIES and the record attribute catch up |
| prov-2026-bf01cd6b | blueprint | sedum regions and a JSON grow report feed provenance evidence |
| prov-2026-7a704b59 | imprint | The go package leaves go.mod and go.sum to go mod |
| prov-2026-5f1e4f0a | blueprint | sedum reads records through the shared provenance record library |
| prov-2026-3c35801c | blueprint | Prepare grow-prep for its PR: standalone build, a proof document, v2 records |
