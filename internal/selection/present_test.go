package selection

import (
	"strings"
	"testing"
)

// The selection prompt lists the regions other records already own in this
// record's files, so the model does not select them (prov-2026-cc73cf61).
func TestPresentRegionsAreListedAsDoNotSelect(t *testing.T) {
	req := railsRequest(t)
	req.Present = []Present{{
		Path:   "app/controllers/users_controller.rb",
		Action: "createControllerMethod",
		Kwargs: map[string]any{"controller": "users", "name": "PRESENT-77c1"},
		Owner:  "prov-2026-bbbbbbbb",
	}}

	user := promptFor(t, req)
	section := strings.Index(user, "Already present")
	if section < 0 {
		t.Fatalf("no already-present section:\n%s", user)
	}
	for _, want := range []string{"do not select", "app/controllers/users_controller.rb", "createControllerMethod", "PRESENT-77c1", "prov-2026-bbbbbbbb"} {
		if i := strings.Index(user, want); i < 0 || (want != "do not select" && i < section) {
			t.Errorf("the already-present section does not carry %q:\n%s", want, user)
		}
	}
	if constraints := strings.Index(user, "No writes."); section < constraints {
		t.Errorf("the already-present section must follow the record's own constraints")
	}
}

func TestNoPresentRegionsMeansNoSection(t *testing.T) {
	if user := promptFor(t, railsRequest(t)); strings.Contains(user, "Already present") {
		t.Errorf("a prompt with nothing present grew a section:\n%s", user)
	}
}
