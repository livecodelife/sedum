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

MODULE="${MODULE:-todo}"
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

  APP_PORT=8080   # the file template hardcodes the listen port
  BASE_URL="http://127.0.0.1:$APP_PORT"
  DATABASE_URL="$PGURL/$DBNAME?sslmode=disable" "$WORK/server" > "$LOGS/server.log" 2>&1 &
  APP_PID=$!

  wait_for "$BASE_URL/todos" 30 || { tail -40 "$LOGS/server.log"; return 1; }
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

  code=$(status PUT "/todos/$id" '{"completed":true}')
  check "partial update returns 200" "$code" "200"
  check "partial update sets completed" "$(body '.completed')" "true"
  check "partial update keeps the title it never sent" "$(body '.title')" "write the harness"

  code=$(status PUT "/todos/$id" '{}')
  check "an empty update is a 400" "$code" "400"

  code=$(status POST /todos '{"completed":true}')
  check "a todo with no title is rejected" "$code" "400"

  code=$(status PUT /todos/999999 '{"completed":true}')
  check "update of a missing record returns 404" "$code" "404"

  code=$(status DELETE "/todos/$id")
  check "delete returns 204" "$code" "204"

  code=$(status GET "/todos/$id")
  check "a deleted record is gone" "$code" "404"

  code=$(status DELETE /todos/999999)
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
