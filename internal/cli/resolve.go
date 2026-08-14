package cli

import (
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
			return runResolve(cmd.OutOrStdout(), cmd.ErrOrStderr(), cfg)
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

// runResolve reports what every authorized path resolved to.
//
// It halts at the resolution boundary, so it consults no output tree: the only
// directories it reads are the two it was pointed at, and its answer does not
// depend on where it was run from (prov-2026-43808a47). Rendering for
// --show-template is Phase 3's, but it is the half that touches nothing.
func runResolve(out, errOut io.Writer, cfg ResolveConfig) error {
	result, err := pipeline.Run(pipeline.Config{
		Generators:     cfg.Generators,
		Records:        cfg.Records,
		Lang:           cfg.Lang,
		Only:           cfg.Only,
		StopAfterPhase: pipeline.PhaseResolve,
		Log:            runlog.Discard(),
	})
	if err != nil {
		return err
	}

	rendered := map[string]string{}
	if cfg.ShowTemplate {
		for _, r := range result.Resolutions {
			text, err := resolve.Render(r)
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

	fmt.Fprintf(out, "\n%d record(s), %d path(s) resolved across %d package(s)\n",
		len(result.Records.Records), len(result.Resolutions), len(result.Packages.Packages))
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

func printWarnings(errOut io.Writer, warnings []string) {
	for _, w := range warnings {
		fmt.Fprintf(errOut, "sedum: warning: %s\n", w)
	}
}
