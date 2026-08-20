package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/livecodelife/sedum/internal/catalog"
	"github.com/livecodelife/sedum/internal/genpkg"
)

func newActionsCommand() *cobra.Command {
	var cfg ActionsConfig

	cmd := &cobra.Command{
		Use:   "actions",
		Short: "Print a package's action catalog exactly as the model would receive it",
		Long: `Prints the exposed actions of a generator package with their kwarg schemas and
variant lists, in the form the model is given.

The authoring feedback loop for exposure and catalog clarity.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runActions(cmd.OutOrStdout(), cfg)
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

// runActions prints one package's catalog.
//
// The catalog is built by the same code Phase 4 builds its prompt's from, which
// is the entire reason this command is worth having: it is evidence of what the
// model receives rather than a second rendering that happens to agree today.
// Only the presentation below is this command's own.
func runActions(out io.Writer, cfg ActionsConfig) error {
	// Only the named package is loaded, so an unrelated package that fails to
	// load elsewhere in the directory does not stop an author inspecting this
	// one.
	set, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{Only: []string{cfg.Package}})
	if err != nil {
		return err
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			return fmt.Errorf("package %s did not load, so it has no catalog to print: %s", cfg.Package, f)
		}
	}

	pkg, ok := set.Lookup(cfg.Package)
	if !ok {
		return fmt.Errorf("no package named %q in %s; it declares %s",
			cfg.Package, cfg.Generators, packageNames(set))
	}

	c := catalog.Build([]*genpkg.Package{pkg}, catalog.Options{IncludeUnexposed: cfg.All})

	if cfg.JSON {
		payload, err := c.JSON()
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(payload))
		return nil
	}

	printCatalog(out, c, cfg)
	return nil
}

// printCatalog renders the catalog for a person.
//
// It states the same facts --json carries, in the order the model sees them. An
// author reading this is answering two questions - is the action here at all,
// and would the description let a model pick it - so exposure and the variant
// list are the things given room.
func printCatalog(out io.Writer, c catalog.Catalog, cfg ActionsConfig) {
	if len(c.Actions) == 0 {
		fmt.Fprintf(out, "package %s exposes no actions\n", cfg.Package)
		if !cfg.All {
			fmt.Fprintln(out, "\n(--all would include any it declares unexposed)")
		}
		return
	}

	for i, a := range c.Actions {
		if i > 0 {
			fmt.Fprintln(out)
		}

		hidden := ""
		if a.Exposed != nil && !*a.Exposed {
			hidden = "  [unexposed]"
		}
		fmt.Fprintf(out, "%s%s\n", a.Name, hidden)

		if a.Description != "" {
			fmt.Fprintf(out, "  %s\n", a.Description)
		}

		if len(a.Composes) > 0 {
			fmt.Fprintf(out, "  composes: %s\n", strings.Join(a.Composes, ", "))
		}

		for _, name := range sortedKwargs(a.Kwargs) {
			k := a.Kwargs[name]
			requirement := "optional"
			if k.Required {
				requirement = "required"
			}
			fmt.Fprintf(out, "  %-14s %-8s %s\n", name, k.Type, requirement)
			// On its own line rather than a fourth column, because a
			// description is a sentence and a column would truncate it or
			// wrap it into the next kwarg's row.
			if k.Description != "" {
				fmt.Fprintf(out, "    %s\n", k.Description)
			}
		}

		if a.Discriminator != "" {
			fallback := "no fallback, so any other value is an error"
			if a.HasDefault {
				fallback = "anything else falls to _default"
			}
			fmt.Fprintf(out, "  %s selects a template: %s (%s)\n",
				a.Discriminator, strings.Join(a.Variants, ", "), fallback)
		}

		printRequirements(out, a)
	}

	exposed := 0
	for _, a := range c.Actions {
		if a.Exposed == nil {
			exposed++
		}
	}
	fmt.Fprintf(out, "\n%d action(s), %d exposed to the model\n", len(c.Actions), exposed)
}

// printRequirements shows what an action's templates need beyond what the kwarg
// schema declares.
//
// This is in the human rendering because the command's whole claim is that it
// is evidence of what the model receives. A field carried in --json and omitted
// here would make the two disagree, which is the one thing this command cannot
// afford (prov-2026-369544c1).
func printRequirements(out io.Writer, a catalog.Action) {
	if len(a.Requires) > 0 {
		fmt.Fprintf(out, "  its template renders: %s\n", strings.Join(a.Requires, ", "))
	}
	if len(a.VariantRequires) == 0 {
		return
	}

	variants := make([]string, 0, len(a.VariantRequires))
	for variant := range a.VariantRequires {
		variants = append(variants, variant)
	}
	sort.Strings(variants)

	fmt.Fprintln(out, "  each template renders, and so needs bound:")
	for _, variant := range variants {
		needs := "nothing beyond the schema"
		if len(a.VariantRequires[variant]) > 0 {
			needs = strings.Join(a.VariantRequires[variant], ", ")
		}
		fmt.Fprintf(out, "    %-12s %s\n", variant, needs)
	}
}

func sortedKwargs(kwargs map[string]catalog.Kwarg) []string {
	out := make([]string, 0, len(kwargs))
	for name := range kwargs {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func packageNames(set *genpkg.Set) string {
	if len(set.Packages) == 0 {
		return "none"
	}
	names := make([]string, 0, len(set.Packages))
	for _, p := range set.Packages {
		names = append(names, p.Name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
