package selection

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/selection/verifyerror_test.go","path":"errors"},"kwargs":"sha256:884841ac8cf5e46d8b73743ff259e4fe2ec1e8766fa281bcf43124212dde588b","body":"sha256:2b718237ba0a56467c8fd494cc131e942ef410e6228c03b644d70476c71862dc"}
import "errors"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/selection/verifyerror_test.go","path":"fmt"},"kwargs":"sha256:65b98de730975314a131991c828115451cae50005c8b2c53265f15bb06adda84","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/selection/verifyerror_test.go","path":"strings"},"kwargs":"sha256:79505515cd9d4bbaced556022fcbd07c6baade9468ea149897cc75ddd6cc6a66","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/selection/verifyerror_test.go","path":"testing"},"kwargs":"sha256:fc4db868c0429051a04887494eda12db54ff2b9e13ae8e432e5f3a9d20f40e2e","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/selection/verifyerror_test.go","name":"TestVerifyErrorMessage"},"kwargs":"sha256:bd0f590ee7def877ee2508d26000f7b38ced80e90c50ef8f640b945fb73c256a","body":"sha256:df63c6c6ae10c03e1fbfb7834f211da6e1a280b9711b75d627bfdcfbf0b3eb46"}
func TestVerifyErrorMessage(t *testing.T) {
	err := &VerifyError{RecordID: "r1", Repairs: 2, Calls: 3, Rejected: 1, Failure: "tsc failed"}
	want := "record r1: the selection did not verify after 2 repair attempt(s) (3 model call(s), 1 answer(s) rejected by validation):\ntsc failed"
	if got := err.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-d3dc76b9","gen":"sha256:18e240ef47acdd31aca64499659c9c069bbac146a60dcf47ea390e4dfa5b43e1","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/selection/verifyerror_test.go","name":"TestVerifyErrorIsFoundThroughWrapping"},"kwargs":"sha256:b9b43cb120de14cb8528506ee5043022f750bb75c20d5b4ab728b2f94c639d88","body":"sha256:96d3b07b8b5998943cab2afef60eb617fd85c834ac97d8e9083cb4d301eca656"}
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
