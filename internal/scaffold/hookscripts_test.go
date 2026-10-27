package scaffold

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The hook scripts unit holds the text of the commit-msg and pre-push hooks and
// the list that names them. The scripts are run for real, in throwaway git
// repositories, with fake sedum and provenance binaries that record how they
// were called.

func TestHooksListsCommitMsgThenPrePush(t *testing.T) {
	hooks := Hooks()
	if len(hooks) != 2 || hooks[0].Name != "commit-msg" || hooks[1].Name != "pre-push" {
		t.Fatalf("Hooks = %+v", hooks)
	}
	for _, h := range hooks {
		lines := strings.Split(h.Body, "\n")
		want := "# sedum " + h.Name + " hook, installed by `sedum hooks install`."
		if lines[0] != "#!/bin/sh" || lines[1] != want || !strings.HasSuffix(h.Body, "\n") || strings.HasSuffix(h.Body, "\n\n") {
			t.Errorf("%s: first lines %q, text must end with exactly one newline", h.Name, lines[:2])
		}
	}
	if &Hooks()[0] == &hooks[0] {
		t.Error("Hooks returned the same slice twice")
	}
}

type hookEnv struct {
	t    *testing.T
	repo string
	bin  string
	log  string
	env  []string
}

func newHookEnv(t *testing.T) *hookEnv {
	t.Helper()
	h := &hookEnv{t: t, repo: t.TempDir(), bin: t.TempDir(), log: filepath.Join(t.TempDir(), "calls.log")}
	home := t.TempDir()
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "GIT_") && !strings.HasPrefix(kv, "HOME=") && !strings.HasPrefix(kv, "PATH=") && !strings.HasPrefix(kv, "SEDUM=") && !strings.HasPrefix(kv, "PROVENANCE=") {
			h.env = append(h.env, kv)
		}
	}
	h.env = append(h.env, "HOME="+home, "GIT_CONFIG_GLOBAL="+filepath.Join(home, "gitconfig"), "GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com",
		"CALLS="+h.log, "PATH="+h.bin+":/usr/bin:/bin")
	h.git("init", "-q", "-b", "main")
	return h
}

func (h *hookEnv) git(args ...string) string {
	h.t.Helper()
	cmd := exec.Command("git", append([]string{"-C", h.repo}, args...)...)
	cmd.Env = h.env
	out, err := cmd.CombinedOutput()
	if err != nil {
		h.t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func (h *hookEnv) commit(rel, content string) string {
	h.t.Helper()
	write(h.t, h.repo, rel, content, 0o644)
	h.git("add", "-A")
	h.git("commit", "-q", "-m", "c "+rel)
	return h.git("rev-parse", "HEAD")
}

// fake installs an executable named name that appends its arguments, and its
// standard input when stdin is true, to the calls log, then exits with status.
func (h *hookEnv) fake(dir, name string, status int, stdin bool) string {
	h.t.Helper()
	body := `echo "` + name + `: $*" >> "$CALLS"`
	if stdin {
		body += "\n" + `cat | sed "s/^/` + name + ` stdin: /" >> "$CALLS"`
	}
	body += "\nexit " + string(rune('0'+status))
	path := filepath.Join(dir, name)
	write(h.t, dir, name, "#!/bin/sh\n"+body+"\n", 0o755)
	return path
}

func (h *hookEnv) calls() []string {
	data, err := os.ReadFile(h.log)
	if err != nil {
		return nil
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n")
}

// run runs a hook script as git would, from the repository top level.
func (h *hookEnv) run(script string, stdin string, args ...string) (int, string) {
	h.t.Helper()
	path := filepath.Join(h.t.TempDir(), "hook")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		h.t.Fatal(err)
	}
	cmd := exec.Command(path, args...)
	cmd.Dir = h.repo
	cmd.Env = h.env
	cmd.Stdin = strings.NewReader(stdin)
	out, err := cmd.CombinedOutput()
	code := 0
	if exit, ok := err.(*exec.ExitError); ok {
		code = exit.ExitCode()
	} else if err != nil {
		h.t.Fatal(err)
	}
	return code, string(out)
}

func bodyOf(name string) string {
	for _, h := range Hooks() {
		if h.Name == name {
			return h.Body
		}
	}
	return ""
}

func TestBothScriptsFailWhenSedumIsFoundNowhere(t *testing.T) {
	for _, name := range []string{"commit-msg", "pre-push"} {
		h := newHookEnv(t)
		code, out := h.run(bodyOf(name), "", "msgfile")
		if code != 1 || !strings.Contains(out, "sedum hook: no sedum binary. Set SEDUM or put sedum on PATH.") {
			t.Errorf("%s: code %d output %q", name, code, out)
		}
	}
}

func TestScriptsPreferTheSEDUMVariableAndNeverRunAFileOfTheRepository(t *testing.T) {
	h := newHookEnv(t)
	h.fake(h.repo, "sedum", 0, false) // a file named sedum in the repository
	other := h.fake(t.TempDir(), "my-sedum", 0, false)
	h.env = append(h.env, "SEDUM="+other, "PATH="+h.env[len(h.env)-1][5:]+":.")

	code, _ := h.run(bodyOf("commit-msg"), "", "msg")

	if code != 0 || len(h.calls()) != 1 || !strings.HasPrefix(h.calls()[0], "my-sedum: check --staged --message-file msg") {
		t.Errorf("code %d calls %v", code, h.calls())
	}
}

func TestCommitMsgWithNoConfigRootRunsSedumCheckOnce(t *testing.T) {
	h := newHookEnv(t)
	h.fake(h.bin, "sedum", 0, false)
	h.commit("a.txt", "1")

	code, _ := h.run(bodyOf("commit-msg"), "", ".git/COMMIT_EDITMSG")

	if code != 0 || strings.Join(h.calls(), "|") != "sedum: check --staged --message-file .git/COMMIT_EDITMSG" {
		t.Errorf("code %d calls %v", code, h.calls())
	}
}

func TestCommitMsgRunsEachConfigRootThatHasStagedFilesAndSkipsTheOthers(t *testing.T) {
	h := newHookEnv(t)
	h.fake(h.bin, "sedum", 0, false)
	h.commit(".provenance.yml", "dir: provenance\n")
	h.commit("sub/.provenance.yml", "dir: provenance\n")
	h.commit("other/.linespec.yml", "provenance:\n  dir: provenance\n")
	h.commit("plain/.linespec.yml", "something: else\n")
	write(h.t, h.repo, "sub/x.txt", "1", 0o644)
	h.git("add", "sub/x.txt")

	code, _ := h.run(bodyOf("commit-msg"), "", "MSG")

	got := strings.Join(h.calls(), "|")
	want := []string{"sedum: check --staged --message-file MSG --config .provenance.yml", "sedum: check --staged --message-file MSG --config sub/.provenance.yml"}
	if code != 0 {
		t.Errorf("code %d", code)
	}
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("calls %q lack %q", got, w)
		}
	}
	if strings.Contains(got, "other/") || strings.Contains(got, "plain/") {
		t.Errorf("a root with no staged file or without a provenance block ran: %q", got)
	}
}

func TestCommitMsgRunsEveryRootAndBlocksWhenAnyCheckFailed(t *testing.T) {
	h := newHookEnv(t)
	h.fake(h.bin, "sedum", 1, false)
	h.commit(".provenance.yml", "dir: provenance\n")
	h.commit("sub/.provenance.yml", "dir: provenance\n")
	write(h.t, h.repo, "sub/x.txt", "1", 0o644)
	h.git("add", "sub/x.txt")

	code, out := h.run(bodyOf("commit-msg"), "", "MSG")

	if code != 1 || len(h.calls()) != 2 || !strings.Contains(out, "\nCommit blocked due to sedum check failures") {
		t.Errorf("code %d calls %v output %q", code, h.calls(), out)
	}
}

func TestPrePushFailsWhenProvenanceIsFoundNowhere(t *testing.T) {
	h := newHookEnv(t)
	h.fake(h.bin, "sedum", 0, false)

	code, out := h.run(bodyOf("pre-push"), "", "origin", "url")

	if code != 1 || !strings.Contains(out, "sedum hook: no provenance binary. Set PROVENANCE or put provenance on PATH.") {
		t.Errorf("code %d output %q", code, out)
	}
}

const zeros = "0000000000000000000000000000000000000000"

func pushLine(local, remote string) string {
	return "refs/heads/main " + local + " refs/heads/main " + remote + "\n"
}

func prePushEnv(t *testing.T, sedumStatus, provStatus int) *hookEnv {
	h := newHookEnv(t)
	h.fake(h.bin, "sedum", sedumStatus, false)
	h.fake(h.bin, "provenance", provStatus, true)
	return h
}

func TestPrePushRunsProvenancesHookFirstWithTheSavedInputThenChecksEachRange(t *testing.T) {
	h := prePushEnv(t, 0, 0)
	first := h.commit("a.txt", "1")
	second := h.commit("b.txt", "1")

	code, _ := h.run(bodyOf("pre-push"), pushLine(second, first)+"refs/heads/gone "+zeros+" refs/heads/gone "+first+"\n", "origin", "url")

	calls := h.calls()
	if code != 0 || len(calls) != 4 {
		t.Fatalf("code %d calls %v", code, calls)
	}
	if calls[0] != "provenance: hooks pre-push origin url" || !strings.HasPrefix(calls[1], "provenance stdin: refs/heads/main "+second) {
		t.Errorf("provenance calls %v", calls[:2])
	}
	if calls[3] != "sedum: check --range "+first+".."+second {
		t.Errorf("sedum call %q (a deletion must be skipped)", calls[3])
	}
}

func TestPrePushChecksTheOldestUnpushedCommitsRangeForANewRef(t *testing.T) {
	h := prePushEnv(t, 0, 0)
	base := h.commit("a.txt", "1")
	h.git("update-ref", "refs/remotes/origin/main", base)
	h.commit("b.txt", "1")
	tip := h.commit("c.txt", "1")

	code, _ := h.run(bodyOf("pre-push"), pushLine(tip, zeros), "origin", "url")

	last := h.calls()[len(h.calls())-1]
	if code != 0 || last != "sedum: check --range "+base+".."+tip {
		t.Errorf("code %d last call %q", code, last)
	}
}

func TestPrePushTreatsARemoteShaItDoesNotKnowLikeANewRef(t *testing.T) {
	h := prePushEnv(t, 0, 0)
	base := h.commit("a.txt", "1")
	h.git("update-ref", "refs/remotes/origin/main", base)
	tip := h.commit("b.txt", "1")

	code, _ := h.run(bodyOf("pre-push"), pushLine(tip, "1111111111111111111111111111111111111111"), "origin", "url")

	last := h.calls()[len(h.calls())-1]
	if code != 0 || last != "sedum: check --range "+base+".."+tip {
		t.Errorf("code %d last call %q", code, last)
	}
}

func TestPrePushChecksARootCommitWithProvenanceAlone(t *testing.T) {
	h := prePushEnv(t, 0, 0)
	root := h.commit("a.txt", "1")

	code, _ := h.run(bodyOf("pre-push"), pushLine(root, zeros), "origin", "url")

	last := h.calls()[len(h.calls())-1]
	if code != 0 || last != "provenance: check --commit "+root {
		t.Errorf("code %d last call %q", code, last)
	}
	for _, c := range h.calls() {
		if strings.HasPrefix(c, "sedum:") {
			t.Errorf("sedum ran for a root commit: %q", c)
		}
	}
}

func TestPrePushChecksNothingWhenEveryCommitIsOnARemote(t *testing.T) {
	h := prePushEnv(t, 0, 0)
	tip := h.commit("a.txt", "1")
	h.git("update-ref", "refs/remotes/origin/main", tip)

	code, _ := h.run(bodyOf("pre-push"), pushLine(tip, zeros), "origin", "url")

	for _, c := range h.calls() {
		if strings.HasPrefix(c, "sedum:") || strings.HasPrefix(c, "provenance: check") {
			t.Errorf("a check ran: %q", c)
		}
	}
	if code != 0 {
		t.Errorf("code %d", code)
	}
}

func TestPrePushSkipsShortLinesAndALocalShaItAlreadyChecked(t *testing.T) {
	h := prePushEnv(t, 0, 0)
	first := h.commit("a.txt", "1")
	second := h.commit("b.txt", "1")
	input := "short line\n" + pushLine(second, first) + "refs/heads/other " + second + " refs/heads/other " + first + "\n"

	code, _ := h.run(bodyOf("pre-push"), input, "origin", "url")

	n := 0
	for _, c := range h.calls() {
		if strings.HasPrefix(c, "sedum:") {
			n++
		}
	}
	if code != 0 || n != 1 {
		t.Errorf("code %d, %d sedum checks, want 1", code, n)
	}
}

func TestPrePushBlocksWhenProvenancesHookOrACheckFailed(t *testing.T) {
	for name, statuses := range map[string][2]int{"provenance": {1, 0}, "sedum": {0, 1}} {
		h := prePushEnv(t, statuses[1], statuses[0])
		first := h.commit("a.txt", "1")
		second := h.commit("b.txt", "1")

		code, out := h.run(bodyOf("pre-push"), pushLine(second, first), "origin", "url")

		if code != 1 || !strings.Contains(out, "\nPush blocked due to failed checks") {
			t.Errorf("%s fails: code %d output %q", name, code, out)
		}
		if name == "provenance" && len(h.calls()) < 3 {
			t.Errorf("the checks must still run after provenance's hook failed: %v", h.calls())
		}
	}
}

func TestPrePushRemovesItsTemporaryFile(t *testing.T) {
	h := prePushEnv(t, 0, 0)
	tmp := t.TempDir()
	h.env = append(h.env, "TMPDIR="+tmp)
	first := h.commit("a.txt", "1")
	second := h.commit("b.txt", "1")

	h.run(bodyOf("pre-push"), pushLine(second, first), "origin", "url")

	if entries, _ := os.ReadDir(tmp); len(entries) != 0 {
		t.Errorf("temporary files left behind: %v", entries)
	}
}
