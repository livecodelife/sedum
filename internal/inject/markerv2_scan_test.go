package inject

import (
	"reflect"
	"testing"
)

// prov-2026-5a43158d
func TestBodyStateType(t *testing.T) {
	if k := reflect.TypeOf(BodyState("")).Kind(); k != reflect.String {
		t.Errorf("BodyState has kind %s, want string", k)
	}
	tests := []struct {
		name string
		got  BodyState
		want string
	}{
		{"BodyIntact", BodyIntact, "intact"},
		{"BodyEdited", BodyEdited, "edited"},
		{"BodyUnverifiable", BodyUnverifiable, "unverifiable"},
	}
	for _, tc := range tests {
		if string(tc.got) != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}

// prov-2026-b00850d0
func TestOwnedRegionBodyState(t *testing.T) {
	const hash = "sha256:aaaa"
	region := func(m Marker, bodyHash string) OwnedRegion {
		return OwnedRegion{Path: "f.rb", Region: Region{Marker: m, Start: 3, End: 9, BodyHash: bodyHash}}
	}
	tests := []struct {
		name string
		in   OwnedRegion
		want BodyState
	}{
		{"v2 with the recorded hash", region(Marker{Version: 2, Body: hash}, hash), BodyIntact},
		{"v2 with another hash", region(Marker{Version: 2, Body: hash}, "sha256:bbbb"), BodyEdited},
		{"v2 whose region hash is empty", region(Marker{Version: 2, Body: hash}, ""), BodyEdited},
		{"v2 with no recorded body", region(Marker{Version: 2}, hash), BodyUnverifiable},
		{"v2 with no recorded body and no hash", region(Marker{Version: 2}, ""), BodyUnverifiable},
		{"v1 marker", region(Marker{Version: 0}, hash), BodyUnverifiable},
		{"v1 marker that happens to carry a body", region(Marker{Version: 0, Body: hash}, hash), BodyUnverifiable},
		{"a version other than 2", region(Marker{Version: 1, Body: hash}, hash), BodyUnverifiable},
		{"foreign marker", region(Marker{Foreign: true}, hash), BodyUnverifiable},
		{"foreign marker with matching hashes", region(Marker{Foreign: true, Version: 2, Body: hash}, hash), BodyUnverifiable},
		{"the zero region", OwnedRegion{}, BodyUnverifiable},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.in.BodyState(); got != tc.want {
				t.Errorf("BodyState = %q, want %q", got, tc.want)
			}
		})
	}
}

// The state depends on the OwnedRegion alone: its path, its package and its
// offsets say nothing about it.
func TestOwnedRegionBodyStateReadsOnlyTheRegion(t *testing.T) {
	a := OwnedRegion{Path: "a.rb", Region: Region{Marker: Marker{Version: 2, Body: "sha256:x"}, Start: 0, End: 10, BodyHash: "sha256:x"}}
	b := OwnedRegion{Path: "b/c.go", Region: Region{Marker: Marker{Version: 2, Body: "sha256:x"}, Start: 500, End: 900, BodyHash: "sha256:x"}}
	if a.BodyState() != BodyIntact || b.BodyState() != BodyIntact {
		t.Errorf("states = %q, %q; want intact for both", a.BodyState(), b.BodyState())
	}
}
