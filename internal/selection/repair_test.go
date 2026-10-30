package selection

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/selection/repair_test.go","path":"context"},"kwargs":"sha256:a6342cab2f88e3555eb3ac126a6b451102accec65ef306d0c42264c7df427f07","body":"sha256:0f7443fb25d7ad6a970617d4d0df6123283e24801b94c1d84fbe431230f31e1e"}
import "context"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/selection/repair_test.go","path":"errors"},"kwargs":"sha256:eb2963e17ec752c5b8d3c494d38a81160129f6e564d4b444ae98769c1f6cb31c","body":"sha256:2b718237ba0a56467c8fd494cc131e942ef410e6228c03b644d70476c71862dc"}
import "errors"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/selection/repair_test.go","path":"strings"},"kwargs":"sha256:4270394a91dcea7f877051cb504b7f8d99a0f6ca5d90e3bdf42deb018a6a6d3c","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/selection/repair_test.go","path":"testing"},"kwargs":"sha256:e919e837ba28a5d2835c0b7f3a226d7a393f5a041f3e5e77bdbadfcf01794727","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/selection/repair_test.go","name":"failing"},"kwargs":"sha256:76c8e9c32c5b25f57dbc4a9d7aa0931742acaf929b397f69e1b74c82969e015b","body":"sha256:555e625612eb49dabaff112b072a2b5b135286d8bc1643b02af82a37659ac18e"}
// failing is a Verify that fails its first n calls with message and then
// passes, remembering the candidates it was given.
type failing struct {
	n          int
	message    string
	candidates []map[string]string
}

// /sedum:setType
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/selection/repair_test.go","name":"verify","receiver":"failing"},"kwargs":"sha256:0ae3fa3aea59ab3858f011144c40e28196d19c7f18593f3e49cd0c460a0481a9","body":"sha256:68491d4bce9072733f7196b0b1026631e94b8def5f0dc1344b14258217be5752"}
func (f *failing) verify(_ context.Context, candidate map[string]string) (string, error) {
	f.candidates = append(f.candidates, candidate)
	if len(f.candidates) <= f.n {
		return f.message, nil
	}
	return "", nil
}

// /sedum:setMethod
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/selection/repair_test.go","name":"selectRepairing"},"kwargs":"sha256:4f1c6811e578a5cafeff5e8d1655e2e6865489f1d80cb3f9d40cf0b1f2d6d536","body":"sha256:3bc7c6a8696b262c680ca8c0ec15a6dfa8bdc71817fccace81260a2592f05b6a"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/selection/repair_test.go","name":"TestRepairSendsAFailedVerificationBack"},"kwargs":"sha256:bcc84c5c480f59e43b262c26e2a22aae93e70dd7a953f2fce61e22d678f5f886","body":"sha256:3b98543e68fd4789f2bba99494acb42d7bf0ccea7c5e2d47b0657c305e121a36"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/selection/repair_test.go","name":"TestRepairVerifiesTheFilesTheAnswerWouldWrite"},"kwargs":"sha256:f1f0d7a22d94e1320f0c4131051f9cd9f49c37f55fb286d83f75981d919887a9","body":"sha256:589921a6dfde2dfd6869f0a59b7c74a19678191240746ff7d7b0ce653479520c"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/selection/repair_test.go","name":"TestRepairExhaustionFailsNamingTheRecordAndTheFailure"},"kwargs":"sha256:9da6e97f8e3aecd4df0f11ba3bccea6051d3c2d0eec0e523eaf5d71264dd8ea8","body":"sha256:d9d49544369fa3f501c04bb0a7f651f6627a381b7c0da45f4187a4da1853f4a0"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/selection/repair_test.go","name":"TestWithoutVerifyNothingIsRepaired"},"kwargs":"sha256:80d2079b8c32d33315435ba5e741b742577252b8b3a2941f8d671f1bebd5e2fd","body":"sha256:cabe5b6b947108659ce8de5e7499d48906b2da0652ae243d3b5f7e2ec5122642"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/selection/repair_test.go","name":"TestAVerifyErrorStopsTheRecordWithoutARepair"},"kwargs":"sha256:c64425280aebfa296d68747290c3b2e06b5f0aab9691941261cc47aa5a0334f2","body":"sha256:4cbda895cbcc3de0cfcb8b4014041c7d83d4ee2f0691c49cc054914c34795bbd"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/selection/repair_test.go","name":"TestRepairExhaustionCountsEveryModelCall"},"kwargs":"sha256:756dbb051b9f5d07628e5a917e23142461b716d2d61ea38873e7395e65616152","body":"sha256:7043516c47ef1b52372812443657e20fa713bc7c40f29d5b8f476c8a7134e57e"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-751790c4","gen":"sha256:5f567aeafced45e116fea6b5b9f826a81121f44699d5a230e602659155ed15c4","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/selection/repair_test.go","name":"TestRepairZeroVerifiesOnceAndFails"},"kwargs":"sha256:44d6044271816c65a87e9104ba96bc45fab2052362773154d000824800b69e98","body":"sha256:b2f0d1af7d53126004d54152f6da40c9bf1767ee37b38a401b39f88c38b5c4fb"}
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
