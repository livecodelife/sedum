# Sedum conformance corpus

Golden cases for the file formats Sedum shares with other tools, so that a tool
Sedum's authors did not write can check that it has complied rather than
discovering the rules one breakage at a time.

It covers the **ownership marker** (v1 and v2), the **`pkg`**, **`kwargs`** and **`body`** hashes a marker carries, and the two **JSON reports** Sedum prints for other tools (`formats/`: `sedum.regions/1` and `sedum.grow/1`, `prov-2026-bf01cd6b`).

## Why the marker and nothing else yet

A tool built above Sedum reads `sedum.yaml` and `actions.yaml`, so a misreading
is its own problem and fails on its own side. It writes recordings, and replay
validates those terminally at ingestion, so a malformed one is refused before
anything is written.

The marker is the only format a foreign tool writes into files that Sedum will
later re-read and rewrite, with no validation gate in between. A writer that
gets it wrong corrupts state silently, and the first breakage looks like a bug
in Sedum.

## What is here

`markers/cases.json` — one array of cases. Every case carries:

| Field | Meaning |
|---|---|
| `id` | Stable identifier. Cite it when reporting a divergence. |
| `kind` | `parse`, `emit`, `round_trip`, `replace`, `close`, or `error`. |
| `decision` | The provenance record that decided this behavior. |
| `why` | What goes wrong for an implementation that gets it wrong. |
| `comment_prefix` | The package's declared comment prefix. Never inferred from a file extension. |
| `comment_suffix` | The package's declared closing delimiter, for a target whose comments do not run to end of line. **Optional, and absent in most cases.** Absent means the marker carries no closing delimiter and no trailing space — byte-identical to what Sedum wrote before the field existed. When present it is separated from what precedes it by one space, on the opening line and the closing line alike, and a reader must strip it before interpreting what remains: the attribute object on an opening marker is not readable as JSON with the delimiter still attached. |

and, depending on `kind`:

- **`parse`** — `input` (one line) and `parsed`. Structural comparison. A
  `parsed.recognized` of `false` means the line is not a marker, which is not an
  error: most lines in a generated file are not markers.
- **`emit`** — `marker` and `output`. **Byte comparison.**
- **`round_trip`** — `input`, `parsed`, and `output`. Parse, then emit, and
  compare bytes. Emitting the parse of `output` must produce `output` again.
- **`replace`** — `existing` (the marker already in the file), `marker` (the one
  being written over it), and `output`. **Byte comparison.** The rule: take
  every modelled field from `marker`, and take the unmodelled keys from
  `existing`.
- **`close`** — `marker` and `output`, for the region's closing line.
- **`error`** — `input`, `error: true`, and `error_names`, a substring the
  failure must contain.

`parsed.extra` and `marker.extra` hold the attribute keys the format does not
model. Compare them as JSON values, not as text.

## `markers/v2cases.json`

The v2 ownership marker (`"v":2`: an `id` object, and `kwargs` and `body` as
`sha256:` hash strings) in the case shapes of `markers/cases.json`, held in the
`cases` array of an object that also carries `format`, `corpus_version` and a
`description`. Its kinds are `parse`, `emit`, `round_trip`, `replace` and
`error`. **The cases in `v2cases.json` carry no `decision` field**; every other
field of the table above applies. In a case's `marker` and `parsed` objects,
`kwargs` holds the `id` object, `kwargs_hash` the hash string, `body` the body
hash, `version` is `2` (or `0` for a v1 or foreign marker), and `foreign` says
whether the marker is another format's.

## `kwargshash/cases.json`

The two hashes a v2 marker carries, so a foreign writer can compute them
without linking Sedum, as it can `pkg` (`pkghash/cases.json`). Its cases are in
the `cases` array of an object, and each has an `id`, a `kind` (`kwargs` or
`body`), a `why` and:

- **`kwargs`** — `input` (the resolved kwargs as JSON text) and either
  `canonical` (the exact canonical text hashed) and `hash`, or `error: true`
  and `error_names`, a substring the refusal must contain. Keys are sorted at
  every depth, there is no insignificant whitespace, and an integer-valued
  number has one spelling.
- **`body`** — `input` (the region text between the marker lines) and `hash`.
  Line endings are normalized to `\n` before hashing.

Hashes are `sha256:` followed by 64 lowercase hex digits.

## `pkghash/cases.json`

The `pkg` attribute on a marker is the hash of everything that decides what an
action renders for fixed kwargs (`prov-2026-66bea9f3`). A tool deciding whether
a region was rendered from an older package needs to compute it, and must not
have to link Sedum to do so. Each case carries:

| Field | Meaning |
|---|---|
| `files` | A complete generator package, keyed by path under the generators directory. |
| `action` | The action whose hash the case pins. |
| `canonical` | The exact bytes hashed. |
| `pkg` | `sha256:` and the lowercase hex SHA-256 of `canonical`. |

`canonical` is one JSON object with byte-sorted keys and no insignificant
whitespace, escaping only `"`, `\`, and U+0000–U+001F, which are always
written `\u00xx` in lowercase hex — the same rules as the shared `gen_hash`
encoding, so a reader learns one encoding. Its keys:

- `format`: `"sedum.pkg/1"`
- `action`: `name`; `kwargs` as `{name: {type, required, default}}` with an
  undeclared default as `null`; `discriminator`; `variants` in declared order;
  `injects_into`; `anchor`, `anchor_start`, `anchor_end`, `anchor_pattern`;
  `identity`. Absent strings are `""` and absent lists `[]`.
- `templates`: each template's content keyed by variant name, `_default`
  included, or by `""` for a simple action
- `transforms` and `op_exceptions`: the package's, as declared

Descriptions are left out, of the action and of each kwarg: they reach the
model, never the renderer. The expected values were computed by an
independent Python implementation of the description above, not by Sedum.

## Running it

There is no runner in this repository, deliberately: a runner in one language is
a runner for one implementation, and the audience for this corpus is a writer
who is not using Go. Load the JSON, drive your own codec, compare.

Sedum's own tests read these files rather than carrying a parallel copy of the
same expectations. That is the property that keeps the corpus honest — a corpus
maintained beside the tests it duplicates goes stale within a release and then
misleads, because a stranger checking against it has no way to know. What is
published here is what fails when Sedum regresses.

## What is not mechanized here

Two of the rules a foreign writer must follow are file-level and cannot be
expressed as a parse-and-re-emit golden. They are still rules, and the corpus's
existence does not close them:

- **Do not remove or relocate a marker a file template planted.** Anchors exist
  because a template created them; moving one breaks the injection that targets
  it.
- **Do not reorder regions within a file.** Repeated injections at one anchor
  accumulate in invocation order, and reordering them makes a rerun's output
  differ from its input for reasons Sedum cannot see.

Transform behavior — rendering an action's `injects_into` through pipelines,
exception tables, and the inflection table — is deliberately absent too. It is
not corpus-shaped: pinning it would mean enumerating every irregular and growing
the corpus whenever the inflection table does. It is exposed as behavior
instead, through `sedum render`, and a caller is expected to invoke that rather
than reimplement it.

## Adding a case

A case is added when a format decision is made, and it cites that decision.
Coverage is not a goal in itself: the cases worth pinning are the ones that were
subtle enough to need working out, because those are what an independent
implementation gets wrong.

Adding a case never changes a marker already on disk. These fixtures record
current behavior; a case that would require rewriting existing markers is a
format change and needs its own record.
