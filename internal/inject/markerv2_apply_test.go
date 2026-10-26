package inject

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/genpkg"
)

// prov-2026-8f5373b0: renderRegion wraps content in markers and hashes the body.

func TestRenderRegion(t *testing.T) {
	pkg := &genpkg.Package{Name: "t", CommentPrefix: "#"}
	marker := Marker{Action: "addStep", Variant: "run", Version: 2, Tier: TierOwned, Record: "PR-1",
		Kwargs: map[string]any{"n": "x"}, KwargsHash: "sha256:k"}

	tests := []struct {
		name    string
		content string
		indent  string
		body    string // exactly the text between the marker lines
	}{
		{"content ending in one line feed", "a\nb\n", "", "a\nb\n"},
		{"content with no trailing line feed", "a\nb", "", "a\nb\n"},
		{"content ending in two line feeds keeps one blank line", "a\n\n", "", "a\n\n"},
		{"content ending in many line feeds keeps one blank line", "a\n\n\n\n", "", "a\n\n"},
		{"a blank line inside is kept", "a\n\nb\n", "", "a\n\nb\n"},
		{"empty content has an empty body", "", "", ""},
		{"content of only line feeds has an empty body", "\n\n\n", "", ""},
		{"a single line feed has an empty body", "\n", "", ""},
		{"markers are indented, the body is not", "def a\nend\n", "  ", "def a\nend\n"},
		{"tab indentation", "x\n", "\t\t", "x\n"},
		{"CRLF in content is written as it is", "a\r\nb\r\n", "", "a\r\nb\r\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inv := Invocation{Package: pkg, Content: tc.content}
			wantHash := BodyHash(tc.body)

			written := marker
			written.Body = wantHash
			open, err := written.Open(pkg.Comment())
			if err != nil {
				t.Fatalf("Open: %v", err)
			}
			want := tc.indent + open + "\n" + tc.body + tc.indent + written.Close(pkg.Comment()) + "\n"

			got, hash, err := renderRegion(inv, marker, tc.indent)
			if err != nil {
				t.Fatalf("renderRegion: %v", err)
			}
			if got != want {
				t.Errorf("region =\n%q\nwant\n%q", got, want)
			}
			if hash != wantHash {
				t.Errorf("hash = %s, want the hash of %q (%s)", hash, tc.body, wantHash)
			}
		})
	}
}

func TestRenderRegionWritesTheBodyHashOnAV2Marker(t *testing.T) {
	pkg := &genpkg.Package{Name: "t", CommentPrefix: "#"}
	inv := Invocation{Package: pkg, Content: "policy :read, admin?\n"}

	// Whatever Body the caller's marker held, the one written is the hash of
	// the body written.
	m := Marker{Action: "addPolicy", Version: 2, Tier: TierOwned, Body: "sha256:stale"}
	got, hash, err := renderRegion(inv, m, "")
	if err != nil {
		t.Fatalf("renderRegion: %v", err)
	}
	const wantHash = "sha256:ca235cb34f9574af0f8880305104f597e91d812fca6b884d347fa0775a6f7076"
	if hash != wantHash {
		t.Errorf("hash = %s, want %s", hash, wantHash)
	}
	want := `# sedum:addPolicy {"v":2,"tier":"owned","body":"` + wantHash + `"}` + "\n" +
		"policy :read, admin?\n" +
		"# /sedum:addPolicy\n"
	if got != want {
		t.Errorf("region =\n%s\nwant\n%s", got, want)
	}
}

// The hash is returned whatever the marker's version, and a v1 marker is
// written without a body key.
func TestRenderRegionHashesAV1Marker(t *testing.T) {
	pkg := &genpkg.Package{Name: "t", CommentPrefix: "#"}
	inv := Invocation{Package: pkg, Content: "x\n"}

	got, hash, err := renderRegion(inv, Marker{Action: "addStep", Tier: TierOwned, Body: "sha256:ignored"}, "")
	if err != nil {
		t.Fatalf("renderRegion: %v", err)
	}
	if hash != BodyHash("x\n") {
		t.Errorf("hash = %s, want the hash of the body", hash)
	}
	want := `# sedum:addStep {"tier":"owned"}` + "\nx\n# /sedum:addStep\n"
	if got != want {
		t.Errorf("region =\n%s\nwant\n%s", got, want)
	}
}

func TestRenderRegionUsesThePackagesComment(t *testing.T) {
	pkg := &genpkg.Package{Name: "html", CommentPrefix: "<!--", CommentSuffix: "-->"}
	got, hash, err := renderRegion(Invocation{Package: pkg, Content: "<td>a</td>\n"},
		Marker{Action: "addCol", Variant: "a", Version: 2, Tier: TierOwned}, "  ")
	if err != nil {
		t.Fatalf("renderRegion: %v", err)
	}
	want := `  <!-- sedum:addCol:a {"v":2,"tier":"owned","body":"` + BodyHash("<td>a</td>\n") + `"} -->` + "\n" +
		"<td>a</td>\n" +
		"  <!-- /sedum:addCol:a -->\n"
	if got != want {
		t.Errorf("region =\n%s\nwant\n%s", got, want)
	}
	if hash != BodyHash("<td>a</td>\n") {
		t.Errorf("hash = %s", hash)
	}
}

func TestRenderRegionOpenErrorReturnsNothing(t *testing.T) {
	pkg := &genpkg.Package{Name: "t", CommentPrefix: "#"}
	inv := Invocation{Package: pkg, Content: "x\n"}
	m := Marker{Action: "addStep", Version: 2, Kwargs: map[string]any{"f": func() {}}}

	_, wantErr := m.Open(pkg.Comment())
	if wantErr == nil {
		t.Fatal("Open accepted a kwarg that cannot be encoded")
	}
	got, hash, err := renderRegion(inv, m, "  ")
	if err == nil {
		t.Fatal("renderRegion accepted a marker Open refuses")
	}
	if got != "" || hash != "" {
		t.Errorf("renderRegion = %q, %q with an error, want an empty region and an empty hash", got, hash)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("error = %q, want Open's: %q", err.Error(), wantErr.Error())
	}
}

// prov-2026-dee37812: applyOne injects one invocation and writes a v2 marker.

const applyScaffold = "class UsersController\n" +
	"  # sedum:anchor:class_body_top\n" +
	"\n" +
	"  # sedum:anchor:class_body\n" +
	"end\n"

const (
	policyKwargsHash = "sha256:498b9bece8c809c60f6d4422e9aa2559c0bfbb2409b335021638015aa3f04880"
	policyBodyHash   = "sha256:ca235cb34f9574af0f8880305104f597e91d812fca6b884d347fa0775a6f7076"
)

func policyKwargs(rule string) map[string]any {
	return map[string]any{"controller": "users", "name": "read", "rule": rule}
}

func policyInvocation(t *testing.T, pkg *genpkg.Package, rule string) Invocation {
	t.Helper()
	return invocation(t, pkg, "addPolicy", "", "policy :read, "+rule+"\n", "PR-1", policyKwargs(rule))
}

func policyRegion(pkg *genpkg.Package, extraKeys string, kwargsHash, bodyHash, body string) string {
	return `  # sedum:addPolicy {"v":2,"tier":"owned","record":"PR-1","pkg":"` + pkg.Actions["addPolicy"].PkgHash +
		`","id":{"controller":"users","name":"read"},"kwargs":"` + kwargsHash + `","body":"` + bodyHash + `"` + extraKeys + `}` + "\n" +
		body +
		"  # /sedum:addPolicy\n"
}

func scaffoldWith(region string) string {
	return "class UsersController\n" +
		"  # sedum:anchor:class_body_top\n" +
		"\n" +
		"  # sedum:anchor:class_body\n" +
		region +
		"end\n"
}

func TestApplyOneInsertsAV2Region(t *testing.T) {
	pkg := loadPackage(t)
	inv := policyInvocation(t, pkg, "admin?")

	got, result, err := applyOne(inv, applyScaffold, newOwnership(nil))
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}

	want := scaffoldWith(policyRegion(pkg, "", policyKwargsHash, policyBodyHash, "policy :read, admin?\n"))
	if got != want {
		t.Errorf("content =\n%s\nwant\n%s", got, want)
	}

	if result.Path != controllerPath || result.Action != "addPolicy" || result.Record != "PR-1" {
		t.Errorf("result identity = %q %q %q", result.Path, result.Action, result.Record)
	}
	if result.Replaced || result.Skipped || !result.Changed {
		t.Errorf("Replaced %v Skipped %v Changed %v; want an insertion reported as changed", result.Replaced, result.Skipped, result.Changed)
	}

	m := result.Marker
	if m.Version != 2 {
		t.Errorf("Marker.Version = %d, want 2", m.Version)
	}
	if want := map[string]any{"controller": "users", "name": "read"}; !reflect.DeepEqual(m.Kwargs, want) {
		t.Errorf("Marker.Kwargs = %v, want only the identity kwargs %v", m.Kwargs, want)
	}
	if m.KwargsHash != policyKwargsHash {
		t.Errorf("Marker.KwargsHash = %s, want the hash of every invocation kwarg (%s)", m.KwargsHash, policyKwargsHash)
	}
	if m.Body != policyBodyHash {
		t.Errorf("Marker.Body = %s, want %s", m.Body, policyBodyHash)
	}
	if m.Action != "addPolicy" || m.Record != "PR-1" {
		t.Errorf("Marker = %+v", m)
	}
	if !reflect.DeepEqual(result.Kwargs, inv.Kwargs) {
		t.Errorf("Result.Kwargs = %v, want every resolved kwarg %v", result.Kwargs, inv.Kwargs)
	}
}

// The marker written is what the region now carries: reading it back gives the
// same Marker, and FindRegions reports the body as intact.
func TestApplyOneResultMarkerIsTheMarkerWritten(t *testing.T) {
	pkg := loadPackage(t)
	got, result, err := applyOne(policyInvocation(t, pkg, "admin?"), applyScaffold, newOwnership(nil))
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}
	regions, err := FindRegions(pkg.Comment(), got)
	if err != nil {
		t.Fatalf("FindRegions: %v", err)
	}
	if len(regions) != 1 {
		t.Fatalf("found %d regions, want 1", len(regions))
	}
	r := regions[0]
	if r.Marker.Version != 2 || r.Marker.KwargsHash != result.Marker.KwargsHash || r.Marker.Body != result.Marker.Body {
		t.Errorf("marker read back = %+v, want the one in Result: %+v", r.Marker, result.Marker)
	}
	if r.BodyHash != r.Marker.Body {
		t.Errorf("a region just written has BodyHash %s but its marker says %s", r.BodyHash, r.Marker.Body)
	}
	if state := (OwnedRegion{Region: r}).BodyState(); state != BodyIntact {
		t.Errorf("BodyState = %q, want intact", state)
	}
}

func TestApplyOneRerunChangesNothing(t *testing.T) {
	pkg := loadPackage(t)
	inv := policyInvocation(t, pkg, "admin?")
	own := newOwnership(nil)

	first, _, err := applyOne(inv, applyScaffold, own)
	if err != nil {
		t.Fatalf("first applyOne: %v", err)
	}
	second, result, err := applyOne(inv, first, own)
	if err != nil {
		t.Fatalf("second applyOne: %v", err)
	}
	if second != first {
		t.Errorf("a rerun changed the file:\n%s\nwant\n%s", second, first)
	}
	if !result.Replaced || result.Changed {
		t.Errorf("Replaced %v Changed %v; want a replacement that changed nothing", result.Replaced, result.Changed)
	}
	if result.Marker.Version != 2 || result.Marker.Body != policyBodyHash {
		t.Errorf("Result.Marker = %+v, want the v2 marker written", result.Marker)
	}
}

// A rule is not identity: changing it replaces the region in place, and the
// kwargs hash and the body hash move with it.
func TestApplyOneReplacementMovesTheHashes(t *testing.T) {
	pkg := loadPackage(t)
	own := newOwnership(nil)
	first, _, err := applyOne(policyInvocation(t, pkg, "admin?"), applyScaffold, own)
	if err != nil {
		t.Fatalf("first applyOne: %v", err)
	}

	got, result, err := applyOne(policyInvocation(t, pkg, "owner?"), first, own)
	if err != nil {
		t.Fatalf("second applyOne: %v", err)
	}
	wantKwargs := sha256Of(`{"controller":"users","name":"read","rule":"owner?"}`)
	wantBody := sha256Of("policy :read, owner?\n")
	want := scaffoldWith(policyRegion(pkg, "", wantKwargs, wantBody, "policy :read, owner?\n"))
	if got != want {
		t.Errorf("content =\n%s\nwant\n%s", got, want)
	}
	if !result.Replaced || !result.Changed {
		t.Errorf("Replaced %v Changed %v; want a changed replacement", result.Replaced, result.Changed)
	}
	if result.Marker.KwargsHash != wantKwargs || result.Marker.Body != wantBody {
		t.Errorf("Result.Marker hashes = %s, %s; want %s, %s", result.Marker.KwargsHash, result.Marker.Body, wantKwargs, wantBody)
	}
}

// A hand-edited region is rewritten from the invocation, and that is a change.
func TestApplyOneRestoresAHandEditedBody(t *testing.T) {
	pkg := loadPackage(t)
	own := newOwnership(nil)
	inv := policyInvocation(t, pkg, "admin?")
	first, _, err := applyOne(inv, applyScaffold, own)
	if err != nil {
		t.Fatalf("first applyOne: %v", err)
	}
	edited := strings.Replace(first, "policy :read, admin?\n", "policy :read, admin? # hand edit\n", 1)
	if edited == first {
		t.Fatal("the fixture did not edit the body")
	}

	got, result, err := applyOne(inv, edited, own)
	if err != nil {
		t.Fatalf("second applyOne: %v", err)
	}
	if got != first {
		t.Errorf("content =\n%s\nwant the generated text restored:\n%s", got, first)
	}
	if !result.Replaced || !result.Changed {
		t.Errorf("Replaced %v Changed %v; want a changed replacement", result.Replaced, result.Changed)
	}
}

// A v1 marker being replaced is replaced by the v2 marker, keeping what the
// format does not model and dropping the v1 kwargs object.
func TestApplyOneReplacesAV1Marker(t *testing.T) {
	pkg := loadPackage(t)
	v1 := `  # sedum:addPolicy {"tier":"owned","record":"PR-1","pkg":"` + pkg.Actions["addPolicy"].PkgHash +
		`","kwargs":{"controller":"users","name":"read","rule":"admin?"},"harness_attempts":2,"verified_by":"spec/users.linespec"}` + "\n" +
		"policy :read, admin?\n" +
		"  # /sedum:addPolicy\n"

	got, result, err := applyOne(policyInvocation(t, pkg, "admin?"), scaffoldWith(v1), newOwnership(nil))
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}
	want := scaffoldWith(policyRegion(pkg, `,"harness_attempts":2,"verified_by":"spec/users.linespec"`,
		policyKwargsHash, policyBodyHash, "policy :read, admin?\n"))
	if got != want {
		t.Errorf("content =\n%s\nwant\n%s", got, want)
	}
	if strings.Contains(got, `"kwargs":{`) {
		t.Errorf("the v2 marker still carries the v1 kwargs object:\n%s", got)
	}
	if !result.Replaced || !result.Changed || result.Marker.Version != 2 {
		t.Errorf("Replaced %v Changed %v Version %d; want a changed replacement written as v2",
			result.Replaced, result.Changed, result.Marker.Version)
	}
}

// A marker whose v this version does not read is somebody else's. It is not
// compared with the invocation, not replaced, and not touched.
func TestApplyOneLeavesAForeignRegionAlone(t *testing.T) {
	pkg := loadPackage(t)
	foreign := "  # sedum:addPolicy " +
		`{"v":3,"tier":"owned","record":"PR-9","kwargs":{"controller":"users","name":"read","rule":"admin?"},"id":{"controller":"users","name":"read"}}` + "\n" +
		"  policy :read, from a newer sedum\n" +
		"  # /sedum:addPolicy\n"

	got, result, err := applyOne(policyInvocation(t, pkg, "admin?"), scaffoldWith(foreign), newOwnership(nil))
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}
	if !strings.Contains(got, foreign) {
		t.Errorf("the foreign region was changed:\n%s", got)
	}
	if result.Replaced || result.Skipped || !result.Changed {
		t.Errorf("Replaced %v Skipped %v Changed %v; want the invocation inserted beside the foreign region",
			result.Replaced, result.Skipped, result.Changed)
	}
	regions, err := FindRegions(pkg.Comment(), got)
	if err != nil {
		t.Fatalf("FindRegions: %v", err)
	}
	if len(regions) != 2 {
		t.Fatalf("found %d regions, want the foreign one and the new one", len(regions))
	}
	foreignCount := 0
	for _, r := range regions {
		if r.Marker.Foreign {
			foreignCount++
		}
	}
	if foreignCount != 1 {
		t.Errorf("%d foreign regions, want 1", foreignCount)
	}
}

// An action none of whose kwargs select has an empty identity, which every
// region of that action would share. A foreign region still is not that
// region.
func TestApplyOneNeverComparesAForeignRegionWithAnEmptyIdentity(t *testing.T) {
	pkg := &genpkg.Package{Name: "t", CommentPrefix: "#"}
	action := &genpkg.Action{Name: "addBanner", Anchor: genpkg.AnchorStartOfFile,
		Kwargs: map[string]genpkg.Kwarg{"text": {Type: "string", Required: false}}}
	inv := Invocation{Package: pkg, Action: action, Kwargs: map[string]any{"text": "hi"},
		Path: "f.txt", RecordID: "PR-1", Content: "banner\n"}

	foreign := `# sedum:addBanner {"v":3,"tier":"owned"}` + "\nsomeone else's banner\n# /sedum:addBanner\n"
	got, result, err := applyOne(inv, foreign+"rest\n", newOwnership(nil))
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}
	if !strings.Contains(got, foreign) {
		t.Errorf("the foreign region was changed or replaced:\n%s", got)
	}
	if result.Replaced {
		t.Error("a foreign region was replaced")
	}

	// And a v2 region of the same action, whose identity is equally empty, is
	// the invocation's own region.
	own := newOwnership(nil)
	first, _, err := applyOne(inv, "rest\n", own)
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}
	_, again, err := applyOne(inv, first, own)
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}
	if !again.Replaced || again.Changed {
		t.Errorf("Replaced %v Changed %v; want the v2 region found again and left unchanged", again.Replaced, again.Changed)
	}
}

func TestApplyOneMatchesARegionByItsRecordedID(t *testing.T) {
	pkg := loadPackage(t)
	// A v2 region whose id names the selecting kwargs and whose kwargs is only
	// a hash: the id is what identifies it.
	existing := policyRegion(pkg, "", "sha256:stale", "sha256:stale", "old text\n")
	got, result, err := applyOne(policyInvocation(t, pkg, "admin?"), scaffoldWith(existing), newOwnership(nil))
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}
	if !result.Replaced {
		t.Errorf("the region with the same id was not replaced:\n%s", got)
	}
	if n := strings.Count(got, "# sedum:addPolicy "); n != 1 {
		t.Errorf("file holds %d addPolicy regions, want 1:\n%s", n, got)
	}
}

func TestApplyOneSeededV2RegionIsSkipped(t *testing.T) {
	pkg := loadPackage(t)
	seeded := `  # sedum:addPolicy {"v":2,"tier":"seeded","record":"PR-1","id":{"controller":"users","name":"read"},"kwargs":"sha256:x","body":"sha256:y"}` + "\n" +
		"policy :read, my own\n" +
		"  # /sedum:addPolicy\n"
	content := scaffoldWith(seeded)

	got, result, err := applyOne(policyInvocation(t, pkg, "admin?"), content, newOwnership(nil))
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}
	if got != content {
		t.Errorf("a seeded region was changed:\n%s", got)
	}
	if !result.Skipped || result.Replaced || result.Changed {
		t.Errorf("Skipped %v Replaced %v Changed %v; want skipped and unchanged", result.Skipped, result.Replaced, result.Changed)
	}
	if !reflect.DeepEqual(result.Marker, Marker{}) || result.Kwargs != nil {
		t.Errorf("Marker = %+v, Kwargs = %v; want the zero Marker and nil for a skipped region", result.Marker, result.Kwargs)
	}
}

// A selecting kwarg the invocation does not hold is left out of the id, and the
// hash is still over every kwarg there is.
func TestApplyOneLeavesOutASelectingKwargTheInvocationLacks(t *testing.T) {
	pkg := loadPackage(t)
	kwargs := map[string]any{"controller": "users", "rule": "admin?"}
	inv := invocation(t, pkg, "addPolicy", "", "policy\n", "PR-1", kwargs)

	_, result, err := applyOne(inv, applyScaffold, newOwnership(nil))
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}
	if want := map[string]any{"controller": "users"}; !reflect.DeepEqual(result.Marker.Kwargs, want) {
		t.Errorf("Marker.Kwargs = %v, want %v", result.Marker.Kwargs, want)
	}
	if want := sha256Of(`{"controller":"users","rule":"admin?"}`); result.Marker.KwargsHash != want {
		t.Errorf("KwargsHash = %s, want %s", result.Marker.KwargsHash, want)
	}
}

// A declared identity narrows the id to its kwargs, and a kwarg the action
// declares with a default is hashed whether or not it is id.
func TestApplyOneIdentityIsTheSelectingKwargsAlone(t *testing.T) {
	pkg := loadPackage(t)
	inv := invocation(t, pkg, "addBeforeFilter", "", "before_action :auth\n", "PR-1",
		map[string]any{"controller": "users", "filter": "auth", "only": []any{"show", "edit"}})

	_, result, err := applyOne(inv, applyScaffold, newOwnership(nil))
	if err != nil {
		t.Fatalf("applyOne: %v", err)
	}
	if want := map[string]any{"controller": "users", "filter": "auth"}; !reflect.DeepEqual(result.Marker.Kwargs, want) {
		t.Errorf("Marker.Kwargs = %v, want %v", result.Marker.Kwargs, want)
	}
	if want := sha256Of(`{"controller":"users","filter":"auth","only":["show","edit"]}`); result.Marker.KwargsHash != want {
		t.Errorf("KwargsHash = %s, want %s", result.Marker.KwargsHash, want)
	}
}

func TestApplyOneKwargsThatCannotBeHashedWriteNothing(t *testing.T) {
	pkg := loadPackage(t)
	kwargs := policyKwargs("admin?")
	kwargs["rule"] = func() {}
	inv := invocation(t, pkg, "addPolicy", "", "policy\n", "PR-1", kwargs)

	got, _, err := applyOne(inv, applyScaffold, newOwnership(nil))
	if err == nil {
		t.Fatalf("applyOne accepted kwargs that cannot be hashed:\n%s", got)
	}
	const want = "action addPolicy: kwargs cannot be recorded on a marker: kwargs cannot be hashed: json: unsupported type"
	if !strings.HasPrefix(err.Error(), want) {
		t.Errorf("error = %q, want it to start with %q", err.Error(), want)
	}
	if strings.Contains(got, "sedum:addPolicy") {
		t.Errorf("content written despite the error:\n%s", got)
	}
}

// A variant is part of the label the error names.
func TestApplyOneKwargsErrorNamesTheLabel(t *testing.T) {
	pkg := loadPackage(t)
	inv := invocation(t, pkg, "createControllerMethod", "index", "def index\nend\n", "PR-1",
		map[string]any{"controller": "users", "name": "index", "collection": func() {}})

	_, _, err := applyOne(inv, applyScaffold, newOwnership(nil))
	if err == nil {
		t.Fatal("applyOne accepted kwargs that cannot be hashed")
	}
	if !strings.HasPrefix(err.Error(), "action createControllerMethod:index: kwargs cannot be recorded on a marker: ") {
		t.Errorf("error = %q", err.Error())
	}
}

// The v2 marker written by Apply is what a file holds after a run.
func TestApplyWritesV2Markers(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)
	results, err := Apply([]Invocation{policyInvocation(t, pkg, "admin?")}, Options{Output: root})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	want := scaffoldWith(policyRegion(pkg, "", policyKwargsHash, policyBodyHash, "policy :read, admin?\n"))
	if got := read(t, root); got != want {
		t.Errorf("file =\n%s\nwant\n%s", got, want)
	}
	if len(results) != 1 || results[0].Marker.Version != 2 || results[0].Marker.Body != policyBodyHash {
		t.Errorf("results = %+v, want one result carrying the v2 marker written", results)
	}
	data, err := json.Marshal(results[0].Marker.Kwargs)
	if err != nil || string(data) != `{"controller":"users","name":"read"}` {
		t.Errorf("Result.Marker.Kwargs = %s, %v", data, err)
	}
}
