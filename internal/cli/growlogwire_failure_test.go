package cli

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogwire_failure_test.go","path":"fmt"},"kwargs":"sha256:93cd10c5e498407113a35cf1a350f5dae38c7e330a4c890a47118a4e457011e3","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogwire_failure_test.go","path":"testing"},"kwargs":"sha256:e4e52e48afec29bffc9a853c335be834cb1af793a8ce6ca0038d350b87accf4a","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogwire_failure_test.go","path":"time"},"kwargs":"sha256:cdb4ecbe18618664ed36e60731337340f4be0d51f0d2f5099fc73c77cec3b209","body":"sha256:e99128ebf4e36998f40b33adf2eef57f4d6c315f1a4d1bdd3b7fc09a0d7bc82d"}
import "time"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogwire_failure_test.go","path":"github.com/livecodelife/sedum/internal/selection"},"kwargs":"sha256:dea9a0ad28f27442340100daacf1a19b9ecac592c8585fc491f2ddf8f15636c7","body":"sha256:2f130b07ae47e7f9e9e5926057408fb954954176e88beaf96984f1a3b7ef3316"}
import "github.com/livecodelife/sedum/internal/selection"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-6edbd718","gen":"sha256:fe5f8ef6698e14a5acdfb5f1e3db7faa0a869ea16b3292dc2c11d3ed3591000e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/growlogwire_failure_test.go","name":"TestGrowLogFailureKeepsTheCostOfAFailedVerification"},"kwargs":"sha256:7d11a351372b9596625ebd7c85d33d64fff04574a119471353ed80c2925eac31","body":"sha256:23ab4ba7fc8943145723e028bd212ce830231ede7948cad347a8852ed7613ca0"}
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
