package pipeline

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/record"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/resolve"
)

// sedum markers migrate (prov-2026-0a398351): every v1 marker it can is
// rewritten to v2, with no model, no change to any region's text, and the
// recordings of the records the markers name. Hash literals were derived with
// a throwaway script (see recordings_test.go); those for the typed action are
//
//	kwargs {"field":"age","kind":"integer","model":"user"}  sha256:259fa590...d86ca8
//	body   "attribute :age, :integer\n"                      sha256:09435243...826938
const (
	hashKwargsAge = "sha256:259fa5909b80cc27ad8edde277f67a11e37a70d86daf015bd79bfd78fcd86ca8"
	hashBodyAge   = "sha256:0943524360de89a64cda46747ac536060a2f1d38de4e1252497b558bfb826938"

	migGen = "sha256:gen-from-the-marker"
	migPkg = "sha256:pkg-from-the-marker"
)

// Types.

func TestMigrateConfigHasTheFieldsTheRecordNames(t *testing.T) {
	cfg := MigrateConfig{
		Generators: "g", Records: "r", Output: "o",
		Variables: map[string]string{"a": "b"},
		DryRun:    true,
	}
	if cfg.Generators != "g" || cfg.Records != "r" || cfg.Output != "o" || cfg.Variables["a"] != "b" || !cfg.DryRun {
		t.Errorf("cfg = %+v", cfg)
	}
	// No field for excluded directories: Migrate passes none.
	if n := reflect.TypeOf(cfg).NumField(); n != 5 {
		t.Errorf("MigrateConfig has %d fields, want exactly Generators, Records, Output, Variables and DryRun", n)
	}
	for name, want := range map[string]reflect.Type{
		"Generators": reflect.TypeOf(""), "Records": reflect.TypeOf(""), "Output": reflect.TypeOf(""),
		"Variables": reflect.TypeOf(map[string]string{}), "DryRun": reflect.TypeOf(true),
	} {
		if f, ok := reflect.TypeOf(cfg).FieldByName(name); !ok || f.Type != want {
			t.Errorf("MigrateConfig.%s = %v, want %v", name, f.Type, want)
		}
	}
}

func TestMigrateNoteHasTheFieldsTheRecordNames(t *testing.T) {
	n := MigrateNote{Path: "a.rb", Action: "addField", Variant: "v", Why: "because"}
	if n.Path != "a.rb" || n.Action != "addField" || n.Variant != "v" || n.Why != "because" {
		t.Errorf("note = %+v", n)
	}
	if got := reflect.TypeOf(n).NumField(); got != 4 {
		t.Errorf("MigrateNote has %d fields, want 4", got)
	}
}

func TestMigrateReportHasTheFieldsTheRecordNames(t *testing.T) {
	note := MigrateNote{Path: "a.rb"}
	r := MigrateReport{
		Migrated: []MigrateNote{note}, Refused: []MigrateNote{note}, Left: []MigrateNote{note},
		Files: []string{"a.rb"}, Recordings: []string{"prov-2026-aaaaaaaa"}, Skipped: 2,
	}
	if len(r.Migrated) != 1 || len(r.Refused) != 1 || len(r.Left) != 1 || r.Files[0] != "a.rb" || r.Recordings[0] != "prov-2026-aaaaaaaa" {
		t.Errorf("report = %+v", r)
	}
	if got := reflect.TypeOf(r).NumField(); got != 6 || r.Skipped != 2 {
		t.Errorf("MigrateReport has %d fields, want 6", got)
	}
}

func TestMigratedRegionHasTheFieldsTheRecordNames(t *testing.T) {
	pkg := &genpkg.Package{Name: "rails"}
	m := migratedRegion{
		Path:    "a.rb",
		Package: pkg,
		Region:  inject.Region{Start: 3},
		Entry:   recording.Region{Path: "a.rb"},
	}
	if m.Path != "a.rb" || m.Package != pkg || m.Region.Start != 3 || m.Entry.Path != "a.rb" {
		t.Errorf("migratedRegion = %+v", m)
	}
	if got := reflect.TypeOf(m).NumField(); got != 4 {
		t.Errorf("migratedRegion has %d fields, want 4", got)
	}
}

// splitOpeningLine.

func TestSplitOpeningLine(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string
		start   int
		indent  string
		end     int
		ending  string
	}{
		{"no indentation, line feed", "abc\ndef\n", 0, "", 3, "\n"},
		{"spaces", "  abc\n", 0, "  ", 5, "\n"},
		{"tabs and spaces", " \t x\n", 0, " \t ", 4, "\n"},
		{"carriage return and line feed", "abc\r\ndef", 0, "", 3, "\r\n"},
		{"indented, carriage return and line feed", "\tabc\r\n", 0, "\t", 4, "\r\n"},
		{"content ends before a line feed", "  abc", 0, "  ", 5, ""},
		{"empty content", "", 0, "", 0, ""},
		{"only an indentation, no line feed", "   ", 0, "   ", 3, ""},
		{"only an indentation, line feed", "   \n", 0, "   ", 3, "\n"},
		{"a line in the middle", "first\n  second\r\nthird", 6, "  ", 14, "\r\n"},
		{"the last line, no ending", "first\n  last", 6, "  ", 12, ""},
		{"a carriage return not before the line feed is part of the text", "a\rb\n", 0, "", 3, "\n"},
		{"an empty line", "\n\n", 0, "", 0, "\n"},
		{"an empty line with carriage return", "\r\nnext", 0, "", 0, "\r\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indent, end, ending := splitOpeningLine(tc.content, tc.start)
			if indent != tc.indent || end != tc.end || ending != tc.ending {
				t.Errorf("splitOpeningLine(%q, %d) = (%q, %d, %q), want (%q, %d, %q)",
					tc.content, tc.start, indent, end, ending, tc.indent, tc.end, tc.ending)
			}
		})
	}
}

// migrateContent.

// migrationPackage is the loaded rails fixture package.
func migrationPackage(t *testing.T, dir string) *genpkg.Package {
	t.Helper()
	set, findings, err := genpkg.Load(dir, genpkg.Options{})
	if err != nil || findings.HasErrors() {
		t.Fatalf("genpkg.Load: %v %v", err, findings)
	}
	if len(set.Packages) != 1 {
		t.Fatalf("loaded %d packages, want 1", len(set.Packages))
	}
	return set.Packages[0]
}

func migOpen(indent, label, attrs string) string {
	return indent + "# sedum:" + label + " " + attrs
}

func migRegion(indent, label, attrs, body string) string {
	return migOpen(indent, label, attrs) + "\n" + body + indent + "# /sedum:" + label + "\n"
}

const (
	migHead = "class User\n  # sedum:anchor:class_body\n"
	migTail = "end\n"
)

func migFile(regions ...string) string { return migHead + strings.Join(regions, "") + migTail }

func migV1Attrs(record, gen, pkg, kwargs string) string {
	var parts []string
	parts = append(parts, `"tier":"owned"`)
	if record != "" {
		parts = append(parts, `"record":"`+record+`"`)
	}
	if gen != "" {
		parts = append(parts, `"gen":"`+gen+`"`)
	}
	if pkg != "" {
		parts = append(parts, `"pkg":"`+pkg+`"`)
	}
	parts = append(parts, `"kwargs":`+kwargs)
	return "{" + strings.Join(parts, ",") + "}"
}

func migV2Attrs(tier, record, gen, pkg, extraBeforeID, id, kwargsHash, body, extraAfter string) string {
	parts := []string{`"v":2`, `"tier":"` + tier + `"`}
	if record != "" {
		parts = append(parts, `"record":"`+record+`"`)
	}
	if gen != "" {
		parts = append(parts, `"gen":"`+gen+`"`)
	}
	if pkg != "" {
		parts = append(parts, `"pkg":"`+pkg+`"`)
	}
	if extraBeforeID != "" {
		parts = append(parts, extraBeforeID)
	}
	parts = append(parts, `"id":`+id, `"kwargs":"`+kwargsHash+`"`, `"body":"`+body+`"`)
	if extraAfter != "" {
		parts = append(parts, extraAfter)
	}
	return "{" + strings.Join(parts, ",") + "}"
}

const (
	migEmailKwargs = `{"field":"email","model":"user"}`
	migPhoneKwargs = `{"field":"phone","model":"user"}`
)

func migEmailV1() string {
	return migV1Attrs(v2RecordA, migGen, migPkg, migEmailKwargs)
}

func migEmailV2() string {
	return migV2Attrs("owned", v2RecordA, migGen, migPkg, "", migEmailKwargs, hashKwargsEmail, hashBodyEmail, "")
}

func migKnown(ids ...string) func(string) bool {
	return func(id string) bool {
		for _, k := range ids {
			if k == id {
				return true
			}
		}
		return false
	}
}

const migPath = "app/models/user.rb"

func TestMigrateContentRewritesAV1MarkerToV2AndNothingElse(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	content := migFile(migRegion("  ", "addField", migEmailV1(), "attribute :email\n"))
	want := migFile(migRegion("  ", "addField", migEmailV2(), "attribute :email\n"))

	got, migrated, refused, left, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatalf("migrateContent: %v", err)
	}
	if got != want {
		t.Errorf("content =\n%s\nwant\n%s", got, want)
	}
	if len(refused) != 0 || len(left) != 0 {
		t.Errorf("refused = %+v, left = %+v", refused, left)
	}
	if len(migrated) != 1 {
		t.Fatalf("migrated = %+v, want one", migrated)
	}

	m := migrated[0]
	if m.Path != migPath || m.Package != pkg {
		t.Errorf("Path = %q, Package = %p (want %p)", m.Path, m.Package, pkg)
	}
	// Region holds the region as it was read, with its v1 marker.
	if m.Region.Marker.Version != 0 || m.Region.Marker.Foreign || m.Region.Marker.Gen != migGen || m.Region.Marker.Pkg != migPkg {
		t.Errorf("Region.Marker = %+v, want the v1 marker as read", m.Region.Marker)
	}
	if jsonOf(t, m.Region.Marker.Kwargs) != migEmailKwargs || m.Region.BodyHash != hashBodyEmail {
		t.Errorf("Region = %+v", m.Region)
	}
	if m.Region.Start != strings.Index(content, "  # sedum:addField") {
		t.Errorf("Region.Start = %d", m.Region.Start)
	}

	e := m.Entry
	if e.Path != migPath || e.Action != "addField" || e.Variant != "" || e.Tier != "owned" || e.Pkg != migPkg {
		t.Errorf("Entry = %+v", e)
	}
	if jsonOf(t, e.ID) != migEmailKwargs || jsonOf(t, e.Kwargs) != migEmailKwargs {
		t.Errorf("Entry id = %v, kwargs = %v", e.ID, e.Kwargs)
	}
	if e.KwargsHash != hashKwargsEmail || e.Body != hashBodyEmail {
		t.Errorf("Entry hashes = %q, %q", e.KwargsHash, e.Body)
	}
}

func TestMigrateContentKeepsGenPkgTierWriterAndExtra(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	v1 := `{"tier":"seeded","record":"prov-2026-aaaaaaaa","gen":"sha256:g","pkg":"sha256:p","writer":"acme","kwargs":` + migEmailKwargs + `,"zzz":1,"aaa":{"x": true}}`
	// The new marker holds the nine keys in their order, then the keys Sedum
	// does not model, in sorted order and compacted.
	v2 := `{"v":2,"tier":"seeded","record":"prov-2026-aaaaaaaa","gen":"sha256:g","pkg":"sha256:p","writer":"acme",` +
		`"id":` + migEmailKwargs + `,"kwargs":"` + hashKwargsEmail + `","body":"` + hashBodyEmail + `","aaa":{"x":true},"zzz":1}`

	got, migrated, _, _, err := migrateContent(migPath, migFile(migRegion("  ", "addField", v1, "attribute :email\n")), pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatalf("migrateContent: %v", err)
	}
	if want := migFile(migRegion("  ", "addField", v2, "attribute :email\n")); got != want {
		t.Errorf("content =\n%s\nwant\n%s", got, want)
	}
	if len(migrated) != 1 || migrated[0].Entry.Tier != "seeded" || migrated[0].Entry.Pkg != "sha256:p" {
		t.Errorf("migrated = %+v", migrated)
	}
}

func TestMigrateContentAV1MarkerWithNoGenOrPkgStaysWithout(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	v1 := migV1Attrs(v2RecordA, "", "", migEmailKwargs)
	v2 := migV2Attrs("owned", v2RecordA, "", "", "", migEmailKwargs, hashKwargsEmail, hashBodyEmail, "")

	got, migrated, _, _, err := migrateContent(migPath, migFile(migRegion("  ", "addField", v1, "attribute :email\n")), pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatalf("migrateContent: %v", err)
	}
	if want := migFile(migRegion("  ", "addField", v2, "attribute :email\n")); got != want {
		t.Errorf("content =\n%s\nwant\n%s", got, want)
	}
	if len(migrated) != 1 || migrated[0].Entry.Pkg != "" {
		t.Errorf("migrated = %+v", migrated)
	}
}

func TestMigrateContentIdentityKwargsAreTheActionsSelectingOnes(t *testing.T) {
	pkg := migrationPackage(t, v2TypedGenerators(t))
	all := `{"field":"age","kind":"integer","model":"user"}`
	id := `{"field":"age","model":"user"}`
	v1 := migV1Attrs(v2RecordA, migGen, migPkg, all)
	v2 := migV2Attrs("owned", v2RecordA, migGen, migPkg, "", id, hashKwargsAge, hashBodyAge, "")

	got, migrated, _, _, err := migrateContent(migPath, migFile(migRegion("  ", "addTyped", v1, "attribute :age, :integer\n")), pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatalf("migrateContent: %v", err)
	}
	if want := migFile(migRegion("  ", "addTyped", v2, "attribute :age, :integer\n")); got != want {
		t.Errorf("content =\n%s\nwant\n%s", got, want)
	}
	e := migrated[0].Entry
	if jsonOf(t, e.ID) != id {
		t.Errorf("Entry.ID = %v, want the identity kwargs only", e.ID)
	}
	if jsonOf(t, e.Kwargs) != all {
		t.Errorf("Entry.Kwargs = %v, want every kwarg of the old marker", e.Kwargs)
	}
	if e.KwargsHash != hashKwargsAge || e.Body != hashBodyAge {
		t.Errorf("Entry hashes = %q, %q", e.KwargsHash, e.Body)
	}
}

func TestMigrateContentKeepsIndentationAndLineEndings(t *testing.T) {
	pkg := migrationPackage(t, generators(t))

	t.Run("tab indentation", func(t *testing.T) {
		content := migFile(migRegion("\t", "addField", migEmailV1(), "\tattribute :email\n"))
		got, migrated, _, _, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
		if err != nil {
			t.Fatal(err)
		}
		want := migFile(migRegion("\t", "addField", migV2Attrs("owned", v2RecordA, migGen, migPkg, "", migEmailKwargs, hashKwargsEmail, bodyHashOf(t, "\tattribute :email\n"), ""), "\tattribute :email\n"))
		if got != want || len(migrated) != 1 {
			t.Errorf("content =\n%q\nwant\n%q", got, want)
		}
	})

	t.Run("CRLF", func(t *testing.T) {
		crlf := func(s string) string { return strings.ReplaceAll(s, "\n", "\r\n") }
		content := crlf(migFile(migRegion("  ", "addField", migEmailV1(), "attribute :email\n")))
		got, migrated, _, _, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
		if err != nil {
			t.Fatal(err)
		}
		want := crlf(migFile(migRegion("  ", "addField", migEmailV2(), "attribute :email\n")))
		if got != want {
			t.Errorf("content =\n%q\nwant\n%q", got, want)
		}
		// The body hash reads each CRLF as LF, so the same region hashes the
		// same in either ending.
		if len(migrated) != 1 || migrated[0].Entry.Body != hashBodyEmail {
			t.Errorf("migrated = %+v", migrated)
		}
	})
}

// bodyHashOf is inject.BodyHash, which a body this package did not write has
// no literal for.
func bodyHashOf(t *testing.T, s string) string {
	t.Helper()
	return inject.BodyHash(s)
}

func TestMigrateContentRewritesSeveralRegionsLastToFirst(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	phoneV1 := migV1Attrs(v2RecordA, migGen, migPkg, migPhoneKwargs)
	phoneV2 := migV2Attrs("owned", v2RecordA, migGen, migPkg, "", migPhoneKwargs, hashKwargsPhone, hashBodyPhone, "")
	content := migFile(
		migRegion("  ", "addField", migEmailV1(), "attribute :email\n"),
		migRegion("    ", "addField", phoneV1, "attribute :phone\n"),
	)

	got, migrated, refused, left, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatal(err)
	}
	want := migFile(
		migRegion("  ", "addField", migEmailV2(), "attribute :email\n"),
		migRegion("    ", "addField", phoneV2, "attribute :phone\n"),
	)
	if got != want {
		t.Errorf("content =\n%s\nwant\n%s", got, want)
	}
	if len(refused)+len(left) != 0 || len(migrated) != 2 {
		t.Fatalf("migrated = %+v, refused = %+v, left = %+v", migrated, refused, left)
	}
	// File order, each holding the Package argument.
	if migrated[0].Entry.KwargsHash != hashKwargsEmail || migrated[1].Entry.KwargsHash != hashKwargsPhone {
		t.Errorf("migrated is not in file order: %+v", migrated)
	}
	if migrated[0].Region.Start >= migrated[1].Region.Start {
		t.Errorf("starts = %d, %d; want file order", migrated[0].Region.Start, migrated[1].Region.Start)
	}
	for _, m := range migrated {
		if m.Package != pkg || m.Path != migPath {
			t.Errorf("migrated region = %+v", m)
		}
	}
}

func TestMigrateContentVariantIsKept(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	content := migFile(migRegion("  ", "addField:strict", migEmailV1(), "attribute :email\n"))

	got, migrated, _, _, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatal(err)
	}
	if want := migFile(migRegion("  ", "addField:strict", migEmailV2(), "attribute :email\n")); got != want {
		t.Errorf("content =\n%s\nwant\n%s", got, want)
	}
	if len(migrated) != 1 || migrated[0].Entry.Variant != "strict" || migrated[0].Entry.Action != "addField" {
		t.Errorf("migrated = %+v", migrated)
	}
}

func TestMigrateContentLeavesARegionWhoseMarkerNamesNoRecord(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	content := migFile(migRegion("  ", "addField", migV1Attrs("", migGen, migPkg, migEmailKwargs), "attribute :email\n"))

	got, migrated, refused, left, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatal(err)
	}
	if got != content || len(migrated) != 0 || len(refused) != 0 {
		t.Errorf("content changed or reported: migrated = %+v, refused = %+v", migrated, refused)
	}
	want := []MigrateNote{{Path: migPath, Action: "addField", Why: "marker names no record"}}
	if !reflect.DeepEqual(left, want) {
		t.Errorf("left = %+v, want %+v", left, want)
	}
}

func TestMigrateContentLeavesARegionWhoseRecordIsNotKnown(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	content := migFile(migRegion("  ", "addField:later", migV1Attrs("prov-2026-gone0000", migGen, migPkg, migEmailKwargs), "attribute :email\n"))

	var asked []string
	known := func(id string) bool {
		asked = append(asked, id)
		return false
	}
	got, migrated, refused, left, err := migrateContent(migPath, content, pkg, known)
	if err != nil {
		t.Fatal(err)
	}
	if got != content || len(migrated) != 0 || len(refused) != 0 {
		t.Errorf("content changed or reported: migrated = %+v, refused = %+v", migrated, refused)
	}
	want := []MigrateNote{{Path: migPath, Action: "addField", Variant: "later", Why: "record not in the records directory"}}
	if !reflect.DeepEqual(left, want) {
		t.Errorf("left = %+v, want %+v", left, want)
	}
	if !reflect.DeepEqual(asked, []string{"prov-2026-gone0000"}) {
		t.Errorf("known was asked about %v, want the marker's record id", asked)
	}
}

func TestMigrateContentRefusesARegionWhoseActionThePackageNoLongerDeclares(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	content := migFile(migRegion("  ", "removedAction", migEmailV1(), "attribute :email\n"))

	got, migrated, refused, left, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatal(err)
	}
	if got != content || len(migrated) != 0 || len(left) != 0 {
		t.Errorf("content changed or reported: migrated = %+v, left = %+v", migrated, left)
	}
	want := []MigrateNote{{Path: migPath, Action: "removedAction", Why: "action no longer declared by its package"}}
	if !reflect.DeepEqual(refused, want) {
		t.Errorf("refused = %+v, want %+v", refused, want)
	}
}

func TestMigrateContentRefusesARegionWhoseKwargsCannotBeHashed(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	kwargs := `{"field":"email","model":"user","n":1152921504606846976}`
	content := migFile(migRegion("  ", "addField", migV1Attrs(v2RecordA, migGen, migPkg, kwargs), "attribute :email\n"))

	got, migrated, refused, left, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatal(err)
	}
	if got != content || len(migrated) != 0 || len(left) != 0 {
		t.Errorf("content changed or reported: migrated = %+v, left = %+v", migrated, left)
	}
	if len(refused) != 1 || refused[0].Path != migPath || refused[0].Action != "addField" {
		t.Fatalf("refused = %+v", refused)
	}
	if !strings.HasPrefix(refused[0].Why, "kwargs cannot be hashed:") || !strings.Contains(refused[0].Why, "integer beyond 2^53") {
		t.Errorf("Why = %q, want the text of inject.KwargsHash's error", refused[0].Why)
	}
}

func TestMigrateContentTouchesAndReportsNoV2OrForeignRegion(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	foreign := `{"v":3,"tier":"owned","record":"prov-2026-aaaaaaaa","kwargs":{"x":1}}`
	content := migFile(
		migRegion("  ", "addField", migEmailV2(), "attribute :email\n"),
		migRegion("  ", "addField:foreign", foreign, "attribute :foreign\n"),
		migRegion("  ", "addField:one", `{"v":1,"tier":"owned"}`, "attribute :one\n"),
	)

	got, migrated, refused, left, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatal(err)
	}
	if got != content {
		t.Errorf("content changed:\n%s", got)
	}
	if len(migrated)+len(refused)+len(left) != 0 {
		t.Errorf("reported: migrated = %+v, refused = %+v, left = %+v", migrated, refused, left)
	}
}

func TestMigrateContentWithNoV1RegionIsReturnedAsItIs(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	for _, content := range []string{"", "class User\nend\n", migFile(), "no newline"} {
		got, migrated, refused, left, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
		if err != nil {
			t.Fatalf("migrateContent(%q): %v", content, err)
		}
		if got != content {
			t.Errorf("content %q became %q", content, got)
		}
		if len(migrated) != 0 || len(refused) != 0 || len(left) != 0 {
			t.Errorf("content %q: migrated = %+v, refused = %+v, left = %+v", content, migrated, refused, left)
		}
	}
}

func TestMigrateContentMigratesWhatItCanBesideWhatItCannot(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	removed := migRegion("  ", "removedAction", migEmailV1(), "attribute :x\n")
	orphan := migRegion("  ", "addField:later", migV1Attrs("", migGen, migPkg, migPhoneKwargs), "attribute :phone\n")
	content := migFile(
		migRegion("  ", "addField", migEmailV1(), "attribute :email\n"),
		removed,
		orphan,
	)

	got, migrated, refused, left, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
	if err != nil {
		t.Fatal(err)
	}
	want := migFile(migRegion("  ", "addField", migEmailV2(), "attribute :email\n"), removed, orphan)
	if got != want {
		t.Errorf("content =\n%s\nwant\n%s", got, want)
	}
	if len(migrated) != 1 || len(refused) != 1 || len(left) != 1 {
		t.Errorf("migrated = %+v, refused = %+v, left = %+v", migrated, refused, left)
	}
}

func TestMigrateContentFailureToFindRegionsNamesTheFile(t *testing.T) {
	pkg := migrationPackage(t, generators(t))
	content := migHead + migOpen("  ", "addField", migEmailV1()) + "\nattribute :email\n" + migTail

	_, _, _, _, err := migrateContent(migPath, content, pkg, migKnown(v2RecordA))
	if err == nil {
		t.Fatal("migrateContent accepted a region that is never closed")
	}
	_, findErr := inject.FindRegions(pkg.Comment(), content)
	if findErr == nil {
		t.Fatal("test setup: FindRegions accepts the content")
	}
	if want := "file " + migPath + ": " + findErr.Error(); err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

// buildMigrationRecording.

type buildFixture struct {
	pkg     *genpkg.Package
	records *record.Set
	files   []resolve.File
	owners  map[string]inject.Owner
}

func newBuildFixture(t *testing.T) buildFixture {
	t.Helper()
	gens := generators(t)
	pkg := migrationPackage(t, gens)
	recs := writeTree(t, map[string]string{
		v2RecordA + ".yml":       recordFile(v2RecordA, userScope, "created_at: \"2026-03-04\"\n"),
		"prov-2026-bbbbbbbb.yml": recordFile("prov-2026-bbbbbbbb", "  - app/models/other.rb\n", "created_at: \"2026-04-05\"\n"),
	})
	set, _, err := record.Load(recs, record.Options{})
	if err != nil {
		t.Fatal(err)
	}

	file := func(rec, path string, p *genpkg.Package, unmanaged bool) resolve.File {
		return resolve.File{Resolution: resolve.Resolution{
			RecordID: rec, Path: path, Package: p, Template: "app/models/{name}.rb",
			Captures: map[string]string{"name": strings.TrimSuffix(filepath.Base(path), ".rb")}, Unmanaged: unmanaged,
		}}
	}
	return buildFixture{
		pkg:     pkg,
		records: set,
		files: []resolve.File{
			file(v2RecordA, "app/models/user.rb", pkg, false),
			file(v2RecordA, "app/models/order.rb", pkg, false),
			file(v2RecordA, "Gemfile", nil, true),
			file(v2RecordA, "notes.txt", nil, false),
			file("prov-2026-bbbbbbbb", "app/models/other.rb", pkg, false),
		},
		owners: ownersOf(set),
	}
}

// migrated builds a region a migration rewrote.
func (f buildFixture) migrated(path string, start int, field, gen string) migratedRegion {
	kwargs := map[string]any{"model": "user", "field": field}
	return migratedRegion{
		Path:    path,
		Package: f.pkg,
		Region: inject.Region{
			Start: start,
			Marker: inject.Marker{
				Action: "addField", Tier: inject.TierOwned, Record: v2RecordA, Gen: gen, Pkg: "sha256:p", Kwargs: kwargs,
			},
		},
		Entry: recording.Region{
			Path: path, Action: "addField", Tier: "owned", Pkg: "sha256:p",
			ID: kwargs, Kwargs: kwargs, KwargsHash: "sha256:k-" + field, Body: "sha256:b-" + field,
		},
	}
}

func TestBuildMigrationRecordingWithNoExistingRecording(t *testing.T) {
	f := newBuildFixture(t)
	later := f.migrated("app/models/user.rb", 40, "second", "sha256:gen-second")
	first := f.migrated("app/models/user.rb", 10, "first", "sha256:gen-first")
	other := f.migrated("app/models/order.rb", 5, "third", "sha256:gen-third")
	vars := map[string]string{"service": "orders"}

	// Given out of order: the regions are taken in path then position order.
	got := buildMigrationRecording(v2RecordA, []migratedRegion{later, other, first}, f.records, f.files, f.owners, vars, nil)

	if got.SedumVersion != recording.Version {
		t.Errorf("SedumVersion = %q, want %q", got.SedumVersion, recording.Version)
	}
	if !reflect.DeepEqual(got.Variables, vars) {
		t.Errorf("Variables = %v", got.Variables)
	}
	if want := map[string]recording.Package{"rails": {Extensions: []string{".rb"}}}; !reflect.DeepEqual(got.Packages, want) {
		t.Errorf("Packages = %v, want only the packages of the record's managed files: %v", got.Packages, want)
	}
	if len(got.Records) != 1 {
		t.Fatalf("Records = %+v, want one", got.Records)
	}

	rec := got.Records[0]
	if rec.RecordID != v2RecordA {
		t.Errorf("RecordID = %q", rec.RecordID)
	}
	if rec.CreatedAt != "2026-03-04" {
		t.Errorf("CreatedAt = %q, want the record's", rec.CreatedAt)
	}
	// The gen is the gen found on the markers, taken from the first region in
	// path then position order, not the record's current gen.
	current, _ := f.records.Lookup(v2RecordA)
	if rec.Gen != "sha256:gen-third" || rec.Gen == current.Gen {
		t.Errorf("Gen = %q, want the first migrated region's marker gen", rec.Gen)
	}

	wantFiles := []recording.File{
		{Path: "app/models/user.rb", Package: "rails", Template: "app/models/{name}.rb", Captures: map[string]string{"name": "user"}},
		{Path: "app/models/order.rb", Package: "rails", Template: "app/models/{name}.rb", Captures: map[string]string{"name": "order"}},
	}
	if !reflect.DeepEqual(rec.Files, wantFiles) {
		t.Errorf("Files = %+v, want %+v (unmanaged and package-less files left out)", rec.Files, wantFiles)
	}

	if len(rec.Phases) != 1 || rec.Phases[0].Name != recording.DefaultPhase {
		t.Fatalf("Phases = %+v, want the one default phase", rec.Phases)
	}
	ordered := []migratedRegion{other, first, later}
	var owned []inject.OwnedRegion
	for _, m := range ordered {
		owned = append(owned, inject.OwnedRegion{Path: m.Path, Package: m.Package, Region: m.Region})
	}
	wantInvocations := priorFromMarkers(v2RecordA, owned, f.owners)
	if len(wantInvocations) != 3 {
		t.Fatalf("test setup: priorFromMarkers = %+v", wantInvocations)
	}
	if !reflect.DeepEqual(rec.Phases[0].Invocations, wantInvocations) {
		t.Errorf("invocations = %+v, want what priorFromMarkers returns: %+v", rec.Phases[0].Invocations, wantInvocations)
	}

	wantRegions := []recording.Region{other.Entry, first.Entry, later.Entry}
	if !reflect.DeepEqual(rec.Regions, wantRegions) {
		t.Errorf("Regions = %+v, want the migrated entries in path then position order", rec.Regions)
	}
}

func TestBuildMigrationRecordingForARecordNotInTheSetHasNoCreatedAt(t *testing.T) {
	f := newBuildFixture(t)
	m := f.migrated("app/models/user.rb", 1, "first", "sha256:g")
	got := buildMigrationRecording("prov-2026-unknown0", []migratedRegion{m}, f.records, f.files, f.owners, nil, nil)
	if len(got.Records) != 1 || got.Records[0].RecordID != "prov-2026-unknown0" || got.Records[0].CreatedAt != "" {
		t.Errorf("Records = %+v", got.Records)
	}
}

func TestBuildMigrationRecordingKeepsWhatAnExistingRecordingHas(t *testing.T) {
	f := newBuildFixture(t)
	m1 := f.migrated("app/models/user.rb", 10, "email", "sha256:gen-marker")
	m2 := f.migrated("app/models/user.rb", 30, "phone", "sha256:gen-marker")

	kept := recording.Region{
		Path: "app/models/user.rb", Action: "addField", Tier: "owned", Pkg: "sha256:old",
		ID: map[string]any{"model": "user", "field": "kept"}, Kwargs: map[string]any{"model": "user", "field": "kept"},
		KwargsHash: "sha256:kept", Body: "sha256:kept",
	}
	stale := m1.Entry
	stale.KwargsHash, stale.Body, stale.Pkg = "sha256:stale", "sha256:stale", "sha256:stale"

	phases := []recording.Phase{{Name: recording.DefaultPhase, Invocations: []recording.Invocation{
		{Action: "addField", Kwargs: map[string]any{"model": "user", "field": "kept-phase"}},
	}}}
	files := []recording.File{{Path: "kept.rb", Package: "rails", Template: "t"}}
	existing := &recording.Recording{
		SedumVersion: "0.8.0",
		Packages:     map[string]recording.Package{"rails": {Extensions: []string{".rb"}}},
		Records: []recording.Record{{
			RecordID: v2RecordA, Gen: "sha256:existing-gen", CreatedAt: "2026-01-01",
			Files: files, Phases: phases, Regions: []recording.Region{kept, stale},
		}},
	}

	got := buildMigrationRecording(v2RecordA, []migratedRegion{m1, m2}, f.records, f.files, f.owners, nil, existing)
	if len(got.Records) != 1 {
		t.Fatalf("Records = %+v", got.Records)
	}
	rec := got.Records[0]
	if rec.RecordID != v2RecordA || rec.Gen != "sha256:existing-gen" || rec.CreatedAt != "2026-01-01" {
		t.Errorf("record = %+v, want the existing recording's id, gen and created_at", rec)
	}
	if !reflect.DeepEqual(rec.Files, files) || !reflect.DeepEqual(rec.Phases, phases) {
		t.Errorf("files = %+v, phases = %+v; want the existing recording's", rec.Files, rec.Phases)
	}

	keys := map[string]recording.Region{}
	for _, r := range rec.Regions {
		keys[recording.RegionKey(r.Path, r.Action, r.Variant, r.ID)] = r
	}
	if len(rec.Regions) != 3 || len(keys) != 3 {
		t.Fatalf("Regions = %+v, want the kept entry and the two migrated ones, each once", rec.Regions)
	}
	if got := keys[recording.RegionKey(kept.Path, kept.Action, kept.Variant, kept.ID)]; !reflect.DeepEqual(got, kept) {
		t.Errorf("the entry of a region not migrated changed: %+v", got)
	}
	if got := keys[recording.RegionKey(m1.Entry.Path, m1.Entry.Action, m1.Entry.Variant, m1.Entry.ID)]; !reflect.DeepEqual(got, m1.Entry) {
		t.Errorf("the entry of a migrated region was not replaced: %+v", got)
	}
	if got := keys[recording.RegionKey(m2.Entry.Path, m2.Entry.Action, m2.Entry.Variant, m2.Entry.ID)]; !reflect.DeepEqual(got, m2.Entry) {
		t.Errorf("the entry of a migrated region was not added: %+v", got)
	}
}

// Migrate.

// migFixture is a rails output tree of v1 markers, its records and generators.
type migFixture struct {
	cfg Config
	mc  MigrateConfig
}

func newMigFixture(t *testing.T) migFixture {
	t.Helper()
	cfg := config(t)
	return migFixture{cfg: cfg, mc: MigrateConfig{Generators: cfg.Generators, Records: cfg.Records, Output: cfg.Output}}
}

// current is the gen of the record and the pkg of addField as they are now.
func (f migFixture) current(t *testing.T) (gen, pkg string) {
	t.Helper()
	records, _, err := record.Load(f.cfg.Records, record.Options{})
	if err != nil {
		t.Fatal(err)
	}
	return records.Records[0].Gen, migrationPackage(t, f.cfg.Generators).Actions["addField"].PkgHash
}

func (f migFixture) writeUser(t *testing.T, regions ...string) {
	t.Helper()
	v2WriteFile(t, f.cfg, migPath, migFile(regions...))
}

func (f migFixture) email(gen, pkg string) string {
	return migRegion("  ", "addField", migV1Attrs(v2RecordA, gen, pkg, migEmailKwargs), "attribute :email\n")
}

func (f migFixture) emailV2(gen, pkg string) string {
	return migRegion("  ", "addField", migV2Attrs("owned", v2RecordA, gen, pkg, "", migEmailKwargs, hashKwargsEmail, hashBodyEmail, ""), "attribute :email\n")
}

func (f migFixture) migrate(t *testing.T) *MigrateReport {
	t.Helper()
	report, err := Migrate(f.mc)
	if err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	if report == nil {
		t.Fatal("Migrate returned no report")
	}
	return report
}

func TestMigrateRewritesTheMarkersAndWritesTheRecordings(t *testing.T) {
	f := newMigFixture(t)
	gen, pkg := f.current(t)
	f.writeUser(t, f.email(gen, pkg))
	if err := os.Chmod(filepath.Join(f.cfg.Output, migPath), 0o640); err != nil {
		t.Fatal(err)
	}

	report := f.migrate(t)

	if want := []MigrateNote{{Path: migPath, Action: "addField"}}; !reflect.DeepEqual(report.Migrated, want) {
		t.Errorf("Migrated = %+v, want %+v", report.Migrated, want)
	}
	if len(report.Refused) != 0 || len(report.Left) != 0 {
		t.Errorf("Refused = %+v, Left = %+v", report.Refused, report.Left)
	}
	if !reflect.DeepEqual(report.Files, []string{migPath}) || !reflect.DeepEqual(report.Recordings, []string{v2RecordA}) {
		t.Errorf("Files = %v, Recordings = %v", report.Files, report.Recordings)
	}

	// The file: the opening marker line is the only thing that changed, and the
	// file keeps its permissions.
	if got, want := userModel(t, f.cfg), migFile(f.emailV2(gen, pkg)); got != want {
		t.Errorf("user.rb =\n%s\nwant\n%s", got, want)
	}
	info, err := os.Stat(filepath.Join(f.cfg.Output, migPath))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Errorf("mode = %v, want the original 0640", info.Mode().Perm())
	}
	region := v2Regions(t, f.cfg, migPath)[0]
	if region.Marker.Version != 2 || region.Marker.Body != region.BodyHash {
		t.Errorf("the new marker's body hash does not match the text: %+v", region)
	}

	// The recording.
	dir := filepath.Join(f.cfg.Output, ".sedum", "recordings")
	file, err := recording.ReadRecord(dir, v2RecordA)
	if err != nil {
		t.Fatalf("ReadRecord: %v", err)
	}
	if file.SedumVersion != recording.Version {
		t.Errorf("SedumVersion = %q", file.SedumVersion)
	}
	rec := file.Records[0]
	if rec.Gen != gen {
		t.Errorf("Gen = %q, want %q", rec.Gen, gen)
	}
	if len(rec.Files) != 1 || rec.Files[0].Path != migPath || rec.Files[0].Package != "rails" ||
		rec.Files[0].Template != "app/models/{name}.rb" || !reflect.DeepEqual(rec.Files[0].Captures, map[string]string{"name": "user"}) {
		t.Errorf("Files = %+v", rec.Files)
	}
	if len(rec.Phases) != 1 || rec.Phases[0].Name != recording.DefaultPhase || len(rec.Phases[0].Invocations) != 1 ||
		rec.Phases[0].Invocations[0].Action != "addField" || jsonOf(t, rec.Phases[0].Invocations[0].Kwargs) != migEmailKwargs {
		t.Errorf("Phases = %+v, want the selection rebuilt from the marker", rec.Phases)
	}
	if len(rec.Regions) != 1 {
		t.Fatalf("Regions = %+v", rec.Regions)
	}
	e := rec.Regions[0]
	if e.Path != migPath || e.Action != "addField" || e.Tier != "owned" || e.Pkg != pkg ||
		jsonOf(t, e.ID) != migEmailKwargs || jsonOf(t, e.Kwargs) != migEmailKwargs ||
		e.KwargsHash != hashKwargsEmail || e.Body != hashBodyEmail {
		t.Errorf("entry = %+v", e)
	}
}

// The recording of a migrated record carries the gen on its markers, not the
// record's current gen, and the marker keeps its gen and pkg.
func TestMigrateKeepsTheGenAndPkgOnTheMarkerAndRecordsTheMarkersGen(t *testing.T) {
	f := newMigFixture(t)
	f.writeUser(t, f.email("sha256:old-gen", "sha256:old-pkg"))

	f.migrate(t)

	if got, want := userModel(t, f.cfg), migFile(f.emailV2("sha256:old-gen", "sha256:old-pkg")); got != want {
		t.Errorf("user.rb =\n%s\nwant\n%s", got, want)
	}
	rec := v2ReadRecord(t, filepath.Join(f.cfg.Output, ".sedum", "recordings"), v2RecordA)
	if rec.Gen != "sha256:old-gen" {
		t.Errorf("Gen = %q, want the marker's, not the record's current gen", rec.Gen)
	}
}

// grow --check reports the same classes after migration as before it.
func TestMigrateKeepsTheClassesGrowCheckReports(t *testing.T) {
	for _, tc := range []struct {
		name     string
		gen, pkg func(current string) string
		want     Mode
	}{
		{"current", func(c string) string { return c }, func(c string) string { return c }, ModeSkip},
		{"stale gen", func(string) string { return "sha256:other-gen" }, func(c string) string { return c }, ModeReselect},
		{"unknown gen", func(string) string { return "" }, func(c string) string { return c }, ModeReselect},
		{"stale pkg", func(c string) string { return c }, func(string) string { return "sha256:other-pkg" }, ModeRerender},
		{"unknown pkg", func(c string) string { return c }, func(string) string { return "" }, ModeRerender},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newMigFixture(t)
			curGen, curPkg := f.current(t)
			f.writeUser(t, f.email(tc.gen(curGen), tc.pkg(curPkg)))

			before := v2Class(t, v2Check(t, f.cfg), v2RecordA)
			f.migrate(t)
			after := v2Class(t, v2Check(t, f.cfg), v2RecordA)

			if before.Mode != tc.want {
				t.Fatalf("test setup: before migration mode = %v, want %v", before.Mode, tc.want)
			}
			if after.Mode != before.Mode || len(after.Stale) != len(before.Stale) || len(after.Rerender) != len(before.Rerender) {
				t.Errorf("after migration: %+v, before: %+v", after, before)
			}
		})
	}
}

func TestMigrateReRenderingAMigratedRegionFromItsRecordingReproducesItsBody(t *testing.T) {
	f := newMigFixture(t)
	gen, _ := f.current(t)
	// Current gen, stale pkg: the migrated region is re-rendered, with no
	// model, from its recording's entry.
	f.writeUser(t, f.email(gen, "sha256:stale-pkg"))
	before := v2Regions(t, f.cfg, migPath)[0]

	f.migrate(t)
	migrated := v2Regions(t, f.cfg, migPath)[0]
	if migrated.Marker.Body != before.BodyHash {
		t.Fatalf("the migrated body hash %q is not the text's %q", migrated.Marker.Body, before.BodyHash)
	}

	result, client := v2Run(t, f.cfg, fields("email"))
	if client.calls != 0 {
		t.Errorf("a re-render made %d model call(s)", client.calls)
	}
	if v2Class(t, result, v2RecordA).Mode != ModeRerender {
		t.Fatalf("classifications = %+v, want the migrated region re-rendered", result.Classifications)
	}
	after := v2Regions(t, f.cfg, migPath)[0]
	if after.BodyHash != before.BodyHash || after.Marker.Body != before.BodyHash {
		t.Errorf("re-rendering from the recording changed the body: %q -> %q (marker %q)", before.BodyHash, after.BodyHash, after.Marker.Body)
	}
	if !strings.Contains(userModel(t, f.cfg), "\nattribute :email\n  # /sedum:addField") {
		t.Errorf("the body text changed:\n%s", userModel(t, f.cfg))
	}
	if after.Marker.Pkg == "sha256:stale-pkg" || after.Marker.Gen != gen {
		t.Errorf("marker = %+v, want the current pkg and the gen kept", after.Marker)
	}
}

func TestMigratingASecondTimeChangesNothing(t *testing.T) {
	f := newMigFixture(t)
	gen, pkg := f.current(t)
	f.writeUser(t, f.email(gen, pkg))
	f.migrate(t)
	ageTree(t, f.cfg.Output)
	before := snap(t, f.cfg.Output)

	report := f.migrate(t)

	sameSnapshot(t, before, snap(t, f.cfg.Output))
	if n := len(report.Migrated) + len(report.Refused) + len(report.Left) + len(report.Files) + len(report.Recordings); n != 0 {
		t.Errorf("a second migration reported %+v", report)
	}
}

func TestMigrateDryRunWritesNothingAndReportsWhatItWouldDo(t *testing.T) {
	f := newMigFixture(t)
	gen, pkg := f.current(t)
	f.writeUser(t, f.email(gen, pkg))
	ageTree(t, f.cfg.Output)
	before := snap(t, f.cfg.Output)

	f.mc.DryRun = true
	report := f.migrate(t)

	sameSnapshot(t, before, snap(t, f.cfg.Output))
	if got := tree(t, f.cfg.Output); !reflect.DeepEqual(got, []string{migPath}) {
		t.Errorf("a dry run left %v", got)
	}
	if len(report.Migrated) != 1 || !reflect.DeepEqual(report.Files, []string{migPath}) || !reflect.DeepEqual(report.Recordings, []string{v2RecordA}) {
		t.Errorf("report = %+v, want the same lists a real run reports", report)
	}
}

func TestMigrateRefusedAndLeftRegionsKeepTheirV1MarkersAndHaveNoEntry(t *testing.T) {
	f := newMigFixture(t)
	f.cfg.Records = recordsDir(t, map[string]string{
		v2RecordA:            "  - app/models/user.rb\n",
		"prov-2026-bbbbbbbb": "  - app/models/order.rb\n",
	})
	f.mc.Records = f.cfg.Records
	gen, pkg := f.current(t)

	removed := migRegion("  ", "removedAction", migV1Attrs(v2RecordA, gen, pkg, migEmailKwargs), "attribute :x\n")
	noRecord := migRegion("  ", "addField", migV1Attrs("", gen, pkg, migPhoneKwargs), "attribute :phone\n")
	gone := migRegion("  ", "addField:later", migV1Attrs("prov-2026-gone0000", gen, pkg, migPhoneKwargs), "attribute :later\n")
	f.writeUser(t, f.email(gen, pkg), removed, noRecord, gone)

	order := migFile(migRegion("  ", "removedAction", migV1Attrs("prov-2026-bbbbbbbb", gen, pkg, migEmailKwargs), "attribute :x\n"))
	if err := os.MkdirAll(filepath.Join(f.cfg.Output, "app", "models"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.cfg.Output, "app", "models", "order.rb"), []byte(order), 0o644); err != nil {
		t.Fatal(err)
	}

	report := f.migrate(t)

	if want := []MigrateNote{{Path: migPath, Action: "addField"}}; !reflect.DeepEqual(report.Migrated, want) {
		t.Errorf("Migrated = %+v", report.Migrated)
	}
	wantRefused := []MigrateNote{
		{Path: "app/models/order.rb", Action: "removedAction", Why: "action no longer declared by its package"},
		{Path: migPath, Action: "removedAction", Why: "action no longer declared by its package"},
	}
	if !reflect.DeepEqual(report.Refused, wantRefused) {
		t.Errorf("Refused = %+v, want %+v", report.Refused, wantRefused)
	}
	wantLeft := []MigrateNote{
		{Path: migPath, Action: "addField", Why: "marker names no record"},
		{Path: migPath, Action: "addField", Variant: "later", Why: "record not in the records directory"},
	}
	if !reflect.DeepEqual(report.Left, wantLeft) {
		t.Errorf("Left = %+v, want %+v", report.Left, wantLeft)
	}
	// Only the file with a migrated region is rewritten, and only the record
	// with one has a recording.
	if !reflect.DeepEqual(report.Files, []string{migPath}) || !reflect.DeepEqual(report.Recordings, []string{v2RecordA}) {
		t.Errorf("Files = %v, Recordings = %v", report.Files, report.Recordings)
	}

	if got, want := userModel(t, f.cfg), migFile(f.emailV2(gen, pkg), removed, noRecord, gone); got != want {
		t.Errorf("user.rb =\n%s\nwant\n%s", got, want)
	}
	if got := readFile(t, filepath.Join(f.cfg.Output, "app", "models", "order.rb")); got != order {
		t.Errorf("a file with only a refused region was rewritten:\n%s", got)
	}

	dir := filepath.Join(f.cfg.Output, ".sedum", "recordings")
	if got := tree(t, dir); !reflect.DeepEqual(got, []string{v2RecordA + ".json"}) {
		t.Errorf("recordings = %v, want only the record with a migrated region", got)
	}
	if rec := v2ReadRecord(t, dir, v2RecordA); len(rec.Regions) != 1 || rec.Regions[0].KwargsHash != hashKwargsEmail {
		t.Errorf("regions = %+v, want only the migrated region's entry", rec.Regions)
	}
}

func TestMigrateAddsToAnExistingRecordingAndKeepsWhatItHas(t *testing.T) {
	cfg, dir := v2Grow(t)
	v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) {
		r.Gen = "sha256:sentinel-gen"
		r.CreatedAt = "2026-05-05"
	})
	before := v2ReadRecord(t, dir, v2RecordA)

	f := migFixture{cfg: cfg, mc: MigrateConfig{Generators: cfg.Generators, Records: cfg.Records, Output: cfg.Output}}
	gen, pkg := f.current(t)
	phone := migRegion("  ", "addField", migV1Attrs(v2RecordA, gen, pkg, migPhoneKwargs), "attribute :phone\n")
	content := userModel(t, cfg)
	v2WriteFile(t, cfg, migPath, strings.TrimSuffix(content, migTail)+phone+migTail)

	report := f.migrate(t)
	if len(report.Migrated) != 1 || report.Migrated[0].Action != "addField" {
		t.Fatalf("Migrated = %+v, want only the v1 region", report.Migrated)
	}

	after := v2ReadRecord(t, dir, v2RecordA)
	if after.Gen != "sha256:sentinel-gen" || after.CreatedAt != "2026-05-05" {
		t.Errorf("gen = %q, created_at = %q; the existing recording's must be kept", after.Gen, after.CreatedAt)
	}
	if jsonOf(t, after.Phases) != jsonOf(t, before.Phases) || jsonOf(t, after.Files) != jsonOf(t, before.Files) {
		t.Errorf("phases or files changed:\n%+v\n%+v", after, before)
	}
	if len(after.Regions) != 2 {
		t.Fatalf("regions = %+v, want the existing entry and the migrated one", after.Regions)
	}
	byHash := map[string]recording.Region{}
	for _, r := range after.Regions {
		byHash[r.KwargsHash] = r
	}
	if byHash[hashKwargsEmail].Body != hashBodyEmail || byHash[hashKwargsPhone].Body != hashBodyPhone {
		t.Errorf("regions = %+v", after.Regions)
	}
}

func TestMigrateReadsNothingInTheGeneratorsOrRecordsDirectories(t *testing.T) {
	f := newMigFixture(t)
	gens := filepath.Join(f.cfg.Output, "generators")
	if err := os.CopyFS(gens, os.DirFS(f.cfg.Generators)); err != nil {
		t.Fatal(err)
	}
	recs := filepath.Join(f.cfg.Output, "records")
	if err := os.CopyFS(recs, os.DirFS(f.cfg.Records)); err != nil {
		t.Fatal(err)
	}
	f.cfg.Generators, f.cfg.Records = gens, recs
	f.mc.Generators, f.mc.Records = gens, recs

	gen, pkg := f.current(t)
	golden := migFile(f.email(gen, pkg))
	fixtures := []string{
		filepath.Join(gens, "rails", "testdata", "golden", "user.rb"),
		filepath.Join(recs, "notes.rb"),
	}
	for _, p := range fixtures {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(golden), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	f.writeUser(t, f.email(gen, pkg))

	report := f.migrate(t)

	if !reflect.DeepEqual(report.Files, []string{migPath}) || len(report.Migrated) != 1 {
		t.Errorf("report = %+v, want only user.rb", report)
	}
	for _, p := range fixtures {
		if got := readFile(t, p); got != golden {
			t.Errorf("%s was migrated:\n%s", p, got)
		}
	}
}

func TestMigrateGroupsRegionsUnderTheSurvivingIdOfAnAlias(t *testing.T) {
	f := newMigFixture(t)
	f.cfg.Records = writeTree(t, map[string]string{
		v2RecordA + ".yml": recordFile(v2RecordA, userScope, "aliases: [prov-2026-01d01d01]\n"),
	})
	f.mc.Records = f.cfg.Records
	gen, pkg := f.current(t)
	f.writeUser(t, migRegion("  ", "addField", migV1Attrs("prov-2026-01d01d01", gen, pkg, migEmailKwargs), "attribute :email\n"))

	report := f.migrate(t)

	if len(report.Migrated) != 1 || len(report.Left) != 0 {
		t.Errorf("Migrated = %+v, Left = %+v; a marker naming an alias is a known record", report.Migrated, report.Left)
	}
	if !reflect.DeepEqual(report.Recordings, []string{v2RecordA}) {
		t.Errorf("Recordings = %v, want the surviving id", report.Recordings)
	}
	// The marker keeps the record it named.
	if got := userModel(t, f.cfg); !strings.Contains(got, `"record":"prov-2026-01d01d01"`) {
		t.Errorf("the alias on the marker was changed:\n%s", got)
	}
	if rec := v2ReadRecord(t, filepath.Join(f.cfg.Output, ".sedum", "recordings"), v2RecordA); len(rec.Regions) != 1 {
		t.Errorf("regions = %+v", rec.Regions)
	}
}

func TestMigrateSortsItsListsByPathActionAndVariantAndByRecordId(t *testing.T) {
	f := newMigFixture(t)
	f.cfg.Records = recordsDir(t, map[string]string{
		v2RecordA:            "  - app/models/user.rb\n",
		"prov-2026-bbbbbbbb": "  - app/models/order.rb\n",
	})
	f.mc.Records = f.cfg.Records
	gen, pkg := f.current(t)

	f.writeUser(t, f.email(gen, pkg))
	orderEmail := migRegion("  ", "addField", migV1Attrs("prov-2026-bbbbbbbb", gen, pkg, migPhoneKwargs), "attribute :phone\n")
	if err := os.WriteFile(filepath.Join(f.cfg.Output, "app", "models", "order.rb"), []byte(migFile(orderEmail)), 0o644); err != nil {
		t.Fatal(err)
	}

	report := f.migrate(t)

	if got, want := report.Files, []string{"app/models/order.rb", migPath}; !reflect.DeepEqual(got, want) {
		t.Errorf("Files = %v, want %v", got, want)
	}
	if got, want := report.Recordings, []string{v2RecordA, "prov-2026-bbbbbbbb"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Recordings = %v, want %v", got, want)
	}
	if len(report.Migrated) != 2 || report.Migrated[0].Path != "app/models/order.rb" || report.Migrated[1].Path != migPath {
		t.Errorf("Migrated = %+v, want sorted by path", report.Migrated)
	}
}

func TestMigrateReportsGeneratorPackagesThatDidNotLoad(t *testing.T) {
	f := newMigFixture(t)
	f.mc.Generators = writeTree(t, map[string]string{
		"rails/sedum.yaml":           "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\nextnesions: [\".erb\"]\n",
		"rails/actions/actions.yaml": "actions: {}\n",
	})

	_, err := Migrate(f.mc)
	if err == nil {
		t.Fatal("Migrate loaded a broken generator package")
	}
	if !strings.HasPrefix(err.Error(), "generator packages did not load:\n  ") {
		t.Errorf("error = %q, want the prefix %q", err.Error(), "generator packages did not load:\n  ")
	}
	if !strings.Contains(err.Error(), "extnesions") {
		t.Errorf("error = %q, want every finding's text", err.Error())
	}
}

func TestMigrateRunsWithNoModel(t *testing.T) {
	// There is no client in MigrateConfig: the type says so. A migration that
	// asked for one would not compile, and one that reached for a global would
	// fail a run with no model configured.
	f := newMigFixture(t)
	gen, pkg := f.current(t)
	f.writeUser(t, f.email(gen, pkg))
	if _, err := Migrate(f.mc); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	if _, ok := reflect.TypeOf(f.mc).FieldByName("Client"); ok {
		t.Error("MigrateConfig has a Client")
	}
}
