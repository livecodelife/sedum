// Package genpkg loads generator packages and runs every load-time check.
//
// A generators directory's top-level subdirectories are generator packages, one
// per target stack. Each declares its name, the extensions it claims, its
// comment prefix, and its transform pipelines in sedum.yaml; its actions in
// actions/actions.yaml; its file templates as a literal mirror of the target
// project's structure under files/; and its action templates under actions/.
//
// Packages are wholly valid or rejected. There is no partial load and no "load
// what parses" fallback: a package with one broken action does not contribute
// its working actions to a run, because a run that silently omits an action
// fails later as a missing feature rather than as a broken package.
//
// Loading needs no records, no model, and no network.
package genpkg

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/calebcowen/sedum/internal/pathpat"
	"github.com/calebcowen/sedum/internal/transform"
)

const (
	manifestFile    = "sedum.yaml"
	actionsDirName  = "actions"
	actionsFileName = "actions.yaml"
	filesDirName    = "files"
)

// Options controls which packages a load considers.
type Options struct {
	// Only names the packages to load. Empty loads every package in the
	// directory.
	Only []string
}

func (o Options) wanted(name string) bool {
	return len(o.Only) == 0 || slices.Contains(o.Only, name)
}

// Set is the loaded generators directory: the valid packages and the
// extension-to-package map built from them.
type Set struct {
	// Packages holds only packages that loaded clean, sorted by name.
	Packages []*Package

	byExtension map[string][]*Package
}

// ForExtension returns the packages claiming ext, in package-name order.
//
// More than one claimant is not an error here. A contested extension becomes an
// error only when a path carrying it appears and no --lang flag disambiguates,
// so a multi-package directory stays legal and fails at the point of real
// ambiguity.
func (s *Set) ForExtension(ext string) []*Package {
	return s.byExtension[strings.ToLower(ext)]
}

// Extensions returns every extension any loaded package claims, sorted.
func (s *Set) Extensions() []string {
	out := make([]string, 0, len(s.byExtension))
	for ext := range s.byExtension {
		out = append(out, ext)
	}
	sort.Strings(out)
	return out
}

// Lookup returns the package with the given name.
func (s *Set) Lookup(name string) (*Package, bool) {
	for _, p := range s.Packages {
		if p.Name == name {
			return p, true
		}
	}
	return nil, false
}

// Load reads every generator package under dir and validates it.
//
// The returned error is for I/O problems that stop the load from happening at
// all - a missing or unreadable generators directory. Everything wrong with a
// package's contents comes back as findings, because reporting them all in one
// pass is what lets an author fix a package once rather than iteratively.
func Load(dir string, opts Options) (*Set, Findings, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("read generators directory %s: %w", dir, err)
	}

	set := &Set{byExtension: map[string][]*Package{}}
	var findings Findings

	for _, entry := range entries {
		if !entry.IsDir() || !opts.wanted(entry.Name()) {
			continue
		}

		pkg, pkgFindings, err := loadPackage(filepath.Join(dir, entry.Name()), entry.Name())
		if err != nil {
			return nil, nil, err
		}
		findings = append(findings, pkgFindings...)
		if pkg == nil {
			continue
		}
		set.Packages = append(set.Packages, pkg)
	}

	sort.Slice(set.Packages, func(i, j int) bool { return set.Packages[i].Name < set.Packages[j].Name })
	for _, p := range set.Packages {
		for _, ext := range p.Extensions {
			ext = strings.ToLower(ext)
			set.byExtension[ext] = append(set.byExtension[ext], p)
		}
	}

	return set, findings.sorted(), nil
}

// loadPackage reads one package directory. It returns a nil package when the
// package was rejected. dirName is the directory's own name, which is what the
// declared name must agree with and what the author types after --package.
func loadPackage(dir, dirName string) (*Package, Findings, error) {
	r := &reporter{pkg: dirName}

	man, ok, err := readManifest(dir, r)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		// Without a manifest there is no comment prefix and no transform
		// table, so no later check can produce a trustworthy answer.
		return nil, r.findings, nil
	}

	pkg := &Package{
		Name:          dirName,
		Dir:           dir,
		Extensions:    man.Extensions,
		CommentPrefix: man.CommentPrefix,
		Transforms:    man.Transforms,
		OpExceptions:  man.OpExceptions,
		Unmanaged:     man.Unmanaged,
		Actions:       map[string]*Action{},
	}

	// The engine is built before anything reads a template, because the
	// vocabulary a reference is checked against is exactly what it declares.
	// When it cannot be built, the reference check is skipped rather than run
	// against an empty vocabulary, which would report every template in the
	// package as broken and bury the one pipeline that actually is.
	engine, err := transform.New(transform.Config{
		Pipelines:  man.Transforms,
		Exceptions: man.OpExceptions,
	})
	if err != nil {
		for _, problem := range unwrap(err) {
			r.errorf(manifestFile, RuleTransformInvalid, "%v", problem)
		}
	}
	pkg.Engine = engine

	if man.Name != dirName {
		r.errorf(manifestFile, RuleNameMismatch,
			"package declares name %q but its directory is named %q; the directory name is what --package and --lang refer to",
			man.Name, dirName)
	}
	for _, ext := range man.Extensions {
		if !strings.HasPrefix(ext, ".") {
			r.errorf(manifestFile, RuleExtensionInvalid,
				"extension %q does not start with a dot; paths are matched by their dotted extension", ext)
		}
	}

	decls, ok, err := readActions(dir, r)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, r.findings, nil
	}

	fileTemplates, err := loadFileTemplates(dir, r)
	if err != nil {
		return nil, nil, err
	}
	pkg.FileTemplates = fileTemplates.patterns
	pkg.fileContents = fileTemplates.contents

	// Actions are built before their templates are resolved, because the
	// declared shape is what says where to look.
	validateActions(pkg, decls, r)

	actionTemplates, err := resolveActionTemplates(dir, pkg, r)
	if err != nil {
		return nil, nil, err
	}
	pkg.actionContents = actionTemplates

	checkTemplates(pkg, fileTemplates.contents, actionTemplates, r)
	deriveRequirements(pkg, actionTemplates, r)
	checkMarkerAnchors(pkg, fileTemplates.bodies(), r)
	checkMarkersFilled(pkg, fileTemplates.bodies(), r)
	checkUnmanaged(pkg, r)
	checkDeadConfig(pkg, r)

	if r.hasErrors() {
		return nil, r.findings, nil
	}
	return pkg, r.findings, nil
}

func readManifest(dir string, r *reporter) (manifest, bool, error) {
	var man manifest

	data, err := os.ReadFile(filepath.Join(dir, manifestFile))
	switch {
	case errors.Is(err, fs.ErrNotExist):
		r.errorf(manifestFile, RuleManifestMissing,
			"package has no %s; every package declares its name, extensions, and comment prefix there", manifestFile)
		return man, false, nil
	case err != nil:
		return man, false, fmt.Errorf("read %s: %w", filepath.Join(dir, manifestFile), err)
	}

	if err := strictUnmarshal(data, &man); err != nil {
		r.errorf(manifestFile, RuleManifestMalformed, "%s could not be read: %v", manifestFile, err)
		return man, false, nil
	}
	return man, true, nil
}

func readActions(dir string, r *reporter) (map[string]*actionDecl, bool, error) {
	rel := filepath.ToSlash(filepath.Join(actionsDirName, actionsFileName))

	data, err := os.ReadFile(filepath.Join(dir, actionsDirName, actionsFileName))
	switch {
	case errors.Is(err, fs.ErrNotExist):
		r.errorf(rel, RuleActionsMissing,
			"package has no %s; a package with no actions declares an empty actions map rather than omitting the file", rel)
		return nil, false, nil
	case err != nil:
		return nil, false, fmt.Errorf("read %s: %w", filepath.Join(dir, actionsDirName, actionsFileName), err)
	}

	var file actionsFile
	if err := strictUnmarshal(data, &file); err != nil {
		r.errorf(rel, RuleActionsMalformed, "%s could not be read: %v", rel, err)
		return nil, false, nil
	}
	if file.Actions == nil {
		file.Actions = map[string]*actionDecl{}
	}
	for name, decl := range file.Actions {
		if decl == nil {
			file.Actions[name] = &actionDecl{}
		}
	}
	return file.Actions, true, nil
}

// strictUnmarshal decodes with KnownFields set, so an unrecognized key is an
// error naming the key rather than a silently ignored field.
func strictUnmarshal(data []byte, into any) error {
	dec := yaml.NewDecoder(strings.NewReader(string(data)))
	dec.KnownFields(true)
	if err := dec.Decode(into); err != nil {
		return err
	}
	return nil
}

// fileTemplateSet is the patterns under files/ and the contents behind them.
// The contents are kept because two checks read them - transform references and
// planted markers - and because Phase 3 renders them.
//
// Contents are keyed by pattern rather than held in a parallel slice. The
// patterns are sorted so findings do not depend on directory order, while a
// walk is depth-first over each directory's entries, and the two orders diverge
// as soon as a file and a directory share a prefix (prov-2026-11af2675).
type fileTemplateSet struct {
	patterns []string
	contents map[string]string
}

// bodies returns the template contents in pattern order, for the checks that
// read every template without caring which is which.
func (s fileTemplateSet) bodies() []string {
	out := make([]string, 0, len(s.patterns))
	for _, p := range s.patterns {
		out = append(out, s.contents[p])
	}
	return out
}

func loadFileTemplates(dir string, r *reporter) (fileTemplateSet, error) {
	out := fileTemplateSet{contents: map[string]string{}}

	root := filepath.Join(dir, filesDirName)
	if _, err := os.Stat(root); errors.Is(err, fs.ErrNotExist) {
		// A package with no file templates is legal: not every target
		// needs boilerplate in newly created files.
		return out, nil
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		pattern := filepath.ToSlash(rel)
		out.patterns = append(out.patterns, pattern)
		out.contents[pattern] = string(data)
		return nil
	})
	if err != nil {
		return out, fmt.Errorf("walk %s: %w", root, err)
	}

	// Sorted so that findings do not depend on directory iteration order.
	sort.Strings(out.patterns)
	checkFileTemplatePatterns(out.patterns, r)
	return out, nil
}

// resolveActionTemplates resolves each action's declared shape against the
// filesystem, records the resolved paths on the action, and returns the
// template contents for transform checking, keyed by their path relative to the
// package directory.
//
// The shape comes from the schema, never from what is on disk. A composite
// triggers no lookup at all - searching the filesystem for one is a bug, not a
// fallback.
func resolveActionTemplates(dir string, pkg *Package, r *reporter) (map[string]string, error) {
	actionsDir := filepath.Join(dir, actionsDirName)
	contents := map[string]string{}

	for _, name := range sortedKeys(pkg.Actions) {
		action := pkg.Actions[name]
		switch action.Kind() {
		case Composite:
			continue
		case Simple:
			if err := resolveSimpleTemplate(actionsDir, action, contents, r); err != nil {
				return nil, err
			}
		case Discriminated:
			if err := resolveVariantTemplates(actionsDir, action, contents, r); err != nil {
				return nil, err
			}
		}
	}
	return contents, nil
}

// resolveSimpleTemplate finds the single file named for the action. Its
// extension is not fixed: a package's templates carry whatever extension its
// target uses, and the package may claim several.
func resolveSimpleTemplate(actionsDir string, action *Action, contents map[string]string, r *reporter) error {
	name := action.Name
	rel := filepath.ToSlash(filepath.Join(actionsDirName, name))

	if info, err := os.Stat(filepath.Join(actionsDir, name)); err == nil && info.IsDir() {
		r.errorf(rel, RuleTemplateWrongForm,
			"action %s declares no discriminator, so its template is a single file named %s.<ext>, but %s is a directory",
			name, name, rel)
		return nil
	}

	matches, err := siblingsNamed(actionsDir, name)
	if err != nil {
		return err
	}
	switch len(matches) {
	case 0:
		r.errorf(rel, RuleTemplateMissing,
			"action %s has no template; a simple action's template is a single file named %s.<ext> under %s/",
			name, name, actionsDirName)
		return nil
	case 1:
	default:
		r.errorf(rel, RuleTemplateAmbiguous,
			"action %s resolves to %d templates (%s); a simple action must name exactly one",
			name, len(matches), strings.Join(matches, ", "))
		return nil
	}

	data, err := os.ReadFile(filepath.Join(actionsDir, matches[0]))
	if err != nil {
		return fmt.Errorf("read %s: %w", filepath.Join(actionsDir, matches[0]), err)
	}
	action.Template = filepath.ToSlash(filepath.Join(actionsDirName, matches[0]))
	contents[action.Template] = string(data)
	return nil
}

func resolveVariantTemplates(actionsDir string, action *Action, contents map[string]string, r *reporter) error {
	name := action.Name
	rel := filepath.ToSlash(filepath.Join(actionsDirName, name))
	dir := filepath.Join(actionsDir, name)

	info, err := os.Stat(dir)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		// The directory may be absent, or a file of that name may stand
		// where the directory belongs. Both are the schema disagreeing
		// with the filesystem, but they are different mistakes.
		siblings, sErr := siblingsNamed(actionsDir, name)
		if sErr != nil {
			return sErr
		}
		if len(siblings) > 0 {
			r.errorf(rel, RuleTemplateWrongForm,
				"action %s declares discriminator %q, so its templates live in a directory named %s, but %s is a file",
				name, action.Discriminator, name, filepath.ToSlash(filepath.Join(actionsDirName, siblings[0])))
			return nil
		}
		r.errorf(rel, RuleTemplateMissing,
			"action %s declares discriminator %q, so its templates live in %s/ with one file per variant, but that directory does not exist",
			name, action.Discriminator, rel)
		return nil
	case err != nil:
		return fmt.Errorf("stat %s: %w", dir, err)
	case !info.IsDir():
		r.errorf(rel, RuleTemplateWrongForm,
			"action %s declares discriminator %q, so its templates live in a directory named %s, but %s is a file",
			name, action.Discriminator, name, rel)
		return nil
	}

	for _, variant := range append(append([]string{}, action.Variants...), DefaultVariant) {
		matches, err := siblingsNamed(dir, variant)
		if err != nil {
			return err
		}
		if len(matches) == 0 {
			if variant == DefaultVariant {
				// _default is optional: an action may cover every
				// value it accepts.
				continue
			}
			r.errorf(rel, RuleVariantTemplateMissing,
				"action %s declares variant %q but %s/%s.<ext> does not exist; a misspelled variant filename must fail here rather than fall through to %s",
				name, variant, rel, variant, DefaultVariant)
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, matches[0]))
		if err != nil {
			return fmt.Errorf("read %s: %w", filepath.Join(dir, matches[0]), err)
		}
		if action.Templates == nil {
			action.Templates = map[string]string{}
		}
		path := filepath.ToSlash(filepath.Join(actionsDirName, name, matches[0]))
		action.Templates[variant] = path
		contents[path] = string(data)
	}
	return nil
}

// siblingsNamed returns the files in dir whose name without extension is stem,
// sorted.
func siblingsNamed(dir, stem string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}

	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.TrimSuffix(name, filepath.Ext(name)) == stem {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out, nil
}

// unwrap flattens a joined error into the problems it was built from, so that
// each one becomes its own finding rather than a paragraph inside one.
func unwrap(err error) []error {
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		return joined.Unwrap()
	}
	return []error{err}
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Unmanaged reports whether any loaded package declares this path one Sedum
// does not write, returning the package that declared it and the entry that
// matched.
//
// The union across packages rather than one package's list, because the check
// runs before a path has been resolved to a package - which is the point, since
// a path with no extension can never reach the package that disowns it.
//
// Two packages declaring the same path do not conflict; they agree. A path one
// package disowns and another would claim by extension is unmanaged, because a
// declaration that a file is hand-owned is a statement about the file, and the
// alternative is a run whose behavior depends on which package was consulted
// first.
func (s *Set) Unmanaged(target string) (pkg, entry string, ok bool) {
	for _, p := range s.Packages {
		if entry, ok := pathpat.MatchAny(p.Unmanaged, target); ok {
			return p.Name, entry, true
		}
	}
	return "", "", false
}
