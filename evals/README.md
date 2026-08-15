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

## Reading the history

The results files are append-only and pinned to commits, and `history` is how
they are read back:

```
go run ./evals/cmd/history                    # every case
go run ./evals/cmd/history todo-chi-defined   # one case
```

```
todo-rails-defined  (7 entries)
  date       commit     n  r  c  valid              first call         calls  tok/call    wall     tok/s
  2026-08-15 3e38d7c*   1  0  1  1/1 [0.21,1.00]    -                  -      -           1m5.5s   -
~ 2026-08-15 90fe770*   5  2  2  5/5 [0.57,1.00]    5/5 [0.57,1.00]    1.00   -           7m1s     -
~ 2026-08-15 42a426f*   4  0  1  4/4 [0.51,1.00]    4/4 [0.51,1.00]    1.00   1799+605    4m24.4s  36.4
  * tree was dirty: not re-runnable from that commit
  ~ interval overlaps the previous entry: these runs do not distinguish each other
  - not recorded: the entry predates that field
```

`n`/`r`/`c` are samples, retry budget, and concurrency. **A dash is not a
zero** — it means the entry was written before that field existed, and the
schema's additive promise is that such an entry keeps meaning what it meant.

The `~` column is the one to read first: it marks entries whose interval
overlaps the entry above, which is the harness saying *these two runs do not
distinguish each other*. A history with `~` down the whole column is a history
in which nothing measurable has changed.

It reads committed files only — no endpoint, no model, no build tag — and never
writes.

## Why it exists

The measurement that corrected `prov-2026-6d87dc11` was nineteen shell loops with
`grep -c`. One run's output was lost before it could be read, two more were
mangled by a bad `sed`, and every number in that record was hand-tallied. The
finding was worth having; the method was not repeatable.

## Reading a report

A real run, annotated:

```
todo-chi-defined  model=qwen2.5-coder-14b-instruct/llama.cpp-q4_k_m  arm=sedum  tightness=defined  n=5  retries=2
  wall 18m46.1s  per-sample fastest 4m22.9s / mean 7m24s / slowest 9m59.3s  concurrency 2  (2.0x over sequential)
  valid within 3 calls: 5/5 [0.57,1.00]
  valid first call: 4/5 [0.38,0.96]
  cost: 9 call(s) over 5 sample(s), mean 1.80  (3 completeness observation(s))
  tokens: 24.1k prompt + 3.2k completion, per call 2680 + 356
  throughput: 24.2 tok/s at concurrency 2 (prompt 21.4, completion 2.8)
  action                    want           selected              exact    mean   first
  addQueryTest                 5    4/5 [0.38,0.96]    4/5 [0.38,0.96]    4.00      0%
```

**Header** — what was measured. `n` is independent runs of the same case;
`retries` is the budget each got. The model is `id/engine-quant` because one
checkpoint under two engines is two rows, not one.

**`wall` / per-sample spread** — total elapsed, then how the individual samples
were distributed. The fastest-to-slowest spread is the *throttling* detector: on
a fanless machine, later samples much slower than earlier ones means the box is
heat-limited and every number is soft.

**`(2.0x over sequential)`** — the sum of sample times over the wall clock. Read
it as *how much the samples overlapped*, *never* as a speedup. If per-sample
latency inflates by the same factor you are overlapping, this reads 2.0x while
the run took exactly as long — which is what the concurrency sweep below
measured on this hardware.

**`valid within N calls` / `valid first call`** — how many samples produced
something Phase 5 accepted, and how many needed no retry to do it. **The gap
between the two lines is what the retry budget rescued.** The denominator
excludes samples that never reached the model.

**`cost:`** — model calls per sample. The decomposition of `9` above is 5 first
calls + 1 retry + 3 completeness observations, and the mean is the number to
compare against another arm's.

**`tokens:`** — prompt billed against completion produced. **Billed is not
computed**: this server reuses a cached prefix, so a 1799-token prompt costs one
token of work. The prompt figure is what a hosted endpoint would charge and what
a rate limit counts; it is not where the local time goes.

**`throughput:`** — *completion* tokens over the wall clock, at that
concurrency. Completion only, because those are the tokens the machine produced;
counting billed prompt tokens made a case look faster the larger its prompt was.
Comparable across runs *only* when the server was run the same way both times.

**The action table** — `want` is what a complete answer contains. `selected` is
samples that included the action at all; `exact` is samples that included
exactly `want` of them. **The gap between those two is the diagnosis**: equal
means all-or-nothing (the model either does it or forgets it), while `selected`
above `exact` means it is showing up short. `first` is how often the action
opened the answer, which sounds trivial and was once the whole story — a dropped
action appeared if and only if it appeared first, meaning the model was not
weighing and rejecting it, it was never arriving there.

### Which number answers which question

| You want to know | Look at |
|---|---|
| Is one call enough? | `valid first call` |
| Are the answers complete? | the `exact` column |
| Missing entirely, or showing up short? | `selected` against `exact` |
| What does a run cost? | `cost:` calls per sample |
| *Why* is this case slow? | the `completion` half of `tokens:` — the prompt is cached |
| Is the box throttling? | the fastest-to-slowest spread |
| Did my change do anything? | the `~` marks in `history` |

## What a number here means

**Rates, not verdicts.** Nothing fails a build on a selection rate. What is being
measured is a sampled distribution, and an assertion over a coin flip is a flaky
test wearing a suit. The runner fails only when a case produced *nothing* to
measure, which is a broken harness rather than a poor model.

**Every rate carries a 95% Wilson interval**, and the fraction stays beside it
because the sample size is what produced the width. At five samples, 4/5 is
`[0.38, 0.96]` and 5/5 is `[0.57, 1.00]` — overlapping, so a run that moved from
one to the other has not shown you anything. `Compare` marks overlapping rows
with `~`. There are deliberately no significance tests: an interval states the
measurement honestly, while a p-value is a verdict wearing a number, and the
decision belongs to the reader.

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

## The model server

**The harness does not configure the server, and cannot see how it was
configured.** The endpoint is a variable on purpose — a local server, a hosted
API, and a colleague's tuned box are all the same code path. The consequence is
that an entry records the endpoint URL and the concurrency the harness asked
for, and nothing about slots, batch size, or prefix caching. So a **throughput
comparison across two entries is only a comparison when the server was run the
same way both times**, and that fact belongs beside the numbers rather than in
someone's memory.

What the entries in `results/` were actually measured against, read off the
running server rather than reconstructed:

```
$ lms ps
IDENTIFIER                  MODEL                       STATUS  SIZE     CONTEXT  PARALLEL  DEVICE
qwen2.5-coder-14b-instruct  qwen2.5-coder-14b-instruct  IDLE    8.99 GB  16292    4         Local

$ curl -s localhost:1234/api/v0/models | jq '.data[] | select(.state=="loaded")'
  quantization: Q4_K_M   arch: qwen2   compatibility_type: gguf
  max_context_length: 131072   loaded_context_length: 16292
```

LM Studio, serving GGUF Q4_K_M over an OpenAI-compatible endpoint, 16292 tokens
of loaded context and 4 parallel slots. `PARALLEL 4` is what bounds
`-eval.concurrency`: asking for more than that queues rather than parallelizes,
which is why the c=8 row in the sweep below is not a different experiment from
c=4.

**Run those two commands rather than trusting this block.** A configuration
written down from memory is worse than none: a wrong number is caught by
re-running, while a wrong account of how a number was produced survives the
re-run, because the re-run reproduces what the document says instead of what
happened (`prov-2026-6f728c95`).

**What the server does not expose, this does not claim.** Batch policy and cache
configuration are not in either output, so they are unknown here rather than
assumed.

### Prefix caching is already on, and there is no flag for it

A case's prompt is the record plus the catalog, and it is **byte-identical
across every sample of that case** — samples differ only in the model's
sampling. The llama.cpp backend under LM Studio exploits that automatically: it
picks a slot by longest-common-prefix similarity and keeps everything up to the
divergence point.

From the server log during a rails run:

```
slot get_availabl: selected slot by LCP similarity, sim_best = 1.000 (> 0.100 thold)
slot update_slots: new prompt, n_keep = 1799, task.n_tokens = 1799
slot update_slots: n_past was set to 1798
prompt eval time =     746.79 ms /     1 tokens
       eval time =  194531.80 ms /   603 tokens (322.61 ms per token, 3.10 tokens per second)
```

**One token of 1799 is evaluated.** Across a day of runs, 63 prompt evaluations
computed a single token, three computed the full prompt as cold starts, and the
rest were completeness re-prompts, which share the prefix and diverge after it.

Two consequences worth keeping straight:

- **Generation is essentially the whole cost** — 746 ms against 194.5 s, or
  99.6%. A slow case is slow because of what the model writes, not what it reads.
- **Billed is not computed.** `usage` reports the full prompt length on a cache
  hit, so the harness cannot see how much was recomputed, and the `tokens:` line
  says `prompt billed` rather than implying work.

There is nothing to enable and no flag to pass: `lms load` exposes context
length, parallelism, GPU ratio, and TTL, and no cache options, because the reuse
is structural.

### What the concurrency sweep found

Six runs of `todo-rails-defined`, four samples each, on the invocation above:

| Order | `-eval.concurrency` | Wall | Per-sample mean | tok/s |
|---|---|---|---|---|
| 1 | 1 | 4m24.4s | 1m06.1s | 36.4 |
| 2 | 2 | 4m15.2s | 2m07.4s | 37.7 |
| 3 | 4 | 4m01.7s | 4m00.9s | 39.8 |
| 4 | 8 | 3m16.8s | 3m16.4s | 48.8 |
| 5 | 4 | 3m04.8s | 3m04.6s | 52.0 |
| 6 | 1 | 5m06.9s | 1m16.7s | 31.3 |

Runs 5 and 6 are controls, repeating earlier levels after everything was warm.
They are why this table says less than it appears to.

**Two identical c=4 runs differed by 31%** (rows 3 and 5), and the slowest run of
all six was the last and warmest one. So there is no warm-up trend, and the
run-to-run noise is as large as any effect the sweep could claim.

What survives: **every parallel run finished faster than every sequential one**
(185–255s against 264–307s, disjoint ranges), and per-sample latency scales
almost exactly with concurrency — at c=4 each sample takes essentially the whole
wall clock. Read that as suggestive rather than established: with two runs per
condition, that ordering arises by chance about 7% of the time.

**The practical advice.** Leave concurrency at 1 while iterating: time-to-first
result is 1m06s instead of 3m05s, so a broken fixture surfaces in a minute
rather than after the whole run. Raise it to about the server's `-np` when you
want a whole measurement sooner and will wait for all of it either way. Do not
quote a wall-clock difference under ~30% from single runs of each condition.

### Timing is soft; the counts are not

The same six runs, by what stayed still:

| Quantity | Spread across six runs |
|---|---|
| Wall clock | 3m04s – 5m07s (±31%) |
| Throughput | 31.3 – 52.0 tok/s (±25%) |
| Tokens per call | 1799+605 → 1799+604 (±0.2%) |
| Calls per sample | 1.00 every time (exact) |

This is the argument for `cost:` being denominated in calls and tokens rather
than seconds. The selection measurements are reproducible on this hardware; the
timing measurements are worth ±30%, and a comparison inside that band is not a
result. A comparison well outside it — the chi arm at 2.7x the rails arm — still
is.

### Throughput is reported, not estimated

The `throughput:` line divides the counted tokens by the wall clock. It is never
a model's advertised token rate, which describes a single stream on the vendor's
hardware and says nothing about this server at this concurrency.

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
