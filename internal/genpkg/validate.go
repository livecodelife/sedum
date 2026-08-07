package genpkg

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/calebcowen/sedum/internal/filetmpl"
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

func buildAction(name string, decl *actionDecl, pkg *Package, r *reporter) *Action {
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

// checkTransformRefs rejects a reference that resolves to neither a built-in
// operation nor a pipeline the package declares. It is a hard error at load
// rather than at render, because a package that cannot render is broken whether
// or not a given run happens to reach the template.
func checkTransformRefs(pkg *Package, fileTemplates []string, actionTemplates map[string]string, r *reporter) {
	report := func(file string, source string, text string) {
		for _, ref := range transformRefs(text) {
			if pkg.resolves(ref) {
				continue
			}
			r.errorf(file, RuleTransformUndefined,
				"%s references transform %q, which is neither a built-in operation (%s) nor a pipeline this package declares (%s)",
				source, ref, strings.Join(BuiltinOperations, ", "), strings.Join(sortedKeys(pkg.Transforms), ", "))
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
	for i, pattern := range pkg.FileTemplates {
		rel := filepath.ToSlash(filepath.Join(filesDirName, pattern))
		report(rel, "file template "+pattern, fileTemplates[i])
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
