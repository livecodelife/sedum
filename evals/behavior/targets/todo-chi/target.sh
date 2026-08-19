# The Go arm of the behaviour harness.
#
# Written second, and deliberately against a stack with nothing in common with
# the Rails one: a different scaffold command, a compiler in place of a bundler,
# a schema file in place of a migration, four generator packages instead of one.
# If the runner needed changing to accommodate it, the runner was not generic.

GENERATORS="evals/testdata/todo-api/generators/defined"
RECORDS="evals/testdata/todo-api/records"

GENERATED_PATHS=(
  main.go
  init.sql
  db/todos.go
  db/todos_test.go
  handlers/todos.go
  handlers/todos_test.go
)

# The case is the single source. RunBehavior exports the case's variables as
# SEDUM_VAR_<NAME>, and both go mod init and --var read this one value, so the
# scaffold and the generation cannot disagree about the module's name
# (prov-2026-1c33a50b).
MODULE="${SEDUM_VAR_MODULE:-${MODULE:-todo}}"
# The module path is a project fact the standard cannot know, so the run supplies
# it rather than the model guessing it (prov-2026-6fc3d13d). This is the same
# value `go mod init` was given, which is the point: one fact, named once.
SEDUM_VARS="--var module=$MODULE"
PGURL="${PGURL:-postgres://$(whoami)@localhost:5432}"
DBNAME="sedum_behave_chi_$$"

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

  # The harness picks the port and tells the service, rather than knowing
  # what the template chose. main.go is in the record's affected_scope, so
  # the baseline arm writes it too, and a hardcoded port here would fail an
  # arm that listened anywhere else (prov-2026-682a079b).
  APP_PORT=$(free_port)
  BASE_URL="http://127.0.0.1:$APP_PORT"
  PORT="$APP_PORT" DATABASE_URL="$PGURL/$DBNAME?sslmode=disable" "$WORK/server" > "$LOGS/server.log" 2>&1 &
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

  code=$(status POST /todos '{"title":"write the harness"}')
  check "create returns 201" "$code" "201"
  id=$(body '.id')
  check "create body carries a title"    "$(body '.title')"     "write the harness"
  check "create body defaults completed" "$(body '.completed')" "false"
  check "create body carries created_at" "$(body 'has("created_at")')" "true"
  check "create body carries updated_at" "$(body 'has("updated_at")')" "true"

  code=$(status GET /todos)
  check "list returns 200" "$code" "200"
  check "list is an array" "$(body 'type')" "array"

  code=$(status GET "/todos/$id")
  check "get returns 200" "$code" "200"

  code=$(status GET /todos/999999)
  check "get of a missing record returns 404" "$code" "404"

  code=$(status PUT "/todos/$id" '{"completed":true, "title":"write the harness"}')
  check "update returns 200" "$code" "200"
  check "update sets completed" "$(body '.completed')" "true"

  code=$(status PUT /todos/999999 '{"completed":true}')
  check "update of a missing record returns 404" "$code" "404"

  code=$(status DELETE "/todos/$id")
  check "delete returns 204" "$code" "204"

  code=$(status GET "/todos/$id")
  check "a deleted record is gone" "$code" "404"

  code=$(status DELETE /todos/999999)
  check "delete of a missing record returns 404" "$code" "404"

  # The schema init.sql actually produced.
  #
  # Both records state the column rules in the same words, and only the Rails
  # target checked them - so the two arms of the language comparison were
  # reading their records to different depths. No HTTP call can observe a
  # default on a NOT NULL column the caller always supplies, which is the
  # blindness the Rails side added these for after a green wall lied once
  # (prov-2026-71940de0).
  local coldefault
  coldefault=$(psql -X -A -t -d "$PGURL/$DBNAME" -c \
    "SELECT coalesce(column_default, 'NULL') FROM information_schema.columns
      WHERE table_name = 'todos' AND column_name = 'title'")
  check "title carries no column default" "$coldefault" "NULL"

  coldefault=$(psql -X -A -t -d "$PGURL/$DBNAME" -c \
    "SELECT coalesce(column_default, 'NULL') FROM information_schema.columns
      WHERE table_name = 'todos' AND column_name = 'completed'")
  check "completed defaults to false" "$coldefault" "false"

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
