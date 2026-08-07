package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/calebcowen/sedum/internal/genpkg"
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runValidate(cmd.OutOrStdout(), cfg)
		},
	}

	f := cmd.Flags()
	f.StringVar(&cfg.Generators, "generators", "", "Generators directory. Required.")
	f.StringSliceVar(&cfg.Packages, "package", nil, "Validate a single package. Repeatable.")
	f.BoolVar(&cfg.Strict, "strict", false, "Treat warnings as errors.")

	mustMarkRequired(cmd, "generators")

	return cmd
}

// runValidate loads every package and prints what loading found. Presentation
// lives here; the checks themselves are genpkg's.
func runValidate(out io.Writer, cfg ValidateConfig) error {
	set, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{Only: cfg.Packages})
	if err != nil {
		return err
	}
	if cfg.Strict {
		// Under --strict a warning is worth failing over, so it is
		// reported as an error rather than merely counted as one.
		findings = findings.Strict()
	}

	for _, f := range findings {
		fmt.Fprintln(out, f)
	}

	var errors, warnings int
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			errors++
			continue
		}
		warnings++
	}

	// A package genpkg rejected is already absent from the set, but --strict
	// can reject one it returned, so the count is taken against the findings
	// rather than reported as the set's length.
	out_ := map[string]bool{}
	for _, name := range rejected(findings) {
		out_[name] = true
	}
	loaded := 0
	for _, p := range set.Packages {
		if !out_[p.Name] {
			loaded++
		}
	}

	fmt.Fprintf(out, "\n%d package(s) loaded, %d error(s), %d warning(s)\n", loaded, errors, warnings)

	if errors > 0 {
		// The findings above carry the detail. This names the packages
		// that were rejected so the failure is legible on its own.
		return fmt.Errorf("rejected %s; see the findings above", strings.Join(rejected(findings), ", "))
	}
	return nil
}

// rejected names the distinct packages carrying at least one error, sorted.
func rejected(findings genpkg.Findings) []string {
	seen := map[string]bool{}
	var out []string
	for _, f := range findings {
		if f.Kind != genpkg.KindError || seen[f.Package] {
			continue
		}
		seen[f.Package] = true
		out = append(out, f.Package)
	}
	sort.Strings(out)
	return out
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
