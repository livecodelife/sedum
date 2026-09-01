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

// TargetedMarkers returns every marker this action can inject at, which is one
// for a marker anchor, two for a region anchor, and none for the rest.
//
// A region names its endpoints through anchor_start and anchor_end rather than
// anchor, so reading only MarkerAnchor would report every region's endpoints as
// targeted by nothing.
//
// It is one definition rather than two because two questions turn on it and
// they must not be able to disagree: whether a package plants a marker nothing
// targets, asked once at load (prov-2026-a9e59197), and whether an unfilled
// anchor is fillable at all, asked per run (prov-2026-206fa618). A run that
// asked the model about an anchor load had already called dead would be exactly
// that disagreement, charged per sample.
func (a *Action) TargetedMarkers() []string {
	if a.Anchor == AnchorRegion {
		var out []string
		if a.AnchorStart != "" {
			out = append(out, a.AnchorStart)
		}
		if a.AnchorEnd != "" {
			out = append(out, a.AnchorEnd)
		}
		return out
	}
	if marker, ok := a.MarkerAnchor(); ok {
		return []string{marker}
	}
	return nil
}

// markerDecl is the text a file template plants to create an anchor point. The
// delimiters are the package's declared ones, since #, // and -- all appear
// across targets and some targets close their comments.
func markerDecl(comment Comment, name string) string {
	return comment.Wrap("sedum:anchor:" + name)
}

// markerDeclOpen is everything a planted declaration carries before the anchor
// name: the comment's opening delimiter and the keyword.
//
// It is separate from markerDecl because the two uses pull in opposite
// directions. Writing a declaration, or naming its shape in a diagnostic, wants
// the whole line including any closing delimiter. Searching for one wants only
// what comes before the name - markerDecl(comment, "") would put the closing
// delimiter immediately where the name belongs, and match nothing.
func markerDeclOpen(comment Comment) string {
	return comment.Prefix + " sedum:anchor:"
}

// MarkersIn returns the marker names planted in one piece of template text,
// in the order they appear, without repeats.
//
// It is exported because Phase 3 reads it back: a file that already exists is
// checked for the markers its template plants rather than re-rendered over. The
// marker's shape is declared here and nowhere else, so that a change to it
// cannot leave the writer and the reader disagreeing.
// The trailing character class is what bounds the name, and it is why a closing
// delimiter needs no handling here: it cannot match a space, so " -->" ends the
// capture rather than joining it. That is convenient rather than designed, and
// TestMarkerNameClassBoundsTheSuffix pins it so that widening the class cannot
// quietly start swallowing a suffix into an anchor name.
func MarkersIn(comment Comment, content string) []string {
	// The delimiter is author-supplied, so it is escaped rather than
	// interpolated into a pattern.
	re := regexp.MustCompile(regexp.QuoteMeta(markerDeclOpen(comment)) + `([A-Za-z0-9_.-]+)`)

	var out []string
	seen := map[string]bool{}
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	return out
}

// MissingMarkers returns the markers a template plants that content does not
// carry, in template order.
//
// Phase 3 asks this of a file that already exists. A file missing the markers
// its template declares was written by something other than Sedum, or its
// template changed shape after it was generated; either way the injections
// aimed at it have nowhere to land.
func MissingMarkers(comment Comment, template, content string) []string {
	present := map[string]bool{}
	for _, name := range MarkersIn(comment, content) {
		present[name] = true
	}

	var out []string
	for _, name := range MarkersIn(comment, template) {
		if !present[name] {
			out = append(out, name)
		}
	}
	return out
}

// plantedMarkers returns every marker name planted across the package's file
// template contents.
func plantedMarkers(comment Comment, contents []string) map[string]bool {
	out := map[string]bool{}
	for _, c := range contents {
		for _, name := range MarkersIn(comment, c) {
			out[name] = true
		}
	}
	return out
}
