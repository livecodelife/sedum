package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd_test.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd_test.go","path":"time"}}
import "time"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd_test.go","path":"github.com/livecodelife/sedum/internal/growlog"}}
import "github.com/livecodelife/sedum/internal/growlog"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestParseSince(t *testing.T) {\n\tnow := time.Date(2026, 10, 16, 12, 0, 0, 0, time.UTC)\n\tcases := []struct {\n\t\tin   string\n\t\twant time.Time\n\t\tbad  bool\n\t}{\n\t\t{\"\", time.Time{}, false},\n\t\t{\"7d\", time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC), false},\n\t\t{\"36h\", time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC), false},\n\t\t{\"2026-10-01\", time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), false},\n\t\t{\"2026-10-01T06:00:00Z\", time.Date(2026, 10, 1, 6, 0, 0, 0, time.UTC), false},\n\t\t{\"yesterday\", time.Time{}, true},\n\t}\n\tfor _, tc := range cases {\n\t\tgot, err := parseSince(tc.in, now)\n\t\tif (err != nil) != tc.bad {\n\t\t\tt.Errorf(\"parseSince(%q) error = %v, want error %v\", tc.in, err, tc.bad)\n\t\t\tcontinue\n\t\t}\n\t\tif !tc.bad && !got.Equal(tc.want) {\n\t\t\tt.Errorf(\"parseSince(%q) = %v, want %v\", tc.in, got, tc.want)\n\t\t}\n\t}\n}","file":"internal/cli/growlogcmd_test.go","name":"TestParseSince"}}
func TestParseSince(t *testing.T) {
	now := time.Date(2026, 10, 16, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		in   string
		want time.Time
		bad  bool
	}{
		{"", time.Time{}, false},
		{"7d", time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC), false},
		{"36h", time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC), false},
		{"2026-10-01", time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), false},
		{"2026-10-01T06:00:00Z", time.Date(2026, 10, 1, 6, 0, 0, 0, time.UTC), false},
		{"yesterday", time.Time{}, true},
	}
	for _, tc := range cases {
		got, err := parseSince(tc.in, now)
		if (err != nil) != tc.bad {
			t.Errorf("parseSince(%q) error = %v, want error %v", tc.in, err, tc.bad)
			continue
		}
		if !tc.bad && !got.Equal(tc.want) {
			t.Errorf("parseSince(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunGrowLog(t *testing.T) {\n\tpath := filepath.Join(t.TempDir(), \"grow.jsonl\")\n\tfor _, e := range []growlog.Entry{\n\t\t{TS: \"2026-10-10T00:00:00Z\", Record: \"r1\", Mode: \"select\", Model: \"haiku\", Outcome: \"ok\", Calls: 2, Rejected: 1},\n\t\t{TS: \"2026-10-15T00:00:00Z\", Record: \"r2\", Mode: \"select\", Model: \"stub\", Outcome: \"ok\", Calls: 1},\n\t} {\n\t\tif err := growlog.Append(path, e); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\tnow := time.Date(2026, 10, 16, 0, 0, 0, 0, time.UTC)\n\tvar out bytes.Buffer\n\tif err := runGrowLog(&out, path, false, \"\", \"\", \"\", now); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif !strings.Contains(out.String(), \"grow log: 2 line(s), 0 malformed\") || !strings.Contains(out.String(), \"haiku\") {\n\t\tt.Fatalf(\"text output:\\n%s\", out.String())\n\t}\n\tout.Reset()\n\tif err := runGrowLog(&out, path, true, \"3d\", \"\", \"\", now); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif !strings.Contains(out.String(), `\"entries\": 1`) {\n\t\tt.Fatalf(\"--since 3d should keep one entry:\\n%s\", out.String())\n\t}\n\tout.Reset()\n\tif err := runGrowLog(&out, path, true, \"\", \"\", \"haiku\", now); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif !strings.Contains(out.String(), `\"entries\": 1`) {\n\t\tt.Fatalf(\"--model haiku should keep one entry:\\n%s\", out.String())\n\t}\n}","file":"internal/cli/growlogcmd_test.go","name":"TestRunGrowLog"}}
func TestRunGrowLog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "grow.jsonl")
	for _, e := range []growlog.Entry{
		{TS: "2026-10-10T00:00:00Z", Record: "r1", Mode: "select", Model: "haiku", Outcome: "ok", Calls: 2, Rejected: 1},
		{TS: "2026-10-15T00:00:00Z", Record: "r2", Mode: "select", Model: "stub", Outcome: "ok", Calls: 1},
	} {
		if err := growlog.Append(path, e); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Date(2026, 10, 16, 0, 0, 0, 0, time.UTC)
	var out bytes.Buffer
	if err := runGrowLog(&out, path, false, "", "", "", now); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "grow log: 2 line(s), 0 malformed") || !strings.Contains(out.String(), "haiku") {
		t.Fatalf("text output:\n%s", out.String())
	}
	out.Reset()
	if err := runGrowLog(&out, path, true, "3d", "", "", now); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"entries": 1`) {
		t.Fatalf("--since 3d should keep one entry:\n%s", out.String())
	}
	out.Reset()
	if err := runGrowLog(&out, path, true, "", "", "haiku", now); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"entries": 1`) {
		t.Fatalf("--model haiku should keep one entry:\n%s", out.String())
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunGrowLogErrors(t *testing.T) {\n\tvar out bytes.Buffer\n\tif err := runGrowLog(&out, filepath.Join(t.TempDir(), \"none.jsonl\"), false, \"\", \"\", \"\", time.Now()); err == nil {\n\t\tt.Fatal(\"a missing log must be an error\")\n\t}\n\tpath := filepath.Join(t.TempDir(), \"grow.jsonl\")\n\tif err := growlog.Append(path, growlog.Entry{TS: \"2026-10-10T00:00:00Z\", Outcome: \"ok\"}); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := runGrowLog(&out, path, false, \"yesterday\", \"\", \"\", time.Now()); err == nil {\n\t\tt.Fatal(\"a bad --since must be an error\")\n\t}\n}","file":"internal/cli/growlogcmd_test.go","name":"TestRunGrowLogErrors"}}
func TestRunGrowLogErrors(t *testing.T) {
	var out bytes.Buffer
	if err := runGrowLog(&out, filepath.Join(t.TempDir(), "none.jsonl"), false, "", "", "", time.Now()); err == nil {
		t.Fatal("a missing log must be an error")
	}
	path := filepath.Join(t.TempDir(), "grow.jsonl")
	if err := growlog.Append(path, growlog.Entry{TS: "2026-10-10T00:00:00Z", Outcome: "ok"}); err != nil {
		t.Fatal(err)
	}
	if err := runGrowLog(&out, path, false, "yesterday", "", "", time.Now()); err == nil {
		t.Fatal("a bad --since must be an error")
	}
}

// /sedum:setFunction
