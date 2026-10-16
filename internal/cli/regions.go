package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/livecodelife/sedum/internal/pipeline"
)

// RegionsConfig is `sedum regions`'s flags.
type RegionsConfig struct {
	Generators string
	Output     string
	Records    string
	Only       []string
	JSON       bool
}

func newRegionsCommand() *cobra.Command {
	var cfg RegionsConfig

	cmd := &cobra.Command{
		Use:   "regions",
		Short: "List every region whose marker names a record, with its gen and pkg hashes",
		Long: `Reads the markers under --output and lists every region that names a record:
its file, action, identifying kwargs, tier, and the gen and pkg hashes it was
written with. An empty hash means the marker has none, which a consumer must
read as unknown; sedum grow --stamp fills them in.

Read-only: no model, nothing written, no records directory required. Hidden
directories, the generators directory and the records directory are never read.
With --json it prints the sedum.regions/1 document described in
conformance/formats, which the provenance library's evidence code hash reads.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runRegions(cmd.OutOrStdout(), cfg)
		},
	}

	f := cmd.Flags()
	f.StringVar(&cfg.Generators, "generators", "", "Generators directory. Required.")
	f.StringVar(&cfg.Output, "output", defaultOutput, "Output directory to read markers from.")
	f.StringVar(&cfg.Records, "records", "", "Provenance records directory. Optional: given, a marker naming an alias is reported under the surviving id.")
	f.StringArrayVar(&cfg.Only, "record", nil, "List only the named record's regions, by id or alias. Repeatable.")
	f.BoolVar(&cfg.JSON, "json", false, "Print the sedum.regions/1 JSON document.")

	mustMarkRequired(cmd, "generators")
	return cmd
}

func runRegions(out io.Writer, cfg RegionsConfig) error {
	report, err := pipeline.Regions(pipeline.RegionsConfig{
		Generators: cfg.Generators,
		Output:     cfg.Output,
		Records:    cfg.Records,
		Only:       cfg.Only,
	})
	if err != nil {
		return err
	}
	if cfg.JSON {
		return writeJSON(out, report)
	}
	for _, r := range report.Regions {
		identity, _ := json.Marshal(r.Identity)
		fmt.Fprintf(out, "%s  %s  %s %s  %s  gen=%s pkg=%s\n", r.Record, r.File,
			describeRef(pipeline.RegionRef{Action: r.Action, Variant: r.Variant}), identity, r.Tier,
			orUnknown(r.Gen), orUnknown(r.Pkg))
	}
	fmt.Fprintf(out, "\n%d region(s)\n", len(report.Regions))
	return nil
}

func orUnknown(h string) string {
	if h == "" {
		return "(unknown)"
	}
	return h
}

// writeJSON prints one indented JSON document, without HTML escaping.
func writeJSON(out io.Writer, v any) error {
	encoder := json.NewEncoder(out)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}
