package inject

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Region ownership and reconciliation (prov-2026-375f5d4a): a region belongs to
// the record named on its marker, only that record rewrites or prunes it, and
// a record's regions are reconciled against what is on disk by identity.

func indexKwargs() map[string]any {
	return map[string]any{"controller": "users", "name": "index", "collection": "users"}
}

func showKwargs() map[string]any {
	return map[string]any{"controller": "users", "name": "show"}
}

func filterKwargs(filter string) map[string]any {
	return map[string]any{"controller": "users", "filter": filter}
}

// reconcileAll names one record's controller file for reconciliation.
func reconcileOne(t *testing.T, record string) Reconcile {
	t.Helper()
	return Reconcile{RecordID: record, Targets: []Target{{Path: controllerPath, Package: loadPackage(t)}}}
}

func mustApply(t *testing.T, invocations []Invocation, opts Options) []Result {
	t.Helper()
	results, err := Apply(invocations, opts)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	return results
}

func TestSelectingAnotherRecordsRegionIsAnError(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	mustApply(t, []Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\nend\n", "prov-2026-aaaaaaaa", indexKwargs()),
	}, Options{Output: root})
	before := read(t, root)

	_, err := Apply([]Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\n  # other\nend\n", "prov-2026-bbbbbbbb", indexKwargs()),
	}, Options{Output: root})
	if err == nil {
		t.Fatal("a record overwrote a region another record owns")
	}
	for _, want := range []string{controllerPath, "prov-2026-aaaaaaaa", "prov-2026-bbbbbbbb"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the diagnostic does not name %q: %v", want, err)
		}
	}
	if got := read(t, root); got != before {
		t.Errorf("a failed run changed the file:\n%s", got)
	}
}

func TestAnAliasStillOwnsItsRegions(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	mustApply(t, []Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\nend\n", "prov-2026-01d01d01", indexKwargs()),
	}, Options{Output: root})

	owners := map[string]Owner{"prov-2026-aaaaaaaa": {ID: "prov-2026-aaaaaaaa", Aliases: []string{"prov-2026-01d01d01"}}}
	results := mustApply(t, []Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\n  # new\nend\n", "prov-2026-aaaaaaaa", indexKwargs()),
	}, Options{Output: root, Owners: owners})
	if len(results) != 1 || !results[0].Replaced {
		t.Fatalf("results = %+v, want the aliased region replaced", results)
	}
	got := read(t, root)
	if !strings.Contains(got, `"record":"prov-2026-aaaaaaaa"`) || strings.Contains(got, "prov-2026-01d01d01") {
		t.Errorf("a rewritten region should name the current id, not the alias:\n%s", got)
	}
}

func TestAnUnownedRegionIsClaimed(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	mustApply(t, []Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\nend\n", "", indexKwargs()),
	}, Options{Output: root})

	results := mustApply(t, []Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\n  # mine\nend\n", "prov-2026-aaaaaaaa", indexKwargs()),
	}, Options{Output: root})
	if len(results) != 1 || !results[0].Replaced {
		t.Fatalf("results = %+v, want the unowned region claimed in place", results)
	}
	if got := read(t, root); !strings.Contains(got, `"record":"prov-2026-aaaaaaaa"`) {
		t.Errorf("the claimed region does not name its new owner:\n%s", got)
	}
}

func TestReconcilePrunesOnlyThisRecordsUnselectedOwnedRegions(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	seeded := invocation(t, pkg, "createControllerMethod", "show", "def show\nend\n", "prov-2026-aaaaaaaa", showKwargs())
	seeded.Tier = TierSeeded
	mustApply(t, []Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\nend\n", "prov-2026-aaaaaaaa", indexKwargs()),
		invocation(t, pkg, "addBeforeFilter", "", "before_action :auth\n", "prov-2026-aaaaaaaa", filterKwargs("auth")),
		invocation(t, pkg, "addBeforeFilter", "", "before_action :log\n", "prov-2026-bbbbbbbb", filterKwargs("log")),
		seeded,
	}, Options{Output: root})

	// The record is re-selected and keeps only its index action.
	results := mustApply(t, []Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\nend\n", "prov-2026-aaaaaaaa", indexKwargs()),
	}, Options{Output: root, Reconcile: []Reconcile{reconcileOne(t, "prov-2026-aaaaaaaa")}})

	got := read(t, root)
	if strings.Contains(got, "before_action :auth") {
		t.Errorf("an owned region the record no longer selects was not pruned:\n%s", got)
	}
	for _, kept := range []string{"def index", "before_action :log", "def show"} {
		if !strings.Contains(got, kept) {
			t.Errorf("%q was removed; only this record's unselected owned regions may be pruned:\n%s", kept, got)
		}
	}

	pruned := 0
	for _, r := range results {
		if r.Pruned {
			pruned++
			if r.Action != "addBeforeFilter" || r.Record != "prov-2026-aaaaaaaa" {
				t.Errorf("pruned result = %+v, want addBeforeFilter owned by prov-2026-aaaaaaaa", r)
			}
		}
	}
	if pruned != 1 {
		t.Errorf("results report %d pruned regions, want 1: %+v", pruned, results)
	}
}

func TestAnEmptySelectionOverOwnedRegionsNeedsForce(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	mustApply(t, []Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\nend\n", "prov-2026-aaaaaaaa", indexKwargs()),
	}, Options{Output: root})
	before := read(t, root)
	reconcile := []Reconcile{reconcileOne(t, "prov-2026-aaaaaaaa")}

	_, err := Apply(nil, Options{Output: root, Reconcile: reconcile})
	if err == nil {
		t.Fatal("an empty selection removed a record's regions without --force")
	}
	if !strings.Contains(err.Error(), "prov-2026-aaaaaaaa") || !strings.Contains(err.Error(), "--force") {
		t.Errorf("the diagnostic should name the record and --force: %v", err)
	}
	if got := read(t, root); got != before {
		t.Errorf("a refused prune changed the file:\n%s", got)
	}

	// A dry run shows what --force would remove, and writes nothing.
	results := mustApply(t, nil, Options{Output: root, Reconcile: reconcile, DryRun: true})
	if len(results) != 1 || !results[0].Pruned || !results[0].Held {
		t.Errorf("dry-run results = %+v, want one pruned region held for --force", results)
	}
	if got := read(t, root); got != before {
		t.Errorf("a dry run changed the file:\n%s", got)
	}

	mustApply(t, nil, Options{Output: root, Reconcile: reconcile, Force: true})
	if got := read(t, root); strings.Contains(got, "def index") {
		t.Errorf("--force did not prune the record's regions:\n%s", got)
	}
}

// Pruning waits for the whole run: a failure anywhere leaves every region alone.
func TestAFailedRunPrunesNothing(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	mustApply(t, []Invocation{
		invocation(t, pkg, "addBeforeFilter", "", "before_action :auth\n", "prov-2026-aaaaaaaa", filterKwargs("auth")),
	}, Options{Output: root})
	before := read(t, root)

	missing := invocation(t, pkg, "addUnplanted", "", "# unplanted\n", "prov-2026-aaaaaaaa", map[string]any{"controller": "users"})
	_, err := Apply([]Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\nend\n", "prov-2026-aaaaaaaa", indexKwargs()),
		missing,
	}, Options{Output: root, Reconcile: []Reconcile{reconcileOne(t, "prov-2026-aaaaaaaa")}})
	if err == nil {
		t.Fatal("the fixture invocation was expected to fail")
	}
	if got := read(t, root); got != before {
		t.Errorf("a failed run pruned or wrote:\n%s", got)
	}
}

// Regions at a shared anchor are ordered by owner - created_at, then id - so
// a file regenerated from scratch is the same whatever order records grew in.
func TestSharedAnchorOrderIsByOwner(t *testing.T) {
	pkg := loadPackage(t)
	owners := map[string]Owner{
		"prov-2026-cccccccc": {ID: "prov-2026-cccccccc", CreatedAt: "2026-08-01"},
		"prov-2026-aaaaaaaa": {ID: "prov-2026-aaaaaaaa", CreatedAt: "2026-08-15"},
		"prov-2026-bbbbbbbb": {ID: "prov-2026-bbbbbbbb", CreatedAt: "2026-08-15"},
	}
	filter := func(name, record string) Invocation {
		return invocation(t, pkg, "addBeforeFilter", "", "before_action :"+name+"\n", record, filterKwargs(name))
	}

	orders := [][]Invocation{
		{filter("a1", "prov-2026-aaaaaaaa"), filter("a2", "prov-2026-aaaaaaaa"), filter("b", "prov-2026-bbbbbbbb"), filter("c", "prov-2026-cccccccc")},
		{filter("c", "prov-2026-cccccccc"), filter("b", "prov-2026-bbbbbbbb"), filter("a1", "prov-2026-aaaaaaaa"), filter("a2", "prov-2026-aaaaaaaa")},
		{filter("b", "prov-2026-bbbbbbbb"), filter("a1", "prov-2026-aaaaaaaa"), filter("c", "prov-2026-cccccccc"), filter("a2", "prov-2026-aaaaaaaa")},
	}

	var first string
	for i, order := range orders {
		root := output(t)
		// One run per record, in the order given, as separate grows would.
		for _, inv := range order {
			mustApply(t, []Invocation{inv}, Options{Output: root, Owners: owners})
		}
		got := read(t, root)

		want := []string{":c", ":a1", ":a2", ":b"}
		last := -1
		for _, name := range want {
			at := strings.Index(got, "before_action "+name+"\n")
			if at < last {
				t.Errorf("order %d: %s is out of place, want c, a1, a2, b:\n%s", i, name, got)
			}
			last = at
		}
		if i == 0 {
			first = got
		} else if got != first {
			t.Errorf("order %d produced a different file:\n%s\nwant:\n%s", i, got, first)
		}
	}
}

// A replaced region keeps its place even when its owner would sort elsewhere:
// existing regions are never moved.
func TestReplacementNeverMovesARegion(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	mustApply(t, []Invocation{
		invocation(t, pkg, "addBeforeFilter", "", "before_action :late\n", "prov-2026-bbbbbbbb", filterKwargs("late")),
		invocation(t, pkg, "addBeforeFilter", "", "before_action :early\n", "prov-2026-aaaaaaaa", filterKwargs("early")),
	}, Options{Output: root})
	before := read(t, root)

	mustApply(t, []Invocation{
		invocation(t, pkg, "addBeforeFilter", "", "before_action :late\n", "prov-2026-bbbbbbbb", filterKwargs("late")),
	}, Options{Output: root})
	if got := read(t, root); got != before {
		t.Errorf("rewriting a region moved it:\n%s\nwant:\n%s", got, before)
	}

	full := filepath.Join(root, filepath.FromSlash(controllerPath))
	if _, err := os.Stat(full); err != nil {
		t.Fatal(err)
	}
}
