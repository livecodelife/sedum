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

> **Everything a harness needs is derivable from artifacts on disk. Sedum exposes
> no runtime API, and a harness never links against Sedum.**

| Surface | Read/write | Who else authors it |
|---|---|---|
| **Ownership markers** | harness reads; may write within its own rules | Sedum writes them |
| **Generator package** (`sedum.yaml`, `actions.yaml`, templates) | harness reads | users author them |
| **Recording** (JSON) | harness **writes**, Sedum reads | Sedum can write them too |
| **`sedum grow --execute`** | harness invokes | — |

**"What exists?"** — grep the markers. ODQ §13's world-state scan is a filesystem
walk requiring nothing from Sedum's process.

**"Could I fix this deterministically?"** — read `actions.yaml`, render each
action's `injects_into` against the kwargs already on the marker, compare to the
region's file. Forward rendering and comparison, never inversion — which matters,
because `snake` and `plural` are not invertible. It works *only because kwargs are
recorded on the marker*, which ODQ §3 asserts and which this makes load-bearing.

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

### 5.3 The writer field is still absent

`attrs` has no writer or authority. Without it:

- Phase 3's diagnostic (*"something other than Sedum wrote it"*) is not merely
  wrong in a two-writer world but undiagnosable.
- Tier demotion under 5.1 is unattributable — a region reading `seeded` gives no
  way to tell whether a package author declared it or a harness demoted it.

Cheap now for the same reason record was: one key in an object that already exists.

Lower-priority companion: **package** on the marker. Derivable from the extension
map, ambiguous only where two packages contest an extension and the run's `--lang`
choice is not recoverable from disk.

### 5.4 What a foreign writer may not do is unwritten

At minimum: do not remove or relocate a template-planted marker; do not reorder
regions within a file; preserve unknown marker keys in both directions. Unwritten,
these become discovered constraints, and the first breakage looks like a Sedum bug.

---

## 6. Audit

### Conflicts

**1. Tier precedence unstated.** See 5.1.

**2. "Multiple records injecting into the same file is out of scope."** Contradicted
by the shipped `record` field. PRD is stale.

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

## 7. Edits to make

Grouped as they would be governed. Both PRD blueprints touch `PRD.md`, which is
already in `prov-2026-529954ab`'s `affected_scope` — worth sequencing so the
unmanaged draft is not left straddling them.

### Record A — the marker is a public format

*Mostly PRD catching up to code, plus two small changes in `internal/inject`.*

| Change | Where |
|---|---|
| Retain and re-emit unrecognised attribute keys (5.2) | `marker.go` |
| Add the writer/authority key (5.3) | `marker.go`, PRD |
| State that the marker's tier governs after first write; the `actions.yaml` declaration is the value used at first write only | PRD, *Ownership tiers* |
| Replace "multiple records… out of scope" with the `record` attribute as shipped | PRD, *Ownership and Idempotency* |
| Restate `owned`'s "hand edits are lost" as edits by any writer | PRD |
| State the attribute object as the extension point, and unrecognised-key preservation as a promise | PRD |
| Fix `authorized()`'s comment and error text: authorized and managed, not created | `expand.go` |

### Record B — the integration surface

*New PRD section. No code.*

- The three formats and one command that constitute the supported surface (§4)
- No runtime API, no hooks, no callbacks — with the reason (§3)
- The recording is an input format, not only a capture artifact
- What a foreign writer may not do (5.4)
- Restate the "scanning existing codebases" non-goal as no structure inference from
  unmarked source

### Open Questions edits

1. **Reframe the preamble.** It currently presents §1–§15 as directions "kept out
   of the PRD because most of them depend on answers M1–M7 will produce," which
   implies they return to Sedum. §4–§14 do not. Note also that the split is not
   clean by section: §1, §2, and §3 each have a Sedum half and a harness half.
2. **Extend the pulled-forward list** — it names ownership tiers, recorded kwargs,
   and the `phases` grouping. Add the record attribute (`prov-2026-36c8a99c`) and
   whatever lands from Record A.
3. **§4** — mark the "third state" resolved by deferral. With unrecognised-key
   preservation, a harness expresses it in its own keys, and it can be standardised
   later once its shape is known from use rather than guessed now.
4. **§13** — fold in the beat-one break (§6, finding 4).
5. **§3** — note that in-parameter-space re-invocation against an existing file
   already validates correctly; the constraint is the marker round trip, not path
   validation.
6. **§15** — `forbidden_scope` becomes load-bearing not only when Sedum edits an
   existing repository but as soon as a harness does.

---

## 8. Open

**Naming the harness.**

**Whether a marker version token is worth adding anyway.** Not needed for additive
evolution (5.1). The question is whether any field's *meaning* is likely to change.

**Whether `--only` needs region granularity.** A harness can synthesise a
single-invocation recording today, which may be sufficient and avoids a new
selector vocabulary. Confirm before M7.

**Where `adopt` lives** (ODQ §3).

**Whether `unmanaged` entries should carry a reason** — not for the human/agent
distinction, which rarely matters, but a harness deciding whether to queue a path
or halt on it may want more than the pattern.
