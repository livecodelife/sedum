package growlog

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/render_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/render_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSummaryText(t *testing.T) {\n\ts := Summarize(summaryEntries(), Filter{})\n\ts.Malformed = 1\n\ttext := s.Text()\n\tlines := strings.Split(strings.TrimSuffix(text, \"\\n\"), \"\\n\")\n\tif lines[0] != \"grow log: 4 line(s), 1 malformed\" {\n\t\tt.Fatalf(\"first line = %q\", lines[0])\n\t}\n\tif len(lines) != 6 {\n\t\tt.Fatalf(\"got %d lines, want 6:\\n%s\", len(lines), text)\n\t}\n\twant := map[string][]string{\n\t\t\"haiku\": {\"2\", \"1\", \"5\", \"4\", \"1\", \"150/25\", \"75%\"},\n\t\t\"stub\":  {\"1\", \"0\", \"1\", \"0\", \"0\", \"0/0\", \"0%\"},\n\t\t\"total\": {\"3\", \"1\", \"6\", \"4\", \"1\", \"150/25\", \"50%\"},\n\t}\n\tfor _, l := range lines[2:] {\n\t\tf := strings.Fields(l)\n\t\tif w, ok := want[f[0]]; ok {\n\t\t\tif strings.Join(f[1:], \" \") != strings.Join(w, \" \") {\n\t\t\t\tt.Errorf(\"row %s = %v, want %v\", f[0], f[1:], w)\n\t\t\t}\n\t\t\tdelete(want, f[0])\n\t\t}\n\t}\n\tif len(want) != 0 {\n\t\tt.Errorf(\"rows missing: %v\\n%s\", want, text)\n\t}\n\tif !strings.Contains(text, \"(none)\") {\n\t\tt.Errorf(\"a run with no model must show as (none):\\n%s\", text)\n\t}\n}","file":"internal/growlog/render_test.go","name":"TestSummaryText"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-13a4d6a9","gen":"sha256:2614c9df23ac9a504543127efa2248e92422009df5a5862f92fc7455a24bbfc2","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestSummaryJSON(t *testing.T) {\n\tb, err := Summarize(nil, Filter{}).JSON()\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\ts := string(b)\n\tif !strings.HasSuffix(s, \"\\n\") || !strings.Contains(s, `\"models\": []`) {\n\t\tt.Fatalf(\"JSON = %q\", s)\n\t}\n\tif strings.Contains(s, \"timestamp\") {\n\t\tt.Fatalf(\"a summary carries no timestamp: %q\", s)\n\t}\n}","file":"internal/growlog/render_test.go","name":"TestSummaryJSON"}}
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
