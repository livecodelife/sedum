package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newValidateCommand() *cobra.Command {
	var cfg ValidateConfig

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Run every load-time check against the generators directory and exit",
		Long: `Loads every generator package and reports all errors and warnings found.

Packages are wholly valid or rejected; there are no partial loads. Requires no
records, no model, and no network.`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return notImplemented("validate", "M1", "generator package loading and load-time validation")
		},
	}

	f := cmd.Flags()
	f.StringVar(&cfg.Generators, "generators", "", "Generators directory. Required.")
	f.StringSliceVar(&cfg.Packages, "package", nil, "Validate a single package. Repeatable.")
	f.BoolVar(&cfg.Strict, "strict", false, "Treat warnings as errors.")

	mustMarkRequired(cmd, "generators")

	return cmd
}

// mustMarkRequired panics if the flag does not exist. A typo here would
// silently drop a required-flag check, so it fails at construction rather than
// letting a run proceed without its inputs.
func mustMarkRequired(cmd *cobra.Command, names ...string) {
	for _, name := range names {
		if err := cmd.MarkFlagRequired(name); err != nil {
			panic(fmt.Sprintf("cli: marking --%s required on %q: %v", name, cmd.Name(), err))
		}
	}
}
