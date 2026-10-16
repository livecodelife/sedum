package pipeline

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/restamp_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/restamp_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/restamp_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/restamp_test.go","path":"github.com/livecodelife/sedum/internal/inject"}}
import "github.com/livecodelife/sedum/internal/inject"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// editedRecords is the record the recordings below were made from, edited\n// since: same id and scope, a new intent, so a new gen.\nfunc editedRecords(t *testing.T) string {\n\tt.Helper()\n\treturn writeTree(t, map[string]string{\n\t\t\"prov-2026-aaaaaaaa.yml\": \"id: prov-2026-aaaaaaaa\\nintent: |\\n  Something else.\\naffected_scope:\\n  - app/models/user.rb\\n\",\n\t})\n}","file":"internal/pipeline/restamp_test.go","name":"editedRecords"}}
// editedRecords is the record the recordings below were made from, edited
// since: same id and scope, a new intent, so a new gen.
func editedRecords(t *testing.T) string {
	t.Helper()
	return writeTree(t, map[string]string{
		"prov-2026-aaaaaaaa.yml": "id: prov-2026-aaaaaaaa\nintent: |\n  Something else.\naffected_scope:\n  - app/models/user.rb\n",
	})
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// userFile is the one file every recording below writes.\nfunc userFile(t *testing.T, output string) string {\n\tt.Helper()\n\treturn readFile(t, filepath.Join(output, \"app\", \"models\", \"user.rb\"))\n}","file":"internal/pipeline/restamp_test.go","name":"userFile"}}
// userFile is the one file every recording below writes.
func userFile(t *testing.T, output string) string {
	t.Helper()
	return readFile(t, filepath.Join(output, "app", "models", "user.rb"))
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRestampGenTakesTheCurrentGen(t *testing.T) {\n\tcfg := config(t)\n\tcfg.Client = &stub{response: withField}\n\trec, grown := recorded(t, cfg)\n\trecordedGen := markerGen(t, grown)\n\n\trcfg := replayConfig(t, cfg)\n\trcfg.Records = editedRecords(t)\n\trcfg.RestampGen = true\n\tresult, err := Replay(rec, rcfg)\n\tif err != nil {\n\t\tt.Fatalf(\"Replay: %v\", err)\n\t}\n\tcurrent := result.Records.Records[0].Gen\n\tif current == recordedGen {\n\t\tt.Fatal(\"the edited record hashes the same as the original\")\n\t}\n\tif got := markerGen(t, rcfg.Output); got != current {\n\t\tt.Errorf(\"gen = %q, want the current record's %q\", got, current)\n\t}\n\tif len(result.Restamps) != 1 || result.Restamps[0].Old != recordedGen || result.Restamps[0].New != current {\n\t\tt.Errorf(\"Restamps = %+v, want one from %q to %q\", result.Restamps, recordedGen, current)\n\t}\n\n\tplain := replayConfig(t, cfg)\n\tplain.Records = rcfg.Records\n\tif _, err := Replay(rec, plain); err != nil {\n\t\tt.Fatalf(\"Replay without the flag: %v\", err)\n\t}\n\tif got := markerGen(t, plain.Output); got != recordedGen {\n\t\tt.Errorf(\"without --restamp-gen the gen = %q, want the recorded %q\", got, recordedGen)\n\t}\n}","file":"internal/pipeline/restamp_test.go","name":"TestRestampGenTakesTheCurrentGen"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRestampGenNeedsRecords(t *testing.T) {\n\tcfg := config(t)\n\tcfg.Client = &stub{response: withField}\n\trec, _ := recorded(t, cfg)\n\n\trcfg := replayConfig(t, cfg)\n\trcfg.RestampGen = true\n\tif _, err := Replay(rec, rcfg); err == nil || !strings.Contains(err.Error(), \"--records\") {\n\t\tt.Fatalf(\"Replay error = %v, want one naming --records\", err)\n\t}\n}","file":"internal/pipeline/restamp_test.go","name":"TestRestampGenNeedsRecords"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRestampGenSkipsARecordThatIsNotCurrent(t *testing.T) {\n\tcfg := config(t)\n\tcfg.Client = &stub{response: withField}\n\trec, grown := recorded(t, cfg)\n\n\trcfg := replayConfig(t, cfg)\n\trcfg.RestampGen = true\n\trcfg.Records = writeTree(t, map[string]string{\n\t\t\"prov-2026-bbbbbbbb.yml\": \"id: prov-2026-bbbbbbbb\\nintent: |\\n  Another record.\\naffected_scope:\\n  - app/models/user.rb\\n\",\n\t})\n\tresult, err := Replay(rec, rcfg)\n\tif err != nil {\n\t\tt.Fatalf(\"Replay: %v\", err)\n\t}\n\tif len(result.Restamps) != 0 {\n\t\tt.Errorf(\"Restamps = %+v, want none for a record the records directory lacks\", result.Restamps)\n\t}\n\tif got, want := markerGen(t, rcfg.Output), markerGen(t, grown); got != want {\n\t\tt.Errorf(\"gen = %q, want the recorded %q\", got, want)\n\t}\n}","file":"internal/pipeline/restamp_test.go","name":"TestRestampGenSkipsARecordThatIsNotCurrent"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRestampCompleteRecordingPrunesWhatItNoLongerSelects(t *testing.T) {\n\tcfg := config(t)\n\tcfg.Client = &stub{response: withField}\n\trec, _ := recorded(t, cfg)\n\n\trcfg := replayConfig(t, cfg)\n\tif _, err := Replay(rec, rcfg); err != nil {\n\t\tt.Fatalf(\"first Replay: %v\", err)\n\t}\n\tif !strings.Contains(userFile(t, rcfg.Output), \"email\") {\n\t\tt.Fatalf(\"the first replay wrote no email region:\\n%s\", userFile(t, rcfg.Output))\n\t}\n\n\trec.Complete = true\n\trec.Records[0].Phases[0].Invocations[0].Kwargs[\"field\"] = \"name\"\n\tresult, err := Replay(rec, rcfg)\n\tif err != nil {\n\t\tt.Fatalf(\"second Replay: %v\", err)\n\t}\n\tcontent := userFile(t, rcfg.Output)\n\tif strings.Contains(content, \"email\") || !strings.Contains(content, \"name\") {\n\t\tt.Errorf(\"a complete recording should leave only the name region:\\n%s\", content)\n\t}\n\tvar pruned int\n\tfor _, r := range result.Injections {\n\t\tif r.Pruned {\n\t\t\tpruned++\n\t\t}\n\t}\n\tif pruned != 1 {\n\t\tt.Errorf(\"%d region(s) pruned, want 1\", pruned)\n\t}\n}","file":"internal/pipeline/restamp_test.go","name":"TestRestampCompleteRecordingPrunesWhatItNoLongerSelects"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRestampPartialRecordingReportsWhatItNoLongerSelects(t *testing.T) {\n\tcfg := config(t)\n\tcfg.Client = &stub{response: withField}\n\trec, _ := recorded(t, cfg)\n\n\trcfg := replayConfig(t, cfg)\n\tif _, err := Replay(rec, rcfg); err != nil {\n\t\tt.Fatalf(\"first Replay: %v\", err)\n\t}\n\n\trec.Records[0].Phases[0].Invocations[0].Kwargs[\"field\"] = \"name\"\n\tresult, err := Replay(rec, rcfg)\n\tif err != nil {\n\t\tt.Fatalf(\"second Replay: %v\", err)\n\t}\n\tcontent := userFile(t, rcfg.Output)\n\tif !strings.Contains(content, \"email\") || !strings.Contains(content, \"name\") {\n\t\tt.Errorf(\"a partial recording must prune nothing:\\n%s\", content)\n\t}\n\tvar unselected int\n\tfor _, r := range result.Injections {\n\t\tif r.Unselected {\n\t\t\tunselected++\n\t\t}\n\t}\n\tif unselected != 1 {\n\t\tt.Errorf(\"%d region(s) unselected, want 1\", unselected)\n\t}\n}","file":"internal/pipeline/restamp_test.go","name":"TestRestampPartialRecordingReportsWhatItNoLongerSelects"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRestampReplayWouldChangeSeesAStaleGenAndAStrayRegion(t *testing.T) {\n\tcfg := config(t)\n\tcfg.Client = &stub{response: withField}\n\trec, _ := recorded(t, cfg)\n\n\trcfg := replayConfig(t, cfg)\n\tif _, err := Replay(rec, rcfg); err != nil {\n\t\tt.Fatalf(\"first Replay: %v\", err)\n\t}\n\tbefore := userFile(t, rcfg.Output)\n\n\tdry := rcfg\n\tdry.DryRun = true\n\tsame, err := Replay(rec, dry)\n\tif err != nil {\n\t\tt.Fatalf(\"dry Replay: %v\", err)\n\t}\n\tif ReplayWouldChange(same.Injections) {\n\t\tt.Error(\"replaying the recording over its own output should change nothing\")\n\t}\n\n\trestamp := dry\n\trestamp.Records = editedRecords(t)\n\trestamp.RestampGen = true\n\tmoved, err := Replay(rec, restamp)\n\tif err != nil {\n\t\tt.Fatalf(\"dry restamping Replay: %v\", err)\n\t}\n\tif !ReplayWouldChange(moved.Injections) {\n\t\tt.Error(\"a restamp to a new gen should be a change\")\n\t}\n\n\trec.Records[0].Phases[0].Invocations[0].Kwargs[\"field\"] = \"name\"\n\tstray, err := Replay(rec, dry)\n\tif err != nil {\n\t\tt.Fatalf(\"dry Replay of the edited recording: %v\", err)\n\t}\n\tif !ReplayWouldChange(stray.Injections) {\n\t\tt.Error(\"a new region beside the old one should be a change\")\n\t}\n\tif userFile(t, rcfg.Output) != before {\n\t\tt.Error(\"a dry replay wrote to the tree\")\n\t}\n\tif got := ReplayWouldChange([]inject.Result{{Skipped: true}}); got {\n\t\tt.Error(\"a seeded region left as it is is not a change\")\n\t}\n}","file":"internal/pipeline/restamp_test.go","name":"TestRestampReplayWouldChangeSeesAStaleGenAndAStrayRegion"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-5176a0b7","gen":"sha256:bda774e8cf8be21135df66b111a07aabae9ab4e706b991c710f4655a16764214","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRestampStringShortensBothGens(t *testing.T) {\n\tr := Restamp{RecordID: \"prov-2026-aaaaaaaa\", Old: \"sha256:0123456789abcdef\", New: \"sha256:fedcba9876543210\"}\n\tif got, want := r.String(), \"restamped prov-2026-aaaaaaaa: gen 01234567 -> fedcba98\"; got != want {\n\t\tt.Errorf(\"String() = %q, want %q\", got, want)\n\t}\n\tif got := (Restamp{RecordID: \"x\", New: \"sha256:abcdef0123\"}).String(); !strings.Contains(got, \"none -> abcdef01\") {\n\t\tt.Errorf(\"String() = %q, want the empty old gen shown as none\", got)\n\t}\n}","file":"internal/pipeline/restamp_test.go","name":"TestRestampStringShortensBothGens"}}
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
