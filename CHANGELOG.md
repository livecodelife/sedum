# Changelog

## Unreleased

- The marker walk (`grow`'s orphan scan, `sedum regions`) skips `node_modules`, `.next`, `vendor`, `dist` and `target` at any depth, so a checkout with large vendored trees is no longer read file by file. Add more with `exclude:` in `.sedum/config.yml` or `SEDUM_EXCLUDE`.
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
