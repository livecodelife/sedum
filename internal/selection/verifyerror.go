package selection

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/selection/verifyerror.go","path":"fmt"},"kwargs":"sha256:47110f466723c64370fcf52135188cc8cf232ddcdba6e16f8ca26a91d5474de6","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/selection/verifyerror.go","name":"VerifyError"},"kwargs":"sha256:01ac86f5ad8700574e10a711f1c2c47969929cc02363804f84fac9aab0d69eb3","body":"sha256:eeb66350abe3311dd05e9f83599a7a6602ec1b7a9fba55b6d489dd368e367ff4"}
// VerifyError is the error for a record whose selection still failed the
// project's verification after every repair it was allowed. It carries what the
// record's calls cost, so a caller that records cost does not lose it.
type VerifyError struct {
	RecordID         string
	Repairs          int
	Calls            int
	Rejected         int
	Completeness     int
	PromptTokens     int
	CompletionTokens int
	Failure          string
}

// /sedum:setType
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/selection/verifyerror.go","name":"Error","receiver":"VerifyError"},"kwargs":"sha256:a9efc37ebe23e5f45447e96618b3ac9f1d18d45cb53efe172e6ac4244d246fca","body":"sha256:f867329e4c43e5a05ffe88ae291163a1073bff1591cb2faff5dea39d9f55871e"}
// Error describes the failure with the repairs, calls and rejections spent.
func (e *VerifyError) Error() string {
	return fmt.Sprintf("record %s: the selection did not verify after %d repair attempt(s) (%d model call(s), %d answer(s) rejected by validation):\n%s",
		e.RecordID, e.Repairs, e.Calls, e.Rejected, e.Failure)
}

// /sedum:setMethod
