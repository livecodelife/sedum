package genpkg

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/calebcowen/sedum/internal/filetmpl"
	"github.com/calebcowen/sedum/internal/pathpat"
	"github.com/calebcowen/sedum/internal/render"
)

// The load-time checks. Each one reports every instance it finds rather than
// the first, so that an author sees everything wrong with a package in one
// pass.

const actionsRel = actionsDirName + "/" + actionsFileName

// validateActions turns declarations into resolved actions, checking each one's
// own coherence first and then the composite rules that need every action to
// have been seen.
func validateActions(pkg *Package, decls map[string]*actionDecl, r *reporter) {
	names := sortedKeys(decls)

	for _, name := range names {
		pkg.Actions[name] = buildAction(name, decls[name], pkg, r)
	}
	for _, name := range names {
		resolveComposite(pkg.Actions[name], pkg, r)
	}
}

// ReservedActionName is the one name an action may not take.
//
// Ownership markers and the anchor declarations file templates plant share the
// "sedum:" namespace and are told apart by the first segment of the label, so
// an action called "anchor" writes markers a reader is obliged to skip. Its
// regions become invisible to the matching that would replace them, and it
// injects a fresh duplicate on every run - silently, which is the one anchor
// mistake that does not announce itself (prov-2026-42f5eedd).
const ReservedActionName = "anchor"

func buildAction(name string, decl *actionDecl, pkg *Package, r *reporter) *Action {
	if name == ReservedActionName {
		r.errorf(actionsRel, RuleActionNameReserved,
			"action %q takes the one name that is reserved: %q is how a file template declares an anchor point, so an action of that name writes markers Sedum cannot tell from injection sites and duplicates its output on every run",
			name, markerDecl(pkg.CommentPrefix, "<name>"))
	}

	a := &Action{
		Name:          name,
		Kwargs:        decl.Kwargs,
		Discriminator: decl.Discriminator,
		Variants:      decl.Variants,
		InjectsInto:   decl.InjectsInto,
		Anchor:        decl.Anchor,
		AnchorStart:   decl.AnchorStart,
		AnchorEnd:     decl.AnchorEnd,
		AnchorPattern: decl.AnchorPattern,
		Composes:      decl.Composes,
		// exposed defaults to true: authoring an action is enough to make
		// it usable, and hiding is the deliberate act.
		Exposed: decl.Exposed == nil || *decl.Exposed,
	}
	if a.Kwargs == nil {
		a.Kwargs = map[string]Kwarg{}
	}

	switch {
	case len(decl.Composes) > 0:
		a.kind = Composite
	case decl.Discriminator != "":
		a.kind = Discriminated
	default:
		a.kind = Simple
	}

	for _, kw := range sortedKeys(a.Kwargs) {
		if !validKwargType(a.Kwargs[kw].Type) {
			r.errorf(actionsRel, RuleKwargTypeUnknown,
				"action %s: kwarg %s declares type %q; the closed set is %s",
				name, kw, a.Kwargs[kw].Type, strings.Join(KwargTypes, "|"))
		}
	}

	if a.kind == Composite {
		checkCompositeShape(a, decl, r)
		return a
	}
	checkSimpleShape(a, r)
	return a
}

// checkCompositeShape rejects a composite that also carries the fields of a
// simple action. An action's kind is read from the schema, so a declaration
// that reads as two kinds at once has no determinable kind at all.
func checkCompositeShape(a *Action, decl *actionDecl, r *reporter) {
	var also []string
	for field, set := range map[string]bool{
		"kwargs":         len(decl.Kwargs) > 0,
		"injects_into":   decl.InjectsInto != "",
		"anchor":         decl.Anchor != "",
		"discriminator":  decl.Discriminator != "",
		"variants":       len(decl.Variants) > 0,
		"anchor_pattern": decl.AnchorPattern != "",
	} {
		if set {
			also = append(also, field)
		}
	}
	if len(also) == 0 {
		return
	}
	sort.Strings(also)
	r.errorf(actionsRel, RuleCompositeMalformed,
		"action %s declares composes and also %s; a composite has no template of its own and takes its kwargs from its children",
		a.Name, strings.Join(also, ", "))
}

func checkSimpleShape(a *Action, r *reporter) {
	if a.InjectsInto == "" {
		r.errorf(actionsRel, RuleActionIncomplete,
			"action %s declares no injects_into; a simple action must name the file it writes into", a.Name)
	}
	if a.Anchor == "" {
		r.errorf(actionsRel, RuleActionIncomplete,
			"action %s declares no anchor; a simple action must say where in the file its output belongs", a.Name)
	}

	if len(a.Variants) > 0 && a.Discriminator == "" {
		r.errorf(actionsRel, RuleActionIncomplete,
			"action %s declares variants %s but no discriminator; variants are the values of the kwarg the discriminator names",
			a.Name, strings.Join(a.Variants, ", "))
	}
	if a.Discriminator != "" {
		if _, ok := a.Kwargs[a.Discriminator]; !ok {
			r.errorf(actionsRel, RuleDiscriminatorUnknown,
				"action %s declares discriminator %q, which is not one of its kwargs (%s)",
				a.Name, a.Discriminator, strings.Join(sortedKeys(a.Kwargs), ", "))
		}
	}

	checkAnchorFields(a, r)
}

// checkAnchorFields enforces that each anchor kind carries exactly the
// companion fields it needs. A field that does not apply to the declared anchor
// is a diagnostic, never an ignored line (prov-2026-d1d61186).
func checkAnchorFields(a *Action, r *reporter) {
	if a.Anchor == AnchorMarker {
		r.errorf(actionsRel, RuleAnchorInvalid,
			"action %s is anchored to %q, which is the kind and not a marker name; name the marker its file template plants",
			a.Name, AnchorMarker)
		return
	}

	needsRegion := a.Anchor == AnchorRegion
	needsPattern := a.Anchor == AnchorAfterMatch || a.Anchor == AnchorBeforeMatch

	if needsRegion && (a.AnchorStart == "" || a.AnchorEnd == "") {
		r.errorf(actionsRel, RuleAnchorInvalid,
			"action %s is anchored to a region but does not declare both anchor_start and anchor_end", a.Name)
	}
	if !needsRegion && (a.AnchorStart != "" || a.AnchorEnd != "") {
		r.errorf(actionsRel, RuleAnchorInvalid,
			"action %s declares anchor_start/anchor_end but its anchor is %q, not %q",
			a.Name, a.Anchor, AnchorRegion)
	}
	if needsPattern && a.AnchorPattern == "" {
		r.errorf(actionsRel, RuleAnchorInvalid,
			"action %s is anchored to %q but declares no anchor_pattern", a.Name, a.Anchor)
	}
	if !needsPattern && a.AnchorPattern != "" {
		r.errorf(actionsRel, RuleAnchorInvalid,
			"action %s declares anchor_pattern but its anchor is %q, not %s or %s",
			a.Name, a.Anchor, AnchorAfterMatch, AnchorBeforeMatch)
	}

	if a.AnchorPattern != "" {
		checkAnchorPattern(a, r)
	}
}

// checkAnchorPattern reads the declared expression here rather than leaving it
// to Phase 7.
//
// An expression that does not parse is a defect in the package, and a regex is
// checkable with nothing but the regex. Leaving it until injection means a
// package is declared wholly valid and then fails partway through a run that
// has already written files - the position load-time transform checking already
// rejected (prov-2026-4675cebe).
func checkAnchorPattern(a *Action, r *reporter) {
	if _, err := regexp.Compile(a.AnchorPattern); err != nil {
		r.errorf(actionsRel, RuleAnchorInvalid,
			"action %s declares anchor_pattern %q, which is not a valid expression: %v",
			a.Name, a.AnchorPattern, err)
		return
	}

	// ^ and $ match the bounds of the whole file unless (?m) is set, so a
	// pattern meant to find a line finds nothing and reports it at Phase 7
	// as a fault in the file rather than in the pattern.
	//
	// This warns rather than erroring, and the expression is never rewritten:
	// whole-text anchoring is legal and occasionally meant, and the pattern
	// written in actions.yaml has to be the pattern that runs.
	if !strings.Contains(a.AnchorPattern, "(?m") && lineAnchored(a.AnchorPattern) {
		r.warnf(actionsRel, RuleAnchorPatternLineMode,
			"action %s declares anchor_pattern %q, whose ^ or $ matches the bounds of the whole file rather than of a line; write (?m) at the start of the pattern if a line was meant",
			a.Name, a.AnchorPattern)
	}
}

// lineAnchored reports whether a pattern uses ^ or $ as an anchor, ignoring the
// ones that are escaped or that sit inside a character class, where they mean
// something else entirely.
func lineAnchored(pattern string) bool {
	var inClass bool

	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '\\':
			i++
		case '[':
			inClass = true
		case ']':
			inClass = false
		case '^', '$':
			if !inClass {
				return true
			}
		}
	}
	return false
}

// resolveComposite enforces the composite rules and builds the kwarg union.
// They are enforced here rather than at generation because a composite that
// cannot be expanded is a defect in the package, not in the run.
func resolveComposite(a *Action, pkg *Package, r *reporter) {
	if a.Kind() != Composite {
		return
	}

	union := map[string]Kwarg{}
	for _, childName := range a.Composes {
		child, ok := pkg.Actions[childName]
		if !ok {
			// actions.yaml holds this package's actions and nothing
			// else, so a name that is not here names an action in
			// another package or no action at all.
			r.errorf(actionsRel, RuleCompositeUnknownChild,
				"composite %s composes %q, which this package does not define; a composite may not reach into another package",
				a.Name, childName)
			continue
		}
		if child.Kind() == Composite {
			r.errorf(actionsRel, RuleCompositeNested,
				"composite %s composes %s, which is itself a composite; composites nest exactly one level",
				a.Name, childName)
			continue
		}

		for _, kw := range sortedKeys(child.Kwargs) {
			declared := child.Kwargs[kw]
			existing, seen := union[kw]
			if seen && existing.Type != declared.Type {
				// Same-name-same-type is assumed intentional.
				// Same-name-different-type is an accident that
				// cannot be recovered from at generation.
				r.errorf(actionsRel, RuleCompositeKwargConflict,
					"composite %s: children declare kwarg %s with different types (%s and %s); the caller supplies it once, so the two cannot both be satisfied",
					a.Name, kw, existing.Type, declared.Type)
				continue
			}
			// Union of names, union of required flags: if any child
			// requires it, the composite requires it.
			union[kw] = Kwarg{Type: declared.Type, Required: existing.Required || declared.Required}
		}
	}
	a.Kwargs = union
}

// checkFileTemplatePatterns rejects a package whose file templates cannot be
// ranked against each other. Silently picking a winner would make which
// boilerplate a file receives depend on directory iteration order.
func checkFileTemplatePatterns(patterns []string, r *reporter) {
	set, err := filetmpl.NewSet(patterns)
	if err != nil {
		r.errorf(filesDirName, RuleFileTemplateInvalid,
			"a file template's path is not a usable pattern: %v", err)
		return
	}
	for _, c := range set.Conflicts() {
		r.errorf(filesDirName, RuleFileTemplateTie,
			"file templates %s and %s tie under the specificity ranking; a path matching both would get its boilerplate from whichever the directory listed first",
			c.A, c.B)
	}
}

// checkTemplates reads every template and path pattern in the package with the
// renderer's own parser, and rejects two different mistakes.
//
// A template that is not Sedum's syntax at all is one. The grammar is {{name}}
// and {{name|op|op:arg}}; a conditional or a range would survive translation
// into text/template and work, which is precisely how a syntax drifts into a
// language nobody designed. Rejecting it here is the only place it can be
// caught before it is depended on (prov-2026-4675cebe).
//
// A well-formed reference to a transform that does not exist is the other. That
// is a hard error at load rather than at render because a package that cannot
// render is broken whether or not a given run happens to reach the template.
func checkTemplates(pkg *Package, fileTemplates map[string]string, actionTemplates map[string]string, r *reporter) {
	report := func(file string, source string, text string) {
		exprs, problems := render.Parse(text)
		for _, problem := range problems {
			r.errorf(file, RuleTemplateSyntaxInvalid, "%s: %v", source, problem)
		}
		if pkg.Engine == nil {
			// The vocabulary itself is broken and has already been
			// reported. Checking references against it would name every
			// template that used a pipeline rather than the pipeline.
			return
		}
		for _, expr := range exprs {
			for _, ref := range expr.Transforms {
				if err := pkg.Engine.Check(ref); err != nil {
					r.errorf(file, RuleTransformUndefined, "%s: %s: %v", source, expr.Source, err)
				}
			}
		}
	}

	for _, name := range sortedKeys(pkg.Actions) {
		a := pkg.Actions[name]
		if a.InjectsInto != "" {
			report(actionsRel, fmt.Sprintf("action %s: injects_into", name), a.InjectsInto)
		}
	}
	for _, path := range sortedKeys(actionTemplates) {
		report(path, "template "+path, actionTemplates[path])
	}
	for _, pattern := range pkg.FileTemplates {
		rel := filepath.ToSlash(filepath.Join(filesDirName, pattern))
		report(rel, "file template "+pattern, fileTemplates[pattern])
		// The pattern itself may carry transforms in its segments.
		report(rel, "file template path "+pattern, pattern)
	}
}

// checkMarkerAnchors warns when an action targets a marker no file template in
// its package plants.
//
// It is a warning, not an error. Template selection is path-dependent, so
// complete verification is impossible: an action may legitimately target a
// marker planted only by a template the current run never selects. A false hard
// error there would block legitimate packages, whereas a marker referenced by
// nothing at all is almost certainly a typo.
func checkMarkerAnchors(pkg *Package, fileTemplates []string, r *reporter) {
	planted := plantedMarkers(pkg.CommentPrefix, fileTemplates)

	for _, name := range sortedKeys(pkg.Actions) {
		marker, ok := pkg.Actions[name].MarkerAnchor()
		if !ok || planted[marker] {
			continue
		}
		r.warnf(actionsRel, RuleAnchorUnplanted,
			"action %s is anchored to marker %q, which no file template in this package plants; injection will fail for any file whose template lacks it",
			name, marker)
	}
}

// checkMarkersFilled warns when a file template plants a marker no action in
// the package targets.
//
// It is the mirror of checkMarkerAnchors, and the same defect checkDeadConfig
// reports about actions: an action nothing can invoke is dead configuration,
// and so is an injection point nothing can fill. A file generated from such a
// template carries a hole where an action was expected, and the package
// validates clean, so the gap surfaces wherever the generated code is first
// built or run - which is the expensive place.
//
// Unlike checkMarkerAnchors this check is complete. Verifying that an action's
// marker exists cannot be, because template selection is path-dependent and the
// template carrying the marker may never be chosen. Verifying that a marker is
// targeted compares two sets that are both fully known at load, with no path
// involved.
//
// A warning rather than an error, with no way to suppress it, for the reasons
// the dead-action warning has neither: the package is usable, the author may be
// mid-edit, and the way to silence it is to write the action - the same answer
// Sedum gives everywhere it asks for a declaration rather than inferring one.
func checkMarkersFilled(pkg *Package, fileTemplates []string, r *reporter) {
	targeted := map[string]bool{}
	for _, action := range pkg.Actions {
		if marker, ok := action.MarkerAnchor(); ok {
			targeted[marker] = true
		}
		// A region anchor names its endpoints through anchor_start and
		// anchor_end rather than anchor. Missing them would report every
		// region's endpoints as unfilled.
		if action.Anchor == AnchorRegion {
			targeted[action.AnchorStart] = true
			targeted[action.AnchorEnd] = true
		}
	}

	planted := plantedMarkers(pkg.CommentPrefix, fileTemplates)
	for _, marker := range sortedKeys(planted) {
		if targeted[marker] {
			continue
		}
		r.warnf(filesDirName, RuleMarkerUnfilled,
			"file templates plant marker %q, which no action in this package targets; nothing can inject there, so every file carrying it is generated with the anchor unused",
			marker)
	}
}

// checkDeadConfig warns about an action that is neither exposed to the model
// nor reachable through a composite, which is almost always a rename that
// missed a call site.
func checkDeadConfig(pkg *Package, r *reporter) {
	composed := map[string]bool{}
	for _, a := range pkg.Actions {
		if a.Kind() != Composite {
			continue
		}
		for _, child := range a.Composes {
			composed[child] = true
		}
	}

	for _, name := range sortedKeys(pkg.Actions) {
		if pkg.Actions[name].Exposed || composed[name] {
			continue
		}
		r.warnf(actionsRel, RuleActionDead,
			"action %s is not exposed and no composite references it, so nothing can invoke it", name)
	}
}

// checkUnmanaged validates a package's unmanaged declarations and rejects the
// two ways a package can contradict itself about them.
//
// An unreadable pattern is an error for the same reason an unreadable scope
// entry is: left alone it matches nothing, so a declaration meant to keep Sedum
// out of a file would quietly keep it out of nothing.
//
// The contradictions are errors rather than warnings because neither has a
// sensible outcome. A package that ships a file template for a path it also
// disowns has said both that it knows how to write the file and that it will
// not; an action whose injects_into names one has been pointed at a file its own
// package says nothing writes. Picking a winner would mean deciding which half
// of the package to believe.
func checkUnmanaged(pkg *Package, r *reporter) {
	for _, entry := range pkg.Unmanaged {
		if strings.TrimSpace(entry) == "" {
			r.errorf(manifestFile, RuleUnmanagedInvalid,
				"unmanaged carries an empty entry; an entry that names nothing keeps Sedum out of nothing")
			continue
		}
		if err := pathpat.Check(entry); err != nil {
			r.errorf(manifestFile, RuleUnmanagedInvalid, "unmanaged entry %v", err)
		}
	}

	for _, pattern := range pkg.FileTemplates {
		if entry, ok := pathpat.MatchAny(pkg.Unmanaged, pattern); ok {
			r.errorf(filesDirName, RuleUnmanagedContradiction,
				"file template %s matches unmanaged entry %q; the package both declares how to write this path and declares that it does not write it",
				pattern, entry)
		}
	}

	// Only a literal injects_into can be checked here. One carrying
	// placeholders renders per invocation, so Phase 6 is where it is caught,
	// and catching it there costs a run rather than a load.
	for _, name := range sortedKeys(pkg.Actions) {
		target := pkg.Actions[name].InjectsInto
		if target == "" || strings.Contains(target, "{{") {
			continue
		}
		if entry, ok := pathpat.MatchAny(pkg.Unmanaged, target); ok {
			r.errorf(actionsRel, RuleUnmanagedContradiction,
				"action %s injects into %s, which matches unmanaged entry %q; the package declares that it does not write that file",
				name, target, entry)
		}
	}
}
