package genhash

import "testing"

// An empty shape is absent from the hash, so a record without one keeps the
// hash it had before gen_hash covered shape; a shape moves it.
func TestShapeIsPartOfTheHashOnlyWhenPresent(t *testing.T) {
	base := Input{Intent: "i", Constraints: []string{"c"}}
	none := Of(base)
	for name, shape := range map[string][]any{"nil": nil, "empty": {}, "empty entry": {map[string]any{}}} {
		in := base
		in.Shape = shape
		if got := Of(in); got != none {
			t.Errorf("%s shape hashes to %s, want the shape-less %s", name, got, none)
		}
	}
	in := base
	in.Shape = []any{map[string]any{"path": "a.go", "package": "a"}}
	if got := Of(in); got == none {
		t.Errorf("a shape did not move the hash (%s)", got)
	}
}
