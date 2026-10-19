package selection

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/selection/verifyerror_test.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/selection/verifyerror_test.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/selection/verifyerror_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/selection/verifyerror_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestVerifyErrorMessage(t *testing.T) {\n\terr := &VerifyError{RecordID: \"r1\", Repairs: 2, Calls: 3, Rejected: 1, Failure: \"tsc failed\"}\n\twant := \"record r1: the selection did not verify after 2 repair attempt(s) (3 model call(s), 1 answer(s) rejected by validation):\\ntsc failed\"\n\tif got := err.Error(); got != want {\n\t\tt.Fatalf(\"Error() = %q, want %q\", got, want)\n\t}\n}","file":"internal/selection/verifyerror_test.go","name":"TestVerifyErrorMessage"}}
func TestVerifyErrorMessage(t *testing.T) {
	err := &VerifyError{RecordID: "r1", Repairs: 2, Calls: 3, Rejected: 1, Failure: "tsc failed"}
	want := "record r1: the selection did not verify after 2 repair attempt(s) (3 model call(s), 1 answer(s) rejected by validation):\ntsc failed"
	if got := err.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestVerifyErrorIsFoundThroughWrapping(t *testing.T) {\n\tvar err error = fmt.Errorf(\"grow: %w\", &VerifyError{RecordID: \"r2\", PromptTokens: 10, CompletionTokens: 4})\n\tvar ve *VerifyError\n\tif !errors.As(err, &ve) || ve.RecordID != \"r2\" || ve.PromptTokens != 10 || ve.CompletionTokens != 4 {\n\t\tt.Fatalf(\"errors.As found %+v\", ve)\n\t}\n\tif !strings.Contains(err.Error(), \"r2\") {\n\t\tt.Fatalf(\"message lost the record: %v\", err)\n\t}\n}","file":"internal/selection/verifyerror_test.go","name":"TestVerifyErrorIsFoundThroughWrapping"}}
func TestVerifyErrorIsFoundThroughWrapping(t *testing.T) {
	var err error = fmt.Errorf("grow: %w", &VerifyError{RecordID: "r2", PromptTokens: 10, CompletionTokens: 4})
	var ve *VerifyError
	if !errors.As(err, &ve) || ve.RecordID != "r2" || ve.PromptTokens != 10 || ve.CompletionTokens != 4 {
		t.Fatalf("errors.As found %+v", ve)
	}
	if !strings.Contains(err.Error(), "r2") {
		t.Fatalf("message lost the record: %v", err)
	}
}

// /sedum:setFunction
