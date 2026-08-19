// Proof for prov-2026-b0bf8584.
//
// The blueprint's constraints are rules about how a number may be read, and a
// rule enforced only by review is a rule waiting to be broken by the next
// column somebody adds. Each one below is asserted.
//
//   node --test docs/stats_test.mjs

import test from "node:test";
import assert from "node:assert/strict";
import { readFileSync, readdirSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";

import {
  wilson,
  fmtInterval,
  overlaps,
  cites,
  comparable,
  caveats,
  summarise,
  compare,
} from "./stats.js";

const here = dirname(fileURLToPath(import.meta.url));
const resultsDir = join(here, "..", "evals", "results");

function storedEntries() {
  const out = [];
  for (const name of readdirSync(resultsDir).filter((f) => f.endsWith(".jsonl"))) {
    const text = readFileSync(join(resultsDir, name), "utf8");
    for (const line of text.split("\n")) {
      if (line.trim()) out.push(JSON.parse(line));
    }
  }
  return out;
}

test("wilson matches the harness's own interval", () => {
  // The four widths prov-2026-3039750e tabulated to argue that sample size is
  // a property of the question. If these drift, the page is telling a reader
  // something different from what the terminal report tells them.
  const cases = [
    [5, 5, 0.57, 1.0],
    [10, 10, 0.72, 1.0],
    [20, 20, 0.84, 1.0],
    [30, 30, 0.89, 1.0],
    [4, 5, 0.38, 0.96],
  ];
  for (const [k, n, low, high] of cases) {
    const iv = wilson(k, n);
    assert.equal(Number(iv.low.toFixed(2)), low, `${k}/${n} low`);
    assert.equal(Number(iv.high.toFixed(2)), high, `${k}/${n} high`);
  }
});

test("a perfect score does not collapse to a point", () => {
  // 5/5 is the reading most likely to be over-believed, and an interval that
  // said [1.00, 1.00] would claim it was evidence of a rate of 1.00.
  const iv = wilson(5, 5);
  assert.ok(iv.low < 0.6, "5/5 must stay consistent with a true rate near 57%");
  assert.equal(iv.high, 1);
});

test("an empty denominator is not a rate of zero", () => {
  assert.equal(wilson(0, 0), null);
  assert.equal(fmtInterval(null), "—");
});

test("the fraction is never hidden behind the interval", () => {
  // prov-2026-0baaa119: the interval goes beside 4/5, never instead of it.
  const rendered = fmtInterval(wilson(4, 5));
  assert.ok(rendered.startsWith("4/5"), `got ${rendered}`);
  assert.match(rendered, /\[0\.38, 0\.96\]/);
});

test("4/5 and 5/5 are not distinguished at five samples", () => {
  // The fact the old table hid, and the reason Compare exists at all.
  assert.ok(overlaps(wilson(4, 5), wilson(5, 5)));
  const c = compare(wilson(4, 5), wilson(5, 5));
  assert.equal(c.distinguished, false);
  assert.match(c.text, /not distinguished/);
});

test("comparison never claims a difference, only a failure to distinguish", () => {
  const c = compare(wilson(0, 30), wilson(30, 30));
  assert.equal(c.distinguished, true);
  // No verdict vocabulary: a p-value is a verdict wearing a number.
  for (const forbidden of [/significant/i, /p\s*[<=]/, /pass/i, /fail/i]) {
    assert.doesNotMatch(c.text, forbidden, `${forbidden} must not appear in a comparison`);
  }
});

test("smoke, dirty and unstated entries are excluded from comparison", () => {
  // The three the history chain skips (prov-2026-3039750e, prov-2026-c5ad54ff).
  assert.equal(comparable({ clean: true, resolution: "smoke" }), false);
  assert.equal(comparable({ clean: true, resolution: "" }), false);
  assert.equal(comparable({ clean: false, resolution: "coarse" }), false);
  assert.equal(comparable({ clean: true, resolution: "coarse" }), true);
  assert.equal(comparable({ clean: true, resolution: "fine" }), true);

  assert.equal(cites("smoke"), false);
  assert.equal(cites(""), false);
});

test("every excluded entry carries the reason it was excluded", () => {
  const smoke = caveats({ clean: true, resolution: "smoke", samples: 2, fixture: "abc" });
  assert.equal(smoke.length, 1);
  assert.match(smoke[0].why, /not a measurement/);

  const dirty = caveats({ clean: false, resolution: "coarse", dirty: ["evals/run.go"], fixture: "abc" });
  assert.equal(dirty.length, 1);
  assert.match(dirty[0].why, /evals\/run\.go/);

  const unstated = caveats({ clean: true, resolution: "", samples: 5 });
  const tags = unstated.map((c) => c.tag);
  assert.deepEqual(tags, ["unstated", "no digest"]);
});

test("failed samples stay outside the denominator and are still counted", () => {
  // An unreachable endpoint is not a model that chose badly.
  const s = summarise({
    case: "x",
    at: "2026-08-17T00:00:00Z",
    valid: 0,
    invalid: 0,
    failed: 2,
    samples: 2,
    runs: [{ outcome: "failed", ms: 1 }, { outcome: "failed", ms: 1 }],
  });
  assert.equal(s.answered, 0);
  assert.equal(s.validity, null, "no answered samples is not a rate of zero");
  assert.equal(s.failed, 2, "the excluded count is still visible");
});

test("the baseline arm reports nothing it has no vocabulary for", () => {
  // prov-2026-a4dbe65c: blank with a reason, never zero.
  const s = summarise({
    case: "todo-rails-baseline",
    at: "2026-08-18T00:00:00Z",
    arm: "baseline",
    valid: 1,
    invalid: 0,
    failed: 0,
    samples: 1,
    runs: [{ outcome: "valid", ms: 1, files: { "a.rb": "x" }, missing: ["b.rb"] }],
  });
  assert.equal(s.selection, null);
  assert.equal(s.fill, null);
  assert.equal(s.idempotent, null);
  assert.equal(s.syntax, null);
  assert.equal(s.paths.wrote, 1);
  assert.equal(s.paths.want, 2);
  assert.deepEqual(s.paths.missing, [{ name: "b.rb", n: 1 }]);
});

test("a signal no sample recorded reads as not recorded, not as zero", () => {
  const s = summarise({
    case: "x",
    at: "2026-08-15T00:00:00Z",
    valid: 1,
    invalid: 0,
    failed: 0,
    samples: 1,
    runs: [{ outcome: "valid", ms: 1, counts: { addColumn: 2 } }],
  });
  assert.equal(s.fill, null, "an entry written before anchor fill existed filled no anchors");
  assert.equal(s.behavior, null, "behaviour off is not behaviour broken");
});

test("behaviour keeps its three outcomes apart", () => {
  // A service that never booted and one that booted and answered wrongly are
  // different findings; one rate would merge a broken package with a wrong one.
  const s = summarise({
    case: "x",
    at: "2026-08-17T00:00:00Z",
    valid: 3,
    invalid: 0,
    failed: 0,
    samples: 3,
    runs: [
      { outcome: "valid", ms: 1, behavior: { outcome: "ok", checks: 20, passed: 20 } },
      { outcome: "valid", ms: 1, behavior: { outcome: "checks_failed", checks: 20, passed: 19, failed: ["x"] } },
      { outcome: "valid", ms: 1, behavior: { outcome: "failed", failed_phase: "build", checks: 0, passed: 0 } },
    ],
  });
  assert.equal(s.behavior.working, 1);
  assert.equal(s.behavior.disagreed, 1);
  assert.equal(s.behavior.broke, 1);
  assert.equal(s.behavior.measured, 3);
  assert.deepEqual(s.behavior.failures, [{ name: "x", n: 1 }]);
  assert.deepEqual(s.behavior.phases, [{ name: "build", n: 1 }]);
});

test("selection is described, never scored", () => {
  // A mean of 1.0 could be every sample selecting once or half selecting
  // twice, so both numbers are kept.
  const s = summarise({
    case: "x",
    at: "2026-08-17T00:00:00Z",
    valid: 2,
    invalid: 0,
    failed: 0,
    samples: 2,
    runs: [
      { outcome: "valid", ms: 1, counts: { addColumn: 2 }, first: "addColumn" },
      { outcome: "valid", ms: 1, counts: { addColumn: 0 }, first: "createEndpoint" },
    ],
  });
  const row = s.selection.actions.find((a) => a.action === "addColumn");
  assert.equal(row.mean, 1);
  assert.equal(row.samplesSelecting, 1, "the mean alone would hide this");
  // Nothing in the summary claims an action was correct or expected.
  assert.equal("selected" in row, false);
  assert.equal("exact" in row, false);
});

test("every stored entry summarises without throwing", () => {
  const entries = storedEntries();
  assert.ok(entries.length > 0, "no stored entries found to read");

  for (const e of entries) {
    const s = summarise(e);
    assert.equal(s.case, e.case);

    // The denominator invariant, on real data.
    assert.equal(s.answered, (e.valid || 0) + (e.invalid || 0));
    if (s.validity) {
      assert.equal(s.validity.samples, s.answered);
      assert.ok(s.validity.low >= 0 && s.validity.high <= 1);
    }

    // A baseline row never carries a number it has no vocabulary for.
    if (s.baseline) {
      assert.equal(s.selection, null);
      assert.equal(s.fill, null);
    }

    // Anything not comparable states at least one reason why.
    if (!s.comparable) assert.ok(s.caveats.length > 0, `${s.key} excluded with no reason given`);
  }
});

test("the stored corpus still contains the entries the page is built to explain", () => {
  // Not a check on the data — a check that the reader has not silently stopped
  // seeing a shape it claims to handle.
  const entries = storedEntries();
  const arms = new Set(entries.map((e) => e.arm || "sedum"));
  assert.ok(arms.has("sedum") && arms.has("baseline"), `arms seen: ${[...arms]}`);

  const outcomes = new Set();
  for (const e of entries) for (const r of e.runs || []) if (r.behavior) outcomes.add(r.behavior.outcome);
  assert.ok(outcomes.has("ok"), "no working behaviour sample to render");
  assert.ok(outcomes.has("checks_failed"), "no disagreeing behaviour sample to render");
});
