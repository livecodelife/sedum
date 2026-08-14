package cli

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/calebcowen/sedum/internal/inject"
	"github.com/calebcowen/sedum/internal/pipeline"
	"github.com/calebcowen/sedum/internal/runlog"
	"github.com/calebcowen/sedum/internal/selection"
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
			return runGrow(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), cfg)
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

// runGrow runs the pipeline.
//
// The model client is built before any phase runs, so a run that cannot reach a
// model fails while nothing has been written. Creating files and then failing
// at the model call would leave the output tree half built for a reason the
// user could have been told first.
func runGrow(ctx context.Context, out, errOut io.Writer, cfg GrowConfig) error {
	if cfg.Replaying() {
		return notImplemented("grow --execute", "M7", "recording and replay")
	}

	// No --stop-after means run every phase. Zero is not a stop point, so it
	// reads as "stop after nothing".
	stopAfter := 0
	name := ""
	if cfg.StopAfter != "" {
		sp, ok := lookupStopPoint(cfg.StopAfter)
		if !ok {
			return fmt.Errorf("--stop-after %q is not a phase boundary; expected one of: %s", cfg.StopAfter, stopPointNames())
		}
		stopAfter, name = sp.afterPhase, sp.name
	}

	// A run halting before Phase 4 never consults a model, so it never needs
	// one configured. That is what makes --stop-after resolution and
	// --stop-after files usable with no endpoint at all.
	var client selection.Client
	if stopAfter == 0 || stopAfter >= pipeline.PhaseSelect {
		built, err := selection.NewOpenAI(cfg.Model)
		if err != nil {
			return err
		}
		client = built
	}

	log, err := runlog.New(cfg.LogPath, cfg.Verbose)
	if err != nil {
		return err
	}
	defer log.Close()

	result, err := pipeline.Run(ctx, pipeline.Config{
		Generators:     cfg.Generators,
		Records:        cfg.Records,
		Output:         cfg.Output,
		Lang:           cfg.Lang,
		Only:           cfg.Only,
		DryRun:         cfg.DryRun,
		StopAfterPhase: stopAfter,
		Client:         client,
		Retries:        cfg.Retries,
		Log:            log,
	})
	if err != nil {
		return err
	}

	printWarnings(errOut, result.Warnings)
	report(out, result, stopAfter, cfg.DryRun)
	if name != "" {
		fmt.Fprintf(out, "\nstopped after %s\n", name)
	}
	return nil
}

// report prints what the run got as far as, at the detail the stop point makes
// useful. A run halted at resolution has decided nothing to show but
// resolutions; one that finished has files, selections, and injections, and the
// interesting part is the last of those.
func report(out io.Writer, result *pipeline.Result, stopAfter int, dryRun bool) {
	if stopAfter == pipeline.PhaseResolve {
		printResolutions(out, result, nil, false)
		return
	}

	printFiles(out, result, dryRun)
	if stopAfter == pipeline.PhaseCreate {
		return
	}

	printSelections(out, result)
	if stopAfter == pipeline.PhaseValidate || stopAfter == pipeline.PhaseExpand {
		return
	}
	printInjections(out, result, dryRun)
}

// printSelections reports what the model chose, grouped by record because a
// record is one model call and the unit a reader is reasoning about.
func printSelections(out io.Writer, result *pipeline.Result) {
	fmt.Fprintln(out)
	var total int
	for _, s := range result.Selections {
		fmt.Fprintf(out, "%s\n", s.RecordID)
		if len(s.Invocations) == 0 {
			fmt.Fprintf(out, "  (no action selected)\n")
			continue
		}
		for _, inv := range s.Invocations {
			fmt.Fprintf(out, "  %s %s\n", inv.Action, describeKwargs(inv.Kwargs))
			total++
		}
	}
	fmt.Fprintf(out, "\n%d invocation(s) selected across %d record(s)\n", total, len(result.Selections))
}

func printInjections(out io.Writer, result *pipeline.Result, dryRun bool) {
	verb := "injected"
	if dryRun {
		verb = "would inject"
	}

	fmt.Fprintln(out)
	var replaced, skipped int
	for _, r := range result.Injections {
		switch {
		case r.Skipped:
			skipped++
			fmt.Fprintf(out, "seeded   %s  %s (left as it is)\n", r.Path, describeRegion(r))
		case r.Replaced:
			replaced++
			fmt.Fprintf(out, "replaced %s  %s\n", r.Path, describeRegion(r))
		default:
			fmt.Fprintf(out, "%s %s  %s\n", verb, r.Path, describeRegion(r))
		}
	}
	fmt.Fprintf(out, "\n%d region(s) %s, %d replaced, %d left as seeded\n",
		len(result.Injections), verb, replaced, skipped)
}

func describeRegion(r inject.Result) string {
	if r.Variant == "" {
		return r.Action
	}
	return r.Action + ":" + r.Variant
}

// describeKwargs renders bound arguments in name order, so that one run's
// report reads the same way as another's.
func describeKwargs(kwargs map[string]any) string {
	names := make([]string, 0, len(kwargs))
	for name := range kwargs {
		names = append(names, name)
	}
	sort.Strings(names)

	pairs := make([]string, 0, len(names))
	for _, name := range names {
		pairs = append(pairs, fmt.Sprintf("%s=%v", name, kwargs[name]))
	}
	return strings.Join(pairs, " ")
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
