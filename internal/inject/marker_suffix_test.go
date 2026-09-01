package inject

import (
	"reflect"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/genpkg"
)

// A target whose only comment form is delimited (prov-2026-a6f6bb81).
//
// These are the cases the marker could not express before comment_suffix, and
// the one alongside them is the compatibility claim: an absent suffix writes
// exactly the bytes Sedum wrote before the field existed. That claim is stated
// against a literal here rather than left to be inferred from the other tests
// continuing to pass.

// html is the delimited comment these cases are written against. Sedum does not
// know it is HTML - a Comment is two pieces of text and nothing more.
var html = genpkg.Comment{Prefix: "<!--", Suffix: "-->"}

// The opening line is the hard one. The attribute object is the tail of the
// line, so a suffix left in place makes it unreadable as JSON.
func TestMarkerRoundTripsThroughADelimitedComment(t *testing.T) {
	want := Marker{
		Action:  "addTextColumn",
		Variant: "required",
		Tier:    TierOwned,
		Record:  "ERP-1",
		Kwargs:  map[string]any{"column": "commodity_code"},
	}

	open, err := want.Open(html)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if !strings.HasPrefix(open, "<!-- sedum:addTextColumn:required ") {
		t.Errorf("opening marker does not open a comment: %s", open)
	}
	if !strings.HasSuffix(open, " -->") {
		t.Errorf("opening marker does not close its comment: %s", open)
	}

	if got := want.Close(html); got != "<!-- /sedum:addTextColumn:required -->" {
		t.Errorf("Close = %q", got)
	}

	got, ok, err := parseOpen(html, open)
	if err != nil || !ok {
		t.Fatalf("parseOpen(%q) = %v, %v", open, ok, err)
	}
	if got.Action != want.Action || got.Variant != want.Variant {
		t.Errorf("label = %q, want %q", got.Label(), want.Label())
	}
	if got.Record != want.Record {
		t.Errorf("record = %q, want %q", got.Record, want.Record)
	}
	if !reflect.DeepEqual(got.Kwargs, want.Kwargs) {
		t.Errorf("kwargs = %v, want %v", got.Kwargs, want.Kwargs)
	}

	// Re-emission is stable, or every rerun churns the file without changing
	// what the marker says.
	again, err := got.Open(html)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if again != open {
		t.Errorf("re-emission changed the marker:\n  %s\n  %s", open, again)
	}
}

// The compatibility claim, against literals.
func TestAnAbsentSuffixChangesNoBytes(t *testing.T) {
	shell := genpkg.Comment{Prefix: "#"}
	m := Marker{Action: "createControllerMethod", Variant: "index", Tier: TierOwned, Record: "PR-014"}

	open, err := m.Open(shell)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	const wantOpen = `# sedum:createControllerMethod:index {"tier":"owned","record":"PR-014"}`
	if open != wantOpen {
		t.Errorf("Open =\n  %s\nwant\n  %s", open, wantOpen)
	}
	if got := m.Close(shell); got != "# /sedum:createControllerMethod:index" {
		t.Errorf("Close = %q", got)
	}
}

// A region in a delimited target is found and bounded the same way a
// line-commented one is.
func TestFindRegionsThroughADelimitedComment(t *testing.T) {
	content := "<app-grid>\n" +
		`  <!-- sedum:addTextColumn {"tier":"owned","kwargs":{"column":"a"}} -->` + "\n" +
		"  <app-grid-col />\n" +
		"  <!-- /sedum:addTextColumn -->\n" +
		"</app-grid>\n"

	regions, err := FindRegions(html, content)
	if err != nil {
		t.Fatalf("FindRegions: %v", err)
	}
	if len(regions) != 1 {
		t.Fatalf("found %d regions, want 1", len(regions))
	}
	if regions[0].Marker.Action != "addTextColumn" {
		t.Errorf("action = %q", regions[0].Marker.Action)
	}
	if body := content[regions[0].Start:regions[0].End]; !strings.Contains(body, "app-grid-col") {
		t.Errorf("region does not span its body: %q", body)
	}
}

// A mismatched pair is still an error when the comments are delimited: the
// suffix comes off before the label is compared, so it is the labels that
// disagree rather than the delimiters.
func TestADelimitedRegionStillRequiresMatchingLabels(t *testing.T) {
	content := `<!-- sedum:addTextColumn {"tier":"owned"} -->` + "\n" +
		"<!-- /sedum:addSelectColumn -->\n"

	_, err := FindRegions(html, content)
	if err == nil {
		t.Fatal("a region closed by a different label was accepted")
	}
	if !strings.Contains(err.Error(), "addTextColumn") || !strings.Contains(err.Error(), "addSelectColumn") {
		t.Errorf("error does not name both labels: %v", err)
	}
}

// A marker whose comment is not closed still parses. Sedum does not police the
// target's syntax: a package author may have hand-written the marker, and a
// missing delimiter is the target language's complaint to make.
//
// It matters because the alternative is worse. Rejecting it would make an
// author's typo in one template abort a run rather than produce a file their
// own toolchain tells them about.
func TestAnUnclosedDelimitedMarkerStillParses(t *testing.T) {
	got, ok, err := parseOpen(html, `<!-- sedum:addTextColumn {"tier":"owned"}`)
	if err != nil || !ok {
		t.Fatalf("parseOpen = %v, %v", ok, err)
	}
	if got.Action != "addTextColumn" {
		t.Errorf("action = %q", got.Action)
	}
}

// A template-planted anchor declaration is read as not-a-marker in a delimited
// target too, or the first region after one would appear to close it.
func TestADelimitedAnchorDeclarationIsNotAMarker(t *testing.T) {
	for _, line := range []string{
		"<!-- sedum:anchor:columns -->",
		"  <!-- sedum:anchor:columns -->",
	} {
		if _, ok, err := parseOpen(html, line); ok || err != nil {
			t.Errorf("parseOpen(%q) = %v, %v; want not-a-marker and no error", line, ok, err)
		}
	}
}
