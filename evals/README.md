# Evals

Measures how well a model selects actions, and — eventually — how well what it
selected works.

**This is not part of Sedum.** Nothing under `internal/` or `cmd/` imports it, it
is not compiled into the binary, and its runner is behind a build tag so
`go test ./...` never reaches it. That separation is deliberate: Sedum does not
run or grade the code it generates, and a harness that did both would make that
non-goal untrue by adjacency if it shipped as part of the tool.

## Running it

```
OPENAI_BASE_URL=http://127.0.0.1:1234/v1 OPENAI_API_KEY=local \
  go test -tags eval ./evals -v -timeout 60m
```

| Flag | Meaning |
|---|---|
| `-eval.case <id>` | Run one case. Default: all. |
| `-eval.samples <n>` | Runs per model. Default 5. |
| `-eval.model <substr>` | Run only models whose label contains this. One row at a time when memory is tight. |
| `-eval.root <dir>` | Where the vendored fixtures live. Default `testdata`. |
| `-eval.retries <n>` | Re-prompts a rejected answer may spend. Default 0. See below before raising it. |

Without an endpoint configured the runner skips rather than fails.

## Why it exists

The measurement that corrected `prov-2026-6d87dc11` was nineteen shell loops with
`grep -c`. One run's output was lost before it could be read, two more were
mangled by a bad `sed`, and every number in that record was hand-tallied. The
finding was worth having; the method was not repeatable.

## What a number here means

**Rates, not verdicts.** Nothing fails a build on a selection rate. What is being
measured is a sampled distribution, and an assertion over a coin flip is a flaky
test wearing a suit. The runner fails only when a case produced *nothing* to
measure, which is a broken harness rather than a poor model.

**A rate is not evidence on its own.** It needs the sample size, the model, and
the build it was taken on — which is why `Report` prints all three on the header
line, and why `Compare` exists. Changing the catalog and observing a better rate
proves nothing without the other arm.

Three outcomes are distinguished, and conflating them was this harness's first
bug:

- **valid** — the model answered and Phase 5 accepted it. Scored.
- **invalid** — the model answered and Phase 5 rejected it. *A measurement*, not a
  failure. At the default budget of zero retries, this is how often one call
  produces something acceptable.
- **failed** — the run never reached the model. Excluded from every denominator,
  because an unreachable endpoint is not a model that chose badly.

The first two are told apart from the third by *type*: Phase 5 returns a
`*selection.Rejection` for an answer it refused, and the harness classifies with
`errors.As`. It matched the text of the diagnostic until `prov-2026-0811425c`,
which meant one renamed error would have silently moved samples between columns.

### The retry budget

Zero by default, and the default cannot move: every entry in `results/` was
drawn at zero and means *validated on the first call*, so a new default would
re-point that word while the field name and the report line stayed put.

Raise it with `-eval.retries` when first-call validity is varying enough to eat
the sample size of a different question. The chi case has come in at 5/5, 3/5
and 2/5 across three runs; a comparison of what a complete answer contains, made
from the two samples that survived, is a worse measurement than the same
comparison made from five.

Raising it no longer costs you first-call validity. Every sample records how
many calls it made and how many answers were rejected, so a sample with no
rejections validated on its first call whatever the budget allowed — the report
prints both lines at a raised budget, and the entry keeps the counts per sample.

What a raised budget still costs is **wall-clock comparability**. A sample that
retries pays for every rejected answer, so seconds at one budget are not seconds
at another. Compare cost in **calls** and **tokens** instead, which is what the counts are
for: a per-sample time is calls multiplied by an unknown and varying per-call
cost, and a per-call cost is tokens multiplied by an unknown token rate. The
prompt/completion split is what says whether a slow case has a long catalog or
a long answer — different problems with different responses. A server that fills
no `usage` block gets no token line rather than a line of zeroes.
The budget is in the entry and on the report header, and the validity line says
`within N calls` rather than `first call` whenever it was not one.

The completeness re-prompt is untouched by this. It draws from its own budget in
`internal/selection` and always has, so retries buy re-prompts of a *rejected*
answer and nothing else.

## The matrix

A single cell says almost nothing; every interesting question is comparative. A
case therefore declares each axis even where only one value exists today, so that
growing the matrix adds files rather than changing the schema.

| Axis | Field | Today |
|---|---|---|
| Model | `models` | one local 14B, two engines |
| Engine / quant | `models[].engine`, `.quant` | `mlx-4bit`, `llama.cpp-q4_k_m` |
| Language | `language` | `ruby`, `go` |
| Framework | `framework` | `rails`, `chi` |
| Application | `application.name` | `todo` |
| Complexity | `application.complexity` | tiers 1 and 2 |
| Generator tightness | `tightness` | `defined` |
| Arm | `arm` | `sedum` |

**Engine and quant are part of a model's identity, not decoration.** One
checkpoint under MLX 4-bit and under llama.cpp Q4_K_M uses different
quantization and does not produce identical output, so folding them into one row
would let a rate measured on one read as a claim about the other. Only the
llama.cpp engine supports continuous batching today, which is the other reason
they differ in practice.

**Language groups differently from framework.** Rails and Sinatra are one
language and two frameworks; a result holding across both says something about
Ruby rather than about either.

Two of these are declared and not yet runnable, and both matter more than the
ones that are:

**`arm: baseline`** asks the same model for the same application with no
generator package and no action vocabulary. Without that column, a good number
here is unfalsifiable — there is nothing it is good *compared to*.

**`expect.behavior`** is the fraction of the application's linespec contracts that
pass once the selection is applied and the target actually runs. Selection
completeness is a proxy; this is the thing it is a proxy *for*. It costs a
container and a database per sample rather than one model call, which is why it
is reserved rather than built.

## Fixtures

Cases read vendored generator packages and provenance records from `testdata/`,
not from projects in your workspace. That directory is named `testdata` because
the Go toolchain ignores it — the packages contain `.go` template files that are
not valid source.

See `testdata/README.md` for the layout and for why swapping generators over one
record is the comparison the structure is built around.

**Only the package and the record are vendored, never the application.**
Selection depends on nothing else: each sample runs with a fresh empty output
directory, so Phase 3 sees the same world every time regardless of where the
harness was invoked from.

## Adding a case

Drop a YAML file in `cases/`. See `cases/todo-rails-defined.yaml`; the schema is
`Case` in `case.go` and decoding is strict, so a typo is an error rather than a
silently ignored key.

Expectations are what a *complete* answer contains — taken from runs where the
model produced one, not from what it usually produces, which is the thing being
measured. **A new fixture therefore starts with none.** Run it, read the observed
selections the report prints, and set the counts from an answer that was
actually complete. Writing them by reading the generator package would make the
first run agree with them by construction.

## Known rough edges

- **Only the `sedum` arm runs.** `baseline` is declared and rejected at run time.
- **Behavior is not measured.** The schema reserves it; nothing implements it.
