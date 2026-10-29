package gitcheck

import "testing"

// Provenance 3.0.0 reports the implemented-record refusal under rule id
// REVISE001, earlier provenance under STATE001, with the same message text. Both
// are the refusal; neither id without the text is (prov-2026-ab0dd291).

func TestIsRefusalAcceptsBothRuleIDs(t *testing.T) {
	cases := []struct {
		name string
		f    Finding
		want bool
	}{
		{"REVISE001 with the text", Finding{RuleID: "REVISE001", Message: refusal}, true},
		{"REVISE001 with a sha prefix", Finding{RuleID: "REVISE001", Message: "abc1234: " + refusal}, true},
		{"STATE001 with the text", Finding{RuleID: "STATE001", Message: refusal}, true},
		{"REVISE001 without the text", Finding{RuleID: "REVISE001", Message: "snapshot is append only"}, false},
		{"STATE001 without the text", Finding{RuleID: "STATE001", Message: "snapshot is append only"}, false},
		{"REVISE001 with an empty message", Finding{RuleID: "REVISE001"}, false},
		{"the text under another id", Finding{RuleID: "SCOPE001", Message: refusal}, false},
	}
	for _, c := range cases {
		if got := IsRefusal(c.f); got != c.want {
			t.Errorf("%s: IsRefusal = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestARevise001ResultIsReadAsTheRefusal(t *testing.T) {
	got := findingOf(t, "REVISE001", refusal, "", "prov-2026-aaaaaaaa")
	if got.RuleID != "REVISE001" || !IsRefusal(got) {
		t.Errorf("finding = %+v, want the REVISE001 refusal", got)
	}
}

func TestAllowedWaivesTheRevise001RefusalOfTheTaggedRecord(t *testing.T) {
	const r = "prov-2026-aaaaaaaa"
	const other = "prov-2026-bbbbbbbb"
	ref := Finding{RuleID: "REVISE001", Message: refusal, RecordID: r}
	cases := []struct {
		name  string
		proof string
		f     Finding
		tag   string
		want  bool
	}{
		{"regeneration waives it", "regeneration", ref, r, true},
		{"migration waives it", "migration", ref, r, true},
		{"regeneration, another record's", "regeneration", Finding{RuleID: "REVISE001", Message: refusal, RecordID: other}, r, false},
		{"migration, another record's", "migration", Finding{RuleID: "REVISE001", Message: refusal, RecordID: other}, r, false},
		{"no record id", "regeneration", Finding{RuleID: "REVISE001", Message: refusal}, r, false},
		{"without the text", "migration", Finding{RuleID: "REVISE001", Message: "another revise finding", RecordID: r}, r, false},
		{"another proof name", "other", ref, r, false},
	}
	for _, c := range cases {
		if got := Allowed(c.proof, c.f, c.tag); got != c.want {
			t.Errorf("%s: Allowed = %v, want %v", c.name, got, c.want)
		}
	}
}
