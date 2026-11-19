package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/livecodelife/sedum/internal/genpkg"
)

func newConvergenceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convergence",
		Short: "Render a package's convergence contract (stubs, helpers, seeds)",
		Long: `A package's optional convergence: block says how its language writes a stub
that delegates to a helper, and how a file of placeholder helpers looks. These
commands apply it. Sedum's core knows no language; it renders the package's own
templates.`,
	}
	cmd.AddCommand(newConvergenceRenderCommand(), newConvergenceInfoCommand())
	return cmd
}

func loadConvergencePackage(generators, name string) (*genpkg.Package, error) {
	set, findings, err := genpkg.Load(generators, genpkg.Options{Only: []string{name}})
	if err != nil {
		return nil, err
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			return nil, fmt.Errorf("package %s did not load: %s", name, f)
		}
	}
	pkg, ok := set.Lookup(name)
	if !ok {
		return nil, fmt.Errorf("no package named %q in %s; it declares %s", name, generators, packageNames(set))
	}
	if pkg.Convergence == nil {
		return nil, fmt.Errorf("package %s declares no convergence block", name)
	}
	return pkg, nil
}

func newConvergenceRenderCommand() *cobra.Command {
	var generators, pkgName string
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render stubs, helper signatures and seeds from JSON on stdin",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			pkg, err := loadConvergencePackage(generators, pkgName)
			if err != nil {
				return err
			}
			var in genpkg.ConvRenderInput
			dec := json.NewDecoder(cmd.InOrStdin())
			dec.DisallowUnknownFields()
			if err := dec.Decode(&in); err != nil {
				return fmt.Errorf("stdin is not a convergence render input: %w", err)
			}
			if err := dec.Decode(new(json.RawMessage)); err != io.EOF {
				return fmt.Errorf("stdin carries more than one JSON value")
			}
			out, err := pkg.RenderConvergence(in)
			if err != nil {
				return err
			}
			return writeJSON(cmd.OutOrStdout(), out)
		},
	}
	f := cmd.Flags()
	f.StringVar(&generators, "generators", "", "Generators directory. Required.")
	f.StringVar(&pkgName, "package", "", "Package declaring the convergence block. Required.")
	mustMarkRequired(cmd, "generators", "package")
	return cmd
}

func newConvergenceInfoCommand() *cobra.Command {
	var generators, pkgName string
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Print the package's convergence facts as JSON",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			pkg, err := loadConvergencePackage(generators, pkgName)
			if err != nil {
				return err
			}
			info, err := pkg.ConvergenceInfo()
			if err != nil {
				return err
			}
			return writeJSON(cmd.OutOrStdout(), info)
		},
	}
	f := cmd.Flags()
	f.StringVar(&generators, "generators", "", "Generators directory. Required.")
	f.StringVar(&pkgName, "package", "", "Package declaring the convergence block. Required.")
	mustMarkRequired(cmd, "generators", "package")
	return cmd
}
