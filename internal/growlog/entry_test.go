package growlog

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry_test.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry_test.go","path":"unicode/utf8"}}
import "unicode/utf8"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAppendWritesOneLinePerEntry(t *testing.T) {\n\tpath := filepath.Join(t.TempDir(), \"sub\", \"grow.jsonl\")\n\tfor _, id := range []string{\"rec-1\", \"rec-2\"} {\n\t\tif err := Append(path, Entry{Record: id, Model: \"m\", Outcome: \"ok\"}); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\tdata, err := os.ReadFile(path)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tlines := strings.Split(strings.TrimSuffix(string(data), \"\\n\"), \"\\n\")\n\tif len(lines) != 2 {\n\t\tt.Fatalf(\"got %d lines, want 2: %q\", len(lines), data)\n\t}\n\tvar e Entry\n\tif err := json.Unmarshal([]byte(lines[0]), &e); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif e.V != Version || e.Record != \"rec-1\" {\n\t\tt.Fatalf(\"first entry = %+v\", e)\n\t}\n\tif strings.Contains(lines[0], `\"error\"`) {\n\t\tt.Fatalf(\"an empty error must be omitted: %s\", lines[0])\n\t}\n}","file":"internal/growlog/entry_test.go","name":"TestAppendWritesOneLinePerEntry"}}
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

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestAppendClipsError(t *testing.T) {\n\tcases := []struct {\n\t\tname, in string\n\t\tcheck    func(string) bool\n\t}{\n\t\t{\"first line only\", \"first\\nsecond\", func(s string) bool { return s == \"first\" }},\n\t\t{\"bounded\", strings.Repeat(\"x\", 400), func(s string) bool { return len(s) == 300 }},\n\t\t{\"whole runes\", strings.Repeat(\"é\", 200), func(s string) bool { return len(s) <= 300 && utf8.ValidString(s) }},\n\t}\n\tfor _, tc := range cases {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\tpath := filepath.Join(t.TempDir(), \"grow.jsonl\")\n\t\t\tif err := Append(path, Entry{Outcome: \"failed\", Error: tc.in}); err != nil {\n\t\t\t\tt.Fatal(err)\n\t\t\t}\n\t\t\tentries, _, err := Read(path)\n\t\t\tif err != nil || len(entries) != 1 {\n\t\t\t\tt.Fatalf(\"Read = %v, %v\", entries, err)\n\t\t\t}\n\t\t\tif !tc.check(entries[0].Error) {\n\t\t\t\tt.Fatalf(\"clipped error = %q\", entries[0].Error)\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/growlog/entry_test.go","name":"TestAppendClipsError"}}
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
			entries, _, err := Read(path)
			if err != nil || len(entries) != 1 {
				t.Fatalf("Read = %v, %v", entries, err)
			}
			if !tc.check(entries[0].Error) {
				t.Fatalf("clipped error = %q", entries[0].Error)
			}
		})
	}
}

// /sedum:setFunction
