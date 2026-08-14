# Sedum

## Product Requirements Document

---

## What Sedum Is

Sedum generates boilerplate code from provenance records using generator packages that teams author themselves.

A provenance record declares intent, constraints, and the files a change is authorized to touch. A generator package declares a team's conventions: the shape of each file type, which code-injection actions exist, what arguments they take, what templates they render, and where in a file the results belong. Sedum reads both, uses a language model to translate intent into a list of action invocations, and executes those invocations deterministically.

**Sedum's core contains no language-specific knowledge.** It knows how to read a directory, match a path against a pattern, apply string operations, render a template, and place text at an anchor. It does not know what a controller is, what a header file is, or how Go capitalizes initialisms. Teams express all of that in configuration. A team switching from Rails to Java writes a new generator package; they do not wait for Sedum to add support.

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

Entries use the same grammar as a record's `affected_scope` — `**`, `*`, `?`, `[…]`, and a trailing slash for a subtree.

An authorized path matching one is skipped, reported, and the run continues. It is not created, not rendered, and not injected into.

This exists because a record must be free to describe the whole change. Moving a Rails service to PostgreSQL means adding the `pg` gem, so `Gemfile` belongs in `affected_scope` — and `Gemfile` has no extension, so without a declaration the run halts on it. The alternative is trimming records down to what a generator happens to reach, which makes the record a description of the tool rather than of the work.

**It authorizes nothing.** `affected_scope` still decides what may be touched; `unmanaged` says only that Sedum is not what touches it. The usual reason is that a person will, or that another tool will be pointed at it, so a run reports these paths as a list rather than passing over them silently. The handoff is the point.

The declaration belongs to the package rather than to a flag: which files a human owns is a property of a stack, not of an invocation. Every Rails service has a Gemfile someone edits; a run that had to remember that would eventually forget.

It is not a step toward structured-document editing. Appending a gem to a `Gemfile` needs a format-aware primitive with merge semantics, which is a non-goal. Declaring the path unmanaged is how a package says so out loud.

Two contradictions are load-time errors: a file template whose own path is declared unmanaged, and an action whose literal `injects_into` names one. An `injects_into` carrying placeholders is caught in Phase 6 instead, with a diagnostic that distinguishes *declared unmanaged* from *never created* — the two have different fixes.

Two packages declaring the same path do not conflict; they agree. A path one package disowns and another would claim by extension is unmanaged, because the alternative is behavior that depends on which package was consulted first.

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

### Phase 2 — Resolve paths to packages

For each authorized path, first check it against the union of every package's `unmanaged` patterns; a match is recorded and skipped. Then resolve its extension to a generator package, applying `--lang` where the extension is contested. A path whose extension no package claims is a hard error. A contested extension with no disambiguating flag is a hard error naming both candidate packages.

The unmanaged check runs *before* extension resolution, not after, because the paths most often declared unmanaged are the ones no extension can reach.

### Phase 3 — Create files from templates

For each authorized path: match it against its package's `files/` tree, bind captured segments, render the matched template, and write the file. Fall back to `_default`, then to an empty file with a log line.

**Phase 3 is create-if-absent.** If a path already exists, Sedum does not re-render its template over it — doing so would destroy any injected regions the file already carries. It verifies that the markers the matched template declares are present and moves on. A file that exists but lacks its template's markers is a hard error: something other than Sedum wrote it, or a template changed shape after the file was generated.

This makes reruns safe, which matters because stopping and resuming a run is a normal workflow rather than an edge case.

**Nothing is created that a provenance record did not authorize.** There is no sibling expansion and no inference of companion files. If a C++ record names `src/user_controller.cpp` but omits `include/user_controller.hpp`, the header is not created, and the injection targeting it fails loudly in Phase 7.

This is a governance position. `forbidden_scope` means Sedum does not touch what it was not authorized to touch; conjuring convenient files would violate that regardless of how obvious the convention seemed. Completeness of `affected_scope` is the record author's responsibility.

### Phase 4 — Model invocation

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

---

## Ownership and Idempotency

Every injected region is wrapped in ownership markers naming the action that produced it, the ownership tier, and the kwargs it was rendered from:

```ruby
# sedum:createControllerMethod:index owned {"controller":"users","collection":"users"}
def index
  @users = User.all
  render json: @users
end
# /sedum:createControllerMethod:index
```

Re-running a generation replaces the region an action owns rather than appending beside it. Reruns and partial regeneration are safe without a sidecar cache or resolution manifest.

This produces the audit trail as a side effect: grepping markers yields file -> action -> variant -> arguments with no maintained state, and ownership is visible in the diff a human reviews.

### Ownership tiers

The tier field declares whether Sedum may overwrite a region.

`owned` — Sedum generated this region and replaces it on every run. Hand edits are lost.

`seeded` — Sedum generated this region once and never touches it again. Present in the file, skipped on rerun.

An action declares its tier in `actions.yaml`, defaulting to `owned`. A template whose body is a stub a human is expected to complete should declare `seeded`; a template that fully determines its output should not.

**Only `owned` is exercised by the milestones in this document.** `seeded` is specified now because markers are written to disk and read back on rerun. Adding the field later would leave every file generated in the interim carrying markers in an older shape, requiring a migration across generated codebases. The cost of reserving it is one token per marker.

### Recorded kwargs

The kwargs serialized on the opening marker make a region self-describing. A reader — human or tooling — can see what the region was parameterized with without consulting a recording or rerunning resolution.

Nothing in this document reads them back. They are written for the same reason the tier field is: markers are durable artifacts, and changing their shape after generated codebases exist is expensive.

Multiple records injecting into the same file is out of scope. If it becomes necessary, the marker gains the record ID and nothing else changes.

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

## CLI Surface

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

```
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

Resuming after `resolution` or `files` is an ordinary rerun, because Phase 3 is create-if-absent and Phases 0–2 are pure. Nothing needs to be preserved between the stop and the resume.

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

Runs Phases 0 through 3 without invoking the model, and reports what each authorized path resolved to: its package, the file template that matched, and the captures bound. The primary debugging tool for package resolution and template specificity.

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

**Scanning existing codebases.** Sedum generates into paths it creates. Deriving state from a pre-existing project is out of scope.

**Behavioral verification.** Sedum does not run or grade the code it generates.

---

## Milestones

**M1 — Package loading and validation.** Multi-package discovery, extension mapping, both YAML schemas, action template resolution, all load-time checks. Delivers `sedum validate`. *Testable:* a deliberately malformed package produces the correct specific error for each validation rule.

**M2 — Transforms.** Built-in operations, declarative pipelines, exception tables, inflection rule table. *Testable:* a table of input/pipeline/expected triples; an undefined pipeline fails at load, not at render.

**M3 — Path resolution and file creation.** Provenance ingestion, package resolution with `--lang`, file template matching with specificity ranking, template rendering, create-if-absent semantics. Delivers `sedum resolve` and `--stop-after resolution|files`. *Testable:* every `affected_scope` path exists with correct boilerplate; nothing else was created; rerunning does not re-render over existing files; a contested extension without a flag errors cleanly.

**M4 — Action rendering and anchored injection.** Phases 6 and 7, driven by hand-written invocation lists, with no model involved. *Testable:* a fixture invocation list produces byte-identical output on repeat runs; re-running replaces owned regions rather than duplicating them.

**M5 — Composite expansion.** Union kwarg mapping, ordered child execution. *Testable:* a C++ `createMethod` fixture injects correctly into both header and implementation; omitting the header from `affected_scope` produces a clear hard error.

**M6 — Model invocation and validation.** Phases 4 and 5, closing the loop. Delivers `sedum grow` and `sedum actions`. *Testable:* end-to-end generation from a Rails provenance record with no hand-written invocation list; each validation rule demonstrably triggers a correct re-prompt.

**M7 — Recording and replay.** Serialization of resolved files and validated invocations, replay entry at Phase 3, terminal validation on replay, optional scope enforcement. Delivers `--record`, `--execute`, and `--stop-after invocations|expansion`. *Testable:* recording a run and replaying it produces byte-identical output with no model invocation; stopping at `invocations` and resuming with `--execute` matches an uninterrupted run; a hand-edited recording with an invalid action halts with a specific error; a recording with an unauthorized path fails scope validation when records are supplied.

M4 precedes M6 deliberately. The deterministic half must be provably correct before a model is introduced, so that any failure after M6 is unambiguously a selection failure rather than a rendering failure.

M7's fixture invocation lists from M4 are effectively hand-written recordings. Keeping M4's fixture format aligned with the recording schema means M7 is largely serialization work rather than new execution machinery.

---

## Success Criterion

A team authors a generator package for their stack. They point Sedum at a provenance directory and that package. They receive the files the records authorized, shaped by their file templates, containing the boilerplate their action templates describe.

They run it again and get the same result.

Sedum's core, throughout, knows nothing about their language.
