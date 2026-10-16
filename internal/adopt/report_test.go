package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestReportTotals(t *testing.T) {\n\tr := NewReport()\n\tr.Add(FileReport{Path: \"a.go\", Status: StatusAdopted, Units: 3})\n\tr.Add(FileReport{Path: \"b.go\", Status: StatusSkipped, Reason: ReasonHeader})\n\tr.Add(FileReport{Path: \"c.go\", Status: StatusSkipped, Reason: ReasonHeader})\n\tr.Add(FileReport{Path: \"d.go\", Status: StatusRejected, Reason: ReasonOtherResidual})\n\ttot := r.Totals\n\tif tot.Files != 4 || tot.Adopted != 1 || tot.Skipped != 2 || tot.Rejected != 1 || tot.Units != 3 {\n\t\tt.Errorf(\"totals = %+v\", tot)\n\t}\n\twant := map[Reason]int{ReasonHeader: 2, ReasonOtherResidual: 1}\n\tif len(tot.SkippedByReason) != len(want) {\n\t\tt.Errorf(\"SkippedByReason = %v, want %v\", tot.SkippedByReason, want)\n\t}\n\tfor k, v := range want {\n\t\tif tot.SkippedByReason[k] != v {\n\t\t\tt.Errorf(\"SkippedByReason[%s] = %d, want %d\", k, tot.SkippedByReason[k], v)\n\t\t}\n\t}\n}","file":"internal/adopt/report_test.go","name":"TestReportTotals"}}
func TestReportTotals(t *testing.T) {
	r := NewReport()
	r.Add(FileReport{Path: "a.go", Status: StatusAdopted, Units: 3})
	r.Add(FileReport{Path: "b.go", Status: StatusSkipped, Reason: ReasonHeader})
	r.Add(FileReport{Path: "c.go", Status: StatusSkipped, Reason: ReasonHeader})
	r.Add(FileReport{Path: "d.go", Status: StatusRejected, Reason: ReasonOtherResidual})
	tot := r.Totals
	if tot.Files != 4 || tot.Adopted != 1 || tot.Skipped != 2 || tot.Rejected != 1 || tot.Units != 3 {
		t.Errorf("totals = %+v", tot)
	}
	want := map[Reason]int{ReasonHeader: 2, ReasonOtherResidual: 1}
	if len(tot.SkippedByReason) != len(want) {
		t.Errorf("SkippedByReason = %v, want %v", tot.SkippedByReason, want)
	}
	for k, v := range want {
		if tot.SkippedByReason[k] != v {
			t.Errorf("SkippedByReason[%s] = %d, want %d", k, tot.SkippedByReason[k], v)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestReportJSONKeys(t *testing.T) {\n\tr := NewReport()\n\tr.Add(FileReport{Path: \"a.go\", Status: StatusAdopted, Units: 1})\n\tb, err := r.JSON()\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\ts := string(b)\n\tfor _, want := range []string{`\"skipped_by_reason\": {}`, `\"status\": \"adopted\"`} {\n\t\tif !strings.Contains(s, want) {\n\t\t\tt.Errorf(\"JSON missing %s:\\n%s\", want, s)\n\t\t}\n\t}\n\tif !strings.HasSuffix(s, \"\\n\") {\n\t\tt.Errorf(\"JSON does not end with a newline\")\n\t}\n\tif strings.Contains(s, \"timestamp\") {\n\t\tt.Errorf(\"JSON has a timestamp field:\\n%s\", s)\n\t}\n}","file":"internal/adopt/report_test.go","name":"TestReportJSONKeys"}}
func TestReportJSONKeys(t *testing.T) {
	r := NewReport()
	r.Add(FileReport{Path: "a.go", Status: StatusAdopted, Units: 1})
	b, err := r.JSON()
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{`"skipped_by_reason": {}`, `"status": "adopted"`} {
		if !strings.Contains(s, want) {
			t.Errorf("JSON missing %s:\n%s", want, s)
		}
	}
	if !strings.HasSuffix(s, "\n") {
		t.Errorf("JSON does not end with a newline")
	}
	if strings.Contains(s, "timestamp") {
		t.Errorf("JSON has a timestamp field:\n%s", s)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestResidualsOf(t *testing.T) {\n\tsrc := []byte(\"# licence\\nimport os\\n\\ndef f():\\n    pass\\n\\n# note\\nx = 1\\n\")\n\tspans := []Span{\n\t\t{Kind: KindResidual, Text: \"# licence\\n\"},\n\t\t{Kind: KindResidual, Text: \"import os\\n\\n\"},\n\t\t{Kind: KindUnit, Text: \"def f():\\n    pass\\n\", Action: \"setFunction\"},\n\t\t{Kind: KindResidual, Text: \"\\n# note\\n\"},\n\t\t{Kind: KindResidual, Text: \"x = 1\\n\"},\n\t}\n\tgot := ResidualsOf(src, spans, genpkg.Comment{Prefix: \"#\"})\n\twant := []ResidualInfo{\n\t\t{StartLine: 1, EndLine: 1, FirstLine: \"# licence\", Reason: ReasonHeader},\n\t\t{StartLine: 2, EndLine: 2, FirstLine: \"import os\", Reason: ReasonOtherResidual},\n\t\t{StartLine: 7, EndLine: 7, FirstLine: \"# note\", Reason: ReasonFloatingComment},\n\t\t{StartLine: 8, EndLine: 8, FirstLine: \"x = 1\", Reason: ReasonOtherResidual},\n\t}\n\tif len(got) != len(want) {\n\t\tt.Fatalf(\"got %d residuals, want %d: %+v\", len(got), len(want), got)\n\t}\n\tfor i := range want {\n\t\tif got[i] != want[i] {\n\t\t\tt.Errorf(\"residual %d = %+v, want %+v\", i, got[i], want[i])\n\t\t}\n\t}\n}","file":"internal/adopt/report_test.go","name":"TestResidualsOf"}}
func TestResidualsOf(t *testing.T) {
	src := []byte("# licence\nimport os\n\ndef f():\n    pass\n\n# note\nx = 1\n")
	spans := []Span{
		{Kind: KindResidual, Text: "# licence\n"},
		{Kind: KindResidual, Text: "import os\n\n"},
		{Kind: KindUnit, Text: "def f():\n    pass\n", Action: "setFunction"},
		{Kind: KindResidual, Text: "\n# note\n"},
		{Kind: KindResidual, Text: "x = 1\n"},
	}
	got := ResidualsOf(src, spans, genpkg.Comment{Prefix: "#"})
	want := []ResidualInfo{
		{StartLine: 1, EndLine: 1, FirstLine: "# licence", Reason: ReasonHeader},
		{StartLine: 2, EndLine: 2, FirstLine: "import os", Reason: ReasonOtherResidual},
		{StartLine: 7, EndLine: 7, FirstLine: "# note", Reason: ReasonFloatingComment},
		{StartLine: 8, EndLine: 8, FirstLine: "x = 1", Reason: ReasonOtherResidual},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d residuals, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("residual %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestResidualsOfEmptyPrefix(t *testing.T) {\n\tsrc := []byte(\"# licence\\nimport os\\n\\ndef f():\\n    pass\\n\\n# note\\nx = 1\\n\")\n\tspans := []Span{\n\t\t{Kind: KindResidual, Text: \"# licence\\n\"},\n\t\t{Kind: KindResidual, Text: \"import os\\n\\n\"},\n\t\t{Kind: KindUnit, Text: \"def f():\\n    pass\\n\", Action: \"setFunction\"},\n\t\t{Kind: KindResidual, Text: \"\\n# note\\n\"},\n\t\t{Kind: KindResidual, Text: \"x = 1\\n\"},\n\t}\n\tgot := ResidualsOf(src, spans, genpkg.Comment{})\n\tif len(got) != 4 {\n\t\tt.Fatalf(\"got %d residuals, want 4: %+v\", len(got), got)\n\t}\n\tfor i, res := range got {\n\t\tif res.Reason != ReasonOtherResidual {\n\t\t\tt.Errorf(\"residual %d reason = %s, want %s\", i, res.Reason, ReasonOtherResidual)\n\t\t}\n\t}\n}","file":"internal/adopt/report_test.go","name":"TestResidualsOfEmptyPrefix"}}
func TestResidualsOfEmptyPrefix(t *testing.T) {
	src := []byte("# licence\nimport os\n\ndef f():\n    pass\n\n# note\nx = 1\n")
	spans := []Span{
		{Kind: KindResidual, Text: "# licence\n"},
		{Kind: KindResidual, Text: "import os\n\n"},
		{Kind: KindUnit, Text: "def f():\n    pass\n", Action: "setFunction"},
		{Kind: KindResidual, Text: "\n# note\n"},
		{Kind: KindResidual, Text: "x = 1\n"},
	}
	got := ResidualsOf(src, spans, genpkg.Comment{})
	if len(got) != 4 {
		t.Fatalf("got %d residuals, want 4: %+v", len(got), got)
	}
	for i, res := range got {
		if res.Reason != ReasonOtherResidual {
			t.Errorf("residual %d reason = %s, want %s", i, res.Reason, ReasonOtherResidual)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestReportText(t *testing.T) {\n\tr := NewReport()\n\tr.Add(FileReport{\n\t\tPath:      \"a.go\",\n\t\tStatus:    StatusSkipped,\n\t\tReason:    ReasonHeader,\n\t\tResiduals: []ResidualInfo{{StartLine: 1, EndLine: 2, FirstLine: \"# licence\", Reason: ReasonHeader}},\n\t})\n\tr.Recording = \".sedum/adopt/x.json\"\n\twant := \"adopt: 1 file(s): 0 adopted, 1 skipped, 0 rejected, 0 unit(s)\\n\" +\n\t\t\"skipped by reason:\\n\" +\n\t\t\"  header: 1\\n\" +\n\t\t\"header:\\n\" +\n\t\t\"  a.go (lines 1-2: # licence)\\n\" +\n\t\t\"recording: .sedum/adopt/x.json\\n\"\n\tif got := r.Text(); got != want {\n\t\tt.Errorf(\"Text() =\\n%s\\nwant\\n%s\", got, want)\n\t}\n}","file":"internal/adopt/report_test.go","name":"TestReportText"}}
func TestReportText(t *testing.T) {
	r := NewReport()
	r.Add(FileReport{
		Path:      "a.go",
		Status:    StatusSkipped,
		Reason:    ReasonHeader,
		Residuals: []ResidualInfo{{StartLine: 1, EndLine: 2, FirstLine: "# licence", Reason: ReasonHeader}},
	})
	r.Recording = ".sedum/adopt/x.json"
	want := "adopt: 1 file(s): 0 adopted, 1 skipped, 0 rejected, 0 unit(s)\n" +
		"skipped by reason:\n" +
		"  header: 1\n" +
		"header:\n" +
		"  a.go (lines 1-2: # licence)\n" +
		"recording: .sedum/adopt/x.json\n"
	if got := r.Text(); got != want {
		t.Errorf("Text() =\n%s\nwant\n%s", got, want)
	}
}

// /sedum:setFunction
