package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"reflect"}}
import "reflect"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"github.com/livecodelife/sedum/internal/record"}}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"github.com/livecodelife/sedum/internal/recording"}}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingTwoFiles(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\trec := &record.Record{ID: \"imp-1\", Gen: \"sha256:aa\", CreatedAt: \"2026-10-16\"}\n\titems := []Adopted{\n\t\tadoptedItem(\"b.py\", pkg, rec, \"B1\", \"B2\"),\n\t\tadoptedItem(\"a.py\", pkg, rec, \"A1\", \"A2\"),\n\t}\n\tgot, err := BuildRecording(items)\n\tif err != nil {\n\t\tt.Fatalf(\"BuildRecording: %v\", err)\n\t}\n\tif got.SedumVersion != recording.Version {\n\t\tt.Errorf(\"SedumVersion = %q, want %q\", got.SedumVersion, recording.Version)\n\t}\n\tif want := []string{\".py\"}; !reflect.DeepEqual(got.Packages[\"py\"].Extensions, want) {\n\t\tt.Errorf(\"Packages[py].Extensions = %v, want %v\", got.Packages[\"py\"].Extensions, want)\n\t}\n\tif len(got.Records) != 1 {\n\t\tt.Fatalf(\"got %d records, want 1\", len(got.Records))\n\t}\n\tr := got.Records[0]\n\tif r.RecordID != \"imp-1\" || r.Gen != \"sha256:aa\" || r.CreatedAt != \"2026-10-16\" {\n\t\tt.Errorf(\"record = %+v\", r)\n\t}\n\tif len(r.Files) != 2 || r.Files[0].Path != \"a.py\" || r.Files[1].Path != \"b.py\" {\n\t\tt.Fatalf(\"files = %+v, want a.py then b.py\", r.Files)\n\t}\n\tif len(r.Phases) != 1 || r.Phases[0].Name != recording.DefaultPhase {\n\t\tt.Fatalf(\"phases = %+v\", r.Phases)\n\t}\n\tinv := r.Phases[0].Invocations\n\twant := []string{\"A1\", \"A2\", \"B1\", \"B2\"}\n\tif len(inv) != len(want) {\n\t\tt.Fatalf(\"got %d invocations, want %d\", len(inv), len(want))\n\t}\n\tfor i, w := range want {\n\t\tif inv[i].Action != \"setFunction\" || inv[i].Kwargs[\"name\"] != w {\n\t\t\tt.Errorf(\"invocation %d = %+v, want name %s\", i, inv[i], w)\n\t\t}\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingTwoFiles"}}
func TestBuildRecordingTwoFiles(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"}
	items := []Adopted{
		adoptedItem("b.py", pkg, rec, "B1", "B2"),
		adoptedItem("a.py", pkg, rec, "A1", "A2"),
	}
	got, err := BuildRecording(items)
	if err != nil {
		t.Fatalf("BuildRecording: %v", err)
	}
	if got.SedumVersion != recording.Version {
		t.Errorf("SedumVersion = %q, want %q", got.SedumVersion, recording.Version)
	}
	if want := []string{".py"}; !reflect.DeepEqual(got.Packages["py"].Extensions, want) {
		t.Errorf("Packages[py].Extensions = %v, want %v", got.Packages["py"].Extensions, want)
	}
	if len(got.Records) != 1 {
		t.Fatalf("got %d records, want 1", len(got.Records))
	}
	r := got.Records[0]
	if r.RecordID != "imp-1" || r.Gen != "sha256:aa" || r.CreatedAt != "2026-10-16" {
		t.Errorf("record = %+v", r)
	}
	if len(r.Files) != 2 || r.Files[0].Path != "a.py" || r.Files[1].Path != "b.py" {
		t.Fatalf("files = %+v, want a.py then b.py", r.Files)
	}
	if len(r.Phases) != 1 || r.Phases[0].Name != recording.DefaultPhase {
		t.Fatalf("phases = %+v", r.Phases)
	}
	inv := r.Phases[0].Invocations
	want := []string{"A1", "A2", "B1", "B2"}
	if len(inv) != len(want) {
		t.Fatalf("got %d invocations, want %d", len(inv), len(want))
	}
	for i, w := range want {
		if inv[i].Action != "setFunction" || inv[i].Kwargs["name"] != w {
			t.Errorf("invocation %d = %+v, want name %s", i, inv[i], w)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingTwoRecords(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\titems := []Adopted{\n\t\tadoptedItem(\"a.py\", pkg, &record.Record{ID: \"imp-1\", Gen: \"sha256:aa\", CreatedAt: \"2026-10-16\"}, \"A\"),\n\t\tadoptedItem(\"z.py\", pkg, &record.Record{ID: \"imp-2\", Gen: \"sha256:bb\", CreatedAt: \"2026-10-16\"}, \"Z\"),\n\t}\n\tgot, err := BuildRecording(items)\n\tif err != nil {\n\t\tt.Fatalf(\"BuildRecording: %v\", err)\n\t}\n\tif len(got.Records) != 2 {\n\t\tt.Fatalf(\"got %d records, want 2\", len(got.Records))\n\t}\n\tif got.Records[0].RecordID != \"imp-1\" || got.Records[1].RecordID != \"imp-2\" {\n\t\tt.Errorf(\"record order = %s, %s; want imp-1, imp-2\", got.Records[0].RecordID, got.Records[1].RecordID)\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingTwoRecords"}}
func TestBuildRecordingTwoRecords(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	items := []Adopted{
		adoptedItem("a.py", pkg, &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"}, "A"),
		adoptedItem("z.py", pkg, &record.Record{ID: "imp-2", Gen: "sha256:bb", CreatedAt: "2026-10-16"}, "Z"),
	}
	got, err := BuildRecording(items)
	if err != nil {
		t.Fatalf("BuildRecording: %v", err)
	}
	if len(got.Records) != 2 {
		t.Fatalf("got %d records, want 2", len(got.Records))
	}
	if got.Records[0].RecordID != "imp-1" || got.Records[1].RecordID != "imp-2" {
		t.Errorf("record order = %s, %s; want imp-1, imp-2", got.Records[0].RecordID, got.Records[1].RecordID)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingErrors(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\trec := &record.Record{ID: \"imp-1\", Gen: \"sha256:aa\", CreatedAt: \"2026-10-16\"}\n\tcases := []struct {\n\t\tname  string\n\t\titems []Adopted\n\t}{\n\t\t{\"no unit spans\", []Adopted{{Candidate: Candidate{Path: \"c.py\", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindHeader, Text: \"x\"}}}}},\n\t\t{\"duplicate path\", []Adopted{adoptedItem(\"d.py\", pkg, rec, \"D\"), adoptedItem(\"d.py\", pkg, rec, \"E\")}},\n\t}\n\tfor _, tc := range cases {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\tif _, err := BuildRecording(tc.items); err == nil {\n\t\t\t\tt.Error(\"want error, got nil\")\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingErrors"}}
func TestBuildRecordingErrors(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"}
	cases := []struct {
		name  string
		items []Adopted
	}{
		{"no unit spans", []Adopted{{Candidate: Candidate{Path: "c.py", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindHeader, Text: "x"}}}}},
		{"duplicate path", []Adopted{adoptedItem("d.py", pkg, rec, "D"), adoptedItem("d.py", pkg, rec, "E")}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := BuildRecording(tc.items); err == nil {
				t.Error("want error, got nil")
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingRoundTripsThroughWrite(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\trec := &record.Record{ID: \"imp-1\", Gen: \"sha256:aa\", CreatedAt: \"2026-10-16\"}\n\titems := []Adopted{\n\t\tadoptedItem(\"b.py\", pkg, rec, \"B1\", \"B2\"),\n\t\tadoptedItem(\"a.py\", pkg, rec, \"A1\", \"A2\"),\n\t}\n\tgot, err := BuildRecording(items)\n\tif err != nil {\n\t\tt.Fatalf(\"BuildRecording: %v\", err)\n\t}\n\tvar buf bytes.Buffer\n\tif err := recording.Write(&buf, got); err != nil {\n\t\tt.Fatalf(\"Write: %v\", err)\n\t}\n\tback, err := recording.Read(&buf)\n\tif err != nil {\n\t\tt.Fatalf(\"Read: %v\", err)\n\t}\n\tif !reflect.DeepEqual(back, got) {\n\t\tt.Errorf(\"round trip mismatch:\\n got %+v\\nwant %+v\", back, got)\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingRoundTripsThroughWrite"}}
func TestBuildRecordingRoundTripsThroughWrite(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"}
	items := []Adopted{
		adoptedItem("b.py", pkg, rec, "B1", "B2"),
		adoptedItem("a.py", pkg, rec, "A1", "A2"),
	}
	got, err := BuildRecording(items)
	if err != nil {
		t.Fatalf("BuildRecording: %v", err)
	}
	var buf bytes.Buffer
	if err := recording.Write(&buf, got); err != nil {
		t.Fatalf("Write: %v", err)
	}
	back, err := recording.Read(&buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if !reflect.DeepEqual(back, got) {
		t.Errorf("round trip mismatch:\n got %+v\nwant %+v", back, got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func adoptedItem(path string, pkg *genpkg.Package, rec *record.Record, names ...string) Adopted {\n\tspans := []Span{{Kind: KindHeader, Text: \"// header\\n\"}}\n\tfor _, n := range names {\n\t\tspans = append(spans, Span{Kind: KindUnit, Text: \"func \" + n + \"() {}\\n\", Action: \"setFunction\", Kwargs: map[string]any{\"file\": path, \"name\": n, \"code\": \"func \" + n + \"() {}\"}})\n\t}\n\treturn Adopted{Candidate: Candidate{Path: path, Package: pkg, Template: \"tmpl\", Record: rec}, Spans: spans}\n}","file":"internal/adopt/build_test.go","name":"adoptedItem"}}
func adoptedItem(path string, pkg *genpkg.Package, rec *record.Record, names ...string) Adopted {
	spans := []Span{{Kind: KindHeader, Text: "// header\n"}}
	for _, n := range names {
		spans = append(spans, Span{Kind: KindUnit, Text: "func " + n + "() {}\n", Action: "setFunction", Kwargs: map[string]any{"file": path, "name": n, "code": "func " + n + "() {}"}})
	}
	return Adopted{Candidate: Candidate{Path: path, Package: pkg, Template: "tmpl", Record: rec}, Spans: spans}
}

// /sedum:setFunction
