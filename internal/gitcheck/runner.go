package gitcheck

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
)

// ResolveBinary picks the provenance binary: the --provenance value when given,
// else the PROVENANCE environment value when not empty, else the first
// `provenance` on PATH. An empty string is not given. A binary whose path, once
// symlinks are resolved, lies under the git top level is refused, so a hook
// never runs whatever is checked in under that name.
func ResolveBinary(flag, env, topLevel string) (string, error) {
	binary := flag
	if binary == "" {
		binary = env
	}
	if binary == "" {
		found, err := exec.LookPath("provenance")
		if err != nil {
			return "", ProvenanceError("provenance not found: set --provenance or PROVENANCE, or put provenance on PATH")
		}
		binary = found
	}
	abs, err := filepath.Abs(binary)
	if err != nil {
		return "", ProvenanceError("cannot resolve " + binary + ": " + err.Error())
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", ProvenanceError("cannot resolve " + binary + ": " + err.Error())
	}
	top, err := filepath.EvalSymlinks(topLevel)
	if err != nil {
		top = topLevel
	}
	if rel, err := filepath.Rel(top, resolved); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", ProvenanceError("refusing " + binary + ": it lies inside the repository")
	}
	return binary, nil
}

// CheckArgs builds the provenance arguments for one commit, or for the staged
// change when staged: check first, the commit or --staged, the message file
// (staged only), the config when given, and --format sarif last.
func CheckArgs(sha string, staged bool, messageFile, config string) []string {
	args := []string{"check"}
	if staged {
		args = append(args, "--staged")
		if messageFile != "" {
			args = append(args, "--message-file", messageFile)
		}
	} else {
		args = append(args, "--commit", sha)
	}
	if config != "" {
		args = append(args, "--config", config)
	}
	return append(args, "--format", "sarif")
}

// Environment returns the environment for git and provenance: every GIT_*
// variable is cleared, except that with staged GIT_DIR, GIT_WORK_TREE and
// GIT_INDEX_FILE pass through unchanged so that the index being committed is read.
func Environment(staged bool, current []string) []string {
	var out []string
	for _, kv := range current {
		if strings.HasPrefix(kv, "GIT_") {
			name, _, _ := strings.Cut(kv, "=")
			if !staged || (name != "GIT_DIR" && name != "GIT_WORK_TREE" && name != "GIT_INDEX_FILE") {
				continue
			}
		}
		out = append(out, kv)
	}
	return out
}

// RunProvenance runs the binary in the working directory with the environment
// and returns what it printed. Exit status 0 or 1 with output is a result (1 is
// provenance reporting findings); empty output, a higher status, or a binary
// that cannot start is a provenance error.
func RunProvenance(binary string, args []string, workDir string, env []string) ([]byte, error) {
	cmd := exec.Command(binary, args...)
	cmd.Dir = workDir
	cmd.Env = env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		exit, ok := err.(*exec.ExitError)
		if !ok {
			return nil, ProvenanceError("cannot run " + binary + ": " + err.Error())
		}
		if exit.ExitCode() > 1 {
			msg := strings.TrimSpace(stderr.String())
			if msg != "" {
				msg = ": " + msg
			}
			return nil, ProvenanceError(binary + " failed: " + exit.Error() + msg)
		}
	}
	if len(bytes.TrimSpace(stdout.Bytes())) == 0 {
		return nil, ProvenanceError(binary + " printed nothing")
	}
	return stdout.Bytes(), nil
}
