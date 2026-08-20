package pipeline

import (
	"fmt"
	"sort"
	"strings"

	"github.com/livecodelife/sedum/internal/expand"
	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/record"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/resolve"
	"github.com/livecodelife/sedum/internal/runlog"
	"github.com/livecodelife/sedum/internal/selection"
)

// ReplayConfig is everything a replay needs.
//
// It has no Client and no Retries. Replay skips Phase 4, so there is no model
// to configure, and its validation is terminal, so there is no budget to bound.
// Their absence is the type saying so rather than a comment saying so.
type ReplayConfig struct {
	Generators string
	Output     string

	// Records is optional. Supplied, every path the recording names is checked
	// against affected_scope and forbidden_scope. Omitted, the recording
	// executes as written - a hand-edited generic scaffold has no
	// corresponding records, and that is a legitimate mode rather than a
	// degraded one.
	Records string

	// Variables override what the recording carries, for a replay generating a
	// new service from a committed scaffold. Unset values fall back to the
	// recording's, then to the packages' declared defaults.
	Variables map[string]string

	DryRun         bool
	StopAfterPhase int
	Log            *runlog.Log
}

// Replay executes a recording.
//
// It enters at Phase 3 with resolution already decided, skips Phase 4 entirely,
// and runs Phases 5 through 7 unchanged. What makes that sound is that every
// phase after model invocation is a pure function of the validated invocation
// list and the generator packages, so a recording holds everything Phase 4
// would have produced.
//
// Failures are terminal. A recording naming a nonexistent action, omitting a
// required kwarg, or targeting a file it did not create fails exactly the
// checks a model response fails, and the run halts because there is nothing to
// re-prompt.
func Replay(rec recording.Recording, cfg ReplayConfig) (*Result, error) {
	log := cfg.Log
	if log == nil {
		log = runlog.Discard()
	}

	result := &Result{StoppedAfter: cfg.StopAfterPhase}

	// Phase 0 - load and validate generator packages, exactly as a grow run
	// does. A recording is not a lockfile: it captures decisions, not
	// templates, so the packages on disk are the ones that render.
	packages, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{})
	if err != nil {
		return nil, err
	}
	var rejected bool
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			rejected = true
			continue
		}
		result.Warnings = append(result.Warnings, f.String())
	}
	if rejected {
		var lines []string
		for _, f := range findings {
			lines = append(lines, "  "+f.String())
		}
		return nil, fmt.Errorf("generator packages did not load:\n%s", strings.Join(lines, "\n"))
	}
	result.Packages = packages
	log.Info("loaded generator packages", "count", len(packages.Packages), "extensions", packages.Extensions())

	if err := checkPackageIdentity(rec, packages); err != nil {
		return nil, err
	}

	// The recording's variables are the run's, unless this replay overrides
	// them. Resolving against the packages fills defaults for a variable a
	// package declared after the recording was written, which is the same
	// "picks up the new templates" property that makes a recording not a
	// lockfile.
	variables, err := packages.ResolveVariables(mergeVariables(rec.Variables, cfg.Variables))
	if err != nil {
		return nil, err
	}
	result.Variables = variables
	if len(variables) > 0 {
		log.Info("bound run variables", "variables", variables)
	}

	// Phase 1 - ingest records, only when scope validation was asked for.
	//
	// No duplicate-path check. Its justification is Phase 4's one-call-per-
	// record shape and replay never reaches Phase 4, so two records
	// legitimately refining regions in one file replay without complaint
	// (prov-2026-dc227be7).
	if cfg.Records != "" {
		loaded, warnings, err := record.Load(cfg.Records, record.Options{})
		result.Warnings = append(result.Warnings, warnings...)
		if err != nil {
			return nil, err
		}
		result.Records = loaded
		if err := loaded.Authorize(recordedPaths(rec)); err != nil {
			return nil, fmt.Errorf("the recording names paths these records do not authorize:\n%w", err)
		}
		log.Info("validated recorded paths against records", "records", len(loaded.Records))
	}

	// Phase 2 is skipped: resolution is what the recording carries.
	resolutions, err := recordedResolutions(rec, packages)
	if err != nil {
		return nil, err
	}
	result.Resolutions = resolutions

	// There is no stop after resolution here. Replay is handed resolution
	// rather than computing it, so a boundary "after" it is not one this path
	// has - which is what the stop-point table means by runByReplay, and why
	// the command refuses the combination rather than silently accepting it.

	// Phase 3 - create files from the recorded resolution.
	files, err := resolve.Create(resolutions, resolve.Options{
		Output:    cfg.Output,
		DryRun:    cfg.DryRun,
		Variables: variables,
		Log:       log,
	})
	if err != nil {
		return nil, err
	}
	result.Files = files

	if cfg.StopAfterPhase == PhaseCreate {
		log.Info("stopping after file creation")
		return result, nil
	}

	// Phase 5 - validation, identical to model-output validation and terminal.
	selections, err := validateRecorded(rec, files, log)
	if err != nil {
		return nil, err
	}
	result.Selections = selections

	if cfg.StopAfterPhase == PhaseValidate {
		log.Info("stopping after validated invocations")
		return result, nil
	}

	// Phase 6 - expansion, unchanged. Invocations were recorded pre-expansion
	// because expansion is deterministic and re-running it is free.
	var resolved []inject.Invocation
	for _, s := range selections {
		expanded, err := expand.Expand(s.RecordID, s.Files, s.Invocations, variables)
		if err != nil {
			return nil, fmt.Errorf("record %s: %w", s.RecordID, err)
		}
		log.Info("expanded invocations", "record", s.RecordID,
			"selected", len(s.Invocations), "resolved", len(expanded))
		resolved = append(resolved, expanded...)
	}

	if cfg.StopAfterPhase == PhaseExpand {
		log.Info("stopping after expansion")
		return result, nil
	}

	// Phase 7 - inject, unchanged.
	applied, err := inject.Apply(resolved, inject.Options{
		Output:    cfg.Output,
		DryRun:    cfg.DryRun,
		Unwritten: unwrittenFiles(cfg.DryRun, files),
		Log:       log,
	})
	if err != nil {
		return nil, err
	}
	result.Injections = applied
	return result, nil
}

// checkPackageIdentity verifies every package a recording names is present and
// still claims the extensions recorded against it.
//
// A recording that resolved .rb to one package, replayed against a directory
// where another claims it, would generate under the wrong conventions and
// succeed at it. That is the failure this exists to catch, and it is worth
// halting for because nothing downstream would notice.
func checkPackageIdentity(rec recording.Recording, packages *genpkg.Set) error {
	var problems []string

	for _, name := range sortedPackageNames(rec.Packages) {
		recorded := rec.Packages[name]
		pkg, ok := packages.Lookup(name)
		if !ok {
			problems = append(problems, fmt.Sprintf(
				"package %s is named by the recording but is not in the generators directory", name))
			continue
		}
		if missing := missingExtensions(recorded.Extensions, pkg.Extensions); len(missing) > 0 {
			problems = append(problems, fmt.Sprintf(
				"package %s no longer claims %s; it was recorded claiming %s and now claims %s",
				name, strings.Join(missing, ", "),
				strings.Join(recorded.Extensions, ", "), strings.Join(pkg.Extensions, ", ")))
		}
	}

	if len(problems) > 0 {
		return fmt.Errorf("the recording does not match the generator packages on disk:\n  %s",
			strings.Join(problems, "\n  "))
	}
	return nil
}

// missingExtensions returns the recorded extensions a package no longer claims.
func missingExtensions(recorded, current []string) []string {
	have := map[string]bool{}
	for _, e := range current {
		have[e] = true
	}
	var missing []string
	for _, e := range recorded {
		if !have[e] {
			missing = append(missing, e)
		}
	}
	return missing
}

// recordedResolutions rebuilds Phase 2's output from the recording.
//
// This is the "enters at Phase 3 with resolution already decided" step: the
// recording says which package claimed each path, which file template matched,
// and what that template's captures bound, so none of it is recomputed.
func recordedResolutions(rec recording.Recording, packages *genpkg.Set) ([]resolve.Resolution, error) {
	var out []resolve.Resolution

	for _, r := range rec.Records {
		for _, f := range r.Files {
			pkg, ok := packages.Lookup(f.Package)
			if !ok {
				// checkPackageIdentity has already reported a package the
				// recording declared. This catches a file naming one the
				// recording never declared at all, which is a hand-edit
				// authored against nothing.
				return nil, fmt.Errorf(
					"record %s: file %s names package %s, which the recording does not declare and the generators directory does not contain",
					r.RecordID, f.Path, f.Package)
			}
			out = append(out, resolve.Resolution{
				RecordID: r.RecordID,
				Path:     f.Path,
				Package:  pkg,
				Template: f.Template,
				Captures: capturesOrEmpty(f.Captures),
			})
		}
	}

	return out, nil
}

// validateRecorded runs Phase 5 over the recording's invocations.
//
// One entry per record, matching the grow path's shape, because a recording is
// that list serialized. The counts a Selection carries are zero throughout: no
// call was made, nothing was rejected and re-prompted, and reporting anything
// else would attribute a cost to a run that did not incur one.
func validateRecorded(rec recording.Recording, files []resolve.File, log *runlog.Log) ([]Selection, error) {
	var out []Selection

	for _, r := range rec.Records {
		mine := filesOf(files, r.RecordID)

		// Phases are executed in order. This implementation writes exactly one,
		// named default, but the level is reserved and replay honors it rather
		// than assuming the count.
		var invocations []recording.Invocation
		for _, phase := range r.Phases {
			invocations = append(invocations, phase.Invocations...)
		}

		if violations := selection.Validate(mine, invocations); len(violations) > 0 {
			return nil, &InvalidRecording{RecordID: r.RecordID, Violations: violations}
		}
		log.Info("recording validated", "record", r.RecordID, "invocations", len(invocations))

		out = append(out, Selection{
			RecordID:    r.RecordID,
			Files:       mine,
			Invocations: invocations,
		})
	}

	return out, nil
}

// InvalidRecording is a recording that did not validate.
//
// It is a type rather than a formatted string because the recording is an input
// format: a caller submitting one wants the violations as data, and matching on
// text is what a submission protocol should never require of its clients.
type InvalidRecording struct {
	RecordID   string
	Violations []selection.Violation
}

func (e *InvalidRecording) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "record %s: the recording did not validate", e.RecordID)
	// No attempt number and no retry count. There was one submission and there
	// is nothing to re-prompt, which is the difference between this and a
	// model's rejection.
	for _, v := range e.Violations {
		fmt.Fprintf(&b, "\n    %s", v)
	}
	return b.String()
}

// recordedPaths returns every path a recording names, for the scope check.
func recordedPaths(rec recording.Recording) []string {
	var out []string
	for _, r := range rec.Records {
		for _, f := range r.Files {
			out = append(out, f.Path)
		}
	}
	return out
}

// mergeVariables layers a replay's overrides over what the recording carries.
func mergeVariables(recorded, overrides map[string]string) map[string]string {
	if len(recorded) == 0 && len(overrides) == 0 {
		return nil
	}
	out := make(map[string]string, len(recorded)+len(overrides))
	for k, v := range recorded {
		out[k] = v
	}
	for k, v := range overrides {
		out[k] = v
	}
	return out
}

func capturesOrEmpty(m map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return m
}

func sortedPackageNames(m map[string]recording.Package) []string {
	out := make([]string, 0, len(m))
	for name := range m {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// unwrittenFiles is Run's unwritten, taking the flag rather than the config so
// both paths share one implementation of the rule.
func unwrittenFiles(dryRun bool, files []resolve.File) map[string]string {
	if !dryRun {
		return nil
	}
	out := map[string]string{}
	for _, f := range files {
		if f.Existed || f.Unmanaged {
			continue
		}
		out[f.Path] = f.Rendered
	}
	return out
}
