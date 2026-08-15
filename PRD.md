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
      collection: { type: string, required: true }
    discriminator: name
    variants: [index, show, create, update, destroy]
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body
```

`kwargs` is the schema the model is held to. Types come from a closed set — `string`, `int`, `bool`, `list` — sufficient for argument binding and nothing more.

### Descriptions

A kwarg may declare a `description`, and so may an action:

```yaml
addColumn:
  description: Adds a column to a migration's change block.
  kwargs:
    name:
      type: string
      required: true
      description: the column name, bare — the template writes the leading colon
    options:
      type: string
      required: false
      description: the rest of the line after the column name, e.g. "null: false"
```

The type set is deliberately tiny, so it cannot express *a comma-separated list of bare names* or *the rest of the line after the column name* — and should not try. An author can say it in a sentence, and until they had somewhere to say it, they said it in a YAML comment the model never sees.

The retry loop is not the alternative. A wrong value here is usually a **valid string for a required string kwarg**, so every Phase 5 check passes and the failure surfaces as an error in the running target, which is past the last point Sedum looks. Where a check can tell, it names a rendered path or a template file rather than the convention that was missed, so the model corrects toward the symptom. The retry loop is a backstop for mistakes, not a channel for information the catalog could have carried.

Both are optional, both are passed through untouched, and absent means absent — nothing is synthesised to fill the gap, because an invented description reads with an authority it has not earned.

**Nothing interprets a description.** Sedum does not parse one, validate one, derive a constraint from one, or check that a bound value agrees with one. A description must never become a place where a rule lives: a rule the model can read and Phase 5 cannot enforce is worse than no rule. If a rule matters, it belongs in the schema where it is checked.

This makes catalog quality an authoring responsibility rather than a property of the model, which is the right place for it. `sedum actions` is where an author sees what the model will see.

A kwarg every one of an action's templates renders unconditionally should be declared `required`, as `collection` is above, and Sedum warns when a single-template action declares one optional. But the declaration is not the only source of the requirement, because it cannot be: a discriminated action shares one kwarg schema across every variant, so a value that `index` needs and `destroy` forbids can only be declared optional.

**So a requirement is also derived from the template.** Sedum's grammar has no conditionals, no loops, and no field access — `{{name}}` and `{{name|op}}` are the whole of it — so a value a template references is a value that template unconditionally needs. There is no shape of template for which *referenced* and *required* differ, which makes the derivation exact rather than approximate.

The effective requirement for an invocation is the **union** of the action's declared required kwargs and the values referenced by the template that invocation selects. A declared requirement holds for every variant; a derived one holds when its template is chosen. Neither replaces the other, and they are reported as separate violations because they have different fixes.

A template referencing a value its action does not declare at all is a load-time error.

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

One call per provenance record. The prompt contains the record's `intent`, its `constraints`, the paths created for it in Phase 3, and the action catalog — the **union of exposed actions across every package the record's paths resolved to**, with their kwarg schemas, authored descriptions, variant lists, per-variant derived requirements, and `injects_into` patterns.

Variant lists are included deliberately. Without them there is an invisible cliff: `name: index` gets a full implementation while `name: search` falls to `_default`, and the model has no way to know it fell off. Exposing the list lets it prefer covered values where intent maps cleanly, and take the fallback knowingly where it does not. Whether a `_default` exists is carried alongside, because *knowingly* is not available to a model that cannot see whether there is a fallback to take.

`injects_into` patterns are included for a harder reason. Without them the catalog names a kwarg and the file list names a path, and nothing connects the two, so a model asked to bind `controller` binds the path it was shown. Recovering the kwarg from a rendered path would require inverting `snake` and `plural`, which is precisely what nothing in this system does. With the pattern present the reasoning runs forwards: match its literal segments against an authorized file and bind what is left. The model still never chooses a path — it binds arguments such that the package's own pattern lands on a file the record authorized.

The response is **structured output, not tool calls** — a JSON array of `{action, kwargs}` objects, wrapped as `{"invocations": [...]}` so that the document's root is an object. This keeps the mechanism working with models that lack tool-calling support, which is most of the open-weight range worth evaluating.

Sedum does not ask the server to constrain the response. Grammar-constrained decoding takes the shortest legal completion at the array's first token, and an empty array is legal, so requesting a schema made a model that had been selecting three correct invocations select none. The contract lives in the prompt and in Phase 5, where a violation can be stated specifically enough to re-prompt with.

### Phase 5 — Validate the model's output

Deterministic checks, each producing a specific re-promptable error:

- Action name exists and is exposed in the record's catalog
- All required kwargs present
- Every value the selected template renders is bound, whatever the schema declares
- No unknown kwargs
- Every kwarg value matches its declared type
- Discriminator value is a declared variant, or `_default` exists
- Rendered `injects_into` path was created in Phase 3

The third check is the one that keeps the phase boundary honest. Without it Phase 5 accepts an invocation Phase 6 then halts on, with the retry loop already skipped because nothing was wrong with the *selection* — and a phase that accepts what a later phase rejects makes the loop that could have fixed it unreachable.

On failure, re-prompt with the specific violations appended, up to a configured retry limit. This loop costs one model call — no compilation, no service startup, no test execution.

### Completeness

Every check above reads what the model returned. None of them can read what it did not return, and **a short list is valid output** — so a run that selects thirteen of fourteen correct invocations passes every rule, first try, with no retry, because nothing about it was wrong. Under-selection is the one failure class this design is otherwise blind to.

A created file states on disk what work it expects: a file template plants the markers an action's anchor targets. So the anchors a run planted, minus the anchors its selections fill, is the set of regions this run made and nothing it chose writes into. That observation is fed back to the model **once**.

**It is never a hard error.** A template may plant a region a given change does not need, or one a later record will fill, and a run that refused to proceed would make every such template a liability. If the model declines again, its answer stands and the run continues — the judgment stays the model's; it simply stops making it blind.

**An anchor nothing in the run can fill is not part of that difference.** Phase 0 already warns when a package's file templates plant a marker no action targets, so a run that reported it again would be asking the model to inject where it has just been told nothing can. It is subtracted before the re-prompt decision, silently — the author was told once, at load, by name and with the reason, and repeating it per run would charge an authoring warning by the run. Fillability is a property of the run rather than of the package that planted the marker: an action's `injects_into` renders to a path that resolves to a package by extension, so an action in one package may target a file another owns, and the set is the union across every package the run loaded.

**A completeness re-prompt is not a validation retry.** One says the answer was wrong, the other says it may be incomplete, so they are reported differently and drawn from separate budgets. A response that leaves nothing unfilled is never re-prompted, so the common case costs nothing.

Only marker anchors participate. `start_of_file`, `end_of_file`, and the match anchors name no region a file can be observed to be missing. Replay runs no completeness pass: a recording is a caller's own selection, so an unfilled anchor there is a decision rather than an oversight, and there is nobody to re-prompt.

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

Every injected region is wrapped in ownership markers naming the action that produced it, the ownership tier, the record that last parameterized it, and the kwargs it was rendered from:

```ruby
# sedum:createControllerMethod:index {"tier":"owned","record":"PR-014","kwargs":{"controller":"users","collection":"users"}}
def index
  @users = User.all
  render json: @users
end
# /sedum:createControllerMethod:index
```

Re-running a generation replaces the region an action owns rather than appending beside it. Reruns and partial regeneration are safe without a sidecar cache or resolution manifest.

This produces the audit trail as a side effect: grepping markers yields file -> action -> variant -> arguments with no maintained state, and ownership is visible in the diff a human reviews.

The action and variant stay literal on the line because that is what makes the audit trail greppable rather than parseable. Everything else is one JSON object.

### The attribute object is the extension point

The marker is a durable artifact and Sedum is not the only tool that reads one. A marker sits in a generated codebase long after the version that wrote it is gone, so a marker's reader and its writer are routinely different versions of Sedum — and, once anything is built on top of Sedum, different tools entirely.

The attributes are a JSON object rather than positional fields for that reason. A field added later is an addition rather than a migration across every repository already carrying markers.

Two promises follow, and Sedum keeps both:

**A marker is never rejected for carrying more than the reading version expects, or less.** An unrecognised key is ignored; a declared key that is absent takes its documented default. An attribute object that is not readable as JSON is still an error, because that is corruption rather than version skew.

**An unrecognised key survives the round trip.** It is retained on read and re-emitted when the region is rewritten, including when a rerun replaces the region's contents. This is the half that does not come free: ignoring a key on read and preserving it on write are different promises, and only the first one is what a decoder gives you.

The second promise is what makes the attribute object usable by anything other than Sedum. A tool built above Sedum annotates a region with its own state under its own keys, and that state is still there after the next run — without a schema change and without a sidecar file keyed by region, which would reintroduce exactly the maintained state markers exist to avoid.

Preservation is transparent. Sedum does not read, validate, interpret, or reorder an unrecognised key, and a carried key may not shadow one Sedum models. The keys Sedum models are `tier`, `record`, `writer`, and `kwargs`.

`writer` names the tool that last wrote the region. It is absent when that tool is Sedum, so every marker written before the key existed reads correctly. It exists because a demoted tier is otherwise unattributable — a region reading `seeded` gives no way to tell whether a package author declared it or another tool demoted it — and because Phase 3's *"something other than Sedum wrote it"* diagnostic is undiagnosable in a two-writer world without it.

No version token. Additive evolution is already free under the object, and a token earns its place only for a breaking change, meaning an existing field whose meaning changes.

### Ownership tiers

The tier field declares whether Sedum may overwrite a region.

`owned` — Sedum generated this region and replaces it on every run. Edits to it are lost, whoever made them. The tier is a statement about the region, not about humans; a tool editing an owned region loses its work exactly as a person does.

`seeded` — Sedum generated this region once and never touches it again. Present in the file, skipped on rerun.

An action declares its tier in `actions.yaml`, defaulting to `owned`. A template whose body is a stub a human is expected to complete should declare `seeded`; a template that fully determines its output should not.

**The marker's tier governs after the first write.** The declaration in `actions.yaml` supplies the value written the first time and has no authority over a region that already exists. A region whose tier was demoted after generation — because it stopped being Sedum's to overwrite — is precisely the region a declaration must not overrule, so the file is what decides and the declaration is only the default it started from.

This is the rule for when the `actions.yaml` key lands. Decoding is strict today and the key is not accepted yet, so the declaration currently supplies nothing; both tiers are already honored when read from a marker.

**Only `owned` is exercised by the milestones in this document.** `seeded` is specified now because markers are written to disk and read back on rerun. Adding the field later would leave every file generated in the interim carrying markers in an older shape, requiring a migration across generated codebases. The cost of reserving it is one token per marker.

### Recorded kwargs

The kwargs serialized on the opening marker make a region self-describing. A reader — human or tooling — can see what the region was parameterized with without consulting a recording or rerunning resolution.

Nothing in this document reads them back. They are written for the same reason the tier field is: markers are durable artifacts, and changing their shape after generated codebases exist is expensive.

Recording all of them, rather than the subset that currently selects a region, keeps open the question of which kwargs select and which merely parameterize. That rule lives in code and can be redefined without migrating generated codebases; the rule today is that an action's required kwargs select and its optional kwargs parameterize.

### The record attribute

The marker carries the ID of the record that last wrote a region, **as an attribute rather than as part of the region's identity**. A region is identified by its action, its variant, and its selecting kwargs; the record ID says who last wrote it.

That distinction is what lets a later record refine a region an earlier record produced, replacing it in place instead of minting a second region beside it. Under record-scoped identity, a record whose intent is *"PUT should support partial updates"* would produce two definitions of one function rather than a replacement.

**One file is still generated from one record within a single run.** Ingestion rejects a path named by two records, because Phase 4 makes one model call per record and two records naming one file would mean two independent calls deciding the same file's contents. Records that share a file are generated one at a time with `--only`, and their regions coexist and survive each other's reruns.

So the two halves are worth stating separately: refinement of a region across records and across runs is supported and is what the record attribute exists for. Two records naming one path in one run is not, and the restriction is on the whole-directory run rather than on the outcome.

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

This section states what such a tool may depend on.

> **Everything a tool built on Sedum needs is derivable from artifacts on disk. Sedum exposes no runtime API, and nothing links against Sedum.**

### The supported surface

Four things, three of them file formats and one a command.

| Surface | Direction | Who else authors it |
|---|---|---|
| **Ownership markers** | read; written within the rules below | Sedum writes them |
| **Generator package** (`sedum.yaml`, `actions.yaml`, templates) | read | users author them |
| **Recording** (JSON) | **written** and read | Sedum can write them too |
| **`sedum grow --execute`** | invoked | — |

Everything else — Sedum's packages, its internal vocabulary, the run log's shape, the `--stop-after` phase names — is internal and may change.

Four questions, and what answers each:

**"What exists?"** Grep the markers. Enumerating what has been generated, with which action, variant, and arguments, is a filesystem walk that requires nothing from Sedum's process.

**"Could this be fixed deterministically?"** Read `actions.yaml`, render each action's `injects_into` against the kwargs already recorded on the marker, and compare the result to the region in the file. Forward rendering and comparison, never inversion — `snake` and `plural` are not invertible, and no part of this surface may imply otherwise. It works only because the kwargs are recorded on the marker, which is what makes that decision load-bearing rather than merely convenient.

**"Do it."** Synthesise a recording containing exactly the invocations wanted and `--execute` it. A caller doing its own selection never touches Phase 4 at all. `--record --dry-run` gives the other direction: Sedum's selection, captured, with nothing written.

**"What did Sedum decline to do?"** The `unmanaged` patterns are in `sedum.yaml`, which is already being read. A path Sedum skipped is reported by the run and derivable from the package.

### The recording is an input format

The recording is not only a capture artifact. It is the way a caller tells Sedum what to do without a model in the loop, and replay's semantics are already right for that: **validation is identical to model-output validation, but failures are terminal.** A recording naming an action that does not exist, omitting a required kwarg, or targeting an unauthorized path fails exactly the checks a model response would fail, and the run halts rather than re-prompting — because there is nothing to re-prompt.

That is a submission protocol, not merely a convenience for hand-editing.

### Direction of control

A tool built on Sedum sits **above** it, not after it. It calls Sedum; Sedum never calls it.

There is no callback, no hook, and no plugin point anywhere in the pipeline. Adding one would invert the direction and make Sedum's guarantees depend on code it does not control — every phase after model invocation is deterministic precisely because nothing foreign runs inside it.

This is also what keeps Sedum independently runnable. Committing a recording as a standard service scaffold and replaying it needs no harness, no container, and no model server, and nothing in this section may make one a precondition.

### What a foreign writer may not do

Sedum tolerates other writers in the files it generates. The tolerance has limits, and they are stated here rather than discovered, because an unwritten constraint's first breakage looks like a Sedum bug.

- **Do not remove or relocate a marker a file template planted.** Anchors exist because a template created them; moving one breaks the injection that targets it, and Phase 7 reports a missing anchor rather than guessing.
- **Do not reorder regions within a file.** Repeated injections at one anchor accumulate in invocation order, and reordering them makes a rerun's output differ from its input for reasons Sedum cannot see.
- **Preserve unknown marker attributes in both directions.** This is the same promise Sedum makes, and it is symmetric: a writer that drops the keys it does not recognise destroys the state of every other writer, Sedum included.
- **Do not edit an `owned` region and expect the edit to survive.** Demote the region to `seeded` if it has stopped being Sedum's to overwrite. That is what the tier is for, and what the `writer` attribute makes attributable.

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

Prints a package's action catalog exactly as the model would receive it — exposed actions, kwarg schemas, descriptions, variant lists, derived requirements. The authoring feedback loop for exposure and catalog clarity.

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

**Structure inference from unmarked source.** Sedum reads structure only where a file template planted a marker. It does not parse a target language, and it does not derive regions, conventions, or injection points from source it did not shape. This is what makes "no parsers anywhere" hold.

Reading is not inference. Phase 3 reads an existing file to verify its markers are present, and enumerating what has been generated by grepping markers is a supported operation — both read what Sedum wrote, which is the distinction that matters.

**Behavioral verification.** Sedum does not run or grade the code it generates, and does not know whether the result works. A loop that generates, verifies, and repairs is a different tool that calls this one; see *The Integration Surface*.

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
