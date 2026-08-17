# The Rails arm of the behaviour harness.
#
# Everything framework-specific lives here: how the empty project is made, how
# it is built, how it is started, and what its API is supposed to do. The runner
# knows none of it, which is what a chi or a react target would replace.

GENERATORS="evals/testdata/todo-rails/generators/defined"
RECORDS="evals/testdata/todo-rails/records"

# Copied out beside the results when the run ends, so a failed assertion can be
# read against the code that produced it after the project is gone.
GENERATED_PATHS=(
  app/controllers/todos_controller.rb
  app/models/todo.rb
  config/routes.rb
  db/migrate/20260814000000_create_todos.rb
  test/controllers/todos_controller_test.rb
  test/models/todo_test.rb
)

RAILS_VERSION="${RAILS_VERSION:-7.2.3}"
PGURL="${PGURL:-postgres://$(whoami)@localhost:5432}"
DBNAME="sedum_behave_$$"

target_scaffold() {
  mkdir -p "$(dirname "$APP")"
  # Not --api. The controller template skips the CSRF filter unconditionally,
  # and ActionController::API has no such filter to skip, so an API-only
  # scaffold raises at class definition time. The standard these generators
  # describe is a full Rails application serving JSON.
  rails "_${RAILS_VERSION}_" new "$APP" \
    --database=postgresql --skip-git --skip-action-mailbox --skip-action-text \
    --skip-action-cable --skip-jbuilder --skip-system-test
}

target_prepare() {
  # The paths the record authorizes, read from the record rather than listed
  # again here - the record is what decides them.
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

  # The record does not authorize database.yml, so the scaffold's is what runs.
  # Rails merges DATABASE_URL over it, which is how the harness points the
  # application at a database it can drop afterwards.
  echo "database: $DBNAME"
}

target_build() {
  cd "$APP" || return 1

  # rails new already ran bundle install, and nothing Sedum wrote adds a
  # dependency: the generator package declares the Gemfile unmanaged and reports
  # it as a handoff, so a build step that had to install something on Sedum's
  # behalf would be a finding rather than a step.
  #
  # This one is not on Sedum's behalf. A fresh 7.2 scaffold resolves minitest 6,
  # whose runner signature railties 7.2 does not match, so `rails test` dies
  # before running a test. The harness is the person doing the Gemfile handoff.
  if ! grep -q '^gem "minitest"' Gemfile; then
    printf '\ngem "minitest", "~> 5.25"\n' >> Gemfile
    bundle install --quiet || return 1
  fi

  bundle check
}

target_boot() {
  cd "$APP" || return 1
  export DATABASE_URL="$PGURL/$DBNAME"
  export RAILS_ENV=development

  bin/rails db:create || return 1
  bin/rails db:migrate || return 1

  APP_PORT=$(free_port)
  BASE_URL="http://127.0.0.1:$APP_PORT"
  bin/rails server -b 127.0.0.1 -p "$APP_PORT" > "$LOGS/server.log" 2>&1 &
  APP_PID=$!

  wait_for "$BASE_URL/up" 60 || { tail -40 "$LOGS/server.log"; return 1; }
  echo "up on $BASE_URL"
}

# The record's constraints, one assertion each. These are written from the
# record, not from what the generated code happens to do - a contract read off
# the implementation is not a contract.
target_verify() {
  local code id title

  # Five endpoints, and the JSON shape the record names.
  code=$(status POST /todos '{"title":"write the harness"}')
  check "create returns 201" "$code" "201"
  id=$(body '.id')
  check "create body carries a title"      "$(body '.title')"      "write the harness"
  check "create body defaults completed"   "$(body '.completed')"  "false"
  check "create body carries created_at"   "$(body 'has("created_at")')" "true"
  check "create body carries updated_at"   "$(body 'has("updated_at")')" "true"

  code=$(status GET /todos)
  check "index returns 200" "$code" "200"
  check "index is an array"  "$(body 'type')" "array"

  code=$(status GET "/todos/$id")
  check "show returns 200" "$code" "200"

  code=$(status GET /todos/999999)
  check "show of a missing record returns 404" "$code" "404"

  code=$(status PUT "/todos/$id" '{"completed":true}')
  check "update returns 200" "$code" "200"
  check "update sets completed" "$(body '.completed')" "true"

  # Only title and completed are accepted from a body.
  code=$(status PUT "/todos/$id" '{"title":"renamed","nonsense":"ignored"}')
  check "an unpermitted attribute is ignored, not assigned" "$code" "200"
  check "the permitted attribute still applied" "$(body '.title')" "renamed"

  code=$(status PUT /todos/999999 '{"completed":true}')
  check "update of a missing record returns 404" "$code" "404"

  code=$(status DELETE "/todos/$id")
  check "destroy returns 204" "$code" "204"

  code=$(status GET "/todos/$id")
  check "a destroyed record is gone" "$code" "404"

  code=$(status DELETE /todos/999999)
  check "destroy of a missing record returns 404" "$code" "404"

  # The schema the migration actually produced.
  #
  # Added because a green wall lied once. Every assertion above passed against a
  # migration that declared `default: ""` on a column the record says carries no
  # default - no HTTP call can observe a default on a NOT NULL column the caller
  # always supplies, so the contract was blind to it by construction.
  #
  # Narrowly a statement about the running database, not a second copy of the
  # eval's binding expectations. The eval scores what the model chose; this
  # scores what the database ended up with.
  local coldefault
  coldefault=$(psql -X -A -t -d "$PGURL/$DBNAME" -c \
    "SELECT coalesce(column_default, 'NULL') FROM information_schema.columns
      WHERE table_name = 'todos' AND column_name = 'title'")
  check "title carries no column default" "$coldefault" "NULL"

  coldefault=$(psql -X -A -t -d "$PGURL/$DBNAME" -c \
    "SELECT coalesce(column_default, 'NULL') FROM information_schema.columns
      WHERE table_name = 'todos' AND column_name = 'completed'")
  check "completed defaults to false" "$coldefault" "false"

  # The tests the record required each endpoint to carry. Whether they pass is
  # a second, independent reading of the same code.
  ( cd "$APP" && DATABASE_URL="$PGURL/${DBNAME}_test" RAILS_ENV=test \
      bin/rails db:create db:schema:load test 2>&1 ) > "$LOGS/minitest.log"
  local rc=$?
  check "the generated minitest suite passes" "$rc" "0"
  grep -E '^[0-9]+ runs' "$LOGS/minitest.log" | tail -1 >&3

  # A failed assertion is a result, not a broken phase. Returning the last
  # check's status would report the two as the same thing.
  return 0
}

target_teardown() {
  [ -n "${APP_PID:-}" ] && kill "$APP_PID" 2>/dev/null
  sleep 1
  dropdb --if-exists "$DBNAME" 2>/dev/null
  dropdb --if-exists "${DBNAME}_test" 2>/dev/null
}
