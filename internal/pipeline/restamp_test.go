package pipeline

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/restamp_test.go","path":"path/filepath"},"kwargs":"sha256:b81af9344703ad90d922cf7b1f1a966dcb324d2f9427e102e8c448c64a205299","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/restamp_test.go","path":"strings"},"kwargs":"sha256:249f86d2b6349c72c9693df0c1ef1d47b507ebc6e7abc911697b0b9734981b8d","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/restamp_test.go","path":"testing"},"kwargs":"sha256:1ee6d2941ef5e4298fe727e6c17bf95b02a7f2d5e64ee5082fc9766bac050e01","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/restamp_test.go","path":"github.com/livecodelife/sedum/internal/inject"},"kwargs":"sha256:7aace62b3f16fe63971219d086034beb79da78d70c4dafd1d03c76a5c1eb9772","body":"sha256:c32cfeba25229dcfd3a66390c7c5b1c1c3cb162181c8360a4df0ff1bdfad701f"}
import "github.com/livecodelife/sedum/internal/inject"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp_test.go","name":"editedRecords"},"kwargs":"sha256:44b73121a12e15d4c32e0d0c44e3bf21e5de7327ee8c3ab60b04f40003013309","body":"sha256:1788cd4293cdb847ac5d7b8915a22f7e2d89d127b8af9c5330311049ed471dbb"}
// editedRecords is the record the recordings below were made from, edited
// since: same id and scope, a new intent, so a new gen.
func editedRecords(t *testing.T) string {
	t.Helper()
	return writeTree(t, map[string]string{
		"prov-2026-aaaaaaaa.yml": "id: prov-2026-aaaaaaaa\nintent: |\n  Something else.\naffected_scope:\n  - app/models/user.rb\n",
	})
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp_test.go","name":"userFile"},"kwargs":"sha256:039ce93cdecabee1bbd275306316bb6a9e5b4a811bb364c8f24e7b892c0a119b","body":"sha256:abd072ba290b4345edc87167fc0640eb534267d5e066ecab2cccb04dc3b605cc"}
// userFile is the one file every recording below writes.
func userFile(t *testing.T, output string) string {
	t.Helper()
	return readFile(t, filepath.Join(output, "app", "models", "user.rb"))
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp_test.go","name":"TestRestampGenTakesTheCurrentGen"},"kwargs":"sha256:23157320348fa61d09ed3422e0b4a9ea3087ab0a404b157b968eb16644241b0a","body":"sha256:8870fff862d6d2c030775d9f65d98602e0d74680cadf8b2bb8ef2eee409c4166"}
func TestRestampGenTakesTheCurrentGen(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, grown := recorded(t, cfg)
	recordedGen := markerGen(t, grown)

	rcfg := replayConfig(t, cfg)
	rcfg.Records = editedRecords(t)
	rcfg.RestampGen = true
	result, err := Replay(rec, rcfg)
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	current := result.Records.Records[0].Gen
	if current == recordedGen {
		t.Fatal("the edited record hashes the same as the original")
	}
	if got := markerGen(t, rcfg.Output); got != current {
		t.Errorf("gen = %q, want the current record's %q", got, current)
	}
	if len(result.Restamps) != 1 || result.Restamps[0].Old != recordedGen || result.Restamps[0].New != current {
		t.Errorf("Restamps = %+v, want one from %q to %q", result.Restamps, recordedGen, current)
	}

	plain := replayConfig(t, cfg)
	plain.Records = rcfg.Records
	if _, err := Replay(rec, plain); err != nil {
		t.Fatalf("Replay without the flag: %v", err)
	}
	if got := markerGen(t, plain.Output); got != recordedGen {
		t.Errorf("without --restamp-gen the gen = %q, want the recorded %q", got, recordedGen)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp_test.go","name":"TestRestampGenNeedsRecords"},"kwargs":"sha256:6dc14d4b867e76019f547e8efdd1e204d08e5881f6132fc02637fde02a55f905","body":"sha256:99cd164a362cf50ba521236f274925449c02831ddc070460bc9a4ff4257fcada"}
func TestRestampGenNeedsRecords(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)

	rcfg := replayConfig(t, cfg)
	rcfg.RestampGen = true
	if _, err := Replay(rec, rcfg); err == nil || !strings.Contains(err.Error(), "--records") {
		t.Fatalf("Replay error = %v, want one naming --records", err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp_test.go","name":"TestRestampGenSkipsARecordThatIsNotCurrent"},"kwargs":"sha256:cab69bca06ad92ccf126644c851dd0a19d96e9020a3732ce9242962ed6d7e0a3","body":"sha256:96a04c3fd771f27d7ccef5df992ec1e0d06bc93376b1acdf36338dc20acd57ad"}
func TestRestampGenSkipsARecordThatIsNotCurrent(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, grown := recorded(t, cfg)

	rcfg := replayConfig(t, cfg)
	rcfg.RestampGen = true
	rcfg.Records = writeTree(t, map[string]string{
		"prov-2026-bbbbbbbb.yml": "id: prov-2026-bbbbbbbb\nintent: |\n  Another record.\naffected_scope:\n  - app/models/user.rb\n",
	})
	result, err := Replay(rec, rcfg)
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if len(result.Restamps) != 0 {
		t.Errorf("Restamps = %+v, want none for a record the records directory lacks", result.Restamps)
	}
	if got, want := markerGen(t, rcfg.Output), markerGen(t, grown); got != want {
		t.Errorf("gen = %q, want the recorded %q", got, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp_test.go","name":"TestRestampCompleteRecordingPrunesWhatItNoLongerSelects"},"kwargs":"sha256:38362133a62525cbefa01d945bd65c1d1279820a496a9f044dbc11955d1d626e","body":"sha256:ca89d7a98e2dac011ecd6299223155152101a35ca7c8deb78e9258fca62c385b"}
func TestRestampCompleteRecordingPrunesWhatItNoLongerSelects(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)

	rcfg := replayConfig(t, cfg)
	if _, err := Replay(rec, rcfg); err != nil {
		t.Fatalf("first Replay: %v", err)
	}
	if !strings.Contains(userFile(t, rcfg.Output), "email") {
		t.Fatalf("the first replay wrote no email region:\n%s", userFile(t, rcfg.Output))
	}

	rec.Complete = true
	rec.Records[0].Phases[0].Invocations[0].Kwargs["field"] = "name"
	result, err := Replay(rec, rcfg)
	if err != nil {
		t.Fatalf("second Replay: %v", err)
	}
	content := userFile(t, rcfg.Output)
	if strings.Contains(content, "email") || !strings.Contains(content, "name") {
		t.Errorf("a complete recording should leave only the name region:\n%s", content)
	}
	var pruned int
	for _, r := range result.Injections {
		if r.Pruned {
			pruned++
		}
	}
	if pruned != 1 {
		t.Errorf("%d region(s) pruned, want 1", pruned)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp_test.go","name":"TestRestampPartialRecordingReportsWhatItNoLongerSelects"},"kwargs":"sha256:cfd463aa96cb513426cd92290270d8454c6175e470eb6c576f963b6ee56b0a12","body":"sha256:54f6532e0584562e772940859d26e89d1b77d41718bd3aa73881554b14a3173a"}
func TestRestampPartialRecordingReportsWhatItNoLongerSelects(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)

	rcfg := replayConfig(t, cfg)
	if _, err := Replay(rec, rcfg); err != nil {
		t.Fatalf("first Replay: %v", err)
	}

	rec.Records[0].Phases[0].Invocations[0].Kwargs["field"] = "name"
	result, err := Replay(rec, rcfg)
	if err != nil {
		t.Fatalf("second Replay: %v", err)
	}
	content := userFile(t, rcfg.Output)
	if !strings.Contains(content, "email") || !strings.Contains(content, "name") {
		t.Errorf("a partial recording must prune nothing:\n%s", content)
	}
	var unselected int
	for _, r := range result.Injections {
		if r.Unselected {
			unselected++
		}
	}
	if unselected != 1 {
		t.Errorf("%d region(s) unselected, want 1", unselected)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp_test.go","name":"TestRestampReplayWouldChangeSeesAStaleGenAndAStrayRegion"},"kwargs":"sha256:d80ac5949c5af748b169c31c4530ba699c04eebb013cf5de784a566b87083f36","body":"sha256:9c7ca921337ecf9496741789636d45a534e32740f5c968c2b00fd27273d9c96c"}
func TestRestampReplayWouldChangeSeesAStaleGenAndAStrayRegion(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{response: withField}
	rec, _ := recorded(t, cfg)

	rcfg := replayConfig(t, cfg)
	if _, err := Replay(rec, rcfg); err != nil {
		t.Fatalf("first Replay: %v", err)
	}
	before := userFile(t, rcfg.Output)

	dry := rcfg
	dry.DryRun = true
	same, err := Replay(rec, dry)
	if err != nil {
		t.Fatalf("dry Replay: %v", err)
	}
	if ReplayWouldChange(same.Injections) {
		t.Error("replaying the recording over its own output should change nothing")
	}

	restamp := dry
	restamp.Records = editedRecords(t)
	restamp.RestampGen = true
	moved, err := Replay(rec, restamp)
	if err != nil {
		t.Fatalf("dry restamping Replay: %v", err)
	}
	if !ReplayWouldChange(moved.Injections) {
		t.Error("a restamp to a new gen should be a change")
	}

	rec.Records[0].Phases[0].Invocations[0].Kwargs["field"] = "name"
	stray, err := Replay(rec, dry)
	if err != nil {
		t.Fatalf("dry Replay of the edited recording: %v", err)
	}
	if !ReplayWouldChange(stray.Injections) {
		t.Error("a new region beside the old one should be a change")
	}
	if userFile(t, rcfg.Output) != before {
		t.Error("a dry replay wrote to the tree")
	}
	if got := ReplayWouldChange([]inject.Result{{Skipped: true}}); got {
		t.Error("a seeded region left as it is is not a change")
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp_test.go","name":"TestRestampStringShortensBothGens"},"kwargs":"sha256:c396ccbb3e2917e42048e569bfde4010cdff6c472a1c97636763e2e35f0ad694","body":"sha256:efc489d724b61b1e9ce2344caee166ff21e664712aa3123fc084255e747951bd"}
func TestRestampStringShortensBothGens(t *testing.T) {
	r := Restamp{RecordID: "prov-2026-aaaaaaaa", Old: "sha256:0123456789abcdef", New: "sha256:fedcba9876543210"}
	if got, want := r.String(), "restamped prov-2026-aaaaaaaa: gen 01234567 -> fedcba98"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	if got := (Restamp{RecordID: "x", New: "sha256:abcdef0123"}).String(); !strings.Contains(got, "none -> abcdef01") {
		t.Errorf("String() = %q, want the empty old gen shown as none", got)
	}
}

// /sedum:setFunction
