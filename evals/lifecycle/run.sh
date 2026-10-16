#!/usr/bin/env bash
# Lifecycle testbed step runner. Plan: linespec-migration-dry-run/docs/LIFECYCLE_TESTBED.md.
#
#   tools/run.sh today          # every step that is runnable today, in order, from a fresh work/
#   tools/run.sh <step> ...     # named steps against the current work/ (step0 resets it)
#
# Each step asserts its checks and exits non-zero on the first failure.
# The work tree is work/ (gitignored); fixtures/ is never modified.
set -euo pipefail

TB="$(cd "$(dirname "$0")" && pwd)"
W="$TB/work"
SEDUM="${SEDUM:-$TB/bin/sedum}"                 # built from sedum-grow-prep
SEDUM_MAIN="${SEDUM_MAIN:-$TB/bin/sedum-main}"  # built from sedum main (v0.8.0), for step 2'
PROVENANCE="${PROVENANCE:-$TB/bin/provenance}"  # built from provenance-cli (branch cli)
# The chi package: the testbed's own copy, or (inside sedum, at evals/lifecycle) the checkout's.
if [ -z "${CHI:-}" ]; then
	if [ -d "$TB/fixtures/generators/chi" ]; then CHI="$TB/fixtures/generators/chi"
	else CHI="$(git -C "$TB" rev-parse --show-toplevel)/testdata/generators/chi"; fi
fi
RESULTS="${TB_RESULTS:-/dev/null}"              # "<step> pass|FAIL|skipped <detail>" per line
CURRENT=""
PORT="${TB_PORT:-18765}"

BRIEF=prov-2026-7e570001 BP_API=prov-2026-7e570002 NOTES_OLD=prov-2026-7e570003
HEALTH=prov-2026-7e570004 LEGACY=prov-2026-7e570005 NOTES=prov-2026-7e570006
BUG1=prov-2026-7e570007 DUP=prov-2026-7e570008
V1_IMPRINTS=("$NOTES_OLD" "$HEALTH" "$LEGACY")  # the imprints step 1 adds
IMPRINTS=("${V1_IMPRINTS[@]}")                  # the imprints that grow now; steps 10 and 11 change it

# Stub model only. Herd rule (COMMON.md, "Sedum model calls"): real model runs go through
# ~/workspace/.herd/sedum-haiku; this harness must never reach a real endpoint.
export OPENAI_BASE_URL="http://127.0.0.1:$PORT/v1" OPENAI_API_KEY=stub
case "$OPENAI_BASE_URL" in http://127.0.0.1:*) ;; *) echo "refusing: model endpoint is not the local stub" >&2; exit 1 ;; esac

# ---- helpers -------------------------------------------------------------------

say()  { printf '\n== %s\n' "$*"; }
ok()   { printf '   ok: %s\n' "$*"; }
fail() {
	printf '   FAIL: %s\n' "$*" >&2
	[ -z "$CURRENT" ] || printf '%s FAIL %s\n' "$CURRENT" "$*" >>"$RESULTS"
	exit 1
}

add_records() {  # add_records <stepN>: fixtures are stored as <hex>.yml (see make_records.py)
	local f
	for f in "$TB/fixtures/$1/records/"*.yml; do cp "$f" "$W/provenance/prov-2026-$(basename "$f")"; done
}

v1_lint() {  # the v1 records are valid v1; linespec is optional so the demo runs without it
	if ! command -v linespec >/dev/null; then ok "linespec not installed: v1 lint skipped"; return; fi
	(cd "$W" && linespec provenance lint -c .linespec.yml >"$W/.stub/last.out" 2>&1) || { cat "$W/.stub/last.out"; fail "linespec provenance lint"; }
	ok "linespec provenance lint"
}

commit() { git -C "$W" add -A && git -C "$W" commit -qm "$1"; }

grow() {  # grow <binary> <record ids...> [-- extra flags]; output in $W/.stub/last.out
	local bin="$1"; shift
	local ids=() extra=()
	while [ $# -gt 0 ]; do
		if [ "$1" = "--" ]; then shift; extra=("$@"); break; fi
		ids+=("$1"); shift
	done
	local args=() id
	for id in "${ids[@]}"; do args+=(--only "$id"); done
	(cd "$W" && "$bin" grow --generators generators --records provenance --model stub \
		"${args[@]}" ${extra[@]+"${extra[@]}"}) >"$W/.stub/last.out" 2>&1
}

check() {  # sedum grow --check; sets CHECK_RC
	local args=() id
	for id in "$@"; do args+=(--only "$id"); done
	set +e
	(cd "$W" && "$SEDUM" grow --check --generators generators --records provenance "${args[@]}") >"$W/.stub/last.out" 2>&1
	CHECK_RC=$?
	set -e
}

calls()      { wc -l <"$W/.stub/calls.log" | tr -d ' '; }
mark_calls() { CALLS0=$(calls); }
expect_calls() {
	local got=$(( $(calls) - CALLS0 ))
	[ "$got" -eq "$1" ] || { cat "$W/.stub/last.out"; fail "expected $1 model call(s), got $got"; }
	ok "$1 model call(s)"
}

# content hash + mtime of every app file
snapshot() { (cd "$W" && find internal -type f | sort | while read -r f; do
	printf '%s %s %s\n' "$(shasum -a 256 "$f" | cut -c1-16)" "$(stat -f %m "$f")" "$f"; done) >"$W/.stub/$1"; }
changed_files() {  # files whose content changed between two snapshots
	join -j 3 -o 1.1,2.1,0 <(sort -k3 "$W/.stub/$1") <(sort -k3 "$W/.stub/$2") | awk '$1 != $2 {print $3}'
}
expect_changed() {  # expect_changed <before> <after> <files...>  (exact set; mtimes of others unchanged)
	local before="$1" after="$2"; shift 2
	local got want
	got=$(changed_files "$before" "$after" | sort | tr '\n' ' ')
	want=$(printf '%s\n' "$@" | sed '/^$/d' | sort | tr '\n' ' ')
	[ "$got" = "$want" ] || fail "changed files: want [$want] got [$got]"
	local untouched; untouched=$(join -j 3 -o 1.2,2.2,0 <(sort -k3 "$W/.stub/$before") <(sort -k3 "$W/.stub/$after") |
		awk -v ch=" $want" '$1 != $2 && index(ch, " " $3 " ") == 0 {print $3}')
	[ -z "$untouched" ] || fail "mtime changed on unchanged files: $untouched"
	ok "changed: [${want% }]"
}
expect_build() { (cd "$W" && go build ./internal/...) || fail "go build ./internal/..."; ok "go build ./internal/..."; }

markers() { grep -rho 'sedum:[A-Za-z]* {[^}]*}[^}]*}' "$W/internal" 2>/dev/null || true; }
owner_of_file() { grep -ho '"record":"[^"]*"' "$W/internal/handlers/$1" | sort -u | cut -d'"' -f4; }

set_intent() {  # set_intent <id> <new intent>   (hand edit of a v1 record)
	python3 -I - "$W/provenance/$1.yml" "$2" <<'EOF'
import json, re, sys
p, intent = sys.argv[1], sys.argv[2]
s = open(p).read()
s, n = re.subn(r'^intent: .*$', 'intent: ' + json.dumps(intent), s, count=1, flags=re.M)
assert n == 1, p
open(p, "w").write(s)
EOF
}

stub_start() {
	stub_stop
	if curl -s -o /dev/null "http://127.0.0.1:$PORT/" 2>/dev/null; then fail "port $PORT is busy (a leftover stub?)"; fi
	python3 -I "$TB/stub_router.py" "$TB/fixtures/answers/routes.json" "$PORT" "$W/.stub/calls.log" &
	echo $! >"$W/.stub/pid"
	for _ in $(seq 50); do curl -s -o /dev/null "http://127.0.0.1:$PORT/" 2>/dev/null && return; sleep 0.1; done
}
stub_stop() {
	if [ -f "$W/.stub/pid" ]; then kill "$(cat "$W/.stub/pid")" 2>/dev/null || true; fi
	rm -f "$W/.stub/pid"
}

# ---- steps ------------------------------------------------------------------------

step0() {
	say "0. scaffold"
	stub_stop; rm -rf "$W"; mkdir -p "$W/.stub"; : >"$W/.stub/calls.log"
	mkdir -p "$W/generators"; cp -R "$CHI" "$W/generators/chi"; chmod -R u+w "$W/generators"
	mkdir "$W/provenance"
	(cd "$W" && go mod init testbed >/dev/null 2>&1)
	# work/ is its own git repo: migrate orders versions by git history and refuses a dirty records dir
	git -C "$W" init -q
	git -C "$W" config user.name testbed; git -C "$W" config user.email testbed@example.invalid
	printf '.stub/\n.sedum/\n' >"$W/.gitignore"
	printf 'provenance:\n  enforcement: none\n  dir: provenance\n' >"$W/.linespec.yml"
	(cd "$W" && "$SEDUM" validate --generators generators >/dev/null) || fail "sedum validate"
	ok "sedum validate"
	stub_start; ok "stub on :$PORT"
}

step1() {
	say "1. v1 records"
	add_records step1
	v1_lint
	! grep -l '^version:' "$W/provenance/"*.yml >/dev/null || fail "a record already has version:"
	ok "5 records, all schema 1"
	commit "step 1: v1 records"
}

step2() {
	say "2. first grow (grow-prep binary)"
	mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow"; }
	expect_calls 3
	for f in note health legacy; do [ -f "$W/internal/handlers/$f.go" ] || fail "missing $f.go"; done
	ok "note.go, health.go, legacy.go created"
	expect_build
	local n; n=$(markers | wc -l | tr -d ' ')
	[ "$n" -eq 11 ] || fail "expected 11 markers, got $n"
	[ "$(markers | grep -c '"gen":"sha256:')" -eq 11 ] && [ "$(markers | grep -c '"pkg":"sha256:')" -eq 11 ] || fail "markers missing gen/pkg"
	ok "11 markers, each with record, gen and pkg"
	[ "$(owner_of_file note.go)" = "$NOTES_OLD" ] && [ "$(owner_of_file health.go)" = "$HEALTH" ] &&
		[ "$(owner_of_file legacy.go)" = "$LEGACY" ] || fail "region owners"
	ok "each file's regions are owned by the imprint that selected them"
	commit "step 2: first grow"
}

step2prime() {
	say "2'. cost without gen: grow with main's sedum (no gen on markers), then with grow-prep"
	step0 >/dev/null; step1 >/dev/null
	IMPRINTS=("${V1_IMPRINTS[@]}")
	mark_calls
	grow "$SEDUM_MAIN" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow (main)"; }
	expect_calls 3
	[ "$(markers | grep -c '"gen":')" -eq 0 ] || fail "main's markers should carry no gen"
	ok "main's markers carry no gen"
	check "${IMPRINTS[@]}"
	[ "$CHECK_RC" -ne 0 ] || fail "--check should fail on gen-less markers"
	ok "--check reports them stale (rc=$CHECK_RC): $(grep -c '^stale' "$W/.stub/last.out" || true) stale"
	mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow (grow-prep)"; }
	expect_calls 3
	ok "first grow-prep grow re-selects every record: the cost grow --stamp (bp4) is meant to remove"
	mark_calls; grow "$SEDUM" "${IMPRINTS[@]}"; expect_calls 0
}

step6() {
	say "6. a second grow does nothing"
	snapshot s6a; mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow"; }
	expect_calls 0
	snapshot s6b; expect_changed s6a s6b
	check "${IMPRINTS[@]}"
	[ "$CHECK_RC" -eq 0 ] || { cat "$W/.stub/last.out"; fail "--check rc=$CHECK_RC"; }
	ok "--check: all current, exit 0"
}

step7v1() {
	say "7 (sedum half). edit one record: only its regions regenerate"
	# provenance revise is pending; a hand edit of the v1 intent drives the same gen_hash change
	set_intent "$HEALTH" "Serve the health resource: GET /health answers 200 when the service is up, and HEAD /health answers the same without a body."
	snapshot s7a; mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow"; }
	expect_calls 1
	snapshot s7b; expect_changed s7a s7b internal/handlers/health.go
	grep -q 'func HeadHealth' "$W/internal/handlers/health.go" || fail "HEAD handler missing"
	ok "health.go gained HeadHealth"
	expect_build
}

step7b() {
	say "7b. package change: re-render without a model call"
	local tpl="$W/generators/chi/actions/addRoute.go"
	local gens_before; gens_before=$(markers | grep -o '"gen":"[^"]*"' | sort | shasum)
	printf '// route %s: {{verb|upper}} /{{resource|kebab}}\n' '(testbed)' >"$tpl"
	snapshot s7ba; mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow"; }
	expect_calls 0
	snapshot s7bb; expect_changed s7ba s7bb internal/handlers/health.go internal/handlers/legacy.go internal/handlers/note.go
	local routes new old
	routes=$(markers | grep -c '^sedum:addRoute ' || true)
	new=$({ grep -rh '// route (testbed):' "$W/internal" || true; } | wc -l | tr -d ' ')
	old=$({ grep -rh '// route: ' "$W/internal" || true; } | wc -l | tr -d ' ')
	[ "$routes" -gt 0 ] && [ "$new" -eq "$routes" ] && [ "$old" -eq 0 ] ||
		fail "addRoute regions not re-rendered ($new new, $old old, $routes markers)"
	ok "all $routes addRoute regions re-rendered"
	[ "$(markers | { grep -o '"gen":"[^"]*"' || true; } | sort | shasum)" = "$gens_before" ] || fail "gen changed on a re-render"
	ok "gen unchanged"
	expect_build
}

step9() {
	say "9. grow --check"
	set_intent "$LEGACY" "Serve the legacy resource: GET /legacy keeps old clients working until they move. It is scheduled for removal."
	snapshot s9a
	check "${IMPRINTS[@]}"
	[ "$CHECK_RC" -ne 0 ] || { cat "$W/.stub/last.out"; fail "--check should exit non-zero"; }
	grep -q "^stale *$LEGACY" "$W/.stub/last.out" || { cat "$W/.stub/last.out"; fail "--check should list $LEGACY as stale"; }
	snapshot s9b; expect_changed s9a s9b
	ok "--check exits $CHECK_RC, lists $LEGACY as stale, writes nothing"
	mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow"; }
	expect_calls 1
	check "${IMPRINTS[@]}"
	[ "$CHECK_RC" -eq 0 ] || { cat "$W/.stub/last.out"; fail "--check after grow rc=$CHECK_RC"; }
	ok "after grow, --check exits 0"
}

step11() {
	say "11. ownership conflict"
	local owner; owner=$(owner_of_file note.go)
	add_records step11
	snapshot s11a
	if grow "$SEDUM" "${IMPRINTS[@]}" "$DUP"; then cat "$W/.stub/last.out"; fail "grow should fail"; fi
	grep -q "$DUP" "$W/.stub/last.out" && grep -q "$owner" "$W/.stub/last.out" ||
		{ cat "$W/.stub/last.out"; fail "error should name $DUP and $owner"; }
	ok "grow fails: $(grep -m1 "$owner" "$W/.stub/last.out" | cut -c1-140)"
	snapshot s11b; expect_changed s11a s11b
	set_intent "$DUP" "A read-only mirror of the note resource for reporting, served as PATCH /note so it never collides with GET."
	mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" "$DUP" || { cat "$W/.stub/last.out"; fail "grow with a new identity"; }
	# exactly 1: the completeness follow-up no longer fires for an anchor another record's
	# region already fills (finding F2, fixed in prov-2026-53f4f82b)
	expect_calls 1
	grep -q 'func PatchNote' "$W/internal/handlers/note.go" || fail "PatchNote missing"
	local last; last=$(grep -o '"record":"[^"]*"' "$W/internal/handlers/note.go" | tail -1 | cut -d'"' -f4)
	[ "$last" = "$DUP" ] || fail "the new region should sit after $owner's (last owner is $last)"
	ok "with a new identity it succeeds; $DUP's region sits after $owner's"
	expect_build
	IMPRINTS+=("$DUP")
}

step12() {
	say "12. --force"
	snapshot s12a; mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" -- --force || { cat "$W/.stub/last.out"; fail "grow --force"; }
	# one per imprint; no completeness follow-up for the shared file (F2 fixed)
	expect_calls "${#IMPRINTS[@]}"
	snapshot s12b
	[ -z "$(changed_files s12a s12b)" ] || fail "content changed under --force: $(changed_files s12a s12b)"
	ok "content byte-identical"
}

step3() {
	say "3. changes made the v1 way"
	add_records step3
	v1_lint
	grep -q '^status: superseded' "$W/provenance/$NOTES_OLD.yml" || fail "$NOTES_OLD not superseded"
	ok "$NOTES supersedes $NOTES_OLD; $BUG1 extends $BP_API; lint passes"
	commit "step 3: v1 supersession and bug"
}

step3grow() {
	say "3x. (observation) growing a v1 supersession before migration"
	snapshot s3a
	local ids=("$NOTES") id
	for id in "${IMPRINTS[@]}"; do [ "$id" = "$NOTES_OLD" ] || ids+=("$id"); done
	if grow "$SEDUM" "${ids[@]}"; then
		ok "grow succeeded (no conflict)"
	else
		ok "grow fails before migration: $(grep -m1 "$NOTES_OLD" "$W/.stub/last.out" | cut -c1-140)"
		snapshot s3b; expect_changed s3a s3b
	fi
}

pending() { say "$1"; printf '   skipped: waiting on %s\n' "$2"; SKIPPED="waiting on $2"; }
set_status() {  # set_status <id> <status>
	sed -i '' "s/^status: .*/status: $2/" "$W/provenance/$1.yml"
}

drop_imprint() { local keep=() id; for id in "${IMPRINTS[@]}"; do [ "$id" = "$1" ] || keep+=("$id"); done; IMPRINTS=("${keep[@]}"); }

strip_markers() { (cd "$W" && find internal -type f | sort | while read -r f; do echo "## $f"; sed '/sedum:/d' "$f"; done); }

orphans_in() { grep -c "^orphan .*record $1 " "$W/.stub/last.out" || true; }

step10() {
	say "10. retire a record and prune its orphans"
	local n; n=$(grep -c '"record":"'"$LEGACY"'"' "$W/internal/handlers/legacy.go")
	[ "$n" -gt 0 ] || fail "$LEGACY owns no regions to retire"
	# (a) deprecated
	set_status "$LEGACY" deprecated
	drop_imprint "$LEGACY"
	if grow "$SEDUM" "$LEGACY"; then fail "growing a deprecated record should be refused"; fi
	ok "grow refuses --only for a deprecated record"
	snapshot s10a; mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow"; }
	expect_calls 0
	[ "$(orphans_in "$LEGACY")" -eq "$n" ] || { cat "$W/.stub/last.out"; fail "expected $n orphans reported"; }
	snapshot s10b; expect_changed s10a s10b
	ok "plain grow reports $n orphans and deletes nothing"
	check "${IMPRINTS[@]}"
	[ "$CHECK_RC" -ne 0 ] || fail "--check should fail on orphans"
	ok "--check exits $CHECK_RC on orphans"
	mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" -- --prune-orphans || { cat "$W/.stub/last.out"; fail "grow --prune-orphans"; }
	expect_calls 0
	[ -f "$W/internal/handlers/legacy.go" ] || fail "legacy.go should stay"
	! grep -q '"record":"'"$LEGACY"'"' "$W/internal/handlers/legacy.go" || fail "$LEGACY regions remain"
	grep -q 'sedum:anchor:imports' "$W/internal/handlers/legacy.go" && grep -q 'sedum:anchor:handlers' "$W/internal/handlers/legacy.go" ||
		fail "legacy.go lost its anchors"
	ok "--prune-orphans removes all $n regions; legacy.go stays with its anchors"
	expect_build
	check "${IMPRINTS[@]}"; [ "$CHECK_RC" -eq 0 ] || { cat "$W/.stub/last.out"; fail "--check after prune rc=$CHECK_RC"; }
	snapshot s10c; grow "$SEDUM" "${IMPRINTS[@]}" -- --prune-orphans; snapshot s10d; expect_changed s10c s10d
	ok "--check exits 0; a second --prune-orphans is a no-op"
	# (b) record file removed
	set_status "$LEGACY" implemented
	mark_calls; grow "$SEDUM" "${IMPRINTS[@]}" "$LEGACY" || { cat "$W/.stub/last.out"; fail "regrow $LEGACY"; }
	expect_calls 1
	rm "$W/provenance/$LEGACY.yml"
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow"; }
	[ "$(orphans_in "$LEGACY")" -eq "$n" ] || { cat "$W/.stub/last.out"; fail "removed record: expected $n orphans"; }
	grow "$SEDUM" "${IMPRINTS[@]}" -- --prune-orphans || { cat "$W/.stub/last.out"; fail "prune"; }
	! grep -q '"record":"'"$LEGACY"'"' "$W/internal/handlers/legacy.go" || fail "$LEGACY regions remain after removal"
	ok "same with the record file removed: $n orphans reported, then pruned"
	expect_build
}

step4b() {
	say "4b. grow --stamp: markers without gen become current with no model call"
	step0 >/dev/null; step1 >/dev/null
	IMPRINTS=("${V1_IMPRINTS[@]}")
	grow "$SEDUM_MAIN" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow (main)"; }
	[ "$(markers | grep -c '"gen":')" -eq 0 ] || fail "main's markers should carry no gen"
	if grow "$SEDUM" "${IMPRINTS[@]}" -- --stamp --force; then fail "--stamp --force should be refused"; fi
	ok "--stamp --force refused"
	strip_markers >"$W/.stub/strip_a"; mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" -- --stamp || { cat "$W/.stub/last.out"; fail "grow --stamp"; }
	expect_calls 0
	local total; total=$(markers | wc -l | tr -d ' ')
	[ "$(markers | grep -c '"gen":"sha256:')" -eq "$total" ] && [ "$(markers | grep -c '"pkg":"sha256:')" -eq "$total" ] ||
		fail "stamp should write gen and pkg on all $total markers"
	ok "gen and pkg stamped onto all $total markers"
	strip_markers >"$W/.stub/strip_b"
	cmp -s "$W/.stub/strip_a" "$W/.stub/strip_b" || fail "--stamp changed something other than marker lines"
	ok "nothing but marker lines changed"
	check "${IMPRINTS[@]}"
	[ "$CHECK_RC" -eq 0 ] || { cat "$W/.stub/last.out"; fail "--check after stamp rc=$CHECK_RC"; }
	ok "--check: all current"
	mark_calls; grow "$SEDUM" "${IMPRINTS[@]}"; expect_calls 0
	ok "next grow makes no model calls (compare step 2': 3)"
}

yget() {  # yget <id> <key>: a top-level scalar from a v2 record
	sed -n "s/^$2: //p" "$W/provenance/$1.yml" | head -1 | sed "s/^[\"']//; s/[\"']\$//"
}
ylist() {  # ylist <id> <key>: items of a top-level list
	awk -v k="$2:" '$0==k {on=1; next} on && /^  - / {sub(/^  - /,""); gsub(/^"|"$/,""); print; next} on {exit}' "$W/provenance/$1.yml"
}

step4() {
	say "4. provenance migrate"
	[ -x "$PROVENANCE" ] || fail "bin/provenance missing (tools/build.sh)"
	(cd "$W" && "$PROVENANCE" migrate --plan --json) >"$W/.stub/plan.json" || fail "migrate --plan"
	python3 -I - "$W/.stub/plan.json" "$NOTES_OLD" "$NOTES" "$BUG1" "$BP_API" <<'PY' || fail "plan"
import json, sys
p = json.load(open(sys.argv[1])); old, new, bug, bp = sys.argv[2:]
got = sorted((e["id"], e["target"], e["rule"], e["kept"]) for e in p["entries"])
want = sorted([(old, new, "1", False), (bug, bp, "4", False)])
assert (p["records_before"], p["records_after"]) == (7, 5), (p["records_before"], p["records_after"])
assert got == want, got
assert not p["suggestions"] and not p["warnings"], (p["suggestions"], p["warnings"])
PY
	[ -z "$(git -C "$W" status --porcelain)" ] || fail "--plan wrote something"
	ok "plan: 7 -> 5; $NOTES_OLD -> $NOTES (rule 1), $BUG1 -> $BP_API (rule 4); no suggestions or warnings; nothing written"
	(cd "$W" && "$PROVENANCE" migrate) >"$W/.stub/last.out" 2>&1 || { cat "$W/.stub/last.out"; fail "migrate"; }
	[ ! -e "$W/provenance/$NOTES_OLD.yml" ] && [ ! -e "$W/provenance/$BUG1.yml" ] || fail "absorbed files remain"
	[ "$(ls "$W/provenance/"*.yml | wc -l | tr -d ' ')" -eq 5 ] || fail "expected 5 records"
	ok "5 records; absorbed files removed"
	[ "$(yget "$NOTES" version)" = 2 ] && [ "$(ylist "$NOTES" aliases)" = "$NOTES_OLD" ] || fail "$NOTES should be v2 with alias $NOTES_OLD"
	ok "$NOTES: v2, alias $NOTES_OLD, reason \"$(yget "$NOTES" reason)\""
	[ "$(yget "$BP_API" version)" = 2 ] && [ "$(ylist "$BP_API" aliases)" = "$BUG1" ] && [ "$(yget "$BP_API" reason_kind)" = fix ] || fail "$BP_API v2 fix"
	ylist "$BP_API" constraints | grep -q 'JSON with an error field' || fail "bug constraint not appended"
	ylist "$BP_API" affected_scope | grep -qx 'internal/handlers/errors.go' || fail "bug scope not unioned"
	ok "$BP_API: v2, alias $BUG1, reason_kind fix, bug constraint appended and scope unioned"
	local id
	for id in "$NOTES" "$HEALTH" "$LEGACY"; do [ "$(yget "$id" implements_version)" = 1 ] || fail "$id should pin $BP_API v1"; done
	ok "imprints pin $BP_API v1 (behind its v2: stale, as intended)"
	[ "$(find "$W/provenance/.versions" -name '*.yml' | wc -l | tr -d ' ')" -eq 7 ] || fail "expected 7 snapshots"
	[ "$(wc -l <"$W/.provenance/versions.jsonl" | tr -d ' ')" -eq 5 ] || fail "expected 5 manifest lines"
	[ -f "$W/.provenance.yml" ] || fail ".provenance.yml missing"
	ok "7 snapshots, 5-line version manifest, config in .provenance.yml"
	commit "step 4: migrate"
	(cd "$W" && "$PROVENANCE" migrate) >"$W/.stub/last.out" 2>&1 || { cat "$W/.stub/last.out"; fail "second migrate"; }
	grep -qi 'nothing to migrate' "$W/.stub/last.out" || { cat "$W/.stub/last.out"; fail "second migrate should be a no-op"; }
	[ -z "$(git -C "$W" status --porcelain)" ] || fail "second migrate wrote something"
	ok "a second migrate: nothing to migrate"
	IMPRINTS=("$NOTES" "$HEALTH" "$LEGACY")
}

step5() {
	say "5. grow after migration"
	check "${IMPRINTS[@]}"
	[ "$CHECK_RC" -ne 0 ] && grep -q "^stale *$NOTES" "$W/.stub/last.out" &&
		grep -q "^current *$HEALTH" "$W/.stub/last.out" && grep -q "^current *$LEGACY" "$W/.stub/last.out" ||
		{ cat "$W/.stub/last.out"; fail "--check: want $NOTES stale, others current"; }
	ok "--check: $NOTES stale (new content); $HEALTH and $LEGACY current (migration left their gen alone)"
	snapshot s5a; mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow"; }
	expect_calls 1
	snapshot s5b; expect_changed s5a s5b internal/handlers/note.go
	[ "$(owner_of_file note.go)" = "$NOTES" ] || fail "note.go markers should all name $NOTES, got: $(owner_of_file note.go | tr '\n' ' ')"
	ok "every note.go marker now names the survivor $NOTES (alias markers rewritten on write)"
	grep -q 'func DeleteNote' "$W/internal/handlers/note.go" || fail "DeleteNote missing"
	ok "note.go gained DeleteNote"
	expect_build
	commit "step 5: grow after migration"
	snapshot s5c; mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow"; }
	expect_calls 0; snapshot s5d; expect_changed s5c s5d
}
step8()  { pending "8. stale / repin" "staleness + repin (Phase 3)"; }

run_step() {  # run_step <step function>: run it and record pass / skipped
	CURRENT="$1"; SKIPPED=""
	"$1"
	if [ -n "$SKIPPED" ]; then printf '%s skipped %s\n' "$1" "$SKIPPED" >>"$RESULTS"
	else printf '%s pass\n' "$1" >>"$RESULTS"; fi
	CURRENT=""
}
run_steps() { local s; for s in "$@"; do run_step "$s"; done; }

# Three independent sequences, each from a fresh work/
seq_sedum()   { IMPRINTS=("${V1_IMPRINTS[@]}"); run_steps step0 step1 step2 step6 step7v1 step7b step9 step11 step12 step10; }
seq_migrate() { IMPRINTS=("${V1_IMPRINTS[@]}"); run_steps step0 step1 step2 step3 step3grow step4 step5 step8; }
seq_pregen()  { IMPRINTS=("${V1_IMPRINTS[@]}"); run_steps step2prime step4b; }

today() {
	seq_sedum; seq_migrate; seq_pregen
	stub_stop
	say "passed: 0-2, 2', 3, 4, 4b, 5, 6, 7 (sedum half), 7b, 9, 10, 11, 12"
}

trap 'stub_stop' EXIT
if [ $# -eq 0 ]; then sed -n '2,8p' "$0"; exit 2; fi
for s in "$@"; do
	case "$s" in
		today) today ;;
		seq_sedum|seq_migrate|seq_pregen) "$s"; stub_stop ;;
		step*) [ -f "$W/.stub/pid" ] || [ "$s" = step0 ] || [ "$s" = step2prime ] || stub_start; run_step "$s" ;;
		*) fail "unknown step $s" ;;
	esac
done
