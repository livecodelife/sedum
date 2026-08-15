# Generator packages — the house service standard

Sedum input. Point at it with `--generators ./_generators`.

The leading underscore is load-bearing: the Go tool ignores directories whose
names start with `_`, and these templates are not valid Go — `{name}.go` is not
a legal filename to the compiler, and an action template is a bare function body
with no package clause. Without it, `go build ./...` at the repo root tries to
compile the generators and fails.

These packages define **how every service is built**, not how this one is. The
stack, the file layout, the database, and the boilerplate in each file are the
standard; resource names, columns, and endpoints belong to the individual
service and arrive as bound arguments. The same directory generates `type Order`
in `db/orders.go` from the template that generates `type Todo` in `db/todos.go`.
Nothing under `chi/` mentions either.

## The standard

**Layout.** Go sources in the repo root: `main.go`, `handlers/<plural>.go`,
`db/<plural>.go`, each with a `_test.go` sibling. One resource per file, named
for the file.

**Runtime.** Chi v5, lib/pq, `DATABASE_URL` from the environment. Startup waits
for the database — 30 attempts, one second apart — then fatals. A container
starting beside postgres loses the race, and a service that cannot reach its
database has nothing to serve.

**Data.** Postgres 16-alpine, md5 auth, host port 5433. Every table carries
`id`, `created_at`, `updated_at`; the model struct carries the matching fields
with an anchor between them for the service's own.

**Query patterns.** Two-step INSERT (`RETURNING id`, then a separate SELECT).
UPDATE via `Exec` with a `RowsAffected` check, then a re-read, with `id` bound
to `$1` so assignments number from `$2`. LIST via `PrepareContext` +
`stmt.QueryContext`. `db.ErrNotFound` is what a missing row returns, so handlers
map 404 without importing `database/sql`.

## Four packages, because comment syntax draws the boundary

A package declares **one** `comment_prefix`, and a marker planted with `//` is
not a SQL comment. Since anchors exist only because a file template planted
them, that decides where one package ends and the next begins.

| Package | Claims | Prefix |
|---|---|---|
| `chi` | `.go` | `//` |
| `postgres` | `.sql` | `--` |
| `compose` | `.yml` | `#` |
| `repo` | `.gitignore` `.mod` `.sum` | `#` |

One record legitimately spans all four, and each path is generated under its own
package's conventions. That is what per-file resolution is for.

## Rules this standard learned the hard way

**A file template must use every import it declares.** Go rejects an unused
import, so a template that imports on behalf of code not yet injected produces a
file that will not compile — and it stays broken if the model never selects the
action that would have used it. Every import here is used by a helper the
template itself declares: `pathID`, `writeJSON`, `decodeBody` in handlers;
`Ready` and `ErrNotFound` in db; `request` and `run` in handler tests;
`requireConn` in db tests. Injected code calls those helpers rather than
repeating them, which is why they earn their place twice.

**An `imports` anchor belongs inside the import block.** Put it after the
closing paren and the first injected import is a syntax error. Nothing catches
this before injection: Sedum verifies a marker is *present*, not that it sits
somewhere an injection would be valid.

**Imports are injected one at a time, and the path is data.** A file template
declares the imports the standard always needs and stops there. It cannot name
the service's own packages, whose paths begin with a module name no template
knows, and it cannot anticipate a dependency some later invocation brings in. So
every file that might grow an import opens its import block with an anchor, and
`addMainImport`, `addHandlerImport`, `addDBImport`, `addHandlerTestImport` and
`addDBTestImport` each take one `path`.

One action per file, not one action taking the file. Which files a service has
is the standard's business and belongs in configuration; which package to import
is the service's business and is an argument. An action that took the path would
hand both decisions to the caller.

Identity is the import path, so invocations accumulate — two paths are two
regions — and re-running with the same path rewrites in place. Repetition comes
from invocations here exactly as it does for fields and queries.

The service's own packages are then just ordinary invocations:
`addMainImport path=<module>/db`, `addMainImport path=<module>/handlers`,
`addHandlerImport path=<module>/db`. They are the part a record cannot skip —
every other action adds a feature, while without these the output does not
compile at all. An anchor no action targets is invisible to Sedum today: Phase 0
warns about an action naming a marker no template plants, not about a template
planting a marker no action fills. This gap was found by building a generated
service, not by validating the package.

**An injected import must be used.** Go rejects an unused import, so an
invocation that adds a path nothing references breaks the build as surely as a
missing one does. Nothing in the pipeline can catch it — whether a package is
used is a property of the whole file after every injection lands.

**Wiring belongs in the file template when it is unconditional.** `main.go`
calls `db.SetConn` itself rather than through an action, because every service
in this standard wires exactly one pool exactly once. An action would make it
skippable, and a service whose pool is never wired compiles and then panics on
its first query, which is worse than failing to build.

**Repetition comes from invocations, not loops.** The grammar is `{{name}}` and
`{{name|transform}}` and nothing else. A model with five fields is five
`addModelField` invocations against the `fields` anchor, not one template that
iterates. Design actions at the granularity of one repetition.

**Generated Go is gofmt-dirty.** `gofmt` aligns struct tags across sibling
fields, which is a whole-block property, and each field is injected
independently with no knowledge of its siblings' widths. The result compiles,
vets, and tests clean — it just is not gofmt output. Run `gofmt -w` afterwards.

## What this standard deliberately does not generate

**LineSpec tests.** Sedum never generates them, and no record here authorizes a
path under `linespecs/`. A linespec is a human-authored behavioral contract; a
generated contract asserts whatever the generator assumed. The linespecs in this
repo remain as `associated_specs` — they gate the work without being subject to
it. A team *could* write a package claiming `.linespec`; it would be bad
practice.

**Unit tests are the opposite case** and are generated. They assert what the
standard promises — the status code a verb returns, the error a missing row
produces — which is boilerplate. Handler tests cover only paths that need no
database, since that is all a generated test can honestly assert.

**`go.mod` and `go.sum`.** Claimed but templateless: the module path is
per-service and the checksum file belongs to the toolchain. Sedum creates them
empty and logs it. Create-if-absent means running `go mod init` first leaves them
alone.

**`Dockerfile`.** Cannot be generated — Sedum resolves by extension and it has
none. Same for `Makefile`, `Procfile`, `LICENSE`. These are human-authored
anyway, so no record here names one. Dotfiles are fine: Go reads `.gitignore` as
having the extension `.gitignore`.

## The seam worth watching

`createQuery` takes six string kwargs — `columns`, `placeholders`,
`insert_values`, `set_clause`, `update_values`, `scan_targets`. Some are data (a
SQL column list); `insert_values` binds to something like `in.Title,
in.Completed`, which is a Go fragment. The PRD says the model binds arguments and
does not write code, and six fragments is a lot of weight on that line.

It is here because a query is one statement whose shape depends on a list, and
the grammar has no loops. Worth revisiting before M6: either accept fragments as
bound data, or scaffold only the db layer's shell and leave query bodies to a
human or to the focused-edit layer.

They are declared optional because one kwarg schema covers every variant of a
discriminated action and `delete` needs none of them. Per-variant requirements
cannot be expressed today.

## Running it

One file, one record: `main.go` is named by both `prov-2026-4317fe6c` and
`prov-2026-addce7fc`, and `handlers/todos.go` by both `4317fe6c` and
`e36db3d5`. Sedum makes one model call per record, so this needs `--only` until
ownership markers carry the record ID.

```
sedum resolve --generators ./generators --records ./provenance \
              --only prov-2026-4317fe6c --show-template

sedum grow    --generators ./generators --records ./provenance \
              --only prov-2026-4317fe6c --output ../build --stop-after files
```

`--stop-after files` is as far as Sedum goes today (M3): correct skeletons with
anchors planted and empty regions. Everything under `actions/` is rendered and
injected by M4. The skeleton compiles as-is, and with the action templates
injected by hand it builds, vets, and passes `go test ./...`.
