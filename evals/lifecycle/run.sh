#!/usr/bin/env bash
# Lifecycle testbed step runner. Plan: linespec-migration-dry-run/docs/LIFECYCLE_TESTBED.md.
#
#   tools/run.sh today          # every step that is runnable today, in order, from a fresh work/
#   tools/run.sh <step> ...     # named steps against the current work/ (step0 resets it)
#
# Each step asserts its checks and exits non-zero on the first failure.
# The work tree is work/ (gitignored); fixtures/ is never modified.
set -euo pipefail

# Run from inside a git hook (e.g. a provenance spec run by pre-commit), the hook's GIT_DIR,
# GIT_INDEX_FILE, GIT_WORK_TREE... would point every git command here at the enclosing repo.
for v in $(env | sed -n 's/^\(GIT_[A-Za-z0-9_]*\)=.*/\1/p'); do unset "$v"; done
# Isolate git from the user's global and system config (herd rule, COMMON.md "Git in tests and
# scripts"). HOME itself is left alone: Go's build and module caches live under it.
GIT_CONFIG_GLOBAL="$(mktemp -t lifecycle-gitconfig.XXXXXX)"; export GIT_CONFIG_GLOBAL GIT_CONFIG_NOSYSTEM=1

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

# Newer sedum prints one summary line per orphan fate and lists each region only with --orphans.
ORPHANS_FLAG=(); "$SEDUM" grow --help 2>/dev/null | grep -q -- '--orphans' && ORPHANS_FLAG=(--orphans)

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
	# cp keeps the source mode; inside a sedum checkout under the v2 hooks tracked files are 0444
	for f in "$TB/fixtures/$1/records/"*.yml; do cp "$f" "$W/provenance/prov-2026-$(basename "$f")"; chmod u+w "$W/provenance/prov-2026-$(basename "$f")"; done
}

v1_lint() {  # the v1 records are valid v1; linespec is optional so the demo runs without it
	if ! command -v linespec >/dev/null; then ok "linespec not installed: v1 lint skipped"; return; fi
	(cd "$W" && linespec provenance lint -c .linespec.yml >"$W/.stub/last.out" 2>&1) || { cat "$W/.stub/last.out"; fail "linespec provenance lint"; }
	ok "linespec provenance lint"
}

work_git_guard() {  # never let a work/ git command land in another repository
	local top; top=$(git -C "$W" rev-parse --show-toplevel 2>/dev/null) || fail "work/ is not a git repo"
	[ "$top" = "$(cd "$W" && pwd -P)" ] || fail "git in work/ resolves to $top, not work/; refusing"
}
commit() { work_git_guard; git -C "$W" add -A && git -C "$W" commit -qm "$1"; }

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
	printf '%s %s %s\n' "$(shasum -a 256 "$f" | cut -c1-16)" "$(stat -c %Y "$f" 2>/dev/null || stat -f %m "$f")" "$f"; done) >"$W/.stub/$1"; }
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
	work_git_guard
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
	sed -i.bak "s/^status: .*/status: $2/" "$W/provenance/$1.yml" && rm -f "$W/provenance/$1.yml.bak"
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
	grow "$SEDUM" "${IMPRINTS[@]}" -- ${ORPHANS_FLAG[@]+"${ORPHANS_FLAG[@]}"} || { cat "$W/.stub/last.out"; fail "grow"; }
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
	grow "$SEDUM" "${IMPRINTS[@]}" -- ${ORPHANS_FLAG[@]+"${ORPHANS_FLAG[@]}"} || { cat "$W/.stub/last.out"; fail "grow"; }
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
	work_git_guard   # migrate reads git history and writes in work/
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
step8() {
	say "8. stale children and repin"
	[ -x "$PROVENANCE" ] || fail "bin/provenance missing (tools/build.sh)"
	"$PROVENANCE" stale --help >/dev/null 2>&1 && "$PROVENANCE" repin --help >/dev/null 2>&1 ||
		{ pending "8. stale / repin" "this provenance CLI has no stale/repin (needs cli 5ca39a7 or later)"; return; }
	local lines0 st
	(cd "$W" && "$PROVENANCE" edit "$BP_API" --reason "clarify error shape" \
		--add-constraint "Error bodies are JSON with a message field.") >"$W/.stub/last.out" 2>&1 ||
		{ cat "$W/.stub/last.out"; fail "edit $BP_API"; }
	[ "$(yget "$BP_API" version)" = 3 ] || fail "$BP_API should be at v3"
	commit "step 8: $BP_API v3"
	st=$(cd "$W" && "$PROVENANCE" stale 2>&1) || fail "stale"
	for id in "$HEALTH" "$LEGACY" "$NOTES"; do
		echo "$st" | grep -q "$id@[0-9]* *pinned to $BP_API@1, parent is at 3 *(direct)" ||
			{ echo "$st"; fail "stale should list $id as directly stale (pin 1, parent 3)"; }
	done
	ok "stale: all three imprints directly stale (pin v1, $BP_API at v3)"
	if (cd "$W" && "$PROVENANCE" stale --exit-code) >/dev/null 2>&1; then fail "stale --exit-code should exit non-zero"; fi
	ok "stale --exit-code exits non-zero"
	if (cd "$W" && "$PROVENANCE" repin "$HEALTH" --to 3) >/dev/null 2>&1; then fail "repin without --reason should fail"; fi
	ok "repin without --reason is refused"
	lines0=$(wc -l <"$W/.provenance/versions.jsonl" | tr -d ' ')
	(cd "$W" && "$PROVENANCE" repin "$HEALTH" --to 3 --reason "still correct") >"$W/.stub/last.out" 2>&1 ||
		{ cat "$W/.stub/last.out"; fail "repin $HEALTH"; }
	[ "$(yget "$HEALTH" implements_version)" = 3 ] || fail "$HEALTH should pin $BP_API v3"
	[ "$(yget "$HEALTH" version)" = 1 ] || fail "repin must not change $HEALTH's version"
	[ "$(wc -l <"$W/.provenance/versions.jsonl" | tr -d ' ')" -eq "$lines0" ] || fail "repin wrote a version"
	[ "$(find "$W/provenance/.versions/$HEALTH" -name '*.yml' | wc -l | tr -d ' ')" -eq 1 ] || fail "repin wrote a snapshot"
	[ -s "$W/provenance/.versions/$HEALTH/pins.jsonl" ] || fail "repin should log to pins.jsonl"
	ok "repin: $HEALTH pin 1 -> 3; version unchanged, no snapshot, no new version line; pins.jsonl logged"
	st=$(cd "$W" && "$PROVENANCE" stale 2>&1)
	echo "$st" | grep -q "$HEALTH@" && { echo "$st"; fail "$HEALTH should no longer be stale"; }
	echo "$st" | grep -q "$LEGACY@" && echo "$st" | grep -q "$NOTES@" || { echo "$st"; fail "stale should still list $LEGACY and $NOTES"; }
	ok "stale now lists only $LEGACY and $NOTES"
	commit "step 8: repin $HEALTH"
	snapshot s8a; mark_calls
	grow "$SEDUM" "${IMPRINTS[@]}" || { cat "$W/.stub/last.out"; fail "grow"; }
	expect_calls 0
	snapshot s8b; expect_changed s8a s8b
	check "${IMPRINTS[@]}"
	[ "$CHECK_RC" -eq 0 ] || { cat "$W/.stub/last.out"; fail "--check should be clean: implements_version is not part of gen_hash"; }
	ok "grow after repin: 0 calls, no file touched; --check clean (a pin is not a content change)"
}

ev_state() {  # ev_state <id>: the state column of `provenance check --evidence` (plus [reason]) for the record
	(cd "$W" && SEDUM_BIN="$SEDUM" "$PROVENANCE" check --evidence 2>/dev/null || true) |
		awk -v id="$1" '$1 ~ "^" id "@" {$1=""; sub(/^ +/,""); print; exit}' | tr -s ' '
}
stale_ids() {  # ids `provenance stale` lists, sorted
	(cd "$W" && "$PROVENANCE" stale 2>/dev/null) | awk '/@[0-9]+ .*pinned to/ {sub(/@.*/,"",$1); print $1}' | sort
}

step13() {
	say "13. evidence flow (DESIGN Phase 3 done-check)"
	"$PROVENANCE" check --help 2>&1 | grep -q -- '--evidence' ||
		{ pending "13. evidence" "this provenance CLI has no check --evidence (needs cli aae06e5 or later)"; return; }
	local child spec=specs/ping.sh want got
	printf 'exclude_paths: ["generators/**"]\nrun_associated_specs_on_complete: true\ncommit_on_status_change: true\n' >>"$W/.provenance.yml"
	mkdir -p "$W/specs"
	printf '#!/bin/sh\ntest -f internal/handlers/ping.go && go build ./internal/...\n' >"$W/$spec"
	commit "step 13: config + spec"
	# a new imprint of BP-API that sedum grows (so it owns regions and has a code hash)
	(cd "$W" && "$PROVENANCE" create --type imprint --implements "$BP_API" --title "Ping endpoint" \
		--intent "Serve the ping resource: GET /ping answers 200 with pong." \
		--scope internal/handlers/ping.go --no-edit) >"$W/.stub/last.out" 2>&1 || { cat "$W/.stub/last.out"; fail "create"; }
	child=$(sed -n 's|.*provenance/\(prov-[0-9]*-[0-9a-f]*\)\.yml.*|\1|p' "$W/.stub/last.out" | head -1)
	CHILD13=$child
	[ -n "$child" ] || { cat "$W/.stub/last.out"; fail "could not read the new record id"; }
	python3 -I - "$W/provenance/$child.yml" "$spec" <<'EOF2'
import sys
p, spec = sys.argv[1], sys.argv[2]
s = open(p).read()
assert "associated_specs: []\n" in s
open(p, "w").write(s.replace("associated_specs: []\n", "associated_specs:\n  - path: %s\n    type: code\n    run_command: sh\n" % spec))
EOF2
	commit "step 13: draft $child"
	(cd "$W" && "$PROVENANCE" open --record "$child") >"$W/.stub/last.out" 2>&1 || { cat "$W/.stub/last.out"; fail "open"; }
	mark_calls; grow "$SEDUM" "$child" || { cat "$W/.stub/last.out"; fail "grow $child"; }
	expect_calls 1
	[ -f "$W/internal/handlers/ping.go" ] || fail "ping.go missing"
	commit "step 13: grow $child"
	(cd "$W" && SEDUM_BIN="$SEDUM" "$PROVENANCE" complete --record "$child") >"$W/.stub/last.out" 2>&1 ||
		{ cat "$W/.stub/last.out"; fail "complete $child"; }
	[ -n "$(ls "$W/provenance/.evidence/$child" 2>/dev/null)" ] || fail "complete wrote no evidence"
	[ -z "$(git -C "$W" status --porcelain)" ] || fail "complete left the tree dirty (evidence not committed with the record)"
	[ "$(ev_state "$child")" = "fresh" ] || { ev_state "$child"; fail "$child evidence should be fresh after complete"; }
	ok "complete: spec ran, evidence committed with the record; check --evidence: fresh"

	# bump the parent: the evidence goes stale for the parent, and stale lists exactly BP-API's direct children
	(cd "$W" && "$PROVENANCE" edit "$BP_API" --reason "name the error body" \
		--add-constraint "Error bodies name the failing field.") >"$W/.stub/last.out" 2>&1 || { cat "$W/.stub/last.out"; fail "edit $BP_API"; }
	commit "step 13: $BP_API bump"
	[ "$(ev_state "$child")" = "stale [parent]" ] || { ev_state "$child"; fail "$child evidence should be stale [parent]"; }
	want=$(for f in "$W"/provenance/prov-*.yml; do id=$(basename "$f" .yml)
		[ "$(yget "$id" implements)" = "$BP_API" ] && echo "$id"; done | sort)
	got=$(stale_ids)
	[ "$got" = "$want" ] && [ "$(echo "$got" | wc -l | tr -d ' ')" -eq 4 ] ||
		{ echo "want: $want"; echo "got: $got"; fail "stale should list exactly $BP_API's 4 descendants"; }
	ok "parent bump: $child evidence stale [parent]; stale lists exactly the 4 descendants of $BP_API"
	if (cd "$W" && SEDUM_BIN="$SEDUM" "$PROVENANCE" check --evidence) >/dev/null 2>&1; then fail "check --evidence should exit non-zero with stale evidence"; fi
	ok "check --evidence exits non-zero while evidence is stale"

	# repin: fresh again, and only that record leaves the stale list
	(cd "$W" && "$PROVENANCE" repin "$child" --reason "still correct") >"$W/.stub/last.out" 2>&1 || { cat "$W/.stub/last.out"; fail "repin"; }
	commit "step 13: repin $child"
	[ "$(ev_state "$child")" = "fresh" ] || { ev_state "$child"; fail "$child evidence should be fresh after repin"; }
	[ "$(stale_ids)" = "$(echo "$want" | grep -v "^$child\$")" ] || fail "stale should now list the other 3"
	ok "repin: $child fresh again; the other 3 still stale"

	# edit the spec: evidence stale for specs
	printf '# edited\n' >>"$W/$spec"
	case "$(ev_state "$child")" in stale*specs*) ;; *) ev_state "$child"; fail "spec edit should make evidence stale (specs)" ;; esac
	ok "spec edit: evidence stale [$(ev_state "$child" | sed 's/.*\[\(.*\)\]/\1/')]"
	git -C "$W" checkout -q -- "$spec"
	[ "$(ev_state "$child")" = "fresh" ] || fail "restoring the spec should make evidence fresh again"

	# check --commit on the tagged complete commit
	local sha; sha=$(git -C "$W" log --format=%H --grep "Complete provenance record \[$child@1\]" -1)
	[ -n "$sha" ] || fail "no tagged complete commit"
	(cd "$W" && "$PROVENANCE" check --commit "$sha") >"$W/.stub/last.out" 2>&1 || { cat "$W/.stub/last.out"; fail "check --commit $sha"; }
	ok "check --commit ${sha:0:7} (tagged $child@1): passes"
}

replay() {  # replay <recording>: sedum grow --execute (no model); output in $W/.stub/last.out
	(cd "$W" && "$SEDUM" grow --execute "$1" --generators generators --records provenance) >"$W/.stub/last.out" 2>&1
}

step14() {
	say "14. template edit + replay: what keeps the hashes right, and where the gap is"
	[ -n "${CHILD13:-}" ] || { pending "14. template replay" "needs step 13's child (this CLI has no check --evidence)"; return; }
	local c=$CHILD13 tpl=generators/chi/actions/addHandlerFunc.go rec="$W/.stub/rec14.json" rec2="$W/.stub/rec14b.json"
	mark_calls; grow "$SEDUM" "$c" -- --force --record "$rec"; expect_calls 1
	[ -z "$(git -C "$W" status --porcelain)" ] || fail "recording the current run should change nothing"
	[ "$(ev_state "$c")" = "fresh" ] || fail "$c should start fresh"

	# A. edit the generator template; replay the recording with no model
	printf '// handled by the template v2\n%s' "$(cat "$W/$tpl")" >"$W/$tpl"
	check "$c"
	[ "$CHECK_RC" -eq 1 ] && grep -q "1 to re-render" "$W/.stub/last.out" || { cat "$W/.stub/last.out"; fail "--check should want a re-render after a template edit"; }
	ok "A: template edit: grow --check exits 1 (1 to re-render)"
	mark_calls; replay "$rec" || { cat "$W/.stub/last.out"; fail "replay"; }; expect_calls 0
	grep -q "handled by the template v2" "$W/internal/handlers/ping.go" || fail "replay should re-render the handler"
	expect_build
	check "$c"; [ "$CHECK_RC" -eq 0 ] || { cat "$W/.stub/last.out"; fail "--check should be clean after replay (pkg re-stamped by sedum)"; }
	ok "A: replay (0 calls) re-renders and re-stamps the markers: grow --check clean, no restamp command needed"
	case "$(ev_state "$c")" in "stale [code]") ;; *) ev_state "$c"; fail "evidence should be stale [code] after the code moved" ;; esac
	ok "A: check --evidence: stale [code] (the code moved; the spec evidence did not)"
	commit "step 14: template v2 + replay"
	(cd "$W" && "$PROVENANCE" run-specs --record "$c") >/dev/null 2>&1 || fail "run-specs"
	[ "$(ev_state "$c")" = "stale [code]" ] || fail "run-specs must not refresh evidence"
	if (cd "$W" && SEDUM_BIN="$SEDUM" "$PROVENANCE" complete --record "$c") >/dev/null 2>&1; then fail "complete on an implemented record should be refused"; fi
	ok "A: GAP: run-specs does not refresh evidence and complete refuses an implemented record: no sanctioned re-evidence without a new version"

	# B. the sanctioned path today: edit (new version, reopened) -> grow -> complete
	(cd "$W" && "$PROVENANCE" edit "$c" --reason "template v2 regenerated" \
		--add-constraint "Handler bodies carry the template comment.") >"$W/.stub/last.out" 2>&1 || { cat "$W/.stub/last.out"; fail "edit $c"; }
	commit "step 14: $c v2"
	check "$c"; [ "$CHECK_RC" -eq 1 ] || fail "--check should be stale after the record edit (gen_hash moved)"
	mark_calls; replay "$rec" || { cat "$W/.stub/last.out"; fail "replay (old recording)"; }; expect_calls 0
	check "$c"; [ "$CHECK_RC" -eq 1 ] || { cat "$W/.stub/last.out"; fail "replaying the OLD recording should leave --check stale: it stamps the recording's gen"; }
	ok "B: record edit: --check stale; replay of the old recording (0 calls) re-stamps the OLD gen, so it stays stale"
	mark_calls; grow "$SEDUM" "$c" -- --force --record "$rec2"; expect_calls 1
	check "$c"; [ "$CHECK_RC" -eq 0 ] || fail "--check should be clean after a fresh grow"
	commit "step 14: regrow $c"
	(cd "$W" && SEDUM_BIN="$SEDUM" "$PROVENANCE" complete --record "$c") >"$W/.stub/last.out" 2>&1 || { cat "$W/.stub/last.out"; fail "complete v2"; }
	[ "$(ev_state "$c")" = "fresh" ] || { ev_state "$c"; fail "evidence should be fresh after complete @2"; }
	ok "B: fresh recording (1 call) + complete @2: --check clean, evidence fresh"

	# C. edit the recorded invocation (kwargs) and replay the fresh recording
	mark_calls; sed -i.bak 's/"verb": "get"/"verb": "post"/' "$rec2"
	replay "$rec2" || { cat "$W/.stub/last.out"; fail "replay (edited recording)"; }; expect_calls 0
	grep -q 'func PostPing' "$W/internal/handlers/ping.go" && grep -q 'func GetPing' "$W/internal/handlers/ping.go" || fail "kwargs edit should add PostPing beside GetPing"
	expect_build
	check "$c"; [ "$CHECK_RC" -eq 0 ] || fail "--check should be clean after the kwargs replay"
	case "$(ev_state "$c")" in "stale [code]") ;; *) ev_state "$c"; fail "evidence should be stale [code] after the kwargs replay" ;; esac
	ok "C: kwargs edit + replay (0 calls): adds PostPing beside GetPing (new region, old kept); --check clean; evidence stale [code]"
	git -C "$W" checkout -q -- . ; rm -f "$rec2.bak"
}

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
seq_migrate() { IMPRINTS=("${V1_IMPRINTS[@]}"); run_steps step0 step1 step2 step3 step3grow step4 step5 step8 step13 step14; }
seq_pregen()  { IMPRINTS=("${V1_IMPRINTS[@]}"); run_steps step2prime step4b; }

today() {
	seq_sedum; seq_migrate; seq_pregen
	stub_stop
	say "passed: 0-2, 2', 3, 4, 4b, 5, 6, 7 (sedum half), 7b, 9, 10, 11, 12"
}

trap 'stub_stop; rm -f "$GIT_CONFIG_GLOBAL"' EXIT
if [ $# -eq 0 ]; then sed -n '2,8p' "$0"; exit 2; fi
for s in "$@"; do
	case "$s" in
		today) today ;;
		seq_sedum|seq_migrate|seq_pregen) "$s"; stub_stop ;;
		step*) [ -f "$W/.stub/pid" ] || [ "$s" = step0 ] || [ "$s" = step2prime ] || stub_start; run_step "$s" ;;
		*) fail "unknown step $s" ;;
	esac
done
