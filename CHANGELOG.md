# Changelog

## 0.10.1

- `sedum check` recognises the REVISE001 refusal of provenance 3.0.0, which reports "is implemented: run 'provenance revise'" under that rule id instead of STATE001, so a pure regeneration or marker migration of an implemented record is waived again. The STATE001 form of earlier provenance is still recognised.

## 0.10.0

- **`sedum check`** runs `provenance check` on each commit of `--range` (or on the staged change) and waives only the findings it can prove come from a pure regeneration or a pure marker migration, so the checks that used to live in provenance now live here. Exit 0 means no finding is left, 1 that some are, 2 that provenance or git could not be run or read.
- **`sedum init`** creates `.sedum/context.md`, makes `.gitignore` ignore `.sedum/run.log`, and adds `generators/**` and `.sedum/**` to the `exclude_paths` of `.provenance.yml`; it is safe to run again.
- **`sedum hooks install`** writes sedum's own `commit-msg` and `pre-push` hooks (they run `sedum check`; the pre-push hook runs provenance's own pre-push hook first). A different hook already there is kept as `<name>.bak`.
- **`sedum evidence`** judges each record's code freshness from the regions and the record's committed evidence file, `.sedum/evidence/<record>.jsonl`; `--refresh` appends an entry for each stale or untracked record. Run it where `provenance evidence refresh --all` was used.
- `sedum grow --check` skips, and counts in its report, the scope paths no loaded package claims instead of failing on them.
- A provenance release that no longer knows about Sedum follows this one.

## 0.9.2

- `sedum markers migrate` replaces each file by renaming a staged replacement over it, with the file's own mode, so a file locked against writing (as `provenance reconcile` leaves the files of sealed records) is migrated and stays locked, and a failure while staging leaves the tree as it was: no file rewritten, no recording written, no temporary file left. A rename that fails after staging can still leave files already renamed in place.

## 0.9.1

- `sedum markers migrate` skips, and counts in its report, the scope paths no loaded package claims instead of failing on them, and judges a file's regions against the package that declares their actions.

## 0.9.0

- **Marker v2.** Every region is written with a marker that carries `"v":2`, an `id` object of its identifying kwargs, and `kwargs` and `body` hashes in place of the kwargs object. Each `grow` that writes also writes `.sedum/recordings/<record-id>.json` (or the directory `--record` names), which holds the kwargs a re-render needs; a missing recording costs a re-select, never a failure. `grow --check` reports a region whose text no longer matches its `body` as `edited` and fails; a re-render over a hand edit is refused, and `--force` overwrites it. v1 regions are listed as `unverifiable` and keep re-rendering from their own kwargs. `--record <file>.json` keeps re-selecting every record. `sedum markers migrate` rewrites an existing tree with no model. A v2 marker cannot be read by an older Sedum, so upgrade before sharing a tree. New corpora: `conformance/markers/v2cases.json` and `conformance/kwargshash/cases.json`.
- The marker walk (`grow`'s orphan scan, `sedum regions`) skips `node_modules`, `.next`, `vendor`, `dist` and `target` at any depth, so a checkout with large vendored trees is no longer read file by file. Add more with `exclude:` in `.sedum/config.yml` or `SEDUM_EXCLUDE`.
- A path pattern segment written exactly as the path's segment now matches it literally, so a path with bracket characters (a Next.js `[id]` or `[...slug]` route folder) is covered by the scope entry that names it; other bracket segments are still read as character classes.
- Two built-in transforms: `trimprefix:X` (strip a leading literal and the blanks after it) and `unparen` (strip one matching outer pair of parentheses).

## 0.8.0

The first release since 0.7.0 that understands provenance v2 records and can
bring an existing codebase under Sedum. `sedum --version` prints `0.8.0`, the
same value a recording carries in `sedum_version`.

### New

- **`sedum adopt`** brings existing files under Sedum without a model
  (design: `docs/ADOPT.md`). A generator package declares an adopter in its
  `sedum.yaml`: declarative (`unit_start`, `unit_end`, `header`, `actions`), an
  external `command:` (run only with `--allow-command-adopters`), or
  `builtin: go`, the Go adopter compiled into sedum (`go/parser`; gofmt-clean
  input, one unit per declaration, grouped imports via the new `setImports`
  action). Files it cannot adopt are skipped and named with a reason and line
  range; a file is kept only if stripping the markers gives back the original.
  The report prints the `provenance edit <imprint> --add-scope ...` lines that
  make `grow --check` see the adopted files. Limit in this release: root-level
  Go files must be `package main`.
- **Hash-aware grow.** Markers carry `gen` (a hash of the record) and `pkg` (a
  hash of the action and its templates); `grow` skips a current record,
  re-renders one whose template changed, and re-selects one whose record
  changed, per record, with `--force`, `--check`, `--prune-orphans`, and
  region ownership by record. Provenance v2 records are read through the
  provenance library (`github.com/livecodelife/provenance`).
- **`grow --repair N`** sends a selection back, with the output, when a project
  verification command (`.sedum/config.yml` or `--verify`) fails it.
- **`sedum generators install`** and a built-in `go` generator package, so a Go
  project needs no generator directory to start.
- **The grow log.** Every `sedum grow` that writes appends one line per record
  to `.sedum/grow.jsonl` (model, calls, retries, repairs, tokens, literal vs
  declarative invocations, region counts, outcome); `sedum grow-log` summarizes
  it per model. `SEDUM_GROW_LOG=<path>|off` and `SEDUM_GROW_LOG_SKIPS=1`
  control it; keep the file out of git (`.sedum/*`).
- `replay --restamp-gen`, `grow --execute --check`, an orphan-region summary
  (`--orphans` for the full list), and a clear error for a file no template
  matches.
- **Evals.** `evals/lifecycle` (the runnable hash-aware-grow demo) and
  `evals/design` (how much design a cheap model does from a behaviour-only
  intent; first result in `evals/design/README.md`). Neither is part of the
  binary.

### Changed

- CI gates on `provenance lint` and on every commit satisfying its records.
- Generated provenance files are collapsed in diffs (`.gitattributes`).

### Upgrading

Install the release with `install.sh` (it verifies `checksums.txt`) or
`go install github.com/livecodelife/sedum/cmd/sedum@v0.8.0`. Markers written before
hash-aware grow carry no `gen`/`pkg`; one `sedum grow --stamp` makes them current
with no model call.
