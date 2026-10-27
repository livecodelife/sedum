package evidence

import (
	"errors"
	"reflect"
	"testing"
)

// The code hash of a record is the hash docs/HASHING.md section 7 defines for
// the gen and pkg of its owned regions; the expected values are the conformance
// values, so an old code_hash still matches.

func TestCodeHashMatchesTheConformanceValues(t *testing.T) {
	aa := Region{Gen: "sha256:aa", Pkg: "sha256:01"}
	bb := Region{Gen: "sha256:bb", Pkg: "sha256:02"}
	cases := []struct {
		name    string
		regions []Region
		want    string
	}{
		{"no regions", nil, "sha256:4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945"},
		{"two regions", []Region{aa, bb}, "sha256:9aaa1fe3e1be19258596ec7d6e2c8185d8797c55ba85159bc80c3fc528fbd904"},
		{"the same two in the other order", []Region{bb, aa}, "sha256:9aaa1fe3e1be19258596ec7d6e2c8185d8797c55ba85159bc80c3fc528fbd904"},
		{"a seeded region with no hashes is ignored", []Region{aa, bb, {Tier: "seeded"}}, "sha256:9aaa1fe3e1be19258596ec7d6e2c8185d8797c55ba85159bc80c3fc528fbd904"},
		{"an owned tier counts like an absent one", []Region{{Gen: "sha256:aa", Pkg: "sha256:01", Tier: "owned"}, bb}, "sha256:9aaa1fe3e1be19258596ec7d6e2c8185d8797c55ba85159bc80c3fc528fbd904"},
		{"duplicates are kept", []Region{aa, aa}, "sha256:2ecdfde41f5c31933ab7587b81c7645940d70ff3f05598717ce6e70d8f8bcf82"},
		{"a pkg change", []Region{{Gen: "sha256:aa", Pkg: "sha256:09"}, bb}, "sha256:f55a1e47c7805e78e7b5ff0114c7163d6993fa62d0bb634c94b2e5d8ad701c89"},
	}
	for _, c := range cases {
		got, err := CodeHash(c.regions)
		if err != nil || got != c.want {
			t.Errorf("%s: CodeHash = %q, %v; want %q", c.name, got, err, c.want)
		}
	}
}

func TestCodeHashRefusesAnOwnedRegionWithoutGenOrPkg(t *testing.T) {
	for name, r := range map[string]Region{
		"no pkg":  {Gen: "sha256:aa"},
		"no gen":  {Pkg: "sha256:01"},
		"neither": {},
	} {
		got, err := CodeHash([]Region{{Gen: "sha256:bb", Pkg: "sha256:02"}, r})
		if got != "" || !errors.Is(err, ErrUnstamped) || err != ErrUnstamped {
			t.Errorf("%s: CodeHash = %q, %v; want empty and ErrUnstamped", name, got, err)
		}
	}
	if _, err := CodeHash([]Region{{Tier: "seeded"}}); err != nil {
		t.Errorf("a dropped region caused %v", err)
	}
}

func TestErrUnstampedSaysWhatIsMissingAndHowToFixIt(t *testing.T) {
	msg := ErrUnstamped.Error()
	for _, want := range []string{"gen", "pkg", "sedum grow --stamp"} {
		if !contains(msg, want) {
			t.Errorf("ErrUnstamped = %q lacks %q", msg, want)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestOwnedKeepsEmptyAndOwnedTiersInOrderAndReturnsAnEmptySliceWhenNoneRemain(t *testing.T) {
	in := []Region{{Gen: "1"}, {Gen: "2", Tier: "seeded"}, {Gen: "3", Tier: "owned"}, {Gen: "4", Tier: "other"}}

	got := Owned(in)

	if !reflect.DeepEqual(got, []Region{{Gen: "1"}, {Gen: "3", Tier: "owned"}}) {
		t.Errorf("Owned = %+v", got)
	}
	none := Owned([]Region{{Tier: "seeded"}})
	if none == nil || len(none) != 0 {
		t.Errorf("Owned of nothing owned = %#v, want an empty slice", none)
	}
}

func TestCodeHashNeverChangesItsArgument(t *testing.T) {
	in := []Region{{Gen: "sha256:bb", Pkg: "sha256:02"}, {Gen: "sha256:aa", Pkg: "sha256:01"}, {Tier: "seeded"}}
	before := append([]Region{}, in...)

	if _, err := CodeHash(in); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(in, before) {
		t.Errorf("argument changed: %+v", in)
	}
}

func TestCanonicalEscapingOnlyEscapesQuoteBackslashAndControlCodes(t *testing.T) {
	got := canonicalRegion(Region{Gen: "a\"b\\c\n\x01", Pkg: "é/<>&  "})

	want := "{\"gen\":\"a\\\"b\\\\c\\u000a\\u0001\",\"pkg\":\"é/<>&  \"}"
	if got != want {
		t.Errorf("canonicalRegion = %q, want %q", got, want)
	}
}
