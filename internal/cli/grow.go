package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/calebcowen/sedum/internal/pipeline"
	"github.com/calebcowen/sedum/internal/runlog"
)

func newGrowCommand() *cobra.Command {
	var cfg GrowConfig

	cmd := &cobra.Command{
		Use:   "grow",
		Short: "Run the full pipeline: load, ingest, resolve, create, invoke, validate, expand, inject",
		Long: `Runs every phase in order: load and validate generator packages, ingest
provenance records, resolve paths to packages, create files from templates,
invoke the model, validate its output, expand and resolve, and inject.

Nothing is created that a provenance record did not authorize.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cfg.Validate(); err != nil {
				return err
			}
			for _, ignored := range cfg.IgnoredFlags() {
				fmt.Fprintf(cmd.ErrOrStderr(), "sedum: ignoring %s\n", ignored)
			}
			return runGrow(cmd.OutOrStdout(), cmd.ErrOrStderr(), cfg)
		},
	}

	f := cmd.Flags()
	f.StringVar(&cfg.Generators, "generators", "", "Generators directory. Required.")
	f.StringVar(&cfg.Records, "records", "", "Provenance records directory. Required unless --execute is given.")
	f.StringVar(&cfg.Output, "output", defaultOutput, "Output directory.")
	f.StringSliceVar(&cfg.Lang, "lang", nil, "Prefer the named package where an extension is contested. Repeatable.")
	f.StringSliceVar(&cfg.Only, "only", nil, "Generate only the named provenance record. Repeatable.")
	f.StringVar(&cfg.RecordTo, "record", "", "Write a recording of the run to the given path.")
	f.StringVar(&cfg.Execute, "execute", "", "Replay a recording. Skips model invocation.")
	f.BoolVar(&cfg.DryRun, "dry-run", false, "Run every phase, write nothing.")
	f.StringVar(&cfg.StopAfter, "stop-after", "", "Halt after the named phase. One of: "+stopPointNames()+".")
	f.IntVar(&cfg.Retries, "retries", defaultRetries, "Model output validation retry limit. Ignored with --execute.")
	f.StringVar(&cfg.Model, "model", "", "Model identifier. Endpoint and credentials come from environment. Ignored with --execute.")
	f.StringVar(&cfg.LogPath, "log", defaultLogPath, "Run log location.")
	f.BoolVarP(&cfg.Verbose, "verbose", "v", false, "Mirror the run log to stdout.")

	mustMarkRequired(cmd, "generators")

	// A recording is either being written or being replayed, never both.
	cmd.MarkFlagsMutuallyExclusive("record", "execute")

	return cmd
}

// runGrow runs as much of the pipeline as has landed.
//
// A run that would reach an unimplemented phase is refused before anything is
// written. Creating files and then failing at the model call would leave the
// output tree half built for a reason the user could have been told first.
func runGrow(out, errOut io.Writer, cfg GrowConfig) error {
	if cfg.Replaying() {
		return notImplemented("grow --execute", "M7", "recording and replay")
	}

	sp, ok := lookupStopPoint(cfg.StopAfter)
	if !ok || sp.afterPhase > pipeline.PhaseCreate {
		return notImplemented("grow", "M6", "model invocation and output validation")
	}

	log, err := runlog.New(cfg.LogPath, cfg.Verbose)
	if err != nil {
		return err
	}
	defer log.Close()

	result, err := pipeline.Run(pipeline.Config{
		Generators:     cfg.Generators,
		Records:        cfg.Records,
		Output:         cfg.Output,
		Lang:           cfg.Lang,
		Only:           cfg.Only,
		DryRun:         cfg.DryRun,
		StopAfterPhase: sp.afterPhase,
		Log:            log,
	})
	if err != nil {
		return err
	}

	printWarnings(errOut, result.Warnings)
	if sp.afterPhase == pipeline.PhaseResolve {
		printResolutions(out, result, nil, false)
	} else {
		printFiles(out, result, cfg.DryRun)
	}
	fmt.Fprintf(out, "\nstopped after %s\n", sp.name)
	return nil
}

// printFiles reports what Phase 3 did, separating the files it created from the
// ones it found already there. Create-if-absent means the second group is
// normal rather than exceptional, so it is reported rather than hidden.
func printFiles(out io.Writer, result *pipeline.Result, dryRun bool) {
	verb := "created"
	if dryRun {
		verb = "would create"
	}

	var existing int
	for _, f := range result.Files {
		if f.Existed {
			existing++
			continue
		}
		fmt.Fprintf(out, "%s  %s  (%s %s)\n", verb, f.Path, f.Package.Name, describeTemplate(f.Resolution))
	}
	for _, f := range result.Files {
		if f.Existed {
			fmt.Fprintf(out, "exists   %s  (left as it is)\n", f.Path)
		}
	}

	fmt.Fprintf(out, "\n%d path(s) authorized, %d %s, %d already present\n",
		len(result.Files), len(result.Files)-existing, verb, existing)
}
