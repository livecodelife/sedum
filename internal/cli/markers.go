package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/livecodelife/sedum/internal/pipeline"
)

// newMarkersCommand builds the markers command group (prov-2026-fa3497af).
func newMarkersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "markers",
		Short: "Work on ownership markers",
	}
	cmd.AddCommand(newMarkersMigrateCommand())
	return cmd
}

// newMarkersMigrateCommand builds the markers migrate command
// (prov-2026-f6643ae2): it rewrites v1 markers to v2 and writes the recordings
// of the records they name, with no model and no change to a region's text.
func newMarkersMigrateCommand() *cobra.Command {
	var (
		generators, records, output string
		dryRun                      bool
		vars                        []string
	)

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Rewrite v1 markers as v2 markers and write the recordings of their records",
		Long: `Rewrites each v1 marker under --output to a v2 marker and writes the
recording of every record the markers name, with no model call. A region's text
is never changed, and a migrated marker keeps the gen and pkg of the marker it
replaces, so grow --check reports the same classes after migration as before.

A region that cannot be migrated is refused and keeps its v1 marker; a region
whose marker names no known record is left as it is. Both are reported, and
refusals make the command exit non-zero. Migrating twice changes nothing.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			variables, err := parseVars(vars)
			if err != nil {
				return err
			}
			report, err := pipeline.Migrate(pipeline.MigrateConfig{
				Generators: generators,
				Records:    records,
				Output:     output,
				Variables:  variables,
				DryRun:     dryRun,
			})
			if err != nil {
				return err
			}
			return printMigration(cmd.OutOrStdout(), report, dryRun)
		},
	}

	f := cmd.Flags()
	f.StringVar(&generators, "generators", "", "Generators directory. Required.")
	f.StringVar(&records, "records", "", "Provenance records directory. Required.")
	f.StringVar(&output, "output", defaultOutput, "Output directory.")
	f.BoolVar(&dryRun, "dry-run", false, "Report what would be migrated, write nothing.")
	f.StringArrayVar(&vars, "var", nil, "Bind a variable a generator package declares, as name=value. Repeatable.")

	mustMarkRequired(cmd, "generators")
	mustMarkRequired(cmd, "records")
	return cmd
}

// printMigration prints a migration's report - migrated, then refused, then
// left - and returns an error when any region was refused.
func printMigration(out io.Writer, report *pipeline.MigrateReport, dryRun bool) error {
	verb := "migrated"
	if dryRun {
		fmt.Fprintln(out, "dry run: nothing is written")
		verb = "would migrate"
	}
	for _, n := range report.Migrated {
		fmt.Fprintf(out, "%s %s  %s\n", verb, n.Path, label(n.Action, n.Variant))
	}
	for _, n := range report.Refused {
		fmt.Fprintf(out, "refused %s  %s: %s\n", n.Path, label(n.Action, n.Variant), n.Why)
	}
	for _, n := range report.Left {
		fmt.Fprintf(out, "left %s  %s: %s\n", n.Path, label(n.Action, n.Variant), n.Why)
	}
	if report.Skipped > 0 {
		fmt.Fprintf(out, "skipped %d scope paths no generator package claims\n", report.Skipped)
	}
	if dryRun {
		fmt.Fprintf(out, "%d regions would be migrated, %d files would be rewritten, %d recordings would be written\n",
			len(report.Migrated), len(report.Files), len(report.Recordings))
	} else {
		fmt.Fprintf(out, "%d regions migrated, %d files rewritten, %d recordings written\n",
			len(report.Migrated), len(report.Files), len(report.Recordings))
	}
	if len(report.Refused) > 0 {
		return fmt.Errorf("%d regions could not be migrated", len(report.Refused))
	}
	return nil
}
