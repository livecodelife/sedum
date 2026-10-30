package cli

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd_test.go","path":"bytes"},"kwargs":"sha256:f0faa2fffa6a7de3306cae5f6311260e18f8e5bf9de12c624c41119916f02655","body":"sha256:468731f6dd513df0a66690d4649dc9e7244c96de227d8e6ead98cc5e9747c57a"}
import "bytes"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd_test.go","path":"path/filepath"},"kwargs":"sha256:49d80ecc24244c53a04721848ae3bd7a63c642d1f23ae6ccf54b89c969d1c822","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd_test.go","path":"strings"},"kwargs":"sha256:578f3ce46021d682b9f1f8a882679d8884823194bc9ba089e456e920d05c631d","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd_test.go","path":"testing"},"kwargs":"sha256:8a87cd78977fdb2a197f9e0a895895c233d6abb575674b2dd6e6b2e87f108316","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd_test.go","path":"time"},"kwargs":"sha256:bef585836c90aea98163b2a83ef1dc83fa34b84696ad10ef9dba55c02a31f78d","body":"sha256:e99128ebf4e36998f40b33adf2eef57f4d6c315f1a4d1bdd3b7fc09a0d7bc82d"}
import "time"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd_test.go","path":"github.com/livecodelife/sedum/internal/growlog"},"kwargs":"sha256:0502c7fda0bed26466f16514da4178cc4aac67ec98810757be916ab15276b2a3","body":"sha256:b7968f24df7d76473f917ee9527158f424f930250a91e5407a0140aa4c56c91f"}
import "github.com/livecodelife/sedum/internal/growlog"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/growlogcmd_test.go","name":"TestParseSince"},"kwargs":"sha256:121ebdaef6b72e0f6024126c20664c8e4516f1d76aa645e7c9bd2a5483dd62cb","body":"sha256:4d6fe39d866e7e821a621a496cbaa2e745bf2c81821d7bc145a223c2d793a380"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/growlogcmd_test.go","name":"TestRunGrowLog"},"kwargs":"sha256:a978e12b901b0e26dbfdeb00dbaea65fd2093a87df925d6d466b6f7120d89a3e","body":"sha256:ed4d1341e289b45e788e244f442239dd3b9c87b7bc853d204ba847431eb8de50"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/growlogcmd_test.go","name":"TestRunGrowLogErrors"},"kwargs":"sha256:ab766affe481fe6e5cc1a0361673d66036354ac96d8d5287d5fabff036ec4f43","body":"sha256:33eed7d495946048ab7bddde2fec28ab531cce98e26754de47d6f5c00bf0c740"}
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
