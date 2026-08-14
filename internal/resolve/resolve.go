// Package resolve runs Phases 2 and 3: it decides what each authorized path
// resolves to, and it creates the files.
//
// Phase 2 is pure. Every authorized path resolves to a generator package
// through the extensions that package declares - there is no built-in extension
// map, no default package, and no rule that infers a target from a path's shape
// or a directory's name - and then to the file template that matches it, with
// the captures that template binds (prov-2026-5696ff65).
//
// Phase 3 renders and writes. It is create-if-absent: a path that already
// exists is never re-rendered over, because doing so would destroy the injected
// regions the file already carries.
package resolve

import (
	"errors"
	"fmt"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/calebcowen/sedum/internal/filetmpl"
	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/record"
)

// defaultTemplate is the stem of the fallback file template. The one that
// applies to a path is the one carrying that path's extension, and no other
// (prov-2026-326598ac).
const defaultTemplate = "_default"

// Resolution is one authorized path with everything Phase 2 decided about it.
type Resolution struct {
	// RecordID is the record that authorized the path. Phase 4 makes one
	// model call per record, so a path carries its record forward.
	RecordID string
	Path     string
	Package  *genpkg.Package

	// Template is the file template that matched, as a pattern relative to
	// the package's files/ directory. Empty when nothing matched and the
	// package ships no default for the path's extension.
	Template string
	// Captures are the values Template's captures bound. Empty, never nil.
	Captures map[string]string
	// Default reports that Template is the package's fallback rather than a
	// pattern that matched the path.
	Default bool
}

// Paths runs Phase 2 over every authorized path.
//
// prefer holds the --lang values: package names to prefer where an extension is
// contested. A preference that cannot be honored is reported as a warning and
// never refuses the run (prov-2026-fe1e68b8).
//
// Errors are joined, because a records directory with three unresolvable paths
// should report three, not the first.
func Paths(set *genpkg.Set, records *record.Set, prefer []string) ([]Resolution, []string, error) {
	warnings := preferenceWarnings(set, records, prefer)

	matchers := map[string]*filetmpl.Set{}
	var (
		out      []Resolution
		problems []error
	)

	for _, rec := range records.Records {
		for _, p := range rec.Paths {
			pkg, err := packageFor(set, p, prefer)
			if err != nil {
				problems = append(problems, err)
				continue
			}

			res := Resolution{RecordID: rec.ID, Path: p, Package: pkg, Captures: map[string]string{}}
			if err := matchTemplate(&res, matchers); err != nil {
				problems = append(problems, err)
				continue
			}
			out = append(out, res)
		}
	}

	if len(problems) > 0 {
		return nil, warnings, errors.Join(problems...)
	}
	return out, warnings, nil
}

// packageFor resolves one path to the package that will generate it.
//
// The only thing that resolves a path is an extension some package declared.
// Every other property of the path - which directory it sits in, what its name
// looks like - is opaque here, which is what keeps target knowledge out of the
// core at the one place it would be easiest to smuggle in.
func packageFor(set *genpkg.Set, target string, prefer []string) (*genpkg.Package, error) {
	ext := strings.ToLower(path.Ext(target))
	if ext == "" {
		return nil, fmt.Errorf(
			"path %s has no extension, so no generator package can claim it; packages claim extensions and Sedum does not infer a target from a path's shape",
			target)
	}

	candidates := set.ForExtension(ext)
	switch len(candidates) {
	case 0:
		return nil, fmt.Errorf(
			"path %s carries extension %s, which no generator package claims; the loaded packages claim %s",
			target, ext, listOr(set.Extensions(), "no extensions at all"))
	case 1:
		return candidates[0], nil
	}

	// A contested extension is legal at load time and becomes an error only
	// here, at the point of real ambiguity.
	preferred := make([]*genpkg.Package, 0, len(candidates))
	for _, c := range candidates {
		if slices.Contains(prefer, c.Name) {
			preferred = append(preferred, c)
		}
	}

	switch len(preferred) {
	case 1:
		return preferred[0], nil
	case 0:
		return nil, fmt.Errorf(
			"path %s carries extension %s, which packages %s all claim; pass --lang to name the one to prefer",
			target, ext, listAnd(names(candidates)))
	default:
		return nil, fmt.Errorf(
			"path %s carries extension %s, claimed by packages %s, and --lang names %s; name exactly one",
			target, ext, listAnd(names(candidates)), listAnd(names(preferred)))
	}
}

// matchTemplate binds the file template for a resolution, falling back to the
// package's default for the path's extension.
func matchTemplate(res *Resolution, matchers map[string]*filetmpl.Set) error {
	pkg := res.Package

	set, ok := matchers[pkg.Name]
	if !ok {
		var err error
		if set, err = filetmpl.NewSet(pkg.FileTemplates); err != nil {
			return fmt.Errorf("package %s: %w", pkg.Name, err)
		}
		matchers[pkg.Name] = set
	}

	result, matched, err := set.Match(res.Path)
	if err != nil {
		// Loading already rejects a package whose templates tie, so this is
		// unreachable for a package that loaded. It is reported rather than
		// ignored because the alternative is picking one silently.
		return fmt.Errorf("package %s: %w", pkg.Name, err)
	}
	if matched {
		res.Template = result.Pattern
		res.Captures = result.Captures
		return nil
	}

	// No match is not an error. The caller decides what it means, and what it
	// means is the package's default for this extension, or an empty file.
	fallback := defaultTemplate + strings.ToLower(path.Ext(res.Path))
	if _, has := pkg.FileTemplate(fallback); has {
		res.Template = fallback
		res.Default = true
	}
	return nil
}

// preferenceWarnings reports --lang values the run cannot use. Neither is an
// error: --lang is a repeatable preference flag that scripts pass
// unconditionally, and a preference that cannot be honored is not a reason to
// refuse work (prov-2026-fe1e68b8).
func preferenceWarnings(set *genpkg.Set, records *record.Set, prefer []string) []string {
	carried := map[string]bool{}
	for _, p := range records.Paths() {
		carried[strings.ToLower(path.Ext(p))] = true
	}

	var out []string
	for _, name := range prefer {
		pkg, ok := set.Lookup(name)
		if !ok {
			out = append(out, fmt.Sprintf(
				"--lang names %q, which is not a package in the generators directory", name))
			continue
		}

		used := false
		for _, ext := range pkg.Extensions {
			if carried[strings.ToLower(ext)] {
				used = true
				break
			}
		}
		if !used {
			out = append(out, fmt.Sprintf(
				"--lang names %q, which claims no extension any authorized path carries; the preference has no effect on this run", name))
		}
	}
	return out
}

func names(packages []*genpkg.Package) []string {
	out := make([]string, 0, len(packages))
	for _, p := range packages {
		out = append(out, p.Name)
	}
	sort.Strings(out)
	return out
}

func listAnd(items []string) string { return join(items, "and") }

func listOr(items []string, empty string) string {
	if len(items) == 0 {
		return empty
	}
	return join(items, "or")
}

func join(items []string, conjunction string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return items[0] + " " + conjunction + " " + items[1]
	default:
		return strings.Join(items[:len(items)-1], ", ") + ", " + conjunction + " " + items[len(items)-1]
	}
}
