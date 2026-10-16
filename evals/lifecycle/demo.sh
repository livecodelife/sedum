#!/usr/bin/env bash
# One-command demo of hash-aware grow and provenance migrate, run from this sedum
# checkout. Builds sedum from the checkout, sedum v0.8.0 from this repo's history,
# and the provenance CLI from github.com/livecodelife/provenance at a pinned
# commit; runs every lifecycle step against a local stub model (no real model is
# ever called); prints a pass/fail table. Exit status is non-zero if any step fails.
#
#   evals/lifecycle/demo.sh
#   PROVENANCE_REF=<commit> evals/lifecycle/demo.sh   # another provenance CLI commit
#
# Needs: go, python3, git, curl. linespec is optional (only the v1 lint uses it).
set -euo pipefail

HERE="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(git -C "$HERE" rev-parse --show-toplevel)"

# ---- pins ------------------------------------------------------------------------
SEDUM_MAIN_REF="${SEDUM_MAIN_REF:-0571de7de06364e9cbc0076098f96c3a4e9ee7cc}"  # sedum v0.8.0, for step 2'
PROVENANCE_REF="${PROVENANCE_REF:-afdc4f2d9cf440a0eb79d038de0be88d2f4f3d83}"  # provenance, branch cli

CACHE="$HERE/.demo"
for tool in go python3 git curl; do command -v "$tool" >/dev/null || { echo "demo: $tool is required" >&2; exit 2; }; done
mkdir -p "$CACHE/bin"

# ---- build ------------------------------------------------------------------------
printf 'Building: sedum from this checkout, sedum v0.8.0, provenance CLI %s...\n' "${PROVENANCE_REF:0:12}"
SEDUM="$CACHE/bin/sedum"
(cd "$ROOT" && go build -o "$SEDUM" ./cmd/sedum)

SEDUM_MAIN="$CACHE/bin/sedum-main-${SEDUM_MAIN_REF:0:12}"
if [ ! -x "$SEDUM_MAIN" ]; then
	src=$(mktemp -d "$CACHE/src.XXXXXX")
	git -C "$ROOT" archive "$SEDUM_MAIN_REF" | tar -x -C "$src"
	(cd "$src" && go build -o "$SEDUM_MAIN" ./cmd/sedum)
	rm -rf "$src"
fi

PROVENANCE="${PROVENANCE_BIN:-$CACHE/bin/provenance-${PROVENANCE_REF:0:12}}"
if [ ! -x "$PROVENANCE" ]; then
	GOBIN="$CACHE/gobin" go install "github.com/livecodelife/provenance/cmd/provenance@$PROVENANCE_REF"
	mv "$CACHE/gobin/provenance" "$PROVENANCE"
fi
export SEDUM SEDUM_MAIN PROVENANCE
SEDUM_HEAD="$(git -C "$ROOT" rev-parse --short=12 HEAD)$(git -C "$ROOT" diff --quiet HEAD -- . 2>/dev/null || echo '+dirty')"

# ---- run --------------------------------------------------------------------------
RESULTS="$CACHE/results.txt"; LOG="$CACHE/run.log"
: >"$RESULTS"; : >"$LOG"
for seq in seq_sedum seq_migrate seq_pregen; do
	TB_RESULTS="$RESULTS" "$HERE/run.sh" "$seq" >>"$LOG" 2>&1 || true
done

# ---- table ------------------------------------------------------------------------
python3 -I "$HERE/table.py" "$RESULTS" "$SEDUM_HEAD" "$PROVENANCE_REF" "$SEDUM_MAIN_REF" "$LOG"
