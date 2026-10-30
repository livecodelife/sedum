package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/report_test.go","path":"strings"},"kwargs":"sha256:d373cfd1037d92def357d45dd0306b74b82777f365eb4d385cbbd3b25278229f","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/report_test.go","path":"testing"},"kwargs":"sha256:88eef00d50653115dae4aa4da4e6294f154dd66143942b9a63d93e4596b5ee5b","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/report_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:0459380aa09d96b25f57f3be46df210e5b81509f090b72befd19bede0d142492","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/report_test.go","name":"TestReportTotals"},"kwargs":"sha256:13cb89689ab161599c860bcfb84c74d15f377a71a00ee73715648e017638b062","body":"sha256:1f5b0649a9c1f2827635d31c798a608e00937044ef1866ce1c61c2525d8a7fc0"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/report_test.go","name":"TestReportJSONKeys"},"kwargs":"sha256:d4107a3b8550a675c9da226ba68a07c4aa70509321684b8e91f9b2a54fde9500","body":"sha256:8ebe122b9df28a3becb54bef70c2d5ea3df4bf6adf49ecc97473faf2a918a6b3"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/report_test.go","name":"TestResidualsOf"},"kwargs":"sha256:9ba712378dae8b54642dda647e8c115df0c61e7ec92ff90e11ab3961754820d6","body":"sha256:a66507220542ea76448740e98449a932718393f154d1f4417ef2ab4cc65ff5bf"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/report_test.go","name":"TestResidualsOfEmptyPrefix"},"kwargs":"sha256:14c795f7507db9b0f0031650af9e26480fff741ebeed9a4df34e885c1997d8e9","body":"sha256:127088ff5b9adc53199644bb22c5926c9978940026414b784ed566b8080927cf"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/report_test.go","name":"TestReportText"},"kwargs":"sha256:c4892ac84bf207069fbd4698d649a2a46922c6b77014c74371dfb44393af3688","body":"sha256:f005ebb8260b9bd877c643c631cc08d76f08f0758225bf847e6e187f0e937c13"}
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
