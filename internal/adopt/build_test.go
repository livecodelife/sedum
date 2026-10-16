package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"reflect"}}
import "reflect"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"github.com/livecodelife/sedum/internal/record"}}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_test.go","path":"github.com/livecodelife/sedum/internal/recording"}}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingTwoFiles(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\trec := &record.Record{ID: \"imp-1\", Gen: \"sha256:aa\", CreatedAt: \"2026-10-16\"}\n\tmk := func(path string) Adopted {\n\t\treturn Adopted{\n\t\t\tCandidate: Candidate{Path: path, Package: pkg, Template: \"t\", Record: rec},\n\t\t\tSpans: []Span{\n\t\t\t\t{Kind: KindHeader, Text: \"# header\\n\"},\n\t\t\t\t{Kind: KindUnit, Text: \"a\", Action: \"one\", Kwargs: map[string]any{\"path\": path}},\n\t\t\t\t{Kind: KindUnit, Text: \"b\", Action: \"two\", Kwargs: map[string]any{\"path\": path}},\n\t\t\t},\n\t\t}\n\t}\n\tgot, err := BuildRecording([]Adopted{mk(\"b.py\"), mk(\"a.py\")})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif got.SedumVersion != recording.Version {\n\t\tt.Errorf(\"version = %q\", got.SedumVersion)\n\t}\n\tif !reflect.DeepEqual(got.Packages[\"py\"].Extensions, []string{\".py\"}) {\n\t\tt.Errorf(\"extensions = %v\", got.Packages[\"py\"].Extensions)\n\t}\n\tif len(got.Records) != 1 {\n\t\tt.Fatalf(\"records = %d, want 1\", len(got.Records))\n\t}\n\tr := got.Records[0]\n\tif r.RecordID != \"imp-1\" || r.Gen != \"sha256:aa\" || r.CreatedAt != \"2026-10-16\" {\n\t\tt.Errorf(\"record = %+v\", r)\n\t}\n\tif len(r.Files) != 2 || r.Files[0].Path != \"a.py\" || r.Files[1].Path != \"b.py\" {\n\t\tt.Fatalf(\"files = %+v\", r.Files)\n\t}\n\tinv := r.Phases[0].Invocations\n\tif len(inv) != 4 {\n\t\tt.Fatalf(\"invocations = %d, want 4\", len(inv))\n\t}\n\twant := []struct{ path, action string }{{\"a.py\", \"one\"}, {\"a.py\", \"two\"}, {\"b.py\", \"one\"}, {\"b.py\", \"two\"}}\n\tfor i, w := range want {\n\t\tif inv[i].Action != w.action || inv[i].Kwargs[\"path\"] != w.path {\n\t\t\tt.Errorf(\"invocation %d = %+v\", i, inv[i])\n\t\t}\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingTwoFiles"}}
func TestBuildRecordingTwoFiles(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"}
	mk := func(path string) Adopted {
		return Adopted{
			Candidate: Candidate{Path: path, Package: pkg, Template: "t", Record: rec},
			Spans: []Span{
				{Kind: KindHeader, Text: "# header\n"},
				{Kind: KindUnit, Text: "a", Action: "one", Kwargs: map[string]any{"path": path}},
				{Kind: KindUnit, Text: "b", Action: "two", Kwargs: map[string]any{"path": path}},
			},
		}
	}
	got, err := BuildRecording([]Adopted{mk("b.py"), mk("a.py")})
	if err != nil {
		t.Fatal(err)
	}
	if got.SedumVersion != recording.Version {
		t.Errorf("version = %q", got.SedumVersion)
	}
	if !reflect.DeepEqual(got.Packages["py"].Extensions, []string{".py"}) {
		t.Errorf("extensions = %v", got.Packages["py"].Extensions)
	}
	if len(got.Records) != 1 {
		t.Fatalf("records = %d, want 1", len(got.Records))
	}
	r := got.Records[0]
	if r.RecordID != "imp-1" || r.Gen != "sha256:aa" || r.CreatedAt != "2026-10-16" {
		t.Errorf("record = %+v", r)
	}
	if len(r.Files) != 2 || r.Files[0].Path != "a.py" || r.Files[1].Path != "b.py" {
		t.Fatalf("files = %+v", r.Files)
	}
	inv := r.Phases[0].Invocations
	if len(inv) != 4 {
		t.Fatalf("invocations = %d, want 4", len(inv))
	}
	want := []struct{ path, action string }{{"a.py", "one"}, {"a.py", "two"}, {"b.py", "one"}, {"b.py", "two"}}
	for i, w := range want {
		if inv[i].Action != w.action || inv[i].Kwargs["path"] != w.path {
			t.Errorf("invocation %d = %+v", i, inv[i])
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingTwoRecords(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\tunit := []Span{{Kind: KindUnit, Text: \"x\", Action: \"one\"}}\n\titems := []Adopted{\n\t\t{Candidate: Candidate{Path: \"z.py\", Package: pkg, Record: &record.Record{ID: \"imp-2\"}}, Spans: unit},\n\t\t{Candidate: Candidate{Path: \"a.py\", Package: pkg, Record: &record.Record{ID: \"imp-1\"}}, Spans: unit},\n\t}\n\tgot, err := BuildRecording(items)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(got.Records) != 2 || got.Records[0].RecordID != \"imp-1\" || got.Records[1].RecordID != \"imp-2\" {\n\t\tt.Errorf(\"records = %+v\", got.Records)\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingTwoRecords"}}
func TestBuildRecordingTwoRecords(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	unit := []Span{{Kind: KindUnit, Text: "x", Action: "one"}}
	items := []Adopted{
		{Candidate: Candidate{Path: "z.py", Package: pkg, Record: &record.Record{ID: "imp-2"}}, Spans: unit},
		{Candidate: Candidate{Path: "a.py", Package: pkg, Record: &record.Record{ID: "imp-1"}}, Spans: unit},
	}
	got, err := BuildRecording(items)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Records) != 2 || got.Records[0].RecordID != "imp-1" || got.Records[1].RecordID != "imp-2" {
		t.Errorf("records = %+v", got.Records)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingErrors(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\trec := &record.Record{ID: \"imp-1\"}\n\tunit := Span{Kind: KindUnit, Text: \"x\", Action: \"one\"}\n\ttests := []struct {\n\t\tname  string\n\t\titems []Adopted\n\t}{\n\t\t{\"no unit spans\", []Adopted{{Candidate: Candidate{Path: \"a.py\", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindHeader, Text: \"h\"}}}}},\n\t\t{\"duplicate path\", []Adopted{\n\t\t\t{Candidate: Candidate{Path: \"a.py\", Package: pkg, Record: rec}, Spans: []Span{unit}},\n\t\t\t{Candidate: Candidate{Path: \"a.py\", Package: pkg, Record: rec}, Spans: []Span{unit}},\n\t\t}},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\tif _, err := BuildRecording(tt.items); err == nil {\n\t\t\t\tt.Fatal(\"want error\")\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingErrors"}}
func TestBuildRecordingErrors(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1"}
	unit := Span{Kind: KindUnit, Text: "x", Action: "one"}
	tests := []struct {
		name  string
		items []Adopted
	}{
		{"no unit spans", []Adopted{{Candidate: Candidate{Path: "a.py", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindHeader, Text: "h"}}}}},
		{"duplicate path", []Adopted{
			{Candidate: Candidate{Path: "a.py", Package: pkg, Record: rec}, Spans: []Span{unit}},
			{Candidate: Candidate{Path: "a.py", Package: pkg, Record: rec}, Spans: []Span{unit}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := BuildRecording(tt.items); err == nil {
				t.Fatal("want error")
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingRoundTripsThroughWrite(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\trec := &record.Record{ID: \"imp-1\", Gen: \"sha256:aa\", CreatedAt: \"2026-10-16\"}\n\tmk := func(path string) Adopted {\n\t\treturn Adopted{\n\t\t\tCandidate: Candidate{Path: path, Package: pkg, Template: \"t\", Record: rec},\n\t\t\tSpans: []Span{\n\t\t\t\t{Kind: KindHeader, Text: \"# header\\n\"},\n\t\t\t\t{Kind: KindUnit, Text: \"a\", Action: \"one\", Kwargs: map[string]any{\"path\": path}},\n\t\t\t\t{Kind: KindUnit, Text: \"b\", Action: \"two\", Kwargs: map[string]any{\"path\": path}},\n\t\t\t},\n\t\t}\n\t}\n\tgot, err := BuildRecording([]Adopted{mk(\"b.py\"), mk(\"a.py\")})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tvar buf bytes.Buffer\n\tif err := recording.Write(&buf, got); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tback, err := recording.Read(&buf)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif !reflect.DeepEqual(back, got) {\n\t\tt.Errorf(\"round trip differs:\\n got %+v\\nwant %+v\", back, got)\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingRoundTripsThroughWrite"}}
func TestBuildRecordingRoundTripsThroughWrite(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"}
	mk := func(path string) Adopted {
		return Adopted{
			Candidate: Candidate{Path: path, Package: pkg, Template: "t", Record: rec},
			Spans: []Span{
				{Kind: KindHeader, Text: "# header\n"},
				{Kind: KindUnit, Text: "a", Action: "one", Kwargs: map[string]any{"path": path}},
				{Kind: KindUnit, Text: "b", Action: "two", Kwargs: map[string]any{"path": path}},
			},
		}
	}
	got, err := BuildRecording([]Adopted{mk("b.py"), mk("a.py")})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := recording.Write(&buf, got); err != nil {
		t.Fatal(err)
	}
	back, err := recording.Read(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, got) {
		t.Errorf("round trip differs:\n got %+v\nwant %+v", back, got)
	}
}

// /sedum:setFunction
