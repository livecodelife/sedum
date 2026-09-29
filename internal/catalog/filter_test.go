package catalog

import (
	"fmt"
	"strings"
	"testing"
)

// checkRuleGenerators is modelled on a real package shape: one action whose
// target is per-record (a rule file named by an id) and one whose target is a
// single file every record shares (a catalog every rule registers itself
// into). It is exactly the shape that produced the measured failure - a
// record authorizing only its own rule file still saw the catalog-entry
// action, bound it, and paid a retry for unauthorized_path.
func checkRuleGenerators() map[string]string {
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

// tableStyleGenerators is modelled on the other real shape measured: one
// action whose target runs a value through a pipeline of transforms
// (pluralization) and one whose target is a free-standing placeholder with no
// transform at all. Both are the case the filter must never narrow past
// "could this match" - a pipeline and a bare placeholder alike stand for a
// value nothing has bound yet.
func tableStyleGenerators() map[string]string {
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
	}
}

func TestFilterDropsAnActionWhoseTargetNoAuthorizedPathCanMatch(t *testing.T) {
	packages := loadPackages(t, checkRuleGenerators(), "checkrule")
	cat := Build(packages, Options{})

	kept, removed := FilterForRecord(cat, []string{"checks/doc001.ts"})

	if _, err := lookupOnly(kept, "addRule"); err != nil {
		t.Errorf("addRule should stay: %v", err)
	}
	if got := kept.Lookup("addCatalogEntry"); len(got) != 0 {
		t.Errorf("addCatalogEntry should have been filtered out, got %+v", got)
	}

	if len(removed) != 1 || removed[0].Action != "addCatalogEntry" {
		t.Fatalf("removed = %+v, want exactly addCatalogEntry", removed)
	}
	if !strings.Contains(removed[0].Reason, "catalog.ts") {
		t.Errorf("removal reason does not name the unreachable target:\n%s", removed[0].Reason)
	}
}

func TestFilterKeepsEveryActionAuthorizedPathsCanReach(t *testing.T) {
	packages := loadPackages(t, checkRuleGenerators(), "checkrule")
	cat := Build(packages, Options{})

	kept, removed := FilterForRecord(cat, []string{"checks/doc001.ts", "catalog.ts"})

	if len(removed) != 0 {
		t.Errorf("removed %+v, want nothing removed when both targets are authorized", removed)
	}
	for _, name := range []string{"addRule", "addCatalogEntry"} {
		if len(kept.Lookup(name)) != 1 {
			t.Errorf("%s should have stayed in the filtered catalog", name)
		}
	}
}

// Pipelines and free placeholders must never cause a removal: nothing has
// bound table or migrationName yet, so both stand for any value a record's
// authorized paths might actually carry.
func TestFilterKeepsPipelinesAndPlainPlaceholders(t *testing.T) {
	packages := loadPackages(t, tableStyleGenerators(), "tables")
	cat := Build(packages, Options{})

	kept, removed := FilterForRecord(cat, []string{
		"lib/identity/humans.ts",
		"supabase/migrations/0001_humans.sql",
	})

	if len(removed) != 0 {
		t.Errorf("removed %+v, want nothing removed", removed)
	}
	for _, name := range []string{"addField", "addColumn"} {
		if len(kept.Lookup(name)) != 1 {
			t.Errorf("%s should have stayed in the filtered catalog", name)
		}
	}
}

// A free-target action - the whole pattern is one placeholder - has no fixed
// text to fail to match, so it is never removed regardless of what a record
// authorizes.
func TestFilterNeverDropsAFreeTargetAction(t *testing.T) {
	files := map[string]string{
		"free/sedum.yaml": `name: free
extensions: [".ts"]
comment_prefix: "//"
`,
		"free/files/src/{name}.ts": "// sedum:anchor:imports\n",
		"free/actions/actions.yaml": `actions:
  addImport:
    kwargs:
      file: { type: string, required: true }
      symbol: { type: string, required: true }
    injects_into: "{{file}}"
    anchor: imports
`,
		"free/actions/addImport.ts": "import { {{symbol}} }\n",
	}
	packages := loadPackages(t, files, "free")
	cat := Build(packages, Options{})

	kept, removed := FilterForRecord(cat, []string{"anything/at/all.ts"})
	if len(removed) != 0 {
		t.Errorf("a free-target action was removed: %+v", removed)
	}
	if len(kept.Lookup("addImport")) != 1 {
		t.Error("addImport should have stayed")
	}
}

// A composite is removed only when at least one child's target is
// unreachable, because expand.Targets requires every child's rendered target
// to be authorized for the composite to validate at all - a structurally
// unreachable child dooms the whole selection regardless of the others.
func TestFilterDropsACompositeWhenAnyChildIsUnreachable(t *testing.T) {
	files := map[string]string{
		"cairn/sedum.yaml": `name: cairn
extensions: [".crn"]
comment_prefix: ";;"
`,
		"cairn/files/Units/{name}/Manifest.crn": ";; sedum:anchor:steps\n",
		"cairn/files/Shared/{name}.crn":         ";; shared\n",
		"cairn/actions/actions.yaml": `actions:
  provisionStep:
    composes: [addStep, declareConstant]

  addStep:
    kwargs:
      unit: { type: string, required: true }
      step: { type: string, required: true }
    injects_into: "Units/{{unit}}/Manifest.crn"
    anchor: steps
    exposed: false

  declareConstant:
    kwargs:
      unit: { type: string, required: true }
      name: { type: string, required: true }
    injects_into: "Shared/{{name}}.crn"
    anchor: end_of_file
    exposed: false
`,
		"cairn/actions/addStep.crn":         ";; step {{step}}\n",
		"cairn/actions/declareConstant.crn": ";; const {{name}}\n",
	}
	packages := loadPackages(t, files, "cairn")
	cat := Build(packages, Options{})

	// Only a Units/ path authorized: declareConstant's Shared/{{name}}.crn can
	// never match, so the whole composite is unreachable.
	kept, removed := FilterForRecord(cat, []string{"Units/billing-runs/Manifest.crn"})
	if len(kept.Lookup("provisionStep")) != 0 {
		t.Error("provisionStep should have been removed: one child can never be authorized")
	}
	if len(removed) != 1 || removed[0].Action != "provisionStep" {
		t.Fatalf("removed = %+v, want exactly provisionStep", removed)
	}

	// Both shapes authorized (not necessarily the exact file a later
	// invocation binds to): the composite stays, because both children could
	// still resolve to an authorized path.
	kept, removed = FilterForRecord(cat, []string{
		"Units/billing-runs/Manifest.crn", "Shared/retry-limits.crn",
	})
	if len(removed) != 0 {
		t.Errorf("removed %+v, want nothing removed once both children can match", removed)
	}
	if len(kept.Lookup("provisionStep")) != 1 {
		t.Error("provisionStep should have stayed")
	}
}

// Filtering with no authorized paths at all removes everything with fixed
// text in its target - the case an empty result maps to a "nothing to offer"
// error before any model call, one level up in selection.Select.
func TestFilterWithNoAuthorizedPathsRemovesEveryFixedTarget(t *testing.T) {
	packages := loadPackages(t, checkRuleGenerators(), "checkrule")
	cat := Build(packages, Options{})

	kept, removed := FilterForRecord(cat, nil)
	if len(kept.Actions) != 0 {
		t.Errorf("kept %+v, want nothing kept with no authorized paths", kept.Actions)
	}
	if len(removed) != 2 {
		t.Errorf("removed %d action(s), want both", len(removed))
	}
}

func lookupOnly(c Catalog, name string) (Action, error) {
	matches := c.Lookup(name)
	if len(matches) != 1 {
		return Action{}, fmt.Errorf("looking up %q found %d entries, want 1", name, len(matches))
	}
	return matches[0], nil
}
