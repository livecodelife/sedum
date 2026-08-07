package cli

import "github.com/spf13/cobra"

func newActionsCommand() *cobra.Command {
	var cfg ActionsConfig

	cmd := &cobra.Command{
		Use:   "actions",
		Short: "Print a package's action catalog exactly as the model would receive it",
		Long: `Prints the exposed actions of a generator package with their kwarg schemas and
variant lists, in the form the model is given.

The authoring feedback loop for exposure and catalog clarity.`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return notImplemented("actions", "M6", "catalog construction shared with model invocation")
		},
	}

	f := cmd.Flags()
	f.StringVar(&cfg.Generators, "generators", "", "Generators directory. Required.")
	f.StringVar(&cfg.Package, "package", "", "Package to inspect. Required.")
	f.BoolVar(&cfg.All, "all", false, "Include unexposed actions, marked as such.")
	f.BoolVar(&cfg.JSON, "json", false, "Emit the raw catalog payload rather than formatted output.")

	mustMarkRequired(cmd, "generators", "package")

	return cmd
}
