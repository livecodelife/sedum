package gitcheck

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/livecodelife/sedum/internal/pipeline"
)

// ExportTree exports the whole tree of the git top level of baseDir into a
// temporary directory and returns its root and a function that removes it. For a
// revision it is that commit's tree; for an empty revision with staged true it is
// the index, written out without writing any git object.
func ExportTree(baseDir, rev string, staged bool, env []string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "sedum-check-")
	if err != nil {
		return "", nil, GitError("cannot create a temporary directory: " + err.Error())
	}
	cleanup := func() { os.RemoveAll(dir) }
	top, err := TopLevel(baseDir, env)
	if err != nil {
		cleanup()
		return "", nil, err
	}
	var cmd *exec.Cmd
	if staged && rev == "" {
		cmd = exec.Command("git", "-C", top, "checkout-index", "-a", "--prefix="+dir+string(os.PathSeparator))
	} else {
		cmd = exec.Command("sh", "-c", `git -C "$1" archive "$2" | tar -x -C "$3"`, "sh", top, rev, dir)
	}
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil {
		cleanup()
		return "", nil, GitError("cannot export the tree: " + strings.TrimSpace(string(out)))
	}
	return dir, cleanup, nil
}

// classOf maps a pipeline mode to the class word `sedum grow --check` prints.
func classOf(mode pipeline.Mode) string {
	switch mode {
	case pipeline.ModeNone:
		return "none"
	case pipeline.ModeSelect:
		return "unknown"
	case pipeline.ModeReselect:
		return "stale"
	case pipeline.ModeSkip:
		return "current"
	case pipeline.ModeRerender:
		return "re-render"
	case pipeline.ModeEdited:
		return "edited"
	}
	return "unknown"
}

// Classes computes the class of each record on the tree in-process, with the
// same pipeline as `sedum grow --check` and no model. A scope path no generator
// package claims is skipped and returned, each distinct path once. A generators
// directory that does not exist means no packages.
func Classes(tree, generators, records string) (map[string]string, []string, error) {
	if info, err := os.Stat(generators); err != nil || !info.IsDir() {
		empty, err := os.MkdirTemp("", "sedum-nogenerators-")
		if err != nil {
			return nil, nil, ProvenanceError("cannot create a temporary directory: " + err.Error())
		}
		defer os.RemoveAll(empty)
		generators = empty
	}
	result, err := pipeline.Run(context.Background(), pipeline.Config{
		Check:              true,
		SkipUnclaimedScope: true,
		Generators:         generators,
		Records:            records,
		Output:             filepath.Clean(tree),
	})
	if err != nil {
		return nil, nil, ProvenanceError("cannot compute the classes: " + err.Error())
	}
	classes := map[string]string{}
	for _, c := range result.Classifications {
		classes[c.RecordID] = classOf(c.Mode)
	}
	return classes, result.SkippedScopePaths, nil
}
