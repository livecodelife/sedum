package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/build_test.go","path":"bytes"},"kwargs":"sha256:4037ef22f542c176c4b00a0578700aeefa8134c8374d693d5673115d80e240b7","body":"sha256:468731f6dd513df0a66690d4649dc9e7244c96de227d8e6ead98cc5e9747c57a"}
import "bytes"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/build_test.go","path":"reflect"},"kwargs":"sha256:284d6e66995913de55bcb3abcda8cb93ca4ed75c95a6790490e7cb4a9538a225","body":"sha256:893cf368d31e1d8fe04189a0aba8627e8e69b39eafe10bea46d39a275e986fa2"}
import "reflect"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/build_test.go","path":"testing"},"kwargs":"sha256:f3f3688668cade5e2d7eec656534ffe02c4b306f6b8b494fc178f35f021d15ba","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/build_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:d10bfa475952973d63eab19966dc7e832bf71ce7716b32bcfa0a60ced389a7bb","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/build_test.go","path":"github.com/livecodelife/sedum/internal/record"},"kwargs":"sha256:8621b2c90afa076527c6eb37bf887bebcb3ed6e81f1fd183f7ce0883a52c36cb","body":"sha256:527c0befadfcad873cb5f9da4ab560755bc3040424ba6241011c3abf87fcc5e0"}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/build_test.go","path":"github.com/livecodelife/sedum/internal/recording"},"kwargs":"sha256:1aae84d7d8154c5f3715630d92bcafdebfef6489d24b644c2177a3b661911100","body":"sha256:77eda6b9872215a21291b42d59b9ec4b021705be869f51a3d61adc6dcb120e03"}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/build_test.go","name":"TestBuildRecordingTwoFiles"},"kwargs":"sha256:1dcfbe1ddd8bf00cb2b90541a84941a3f966e4deac6f15ec2d33ede57c0c17df","body":"sha256:0bac2d952b88994a60d4ba46e9c80baad0911879a56d46ee01eeb2a8766e2912"}
func TestBuildRecordingTwoFiles(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"}
	mk := func(path string) Adopted {
		return Adopted{
			Candidate: Candidate{Path: path, Package: pkg, Template: "t", Record: rec},
			Spans: []Span{
				{Kind: KindHeader, Text: "# header\n"},
				{Kind: KindUnit, Text: "a", Action: "one", Kwargs: map[string]any{"path": path}},
				{Kind: KindUnit, Text: "b", Action: "two", Kwargs: map[string]any{"path": path}},
			},
		}
	}
	got, err := BuildRecording([]Adopted{mk("b.py"), mk("a.py")})
	if err != nil {
		t.Fatal(err)
	}
	if got.SedumVersion != recording.Version {
		t.Errorf("version = %q", got.SedumVersion)
	}
	if !reflect.DeepEqual(got.Packages["py"].Extensions, []string{".py"}) {
		t.Errorf("extensions = %v", got.Packages["py"].Extensions)
	}
	if len(got.Records) != 1 {
		t.Fatalf("records = %d, want 1", len(got.Records))
	}
	r := got.Records[0]
	if r.RecordID != "imp-1" || r.Gen != "sha256:aa" || r.CreatedAt != "2026-10-16" {
		t.Errorf("record = %+v", r)
	}
	if len(r.Files) != 2 || r.Files[0].Path != "a.py" || r.Files[1].Path != "b.py" {
		t.Fatalf("files = %+v", r.Files)
	}
	inv := r.Phases[0].Invocations
	if len(inv) != 4 {
		t.Fatalf("invocations = %d, want 4", len(inv))
	}
	want := []struct{ path, action string }{{"a.py", "one"}, {"a.py", "two"}, {"b.py", "one"}, {"b.py", "two"}}
	for i, w := range want {
		if inv[i].Action != w.action || inv[i].Kwargs["path"] != w.path {
			t.Errorf("invocation %d = %+v", i, inv[i])
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/build_test.go","name":"TestBuildRecordingTwoRecords"},"kwargs":"sha256:ba790148683fcc7313f1940c310b0cbae2ac0107492c608905912734236ba370","body":"sha256:af386712ddb1fc01e3327d5114e39efa954e219c9e34a0c88bd3c79994757b6a"}
func TestBuildRecordingTwoRecords(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	unit := []Span{{Kind: KindUnit, Text: "x", Action: "one"}}
	items := []Adopted{
		{Candidate: Candidate{Path: "z.py", Package: pkg, Record: &record.Record{ID: "imp-2"}}, Spans: unit},
		{Candidate: Candidate{Path: "a.py", Package: pkg, Record: &record.Record{ID: "imp-1"}}, Spans: unit},
	}
	got, err := BuildRecording(items)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Records) != 2 || got.Records[0].RecordID != "imp-1" || got.Records[1].RecordID != "imp-2" {
		t.Errorf("records = %+v", got.Records)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/build_test.go","name":"TestBuildRecordingErrors"},"kwargs":"sha256:eeeb8af8c0f4890605e892843933463ed0711366aae1026b670c8948ed425ebf","body":"sha256:73b2cec6c76b6888e4947111e451082babc53eef3d033608d15782b2fe8436af"}
func TestBuildRecordingErrors(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1"}
	unit := Span{Kind: KindUnit, Text: "x", Action: "one"}
	tests := []struct {
		name  string
		items []Adopted
	}{
		{"no unit spans", []Adopted{{Candidate: Candidate{Path: "a.py", Package: pkg, Record: rec}, Spans: []Span{{Kind: KindHeader, Text: "h"}}}}},
		{"duplicate path", []Adopted{
			{Candidate: Candidate{Path: "a.py", Package: pkg, Record: rec}, Spans: []Span{unit}},
			{Candidate: Candidate{Path: "a.py", Package: pkg, Record: rec}, Spans: []Span{unit}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := BuildRecording(tt.items); err == nil {
				t.Fatal("want error")
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1e4a5690","gen":"sha256:dbab15369269c23a1776e610ab2bf128764a026cc1647f1f7b74bc855cc4a4c7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/build_test.go","name":"TestBuildRecordingRoundTripsThroughWrite"},"kwargs":"sha256:09d8c95640282cf53ecffca71a8f2f6c0aceabab02a757f90e23d4c2dac1411b","body":"sha256:556173ba57f318175804bc2871a1762788446321fcea4b665caabf3fee394fe2"}
func TestBuildRecordingRoundTripsThroughWrite(t *testing.T) {
	pkg := &genpkg.Package{Name: "py", Extensions: []string{".py"}}
	rec := &record.Record{ID: "imp-1", Gen: "sha256:aa", CreatedAt: "2026-10-16"}
	mk := func(path string) Adopted {
		return Adopted{
			Candidate: Candidate{Path: path, Package: pkg, Template: "t", Record: rec},
			Spans: []Span{
				{Kind: KindHeader, Text: "# header\n"},
				{Kind: KindUnit, Text: "a", Action: "one", Kwargs: map[string]any{"path": path}},
				{Kind: KindUnit, Text: "b", Action: "two", Kwargs: map[string]any{"path": path}},
			},
		}
	}
	got, err := BuildRecording([]Adopted{mk("b.py"), mk("a.py")})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := recording.Write(&buf, got); err != nil {
		t.Fatal(err)
	}
	back, err := recording.Read(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, got) {
		t.Errorf("round trip differs:\n got %+v\nwant %+v", back, got)
	}
}

// /sedum:setFunction
