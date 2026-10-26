# Marker v2 and per-record recordings

Status: draft for review. Decided by Caleb (2026-10-20); this document turns the
decision into a spec. Nothing here is implemented yet.

## Why

A marker today carries `kwargs`, every kwarg the region was rendered from,
including `code`, the whole text of the region a second time. That makes
markers large (on provisio-backend-ts, 477 markers are 300 kB of marker lines),
makes the generated code unreadable in a diff, and still cannot tell whether
the region was hand-edited: nothing hashes the region's body, so `grow --check`
says "current" for an edited region and the next re-render overwrites the edit
without a word.

Marker v2 keeps in the file only what is needed to *identify* a region and to
*detect* drift, and moves what is needed to *rebuild* a region into a
recording under `.sedum/`.

## What changes

1. Markers carry hashes of the kwargs and of the region text, not the kwargs.
2. Every grow writes one recording per record under `.sedum/recordings/`.
3. A re-render reads its kwargs from that recording.
4. `grow --check` reports a region whose text no longer matches its `body` hash.
5. A one-off, model-free migration rewrites an existing tree.

Old markers (no `v`) keep reading exactly as today.

## Marker v2

```
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-ab12cd34","gen":"sha256:…","pkg":"sha256:…","id":{"file":"internal/x/y.go","name":"parse"},"kwargs":"sha256:…","body":"sha256:…"}
```

The action (and variant) stay literal on the line, as today. The attribute
object is still one JSON object; unknown keys are still carried, never read.

| Key | Meaning |
|---|---|
| `v` | Marker version. `2` for this format. Absent means v1 (today's). A reader meeting a `v` it does not know treats the region as foreign: it does not rewrite it, and reports it. |
| `tier`, `record`, `gen`, `pkg`, `writer` | Unchanged. |
| `id` | The identity kwargs: exactly the kwargs `selectingKwargs(action)` names, with their resolved values (the action's declared `identity`, else its required kwargs). This is what makes region identity, orphan detection and `sedum regions` work without the recording. It is an object, so a tool reading markers by `grep` can still see file, name and receiver. |
| `kwargs` | `sha256:` hash of the canonical form (below) of **all** resolved kwargs the region was rendered from. A string in v2; an object in v1. A reader tells the versions apart by `v`, and tolerates the type mismatch by type. |
| `body` | `sha256:` hash of the region text (below). |

No kwarg name is special. `code` is hashed with every other kwarg and is
recovered from the recording like every other kwarg.

A v2 marker never carries the v1 `kwargs` object, so `sedum` before this change
fails to read a v2 marker (a string where it expects an object). This is a
format break in one direction only; see open question 1.

### Canonical kwargs

The input is the resolved kwargs object, the value `Marker.Kwargs` holds today:
every kwarg of the action that has a value after variable substitution and
defaults, nothing else.

- **Absent vs defaulted.** A kwarg the action declares with a default appears
  with that default when the invocation omits it. A kwarg with no default that
  the invocation omits is absent, not `null`. An explicit `null` is present.
  So the hash does not change when a record spells a default out, and does
  change when a package changes a default (which is how a changed default
  reaches a region only through `--force`, as documented today).
- **Encoding.** JSON, UTF-8, no insignificant whitespace, no trailing newline.
  Object keys sorted by byte order of their UTF-8 form, at every depth. Arrays
  keep their order. Strings escape only `"`, `\`, and control characters
  below U+0020 (`\b \f \n \r \t` short forms, otherwise `\u00xx` lower-case
  hex); `<`, `>`, `&`, U+2028 and U+2029 are written literally. No Unicode
  normalization.
- **Numbers.** A number whose value is an integer is written as that integer in
  decimal with no sign for zero, no leading zeros, no fraction and no
  exponent, so `1`, `1.0` and `1e0` all hash as `1`, and `-0` as `0`. Any other
  number is written in the ECMAScript `Number::toString` form of its nearest
  double. `NaN`, infinities, and integers beyond 2^53 are refused: the region
  cannot be recorded (an error, not a silent rounding). A tool reading a
  recording must parse numbers as arbitrary-precision decimals before this
  step; `json.Number` does that today.
- **Hash.** `sha256:` followed by 64 lower-case hex digits of SHA-256 over the
  canonical bytes.

### Body

The body is the text between the marker lines:

- It starts at the first byte after the newline that ends the opening marker
  line and ends after the newline that ends the last line before the closing
  marker line. The marker lines themselves are excluded (they contain the
  hash). The closing marker's own indentation is excluded.
- Bytes are hashed as in the file, indentation included: a re-indent is an
  edit. Line endings are normalized first: every `\r\n` becomes `\n`, so a
  checkout with autocrlf does not read as a hand edit. A body with no lines
  (an empty region) hashes the empty string.
- Trailing blank lines and trailing spaces inside the region count. The blank
  line Sedum writes *after* the closing marker is outside the region.

## Recordings

`sedum grow` writes `.sedum/recordings/<record-id>.json` for every record whose
regions it handled, as a side effect of a normal run:

- **Re-selected record:** written from the selection and what it resolved.
- **Skipped record (everything current):** the file is left as it is, but is
  verified (below); if it is missing or does not match, the record is stale.
- **Re-rendered record (only `pkg` moved):** the regions are re-rendered from
  the recording, and the file is rewritten with the new region hashes.
- **`--dry-run`, `--check`:** nothing is written.

`--record <dir>` relocates the directory (default `.sedum/recordings`). A
`--record` path ending in `.json` keeps the old behaviour, which forces a
re-select of every record; any other path is a directory and forces nothing. The files are
meant to be committed; a team that does not want them gitignores `.sedum/`
(`.sedum/` is already in `.provenance.yml` `exclude_paths` here).

### File format

A per-record file is a valid `recording.Recording` holding exactly one
`Record`, so `sedum grow --execute` and every existing reader accept it
unchanged. One additive field is new on `Record`:

```json
{
  "sedum_version": "0.9.0",
  "packages": { … },
  "variables": { … },
  "records": [{
    "record_id": "prov-2026-ab12cd34",
    "gen": "sha256:…",
    "created_at": "2026-10-16",
    "files": [ … ],
    "phases": [ { "name": "default", "invocations": [ … ] } ],
    "regions": [
      { "path": "internal/x/y.go", "action": "setFunction", "variant": "",
        "tier": "owned", "pkg": "sha256:…",
        "id": { "file": "internal/x/y.go", "name": "parse" },
        "kwargs": { "code": "…", "file": "…", "name": "parse" },
        "kwargs_hash": "sha256:…", "body": "sha256:…" }
    ]
  }]
}
```

`phases` stays what the model selected (pre-expansion, as today): it is what a
re-select shows the model as its prior selection and what `--execute` replays.
`regions` is new: one entry per region after expansion, with the resolved
kwargs. It is what a re-render reads, and it is sorted by (path, action,
variant, canonical `id`) so the file diffs cleanly. The two overlap on `code`;
see open question 3.

The file carries no timestamp, run id, or other volatile field (unchanged).

### Verifying a region against its recording

A v2 region is **current** when all hold:

1. `gen` equals the record's current gen, and `pkg` equals its action's pkg hash
   (as today);
2. the recording file exists, parses, and has a `regions` entry with the
   region's identity;
3. that entry's `kwargs_hash` equals the marker's `kwargs`, and the hash of its
   `kwargs` object equals `kwargs_hash`.

Otherwise the region is **stale** (the record is re-selected), with a reason:
`no recording`, `recording unreadable`, `recording has no entry for this region`,
`recording does not match the marker`. A missing or unreadable recording is
never an error by itself: it costs a re-select, not a failure.

### Re-render

When only `pkg` moved, the region is re-rendered from the recording entry's
`kwargs` instead of the marker's, then rewritten with a new `pkg`, `kwargs`
(unchanged unless the action's kwarg set grew a default) and `body`.

**Hand edits.** If the region's current text does not hash to the marker's
`body`, a re-render refuses to overwrite it: the run reports the region and
fails, and `--force` overwrites. (A re-select overwrites, as it always has.)

## `grow --check`

`--check` reports a fourth class beside current, stale, and re-render:

```
edited     prov-2026-ab12cd34
             internal/x/y.go  setFunction  (hand-edited: body hash differs)
```

and exits non-zero when any region is edited, stale, or needs re-render. v1
regions have no `body`, so they cannot be reported as edited; they are reported
`unverifiable` (not a failure) until migrated.

## Migration

`sedum markers migrate --generators <dir> --records <dir> [--output <dir>]
[--dry-run]`, no model. For each file holding v1 markers:

1. read every v1 marker (its `kwargs` object is the source of truth);
2. build the record's recording from the markers: `regions` from the kwargs;
   `phases` rebuilt as today's prior-selection logic does (pipeline.go: markers
   to invocations), `files` from resolution, `gen` from the marker (not the
   current record: a stale marker stays stale);
3. rewrite each marker to v2, with `body` computed over the file's present
   text, and `pkg`/`gen` kept as they were;
4. write the recordings.

It refuses a region whose canonical form cannot be computed (open numbers) and
reports it; it never changes a region's body. **Acceptance:** after migration,
`sedum grow --check` reports the same current/stale classes as before, and
re-rendering any migrated region from its recording produces text equal to its
`body` (zero regions changed), on provisio-backend-ts and on this repository.
A second migration is a no-op.

## Conformance

`conformance/markers/cases.json` gains v2 cases: parse, emit and round-trip of a
v2 marker; v1 marker still parsing; unknown `v`; the canonical-kwargs cases
(`1` vs `1.0`, key order at depth, escaping, absent vs defaulted, the refusal
cases); and body-hash cases (CRLF, indentation, empty region, trailing blank
lines). A new `conformance/kwargshash/cases.json` holds input kwargs and the
expected hash, so a foreign writer can compute `kwargs` and `body` without
linking Sedum, as it can `pkg` today.

## Size

Measured by rewriting the existing markers' attribute objects to the v2 shape
(estimate; the migration will replace these with measured values):

| Tree | Markers | v1 marker bytes | v2 (full `sha256:` hashes) | v2 (16 hex digit `kwargs`/`body`) |
|---|---|---|---|---|
| provisio-backend-ts | 477 | 300,343 | 222,636 (-26%) | ≈170,000 (-43%) |
| sedum | 868 | 538,902 | 399,766 (-26%) | ≈304,000 (-44%) |

Most of what stays is the two existing 71-byte hashes (`gen`, `pkg`) and
`record`. See open question 2.

## Implementation units

New or revised, each an imprint under one blueprint, grown (no handwriting
except tests where a fixture is exact):

| Unit | Files |
|---|---|
| Canonical kwargs and body hashing | `internal/inject/canon.go` (new) |
| Marker v2 read/write, identity, unknown `v` | `internal/inject/marker.go` |
| Hash the body on inject; verify on scan | `internal/inject/inject.go`, `scan.go` |
| Recording gains `regions`; per-record read/write | `internal/recording/…` |
| Write recordings from grow; `--record <dir>` | `internal/pipeline/…`, `internal/cli/grow.go` |
| Classify with recordings; edited class; re-render from recording | `internal/pipeline/classify.go` |
| `--check` report | `internal/cli/report.go` |
| `sedum markers migrate` | `internal/cli/markers.go` (new), `internal/pipeline/migrate.go` (new) |
| Conformance cases and corpus | `conformance/` |

## Open questions

1. **v2 markers are unreadable to older sedum.** `kwargs` changes type. Keep the
   name `kwargs` (as decided) and accept the break, or use another key
   (`kwargs_hash`) so an older sedum reads the marker and ignores the hash?
   An older reader would then read `id` as an unknown key and find no `kwargs`.
2. **Hash length.** Full `sha256:` + 64 hex (consistent with `gen` and `pkg`)
   or 16 hex digits for the two new ones? The second saves another ~17 points
   of marker size; the collision risk for a drift detector is negligible at 64
   bits.
3. **Recording duplicates `code`.** `phases` and `regions` both carry it.
   Options: keep both (simple; the file is about twice the generated code);
   store `regions` only and rebuild `phases` for `--execute`; or store `phases`
   only and re-expand to get `regions`, which loses the "changed default does
   not reach the region" behaviour.
4. **Re-render over a hand edit.** Refuse (spec'd), overwrite with a warning, or
   write the new text beside the old?
5. **`--record <file>`.** Decided: a path ending in `.json` keeps the old
   force-re-select behaviour, so scripts passing a file path do not break; any
   other path is a recordings directory.
6. **Gitignored recordings.** With `.sedum/` ignored on a clean checkout every
   v2 region is stale (`no recording`) and the first grow re-selects everything
   with a model. A fallback that rebuilds the recording from the code is
   impossible once markers drop `code`. Is that acceptable, or should the first
   cut keep `code` in markers when no recording exists?
7. **Unknown `v`.** Foreign-and-reported (spec'd) vs error.
8. **v1 regions in `--check`.** `unverifiable` is not a failure until a cutoff
   version; set a date after which it fails?
9. **Float kwargs.** The refusal of integers beyond 2^53 and non-finite numbers
   may be stricter than needed; no kwarg in provisio-backend or sedum uses one.
