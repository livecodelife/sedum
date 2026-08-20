# Sedum

Sedum generates boilerplate code from provenance records using generator packages that teams author themselves.

A provenance record declares intent, constraints, and the files a change is authorized to touch. A generator package declares a team's conventions: the shape of each file type, which code-injection actions exist, what arguments they take, what templates they render, and where in a file the results belong. Sedum reads both, uses a language model to translate intent into a list of action invocations, and executes those invocations deterministically.

**Sedum's core contains no language-specific knowledge.** It knows how to read a directory, match a path against a pattern, apply string operations, render a template, and place text at an anchor. It does not know what a controller is, what a header file is, or how Go capitalizes initialisms. Teams express all of that in configuration. A team switching from Rails to Java writes a new generator package; they do not wait for Sedum to add support.

---

## Evaluation results

Measured runs of the harness under `evals/` — what a model selected, what it
bound, and whether the application it produced builds, boots and answers — are
published at **<https://livecodelife.github.io/sedum/>**, read directly from
`evals/results/`.

---

## Status

Sedum is under active development. **The whole pipeline has landed**, model
invocation included: package loading and validation, transforms and template
rendering, record ingestion, path resolution, file creation and reconciliation,
model invocation with validation and re-prompting, composite expansion, and
anchored injection with ownership markers.

The deterministic half was built first, deliberately. The part that does not
involve a model had to be provably correct before one was introduced, so that a
failure now is unambiguously a selection failure rather than a rendering
failure. That property is what the harness under `evals/` measures.

| Phase | State |
|---|---|
| 0 — Load and validate generator packages | Working |
| 1 — Ingest provenance records | Working |
| 2 — Resolve paths to packages and templates | Working, including unmanaged paths |
| 3 — Create files from templates | Working, including reconciling a file that already exists |
| 4 — Model invocation | Working |
| 5 — Validate the model's output | Working, including re-prompting to the retry limit |
| 6 — Expand and resolve | Working, including composite expansion |
| 7 — Inject | Working |

| Command | State |
|---|---|
| `sedum grow` | Working, end to end and at every stop point |
| `sedum validate` | Working |
| `sedum resolve` | Working |
| `sedum actions` | Working |
| `sedum render` | Working |

Recording and replay are live: `--record` captures what a run resolved and what
the model decided, and `--execute` replays that with no model involved. A
recording is plain JSON and hand-editing one is a supported workflow rather than
a trick — see [Recorded Executions](#recorded-executions).

Two limits are worth knowing before pointing this at a real project.

A path named by two provenance records is rejected, now at Phase 4's entry
rather than at ingestion. Phase 4 makes one model call per record, so two
records naming one file would mean two independent calls deciding one file's
contents. Under `--execute` there is no model call, so the check does not apply
and a recording may legitimately cover a shared file. Records that share a file
generate correctly one at a time with `--only`, and their regions coexist and
survive each other's reruns.

And a model's selection is the one part of a run that is not reproducible by
construction. Temperature is pinned to zero, but the open-weight range this is
built for still varies between runs — dropping an invocation, adding a
defensible extra, or binding a kwarg in a form the package did not intend. That
is the failure mode recordings exist for: capture once, review the JSON, correct
it by hand, and replay it forever.

---

## Building and running

Requires Go 1.25 or newer.

```sh
go build -o sedum ./cmd/sedum
```

Run the test suite:

```sh
go test ./...
```

### Commands you can run today

**Validate the example generator packages** — runs every load-time check against a generators directory and exits. Requires no records, no model, and no network.

```sh
./sedum validate --generators ./testdata/generators
```

```
3 package(s) loaded, 0 error(s), 0 warning(s)
```

Validate a single package, and treat warnings as failures:

```sh
./sedum validate --generators ./testdata/generators --package rails --strict
```

**Bind a project fact a package cannot know** — a generator package may declare
`variables:` for things like a Go module path or a .NET root namespace. The
package declares the name; the run supplies the value, and Sedum never learns
what it means. A declared variable with no value and no default stops the run
before anything is written.

```sh
./sedum grow --generators ./generators --records ./provenance \
             --var module=github.com/acme/todo
```

**See what each authorized path resolves to** — the generator package that
claims it, the file template that matched, and the captures that template bound.
Writes nothing.

```sh
./sedum resolve --generators ./testdata/generators --records ./your-records
./sedum resolve --generators ./testdata/generators --records ./your-records --show-template
```

```
prov-2026-aaaaaaaa
  app/controllers/users_controller.rb
    package   rails
    template  app/controllers/{name}_controller.rb
    captures  name=users
```

**Scaffold the files a record authorizes** — every path created with its
boilerplate rendered and its anchors planted, then stop.

```sh
./sedum grow --generators ./testdata/generators --records ./your-records \
             --output ./build --stop-after files
```

Rerunning is safe and is the ordinary way to resume: a path that already
carries everything its template declares is left exactly as it is, and one that
does not is reconciled to it. Injected regions are never overwritten by Phase 3.

Add `--dry-run` to decide everything and write nothing, or `-v` to mirror the run
log to stdout.

**Run the whole thing** — the same command without a stop point, which invokes a
model. The endpoint and credentials come from the environment
(`OPENAI_BASE_URL`, `OPENAI_API_KEY`), so a local server and a hosted API are
the same code path.

```sh
./sedum grow --generators ./generators --records ./provenance \
             --model qwen2.5-coder-14b-instruct --record ./scaffold.json
```

**Replay it with no model at all** — deterministic, free, and the way to
re-generate a reviewed scaffold:

```sh
./sedum grow --generators ./generators --execute ./scaffold.json --output ./build
```

**See what the model will be shown** — the exposed actions of a package with
their kwarg schemas and variant lists, in the form Phase 4 hands over. The
authoring loop for catalog quality:

```sh
./sedum actions --generators ./testdata/generators --package rails
./sedum actions --generators ./testdata/generators --package rails --all --json
```

**Ask where an invocation would land** — an action's target rendered against
kwargs, by the package's own transform engine, so the answer agrees with what a
run would do instead of approximating it:

```sh
./sedum render --generators ./testdata/generators --package rails \
               --action createControllerMethod --kwargs '{"controller":"person"}'
```

**Inspect the command surface** — every command, its flags, and its documented behavior:

```sh
./sedum --help
./sedum grow --help
```

`testdata/generators` holds three worked example packages covering file
templates, discriminated actions, composites, declarative transform pipelines,
and exception tables. `rails` and `chi` are the familiar cases. `cairn` is a
package for a target that does not exist — its own extension, comment prefix,
pipeline vocabulary, and directory shape, sharing none of them with the others.
It is there so that "no language knowledge in the core" is a thing the test
suite checks rather than a thing the README claims.

---

## Design Principles

**Configuration carries all target knowledge.** Adding a language or framework means authoring a generator package, never modifying Sedum.

**The model does one bounded job.** It selects actions from a closed, declared vocabulary and binds their arguments. It does not write code, choose file paths, or make structural decisions. Selection over a declared schema is machine-checkable; open-ended synthesis is not.

**Everything after the model's response is deterministic.** File creation, template rendering, transform resolution, composite expansion, path resolution, and injection are pure functions of the model's validated output and the generator package.

**Sedum touches only what a provenance record authorizes.** No file is created that `affected_scope` did not name. No convention rule infers companion files into existence.

**Failures are loud and specific.** A missing anchor, an undefined transform, or an unauthorized path halts the run with a diagnostic naming the action, the file, and the rule violated.

---

## The Generators Directory

The user points Sedum at a generators directory. Its top-level subdirectories are **generator packages**, one per target stack.

```
generators/
  rails/
    sedum.yaml
    files/
    actions/
  chi/
    sedum.yaml
    files/
    actions/
  react/
    sedum.yaml
    files/
    actions/
```

Package directories are named by the author. The name is a label, not a language identifier — `rails` and `sinatra` are both Ruby with different conventions, as `chi` and `echo` are both Go.

Each package declares which file extensions it claims:

```yaml
# generators/rails/sedum.yaml
name: rails
extensions: [.rb, .erb]
comment_prefix: "#"
```

### Package resolution

Sedum builds an extension-to-package map when the generators directory loads. Every path in a provenance record's `affected_scope` resolves to a package by its extension.

**Resolution is per file, not per run.** A single record may legitimately touch `app/controllers/users_controller.rb` and `app/javascript/users.ts`; those paths resolve to different packages and each is generated under its own conventions.

If two packages in the directory claim the same extension, that is not an error at load time. It becomes an error only when a path with that extension appears and no `--lang` flag disambiguates it. This keeps multi-package directories legal while failing loudly at the point of real ambiguity.

The `--lang` flag names a package to prefer. It may be repeated. Sedum infers which extension each named package resolves; a flag naming a package that claims no extension in the current record set is a warning, not an error.

A path whose extension no package claims is a hard error. Sedum does not guess.

### Unmanaged paths

A package may declare paths it does not write:

```yaml
unmanaged:
  - Gemfile
  - config/credentials/
```

Entries use the same grammar as a record's `affected_scope` — `**`, `*`, `?`, `[…]`, and a trailing slash for a subtree. An authorized path matching one is skipped, reported, and the run continues: not created, not rendered, not injected into.

This exists so a record can describe the whole change. Moving a Rails service to PostgreSQL means adding the `pg` gem, so `Gemfile` belongs in `affected_scope` — and `Gemfile` has no extension, so without a declaration the run halts on it. The alternative is trimming records to what a generator happens to reach, which makes the record a description of the tool rather than of the work.

**It authorizes nothing.** `affected_scope` still decides what may be touched; this says only that Sedum is not what touches it. Usually a person is, or another tool pointed at it separately, so a run gathers these paths into a list rather than passing over them silently:

```
1 path(s) left unmanaged, for a person or another tool:
  Gemfile (rails declares "Gemfile")
```

An unclaimed extension nobody disowned is still a hard error. Silence is not a declaration.

### Test paths

A package may also declare which of its paths hold tests:

```yaml
test_paths:
  - "**/*_test.go"
  - spec/
```

Same grammar as `unmanaged` and a record's scope entries. Patterns rather than a
directory, because Go puts `user_test.go` beside `user.go` — same extension,
same package, same directory, and nothing directory-shaped can separate them.

**Sedum validates this and does not read it.** That is the whole of the field
and it is deliberate. The consumers are a caller partitioning a record's regions
into assertions and logic, and Sedum's own possible future two-pass generation;
neither is implemented, and an invalid entry is still rejected at load.

It lives in `sedum.yaml` rather than in a caller's own configuration because
`sedum.yaml` is decoded strictly. A package declaring a key Sedum has not
modelled does not load at all — not "loads and is ignored", but rejected with
every action in it unavailable. So a tool told to read `sedum.yaml` directly
cannot ask a package author to declare anything Sedum has not declared first,
whoever ends up doing the reading.

---

## File Templates

Creating a file empty is rarely useful. A Rails controller needs its class definition; a Go handler needs its package declaration and imports; a C++ header needs its include guard and class block. File templates supply that boilerplate.

A package's `files/` directory is a **literal mirror of the target project's structure**, with captures in the path segments:

```
generators/rails/files/
  app/controllers/{name}_controller.rb
  app/models/{name}.rb
  config/initializers/{name}.rb
  _default.rb
```

An authorized path is matched against this tree. `app/controllers/users_controller.rb` matches the first template, capturing `name=users`, and renders:

```erb
class {{name|constantize}}Controller < ApplicationController
  # sedum:anchor:class_body_top

  # sedum:anchor:class_body
end
```

The template's own path is the pattern. There is no mapping table, no separate registry of file types. Captured segments are available to the template as bound values, and transforms apply to them exactly as they do in action templates.

File templates are **not part of the action catalog**. They have a different lifecycle — invoked by path match during file creation, never selected. The model does not see them and cannot reference them.

### Match specificity

Multiple templates may match a path. Ranking, applied leftmost segment first:

1. Literal segment beats capture
2. Capture beats glob
3. Longer literal prefix beats shorter

If two templates tie under this ranking, the package is rejected at load. Sedum does not silently pick one.

### No match

If no template matches, Sedum uses the package's `_default` template if one exists. If not, it creates the file empty and writes a log line.

This is deliberately not a hard error. Not every authorized path needs boilerplate — a SQL migration or a plain config file may legitimately start blank, and forcing authors to declare a template for every extension would be noise.

### Anchors and marker planting

File templates plant the markers that action anchors target. This closes the loop: injection points exist because a file template created them, not because Sedum parsed the file for structure. It is the reason no language parser is required anywhere in the system.

At load time, Sedum checks each action's marker-based anchor against the markers present in its package's file templates. Because template selection is path-dependent, this cannot be a complete verification — but an action referencing a marker that appears in no template in its package is almost certainly a typo. This produces a warning.

---

## Action Definitions

### Simple actions

```yaml
actions:
  addBeforeFilter:
    kwargs:
      controller: { type: string, required: true }
      filter:     { type: string, required: true }
      only:       { type: list,   required: false }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body_top

  createControllerMethod:
    kwargs:
      controller: { type: string, required: true }
      name:       { type: string, required: true }
      collection: { type: string, required: false }
    discriminator: name
    variants: [index, show, create, update, destroy]
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body
```

`kwargs` is the schema the model is held to. Types come from a closed set — `string`, `int`, `bool`, `list` — sufficient for argument binding and nothing more.

`discriminator` names the kwarg whose value selects a template. `variants` enumerates the values that have dedicated templates. Both are declared explicitly rather than inferred from directory structure, so that specializing on a second argument later cannot create an undocumented precedence rule, and so a misspelled variant filename fails at load rather than silently falling through to `_default`.

`injects_into` is a path pattern rendered against the bound kwargs. It must resolve to exactly one file.

**An action may declare its target to be a kwarg instead.** `injects_into: "{{file}}"` with a `file` kwarg makes the target a parameter the caller aims, rather than a convention the action encodes:

```yaml
  addImportTo:
    kwargs:
      file: { type: string, required: true, description: "The file to import into. It must be one of the authorized files, and it must carry an imports injection point." }
      path: { type: string, required: true, description: "The import path, without quotes." }
    injects_into: "{{file}}"
    anchor: imports
```

Both forms are legitimate and the choice is the author's. For a structure-creating action the target is implied by the action's identity — `createControllerMethod` only ever means a controller, and the pattern is exactly right. For a file-agnostic action the file is genuinely a parameter, and encoding it in the action definition forces one action per file group: `addImport` becomes `addModelImport`, `addControllerImport`, `addSerializerImport`, with identical kwargs and identical template bodies differing only in a path.

The authoring duplication is the visible cost. The real cost is catalog quality — Phase 4 would hand the model a set of near-identical entries whose only discriminator is a name encoding a fact it has to infer, which is the worst possible shape for selection accuracy.

What does not change is that the model binds arguments and never inverts a path. A pattern-targeted action's kwargs are bound, never derived from a file name.

**Applicability comes from the anchor.** A pattern-targeted action can only reach the files its pattern describes, so where it applies is implied. A free-target action can be aimed anywhere, so the anchor states it instead: an invocation aimed at a file that does not carry the action's injection point is a Phase 5 error, re-promptable, rather than a Phase 7 crash. The prompt lists each authorized file with the anchors it carries, built from the same reading as the check, so the catalog cannot describe a file differently from the rule that will judge what the model does with it.

### Action template layout

**A simple action is a file or a directory, depending on whether it declares variants.** With no discriminator, its template is a single file named for the action. With a discriminator, its templates live in a directory named for the action, one file per variant value, plus an optional `_default` fallback.

```
generators/rails/actions/
  actions.yaml
  createControllerMethod/     # directory form — action has variants
    index.rb
    show.rb
    create.rb
    _default.rb
  addBeforeFilter.rb          # file form — action has no variants
  createModelClass.rb
```

The schema determines the shape. Sedum resolves the declared shape and errors if the filesystem disagrees; it does not infer an action's kind from what it finds on disk.

A variant template for `index`:

```erb
def index
  {{collection|instantize}} = {{collection|constantize}}.all
  render json: {{collection|instantize}}
end
```

**A composite action is neither, and triggers no filesystem lookup.** A composite has no template. It exists only in `actions.yaml` and resolves to an ordered list of simple actions. Searching the filesystem for a composite is a bug, not a fallback.

### Composite actions

A composite bundles simple actions that are always performed together. The motivating case is a language where one logical change spans two files — a C++ method requires a declaration in the header and a definition in the implementation file, and neither is valid alone.

```yaml
actions:
  createMethod:
    composes:
      - addMethodDeclaration
      - addMethodDefinition
    exposed: true

  addMethodDeclaration:
    kwargs:
      class: { type: string, required: true }
      name:  { type: string, required: true }
      args:  { type: list,   required: false }
    injects_into: "include/{{class|snake}}.hpp"
    anchor: class_body
    exposed: false

  addMethodDefinition:
    kwargs:
      class: { type: string, required: true }
      name:  { type: string, required: true }
      args:  { type: list,   required: false }
    injects_into: "src/{{class|snake}}.cpp"
    anchor: end_of_file
    exposed: false
```

Composites nest exactly one level. A composite may not compose another composite, and may not compose an action from a different package.

A composite's kwarg schema is the **union of its children's schemas** — union of names, union of required flags. A kwarg shared by two children is supplied once by the caller and passed to both. This is the mechanism's main ergonomic payoff: the model binds `class` and `name` once, and two files receive correctly shaped injections.

Two load-time checks:

- If two children declare the same kwarg name with **different types**, reject the package. Same-name-same-type is assumed intentional; same-name-different-type is an accident that cannot be recovered from at generation time.
- If any child requires a kwarg, the composite requires it.

Children execute in declaration order. There is no reordering and no dependency resolution between them.

### Exposure

Every action carries `exposed`, defaulting to `true`. Only exposed actions appear in the catalog shown to the model.

The default is permissive so that authoring an action is enough to make it usable — hiding is the deliberate act. Marking a sub-action `exposed: false` removes it from the model's option set entirely, which makes an entire class of invalid invocation unrepresentable rather than merely rejected.

An action that is unexposed and referenced by no composite is dead configuration: warn at load, do not error. It is almost always a rename that missed a call site.

A composite may reference an exposed child. This is legal, but it gives the model overlapping options and makes selection harder. Address it in authoring guidance, not machinery.

---

## Transforms

Templates and path patterns reference transforms with pipe syntax: `{{collection|constantize}}`, `{{class|snake}}`, `{{collection|plural|prefix:@}}`.

### Built-in operations

Shipped in Sedum's core, always available, pure `string -> string`:

`pascal`, `camel`, `snake`, `kebab`, `upper`, `lower`, `plural`, `singular`, `prefix:X`, `suffix:X`

Operation arguments are string literals only. Dynamic arguments (`prefix:{{other_arg}}`) are not supported — supporting them starts the construction of an expression language.

### Declarative pipelines

Named compositions of built-in operations, declared in a package's `sedum.yaml`:

```yaml
transforms:
  instantize:  [plural, "prefix:@"]
  constantize: [singular, pascal]
  pathify:     [plural, "suffix:_path"]
  tablename:   [plural, snake]
```

A Go package declares different pipelines over the identical operation set. Nothing language-specific enters Sedum's core.

Exception tables handle irregular mappings without new machinery:

```yaml
op_exceptions:
  pascal:
    url: URL
    id:  ID
```

A template referencing an undefined pipeline is a **hard error at package load**, not at render time.

### Inflection

`plural` and `singular` are the only operations that cannot be expressed as pattern rules alone, because English morphology is irregular. Two mechanisms cover it:

**A declarative rule table**, shipped per language rather than per package — ordered regex/replacement pairs, an irregulars map, an uncountables list.

**Model-supplied forms.** Since the model is already binding arguments, it may return explicit forms where they matter: `collection: { plural: users, singular: user }`. Morphology is something models handle reliably. The rule table is the default; a model-supplied form overrides it when both are present.

### User-defined transforms

A package may ship a `transforms` file in the target language exporting named `string -> string` functions, invoked when a template references a transform not resolvable among built-ins or pipelines.

This is an escape hatch with real costs: it introduces a runtime dependency, it means Sedum cannot validate a package without that runtime, and it breaks determinism, since nothing prevents user code from reading a file or hitting the network.

Three requirements make it tolerable:

- **Declared examples.** The transforms file ships `input -> expected` pairs. Sedum executes them at package load and refuses to proceed on mismatch.
- **A long-lived process.** One process with a line-oriented stdin/stdout protocol, not a spawn per invocation.
- **Explicit marking.** Packages using it are marked impure in the run log.

**This tier may be deferred out of the first implementation.** Built-ins, pipelines, and exception tables cover Rails, Go, and C++ conventions. Build it when a target demands it.

---

## The Pipeline

Seven phases, strictly ordered.

### Phase 0 — Load and validate generator packages

For every package in the generators directory: parse `sedum.yaml` and `actions.yaml`. Build the extension-to-package map. Resolve every simple action to its template file or variant directory. Verify every referenced transform pipeline exists. Verify composite kwarg unions are type-consistent. Verify every declared variant has a template file. Verify no composite composes a composite or crosses packages. Verify no two file templates tie under the specificity ranking. Check marker-based anchors against markers present in file templates. Run declared transform examples if user-defined transforms are present.

Packages are wholly valid or rejected. No partial loads.

### Phase 1 — Ingest provenance records

Read the provenance directory. Parse each record: `intent`, `constraints`, `affected_scope`, `forbidden_scope`. Validate schema conformance. Collect the authorized path set.

### Phase 2 — Resolve paths to packages and templates

For each authorized path, first check it against the union of every package's `unmanaged` patterns; a match is recorded and skipped. Then resolve its extension to a generator package, applying `--lang` where the extension is contested. A path whose extension no package claims is a hard error. A contested extension with no disambiguating flag is a hard error naming both candidate packages.

The unmanaged check runs *before* extension resolution, because the paths most often declared unmanaged are the ones no extension can reach.

Then match the path against that package's `files/` tree and bind the captured segments. Matching lives here rather than in Phase 3 because it is pure — it reads the loaded packages and decides, touching nothing — which is what lets `--stop-after resolution` report a matched template and its captures without any phase having written or read a file.

### Phase 3 — Create files from templates

For each resolved path: render the matched template against its captures and write the file. Fall back to the package's `_default` for that path's extension, then to an empty file with a log line.

Phase 3 and Phase 7 are the only phases that touch the output tree.

**A path that already exists is reconciled with its template, or the run halts.** Sedum never re-renders a template over a file, which would destroy any injected regions it carries. Instead it applies whatever the template *adds* and leaves everything else alone. A file that already carries everything its template declares is untouched.

The rule is one sentence: the file's structural lines must be a subsequence of the template's. If every structural line of the file appears in the template, in order, then the template contains everything the file does and the difference is only what the template adds — so adding it is safe. If one does not, the file carries content the template does not account for, and what to do with that is a judgment about someone's code rather than Sedum's to make, so the run halts naming the line.

Marked regions are removed before the comparison, so what is compared is the boilerplate rather than the boilerplate plus whatever has been injected into it since. A line is structural unless it is blank or begins with the package's declared `comment_prefix` — that prefix is the only language knowledge involved, and the package supplies it.

This replaced a stricter rule that verified an existing file carried its template's markers and halted otherwise. That refused three situations it could have handled: an empty file, a file another tool generated, and a file whose template has grown since it was written. Adopting Sedum into a repository that already has files is the ordinary case, not the exception.

Reruns are safe, which matters because stopping and resuming a run is a normal workflow rather than an edge case. The reconciliation cases are golden fixtures under `testdata/reconcile/`.

**Nothing is created that a provenance record did not authorize.** There is no sibling expansion and no inference of companion files. If a C++ record names `src/user_controller.cpp` but omits `include/user_controller.hpp`, the header is not created, and the injection targeting it fails loudly in Phase 7.

This is a governance position. `forbidden_scope` means Sedum does not touch what it was not authorized to touch; conjuring convenient files would violate that regardless of how obvious the convention seemed. Completeness of `affected_scope` is the record author's responsibility.

### Phase 4 — Model invocation

A path named by two records is rejected here, at the phase's entry. One model call per record means two records naming one file would be two independent calls deciding one file's contents. The check lives here rather than at ingestion because replay has no model call and so no such conflict — under `--execute` a recording may legitimately cover a shared file.

One call per provenance record. The prompt contains the record's `intent`, its `constraints`, the paths created for it in Phase 3, and the action catalog — the **union of exposed actions across every package the record's paths resolved to**, with their kwarg schemas and variant lists.

Variant lists are included deliberately. Without them there is an invisible cliff: `name: index` gets a full implementation while `name: search` falls to `_default`, and the model has no way to know it fell off. Exposing the list lets it prefer covered values where intent maps cleanly, and take the fallback knowingly where it does not.

The response is **structured output, not tool calls** — a JSON array of `{action, kwargs}` objects. This keeps the mechanism working with models that lack tool-calling support, which is most of the open-weight range worth evaluating.

### Phase 5 — Validate the model's output

Deterministic checks, each producing a specific re-promptable error:

- Action name exists and is exposed in the record's catalog
- All required kwargs present
- No unknown kwargs
- Every kwarg value matches its declared type
- Discriminator value is a declared variant, or `_default` exists
- Rendered `injects_into` path was created in Phase 3
- The target file carries the anchor the action injects into

On failure, re-prompt with the specific violations appended, up to a configured retry limit. This loop costs one model call — no compilation, no service startup, no test execution.

### Phase 6 — Expand and resolve

Expand composites into ordered children, mapping union kwargs to each child's schema. Render every `injects_into` pattern. Select the variant template for each discriminated action. Apply transforms.

Fully deterministic. The model does not participate.

### Phase 7 — Inject

For each resolved invocation: render the template with bound kwargs, locate the anchor in the target file, write the rendered content into the anchored region.

**A missing anchor is a hard error.** It means the file is not shaped the way the action assumed — a disagreement between configuration and reality that the author must resolve. Auto-creating an anchor would paper over exactly the mistake worth surfacing.

---

## Anchors

Anchors are a small closed vocabulary, declared per action, evaluated at the text level. No parsers, no per-language AST work.

`marker` — a named Sedum comment planted by a file template
`region` — between a named start and end marker
`start_of_file`
`end_of_file`
`after_match` / `before_match` — a regex declared in the action definition

Marker comments are the load-bearing case. Marker syntax uses the package's declared `comment_prefix`, since `#`, `//`, and `--` all appear across targets.

### Placement

Where an anchor puts content is part of the authoring contract, so it is stated rather than discovered:

- `marker` inserts after the line carrying the marker.
- `region` accumulates at the end of the region, just inside the marker that closes it, so repeated injections stay in invocation order.
- `after_match` and `before_match` snap to the line boundary around the match rather than to the match's own bounds, so injected content never lands inside an existing line.
- A marker name is compared for equality rather than as a prefix, so an action anchored to `class_body` does not land at `class_body_top`.

### What Phase 0 checks

Each anchor kind must carry exactly its companion fields — `region` needs both `anchor_start` and `anchor_end`, `after_match` and `before_match` need `anchor_pattern`, and bare `marker` is rejected because it names the kind rather than a marker.

`anchor_pattern` is compiled at load. An expression that does not parse is a defect in the package, and a regex is checkable with nothing but the regex, so a package carrying a broken one is rejected before anything is written.

A pattern using `^` or `$` without `(?m)` **warns**. Those anchor to the bounds of the whole file rather than of a line, so a pattern meant to find a line finds nothing and the failure surfaces at injection time as a fault in the file rather than in the pattern. It warns rather than erroring because whole-text anchoring is legal and occasionally meant — and the expression is never rewritten on the author's behalf, since the pattern in `actions.yaml` has to be the pattern that runs.

An action targeting a marker that no file template in its package plants also warns. Because template selection is path-dependent this cannot be a complete verification, but it is almost always a typo.

---

## Ownership and Idempotency

Every injected region is wrapped in ownership markers naming the action that produced it, the ownership tier, the record that last parameterized it, and the kwargs it was rendered from:

```ruby
  # sedum:createControllerMethod:index {"tier":"owned","record":"PR-014","kwargs":{"collection":"users","controller":"users"}}
def index
  @users = User.all
  render json: @users
end
  # /sedum:createControllerMethod:index
```

The action and variant stay literal on the line because they are the audit trail, and grepping only works if it is grep rather than a parser. Everything else is one JSON object rather than positional fields, because the parser has to ignore fields it does not recognize and default the ones that are absent. Positional slots would make every field added later a migration across every repository already carrying markers; an object makes it an addition.

Anchor declarations planted by file templates share the `sedum:` namespace with ownership markers, so the reader tells the two apart. The consequence is that `anchor` is not available as an action name, and a package declaring one is rejected at load.

Sedum's own marker lines take the indentation of the anchor they land at, copied verbatim from the file. The rendered body is never re-indented: a template author writes each template at the depth its anchor sits at, and one package legitimately mixes a fragment indented to sit inside a struct with a top-level declaration starting at column zero. The author owns the body; Sedum owns its own lines.

Re-running a generation replaces the region an action owns rather than appending beside it. Reruns and partial regeneration are safe without a sidecar cache or resolution manifest.

This produces the audit trail as a side effect: grepping markers yields file -> action -> variant -> arguments with no maintained state, and ownership is visible in the diff a human reviews.

### Ownership tiers

The tier field declares whether Sedum may overwrite a region.

`owned` — Sedum generated this region and replaces it on every run. Edits to it are lost, whoever made them. The tier is a statement about the region, not about humans: a tool that edits an owned region loses its work exactly as a person does.

`seeded` — Sedum generated this region once and never touches it again. Present in the file, skipped on rerun.

An action declares its tier in `actions.yaml`, defaulting to `owned`. A template whose body is a stub a human is expected to complete should declare `seeded`; a template that fully determines its output should not.

**The marker's tier governs after the first write.** The `actions.yaml` declaration supplies the value written the first time and has no authority over a region that already exists. A region demoted to `seeded` after generation — because it stopped being Sedum's to overwrite — is exactly the region a declaration must not overrule, so the file decides and the declaration is only the default it started from.

**The `tier` key is not yet accepted in `actions.yaml`** — decoding is strict, so declaring one is currently an error, and the precedence rule above is the rule for when it lands. Both tiers are honored when *read* from a marker, which is what makes adding the key later an addition rather than a migration: the field is already on every marker Sedum has ever written.

### Recorded kwargs

The kwargs serialized on the opening marker make a region self-describing. A reader — human or tooling — can see what the region was parameterized with without consulting a recording or rerunning resolution.

Nothing currently reads them back. They are written for the same reason the tier field is: markers are durable artifacts, and changing their shape after generated codebases exist is expensive.

The marker carries the record ID from the first generated file — as an attribute, not as part of a region's identity. A region is identified by its action, its variant, and the kwargs it was rendered from; the record ID records who last wrote it. That distinction is what lets a later record refine a region an earlier one produced, rather than minting a second region beside it.

Refinement of a region across records and across runs is what that attribute is for, and it works. Two records naming **one path in one run** is a different thing and is still rejected at ingestion, because Phase 4 makes one model call per record — see the note in [Status](#status) about `--only`.

### The attribute object is the extension point

Sedum owns the marker schema but not time: a marker sits in a generated codebase long after the version that wrote it is gone, so a reader and a writer of the same marker are routinely different versions — and, once anything is built on top of Sedum, different tools entirely.

Two promises follow, and Sedum keeps both.

**A marker is never rejected for carrying more than the reading version expects, or less.** An unrecognised key is ignored; a declared key that is absent takes its documented default. An attribute object that is not readable as JSON is still an error — that is corruption rather than version skew.

**An unrecognised key survives the round trip.** It is retained on read and re-emitted when the region is rewritten, including when a rerun replaces the region's contents:

```ruby
  # sedum:createControllerMethod:index {"tier":"owned","record":"PR-014","kwargs":{"controller":"users"},"verified_by":"spec/users.linespec"}
```

Rerun that region and `verified_by` is still there. This is the half that does not come free — ignoring a key on read is what a JSON decoder does anyway; preserving it on write is a separate promise, and Sedum's own writer had to be taught not to drop it.

That promise is what makes the object usable by anything other than Sedum. A tool built above Sedum records its own per-region state under its own keys and finds it intact after the next run, with no schema change and no sidecar file keyed by region — which would reintroduce exactly the maintained state markers exist to avoid.

Preservation is transparent: Sedum does not read, validate, interpret, or reorder a carried key, and a carried key may not shadow one Sedum models. The modelled keys are `tier`, `record`, `writer`, and `kwargs`. Carried keys are written after them in sorted order, so a marker with none is byte-identical to one written before any of this existed.

`writer` names the tool that last wrote the region, and is omitted when that tool is Sedum. It exists because a demoted tier is otherwise unattributable — a region reading `seeded` cannot say whether a package author declared it or another tool demoted it.

There is no version token. Additive evolution is already free under the object; a token would earn its place only for a breaking change, meaning an existing field whose meaning changes.

---

## Recorded Executions

Every phase after model invocation is deterministic. That means the model's contribution to a run can be captured once and replayed forever.

A **recording** is a JSON file describing everything Sedum resolved and everything it was told to do: which files to create, which package and template each resolved to, and which action invocations to apply. Replaying a recording produces the same result without invoking a model at all.

This makes three things possible. Runs become reproducible in the strong sense — same recording, same output, no sampling involved. Teams can commit a recording as a standard service scaffold and generate new services from it deterministically. And a recording is plain JSON, so it can be hand-edited when a situation needs something the model did not produce or produced wrongly.

### Format

```json
{
  "sedum_version": "0.1.0",
  "packages": {
    "rails": { "extensions": [".rb", ".erb"] }
  },
  "records": [
    {
      "record_id": "PR-014",
      "files": [
        {
          "path": "app/controllers/users_controller.rb",
          "package": "rails",
          "template": "app/controllers/{name}_controller.rb",
          "captures": { "name": "users" }
        }
      ],
      "phases": [
        {
          "name": "default",
          "invocations": [
            {
              "action": "createControllerMethod",
              "kwargs": {
                "controller": "users",
                "name": "index",
                "collection": "users"
              }
            }
          ]
        }
      ]
    }
  ]
}
```

Invocations are grouped under `phases` rather than listed flat. Every recording produced by this implementation contains exactly one phase, named `default`, and replay executes phases in order.

The grouping exists because recordings are committed artifacts. A team's standard service scaffold lives in version control, and changing the schema after those files exist means migrating them. Reserving the level costs one array nesting; adding it later costs a migration.

Invocations are recorded **pre-expansion**. A composite is stored as the composite, not as its children, because expansion is deterministic and re-running it costs nothing. This keeps recordings compact and keeps them at the abstraction level an author edits in — changing a `createMethod` call means editing one entry, not keeping two injections in sync.

Recordings carry no timestamps, run identifiers, or other volatile fields. Two recordings of equivalent runs should diff cleanly, so that a change in model output is visible as a change in the file.

### Replay semantics

Replaying enters the pipeline at Phase 3 with resolution already decided, skips Phase 4 entirely, and runs Phases 5 through 7 unchanged.

**Validation is identical, but failures are terminal.** A hand-edited recording naming an action that does not exist, omitting a required kwarg, or targeting a file it did not create fails exactly the checks a model response would fail. There is no re-prompt loop, because there is nothing to re-prompt — the error is reported and the run halts.

Sedum verifies that every package named in the recording is present in the generators directory and still claims the extensions recorded against it. A mismatch halts the run.

**A recording is not a lockfile.** It captures decisions, not templates. Templates live in the generator package, so replaying after a package change picks up the new templates. This is usually what teams want — boilerplate improvements propagate to every scaffold — but it means byte-identical output across time requires versioning the generator package alongside the recording.

### Scope validation on replay

Provenance records are optional when replaying. Supplied, Sedum verifies that every path in the recording is authorized by `affected_scope` and violates no `forbidden_scope`, and halts on any unauthorized path. Omitted, the recording executes as written.

Both modes are legitimate. A team replaying a recording against the records that produced it wants the check. A team using a hand-edited recording as a generic service scaffold has no corresponding records to check against.

---

## The Integration Surface

Sedum is a component. Tools will be built on top of it — a loop that drives generation, verifies the result, and repairs what did not converge is the obvious one, and Sedum is deliberately not that tool.

> **Everything a tool built on Sedum needs is derivable from artifacts on disk. Sedum exposes no runtime API, and nothing links against Sedum.**

Five things are supported, three of them file formats and two of them commands:

| Surface | Direction | Who else authors it |
|---|---|---|
| **Ownership markers** | read; written within the rules below | Sedum writes them |
| **Generator package** (`sedum.yaml`, `actions.yaml`, templates) | read | users author them |
| **Recording** (JSON) | **written** and read | Sedum can write them too |
| **`sedum grow --execute`** | invoked | — |
| **`sedum render`** | invoked | — |

Everything else — Sedum's Go packages, the run log's shape, the `--stop-after` phase names — is internal and may change.

**All five are live today.** The marker's guarantees below — unrecognised-key preservation in both directions, and the `writer` attribute — are implemented and tested now, because they are the ones whose absence would be expensive to correct after markers exist in generated repositories.

**"What exists?"** Grep the markers. Enumerating what has been generated, with which action, variant, and arguments, is a filesystem walk needing nothing from Sedum's process.

**"Could this be fixed deterministically?"** Ask `sedum render`, which renders an action's `injects_into` against the kwargs already on the marker so the caller can compare the result to the region in the file. Forward rendering and comparison, never inversion — `snake` and `plural` are not invertible. It works only because the kwargs are recorded on the marker.

It is a command rather than an instruction to read `actions.yaml` and do it yourself: the inflection table ships per language rather than per package, so a caller rendering the pattern by hand reimplements `plural` and diverges on the first irregular.

**"Do it."** Synthesise a recording of exactly the invocations wanted and `--execute` it, never touching Phase 4. `--record --dry-run` gives the other direction.

**"What did Sedum decline to do?"** The `unmanaged` patterns are in `sedum.yaml`, which is already being read.

### The recording is an input format

The recording is not only a capture artifact — it is how a caller tells Sedum what to do with no model in the loop, and replay's semantics are already right for that. **Validation is identical to model-output validation, but failures are terminal.** A recording naming a nonexistent action, omitting a required kwarg, or targeting an unauthorized path fails exactly the checks a model response would fail, and the run halts rather than re-prompting, because there is nothing to re-prompt.

That is a submission protocol, not merely a convenience for hand-editing.

### Direction of control

A tool built on Sedum sits **above** it, not after it. It calls Sedum; Sedum never calls it. There is no callback, no hook, and no plugin point anywhere in the pipeline — every phase after model invocation is deterministic precisely because nothing foreign runs inside it, and adding an extension point would make Sedum's guarantees depend on code it does not control.

This is also what keeps Sedum independently runnable: committing a recording as a standard service scaffold and replaying it needs no harness, no container, and no model server.

### What a foreign writer may not do

Sedum tolerates other writers in the files it generates. The limits are stated here rather than discovered, because an unwritten constraint's first breakage looks like a Sedum bug.

- **Do not remove or relocate a marker a file template planted.** Anchors exist because a template created them; moving one breaks the injection targeting it, and Phase 7 reports a missing anchor rather than guessing.
- **Do not reorder regions within a file.** Repeated injections at one anchor accumulate in invocation order.
- **Preserve unknown marker attributes in both directions.** The same promise Sedum makes, and it is symmetric — a writer that drops keys it does not recognise destroys every other writer's state, Sedum's included.
- **Do not edit an `owned` region and expect the edit to survive.** Demote it to `seeded` if it has stopped being Sedum's to overwrite; that is what the tier is for and what `writer` makes attributable.

### The conformance corpus

`conformance/` holds golden cases for the formats Sedum shares with other tools,
so a tool Sedum's authors did not write can check that it complies rather than
discovering the rules one breakage at a time. Sedum's own test suite reads the
same corpus, which is what keeps it honest.

Today it covers one format, the ownership marker, and that is a deliberate
scope. A tool above Sedum reads `sedum.yaml` and `actions.yaml`, so a misreading
is its own problem and fails on its own side. It writes recordings, and replay
validates those terminally at ingestion, so a malformed one is refused before
anything is written. The marker is the only format a foreign tool writes into
files that Sedum will later re-read and rewrite with no validation gate in
between — a writer that gets it wrong corrupts state silently, and the first
breakage looks like a bug in Sedum.

Cases live in `conformance/markers/cases.json`, each citing the provenance
record that decided the behavior and what goes wrong for an implementation that
gets it wrong. See `conformance/README.md`.

---

## CLI Surface

This is the full designed surface. See [Status](#status) for what is wired up today.

### `sedum grow`

Runs the full pipeline: load packages, ingest records, resolve, create, invoke, validate, expand, inject.

```
sedum grow --generators ./generators --records ./provenance --output ./build
```

| Flag | Description |
|---|---|
| `--generators <dir>` | Generators directory. Required. |
| `--records <dir>` | Provenance records directory. Required. |
| `--output <dir>` | Output directory. Defaults to the current directory. |
| `--lang <name>` | Prefer the named package where an extension is contested. Repeatable. |
| `--only <id>` | Generate only the named provenance record. Repeatable. |
| `--record <path>` | Write a recording of the run to the given path. |
| `--execute <path>` | Replay a recording. Skips model invocation. Mutually exclusive with `--record`. |
| `--dry-run` | Run every phase, write nothing. Reports the files that would be created and the injections that would be applied. |
| `--stop-after <phase>` | Halt after the named phase. One of `resolution`, `files`, `invocations`, `expansion`. |
| `--retries <n>` | Model output validation retry limit. Default 3. Ignored with `--execute`. |
| `--model <name>` | Model identifier. Endpoint and credentials come from environment. Ignored with `--execute`. |
| `--log <path>` | Run log location. Defaults to `.sedum/run.log`. |
| `-v, --verbose` | Mirror the run log to stdout. |

`--record` and `--dry-run` compose. Together they capture a recording without writing any generated files — the model runs, its output is validated and saved, and nothing is created. Without `--dry-run`, the run records and executes in the same pass.

`--execute` and `--dry-run` also compose, validating a recording against the generator packages and reporting what it would produce without writing.

With `--execute`, `--records` becomes optional and enables scope validation when supplied. `--lang` is ignored, since package resolution is recorded per file.

```sh
# generate and capture
sedum grow --generators ./generators --records ./provenance --record ./scaffold.json

# capture only, write nothing
sedum grow --generators ./generators --records ./provenance --record ./scaffold.json --dry-run

# replay, no model
sedum grow --generators ./generators --execute ./scaffold.json --output ./build

# replay with scope enforcement
sedum grow --generators ./generators --records ./provenance --execute ./scaffold.json
```

### Stop points

`--stop-after` halts a run at a phase boundary so a person can inspect what Sedum decided before it acts on the decision. Each stop point has a resume path; stopping is not abandoning the run.

| Value | Halts after | What is available to review | Resume |
|---|---|---|---|
| `resolution` | Phase 2 | Every authorized path with its resolved package, matched file template, and bound captures. Nothing written. | Rerun `grow` |
| `files` | Phase 3 | The scaffolded files on disk, with markers planted and no injections applied. | Rerun `grow` |
| `invocations` | Phase 5 | The validated action list the model produced, written to the recording. Files exist; nothing injected. | `--execute` the recording |
| `expansion` | Phase 6 | Composites expanded to children, paths rendered, variants selected, transforms applied. Diagnostic. | `--execute` the recording |

`invocations` is the load-bearing one. It is the point where a person reads what the model decided, corrects it in a text editor, and hands the corrected version back — the whole reason recordings are plain JSON.

**`--stop-after invocations` and `--stop-after expansion` require `--record`.** Without it the model's output is discarded and there is nothing to resume from, having already paid for the call. Sedum rejects the combination rather than running it.

Resuming after `resolution` or `files` is an ordinary rerun, because Phases 0–2 are pure and Phase 3 leaves an already-reconciled file untouched. Nothing needs to be preserved between the stop and the resume.

With `--execute`, a `--stop-after` value naming a phase that replay does not run is an error. Only `expansion` is meaningful, and `--execute --dry-run` covers most of what it offers.

`--stop-after resolution` overlaps `sedum resolve` by design. The command is the read-only inspection tool; the flag is a checkpoint inside a run that will continue.

### `sedum validate`

Runs Phase 0 against the generators directory and exits. Reports every error and warning found. Requires no records, no model, and no network.

```
sedum validate --generators ./generators
```

| Flag | Description |
|---|---|
| `--generators <dir>` | Generators directory. Required. |
| `--package <name>` | Validate a single package. Repeatable. |
| `--strict` | Treat warnings as errors. |

### `sedum resolve`

Reports what each authorized path resolved to: its package, the file template that matched, and the captures bound. The primary debugging tool for package resolution and template specificity.

It consults **no output tree**. The only directories it reads are the two it is pointed at, so its answer never depends on where you run it from. `--show-template` renders each matched template, which needs nothing but the template and its captures. To ask instead whether a particular output tree already holds these files and whether their markers are intact, use `grow --stop-after files --dry-run --output <dir>`.

```
sedum resolve --generators ./generators --records ./provenance
```

| Flag | Description |
|---|---|
| `--generators <dir>` | Generators directory. Required. |
| `--records <dir>` | Provenance records directory. Required. |
| `--lang <name>` | Prefer the named package where an extension is contested. Repeatable. |
| `--only <id>` | Resolve only the named provenance record. Repeatable. |
| `--show-template` | Include rendered template output for each path. |

### `sedum actions`

Prints a package's action catalog exactly as the model would receive it — exposed actions, kwarg schemas, variant lists. The authoring feedback loop for exposure and catalog clarity.

```
sedum actions --generators ./generators --package rails
```

| Flag | Description |
|---|---|
| `--generators <dir>` | Generators directory. Required. |
| `--package <name>` | Package to inspect. Required. |
| `--all` | Include unexposed actions, marked as such. |
| `--json` | Emit the raw catalog payload rather than formatted output. |

---

## Logging

There is no plan artifact. The execution sequence is a consequence of the record and the configuration, not a decision any component makes.

The run log records package resolution, file template matches and captures, the model's raw response, validation failures and retries, composite expansions, resolved paths, selected variants, transform resolutions, and anchor matches. It is diagnostic output, cleared per run. Nothing depends on it — idempotency state lives in the ownership markers, in the generated files.

---

## Non-Goals

**Business logic synthesis.** Templates produce boilerplate. Domain-specific logic is not in scope.

**Structured-document editing.** Appending text to `pom.xml`, `package.json`, or `CMakeLists.txt` is wrong. Doing it correctly requires a format-aware primitive with document selectors and merge semantics. Deferred.

**Cross-action value flow.** No action's template may depend on a value produced by another action. Every value comes from the model's bound kwargs or from transforms over them. If this constraint becomes limiting, it is a signal that a planning layer is warranted — a substantially different architecture, not an increment on this one.

**Structure inference from unmarked source.** Sedum reads structure only where a file template planted a marker. It does not parse a target language, and it does not derive regions, conventions, or injection points from source it did not shape.

Reading is not inference. Phase 3 reads an existing file to verify its markers are present, and enumerating what has been generated by grepping markers is a supported operation — both read what Sedum wrote, which is the distinction that matters.

**Behavioral verification.** Sedum does not run or grade the code it generates, and does not know whether the result works. A loop that generates, verifies, and repairs is a different tool that calls this one — see [The Integration Surface](#the-integration-surface).

---

## Repository layout

```
cmd/sedum/            Entry point.
internal/cli/         Command surface: flags, interdependence checks, config structs.
internal/pipeline/    Phase ordering and stop points.
internal/genpkg/      Package loading and every Phase 0 check.
internal/record/      Provenance record ingestion and the authorized path set.
internal/resolve/     Path-to-package resolution, template matching, file creation and reconciliation.
internal/selection/   Phase 4: the model client, the prompt, and Phase 5's validation.
internal/catalog/     The action catalog a record's packages present to the model.
internal/pathpat/     The scope and unmanaged pattern grammar.
internal/expand/      Phase 6: composite expansion, injects_into rendering, variant selection, transforms.
internal/inject/      Phase 7: the marker format, anchor location, region replacement.
internal/recording/   The recording schema, and reading and writing it.
internal/transform/   Built-in operations, declarative pipelines, inflection.
internal/filetmpl/    File template patterns, matching, and specificity ranking.
internal/render/      Template rendering with transform pipes.
internal/runlog/      Run log.
testdata/generators/  Worked example packages: rails, chi, cairn.
testdata/reconcile/   Golden cases for Phase 3 reconciliation.
conformance/          Golden cases for the formats Sedum shares with other tools.
evals/                The harness measuring what models select, and its results.
docs/                 The published results site.
provenance/           Provenance records governing Sedum's own development.
PRD.md                Product requirements, including the milestone plan.
OPEN_QUESTIONS.md     Design directions explored but deliberately out of scope.
TOOL_BOUNDARIES.md    What belongs in Sedum and what belongs in a tool above it.
```

Sedum's own development is governed by the provenance records in `provenance/`.
Every design decision behind the code — including the ones that were considered
and rejected — is recorded there rather than in commit messages, so
`linespec provenance status` is the fastest way to see why something is the way
it is.

---

## Success Criterion

A team authors a generator package for their stack. They point Sedum at a provenance directory and that package. They receive the files the records authorized, shaped by their file templates, containing the boilerplate their action templates describe.

They run it again and get the same result.

Sedum's core, throughout, knows nothing about their language.
