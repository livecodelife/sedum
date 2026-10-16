package growlog

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/read_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/read_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/read_test.go","path":"sync"}}
import "sync"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/read_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestReadSkipsMalformedLines(t *testing.T) {\n\tpath := filepath.Join(t.TempDir(), \"grow.jsonl\")\n\tdata := \"{\\\"v\\\":1,\\\"record\\\":\\\"a\\\"}\\n\\nnot json\\n{\\\"v\\\":1,\\\"record\\\":\\\"b\\\"}\\n\"\n\tif err := os.WriteFile(path, []byte(data), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tentries, malformed, err := Read(path)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(entries) != 2 || malformed != 1 {\n\t\tt.Fatalf(\"entries = %d, malformed = %d, want 2 and 1\", len(entries), malformed)\n\t}\n\tif entries[0].Record != \"a\" || entries[1].Record != \"b\" {\n\t\tt.Fatalf(\"entries out of order: %+v\", entries)\n\t}\n}","file":"internal/growlog/read_test.go","name":"TestReadSkipsMalformedLines"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestReadMissingFile(t *testing.T) {\n\tif _, _, err := Read(filepath.Join(t.TempDir(), \"none.jsonl\")); err == nil {\n\t\tt.Fatal(\"want an error for a missing file\")\n\t}\n}","file":"internal/growlog/read_test.go","name":"TestReadMissingFile"}}
func TestReadMissingFile(t *testing.T) {
	if _, _, err := Read(filepath.Join(t.TempDir(), "none.jsonl")); err == nil {
		t.Fatal("want an error for a missing file")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAppendConcurrent(t *testing.T) {\n\tpath := filepath.Join(t.TempDir(), \"grow.jsonl\")\n\tvar wg sync.WaitGroup\n\tfor i := 0; i < 20; i++ {\n\t\twg.Add(1)\n\t\tgo func() {\n\t\t\tdefer wg.Done()\n\t\t\tif err := Append(path, Entry{Record: \"r\", Model: \"m\", Outcome: \"ok\"}); err != nil {\n\t\t\t\tt.Error(err)\n\t\t\t}\n\t\t}()\n\t}\n\twg.Wait()\n\tentries, malformed, err := Read(path)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(entries) != 20 || malformed != 0 {\n\t\tt.Fatalf(\"entries = %d, malformed = %d, want 20 and 0\", len(entries), malformed)\n\t}\n}","file":"internal/growlog/read_test.go","name":"TestAppendConcurrent"}}
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
