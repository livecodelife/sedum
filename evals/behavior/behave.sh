#!/usr/bin/env bash
#
# Behaviour harness spike.
#
# Scaffolds a target application the way its framework's own generator does,
# points Sedum at it with an eval provenance record, boots what comes out and
# asserts against it over HTTP. The project is temporary and is deleted on the
# way out; what survives is a results file.
#
# The eval measures which actions a model selects. This measures whether what
# those actions rendered actually runs, which is a different question and needs
# a different fixture: not a checked-in application with the answer already in
# it, but an empty one from `rails new`.
#
#   ./behave.sh todo-rails                       # the target's own answer.json
#   ./behave.sh todo-rails --answer /tmp/s.json  # one eval sample's selection
#   ./behave.sh todo-rails --model qwen2.5-coder-14b-instruct
#   ./behave.sh todo-rails --keep                # leave the project behind
#
# The eval runner drives the --answer form, one call per valid sample, which is
# what makes a behaviour rate a statement about what the model chose rather than
# about a list somebody checked in.
#
set -uo pipefail

# A phase's own output goes to its log; the harness's narration goes here, so
# that both are legible and neither is mixed into the other.
exec 3>&1

HARNESS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# The repo is two levels up from evals/behavior. Derived rather than
# configured, so a clone works without anybody exporting anything.
SEDUM_REPO="${SEDUM_REPO:-$(cd "$HARNESS_DIR/../.." && pwd)}"
SEDUM_BIN="${SEDUM_BIN:-}"

TARGET="${1:-}"
[ -n "$TARGET" ] || { echo "usage: behave.sh <target> [--model <id>] [--keep]" >&2; exit 2; }
shift

MODEL=""          # empty means a replayed selection rather than a live model
KEEP=0
ANSWER=""
while [ $# -gt 0 ]; do
  case "$1" in
    --model)  MODEL="$2";  shift 2 ;;
    --answer) ANSWER="$2"; shift 2 ;;
    --keep)   KEEP=1; shift ;;
    *) echo "unknown flag: $1" >&2; exit 2 ;;
  esac
done

TARGET_DIR="$HARNESS_DIR/targets/$TARGET"
[ -d "$TARGET_DIR" ] || { echo "no such target: $TARGET" >&2; exit 2; }

# The target's own answer.json is the fallback, not the norm. It is a
# hand-authored assertion about what a correct selection looks like, useful for
# exercising the harness itself; a measurement supplies the sample's own.
[ -n "$ANSWER" ] || ANSWER="$TARGET_DIR/answer.json"
[ -f "$ANSWER" ] || { echo "no such answer file: $ANSWER" >&2; exit 2; }

# ---------------------------------------------------------------- bookkeeping

# Dot-prefixed, because the Go toolchain ignores directories beginning with a
# dot or an underscore and this one holds generated .go sources. Named results/
# it was inside the module, so `go build ./...` tried to compile a half-built
# chi service and failed on imports the harness had deliberately not resolved.
RESULTS_DIR="${RESULTS_DIR:-$HARNESS_DIR/.results}"
mkdir -p "$RESULTS_DIR"
# The pid is in the run id because samples run concurrently, and two starting
# in the same second would otherwise write one results file between them.
RUN_ID="$TARGET-$(date +%Y%m%dT%H%M%S)-$$"
RESULTS="$RESULTS_DIR/$RUN_ID.json"

WORK="$(mktemp -d "${TMPDIR:-/tmp}/sedum-behave-XXXXXX")"
APP="$WORK/app"
LOGS="$WORK/logs"
mkdir -p "$LOGS"

PHASES_JSON="[]"
CHECKS_JSON="[]"
OUTCOME="ok"
FAILED_PHASE=""
STUB_PID=""
APP_PID=""

# A scaffold may refuse a request whose user agent it does not recognise, which
# has nothing to do with the code under test.
UA="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0 Safari/537.36"

# phase <name> <command...> - runs it, times it, records it, and stops the run
# on failure. Every phase's output goes to its own log so a failure is
# explicable after the project has been deleted.
phase() {
  local name="$1"; shift
  [ "$OUTCOME" = ok ] || return 0

  local started ended status
  started=$(date +%s)
  echo "  → $name" >&3
  # Not a subshell. A phase sets variables the next one reads - the port the
  # application came up on, the pid to kill - and a subshell would discard them
  # along with every assertion the verify phase recorded.
  "$@" > "$LOGS/$name.log" 2>&1
  status=$?
  ended=$(date +%s)

  PHASES_JSON=$(jq -c --arg n "$name" --argjson s "$status" --argjson d "$((ended-started))" \
    '. + [{phase:$n, status:$s, seconds:$d}]' <<<"$PHASES_JSON")

  if [ $status -ne 0 ]; then
    OUTCOME="failed"
    FAILED_PHASE="$name"
    echo "    ✗ $name failed (exit $status)" >&3
    tail -25 "$LOGS/$name.log" | sed 's/^/      /' >&3
  fi
  return 0
}

# check <name> <actual> <expected> - one behavioural assertion.
check() {
  local name="$1" actual="$2" expected="$3" ok=false
  [ "$actual" = "$expected" ] && ok=true
  CHECKS_JSON=$(jq -c --arg n "$name" --arg a "$actual" --arg e "$expected" --argjson ok "$ok" \
    '. + [{check:$n, expected:$e, actual:$a, pass:$ok}]' <<<"$CHECKS_JSON")
  if [ "$ok" = true ]; then
    echo "    ✓ $name" >&3
  else
    echo "    ✗ $name: expected $expected, got $actual" >&3
  fi
}

# ------------------------------------------------------------------ utilities

free_port() {
  python3 -c 'import socket;s=socket.socket();s.bind(("127.0.0.1",0));print(s.getsockname()[1]);s.close()'
}

# wait_for <url> <seconds> - polls until something answers at all.
wait_for() {
  local url="$1" limit="${2:-60}" i=0
  while [ "$i" -lt "$limit" ]; do
    if curl -s -o /dev/null --max-time 2 -H "User-Agent: $UA" "$url"; then return 0; fi
    sleep 1; i=$((i+1))
  done
  return 1
}

# status <method> <path> [body] - the HTTP status of one request. The response
# body is left in $WORK/last-body.json for body() to read.
status() {
  local method="$1" path="$2" body="${3:-}"
  if [ -n "$body" ]; then
    curl -s -o "$WORK/last-body.json" -w '%{http_code}' -X "$method" \
      -H 'Content-Type: application/json' -H "User-Agent: $UA" \
      --max-time 15 -d "$body" "$BASE_URL$path"
  else
    curl -s -o "$WORK/last-body.json" -w '%{http_code}' -X "$method" \
      -H "User-Agent: $UA" --max-time 15 "$BASE_URL$path"
  fi
}

body() { jq -r "$1" < "$WORK/last-body.json" 2>/dev/null || echo "<unparseable>"; }

# --------------------------------------------------------------- sedum itself

run_sedum() {
  local base key model_arg
  # Built from the tree under test rather than found on PATH. A behaviour run
  # is a claim about this commit's Sedum, and an installed binary of unknown
  # vintage would make it a claim about nothing in particular.
  if [ -z "$SEDUM_BIN" ]; then
    SEDUM_BIN="$WORK/sedum"
    (cd "$SEDUM_REPO" && go build -o "$SEDUM_BIN" ./cmd/sedum) || return 1
  fi

  if [ -z "$MODEL" ]; then
    STUB_PORT=$(free_port)
    python3 "$HARNESS_DIR/stub_model.py" "$ANSWER" "$STUB_PORT" &
    STUB_PID=$!
    wait_for "http://127.0.0.1:$STUB_PORT/v1/models" 20 || return 1
    base="http://127.0.0.1:$STUB_PORT/v1"; key=stub; model_arg=canned
  else
    base="${OPENAI_BASE_URL:-http://localhost:1234/v1}"
    key="${OPENAI_API_KEY:-x}"
    model_arg="$MODEL"
  fi

  OPENAI_BASE_URL="$base" OPENAI_API_KEY="$key" \
    "$SEDUM_BIN" grow \
      --generators "$SEDUM_REPO/$GENERATORS" \
      --records "$SEDUM_REPO/$RECORDS" \
      --output "$APP" \
      --model "$model_arg" \
      ${SEDUM_VARS:-} \
      --retries "${RETRIES:-2}" \
      --record "$LOGS/recording.json" \
      --log "$LOGS/sedum.log"
      # Both under $LOGS, which is copied out. The run log holds the prompt and
      # every response, and is the only record of what the model was told - so
      # writing it into the project that gets deleted threw away the evidence
      # for the one phase that cannot be re-derived from the output.
}

# ------------------------------------------------------------------- teardown

cleanup() {
  [ -n "$APP_PID" ] && kill "$APP_PID" 2>/dev/null
  [ -n "$STUB_PID" ] && kill "$STUB_PID" 2>/dev/null
  if declare -f target_teardown >/dev/null; then target_teardown >"$LOGS/teardown.log" 2>&1; fi

  local passed total
  passed=$(jq '[.[]|select(.pass)]|length' <<<"$CHECKS_JSON")
  total=$(jq 'length' <<<"$CHECKS_JSON")
  if [ "$OUTCOME" = ok ] && [ "$passed" != "$total" ]; then OUTCOME="checks_failed"; fi

  jq -n --arg run "$RUN_ID" --arg target "$TARGET" \
        --arg model "${MODEL:-canned}" --arg outcome "$OUTCOME" \
        --arg failed "$FAILED_PHASE" --arg at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
        --argjson phases "$PHASES_JSON" --argjson checks "$CHECKS_JSON" \
        --argjson passed "$passed" --argjson total "$total" \
    '{run:$run, target:$target, selection:$model, at:$at, outcome:$outcome,
      failed_phase:(if $failed=="" then null else $failed end),
      checks_passed:$passed, checks_total:$total,
      phases:$phases, checks:$checks}' > "$RESULTS"

  # The generated sources are the evidence for a failed check, so they are kept
  # beside the results rather than deleted with the project.
  if [ -d "$APP" ]; then
    mkdir -p "$RESULTS_DIR/$RUN_ID-generated"
    (cd "$APP" && for p in $(printf '%s\n' "${GENERATED_PATHS[@]:-}"); do
       [ -f "$p" ] && { mkdir -p "$RESULTS_DIR/$RUN_ID-generated/$(dirname "$p")"; cp "$p" "$RESULTS_DIR/$RUN_ID-generated/$p"; }
     done) 2>/dev/null
  fi
  cp -R "$LOGS" "$RESULTS_DIR/$RUN_ID-logs" 2>/dev/null

  if [ "$KEEP" = 1 ]; then
    echo; echo "project kept at $WORK"
  else
    rm -rf "$WORK"
  fi

  echo
  echo "outcome: $OUTCOME   checks: $passed/$total"
  echo "results: $RESULTS"
}
trap cleanup EXIT

# ------------------------------------------------------------------- the run

# shellcheck source=/dev/null
source "$TARGET_DIR/target.sh"

echo "behaviour run $RUN_ID"
echo "  project:   $APP"
echo "  selection: ${MODEL:-canned}"

phase scaffold target_scaffold

# What the framework's generator wrote for a path the record authorizes is
# removed, so Sedum creates it from its own file template. A scaffold's
# config/routes.rb carries no anchor, and a file already present is left as it
# is - so without this the run either fails its marker check or injects into
# nothing.
phase prepare  target_prepare
phase generate run_sedum
phase build    target_build
phase boot     target_boot
phase verify   target_verify
