package gitcheck

import (
	"reflect"
	"testing"
)

// The finding unit reads a result as a finding, recognises the
// implemented-record refusal, parses commit tags and says which findings a
// proof may waive.

const refusal = "prov-2026-aaaaaaaa is implemented: run 'provenance revise prov-2026-aaaaaaaa' to start a new version, then tag the commit with it."

func findingOf(t *testing.T, rule, msg, file, record string) Finding {
	t.Helper()
	log := parse(t, logText(rulesAB, "["+result(rule, 0, msg, file, record)+"]", "[]"))
	return FindingOf(log.Results()[0])
}

func TestFindingOfReadsTheFieldsOfAResult(t *testing.T) {
	got := findingOf(t, "SCOPE001", "forbids", "a/b.go", "prov-2026-aaaaaaaa")
	want := Finding{RuleID: "SCOPE001", Message: "forbids", File: "a/b.go", RecordID: "prov-2026-aaaaaaaa"}
	if got != want {
		t.Errorf("Finding = %+v, want %+v", got, want)
	}
	if bare := findingOf(t, "TAG001", "m", "", ""); bare.File != "" || bare.RecordID != "" {
		t.Errorf("a result without location or recordId = %+v", bare)
	}
}

func TestIsRefusalNeedsTheRuleIDAndTheMessageText(t *testing.T) {
	cases := []struct {
		name string
		f    Finding
		want bool
	}{
		{"both", Finding{RuleID: "STATE001", Message: refusal}, true},
		{"with a sha prefix", Finding{RuleID: "STATE001", Message: "abc1234: " + refusal}, true},
		{"another state finding", Finding{RuleID: "STATE001", Message: "snapshot is append only"}, false},
		{"the text under another id", Finding{RuleID: "SCOPE001", Message: refusal}, false},
	}
	for _, c := range cases {
		if got := IsRefusal(c.f); got != c.want {
			t.Errorf("%s: IsRefusal = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestAllowed(t *testing.T) {
	const r = "prov-2026-aaaaaaaa"
	ref := Finding{RuleID: "STATE001", Message: refusal, RecordID: r}
	scope := Finding{RuleID: "SCOPE001", Message: "forbids", File: "x.go", RecordID: r}
	cases := []struct {
		name  string
		proof string
		f     Finding
		tag   string
		want  bool
	}{
		{"regeneration waives the refusal of R", "regeneration", ref, r, true},
		{"regeneration does not waive scope", "regeneration", scope, r, false},
		{"migration waives the refusal of the tag", "migration", ref, r, true},
		{"migration waives scope of the tag", "migration", scope, r, true},
		{"another record's refusal", "regeneration", Finding{RuleID: "STATE001", Message: refusal, RecordID: "prov-2026-bbbbbbbb"}, r, false},
		{"another record's scope under migration", "migration", Finding{RuleID: "SCOPE001", RecordID: "prov-2026-bbbbbbbb"}, r, false},
		{"no record id", "migration", Finding{RuleID: "SCOPE001", Message: "forbids"}, "", false},
		{"a tag error", "migration", Finding{RuleID: "TAG001", Message: "tag", RecordID: r}, r, false},
		{"a lint finding", "migration", Finding{RuleID: "PROV010", RecordID: r}, r, false},
		{"an evidence finding", "regeneration", Finding{RuleID: "EVID001", RecordID: r}, r, false},
		{"an unknown id", "migration", Finding{RuleID: "NEW999", RecordID: r}, r, false},
		{"another proof name", "other", ref, r, false},
	}
	for _, c := range cases {
		if got := Allowed(c.proof, c.f, c.tag); got != c.want {
			t.Errorf("%s: Allowed = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestParseTags(t *testing.T) {
	msg := "Fix the thing [prov-2026-aaaaaaaa@3]\n\nAlso [prov-2026-bbbbbbbb] and [prov-2026-123] and [prov-2026-cccccccc-x1@12], but not [prov-2026-ZZ] or [prov-2026-dddddddd@0] or prov-2026-eeeeeeee"
	got := ParseTags(msg)
	want := []Tag{
		{ID: "prov-2026-aaaaaaaa", Version: 3, HasVersion: true},
		{ID: "prov-2026-bbbbbbbb"},
		{ID: "prov-2026-123"},
		{ID: "prov-2026-cccccccc-x1", Version: 12, HasVersion: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseTags = %+v, want %+v", got, want)
	}
	if got := ParseTags("no tags here"); len(got) != 0 {
		t.Errorf("ParseTags = %+v, want none", got)
	}
}
