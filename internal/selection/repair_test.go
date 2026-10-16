package selection

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/selection/repair_test.go","path":"context"}}
import "context"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/selection/repair_test.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/selection/repair_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/selection/repair_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// failing is a Verify that fails its first n calls with message and then\n// passes, remembering the candidates it was given.\ntype failing struct {\n\tn          int\n\tmessage    string\n\tcandidates []map[string]string\n}","file":"internal/selection/repair_test.go","name":"failing"}}
// failing is a Verify that fails its first n calls with message and then
// passes, remembering the candidates it was given.
type failing struct {
	n          int
	message    string
	candidates []map[string]string
}

// /sedum:setType
// sedum:setMethod {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"func (f *failing) verify(_ context.Context, candidate map[string]string) (string, error) {\n\tf.candidates = append(f.candidates, candidate)\n\tif len(f.candidates) <= f.n {\n\t\treturn f.message, nil\n\t}\n\treturn \"\", nil\n}","file":"internal/selection/repair_test.go","name":"verify","receiver":"failing"}}
func (f *failing) verify(_ context.Context, candidate map[string]string) (string, error) {
	f.candidates = append(f.candidates, candidate)
	if len(f.candidates) <= f.n {
		return f.message, nil
	}
	return "", nil
}

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func selectRepairing(t *testing.T, repair int, verify func(context.Context, map[string]string) (string, error), responses ...string) (Answer, *stub, error) {\n\tt.Helper()\n\tclient := &stub{responses: responses}\n\treq := railsRequest(t)\n\treq.Files[0].Rendered = \"class UsersController\\n  # sedum:anchor:class_body\\nend\\n\"\n\tanswer, err := Select(context.Background(), client, req, Options{\n\t\tRetries: 3,\n\t\tRepair:  repair,\n\t\tVerify:  verify,\n\t\tOutput:  t.TempDir(),\n\t})\n\treturn answer, client, err\n}","file":"internal/selection/repair_test.go","name":"selectRepairing"}}
func selectRepairing(t *testing.T, repair int, verify func(context.Context, map[string]string) (string, error), responses ...string) (Answer, *stub, error) {
	t.Helper()
	client := &stub{responses: responses}
	req := railsRequest(t)
	req.Files[0].Rendered = "class UsersController\n  # sedum:anchor:class_body\nend\n"
	answer, err := Select(context.Background(), client, req, Options{
		Retries: 3,
		Repair:  repair,
		Verify:  verify,
		Output:  t.TempDir(),
	})
	return answer, client, err
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRepairSendsAFailedVerificationBack(t *testing.T) {\n\tcheck := &failing{n: 1, message: \"These checks failed: go build ./...\"}\n\tanswer, client, err := selectRepairing(t, 1, check.verify, validResponse, validResponse)\n\tif err != nil {\n\t\tt.Fatalf(\"the repair should have succeeded: %v\", err)\n\t}\n\tif answer.Repairs != 1 || answer.Calls != 2 || answer.Rejected != 0 {\n\t\tt.Errorf(\"Repairs = %d, Calls = %d, Rejected = %d, want 1, 2 and 0\", answer.Repairs, answer.Calls, answer.Rejected)\n\t}\n\tsecond := client.seen[1]\n\tlast := second[len(second)-1]\n\tif last.Role != RoleUser || last.Content != check.message {\n\t\tt.Errorf(\"the failure was not put to the model as the next message: %+v\", last)\n\t}\n\tif prev := second[len(second)-2]; prev.Role != RoleAssistant || prev.Content != validResponse {\n\t\tt.Errorf(\"the answer that failed was not kept in the conversation: %+v\", prev)\n\t}\n}","file":"internal/selection/repair_test.go","name":"TestRepairSendsAFailedVerificationBack"}}
func TestRepairSendsAFailedVerificationBack(t *testing.T) {
	check := &failing{n: 1, message: "These checks failed: go build ./..."}
	answer, client, err := selectRepairing(t, 1, check.verify, validResponse, validResponse)
	if err != nil {
		t.Fatalf("the repair should have succeeded: %v", err)
	}
	if answer.Repairs != 1 || answer.Calls != 2 || answer.Rejected != 0 {
		t.Errorf("Repairs = %d, Calls = %d, Rejected = %d, want 1, 2 and 0", answer.Repairs, answer.Calls, answer.Rejected)
	}
	second := client.seen[1]
	last := second[len(second)-1]
	if last.Role != RoleUser || last.Content != check.message {
		t.Errorf("the failure was not put to the model as the next message: %+v", last)
	}
	if prev := second[len(second)-2]; prev.Role != RoleAssistant || prev.Content != validResponse {
		t.Errorf("the answer that failed was not kept in the conversation: %+v", prev)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRepairVerifiesTheFilesTheAnswerWouldWrite(t *testing.T) {\n\tcheck := &failing{}\n\tif _, _, err := selectRepairing(t, 1, check.verify, validResponse); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(check.candidates) != 1 {\n\t\tt.Fatalf(\"Verify was called %d times, want 1\", len(check.candidates))\n\t}\n\tcontent, ok := check.candidates[0][\"app/controllers/users_controller.rb\"]\n\tif !ok || !strings.Contains(content, \"index\") {\n\t\tt.Errorf(\"the candidate lacks the controller the answer writes: %v\", check.candidates[0])\n\t}\n}","file":"internal/selection/repair_test.go","name":"TestRepairVerifiesTheFilesTheAnswerWouldWrite"}}
func TestRepairVerifiesTheFilesTheAnswerWouldWrite(t *testing.T) {
	check := &failing{}
	if _, _, err := selectRepairing(t, 1, check.verify, validResponse); err != nil {
		t.Fatal(err)
	}
	if len(check.candidates) != 1 {
		t.Fatalf("Verify was called %d times, want 1", len(check.candidates))
	}
	content, ok := check.candidates[0]["app/controllers/users_controller.rb"]
	if !ok || !strings.Contains(content, "index") {
		t.Errorf("the candidate lacks the controller the answer writes: %v", check.candidates[0])
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRepairExhaustionFailsNamingTheRecordAndTheFailure(t *testing.T) {\n\tcheck := &failing{n: 5, message: \"These checks failed: make test\"}\n\t_, client, err := selectRepairing(t, 1, check.verify, validResponse, validResponse)\n\tif err == nil {\n\t\tt.Fatal(\"an answer that never verifies should fail\")\n\t}\n\tfor _, want := range []string{\"PR-014\", \"did not verify after 1 repair attempt(s)\", check.message} {\n\t\tif !strings.Contains(err.Error(), want) {\n\t\t\tt.Errorf(\"the error lacks %q: %v\", want, err)\n\t\t}\n\t}\n\tif len(client.seen) != 2 {\n\t\tt.Errorf(\"the model was called %d times, want the first answer and one repair\", len(client.seen))\n\t}\n}","file":"internal/selection/repair_test.go","name":"TestRepairExhaustionFailsNamingTheRecordAndTheFailure"}}
func TestRepairExhaustionFailsNamingTheRecordAndTheFailure(t *testing.T) {
	check := &failing{n: 5, message: "These checks failed: make test"}
	_, client, err := selectRepairing(t, 1, check.verify, validResponse, validResponse)
	if err == nil {
		t.Fatal("an answer that never verifies should fail")
	}
	for _, want := range []string{"PR-014", "did not verify after 1 repair attempt(s)", check.message} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error lacks %q: %v", want, err)
		}
	}
	if len(client.seen) != 2 {
		t.Errorf("the model was called %d times, want the first answer and one repair", len(client.seen))
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestWithoutVerifyNothingIsRepaired(t *testing.T) {\n\tanswer, client, err := selectRepairing(t, 3, nil, validResponse)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif answer.Repairs != 0 || len(client.seen) != 1 {\n\t\tt.Errorf(\"Repairs = %d after %d call(s), want 0 after 1\", answer.Repairs, len(client.seen))\n\t}\n}","file":"internal/selection/repair_test.go","name":"TestWithoutVerifyNothingIsRepaired"}}
func TestWithoutVerifyNothingIsRepaired(t *testing.T) {
	answer, client, err := selectRepairing(t, 3, nil, validResponse)
	if err != nil {
		t.Fatal(err)
	}
	if answer.Repairs != 0 || len(client.seen) != 1 {
		t.Errorf("Repairs = %d after %d call(s), want 0 after 1", answer.Repairs, len(client.seen))
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAVerifyErrorStopsTheRecordWithoutARepair(t *testing.T) {\n\tbroken := func(context.Context, map[string]string) (string, error) { return \"\", errors.New(\"no shell\") }\n\t_, client, err := selectRepairing(t, 3, broken, validResponse, validResponse)\n\tif err == nil || !strings.Contains(err.Error(), \"verifying the selection\") || !strings.Contains(err.Error(), \"no shell\") {\n\t\tt.Fatalf(\"error = %v, want one naming the verification and its cause\", err)\n\t}\n\tif len(client.seen) != 1 {\n\t\tt.Errorf(\"the model was called %d times, want 1: a broken verifier is not the model's mistake\", len(client.seen))\n\t}\n}","file":"internal/selection/repair_test.go","name":"TestAVerifyErrorStopsTheRecordWithoutARepair"}}
func TestAVerifyErrorStopsTheRecordWithoutARepair(t *testing.T) {
	broken := func(context.Context, map[string]string) (string, error) { return "", errors.New("no shell") }
	_, client, err := selectRepairing(t, 3, broken, validResponse, validResponse)
	if err == nil || !strings.Contains(err.Error(), "verifying the selection") || !strings.Contains(err.Error(), "no shell") {
		t.Fatalf("error = %v, want one naming the verification and its cause", err)
	}
	if len(client.seen) != 1 {
		t.Errorf("the model was called %d times, want 1: a broken verifier is not the model's mistake", len(client.seen))
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRepairExhaustionCountsEveryModelCall(t *testing.T) {\n\tbad := `{\"invocations\":[{\"action\":\"createControllerMethod\",\"kwargs\":{\"controller\":\"users\"}}]}`\n\tcheck := &failing{n: 5, message: \"These checks failed: make test\"}\n\t_, client, err := selectRepairing(t, 1, check.verify, validResponse, bad, validResponse)\n\tif err == nil {\n\t\tt.Fatal(\"an answer that never verifies should fail\")\n\t}\n\tfor _, want := range []string{\"1 repair attempt(s)\", \"3 model call(s)\", \"1 answer(s) rejected by validation\"} {\n\t\tif !strings.Contains(err.Error(), want) {\n\t\t\tt.Errorf(\"the error lacks %q: %v\", want, err)\n\t\t}\n\t}\n\tif len(client.seen) != 3 {\n\t\tt.Errorf(\"the model was called %d times, want 3\", len(client.seen))\n\t}\n}","file":"internal/selection/repair_test.go","name":"TestRepairExhaustionCountsEveryModelCall"}}
func TestRepairExhaustionCountsEveryModelCall(t *testing.T) {
	bad := `{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users"}}]}`
	check := &failing{n: 5, message: "These checks failed: make test"}
	_, client, err := selectRepairing(t, 1, check.verify, validResponse, bad, validResponse)
	if err == nil {
		t.Fatal("an answer that never verifies should fail")
	}
	for _, want := range []string{"1 repair attempt(s)", "3 model call(s)", "1 answer(s) rejected by validation"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error lacks %q: %v", want, err)
		}
	}
	if len(client.seen) != 3 {
		t.Errorf("the model was called %d times, want 3", len(client.seen))
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRepairZeroVerifiesOnceAndFails(t *testing.T) {\n\tcheck := &failing{n: 5, message: \"These checks failed: make test\"}\n\t_, client, err := selectRepairing(t, 0, check.verify, validResponse, validResponse)\n\tif err == nil || !strings.Contains(err.Error(), \"after 0 repair attempt(s)\") {\n\t\tt.Fatalf(\"error = %v, want one saying no repair was allowed\", err)\n\t}\n\tif len(client.seen) != 1 || len(check.candidates) != 1 {\n\t\tt.Errorf(\"%d model call(s) and %d verification(s), want 1 and 1\", len(client.seen), len(check.candidates))\n\t}\n}","file":"internal/selection/repair_test.go","name":"TestRepairZeroVerifiesOnceAndFails"}}
func TestRepairZeroVerifiesOnceAndFails(t *testing.T) {
	check := &failing{n: 5, message: "These checks failed: make test"}
	_, client, err := selectRepairing(t, 0, check.verify, validResponse, validResponse)
	if err == nil || !strings.Contains(err.Error(), "after 0 repair attempt(s)") {
		t.Fatalf("error = %v, want one saying no repair was allowed", err)
	}
	if len(client.seen) != 1 || len(check.candidates) != 1 {
		t.Errorf("%d model call(s) and %d verification(s), want 1 and 1", len(client.seen), len(check.candidates))
	}
}

// /sedum:setFunction
