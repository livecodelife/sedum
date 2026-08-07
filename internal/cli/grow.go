package cli

import (
	"fmt"

	"github.com/spf13/cobra"
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
			return notImplemented("grow", "M6", "model invocation and output validation")
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
