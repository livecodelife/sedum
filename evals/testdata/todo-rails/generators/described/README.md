# Generator packages — the described arm

The `defined/` set with descriptions added, and nothing else changed. It is one
half of prov-2026-ac15ed2b: run the harness over todo-rails twice, once with a
package whose actions and kwargs carry `description` fields and once without,
and record whether the selection rate moved.

The standard itself — the layout, the runtime, the rules this package learned
the hard way — is documented once, in [`../defined/README.md`](../defined/README.md).
Nothing here restates it, because a second copy of that prose is a second thing
to keep true.

## What differs, and what enforces it

Only `rails/actions/actions.yaml`, and only by added `description` fields. Every
other file in the set is a byte-for-byte copy: the manifest, the eight action
templates, the six file templates.

`TestTheDescribedSetDiffersOnlyByDescriptions` in `evals/run_test.go` loads both
sets and compares them field by field with descriptions excluded. A renamed
kwarg, a flipped `required`, a new variant, a changed `injects_into`, an edited
template — any of those fails the default test suite.

That test is the guard, not a rule about who may edit this directory. A
hand-written second package drifts, and a drifted one measures something nobody
chose; a diff catches that and a convention does not.

## The descriptions are the comments, moved

Every description here says what a YAML comment in `defined/` already said. That
is the whole point of the feature under test: prov-2026-c5697387 was written
from four wrong bindings whose answers were each sitting in a comment beside the
declaration, where the model cannot see them. The comments are kept in both
sets, so what varies between the arms is only what reaches the catalog.

## What a result here does and does not mean

This measures **selection** — whether the model picks the right actions the
right number of times. It does not measure binding, which is what descriptions
were added for. Every wrong binding behind prov-2026-c5697387 produced a
correctly *selected* action with a bad *argument*, and would have scored a
perfect catalog count. Reading a moved rate here as validation of that record's
feature is the specific mistake to avoid.
