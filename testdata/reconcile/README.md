# Phase 3 reconciliation fixtures

Inputs and expected outputs for `prov-2026-4c49ca46`: what Phase 3 does when an
authorized path already exists.

Each case directory holds `existing` (the file on disk), `template` (the file
template **as rendered**, never the template source), and — for a case that
applies — `expected`. A case that halts has no `expected`; `cases.json` carries
the diagnostic it must produce instead.

These are internal fixtures rather than a published surface, which is why they
live under `testdata/` and not beside `conformance/`. The marker corpus is
published because a third party writes that format; nothing outside Sedum
implements this.

## Where the inputs come from

The Rails inputs are real, not reconstructed. `rails-new-routes` is
`config/routes.rb` from `rails _7.2.3_ new --api`. The `rails-generate-*` cases
are rendered from railties' own `controller.rb.tt`. The rendered templates come
from Sedum itself, via `grow --stop-after files` against the packages in this
repository.

## What each case is evidence for

`cases.json` carries a `why` per case. Two are worth knowing before reading the
data:

- **`rails-generate-bare` applies and `rails-generate-actions` halts.** The
  boundary is not "controller templates cannot be adopted" — an empty generated
  controller adopts exactly, one carrying hand-written actions does not.
- **`private-with-method` halts because of the method, not the keyword.**
  `private-empty` is the same file without it, and applies.

## Expected outputs state required behaviour

They are not transcripts of a prototype. `private-empty` carries a blank line
between the `actions` anchor and `private` that the spike this record was
written from does not yet emit — the fixture is right and the implementation
has to meet it.
