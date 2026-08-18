# Behaviour

Answers a question the rest of the eval cannot: **does what Sedum generated
actually run?**

Every other number here is about the answer — which actions were selected, what
their arguments were bound to. None of them is about the application. A sample
can select perfectly, bind correctly, and produce a service that violates the
record it was generated from, because what the generator package can express is
not something the model has any say in.

This scaffolds an empty project the way the framework's own generator does,
applies a selection to it with Sedum, builds it, boots it, and asserts against
it over HTTP. The project is temporary and is deleted on the way out.

## Running it

```
go run ./evals/cmd/eval -behavior todo-rails-defined     # as part of a measurement
```

Or directly, from `evals/`:

```
bash behavior/behave.sh todo-rails                       # the target's own answer.json
bash behavior/behave.sh todo-rails --answer /tmp/s.json  # a specific selection
bash behavior/behave.sh todo-rails --files /tmp/src      # the baseline arm
bash behavior/behave.sh todo-rails --model qwen2.5-coder-14b-instruct
bash behavior/behave.sh todo-rails --keep                # leave the project behind
```

`--files` is the baseline arm: the sources are already written and `generate`
copies them in rather than running Sedum. It is exclusive with `--answer`,
because a run either used Sedum or is the measurement of not having used it.

A measurement drives the `--answer` form, one call per valid sample, with each
sample's own invocations. That is what makes a behaviour rate a statement about
what the model chose rather than about a list somebody checked in — a target's
`answer.json` is a hand-authored assertion of what a correct selection looks
like, useful for exercising the harness itself.

## Shape

Six phases, run in order, stopping at the first that fails:

| Phase | What it is |
|---|---|
| `scaffold` | `rails new` / `go mod init` — the framework's own generator |
| `prepare` | Delete the paths the record authorizes, so Sedum's file templates own them |
| `generate` | `sedum grow`, built from the tree under test — or, under `--files`, a copy of sources Sedum did not write |
| `build` | `bundle install` / `go build` — the dependency handoff a person does |
| `boot` | Create a database, migrate or load the schema, start the server, wait for it |
| `verify` | One assertion per constraint in the record, then the generated test suite |

`behave.sh` is the runner and knows nothing about any framework. Everything
stack-specific is `targets/<name>/target.sh`, defining those six functions plus
`target_teardown`. **Adding a stack is a new directory** — not a change to the
runner, the case schema, or the report.

The baseline arm changes no target for the same reason. Five of the six phases
work on a directory of source and have no opinion about how it got there —
`verify` in particular is HTTP assertions against a running service, which is
the same question whoever wrote the code.

`prepare` matters more than it looks. A scaffold's `config/routes.rb` carries no
Sedum anchor, and Phase 3 is create-if-absent, so a run against an untouched
scaffold either fails its marker check or injects into nothing. The list comes
from the record's `affected_scope` rather than being written out again.

## Outcomes

Three, kept apart on purpose:

| Outcome | Meaning |
|---|---|
| `ok` | Every assertion held |
| `checks_failed` | It built, booted, answered — and disagreed with the record |
| `failed` | A phase died, so the assertions never ran |

A service that never booted and one that booted and answered wrongly are
different findings with different fixes. A single pass rate would report a
broken generator package and a wrong one as the same number.

A `failed` run carries the tail of the dead phase's log, in the result file and
in the eval's report. Without it a failure said "build" and the log that said
why went with the temporary project, which cost three hand reconstructions in
one session — and a reconstruction is not the sample that failed. Use `--keep`
when the tail is not enough.

## The stub model

`stub_model.py` is an OpenAI-compatible endpoint that always returns the same
answer. Applying a recorded selection means running `sedum grow` without asking
a model, and `grow --execute` — designed for exactly this — is not built (M7).

**It is a stand-in.** When `--execute` lands, delete this rather than keeping it
as an alternative: two ways to replay a selection is one more than the property
being measured can survive.

## What it costs

About 20 seconds per sample on a warm machine — 7s scaffold, 10s generate, 3s
boot, 2s verify. A cold `bundle install` adds a couple of minutes to the first
Rails run. That is why `-behavior` is off by default, and why the wrapper's
derived timeout budgets extra time per sample when it is on.

## Requirements

A PostgreSQL the current user can `createdb` on (Homebrew `postgresql@14` on
5432 is what this was built against), Ruby 3.2 with Rails 7.2 for the Rails
target, a Go toolchain for the Chi one, plus `jq` and `python3`.

No Docker. The generators' `docker-compose.yml` is not in either record's
`affected_scope`, so nothing generated needs it, and a local server boots in
three seconds rather than building an image.

## What it found

- **The todo-chi package could not satisfy three constraints of its own
  record** — no validation action, no way to express partial update, nothing
  emitting a 400 for an empty body. The model selected perfectly and the service
  was still wrong, and no selection could have fixed it. Both records were cut
  to what both stacks can express (`prov-2026-40bdd9ac`); restoring those
  constraints means giving the Chi package the actions it lacks.
- **The Chi service did not compile at all** until the module path stopped being
  something the model was asked to guess (`prov-2026-6fc3d13d`). It bound the
  five import actions to chi, sqlx, net/http, net/http/httptest and testify —
  four unused, one already imported — because go.mod is in no record and no file
  list.
- **A migration with `default: ""` on a column the record says carries no
  default passed twenty-one HTTP assertions.** No HTTP call can observe a
  column default on a `NOT NULL` column the caller always supplies. The schema
  assertion in the Rails target exists because that was missed.
- **The Rails controller template rules out `rails new --api`** — it calls
  `skip_before_action :verify_authenticity_token` unconditionally, and
  `ActionController::API` has no such filter, so the class raises at definition.

## Known rough edges

- `answer.json` is hand-authored, so it is an assertion about what a correct
  selection looks like and should be reviewed as one.
- The Rails target appends `gem "minitest", "~> 5.25"` to the scaffold's
  Gemfile. A fresh 7.2 scaffold resolves minitest 6, whose runner signature
  railties 7.2 does not match, so `rails test` dies before running a test. A
  scaffold problem rather than a Sedum one, but it has to be handled somewhere.
- The Chi target hardcodes port 8080 because its `main.go` file template does,
  so two Chi runs cannot go at once. Rails picks a free port and is safe under
  concurrency.
