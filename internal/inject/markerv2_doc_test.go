package inject

import (
	"os"
	"strings"
	"testing"
)

// TestMarkerV2SpecSections keeps MARKER_V2.md honest about what it must define
// before any of it is implemented: the marker keys, both hashes, the recording
// layout, the check class, the migration and the open questions.
func TestMarkerV2SpecSections(t *testing.T) {
	b, err := os.ReadFile("../../MARKER_V2.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	for _, heading := range []string{
		"## Marker v2",
		"### Canonical kwargs",
		"### Body",
		"## Recordings",
		"### Verifying a region against its recording",
		"### Re-render",
		"## `grow --check`",
		"## Migration",
		"## Conformance",
		"## Size",
		"## Open questions",
	} {
		if !strings.Contains(doc, "\n"+heading+"\n") {
			t.Errorf("MARKER_V2.md has no %q section", heading)
		}
	}
	for _, key := range []string{"`v`", "`id`", "`kwargs`", "`body`", ".sedum/recordings/<record-id>.json"} {
		if !strings.Contains(doc, key) {
			t.Errorf("MARKER_V2.md never mentions %s", key)
		}
	}
}
