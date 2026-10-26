package recording

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// WriteRecords and ReadRecord keep one recording file per record
// (prov-2026-2a2edf2b, prov-2026-296d6901).

func regionOf(path, action, variant string, id map[string]any) Region {
	return Region{
		Path: path, Action: action, Variant: variant, Tier: "owned",
		Pkg: "sha256:pkg", ID: id, Kwargs: id, KwargsHash: "sha256:kw", Body: "sha256:body",
	}
}

func recordOf(id string, pkg string, regions ...Region) Record {
	return Record{
		RecordID:  id,
		Gen:       "sha256:gen-" + id,
		CreatedAt: "2026-01-02",
		Files:     []File{{Path: id + ".rb", Package: pkg, Template: "{name}.rb"}},
		Phases: []Phase{{Name: DefaultPhase, Invocations: []Invocation{
			{Action: "addField", Kwargs: map[string]any{"field": id}},
		}}},
		Regions: regions,
	}
}

func fixture() Recording {
	return Recording{
		SedumVersion: "9.9.9",
		Packages: map[string]Package{
			"rails": {Extensions: []string{".rb"}},
			"go":    {Extensions: []string{".go"}},
			"other": {Extensions: []string{".x"}},
		},
		Records: []Record{
			recordOf("prov-2026-aaaaaaaa", "rails"),
			recordOf("prov-2026-bbbbbbbb", "go"),
		},
		Variables: map[string]string{"service": "orders"},
		Complete:  true,
		Prompt:    "carried text",
	}
}

func listDir(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir %s: %v", dir, err)
	}
	var out []string
	for _, e := range entries {
		out = append(out, e.Name())
	}
	sort.Strings(out)
	return out
}

func readBack(t *testing.T, dir, id string) Recording {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, id+".json"))
	if err != nil {
		t.Fatalf("read %s: %v", id, err)
	}
	rec, err := Read(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode %s: %v", id, err)
	}
	return rec
}

func TestWriteRecordsWritesOneFilePerRecord(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "recordings")
	in := fixture()

	if err := WriteRecords(dir, in); err != nil {
		t.Fatalf("WriteRecords: %v", err)
	}

	// Nothing but the two record files: the temporary file each is written
	// through is renamed away.
	if got, want := listDir(t, dir), []string{"prov-2026-aaaaaaaa.json", "prov-2026-bbbbbbbb.json"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("directory holds %v, want %v", got, want)
	}

	for _, rec := range in.Records {
		got := readBack(t, dir, rec.RecordID)
		if len(got.Records) != 1 || got.Records[0].RecordID != rec.RecordID {
			t.Errorf("%s holds %+v, want exactly the one record it is named for", rec.RecordID, got.Records)
		}
		if got.SedumVersion != "9.9.9" {
			t.Errorf("%s: SedumVersion = %q", rec.RecordID, got.SedumVersion)
		}
		if !reflect.DeepEqual(got.Variables, map[string]string{"service": "orders"}) {
			t.Errorf("%s: Variables = %v", rec.RecordID, got.Variables)
		}
		if got.Complete {
			t.Errorf("%s: Complete is true; a per-record file never says it holds a complete run", rec.RecordID)
		}
		if got.Prompt != "" {
			t.Errorf("%s: Prompt = %q, want empty", rec.RecordID, got.Prompt)
		}
	}
}

func TestWriteRecordsPackagesAreOnlyThoseTheRecordsFilesName(t *testing.T) {
	dir := t.TempDir()
	if err := WriteRecords(dir, fixture()); err != nil {
		t.Fatalf("WriteRecords: %v", err)
	}

	a := readBack(t, dir, "prov-2026-aaaaaaaa")
	if want := map[string]Package{"rails": {Extensions: []string{".rb"}}}; !reflect.DeepEqual(a.Packages, want) {
		t.Errorf("aaaaaaaa Packages = %v, want %v", a.Packages, want)
	}
	b := readBack(t, dir, "prov-2026-bbbbbbbb")
	if want := map[string]Package{"go": {Extensions: []string{".go"}}}; !reflect.DeepEqual(b.Packages, want) {
		t.Errorf("bbbbbbbb Packages = %v, want %v", b.Packages, want)
	}
}

func TestWriteRecordsWritesRegionsSortedByRegionKey(t *testing.T) {
	r1 := regionOf("b.rb", "addField", "", map[string]any{"f": "1"})
	r2 := regionOf("a.rb", "addSetting", "strict", map[string]any{"f": "1"})
	r3 := regionOf("a.rb", "addField", "", map[string]any{"f": "2"})
	r4 := regionOf("a.rb", "addField", "", map[string]any{"f": "1"})
	unsorted := []Region{r1, r2, r3, r4}

	rec := recordOf("prov-2026-aaaaaaaa", "rails", unsorted...)
	in := Recording{SedumVersion: Version, Packages: fixture().Packages, Records: []Record{rec}}

	dir := t.TempDir()
	if err := WriteRecords(dir, in); err != nil {
		t.Fatalf("WriteRecords: %v", err)
	}

	got := readBack(t, dir, "prov-2026-aaaaaaaa").Records[0].Regions
	want := []Region{r4, r3, r2, r1}
	if !reflect.DeepEqual(keysOf(got), keysOf(want)) {
		t.Errorf("regions written in order\n%v\nwant\n%v", keysOf(got), keysOf(want))
	}

	// The argument's own slice is left in the order it came in.
	if !reflect.DeepEqual(in.Records[0].Regions, unsorted) {
		t.Errorf("the argument's Regions were reordered: %v", keysOf(in.Records[0].Regions))
	}
	if !reflect.DeepEqual(rec.Regions, []Region{r1, r2, r3, r4}) {
		t.Errorf("the caller's slice was reordered: %v", keysOf(rec.Regions))
	}
}

func keysOf(regions []Region) []string {
	var out []string
	for _, r := range regions {
		out = append(out, RegionKey(r.Path, r.Action, r.Variant, r.ID))
	}
	return out
}

// The file's bytes are those of Write, so a reader that parses a per-record
// file and one that parses a whole recording read the same shape.
func TestWriteRecordsFileBytesAreThoseOfWrite(t *testing.T) {
	region := regionOf("a.rb", "addField", "", map[string]any{"f": "1"})
	rec := recordOf("prov-2026-aaaaaaaa", "rails", region)
	in := Recording{
		SedumVersion: "9.9.9",
		Packages:     fixture().Packages,
		Records:      []Record{rec},
		Variables:    map[string]string{"service": "orders"},
		Complete:     true,
		Prompt:       "dropped",
	}

	dir := t.TempDir()
	if err := WriteRecords(dir, in); err != nil {
		t.Fatalf("WriteRecords: %v", err)
	}

	var want bytes.Buffer
	if err := Write(&want, Recording{
		SedumVersion: "9.9.9",
		Packages:     map[string]Package{"rails": fixture().Packages["rails"]},
		Records:      []Record{rec},
		Variables:    map[string]string{"service": "orders"},
	}); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "prov-2026-aaaaaaaa.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want.String() {
		t.Errorf("file bytes differ from Write's:\n got %s\nwant %s", got, want.String())
	}
}

func TestWriteRecordsReplacesAnExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prov-2026-aaaaaaaa.json")
	if err := os.WriteFile(path, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	in := Recording{SedumVersion: Version, Packages: fixture().Packages, Records: []Record{recordOf("prov-2026-aaaaaaaa", "rails")}}
	if err := WriteRecords(dir, in); err != nil {
		t.Fatalf("WriteRecords: %v", err)
	}
	if got := readBack(t, dir, "prov-2026-aaaaaaaa"); len(got.Records) != 1 {
		t.Errorf("stale file was not replaced: %+v", got)
	}
	if got := listDir(t, dir); len(got) != 1 {
		t.Errorf("directory holds %v after replacing", got)
	}
}

func TestWriteRecordsOfNoRecordsCreatesTheDirectoryAndWritesNoFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b")
	if err := WriteRecords(dir, Recording{SedumVersion: Version}); err != nil {
		t.Fatalf("WriteRecords: %v", err)
	}
	if got := listDir(t, dir); len(got) != 0 {
		t.Errorf("directory holds %v, want nothing", got)
	}
}

func TestWriteRecordsReportsADirectoryItCannotCreate(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(blocker, "recordings")

	err := WriteRecords(dir, fixture())
	if err == nil {
		t.Fatal("WriteRecords under a regular file succeeded")
	}
	if prefix := "creating recordings directory " + dir + ": "; !strings.HasPrefix(err.Error(), prefix) {
		t.Errorf("error = %q, want the prefix %q", err.Error(), prefix)
	}
}

func TestWriteRecordsRefusesARecordIdThatCannotNameAFile(t *testing.T) {
	for _, tc := range []struct {
		name, id, wantText string
	}{
		{"empty", "", `record id "" cannot name a recording file`},
		{"slash", "a/b", `record id "a/b" cannot name a recording file`},
		{"backslash", `a\b`, `record id "a\\b" cannot name a recording file`},
		{"dot dot alone", "..", `record id ".." cannot name a recording file`},
		{"dot dot inside", "a..b", `record id "a..b" cannot name a recording file`},
		{"escaping", "../escape", `record id "../escape" cannot name a recording file`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "recordings")
			in := Recording{
				SedumVersion: Version,
				Packages:     fixture().Packages,
				Records:      []Record{recordOf(tc.id, "rails")},
			}

			err := WriteRecords(dir, in)
			if err == nil {
				t.Fatal("WriteRecords accepted the id")
			}
			if err.Error() != tc.wantText {
				t.Errorf("error = %q, want %q", err.Error(), tc.wantText)
			}

			// No file was written for it, in the directory or anywhere the id
			// could have pointed.
			var written []string
			_ = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
				if err == nil && !d.IsDir() {
					written = append(written, p)
				}
				return nil
			})
			if len(written) != 0 {
				t.Errorf("files were written for an id that cannot name one: %v", written)
			}
		})
	}
}

func TestWriteRecordsReportsAWriteFailure(t *testing.T) {
	dir := t.TempDir()
	// A directory where the file must go: the rename cannot replace it.
	path := filepath.Join(dir, "prov-2026-aaaaaaaa.json")
	if err := os.MkdirAll(filepath.Join(path, "keep"), 0o755); err != nil {
		t.Fatal(err)
	}
	in := Recording{SedumVersion: Version, Packages: fixture().Packages, Records: []Record{recordOf("prov-2026-aaaaaaaa", "rails")}}

	err := WriteRecords(dir, in)
	if err == nil {
		t.Fatal("WriteRecords over a directory succeeded")
	}
	if prefix := "writing recording " + path + ": "; !strings.HasPrefix(err.Error(), prefix) {
		t.Errorf("error = %q, want the prefix %q", err.Error(), prefix)
	}
}

// ReadRecord.

func TestReadRecordReturnsOnlyTheRequestedRecordWithTheFilesFields(t *testing.T) {
	dir := t.TempDir()
	region := regionOf("a.rb", "addField", "", map[string]any{"f": "1"})
	in := fixture()
	in.Records[0].Regions = []Region{region}
	if err := WriteRecords(dir, in); err != nil {
		t.Fatalf("WriteRecords: %v", err)
	}

	got, err := ReadRecord(dir, "prov-2026-aaaaaaaa")
	if err != nil {
		t.Fatalf("ReadRecord: %v", err)
	}
	if len(got.Records) != 1 || got.Records[0].RecordID != "prov-2026-aaaaaaaa" {
		t.Fatalf("Records = %+v, want only prov-2026-aaaaaaaa", got.Records)
	}
	if got.SedumVersion != "9.9.9" {
		t.Errorf("SedumVersion = %q", got.SedumVersion)
	}
	if want := map[string]Package{"rails": {Extensions: []string{".rb"}}}; !reflect.DeepEqual(got.Packages, want) {
		t.Errorf("Packages = %v, want %v", got.Packages, want)
	}
	if !reflect.DeepEqual(got.Variables, map[string]string{"service": "orders"}) {
		t.Errorf("Variables = %v", got.Variables)
	}
	if r := got.Records[0]; r.Gen != "sha256:gen-prov-2026-aaaaaaaa" || r.CreatedAt != "2026-01-02" ||
		len(r.Files) != 1 || len(r.Phases) != 1 || !reflect.DeepEqual(r.Regions, []Region{region}) {
		t.Errorf("record = %+v", r)
	}
}

func TestReadRecordKeepsOnlyTheRequestedRecordOfAMultiRecordFile(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	if err := Write(&buf, fixture()); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "prov-2026-bbbbbbbb.json"), buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := ReadRecord(dir, "prov-2026-bbbbbbbb")
	if err != nil {
		t.Fatalf("ReadRecord: %v", err)
	}
	if len(got.Records) != 1 || got.Records[0].RecordID != "prov-2026-bbbbbbbb" {
		t.Errorf("Records = %+v, want only prov-2026-bbbbbbbb", got.Records)
	}
	if len(got.Packages) != 3 {
		t.Errorf("Packages = %v, want the file's Packages as they are", got.Packages)
	}
}

func TestReadRecordOfAMissingFileSatisfiesErrNotExist(t *testing.T) {
	_, err := ReadRecord(t.TempDir(), "prov-2026-aaaaaaaa")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("error = %v, want one satisfying errors.Is(fs.ErrNotExist)", err)
	}

	_, err = ReadRecord(filepath.Join(t.TempDir(), "no", "such", "dir"), "prov-2026-aaaaaaaa")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("a missing directory: error = %v, want one satisfying errors.Is(fs.ErrNotExist)", err)
	}
}

func TestReadRecordReportsAFileThatCannotBeDecoded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prov-2026-aaaaaaaa.json")
	if err := os.WriteFile(path, []byte("this is not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := ReadRecord(dir, "prov-2026-aaaaaaaa")
	if err == nil {
		t.Fatal("ReadRecord decoded garbage")
	}
	_, readErr := Read(strings.NewReader("this is not json"))
	if want := path + ": " + readErr.Error(); err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
	if errors.Is(err, fs.ErrNotExist) {
		t.Error("an undecodable file reads as a missing one")
	}
}

func TestReadRecordReportsAFileHoldingNoSuchRecord(t *testing.T) {
	dir := t.TempDir()
	other := Recording{SedumVersion: Version, Packages: map[string]Package{}, Records: []Record{recordOf("prov-2026-bbbbbbbb", "rails")}}
	var buf bytes.Buffer
	if err := Write(&buf, other); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "prov-2026-aaaaaaaa.json")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := ReadRecord(dir, "prov-2026-aaaaaaaa")
	if err == nil {
		t.Fatal("ReadRecord returned a recording holding no such record")
	}
	if want := path + ` holds no record "prov-2026-aaaaaaaa"`; err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestReadRecordRefusesARecordIdThatCannotNameAFileAndOpensNothing(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "recordings")

	// Valid recordings sit where a traversing or nested id would reach, so a
	// ReadRecord that opened anything would succeed.
	var buf bytes.Buffer
	in := Recording{SedumVersion: Version, Packages: map[string]Package{}, Records: []Record{
		recordOf("../escape", "rails"), recordOf("a/b", "rails"), recordOf("a..b", "rails"), recordOf("", "rails"),
	}}
	if err := Write(&buf, in); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{
		filepath.Join(root, "escape.json"),
		filepath.Join(dir, "a", "b.json"),
		filepath.Join(dir, "a..b.json"),
		filepath.Join(dir, ".json"),
	} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, buf.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	for _, tc := range []struct{ id, want string }{
		{"", `record id "" cannot name a recording file`},
		{"a/b", `record id "a/b" cannot name a recording file`},
		{`a\b`, `record id "a\\b" cannot name a recording file`},
		{"../escape", `record id "../escape" cannot name a recording file`},
		{"a..b", `record id "a..b" cannot name a recording file`},
	} {
		_, err := ReadRecord(dir, tc.id)
		if err == nil {
			t.Errorf("ReadRecord(%q) succeeded", tc.id)
			continue
		}
		if err.Error() != tc.want {
			t.Errorf("ReadRecord(%q) error = %q, want %q", tc.id, err.Error(), tc.want)
		}
	}
}

// What WriteRecords wrote, ReadRecord reads back, and the file is a valid
// recording on its own.
func TestWriteThenReadRecordRoundTripsARecord(t *testing.T) {
	dir := t.TempDir()
	region := regionOf("a.rb", "addField", "", map[string]any{"f": "1"})
	in := fixture()
	in.Records[1].Regions = []Region{region}
	if err := WriteRecords(dir, in); err != nil {
		t.Fatal(err)
	}

	got, err := ReadRecord(dir, "prov-2026-bbbbbbbb")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Records[0].Regions, []Region{region}) {
		t.Errorf("regions = %#v", got.Records[0].Regions)
	}
	if got.Records[0].Gen != in.Records[1].Gen {
		t.Errorf("gen = %q", got.Records[0].Gen)
	}
}
