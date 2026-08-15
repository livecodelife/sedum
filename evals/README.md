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
| `-eval.root <dir>` | Where the fixture applications live. Default `../..` — `go test` runs with cwd at `evals/`, so this is the workspace. |

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
  failure. Retries are zero here on purpose, so this is how often one call
  produces something acceptable.
- **failed** — the run never reached the model. Excluded from every denominator,
  because an unreachable endpoint is not a model that chose badly.

## The matrix

A single cell says almost nothing; every interesting question is comparative. A
case therefore declares each axis even where only one value exists today, so that
growing the matrix adds files rather than changing the schema.

| Axis | Field | Today |
|---|---|---|
| Model | `models` | one local 14B |
| Framework | `framework` | `rails` |
| Application | `application.name` | `todo` |
| Complexity | `application.complexity` | tier 1 |
| Generator tightness | `tightness` | `defined` |
| Arm | `arm` | `sedum` |

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

## Adding a case

Drop a YAML file in `cases/`. See `cases/todo-rails.yaml`; the schema is
`Case` in `case.go` and decoding is strict, so a typo is an error rather than a
silently ignored key.

Expectations are what a *complete* answer contains — taken from runs where the
model produced one, not from what it usually produces, which is the thing being
measured.

## Known rough edges

- **Invalid-vs-failed is classified by matching an error string.**
  `internal/selection` exports no error type for a rejected response. The honest
  fix is a sentinel there, and it should land before anything depends on the
  validity rate.
- **Only the `sedum` arm runs.** `baseline` is declared and rejected at run time.
- **Behavior is not measured.** The schema reserves it; nothing implements it.
