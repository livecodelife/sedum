package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// testTemp is the root every fixture repository must live under.
var testTemp string

// TestMain isolates git for the whole test process. Run from a git hook
// (linespec's pre-commit re-runs specs), git exports GIT_DIR and
// GIT_INDEX_FILE, and git commands in a fixture would act on the real
// repository instead: once that rewrote the shared .git/config and committed
// a fixture to the real branch. So every GIT_* variable is cleared, HOME and
// the global config point into a temp dir that holds the test identity and an
// empty hooks dir (fixtures never run `git config`), system config is off,
// and discovery cannot climb out of the temp root. The CLI itself keeps
// honouring GIT_* in real use, since a hook-invoked command reads the index
// being committed.
func TestMain(m *testing.M) {
	for _, kv := range os.Environ() {
		if name, _, _ := strings.Cut(kv, "="); strings.HasPrefix(name, "GIT_") {
			os.Unsetenv(name)
		}
	}
	tmp, err := os.MkdirTemp("", "provenance-cli-test-")
	if err != nil {
		panic(err)
	}
	if testTemp, err = filepath.EvalSymlinks(tmp); err != nil {
		panic(err)
	}
	home := filepath.Join(testTemp, "home")
	hooks := filepath.Join(home, "no-hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		panic(err)
	}
	gitconfig := filepath.Join(home, ".gitconfig")
	body := "[user]\n\temail = test@example.com\n\tname = Test\n[core]\n\thooksPath = " + hooks + "\n[init]\n\tdefaultBranch = main\n"
	if err := os.WriteFile(gitconfig, []byte(body), 0o644); err != nil {
		panic(err)
	}
	os.Setenv("HOME", home)
	os.Setenv("XDG_CONFIG_HOME", home)
	os.Setenv("GIT_CONFIG_GLOBAL", gitconfig)
	os.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	os.Setenv("GIT_CEILING_DIRECTORIES", testTemp)
	os.Setenv("TMPDIR", testTemp)
	code := m.Run()
	os.RemoveAll(testTemp)
	os.Exit(code)
}

// run executes the provenance command tree with args and returns its stdout.
func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := newRootCmd()
	var out, stderr bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&stderr)
	root.SetArgs(args)
	err := root.Execute()
	if err != nil && stderr.Len() > 0 {
		t.Logf("provenance %s stderr:\n%s", strings.Join(args, " "), stderr.String())
	}
	return out.String(), err
}

func mustRun(t *testing.T, args ...string) string {
	t.Helper()
	out, err := run(t, args...)
	if err != nil {
		t.Fatalf("provenance %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return out
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// gitIn runs git on the fixture repository at dir, which must live under
// the test temp root. GIT_DIR and GIT_WORK_TREE name it explicitly, so no
// inherited state can redirect the command to another repository.
func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	real, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(real+string(filepath.Separator), testTemp+string(filepath.Separator)) {
		t.Fatalf("refusing git %s outside the test temp root: %s", strings.Join(args, " "), dir)
	}
	c := exec.Command("git", append([]string{"-C", real}, args...)...)
	c.Dir = real
	c.Env = append(os.Environ(), "GIT_DIR="+filepath.Join(real, ".git"), "GIT_WORK_TREE="+real)
	return runGit(t, c)
}

// gitClone clones src into dst, which must live under the test temp root.
func gitClone(t *testing.T, src, dst string) {
	t.Helper()
	parent, err := filepath.EvalSymlinks(filepath.Dir(dst))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(parent+string(filepath.Separator), testTemp+string(filepath.Separator)) {
		t.Fatalf("refusing to clone outside the test temp root: %s", dst)
	}
	c := exec.Command("git", "-C", parent, "clone", "-q", src, filepath.Join(parent, filepath.Base(dst)))
	c.Dir = parent
	runGit(t, c)
}

// gitStatus reads `git status --porcelain` of a real repository. It is the
// only git command tests run outside the temp root, and it writes nothing.
func gitStatus(t *testing.T, repo string, paths ...string) string {
	t.Helper()
	c := exec.Command("git", append([]string{"-C", repo, "--no-optional-locks", "status", "--porcelain", "--"}, paths...)...)
	c.Dir = repo
	return runGit(t, c)
}

func runGit(t *testing.T, c *exec.Cmd) string {
	t.Helper()
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v\n%s", strings.Join(c.Args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

// gitRepo makes dir a git repository and commits everything in it. Identity
// and the empty hooks dir come from TestMain's global config.
func gitRepo(t *testing.T, dir string) {
	t.Helper()
	gitIn(t, dir, "init", "-q")
	gitIn(t, dir, "add", "-A")
	gitIn(t, dir, "commit", "-q", "-m", "fixture")
}

// v1 is a v1 record file, as linespec writes it.
func v1(id, typ, status, extra string) string {
	return "id: " + id + "\ntitle: T " + id + "\nstatus: " + status + "\ncreated_at: \"2026-01-01\"\nauthor: a@b\n" +
		"intent: Intent of " + id + "\nconstraints:\n  - Shared\naffected_scope:\n  - a.go\nforbidden_scope: []\ntype: " + typ + "\n" +
		"supersedes: \"\"\nsuperseded_by: \"\"\nrelated: []\nsealed_at_sha: \"\"\nassociated_specs: []\nassociated_traces: []\nmonitors: []\ntags: []\n" + extra
}

// TestGitEnvFromHook: a hook exports GIT_DIR and GIT_INDEX_FILE for the
// repository being committed. Acting on another repository, the CLI's git
// commands must ignore them; acting on the hooked repository itself, they
// must keep them.
func TestGitEnvFromHook(t *testing.T) {
	hooked, other := t.TempDir(), t.TempDir()
	for _, d := range []string{hooked, other} {
		writeFile(t, filepath.Join(d, "f"), d+"\n")
		gitRepo(t, d)
	}
	hookedHead := gitIn(t, hooked, "rev-parse", "HEAD")
	t.Setenv("GIT_DIR", filepath.Join(hooked, ".git"))
	t.Setenv("GIT_INDEX_FILE", filepath.Join(hooked, ".git", "index"))

	r := &repo{Root: other}
	writeFile(t, filepath.Join(other, "g"), "new\n")
	if err := r.commit("other commit", "g"); err != nil {
		t.Fatal(err)
	}
	// Read both repos with an environment free of the hook's variables.
	os.Unsetenv("GIT_DIR")
	os.Unsetenv("GIT_INDEX_FILE")
	if got := gitIn(t, other, "log", "-1", "--format=%s"); got != "other commit" {
		t.Fatalf("commit did not land in the repo acted on: %q", got)
	}
	if got := gitIn(t, hooked, "rev-parse", "HEAD"); got != hookedHead {
		t.Fatalf("hooked repo moved: %s -> %s", hookedHead, got)
	}

	t.Setenv("GIT_DIR", filepath.Join(hooked, ".git"))
	t.Setenv("GIT_INDEX_FILE", filepath.Join(hooked, ".git", "next-index"))
	env := (&repo{Root: hooked}).gitEnv()
	if !slices.Contains(env, "GIT_INDEX_FILE="+filepath.Join(hooked, ".git", "next-index")) {
		t.Fatal("acting on the hooked repo dropped its GIT_INDEX_FILE")
	}
}
