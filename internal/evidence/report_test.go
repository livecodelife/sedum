package evidence

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// Judge gives a record's freshness verdict from its regions and its evidence
// entries, and WriteText and WriteJSON write the reports.

const stampedHash = "sha256:9aaa1fe3e1be19258596ec7d6e2c8185d8797c55ba85159bc80c3fc528fbd904"

func stampedRegions() []Region {
	return []Region{{Gen: "sha256:aa", Pkg: "sha256:01"}, {Gen: "sha256:bb", Pkg: "sha256:02"}}
}

func TestJudgeVerdicts(t *testing.T) {
	now := stampedHash
	cases := []struct {
		name    string
		regions []Region
		entries []Entry
		want    Result
	}{
		{"no regions", nil, nil, Result{Record: "r", Version: 3, Verdict: NoRegions}},
		{"only a seeded region", []Region{{Tier: "seeded"}}, nil, Result{Record: "r", Version: 3, Verdict: NoRegions}},
		{"an unstamped region", []Region{{Gen: "sha256:aa"}}, nil, Result{Record: "r", Version: 3, Verdict: Unknown, Regions: 1}},
		{"unknown wins over none", []Region{{Gen: "sha256:aa"}}, []Entry{entry("r", 3, "sha256:old")}, Result{Record: "r", Version: 3, Verdict: Unknown, Regions: 1}},
		{"none without an entry", stampedRegions(), nil, Result{Record: "r", Version: 3, Verdict: None, CodeHash: now, Regions: 2}},
		{"none when only an earlier version has an entry", stampedRegions(), []Entry{entry("r", 2, now)}, Result{Record: "r", Version: 3, Verdict: None, CodeHash: now, Regions: 2}},
		{"fresh", stampedRegions(), []Entry{entry("r", 3, now)}, Result{Record: "r", Version: 3, Verdict: Fresh, CodeHash: now, Recorded: now, Regions: 2}},
		{"stale", stampedRegions(), []Entry{entry("r", 3, "sha256:old")}, Result{Record: "r", Version: 3, Verdict: Stale, CodeHash: now, Recorded: "sha256:old", Regions: 2}},
		{"the latest entry of the version decides", stampedRegions(), []Entry{entry("r", 3, now), entry("r", 3, "sha256:newer")}, Result{Record: "r", Version: 3, Verdict: Stale, CodeHash: now, Recorded: "sha256:newer", Regions: 2}},
	}
	for _, c := range cases {
		if got := Judge("r", 3, c.regions, c.entries); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s:\n got %+v\nwant %+v", c.name, got, c.want)
		}
	}
}

func TestFailingIsTrueForStaleOrUnknownOnly(t *testing.T) {
	for _, c := range []struct {
		verdicts []Verdict
		want     bool
	}{
		{nil, false},
		{[]Verdict{Fresh, None, NoRegions}, false},
		{[]Verdict{Fresh, Stale}, true},
		{[]Verdict{Unknown}, true},
	} {
		var rs []Result
		for _, v := range c.verdicts {
			rs = append(rs, Result{Verdict: v})
		}
		if got := Failing(rs); got != c.want {
			t.Errorf("Failing(%v) = %v, want %v", c.verdicts, got, c.want)
		}
	}
}

func results() []Result {
	return []Result{
		{Record: "prov-2026-bbbbbbbb", Version: 2, Verdict: Stale, CodeHash: "sha256:0123456789abcdef0123", Recorded: "sha256:fedcba9876543210ffff", Regions: 2},
		{Record: "prov-2026-aaaaaaaa", Version: 3, Verdict: None, CodeHash: "sha256:00112233445566778899", Regions: 1},
		{Record: "prov-2026-cccccccc", Version: 1, Verdict: NoRegions},
		{Record: "prov-2026-dddddddd", Version: 1, Verdict: Unknown, Regions: 1},
		{Record: "prov-2026-eeeeeeee", Version: 1, Verdict: Fresh, CodeHash: "sha256:abababababababababab", Recorded: "sha256:abababababababababab", Regions: 1},
	}
}

func TestWriteTextWritesOneLinePerRecordInIdOrderThenTheCounts(t *testing.T) {
	in := results()
	var out bytes.Buffer

	if err := WriteText(&out, in); err != nil {
		t.Fatal(err)
	}

	want := "prov-2026-aaaaaaaa  none  now=001122334455 recorded=-\n" +
		"prov-2026-bbbbbbbb  stale  now=0123456789ab recorded=fedcba987654\n" +
		"prov-2026-cccccccc  no regions  now=- recorded=-\n" +
		"prov-2026-dddddddd  unknown  now=- recorded=-\n" +
		"prov-2026-eeeeeeee  fresh  now=abababababab recorded=abababababab\n" +
		"\n" +
		"1 fresh, 1 stale, 1 none, 1 unknown, 1 no regions\n"
	if out.String() != want {
		t.Errorf("output =\n%q\nwant\n%q", out.String(), want)
	}
	if in[0].Record != "prov-2026-bbbbbbbb" {
		t.Error("WriteText sorted the caller's slice")
	}
}

func TestWriteTextCountsEmptyInput(t *testing.T) {
	var out bytes.Buffer
	if err := WriteText(&out, nil); err != nil {
		t.Fatal(err)
	}
	if out.String() != "\n0 fresh, 0 stale, 0 none, 0 unknown, 0 no regions\n" {
		t.Errorf("output = %q", out.String())
	}
}

func TestWriteJSONWritesTheDocumentIndentedSortedAndWithoutHTMLEscaping(t *testing.T) {
	in := results()
	in[0].Recorded = "sha256:<&>"
	var out bytes.Buffer

	if err := WriteJSON(&out, in); err != nil {
		t.Fatal(err)
	}

	text := out.String()
	if !strings.HasPrefix(text, "{\n  \"format\": \"sedum.evidence/1\",\n  \"records\": [\n    {\n      \"record\": \"prov-2026-aaaaaaaa\",") || !strings.HasSuffix(text, "}\n") {
		t.Errorf("document = %s", text)
	}
	if !strings.Contains(text, "sha256:<&>") {
		t.Error("HTML was escaped")
	}
	var doc Document
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil || doc.Format != Format || len(doc.Records) != 5 || doc.Records[4].Record != "prov-2026-eeeeeeee" {
		t.Errorf("doc = %+v err = %v", doc, err)
	}
	if doc.Records[1].Verdict != Stale || doc.Records[1].Regions != 2 || doc.Records[1].Recorded != "sha256:<&>" {
		t.Errorf("record = %+v", doc.Records[1])
	}
	if in[0].Record != "prov-2026-bbbbbbbb" {
		t.Error("WriteJSON sorted the caller's slice")
	}
}

func TestWriteJSONWritesAnEmptyArrayNotNull(t *testing.T) {
	var out bytes.Buffer
	if err := WriteJSON(&out, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"records": []`) {
		t.Errorf("document = %s", out.String())
	}
}

func TestResultFieldsAreNamedInJSON(t *testing.T) {
	data, err := json.Marshal(Result{Record: "r", Version: 1, Verdict: Fresh, CodeHash: "c", Recorded: "d", Regions: 2})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"record":"r","version":1,"verdict":"fresh","code_hash":"c","recorded":"d","regions":2}` {
		t.Errorf("json = %s", data)
	}
}
