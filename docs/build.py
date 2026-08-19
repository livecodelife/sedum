#!/usr/bin/env python3
"""Assemble the published dashboard from the stored eval results.

The site is a reader. This script copies `evals/results/*.jsonl` verbatim,
writes a manifest naming what it found, and pulls the descriptive half of each
case file across so the page can say what a case asked for. It never edits,
normalises or backfills an entry: results are append-only observations
(prov-2026-eb283c56), and a build step that tidied one would be editing a
measurement.

Nothing it writes is committed. The output directory is derived from the
results on every build, so the site cannot drift from the file it reads and a
new case needs no edit here (prov-2026-b0bf8584).

    python3 docs/build.py                 # build into _site/
    python3 docs/build.py -o /tmp/site    # build somewhere else
    python3 docs/build.py --check         # build to a temp dir and verify it
"""

from __future__ import annotations

import argparse
import filecmp
import hashlib
import json
import pathlib
import shutil
import subprocess
import sys
import tempfile

ROOT = pathlib.Path(__file__).resolve().parent.parent
DOCS = ROOT / "docs"
RESULTS = ROOT / "evals" / "results"
CASES = ROOT / "evals" / "cases"

# The static half of the site. Everything else under docs/ is source for the
# build or proof for the record, and is not published.
STATIC = ["index.html", "style.css", "app.js", "stats.js"]

# The half of a case file that describes what was asked. The expectations are
# carried so the page can show what a complete answer looks like beside what
# was selected - as a description, never as a score the site computed.
CASE_FIELDS = [
    "id",
    "application",
    "language",
    "framework",
    "tightness",
    "arm",
    "records",
    "generators",
    "models",
    "check",
    "expect",
]


def die(msg: str) -> None:
    print(f"build.py: {msg}", file=sys.stderr)
    sys.exit(1)


def read_entries(path: pathlib.Path) -> tuple[list[dict], list[str]]:
    """Parse one results file, keeping unreadable lines as unreadable.

    A line that will not parse is reported rather than dropped or repaired.
    The file is the record of what happened and this is downstream of it.
    """
    entries: list[dict] = []
    broken: list[str] = []
    for n, line in enumerate(path.read_text(encoding="utf-8").splitlines(), start=1):
        if not line.strip():
            continue
        try:
            entries.append(json.loads(line))
        except json.JSONDecodeError as e:
            broken.append(f"{path.name}:{n}: {e}")
    return entries, broken


def load_cases() -> dict:
    """The descriptive half of every case file, keyed by case id."""
    if not CASES.is_dir():
        return {}
    try:
        import yaml
    except ImportError:
        die("PyYAML is needed to read evals/cases (pip install pyyaml)")

    out = {}
    for path in sorted(CASES.glob("*.yaml")):
        doc = yaml.safe_load(path.read_text(encoding="utf-8")) or {}
        case_id = doc.get("id") or path.stem
        out[case_id] = {k: doc[k] for k in CASE_FIELDS if k in doc}
    return out


def source_commit() -> dict:
    """What the site was built from, so a reader can pin the page itself."""
    def git(*args: str) -> str:
        try:
            return subprocess.check_output(
                ["git", *args], cwd=ROOT, stderr=subprocess.DEVNULL
            ).decode().strip()
        except Exception:
            return ""

    return {"commit": git("rev-parse", "--short", "HEAD"), "ref": git("rev-parse", "--abbrev-ref", "HEAD")}


def build(out: pathlib.Path) -> dict:
    if not RESULTS.is_dir():
        die(f"no results directory at {RESULTS}")

    files = sorted(RESULTS.glob("*.jsonl"))
    if not files:
        die(f"no *.jsonl under {RESULTS}")

    if out.exists():
        shutil.rmtree(out)
    data = out / "data"
    data.mkdir(parents=True)

    cases = load_cases()
    manifest = {"source": source_commit(), "cases": [], "unreadable": []}

    for path in files:
        entries, broken = read_entries(path)
        manifest["unreadable"].extend(broken)

        # Copied byte for byte. The page parses the same lines the harness
        # wrote, so there is no transformation here to be wrong.
        shutil.copy2(path, data / path.name)

        case_id = path.stem
        manifest["cases"].append(
            {
                "id": case_id,
                "file": f"data/{path.name}",
                "entries": len(entries),
                "sha256": hashlib.sha256(path.read_bytes()).hexdigest()[:16],
                "case": cases.get(case_id),
            }
        )

    (data / "manifest.json").write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")

    for name in STATIC:
        src = DOCS / name
        if not src.is_file():
            die(f"missing {src}")
        shutil.copy2(src, out / name)

    # upload-pages-artifact does not run Jekyll, but a branch deployment would.
    (out / ".nojekyll").write_text("", encoding="utf-8")

    return manifest


def check() -> None:
    """Build to a temp dir and assert the site is a faithful copy."""
    with tempfile.TemporaryDirectory() as tmp:
        out = pathlib.Path(tmp) / "site"
        before = {p.name: p.read_bytes() for p in sorted(RESULTS.glob("*.jsonl"))}
        manifest = build(out)

        problems: list[str] = []

        # 1. Every results file reached the site, byte for byte.
        for name in before:
            copied = out / "data" / name
            if not copied.is_file():
                problems.append(f"{name} was not published")
            elif not filecmp.cmp(RESULTS / name, copied, shallow=False):
                problems.append(f"{name} was altered on the way to the site")

        # 2. The manifest names all of them and invents none.
        named = {c["file"].removeprefix("data/") for c in manifest["cases"]}
        for missing in sorted(set(before) - named):
            problems.append(f"{missing} is published but not in the manifest")
        for extra in sorted(named - set(before)):
            problems.append(f"{extra} is in the manifest but has no source file")

        # 3. The source results are untouched. The site is a reader.
        after = {p.name: p.read_bytes() for p in sorted(RESULTS.glob("*.jsonl"))}
        if before != after:
            problems.append("evals/results changed during a build; the site must never write to it")

        # 4. Every static file the page needs is present.
        for name in [*STATIC, ".nojekyll"]:
            if not (out / name).is_file():
                problems.append(f"{name} is missing from the built site")

        # 5. An unreadable line is reported rather than silently dropped.
        if manifest["unreadable"]:
            for line in manifest["unreadable"]:
                print(f"  unreadable: {line}")

        if problems:
            for p in problems:
                print(f"build.py: {p}", file=sys.stderr)
            sys.exit(1)

        total = sum(c["entries"] for c in manifest["cases"])
        print(f"ok: {len(manifest['cases'])} case file(s), {total} entries, copied unaltered")


def main() -> None:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("-o", "--out", default=str(ROOT / "_site"), help="output directory (default: _site)")
    ap.add_argument("--check", action="store_true", help="build to a temp dir and verify, then discard")
    args = ap.parse_args()

    if args.check:
        check()
        return

    manifest = build(pathlib.Path(args.out))
    total = sum(c["entries"] for c in manifest["cases"])
    print(f"built {args.out}: {len(manifest['cases'])} case file(s), {total} entries")
    for line in manifest["unreadable"]:
        print(f"  unreadable: {line}", file=sys.stderr)


if __name__ == "__main__":
    main()
