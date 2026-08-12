package genpkg

import "regexp"

// Anchors: where in a file an action's output belongs.
//
// anchor is a scalar; the values below are reserved and every other value names
// a marker planted by a file template (prov-2026-d1d61186). Marker syntax uses
// the package's declared comment prefix, since #, // and -- all appear across
// targets.
const (
	AnchorStartOfFile = "start_of_file"
	AnchorEndOfFile   = "end_of_file"
	AnchorRegion      = "region"
	AnchorAfterMatch  = "after_match"
	AnchorBeforeMatch = "before_match"
	// AnchorMarker is reserved but not usable bare: an action anchored to
	// "marker" has not said which marker.
	AnchorMarker = "marker"
)

func reservedAnchor(anchor string) bool {
	switch anchor {
	case AnchorStartOfFile, AnchorEndOfFile, AnchorRegion, AnchorAfterMatch, AnchorBeforeMatch, AnchorMarker:
		return true
	}
	return false
}

// MarkerAnchor returns the marker name an action is anchored to, and whether it
// is anchored to a marker at all. A reserved keyword is not a marker name.
func (a *Action) MarkerAnchor() (string, bool) {
	if a.Anchor == "" || reservedAnchor(a.Anchor) {
		return "", false
	}
	return a.Anchor, true
}

// markerDecl is the text a file template plants to create an anchor point. The
// prefix is the package's declared comment_prefix, since #, // and -- all
// appear across targets.
func markerDecl(commentPrefix, name string) string {
	return commentPrefix + " sedum:anchor:" + name
}

// plantedMarkers returns every marker name planted across the package's file
// template contents.
func plantedMarkers(commentPrefix string, contents []string) map[string]bool {
	// The prefix is author-supplied, so it is escaped rather than
	// interpolated into a pattern.
	re := regexp.MustCompile(regexp.QuoteMeta(markerDecl(commentPrefix, "")) + `([A-Za-z0-9_.-]+)`)

	out := map[string]bool{}
	for _, c := range contents {
		for _, m := range re.FindAllStringSubmatch(c, -1) {
			out[m[1]] = true
		}
	}
	return out
}
