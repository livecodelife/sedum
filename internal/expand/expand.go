// Package expand runs Phase 6: it turns a validated list of action invocations
// into injections ready to write.
//
// It renders every injects_into pattern, requires each to resolve to exactly
// one file the run created, selects the variant template for a discriminated
// action, and applies transforms. It is fully deterministic - the model does not
// participate, and in this milestone no model has run at all: the invocation
// list is a hand-written fixture shaped exactly like a recording.
//
// Composite expansion is M5. A composite reaching this package is reported as
// unsupported rather than silently treated as a simple action.
package expand

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/inject"
	"github.com/calebcowen/sedum/internal/recording"
	"github.com/calebcowen/sedum/internal/render"
	"github.com/calebcowen/sedum/internal/resolve"
)

// Expand resolves one record's invocations against the files that record
// created.
//
// The files are the run's Phase 3 output. They decide two things: which
// packages the record's catalog is drawn from, and which paths an injects_into
// pattern is allowed to name.
//
// Every problem is reported rather than the first, because a fixture invocation
// list with three mistakes should report three.
func Expand(recordID string, files []resolve.File, invocations []recording.Invocation) ([]inject.Invocation, error) {
	packages := packagesOf(files)

	var (
		out      []inject.Invocation
		problems []error
	)
	for _, inv := range invocations {
		resolved, err := expandOne(recordID, packages, files, inv)
		if err != nil {
			problems = append(problems, err)
			continue
		}
		out = append(out, resolved)
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	return out, nil
}

func expandOne(recordID string, packages []*genpkg.Package, files []resolve.File, inv recording.Invocation) (inject.Invocation, error) {
	pkg, action, err := lookupAction(packages, inv.Action)
	if err != nil {
		return inject.Invocation{}, err
	}

	if action.Kind() == genpkg.Composite {
		return inject.Invocation{}, fmt.Errorf(
			"action %s is a composite; composite expansion is not implemented", action.Name)
	}

	path, err := renderPath(pkg, action, inv.Kwargs)
	if err != nil {
		return inject.Invocation{}, err
	}
	if err := authorized(action, files, path); err != nil {
		return inject.Invocation{}, err
	}

	variant, template, err := selectTemplate(action, inv.Kwargs)
	if err != nil {
		return inject.Invocation{}, err
	}

	content, err := renderTemplate(pkg, action, template, inv.Kwargs)
	if err != nil {
		return inject.Invocation{}, err
	}

	return inject.Invocation{
		Package:  pkg,
		Action:   action,
		Variant:  variant,
		Kwargs:   inv.Kwargs,
		Path:     path,
		RecordID: recordID,
		Content:  content,
	}, nil
}

// lookupAction finds an action by name across the packages the record's paths
// resolved to.
//
// Resolution is per file rather than per run, so a record legitimately spans
// packages and its catalog is their union. An action named by two of them is
// ambiguous and is reported rather than resolved by an ordering rule nothing
// declared.
func lookupAction(packages []*genpkg.Package, name string) (*genpkg.Package, *genpkg.Action, error) {
	var (
		found   *genpkg.Action
		owner   *genpkg.Package
		clashes []string
	)
	for _, pkg := range packages {
		action, ok := pkg.Actions[name]
		if !ok {
			continue
		}
		if found != nil {
			clashes = append(clashes, pkg.Name)
			continue
		}
		found, owner = action, pkg
	}

	switch {
	case found == nil:
		return nil, nil, fmt.Errorf("no action named %q is declared by %s",
			name, packageList(packages))
	case len(clashes) > 0:
		return nil, nil, fmt.Errorf("action %q is declared by more than one package (%s); nothing says which one is meant",
			name, strings.Join(append([]string{owner.Name}, clashes...), ", "))
	}
	return owner, found, nil
}

// renderPath renders an action's injects_into pattern against its bound kwargs.
func renderPath(pkg *genpkg.Package, action *genpkg.Action, kwargs map[string]any) (string, error) {
	if action.InjectsInto == "" {
		return "", fmt.Errorf("action %s declares no injects_into, so there is no file to inject into", action.Name)
	}

	tmpl, err := render.Compile(pkg.Engine, action.InjectsInto)
	if err != nil {
		return "", fmt.Errorf("action %s: injects_into %q: %w", action.Name, action.InjectsInto, err)
	}
	path, err := tmpl.Render(kwargs)
	if err != nil {
		return "", fmt.Errorf("action %s: injects_into %q: %w", action.Name, action.InjectsInto, err)
	}
	return path, nil
}

// authorized requires a rendered injects_into to name exactly one file the run
// created.
//
// Zero is the interesting case. It means the action targets a path no
// provenance record authorized, which is the failure a record omitting a
// companion file produces - a C++ record naming the implementation but not the
// header lands exactly here, naming the file the author forgot.
func authorized(action *genpkg.Action, files []resolve.File, path string) error {
	var matches int
	for _, f := range files {
		if f.Path == path {
			matches++
		}
	}

	switch matches {
	case 1:
		return nil
	case 0:
		return fmt.Errorf(
			"action %s injects into %q, which is not one of the paths this record created (%s)",
			action.Name, path, pathList(files))
	default:
		return fmt.Errorf(
			"action %s injects into %q, which %d created files claim; injects_into must resolve to exactly one file",
			action.Name, path, matches)
	}
}

// selectTemplate chooses the template an invocation renders and the variant
// name that selected it.
//
// A discriminated action's variants are declared explicitly rather than
// inferred from the directory, so a discriminator value with no template falls
// to _default when the package ships one and is an error when it does not. The
// alternative - falling through silently - is the invisible cliff the catalog's
// variant list exists to make visible.
func selectTemplate(action *genpkg.Action, kwargs map[string]any) (variant, template string, err error) {
	if action.Kind() != genpkg.Discriminated {
		return "", action.Template, nil
	}

	raw, bound := kwargs[action.Discriminator]
	if !bound {
		return "", "", fmt.Errorf(
			"action %s selects its template with the %q kwarg, which nothing bound",
			action.Name, action.Discriminator)
	}
	value, ok := raw.(string)
	if !ok {
		return "", "", fmt.Errorf(
			"action %s selects its template with the %q kwarg, whose value %v is not a string",
			action.Name, action.Discriminator, raw)
	}

	if path, ok := action.Templates[value]; ok {
		return value, path, nil
	}
	if path, ok := action.Templates[genpkg.DefaultVariant]; ok {
		return value, path, nil
	}
	return "", "", fmt.Errorf(
		"action %s has no template for %s %q and the package ships no %s; the templates it declares are %s",
		action.Name, action.Discriminator, value, genpkg.DefaultVariant, variantList(action))
}

// renderTemplate renders an action's template against its bound kwargs.
func renderTemplate(pkg *genpkg.Package, action *genpkg.Action, template string, kwargs map[string]any) (string, error) {
	source, ok := pkg.ActionTemplate(template)
	if !ok {
		return "", fmt.Errorf("package %s: action template %s vanished between load and generation",
			pkg.Name, template)
	}

	tmpl, err := render.Compile(pkg.Engine, source)
	if err != nil {
		return "", fmt.Errorf("package %s: action template %s: %w", pkg.Name, template, err)
	}
	out, err := tmpl.Render(kwargs)
	if err != nil {
		return "", fmt.Errorf("action %s: template %s: %w", action.Name, template, err)
	}
	return out, nil
}

// packagesOf returns the distinct packages the record's files resolved to, in
// name order so that a diagnostic listing them reads the same way twice.
func packagesOf(files []resolve.File) []*genpkg.Package {
	seen := map[string]*genpkg.Package{}
	for _, f := range files {
		seen[f.Package.Name] = f.Package
	}

	out := make([]*genpkg.Package, 0, len(seen))
	for _, name := range sortedKeys(seen) {
		out = append(out, seen[name])
	}
	return out
}

func packageList(packages []*genpkg.Package) string {
	if len(packages) == 0 {
		return "any package this record's paths resolved to"
	}
	names := make([]string, 0, len(packages))
	for _, p := range packages {
		names = append(names, p.Name)
	}
	return strings.Join(names, ", ")
}

func pathList(files []resolve.File) string {
	if len(files) == 0 {
		return "none"
	}
	paths := make([]string, 0, len(files))
	for _, f := range files {
		paths = append(paths, strconv.Quote(f.Path))
	}
	sort.Strings(paths)
	return strings.Join(paths, ", ")
}

func variantList(action *genpkg.Action) string {
	names := make([]string, 0, len(action.Templates))
	for name := range action.Templates {
		names = append(names, strconv.Quote(name))
	}
	if len(names) == 0 {
		return "none"
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
