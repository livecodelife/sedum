package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/orphans_test.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/orphans_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/orphans_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/orphans_test.go","path":"github.com/livecodelife/sedum/internal/pipeline"}}
import "github.com/livecodelife/sedum/internal/pipeline"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// orphanResult is a run with five removable orphans in three directories, one\n// pruned orphan, and one orphan in a file no current record authorizes.\nfunc orphanResult() *pipeline.Result {\n\tremovable := func(file, record string) pipeline.Orphan {\n\t\treturn pipeline.Orphan{Path: file, Action: \"setFunction\", Record: record, Authorized: true}\n\t}\n\treturn &pipeline.Result{Orphans: []pipeline.Orphan{\n\t\tremovable(\"vendor/a/x.go\", \"prov-2026-aaaaaaaa\"),\n\t\tremovable(\"vendor/a/x.go\", \"prov-2026-aaaaaaaa\"),\n\t\tremovable(\"vendor/a/y.go\", \"prov-2026-bbbbbbbb\"),\n\t\tremovable(\"vendor/b/z.go\", \"prov-2026-aaaaaaaa\"),\n\t\tremovable(\"lib/w.go\", \"prov-2026-cccccccc\"),\n\t\t{Path: \"lib/p.go\", Action: \"setFunction\", Record: \"prov-2026-dddddddd\", Authorized: true, Pruned: true},\n\t\t{Path: \"docs/q.go\", Action: \"setFunction\", Record: \"prov-2026-eeeeeeee\"},\n\t}}\n}","file":"internal/cli/orphans_test.go","name":"orphanResult"}}
// orphanResult is a run with five removable orphans in three directories, one
// pruned orphan, and one orphan in a file no current record authorizes.
func orphanResult() *pipeline.Result {
	removable := func(file, record string) pipeline.Orphan {
		return pipeline.Orphan{Path: file, Action: "setFunction", Record: record, Authorized: true}
	}
	return &pipeline.Result{Orphans: []pipeline.Orphan{
		removable("vendor/a/x.go", "prov-2026-aaaaaaaa"),
		removable("vendor/a/x.go", "prov-2026-aaaaaaaa"),
		removable("vendor/a/y.go", "prov-2026-bbbbbbbb"),
		removable("vendor/b/z.go", "prov-2026-aaaaaaaa"),
		removable("lib/w.go", "prov-2026-cccccccc"),
		{Path: "lib/p.go", Action: "setFunction", Record: "prov-2026-dddddddd", Authorized: true, Pruned: true},
		{Path: "docs/q.go", Action: "setFunction", Record: "prov-2026-eeeeeeee"},
	}}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func renderOrphans(result *pipeline.Result, list bool) string {\n\tvar out bytes.Buffer\n\tprintOrphans(&out, result, list)\n\treturn out.String()\n}","file":"internal/cli/orphans_test.go","name":"renderOrphans"}}
func renderOrphans(result *pipeline.Result, list bool) string {
	var out bytes.Buffer
	printOrphans(&out, result, list)
	return out.String()
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestOrphanSummaryGroupsByFate(t *testing.T) {\n\twant := \"\\n\" + strings.Join([]string{\n\t\t\"pruned 1 orphan region(s) in 1 file(s) from 1 gone or deprecated record(s); top directories: lib (1)\",\n\t\t\"5 orphan region(s) in 4 file(s) from 3 gone or deprecated record(s); --prune-orphans removes them; top directories: vendor/a (3), lib (1), vendor/b (1)\",\n\t\t\"1 orphan region(s) in 1 file(s) from 1 gone or deprecated record(s); left alone, no current record authorizes the file; top directories: docs (1)\",\n\t\t\"--orphans lists every region\",\n\t}, \"\\n\") + \"\\n\"\n\tif got := renderOrphans(orphanResult(), false); got != want {\n\t\tt.Errorf(\"summary differs\\n got: %q\\nwant: %q\", got, want)\n\t}\n}","file":"internal/cli/orphans_test.go","name":"TestOrphanSummaryGroupsByFate"}}
func TestOrphanSummaryGroupsByFate(t *testing.T) {
	want := "\n" + strings.Join([]string{
		"pruned 1 orphan region(s) in 1 file(s) from 1 gone or deprecated record(s); top directories: lib (1)",
		"5 orphan region(s) in 4 file(s) from 3 gone or deprecated record(s); --prune-orphans removes them; top directories: vendor/a (3), lib (1), vendor/b (1)",
		"1 orphan region(s) in 1 file(s) from 1 gone or deprecated record(s); left alone, no current record authorizes the file; top directories: docs (1)",
		"--orphans lists every region",
	}, "\n") + "\n"
	if got := renderOrphans(orphanResult(), false); got != want {
		t.Errorf("summary differs\n got: %q\nwant: %q", got, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestOrphanListKeepsALinePerRegion(t *testing.T) {\n\tgot := renderOrphans(orphanResult(), true)\n\tlines := strings.Split(strings.TrimSpace(got), \"\\n\")\n\tif len(lines) != 7 {\n\t\tt.Fatalf(\"got %d lines, want one per orphan (7):\\n%s\", len(lines), got)\n\t}\n\tif !strings.HasPrefix(lines[0], \"orphan        vendor/a/x.go  setFunction (record prov-2026-aaaaaaaa is gone or deprecated; --prune-orphans removes it)\") {\n\t\tt.Errorf(\"first line = %q\", lines[0])\n\t}\n\tif strings.Contains(got, \"--orphans lists\") {\n\t\tt.Errorf(\"the full list must not carry the summary hint:\\n%s\", got)\n\t}\n}","file":"internal/cli/orphans_test.go","name":"TestOrphanListKeepsALinePerRegion"}}
func TestOrphanListKeepsALinePerRegion(t *testing.T) {
	got := renderOrphans(orphanResult(), true)
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != 7 {
		t.Fatalf("got %d lines, want one per orphan (7):\n%s", len(lines), got)
	}
	if !strings.HasPrefix(lines[0], "orphan        vendor/a/x.go  setFunction (record prov-2026-aaaaaaaa is gone or deprecated; --prune-orphans removes it)") {
		t.Errorf("first line = %q", lines[0])
	}
	if strings.Contains(got, "--orphans lists") {
		t.Errorf("the full list must not carry the summary hint:\n%s", got)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestOrphanNoneIsSilent(t *testing.T) {\n\tfor _, list := range []bool{false, true} {\n\t\tif got := renderOrphans(&pipeline.Result{}, list); got != \"\" {\n\t\t\tt.Errorf(\"list=%v printed %q with no orphans\", list, got)\n\t\t}\n\t}\n}","file":"internal/cli/orphans_test.go","name":"TestOrphanNoneIsSilent"}}
func TestOrphanNoneIsSilent(t *testing.T) {
	for _, list := range []bool{false, true} {
		if got := renderOrphans(&pipeline.Result{}, list); got != "" {
			t.Errorf("list=%v printed %q with no orphans", list, got)
		}
	}
}

// /sedum:setFunction
