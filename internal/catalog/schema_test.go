package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/recording"
)

// prov-2026-91c54941: --response-schema compiles what survives Phase 4's
// filter into a JSON Schema for {"invocations": [...]}, so an endpoint that
// supports structured output can be asked to constrain the model to it.
//
// This file holds Catalog.ResponseSchema() - the compiler - to the shape the
// record specifies. It does not exist yet; every test below is written to
// fail against today's catalog package (a compile failure, since the method
// is new) until code-author adds it.
//
// A schema that rejects a valid binding is the worst possible bug this record
// can ship (a false rejection burns a retry, or worse, makes a legal answer
// unreachable through the one channel meant to narrow the model toward it),
// so the fixtures below are deliberately modelled on the two real package
// shapes this feature was measured against, and the "known-good recording
// must validate" test carries a real accepted sedum grow recording rather
// than a synthetic one.
//
// No third-party JSON Schema library is used. This repo's go.mod is
// read-only and governed entirely by other, already-implemented records
// (none of them this one), so adding a dependency is out of this record's
// affected_scope; validateAgainstSchema below implements just the keywords
// ResponseSchema is specified to emit (type, const, properties, required,
// additionalProperties, items, oneOf) against real JSON Schema semantics.

// checkRuleStyleGenerators is modelled on the check-rule real package shape:
// a rule action whose target is per-record (a file named by an id) and a
// catalog-entry action whose target every record shares, plus - what this
// record's test requirements ask this fixture to add beyond
// checkRuleGenerators in filter_test.go - a kwarg with a prose description
// and a list kwarg, so the schema compiler's type mapping and its handling of
// Description (which must not leak into the schema; only the wire type
// matters to a validator) both have something to prove.
func checkRuleStyleGenerators() map[string]string {
	return map[string]string{
		"checkrule/sedum.yaml": `name: checkrule
extensions: [".ts"]
comment_prefix: "//"
`,
		"checkrule/files/checks/{id}.ts": "// sedum:anchor:body\n",
		"checkrule/files/catalog.ts":     "// sedum:anchor:entries\n",
		"checkrule/actions/actions.yaml": `actions:
  addRule:
    kwargs:
      id: { type: string, required: true }
      body: { type: string, required: true }
      severity:
        type: string
        required: false
        description: >-
          Prose. One of info, warning, or error, spelled the way the target
          linter's own configuration spells its severities - not a number and
          not this schema's own idea of an enum key.
      tags:
        type: list
        required: false
    injects_into: "checks/{{id}}.ts"
    anchor: body

  addCatalogEntry:
    kwargs:
      id: { type: string, required: true }
    injects_into: "catalog.ts"
    anchor: entries
`,
		"checkrule/actions/addRule.ts":         "// {{body}}\n",
		"checkrule/actions/addCatalogEntry.ts": "// entry {{id}}\n",
	}
}

// tableStyleWithDiscriminatorGenerators is modelled on the other real shape
// measured, extended with what this record's test requirements ask for that
// tableStyleGenerators in filter_test.go does not carry: three actions across
// two packages, a discriminator with two variants and variant_requires, and a
// pipeline placeholder in injects_into (addField's {{table|pluralTable}}).
//
// addNotify's two variant templates are what gives it a variant_requires:
// created.nts renders {{extra}}, which addNotify's own schema calls optional,
// so createControllerMethod's mechanism (prov-2026-369544c1) is exercised
// here the same way TestVariantRequirementsReachBothConsumers exercises it in
// catalog_test.go. deleted.nts does not, so "deleted" carries no extra
// requirement beyond the action's own declared-required kwargs.
func tableStyleWithDiscriminatorGenerators() map[string]string {
	return map[string]string{
		"tables/sedum.yaml": `name: tables
extensions: [".ts", ".sql"]
comment_prefix: "--"
transforms:
  pluralTable: [plural]
`,
		"tables/files/lib/identity/{table}.ts":        "// sedum:anchor:fields\n",
		"tables/files/supabase/migrations/{name}.sql": "-- sedum:anchor:columns\n",
		"tables/actions/actions.yaml": `actions:
  addField:
    kwargs:
      table: { type: string, required: true }
      field: { type: string, required: true }
    injects_into: "lib/identity/{{table|pluralTable}}.ts"
    anchor: fields

  addColumn:
    kwargs:
      migrationName: { type: string, required: true }
      column: { type: string, required: true }
    injects_into: "supabase/migrations/{{migrationName}}.sql"
    anchor: columns
`,
		"tables/actions/addField.ts":  "// {{field}}\n",
		"tables/actions/addColumn.ts": "-- {{column}}\n",

		"notify/sedum.yaml": `name: notify
extensions: [".nts"]
comment_prefix: "//"
`,
		"notify/files/lib/notifiers/{resource}.nts": "// sedum:anchor:notifications\n",
		"notify/actions/actions.yaml": `actions:
  addNotify:
    kwargs:
      resource: { type: string, required: true }
      event: { type: string, required: true }
      extra: { type: string, required: false }
    discriminator: event
    variants: [created, deleted]
    injects_into: "lib/notifiers/{{resource}}.nts"
    anchor: notifications
`,
		"notify/actions/addNotify/created.nts": "// created {{resource}} extra={{extra}}\n",
		"notify/actions/addNotify/deleted.nts": "// deleted {{resource}}\n",
	}
}

// mustResponseSchema builds and decodes the schema, failing the test with
// context rather than a bare panic if the compiler cannot produce one.
func mustResponseSchema(t *testing.T, c Catalog) (raw []byte, decoded any) {
	t.Helper()
	raw, err := c.ResponseSchema()
	if err != nil {
		t.Fatalf("ResponseSchema: %v", err)
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("ResponseSchema produced invalid JSON: %v\n%s", err, raw)
	}
	return raw, decoded
}

func mustDecode(t *testing.T, s string) any {
	t.Helper()
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		t.Fatalf("fixture does not parse as JSON: %v\n%s", err, s)
	}
	return v
}

// The check-rule fixture's exact expected schema, authored by hand from the
// record rather than derived from any implementation: a plain oneOf branch
// per action (neither is discriminated), additionalProperties: false
// throughout, required folding entry.Requires into the declared-required
// kwargs, list mapped to array-of-string, and no description anywhere - the
// wire contract is the type, not the prose.
const checkRuleExpectedSchema = `{
  "type": "object",
  "additionalProperties": false,
  "required": ["invocations"],
  "properties": {
    "invocations": {
      "type": "array",
      "items": {
        "oneOf": [
          {
            "type": "object",
            "additionalProperties": false,
            "required": ["action", "kwargs"],
            "properties": {
              "action": {"const": "addCatalogEntry"},
              "kwargs": {
                "type": "object",
                "additionalProperties": false,
                "required": ["id"],
                "properties": {
                  "id": {"type": "string"}
                }
              }
            }
          },
          {
            "type": "object",
            "additionalProperties": false,
            "required": ["action", "kwargs"],
            "properties": {
              "action": {"const": "addRule"},
              "kwargs": {
                "type": "object",
                "additionalProperties": false,
                "required": ["body", "id"],
                "properties": {
                  "id": {"type": "string"},
                  "body": {"type": "string"},
                  "severity": {"type": "string"},
                  "tags": {"type": "array", "items": {"type": "string"}}
                }
              }
            }
          }
        ]
      }
    }
  }
}`

func TestResponseSchemaSnapshotCheckRuleStyle(t *testing.T) {
	packages := loadPackages(t, checkRuleStyleGenerators(), "checkrule")
	cat := Build(packages, Options{})

	_, got := mustResponseSchema(t, cat)
	want := mustDecode(t, checkRuleExpectedSchema)

	if !reflect.DeepEqual(got, want) {
		gotJSON, _ := json.MarshalIndent(got, "", "  ")
		t.Errorf("ResponseSchema for the check-rule fixture does not match the hand-authored schema:\ngot:\n%s\n\nwant:\n%s",
			gotJSON, checkRuleExpectedSchema)
	}
}

// The table-style fixture's exact expected schema: three actions across two
// packages, addNotify expanding into one oneOf branch per declared variant
// with the discriminator pinned by const (so oneOf stays exclusive - without
// the const, an instance naming "created" could also satisfy the "deleted"
// branch, since additionalProperties: false does not forbid a property the
// schema still declares, only one the model bound and no branch declares at
// all) and that branch's variant_requires folded into required alongside the
// action's own declared-required kwargs.
const tableStyleExpectedSchema = `{
  "type": "object",
  "additionalProperties": false,
  "required": ["invocations"],
  "properties": {
    "invocations": {
      "type": "array",
      "items": {
        "oneOf": [
          {
            "type": "object",
            "additionalProperties": false,
            "required": ["action", "kwargs"],
            "properties": {
              "action": {"const": "addColumn"},
              "kwargs": {
                "type": "object",
                "additionalProperties": false,
                "required": ["column", "migrationName"],
                "properties": {
                  "migrationName": {"type": "string"},
                  "column": {"type": "string"}
                }
              }
            }
          },
          {
            "type": "object",
            "additionalProperties": false,
            "required": ["action", "kwargs"],
            "properties": {
              "action": {"const": "addField"},
              "kwargs": {
                "type": "object",
                "additionalProperties": false,
                "required": ["field", "table"],
                "properties": {
                  "table": {"type": "string"},
                  "field": {"type": "string"}
                }
              }
            }
          },
          {
            "type": "object",
            "additionalProperties": false,
            "required": ["action", "kwargs"],
            "properties": {
              "action": {"const": "addNotify"},
              "kwargs": {
                "type": "object",
                "additionalProperties": false,
                "required": ["event", "extra", "resource"],
                "properties": {
                  "resource": {"type": "string"},
                  "event": {"const": "created"},
                  "extra": {"type": "string"}
                }
              }
            }
          },
          {
            "type": "object",
            "additionalProperties": false,
            "required": ["action", "kwargs"],
            "properties": {
              "action": {"const": "addNotify"},
              "kwargs": {
                "type": "object",
                "additionalProperties": false,
                "required": ["event", "resource"],
                "properties": {
                  "resource": {"type": "string"},
                  "event": {"const": "deleted"},
                  "extra": {"type": "string"}
                }
              }
            }
          }
        ]
      }
    }
  }
}`

func TestResponseSchemaSnapshotTableStyle(t *testing.T) {
	packages := loadPackages(t, tableStyleWithDiscriminatorGenerators(), "tables", "notify")
	cat := Build(packages, Options{})

	_, got := mustResponseSchema(t, cat)
	want := mustDecode(t, tableStyleExpectedSchema)

	if !reflect.DeepEqual(got, want) {
		gotJSON, _ := json.MarshalIndent(got, "", "  ")
		t.Errorf("ResponseSchema for the table-style fixture does not match the hand-authored schema:\ngot:\n%s\n\nwant:\n%s",
			gotJSON, tableStyleExpectedSchema)
	}
}

// Generation is deterministic: the same catalog produces byte-identical
// schema bytes on every run, which is what makes "the schema, or its hash and
// the path it was written to" (the run.log constraint) a meaningful thing to
// log.
func TestResponseSchemaIsDeterministic(t *testing.T) {
	packages := loadPackages(t, tableStyleWithDiscriminatorGenerators(), "tables", "notify")
	cat := Build(packages, Options{})

	a, err := cat.ResponseSchema()
	if err != nil {
		t.Fatalf("ResponseSchema (1st): %v", err)
	}
	b, err := cat.ResponseSchema()
	if err != nil {
		t.Fatalf("ResponseSchema (2nd): %v", err)
	}
	if string(a) != string(b) {
		t.Errorf("ResponseSchema is not deterministic:\n1st:\n%s\n\n2nd:\n%s", a, b)
	}
}

// prov-2026-4bcabb2f found that a schema forcing at least one invocation
// collapsed selection. This record's standing constraint is that
// {"invocations": []} stays legal under --response-schema exactly as without
// it, so the array carries no minItems and an empty list must validate.
func TestResponseSchemaNeverForcesNonEmptyInvocations(t *testing.T) {
	packages := loadPackages(t, tableStyleWithDiscriminatorGenerators(), "tables", "notify")
	cat := Build(packages, Options{})

	raw, schema := mustResponseSchema(t, cat)

	if containsKey(schema, "minItems") {
		t.Error("the schema declares minItems; prov-2026-4bcabb2f already measured this exact family of constraint collapsing selection")
	}

	if err := validateAgainstSchema(schema, mustDecode(t, `{"invocations": []}`)); err != nil {
		t.Errorf("{\"invocations\": []} must validate, and does not: %v\n%s", err, raw)
	}
}

// Optional kwargs stay optional properties and are never modeled as
// nullable - Phase 5 rejects null for an absent optional kwarg
// (prov-2026-9a554c93/prov-2026-9a491128), so a schema legalizing null would
// legalize an output shape the validator throws out anyway, at the cost of a
// wasted retry instead of a clean rejection at the decoding boundary.
func TestResponseSchemaNeverModelsOptionalKwargsAsNullable(t *testing.T) {
	packages := loadPackages(t, checkRuleStyleGenerators(), "checkrule")
	cat := Build(packages, Options{})
	_, schema := mustResponseSchema(t, cat)

	// severity and tags are addRule's optional kwargs. Binding either to null
	// must fail exactly as binding a wrong type would.
	for _, instance := range []string{
		`{"invocations":[{"action":"addRule","kwargs":{"id":"x","body":"y","severity":null}}]}`,
		`{"invocations":[{"action":"addRule","kwargs":{"id":"x","body":"y","tags":null}}]}`,
	} {
		if err := validateAgainstSchema(schema, mustDecode(t, instance)); err == nil {
			t.Errorf("null bound to an optional kwarg validated, want rejected: %s", instance)
		}
	}
}

// The catalog is the single source both Phase 5's derivedRequirements and the
// schema compiler must read required kwargs from (no second, independently
// maintained copy of the rule). This does not know the compiler's internals;
// it recomputes what a correct compiler's "required" set has to equal,
// directly from catalog.Action's own Requires/VariantRequires/Kwargs fields -
// the same fields derivedRequirements in internal/selection/validate.go
// reads - and holds every branch of both fixtures to it.
func TestResponseSchemaRequiredMatchesCatalogsOwnDerivation(t *testing.T) {
	check := Build(loadPackages(t, checkRuleStyleGenerators(), "checkrule"), Options{})
	table := Build(loadPackages(t, tableStyleWithDiscriminatorGenerators(), "tables", "notify"), Options{})

	for _, cat := range []Catalog{check, table} {
		_, schema := mustResponseSchema(t, cat)
		branches := oneOfBranches(t, schema)

		for _, action := range cat.Actions {
			variants := action.Variants
			if len(variants) == 0 {
				variants = []string{""}
			}
			for _, variant := range variants {
				want := expectedRequired(action, variant)
				got, ok := requiredForBranch(branches, action.Name, action.Discriminator, variant)
				if !ok {
					t.Errorf("no oneOf branch found for action %q variant %q", action.Name, variant)
					continue
				}
				if !sortedEqual(got, want) {
					t.Errorf("action %q variant %q: schema requires %v, want %v (declared-required kwargs union Requires union VariantRequires[variant])",
						action.Name, variant, got, want)
				}
			}
		}
	}
}

// expectedRequired computes, straight from catalog.Action, the set
// derivedRequirements (internal/selection/validate.go) already treats as
// required for a given variant - the declared-required kwargs, plus
// entry.Requires, plus entry.VariantRequires[variant] when the action is
// discriminated.
func expectedRequired(action Action, variant string) []string {
	set := map[string]bool{}
	for name, k := range action.Kwargs {
		if k.Required {
			set[name] = true
		}
	}
	for _, name := range action.Requires {
		set[name] = true
	}
	if variant != "" {
		for _, name := range action.VariantRequires[variant] {
			set[name] = true
		}
	}
	out := make([]string, 0, len(set))
	for name := range set {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func sortedEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// oneOfBranches walks the compiled schema down to invocations.items.oneOf.
func oneOfBranches(t *testing.T, schema any) []any {
	t.Helper()
	m, ok := schema.(map[string]any)
	if !ok {
		t.Fatalf("schema root is not an object: %#v", schema)
	}
	props, _ := m["properties"].(map[string]any)
	invocations, _ := props["invocations"].(map[string]any)
	items, _ := invocations["items"].(map[string]any)
	branches, ok := items["oneOf"].([]any)
	if !ok {
		t.Fatalf("schema does not carry properties.invocations.items.oneOf: %#v", schema)
	}
	return branches
}

// requiredForBranch finds the branch naming action (and, for a discriminated
// action, pinning discriminator to variant via const) and returns its
// kwargs.required.
func requiredForBranch(branches []any, action, discriminator, variant string) ([]string, bool) {
	for _, b := range branches {
		branch, _ := b.(map[string]any)
		props, _ := branch["properties"].(map[string]any)
		actionSchema, _ := props["action"].(map[string]any)
		if actionSchema["const"] != action {
			continue
		}
		if discriminator != "" {
			kwargs, _ := props["kwargs"].(map[string]any)
			kwargProps, _ := kwargs["properties"].(map[string]any)
			discSchema, _ := kwargProps[discriminator].(map[string]any)
			if discSchema["const"] != variant {
				continue
			}
		}
		kwargs, _ := props["kwargs"].(map[string]any)
		reqRaw, _ := kwargs["required"].([]any)
		out := make([]string, 0, len(reqRaw))
		for _, r := range reqRaw {
			out = append(out, r.(string))
		}
		return out, true
	}
	return nil, false
}

// containsKey reports whether key appears anywhere in a decoded JSON value,
// at any depth - used to assert the compiler never reaches for minItems
// (prov-2026-4bcabb2f) regardless of where a future change might be tempted
// to add it.
func containsKey(v any, key string) bool {
	switch t := v.(type) {
	case map[string]any:
		if _, ok := t[key]; ok {
			return true
		}
		for _, sub := range t {
			if containsKey(sub, key) {
				return true
			}
		}
	case []any:
		for _, sub := range t {
			if containsKey(sub, key) {
				return true
			}
		}
	}
	return false
}

// Every case the schema itself must reject, per this record's test
// requirements: an unknown action, an action the Phase 4 filter removed, an
// unknown kwarg, a missing required kwarg, a missing variant-required kwarg,
// a wrong type, and an undeclared variant value.
func TestResponseSchemaRejectsWhatItMust(t *testing.T) {
	checkPackages := loadPackages(t, checkRuleStyleGenerators(), "checkrule")
	checkCat := Build(checkPackages, Options{})
	_, checkSchema := mustResponseSchema(t, checkCat)

	tablePackages := loadPackages(t, tableStyleWithDiscriminatorGenerators(), "tables", "notify")
	tableCat := Build(tablePackages, Options{})
	_, tableSchema := mustResponseSchema(t, tableCat)

	// An action the Phase 4 filter already removed: only checks/x.ts is
	// authorized, so addCatalogEntry (target catalog.ts) is filtered out
	// before the schema ever sees it - the same case
	// TestFilterDropsAnActionWhoseTargetNoAuthorizedPathCanMatch exercises for
	// the prompt.
	filtered, removed := FilterForRecord(checkCat, []string{"checks/x.ts"})
	if len(removed) != 1 || removed[0].Action != "addCatalogEntry" {
		t.Fatalf("fixture assumption broken: removed = %+v, want exactly addCatalogEntry", removed)
	}
	_, filteredSchema := mustResponseSchema(t, filtered)

	cases := []struct {
		name     string
		schema   any
		instance string
	}{
		{
			"unknown action",
			checkSchema,
			`{"invocations":[{"action":"noSuchAction","kwargs":{}}]}`,
		},
		{
			"action the filter removed",
			filteredSchema,
			`{"invocations":[{"action":"addCatalogEntry","kwargs":{"id":"x"}}]}`,
		},
		{
			"unknown kwarg",
			checkSchema,
			`{"invocations":[{"action":"addRule","kwargs":{"id":"x","body":"y","notDeclared":"z"}}]}`,
		},
		{
			"missing required kwarg",
			checkSchema,
			`{"invocations":[{"action":"addRule","kwargs":{"id":"x"}}]}`,
		},
		{
			"missing variant-required kwarg",
			tableSchema,
			`{"invocations":[{"action":"addNotify","kwargs":{"resource":"widgets","event":"created"}}]}`,
		},
		{
			"wrong type: a string where a list is declared",
			checkSchema,
			`{"invocations":[{"action":"addRule","kwargs":{"id":"x","body":"y","tags":"not-a-list"}}]}`,
		},
		{
			"undeclared variant value",
			tableSchema,
			`{"invocations":[{"action":"addNotify","kwargs":{"resource":"widgets","event":"archived"}}]}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := validateAgainstSchema(c.schema, mustDecode(t, c.instance)); err == nil {
				t.Errorf("%s: instance validated, want rejected:\n%s", c.name, c.instance)
			}
		})
	}
}

// Every branch either fixture legitimately offers must accept a correctly
// bound invocation, including both of addNotify's variant branches - a
// schema this strict is only trustworthy if it is not so strict it also
// rejects what it should accept.
func TestResponseSchemaAcceptsEveryLegalShape(t *testing.T) {
	checkCat := Build(loadPackages(t, checkRuleStyleGenerators(), "checkrule"), Options{})
	_, checkSchema := mustResponseSchema(t, checkCat)

	tableCat := Build(loadPackages(t, tableStyleWithDiscriminatorGenerators(), "tables", "notify"), Options{})
	_, tableSchema := mustResponseSchema(t, tableCat)

	cases := []struct {
		name     string
		schema   any
		instance string
	}{
		{"addRule with only its required kwargs", checkSchema,
			`{"invocations":[{"action":"addRule","kwargs":{"id":"x","body":"y"}}]}`},
		{"addRule with its optional kwargs bound", checkSchema,
			`{"invocations":[{"action":"addRule","kwargs":{"id":"x","body":"y","severity":"warning","tags":["a","b"]}}]}`},
		{"addCatalogEntry", checkSchema,
			`{"invocations":[{"action":"addCatalogEntry","kwargs":{"id":"x"}}]}`},
		{"addField with a pipeline placeholder in its target", tableSchema,
			`{"invocations":[{"action":"addField","kwargs":{"table":"widget","field":"name"}}]}`},
		{"addNotify, created variant, extra bound", tableSchema,
			`{"invocations":[{"action":"addNotify","kwargs":{"resource":"widgets","event":"created","extra":"note"}}]}`},
		{"addNotify, deleted variant, extra omitted", tableSchema,
			`{"invocations":[{"action":"addNotify","kwargs":{"resource":"widgets","event":"deleted"}}]}`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := validateAgainstSchema(c.schema, mustDecode(t, c.instance)); err != nil {
				t.Errorf("%s: %v\n%s", c.name, err, c.instance)
			}
		})
	}
}

// The schema is compiled only from what actions.yaml already declares.
// Nothing about severity's prose description - the part of the catalog that
// exists purely to be read by the model, per prov-2026-c5697387 - may reach
// the wire contract a validator holds the model to.
func TestResponseSchemaCarriesNoDescriptions(t *testing.T) {
	cat := Build(loadPackages(t, checkRuleStyleGenerators(), "checkrule"), Options{})
	raw, _ := mustResponseSchema(t, cat)
	if containsSubstring(raw, "Prose.") {
		t.Errorf("the schema leaked a kwarg's description:\n%s", raw)
	}
	if containsKeyBytes(raw, `"description"`) {
		t.Errorf("the schema carries a description key; it is a wire validation contract, not the prompt:\n%s", raw)
	}
}

func containsSubstring(b []byte, s string) bool {
	return len(s) > 0 && indexOf(string(b), s) >= 0
}

func containsKeyBytes(b []byte, key string) bool {
	return indexOf(string(b), key) >= 0
}

func indexOf(haystack, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

// The most important test this record can ship: a real accepted sedum grow
// recording, produced against evals/testdata/todo-rails's real generator
// package (evals/behavior/.results/todo-rails-20260828T151853-73118-logs/
// recording.json, a validated run this eval harness actually wrote - Sedum
// only writes a recording after Phase 5 accepts it), must validate against
// its own record's compiled schema. If this fails, the schema is rejecting a
// binding the model was already correctly rewarded for making, which is
// worse than shipping nothing.
//
// The catalog is filtered to the recording's own authorized paths first, so
// the schema under test is exactly the "what survives Phase 4's filter" one
// the record specifies - not the broader, unfiltered one.
func TestResponseSchemaValidatesARealAcceptedRecording(t *testing.T) {
	root := filepath.Join("..", "..", "evals", "testdata", "todo-rails", "generators", "defined")
	set, findings, err := genpkg.Load(root, genpkg.Options{})
	if err != nil {
		t.Fatalf("loading the real todo-rails generator package: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("todo-rails generator package does not load: %s", f)
		}
	}
	pkg, ok := set.Lookup("rails")
	if !ok {
		t.Fatal(`todo-rails does not define a "rails" package; fixture layout changed`)
	}

	cat := Build([]*genpkg.Package{pkg}, Options{})

	// The six paths the real recording authorized, copied from
	// evals/behavior/.results/todo-rails-20260828T151853-73118-logs/
	// recording.json's own "files" list.
	authorized := []string{
		"app/controllers/todos_controller.rb",
		"app/models/todo.rb",
		"config/routes.rb",
		"db/migrate/20260814000000_create_todos.rb",
		"test/controllers/todos_controller_test.rb",
		"test/fixtures/todos.yml",
	}
	filtered, _ := FilterForRecord(cat, authorized)

	raw, schema := mustResponseSchema(t, filtered)

	invocationsJSON, err := os.ReadFile(filepath.Join("testdata", "response_schema", "todo-rails-accepted-invocations.json"))
	if err != nil {
		t.Fatalf("reading the real accepted recording fixture: %v", err)
	}
	var invocations []recording.Invocation
	if err := json.Unmarshal(invocationsJSON, &invocations); err != nil {
		t.Fatalf("fixture does not decode as []recording.Invocation: %v", err)
	}
	if len(invocations) == 0 {
		t.Fatal("fixture carries no invocations; nothing would be proven")
	}

	instanceJSON, err := json.Marshal(map[string]any{"invocations": invocations})
	if err != nil {
		t.Fatalf("re-encoding the recorded invocations: %v", err)
	}
	instance := mustDecode(t, string(instanceJSON))

	if err := validateAgainstSchema(schema, instance); err != nil {
		t.Errorf("a real accepted recording was rejected by its own record's generated schema: %v\n\nschema:\n%s\n\nrecording:\n%s",
			err, raw, instanceJSON)
	}
}

// validateAgainstSchema checks instance against the subset of JSON Schema
// this record specifies ResponseSchema to emit: type, const, properties,
// required, additionalProperties, items, and oneOf. It is test scaffolding,
// not a general validator, and it exists only because this repo's go.mod
// cannot be edited under this record's affected_scope to add a third-party
// one (go.mod is read-only, governed entirely by other, already-implemented
// records).
func validateAgainstSchema(schema, instance any) error {
	s, ok := schema.(map[string]any)
	if !ok {
		return fmt.Errorf("schema is not an object: %#v", schema)
	}

	if raw, ok := s["oneOf"]; ok {
		branches, ok := raw.([]any)
		if !ok {
			return fmt.Errorf("oneOf is not an array")
		}
		var matched int
		var errs []error
		for i, b := range branches {
			if err := validateAgainstSchema(b, instance); err != nil {
				errs = append(errs, fmt.Errorf("branch %d: %v", i, err))
				continue
			}
			matched++
		}
		switch matched {
		case 0:
			return fmt.Errorf("matched no oneOf branch: %v", errs)
		case 1:
			return nil
		default:
			return fmt.Errorf("matched %d oneOf branches, want exactly 1", matched)
		}
	}

	if raw, ok := s["const"]; ok {
		a, _ := json.Marshal(raw)
		b, _ := json.Marshal(instance)
		if string(a) != string(b) {
			return fmt.Errorf("const %s does not match %s", a, b)
		}
		return nil
	}

	if raw, ok := s["type"]; ok {
		typeName, _ := raw.(string)
		if !schemaTypeMatches(typeName, instance) {
			return fmt.Errorf("type %q does not match value %#v", typeName, instance)
		}
	}

	switch v := instance.(type) {
	case map[string]any:
		propsRaw, _ := s["properties"].(map[string]any)
		if reqRaw, ok := s["required"]; ok {
			req, _ := reqRaw.([]any)
			for _, r := range req {
				name, _ := r.(string)
				if _, present := v[name]; !present {
					return fmt.Errorf("missing required property %q", name)
				}
			}
		}
		if addl, ok := s["additionalProperties"]; ok {
			if allowed, isBool := addl.(bool); isBool && !allowed {
				for name := range v {
					if _, declared := propsRaw[name]; !declared {
						return fmt.Errorf("undeclared property %q (additionalProperties: false)", name)
					}
				}
			}
		}
		for name, sub := range propsRaw {
			val, present := v[name]
			if !present {
				continue
			}
			if err := validateAgainstSchema(sub, val); err != nil {
				return fmt.Errorf("property %q: %v", name, err)
			}
		}
	case []any:
		if itemsRaw, ok := s["items"]; ok {
			for i, elem := range v {
				if err := validateAgainstSchema(itemsRaw, elem); err != nil {
					return fmt.Errorf("item %d: %v", i, err)
				}
			}
		}
	}
	return nil
}

func schemaTypeMatches(t string, v any) bool {
	switch t {
	case "object":
		_, ok := v.(map[string]any)
		return ok
	case "array":
		_, ok := v.([]any)
		return ok
	case "string":
		_, ok := v.(string)
		return ok
	case "boolean":
		_, ok := v.(bool)
		return ok
	case "integer":
		f, ok := v.(float64)
		return ok && f == float64(int64(f))
	case "number":
		_, ok := v.(float64)
		return ok
	default:
		return true
	}
}
