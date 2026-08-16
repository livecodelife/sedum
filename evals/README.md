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
go run ./evals/cmd/eval                                    # every case, coarse
go run ./evals/cmd/eval todo-rails-described               # one case
go run ./evals/cmd/eval -res fine -model llama.cpp \
  todo-rails-defined todo-rails-described                  # the description A/B
go run ./evals/cmd/eval -dry -res fine todo-rails-defined  # print the plan, run nothing
```

| Flag | Meaning |
|---|---|
| `-res <r>` | The question being asked: `smoke`, `coarse` or `fine`. Default `coarse`. Sets the sample size. |
| `-model <substr>` | Run only models whose label contains this. One row at a time when memory is tight. |
| `-n <n>` | Runs per model. Default: whatever the resolution calls for. A count *below* it is refused. |
| `-c <n>` | Samples in flight at once. Default 1. |
| `-retries <n>` | Re-prompts a rejected answer may spend. Default 0. See below before raising it. |
| `-timeout <d>` | What one case is allowed to take. Default: derived from the samples about to be drawn. |
| `-dirty` | Run against a dirty tree anyway; the entry records as not re-runnable. |
| `-dry` | Print the plan and the invocation, run nothing. |

It prints what the run will cost before spending it, and the exact `go test`
line it is about to execute — it is a wrapper, and every flag it passes is one
the runner already accepts.

```
fine resolution, 2 case(s):
  todo-rails-defined       1 model(s) x 30 samples   ~45m     (timeout 1h30m)
  todo-rails-described     1 model(s) x 30 samples   ~45m     (timeout 1h30m)
                           ~1h30m in total
```

**`~45m` is what the run takes; `1h30m` is when it should be considered hung.**
Two numbers because they answer different questions, and the total is of the
expectations — a total of the ceilings said three hours for a pair that costs
ninety minutes. The estimate is 90 seconds a sample, the top of the 66–90s
observed on the local 14B row, and the ceiling is twice it. `-timeout` replaces
the ceiling in either direction and the plan then says `given`, so a number
nobody derived is never printed as though it had been.

The estimate is a constant rather than a rate read back from `results/`,
which does hold the wall clock of every previous run. A plan that changed with
the last run could not be reproduced from the invocation that printed it, and a
new case has no history to read — so the constant has to exist anyway, and a
second path used only sometimes is a second thing to be wrong.

Three things it does that the raw invocation does not:

**It sizes the timeout from the run.** `go test` defaults to ten minutes and a
fine run takes forty-five. A forgotten `-timeout` does not fail fast; it kills
the run partway and no entry is written.

**It refuses a dirty tree.** An entry taken mid-edit is recorded as not
re-runnable from its commit, which is a thing to discover before spending
forty-five minutes rather than after. `-dirty` says otherwise deliberately.

**It stops between cases when the previous entry is uncommitted.** The first
arm's own result file is untracked when it finishes, and `clean` is read from
`git status --porcelain`, so the second arm would record dirty *because the
first one succeeded*. It names the file to commit and stops. It commits nothing
itself — that commit needs a provenance tag the command has no way to choose.

<details>
<summary>The underlying invocation</summary>

```
OPENAI_BASE_URL=http://127.0.0.1:1234/v1 OPENAI_API_KEY=local \
  go test -tags eval ./evals -v -timeout 60m
```

| Flag | Meaning |
|---|---|
| `-eval.case <id>` | Run one case. Default: all. |
| `-eval.resolution <r>` | The question being asked: `smoke`, `coarse` or `fine`. Default `coarse`. |
| `-eval.samples <n>` | Runs per model. Default: whatever the resolution calls for. |
| `-eval.model <substr>` | Run only models whose label contains this. |
| `-eval.root <dir>` | Where the vendored fixtures live. Default `testdata`. |
| `-eval.retries <n>` | Re-prompts a rejected answer may spend. Default 0. |
| `-eval.concurrency <n>` | Samples in flight at once. Default 1. |
| `-eval.results <dir>` | Where results are appended. Empty disables recording. |

</details>

The endpoint defaults to `http://127.0.0.1:1234/v1` and the key to `local`, and
neither overrides what the environment already sets. Without an endpoint
configured the runner skips rather than fails.

## Reading the history

The results files are append-only and pinned to commits, and `history` is how
they are read back:

```
go run ./evals/cmd/history                    # every case
go run ./evals/cmd/history todo-chi-defined   # one case
```

```
todo-rails-defined  (7 entries)
  date       commit    res       n  r  c  valid              first call         calls  tok/call    wall     tok/s
* 2026-08-15 3e38d7c*  -         1  0  1  1/1 [0.21,1.00]    -                  -      -           1m5.5s   -
? 2026-08-15 90fe770   -         5  2  2  5/5 [0.57,1.00]    5/5 [0.57,1.00]    1.00   -           7m1s     -
s 2026-08-15 1c4d0aa   smoke     2  0  1  2/2 [0.34,1.00]    2/2 [0.34,1.00]    1.00   1802+601    2m9.8s   35.9
  2026-08-15 42a426f   coarse    4  0  1  4/4 [0.51,1.00]    4/4 [0.51,1.00]    1.00   1799+605    4m24.4s  36.4
~ 2026-08-16 302f081   coarse    5  0  1  5/5 [0.57,1.00]    5/5 [0.57,1.00]    1.00   1799+604    5m6.9s   31.3
  * tree was dirty: not re-runnable from that commit, and not compared
  ? no resolution stated: a sample size nobody chose, and not compared
  s smoke: plumbing only at that sample size, not a measurement and not compared
  ~ interval overlaps the previous comparable entry: these runs do not distinguish each other
  - not recorded: the entry predates that field
```

`res`/`n`/`r`/`c` are the resolution, sample size, retry budget, and
concurrency. **A dash is not a zero** — it means the entry was written before
that field existed, and the schema's additive promise is that such an entry
keeps meaning what it meant. A dash under `res` is *unstated*, not `coarse`:
those runs drew five samples because five was the default, and stamping a
decision on them now would invent one.

The `~` column is the one to read first: it marks entries whose interval
overlaps the previous comparable entry, which is the harness saying *these two
runs do not distinguish each other*. A history with `~` down the whole column is
a history in which nothing measurable has changed.

### Only some entries are compared

An entry joins the comparison chain when it is **clean and states a
resolution**. Everything else is printed in full, marked, and skipped — the mark
withholds the comparison, not the measurement, and the chain skips an excluded
entry rather than resetting at it, so a run is still compared against the last
entry that could be stood on.

| Mark | Why it is not compared |
|---|---|
| `*` | The tree was dirty, so the entry is not re-runnable from its commit. |
| `?` | No resolution was stated, so the sample size was nobody's decision. |
| `s` | A smoke run. Two samples overlap almost everything. |

Each was a claim being made from evidence the entry itself records as unusable.
Before this held, every entry in `results/` was dirty *and* unlabelled and each
one was compared against the last, so the `~` marks between them said "these
runs do not distinguish each other" about runs that could not be re-run and
whose sample size nobody chose (`prov-2026-c5ad54ff`).

The `*` in the commit column and the `*` in the mark column are separate
statements: the first says an entry cannot be re-run from its SHA, the second
says it took no part in the comparison. A row can carry both, and the two
answer different questions.

**Nothing is deleted to achieve this.** Results are append-only
(`prov-2026-eb283c56`) — a dirty entry is an honest record of a run taken
mid-edit, and the ones in this repository are the evidence behind the
concurrency sweep below and the sample size fixed for the description
comparison.

It reads committed files only — no endpoint, no model, no build tag — and never
writes.

### What a sample keeps

Beyond the counts, each stored sample carries what produced it:

| Field | On | What it is |
|---|---|---|
| `counts` / `first` | valid | invocations per action, and which one opened |
| `rules` | invalid | the rule slugs it was rejected under, attempt order, repeats kept |
| `invocations` | valid | every action **with the arguments bound to it** |
| `calls` / `rejected` / tokens | both | what it cost |

`rules` exists because `details` cannot be counted. The description A/B produced
7 rejections in one arm and 11 in the other, and every one of them rendered as
*"the model's output did not validate within 1 attempt(s)"* — so the run could
say the model failed more often and not whether it failed *differently*. A slug
can be tallied; the prose stays in `details` for reading.

`invocations` exists because counts are a projection. Every failure that
motivated kwarg descriptions was a correctly *selected* action with a *wrong
argument*, which a count cannot see — and the arguments were in memory one line
before being discarded. Keeping them also makes an entry re-scorable: a rule
invented later runs against samples drawn before it existed, instead of
re-drawing two arms at eighty minutes a pair every time a question is sharpened.

**Both are stored and neither is scored.** No expectation, report column or rate
is derived from them yet — a scoring rule is its own decision, and the change
that first makes data visible should not also decide what to conclude from it
(`prov-2026-2256e6fa`). Entries written before these fields carry neither and
read as *not recorded*, exactly as the earlier additions do.

## Why it exists

The measurement that corrected `prov-2026-6d87dc11` was nineteen shell loops with
`grep -c`. One run's output was lost before it could be read, two more were
mangled by a bad `sed`, and every number in that record was hand-tallied. The
finding was worth having; the method was not repeatable.

## Reading a report

A real run, annotated:

```
todo-chi-defined  model=qwen2.5-coder-14b-instruct/llama.cpp-q4_k_m  arm=sedum  tightness=defined  res=coarse  n=5  retries=2
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
`retries` is the budget each got; `res` is the question that `n` was drawn for,
and `res=smoke` means the numbers below it are a plumbing check rather than a
measurement. The model is `id/engine-quant` because one checkpoint under two
engines is two rows, not one.

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
| How many samples does my question need? | the resolution table below |

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

### Sample size is a property of the question

`-eval.samples` defaulted to 5 and every early run used it, for no reason anyone
chose — it was the first run's number and everything since inherited it. What
five buys is less than it looks:

| Perfect cell | 95% interval | Width |
|---|---|---|
| 5/5 | `[0.57, 1.00]` | 0.43 |
| 10/10 | `[0.72, 1.00]` | 0.28 |
| 20/20 | `[0.84, 1.00]` | 0.16 |
| 30/30 | `[0.89, 1.00]` | 0.11 |

**A perfect score at five samples is consistent with a true rate of 57%.** That
is not an argument for always drawing thirty. It is an argument for choosing the
number against the question, because these questions differ by an order of
magnitude in the effect they are looking for. So a run states its **resolution**
and the sample size follows:

| Resolution | n | The questions it answers |
|---|---|---|
| `smoke` | 2 | Does the plumbing work — a new case loads, a new package resolves, an endpoint answers. **Not a measurement.** |
| `coarse` | 5 | Differences that are enormous. Does a 4b model select usefully where a 14b does; does a framework's package work at all; is one arm several times slower than another. |
| `fine` | 30 | Moving a rate that is already high — 4/5 to 5/5, or a selection rate from 80% to 100%. |

The chi arm coming in at 2.7× the rails arm is a coarse result, and five samples
established it comfortably. Whether kwarg descriptions move a selection rate is
a fine one, and five samples cannot see it either way.

*Resolution*, not *tier*, because a case already has a tier: `complexity` is an
ordinal tier of the application. Two axes sharing a word read as one axis in a
results file.

**A smoke rate is never cited as a measurement.** It says so on its own report,
in any `Compare` it appears in, and in `history` long after whoever ran it has
forgotten which it was.

**A run drawn below its resolution is refused before the first call**, because
both ways of getting this wrong are only expensive afterwards. A comparison that
publishes "no effect" from a run too small to have seen one is worse than no
comparison — the temptation to then keep sampling until the number moves is the
specific failure the refusal exists to prevent. Oversampling is honoured rather
than clamped: it is a cost decision and misleads nobody.

**Raise `n` per question, never globally.** At n=5 on the local 14b row,
`todo-rails-defined` is about 7 minutes and `todo-chi-defined` about 19, and a
cell scales linearly in n:

| | n=5 (coarse) | n=30 (fine) |
|---|---|---|
| `todo-rails-defined` | ~7 min | ~45 min |
| `todo-chi-defined` | ~19 min | ~2 hr |

Four local model rows over both cases is under 2 hours at the coarse tier and
most of a day at the fine one. That arithmetic is the real constraint on
widening the matrix: a matrix run at n=30 everywhere buys resolution for
questions that never needed it, and spends the hours the questions that do need
it were going to require.

**Thirty is an upper bound, not a threshold.** It comes from comparing
intervals, which is deliberately conservative — a proper two-proportion
comparison separates two rates earlier than their intervals stop overlapping. A
`~` mark means *not distinguished by these samples*, never *shown to be the
same*. What the conservatism does not change is the direction of the error at
five: too few for a fine question by a wide margin either way.

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
| Catalog descriptions | which set `generators` names | present, absent |
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

### The description comparison

`todo-rails-defined` and `todo-rails-described` name one record and two package
sets, and the second set is the first with `description` fields added to its
actions and kwargs. It is the harness's first real job (prov-2026-ac15ed2b),
because descriptions shipped on four hand-observed wrong bindings and have had
no measurement behind them since.

One model row at a time, because a 24GB machine holds one 14B resident and the
two arms have to face the same one:

```
go run ./evals/cmd/eval -res fine -model llama.cpp \
  todo-rails-defined todo-rails-described

# it stops after the first arm; commit that entry, then:
go run ./evals/cmd/eval -res fine -model llama.cpp todo-rails-described

go run ./evals/cmd/history todo-rails-described
```

About 45 minutes an arm on the local 14B row, by the timing table above.

**Commit each arm's entry before running the next.** `clean` comes from `git
status --porcelain`, which counts untracked files, so the first arm's own result
file is what makes the tree dirty for the second — and an arm that records
`clean: false` is not pinned to a commit and cannot be re-run against the same
fixture. This is how every entry already in `results/` came to be dirty. The
command stops rather than letting it happen; it does not commit for you.

**Thirty per arm, fixed before the run.** The rails arm has come in at 4/4 and
5/5 across its whole history, so this is a fine question by the taxonomy below:
moving a rate that is already at the ceiling. At five samples 4/5 is [0.38,0.96]
and 5/5 is [0.57,1.00] — overlapping, so a run that size cannot see the effect
in either direction while its null result reads exactly like a finding.

**Both arms are drawn fresh.** The entries already in `results/` are not the
undescribed arm: every one is five samples or fewer, none states a resolution,
and all were taken against a dirty tree.

#### What it found

Both arms drawn 2026-08-16 at n=30, fine, retries 0, concurrency 1, on
`qwen2.5-coder-14b-instruct/llama.cpp-q4_k_m`, against clean trees —
`634f02a` and `e1849dd`.

| | undescribed | described |
|---|---|---|
| First-call validity | 23/30 `[0.59,0.88]` | 19/30 `[0.46,0.78]` |
| Exact selection, among valid | 23/23 `[0.86,1.00]` | 18/19 `[0.75,0.99]` |
| Prompt tokens per call | 1799 | 2677 |
| Wall clock | 39m40s | 40m29s |

**Nothing moved that these samples can distinguish.** Both pairs of intervals
overlap, so on the question as posed — do descriptions move the selection rate
— this is a null result, and it is recorded as readily as a positive one would
have been.

**Selection had nowhere to go.** Every valid sample in the undescribed arm
matched the expectation exactly, all six action counts. That is the ceiling the
sample size was chosen for, and it means a positive result was only ever
available as an improvement in *validity*, not in selection.

**Both numbers are lower with descriptions, not higher.** The described arm's
one off-target valid sample asked for three endpoints and three controller
tests instead of five. Four fewer answers validated on the first call. Neither
gap clears its interval, so the honest statement is "not distinguished" rather
than "worse" — but nothing here points the way the feature was argued for.

**They are not free.** 878 more prompt tokens a call, about half again, for a
rate that did not move.

**What this does not say.** One application, one model, one framework: a
difference here would have been a difference *there*, not a property of
descriptions. And selection is not binding — every wrong binding that motivated
`prov-2026-c5697387` produced a correctly selected action with a bad argument
and would have scored a perfect count in this table. This is evidence for
`prov-2026-6d87dc11`'s claim that enriching the catalog does not fix
under-selection. It is not a verdict on the feature descriptions were added
for, which needs the rendered output checked and is deferred with
`expect.behavior`.

**It measures selection, not binding.** Every wrong binding that motivated
descriptions produced a correctly *selected* action with a bad *argument*, and
would have scored a perfect catalog count. Whether descriptions help the model
bind needs the rendered output checked, which is deferred with `expect.behavior`.
A moved rate here is not validation of that feature.

The two sets are held to differing only by descriptions mechanically:
`TestTheDescribedSetDiffersOnlyByDescriptions` loads both and compares them field
by field with descriptions excluded, so a renamed kwarg, a flipped `required`, a
new variant, or an edited template fails the default suite.

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

### `expect.bindings` — and the opposite discipline

Counts say the model reached for the right actions. `expect.bindings` says
whether it reached for them *correctly*:

```yaml
expect:
  actions: {addColumn: 2}
  bindings:
    addColumn:
      because: >-
        The record: "the todos table carries exactly two columns of its own:
        title, a string that is not null and carries no default, and completed,
        a boolean that is not null and defaults to false."
      key: [name]
      invocations:
        - {name: title,     type: string,  nullable: false, default: "nil"}
        - {name: completed, type: boolean, nullable: false, default: "false"}
```

Scored **per kwarg**, so the report names the broken argument rather than saying
an invocation was wrong — and printed as its own table, never folded into
selection. Four smoke samples scored a perfect 6/6 on selection while every one
of them bound `default` wrong; a combined rate would have hidden exactly that.

| Rule | Why |
|---|---|
| `key` declares what identifies an invocation | Order carries no meaning, and best-match pairing would pair a wrong key with the expectation it least resembles — reporting a fumbled argument where the model addressed the wrong thing. A wrong key reads as one miss plus one unexpected invocation. |
| Only named kwargs are scored | Silence is *not measured*, not *passing*. Adopt one kwarg at a time instead of specifying every fixture first. |
| Comparison never stringifies | The string `"false"` is not the boolean `false`, and `""` is not `nil`. Rendering both sides to text would erase the distinction the whole thing exists for. Numbers are the one normalisation: JSON decodes to float64, YAML to int. |
| `because` is required | Loading refuses a binding without it. |

**Bindings are authored from the record; counts are established from a run.**
Two opposite disciplines in one block, and which applies turns on whether the
correct answer exists independently of the model. How many invocations a
complete answer contains is nobody's computation — observe it. What `default`
belongs on a not-null boolean column is stated in the record and settled by
Rails — asking a run would be letting the thing under test grade itself. No
sample has ever bound `default` correctly, so there was never an observation to
establish these from; that is the finding, not a reason to wait.

**Nothing mechanical can check that an expectation is right.** The diff test
holds the described package to one difference and cannot say the descriptions
are good. Type comparison says `"false"` is not `false` and cannot say `nil` was
correct — only the record says that. The guards here catch drift and accident,
not bad judgment. What protects a judgment is that it is written where someone
can disagree with it, which is why `because` lives in the case file rather than
only in `prov-2026-2b121b62`.

## Known rough edges

- **Only the `sedum` arm runs.** `baseline` is declared and rejected at run time.
- **Behavior is not measured.** The schema reserves it; nothing implements it.
