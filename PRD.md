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

*Amended by `prov-2026-5403d162`:* Sedum's **pipeline** contains no language-specific knowledge. Where a target language has to be understood — generated Go must parse and be gofmt-clean, because model-written code must never land broken — the knowledge lives only in a closed, built-in **syntax registry**: each checker is named, opted into by a package's `syntax:` key, ships with Sedum, and runs no external program. The pipeline asks the registry by name and never branches on a language. A package-declared formatter command was rejected, because it would let a generator package run an executable on every grow, and nothing from a copied or shared source may run automatically.

**The model does one bounded job.** It selects actions from a closed, declared vocabulary and binds their arguments. It does not choose file paths or make structural decisions, and it writes code only into a `literal` kwarg an action declares, which a package that opts into a syntax checker has checked before anything is written. Selection over a declared schema is machine-checkable; open-ended synthesis is not.

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

A target whose comments do not run to end of line declares how they close:

```yaml
# generators/web/sedum.yaml
name: web
extensions: [.html]
comment_prefix: "<!--"
comment_suffix: "-->"
```

`comment_suffix` is optional and defaults to none. An absent suffix produces exactly the marker bytes Sedum wrote before the field existed — no trailing space and no empty delimiter — so adding it rewrites nothing already generated. A suffix declared without a prefix is rejected at load, because a closing delimiter with nothing to close is a package error rather than a default.

One package declares one comment shape. A project generating both C# and Angular templates writes two packages, which is already how a project targeting two languages is expressed.

A suffix is text that closes a comment and nothing more. There is no vocabulary of comment styles, no nesting rule, and no escaping rule — a target that forbids a sequence inside its comments, as HTML forbids `--`, leaves that constraint to the package author.

### Package resolution

Sedum builds an extension-to-package map when the generators directory loads. Every path in a provenance record's `affected_scope` resolves to a package by its extension.

**Resolution is per file, not per run.** A single record may legitimately touch `app/controllers/users_controller.rb` and `app/javascript/users.ts`; those paths resolve to different packages and each is generated under its own conventions.

If two packages in the directory claim the same extension, that is not an error at load time. Where a path with that extension appears, the package whose file template matches it most specifically wins, by the ranking in *Match specificity* applied across packages (`prov-2026-5403d162`): a command package's `cmd/{bin}/{name}_cmd.go` takes its files from a general package's `cmd/{bin}/{name}.go`. Only a path no candidate's templates describe, or one two candidates describe equally closely, is an error without a `--lang` flag, and `--lang` overrides the tie-break. The path's shape decides only through templates the packages declare, never through anything Sedum knows.

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

### Syntax checks

`syntax:` names a built-in checker every file the package touches must pass (`prov-2026-5403d162`). The registry is closed and ships with Sedum; today it holds one member, `go`, which accepts a file only if it parses, is exactly what gofmt would write, and uses no package qualifier that nothing in reach supplies — no import of the file, and no top-level declaration in a sibling file (`prov-2026-0ce03d1c`). It also rejects a function-local variable declared and never used, and an import nothing uses, each named with its line (`prov-2026-ed9124e1`); a plain `x = 1` is not a use, while `x += 1`, `x++`, `s.f = 1` and `m[k] = v` are, as the compiler counts them. That last check uses `go/parser` alone and errs toward silence: a dot import turns it off, and an import whose name cannot be inferred from its path still counts, because a false alarm would re-prompt a correct answer until its retries ran out. An unknown name rejects the package at load.

The check runs twice. In Phase 5 the validated answer is applied to its files in memory and checked; a failure is a violation like any other, fed back to the model within the retry budget — with the gofmt form of the file when it parses, so the model can bind exactly that — and terminal under replay. In Phase 7 every touched file is checked again before anything is written, so the guarantee holds for every path into injection.

A package that checks syntax shapes its markers to pass: a template ending in a blank line keeps one before its closing marker (`prov-2026-ab24df5a`), which is what gofmt wants after a top-level declaration. The reference package is `testdata/generators/go`: `setFunction`, `setMethod`, `setType`, `setVar`, `setConst` and `addImport`, each writing one whole declaration the model supplies as a `literal` into the file its `file` kwarg names, with identity file and name (and receiver). It imports with one `import "path"` declaration per region, because no marker layout inside an `import ( )` block is gofmt-clean both before and after the first import.

### Optional anchors

`optional_anchors` names anchors the package's templates plant that Phase 5's completeness check never asks about (`prov-2026-73d32e29`). They remain injection points; they are simply not work a file expects. The go package declares its `imports` anchor optional, because most declarations need no new import and asking would cost a model call on nearly every answer. An entry no file template in the package plants is a load error.

### Test paths

A package may declare which of its paths hold tests:

```yaml
test_paths:
  - "spec/**"
  - "**/*_test.go"
```

Same grammar as `unmanaged` and a record's `affected_scope`. Entries are validated at load — an empty or unreadable one is an error naming it — and **Sedum does not read them.** No phase consults the field, no output partitions on it, and no other check refers to it.

That is the unusual part, so the reason goes with it. Two consumers want the fact and neither is Sedum today: a tool above Sedum partitioning a record's regions into assertions and logic, and Sedum's own two-pass generation, for which the recording's `phases` grouping is already reserved. Both read a package's conventions, and which files hold tests is a convention — a property of a stack, not of an invocation, by the same argument that puts `unmanaged` on the package.

It has to live here rather than in either consumer's own configuration, and that is forced rather than chosen. `sedum.yaml` is decoded strictly: a package declaring a key Sedum has not modelled does not load at all, with every action in it unavailable. So a caller told to read `sedum.yaml` directly cannot ask a package author to declare anything Sedum has not declared first. Sedum governs what may appear in that file no matter who does the reading, which makes declaring a field it does not use the only way to let someone else use one.

Patterns rather than a directory, for the Go case exactly: `user_test.go` sits beside `user.go` with the same extension in the same package, so no directory-shaped declaration can separate them.

A path may match both `test_paths` and `unmanaged`, and that is not a contradiction — Sedum does not write an unmanaged path, so whether it also holds tests is the reader's question. Nor is a test path checked against the file templates: a package is not obliged to know how to write its own tests.

The risk that a declaration nothing enforces drifts from reality is real and is accepted. The mitigation is that the grammar is checked, not that the meaning is.

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

The empty file is only a problem when an action then targets it: the markers a template would have planted are not there. Phase 7 then names the cause instead of only the missing marker (`prov-2026-8394826c`): `no file template in package go matches this path, so the file was created empty and has none of the markers a template plants`, followed by the package's file templates. The go package has no root-level template on purpose, since a root file's package name is not knowable from its path. A run that fails after Phase 3 removes the files it created and the directories it had to make for them; files that existed before the run are never touched, and a successful run keeps everything.

### Anchors and marker planting

File templates plant the markers that action anchors target. This closes the loop: injection points exist because a file template created them, not because Sedum parsed the file for structure. It is the reason no language parser is required anywhere in the system.

At load time, Sedum checks each action's marker-based anchor against the markers present in its package's file templates. Because template selection is path-dependent, this cannot be a complete verification — but an action referencing a marker that appears in no template in its package is almost certainly a typo. This produces a warning.

That is the coarsest form of a question Phase 5 asks precisely. Load time knows only the package, so it can ask whether *any* template plants the marker; once an invocation names a target the same question narrows to that file, and the answer is an error rather than a warning. The load-time check is not superseded by it — a package with a mistyped anchor should not have to reach a run to find out.

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

### Targets: a convention or an argument

`injects_into` answers two questions at once — *which file does this action write to*, and *what naming convention identifies that file*. For a structure-creating action the two coincide: `createControllerMethod` only ever means a controller, and the pattern is exactly right.

For a file-agnostic action they come apart, and the pattern is a cost. An `addImport` whose template is target-independent still needs one definition per file group — `addModelImport`, `addControllerImport`, `addSerializerImport` — identical in kwargs and template body, differing only in a path. The authoring duplication is the visible cost; the real one is catalog quality, because Phase 4 is handed near-identical entries whose only discriminator is a name encoding a fact the model has to infer. It scales multiplicatively: file-agnostic actions times file groups.

So an action may declare its target to be a kwarg:

```yaml
addImport:
  kwargs:
    file:   { type: string, required: true }
    symbol: { type: string, required: true }
  injects_into: "{{file}}"
  anchor: imports
```

Both forms coexist, and no `sedum.yaml` key gates the choice. A package that wants a purely convention-driven layout declares no action shaped like this; a switch protecting an author from an action that author writes themselves would be configuration with no adversary.

**What the pattern implicitly guarded, the anchor states.** A path pattern restricted applicability as a side effect: `addBeforeFilter` could not be aimed at a migration, because the path would not render to one. A free target removes that, and what replaces it is stronger. An action's `anchor` already declares the region kind it needs, and a file template declares which regions the files it creates carry, so **an action is applicable to a file whose matched template plants its marker, and to no other.** That states the precondition rather than approximating it by directory, and it catches the case a pattern cannot: a path that matches the pattern but whose file carries no such region.

It reframes an action's applicability as anchor vocabulary rather than path shape, which is a more honest description of what these actions are. They operate on a region kind, not on a directory.

**`affected_scope` is still the safety boundary, and the path pattern never was.** A free target changes nothing about Phase 5's authorized-path rule.

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

`pascal`, `camel`, `snake`, `kebab`, `upper`, `lower`, `plural`, `singular`, `prefix:X`, `suffix:X`, `join:X`, `quote`

`quote` renders its input as a complete double-quoted string literal: the value JSON-encoded, without HTML escaping (`prov-2026-f63c340a`). That is a valid string literal in Go, JavaScript, TypeScript, Java and C# for any input, so a template writes `{{usage|quote}}` with **no quotes of its own**, and a quote or backslash the model bound cannot break the generated file. Writing `"{{usage}}"` instead leaves the file one model-supplied `"` away from not compiling, and no Phase 5 check can see it; a kwarg description asking the model to avoid quotes would be a rule in a description, which this document rules out.

`join:X` is the one exception to `string -> string`: it takes a **list** and renders its members joined by the separator `X`, which is the only way a `list` kwarg reaches a file. A value still holding a list when it reaches the template is an error — including a bare `{{words}}`, which would otherwise emit Go's own slice formatting. `join` over an empty list is an error too, since the empty string is a plausible-looking value everywhere a joined list is going.

The separator cannot be a pipe, because the pipe separates transforms. Join into an array literal and let the target language build an alternation, which is also the only place each member can be escaped: `['{{words|join:','}}']`.

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

Records are decoded by the shared provenance record library (`github.com/livecodelife/provenance/record`), v1 and v2 alike, so every tool reads a record the same way (`prov-2026-5f1e4f0a`). The library owns the schema and rejects a malformed value in any field it models; Sedum keeps one leniency on top, because the schema is not its own (`prov-2026-dad4d877`): it hands the library only the keys it reads — `id`, `intent`, `constraints`, `affected_scope`, `forbidden_scope`, `created_at`, `aliases`, `status`, `type`, `implements`, and `version`, which decides the schema — so a field it does not read can never fail a run. A malformed value in one it reads is an error naming the file and the line. Sedum also requires an `id`, which the library does not.

### Phase 2 — Resolve paths to packages

For each authorized path, first check it against the union of every package's `unmanaged` patterns; a match is recorded and skipped. Then resolve its extension to a generator package, applying `--lang` where the extension is contested. A path whose extension no package claims is a hard error. A contested extension with no disambiguating flag is a hard error naming both candidate packages.

The unmanaged check runs *before* extension resolution, not after, because the paths most often declared unmanaged are the ones no extension can reach.

### Phase 3 — Create files from templates

For each authorized path: match it against its package's `files/` tree, bind captured segments, render the matched template, and write the file. Fall back to `_default`, then to an empty file with a log line.

**A file that already exists is reconciled with its template, or the run halts.** Sedum never re-renders a template over an existing file — that would destroy whatever has been injected into it. What it does is compare the two and write in whatever the template adds.

The file is decomposed first: its **skeleton** is the file with every marked region removed, so what is compared is the boilerplate rather than the boilerplate plus what has been injected into it. The skeleton is aligned against the rendered template on **structural** lines — a line is structural unless it is blank or begins with the package's declared `comment_prefix`, and a marker is structural because it is the thing being planted. Only the opening delimiter takes part: a closing one says nothing about whether a line is boilerplate. That declaration is the only knowledge involved; nothing here parses a line.

**If every difference is something the template adds, the template is applied. Otherwise the run halts**, naming the file and the first line the template has no counterpart for. Formally the condition is that the file's structural lines are a subsequence of the template's: if each appears in the template in order, the template contains everything the file does and the difference is what it adds.

This covers four situations that used to be one error:

- A file carrying everything its template declares is left exactly as it is.
- An **empty** file takes its whole template. It has nothing to destroy, which is the reason this phase is create-if-absent to begin with.
- A file **another tool wrote** takes the template where the template accounts for it — a scaffolded routes file gains its anchor; a controller carrying hand-written actions does not, because whether those should become owned regions is not Sedum's to decide.
- A file an **earlier run generated** takes what its template has gained since. Without this a package author's anchor vocabulary is frozen after first use.

The unit is the template, not the marker. A file template is the boilerplate its author wants every time, so planting an anchor while refusing the keyword it depends on would produce a file worse than either applying the template or refusing it.

Two properties follow from applying only insertions. Whitespace between two lines the file and the template share is the same whitespace, so the template's arrangement of it governs; where the file has anything else there — a comment its author wrote — that is theirs, and what the template adds goes after it rather than over it. And an edit script that only inserts cannot tell a template that grew from a file that lost a line, so a line deleted from a file's boilerplate is restored. That is the intended reading: the operation brings a file's skeleton back to its template, and it is what enforces the rule against removing a planted marker.

**What was written into each existing file is reported.** Applying a template can add lines a team did not write, and a team adopting Sedum into a repository has to be able to read what it did rather than discover it.

This makes reruns safe, which matters because stopping and resuming a run is a normal workflow rather than an edge case. A second run over a reconciled file changes nothing.

**Nothing is created that a provenance record did not authorize.** There is no sibling expansion and no inference of companion files. If a C++ record names `src/user_controller.cpp` but omits `include/user_controller.hpp`, the header is not created, and the injection targeting it fails loudly in Phase 7.

This is a governance position. `forbidden_scope` means Sedum does not touch what it was not authorized to touch; conjuring convenient files would violate that regardless of how obvious the convention seemed. Completeness of `affected_scope` is the record author's responsibility.

### Phase 4 — Model invocation

One call per provenance record. The prompt contains the record's `intent`, its `constraints`, the paths created for it in Phase 3, and the action catalog — the **union of exposed actions across every package the record's paths resolved to**, with their kwarg schemas, authored descriptions, variant lists, per-variant derived requirements, and `injects_into` patterns.

That union is then narrowed to what the record could legally use before it is embedded: an action is removed when no binding of its kwargs could render `injects_into` onto one of the record's authorized paths. The match is conservative in one direction only — every `{{...}}` expression, transformed or bare, stands for any value, so only an action's fixed, literal text has to fail to appear in any authorized path before it is removed. Nothing has been bound yet, so nothing is removed on the strength of a guess about what will be. A free-target action (`injects_into: "{{file}}"`) has no fixed text and is never removed by this rule; a composite is removed only if at least one child's pattern is unreachable, since every child's rendered target has to be authorized for the composite to validate at all. If filtering leaves nothing to offer, Sedum fails before the call, naming the record and its authorized paths, rather than sending an empty catalog. Every removal is written to the run log. **Phase 5 validates against the full catalog regardless** — this narrows what is shown, never what is checked — and `--full-catalog` turns the narrowing off, reproducing the unfiltered prompt.

Variant lists are included deliberately. Without them there is an invisible cliff: `name: index` gets a full implementation while `name: search` falls to `_default`, and the model has no way to know it fell off. Exposing the list lets it prefer covered values where intent maps cleanly, and take the fallback knowingly where it does not. Whether a `_default` exists is carried alongside, because *knowingly* is not available to a model that cannot see whether there is a fallback to take.

`injects_into` patterns are included for a harder reason. Without them the catalog names a kwarg and the file list names a path, and nothing connects the two, so a model asked to bind `controller` binds the path it was shown. Recovering the kwarg from a rendered path would require inverting `snake` and `plural`, which is precisely what nothing in this system does. With the pattern present the reasoning runs forwards: match its literal segments against an authorized file and bind what is left. For a pattern-targeted action the model still never chooses a path — it binds arguments such that the package's own pattern lands on a file the record authorized. An action that declares its target to be a kwarg has no pattern to invert and is outside that rule rather than an exception to it.

**Each entry carries its `anchor`, and the file list carries the injection points each created file plants.** This is the same move one layer over: there, the model matched a pattern against the authorized paths and bound what was left; here it matches an anchor against the markers those paths carry and names the file that has one. Forward and mechanical in both cases, and in neither does Sedum contribute target knowledge — the markers are the package author's strings, passed through exactly as the pattern and the kwarg names already are. The planted set is read from what Phase 3 actually rendered, through the same computation the completeness pass uses.

**The prompt suppresses a bare `injects_into`; the JSON keeps it.** `["{{file}}"]` tells the model nothing it can act on and invites the inversion the pattern exists to prevent. It is not noise to a caller reading `sedum actions --json`, though: a present-but-bare pattern is how that caller tells *this action's target is a parameter* apart from *this action has no target pattern*, and dropping it deletes the only signal separating them. So the catalog is the contract and the prompt is a view of it — which was already true, and this is the first change to make it explicit. It sets the rule for later ones: **additions to the catalog are additive to the JSON, and what any single consumer is shown is that consumer's rendering. A prompt-side omission is never implemented as a schema change.**

The response is **structured output, not tool calls** — a JSON array of `{action, kwargs}` objects, wrapped as `{"invocations": [...]}` so that the document's root is an object. This keeps the mechanism working with models that lack tool-calling support, which is most of the open-weight range worth evaluating.

Sedum does not ask the server to constrain the response by default. Grammar-constrained decoding takes the shortest legal completion at the array's first token, and an empty array is legal, so requesting a schema made a model that had been selecting three correct invocations select none. The contract lives in the prompt and in Phase 5, where a violation can be stated specifically enough to re-prompt with.

`--response-schema` makes that constraint available, opt-in, for an endpoint a caller has chosen to point it at. On, the same catalog Phase 4 already filtered to the record's authorized paths is compiled into a JSON Schema for `{"invocations": [...]}` — one `oneOf` branch per exposed action, discriminated actions expanding to one branch per declared variant with the discriminator pinned by `const`, `required` read from exactly what Phase 5's own `derivedRequirements` already treats as required (declared-required kwargs, `requires`, and, per variant, `variant_requires`) so the two cannot drift into independently maintained copies of the same rule — and sent as `response_format: {"type": "json_schema", "json_schema": {name, schema}}` alongside the request. The array itself carries no `minItems`: forcing at least one invocation is exactly the constraint that collapsed selection above, so `{"invocations": []}` stays legal under `--response-schema` exactly as it is without it. An optional kwarg is never modeled as nullable, since Phase 5 rejects `null` for one left unbound. Off, the request is byte-for-byte what it is today. On, and the endpoint rejects `response_format`, `sedum grow` fails clearly, naming `--response-schema`, rather than silently retrying unconstrained. Phase 5 validates every response exactly as it does today either way — the schema constrains the model, it does not replace the check. Whether the bundled `goinfer-serve` local backend accepts `response_format`, or reproduces this section's own collapse if it does, is unverified; combined with `--local-model` or the bundled default, `--response-schema` fails clearly rather than silently running unconstrained.

### Phase 5 — Validate the model's output

Deterministic checks, each producing a specific re-promptable error:

- Action name exists and is exposed in the record's catalog
- All required kwargs present
- Every value the selected template renders is bound, whatever the schema declares
- No unknown kwargs
- Every kwarg value matches its declared type
- Discriminator value is a declared variant, or `_default` exists
- Rendered `injects_into` path was created in Phase 3
- The target file's matched template plants the marker the action is anchored to

The last check is applicability, and it cannot tighten anything: Phase 7 already hard-errors on that exact case, so checking it here is the same failure, earlier, and re-promptable. Its diagnostic distinguishes a file whose template plants no injection points at all from one that plants some but not this one — the two have different fixes, and a `_default` template or an empty-file fallback makes the first case legitimate. It speaks only about files this run rendered: a path Sedum produced no content for is one it has observed nothing about, and Phase 7 remains the backstop there.

The third check is the one that keeps the phase boundary honest. Without it Phase 5 accepts an invocation Phase 6 then halts on, with the retry loop already skipped because nothing was wrong with the *selection* — and a phase that accepts what a later phase rejects makes the loop that could have fixed it unreachable.

On failure, re-prompt with the specific violations appended, up to a configured retry limit. This loop costs one model call — no compilation, no service startup, no test execution.

**Phase 5 reports what it spent**: how many calls the record cost, how many answers were rejected, whether the completeness observation was one of them, and the prompt and completion tokens the server accounted for. The loop is the only thing that can tell those apart after the fact, and a caller left to infer cost from elapsed time is reading calls multiplied by an unknown per-call cost. A record with no rejections validated on its first call, whatever the retry limit allowed — which is what keeps first-call validity measurable at any budget. An answer refused every time it was asked is returned as a typed rejection carrying those counts and every attempt's violations, so a rejected answer and an unreachable endpoint are distinguishable without matching on text.

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

Marker comments are the load-bearing case. Marker syntax uses the package's declared `comment_prefix`, since `#`, `//`, and `--` all appear across targets, plus its `comment_suffix` where the target closes its comments. A reader strips the closing delimiter before interpreting what remains: the attribute object is the tail of an opening marker, and it is not readable as JSON with `-->` still attached.

---

## Ownership and Idempotency

Every injected region is wrapped in ownership markers naming the action that produced it, the ownership tier, the record that owns it, the hashes of that record and of the package it was rendered from, and the kwargs it was rendered from:

```ruby
# sedum:createControllerMethod:index {"tier":"owned","record":"PR-014","gen":"sha256:…","pkg":"sha256:…","kwargs":{"controller":"users","collection":"users"}}
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

Preservation is transparent. Sedum does not read, validate, interpret, or reorder an unrecognised key, and a carried key may not shadow one Sedum models. The keys Sedum models are `tier`, `record`, `gen`, `pkg`, `writer`, and `kwargs`, emitted in that order.

`writer` names the tool that last wrote the region. It is absent when that tool is Sedum, so every marker written before the key existed reads correctly. It exists because a demoted tier is otherwise unattributable — a region reading `seeded` gives no way to tell whether a package author declared it or another tool demoted it — and because Phase 3's *"something other than Sedum wrote it"* diagnostic is undiagnosable in a two-writer world without it.

### The gen and pkg hashes

`gen` and `pkg` say whether a region is current, from the marker alone, so that a later `grow` can decide per record whether to skip, re-render, or re-select without a manifest or lockfile (`prov-2026-66bea9f3`).

- **`gen`** is the `gen_hash` of the record that wrote the region, as it was at write time: SHA-256 over the canonical encoding of exactly the fields Sedum reads — `intent`, `constraints`, `affected_scope`, `forbidden_scope`. Status, sealing and every other field stay out, so a lifecycle change never looks like a content change. The encoding is defined by the shared provenance record library and its corpus (`conformance/hashing/cases.json` in the provenance repository), not by Sedum.
- **`pkg`** is the hash of everything that decides what the action renders for fixed kwargs: its resolved declaration (kwargs with types, required flags and defaults; discriminator; variants; `injects_into`; anchor fields; identity), the content of each of its templates, and the package's `transforms` and `op_exceptions`. Descriptions are left out, because they reach the model and never the renderer. A region written by a composite carries its child's `pkg`, since the marker names the child. The encoding is pinned in `conformance/pkghash/cases.json`.

Both are written as `sha256:` and 64 lowercase hex digits, and both are omitted when empty, so a marker written before they existed is unchanged by being read and rewritten. **Empty means unknown, and unknown must be read as stale, never as current.** A replay with no recorded `gen` and no `--records` writes none. Neither hash takes part in a region's identity: a hash says which version of a region this is, never which region.

The inflection table ships with Sedum rather than with a package, so an upgrade that changes an irregular plural is not reflected in `pkg`.

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

Recording all of them, rather than the subset that selects a region, is what kept the question of which kwargs select and which merely parameterize open long enough to answer it. An action may declare an `identity`: the kwargs that decide which region an invocation *is*. Where it declares none, every required kwarg selects and every optional one parameterizes, which is what every package authored before the field existed was written against.

The declaration exists because required and selecting are different properties that happened to coincide. A kwarg is required because the model must state it; it selects because it says which region is meant. For an action whose required kwargs name a target — the symbol an import binds — those are the same set. For one whose required kwargs are mostly content — a pattern, a selector, a count — they are not, and selecting on all of them means changing a value mints a second region beside the first instead of refining it. The seeded tier and retained marker attributes are both read only after a region matches, so they are lost with it (prov-2026-d71d6f76).

### The record attribute: ownership

The marker carries the ID of the record that **owns** the region, **as an attribute rather than as part of the region's identity**. A region is identified by its action, its variant, and its selecting kwargs; the record ID says whose it is (`prov-2026-375f5d4a`, replacing the earlier last-writer rule of `prov-2026-36c8a99c`).

Keeping the record out of identity is still what stops two records producing two definitions of one function. What changed is what happens when a second record selects a region the first one owns: it is an **error naming both records**, not a silent takeover. Under versioned records a record keeps its id when it changes, so another record selecting the same region means someone wrote a new record where they should have revised the existing one. Refining a region is done by revising its owner.

- A record owns a region when the marker names its id or one of its `aliases` (ids folded into it by migration). A rewritten region's marker names the current id.
- A region whose marker names no record is unowned, and the first record to select its identity claims it.
- Ownership changes only by an explicit handover, which is not yet implemented.

**Reconciliation.** After a record's selection validates and expands, `grow` reconciles it against what is on disk by identity: a region it selects again is replaced in place, a new identity is inserted, and an `owned` region of this record that it did not select again is **pruned**. Seeded regions, other records' regions and unowned regions are never pruned, and only the files the record authorizes in this run are considered. Pruning is part of Phase 7's single all-or-nothing write, so any failure anywhere in the run leaves every region alone. A record skipped without a model call reconciles nothing.

An empty selection from a record that already owns regions is an **error unless `--force`**, because an empty reply is more often a flaky model than an intent to delete a record's code. `--dry-run` reports those regions as held for `--force`.

Replay (`--execute`) checks ownership but never prunes: a recording may be partial, and turning replay into reconciliation is opt-in (*Open Design Questions* §16).

**Records sharing a file grow in one run.** The rule that rejected a path named by two records is gone; ownership now protects the code region by region. Two records selecting one identity in the same run is still an error naming both.

**Order at a shared anchor.** A new region goes after every region at its anchor whose owner sorts at or before its own owner, and before the rest. Owners sort by the record's `created_at`, then id; a region naming no record sorts first, and an owner the run does not know sorts by id alone. One owner's regions keep the order they were applied in, and no existing region is moved, so regenerating a file from scratch reproduces it byte for byte whatever order its records were grown in.

### What grow does per record

`grow` is safe to run on every record, every time (`prov-2026-9621daed`). After Phase 3 it reads the regions each record owns in the managed files it authorizes and decides, from the markers alone:

| The record's owned regions | `grow` | Model call |
| --- | --- | --- |
| none yet | select (a normal grow) | yes |
| some region's `gen` differs from the record's, or is unknown | re-select, then reconcile | yes |
| `--force`, or `--record` | re-select, then reconcile | yes |
| every `gen` and every `pkg` matches | skip: nothing expanded, no file written | no |
| every `gen` matches, some `pkg` differs | re-render those regions in place from their markers | no |

A re-rendered region is rebuilt from the action and kwargs on its own marker, which is a one-line recording; nothing is pruned, because nothing was re-selected. A region that can no longer be re-rendered - its action is gone, or it would now render into another file - sends its record to re-select instead. Seeded regions decide nothing and are never rewritten, but a seeded region whose record's `gen` has moved is reported so a person can look.

So a second `grow` over unchanged records makes no model call and writes no file.

**`--check`** classifies every record the same way and reports it as current, stale, to re-render, or unknown (no regions yet), with no model call and nothing written, Phase 3's files included. It exits non-zero when any record is stale or needs re-rendering; unknown does not fail. It suits CI and pre-commit hooks, and cannot be combined with `--execute`, `--record`, `--dry-run`, `--force` or `--stop-after`.

**`--record` re-selects every record.** A recording captures selections, and a record that was skipped or re-rendered made none. Its invocations cannot be rebuilt from markers either: a composite's regions name its unexposed children, which recording validation rejects.

**A limit, by decision.** A marker records every kwarg its region was rendered from, including package defaults and `--var` values bound at the time. Re-rendering replays them verbatim, and neither `gen` nor `pkg` covers them, so a changed kwarg default or `--var` value reaches an existing region only with `--force`.

### Retired records, orphans, and stamping

**Deprecated records are retired** (`prov-2026-cc73cf61`). A record's `status` is the one lifecycle field Sedum reads, and only the value `deprecated`: such a record takes no new versions, is never grown, and owns nothing. `--only` naming one is an error.

**Orphans.** An orphan is an `owned` region whose marker names a record that is not in the records directory, by id or alias, or that is deprecated. Every `grow` reports orphans: a summary per outcome (counts, and the directories holding the most), or with `--orphans` a line per region. `--prune-orphans` removes them, inside Phase 7's all-or-nothing write, but only in files some current record authorizes; Sedum writes nothing else, so an orphan in any other file is reported and left alone. Seeded and unowned regions are never orphans. `--check` lists orphans and fails on them: an orphan is code no current record stands behind.

**`--stamp`** writes the current `gen` onto a record's regions whose marker has none, and the action's current `pkg` onto those with none, where the owning record exists, and changes nothing else — no body, no value already present, no record id. It makes no model call and runs nothing else. Running it asserts that the code is current, so a repository whose markers predate the hashes can skip a full re-select after migration. Because it writes only missing values, it never hides a real change. It cannot be combined with `--check`, `--execute`, `--record` or `--force`; with `--dry-run` it reports what it would stamp. A marker naming an alias keeps it until its region is next written, when the surviving id replaces it.

**Prompt context.** For each file a record may write, the selection prompt lists the regions other records already own there — action, identifying kwargs, owner — under *Already present*, as not to be selected; selecting one would be refused in Phase 7. When a record is re-selected and no `--prior-recording` supplies its earlier selections, they are rebuilt from its own markers as context, keeping exposed actions and declared kwargs only. Hashes never reach the prompt.

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
      "gen": "sha256:…",
      "created_at": "2026-08-14",
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

A record's `gen` is the `gen_hash` of the record version its invocations were selected from. It is optional: replay writes the recorded value on every marker, takes it from `--records` when the recording has none, and otherwise writes none (`prov-2026-b9298192`). The recorded value wins even when the record has changed since, because it is the version the selection was made against; a later `grow` sees the change and re-selects.

A record's `created_at` is recorded too, also optional, because owners at a shared anchor are ordered by it: replay takes it from the recording, overlaid by `--records` when given, so a replay without `--records` lays out a shared anchor exactly as the grow it recorded (`prov-2026-43b7ec39`).

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

Read strictly: the first-party harness is included, and a shared module that Sedum and a consumer both link against is not an available answer to any gap in the formats. Proposing one is a signal that a format is insufficient.

The reason travels with the rule, because a constraint recorded without its justification reads as fastidiousness and is relaxed by whoever first finds it inconvenient. **The formats must be sufficient without linking, and the only thing that forces them to be is that nobody gets an exemption.** Once a first-party consumer can link Go, every awkward format is cheaper to route around with one more exported helper than to fix. Third-party consumers inherit the gap, and the gap is invisible to the party best placed to notice it, because that party is not experiencing it. The same party authors these formats and is their most privileged consumer, so the rule that we eat the same food is the only thing standing between the surface and quiet insufficiency.

### The supported surface

Six things: three file formats, two commands, and the corpus that makes one of the formats checkable. Counting an addition rather than folding it into an existing row is deliberate: a surface that grows silently is one nobody audits.

| Surface | Direction | Who else authors it |
|---|---|---|
| **Ownership markers** | read; written within the rules below | Sedum writes them |
| **Generator package** (`sedum.yaml`, `actions.yaml`, templates) | read | users author them |
| **Recording** (JSON) | **written** and read | Sedum can write them too |
| **`sedum grow --execute`** | invoked | — |
| **`sedum render`** | invoked | — |
| **`conformance/`** | read | Sedum's own tests read it too |

Everything else — Sedum's packages, its internal vocabulary, the run log's shape, the `--stop-after` phase names — is internal and may change.

Five questions, and what answers each:

**"What exists?"** Grep the markers. Enumerating what has been generated, with which action, variant, and arguments, is a filesystem walk that requires nothing from Sedum's process.

**"Could this be fixed deterministically?"** Ask Sedum where an invocation lands: `sedum render --package <name> --action <name> --kwargs <json>` renders that action's `injects_into` against the kwargs already recorded on the marker, and the caller compares the result to the region in the file. Forward rendering and comparison, never inversion — `snake` and `plural` are not invertible, and no part of this surface may imply otherwise. It works only because the kwargs are recorded on the marker, which is what makes that decision load-bearing rather than merely convenient.

The rendering is a command rather than an instruction to reimplement, and the difference is the whole of the previous paragraph. `sedum.yaml` carries the declarative pipelines and the `op_exceptions` table, so those are derivable — but the inflection table is not, because it ships per language rather than per package, and neither are the word-splitting semantics the case operations share nor the leading-case match applied to an inflected word. A caller answering this question by hand must reimplement `plural`. It diverges on the first irregular, renders a path that does not match the region, and concludes that no deterministic fix exists when one does: a wrong answer wearing the shape of a correct refusal, which is worse than an error because nothing reports it.

Transform behaviour is therefore exposed as behaviour and not as a data dump. Publishing the inflection table would sanction the reimplementation this exists to prevent.

**"Do it."** Synthesise a recording containing exactly the invocations wanted and `--execute` it. A caller doing its own selection never touches Phase 4 at all. `--record --dry-run` gives the other direction: Sedum's selection, captured, with nothing written.

**"What did Sedum decline to do?"** The `unmanaged` patterns are in `sedum.yaml`, which is already being read. A path Sedum skipped is reported by the run and derivable from the package.

**"Which of these paths hold tests?"** The `test_paths` patterns, from the same file. Sedum declares the key and validates its grammar without reading it, because strict decoding means a package cannot carry a key Sedum has not modelled — so a field Sedum does not use is the only way a caller gets one. It is the first surface here that exists purely for a consumer, and the shape it takes is the shape any future one will: a declaration Sedum checks and does not interpret, never an opaque namespace.

### The recording is an input format

The recording is not only a capture artifact. It is the way a caller tells Sedum what to do without a model in the loop, and replay's semantics are already right for that: **validation is identical to model-output validation, but failures are terminal.** A recording naming an action that does not exist, omitting a required kwarg, or targeting an unauthorized path fails exactly the checks a model response would fail, and the run halts rather than re-prompting — because there is nothing to re-prompt.

That is a submission protocol, not merely a convenience for hand-editing.

### Direction of control

A tool built on Sedum sits **above** it, not after it. It calls Sedum; Sedum never calls it.

There is no callback, no hook, and no plugin point anywhere in the pipeline. Adding one would invert the direction and make Sedum's guarantees depend on code it does not control — every phase after model invocation is deterministic precisely because nothing foreign runs inside it.

This is also what keeps Sedum independently runnable. Committing a recording as a standard service scaffold and replaying it needs no harness, no container, and no model server, and nothing in this section may make one a precondition.

### The conformance corpus

The rules below are prose, and prose constraining a durable on-disk format with more than one writer is the kind of constraint that gets discovered rather than followed. `conformance/markers/cases.json` states the checkable ones as data: golden marker lines, their expected parse results, and the expected **bytes** after emitting, re-emitting, and replacing. `conformance/pkghash/cases.json` pins the `pkg` hash's encoding, so a tool above Sedum can tell a region rendered from an older package without linking Sedum.

The marker is the only format in this table that a foreign tool writes into files Sedum will later re-read and rewrite with no validation gate in between. A misread `sedum.yaml` fails on the reader's side; a malformed recording is refused terminally at ingestion. A malformed marker corrupts state silently, and the first breakage looks like a bug in Sedum.

Three properties make it worth being a surface rather than a test fixture. It is at the repository root rather than under `testdata/`, because `testdata/` is a Go convention that reads as one project's internal business and a stranger has no reason to look in it. It is language-neutral, because the audience is a writer who is not using Go — a corpus only a Go implementation can consume checks the one implementation that does not need checking. And **Sedum's own tests read it** rather than carrying a parallel copy, because a corpus that can go stale without failing a test will, and a stale corpus is worse than none: the party checking against it cannot tell.

Every case cites the record that decided the behavior it pins, and a case is added when a decision is made rather than to raise coverage. Two of the rules below are file-level and cannot be expressed as a parse-and-emit golden; the corpus says so out loud rather than letting its existence imply they are covered.

### What a foreign writer may not do

Sedum tolerates other writers in the files it generates. The tolerance has limits, and they are stated here rather than discovered, because an unwritten constraint's first breakage looks like a Sedum bug.

- **Do not remove or relocate a marker a file template planted.** Anchors exist because a template created them; moving one breaks the injection that targets it, and Phase 7 reports a missing anchor rather than guessing.
- **Do not reorder regions within a file.** Regions at one anchor are ordered by owner and then by invocation order, and reordering them makes a rerun's output differ from its input for reasons Sedum cannot see.
- **Do not change a region's `record`.** It names the region's owner, which is the only record that may rewrite or prune it. Ownership changes by an explicit handover, never by editing the marker.
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
| `--var <name>=<value>` | Bind a variable a generator package declares. Repeatable. |
| `--record <path>` | Write a recording of the run to the given path. Re-selects every record, so the recording holds a selection for each. |
| `--execute <path>` | Replay a recording. Skips model invocation. Mutually exclusive with `--record`. |
| `--dry-run` | Run every phase, write nothing. Reports the files that would be created, the injections that would be applied, and the regions that would be pruned. |
| `--force` | Re-select every record even when its markers say it is current, and let a record whose selection comes back empty prune every region it owns (without `--force` that is an error). Ignored with `--execute`. |
| `--check` | Report each record as current, stale, to re-render, or unknown, and any orphan regions, from markers. No model call, nothing written; exits non-zero when any record is stale or needs re-rendering, or any orphan exists. With `--execute`, `--check` is a dry replay: it writes nothing and exits non-zero when replaying the recording would change the tree (a stale `gen`, a changed region, or a region the recording no longer selects). A recording written by `--record` without `--only` is complete: replaying it prunes the regions of the records it names that it no longer selects. A partial recording prunes nothing and reports them as `unselected`. |
| `--prune-orphans` | Remove regions whose record was removed or deprecated, in files some current record authorizes; report the rest. Ignored with `--execute`. |
| `--restamp-gen` | With `--execute` and `--records`: stamp each recorded record's current `gen` onto its regions in place of the `gen` the recording carries, so a replay after a record edit reads as current. The recorded selection must still pass replay's own checks. Asserts the selection still satisfies the edited record. |
| `--orphans` | Print every orphan region, one line each, in place of the per-fate summary (counts and the busiest directories) that `grow` and `--check` print by default. The `--json` report always lists them all. |
| `--json` | Print the run's report as the `sedum.grow/1` JSON document (`conformance/formats/`) instead of text. Each record says whether its regions changed — inserted, pruned, or rewritten with different bytes — so a caller runs only the specs whose code moved. |
| `--stamp` | Write the current `gen` and `pkg` onto markers missing them, where the owning record exists, and nothing else. No model call. Cannot be combined with `--check`, `--execute`, `--record` or `--force`. |
| `--stop-after <phase>` | Halt after the named phase. One of `resolution`, `files`, `invocations`, `expansion`. |
| `--retries <n>` | Model output validation retry limit. Default 3. Ignored with `--execute`. |
| `--model <name>` | Model identifier. Endpoint and credentials come from environment. Ignored with `--execute`. Mutually exclusive with `--local-model`. |
| `--local-model <path>` | Path to a gguf file to serve locally instead of `--model`. Starts a goinfer-serve subprocess for the run and stops it when the run ends. Ignored with `--execute`. Mutually exclusive with `--model`. |
| `--local-model-server <path>` | Path to the `goinfer-serve` binary. Defaults to one found on `PATH`. Requires `--local-model`. |
| `--local-model-backend <name>` | Compute backend passed to goinfer-serve's own `-backend` flag (`cpu`, `metal`, `cuda`, `webgpu`). Applies to `--local-model` when given, and to the bundled default model otherwise. Ignored with `--execute`. |
| `--log <path>` | Run log location. Defaults to `.sedum/run.log`. |
| `-v, --verbose` | Mirror the run log to stdout. |
| `--full-catalog` | Show the model every exposed action, including ones no authorized path can reach. Ignored with `--execute`. |
| `--response-schema` | Compile the record's filtered catalog into a JSON Schema and send it as `response_format`, constraining the model's completion. Off by default; off, the request is byte-for-byte what it is today. Fails clearly if the endpoint rejects `response_format`, or if the client in use (the bundled local default among them) does not support it at all — never falls back unconstrained. Ignored with `--execute`. |

`--record` and `--dry-run` compose. Together they capture a recording without writing any generated files — the model runs, its output is validated and saved, and nothing is created. Without `--dry-run`, the run records and executes in the same pass.

`--execute` and `--dry-run` also compose, validating a recording against the generator packages and reporting what it would produce without writing.

With `--execute`, `--records` becomes optional and enables scope validation when supplied. `--lang` is ignored, since package resolution is recorded per file.

With neither `--model` nor `--local-model`, a run built from an official release archive (via `install.sh` or a downloaded archive) falls back to a `goinfer-serve` binary and a fine-tuned gguf bundled as siblings of the `sedum` binary itself, starting and stopping that server the same way `--local-model` does. `go install` and a from-source build carry neither file, so that fallback fails there with an error naming `--model`/`--local-model` rather than silently picking a wrong path.

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

### `sedum regions`

```
sedum regions --generators ./generators [--output .] [--records ./provenance] [--record <id>]... [--json]
```

Lists every region whose marker names a record: its file, action, identifying kwargs, tier, and the `gen` and `pkg` hashes it was written with (`prov-2026-bf01cd6b`). Read-only — no model, nothing written, no records directory required — and it never reads hidden directories, the generators directory or the records directory. With `--records`, a marker naming an alias is reported under the surviving id, and `--record` matches either. An empty hash means the marker has none, which a consumer must read as unknown; `grow --stamp` fills them in.

`--json` prints the `sedum.regions/1` document described in `conformance/formats/`. It is what the provenance library's evidence code hash is computed from, since that library never parses a marker.

| Flag | Description |
|---|---|
| `--generators <path>` | Generators directory. Required. |
| `--output <path>` | Directory to read markers from. Default `.`. |
| `--records <path>` | Records directory. Optional. |
| `--record <id>` | List only this record's regions. Repeatable. |
| `--json` | Print the `sedum.regions/1` JSON document. |

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
| `--var <name>=<value>` | Bind a variable a generator package declares. Repeatable. |
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

### `sedum render`

Prints the file an action's invocation would land in, by rendering that action's `injects_into` pattern against supplied kwargs with the package's own transform engine. This is how a tool built on Sedum answers "could this be fixed deterministically?" without reimplementing the transforms; see *The Integration Surface*.

Rendering is forward only. `snake` and `plural` are not invertible, so there is no way to ask which kwargs produce a given path.

A composite has no pattern of its own and reports its children's, in execution order — one selection visibly touching two files.

```
sedum render --generators ./generators --package rails --action createControllerMethod \
  --kwargs '{"controller":"Order","name":"index"}'
```

| Flag | Description |
|---|---|
| `--generators <dir>` | Generators directory. Required. |
| `--package <name>` | Package declaring the action. Required. |
| `--action <name>` | Action to render the target of. Required. |
| `--kwargs <json>` | Bound arguments, as a JSON object. Defaults to `{}`. |
| `--json` | Emit the result as JSON rather than formatted output. |

It is its own command rather than a mode of `sedum actions`. That command answers "what can this package do", which is a per-package listing; this one answers "where does this invocation land", which is per-invocation and takes kwargs. Overloading the first with the second would make a listing command require an action and a kwargs blob to answer its second question.

### `sedum --version`

Emits a bare semver string. It is the same value a recording carries in `sedum_version`, so the binary's answer and the artifact's cannot drift.

```
sedum --version
```

A caller that depends on a command depends on a version of it. Under the linking rule it cannot read a constant out of the package, so the surface has to answer the question.

---

## Logging

There is no plan artifact. The execution sequence is a consequence of the record and the configuration, not a decision any component makes.

The run log records package resolution, file template matches and captures, the model's raw response, validation failures and retries, composite expansions, resolved paths, selected variants, transform resolutions, and anchor matches. It also records which actions Phase 4 removed from a record's prompt before the call — the record, the action, and why — so a package author can see why an action never reached the model without reading the prompt itself. It is diagnostic output, cleared per run. Nothing depends on it — idempotency state lives in the ownership markers, in the generated files.

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
