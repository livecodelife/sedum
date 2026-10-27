package gitcheck

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/livecodelife/sedum/internal/pathpat"
)

// git runs git with -C dir and the environment env and returns its trimmed
// standard output. Every failure is a GitError carrying git's own message.
func git(dir string, env []string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = env
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", GitError(msg)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

// TopLevel returns the git top level of the directory.
func TopLevel(baseDir string, env []string) (string, error) {
	return git(baseDir, env, "rev-parse", "--show-toplevel")
}

// Commits lists the commits of a range, oldest first.
func Commits(baseDir, rng string, env []string) ([]string, error) {
	out, err := git(baseDir, env, "rev-list", "--reverse", rng)
	if err != nil {
		return nil, err
	}
	return strings.Fields(out), nil
}

// Base returns the revision a commit is compared against: its first parent, or
// HEAD when staged. The bool is false when there is none: a root commit, or a
// staged change in a repository with no HEAD.
func Base(baseDir, rev string, staged bool, env []string) (string, bool, error) {
	if staged {
		out, err := git(baseDir, env, "rev-parse", "--verify", "-q", "HEAD")
		if err != nil {
			return "", false, nil
		}
		return out, true, nil
	}
	if _, err := git(baseDir, env, "rev-parse", "--verify", "-q", rev+"^{commit}"); err != nil {
		return "", false, GitError("cannot read commit " + rev)
	}
	out, err := git(baseDir, env, "rev-parse", "--verify", "-q", rev+"^1")
	if err != nil {
		return "", false, nil
	}
	return out, true, nil
}

// CommitMessage returns the message of a commit. Staged, it is the message file
// when one is given, otherwise the message of HEAD.
func CommitMessage(baseDir, rev string, staged bool, messageFile string, env []string) (string, error) {
	if staged && messageFile != "" {
		data, err := os.ReadFile(messageFile)
		if err != nil {
			return "", GitError("cannot read the commit message: " + err.Error())
		}
		return strings.TrimSpace(string(data)), nil
	}
	if staged {
		rev = "HEAD"
	}
	out, err := git(baseDir, env, "log", "-1", "--format=%B", rev)
	if err != nil {
		return "", GitError("cannot read the commit message of " + rev)
	}
	return strings.TrimSpace(out), nil
}

// prefixOf returns the path of baseDir below the git top level, with a trailing
// slash, or empty when baseDir is the top level.
func prefixOf(baseDir string, env []string) (string, error) {
	return git(baseDir, env, "rev-parse", "--show-prefix")
}

func nameList(out string) []string {
	var names []string
	for _, n := range strings.Split(out, "\x00") {
		if n != "" {
			names = append(names, n)
		}
	}
	return names
}

// ChangedPaths returns the paths a commit changed, relative to baseDir; paths
// outside baseDir are dropped. A commit with a first parent changes the paths
// that differ from it, a rename counting as one deleted and one added path; a
// root commit changes every path of its tree; a merge changes the paths that
// differ from every parent; staged, the paths that differ between the index and
// HEAD.
func ChangedPaths(baseDir, rev string, staged bool, env []string) ([]string, error) {
	prefix, err := prefixOf(baseDir, env)
	if err != nil {
		return nil, err
	}
	var names []string
	switch {
	case staged:
		out, err := git(baseDir, env, "diff", "--cached", "--no-renames", "--name-only", "-z", "--no-relative")
		if err != nil {
			return nil, err
		}
		names = nameList(out)
	default:
		parents, err := git(baseDir, env, "rev-list", "--parents", "-n", "1", rev)
		if err != nil {
			return nil, GitError("cannot read commit " + rev)
		}
		fields := strings.Fields(parents)
		if len(fields) == 0 {
			return nil, GitError("cannot read commit " + rev)
		}
		if len(fields) == 1 {
			out, err := git(baseDir, env, "diff-tree", "--root", "-r", "--no-commit-id", "--name-only", "-z", "--no-renames", rev)
			if err != nil {
				return nil, err
			}
			names = nameList(out)
			break
		}
		var common map[string]bool
		for _, parent := range fields[1:] {
			out, err := git(baseDir, env, "diff-tree", "-r", "--no-commit-id", "--name-only", "-z", "--no-renames", parent, rev)
			if err != nil {
				return nil, err
			}
			here := map[string]bool{}
			for _, n := range nameList(out) {
				if common == nil || common[n] {
					here[n] = true
				}
			}
			common = here
		}
		out, err := git(baseDir, env, "diff-tree", "-r", "--no-commit-id", "--name-only", "-z", "--no-renames", fields[1], rev)
		if err != nil {
			return nil, err
		}
		for _, n := range nameList(out) {
			if common[n] {
				names = append(names, n)
			}
		}
	}
	var rel []string
	for _, n := range names {
		if prefix == "" {
			rel = append(rel, n)
		} else if strings.HasPrefix(n, prefix) {
			rel = append(rel, strings.TrimPrefix(n, prefix))
		}
	}
	return rel, nil
}

// TreeReader reads the files of a tree: the files at a revision, or in a
// directory. Paths are relative to the tree's base directory.
type TreeReader struct {
	read    func(path string) ([]byte, bool, error)
	regular func(path string) bool
}

// Read returns the content of the file at path and whether it exists.
func (t TreeReader) Read(path string) ([]byte, bool, error) { return t.read(path) }

// Regular reports whether path is a regular file.
func (t TreeReader) Regular(path string) bool { return t.regular(path) }

// BaseReader reads the files of baseDir at a revision.
func BaseReader(baseDir, rev string, env []string) TreeReader {
	return TreeReader{
		read: func(p string) ([]byte, bool, error) {
			if _, err := git(baseDir, env, "cat-file", "-e", rev+":./"+p); err != nil {
				return nil, false, nil
			}
			cmd := exec.Command("git", "-C", baseDir, "cat-file", "blob", rev+":./"+p)
			cmd.Env = env
			out, err := cmd.Output()
			if err != nil {
				return nil, false, GitError("cannot read " + p + " at " + rev)
			}
			return out, true, nil
		},
		regular: func(p string) bool {
			out, err := git(baseDir, env, "ls-tree", rev, "./"+p)
			return err == nil && (strings.HasPrefix(out, "100644 blob") || strings.HasPrefix(out, "100755 blob"))
		},
	}
}

// DirReader reads the files of a directory.
func DirReader(dir string) TreeReader {
	return TreeReader{
		read: func(p string) ([]byte, bool, error) {
			data, err := os.ReadFile(dir + string(os.PathSeparator) + p)
			if os.IsNotExist(err) {
				return nil, false, nil
			}
			if err != nil {
				return nil, false, err
			}
			return data, true, nil
		},
		regular: func(p string) bool {
			info, err := os.Stat(dir + string(os.PathSeparator) + p)
			return err == nil && info.Mode().IsRegular()
		},
	}
}

// GovernedFile is a governed path with its content at the base and at the commit.
type GovernedFile struct {
	Path   string
	Base   []byte
	Commit []byte
}

// GovernedChanges reads the governed files at the base and at the commit. The
// bool is false when there is no governed path, or when one was deleted by the
// commit, cannot be read at the commit, or did not exist at the base.
func GovernedChanges(governed []string, base, commit TreeReader) ([]GovernedFile, bool) {
	if len(governed) == 0 {
		return nil, false
	}
	var files []GovernedFile
	for _, p := range governed {
		after, found, err := commit.Read(p)
		if err != nil || !found {
			return nil, false
		}
		before, found, err := base.Read(p)
		if err != nil || !found {
			return nil, false
		}
		files = append(files, GovernedFile{Path: p, Base: before, Commit: after})
	}
	return files, true
}

// PathConfig is the part of the provenance config the path rules read: the
// records directory and the exclude_paths entries, both relative to the base
// directory.
type PathConfig struct {
	RecordsDir   string
	ExcludePaths []string
}

// Classified is a set of changed paths sorted by how the proofs treat them.
type Classified struct {
	Governed           []string
	Excluded           []string
	Recordings         []string
	ExcludedRecordings []string
	RecordFiles        []string
	StateFiles         []string
}

// Blocked reports whether a record file or a provenance state file changed.
func (c Classified) Blocked() bool {
	return len(c.RecordFiles) > 0 || len(c.StateFiles) > 0
}

const recordingsDir = ".sedum/recordings/"

// Classify sorts changed paths, relative to the base directory. A path outside
// the base directory is dropped first. Then a path is a record file (a .yml
// file with a name directly in the records directory), a provenance state file
// (under .provenance/ or elsewhere under the records directory), a recording
// (a file directly in .sedum/recordings/, also listed as an excluded recording
// when exclude_paths excludes it), excluded, or governed, tested in that order.
func Classify(changed []string, cfg PathConfig) Classified {
	var c Classified
	records := strings.Trim(cfg.RecordsDir, "/")
	for _, p := range changed {
		if p == ".." || strings.HasPrefix(p, "../") || path.IsAbs(p) {
			continue
		}
		_, excluded := pathpat.MatchAny(cfg.ExcludePaths, p)
		switch {
		case records != "" && isRecordFile(p, records):
			c.RecordFiles = append(c.RecordFiles, p)
		case strings.HasPrefix(p, ".provenance/") || (records != "" && strings.HasPrefix(p, records+"/")):
			c.StateFiles = append(c.StateFiles, p)
		case isRecording(p):
			c.Recordings = append(c.Recordings, p)
			if excluded {
				c.ExcludedRecordings = append(c.ExcludedRecordings, p)
			}
		case excluded:
			c.Excluded = append(c.Excluded, p)
		default:
			c.Governed = append(c.Governed, p)
		}
	}
	return c
}

func isRecordFile(p, records string) bool {
	rest, ok := strings.CutPrefix(p, records+"/")
	return ok && !strings.Contains(rest, "/") && strings.HasSuffix(rest, ".yml") && len(rest) > len(".yml")
}

func isRecording(p string) bool {
	rest, ok := strings.CutPrefix(p, recordingsDir)
	return ok && rest != "" && !strings.Contains(rest, "/")
}
