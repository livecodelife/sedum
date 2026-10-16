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
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingTwoFiles(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\trec := &record.Record{ID: \"imp-1\", Gen: \"sha256:aa\", CreatedAt: \"2026-10-16\"}\n\titems := []Adopted{\n\t\t{Candidate: Candidate{Path: \"b.py\", Package: pkg, Template: \"t\", Record: rec}, Spans: []Span{\n\t\t\t{Kind: KindHeader, Text: \"# b\\n\"},\n\t\t\t{Kind: KindUnit, Text: \"b1\", Action: \"b1\", Kwargs: map[string]any{\"name\": \"b1\"}},\n\t\t\t{Kind: KindUnit, Text: \"b2\", Action: \"b2\", Kwargs: map[string]any{\"name\": \"b2\"}},\n\t\t}},\n\t\t{Candidate: Candidate{Path: \"a.py\", Package: pkg, Template: \"t\", Record: rec}, Spans: []Span{\n\t\t\t{Kind: KindHeader, Text: \"# a\\n\"},\n\t\t\t{Kind: KindUnit, Text: \"a1\", Action: \"a1\", Kwargs: map[string]any{\"name\": \"a1\"}},\n\t\t\t{Kind: KindUnit, Text: \"a2\", Action: \"a2\", Kwargs: map[string]any{\"name\": \"a2\"}},\n\t\t}},\n\t}\n\tgot, err := BuildRecording(items)\n\tif err != nil {\n\t\tt.Fatalf(\"BuildRecording: %v\", err)\n\t}\n\tif got.SedumVersion != recording.Version {\n\t\tt.Errorf(\"SedumVersion = %q, want %q\", got.SedumVersion, recording.Version)\n\t}\n\tif !reflect.DeepEqual(got.Packages[\"py\"].Extensions, []string{\".py\"}) {\n\t\tt.Errorf(\"Packages[py].Extensions = %v\", got.Packages[\"py\"].Extensions)\n\t}\n\tif len(got.Records) != 1 {\n\t\tt.Fatalf(\"got %d records, want 1\", len(got.Records))\n\t}\n\tr := got.Records[0]\n\tif r.RecordID != \"imp-1\" || r.Gen != \"sha256:aa\" || r.CreatedAt != \"2026-10-16\" {\n\t\tt.Errorf(\"record = %+v\", r)\n\t}\n\tif len(r.Files) != 2 || r.Files[0].Path != \"a.py\" || r.Files[1].Path != \"b.py\" {\n\t\tt.Fatalf(\"files = %+v, want a.py then b.py\", r.Files)\n\t}\n\tif len(r.Phases) != 1 || r.Phases[0].Name != recording.DefaultPhase {\n\t\tt.Fatalf(\"phases = %+v\", r.Phases)\n\t}\n\tvar acts []string\n\tfor _, inv := range r.Phases[0].Invocations {\n\t\tacts = append(acts, inv.Action)\n\t}\n\tif !reflect.DeepEqual(acts, []string{\"a1\", \"a2\", \"b1\", \"b2\"}) {\n\t\tt.Errorf(\"actions = %v\", acts)\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingTwoFiles"}}
func TestBuildRecordingTwoFiles(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"}
	items := []Adopted{
		{Candidate: Candidate{Path: "b.py", Package: pkg, Template: "t", Record: rec}, Spans: []Span{
			{Kind: KindHeader, Text: "# b\n"},
			{Kind: KindUnit, Text: "b1", Action: "b1", Kwargs: map[string]any{"name": "b1"}},
			{Kind: KindUnit, Text: "b2", Action: "b2", Kwargs: map[string]any{"name": "b2"}},
		}},
		{Candidate: Candidate{Path: "a.py", Package: pkg, Template: "t", Record: rec}, Spans: []Span{
			{Kind: KindHeader, Text: "# a\n"},
			{Kind: KindUnit, Text: "a1", Action: "a1", Kwargs: map[string]any{"name": "a1"}},
			{Kind: KindUnit, Text: "a2", Action: "a2", Kwargs: map[string]any{"name": "a2"}},
		}},
	}
	got, err := BuildRecording(items)
	if err != nil {
		t.Fatalf("BuildRecording: %v", err)
	}
	if got.SedumVersion != recording.Version {
		t.Errorf("SedumVersion = %q, want %q", got.SedumVersion, recording.Version)
	}
	if !reflect.DeepEqual(got.Packages["py"].Extensions, []string{".py"}) {
		t.Errorf("Packages[py].Extensions = %v", got.Packages["py"].Extensions)
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
	var acts []string
	for _, inv := range r.Phases[0].Invocations {
		acts = append(acts, inv.Action)
	}
	if !reflect.DeepEqual(acts, []string{"a1", "a2", "b1", "b2"}) {
		t.Errorf("actions = %v", acts)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingTwoRecords(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\tr1 := &record.Record{ID: \"imp-1\"}\n\tr2 := &record.Record{ID: \"imp-2\"}\n\titems := []Adopted{\n\t\t{Candidate: Candidate{Path: \"z.py\", Package: pkg, Record: r2}, Spans: []Span{{Kind: KindUnit, Action: \"z\"}}},\n\t\t{Candidate: Candidate{Path: \"a.py\", Package: pkg, Record: r1}, Spans: []Span{{Kind: KindUnit, Action: \"a\"}}},\n\t}\n\tgot, err := BuildRecording(items)\n\tif err != nil {\n\t\tt.Fatalf(\"BuildRecording: %v\", err)\n\t}\n\tif len(got.Records) != 2 {\n\t\tt.Fatalf(\"got %d records, want 2\", len(got.Records))\n\t}\n\tif got.Records[0].RecordID != \"imp-1\" || got.Records[1].RecordID != \"imp-2\" {\n\t\tt.Errorf(\"record order = %s, %s\", got.Records[0].RecordID, got.Records[1].RecordID)\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingTwoRecords"}}
func TestBuildRecordingTwoRecords(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	r1 := &record.Record{ID: "imp-1"}
	r2 := &record.Record{ID: "imp-2"}
	items := []Adopted{
		{Candidate: Candidate{Path: "z.py", Package: pkg, Record: r2}, Spans: []Span{{Kind: KindUnit, Action: "z"}}},
		{Candidate: Candidate{Path: "a.py", Package: pkg, Record: r1}, Spans: []Span{{Kind: KindUnit, Action: "a"}}},
	}
	got, err := BuildRecording(items)
	if err != nil {
		t.Fatalf("BuildRecording: %v", err)
	}
	if len(got.Records) != 2 {
		t.Fatalf("got %d records, want 2", len(got.Records))
	}
	if got.Records[0].RecordID != "imp-1" || got.Records[1].RecordID != "imp-2" {
		t.Errorf("record order = %s, %s", got.Records[0].RecordID, got.Records[1].RecordID)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingErrors(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\trec := &record.Record{ID: \"imp-1\"}\n\ttests := []struct {\n\t\tname  string\n\t\titems []Adopted\n\t}{\n\t\t{\n\t\t\tname: \"no unit spans\",\n\t\t\titems: []Adopted{\n\t\t\t\t{Candidate: Candidate{Path: \"a.py\", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindHeader, Text: \"x\"}}},\n\t\t\t},\n\t\t},\n\t\t{\n\t\t\tname: \"duplicate path\",\n\t\t\titems: []Adopted{\n\t\t\t\t{Candidate: Candidate{Path: \"a.py\", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindUnit, Action: \"a\"}}},\n\t\t\t\t{Candidate: Candidate{Path: \"a.py\", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindUnit, Action: \"a\"}}},\n\t\t\t},\n\t\t},\n\t}\n\tfor _, tc := range tests {\n\t\tif _, err := BuildRecording(tc.items); err == nil {\n\t\t\tt.Errorf(\"%s: want error, got nil\", tc.name)\n\t\t}\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingErrors"}}
func TestBuildRecordingErrors(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1"}
	tests := []struct {
		name  string
		items []Adopted
	}{
		{
			name: "no unit spans",
			items: []Adopted{
				{Candidate: Candidate{Path: "a.py", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindHeader, Text: "x"}}},
			},
		},
		{
			name: "duplicate path",
			items: []Adopted{
				{Candidate: Candidate{Path: "a.py", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindUnit, Action: "a"}}},
				{Candidate: Candidate{Path: "a.py", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindUnit, Action: "a"}}},
			},
		},
	}
	for _, tc := range tests {
		if _, err := BuildRecording(tc.items); err == nil {
			t.Errorf("%s: want error, got nil", tc.name)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingRoundTripsThroughWrite(t *testing.T) {\n\tpkg := &genpkg.Package{Name: \"py\", Extensions: []string{\".py\"}}\n\trec := &record.Record{ID: \"imp-1\", Gen: \"sha256:aa\", CreatedAt: \"2026-10-16\"}\n\titems := []Adopted{\n\t\t{Candidate: Candidate{Path: \"b.py\", Package: pkg, Template: \"t\", Record: rec}, Spans: []Span{\n\t\t\t{Kind: KindHeader, Text: \"# b\\n\"},\n\t\t\t{Kind: KindUnit, Text: \"b1\", Action: \"b1\", Kwargs: map[string]any{\"name\": \"b1\"}},\n\t\t\t{Kind: KindUnit, Text: \"b2\", Action: \"b2\", Kwargs: map[string]any{\"name\": \"b2\"}},\n\t\t}},\n\t\t{Candidate: Candidate{Path: \"a.py\", Package: pkg, Template: \"t\", Record: rec}, Spans: []Span{\n\t\t\t{Kind: KindHeader, Text: \"# a\\n\"},\n\t\t\t{Kind: KindUnit, Text: \"a1\", Action: \"a1\", Kwargs: map[string]any{\"name\": \"a1\"}},\n\t\t\t{Kind: KindUnit, Text: \"a2\", Action: \"a2\", Kwargs: map[string]any{\"name\": \"a2\"}},\n\t\t}},\n\t}\n\twant, err := BuildRecording(items)\n\tif err != nil {\n\t\tt.Fatalf(\"BuildRecording: %v\", err)\n\t}\n\tvar buf bytes.Buffer\n\tif err := recording.Write(&buf, want); err != nil {\n\t\tt.Fatalf(\"Write: %v\", err)\n\t}\n\tgot, err := recording.Read(&buf)\n\tif err != nil {\n\t\tt.Fatalf(\"Read: %v\", err)\n\t}\n\tif !reflect.DeepEqual(got, want) {\n\t\tt.Errorf(\"round trip mismatch:\\n got %+v\\nwant %+v\", got, want)\n\t}\n}","file":"internal/adopt/build_test.go","name":"TestBuildRecordingRoundTripsThroughWrite"}}
func TestBuildRecordingRoundTripsThroughWrite(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"}
	items := []Adopted{
		{Candidate: Candidate{Path: "b.py", Package: pkg, Template: "t", Record: rec}, Spans: []Span{
			{Kind: KindHeader, Text: "# b\n"},
			{Kind: KindUnit, Text: "b1", Action: "b1", Kwargs: map[string]any{"name": "b1"}},
			{Kind: KindUnit, Text: "b2", Action: "b2", Kwargs: map[string]any{"name": "b2"}},
		}},
		{Candidate: Candidate{Path: "a.py", Package: pkg, Template: "t", Record: rec}, Spans: []Span{
			{Kind: KindHeader, Text: "# a\n"},
			{Kind: KindUnit, Text: "a1", Action: "a1", Kwargs: map[string]any{"name": "a1"}},
			{Kind: KindUnit, Text: "a2", Action: "a2", Kwargs: map[string]any{"name": "a2"}},
		}},
	}
	want, err := BuildRecording(items)
	if err != nil {
		t.Fatalf("BuildRecording: %v", err)
	}
	var buf bytes.Buffer
	if err := recording.Write(&buf, want); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := recording.Read(&buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("round trip mismatch:\n got %+v\nwant %+v", got, want)
	}
}

// /sedum:setFunction
