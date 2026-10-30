package growlog

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/read_test.go","path":"os"},"kwargs":"sha256:2626d1c37a5f61a33e5e9f56a8fd57fc09b79310d342289f966eb753957febfa","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/read_test.go","path":"path/filepath"},"kwargs":"sha256:b0fc956c38abc5172bb9a8a860c39082bffdc3b0d8cd2e53349c34aabebae3ab","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/read_test.go","path":"sync"},"kwargs":"sha256:ebc3d455bbc2e801723096e510a32ddcdf0e52566e800ec84f60ac562d1be748","body":"sha256:39e6702e4d5d7eda3a5f73722c09965d2dd3e73263f46b4f7f6df5a869228691"}
import "sync"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/read_test.go","path":"testing"},"kwargs":"sha256:7a4344b3780cd9242487b4e6919ba084e01342929da5ca1181e6aa69ba406775","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/read_test.go","name":"TestReadSkipsMalformedLines"},"kwargs":"sha256:d08b64811cb67da88a0538c5695b7356c22aaa5fdc3cac22465204b305a91e2d","body":"sha256:5ec48d3bc1c368326b64a486fa1cc704a533e6cf8fb5422b412f6c22123bb287"}
func TestReadSkipsMalformedLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "grow.jsonl")
	data := "{\"v\":1,\"record\":\"a\"}\n\nnot json\n{\"v\":1,\"record\":\"b\"}\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	entries, malformed, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 || malformed != 1 {
		t.Fatalf("entries = %d, malformed = %d, want 2 and 1", len(entries), malformed)
	}
	if entries[0].Record != "a" || entries[1].Record != "b" {
		t.Fatalf("entries out of order: %+v", entries)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/read_test.go","name":"TestReadMissingFile"},"kwargs":"sha256:bf8e31955924ea9f32569551e62c13336d6d6699f77738444e70e4cb26a8afea","body":"sha256:f6c7b8e65becadede7e321a191c8f97719cb4427a8b6d560a4efd42fb69fb60f"}
func TestReadMissingFile(t *testing.T) {
	if _, _, err := Read(filepath.Join(t.TempDir(), "none.jsonl")); err == nil {
		t.Fatal("want an error for a missing file")
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/read_test.go","name":"TestAppendConcurrent"},"kwargs":"sha256:7c7ad133f90c4ab8cfdc44088fa4bfbed5950cd4014e3c2371778b978c9d04fe","body":"sha256:c8e037e74944f52ca0647ab60509b5cba7ece1737a2f7c38db0302042374db69"}
func TestAppendConcurrent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "grow.jsonl")
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := Append(path, Entry{Record: "r", Model: "m", Outcome: "ok"}); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	entries, malformed, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 20 || malformed != 0 {
		t.Fatalf("entries = %d, malformed = %d, want 20 and 0", len(entries), malformed)
	}
}

// /sedum:setFunction
