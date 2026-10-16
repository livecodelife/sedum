package verify

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify_test.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify_test.go","path":"context"}}
import "context"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify_test.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify_test.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify_test.go","path":"time"}}
import "time"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// tree writes files (slash paths under a fresh directory) and returns it.\nfunc tree(t *testing.T, files map[string]string) string {\n\tt.Helper()\n\troot := t.TempDir()\n\tfor rel, content := range files {\n\t\tpath := filepath.Join(root, filepath.FromSlash(rel))\n\t\tif err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t\tif err := os.WriteFile(path, []byte(content), 0o644); err != nil {\n\t\t\tt.Fatal(err)\n\t\t}\n\t}\n\treturn root\n}","file":"internal/verify/verify_test.go","name":"tree"}}
// tree writes files (slash paths under a fresh directory) and returns it.
func tree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestLoadConfigReadsCommandsAndTimeout(t *testing.T) {\n\troot := tree(t, map[string]string{\n\t\t\".sedum/config.yml\": \"verify:\\n  - go build ./...\\n  - go vet ./...\\nverify_timeout: 30s\\n\",\n\t})\n\tcfg, err := LoadConfig(root)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(cfg.Commands) != 2 || cfg.Commands[0] != \"go build ./...\" || cfg.Timeout != 30*time.Second {\n\t\tt.Errorf(\"LoadConfig = %+v\", cfg)\n\t}\n\tif !cfg.Enabled() {\n\t\tt.Error(\"a config with commands should be enabled\")\n\t}\n}","file":"internal/verify/verify_test.go","name":"TestLoadConfigReadsCommandsAndTimeout"}}
func TestLoadConfigReadsCommandsAndTimeout(t *testing.T) {
	root := tree(t, map[string]string{
		".sedum/config.yml": "verify:\n  - go build ./...\n  - go vet ./...\nverify_timeout: 30s\n",
	})
	cfg, err := LoadConfig(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Commands) != 2 || cfg.Commands[0] != "go build ./..." || cfg.Timeout != 30*time.Second {
		t.Errorf("LoadConfig = %+v", cfg)
	}
	if !cfg.Enabled() {
		t.Error("a config with commands should be enabled")
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestLoadConfigWithoutAFileVerifiesNothing(t *testing.T) {\n\tcfg, err := LoadConfig(t.TempDir())\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif cfg.Enabled() || cfg.Timeout != DefaultTimeout {\n\t\tt.Errorf(\"LoadConfig = %+v, want nothing to run and the default timeout\", cfg)\n\t}\n\tif got := cfg.With([]string{\"make test\"}); len(got.Commands) != 1 || cfg.Enabled() {\n\t\tt.Errorf(\"With = %+v, and the original must not change: %+v\", got, cfg)\n\t}\n}","file":"internal/verify/verify_test.go","name":"TestLoadConfigWithoutAFileVerifiesNothing"}}
func TestLoadConfigWithoutAFileVerifiesNothing(t *testing.T) {
	cfg, err := LoadConfig(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Enabled() || cfg.Timeout != DefaultTimeout {
		t.Errorf("LoadConfig = %+v, want nothing to run and the default timeout", cfg)
	}
	if got := cfg.With([]string{"make test"}); len(got.Commands) != 1 || cfg.Enabled() {
		t.Errorf("With = %+v, and the original must not change: %+v", got, cfg)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestLoadConfigRejectsABadTimeout(t *testing.T) {\n\troot := tree(t, map[string]string{\".sedum/config.yml\": \"verify: [make]\\nverify_timeout: soon\\n\"})\n\tif _, err := LoadConfig(root); err == nil || !strings.Contains(err.Error(), \"verify_timeout\") {\n\t\tt.Errorf(\"LoadConfig error = %v, want one naming verify_timeout\", err)\n\t}\n}","file":"internal/verify/verify_test.go","name":"TestLoadConfigRejectsABadTimeout"}}
func TestLoadConfigRejectsABadTimeout(t *testing.T) {
	root := tree(t, map[string]string{".sedum/config.yml": "verify: [make]\nverify_timeout: soon\n"})
	if _, err := LoadConfig(root); err == nil || !strings.Contains(err.Error(), "verify_timeout") {
		t.Errorf("LoadConfig error = %v, want one naming verify_timeout", err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunSeesTheCandidateAndLeavesTheOutputAlone(t *testing.T) {\n\troot := tree(t, map[string]string{\"keep.txt\": \"kept\\n\"})\n\tvar log bytes.Buffer\n\tcfg := Config{Commands: []string{\"grep -q hello new/file.txt && test -f keep.txt\"}, Timeout: 10 * time.Second}\n\tfailures, err := Run(context.Background(), cfg, root, map[string]string{\"new/file.txt\": \"hello\\n\"}, &log)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(failures) != 0 {\n\t\tt.Errorf(\"failures = %+v, want none\", failures)\n\t}\n\tif _, err := os.Stat(filepath.Join(root, \"new\")); err == nil {\n\t\tt.Error(\"the candidate was written into the output directory\")\n\t}\n\tif !strings.Contains(log.String(), \"verify: grep -q hello\") {\n\t\tt.Errorf(\"the command was not printed before it ran: %q\", log.String())\n\t}\n}","file":"internal/verify/verify_test.go","name":"TestRunSeesTheCandidateAndLeavesTheOutputAlone"}}
func TestRunSeesTheCandidateAndLeavesTheOutputAlone(t *testing.T) {
	root := tree(t, map[string]string{"keep.txt": "kept\n"})
	var log bytes.Buffer
	cfg := Config{Commands: []string{"grep -q hello new/file.txt && test -f keep.txt"}, Timeout: 10 * time.Second}
	failures, err := Run(context.Background(), cfg, root, map[string]string{"new/file.txt": "hello\n"}, &log)
	if err != nil {
		t.Fatal(err)
	}
	if len(failures) != 0 {
		t.Errorf("failures = %+v, want none", failures)
	}
	if _, err := os.Stat(filepath.Join(root, "new")); err == nil {
		t.Error("the candidate was written into the output directory")
	}
	if !strings.Contains(log.String(), "verify: grep -q hello") {
		t.Errorf("the command was not printed before it ran: %q", log.String())
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunReportsAFailureWithItsStatusAndOutput(t *testing.T) {\n\troot := tree(t, nil)\n\tcfg := Config{Commands: []string{\"echo boom; exit 3\", \"true\"}, Timeout: 10 * time.Second}\n\tfailures, err := Run(context.Background(), cfg, root, nil, &bytes.Buffer{})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(failures) != 1 || failures[0].Status != \"exit status 3\" || !strings.Contains(failures[0].Output, \"boom\") {\n\t\tt.Fatalf(\"failures = %+v\", failures)\n\t}\n\tmessage := Message(failures)\n\tfor _, want := range []string{\"These checks failed\", \"$ echo boom; exit 3\", \"exit status 3\", \"boom\"} {\n\t\tif !strings.Contains(message, want) {\n\t\t\tt.Errorf(\"Message lacks %q:\\n%s\", want, message)\n\t\t}\n\t}\n}","file":"internal/verify/verify_test.go","name":"TestRunReportsAFailureWithItsStatusAndOutput"}}
func TestRunReportsAFailureWithItsStatusAndOutput(t *testing.T) {
	root := tree(t, nil)
	cfg := Config{Commands: []string{"echo boom; exit 3", "true"}, Timeout: 10 * time.Second}
	failures, err := Run(context.Background(), cfg, root, nil, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if len(failures) != 1 || failures[0].Status != "exit status 3" || !strings.Contains(failures[0].Output, "boom") {
		t.Fatalf("failures = %+v", failures)
	}
	message := Message(failures)
	for _, want := range []string{"These checks failed", "$ echo boom; exit 3", "exit status 3", "boom"} {
		if !strings.Contains(message, want) {
			t.Errorf("Message lacks %q:\n%s", want, message)
		}
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunKeepsOnlyTheEndOfALongOutput(t *testing.T) {\n\tcfg := Config{Commands: []string{\"yes xxxxxxxxx | head -n 2000; echo the-end; exit 1\"}, Timeout: 10 * time.Second}\n\tfailures, err := Run(context.Background(), cfg, tree(t, nil), nil, &bytes.Buffer{})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif len(failures) != 1 || len(failures[0].Output) > MaxOutput || !strings.HasSuffix(failures[0].Output, \"the-end\\n\") {\n\t\tt.Errorf(\"output length %d, want at most %d ending with the-end\", len(failures[0].Output), MaxOutput)\n\t}\n}","file":"internal/verify/verify_test.go","name":"TestRunKeepsOnlyTheEndOfALongOutput"}}
func TestRunKeepsOnlyTheEndOfALongOutput(t *testing.T) {
	cfg := Config{Commands: []string{"yes xxxxxxxxx | head -n 2000; echo the-end; exit 1"}, Timeout: 10 * time.Second}
	failures, err := Run(context.Background(), cfg, tree(t, nil), nil, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if len(failures) != 1 || len(failures[0].Output) > MaxOutput || !strings.HasSuffix(failures[0].Output, "the-end\n") {
		t.Errorf("output length %d, want at most %d ending with the-end", len(failures[0].Output), MaxOutput)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestRunSkipsHiddenAndNamedDirectories(t *testing.T) {\n\troot := tree(t, map[string]string{\".git/HEAD\": \"x\", \"generators/a.txt\": \"x\", \"src/a.txt\": \"x\"})\n\tcfg := Config{Commands: []string{\"test ! -e .git && test ! -e generators && test -f src/a.txt\"}, Timeout: 10 * time.Second, Skip: []string{\"generators\"}}\n\tfailures, err := Run(context.Background(), cfg, root, nil, &bytes.Buffer{})\n\tif err != nil || len(failures) != 0 {\n\t\tt.Errorf(\"failures = %+v, err = %v\", failures, err)\n\t}\n}","file":"internal/verify/verify_test.go","name":"TestRunSkipsHiddenAndNamedDirectories"}}
func TestRunSkipsHiddenAndNamedDirectories(t *testing.T) {
	root := tree(t, map[string]string{".git/HEAD": "x", "generators/a.txt": "x", "src/a.txt": "x"})
	cfg := Config{Commands: []string{"test ! -e .git && test ! -e generators && test -f src/a.txt"}, Timeout: 10 * time.Second, Skip: []string{"generators"}}
	failures, err := Run(context.Background(), cfg, root, nil, &bytes.Buffer{})
	if err != nil || len(failures) != 0 {
		t.Errorf("failures = %+v, err = %v", failures, err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// A timeout kills the command and everything it started, reports it as a\n// failure, and still removes the scratch directory.\nfunc TestRunKillsTheWholeProcessGroupOnTimeoutAndCleansUp(t *testing.T) {\n\tmarker := filepath.Join(t.TempDir(), \"cwd\")\n\tcfg := Config{Commands: []string{\"pwd > \" + marker + \"; sleep 30 & sleep 30\"}, Timeout: 300 * time.Millisecond}\n\tstart := time.Now()\n\tfailures, err := Run(context.Background(), cfg, tree(t, nil), nil, &bytes.Buffer{})\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif elapsed := time.Since(start); elapsed > 10*time.Second {\n\t\tt.Errorf(\"Run took %s: the background process held the run open\", elapsed)\n\t}\n\tif len(failures) != 1 || !strings.HasPrefix(failures[0].Status, \"timed out\") {\n\t\tt.Fatalf(\"failures = %+v, want one timeout\", failures)\n\t}\n\tdata, err := os.ReadFile(marker)\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif _, err := os.Stat(strings.TrimSpace(string(data))); err == nil {\n\t\tt.Errorf(\"the scratch directory %s was not removed\", strings.TrimSpace(string(data)))\n\t}\n}","file":"internal/verify/verify_test.go","name":"TestRunKillsTheWholeProcessGroupOnTimeoutAndCleansUp"}}
// A timeout kills the command and everything it started, reports it as a
// failure, and still removes the scratch directory.
func TestRunKillsTheWholeProcessGroupOnTimeoutAndCleansUp(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "cwd")
	cfg := Config{Commands: []string{"pwd > " + marker + "; sleep 30 & sleep 30"}, Timeout: 300 * time.Millisecond}
	start := time.Now()
	failures, err := Run(context.Background(), cfg, tree(t, nil), nil, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Errorf("Run took %s: the background process held the run open", elapsed)
	}
	if len(failures) != 1 || !strings.HasPrefix(failures[0].Status, "timed out") {
		t.Fatalf("failures = %+v, want one timeout", failures)
	}
	data, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(strings.TrimSpace(string(data))); err == nil {
		t.Errorf("the scratch directory %s was not removed", strings.TrimSpace(string(data)))
	}
}

// /sedum:setFunction
