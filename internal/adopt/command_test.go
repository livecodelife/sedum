package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestParseSpans(t *testing.T) {\n\ttests := []struct {\n\t\tname    string\n\t\tinput   string\n\t\twantErr bool\n\t\twantLen int\n\t}{\n\t\t{\"valid\", `{\"spans\":[{\"kind\":\"header\",\"text\":\"x\"},{\"kind\":\"unit\",\"text\":\"f()\",\"action\":\"setFunction\",\"kwargs\":{\"name\":\"f\"}}]}`, false, 2},\n\t\t{\"unknown field\", `{\"spans\":[{\"kind\":\"header\",\"text\":\"x\"}],\"extra\":1}`, true, 0},\n\t\t{\"trailing data\", `{\"spans\":[{\"kind\":\"header\",\"text\":\"x\"}]} {}`, true, 0},\n\t\t{\"empty spans\", `{\"spans\":[]}`, true, 0},\n\t\t{\"missing spans\", `{}`, true, 0},\n\t\t{\"invalid json\", `{\"spans\":`, true, 0},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\tspans, err := ParseSpans(strings.NewReader(tt.input))\n\t\t\tif (err != nil) != tt.wantErr {\n\t\t\t\tt.Fatalf(\"ParseSpans error = %v, wantErr %v\", err, tt.wantErr)\n\t\t\t}\n\t\t\tif len(spans) != tt.wantLen {\n\t\t\t\tt.Fatalf(\"got %d spans, want %d\", len(spans), tt.wantLen)\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/adopt/command_test.go","name":"TestParseSpans"}}
func TestParseSpans(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantLen int
	}{
		{"valid", `{"spans":[{"kind":"header","text":"x"},{"kind":"unit","text":"f()","action":"setFunction","kwargs":{"name":"f"}}]}`, false, 2},
		{"unknown field", `{"spans":[{"kind":"header","text":"x"}],"extra":1}`, true, 0},
		{"trailing data", `{"spans":[{"kind":"header","text":"x"}]} {}`, true, 0},
		{"empty spans", `{"spans":[]}`, true, 0},
		{"missing spans", `{}`, true, 0},
		{"invalid json", `{"spans":`, true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spans, err := ParseSpans(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseSpans error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(spans) != tt.wantLen {
				t.Fatalf("got %d spans, want %d", len(spans), tt.wantLen)
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestCommandAdopterSplit(t *testing.T) {\n\tdir := t.TempDir()\n\tfile := filepath.Join(dir, \"input.txt\")\n\tif err := os.WriteFile(file, []byte(\"x\"), 0o644); err != nil {\n\t\tt.Fatal(err)\n\t}\n\ttests := []struct {\n\t\tname    string\n\t\tscript  string\n\t\twantErr string\n\t}{\n\t\t{\n\t\t\tname:   \"success\",\n\t\t\tscript: `printf '{\"spans\":[{\"kind\":\"residual\",\"text\":\"%s\"}]}' \"$1\"`,\n\t\t},\n\t\t{\n\t\t\tname:    \"non-zero exit\",\n\t\t\tscript:  `echo boom >&2; exit 3`,\n\t\t\twantErr: \"exit status 3: boom\",\n\t\t},\n\t\t{\n\t\t\tname:    \"malformed stdout\",\n\t\t\tscript:  `echo 'not json'`,\n\t\t\twantErr: \"decoding spans\",\n\t\t},\n\t}\n\tfor _, tt := range tests {\n\t\tt.Run(tt.name, func(t *testing.T) {\n\t\t\tscript := filepath.Join(dir, tt.name+\".sh\")\n\t\t\tif err := os.WriteFile(script, []byte(tt.script+\"\\n\"), 0o644); err != nil {\n\t\t\t\tt.Fatal(err)\n\t\t\t}\n\t\t\tc := CommandAdopter{Argv: []string{\"sh\", script}}\n\t\t\tspans, err := c.Split(file, nil)\n\t\t\tif tt.wantErr != \"\" {\n\t\t\t\tif err == nil || !strings.Contains(err.Error(), tt.wantErr) {\n\t\t\t\t\tt.Fatalf(\"Split error = %v, want it to contain %q\", err, tt.wantErr)\n\t\t\t\t}\n\t\t\t\treturn\n\t\t\t}\n\t\t\tif err != nil {\n\t\t\t\tt.Fatalf(\"Split: %v\", err)\n\t\t\t}\n\t\t\tif len(spans) != 1 || spans[0].Text != file {\n\t\t\t\tt.Fatalf(\"got spans %+v, want one span with the path %q\", spans, file)\n\t\t\t}\n\t\t})\n\t}\n\tif _, err := (CommandAdopter{}).Split(file, nil); err == nil {\n\t\tt.Fatal(\"empty Argv: want error\")\n\t}\n}","file":"internal/adopt/command_test.go","name":"TestCommandAdopterSplit"}}
func TestCommandAdopterSplit(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "input.txt")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		script  string
		wantErr string
	}{
		{
			name:   "success",
			script: `printf '{"spans":[{"kind":"residual","text":"%s"}]}' "$1"`,
		},
		{
			name:    "non-zero exit",
			script:  `echo boom >&2; exit 3`,
			wantErr: "exit status 3: boom",
		},
		{
			name:    "malformed stdout",
			script:  `echo 'not json'`,
			wantErr: "decoding spans",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := filepath.Join(dir, tt.name+".sh")
			if err := os.WriteFile(script, []byte(tt.script+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			c := CommandAdopter{Argv: []string{"sh", script}}
			spans, err := c.Split(file, nil)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Split error = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Split: %v", err)
			}
			if len(spans) != 1 || spans[0].Text != file {
				t.Fatalf("got spans %+v, want one span with the path %q", spans, file)
			}
		})
	}
	if _, err := (CommandAdopter{}).Split(file, nil); err == nil {
		t.Fatal("empty Argv: want error")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestCommandString(t *testing.T) {\n\ttests := []struct {\n\t\tname string\n\t\targv []string\n\t\twant string\n\t}{\n\t\t{\"single\", []string{\"adopt\"}, \"adopt\"},\n\t\t{\"multiple\", []string{\"sh\", \"script\", \"-v\"}, \"sh script -v\"},\n\t\t{\"empty\", nil, \"\"},\n\t}\n\tfor _, tt := range tests {\n\t\tc := CommandAdopter{Argv: tt.argv}\n\t\tif got := c.String(); got != tt.want {\n\t\t\tt.Errorf(\"%s: String() = %q, want %q\", tt.name, got, tt.want)\n\t\t}\n\t}\n}","file":"internal/adopt/command_test.go","name":"TestCommandString"}}
func TestCommandString(t *testing.T) {
	tests := []struct {
		name string
		argv []string
		want string
	}{
		{"single", []string{"adopt"}, "adopt"},
		{"multiple", []string{"sh", "script", "-v"}, "sh script -v"},
		{"empty", nil, ""},
	}
	for _, tt := range tests {
		c := CommandAdopter{Argv: tt.argv}
		if got := c.String(); got != tt.want {
			t.Errorf("%s: String() = %q, want %q", tt.name, got, tt.want)
		}
	}
}

// /sedum:setFunction
