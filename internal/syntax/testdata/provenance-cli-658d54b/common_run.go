package main

// Helpers every command shares: config, the record store, alias-aware id
// resolution, the version manifest, git and output. Temporary bridge: these
// move into `go` generator regions once that package lands (prov-2026-60e72d0e).

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/livecodelife/provenance/record"
)

const (
	configFile       = ".provenance.yml"
	legacyConfigFile = ".linespec.yml"
	manifestPath     = ".provenance/versions.jsonl"
	legacyManifest   = ".linespec/hash_manifest.json"
)

// sharedRepo is one entry of shared_repos.
type sharedRepo struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Ref  string `yaml:"ref"`
	Dir  string `yaml:"dir"`
}

// config is the provenance configuration: the top level of .provenance.yml,
// or the provenance block of .linespec.yml before a repo is migrated.
type config struct {
	Enforcement                  string       `yaml:"enforcement"`
	Dir                          string       `yaml:"dir"`
	ExcludePaths                 []string     `yaml:"exclude_paths"`
	SharedRepos                  []sharedRepo `yaml:"shared_repos"`
	CacheTTLMinutes              int          `yaml:"cache_ttl_minutes"`
	CommitTagRequired            bool         `yaml:"commit_tag_required"`
	AutoAffectedScope            bool         `yaml:"auto_affected_scope"`
	RunAssociatedSpecsOnComplete bool         `yaml:"run_associated_specs_on_complete"`
	OverlapSpecsOnComplete       string       `yaml:"overlap_specs_on_complete"`
	CommitOnStatusChange         bool         `yaml:"commit_on_status_change"`
	WriteRestriction             *bool        `yaml:"write_restriction,omitempty"`
}

// repo is the repository a command acts on.
type repo struct {
	Root   string
	Config config
	// ConfigPath is the file the config came from; empty when defaults apply.
	ConfigPath string
	// Legacy is set when the config came from .linespec.yml.
	Legacy bool
	Store  *record.FilesStore
	out    *os.File
	// env caches gitEnv.
	env []string
}

// openRepo finds the repository and its config. -c/--config names the file;
// otherwise the repo root is git's top level (or the working directory) and
// .provenance.yml wins over .linespec.yml's provenance block.
func openRepo(cmd *cobra.Command) (*repo, error) {
	r := &repo{out: os.Stdout}
	explicit, _ := cmd.Flags().GetString("config")
	if explicit != "" {
		abs, err := filepath.Abs(explicit)
		if err != nil {
			return nil, err
		}
		r.Root = filepath.Dir(abs)
		r.ConfigPath = abs
		r.Legacy = filepath.Base(abs) == legacyConfigFile
	} else {
		r.Root = repoRoot()
		for _, name := range []string{configFile, legacyConfigFile} {
			p := filepath.Join(r.Root, name)
			if _, err := os.Stat(p); err == nil {
				r.ConfigPath, r.Legacy = p, name == legacyConfigFile
				break
			}
		}
	}
	if r.ConfigPath != "" {
		cfg, err := readConfig(r.ConfigPath, r.Legacy)
		if err != nil {
			return nil, err
		}
		r.Config = cfg
	}
	if r.Config.Dir == "" {
		r.Config.Dir = "provenance"
	}
	r.Config.Dir = strings.TrimRight(r.Config.Dir, "/")
	r.Store = &record.FilesStore{Root: r.Root, Dir: r.Config.Dir}
	return r, nil
}

func repoRoot() string {
	if out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output(); err == nil {
		return strings.TrimSpace(string(out))
	}
	wd, _ := os.Getwd()
	return wd
}

func readConfig(path string, legacy bool) (config, error) {
	var cfg config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if legacy {
		var wrap struct {
			Provenance config `yaml:"provenance"`
		}
		if err := yaml.Unmarshal(data, &wrap); err != nil {
			return cfg, fmt.Errorf("%s: %w", path, err)
		}
		return wrap.Provenance, nil
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("%s: %w", path, err)
	}
	return cfg, nil
}

// errNotFound is returned, after printing linespec's message, when an id
// matches neither a record nor any record's aliases.
var errNotFound = errors.New("record not found")

// resolve loads the record named by id, or the record that lists id among
// its aliases.
func (r *repo) resolve(id string) (*record.Record, *record.Document, error) {
	rec, doc, err := r.Store.Load(id)
	if err == nil {
		return rec, doc, nil
	}
	if !os.IsNotExist(errors.Unwrap(err)) && !os.IsNotExist(err) {
		return nil, nil, err
	}
	ids, err := r.Store.List()
	if err != nil {
		return nil, nil, err
	}
	for _, other := range ids {
		rec, doc, err := r.Store.Load(other)
		if err != nil {
			continue
		}
		for _, a := range rec.Aliases {
			if a == id {
				return rec, doc, nil
			}
		}
	}
	r.fail("Record not found: %s", id)
	return nil, nil, errNotFound
}

// recordPath is the record file's path relative to the repo root.
func (r *repo) recordPath(id string) string {
	return filepath.Join(r.Config.Dir, id+".yml")
}

// manifestLine is one line of the version manifest.
type manifestLine struct {
	ID          string `json:"id"`
	Version     int    `json:"version"`
	ContentHash string `json:"content_hash"`
}

// writeManifest rewrites .provenance/versions.jsonl from the latest snapshot
// of every record that has one, sorted by id.
func (r *repo) writeManifest() error {
	ids, err := r.Store.List()
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	for _, id := range ids {
		vs, err := r.Store.Versions(id)
		if err != nil {
			return err
		}
		if len(vs) == 0 {
			continue
		}
		snap, err := r.Store.Snapshot(id, vs[len(vs)-1])
		if err != nil {
			return err
		}
		line, err := json.Marshal(manifestLine{ID: id, Version: snap.Version, ContentHash: snap.ContentHash})
		if err != nil {
			return err
		}
		buf.Write(line)
		buf.WriteByte('\n')
	}
	path := filepath.Join(r.Root, manifestPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

// versionsDir is a record's snapshot directory relative to the repo root.
func (r *repo) versionsDir(id string) string {
	return filepath.Join(r.Config.Dir, ".versions", id)
}

// gitEnv is the environment for git commands on r.Root. A hook exports
// GIT_DIR, GIT_INDEX_FILE and friends for the repository being committed.
// They are kept when that repository is r.Root's own, because check --staged
// must read the index being committed (git commit -a and git commit <paths>
// use a temporary one), and dropped otherwise, so a command acting on another
// repository never reads or writes the hooked one (the 2026-10-15 incident:
// fixture git commands rewrote the real .git/config).
func (r *repo) gitEnv() []string {
	if r.env != nil {
		return r.env
	}
	env := os.Environ()
	var clean []string
	inherited := false
	for _, kv := range env {
		if strings.HasPrefix(kv, "GIT_") {
			inherited = true
			continue
		}
		clean = append(clean, kv)
	}
	r.env = clean
	if !inherited {
		r.env = env
		return r.env
	}
	hooked := os.Getenv("GIT_DIR")
	if hooked == "" {
		wd, _ := os.Getwd()
		hooked = gitDirOf(wd, clean)
	} else if !filepath.IsAbs(hooked) {
		wd, _ := os.Getwd()
		hooked = filepath.Join(wd, hooked)
	}
	if own := gitDirOf(r.Root, clean); own != "" && samePath(own, hooked) {
		r.env = env
	}
	return r.env
}

// gitDirOf is the absolute git dir of the repository holding dir, found with
// env (which carries no GIT_* variables), or "" when there is none.
func gitDirOf(dir string, env []string) string {
	c := exec.Command("git", "-C", dir, "rev-parse", "--absolute-git-dir")
	c.Env = env
	out, err := c.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func samePath(a, b string) bool {
	ra, errA := filepath.EvalSymlinks(a)
	rb, errB := filepath.EvalSymlinks(b)
	return errA == nil && errB == nil && ra == rb
}

// git runs git in the repo root and returns trimmed stdout.
func (r *repo) git(args ...string) (string, error) {
	c := exec.Command("git", args...)
	c.Dir = r.Root
	c.Env = r.gitEnv()
	var stderr bytes.Buffer
	c.Stderr = &stderr
	out, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, stderr.String())
	}
	return strings.TrimSpace(string(out)), nil
}

// author is git's user.email, as linespec records it.
func (r *repo) author() string {
	if email, err := r.git("config", "user.email"); err == nil && email != "" {
		return email
	}
	return "unknown@example.com"
}

// commit stages paths and commits them, letting hooks print to the terminal.
func (r *repo) commit(message string, paths ...string) error {
	if _, err := r.git(append([]string{"add", "--"}, paths...)...); err != nil {
		return err
	}
	c := exec.Command("git", "commit", "-m", message)
	c.Dir = r.Root
	c.Env = r.gitEnv()
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}
	return nil
}

// unstage undoes commit's `git add` after a rejected commit.
func (r *repo) unstage(paths ...string) {
	if _, err := r.git(append([]string{"reset", "--quiet", "HEAD", "--"}, paths...)...); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to unstage during rollback: %v\n", err)
	}
}

// fileState is a file's bytes before a transition, for rollback. A missing
// file is restored by removing it.
type fileState struct {
	path   string
	data   []byte
	exists bool
}

func snapshotFiles(paths ...string) []fileState {
	out := make([]fileState, 0, len(paths))
	for _, p := range paths {
		data, err := os.ReadFile(p)
		out = append(out, fileState{path: p, data: data, exists: err == nil})
	}
	return out
}

func restoreFiles(states []fileState) {
	for _, s := range states {
		if s.exists {
			_ = os.WriteFile(s.path, s.data, 0o644)
		} else {
			_ = os.RemoveAll(s.path)
		}
	}
}

// fail prints an error the way linespec's formatter does.
func (r *repo) fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "\n✗ "+format+"\n\n", args...)
}

// tag is the commit tag for a record version.
func tag(id string, version int) string {
	return fmt.Sprintf("[%s@%d]", id, version)
}
