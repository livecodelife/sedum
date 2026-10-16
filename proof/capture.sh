#!/usr/bin/env bash
# Reproduces every captured output in proof/GROW_PREP.md (prov-2026-3c35801c).
#
#   proof/capture.sh           rebuild sedum at HEAD, run each scenario, write
#                              proof/out/*.txt, and assert what each shows
#   proof/capture.sh --check   the same runs and assertions, writing nothing
#
# Each output file starts with the commit and the exact command that produced
# it. The model is proof/stub_model.py on 127.0.0.1, which answers with canned
# selections and logs every call: no real model is involved, so call counts are
# exact. PROVENANCE_REAL_REPOS, when set, adds the gen_hash parity run.
set -euo pipefail

ROOT=$(git -C "$(dirname "$0")" rev-parse --show-toplevel)
COMMIT=$(git -C "$ROOT" rev-parse --short HEAD)
T=$(mktemp -d)
OUT="$ROOT/proof/out"
[ "${1:-}" = "--check" ] && OUT="$T/out"
mkdir -p "$OUT"
STUB_PID=""
cleanup() { [ -n "$STUB_PID" ] && kill "$STUB_PID" 2>/dev/null; rm -rf "$T"; }
trap cleanup EXIT

fail() { echo "FAIL: $*" >&2; exit 1; }

(cd "$ROOT" && GOTOOLCHAIN=local go build -o "$T/sedum" ./cmd/sedum)
SEDUM="$T/sedum"

PORT=$((18000 + RANDOM % 1000))
python3 "$ROOT/proof/stub_model.py" "$PORT" "$T/routes.json" "$T/calls.log" &
STUB_PID=$!
export OPENAI_BASE_URL="http://127.0.0.1:$PORT/v1" OPENAI_API_KEY=stub
for _ in $(seq 50); do curl -s -o /dev/null "http://127.0.0.1:$PORT/" && break; sleep 0.1; done

W="$T/work"
mkdir -p "$W/provenance"
cp -R "$ROOT/testdata/generators" "$W/generators"
chmod -R u+w "$W/generators"
CTRL="app/controllers/users_controller.rb"

filter() { # filter <name>: a selection adding a before_action to the users controller
	printf '{"invocations": [{"action": "addBeforeFilter", "kwargs": {"controller": "users", "filter": "%s"}}]}' "$1"
}
cat >"$T/routes.json" <<EOF
{"Authenticate every request.": $(python3 -c 'import json,sys; print(json.dumps(sys.argv[1]))' "$(filter authenticate)"),
 "Log every request.": $(python3 -c 'import json,sys; print(json.dumps(sys.argv[1]))' "$(filter log_request)"),
 "Audit every request.": $(python3 -c 'import json,sys; print(json.dumps(sys.argv[1]))' "$(filter authenticate)")}
EOF

record() { # record <id> <created_at> <intent> [extra yaml]
	printf 'id: %s\ncreated_at: "%s"\n%sintent: |\n  %s\naffected_scope:\n  - %s\n' "$1" "$2" "${4:-}" "$3" "$CTRL" >"$W/provenance/$1.yml"
}
record prov-2026-aaaaaaaa 2026-10-01 "Authenticate every request."
record prov-2026-bbbbbbbb 2026-10-02 "Log every request."

LAST_EXIT=0
capture() { # capture <name> <args...>: run in the work dir, write the output file
	local name=$1; shift
	set +e
	(cd "$W" && "$@") >"$T/last.out" 2>&1
	LAST_EXIT=$?
	set -e
	{
		echo "# grow-prep $COMMIT"
		echo "\$ ${*/#$SEDUM/sedum}" | sed "s#$T/##g"
		echo
		sed "s#$T/##g" "$T/last.out"
		echo "[exit $LAST_EXIT]"
	} >"$OUT/$name.txt"
}
calls() { [ -f "$T/calls.log" ] && wc -l <"$T/calls.log" | tr -d ' ' || echo 0; }
reset_calls() { : >"$T/calls.log"; }
grows() { capture "$1" "$SEDUM" grow --generators generators --records provenance --model stub --log "$T/run.log" "${@:2}"; }
seen() { grep -q -- "$1" "$OUT/$2.txt" || fail "$2: expected to see: $1"; }
snapshot() { (cd "$W" && find app -type f -exec sh -c 'printf "%s %s %s\n" "$(shasum -a 256 "$1" | cut -c1-16)" "$(stat -f %m "$1" 2>/dev/null || stat -c %Y "$1")" "$1"' _ {} \;) | sort; }

# 01 - first grow: both records select; markers carry gen and pkg.
reset_calls
grows 01-first-grow
[ "$LAST_EXIT" -eq 0 ] || fail "01: grow exited $LAST_EXIT"
cp "$T/run.log" "$T/first.log"; [ "$(calls)" -eq 4 ] || fail "01: expected 4 model calls, got $(calls)"
grep -c "selection leaves anchors unfilled" "$T/first.log" | grep -qx 2 || fail "01: expected the 2 extra calls to be completeness questions"
echo "model calls: $(calls) = 2 selections + 2 completeness questions about the empty class_body anchor" >>"$OUT/01-first-grow.txt"
capture 01-markers cat "$CTRL"
seen '"gen":"sha256:' 01-markers
seen '"pkg":"sha256:' 01-markers
seen '"record":"prov-2026-aaaaaaaa"' 01-markers

# 02 - second grow over unchanged records: no model call, no file written.
sleep 1
before=$(snapshot)
reset_calls
grows 02-second-grow
[ "$LAST_EXIT" -eq 0 ] || fail "02: grow exited $LAST_EXIT"
[ "$(calls)" -eq 0 ] || fail "02: expected 0 model calls, got $(calls)"
[ "$before" = "$(snapshot)" ] || fail "02: a file's bytes or mtime changed"
seen "current; no model call" 02-second-grow
echo "model calls: $(calls); files unchanged (bytes and mtime)" >>"$OUT/02-second-grow.txt"

# 03 - a template change re-renders from the markers, still with no model call.
printf 'before_action :{{filter}}, prepend: true\n' >"$W/generators/rails/actions/addBeforeFilter.rb"
reset_calls
grows 03-template-change
[ "$(calls)" -eq 0 ] || fail "03: expected 0 model calls, got $(calls)"
seen "re-rendered from their markers" 03-template-change
capture 03-markers cat "$CTRL"
seen "prepend: true" 03-markers
echo "model calls: $(calls)" >>"$OUT/03-template-change.txt"

# 04 - --check reports a changed record as stale, writes nothing, exits 1.
record prov-2026-aaaaaaaa 2026-10-01 "Authenticate every request. Keep it fast."
capture 04-check-stale "$SEDUM" grow --check --generators generators --records provenance --log "$T/run.log"
[ "$LAST_EXIT" -ne 0 ] || fail "04: --check should fail on a stale record"
seen "stale      prov-2026-aaaaaaaa" 04-check-stale
reset_calls
grows 04-regrow
[ "$(calls)" -eq 2 ] || fail "04: re-selecting one changed record should cost 2 calls (selection + completeness), got $(calls)"
echo "model calls: $(calls), both for the changed record (its selection and its completeness question); the unchanged record made none" >>"$OUT/04-regrow.txt"

# 05 - a record selecting another record's region is an error naming both.
record prov-2026-cccccccc 2026-10-03 "Audit every request."
grows 05-ownership --only prov-2026-cccccccc
[ "$LAST_EXIT" -ne 0 ] || fail "05: the ownership conflict should fail"
seen "prov-2026-aaaaaaaa" 05-ownership
seen "prov-2026-cccccccc" 05-ownership
rm "$W/provenance/prov-2026-cccccccc.yml"

# 06 - --stamp writes missing hashes and nothing else; the record is current after.
sed -E -i.bak 's/"gen":"[^"]*",//; s/"pkg":"[^"]*",//' "$W/$CTRL" && rm "$W/$CTRL.bak"
capture 06-check-unstamped "$SEDUM" grow --check --generators generators --records provenance --log "$T/run.log"
[ "$LAST_EXIT" -ne 0 ] || fail "06: markers without hashes should read as stale"
reset_calls
capture 06-stamp "$SEDUM" grow --stamp --generators generators --records provenance --log "$T/run.log"
[ "$(calls)" -eq 0 ] || fail "06: --stamp made a model call"
seen "marker(s) stamped" 06-stamp
capture 06-check-stamped "$SEDUM" grow --check --generators generators --records provenance --log "$T/run.log"
[ "$LAST_EXIT" -eq 0 ] || fail "06: after --stamp the records should be current"

# 07 - a deprecated record's region is an orphan; --prune-orphans removes it.
record prov-2026-bbbbbbbb 2026-10-02 "Log every request." "status: deprecated
"
capture 07-check-orphan "$SEDUM" grow --check --generators generators --records provenance --log "$T/run.log"
[ "$LAST_EXIT" -ne 0 ] || fail "07: --check should fail on an orphan"
seen "orphan" 07-check-orphan
grows 07-prune-orphans --prune-orphans
seen "pruned orphan" 07-prune-orphans
capture 07-markers cat "$CTRL"
grep -q log_request "$OUT/07-markers.txt" && fail "07: the orphan region is still there"

# 08 - sedum regions --json: what provenance's evidence code hash reads.
capture 08-regions-json "$SEDUM" regions --generators generators --output . --json
seen '"format": "sedum.regions/1"' 08-regions-json

# 09 - grow --json: which records' regions changed.
record prov-2026-aaaaaaaa 2026-10-01 "Authenticate every request."
reset_calls
grows 09-grow-json --json
seen '"format": "sedum.grow/1"' 09-grow-json
seen '"changed": true' 09-grow-json

# 10 - the go package: Go that is not gofmt-clean is rejected; clean Go lands.
G="$T/gowork"
mkdir -p "$G"
gorec() { # gorec <code>: a recording of one setFunction into record/store.go
	python3 - "$1" <<'PY'
import json, sys
f = "record/store.go"
print(json.dumps({"sedum_version": "0.8.0", "packages": {"go": {"extensions": [".go"]}},
  "records": [{"record_id": "prov-2026-aaaaaaaa", "files": [{"path": f, "package": "go",
    "template": "{pkg}/{name}.go", "captures": {"pkg": "record", "name": "store"}}],
    "phases": [{"name": "default", "invocations": [
      {"action": "addImport", "kwargs": {"file": f, "path": "strings"}},
      {"action": "setFunction", "kwargs": {"file": f, "name": "Clean", "code": sys.argv[1]}}]}]}]}))
PY
}
gorec $'func Clean(s string) string  {\n    return strings.TrimSpace(s)\n}' >"$G/bad.json"
gorec $'func Clean(s string) string {\n\treturn strings.TrimSpace(s)\n}' >"$G/good.json"
W_SAVE=$W; W=$G
capture 10-go-rejected "$SEDUM" grow --generators "$ROOT/testdata/generators" --output . --execute bad.json --log "$T/run.log"
[ "$LAST_EXIT" -ne 0 ] || fail "10: unformatted Go should be rejected"
seen "gofmt" 10-go-rejected
capture 10-go-accepted "$SEDUM" grow --generators "$ROOT/testdata/generators" --output . --execute good.json --log "$T/run.log"
[ "$LAST_EXIT" -eq 0 ] || fail "10: clean Go should land"
capture 10-go-file cat record/store.go
(cd "$G" && gofmt -l record) | grep -q . && fail "10: the grown file is not gofmt-clean"
W=$W_SAVE

# 11 - gen_hash parity with the shared library, when the real repositories are present.
if [ -n "${PROVENANCE_REAL_REPOS:-}" ]; then
	W_SAVE=$W; W=$ROOT
	capture 11-parity env GOTOOLCHAIN=local go test -count=1 -run TestGenHashMatchesTheLibraryOnRealRecords -v ./internal/record/
	[ "$LAST_EXIT" -eq 0 ] || fail "11: parity test failed"
	W=$W_SAVE
fi

echo "all scenarios passed at grow-prep $COMMIT; outputs in $OUT"
