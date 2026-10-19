package selection

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/selection/verifyerror.go","path":"fmt"}}
import "fmt"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// VerifyError is the error for a record whose selection still failed the\n// project's verification after every repair it was allowed. It carries what the\n// record's calls cost, so a caller that records cost does not lose it.\ntype VerifyError struct {\n\tRecordID         string\n\tRepairs          int\n\tCalls            int\n\tRejected         int\n\tCompleteness     int\n\tPromptTokens     int\n\tCompletionTokens int\n\tFailure          string\n}","file":"internal/selection/verifyerror.go","name":"VerifyError"}}
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
// sedum:setMethod {"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Error describes the failure with the repairs, calls and rejections spent.\nfunc (e *VerifyError) Error() string {\n\treturn fmt.Sprintf(\"record %s: the selection did not verify after %d repair attempt(s) (%d model call(s), %d answer(s) rejected by validation):\\n%s\",\n\t\te.RecordID, e.Repairs, e.Calls, e.Rejected, e.Failure)\n}","file":"internal/selection/verifyerror.go","name":"Error","receiver":"VerifyError"}}
// Error describes the failure with the repairs, calls and rejections spent.
func (e *VerifyError) Error() string {
	return fmt.Sprintf("record %s: the selection did not verify after %d repair attempt(s) (%d model call(s), %d answer(s) rejected by validation):\n%s",
		e.RecordID, e.Repairs, e.Calls, e.Rejected, e.Failure)
}

// /sedum:setMethod
