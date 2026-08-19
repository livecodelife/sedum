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
  # --api, because that is what an engineer building a JSON API in Rails runs.
  #
  # It used to be a full application, and the reason was circular: the
  # controller template skipped the CSRF filter, which ActionController::API
  # does not have, so an API-only scaffold raised at class definition time. The
  # template skipped it because the scaffold was not --api. The first baseline
  # run is what broke the tie - it scored 0 of 5 with the server log naming
  # InvalidAuthenticityToken fifteen times, because the arm without the package
  # wrote the controller an engineer would write (prov-2026-f5e64f22).
  rails "_${RAILS_VERSION}_" new "$APP" --api \
    --database=postgresql --skip-git --skip-action-mailbox --skip-action-text \
    --skip-action-cable --skip-jbuilder --skip-system-test
}

# plant_routes_anchor inserts the marker Sedum's routes template declares into
# the routes file `rails new` wrote, immediately before the closing `end` of the
# draw block.
#
# The indentation matches the template's, because injection takes its
# indentation from the marker it lands on. Planting twice would leave two
# regions for one anchor to choose between, so it is idempotent.
plant_routes_anchor() {
  ruby -e '
    path = ARGV[0]
    lines = File.readlines(path)
    if lines.any? { |l| l.include?("sedum:anchor:routes") }
      puts "  config/routes.rb already carries its anchor"
      exit 0
    end
    idx = lines.rindex { |l| l.rstrip == "end" }
    abort "config/routes.rb has no closing end to plant before" if idx.nil?
    lines.insert(idx, "\n", "  # sedum:anchor:routes\n")
    File.write(path, lines.join)
    puts "  planted sedum:anchor:routes in the scaffolded config/routes.rb"
  ' "$1"
}

target_prepare() {
  # The paths the record authorizes, read from the record rather than listed
  # again here - the record is what decides them.
  local scope
  scope=$(ruby -ryaml -e '
    Dir[File.join(ARGV[0], "*.yml"), File.join(ARGV[0], "*.yaml")].each do |f|
      puts (YAML.safe_load_file(f)["affected_scope"] || [])
    end' "$SEDUM_REPO/$RECORDS")

  # An authorized path the scaffold already wrote is adopted, not deleted.
  #
  # Phase 3 is create-if-absent: a file already on disk is left alone and
  # checked for the markers its template declares. So an existing file is not a
  # problem for Sedum - an existing file *without the anchors* is, and the fix
  # is one comment line rather than handing the whole file over. A routes file
  # accumulates hand-written routes over the life of a service, and a generator
  # that owned it outright would be hostile to how the file is actually used
  # (prov-2026-ee1b3e46).
  for p in $scope; do
    [ -e "$APP/$p" ] || continue
    case "$p" in
      config/routes.rb) plant_routes_anchor "$APP/$p" || return 1 ;;
      *)
        # The harness cannot know where another file's anchors belong. That is
        # what `adopt` is for (OPEN_QUESTIONS section 3) and it is unbuilt, so
        # this refuses rather than guesses.
        echo "authorized path $p already exists and the harness cannot place its anchors" >&2
        return 1
        ;;
    esac
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
