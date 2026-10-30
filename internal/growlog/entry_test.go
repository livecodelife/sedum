package growlog

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry_test.go","path":"encoding/json"},"kwargs":"sha256:e7d5b2502bcda5e8f4d4342b1e1d25aff90e715405b448f126566ba0fe1d088b","body":"sha256:872d7c3e9dee608e79a7f8ef57930d24b5e30f2a57d891a52549179ace45946a"}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry_test.go","path":"os"},"kwargs":"sha256:f66c0a387a9ceebe8a459c9714bdb9d52c66a5b8f36238c6eb37e7ede8e1f892","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry_test.go","path":"path/filepath"},"kwargs":"sha256:b5d836cdc768310f27998e028410386e5cc7a924deb78a8ea4875899a9af0cd2","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry_test.go","path":"strings"},"kwargs":"sha256:861f8d78f2fa7519634fa5cbfff48fc6ca6e833cc91e02060277cc7dbe06c7d1","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry_test.go","path":"testing"},"kwargs":"sha256:a5c283ea018594e5ce7cb432d58fa5cf7a3c182ea1968eb969b3d67ae4592d4c","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry_test.go","path":"unicode/utf8"},"kwargs":"sha256:0485f123a9a956e3dccaba9bfba10bf231a6673484e13355bcfbc1b2437cee1a","body":"sha256:329f7334124c7cfe72d8208daedf694eae268e0a4bd205bdde0ead7fbb2a0037"}
import "unicode/utf8"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/entry_test.go","name":"TestAppendWritesOneLinePerEntry"},"kwargs":"sha256:aca95916e635030da6ecfeedde67366fb98a3d424e014e709b045397ada4c4ee","body":"sha256:fde8c4d5e5be8c7b535fb2118e1508d2837724855f8f93464c9cc221287fad38"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/entry_test.go","name":"TestAppendClipsError"},"kwargs":"sha256:43e1875e6d6454fc9299e77729b1a8085d4d3e39d9dd18be8b62481a20054a7a","body":"sha256:4275914965102566940b992089e463882faa9cb8c38c2703a4578a0e0a75cdd9"}
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

// /sedum:setFunction
