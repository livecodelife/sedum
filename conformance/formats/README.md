# Sedum's JSON report formats

Three commands print JSON for tools that must not link Sedum - the provenance
library and its TypeScript backend among them (`prov-2026-bf01cd6b`). Each
document names its format and version in `format`. A version changes only when
a field changes meaning; fields are added without one, so a reader must ignore
keys it does not know.

Sedum's own tests decode real output and check every key against the examples
here, so a key documented here is a key Sedum emits.

## `sedum.regions/1` — `sedum regions --json`

Every region whose marker names a record, read from markers under `--output`.
Read-only: no model, nothing written. Hidden directories, the generators
directory and the records directory are never read.

| Key | Meaning |
|---|---|
| `format` | `"sedum.regions/1"` |
| `regions` | Array, sorted by `file`, then by position in the file. |
| `regions[].record` | The owning record. With `--records`, a marker naming an alias is reported under the surviving id. |
| `regions[].file` | Slash-separated path relative to `--output`. |
| `regions[].action`, `regions[].variant` | As on the marker; `variant` is `""` when the action has none. |
| `regions[].tier` | `"owned"` or `"seeded"`. Only owned regions are regenerated. |
| `regions[].identity` | The kwargs that identify the region: the action's declared identity, or its required kwargs; every recorded kwarg when the package no longer declares the action. |
| `regions[].gen` | The record's `gen_hash` when the region was written, or `""` when the marker has none. |
| `regions[].pkg` | The action's `pkg` hash when the region was rendered, or `""`. |

An empty `gen` or `pkg` means unknown. A consumer must treat that code as
unknown, never as current; `sedum grow --stamp` fills them in.

## `sedum.grow/1` — `sedum grow --json`

The run's report, in place of the text report.

| Key | Meaning |
|---|---|
| `format` | `"sedum.grow/1"` |
| `dry_run` | Whether anything was written. |
| `records` | Array, one per record the run considered, sorted by id. |
| `records[].record` | The record id. |
| `records[].mode` | `"select"`, `"re-select"`, `"skip"`, `"re-render"`, `"none"` (nothing Sedum writes), or `"replay"` under `--execute`. |
| `records[].changed` | True exactly when one of the record's regions was inserted, pruned, or rewritten with different bytes. |
| `records[].repair_attempts` | How many times the project's verification sent the record's selection back to the model (`--repair`); omitted when none. |
| `records[].regions` | What happened to each region the run touched for the record. |
| `records[].regions[].file`, `.action`, `.variant` | Which region. |
| `records[].regions[].change` | `"inserted"`, `"replaced"` (bytes differ), `"unchanged"` (rewritten identically), `"pruned"`, `"held"` (a prune only `--force` would make), `"unselected"` (a partial recording's replay left a region of a named record it no longer selects), `"seeded"` (left alone). |
| `orphans` | Regions whose record is gone or deprecated: `file`, `action`, `variant`, `record`, `authorized`, `pruned`. |
| `stamped` | Markers `--stamp` wrote a hash onto: `file`, `action`, `variant`, `record`, `gen`, `pkg` (booleans: which hash was added). |
| `check_failed` | Under `--check`: some record is stale, needs re-rendering, or has left orphans. Always false otherwise. |

`changed` is what lets a caller run exactly the specs whose code moved.

## `sedum.evidence/1` — `sedum evidence --json`

Whether each record's code is the code its latest evidence entry saw. The
evidence files are `.sedum/evidence/<record>.jsonl` under `--output`, one entry
per line, appended by `sedum evidence --refresh` and never rewritten. Read-only
without `--refresh`.

| Key | Meaning |
|---|---|
| `format` | `"sedum.evidence/1"` |
| `records` | Array, one per non-deprecated record in the records directory, sorted by `record`; `[]` when there are none. |
| `records[].record` | The record id. |
| `records[].version` | The record's current version, read from its record file. |
| `records[].verdict` | `"fresh"` (the latest entry for this version has the current code hash), `"stale"` (it has another), `"none"` (no entry covers this version), `"unknown"` (an owned region has no `gen` or `pkg`), or `"no regions"` (the record owns no region, so there is nothing to track). |
| `records[].code_hash` | The code hash of the owned regions now (`docs/HASHING.md` section 7 of the provenance library), or `""` when the verdict is `unknown` or `no regions`. |
| `records[].recorded` | The code hash of the latest entry for this version, or `""` unless the verdict is `fresh` or `stale`. |
| `records[].regions` | How many owned regions the record has. |

A reader ignores keys it does not know. `unknown` must never be read as fresh.
