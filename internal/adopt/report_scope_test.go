package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/report_scope_test.go","path":"testing"},"kwargs":"sha256:7fb1ef572dc134ba2e12e17b777fac243c68c2f44998d5ccef7e53183837252c","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9ac6951e","gen":"sha256:f3bd9f66cf02c71043301cb6b6a6dc79e6a4fc1e4c33cf181069eaa26197a65b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/report_scope_test.go","name":"TestReportTextScopeLines"},"kwargs":"sha256:fc1b2868ec39e2e03b58145c7f9f1e186d8afcf3b1834d677bcd4c7533d02ce5","body":"sha256:2a393f4182f3cff7be3bbfe1a31db1814889fe7b3c03c38f712aa1531a7cbe81"}
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
