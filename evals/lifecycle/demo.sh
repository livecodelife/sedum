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

# Run from inside a git hook (e.g. a provenance spec run by pre-commit), the hook's GIT_DIR,
# GIT_INDEX_FILE, GIT_WORK_TREE... would point every git command here at the enclosing repo.
for v in $(env | sed -n 's/^\(GIT_[A-Za-z0-9_]*\)=.*/\1/p'); do unset "$v"; done
# Isolate git from the user's global and system config (herd rule, COMMON.md "Git in tests and
# scripts"). HOME itself is left alone: Go's build and module caches live under it.
GIT_CONFIG_GLOBAL="$(mktemp -t lifecycle-gitconfig.XXXXXX)"; export GIT_CONFIG_GLOBAL GIT_CONFIG_NOSYSTEM=1
trap 'rm -f "$GIT_CONFIG_GLOBAL"' EXIT

HERE="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(git -C "$HERE" rev-parse --show-toplevel)"

# ---- pins ------------------------------------------------------------------------
SEDUM_MAIN_REF="${SEDUM_MAIN_REF:-0571de7de06364e9cbc0076098f96c3a4e9ee7cc}"  # sedum v0.8.0, for step 2'
PROVENANCE_REF="${PROVENANCE_REF:-1b0412309fd28ee59bd61752df3f0c41fca02148}"  # provenance, branch cli

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
