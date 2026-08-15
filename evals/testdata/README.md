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
  records/               one record: what to build
  generators/
    defined/             a generator package set: how to build it
    loose/               a second package set over the same record
```

`generators/` holds package *sets* so the package can be varied while the record
is held fixed. `records/` is flat because each fixture needs exactly one — Sedum
makes one model call per record, so a second would measure a second selection
and charge the eval for it.

The record directory is `records/` and the files are named for what they
describe rather than `prov-<id>.yml`. Both are deliberate: LineSpec's hooks
treat a `prov-*.yml` as one of *this* repository's records and lint it, and
these belong to other repositories. Left as they were, this repo claimed four
foreign records and refused to commit.

Sedum reads the `id:` field rather than the filename, so nothing about how a
record is ingested changes.

## Fixtures may deviate from their source projects

`todo-api`'s three records are condensed here into one. Sedum makes **one model
call per record**, so a fixture carrying three measures three different
selections and charges the eval for all of them to learn about one. Sequencing
records against a converging codebase is the harness's job, not Sedum's.

The condensed record therefore has its own id — `eval-todo-chi`, not the source
record's — because it is no longer that record and should not claim to be. Where
a fixture diverges, the reason is written at the top of the file it diverges in.

The rule: a fixture exists to make a measurement mean something, not to mirror a
project. Where those pull apart, the measurement wins and the deviation is
recorded.

## Records carry functionality, and nothing else

Both records were cut back to what the API does. The source records also
specified the database adapter, how the connection is configured, whether the
service retries at boot or waits on a compose healthcheck, CSRF handling, the
Docker build, and where the test suite lives.

All of that is real, and none of it belongs in a record used for this
measurement, for one reason: **no action in either package can act on any of
it.** Those files are produced whole by file templates, so a constraint about
them is text in the prompt that no selection can satisfy or violate.

The boot-ordering constraint was the clearest case. Rails put the wait in
`docker-compose.yml`, Chi put it in `main.go`, both hardcoded in file templates
— and the two records prescribed *opposite* answers, so a model satisfying one
had violated the other. That looked like a reason for a third "parity" record.
It was not: stripping both records to functionality removed the disagreement
along with everything else no action served, and two records now reach parity
where three were needed before.

`affected_scope` is trimmed the same way, to the paths an action injects into —
six in each. Authorizing a file no invocation can target only puts it in the
prompt's file list, and in the Chi fixture it also planted two anchors nothing
could fill, which cost a completeness re-prompt on every sample.

The test for this is `TestTheLanguageArmsAreControlled`: it fails if the two
cases ever differ in more than the stack.

