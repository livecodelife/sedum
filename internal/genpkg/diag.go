package genpkg

import (
	"fmt"
	"sort"
)

// Diagnostics from load-time validation.
//
// Every finding names the package, the file it concerns, and the rule it
// violated, and validation reports everything wrong with a package rather than
// halting on the first problem, so an author fixes a package in one pass.

// Kind separates findings that reject a package from findings that only warn.
type Kind int

const (
	KindError Kind = iota
	KindWarning
)

func (k Kind) String() string {
	if k == KindWarning {
		return "warning"
	}
	return "error"
}

// Rule identifiers. They are stable slugs rather than message text, so that a
// caller can suppress or test for one rule without matching on prose.
const (
	RuleManifestMissing        = "manifest_missing"
	RuleManifestMalformed      = "manifest_malformed"
	RuleActionsMissing         = "actions_missing"
	RuleActionsMalformed       = "actions_malformed"
	RuleNameMismatch           = "package_name_mismatch"
	RuleExtensionInvalid       = "extension_invalid"
	RuleKwargTypeUnknown       = "kwarg_type_unknown"
	RuleActionIncomplete       = "action_incomplete"
	RuleDiscriminatorUnknown   = "discriminator_unknown_kwarg"
	RuleAnchorInvalid          = "anchor_invalid"
	RuleTemplateMissing        = "action_template_missing"
	RuleTemplateWrongForm      = "action_template_wrong_form"
	RuleTemplateAmbiguous      = "action_template_ambiguous"
	RuleVariantTemplateMissing = "variant_template_missing"
	RuleTransformUndefined     = "transform_undefined"
	RuleTransformInvalid       = "transform_invalid"
	RuleTemplateSyntaxInvalid  = "template_syntax_invalid"
	RuleCompositeNested        = "composite_nested"
	RuleCompositeUnknownChild  = "composite_unknown_child"
	RuleCompositeKwargConflict = "composite_kwarg_conflict"
	RuleCompositeMalformed     = "composite_malformed"
	RuleFileTemplateInvalid    = "file_template_invalid"
	RuleFileTemplateTie        = "file_template_tie"

	RuleActionNameReserved = "action_name_reserved"

	RuleAnchorUnplanted       = "anchor_marker_unplanted"
	RuleMarkerUnfilled        = "anchor_marker_unfilled"
	RuleAnchorPatternLineMode = "anchor_pattern_line_mode"
	RuleActionDead            = "action_dead"

	RuleUnmanagedInvalid       = "unmanaged_invalid"
	RuleUnmanagedContradiction = "unmanaged_contradiction"
)

// Finding is one load-time diagnostic.
type Finding struct {
	// Package is the package directory's name, which is what the author
	// types after --package.
	Package string
	// File is the path the finding concerns, relative to the package
	// directory. Empty when the finding is about the package as a whole.
	File    string
	Rule    string
	Message string
	Kind    Kind
}

func (f Finding) String() string {
	where := f.Package
	if f.File != "" {
		where += "/" + f.File
	}
	return fmt.Sprintf("%s: %s [%s] %s", f.Kind, where, f.Rule, f.Message)
}

// Findings is the full set of diagnostics from a load.
type Findings []Finding

func (fs Findings) HasErrors() bool {
	for _, f := range fs {
		if f.Kind == KindError {
			return true
		}
	}
	return false
}

// Strict returns the findings with every warning promoted to an error, which
// is what --strict asks for. The original is not modified.
func (fs Findings) Strict() Findings {
	out := make(Findings, len(fs))
	copy(out, fs)
	for i := range out {
		out[i].Kind = KindError
	}
	return out
}

// sorted orders findings so that output does not depend on map iteration or
// directory order: errors before warnings, then by package, file, and rule.
func (fs Findings) sorted() Findings {
	out := make(Findings, len(fs))
	copy(out, fs)
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		switch {
		case a.Kind != b.Kind:
			return a.Kind < b.Kind
		case a.Package != b.Package:
			return a.Package < b.Package
		case a.File != b.File:
			return a.File < b.File
		case a.Rule != b.Rule:
			return a.Rule < b.Rule
		default:
			return a.Message < b.Message
		}
	})
	return out
}

// reporter accumulates findings for one package.
type reporter struct {
	pkg      string
	findings Findings
}

func (r *reporter) errorf(file, rule, format string, args ...any) {
	r.add(KindError, file, rule, format, args...)
}

func (r *reporter) warnf(file, rule, format string, args ...any) {
	r.add(KindWarning, file, rule, format, args...)
}

func (r *reporter) add(kind Kind, file, rule, format string, args ...any) {
	r.findings = append(r.findings, Finding{
		Package: r.pkg,
		File:    file,
		Rule:    rule,
		Message: fmt.Sprintf(format, args...),
		Kind:    kind,
	})
}

func (r *reporter) hasErrors() bool { return r.findings.HasErrors() }
