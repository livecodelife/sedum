# Sedum — Tool Boundaries

Decision doc. Records where Sedum ends and adjacent tools begin, why the boundary
falls there, and what must be stable before M6/M7 so that a harness — ours or a
third party's — can be built on top without changing Sedum.

It does not supersede the PRD. It relocates §4–§14 of *Open Design Questions* out
of Sedum and into a sibling tool, and lists the PRD edits that relocation requires.

**Distribution intent:** Sedum open source, our harness closed. Third-party
harnesses supported. This is the constraint that drives §4–§6 — a private harness
could co-evolve with Sedum, a public one cannot.

**Revision note:** §5 was rewritten after reading `internal/inject/marker.go` and
`internal/expand/expand.go` at M5. Three of the five items in the prior revision
were already settled in code, and one of them (a marker version token) was
misjudged. What survives is smaller and sharper.

**Landed under `prov-2026-72775ae5` (Record A) and `prov-2026-59535684`
(Record B).** Everything in §7 is done, and implementing it corrected three
things this document got wrong; those are marked ⚠︎ inline below. §8's writer-key
and version-token questions are answered. What remains open is listed in §8.

---

## 1. The decision

Three tools.

**LineSpec** — behavioral contracts, provenance records in both domains, record
drafting and discovery, test execution. Exists today.

**Sedum** — deterministic generation from provenance records and generator
packages. The PRD as written, M1–M7. Unchanged in scope.

**The convergence loop** — unnamed. Owns the two-beat lifecycle (ODQ §13),
topological ordering (§7), the focused edit layer (§4), the paired update loop
(§5), failure traceback (§9), the two repair modes (§10), and the escalation
ladder (§11). Ours is one instance of this, not the only possible one.

ODQ §12's record *drafter* belongs in LineSpec, which already owns record
authoring; `provenance create` / `discover` / `next` are the same shape of
operation, and drafting is `discover` pointed at intent rather than at source.

---

## 2. Why the boundary falls here

**Not determinism.** That line cuts wrong. Sedum's Phase 4 is a model call, and
purifying Sedum into a zero-model tool would gut M6 for no benefit.

The line is **bounded selection over a closed, declared vocabulary** versus
**open synthesis validated by execution**. Sedum's model emits `{action, kwargs}`
drawn from a catalog; every failure mode is machine-checkable before anything
runs, and Phase 5's retry costs one model call. A harness's model writes freeform
content and can only be validated by running something. A single binary containing
both takes the weaker guarantee as its own.

Three supporting arguments.

**Target knowledge.** "Sedum's core contains no language-specific knowledge" is
the PRD's load-bearing principle. A harness is irreducibly target-aware: it boots
the service, runs the suite, isolates per-region assertions. That knowledge already
has a home — LineSpec's `.linespec.yml`. Putting the harness in Sedum either drags
that into core or pushes `run_command` shell into generator packages, which ODQ §2
flags as the user-defined-transforms trust problem.

**Runtime shape.** ODQ §4 makes warm-harness design the precondition for the edit
layer existing at all. A warm harness is a long-lived process holding a resident
model and a running container. Sedum is a batch CLI that exits.

**Structural rather than disciplinary guarantee.** ODQ §8 puts integration specs
out-of-tree so the model *cannot* reach them. The split is that move applied to
write authority: Sedum structurally cannot perform a freeform edit.

---

## 3. Direction of control

The harness sits **above** Sedum, not after it.

```
harness
  ├─ beat one: structural completion  → calls Sedum
  ├─ unmanaged-path satisfaction      → human or harness   (§6, finding 4)
  ├─ beat two: behavioral convergence
  │     ├─ in-parameter-space update  → calls Sedum (re-invocation)
  │     ├─ out-of-parameter-space     → harness's own edit layer
  │     └─ verification               → calls LineSpec
  └─ imprint emission                 → calls LineSpec
```

Sedum never calls the harness. No callback, hook, or plugin point in the pipeline —
adding one inverts this direction and makes Sedum's guarantees depend on code it
does not control.

Sedum stays independently runnable, which keeps "commit a recording as a standard
service scaffold and replay it" shippable with no harness, no containers, and no
model server.

---

## 4. The integration surface is file formats, not an API

> **Everything a tool built on Sedum needs is derivable from artifacts on disk.
> Sedum exposes no runtime API, and nothing links against Sedum.**

This used to read "a harness never links against Sedum", which permitted a
standalone format module that Sedum and a harness both linked against. The PRD
said "nothing links against Sedum" and is authoritative, so the strict reading
governs and the first-party harness is included (`prov-2026-b5465dfa`).

The reason travels with the rule, because a constraint recorded without its
justification reads as fastidiousness and is relaxed by whoever first finds it
inconvenient. **The formats must be sufficient without linking, and the only
thing that forces them to be is that nobody gets an exemption.** Once a
first-party harness can link Go, every awkward format is cheaper to route around
with one more exported helper than to fix. Third-party harnesses inherit the
gap, and the gap is invisible to the party best placed to notice it, because
that party is not experiencing it. The same party authors these formats and is
their most privileged consumer, so the rule that we eat the same food is the
only thing standing between the surface and quiet insufficiency.

Six surfaces, counted as six rather than folded into fewer, because a surface
that grows silently is one nobody audits. `sedum render` was added under
`prov-2026-b5465dfa`; `conformance/` under `prov-2026-c6580f1e`, which is the
answer to the objection this section's own rule invites — if nobody may link
Sedum to borrow its implementation of a format, checkability is the whole of
what a stranger gets, and leaving it as prose is a gap in the surface rather
than in the documentation.

| Surface | Read/write | Who else authors it |
|---|---|---|
| **Ownership markers** | harness reads; may write within its own rules | Sedum writes them |
| **Generator package** (`sedum.yaml`, `actions.yaml`, templates) | harness reads | users author them |
| **Recording** (JSON) | harness **writes**, Sedum reads | Sedum can write them too |
| **`sedum grow --execute`** | harness invokes | — |
| **`sedum render`** | harness invokes | — |
| **`conformance/`** | harness reads | Sedum's own tests read it too |

**"What exists?"** — grep the markers. ODQ §13's world-state scan is a filesystem
walk requiring nothing from Sedum's process.

**"Could I fix this deterministically?"** — `sedum render --package <name>
--action <name> --kwargs <json>` renders that action's `injects_into` against the
kwargs already on the marker, and the harness compares the result to the region's
file. Forward rendering and comparison, never inversion — which matters, because
`snake` and `plural` are not invertible. It works *only because kwargs are
recorded on the marker*, which ODQ §3 asserts and which this makes load-bearing.

This is a command because the instruction it replaces was not executable from
artifacts on disk. `sedum.yaml` carries the pipelines and `op_exceptions`, but
the inflection table ships per language rather than per package, and the
word-splitting semantics and the leading-case match are not published at all. A
harness answering this by hand reimplements `plural`, diverges on the first
irregular, and concludes no deterministic fix exists when one does — a wrong
answer wearing the shape of a correct refusal. That this section's own claim was
false for its own second question, before any foreign consumer existed, is the
demonstration the strict reading was argued from.

**"Do it."** — synthesise a recording of just those invocations and `--execute`
it. A harness doing its own selection never touches Phase 4; `--record --dry-run`
gives the other direction.

**"What did Sedum decline to do?"** — `unmanaged` patterns are in `sedum.yaml`,
which the harness already reads. Pattern entries authorizing without naming
(`prov-2026-e8671c88`) are derivable the same way.

Consequence worth stating in the PRD: **the recording is an input format, not only
an output artifact.** M7's "validation is identical, but failures are terminal" is
already the right semantics for a submission protocol.

---

## 5. What is actually still open

Rewritten against the M5 implementation. The marker design has already absorbed
most of what the prior revision asked for, and did it better.

### 5.1 Already settled in code — the PRD is what is stale

**Forward compatibility.** `marker.go` carries the attributes as a JSON object
precisely so fields can be added later: *"Positional fields would make every field
added later a migration across every repository that already carries markers; a
JSON object makes it an addition."* `parseOpen` treats an unreadable object as
corruption rather than version skew, on the same reasoning.

*This supersedes the prior revision's call for a version token, which was
misjudged.* A version token earns its place only for **breaking** changes — an
existing field changing meaning. Additive evolution is already free. It remains
cheap insurance, but it is not the priority the prior revision claimed.

**Record ID on the marker.** Already implemented, under
`prov-2026-36c8a99c` — *"The record ID is an attribute of a region, not part of
its identity."* The reasoning there (a later record refining a region replaces it
in place rather than minting a duplicate beside it) is stronger than the PRD's
deferral.

**Tier precedence.** The marker carries the tier and `seeded` is honored on read,
*"so that a region which has stopped being Sedum's to overwrite can say so."* That
is marker-governs in behavior. What is missing is the *stated rule*, so that a
future implementer does not "correct" it back toward the `actions.yaml`
declaration.

**Phase 5's path check.** `expand.authorized` matches against `[]resolve.File`,
and `resolve.File` carries `Existed bool` — pre-existing files are in the set. The
behavior is right; incremental re-invocation against an existing file passes. Only
the wording is wrong: the comment says *"exactly one file the run created"* and the
error says *"not one of the paths this record created."* The prior revision called
this a spec bug that would block the harness. It does not — but the phrasing will
mislead a harness author reading the source, which under an open-source release is
the same cost more slowly.

⚠︎ **Correction — the `actions.yaml` tier key does not exist yet.** Strict
decoding rejects it, so "the declaration is the value used at first write" is
the rule for when the key lands, not a description of today. Today the marker is
the only source of a region's tier. Stated that way in both the PRD and README.

### 5.2 The one real gap: unknown keys survive decoding but not the round trip

`attrs` is a fixed struct — `Tier`, `Record`, `Kwargs`. Decoding tolerates keys it
does not declare, because `encoding/json` ignores them. But `Marker` does not
retain them, and `Open` re-encodes a fresh `attrs`. So a key Sedum does not
recognise is **silently dropped the next time the region is rewritten.**

The doc comment's promise — that reader and writer of a marker are routinely
different versions of Sedum — holds for reads and fails for round trips. Today the
only cost is skew between Sedum versions. Under a third-party harness it is the
difference between markers being an extension point and not:

- With preservation, a harness annotates a region with its own state — escalation
  attempt count, last-verified spec, the record that modified it — and gets a
  private extension point that requires no cooperation from Sedum, ever.
- Without it, every harness needing per-region state must either petition for a
  schema change or maintain a sidecar file keyed by region, which reintroduces
  exactly the maintained state the marker design exists to avoid.

The fix is small and local: a catch-all (`Extra map[string]json.RawMessage`,
`json:"-"` plus explicit merge, or a two-pass decode) retained on `Marker` and
re-emitted by `Open`. Cheap at M5. Expensive once harnesses exist and their
annotations start disappearing.

**This is the single item most worth doing before M6.**

⚠︎ **Correction — retention and re-emission are not sufficient.** `applyOne`
builds a *fresh* `Marker` from the invocation and renders it over the region,
using the parsed `region.Marker` only for identity and tier before discarding
it. So `Extra` never reaches the writer, and a catch-all retained on `Marker`
and re-emitted by `Open` would read correctly and still drop every annotation on
the first rerun — worse than not doing it, because it would look done.

The replacement path has to carry the parsed region's unknown keys forward onto
the marker it writes. Shipped as a two-pass decode plus that carry-forward, with
a test that fails when the carry-forward line is removed.

Two details settled in implementation. Declared keys keep their struct order and
carried keys are appended sorted, rather than merging everything into one map
and marshalling — a map sorts its keys and would have reshuffled the attribute
object on every marker already written. And a carried key may not shadow a
declared one, which cannot arise from a marker Sedum parsed but can from a
`Marker` a caller builds.

### 5.3 The writer field is still absent

`attrs` has no writer or authority. Without it:

- Phase 3's diagnostic (*"something other than Sedum wrote it"*) is not merely
  wrong in a two-writer world but undiagnosable.
- Tier demotion under 5.1 is unattributable — a region reading `seeded` gives no
  way to tell whether a package author declared it or a harness demoted it.

Cheap now for the same reason record was: one key in an object that already exists.

**Shipped** as one free-form string naming the tool that last wrote the region,
omitted when that tool is Sedum. Omitting rather than writing `"sedum"` keeps a
marker Sedum writes byte-identical to one written before the key existed, so
introducing it rewrites nothing and no marker on disk needs migrating. A
writer/authority pair was considered and rejected: authority has no consumer
yet, and its vocabulary would have been guessed now rather than learned from use.

Lower-priority companion: **package** on the marker. Derivable from the extension
map, ambiguous only where two packages contest an extension and the run's `--lang`
choice is not recoverable from disk. Still open.

### 5.4 What a foreign writer may not do is unwritten

At minimum: do not remove or relocate a template-planted marker; do not reorder
regions within a file; preserve unknown marker keys in both directions. Unwritten,
these become discovered constraints, and the first breakage looks like a Sedum bug.

**Written, and the checkable ones mechanized** (`prov-2026-59535684` for the
prose, `prov-2026-c6580f1e` for the corpus). `conformance/markers/cases.json`
pins the key-preservation rule at byte level in both directions, plus the
carry-forward through the replacement path, declared-then-sorted emission, the
shadowing rejection, and corruption-versus-skew. The two file-level rules —
relocating a planted marker, reordering regions — are not expressible as a
parse-and-emit golden and stay prose; the corpus's README names them as
unmechanized rather than letting the corpus's existence imply otherwise.

---

## 6. Audit

### Conflicts

**1. Tier precedence unstated.** See 5.1.

**2. "Multiple records injecting into the same file is out of scope."** Contradicted
by the shipped `record` field. PRD is stale.

⚠︎ **Correction — stale in one direction only.** `record.checkDuplicatePaths`
still rejects a path named by two records at ingestion, on the grounds that
Phase 4 makes one model call per record. So refinement of a region across
records and across runs is shipped and is what the `record` attribute exists
for; two records naming one path *in one run* still halts, and records sharing a
file are generated one at a time with `--only`. Replacing the sentence with
"multiple records are supported" would have handed a harness author a guarantee
that is not there, so the PRD and README now state both halves.

This has a consequence for M7, folded into `prov-2026-dc227be7`: under
`--execute` there is no model call, so the check's justification does not hold
on the replay path — yet a caller supplying records purely for scope validation
routes through the same ingestion and would halt. Whether the check is skipped
under replay or moved to where its justification lives is M7's to settle.

**3. Phase 3's "something other than Sedum wrote it."** Survives only via a
constraint written down nowhere. See 5.3, 5.4.

**4. `unmanaged` breaks beat one's guarantee.** *(From the in-flight draft
`prov-2026-529954ab`.)* ODQ §13's justification for two beats is that after
structural completion no failure can mean "not built yet." An unmanaged `Gemfile`
without the pg gem means the service does not boot, and a linespec test fails for
exactly the structural reason beat two was supposed to have eliminated.

Beat one is therefore three steps: Sedum scaffolds → the unmanaged paths are
satisfied → beat two runs. ODQ §13 has no slot for the middle step, and it is not
Sedum's to fill.

This is not an argument against `unmanaged`. "The handoff is the point" is right,
and `unmanaged` together with `prov-2026-e8671c88` turns `affected_scope` from "the
list of files Sedum will make" into an authorization surface with several
categories, only one of which Sedum acts on — exactly the shape a record needs when
more than one tool works under it. What it means is that ODQ §13's beat boundary
moves, and a harness must treat an unmanaged path as a precondition.

### Wording

**5. "Scanning existing codebases"** is broader than its intent — read literally it
forbids §4's marker scan, and already contradicts Phase 3. The principle worth
keeping: **no structure inference from unmarked source.**

**6. `owned` says "hand edits are lost."** Restate as edits by any writer.

**7. `authorized()`'s comment and error message say "created."** Should say
authorized and managed by this run. See 5.1.

### The PRD supports the split

**8.** The cross-action value flow non-goal already says a planning layer would be
*"a substantially different architecture, not an increment on this one."* ODQ §7's
sort is that layer; the split gives it a home that is not Sedum.

**9.** The behavioral verification non-goal is the split, stated.

**10.** "Recordings carry no volatile fields" resolves ODQ §2's spec-pass-set
tension: pass sets are the harness's artifact.

**11.** `--stop-after`'s Sedum-internal vocabulary stays closed; ODQ §1's
`--stop-after tests` belongs to the harness's CLI.

---

## 7. Edits to make — done

Grouped as they were governed. Both records are `implemented`.

### Record A — the marker is a public format — `prov-2026-72775ae5`

*Mostly PRD catching up to code, plus three small changes in `internal/inject`
and `internal/expand`.*

| Change | Where | |
|---|---|---|
| Retain and re-emit unrecognised attribute keys (5.2) | `marker.go` | ✅ |
| Carry them forward through the replacement path (5.2 ⚠︎) | `inject.go` | ✅ |
| Add the writer key (5.3) | `marker.go`, PRD | ✅ |
| State that the marker's tier governs after first write; the `actions.yaml` declaration is the value used at first write only — *scoped as the rule for when that key lands* (5.1 ⚠︎) | PRD, *Ownership tiers* | ✅ |
| Replace "multiple records… out of scope" with the `record` attribute as shipped — *stating both halves* (§6.2 ⚠︎) | PRD, *Ownership and Idempotency* | ✅ |
| Restate `owned`'s "hand edits are lost" as edits by any writer | PRD | ✅ |
| State the attribute object as the extension point, and unrecognised-key preservation as a promise | PRD | ✅ |
| Fix `authorized()`'s comment and error text: authorized and managed, not created | `expand.go` | ✅ |

### Record B — the integration surface — `prov-2026-59535684`

*New PRD section. No code.*

- ✅ The three formats and one command that constitute the supported surface (§4)
- ✅ No runtime API, no hooks, no callbacks — with the reason (§3)
- ✅ The recording is an input format, not only a capture artifact
- ✅ What a foreign writer may not do (5.4)
- ✅ Restate the "scanning existing codebases" non-goal as no structure inference
  from unmarked source

### Open Questions edits

1. ✅ **Reframe the preamble.** It presented §1–§15 as directions "kept out of
   the PRD because most of them depend on answers M1–M7 will produce," which
   implies they return to Sedum. §4–§14 do not. The split is also not clean by
   section — §1, §2, §3 and §12 are marked inline where the halves diverge.
2. ✅ **Extend the pulled-forward list** with the record attribute
   (`prov-2026-36c8a99c`) and the marker extension point (`prov-2026-72775ae5`).
3. ✅ **§4** — the "third state" resolved by deferral.
4. ✅ **§13** — the beat-one break (§6, finding 4).
5. ✅ **§3** — in-parameter-space re-invocation against an existing file already
   validates correctly; the constraint is the marker round trip, not path
   validation.
6. ✅ **§15** — `forbidden_scope` becomes load-bearing as soon as anything working
   under a record edits an existing repository, not only when Sedum does.

### Also carried out

- **README** mirrored on both records. It was ahead of the PRD on the record
  attribute and marker tolerance, behind on everything else.
- **M7's blueprint** (`prov-2026-dc227be7`, still `draft`) gained the replay
  duplicate-path question from §6.2's correction, and the recording-as-input
  constraint.

---

## 8. Open

**Naming the harness.**

~~**Whether a marker version token is worth adding anyway.**~~ **Decided: no.**
Additive evolution is free under the object, and unrecognised-key preservation
now makes it free for foreign keys too. A token earns its place only for a
breaking change — an existing field whose *meaning* changes — and nothing on the
horizon is one. Recorded as a constraint on Record A so it is a decision rather
than an omission.

~~**Whether `--only` needs region granularity.**~~ **Decided: no.** A harness
synthesises a single-invocation recording, which works today. The consumer that
arrived needs no selector vocabulary, and a surface grows when a caller
demonstrates it must, not before (`prov-2026-b5465dfa`).

**Whether the duplicate-path check applies under `--execute`** (§6.2 ⚠︎). In
M7's scope, in `prov-2026-dc227be7`.

**Where `adopt` lives** (ODQ §3).

**Whether `unmanaged` entries should carry a reason** — not for the human/agent
distinction, which rarely matters, but a harness deciding whether to queue a path
or halt on it may want more than the pattern. ODQ §13's beat-one break makes this
sharper than it was: an unmanaged path is a *precondition* of beat two, and a
caller that must decide whether to block on one has only the pattern to go on.

**A reserved consumer namespace in the generator package is not open, and is
recorded here because it was never written down.** It was proposed as somewhere
a caller could keep per-stack configuration Sedum validates the shape of but
never interprets. The consumer that arrived derives that configuration instead,
so it is not being added on a consumer's account. It remains defensible on
Sedum's own extensibility grounds and is unclaimed rather than refused
(`prov-2026-b5465dfa`).

**Whether `package` belongs on the marker** (5.3).

**Whether `forbidden_scope` should be checkable independently of who is acting**
(ODQ §15). Sedum enforces it over what Sedum does; nothing enforces it over what
a tool above Sedum does except that tool.
