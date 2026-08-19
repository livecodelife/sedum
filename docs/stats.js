// Pure readings of a stored eval entry: no DOM, no fetch, no state.
//
// Imported by app.js in the browser and by stats_test.mjs under `node --test`,
// so every rule about how a number may be read is asserted rather than
// reviewed (prov-2026-b0bf8584).
//
// Nothing here scores an answer. The harness scores selections and bindings at
// run time against a case's expectations; a stored entry does not carry that
// score, because prov-2026-2256e6fa stores invocations "rather than scored" on
// purpose - the change that first makes data visible should not also decide
// what to conclude from it. Re-deriving a score here would be a second
// implementation of the harness's scoring, and prov-2026-9dcf2658 already
// argues in a different place why two implementations of one judgement
// eventually disagree. So these functions summarise what an entry recorded and
// never grade it.

export const Z = 1.96;

// wilson is the Wilson score interval at 95%, matching evals/report.go exactly.
//
// Chosen over the normal approximation because the approximation collapses to
// zero width at 0/n and n/n, and n/n is the reading most likely to be
// over-believed: 5/5 is not evidence of a rate of 1.00 (prov-2026-0baaa119).
//
// Returns null rather than a zero interval when there is nothing to divide by,
// so a caller renders "not measured" instead of a rate of zero.
export function wilson(successes, samples) {
  if (!samples || samples <= 0) return null;

  const n = samples;
  const p = successes / n;
  const denom = 1 + (Z * Z) / n;
  const center = (p + (Z * Z) / (2 * n)) / denom;
  const spread = (Z * Math.sqrt((p * (1 - p)) / n + (Z * Z) / (4 * n * n))) / denom;

  return {
    successes,
    samples,
    rate: p,
    low: Math.max(0, center - spread),
    high: Math.min(1, center + spread),
  };
}

// fmtInterval renders the fraction and the interval together.
//
// The fraction never disappears behind the interval: the sample count is the
// fact that produced the width, and a reader who sees only the bounds has lost
// it (prov-2026-0baaa119).
export function fmtInterval(iv) {
  if (!iv) return "—";
  return `${iv.successes}/${iv.samples} [${iv.low.toFixed(2)}, ${iv.high.toFixed(2)}]`;
}

// overlaps reports whether two intervals leave open that the rates are the
// same. It is not a test and decides nothing; it is the condition under which
// a reported difference is not distinguishable from sampling.
export function overlaps(a, b) {
  if (!a || !b) return false;
  return a.low <= b.high && b.low <= a.high;
}

// cites reports whether a rate drawn at this resolution may be cited as a
// measurement. Smoke proves plumbing, and an entry drawn before resolutions
// existed states no question at all (prov-2026-3039750e).
export function cites(resolution) {
  return resolution === "coarse" || resolution === "fine";
}

// comparable is the test the history chain applies: an entry joins a
// comparison only when it is clean and its resolution cites. A dirty tree pins
// nothing, and an unstated sample size is a question nobody chose
// (prov-2026-c5ad54ff).
export function comparable(entry) {
  return entry.clean === true && cites(entry.resolution);
}

// caveats names every reason an entry should not be read as a measurement, so
// the reason travels with the number wherever it is shown.
export function caveats(entry) {
  const out = [];

  if (entry.resolution === "smoke") {
    out.push({
      tag: "smoke",
      short: "plumbing only",
      why:
        `Drawn at n=${entry.samples} to prove the harness runs at all — a new case loads, ` +
        `an endpoint answers. A rate from a smoke run is not a measurement and is not cited as one.`,
    });
  }

  if (!cites(entry.resolution) && entry.resolution !== "smoke") {
    out.push({
      tag: "unstated",
      short: "sample size nobody chose",
      why:
        "This entry names no resolution. It was drawn before the harness asked what question a " +
        "sample size was for, so its n was a default rather than a decision. Stamping it " +
        "\"coarse\" now would invent a choice nobody made.",
    });
  }

  if (entry.clean !== true) {
    out.push({
      tag: "dirty",
      short: "not re-runnable",
      why:
        "The working tree carried uncommitted changes, so the commit does not pin the generator " +
        "package, the record, or Sedum's own code. The numbers are what was measured; they " +
        "cannot be reproduced from this commit." +
        (entry.dirty && entry.dirty.length
          ? ` Changed: ${entry.dirty.join(", ")}.`
          : ""),
    });
  }

  if (!entry.fixture) {
    out.push({
      tag: "no digest",
      short: "inputs unknown",
      why:
        "Written before entries recorded a digest of the packages, records and behaviour target " +
        "they were drawn against. There is no way to recover what this run used beyond its " +
        "commit, so it is compared against nothing.",
    });
  }

  return out;
}

// The arms, as a ladder of what the model was given (prov-2026-672c6471).
export const ARMS = {
  sedum: {
    given: "the provenance record and the generator package's action catalog",
    withoutPackage: false,
    heldToPaths: true,
  },
  baseline: {
    given: "the record — its intent, its constraints, and the exact paths it authorizes — and no catalog",
    withoutPackage: true,
    heldToPaths: true,
  },
  intent: {
    given: "the record's intent, and nothing else: no constraints and no list of files",
    withoutPackage: true,
    heldToPaths: false,
  },
};

// withoutPackage mirrors Case.WithoutPackage() in evals/case.go, which exists as
// a predicate rather than a comparison at each site precisely so that a third
// package-free arm does not have to be remembered in four places. This file
// remembered it in four places until the intent arm landed; now it does not.
//
// An unknown arm is treated as packaged, because reporting nothing is the
// recoverable mistake and reporting a number the arm cannot have is not.
export function withoutPackage(arm) {
  return ARMS[arm]?.withoutPackage === true;
}

// heldToPaths reports whether the arm was given the paths it is scored against.
// The intent arm is not: its answer is parsed with no allowed list, so every
// path it writes is accepted and none can be missing or unexpected. The rate
// would be wrote/wrote on every run forever (prov-2026-18a6e7a5).
export function heldToPaths(arm) {
  return ARMS[arm]?.heldToPaths !== false;
}

// UNDEFINED_HERE is why a cell is blank rather than zero, per arm.
// Blank with a reason, never zero: a zero would read as a score, and a 1.00
// that cannot be anything else ends a question instead of inviting one
// (prov-2026-a4dbe65c, prov-2026-18a6e7a5).
export const UNDEFINED_HERE = {
  signals: (arm) =>
    `arm=${arm} has no generator package and no action vocabulary, so there is no selection to ` +
    "count, no binding to check, no anchor to fill and no injection to repeat. All four are " +
    "undefined here rather than zero.",
  paths:
    "arm=intent is given no list of files. Its answer is accepted whatever paths it chooses, so " +
    "nothing can be missing and nothing can be unexpected — the rate would be 1.00 on every run " +
    "by construction. That is arithmetic with no question behind it, not a perfect score. The " +
    "terminal report does print it; this page declines to.",
};

function byCount(counts) {
  return Object.entries(counts)
    .map(([name, n]) => ({ name, n }))
    .sort((a, b) => b.n - a.n || a.name.localeCompare(b.name));
}

// summarise reduces one stored entry to what a page can show, and to nothing
// the entry did not record.
//
// Every denominator here is stated as a field beside its numerator, so a
// reader is never asked to guess what a rate was over.
export function summarise(entry) {
  const runs = entry.runs || [];
  const arm = entry.arm || "sedum";
  const packageless = withoutPackage(arm);

  // Failed samples never reached the model, so they stay outside every
  // denominator - an unreachable endpoint is not a model that chose badly
  // (prov-2026-e6969eb3). The count is carried so a reader can see what the
  // denominator excluded.
  const answered = (entry.valid || 0) + (entry.invalid || 0);

  const s = {
    key: `${entry.case}@${entry.at}`,
    case: entry.case,
    at: entry.at,
    commit: entry.commit,
    clean: entry.clean === true,
    dirty: entry.dirty || [],
    fixture: entry.fixture || "",
    resolution: entry.resolution || "",
    arm,
    packageless,
    heldToPaths: heldToPaths(arm),
    model: entry.model || {},
    endpoint: entry.endpoint || "",
    language: entry.language || "",
    framework: entry.framework || "",
    tightness: entry.tightness || "",
    samples: entry.samples || 0,
    concurrency: entry.concurrency || 1,
    retries: entry.retries || 0,
    valid: entry.valid || 0,
    invalid: entry.invalid || 0,
    failed: entry.failed || 0,
    answered,
    validity: wilson(entry.valid || 0, answered),
    wallMs: entry.wall_ms || 0,
    details: entry.details || [],
    comparable: comparable(entry),
    caveats: caveats(entry),
    runs,
  };

  // First-call validity survives a raised retry budget, and only means
  // something when there was a budget to spend (prov-2026-0811425c).
  if (s.retries > 0) {
    const firstTry = runs.filter(
      (r) => r.outcome === "valid" && !(r.rejected > 0),
    ).length;
    s.firstCall = wilson(firstTry, answered);
  }

  const calls = runs.reduce((a, r) => a + (r.calls || 0), 0);
  s.calls = calls > 0 ? { total: calls, over: runs.length, mean: calls / runs.length } : null;

  const prompt = runs.reduce((a, r) => a + (r.prompt_tokens || 0), 0);
  const completion = runs.reduce((a, r) => a + (r.completion_tokens || 0), 0);
  s.tokens =
    prompt || completion
      ? {
          prompt,
          completion,
          perCall: calls > 0 ? { prompt: prompt / calls, completion: completion / calls } : null,
        }
      : null;

  const durations = runs.map((r) => r.ms || 0).filter((ms) => ms > 0);
  s.timing = durations.length
    ? {
        fastest: Math.min(...durations),
        slowest: Math.max(...durations),
        mean: durations.reduce((a, b) => a + b, 0) / durations.length,
      }
    : null;

  // A slug can be tallied and a sentence cannot, so both are kept: the tally
  // says how often, the sentence says what the answer got wrong. Violations are
  // stored in the same order as the slugs they render (prov-2026-986ac4ca).
  const ruleCounts = {};
  const violationsByRule = {};
  for (const r of runs) {
    (r.rules || []).forEach((rule, i) => {
      ruleCounts[rule] = (ruleCounts[rule] || 0) + 1;
      const text = (r.violations || [])[i];
      if (text) (violationsByRule[rule] ||= new Set()).add(text);
    });
  }
  s.rules = byCount(ruleCounts).map((d) => ({
    ...d,
    violations: [...(violationsByRule[d.name] || [])],
  }));

  s.behavior = behaviorOf(runs);

  // Selection, binding, anchor fill and idempotency are undefined on any arm
  // with no package, and are not reported there (prov-2026-a4dbe65c).
  if (!packageless) {
    s.fill = signalOf(runs, "fill", "planted", "filled", "missed");
    s.syntax = signalOf(runs, "syntax", "checked", "parsed");
    s.idempotent = signalOf(runs, "idempotent", "files", "stable", "unstable");
    s.selection = selectionOf(runs);
    s.paths = null;
  } else {
    s.fill = s.syntax = s.idempotent = s.selection = null;
    // And the path rate only for an arm that was given paths to be held to.
    s.paths = s.heldToPaths ? pathsOf(runs) : null;
  }

  return s;
}

// behaviorOf keeps the three outcomes apart rather than reducing them to a
// rate. A service that never booted and one that booted and answered wrongly
// are different findings, and a single number would merge a broken generator
// package with a wrong one (evals/behavior.go).
function behaviorOf(runs) {
  const applied = runs.filter((r) => r.behavior);
  if (!applied.length) return null;

  let working = 0, disagreed = 0, broke = 0, checks = 0, passed = 0, elapsed = 0;
  const failures = {}, phases = {};

  for (const r of applied) {
    const b = r.behavior;
    elapsed += b.elapsed || 0;
    checks += b.checks || 0;
    passed += b.passed || 0;
    for (const f of b.failed || []) failures[f] = (failures[f] || 0) + 1;

    if (b.outcome === "ok") working++;
    else if (b.outcome === "checks_failed") disagreed++;
    else {
      broke++;
      if (b.failed_phase) phases[b.failed_phase] = (phases[b.failed_phase] || 0) + 1;
    }
  }

  return {
    measured: applied.length,
    working,
    disagreed,
    broke,
    workingRate: wilson(working, applied.length),
    checks,
    passed,
    elapsedNs: elapsed,
    failures: byCount(failures),
    phases: byCount(phases),
  };
}

// signalOf sums one derived signal across the samples that recorded it.
//
// A sample written before the signal existed carries nothing rather than a
// zero, so `measured` counts the samples that actually contributed and a run
// with none reports "not recorded" instead of a rate of zero
// (prov-2026-d61010a4, prov-2026-eb283c56).
//
// The denominator is the signal's own unit — planted anchors, files checked,
// files re-applied — not samples, so every caller has to name it. It matches
// evals/report.go, which puts the interval on the same unit.
function signalOf(runs, field, denomKey, numerKey, listKey) {
  const present = runs.filter((r) => r[field]);
  if (!present.length) return null;

  let denom = 0, numer = 0;
  const listed = {};
  for (const r of present) {
    denom += r[field][denomKey] || 0;
    numer += r[field][numerKey] || 0;
    for (const name of (listKey && r[field][listKey]) || []) listed[name] = (listed[name] || 0) + 1;
  }

  return {
    measured: present.length,
    denom,
    numer,
    unit: denomKey,
    interval: wilson(numer, denom),
    listed: byCount(listed),
  };
}

// selectionOf describes what was selected. It does not score it.
//
// The mean is over every sample that recorded counts, and `samplesSelecting` is
// how many of them reached for the action at all — because a mean of 1.0 could
// be every sample selecting once or half of them selecting twice, and those are
// different findings (the observed means of 0.33, 0.67 and 3.33 that
// prov-2026-eb283c56 kept per-sample rows for are exactly this case).
function selectionOf(runs) {
  const scored = runs.filter((r) => r.counts);
  if (!scored.length) return null;

  const total = {}, selecting = {};
  for (const r of scored) {
    for (const [action, n] of Object.entries(r.counts)) {
      total[action] = (total[action] || 0) + n;
      if (n > 0) selecting[action] = (selecting[action] || 0) + 1;
    }
  }

  const actions = Object.keys(total)
    .map((action) => ({
      action,
      mean: total[action] / scored.length,
      total: total[action],
      samplesSelecting: selecting[action] || 0,
    }))
    .sort((a, b) => b.mean - a.mean || a.action.localeCompare(b.action));

  const firstCounts = {};
  for (const r of scored) if (r.first) firstCounts[r.first] = (firstCounts[r.first] || 0) + 1;

  return { scored: scored.length, actions, first: byCount(firstCounts) };
}

// pathsOf is the baseline arm's only rate: how much of what the record
// authorized the model actually wrote.
//
// A path nothing authorized never reached disk, and saying so matters — code at
// the wrong path is a different finding from no code at all
// (prov-2026-a4dbe65c).
function pathsOf(runs) {
  let wrote = 0, want = 0;
  const missing = {}, unexpected = {};

  for (const r of runs) {
    if (r.outcome !== "valid") continue;
    const files = Object.keys(r.files || {});
    wrote += files.length;
    want += files.length + (r.missing || []).length;
    for (const p of r.missing || []) missing[p] = (missing[p] || 0) + 1;
    for (const p of r.unexpected || []) unexpected[p] = (unexpected[p] || 0) + 1;
  }

  if (want === 0) return null;
  return {
    wrote,
    want,
    interval: wilson(wrote, want),
    missing: byCount(missing),
    unexpected: byCount(unexpected),
  };
}

// compare states the one thing a pair of rates supports and refuses the rest.
//
// Overlap is a conservative signal, not a test: it says these runs do not
// distinguish the arms, never that there is no difference. No p-value and no
// verdict — a p-value is a verdict wearing a number (prov-2026-0baaa119).
export function compare(a, b) {
  if (!a || !b) return { verdict: null, text: "one side has no scoreable samples; nothing to compare" };
  if (overlaps(a, b)) {
    return {
      distinguished: false,
      text: "not distinguished — the intervals overlap, so these runs do not tell these rates apart",
    };
  }
  return {
    distinguished: true,
    delta: b.rate - a.rate,
    text: "the intervals do not overlap at this sample size",
  };
}
