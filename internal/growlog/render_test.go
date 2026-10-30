package growlog

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/render_test.go","path":"strings"},"kwargs":"sha256:6044d8733e0b407295da5117daeb961975c8ce9bb48bf79809e44e1f160e80ef","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/render_test.go","path":"testing"},"kwargs":"sha256:f6c4b28e89c3f71030149cad8fa48b33b9b427c0dd40abdb410b679ce717aee3","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/render_test.go","name":"TestSummaryText"},"kwargs":"sha256:9f4567fef1208b2bf12d0a1e185d621f27af0943d003d49926ceb4176c787f09","body":"sha256:4717d6a92ff9c6c3fb38fb3d0d2e2877278f9ce1fa1832e73db004cb2612cff0"}
func TestSummaryText(t *testing.T) {
	s := Summarize(summaryEntries(), Filter{})
	s.Malformed = 1
	text := s.Text()
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	if lines[0] != "grow log: 4 line(s), 1 malformed" {
		t.Fatalf("first line = %q", lines[0])
	}
	if len(lines) != 6 {
		t.Fatalf("got %d lines, want 6:\n%s", len(lines), text)
	}
	want := map[string][]string{
		"haiku": {"2", "1", "5", "4", "1", "150/25", "75%"},
		"stub":  {"1", "0", "1", "0", "0", "0/0", "0%"},
		"total": {"3", "1", "6", "4", "1", "150/25", "50%"},
	}
	for _, l := range lines[2:] {
		f := strings.Fields(l)
		if w, ok := want[f[0]]; ok {
			if strings.Join(f[1:], " ") != strings.Join(w, " ") {
				t.Errorf("row %s = %v, want %v", f[0], f[1:], w)
			}
			delete(want, f[0])
		}
	}
	if len(want) != 0 {
		t.Errorf("rows missing: %v\n%s", want, text)
	}
	if !strings.Contains(text, "(none)") {
		t.Errorf("a run with no model must show as (none):\n%s", text)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/render_test.go","name":"TestSummaryJSON"},"kwargs":"sha256:75feabbebe9241e4a8d727d22e6517ffd62226110934e6b74316419081004655","body":"sha256:76f4d1f1f24d535cd59d06974966e575f8bfbb6d83b9f8a347637b7520f01e28"}
func TestSummaryJSON(t *testing.T) {
	b, err := Summarize(nil, Filter{}).JSON()
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.HasSuffix(s, "\n") || !strings.Contains(s, `"models": []`) {
		t.Fatalf("JSON = %q", s)
	}
	if strings.Contains(s, "timestamp") {
		t.Fatalf("a summary carries no timestamp: %q", s)
	}
}

// /sedum:setFunction
