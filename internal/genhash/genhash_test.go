package genhash

import "testing"

// The expected values come from an independent implementation of DESIGN.md's
// canonical encoding (Python's json and unicodedata), not from this package,
// so a test passing means two implementations agree rather than that this one
// agrees with itself.
func TestOfMatchesTheCanonicalEncoding(t *testing.T) {
	cases := []struct {
		name string
		in   Input
		want string
	}{
		{
			name: "an empty record hashes its empty fields, not nothing",
			in:   Input{},
			want: "sha256:7923a2ea36ef30a60c1a9cee38850f45e9dc2c661e1a7f052d4e03f784e69b7d",
		},
		{
			name: "every field takes part",
			in: Input{
				Intent:         "Add a users index.",
				Constraints:    []string{"Use the existing model."},
				AffectedScope:  []string{"app/controllers/users_controller.rb"},
				ForbiddenScope: []string{"db/**"},
			},
			want: "sha256:c979c75b0a33724658799be01536c4eb7007b88063267e2cd1ac2db477146941",
		},
		{
			name: "whitespace is trimmed and collapsed, and list order is kept",
			in: Input{
				Intent:         "  Add\n\n a\tusers  index. \n",
				Constraints:    []string{"  keep   order ", "b"},
				AffectedScope:  []string{"z.rb", "a.rb"},
				ForbiddenScope: []string{},
			},
			want: "sha256:ac2fe4bb5a199f0f10e08c2145d489896b0536d83c16fb028c4e96df0abbce98",
		},
		{
			name: "strings are NFC-normalized",
			in:   Input{Intent: "café é"},
			want: "sha256:91167b36f65abe6bbab3a5869f2cf2768dd8042f44a4e5ec27d54e09252be93f",
		},
		{
			name: "only quote, backslash and control characters are escaped",
			in:   Input{Intent: "quote \" back \\ ctl \u0001 html <&>  sep"},
			want: "sha256:e15378f503a2834a91695b9264329fd97350120095529c45b57d660737f22bbf",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Of(c.in); got != c.want {
				t.Errorf("Of() = %s, want %s", got, c.want)
			}
		})
	}
}

// A status flip or a re-pin must not look like a content change, so the hash
// depends on nothing but the four fields - and a nil list and an empty one are
// the same record.
func TestOfTreatsNilAndEmptyListsAlike(t *testing.T) {
	a := Of(Input{Intent: "x"})
	b := Of(Input{Intent: "x", Constraints: []string{}, AffectedScope: []string{}, ForbiddenScope: []string{}})
	if a != b {
		t.Errorf("nil lists hash to %s, empty lists to %s", a, b)
	}
}

func TestOfChangesWithEachField(t *testing.T) {
	base := Input{Intent: "i", Constraints: []string{"c"}, AffectedScope: []string{"a"}, ForbiddenScope: []string{"f"}}
	seen := map[string]string{"base": Of(base)}

	variants := map[string]Input{
		"intent":          {Intent: "j", Constraints: base.Constraints, AffectedScope: base.AffectedScope, ForbiddenScope: base.ForbiddenScope},
		"constraints":     {Intent: base.Intent, Constraints: []string{"d"}, AffectedScope: base.AffectedScope, ForbiddenScope: base.ForbiddenScope},
		"affected_scope":  {Intent: base.Intent, Constraints: base.Constraints, AffectedScope: []string{"b"}, ForbiddenScope: base.ForbiddenScope},
		"forbidden_scope": {Intent: base.Intent, Constraints: base.Constraints, AffectedScope: base.AffectedScope, ForbiddenScope: []string{"g"}},
		"list order":      {Intent: base.Intent, Constraints: []string{"c", "d"}, AffectedScope: base.AffectedScope, ForbiddenScope: base.ForbiddenScope},
		"list order flip": {Intent: base.Intent, Constraints: []string{"d", "c"}, AffectedScope: base.AffectedScope, ForbiddenScope: base.ForbiddenScope},
	}
	for name, in := range variants {
		h := Of(in)
		for other, prev := range seen {
			if h == prev {
				t.Errorf("changing %s hashes the same as %s: %s", name, other, h)
			}
		}
		seen[name] = h
	}
}
