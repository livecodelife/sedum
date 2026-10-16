package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_header_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_header_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/build_header_test.go","path":"github.com/livecodelife/sedum/internal/record"}}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:0ba9679009ce9401281269fe297403e8a58066d05e01c7d30e9e4a87bb5aaf30","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestBuildRecordingHeaderCaptures(t *testing.T) {\n\tsrc := map[string]string{\"pkg\": \"dir\", \"name\": \"a\"}\n\titem := Adopted{\n\t\tCandidate: Candidate{\n\t\t\tPath:     \"a.go\",\n\t\t\tPackage:  &genpkg.Package{Name: \"go\", Extensions: []string{\".go\"}},\n\t\t\tRecord:   &record.Record{ID: \"imp-1\", Gen: \"sha256:aa\", CreatedAt: \"2026-10-16\"},\n\t\t\tCaptures: src,\n\t\t},\n\t\tSpans: []Span{\n\t\t\t{Kind: KindHeader, Text: \"package foo\\n\", Kwargs: map[string]any{\"pkg\": \"foo\"}},\n\t\t\t{Kind: KindUnit, Text: \"func F() {}\\n\", Action: \"setFunction\", Kwargs: map[string]any{\"file\": \"a.go\", \"name\": \"F\", \"code\": \"func F() {}\"}},\n\t\t},\n\t}\n\tgot, err := BuildRecording([]Adopted{item})\n\tif err != nil {\n\t\tt.Fatalf(\"BuildRecording: %v\", err)\n\t}\n\tcaps := got.Records[0].Files[0].Captures\n\twant := map[string]string{\"pkg\": \"foo\", \"name\": \"a\"}\n\tfor k, v := range want {\n\t\tif caps[k] != v {\n\t\t\tt.Errorf(\"Captures[%s] = %q, want %q\", k, caps[k], v)\n\t\t}\n\t}\n\tif len(caps) != len(want) {\n\t\tt.Errorf(\"got %d captures, want %d\", len(caps), len(want))\n\t}\n\tif src[\"pkg\"] != \"dir\" {\n\t\tt.Errorf(\"candidate captures modified: pkg = %q\", src[\"pkg\"])\n\t}\n}","file":"internal/adopt/build_header_test.go","name":"TestBuildRecordingHeaderCaptures"}}
func TestBuildRecordingHeaderCaptures(t *testing.T) {
	src := map[string]string{"pkg": "dir", "name": "a"}
	item := Adopted{
		Candidate: Candidate{
			Path:     "a.go",
			Package:  &genpkg.Package{Name: "go", Extensions: []string{".go"}},
			Record:   &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"},
			Captures: src,
		},
		Spans: []Span{
			{Kind: KindHeader, Text: "package foo\n", Kwargs: map[string]any{"pkg": "foo"}},
			{Kind: KindUnit, Text: "func F() {}\n", Action: "setFunction", Kwargs: map[string]any{"file": "a.go", "name": "F", "code": "func F() {}"}},
		},
	}
	got, err := BuildRecording([]Adopted{item})
	if err != nil {
		t.Fatalf("BuildRecording: %v", err)
	}
	caps := got.Records[0].Files[0].Captures
	want := map[string]string{"pkg": "foo", "name": "a"}
	for k, v := range want {
		if caps[k] != v {
			t.Errorf("Captures[%s] = %q, want %q", k, caps[k], v)
		}
	}
	if len(caps) != len(want) {
		t.Errorf("got %d captures, want %d", len(caps), len(want))
	}
	if src["pkg"] != "dir" {
		t.Errorf("candidate captures modified: pkg = %q", src["pkg"])
	}
}

// /sedum:setFunction
