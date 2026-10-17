# Design eval

How much *design* does the cheap model do when an intent states only behaviour?
Today every Haiku grow in the herd is a transcription: the intent carries the
exact code. This eval asks what happens when it does not.

Not part of Sedum: nothing under `internal/` or `cmd/` imports it.

## What it runs

Five small tasks (`testdata/tasks/<id>/`), each a Go file the model must write,
graded by hidden tests it never sees (the exact tests the herd's own code already
passes). Three intent variants per task:

- `behaviour`: what the code does, in prose, with the identifiers it must provide.
- `signatures`: `behaviour` plus the declarations, bodies elided.
- `code`: the full code to transcribe (the control, today's practice).

For each sample the harness builds a throwaway Go module (the task's context
files, the `go` generator package, one open record whose intent is the variant),
runs `sedum grow` through the herd's `sedum-haiku` wrapper (claude-haiku-5-5 only;
there is no model flag), then copies the hidden tests in and runs `go build` and
`go test`. Retries, repairs and tokens come from the grow log the grow wrote.

## Running it

```
go build -o /tmp/sedum ./cmd/sedum
go run ./evals/cmd/design -sedum /tmp/sedum -dry      # the plan, nothing run
go run ./evals/cmd/design -sedum /tmp/sedum           # 75 grows, K=5
go run ./evals/cmd/design -sedum /tmp/sedum -n 1 -task since -variant code
```

`-n` is capped at 5: raising K needs approval. Results go to
`evals/results/design-<time>.json` with the commit, the model, K, every sample
and a hash of each task, so a result can be re-run and compared.

A default `go test` never reaches a model: the tests use a fake grow and a fake
wrapper script. (There is no `eval` build tag: the go generator package's file
template cannot write a build constraint, so the guarantee is that no test uses
the real wrapper.)

## First result (2026-10-16, 75 grows, 533,861 tokens)

| variant | passed | retries | repairs | tokens |
|---|---|---|---|---|
| behaviour | 25/25 | 5 | 5 | 208,148 |
| signatures | 24/25 | 4 | 4 | 189,175 |
| code | 25/25 | 0 | 0 | 136,538 |

Read it as a direction, not a number (5 samples a cell). Behaviour-only intents
passed every sample, at about 1.5x the tokens of transcription, and almost all of
the retries and repairs sit in one task (`render`, whose spec is a precise table
format) plus one `append`. The caveat that matters: these behaviour intents are
detailed (field names, JSON keys, column order, rounding), written by the same
author as the tests, so this measures how well Haiku turns a precise prose spec
into code, not how well it invents a design from a vague one. A vaguer-intent arm
is the obvious next measurement.
