package verify

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_test.go","path":"bytes"},"kwargs":"sha256:e9b024c4bddfbe869daae1c2dfac6534f9acd3a5ac08a2980336af4d4678d314","body":"sha256:468731f6dd513df0a66690d4649dc9e7244c96de227d8e6ead98cc5e9747c57a"}
import "bytes"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_test.go","path":"context"},"kwargs":"sha256:0d1bb9702c80d00f1bd8429074e07cf6e16521c6b7c492efa67a9163649796b3","body":"sha256:0f7443fb25d7ad6a970617d4d0df6123283e24801b94c1d84fbe431230f31e1e"}
import "context"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_test.go","path":"os"},"kwargs":"sha256:de12497a0f0733a80cea565477d30c2880efaf3b0f710a5212a35e2435acac3e","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_test.go","path":"path/filepath"},"kwargs":"sha256:1c9d67b939574c3db522c781af30478aac3cf59f258f5f50cc1a93c27965c6e9","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_test.go","path":"strings"},"kwargs":"sha256:db8d65d409989410014856dba785c122aca8a318fc3850ead1618831707aa151","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_test.go","path":"testing"},"kwargs":"sha256:401b0e8dd1ea231293aa2de19fe7dbec18df5e96a943d23568d0736dee02cbeb","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_test.go","path":"time"},"kwargs":"sha256:a9e06157e41549642b11d1cd243c9e9e1048b62cf447d49895c9861576fa30bd","body":"sha256:e99128ebf4e36998f40b33adf2eef57f4d6c315f1a4d1bdd3b7fc09a0d7bc82d"}
import "time"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_test.go","name":"tree"},"kwargs":"sha256:ab048bc70ad55ba428f526ca24f45fd1dd211f1e4a6ead1edb54c972201a47f4","body":"sha256:7979762f661d2f772df95d35ce2f1efbbf1679ee06c448056df12f1d1af559b6"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_test.go","name":"TestLoadConfigReadsCommandsAndTimeout"},"kwargs":"sha256:5873d5cd5b7d19a5c36a9e997c08548d9349b12f3b772c8586d70d4f789ff6e6","body":"sha256:8be66d1dd68ea9cb02e8b69bf77bfb027a256a7bd2e7be96989f81f721162f99"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_test.go","name":"TestLoadConfigWithoutAFileVerifiesNothing"},"kwargs":"sha256:2bd9311fa307592535aac5968a1b9b479c5380bb53f5fa8c3b6c05e92ac506f1","body":"sha256:2414bd3a3d97362f712e85cce8436f1b12a458891da8b7894f782b90ca55ab79"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_test.go","name":"TestLoadConfigRejectsABadTimeout"},"kwargs":"sha256:aa1484fc2463a213786d2f1a20fcbd1c8f645f5400085d1f340e1db89ba74570","body":"sha256:f855d5208a9a54eb80560feffe2a4674a4fc28bf40fe85004d0f14084fcbe978"}
func TestLoadConfigRejectsABadTimeout(t *testing.T) {
	root := tree(t, map[string]string{".sedum/config.yml": "verify: [make]\nverify_timeout: soon\n"})
	if _, err := LoadConfig(root); err == nil || !strings.Contains(err.Error(), "verify_timeout") {
		t.Errorf("LoadConfig error = %v, want one naming verify_timeout", err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_test.go","name":"TestRunSeesTheCandidateAndLeavesTheOutputAlone"},"kwargs":"sha256:b76456f283155cf617ed96a56d1d4f246aa4ddfda0fae839f2355f86d4a65fc5","body":"sha256:16a6b594facf403124abd38599dc81d98e1c6a9acaba8cf81d2b5667724ef670"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_test.go","name":"TestRunReportsAFailureWithItsStatusAndOutput"},"kwargs":"sha256:2c30cd70fa9a169e13a9c62be10e2e537355bf5b991f7cdd8d5d1e0de7fbd42e","body":"sha256:4fd8c8ac6663a65347fd27f5db0d4d87bc3c0cb44b8c1d2e40666bda37bb8e5b"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_test.go","name":"TestRunKeepsOnlyTheEndOfALongOutput"},"kwargs":"sha256:9af519bbb6e97eccc6d0d0b85eea4941f72a9c843b7fc5acee42d4ccc4a35873","body":"sha256:572da443ea20e3a2b4a77a48da20f219ef453fa477ae2b86dababf6d6c76705d"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_test.go","name":"TestRunSkipsHiddenAndNamedDirectories"},"kwargs":"sha256:1eda1ef9b10150d51441dfea3b1ffc681a1c1eab4ac0c93875a8a6ad4215c28c","body":"sha256:d81a432a9b347cccddd09690e8398e591f40d37e1024b97e414557dd387bfecd"}
func TestRunSkipsHiddenAndNamedDirectories(t *testing.T) {
	root := tree(t, map[string]string{".git/HEAD": "x", "generators/a.txt": "x", "src/a.txt": "x"})
	cfg := Config{Commands: []string{"test ! -e .git && test ! -e generators && test -f src/a.txt"}, Timeout: 10 * time.Second, Skip: []string{"generators"}}
	failures, err := Run(context.Background(), cfg, root, nil, &bytes.Buffer{})
	if err != nil || len(failures) != 0 {
		t.Errorf("failures = %+v, err = %v", failures, err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-6ab86198","gen":"sha256:aaf13f1fa8f54e47d5b1b268a14d9f7ab32f03b73759ac545746dcf9b58d05e3","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_test.go","name":"TestRunKillsTheWholeProcessGroupOnTimeoutAndCleansUp"},"kwargs":"sha256:4512d7e26e7e390b14ff3ac57fca958f3eecf12c5e13cd41103868750118cf08","body":"sha256:ab510181b38c4dd01719200e0e157b7033bc42d42c4977a238d6011a1d9f8841"}
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
