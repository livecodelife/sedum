package cli

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/orphans_test.go","path":"bytes"},"kwargs":"sha256:a75a8e140971f5670bf3852637203264b0bbcef1a8a772982ae7260c121c6206","body":"sha256:468731f6dd513df0a66690d4649dc9e7244c96de227d8e6ead98cc5e9747c57a"}
import "bytes"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/orphans_test.go","path":"strings"},"kwargs":"sha256:2d5b54402c6e20e39a7639e55349821d1cd6f98b810bf8265217339c5a747c94","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/orphans_test.go","path":"testing"},"kwargs":"sha256:dcd5228c282aded7cbbe1e90246da2f3e2deac094229104b631d83ee5f327a19","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/orphans_test.go","path":"github.com/livecodelife/sedum/internal/pipeline"},"kwargs":"sha256:814b56ebb550fadb35a0a8906323f4df2f3074a89ddcd2cbec498147d76f5bce","body":"sha256:4dfadc322ca3e8f8dc06b68db404b84dcb9cf1254d835d4ed3eedc1206db5fe8"}
import "github.com/livecodelife/sedum/internal/pipeline"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/orphans_test.go","name":"orphanResult"},"kwargs":"sha256:a79ba9214373eeffdba521ee824074fdbcfd19195c03952e6b02a93fc04fdcbf","body":"sha256:9f948fea50cb1df870e93e85078c5e03420232b28b4bb37caeb3a1c790ae367b"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/orphans_test.go","name":"renderOrphans"},"kwargs":"sha256:993ab0d446bec8a670ef2d3c0a26a9f330bc2784249670d3105845924a701fc8","body":"sha256:658091ec470c5fbea6113e7ef43e36919b3ea98a5b68ac47f97fc93d2c484097"}
func renderOrphans(result *pipeline.Result, list bool) string {
	var out bytes.Buffer
	printOrphans(&out, result, list)
	return out.String()
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/orphans_test.go","name":"TestOrphanSummaryGroupsByFate"},"kwargs":"sha256:cd400eb27753a0ad43909480aa02e8927cacc2674d3edb2ec65886f995b7c80d","body":"sha256:5abf65c5e124e48603f437cafbd1f9a0d5f97a6f29edce12f6f229f105ecf0e4"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/orphans_test.go","name":"TestOrphanListKeepsALinePerRegion"},"kwargs":"sha256:b0dcda6e2e5d7e71ba62e934ee54b284ee35850264609da487eadd39a42a8b8e","body":"sha256:0a03d5a44091ee1d3977e15bb897edf7c605926e0457d6472f407e268c07e2f2"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-94c524b5","gen":"sha256:6621721a11158ffecb44372d47c51eb127b6577c8f430f297fbdfc07ffb8fb0f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/orphans_test.go","name":"TestOrphanNoneIsSilent"},"kwargs":"sha256:d6138d2f2392b2cb586ff610d987b93f022e5babb95a04faa31b624328bbcb38","body":"sha256:729d381f74c77e431e3e955ea59e72919a41aa397df13ccdac3f25326b7669f2"}
func TestOrphanNoneIsSilent(t *testing.T) {
	for _, list := range []bool{false, true} {
		if got := renderOrphans(&pipeline.Result{}, list); got != "" {
			t.Errorf("list=%v printed %q with no orphans", list, got)
		}
	}
}

// /sedum:setFunction
