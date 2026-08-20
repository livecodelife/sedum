# Sedum conformance corpus

Golden cases for the file formats Sedum shares with other tools, so that a tool
Sedum's authors did not write can check that it has complied rather than
discovering the rules one breakage at a time.

Today it covers one format: the **ownership marker**.

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
