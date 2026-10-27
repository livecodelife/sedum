package cli

import (
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// `sedum check` registers the command, validates its flags, reads the config and
// maps the result of the check to the process exit code.

func ckEnv(t *testing.T) []string {
	t.Helper()
	home := t.TempDir()
	var env []string
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "GIT_") && !strings.HasPrefix(kv, "HOME=") {
			env = append(env, kv)
		}
	}
	return append(env, "HOME="+home, "GIT_CONFIG_GLOBAL="+filepath.Join(home, "gitconfig"), "GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com")
}

func ckGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := osexec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = ckEnv(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// ckRepo is a repository with two commits and a fake provenance on PROVENANCE
// that prints the given findings (a JSON array of results) and exits 1.
func ckRepo(t *testing.T, results string) string {
	t.Helper()
	dir := t.TempDir()
	env := ckEnv(t)
	for _, kv := range env {
		if k, v, ok := strings.Cut(kv, "="); ok {
			t.Setenv(k, v)
		}
	}
	ckGit(t, dir, "init", "-q", "-b", "main")
	for i, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(name), 0o644); err != nil {
			t.Fatal(err)
		}
		ckGit(t, dir, "add", "-A")
		ckGit(t, dir, "commit", "-q", "-m", "c"+string(rune('0'+i)))
	}
	bin := filepath.Join(t.TempDir(), "provenance")
	script := "#!/bin/sh\ncat <<'EOF'\n" + `{"runs":[{"tool":{"driver":{"rules":[]}},"results":` + results + `}]}` + "\nEOF\nexit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PROVENANCE", bin)
	return dir
}

func wantExit(t *testing.T, err error, code int, prefix string) {
	t.Helper()
	got, ok := CodeOf(err)
	if !ok || got != code || !strings.HasPrefix(err.Error(), prefix) {
		t.Errorf("err = %v (code %d, ok %v), want code %d starting %q", err, got, ok, code, prefix)
	}
}

func TestCheckUsageErrorsExitTwo(t *testing.T) {
	dir := ckRepo(t, "[]")
	cases := map[string][]string{
		"a positional argument":     {"check", "--range", "a..b", "extra"},
		"neither range nor staged":  {"check", "--output", dir},
		"both range and staged":     {"check", "--range", "a..b", "--staged", "--output", dir},
		"message file with a range": {"check", "--range", "HEAD~1..HEAD", "--message-file", "m", "--output", dir},
		"an unknown format":         {"check", "--range", "HEAD~1..HEAD", "--format", "xml", "--output", dir},
		"an unknown flag":           {"check", "--nope"},
	}
	for name, args := range cases {
		_, err := exec(t, args...)
		if err == nil {
			t.Errorf("%s: no error", name)
			continue
		}
		wantExit(t, err, 2, "usage:")
	}
}

func TestCheckReturnsNilWhenNothingIsLeft(t *testing.T) {
	dir := ckRepo(t, "[]")

	out, err := exec(t, "check", "--range", "HEAD~1..HEAD", "--output", dir)

	if err != nil {
		t.Fatalf("check: %v\n%s", err, out)
	}
	if out != "skipped 0 scope paths that no generator package claims\nwaived: 0\n" {
		t.Errorf("output = %q", out)
	}
}

func TestCheckReturnsASilentExitOneWhenAFindingIsLeft(t *testing.T) {
	dir := ckRepo(t, `[{"ruleId":"TAG001","message":{"text":"Commit tag required"}}]`)

	out, err := exec(t, "check", "--range", "HEAD~1..HEAD", "--output", dir)

	exit, ok := err.(ExitError)
	if !ok || exit.Code != 1 || exit.Message != "" {
		t.Fatalf("err = %#v, want ExitError{Code: 1}", err)
	}
	if !strings.HasPrefix(out, "TAG001: Commit tag required\n") {
		t.Errorf("output = %q", out)
	}
}

func TestCheckWritesSARIFWithFormatSarif(t *testing.T) {
	dir := ckRepo(t, "[]")

	out, err := exec(t, "check", "--range", "HEAD~1..HEAD", "--output", dir, "--format", "sarif")

	if err != nil || !strings.Contains(out, `"name": "sedum check"`) {
		t.Errorf("err = %v output = %q", err, out)
	}
}

func TestCheckExitsTwoWhenProvenanceCannotRun(t *testing.T) {
	dir := ckRepo(t, "[]")
	t.Setenv("PROVENANCE", filepath.Join(t.TempDir(), "missing"))

	_, err := exec(t, "check", "--range", "HEAD~1..HEAD", "--output", dir)

	wantExit(t, err, 2, "provenance:")
}

func TestCheckExitsTwoOutsideAGitRepository(t *testing.T) {
	ckRepo(t, "[]")

	_, err := exec(t, "check", "--range", "HEAD~1..HEAD", "--output", t.TempDir())

	wantExit(t, err, 2, "git:")
}

func TestCheckRefusesAnOutputInAnotherRepositoryThanTheConfig(t *testing.T) {
	dir := ckRepo(t, "[]")
	other := t.TempDir()
	ckGit(t, other, "init", "-q")
	cfg := filepath.Join(other, ".provenance.yml")
	if err := os.WriteFile(cfg, []byte("dir: provenance\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := exec(t, "check", "--range", "HEAD~1..HEAD", "--output", dir, "--config", cfg)

	wantExit(t, err, 2, "usage:")
}

func TestReadConfig(t *testing.T) {
	cases := []struct {
		name     string
		files    map[string]string
		flag     string
		wantDir  string
		wantEx   []string
		wantFail bool
	}{
		{"no config", nil, "", "provenance", nil, false},
		{"a .provenance.yml", map[string]string{".provenance.yml": "dir: records/\nexclude_paths: [a, 'b/**']\n"}, "", "records", []string{"a", "b/**"}, false},
		{"slashes trimmed and an empty dir", map[string]string{".provenance.yml": "dir: \"\"\n"}, "", "provenance", nil, false},
		{"a .linespec.yml block", map[string]string{".linespec.yml": "other: 1\nprovenance:\n  dir: prov\n  exclude_paths:\n    - x\n"}, "", "prov", []string{"x"}, false},
		{".provenance.yml wins over .linespec.yml", map[string]string{".provenance.yml": "dir: one\n", ".linespec.yml": "provenance:\n  dir: two\n"}, "", "one", nil, false},
		{"a .linespec.yml without a block", map[string]string{".linespec.yml": "other: 1\n"}, "", "provenance", nil, false},
		{"--config naming a .linespec.yml", map[string]string{"conf/.linespec.yml": "provenance:\n  dir: viaflag\n"}, "conf/.linespec.yml", "viaflag", nil, false},
		{"--config naming another file is read from its top level", map[string]string{"conf/my.yml": "dir: top\nexclude_paths: [z]\n"}, "conf/my.yml", "top", []string{"z"}, false},
		{"unparseable", map[string]string{".provenance.yml": "dir: [unclosed\n"}, "", "", nil, true},
	}
	for _, c := range cases {
		base := t.TempDir()
		for rel, content := range c.files {
			full := filepath.Join(base, filepath.FromSlash(rel))
			if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		flag := ""
		if c.flag != "" {
			flag = filepath.Join(base, filepath.FromSlash(c.flag))
		}
		root := base
		if flag != "" {
			root = filepath.Dir(flag)
		}

		dir, excludes, err := readConfig(root, flag)

		if c.wantFail {
			if err == nil || !strings.HasPrefix(err.Error(), "provenance:") {
				t.Errorf("%s: err = %v, want a provenance: error", c.name, err)
			}
			continue
		}
		if err != nil || dir != c.wantDir || strings.Join(excludes, ",") != strings.Join(c.wantEx, ",") {
			t.Errorf("%s: dir %q excludes %v err %v; want %q %v", c.name, dir, excludes, err, c.wantDir, c.wantEx)
		}
	}
}

func TestCheckCommandHasTheDocumentedFlags(t *testing.T) {
	cmd, _, err := NewRootCommand().Find([]string{"check"})
	if err != nil || cmd.Name() != "check" {
		t.Fatalf("check command = %v, %v", cmd, err)
	}
	defaults := map[string]string{"output": ".", "format": "human", "range": "", "message-file": "", "config": "", "generators": "", "records": "", "provenance": "", "staged": "false"}
	for name, def := range defaults {
		f := cmd.Flags().Lookup(name)
		if f == nil || f.DefValue != def {
			t.Errorf("flag --%s = %+v, want default %q", name, f, def)
		}
	}
	for _, text := range []string{"provenance check", "exit 2", "skew", "classes"} {
		if !strings.Contains(cmd.Long, text) {
			t.Errorf("help text lacks %q", text)
		}
	}
}
