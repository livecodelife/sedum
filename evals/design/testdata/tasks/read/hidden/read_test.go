package growlog

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

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

func TestReadMissingFile(t *testing.T) {
	if _, _, err := Read(filepath.Join(t.TempDir(), "none.jsonl")); err == nil {
		t.Fatal("want an error for a missing file")
	}
}

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
