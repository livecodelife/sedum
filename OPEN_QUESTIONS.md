# Sedum — Open Design Questions

Scratch document. Nothing here is a requirement. These are directions explored after the MVP scope was settled, recorded so the reasoning isn't lost, and deliberately kept out of the PRD because most of them depend on answers M1–M7 will produce.

Two decisions from this exploration were pulled forward into the PRD, because both write durable artifacts whose shape is expensive to change later: **ownership tiers** (`owned` / `seeded`) and **recorded kwargs** on markers, and the **`phases` grouping** in recordings. Everything else below is open.

Sections 7–14 come from a later session on imprint integration, cross-region logic, and the record-authoring workflow. Several of those threads reached a settled position rather than staying open; they are marked as such where they did.

---

## 1. Generated tests as a first pass

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

Specs as **selection evidence** is the cheapest available accuracy lever and requires no execution: a spec asserting `GET /users → 200`, `POST /users → 201`, `DELETE /users/:id → 204` tells the model which methods to generate far more reliably than prose intent does. Sedum never parses it — for prompt purposes any spec is text, which preserves the no-target-knowledge principle.

Specs as **signal** — run once after generation, report, never retry. Boilerplate mostly won't pass behavioral specs and shouldn't be expected to, but a pass count is useful information and one run is not a verify loop.

Specs as **generator-package regression fixtures** is the most interesting use. Record the pass set at record time; compare on replay. A scaffold that passed 3 and now passes 2 means a template regressed. This makes recordings the test suite for generator packages, which is otherwise hard to test.

`run_command` is arbitrary shell — same trust posture as user-defined transforms. Gate behind an explicit flag, never let it block.

---

## 3. Updates to existing code

### The easy cases

**Additive change to a Sedum-generated file** is nearly free: Phase 3 becomes create-if-absent, verify-markers-if-present. A record saying "add destroy" resolves to `createControllerMethod(name: destroy)` injected at the existing marker. This likely covers most real update traffic.

**Additive change to a file Sedum didn't create** has no markers. Hard-error, and offer an `adopt` path that inserts markers only where the file's skeleton still literally matches a file template. Anything fuzzier needs structure inference, which means a parser, which is the thing avoided everywhere else.

### Changes inside the parameter space

A surprising amount of update traffic is a kwarg change. "Include addresses in the index response" is `render json: @users, include: :addresses` — if the `index` variant accepts an `includes` kwarg, the update is a re-invocation against an owned region. Fully deterministic, no code editing at all.

This is why kwargs are recorded on the marker: current state becomes readable from the file, and desired-vs-actual becomes a diff over kwargs.

### Changes outside the parameter space

"Exclude soft-deleted users unless the requester is an admin" is business logic. The useful behavior at this boundary is **refusing well**: detect that the record targets a `seeded` region or requires kwargs the template doesn't accept, and emit a work item naming the file, the region, and the constraint.

---

## 4. The focused edit layer

A work item from the boilerplate layer is a much better-conditioned task than open-ended synthesis. It carries a file that already exists with correct structure, a named region with known bounds, current contents that are usually close to right, a specific constraint describing the delta, and the marker kwargs documenting what the region was parameterized to do. That is a scoped edit against working code, not synthesis — a different difficulty class, and plausible for a local open-weight model in a way the general problem is not.

### Design positions

**Structured edits, not region rewrites.** Ask for a replacement region and the model rewrites the whole thing with incidental drift. Ask for a small list of `{after, insert}` / `{replace, with}` operations validated against current text, and the blast radius is bounded, invalid edits are rejectable without execution, and the human reviews the actual change. Same structured-output discipline as Phase 5, one layer down.

**Edited regions change tier.** A region the edit layer touches can't stay `owned` — the next replay would destroy the edit. It becomes `seeded`, or a third state recording both the originating invocation and the modifying record.

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

## 15. Other unresolved items

**Removal.** If a record drops `destroy`, does replay delete the owned region? Saying yes turns recordings from action logs into declarative desired state and replay into reconciliation — powerful, converges toward idempotent sync. It also means an incomplete hand-edited recording silently deletes code. If pursued: `--prune` opt-in, never default.

**`forbidden_scope` becomes load-bearing** the moment Sedum edits an existing repository rather than generating into paths it created. It is nearly decorative under the MVP scope.

**Two-phase generation requires the recording `phases` grouping** already reserved in the PRD, with test and implementation as separate named phases.

**Spec-in-`affected_scope` creates a within-record ordering dependency** — spec actions must run before implementation actions so the rendered spec can inform the second pass. A hardcoded two-phase split handles it without becoming a planner, but it is the first real pressure on the no-planner position.
