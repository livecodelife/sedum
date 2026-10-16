package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/report_scope_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestReportTextScopeLines(t *testing.T) {\n\tr := NewReport()\n\tr.Add(FileReport{Path: \"a.go\", Record: \"imp-1\", Status: StatusAdopted, Units: 1})\n\tr.Add(FileReport{Path: \"b.go\", Record: \"imp-1\", Status: StatusAdopted, Units: 2})\n\tr.Add(FileReport{Path: \"c.go\", Record: \"imp-2\", Status: StatusAdopted, Units: 1})\n\tr.Add(FileReport{Path: \"d.go\", Record: \"imp-2\", Status: StatusSkipped, Reason: ReasonNoUnits})\n\tr.Recording = \"rec.json\"\n\twant := \"adopt: 4 file(s): 3 adopted, 1 skipped, 0 rejected, 4 unit(s)\\n\" +\n\t\t\"skipped by reason:\\n\" +\n\t\t\"  no_units: 1\\n\" +\n\t\t\"no_units:\\n\" +\n\t\t\"  d.go\\n\" +\n\t\t\"list the adopted files in their imprints so grow --check sees them:\\n\" +\n\t\t\"  provenance edit imp-1 --add-scope a.go --add-scope b.go\\n\" +\n\t\t\"  provenance edit imp-2 --add-scope c.go\\n\" +\n\t\t\"recording: rec.json\\n\"\n\tif got := r.Text(); got != want {\n\t\tt.Errorf(\"Text() =\\n%s\\nwant\\n%s\", got, want)\n\t}\n}","file":"internal/adopt/report_scope_test.go","name":"TestReportTextScopeLines"}}
func TestReportTextScopeLines(t *testing.T) {
	r := NewReport()
	r.Add(FileReport{Path: "a.go", Record: "imp-1", Status: StatusAdopted, Units: 1})
	r.Add(FileReport{Path: "b.go", Record: "imp-1", Status: StatusAdopted, Units: 2})
	r.Add(FileReport{Path: "c.go", Record: "imp-2", Status: StatusAdopted, Units: 1})
	r.Add(FileReport{Path: "d.go", Record: "imp-2", Status: StatusSkipped, Reason: ReasonNoUnits})
	r.Recording = "rec.json"
	want := "adopt: 4 file(s): 3 adopted, 1 skipped, 0 rejected, 4 unit(s)\n" +
		"skipped by reason:\n" +
		"  no_units: 1\n" +
		"no_units:\n" +
		"  d.go\n" +
		"list the adopted files in their imprints so grow --check sees them:\n" +
		"  provenance edit imp-1 --add-scope a.go --add-scope b.go\n" +
		"  provenance edit imp-2 --add-scope c.go\n" +
		"recording: rec.json\n"
	if got := r.Text(); got != want {
		t.Errorf("Text() =\n%s\nwant\n%s", got, want)
	}
}

// /sedum:setFunction
