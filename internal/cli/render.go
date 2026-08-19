package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/render"
)

// renderResult is what one invocation's target rendering answers.
//
// It names the package and action it was asked about as well as the targets,
// because a caller asking speculatively across candidates is correlating
// answers and a bare list of paths would not say which question produced which.
type renderResult struct {
	Package string         `json:"package"`
	Action  string         `json:"action"`
	Targets []renderTarget `json:"targets"`
}

// renderTarget is one file an invocation lands in.
//
// The pattern is carried beside the path deliberately. A caller comparing a
// candidate against a region in a file needs to know which rule produced the
// answer - a path that matches for the wrong reason is the failure mode this
// command exists to prevent.
type renderTarget struct {
	Action      string `json:"action"`
	InjectsInto string `json:"injects_into"`
	Path        string `json:"path"`
}

func newRenderCommand() *cobra.Command {
	var cfg RenderConfig

	cmd := &cobra.Command{
		Use:   "render",
		Short: "Print the file an action's invocation would land in",
		Long: `Renders an action's injects_into pattern against supplied kwargs and prints the
path it resolves to.

This answers "where does this invocation land" for a tool built on Sedum. The
pattern is rendered by the package's own transform engine, so the answer agrees
with what a run would do rather than approximating it - a caller that
reimplemented the transforms would diverge on the first irregular plural and
report a path no run would ever produce.

Rendering is forward only. snake and plural are not invertible, so there is no
way to ask which kwargs produce a given path.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runRender(cmd.OutOrStdout(), cfg)
		},
	}

	f := cmd.Flags()
	f.StringVar(&cfg.Generators, "generators", "", "Generators directory. Required.")
	f.StringVar(&cfg.Package, "package", "", "Package declaring the action. Required.")
	f.StringVar(&cfg.Action, "action", "", "Action to render the target of. Required.")
	f.StringVar(&cfg.Kwargs, "kwargs", "{}", "Bound arguments, as a JSON object.")
	f.BoolVar(&cfg.JSON, "json", false, "Emit the result as JSON rather than formatted output.")

	mustMarkRequired(cmd, "generators", "package", "action")

	return cmd
}

// runRender resolves one invocation to the files it touches.
func runRender(out io.Writer, cfg RenderConfig) error {
	kwargs, err := parseKwargs(cfg.Kwargs)
	if err != nil {
		return err
	}

	// Only the named package is loaded, for the same reason sedum actions
	// loads one: an unrelated package that fails elsewhere in the directory
	// must not stop a caller asking about this one.
	set, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{Only: []string{cfg.Package}})
	if err != nil {
		return err
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			return fmt.Errorf("package %s did not load, so nothing can be rendered against it: %s", cfg.Package, f)
		}
	}

	pkg, ok := set.Lookup(cfg.Package)
	if !ok {
		return fmt.Errorf("no package named %q in %s; it declares %s",
			cfg.Package, cfg.Generators, packageNames(set))
	}

	action, ok := pkg.Actions[cfg.Action]
	if !ok {
		return fmt.Errorf("package %s declares no action named %q; it declares %s",
			cfg.Package, cfg.Action, actionNames(pkg))
	}

	targets, err := renderTargets(pkg, action, kwargs)
	if err != nil {
		return err
	}

	result := renderResult{Package: cfg.Package, Action: cfg.Action, Targets: targets}

	if cfg.JSON {
		payload, err := json.Marshal(result)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(payload))
		return nil
	}

	printRender(out, result)
	return nil
}

// renderTargets renders every pattern this invocation will land on.
//
// A composite has no pattern of its own and takes its children's in execution
// order, which is what makes one selection visibly touch two files. A child a
// package does not define cannot reach here - load rejects the package - so the
// lookup filters rather than restates a check.
func renderTargets(pkg *genpkg.Package, action *genpkg.Action, kwargs map[string]any) ([]renderTarget, error) {
	children := []*genpkg.Action{action}
	if action.Kind() == genpkg.Composite {
		children = children[:0]
		for _, name := range action.Composes {
			child, ok := pkg.Actions[name]
			if !ok {
				continue
			}
			children = append(children, child)
		}
	}

	targets := make([]renderTarget, 0, len(children))
	for _, child := range children {
		if child.InjectsInto == "" {
			continue
		}
		tmpl, err := render.Compile(pkg.Engine, child.InjectsInto)
		if err != nil {
			return nil, fmt.Errorf("action %s: injects_into %q: %w", child.Name, child.InjectsInto, err)
		}
		path, err := tmpl.Render(kwargs)
		if err != nil {
			return nil, fmt.Errorf("action %s: injects_into %q: %w", child.Name, child.InjectsInto, err)
		}
		targets = append(targets, renderTarget{
			Action:      child.Name,
			InjectsInto: child.InjectsInto,
			Path:        path,
		})
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("action %s declares no injects_into, so there is no file for it to land in", action.Name)
	}
	return targets, nil
}

// parseKwargs decodes the --kwargs object.
//
// Numbers are kept as json.Number rather than decoded to float64. A path is not
// where a caller should discover that 1000000 became 1e+06, and a value the
// caller wrote is a value the caller should see rendered.
func parseKwargs(src string) (map[string]any, error) {
	dec := json.NewDecoder(strings.NewReader(src))
	dec.UseNumber()

	var kwargs map[string]any
	if err := dec.Decode(&kwargs); err != nil {
		return nil, fmt.Errorf("--kwargs is not a JSON object: %w", err)
	}
	if err := dec.Decode(new(json.RawMessage)); err != io.EOF {
		return nil, fmt.Errorf("--kwargs carries more than one JSON value; an invocation binds one set of arguments")
	}
	if kwargs == nil {
		return map[string]any{}, nil
	}
	return kwargs, nil
}

// printRender states the same facts --json carries, for a person checking a
// pattern by hand.
func printRender(out io.Writer, result renderResult) {
	var buf bytes.Buffer
	for _, t := range result.Targets {
		if len(result.Targets) > 1 {
			fmt.Fprintf(&buf, "%s\n  ", t.Action)
		}
		fmt.Fprintf(&buf, "%s\n", t.Path)
		fmt.Fprintf(&buf, "  from %s\n", t.InjectsInto)
	}
	io.Copy(out, &buf)
}

func actionNames(pkg *genpkg.Package) string {
	if len(pkg.Actions) == 0 {
		return "none"
	}
	names := make([]string, 0, len(pkg.Actions))
	for name := range pkg.Actions {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
