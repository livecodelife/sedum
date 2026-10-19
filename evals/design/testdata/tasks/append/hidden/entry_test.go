package growlog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestAppendWritesOneLinePerEntry(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "grow.jsonl")
	for _, id := range []string{"rec-1", "rec-2"} {
		if err := Append(path, Entry{Record: id, Model: "m", Outcome: "ok"}); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %q", len(lines), data)
	}
	var e Entry
	if err := json.Unmarshal([]byte(lines[0]), &e); err != nil {
		t.Fatal(err)
	}
	if e.V != Version || e.Record != "rec-1" {
		t.Fatalf("first entry = %+v", e)
	}
	if strings.Contains(lines[0], `"error"`) {
		t.Fatalf("an empty error must be omitted: %s", lines[0])
	}
}

func TestAppendClipsError(t *testing.T) {
	cases := []struct {
		name, in string
		check    func(string) bool
	}{
		{"first line only", "first\nsecond", func(s string) bool { return s == "first" }},
		{"bounded", strings.Repeat("x", 400), func(s string) bool { return len(s) == 300 }},
		{"whole runes", strings.Repeat("é", 200), func(s string) bool { return len(s) <= 300 && utf8.ValidString(s) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "grow.jsonl")
			if err := Append(path, Entry{Outcome: "failed", Error: tc.in}); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var e Entry
			if err := json.Unmarshal(data, &e); err != nil {
				t.Fatal(err)
			}
			if !tc.check(e.Error) {
				t.Fatalf("clipped error = %q", e.Error)
			}
		})
	}
}
