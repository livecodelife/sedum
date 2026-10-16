#!/usr/bin/env python3
"""Print the lifecycle demo's pass/fail table from run.sh's results file.

    table.py <results> <sedum label> <provenance ref> <sedum v0.8.0 ref> <log>

Exits 1 if any step failed or did not run."""
import sys
results, sedum, prov, main, log = sys.argv[1:]
TITLES = {
    "step0": "Scaffold: chi package, records dir, stub model",
    "step1": "v1 records (5)",
    "step2": "First grow: 3 model calls, markers carry gen + pkg, code builds",
    "step6": "Second grow: 0 model calls, no file touched",
    "step7v1": "Edit one record: 1 call, only its file changes",
    "step7b": "Template change: re-render with 0 calls, gen unchanged",
    "step9": "grow --check: flags the edited record, writes nothing",
    "step11": "Ownership: a second record claiming a region is refused",
    "step12": "grow --force: re-selects everything, same bytes",
    "step10": "Deprecate / delete a record: orphans reported, --prune-orphans removes them",
    "step3": "v1 supersession + bug",
    "step3grow": "Growing a v1 supersession before migrating is refused",
    "step4": "provenance migrate: 7 -> 5 records, aliases, versions, manifest",
    "step5": "Grow after migration: 1 call, alias markers rewritten to the survivor",
    "step8": "Stale children + repin",
    "step13": "Evidence: complete, parent bump, repin, spec edit; check --commit",
    "step14": "Template edit + replay vs restamp: hashes and evidence",
    "step2prime": "Markers without gen (sedum v0.8.0): one full re-select",
    "step4b": "grow --stamp: same markers become current with 0 calls",
}
rows, seq = [], "sedum"
seen0 = 0
for line in open(results):
    step, status, *rest = line.rstrip("\n").split(" ", 2)
    if step == "step0":
        seen0 += 1
        seq = "sedum" if seen0 == 1 else "migrate"
    if step == "step2prime":
        seq = "pre-gen"
    rows.append((seq, step, status, rest[0] if rest else ""))
expected = ["step0", "step1", "step2", "step6", "step7v1", "step7b", "step9", "step11", "step12", "step10",
            "step0", "step1", "step2", "step3", "step3grow", "step4", "step5", "step8", "step13", "step14", "step2prime", "step4b"]
print(f"\nsedum {sedum} · provenance {prov[:12]} · sedum v0.8.0 {main[:12]} · stub model\n")
w = max(len(t) for t in TITLES.values())
print(f"  {'sequence':8}  {'step':10}  {'what it shows':{w}}  result")
print(f"  {'-'*8}  {'-'*10}  {'-'*w}  ------")
fails = 0
for seqn, step, status, detail in rows:
    mark = {"pass": "PASS", "FAIL": "FAIL", "skipped": "skip"}[status]
    fails += status == "FAIL"
    print(f"  {seqn:8}  {step:10}  {TITLES.get(step, step):{w}}  {mark}" + (f"  ({detail})" if detail else ""))
missing = len(expected) - len(rows)
passed = sum(r[2] == "pass" for r in rows)
skipped = sum(r[2] == "skipped" for r in rows)
print(f"\n{passed} passed, {fails} failed, {skipped} skipped" + (f", {missing} not run (after a failure)" if missing > 0 else ""))
print(f"Full log: {log}")
sys.exit(1 if fails or missing > 0 else 0)
