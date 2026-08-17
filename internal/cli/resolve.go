package cli

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/calebcowen/sedum/internal/pipeline"
	"github.com/calebcowen/sedum/internal/resolve"
	"github.com/calebcowen/sedum/internal/runlog"
)

func newResolveCommand() *cobra.Command {
	var cfg ResolveConfig

	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Report what each authorized path resolves to, without invoking the model",
		Long: `Runs Phases 0 through 3 and reports, for every authorized path, the generator
package it resolved to, the file template that matched, and the captures bound.

Writes nothing. The primary debugging tool for package resolution and template
specificity.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResolve(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), cfg)
		},
	}

	f := cmd.Flags()
	f.StringVar(&cfg.Generators, "generators", "", "Generators directory. Required.")
	f.StringVar(&cfg.Records, "records", "", "Provenance records directory. Required.")
	f.StringArrayVar(&cfg.Vars, "var", nil, "Bind a variable a generator package declares, as name=value. Repeatable.")
	f.StringSliceVar(&cfg.Lang, "lang", nil, "Prefer the named package where an extension is contested. Repeatable.")
	f.StringSliceVar(&cfg.Only, "only", nil, "Resolve only the named provenance record. Repeatable.")
	f.BoolVar(&cfg.ShowTemplate, "show-template", false, "Include rendered template output for each path.")

	mustMarkRequired(cmd, "generators", "records")

	return cmd
}

// runResolve reports what every authorized path resolved to.
//
// It halts at the resolution boundary, so it consults no output tree: the only
// directories it reads are the two it was pointed at, and its answer does not
// depend on where it was run from (prov-2026-43808a47). Rendering for
// --show-template is Phase 3's, but it is the half that touches nothing.
func runResolve(ctx context.Context, out, errOut io.Writer, cfg ResolveConfig) error {
	variables, err := parseVars(cfg.Vars)
	if err != nil {
		return err
	}

	result, err := pipeline.Run(ctx, pipeline.Config{
		Generators:     cfg.Generators,
		Records:        cfg.Records,
		Lang:           cfg.Lang,
		Only:           cfg.Only,
		StopAfterPhase: pipeline.PhaseResolve,
		Variables:      variables,
		Log:            runlog.Discard(),
	})
	if err != nil {
		return err
	}

	rendered := map[string]string{}
	if cfg.ShowTemplate {
		for _, r := range result.Resolutions {
			text, err := resolve.Render(r, result.Variables)
			if err != nil {
				return err
			}
			rendered[r.Path] = text
		}
	}

	printWarnings(errOut, result.Warnings)
	printResolutions(out, result, rendered, cfg.ShowTemplate)
	return nil
}

// printResolutions reports every authorized path grouped under the record that
// authorized it, since a record is the unit a reader is reasoning about.
func printResolutions(out io.Writer, result *pipeline.Result, rendered map[string]string, showTemplate bool) {
	byRecord := map[string][]resolve.Resolution{}
	for _, r := range result.Resolutions {
		byRecord[r.RecordID] = append(byRecord[r.RecordID], r)
	}

	ids := make([]string, 0, len(byRecord))
	for id := range byRecord {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		fmt.Fprintf(out, "%s\n", id)
		for _, r := range byRecord[id] {
			fmt.Fprintf(out, "  %s\n", r.Path)
			if r.Unmanaged {
				// Named, not passed over in silence: the record
				// authorized this path and something other than Sedum
				// is what changes it.
				fmt.Fprintf(out, "    unmanaged %s declares %q; left for a person or another tool\n",
					r.UnmanagedBy, r.UnmanagedAs)
				continue
			}
			fmt.Fprintf(out, "    package   %s\n", r.Package.Name)
			fmt.Fprintf(out, "    template  %s\n", describeTemplate(r))
			if len(r.Captures) > 0 {
				fmt.Fprintf(out, "    captures  %s\n", describeCaptures(r.Captures))
			}
			if showTemplate {
				fmt.Fprintf(out, "    rendered\n%s", indentBlock(rendered[r.Path]))
			}
		}
	}

	var unmanaged []resolve.Resolution
	for _, r := range result.Resolutions {
		if r.Unmanaged {
			unmanaged = append(unmanaged, r)
		}
	}

	fmt.Fprintf(out, "\n%d record(s), %d path(s) resolved across %d package(s)\n",
		len(result.Records.Records), len(result.Resolutions)-len(unmanaged), len(result.Packages.Packages))

	// The handoff, gathered rather than scattered through the listing: these
	// are authorized paths this run did not touch and something else must.
	if len(unmanaged) > 0 {
		fmt.Fprintf(out, "\n%d path(s) left unmanaged, for a person or another tool:\n", len(unmanaged))
		for _, r := range unmanaged {
			fmt.Fprintf(out, "  %s (%s declares %q)\n", r.Path, r.UnmanagedBy, r.UnmanagedAs)
		}
	}
}

func describeTemplate(r resolve.Resolution) string {
	switch {
	case r.Template == "":
		return "(none matched; the file is created empty)"
	case r.Default:
		return r.Template + " (package default)"
	default:
		return r.Template
	}
}

func describeCaptures(captures map[string]string) string {
	names := make([]string, 0, len(captures))
	for name := range captures {
		names = append(names, name)
	}
	sort.Strings(names)

	pairs := make([]string, 0, len(names))
	for _, name := range names {
		pairs = append(pairs, name+"="+captures[name])
	}
	return strings.Join(pairs, " ")
}

// indentBlock offsets rendered output so that a template's own blank lines and
// indentation stay readable inside the report.
func indentBlock(text string) string {
	if text == "" {
		return "      (empty)\n"
	}

	var b strings.Builder
	for _, line := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		b.WriteString("      " + line + "\n")
	}
	return b.String()
}

// printWarnings reports what the phases collected.
//
// A warning that already says it is one is not told so twice: a load finding
// carries its own "warning:" prefix, and a phase's own warning does not.
func printWarnings(errOut io.Writer, warnings []string) {
	for _, w := range warnings {
		if strings.HasPrefix(w, "warning:") {
			fmt.Fprintf(errOut, "sedum: %s\n", w)
			continue
		}
		fmt.Fprintf(errOut, "sedum: warning: %s\n", w)
	}
}
