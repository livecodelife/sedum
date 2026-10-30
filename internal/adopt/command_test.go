package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:450c20a9cdeb291d35f0f29c47d40a5b6c6a10105ecdcf8fb4a7e01872e8c90f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/command_test.go","path":"os"},"kwargs":"sha256:64e56ebad073686d7a0049d374e272283ec7461b09361cc8f71b5065227e3218","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:450c20a9cdeb291d35f0f29c47d40a5b6c6a10105ecdcf8fb4a7e01872e8c90f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/command_test.go","path":"path/filepath"},"kwargs":"sha256:1b07d7f6d58d5ca44979aa2ef6a5b165ce8c393d64eab8532aa9ae66f3662996","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:450c20a9cdeb291d35f0f29c47d40a5b6c6a10105ecdcf8fb4a7e01872e8c90f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/command_test.go","path":"strings"},"kwargs":"sha256:d5b18283099f79e2b38109e0c43a6936b9cd713b1bf79ae91abb7a3c70892cf9","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:450c20a9cdeb291d35f0f29c47d40a5b6c6a10105ecdcf8fb4a7e01872e8c90f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/command_test.go","path":"testing"},"kwargs":"sha256:eeccbe819056c054b52f3dac78d0ae367de6b4e2343c79dccda69553285ef9a9","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:450c20a9cdeb291d35f0f29c47d40a5b6c6a10105ecdcf8fb4a7e01872e8c90f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/command_test.go","name":"TestParseSpans"},"kwargs":"sha256:c70bc485fc05cf8ecec1251ce1fbd42fe2c8c31d98107a91322875fdfb81a03a","body":"sha256:8215149f16f9a468f502e88d7d2bb0e6f23bed9c3e4a714f3aac4e28f389b426"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:450c20a9cdeb291d35f0f29c47d40a5b6c6a10105ecdcf8fb4a7e01872e8c90f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/command_test.go","name":"TestCommandAdopterSplit"},"kwargs":"sha256:4e86c9058c495369fcf1affd43f2c6784ab7855e9ebc098604b02468fd4c0c39","body":"sha256:1c5ad0d0112a9921ade193726361e84e47fd38dc21d152c8f4d8622847797a41"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:450c20a9cdeb291d35f0f29c47d40a5b6c6a10105ecdcf8fb4a7e01872e8c90f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/command_test.go","name":"TestCommandString"},"kwargs":"sha256:f1e0f4cec7389c9fd43f74e61f940cceabb150a76c1ba4dfc2699adf5cc0656c","body":"sha256:cc8066d4e5b9b0573eb76728149132d3f326925ee59b52396b28a2f5b372c4e6"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:450c20a9cdeb291d35f0f29c47d40a5b6c6a10105ecdcf8fb4a7e01872e8c90f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/command_test.go","name":"TestCommandAdopterRunsInDir"},"kwargs":"sha256:dc449f929290d2a74b150c91570c8a90714a66055cc7eb2d8f1da56713e0a32b","body":"sha256:3ba5e8a0ea7b0c9ee34b2a58cd5b55e912b7f8d75786197c0617cc3eb3b2722c"}
func TestCommandAdopterRunsInDir(t *testing.T) {
	work := t.TempDir()
	if err := os.WriteFile(filepath.Join(work, "x.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	scripts := t.TempDir()
	script := filepath.Join(scripts, "pwd.sh")
	body := `printf '{"spans":[{"kind":"residual","text":"%s"}]}' "$(pwd)"` + "\n"
	if err := os.WriteFile(script, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	c := CommandAdopter{Argv: []string{"sh", script}, Dir: work}
	spans, err := c.Split("x.txt", nil)
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	if len(spans) != 1 || !strings.HasSuffix(spans[0].Text, filepath.Base(work)) {
		t.Fatalf("got spans %+v, want text ending in %q", spans, filepath.Base(work))
	}
}

// /sedum:setFunction
