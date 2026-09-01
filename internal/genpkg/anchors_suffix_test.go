package genpkg

import (
	"regexp"
	"testing"
)

// Anchor declarations in a target whose comments are delimited
// (prov-2026-a6f6bb81).

func TestMarkersInReadsADelimitedDeclaration(t *testing.T) {
	html := Comment{Prefix: "<!--", Suffix: "-->"}
	template := "<app-grid>\n" +
		"  <!-- sedum:anchor:columns -->\n" +
		"</app-grid>\n" +
		"<!-- sedum:anchor:templates -->\n"

	got := MarkersIn(html, template)
	if len(got) != 2 || got[0] != "columns" || got[1] != "templates" {
		t.Fatalf("MarkersIn = %v, want [columns templates]", got)
	}
}

// MarkersIn needs no suffix handling because its trailing character class
// cannot match the space before one. That is convenient rather than designed,
// so it is pinned: widening the class to include whitespace would silently
// start swallowing " -->" into an anchor name, and every affected package would
// fail at injection with a marker that "is not in the file" rather than here.
//
// The assertion is on the pattern's behaviour rather than on its text, so a
// rewrite that keeps the property passes and one that loses it does not.
func TestMarkerNameClassBoundsTheSuffix(t *testing.T) {
	html := Comment{Prefix: "<!--", Suffix: "-->"}

	// Every character a suffix could begin with, after the separating space.
	for _, suffix := range []string{"-->", "*/", "#}", "%}", "]]>"} {
		comment := Comment{Prefix: "<!--", Suffix: suffix}
		got := MarkersIn(comment, comment.Wrap("sedum:anchor:columns"))
		if len(got) != 1 || got[0] != "columns" {
			t.Errorf("suffix %q: MarkersIn = %v, want [columns]", suffix, got)
		}
	}

	// And the property itself: the name pattern must not admit a space, or the
	// bound above is accidental and the next edit can remove it.
	name := regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)
	if name.MatchString("columns -->") {
		t.Error("the marker-name character class admits whitespace; a closing delimiter would be read as part of the name")
	}

	// A declaration with no suffix in a package that declares one is still
	// read. The name ends at end-of-line rather than at the delimiter.
	if got := MarkersIn(html, "<!-- sedum:anchor:columns"); len(got) != 1 || got[0] != "columns" {
		t.Errorf("MarkersIn on an unclosed declaration = %v, want [columns]", got)
	}
}

// MissingMarkers compares a template against a file through the same shape, so
// a delimited package's Phase 3 check does not report every anchor missing.
func TestMissingMarkersThroughADelimitedComment(t *testing.T) {
	html := Comment{Prefix: "<!--", Suffix: "-->"}
	template := "<!-- sedum:anchor:columns -->\n<!-- sedum:anchor:templates -->\n"

	if got := MissingMarkers(html, template, template); len(got) != 0 {
		t.Errorf("a file carrying its template's markers reports %v missing", got)
	}

	partial := "<!-- sedum:anchor:columns -->\n"
	got := MissingMarkers(html, template, partial)
	if len(got) != 1 || got[0] != "templates" {
		t.Errorf("MissingMarkers = %v, want [templates]", got)
	}
}

// A suffix with nothing to close is a package error rather than a default.
func TestCommentWrapAndUnwrapAreInverse(t *testing.T) {
	for _, comment := range []Comment{
		{Prefix: "#"},
		{Prefix: "//"},
		{Prefix: "<!--", Suffix: "-->"},
		{Prefix: "/*", Suffix: "*/"},
	} {
		const body = "sedum:anchor:columns"

		wrapped := comment.Wrap(body)
		// Unwrap operates on text with the prefix already removed, which is
		// how every reader reaches it.
		stripped := wrapped[len(comment.Prefix):]
		if got := comment.Unwrap(stripped); got != body {
			t.Errorf("comment %+v: Unwrap(Wrap(%q)) = %q", comment, body, got)
		}
		if comment.Delimited() != (comment.Suffix != "") {
			t.Errorf("comment %+v: Delimited disagrees with Suffix", comment)
		}
	}
}
