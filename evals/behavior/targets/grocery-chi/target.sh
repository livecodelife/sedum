# A second resource over the Go package set todo-chi already uses.
#
# Same generators, different record. It exists to tell a rule from a copied
# string: the fragment descriptions in that package give their examples in terms
# of todos, and the run that validated them generated todos - so it could not
# say whether the model read the rule or matched the example
# (prov-2026-b585fbbc).
#
# groceries shares no field name and no type with todos, so a copied binding
# does not compile rather than quietly working.

GENERATORS="evals/testdata/todo-api/generators/defined"
RECORDS="evals/testdata/grocery-api/records"

GENERATED_PATHS=(
  main.go
  init.sql
  db/groceries.go
  db/groceries_test.go
  handlers/groceries.go
  handlers/groceries_test.go
)

# The case is the single source. RunBehavior exports the case's variables as
# SEDUM_VAR_<NAME>, and both go mod init and --var read this one value, so the
# scaffold and the generation cannot disagree about the module's name
# (prov-2026-1c33a50b).
MODULE="${SEDUM_VAR_MODULE:-${MODULE:-grocery}}"
# The module path is a project fact the standard cannot know, so the run supplies
# it rather than the model guessing it (prov-2026-6fc3d13d). This is the same
# value `go mod init` was given, which is the point: one fact, named once.
SEDUM_VARS="--var module=$MODULE"
PGURL="${PGURL:-postgres://$(whoami)@localhost:5432}"
DBNAME="sedum_behave_grocery_$$"

target_scaffold() {
  mkdir -p "$APP"
  cd "$APP" || return 1
  # `go mod init` is the whole of a Go scaffold. Everything else this service
  # has - main.go, its packages, its schema - comes from Sedum's file templates,
  # which is the opposite balance from Rails and is the point of running both.
  go mod init "$MODULE"
}

target_prepare() {
  local scope
  scope=$(ruby -ryaml -e '
    Dir[File.join(ARGV[0], "*.yml"), File.join(ARGV[0], "*.yaml")].each do |f|
      puts (YAML.safe_load_file(f)["affected_scope"] || [])
    end' "$SEDUM_REPO/$RECORDS")

  for p in $scope; do
    if [ -e "$APP/$p" ]; then
      echo "removing scaffolded $p so Sedum's file template owns it"
      rm -f "$APP/$p"
    fi
  done
  echo "module $MODULE, database $DBNAME"
}

target_build() {
  cd "$APP" || return 1
  # The dependency handoff, in Go's spelling. The generator package writes the
  # imports; resolving them to modules is the thing a person or a tool does
  # afterwards, the same role bundle install plays on the Rails side.
  go mod tidy || return 1
  go build -o "$WORK/server" . || return 1
  go vet ./...
}

target_boot() {
  cd "$APP" || return 1

  createdb "$DBNAME" || return 1
  psql -q -d "$PGURL/$DBNAME" -f "$APP/init.sql" || return 1

  APP_PORT=8080   # the file template hardcodes the listen port
  BASE_URL="http://127.0.0.1:$APP_PORT"
  DATABASE_URL="$PGURL/$DBNAME?sslmode=disable" "$WORK/server" > "$LOGS/server.log" 2>&1 &
  APP_PID=$!

  # The base URL and no path: a boot gate waits for a listener, and any HTTP
  # response proves one (prov-2026-79e94e7c).
  wait_for "$BASE_URL/" 30 || { tail -40 "$LOGS/server.log"; return 1; }
  echo "up on $BASE_URL"
}

# The same contract the Rails arm asserts, because the two records ask for the
# same functionality constraint for constraint. Where the two arms disagree, the
# generator packages disagree - which is the comparison worth having.
target_verify() {
  local code id

  code=$(status POST /groceries '{"name":"oat milk","quantity":2}')
  check "create returns 201" "$code" "201"
  id=$(body '.id')
  check "create body carries a name"     "$(body '.name')"      "oat milk"
  check "create body carries a quantity" "$(body '.quantity')"  "2"
  check "create body carries created_at" "$(body 'has("created_at")')" "true"
  check "create body carries updated_at" "$(body 'has("updated_at")')" "true"

  code=$(status GET /groceries)
  check "list returns 200" "$code" "200"
  check "list is an array" "$(body 'type')" "array"

  code=$(status GET "/groceries/$id")
  check "get returns 200" "$code" "200"

  code=$(status GET /groceries/999999)
  check "get of a missing record returns 404" "$code" "404"

  code=$(status PUT "/groceries/$id" '{"name":"oat milk","quantity":5}')
  check "update returns 200" "$code" "200"
  check "update sets quantity" "$(body '.quantity')" "5"

  code=$(status PUT /groceries/999999 '{"quantity":9}')
  check "update of a missing record returns 404" "$code" "404"

  code=$(status DELETE "/groceries/$id")
  check "delete returns 204" "$code" "204"

  code=$(status GET "/groceries/$id")
  check "a deleted record is gone" "$code" "404"

  code=$(status DELETE /groceries/999999)
  check "delete of a missing record returns 404" "$code" "404"

  ( cd "$APP" && go test ./... 2>&1 ) > "$LOGS/gotest.log"
  local rc=$?
  check "the generated go test suite passes" "$rc" "0"

  return 0
}

target_teardown() {
  [ -n "${APP_PID:-}" ] && kill "$APP_PID" 2>/dev/null
  sleep 1
  dropdb --if-exists "$DBNAME" 2>/dev/null
}
