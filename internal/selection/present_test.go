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

// An anchor a region owned by another record already fills is not unfilled, so
// it earns no completeness follow-up: the model could only answer by selecting
// that record's region, which Phase 7 would refuse (prov-2026-53f4f82b).
func TestAnAnchorAnotherRecordFillsEarnsNoFollowUp(t *testing.T) {
	req := railsRequest(t)
	req.Files[0].Rendered = "class UsersController\n  # sedum:anchor:class_body\nend\n"
	req.Present = []Present{{
		Path:   req.Files[0].Path,
		Action: "createControllerMethod",
		Kwargs: map[string]any{"controller": "users", "name": "index"},
		Owner:  "prov-2026-bbbbbbbb",
	}}

	_, client, err := selectWith(t, req, 0, `{"invocations":[]}`)
	if err != nil {
		t.Fatalf("select: %v", err)
	}
	if len(client.seen) != 1 {
		t.Errorf("the model was called %d time(s), want 1 - another record already fills class_body", len(client.seen))
	}
}

// The follow-up names the file and the anchor and claims nothing about who
// created the file, which may well have existed before the run.
func TestTheCompletenessNoteClaimsNoCreation(t *testing.T) {
	req := railsRequest(t)
	req.Files[0].Rendered = "class UsersController\n  # sedum:anchor:class_body\nend\n"

	_, client, err := selectWith(t, req, 0, `{"invocations":[]}`, `{"invocations":[]}`)
	if err != nil {
		t.Fatalf("select: %v", err)
	}
	followUp := client.seen[1][len(client.seen[1])-1].Content
	if strings.Contains(followUp, "created") {
		t.Errorf("the note claims the run created a file:\n%s", followUp)
	}
	if !strings.Contains(followUp, req.Files[0].Path) || !strings.Contains(followUp, "class_body") {
		t.Errorf("the note should name the file and the anchor:\n%s", followUp)
	}
}
