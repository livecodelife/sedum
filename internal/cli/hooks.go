package cli

import (
	"fmt"
	"io"
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/livecodelife/sedum/internal/scaffold"
)

func newHooksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Manage the sedum git hooks",
	}
	cmd.AddCommand(newHooksInstallCommand())
	return cmd
}

func newHooksInstallCommand() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the commit-msg and pre-push hooks",
		Long: `Writes sedum's commit-msg and pre-push hooks into the hooks directory of the git
repository that contains --dir. A different hook already there is kept as
<name>.bak the first time. The commit-msg hook runs sedum check on the staged
change; the pre-push hook runs provenance's own pre-push hook and then sedum
check on each pushed ref's commits.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runHooksInstall(cmd.OutOrStdout(), dir)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "A directory inside the git repository")
	return cmd
}

func runHooksInstall(out io.Writer, dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if info, err := os.Stat(abs); err != nil || !info.IsDir() {
		return fmt.Errorf("--dir %q is not an existing directory", dir)
	}
	hooksDir, err := findHooksDir(abs)
	if err != nil {
		return err
	}
	hooks := scaffold.Hooks()
	results := make([]string, len(hooks))
	for i, hook := range hooks {
		result, err := scaffold.InstallHook(hooksDir, hook.Name, hook.Body)
		if err != nil {
			return err
		}
		results[i] = result
	}
	fmt.Fprintln(out, "hooks directory: "+hooksDir)
	for i, hook := range hooks {
		fmt.Fprintln(out, hook.Name+": "+results[i])
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "note: a later provenance hooks install replaces these two hooks again, so run sedum hooks install after it")
	fmt.Fprintln(out, "note: the hooks run $SEDUM else sedum on PATH, and $PROVENANCE else provenance on PATH, never a file in the repository")
	return nil
}

// findHooksDir returns the hooks directory of the git repository containing
// dir, found from the top level so that a subdirectory resolves the same way.
func findHooksDir(dir string) (string, error) {
	var env []string
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "GIT_") {
			env = append(env, kv)
		}
	}
	notInside := fmt.Errorf("%s is not inside a git repository", dir)
	top := osexec.Command("git", "rev-parse", "--show-toplevel")
	top.Dir, top.Env = dir, env
	topOut, err := top.Output()
	if err != nil {
		return "", notInside
	}
	topLevel := strings.TrimSpace(string(topOut))
	hooks := osexec.Command("git", "rev-parse", "--git-path", "hooks")
	hooks.Dir, hooks.Env = topLevel, env
	hooksOut, err := hooks.Output()
	if err != nil {
		return "", notInside
	}
	path := strings.TrimSpace(string(hooksOut))
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Join(topLevel, path), nil
}
