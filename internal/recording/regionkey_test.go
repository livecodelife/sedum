package recording

import (
	"encoding/json"
	"sort"
	"testing"
)

// RegionKey identifies a region among a record's regions (prov-2026-0e23ee6a).

func TestRegionKeyIsTheFourPartsJoinedByNUL(t *testing.T) {
	got := RegionKey("app/models/user.rb", "addSetting", "strict", map[string]any{"name": "x"})
	want := "app/models/user.rb\x00addSetting\x00strict\x00{\"name\":\"x\"}"
	if got != want {
		t.Errorf("RegionKey = %q, want %q", got, want)
	}
}

func TestRegionKeyWithAnEmptyVariantKeepsTheEmptyPart(t *testing.T) {
	got := RegionKey("a.rb", "addField", "", map[string]any{"f": "x"})
	want := "a.rb\x00addField\x00\x00{\"f\":\"x\"}"
	if got != want {
		t.Errorf("RegionKey = %q, want %q", got, want)
	}
}

func TestRegionKeySortsTheIdKeys(t *testing.T) {
	got := RegionKey("p", "a", "v", map[string]any{"b": 1, "a": 2, "c": "z"})
	want := "p\x00a\x00v\x00{\"a\":2,\"b\":1,\"c\":\"z\"}"
	if got != want {
		t.Errorf("RegionKey = %q, want %q", got, want)
	}

	// The same id built in another order is the same key.
	other := RegionKey("p", "a", "v", map[string]any{"c": "z", "a": 2, "b": 1})
	if got != other {
		t.Errorf("key order of the id changed the key:\n%q\n%q", got, other)
	}
}

func TestRegionKeyEncodesNestedIdsAndNumbersLikeEncodingJSON(t *testing.T) {
	id := map[string]any{
		"n":    json.Number("1000000"),
		"list": []any{"a", "b"},
		"obj":  map[string]any{"y": true, "x": nil},
	}
	got := RegionKey("p", "a", "", id)
	encoded, err := json.Marshal(id)
	if err != nil {
		t.Fatal(err)
	}
	want := "p\x00a\x00\x00" + string(encoded)
	if got != want {
		t.Errorf("RegionKey = %q, want %q", got, want)
	}
}

func TestRegionKeyEncodesANilIdAsAnEmptyObject(t *testing.T) {
	if got, want := RegionKey("p", "a", "v", nil), "p\x00a\x00v\x00{}"; got != want {
		t.Errorf("nil id: RegionKey = %q, want %q", got, want)
	}
	if got, want := RegionKey("p", "a", "v", map[string]any{}), "p\x00a\x00v\x00{}"; got != want {
		t.Errorf("empty id: RegionKey = %q, want %q", got, want)
	}
}

func TestRegionKeyEncodesAnIdThatCannotBeEncodedAsAnEmptyObject(t *testing.T) {
	got := RegionKey("p", "a", "v", map[string]any{"bad": make(chan int)})
	if want := "p\x00a\x00v\x00{}"; got != want {
		t.Errorf("RegionKey = %q, want %q", got, want)
	}
	got = RegionKey("p", "a", "v", map[string]any{"bad": func() {}})
	if want := "p\x00a\x00v\x00{}"; got != want {
		t.Errorf("RegionKey = %q, want %q", got, want)
	}
}

func TestRegionKeysAreEqualExactlyWhenTheFourPartsAre(t *testing.T) {
	base := func() (string, string, string, map[string]any) {
		return "p.rb", "act", "var", map[string]any{"k": "v"}
	}
	p, a, v, id := base()
	want := RegionKey(p, a, v, id)

	p2, a2, v2, id2 := base()
	if got := RegionKey(p2, a2, v2, id2); got != want {
		t.Errorf("equal parts gave different keys:\n%q\n%q", got, want)
	}

	different := map[string]string{
		"path":     RegionKey("q.rb", a, v, id),
		"action":   RegionKey(p, "act2", v, id),
		"variant":  RegionKey(p, a, "var2", id),
		"id key":   RegionKey(p, a, v, map[string]any{"j": "v"}),
		"id val":   RegionKey(p, a, v, map[string]any{"k": "w"}),
		"id extra": RegionKey(p, a, v, map[string]any{"k": "v", "z": 1}),
		"id none":  RegionKey(p, a, v, nil),
	}
	for what, got := range different {
		if got == want {
			t.Errorf("a different %s gave the same key %q", what, got)
		}
	}
}

// Sorting keys as strings orders regions by path, then action, then variant,
// then id. The NUL separator sorts below every other character, so a path that
// is a prefix of another sorts first rather than interleaving with it.
func TestSortingKeysOrdersByPathActionVariantThenId(t *testing.T) {
	type region struct {
		path, action, variant string
		id                    map[string]any
	}
	ordered := []region{
		{"a", "x", "", map[string]any{"n": 1}},
		{"a", "x", "", map[string]any{"n": 2}},
		{"a", "x", "v", map[string]any{}},
		{"a", "xa", "", map[string]any{}},
		{"a", "y", "", map[string]any{}},
		{"a.rb", "a", "", map[string]any{}},
		{"a/b", "a", "", map[string]any{}},
		{"b", "a", "", map[string]any{}},
	}

	keys := make([]string, len(ordered))
	for i, r := range ordered {
		keys[i] = RegionKey(r.path, r.action, r.variant, r.id)
	}
	// Shuffle deterministically: reverse.
	shuffled := make([]string, len(keys))
	for i, k := range keys {
		shuffled[len(keys)-1-i] = k
	}
	sort.Strings(shuffled)

	for i := range keys {
		if shuffled[i] != keys[i] {
			t.Errorf("position %d after sorting is %q, want %q", i, shuffled[i], keys[i])
		}
	}
}
