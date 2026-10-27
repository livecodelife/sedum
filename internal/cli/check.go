package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/livecodelife/sedum/internal/gitcheck"
)

type checkFlags struct {
	rng         string
	staged      bool
	messageFile string
	output      string
	config      string
	generators  string
	records     string
	format      string
	provenance  string
}

func newCheckCommand() *cobra.Command {
	var f checkFlags

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check commits with provenance, waiving only proven regenerations and marker migrations",
		Long: `Runs provenance check on each commit of --range, or on the staged change, and
waives only the findings that sedum can prove come from a pure regeneration or a
pure marker migration. It complements provenance check: provenance reports every
finding, and this command keeps the ones it cannot prove are its own.

The classes of the records are computed in-process, so a failure to compute them
leaves the findings in place and exits 0 or 1 by what remains.

Exit 0 means no finding is left, exit 1 that some are, and exit 2 that provenance
or git could not be run or read, that provenance printed output the command
cannot use, or that provenance output lacks fields the command needs (skew).`,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 {
				return usageExit("check takes no arguments")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheck(cmd.OutOrStdout(), f)
		},
	}
	cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return usageExit(err.Error()) })
	cmd.Flags().StringVar(&f.rng, "range", "", "Check every commit in a..b, oldest first")
	cmd.Flags().BoolVar(&f.staged, "staged", false, "Check the staged change")
	cmd.Flags().StringVar(&f.messageFile, "message-file", "", "The commit message file (with --staged)")
	cmd.Flags().StringVar(&f.output, "output", ".", "The directory that holds the provenance config")
	cmd.Flags().StringVar(&f.config, "config", "", "The provenance config file")
	cmd.Flags().StringVar(&f.generators, "generators", "", "The generator packages (default: generators under the base directory)")
	cmd.Flags().StringVar(&f.records, "records", "", "The records directory (default: the config's dir)")
	cmd.Flags().StringVar(&f.format, "format", "human", "The output format: sarif or human")
	cmd.Flags().StringVar(&f.provenance, "provenance", "", "The provenance binary (default: $PROVENANCE, else provenance on PATH)")
	return cmd
}

// usageExit is the error for a command line the command cannot use: exit 2.
func usageExit(message string) error {
	return ExitError{Code: 2, Message: gitcheck.FormatError(gitcheck.UsageError(message))}
}

func runCheck(out io.Writer, f checkFlags) error {
	switch {
	case f.rng != "" && f.staged:
		return usageExit("--range and --staged cannot be used together")
	case f.rng == "" && !f.staged:
		return usageExit("pass --range <a..b> or --staged")
	case f.messageFile != "" && !f.staged:
		return usageExit("--message-file is accepted only with --staged")
	case f.format != "sarif" && f.format != "human":
		return usageExit(fmt.Sprintf("unknown format %q (sarif|human)", f.format))
	}

	output, err := filepath.Abs(f.output)
	if err != nil {
		return usageExit(err.Error())
	}
	base, config := output, ""
	if f.config != "" {
		config = f.config
		if !filepath.IsAbs(config) {
			config = filepath.Join(output, config)
		}
		base = filepath.Dir(config)
	}

	env := gitcheck.Environment(f.staged, os.Environ())
	top, err := gitcheck.TopLevel(base, env)
	if err != nil {
		return ExitError{Code: 2, Message: gitcheck.FormatError(err)}
	}
	outputTop, err := gitcheck.TopLevel(output, env)
	if err != nil {
		return ExitError{Code: 2, Message: gitcheck.FormatError(err)}
	}
	if !sameDir(top, outputTop) {
		return usageExit("--output must lie in the same git repository as the base directory " + base)
	}

	recordsDir, excludes, err := readConfig(base, config)
	if err != nil {
		return ExitError{Code: 2, Message: gitcheck.FormatError(err)}
	}
	generators := f.generators
	if generators == "" {
		generators = "generators"
	}
	records := f.records
	if records == "" {
		records = recordsDir
	}

	outcome, err := gitcheck.Check(gitcheck.Options{
		Range:       f.rng,
		MessageFile: f.messageFile,
		Config:      config,
		Output:      base,
		Generators:  generators,
		Records:     records,
		Format:      f.format,
		Provenance:  f.provenance,
		TopLevel:    top,
		Staged:      f.staged,
		Paths:       gitcheck.PathConfig{RecordsDir: recordsDir, ExcludePaths: excludes},
	}, out)
	if err != nil {
		return ExitError{Code: 2, Message: gitcheck.FormatError(err)}
	}
	if outcome.Unwaived > 0 {
		return ExitError{Code: 1}
	}
	return nil
}

// sameDir reports whether two paths name the same directory once symlinks are
// resolved.
func sameDir(a, b string) bool {
	ra, errA := filepath.EvalSymlinks(a)
	rb, errB := filepath.EvalSymlinks(b)
	if errA != nil || errB != nil {
		return a == b
	}
	return ra == rb
}

type provenanceConfig struct {
	Dir          string   `yaml:"dir"`
	ExcludePaths []string `yaml:"exclude_paths"`
}

// readConfig returns the records directory and the exclude_paths of the
// provenance config: the file named by config when given, else .provenance.yml
// in base, else the provenance block of .linespec.yml in base. With no config the
// records directory is "provenance" and nothing is excluded.
func readConfig(base, config string) (string, []string, error) {
	path, block := config, false
	if path != "" {
		block = filepath.Base(path) == ".linespec.yml"
	} else if _, err := os.Stat(filepath.Join(base, ".provenance.yml")); err == nil {
		path = filepath.Join(base, ".provenance.yml")
	} else if _, err := os.Stat(filepath.Join(base, ".linespec.yml")); err == nil {
		path, block = filepath.Join(base, ".linespec.yml"), true
	}
	if path == "" {
		return "provenance", nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, gitcheck.ProvenanceError("cannot read the config " + path + ": " + err.Error())
	}
	var cfg provenanceConfig
	if block {
		var doc struct {
			Provenance provenanceConfig `yaml:"provenance"`
		}
		err = yaml.Unmarshal(data, &doc)
		cfg = doc.Provenance
	} else {
		err = yaml.Unmarshal(data, &cfg)
	}
	if err != nil {
		return "", nil, gitcheck.ProvenanceError("cannot parse the config " + path + ": " + err.Error())
	}
	dir := cfg.Dir
	for len(dir) > 0 && dir[len(dir)-1] == '/' {
		dir = dir[:len(dir)-1]
	}
	if dir == "" {
		dir = "provenance"
	}
	return dir, cfg.ExcludePaths, nil
}
