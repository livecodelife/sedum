# Sedum — Open Design Questions

Scratch document. Nothing here is a requirement. These are directions explored after the MVP scope was settled, recorded so the reasoning isn't lost.

**Most of this document is no longer about Sedum.** It was originally framed as directions kept out of the PRD because most of them depend on answers M1–M7 will produce — which implied they would return to Sedum once those answers arrived. `TOOL_BOUNDARIES.md` decides that they will not. §4–§14 describe a **convergence loop**: a tool that sits above Sedum, drives generation, verifies the result by running something, and repairs what did not converge. That is a different tool with a weaker guarantee — open synthesis validated by execution, rather than bounded selection over a closed vocabulary — and a single binary containing both would take the weaker guarantee as its own.

So these sections are not deferred Sedum work. They are the design notes for a tool that will be built on Sedum's [integration surface](PRD.md), and the useful question about each is no longer "when does Sedum do this" but "what does Sedum have to expose so something else can."

**The split is not clean by section.** §1, §2, and §3 each have a Sedum half and a harness half, and they are marked inline where the halves diverge. §12's record *drafter* belongs to LineSpec, which already owns record authoring — `provenance create` / `discover` / `next` are the same shape of operation, and drafting is `discover` pointed at intent rather than at source. §16 is a mix of both.

Several decisions from this exploration were pulled forward into the PRD, because they write durable artifacts whose shape is expensive to change later:

- **ownership tiers** (`owned` / `seeded`) and **recorded kwargs** on markers
- the **`phases` grouping** in recordings
- the **record attribute** on the marker, as an attribute rather than as part of a region's identity (`prov-2026-36c8a99c`)
- the **attribute object as an extension point**, with unrecognised-key preservation promised in both directions, plus the **`writer`** attribute (`prov-2026-72775ae5`)

That last one changes what several sections below need from Sedum, because a tool above Sedum can now record its own per-region state on the marker without a schema change.

Sections 7–14 come from a later session on imprint integration, cross-region logic, and the record-authoring workflow. Several of those threads reached a settled position rather than staying open; they are marked as such where they did.

---

## 1. Generated tests as a first pass

> **Split.** Generating test files from templates and actions is Sedum's — it is
> ordinary generation into authorized paths, and the `phases` grouping in
> recordings is already reserved for it. Feeding pass one's rendered output into
> pass two's prompt, and stopping between them for review, belongs to the caller:
> `--stop-after tests` is not a Sedum phase name, and Sedum's `--stop-after`
> vocabulary stays closed to its own seven phases.

### The idea

Add a config value per package declaring which paths hold tests. Sedum generates those files first, in their own pass, using ordinary file templates and actions the author defines (`createDescribe`, `createContext`, `createItBlock` — whatever vocabulary the target's test framework wants). The rendered test files then go into the prompt for the implementation pass, alongside intent and constraints.

### Why it might work

A generated test is a **reified interpretation of the constraints**. Today the model reads constraints and picks implementation actions in one step, with the interpretation implicit and unverifiable. Splitting it makes that interpretation a file — inspectable, diffable, correctable — before anything depends on it. Same move as separating `resolve` from `grow`, applied to intent rather than to paths.

### What has to be true

**Config is a pattern list, not a directory.** Go puts `user_test.go` beside `user.go`: same extension, same package, no separate tree. So `test_paths: ["spec/**", "**/*_test.go"]`, matched against `affected_scope` per package.

**Pass two must receive rendered test files, not just ordering.** The ordering by itself buys nothing. The wiring change is that pass one's output files become pass two's prompt context.

**Intent and constraints stay in pass two.** A test can't express "use the existing UserSerializer" or a negative constraint about what not to touch. Cheap to include; dropping them inherits every gap in pass one.

**Review has to be possible between passes.** `--stop-after tests` or equivalent. The entire value is a human correcting the interpretation at the cheapest point. If the default runs straight through, this adds a step and keeps the same failure mode.

### The honest risk

This is the one place the model writes freeform content rather than binding kwargs. `it` block bodies aren't enumerable variants. The failure mode differs from everything else in the pipeline: a plausible-looking wrong test doesn't error, it silently steers the implementation wrong. Whether the structure is templatable at all is entirely the package author's call — `describe`/`context` blocks may be more constrainable than assertions, but that's their design decision, not Sedum's.

### Relationship to linespecs

These are different artifacts with different authority, and will get conflated by anyone reading this later.

**Linespecs** are behavioral contracts authored ahead of implementation, under their own provenance record, used as the gate.

**Generated unit tests** are model-produced interpretation, used as generation context and as a local validation signal.

A generated test never carries the authority of a hand-authored contract.

---

## 2. Associated specs

> **Split.** Specs as *selection evidence* is Sedum's: a spec is text in the
> Phase 4 prompt, Sedum never parses it, and the no-target-knowledge principle
> survives. Specs as *signal* and as *regression fixtures* are the caller's —
> both require running something, which is behavioral verification and a stated
> non-goal. `run_command` is arbitrary shell and belongs in `.linespec.yml`,
> where target knowledge already lives, not in a generator package.
>
> The PRD settles the fixture tension in passing: recordings carry no volatile
> fields, so a recorded pass set could never have lived in one. Pass sets are the
> caller's artifact.

Specs as **selection evidence** is the cheapest available accuracy lever and requires no execution: a spec asserting `GET /users → 200`, `POST /users → 201`, `DELETE /users/:id → 204` tells the model which methods to generate far more reliably than prose intent does. Sedum never parses it — for prompt purposes any spec is text, which preserves the no-target-knowledge principle.

Specs as **signal** — run once after generation, report, never retry. Boilerplate mostly won't pass behavioral specs and shouldn't be expected to, but a pass count is useful information and one run is not a verify loop.

Specs as **generator-package regression fixtures** is the most interesting use. Record the pass set at record time; compare on replay. A scaffold that passed 3 and now passes 2 means a template regressed. This makes recordings the test suite for generator packages, which is otherwise hard to test.

`run_command` is arbitrary shell — same trust posture as user-defined transforms. Gate behind an explicit flag, never let it block.

---

## 3. Updates to existing code

> **Split.** Additive change to a Sedum-generated file and changes inside the
> parameter space are Sedum's, and both already work. Changes outside the
> parameter space — refusing well, and then the edit layer that acts on the
> refusal — are the caller's. Where `adopt` lives is open.

### The easy cases

**Additive change to a Sedum-generated file** is nearly free: Phase 3 becomes create-if-absent, verify-markers-if-present. A record saying "add destroy" resolves to `createControllerMethod(name: destroy)` injected at the existing marker. This likely covers most real update traffic.

**Additive change to a file Sedum didn't create** has no markers. Hard-error, and offer an `adopt` path that inserts markers only where the file's skeleton still literally matches a file template. Anything fuzzier needs structure inference, which means a parser, which is the thing avoided everywhere else.

### Changes inside the parameter space

A surprising amount of update traffic is a kwarg change. "Include addresses in the index response" is `render json: @users, include: :addresses` — if the `index` variant accepts an `includes` kwarg, the update is a re-invocation against an owned region. Fully deterministic, no code editing at all.

This is why kwargs are recorded on the marker: current state becomes readable from the file, and desired-vs-actual becomes a diff over kwargs.

**This half is Sedum's, and it already works.** Re-invoking an action against a file an earlier run created validates correctly: Phase 3 is create-if-absent, so an existing path stays in the authorized set rather than dropping out of it, and Phase 5's path check accepts it. `prov-2026-72775ae5` corrected the wording that suggested otherwise — the check is authorized-and-managed, not created-by-this-run. The constraint worth watching is the marker round trip, not path validation.

### Changes outside the parameter space

"Exclude soft-deleted users unless the requester is an admin" is business logic. The useful behavior at this boundary is **refusing well**: detect that the record targets a `seeded` region or requires kwargs the template doesn't accept, and emit a work item naming the file, the region, and the constraint.

---

## 4. The focused edit layer

A work item from the boilerplate layer is a much better-conditioned task than open-ended synthesis. It carries a file that already exists with correct structure, a named region with known bounds, current contents that are usually close to right, a specific constraint describing the delta, and the marker kwargs documenting what the region was parameterized to do. That is a scoped edit against working code, not synthesis — a different difficulty class, and plausible for a local open-weight model in a way the general problem is not.

### Design positions

**Structured edits, not region rewrites.** Ask for a replacement region and the model rewrites the whole thing with incidental drift. Ask for a small list of `{after, insert}` / `{replace, with}` operations validated against current text, and the blast radius is bounded, invalid edits are rejectable without execution, and the human reviews the actual change. Same structured-output discipline as Phase 5, one layer down.

**Edited regions change tier.** A region the edit layer touches can't stay `owned` — the next replay would destroy the edit. It becomes `seeded`, or a third state recording both the originating invocation and the modifying record.

*The third state is resolved by deferral.* With unrecognised-key preservation, the edit layer records whatever it needs about the edit in its own marker attributes — the modifying record, the attempt count, the spec that last verified it — and `seeded` plus `writer` covers the part Sedum has to understand: this region is no longer Sedum's to overwrite, and here is who took it. A standardised third tier can be added later from what the keys turn out to be used for, rather than guessed at now. Guessing it now would put a tier in every marker on disk before anything reads it.

**`modifyContent` takes a region marker, not a file and method name.** "Method name" requires knowing how methods are delimited in the target language. A marker reference is target-agnostic and already carries the kwargs.

**Guards:** opt-in per package, so teams can choose purely deterministic generation. A required `reason` field, since it's the easy answer to everything and a model that can't find the right action will reach for it. Atomic failure — restore the region and emit a work item rather than leaving half-edited logic.

### The precondition

Focused edits are only worth attempting if you can tell whether they worked, which means running the relevant assertions per attempt. Validation must target **the specific blocks covering the change**, not the file's whole suite, or the signal doesn't isolate and the latency is unusable.

**Warm-harness design is the gate on this layer existing at all.** Build it before the model integration, not after.

---

## 5. Two-phase updates and the ordering question

For an update record, both passes use the edit mechanism: an `updateTest` action editing a marked region of the test file, then an `updateLogic` action editing the implementation with the updated test in context.

### This is not a planning problem

`updateLogic` depends on its `updateTest`. That's one fixed dependency between two things, known at design time. A planner discovers orderings you couldn't enumerate; this one is enumerable in a sentence. Hardcoding test-then-logic is the honest implementation.

### The real argument for pairing is failure isolation

Batch all test updates then all logic updates, and a failure on the fourth logic edit leaves four test regions and three logic regions changed with no clean rollback point. Pair them and each unit is atomic: edit test region, edit logic region, run the relevant assertions, pass or restore both.

### The unit is the constraint, not the region

One constraint may add three assertions that a single logic edit satisfies. So the model's response groups edits under the constraint that motivated them, and Sedum iterates constraint by constraint. Partial success becomes meaningful — three of five constraints satisfied, two rolled back with work items naming exactly which failed.

### Open

**Does a failed pair block subsequent pairs?** Probably not — constraints are usually independent and stopping at the first failure wastes the run. But that assumption breaks quietly rather than loudly if constraints do interact. Wants a `--fail-fast` escape and a log warning when multiple constraints touch the same marker.

**Where a planner would genuinely earn its place** is if edits started conflicting — two constraints touching one region where order changes the outcome. That's not this. And if it happens, rejecting the ambiguity is probably better than searching it.

---

## 6. Test regeneration on update records — unresolved

The sharpest open problem in this document.

A **create** record's intent describes complete behavior, so generating a spec from it is sound.

An **update** record's intent describes a delta — "also include addresses." Regenerate the spec file from that and you get a spec covering only addresses, silently dropping every existing assertion. Then the edit validates against a test that no longer checks what the method already did, and a regression passes clean.

So pass one on updates can't be regeneration. It has to be **amendment**: existing spec file in context, model emits actions that add or modify specific blocks, existing blocks untouched. A spec block generated by a prior record is owned by that record; a new record adds its own rather than replacing the file.

This makes marker kwargs on spec blocks load-bearing — it's how you determine what's already asserted.

Unresolved: what happens when an update record genuinely invalidates a prior assertion rather than adding to it.

---

## 7. Cross-region ordering — a sort, not a planner

Two problems hide inside "cross-region," and they have different answers.

**Ordering** is deterministic and tractable. Dependencies between regions are stated or extractable, which makes ordering a topological sort over a dependency DAG — linear time. Reaching for GOAP here would be building a planner because one was already lying around.

**Content correctness** is not a planning problem at all. It is a test problem, addressed by the loop in §5 and the traceback in §9.

### Where the edges come from

Three options, escalating in cost and in how much target knowledge they require.

**Declared dependencies** in the record or kwargs. Cheapest, least complete.

**Marker-level metadata** — the chosen first step. Ownership markers record what a region exposes and what it consumes. The graph is built by scanning markers: no language parsing, no AST, target-agnostic, consistent with every other position in this document.

**AST extraction.** Precise, and the wrong shape. If it ever becomes necessary it lives in the generator package as a per-language dependency extractor, never in core.

### What the ordering actually buys

Not correctness by itself. The unlock is **context assembly**. Because the sort guarantees upstream region A is proven before B is generated, B's prompt can include A's proven implementation as context. That collapses a non-local problem into a local one a small model can handle. The sort's entire contribution is the ordering plus inject-upstream-as-context.

Test ordering follows the same topological order as code, because a downstream spec is only writable once the upstream contract exists.

Downstream specs are **amended, never regenerated** — same position as §6, and for the same reason: regeneration silently drops the assertion protecting the coupling.

---

## 8. Integration specs live outside the codebase

Two tiers of test with two custodians.

**In-tree**: generated per-region unit specs. Model-authored, model-amendable. They pin local behavior.

**Out-of-tree**: human-authored linespec integration specs. Model-*readable* as context, never writable. They pin cross-region and whole-service behavior. A human may use a model to draft one, but review and submission are the human's, and the generation pipeline never touches the file.

Placing them out-of-tree makes the prohibition **structural rather than disciplinary**. The model cannot edit what does not appear in its operating world of markers and `affected_scope`. That is a stronger guarantee than a rule the pipeline is merely asked to respect.

It maps cleanly onto LineSpec: published manifests ship briefs and blueprints, while integration specs are proof artifacts referenced by path and need not sit in the service tree. They belong at blueprint level, outliving any single region. The contract that ships to another team *is* the artifact that gates generation.

This works only because linespec asserts at the protocol level — behavior, not implementation — so file layout and language are irrelevant to it.

The payoff: the blast radius of a bad spec amendment shrinks to one region, because the assertion protecting the coupling lives somewhere the model cannot reach.

---

## 9. Failure traceback — from a red integration spec to a region

### How linespec actually works

Entry-point-centric, not a sequence of interactions. One test takes one trigger — an API endpoint, an event consumer, a job — and asserts over the full fan-out: downstream service calls, database interactions, response. There is no "which step failed." A whole path lights up.

Multiple tests per entry point are distinguished by **branch**. An expiry check produces two: one where expiry has not been met (DB read shows valid, response A), one where it has (DB read shows expired, response B).

### The mechanism

Markers alone cannot localize, because a protocol-level test observes the wire, not the source. Three signals combine.

**Markers give the suspect set.** From the failing entry point, walk markers to enumerate every region on the path.

**The pass/fail pattern across branch tests gives the bias.** Both expiry-branch tests red → expiry was never stored → suspect the create path. Only the expired case red → storage works, the check is wrong → suspect the conditional in redirect. Only the valid case red → the happy path regressed. The discriminator is the *diff between which branch tests pass*, which requires nothing from linespec beyond what it already reports.

**Per-region unit specs give confirmation.** Re-run the candidates; the one whose own spec fails is localized.

### Design pressure

This rewards shallow paths, few branches per entry point, and granular individually-tested methods. Wide fan-out on a single branch produces large suspect sets and more human triage. Path width is a property of the user's design, not something Sedum can fix — worth stating plainly in documentation rather than letting users discover it.

---

## 10. The all-green case, and two repair modes

### The case

Integration spec red, candidate regions enumerated, every candidate's unit spec **passes**. By definition this is not a bug in any region — each is locally correct. It is a coupling failure, or equivalently a gap: some region's unit spec is silent about the thing the integration spec depends on.

The naive move — hand the candidates to a model, let it name a culprit, generate against the integration spec — reintroduces the silent-gap failure. The model adds logic, integration goes green, but no *assertion* pins the new behavior, so the next regeneration of that region can drop it and the unit spec will not notice.

### The repair loop

The fix is not code, it is an assertion. Whatever is missing has to be added to a unit spec so the gap becomes locally visible, at which point the ordinary loop from §5 has an honest target: spec now red, logic edit follows, integration re-run confirms.

The model's useful role is producing a **shortlist plus a proposed assertion per candidate**. What varies between the two modes is who ratifies that assertion.

### Gated mode (default)

A human approves the proposed assertion before the loop runs — one sentence, reviewed, at the single moment a new local contract is born. Everything after that is automatic. When a unit spec was *already* red, no new contract is being created and no gate applies at all.

### Automatic mode (opt-in)

The loop closes without the gate. This is defensible, because the human-authored integration spec is still the final gate. If the loop converges and integration passes, either the fix is genuine or the integration suite has a hole — and a hole is a *visible gap in a human contract*, not a silent code failure. Every escape is diagnosable and permanently closable by adding an integration assertion.

So the two modes are not safe versus unsafe. They differ in **where human attention is spent**: per-incident ratification, or one-time investment in integration coverage. The honest documentation line is that automatic mode is exactly as trustworthy as the integration suite, and every surprise it produces is a missing integration assertion.

Sparse suite → gated. Dense, battle-tested suite → automatic is reasonable.

---

## 11. The escalation ladder and its terminal rung

A small model selects from the closed action catalog. A mid-size model authors logic when selection comes up short. On repeated per-region failure, escalate to a larger open-weight model on infrastructure. Gating per region means the expensive tier only fires on regions that earned it, which keeps the common case on-device.

**The terminal rung is an unmet imprint handed to a human, not an automatic frontier call.** A region that defeated the largest local tier is usually a non-local problem, and more parameters do not reliably fix non-local problems. A frontier call is a human's opt-in from there.

### The load-bearing caveat

The ladder's stopping condition is "tests pass," so it is sound only if the test is trustworthy. A model-authored spec that is subtly wrong turns the ladder into a machine for **converging on a bad contract** — and larger models are better at satisfying a bad contract, not worse. This makes spec authoring the sharpest silent-failure seam in the design (see §1).

Which is why `--stop-after` on the test pass should not be optional for anything beyond input-output-shaped logic. The deterministic machinery cannot lie; the generated test is the one artifact that can. Human review converts the ladder from "trust the small model's contract" into "trust the engineer's contract, executed by small models."

Secondary benefit: the constraint pressures design toward simple, granular, input-output-shaped logic, which is the same gradient that makes traceback work (§9).

---

## 12. Two provenance domains, and who authors which record

> **Where this lands.** The record *drafter* belongs to LineSpec, which already
> owns record authoring — `provenance create`, `discover`, and `next` are the
> same shape of operation, and drafting is `discover` pointed at intent rather
> than at source. Nothing in this section is Sedum's: Sedum reads five fields out
> of a record and never writes one.

### The domains

The linespec suite has **its own** provenance hierarchy, separate from the code repository's. A blueprint there says "we want this behavior"; each discrete test edit — update this spec, create that one — is an imprint under it. This is where intent originates and where it is captured: at the moment the behavior is decided, by a human.

The code repository has its own hierarchy governing implementation.

### The handoff

The linespec-side blueprint and its imprints are the **input** to drafting a code-side blueprint — a high-level technical plan of what changes in the implementation to satisfy the behavior just declared.

This matters because it makes the task **translation, not origination**. The model is not abstracting intent from a bare test diff; it is projecting a fully articulated intent from one domain into a technical plan in the other. Substantially easier, and squarely within reach of a mid-size open-weight model.

A human reviews the drafted blueprint, iterates with the model or edits directly, and approves. Opening it is the trigger: specs run, they fail because behavior now outruns code, and the loop fires.

### Which tier authors what

The split follows the **cognitive act**, not parameter count.

**Briefs and blueprints are abstraction-up.** They require system-wide judgment — extend or supersede, does this conflict with a constraint three blueprints away. Capable model as *drafter*, human as author of record. This is not a new gate invented for the AI workflow: briefs and blueprints have always warranted human review in the provenance model. The model's contribution is turning a blank-page authoring task into a review-and-edit task on an artifact whose review was never optional.

**Imprints are description-down.** Everything an imprint records is observable from the generation event itself — the recording's invocation list, the marker written, the region changed. No climbing, no system-wide judgment. An imprint cannot really mis-govern, because it is reporting rather than governing. Small local models own these, via deterministic actions plus the free model for the residual, ungated.

This also reframes the capability question. A drafter does not need to be *right*; it needs to produce a draft whose review-and-edit costs less than authoring from scratch. Much lower bar than "governs unreviewed."

### Open

**Where rationale is captured** determines whether a mid-size model suffices. If the *why* is written at the linespec-side imprint, the code-side draft is propagation and a ~30B tier is credible. If rationale is expected to be reconstructed at drafting time, that is genuine abstraction and wants a frontier model plus review. Decide deliberately rather than discovering it.

**Frontier at authoring, open-weight at generation** is a coherent split — authoring is occasional and human-attended, generation is the high-volume inner loop — but it sits in mild tension with the on-device goal and should be recorded as a deliberate boundary.

---

## 13. TDD from records — structure pass, then behavior pass

### The problem

If opening a record runs the specs and failures trigger the loop, some failures will mean "this region does not exist yet" rather than "this region is wrong." A spec red because its *upstream dependency* has not been generated looks like a logic failure and is not. Feeding that to the judgment layer invites a model to fix something that is merely early.

### Classification is structural, not inferential

The world-state scan already answers this. Whether the region a spec targets exists is *read from the scan*, not inferred from the failure. Region absent → add. Region present → update.

An empty marked stub counts as **present**. It rides the update path with an empty starting point, consistent with the position that new-generation and update unify: a marked region exists, its content is insufficient against the spec, hand it to the judgment action with a delta. The only difference is whether the region held content.

### Two beats

**Beat one — structural completion.** Run the specs, then discard their pass/fail entirely and use the run only as *discovery*: which structure is referenced but absent. Create every missing file, scaffold every absent region, populate all deterministic boilerplate. No logic authoring in this beat. It ends with the world structurally complete.

**Beat two — behavioral convergence.** Re-run the specs. Now every failure is trustworthy, because no failure can mean "not built yet." The topological order, the update-test-then-logic loop, and the escalation ladder all operate on a signal with the structural noise removed.

### `unmanaged` moves the beat boundary

Beat one as written above is not enough to deliver its own guarantee, and `unmanaged` is why.

Sedum declines to write the paths a package declares unmanaged. Those paths are authorized, reported, and skipped — the handoff is the point. But a Rails service whose `Gemfile` never gained the `pg` gem does not boot, and a linespec test against it fails for exactly the structural reason beat two was supposed to have eliminated. Beat one ran, Sedum did everything it does, and the world is still not structurally complete.

So beat one is three steps, not one:

1. Sedum scaffolds what it manages.
2. The unmanaged paths are satisfied — by a person, or by a tool pointed at them.
3. Beat two runs.

The middle step has no slot in the two-beat framing and **is not Sedum's to fill**. Appending a gem to a `Gemfile` needs a format-aware primitive with merge semantics, which is a stated non-goal; declaring the path unmanaged is how a package says so out loud.

This is not an argument against `unmanaged`. Together with pattern entries that authorize without naming (`prov-2026-e8671c88`), it turns `affected_scope` from "the list of files Sedum will make" into an authorization surface with several categories, only one of which Sedum acts on — which is the shape a record needs as soon as more than one tool works under it. What it means is that a tool driving the loop must treat an unmanaged path as a **precondition** of beat two rather than as a report it can read and move past. Sedum tells it which paths those are; deciding whether to queue them, prompt for them, or halt is the caller's.

### Why this ordering rather than interleaving

The judgment layer never sees a failure it should not act on. The cheap deterministic beat absorbs all structural noise before the expensive fallible beat begins — the same instinct as keeping the model's surface small, applied to the *input* of the loop rather than its scope.

It also narrows the topological sort's job. All adds happen unconditionally in beat one, so ordering only has to sequence updates among existing regions, which is what it was good at.

### Worth stating

Running specs in beat one only to discard the result looks redundant — the needed structure could in principle come from `affected_scope` alone. It is worth running anyway, because specs reference regions and entry points that a record's `affected_scope` may describe only coarsely. On the first pass the run is a discovery mechanism, not a verification one.

---

## 14. Where this lands — the linespec suite as the authoring interface

The mature state of a Sedum codebase is one where engineers do not edit code to change behavior. They edit the linespec suite, under its own provenance, optionally with a capable model drafting the records. The changed spec describes behavior the code does not have, so it fails, and the failure *is* the trigger for everything in §7–§13.

The division this produces is the point. The highest-leverage act — stating the new truth at the contract level — stays with humans. The lowest-leverage act — writing code to match — goes to models. Most AI coding tools invert this, leaving humans to babysit code while the model roams over intent.

Editing-the-spec-as-programming is safe *because* the integration layer is human-owned and immovable. That is what makes the north star legitimate rather than merely appealing.

---

## 15. File-agnostic actions and the `injects_into` coupling

> **Sedum's.** This is action authoring and catalog shape. Nothing here needs
> the harness, and none of it is deferred behind M1–M7.

### The problem

`injects_into` is declared on the action, so an action whose template is
target-independent still needs one definition per target file group. `addImport`
becomes `addModelImport`, `addControllerImport`, `addSerializerImport` — same
kwargs, same template body, differing only in a path pattern.

The authoring duplication is the visible cost. The real cost is **catalog
quality**. Phase 4 hands the model a set of near-identical entries with identical
kwarg schemas, whose only discriminator is a name encoding a fact the model has
to infer. That is the worst possible shape for selection accuracy, which is the
thing the closed vocabulary exists to protect. And it scales multiplicatively:
file-agnostic actions × file groups.

### Root cause

`injects_into` conflates two things — *which file this action targets*, and
*what naming convention identifies that file*.

For a structure-creating action the target is implied by the action's identity.
`createControllerMethod` only ever means a controller, and the pattern is exactly
right. For a file-agnostic action the file is genuinely a parameter, and encoding
it in the action definition forces one action per binding of that parameter. The
convention is also already stated once, in `files/`; the action restates it.

### The direction: the target is a kwarg

Nothing forbids this today. `injects_into` is a pattern rendered against bound
kwargs, and load-time validation only inspects it when it is literal — a pattern
containing `{{` is deferred to Phase 6 by construction. So this is already legal:

```yaml
addImport:
  kwargs:
    file:   { type: string, required: true }
    symbol: { type: string, required: true }
    from:   { type: string, required: false }
  injects_into: "{{file}}"
  anchor: imports
```

Phase 5's `unauthorized_path` rule still holds: the rendered path must be one
this record authorized and this run created. **The path pattern was never the
safety boundary — `affected_scope` is.** And Phase 4 already hands the model the
list of paths created for the record, so it is selecting a path from a given set
rather than constructing one from convention.

### What is lost, and the better replacement

The pattern implicitly restricted applicability: `addBeforeFilter` could not be
aimed at a migration, because the path would not render to one. A free target
removes that guard. It should be replaced by **anchor applicability** rather than
reinstated.

An action's anchor already declares the region kind it needs. The check is that
the file template which matched the target path must plant the marker the action
is anchored to. Everything it requires exists:

- `resolve.Resolution.Template` names the file template that matched each path
- `genpkg.MarkersIn(commentPrefix, content)` extracts the markers one template plants
- `Action.MarkerAnchor()` gives the marker an action is anchored to

Today `plantedMarkers` does this package-wide, at load, as a warning — the best
available check when the target is a pattern. With the target known per
invocation it narrows to a per-file error at Phase 5:

```
addImport injects into db/migrate/001_create_users.rb, whose file template
db/migrate/*.rb plants no "imports" marker
```

This is strictly better than what the path pattern gave. It states the actual
precondition rather than approximating it by directory; it catches the case where
a path matches the pattern but the file has no such region; and it is the check
Phase 7 would have failed on anyway, moved from a hard error mid-write to a
re-promptable Phase 5 diagnostic.

It also reframes an action's applicability contract as **anchor vocabulary rather
than path shape**, which is a more honest description of what these actions are.
They operate on a region kind, not on a directory.

### What this depends on, and already has

Many-per-file invocation. `addImport` has no discriminator, so every invocation
would collide on `sedum:addImport` under label-only identity. `inject.IdentityOf`
already handles it: required kwargs select, optional kwargs parameterize. Two
imports with different `symbol` are distinct regions, and re-invoking one with a
different `from` replaces it in place.

Making `file` required puts it in the identity key, which is correct. The
consequence is that changing an invocation's target path orphans the region in
the old file rather than moving it — but that is already true of any required
kwarg feeding a path pattern (`controller` in `addBeforeFilter`), so it is not
introduced here.

### Verified, not reasoned

The claims above were checked by building the package and running it, not by
reading the code. A `chi2` fixture declares `addImport` with
`injects_into: "{{file}}"`, three file-template groups — handlers and models
plant an `imports` marker, store deliberately does not — and one record
authorizing one path in each. A stub OpenAI-compatible endpoint supplied the
selections.

| What was tested | Result |
|---|---|
| `sedum validate` on a bare-kwarg `injects_into` | 0 errors, 0 warnings |
| One `addImport` across two file groups | Both injected correctly |
| Two `addImport` calls into one file | Two distinct regions, no collision |
| Three consecutive runs | 2 injected, then 2 replaced, then 2 replaced — no drift |
| Target a path the record did not authorize | Phase 5 `unauthorized_path`, re-promptable |
| Target a file whose template plants no `imports` | Phase 7 hard error |
| Valid invocation *before* a failing one | **Nothing written.** Phase 7 is atomic |

So the mechanism needs no code to work, and the missing-anchor case is not a
safety hole — Phase 7 refuses the whole record rather than leaving a partial
write. The only defect is *when* it refuses: a Phase 7 halt is terminal where a
Phase 5 violation is re-promptable, which is exactly the shape
`prov-2026-9dcf2658` and `prov-2026-369544c1` both exist to close. **That
diagnostic move is the entire implementation cost of this section.**

### The real obstacle is a sealed constraint, not the code

`prov-2026-1bbb8e2e` — implemented, sealed — puts the target pattern in the
catalog, under a constraint that free targets contradict directly:

> The model still binds arguments and never chooses a path. The pattern is shown
> so that its bindings can land on an authorized file, not so that it can name
> one.

That record was written from an observed failure: qwen2.5-coder-14b, shown a
kwarg named `controller` and an authorized file
`app/controllers/users_controller.rb`, bound `controller` to the whole path. The
fix was to show the pattern so the model could match literal segments and bind
`users`.

**Worth sitting with: the model's "mistake" was this section's proposal.** It
tried to name the file it had been handed. `prov-2026-1bbb8e2e` and free targets
are two answers to one observation — teach the model to invert the pattern by
forward matching, or stop requiring the inversion. The first is strictly
necessary for structure-creating actions, where the path genuinely is the
package's. It is unnecessary work for file-agnostic ones.

Two consequences:

**The constraint needs rescoping, not overriding.** As written it is absolute.
What it should say is that a *pattern-targeted* action's kwargs are bound, never
inverted from a path — which is the failure it was actually built from. An action
that declares its target to be a kwarg has no pattern to invert and no inversion
to prohibit. Adopting free targets is a record amending that scope, and should be
recorded as one rather than smuggled in as a package-authoring convention.

**The catalog entry degenerates.** With `injects_into: "{{file}}"` the catalog
shows the model `["{{file}}"]`, which carries none of the information
`prov-2026-1bbb8e2e` added it to carry. It is not harmful, but it is noise, and
it suggests a free-target action should be presented differently — the anchor it
requires is the useful fact, not its target pattern.

### The alternative that keeps the model out of paths

A discriminated target map: a `target_kind` kwarg selecting among declared
patterns, mirroring `discriminator`/`variants` for the path instead of for the
template.

```yaml
addImport:
  kwargs:
    target_kind: { type: string, required: true }
    name:        { type: string, required: true }
  target:
    discriminator: target_kind
    paths:
      model:      "app/models/{{name|snake}}.rb"
      controller: "app/controllers/{{name|snake}}_controller.rb"
```

It preserves the "Sedum knows the layout, the model says what" property, keeps
transforms in play, and — the argument that got stronger once
`prov-2026-1bbb8e2e` was read — needs no constraint amended, because the model
still never names a path. Against it: it restates the `files/` layout in a second
place where the two can drift, it adds a second discriminator concept to a schema
that already has one, and it makes the model name a category when it could name
the file it was already handed.

Between the two, the free target is the smaller mechanism and the discriminated
map is the smaller governance change. That is the trade, and it is the decision
this section is actually asking for.

### Open

**Whether `prov-2026-1bbb8e2e`'s constraint is rescoped or left standing.** This
is the decision; everything else here follows from it. Rescoping it costs one
record. Leaving it standing means the discriminated target map is the only
available answer.

**Whether both forms coexist,** or whether a free target requires an opt-in in
`sedum.yaml`. A package author who wants a purely convention-driven layout may
reasonably not want any action taking a path from the model — the same grounds as
the opt-in guard on §4's edit layer.

**How a free-target action is presented in the catalog,** given that its
`injects_into` entry conveys nothing. Showing the required anchor in its place is
the obvious candidate and would be the first time the catalog carried an anchor
at all.

**Whether the anchor-applicability check is an error or a warning at Phase 5.**
An error is re-promptable and cheap, which argues for error. But `_default` file
templates and empty-file fallbacks mean a legitimately blank file plants no
markers at all, and that case should not be indistinguishable from aiming an
action at the wrong file.

---

## 16. Other unresolved items

**Removal.** If a record drops `destroy`, does replay delete the owned region? Saying yes turns recordings from action logs into declarative desired state and replay into reconciliation — powerful, converges toward idempotent sync. It also means an incomplete hand-edited recording silently deletes code. If pursued: `--prune` opt-in, never default.

**`forbidden_scope` becomes load-bearing** the moment Sedum edits an existing repository rather than generating into paths it created. It is nearly decorative under the MVP scope.

That moment arrives sooner than the sentence above assumes. It is not gated on Sedum gaining the ability to edit — it arrives as soon as **anything** working under a record edits an existing repository, and the edit layer in §4 is exactly that. Sedum enforces `forbidden_scope` over what Sedum does; nothing enforces it over what a tool above Sedum does except that tool. Whether the enforcement belongs there, or whether a record's forbidden paths should be checkable independently of who is acting, is open.

**Two-phase generation requires the recording `phases` grouping** already reserved in the PRD, with test and implementation as separate named phases.

**Spec-in-`affected_scope` creates a within-record ordering dependency** — spec actions must run before implementation actions so the rendered spec can inform the second pass. A hardcoded two-phase split handles it without becoming a planner, but it is the first real pressure on the no-planner position.
