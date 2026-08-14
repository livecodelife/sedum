// Package expand runs Phase 6: it turns a validated list of action invocations
// into injections ready to write.
//
// It renders every injects_into pattern, requires each to resolve to exactly
// one file the run created, selects the variant template for a discriminated
// action, and applies transforms. It is fully deterministic - the model does not
// participate, and in this milestone no model has run at all: the invocation
// list is a hand-written fixture shaped exactly like a recording.
//
// A composite is expanded here into its children, in declaration order, each
// receiving the kwargs it declares out of the union the caller bound once. That
// is where a composite stops existing: everything after this point sees a list
// of simple invocations and contains no composite-aware branch.
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
//
// One invocation may produce more than one resolved injection: a composite
// produces one per child. They are returned in the order the composite declares
// them, in the position the invocation held.
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
		out = append(out, resolved...)
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	return out, nil
}

func expandOne(recordID string, packages []*genpkg.Package, files []resolve.File, inv recording.Invocation) ([]inject.Invocation, error) {
	pkg, action, err := lookupAction(packages, inv.Action)
	if err != nil {
		return nil, err
	}

	if action.Kind() == genpkg.Composite {
		return expandComposite(recordID, pkg, files, action, inv.Kwargs)
	}

	resolved, err := resolveOne(recordID, pkg, files, action, inv.Kwargs)
	if err != nil {
		return nil, err
	}
	return []inject.Invocation{resolved}, nil
}

// expandComposite turns one composite invocation into one resolved injection
// per child.
//
// Children run in declaration order. There is no reordering, no dependency
// resolution, and no data flow between them - a child's every value comes from
// the kwargs the caller bound, exactly as a directly invoked action's does.
//
// The structural rules a composite has to satisfy are enforced at load
// (prov-2026-8e6dac6c): one level of nesting, no reaching into another package,
// type-consistent kwarg unions. A package that fails any of them is rejected
// whole, so this consumes an already-valid composite and re-checks none of it.
//
// Every child is attempted rather than stopping at the first that fails, for the
// same reason the invocation list is: a composite with two mistakes in it should
// report two.
func expandComposite(recordID string, pkg *genpkg.Package, files []resolve.File, composite *genpkg.Action, kwargs map[string]any) ([]inject.Invocation, error) {
	var (
		out      []inject.Invocation
		problems []error
	)
	for _, childName := range composite.Composes {
		child, ok := pkg.Actions[childName]
		if !ok {
			// Load rejects a composite naming an action its package does
			// not define, so reaching here means a package was accepted
			// that should not have been. It is reported rather than
			// panicked on, and it is not a validation rule stated twice.
			problems = append(problems, fmt.Errorf(
				"composite %s composes %q, which package %s does not define; the package should not have loaded",
				composite.Name, childName, pkg.Name))
			continue
		}

		// The composite is named alongside whatever went wrong with the
		// child: the child is where the defect is, and the composite is
		// the invocation the author has to fix - a child may not even be
		// exposed (prov-2026-a0e37dae).
		resolved, err := resolveOne(recordID, pkg, files, child, project(child, kwargs))
		if err != nil {
			problems = append(problems, fmt.Errorf("composite %s: %w", composite.Name, err))
			continue
		}
		out = append(out, resolved)
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	return out, nil
}

// project maps a composite's union kwargs onto one child's declared schema.
//
// A composite's schema is the union of its children's, so the caller binds a
// shared kwarg once and both children receive it. Each child receives exactly
// what it declares: passing the whole union through would leave a region's
// recorded kwargs describing the composite rather than the region, and a
// declaration would claim to have been parameterized by an argument only the
// definition took.
//
// A kwarg the child declares and nothing bound is left absent rather than
// defaulted. Whether a required argument is missing is Phase 5's judgment, and
// a template that needs the value fails loudly at render either way.
func project(child *genpkg.Action, kwargs map[string]any) map[string]any {
	out := make(map[string]any, len(child.Kwargs))
	for name := range child.Kwargs {
		if value, bound := kwargs[name]; bound {
			out[name] = value
		}
	}
	return out
}

// resolveOne resolves one simple or discriminated action to the point where
// only writing is left. A composite never reaches it; by the time anything here
// runs, expansion has already happened.
func resolveOne(recordID string, pkg *genpkg.Package, files []resolve.File, action *genpkg.Action, kwargs map[string]any) (inject.Invocation, error) {
	path, err := renderPath(pkg, action, kwargs)
	if err != nil {
		return inject.Invocation{}, err
	}
	if err := authorized(action, files, path); err != nil {
		return inject.Invocation{}, err
	}

	variant, template, err := selectTemplate(action, kwargs)
	if err != nil {
		return inject.Invocation{}, err
	}

	content, err := renderTemplate(pkg, action, template, kwargs)
	if err != nil {
		return inject.Invocation{}, err
	}

	// The marker names this action, whether it was invoked directly or
	// reached through a composite (prov-2026-a0e37dae).
	return inject.Invocation{
		Package:  pkg,
		Action:   action,
		Variant:  variant,
		Kwargs:   kwargs,
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
		if f.Path != path {
			continue
		}
		// Authorized, but declared not-ours by a package. The record did its
		// job and the file is somebody else's, so the diagnostic says that
		// rather than reporting a path the record forgot.
		if f.Unmanaged {
			return fmt.Errorf(
				"action %s injects into %q, which package %s declares unmanaged (matching %q); that file is not written by Sedum",
				action.Name, path, f.UnmanagedBy, f.UnmanagedAs)
		}
		matches++
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
		// An unmanaged path resolved to no package, so it contributes
		// nothing to the record's catalog.
		if f.Unmanaged {
			continue
		}
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
		if f.Unmanaged {
			continue
		}
		paths = append(paths, strconv.Quote(f.Path))
	}
	if len(paths) == 0 {
		return "none"
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
