# Generator packages — the house service standard, Rails edition

Sedum input. Point at it with `--generators ./_generators`.

The leading underscore keeps the templates out of everything that walks the
tree: `{name}_controller.rb` is not a filename Ruby can require, and an action
template is a bare method body with no class around it. Rubocop and Zeitwerk
both ignore a directory named this way.

These packages define **how every service is built**, not how this one is. The
stack, the file layout, the database, and the boilerplate in each file are the
standard; resource names, columns, and endpoints belong to the individual
service and arrive as bound arguments. The same directory generates
`class Order` in `app/models/order.rb` from the template that generates
`class Todo` in `app/models/todo.rb`. Nothing under `rails/` mentions either.

## The standard

**Layout.** Stock Rails: `app/controllers/<plural>_controller.rb`,
`app/models/<singular>.rb`, one migration per resource under `db/migrate/`,
minitest siblings under `test/`. One resource per file, named for the file.

**Runtime.** Rails 7.2 on Ruby 3.2.2, PostgreSQL through the `pg` adapter,
connection from `DATABASE_URL`. The web container starts only after compose
reports postgres healthy, so the application carries no retry loop for something
the orchestrator already knows how to express.

**Requests.** JSON in, JSON out, from callers that are not browsers. The
controller skips the CSRF token check once for the whole class rather than
per action, and renders `ActiveRecord::RecordNotFound` as a 404 through a single
`rescue_from`.

**Updates are partial.** `params.permit` drops what was not sent and
`ActiveRecord#update` assigns only what it is given, so an omitted attribute
keeps its stored value without any action arranging it. An empty body is a 400
rather than a write of nothing.

## One package, where the Go standard has four

The service next door splits into `chi`, `postgres`, `compose`, and `repo`
because a package declares one `comment_prefix`, and a marker planted with `//`
is not a SQL comment. Ruby and YAML both comment with `#`, so that boundary does
not exist here and one package claims `.rb` and `.yml` together.

| Package | Claims | Prefix |
|---|---|---|
| `rails` | `.rb` `.yml` | `#` |

Package boundaries follow comment syntax, not file type. A second package here
would be a directory split with nothing behind it.

## Rules this standard learned the hard way

**A file template owns the whole file or none of it.** `config/routes.rb` and
`config/database.yml` already existed — `rails new` wrote them. Phase 3 is
create-if-absent and verifies that an existing file carries the markers its
template declares, so a file Sedum is to inject into either gets its anchor
planted by hand or gets replaced by the template outright.

This standard takes the first path for `config/routes.rb`: `rails new` writes
that file, a service accumulates hand-written routes in it over its life, and a
generator that replaced it outright would be hostile to how the file is
actually used. Adopting it costs one comment line. The template still carries
Rails' own `/up` route because it is what renders when the file is genuinely
absent — owning the file means owning everything in it — and it is what
declares the marker Phase 3 verifies when the file is not.

**Rails' own routes belong in the template, not in an action.** They do not vary
between services, and an action that every record must remember to invoke is an
action some record will forget.

**A filename that carries data needs that data as a kwarg.** A migration is
`db/migrate/<stamp>_create_<table>.rb`, and `injects_into` is rendered from
bound kwargs, so `addColumn` takes `stamp`. It is the one place in this standard
where an argument names a file rather than describing the code going into it.
The record author picked the stamp when they wrote the path into
`affected_scope`; nothing downstream can rediscover it, because Sedum resolves
paths it is given rather than scanning for files that look close enough.

**The template pattern splits it back out.** `{stamp}_create_{name}.rb` binds
both, so the template writes `class Create{{name|models}}` without the timestamp
leaking into the class name. Captures are sub-segment, and two of them in one
segment are legal as long as literal text separates them.

**Repetition comes from invocations, not loops.** The grammar is `{{name}}` and
`{{name|transform}}` and nothing else. Two columns are two `addColumn`
invocations, not one template that iterates a list.

**An anchor no action targets is a trap.** Sedum warns when an action names a
marker no template plants; it cannot warn about the reverse, because a marker
with no action today may be a documented extension point tomorrow. Every anchor
in this package has an action that fills it. If you add one, add its action in
the same change.

**Unconditional wiring belongs in the template.** `skip_before_action`,
`rescue_from`, and `t.timestamps` are in file templates rather than actions
because every service in this standard has them and none of them varies. An
action would make them skippable, and a controller that forgot its `rescue_from`
returns a 500 where the contract says 404.

## What this standard deliberately does not generate

**LineSpec tests.** Never generated, and no record here authorizes a path under
`linespec/`. A linespec is a human-authored behavioral contract; a generated
contract asserts whatever the generator assumed. They remain as
`associated_specs` — gating the work without being subject to it.

**Minitest files are the opposite case** and are generated. They assert what the
standard promises: the status a verb returns, the 404 a missing record produces.
That is boilerplate. Note that Rails integration tests reach the database even
when asserting a 404, so they need a prepared test database — unlike the Go
standard's handler tests, which stay in-process by construction.

**The `Gemfile`.** Sedum resolves by extension and `Gemfile` has none, the same
way `Dockerfile`, `Rakefile`, and `config.ru` have none. Adding the `pg` gem is
a human step. This is not a limitation worth routing around: a dependency list
is a structured document, and appending text to one is a stated non-goal.

**`db/schema.rb`, `bin/`, `config/master.key`.** Written by running Rails, not by
describing it.

## The seam worth watching

`definePermittedParams` takes `attributes` as a string fragment — `":title,
:completed"` — and `addColumn` takes `options` the same way. Some of that is
data and some of it is Ruby, which is the same seam the Go standard flags around
its query kwargs.

It is here because the permitted list is one method whose shape depends on a
list, and the grammar has no loops. `options` is required with no way to omit
it, since a conditional would be a second grammar; `null: true` is the neutral
value written out rather than defaulted.

Worth revisiting before the model is in the loop: either accept fragments as
bound data, or generate one method per attribute and compose them.
