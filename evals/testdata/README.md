# Fixtures

Named `testdata/` rather than `fixtures/` because the Go toolchain ignores
directories with that name. These hold `.go` and `.rb` template files that are
not valid source — a vendored Chi package broke `go build ./...` immediately,
since the toolchain tried to compile `{{path}}` as a Go identifier.

Vendored copies of the generator packages and provenance records the eval cases
run against. **Not the applications themselves** — selection depends only on the
package and the record, so the application source is neither copied nor read.

```
<application>/
  records/               the record(s): what to build
  generators/
    defined/             a generator package set: how to build it
    loose/               a second package set over the same record
```

The record directory is `records/` and the files are named for what they
describe rather than `prov-<id>.yml`. Both are deliberate: LineSpec's hooks
treat a `prov-*.yml` as one of *this* repository's records and lint it, and
these belong to other repositories. Left as they were, this repo claimed four
foreign records and refused to commit.

Sedum reads the `id:` field rather than the filename, so the records are
unchanged in every way that matters — `only:` in a case still names the real
record ID.

## Why these are vendored rather than referenced

They were previously read from sibling projects in the workspace, which made
every measurement unreproducible on any other machine and unversioned against
the package that produced it.

That already cost us once. `prov-2026-6d87dc11` records a finding that could not
be re-measured because the package had been revised in the meantime:

> The package on disk is not the one the original finding was measured
> against… Neither the before nor the after reproduces the original conditions.

A vendored fixture is a snapshot. It changes when someone deliberately changes
it, in a commit, alongside the results that cite it. The live projects stay
where real usage happens; these exist so a number can be re-run.

## Swapping generators over one record

`generators/` holds package *sets*, not packages — each subdirectory is what
`--generators` is pointed at, and the packages live inside it. That is the axis
for asking whether a tighter package selects better than a looser one over
identical intent:

```yaml
records:    todo-rails/records
generators: todo-rails/generators/defined   # change this line only
```

Both arms then differ in exactly one thing, which is what makes the comparison
mean anything.

## What can and cannot be shared

A record's `affected_scope` names real paths — `app/controllers/todos_controller.rb`
is Rails, `internal/handlers/todo.go` is Chi. So a record is **tied to its
framework** and cannot be pointed at a package for another one.

What is shareable is the layer above: *the same application idea* expressed once
per framework. `todo-rails` and `todo-api` are the same todo service under Ruby
and Go, which is what makes their rows comparable on the application axis even
though neither record could run against the other's package.

## Adding a framework

Copy a package set and its record under a new application directory, then add a
case naming them. Nothing else in the harness needs to know the framework
exists — `framework` and `language` are labels for grouping results, and Sedum
reads neither.
