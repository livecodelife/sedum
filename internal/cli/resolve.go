package cli

import "github.com/spf13/cobra"

func newResolveCommand() *cobra.Command {
	var cfg ResolveConfig

	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Report what each authorized path resolves to, without invoking the model",
		Long: `Runs Phases 0 through 3 and reports, for every authorized path, the generator
package it resolved to, the file template that matched, and the captures bound.

The primary debugging tool for package resolution and template specificity.`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return notImplemented("resolve", "M3", "path resolution and file creation from file templates")
		},
	}

	f := cmd.Flags()
	f.StringVar(&cfg.Generators, "generators", "", "Generators directory. Required.")
	f.StringVar(&cfg.Records, "records", "", "Provenance records directory. Required.")
	f.StringSliceVar(&cfg.Lang, "lang", nil, "Prefer the named package where an extension is contested. Repeatable.")
	f.StringSliceVar(&cfg.Only, "only", nil, "Resolve only the named provenance record. Repeatable.")
	f.BoolVar(&cfg.ShowTemplate, "show-template", false, "Include rendered template output for each path.")

	mustMarkRequired(cmd, "generators", "records")

	return cmd
}
