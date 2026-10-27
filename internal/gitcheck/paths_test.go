package gitcheck

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// The paths unit lists commits, the changed paths, the base and the message
// through git, and classifies paths as governed, record files, state files or
// recordings.

// testEnv is an environment with no GIT_* variable and a private HOME, so no
// test reads or writes the developer's git configuration or repository.
func testEnv(t *testing.T) []string {
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

type repoT struct {
	t   *testing.T
	dir string
	env []string
}

func newRepo(t *testing.T) *repoT {
	t.Helper()
	r := &repoT{t: t, dir: t.TempDir(), env: testEnv(t)}
	r.git("init", "-q", "-b", "main")
	return r
}

func (r *repoT) git(args ...string) string {
	r.t.Helper()
	cmd := exec.Command("git", append([]string{"-C", r.dir}, args...)...)
	cmd.Env = r.env
	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func (r *repoT) write(rel, content string) {
	r.t.Helper()
	full := filepath.Join(r.dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		r.t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		r.t.Fatal(err)
	}
}

func (r *repoT) commit(msg string, files map[string]string) string {
	r.t.Helper()
	for rel, content := range files {
		r.write(rel, content)
	}
	r.git("add", "-A")
	r.git("commit", "-q", "-m", msg)
	return r.git("rev-parse", "HEAD")
}

func TestCommitsListsTheRangeOldestFirst(t *testing.T) {
	r := newRepo(t)
	a := r.commit("a", map[string]string{"a.txt": "1"})
	b := r.commit("b", map[string]string{"b.txt": "1"})
	c := r.commit("c", map[string]string{"c.txt": "1"})

	got, err := Commits(r.dir, a+".."+c, r.env)

	if err != nil || !reflect.DeepEqual(got, []string{b, c}) {
		t.Errorf("Commits = %v, %v; want [%s %s]", got, err, b, c)
	}
}

func TestCommitsFailsWithAGitErrorOnABadRange(t *testing.T) {
	r := newRepo(t)
	r.commit("a", map[string]string{"a.txt": "1"})
	if _, err := Commits(r.dir, "nope..HEAD", r.env); err == nil || !strings.HasPrefix(err.Error(), "git:") {
		t.Errorf("err = %v, want a git: error", err)
	}
}

func sorted(s []string) []string {
	out := append([]string{}, s...)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

func TestChangedPathsOfACommitAreTheDifferencesAgainstItsFirstParent(t *testing.T) {
	r := newRepo(t)
	r.commit("a", map[string]string{"keep.txt": "k", "old.txt": "o\nmore\nlines\nhere\n", "gone.txt": "g", "mod.txt": "1"})
	r.git("mv", "old.txt", "new.txt")
	r.git("rm", "-q", "gone.txt")
	r.write("mod.txt", "2")
	r.write("added.txt", "x")
	r.git("add", "-A")
	r.git("commit", "-q", "-m", "change")

	got, err := ChangedPaths(r.dir, "HEAD", false, r.env)

	want := []string{"added.txt", "gone.txt", "mod.txt", "new.txt", "old.txt"}
	if err != nil || !reflect.DeepEqual(sorted(got), want) {
		t.Errorf("ChangedPaths = %v, %v; want %v (a rename is one deleted and one added path)", got, err, want)
	}
}

func TestChangedPathsOfARootCommitAreAllThePathsOfItsTree(t *testing.T) {
	r := newRepo(t)
	r.commit("root", map[string]string{"a.txt": "1", "d/b.txt": "2"})

	got, err := ChangedPaths(r.dir, "HEAD", false, r.env)

	if err != nil || !reflect.DeepEqual(sorted(got), []string{"a.txt", "d/b.txt"}) {
		t.Errorf("ChangedPaths = %v, %v", got, err)
	}
}

func TestChangedPathsOfAMergeDifferFromEveryParent(t *testing.T) {
	r := newRepo(t)
	r.commit("base", map[string]string{"shared.txt": "0", "main-only.txt": "0", "side-only.txt": "0"})
	r.git("checkout", "-q", "-b", "side")
	r.commit("side", map[string]string{"side-only.txt": "1"})
	r.git("checkout", "-q", "main")
	r.commit("main", map[string]string{"main-only.txt": "1"})
	r.git("merge", "-q", "--no-ff", "-m", "merge", "side")
	// A path the merge itself changes differs from both parents.
	r.write("shared.txt", "merged")
	r.git("add", "-A")
	r.git("commit", "-q", "--amend", "-m", "merge")

	got, err := ChangedPaths(r.dir, "HEAD", false, r.env)

	if err != nil || !reflect.DeepEqual(got, []string{"shared.txt"}) {
		t.Errorf("ChangedPaths = %v, %v; want only the path that differs from both parents", got, err)
	}
}

func TestChangedPathsStagedAreTheDifferencesBetweenTheIndexAndHEAD(t *testing.T) {
	r := newRepo(t)
	r.commit("a", map[string]string{"a.txt": "1", "b.txt": "1"})
	r.write("a.txt", "2")
	r.write("c.txt", "3")
	r.git("add", "a.txt", "c.txt")
	r.write("b.txt", "unstaged edit")

	got, err := ChangedPaths(r.dir, "", true, r.env)

	if err != nil || !reflect.DeepEqual(sorted(got), []string{"a.txt", "c.txt"}) {
		t.Errorf("ChangedPaths = %v, %v", got, err)
	}
}

func TestChangedPathsStagedWithNoHEADAreEveryStagedPath(t *testing.T) {
	r := newRepo(t)
	r.write("a.txt", "1")
	r.git("add", "a.txt")

	got, err := ChangedPaths(r.dir, "", true, r.env)

	if err != nil || !reflect.DeepEqual(got, []string{"a.txt"}) {
		t.Errorf("ChangedPaths = %v, %v", got, err)
	}
}

func TestChangedPathsAreRelativeToTheBaseDirectoryAndDropThoseOutsideIt(t *testing.T) {
	r := newRepo(t)
	r.commit("a", map[string]string{"app/x.txt": "1", "app/sub/y.txt": "1", "other/z.txt": "1"})
	r.commit("b", map[string]string{"app/x.txt": "2", "app/sub/y.txt": "2", "other/z.txt": "2"})

	got, err := ChangedPaths(filepath.Join(r.dir, "app"), "HEAD", false, r.env)

	if err != nil || !reflect.DeepEqual(sorted(got), []string{"sub/y.txt", "x.txt"}) {
		t.Errorf("ChangedPaths = %v, %v; want paths relative to app without other/z.txt", got, err)
	}
}

func TestBase(t *testing.T) {
	r := newRepo(t)
	first := r.commit("a", map[string]string{"a.txt": "1"})
	second := r.commit("b", map[string]string{"b.txt": "1"})

	if got, ok, err := Base(r.dir, second, false, r.env); err != nil || !ok || got != first {
		t.Errorf("Base of a commit = %q %v %v, want its first parent", got, ok, err)
	}
	if _, ok, err := Base(r.dir, first, false, r.env); err != nil || ok {
		t.Errorf("Base of a root commit = ok %v err %v, want no base", ok, err)
	}
	if got, ok, err := Base(r.dir, "", true, r.env); err != nil || !ok || got != second {
		t.Errorf("Base for --staged = %q %v %v, want HEAD", got, ok, err)
	}
	if _, _, err := Base(r.dir, "deadbeef", false, r.env); err == nil || !strings.HasPrefix(err.Error(), "git:") {
		t.Errorf("Base of an unknown revision: err = %v, want a git: error", err)
	}
}

func TestBaseForStagedWithNoHEADIsNoBase(t *testing.T) {
	r := newRepo(t)
	if _, ok, err := Base(r.dir, "", true, r.env); err != nil || ok {
		t.Errorf("ok = %v, err = %v; want no base and no error", ok, err)
	}
}

func TestBaseOfAMergeIsItsFirstParent(t *testing.T) {
	r := newRepo(t)
	r.commit("base", map[string]string{"a.txt": "0"})
	r.git("checkout", "-q", "-b", "side")
	r.commit("side", map[string]string{"s.txt": "1"})
	r.git("checkout", "-q", "main")
	main := r.commit("main", map[string]string{"m.txt": "1"})
	r.git("merge", "-q", "--no-ff", "-m", "merge", "side")

	if got, ok, err := Base(r.dir, "HEAD", false, r.env); err != nil || !ok || got != main {
		t.Errorf("Base = %q %v %v, want %s", got, ok, err, main)
	}
}

func TestCommitMessage(t *testing.T) {
	r := newRepo(t)
	r.commit("first subject\n\nbody [prov-2026-aaaaaaaa@2]", map[string]string{"a.txt": "1"})
	r.git("checkout", "-q", "-b", "side")
	r.commit("side", map[string]string{"s.txt": "1"})
	r.git("checkout", "-q", "main")
	r.commit("main", map[string]string{"m.txt": "1"})
	r.git("merge", "-q", "--no-ff", "-m", "the merge's own message", "side")
	file := filepath.Join(t.TempDir(), "MSG")
	if err := os.WriteFile(file, []byte("  from the file \n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if got, err := CommitMessage(r.dir, "HEAD~1", false, "", r.env); err != nil || got != "main" {
		t.Errorf("message of a commit = %q, %v", got, err)
	}
	if got, err := CommitMessage(r.dir, "HEAD", false, "", r.env); err != nil || got != "the merge's own message" {
		t.Errorf("message of a merge = %q, %v; want its own", got, err)
	}
	if got, err := CommitMessage(r.dir, "", true, file, r.env); err != nil || got != "from the file" {
		t.Errorf("staged message with a file = %q, %v", got, err)
	}
	if got, err := CommitMessage(r.dir, "", true, "", r.env); err != nil || got != "the merge's own message" {
		t.Errorf("staged message without a file = %q, %v; want HEAD's message", got, err)
	}
	if _, err := CommitMessage(r.dir, "deadbeef", false, "", r.env); err == nil || !strings.HasPrefix(err.Error(), "git:") {
		t.Errorf("unknown commit: err = %v, want a git: error", err)
	}
}

func TestCommitMessageStagedWithoutAFileAndWithoutHEADIsAGitError(t *testing.T) {
	r := newRepo(t)
	if _, err := CommitMessage(r.dir, "", true, "", r.env); err == nil || !strings.HasPrefix(err.Error(), "git:") {
		t.Errorf("err = %v, want a git: error", err)
	}
}

func TestTopLevel(t *testing.T) {
	r := newRepo(t)
	r.commit("a", map[string]string{"app/x.txt": "1"})
	want, err := filepath.EvalSymlinks(r.dir)
	if err != nil {
		t.Fatal(err)
	}
	got, err := TopLevel(filepath.Join(r.dir, "app"), r.env)
	if err != nil {
		t.Fatal(err)
	}
	if got, _ = filepath.EvalSymlinks(got); got != want {
		t.Errorf("TopLevel = %q, want %q", got, want)
	}
	if _, err := TopLevel(t.TempDir(), r.env); err == nil || !strings.HasPrefix(err.Error(), "git:") {
		t.Errorf("outside a repository: err = %v, want a git: error", err)
	}
}

func TestTreeReaders(t *testing.T) {
	r := newRepo(t)
	first := r.commit("a", map[string]string{"app/x.txt": "old", "app/link-target.txt": "t"})
	r.commit("b", map[string]string{"app/x.txt": "new"})
	base := BaseReader(filepath.Join(r.dir, "app"), first, r.env)
	dir := DirReader(filepath.Join(r.dir, "app"))

	if got, found, err := base.Read("x.txt"); err != nil || !found || string(got) != "old" {
		t.Errorf("base Read = %q %v %v", got, found, err)
	}
	if _, found, err := base.Read("missing.txt"); err != nil || found {
		t.Errorf("base Read of a missing file: found %v err %v", found, err)
	}
	if !base.Regular("x.txt") || base.Regular("missing.txt") {
		t.Error("base Regular is wrong")
	}
	if got, found, err := dir.Read("x.txt"); err != nil || !found || string(got) != "new" {
		t.Errorf("dir Read = %q %v %v", got, found, err)
	}
	if _, found, _ := dir.Read("missing.txt"); found {
		t.Error("dir Read found a missing file")
	}
	if !dir.Regular("x.txt") || dir.Regular("missing.txt") || dir.Regular(".") {
		t.Error("dir Regular is wrong")
	}
}

func TestGovernedChanges(t *testing.T) {
	r := newRepo(t)
	first := r.commit("a", map[string]string{"a.txt": "old", "b.txt": "old", "gone.txt": "g"})
	r.git("rm", "-q", "gone.txt")
	r.write("a.txt", "new")
	r.write("added.txt", "x")
	r.git("add", "-A")
	r.git("commit", "-q", "-m", "b")
	base := BaseReader(r.dir, first, r.env)
	commit := DirReader(r.dir)

	files, ok := GovernedChanges([]string{"a.txt"}, base, commit)
	if !ok || len(files) != 1 || files[0].Path != "a.txt" || string(files[0].Base) != "old" || string(files[0].Commit) != "new" {
		t.Errorf("files = %+v ok = %v", files, ok)
	}
	for name, paths := range map[string][]string{
		"empty":             nil,
		"deleted":           {"a.txt", "gone.txt"},
		"absent at base":    {"a.txt", "added.txt"},
		"unreadable commit": {"a.txt", "never-existed.txt"},
	} {
		if _, ok := GovernedChanges(paths, base, commit); ok {
			t.Errorf("%s: GovernedChanges returned true", name)
		}
	}
}

func TestClassify(t *testing.T) {
	cfg := PathConfig{RecordsDir: "provenance", ExcludePaths: []string{"README.md", "generators/**", ".sedum/**", "docs/**"}}
	changed := []string{
		"app/main.go",                       // governed
		"README.md",                         // excluded
		"generators/go/sedum.yaml",          // excluded
		"provenance/prov-2026-aaaaaaaa.yml", // record file
		"provenance/.yml",                   // a state file: no name before the suffix
		"provenance/.versions/prov-2026-aaaaaaaa/1-x.yml", // state file
		".provenance/versions.jsonl",                      // state file
		".sedum/recordings/prov-2026-aaaaaaaa.json",       // recording, excluded too
		".sedum/recordings/sub/x.json",                    // not a recording: excluded
		"../outside/x.go",                                 // outside the base directory
	}

	got := Classify(changed, cfg)

	want := Classified{
		Governed:           []string{"app/main.go"},
		Excluded:           []string{"README.md", "generators/go/sedum.yaml", ".sedum/recordings/sub/x.json"},
		Recordings:         []string{".sedum/recordings/prov-2026-aaaaaaaa.json"},
		ExcludedRecordings: []string{".sedum/recordings/prov-2026-aaaaaaaa.json"},
		RecordFiles:        []string{"provenance/prov-2026-aaaaaaaa.yml"},
		StateFiles:         []string{"provenance/.yml", "provenance/.versions/prov-2026-aaaaaaaa/1-x.yml", ".provenance/versions.jsonl"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Classify =\n%+v\nwant\n%+v", got, want)
	}
}

func TestARecordFileOrStateFileStaysBlockedEvenWhenExcluded(t *testing.T) {
	cfg := PathConfig{RecordsDir: "provenance", ExcludePaths: []string{"provenance/**", ".provenance/**"}}

	record := Classify([]string{"provenance/prov-2026-aaaaaaaa.yml"}, cfg)
	state := Classify([]string{".provenance/versions.jsonl"}, cfg)

	if !record.Blocked() || !state.Blocked() || len(record.Excluded) != 0 || len(state.Excluded) != 0 {
		t.Errorf("record = %+v state = %+v; want both blocked and neither excluded", record, state)
	}
	if Classify([]string{"app/x.go", "README.md"}, PathConfig{ExcludePaths: []string{"README.md"}}).Blocked() {
		t.Error("ordinary paths blocked")
	}
}

func TestAnEmptyRecordsDirectoryTreatsNothingAsARecordFile(t *testing.T) {
	got := Classify([]string{"a.yml"}, PathConfig{})
	if len(got.RecordFiles) != 0 || !reflect.DeepEqual(got.Governed, []string{"a.yml"}) {
		t.Errorf("Classify = %+v", got)
	}
}
