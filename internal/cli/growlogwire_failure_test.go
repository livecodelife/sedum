package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire_failure_test.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire_failure_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire_failure_test.go","path":"time"}}
import "time"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogwire_failure_test.go","path":"github.com/livecodelife/sedum/internal/selection"}}
import "github.com/livecodelife/sedum/internal/selection"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGrowLogFailureKeepsTheCostOfAFailedVerification(t *testing.T) {\n\tver := &selection.VerifyError{RecordID: \"r7\", Repairs: 3, Calls: 4, Rejected: 1, Completeness: 2, PromptTokens: 900, CompletionTokens: 120, Failure: \"tsc\"}\n\te := growLogFailure(fmt.Errorf(\"grow: %w\", ver), \"m\", time.Now(), time.Second)\n\tif e.Outcome != \"failed\" || e.Record != \"r7\" || e.Model != \"m\" || e.Calls != 4 || e.Rejected != 1 || e.Completeness != 2 || e.Repairs != 3 {\n\t\tt.Fatalf(\"entry = %+v\", e)\n\t}\n\tif e.PromptTokens != 900 || e.CompletionTokens != 120 {\n\t\tt.Fatalf(\"tokens = %d/%d\", e.PromptTokens, e.CompletionTokens)\n\t}\n}","file":"internal/cli/growlogwire_failure_test.go","name":"TestGrowLogFailureKeepsTheCostOfAFailedVerification"}}
func TestGrowLogFailureKeepsTheCostOfAFailedVerification(t *testing.T) {
	ver := &selection.VerifyError{RecordID: "r7", Repairs: 3, Calls: 4, Rejected: 1, Completeness: 2, PromptTokens: 900, CompletionTokens: 120, Failure: "tsc"}
	e := growLogFailure(fmt.Errorf("grow: %w", ver), "m", time.Now(), time.Second)
	if e.Outcome != "failed" || e.Record != "r7" || e.Model != "m" || e.Calls != 4 || e.Rejected != 1 || e.Completeness != 2 || e.Repairs != 3 {
		t.Fatalf("entry = %+v", e)
	}
	if e.PromptTokens != 900 || e.CompletionTokens != 120 {
		t.Fatalf("tokens = %d/%d", e.PromptTokens, e.CompletionTokens)
	}
}

// /sedum:setFunction
