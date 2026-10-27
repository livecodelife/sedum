package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/livecodelife/sedum/internal/scaffold"
)

func newInitCommand() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Prepare a directory for sedum grow",
		Long: `Creates .sedum/context.md (the advisory context file) when it is absent, makes sure
.gitignore ignores .sedum/run.log, and adds generators/** and .sedum/** to the
exclude_paths of .provenance.yml when that file exists. Safe to run again: it
never overwrites a context file and changes nothing that is already in place.

Prints one line per item: created, updated, unchanged or skipped.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runInit(cmd.OutOrStdout(), dir)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "The directory to prepare")
	return cmd
}

func runInit(out io.Writer, dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if info, err := os.Stat(abs); err != nil || !info.IsDir() {
		return fmt.Errorf("--dir %q is not an existing directory", dir)
	}
	items, err := scaffold.Init(abs)
	if err != nil {
		return err
	}
	for _, item := range items {
		line := fmt.Sprintf("%-10s %s", item.Status, item.Name)
		if item.Detail != "" {
			line += " (" + item.Detail + ")"
		}
		fmt.Fprintln(out, line)
	}
	return nil
}
