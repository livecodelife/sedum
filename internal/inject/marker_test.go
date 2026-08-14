package inject

import (
	"encoding/json"
	"strings"
	"testing"
)

// The marker is the one durable artifact this milestone writes. It sits in a
// generated codebase long after the version that wrote it is gone, so the cases
// below are about what a marker means to a reader that is not the writer:
// what identifies a region, what is merely recorded on it, and what a parser
// does with a field it has never heard of (prov-2026-36c8a99c).

func TestMarkerRoundTrips(t *testing.T) {
	want := Marker{
		Action:  "createControllerMethod",
		Variant: "index",
		Tier:    TierOwned,
		Record:  "PR-014",
		Kwargs:  map[string]any{"controller": "users", "collection": "users"},
	}

	open, err := want.Open("#")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	// The action and variant stay literal on the line: grepping markers is
	// the audit trail, and it only works if it is grep and not a parser.
	if !strings.HasPrefix(open, "# sedum:createControllerMethod:index ") {
		t.Errorf("opening marker does not lead with a greppable action:variant label:\n%s", open)
	}

	got, ok, err := parseOpen("#", open)
	if err != nil {
		t.Fatalf("parseOpen: %v", err)
	}
	if !ok {
		t.Fatalf("parseOpen did not recognize the marker it wrote:\n%s", open)
	}

	if got.Action != want.Action || got.Variant != want.Variant {
		t.Errorf("label round-trip = %s:%s, want %s:%s", got.Action, got.Variant, want.Action, want.Variant)
	}
	if got.Tier != want.Tier {
		t.Errorf("tier round-trip = %q, want %q", got.Tier, want.Tier)
	}
	if got.Record != want.Record {
		t.Errorf("record round-trip = %q, want %q", got.Record, want.Record)
	}
	if got.Kwargs["controller"] != "users" || got.Kwargs["collection"] != "users" {
		t.Errorf("kwargs round-trip = %v, want the two it was rendered from", got.Kwargs)
	}

	if close := want.Close("#"); close != "# /sedum:createControllerMethod:index" {
		t.Errorf("closing marker = %q, want it to name the same action and variant", close)
	}
}

// A marker records what the region was rendered from, verbatim. encoding/json
// escapes &, < and > by default, which would record a Go address-of expression
// or a C++ template argument as a run of escapes on a line meant to be read.
func TestKwargsAreRecordedUnescaped(t *testing.T) {
	m := Marker{Action: "createQuery", Variant: "insert", Kwargs: map[string]any{
		"scan_targets": "&t.ID, &t.Title",
		"bound":        "map[string]Handler<T>",
	}}

	open, err := m.Open("//")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	for _, want := range []string{`&t.ID, &t.Title`, `map[string]Handler<T>`} {
		if !strings.Contains(open, want) {
			t.Errorf("marker does not carry %q verbatim:\n%s", want, open)
		}
	}
	for _, escape := range []string{`\u0026`, `\u003c`, `\u003e`} {
		if strings.Contains(open, escape) {
			t.Errorf("marker carries the HTML escape %s:\n%s", escape, open)
		}
	}

	// The marker still has to be one line, and still has to read back.
	if strings.Contains(open, "\n") {
		t.Errorf("marker spans more than one line:\n%s", open)
	}
	parsed, ok, err := parseOpen("//", open)
	if err != nil || !ok {
		t.Fatalf("parseOpen = %v, %v", ok, err)
	}
	if parsed.Kwargs["scan_targets"] != "&t.ID, &t.Title" {
		t.Errorf("kwargs did not round-trip: %v", parsed.Kwargs)
	}
}

// A simple action has no variant, so its label is the action alone rather than
// an action with an empty variant hanging off it.
func TestMarkerWithoutVariant(t *testing.T) {
	m := Marker{Action: "addBeforeFilter", Kwargs: map[string]any{"filter": "authenticate"}}

	open, err := m.Open("#")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if !strings.HasPrefix(open, "# sedum:addBeforeFilter {") {
		t.Errorf("opening marker = %q, want no empty variant on the label", open)
	}
	if got := m.Close("#"); got != "# /sedum:addBeforeFilter" {
		t.Errorf("closing marker = %q, want no empty variant on the label", got)
	}

	parsed, ok, err := parseOpen("#", open)
	if err != nil || !ok {
		t.Fatalf("parseOpen(%q) = %v, %v", open, ok, err)
	}
	if parsed.Variant != "" {
		t.Errorf("variant = %q, want empty", parsed.Variant)
	}
}

// The prefix is the package's declared comment_prefix and is never hardcoded
// and never inferred from a file extension.
func TestMarkerUsesDeclaredCommentPrefix(t *testing.T) {
	for _, prefix := range []string{"#", "//", "--", ";;"} {
		m := Marker{Action: "addStep", Kwargs: map[string]any{"step": "build"}}

		open, err := m.Open(prefix)
		if err != nil {
			t.Fatalf("Open(%q): %v", prefix, err)
		}
		if !strings.HasPrefix(open, prefix+" sedum:") {
			t.Errorf("prefix %q: opening marker = %q", prefix, open)
		}

		// A marker written with one prefix is not a marker to a package
		// that declares another.
		if _, ok, _ := parseOpen("@@", open); ok {
			t.Errorf("prefix %q: marker was recognized under an unrelated comment prefix", prefix)
		}
	}
}

// The parser tolerates fields it does not recognize and defaults fields that are
// absent. The version that wrote a marker and the version that reads it are
// routinely different, so a codebase generated by a newer Sedum has to stay
// readable by an older one, and a field added later has to be an addition
// rather than a migration across every repository that already carries markers.
func TestParserIsLenient(t *testing.T) {
	cases := []struct {
		name       string
		line       string
		wantTier   Tier
		wantRecord string
	}{{
		name:       "a field this version has never heard of is ignored",
		line:       `# sedum:createControllerMethod:index {"tier":"seeded","record":"PR-9","exposes":["User#index"],"consumes":[]}`,
		wantTier:   TierSeeded,
		wantRecord: "PR-9",
	}, {
		name:       "an absent tier takes its default",
		line:       `# sedum:createControllerMethod:index {"record":"PR-9"}`,
		wantTier:   DefaultTier,
		wantRecord: "PR-9",
	}, {
		name:       "an absent record is empty rather than an error",
		line:       `# sedum:createControllerMethod:index {"tier":"owned"}`,
		wantTier:   TierOwned,
		wantRecord: "",
	}, {
		name:     "an attribute object that is absent entirely defaults every field",
		line:     `# sedum:createControllerMethod:index`,
		wantTier: DefaultTier,
	}, {
		name:     "an indented marker is still a marker",
		line:     `    # sedum:createControllerMethod:index {"tier":"owned"}`,
		wantTier: TierOwned,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok, err := parseOpen("#", tc.line)
			if err != nil {
				t.Fatalf("parseOpen: %v", err)
			}
			if !ok {
				t.Fatalf("line was not recognized as a marker: %s", tc.line)
			}
			if got.Tier != tc.wantTier {
				t.Errorf("tier = %q, want %q", got.Tier, tc.wantTier)
			}
			if got.Record != tc.wantRecord {
				t.Errorf("record = %q, want %q", got.Record, tc.wantRecord)
			}
		})
	}
}

// Ignoring an unknown key on read and dropping it on write are different
// promises, and only the first one comes free. The marker is a public format:
// another tool annotates a region with its own state and needs that state to
// still be there after Sedum rewrites the region (prov-2026-72775ae5).
func TestUnrecognisedAttributesSurviveTheRoundTrip(t *testing.T) {
	const line = `# sedum:createControllerMethod:index ` +
		`{"tier":"owned","record":"PR-9","kwargs":{"controller":"users"},` +
		`"harness_attempts":3,"verified_by":"spec/users.linespec"}`

	parsed, ok, err := parseOpen("#", line)
	if err != nil {
		t.Fatalf("parseOpen: %v", err)
	}
	if !ok {
		t.Fatalf("line was not recognized as a marker: %s", line)
	}

	if got := string(parsed.Extra["harness_attempts"]); got != "3" {
		t.Errorf("harness_attempts retained as %q, want 3", got)
	}
	if got := string(parsed.Extra["verified_by"]); got != `"spec/users.linespec"` {
		t.Errorf("verified_by retained as %q, want the spec path it was written with", got)
	}

	// A key Sedum models is bound to its field and never lands in Extra,
	// where it would be written a second time.
	for _, declared := range []string{"tier", "record", "kwargs", "writer"} {
		if _, carried := parsed.Extra[declared]; carried {
			t.Errorf("declared key %q was also carried as an unrecognised one", declared)
		}
	}

	written, err := parsed.Open("#")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	for _, want := range []string{`"harness_attempts":3`, `"verified_by":"spec/users.linespec"`} {
		if !strings.Contains(written, want) {
			t.Errorf("rewritten marker dropped %s:\n%s", want, written)
		}
	}

	// The second trip has to be stable too, or a rerun would churn the file
	// without changing what the marker says.
	again, ok, err := parseOpen("#", written)
	if err != nil || !ok {
		t.Fatalf("parseOpen on the rewritten marker = %v, %v", ok, err)
	}
	stable, err := again.Open("#")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if stable != written {
		t.Errorf("rewriting a marker twice is not stable:\n%s\n%s", written, stable)
	}
}

// Sedum does not know what a carried key means, which is the point of it. It is
// written back exactly as it was read rather than re-encoded, and the keys
// Sedum does model keep the order they have always been written in, so that
// introducing preservation rewrites no marker that has none.
func TestCarriedAttributesAreWrittenVerbatimAndLast(t *testing.T) {
	parsed, _, err := parseOpen("#", `# sedum:createControllerMethod:index `+
		`{"tier":"owned","zeta":1,"alpha":{"b":2,"a":1}}`)
	if err != nil {
		t.Fatalf("parseOpen: %v", err)
	}

	written, err := parsed.Open("#")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	const want = `# sedum:createControllerMethod:index {"tier":"owned","alpha":{"b":2,"a":1},"zeta":1}`
	if written != want {
		t.Errorf("marker written as\n  %s\nwant\n  %s", written, want)
	}
}

// A carried key may not shadow one Sedum models. It cannot arise from a marker
// Sedum parsed, but Marker is an ordinary struct any caller may build, and an
// attribute object naming tier twice is worse than one naming it once.
func TestCarriedAttributesCannotShadowDeclaredOnes(t *testing.T) {
	m := Marker{
		Action: "createControllerMethod",
		Tier:   TierOwned,
		Extra: map[string]json.RawMessage{
			"tier":   json.RawMessage(`"seeded"`),
			"kept":   json.RawMessage(`true`),
			"kwargs": json.RawMessage(`{"controller":"impostor"}`),
		},
	}

	written, err := m.Open("#")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if strings.Count(written, `"tier"`) != 1 {
		t.Errorf("tier appears more than once in the attribute object:\n%s", written)
	}
	if strings.Contains(written, "impostor") {
		t.Errorf("a carried kwargs key overwrote the modelled one:\n%s", written)
	}
	if !strings.Contains(written, `"kept":true`) {
		t.Errorf("a carried key that shadows nothing was dropped:\n%s", written)
	}

	parsed, _, err := parseOpen("#", written)
	if err != nil {
		t.Fatalf("parseOpen: %v", err)
	}
	if parsed.Tier != TierOwned {
		t.Errorf("tier = %q, want the modelled value rather than the shadowing one", parsed.Tier)
	}
}

// The writer names the tool that last wrote a region, so that a demoted tier is
// attributable: seeded alone cannot say whether a package author declared it or
// another tool demoted it. Absent means Sedum, and Sedum omits the key, so
// every marker written before the field existed reads correctly and none of
// them is rewritten by its introduction.
func TestWriterDefaultsToSedumAndIsOmittedWhenItIs(t *testing.T) {
	written, err := Marker{Action: "createControllerMethod", Tier: TierOwned}.Open("#")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if strings.Contains(written, "writer") {
		t.Errorf("Sedum wrote a writer key rather than omitting it:\n%s", written)
	}

	parsed, _, err := parseOpen("#", written)
	if err != nil {
		t.Fatalf("parseOpen: %v", err)
	}
	if parsed.Writer != DefaultWriter {
		t.Errorf("writer = %q, want %q when the key is absent", parsed.Writer, DefaultWriter)
	}

	// A marker read and written back unchanged stays unchanged, which is what
	// makes the default safe to fill in on read.
	again, err := parsed.Open("#")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if again != written {
		t.Errorf("filling in the default writer changed the marker:\n%s\n%s", written, again)
	}

	foreign, _, err := parseOpen("#", `# sedum:createControllerMethod:index {"tier":"seeded","writer":"harness"}`)
	if err != nil {
		t.Fatalf("parseOpen: %v", err)
	}
	if foreign.Writer != "harness" {
		t.Errorf("writer = %q, want the tool that named itself", foreign.Writer)
	}
}

// Neither the writer nor a carried attribute may take part in identity. If one
// could, a tool annotating a region would silently split it in two: the next
// run would not recognize the region as the one its invocation owns, and would
// inject a second copy beside it.
func TestWriterAndCarriedAttributesAreNotPartOfIdentity(t *testing.T) {
	plain := Marker{
		Action:  "createControllerMethod",
		Variant: "index",
		Kwargs:  map[string]any{"controller": "users"},
	}
	annotated := plain
	annotated.Writer = "harness"
	annotated.Tier = TierSeeded
	annotated.Record = "PR-99"
	annotated.Extra = map[string]json.RawMessage{"harness_attempts": json.RawMessage(`4`)}

	selecting := []string{"controller"}
	a, err := IdentityOf(plain, selecting)
	if err != nil {
		t.Fatalf("IdentityOf: %v", err)
	}
	b, err := IdentityOf(annotated, selecting)
	if err != nil {
		t.Fatalf("IdentityOf: %v", err)
	}
	if a != b {
		t.Errorf("annotating a region changed its identity:\n%+v\n%+v", a, b)
	}
}

// Leniency is about version skew, not about corruption. An attribute object
// that is not JSON at all is a broken marker rather than a marker from the
// future, because the object's encoding is the committed shape.
func TestUnreadableAttributesAreAnError(t *testing.T) {
	_, _, err := parseOpen("#", `# sedum:createControllerMethod:index {"tier":`)
	if err == nil {
		t.Fatal("a marker carrying unreadable attributes was accepted")
	}
	if !strings.Contains(err.Error(), "createControllerMethod:index") {
		t.Errorf("error does not name the marker it could not read: %v", err)
	}
}

// A line that is not a marker is not an error. Most lines in a generated file
// are not markers, including comments that merely start the same way.
func TestNonMarkerLines(t *testing.T) {
	for _, line := range []string{
		"def index",
		"# a comment",
		"# sedum is a generator",
		"",
	} {
		if _, ok, err := parseOpen("#", line); ok || err != nil {
			t.Errorf("parseOpen(%q) = %v, %v; want not-a-marker and no error", line, ok, err)
		}
	}
}

// The record ID is written on the marker and takes no part in matching.
//
// Under record-scoped identity a later record could never claim a region an
// earlier one wrote: it would always mint a second region beside the first, so
// a record whose intent is "PUT should support partial updates" would produce
// two definitions of one method rather than a replacement.
func TestRecordIDIsNotPartOfIdentity(t *testing.T) {
	selecting := []string{"controller", "name"}

	earlier := Marker{
		Action: "createControllerMethod", Variant: "index", Record: "PR-014",
		Kwargs: map[string]any{"controller": "users", "name": "index", "collection": "users"},
	}
	later := Marker{
		Action: "createControllerMethod", Variant: "index", Record: "PR-092",
		Kwargs: map[string]any{"controller": "users", "name": "index", "collection": "admins"},
	}

	a, err := IdentityOf(earlier, selecting)
	if err != nil {
		t.Fatalf("IdentityOf: %v", err)
	}
	b, err := IdentityOf(later, selecting)
	if err != nil {
		t.Fatalf("IdentityOf: %v", err)
	}

	if a != b {
		t.Errorf("two records invoking the same action differ in identity:\n  %+v\n  %+v", a, b)
	}
}

// A kwarg that selects what a region targets makes two regions distinct; a
// kwarg that merely parameterizes one does not.
func TestSelectingKwargsDistinguishRegions(t *testing.T) {
	selecting := []string{"controller", "filter"}

	base := func(filter, only string) Marker {
		return Marker{
			Action: "addBeforeFilter",
			Kwargs: map[string]any{"controller": "users", "filter": filter, "only": only},
		}
	}

	authenticate, err := IdentityOf(base("authenticate", "index"), selecting)
	if err != nil {
		t.Fatalf("IdentityOf: %v", err)
	}
	authorize, err := IdentityOf(base("authorize", "index"), selecting)
	if err != nil {
		t.Fatalf("IdentityOf: %v", err)
	}
	reparameterized, err := IdentityOf(base("authenticate", "show"), selecting)
	if err != nil {
		t.Fatalf("IdentityOf: %v", err)
	}

	if authenticate == authorize {
		t.Error("two different before_action filters share one identity; the second would overwrite the first")
	}
	if authenticate != reparameterized {
		t.Error("changing a non-selecting kwarg changed the identity; the update would append beside the region instead of replacing it")
	}
}

// Identity compares what a value is, not what it was built as. A fixture's
// []string and the []any the same list parses back as have to match, or a
// region would fail to recognize itself on the run after it was written.
func TestIdentityNormalizesValueShapes(t *testing.T) {
	authored := Marker{Action: "addBeforeFilter", Kwargs: map[string]any{
		"filter": "authenticate", "only": []string{"index", "show"},
	}}

	open, err := authored.Open("#")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	parsed, ok, err := parseOpen("#", open)
	if err != nil || !ok {
		t.Fatalf("parseOpen = %v, %v", ok, err)
	}

	selecting := []string{"filter", "only"}
	a, err := IdentityOf(authored, selecting)
	if err != nil {
		t.Fatalf("IdentityOf: %v", err)
	}
	b, err := IdentityOf(parsed, selecting)
	if err != nil {
		t.Fatalf("IdentityOf: %v", err)
	}

	if a != b {
		t.Errorf("a region does not recognize itself after a round-trip:\n  %+v\n  %+v", a, b)
	}
}

func TestFindRegions(t *testing.T) {
	content := "class UsersController\n" +
		"  # sedum:anchor:class_body\n" +
		`# sedum:createControllerMethod:index {"tier":"owned","kwargs":{"name":"index"}}` + "\n" +
		"def index\nend\n" +
		"# /sedum:createControllerMethod:index\n" +
		`# sedum:createControllerMethod:show {"tier":"seeded","kwargs":{"name":"show"}}` + "\n" +
		"def show\nend\n" +
		"# /sedum:createControllerMethod:show\n" +
		"end\n"

	regions, err := FindRegions("#", content)
	if err != nil {
		t.Fatalf("FindRegions: %v", err)
	}
	if len(regions) != 2 {
		t.Fatalf("found %d regions, want 2", len(regions))
	}

	// An anchor marker is a point, not a region, and must not be mistaken
	// for one.
	if regions[0].Marker.Variant != "index" || regions[1].Marker.Variant != "show" {
		t.Errorf("regions = %q, %q; want index then show in file order",
			regions[0].Marker.Variant, regions[1].Marker.Variant)
	}
	if regions[1].Marker.Tier != TierSeeded {
		t.Errorf("second region tier = %q, want seeded", regions[1].Marker.Tier)
	}

	// The extent has to cover the markers themselves, or replacing a region
	// would leave its old markers behind.
	first := content[regions[0].Start:regions[0].End]
	if !strings.HasPrefix(first, "# sedum:createControllerMethod:index") {
		t.Errorf("region does not start at its opening marker:\n%s", first)
	}
	if !strings.HasSuffix(first, "# /sedum:createControllerMethod:index\n") {
		t.Errorf("region does not end at its closing marker:\n%s", first)
	}
}

// A region whose extent is unknown cannot be replaced: doing so would either
// destroy the rest of the file or append beside it, and both are worse than
// halting.
func TestUnterminatedRegionIsAnError(t *testing.T) {
	content := `# sedum:createControllerMethod:index {"tier":"owned"}` + "\ndef index\nend\n"

	if _, err := FindRegions("#", content); err == nil {
		t.Fatal("an unterminated region was accepted")
	}
}

func TestMismatchedClosingMarkerIsAnError(t *testing.T) {
	content := `# sedum:createControllerMethod:index {}` + "\ndef index\nend\n" +
		"# /sedum:createControllerMethod:show\n"

	_, err := FindRegions("#", content)
	if err == nil {
		t.Fatal("a region closed by a marker naming a different variant was accepted")
	}
	if !strings.Contains(err.Error(), "index") || !strings.Contains(err.Error(), "show") {
		t.Errorf("error does not name both markers: %v", err)
	}
}

// An anchor a file template planted is not an ownership marker.
//
// The two vocabularies share the "sedum:" namespace, so a reader that does not
// tell them apart sees "sedum:anchor:class_body" as an ownership marker for an
// action named "anchor", and then reads the first region that follows as the
// one closing it.
func TestAnchorDeclarationsAreNotOwnershipMarkers(t *testing.T) {
	for _, line := range []string{
		"# sedum:anchor:class_body",
		"  # sedum:anchor:class_body_top",
	} {
		if _, ok, err := parseOpen("#", line); ok || err != nil {
			t.Errorf("parseOpen(%q) = %v, %v; want an anchor declaration to be read as not-a-marker",
				line, ok, err)
		}
	}
}
