// Package pipeline runs Sedum's phases in order.
//
// There is no plan artifact and nothing here decides anything. The execution
// sequence is a consequence of the record and the configuration, so this
// package's whole job is to run each phase, hand its output to the next as that
// phase's only input, and stop where it was told to stop.
//
// A phase that fails halts the run, which is what makes a stop point mean
// "everything before this is complete and nothing after it started" rather than
// "some of it happened".
//
// Phase 4 is the only phase that consults a model, and it is the only place
// this package is not a pure function of its inputs. Everything after it
// consumes the validated invocation list and the generator packages, which is
// why a stop at Phase 5 can be resumed from a recording and a stop before it
// cannot.
package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/calebcowen/sedum/internal/expand"
	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/inject"
	"github.com/calebcowen/sedum/internal/record"
	"github.com/calebcowen/sedum/internal/recording"
	"github.com/calebcowen/sedum/internal/resolve"
	"github.com/calebcowen/sedum/internal/runlog"
	"github.com/calebcowen/sedum/internal/selection"
)

// The phase boundaries a run can be halted at. The values are the phase numbers
// the PRD uses, so that a stop point named on the command line maps to one
// without a translation table.
const (
	PhaseLoad = iota
	PhaseIngest
	PhaseResolve
	PhaseCreate
	PhaseSelect
	PhaseValidate
	PhaseExpand
	PhaseInject
)

// Config is everything a run needs. It is filled in by the command and read
// nowhere else.
type Config struct {
	Generators string
	Records    string
	Output     string
	Lang       []string
	Only       []string

	// Variables are the run's values for the project facts packages declare,
	// as supplied on the command line. They are resolved against the loaded
	// packages' declarations before any phase writes anything
	// (prov-2026-6fc3d13d).
	Variables map[string]string

	// DryRun runs every phase and writes nothing.
	DryRun bool

	// StopAfterPhase halts the run after the named phase. Zero runs every
	// phase, since Phase 0 is not a stop point.
	StopAfterPhase int

	// Client is the model Phase 4 consults. It may be nil only for a run that
	// stops before Phase 4, which is what sedum resolve is.
	Client selection.Client

	// Retries bounds Phase 5's re-prompt loop.
	Retries int

	// Log is the run log. A nil log discards.
	Log *runlog.Log
}

// Result is what a run decided and what it did.
type Result struct {
	Packages    *genpkg.Set
	Records     *record.Set
	Resolutions []resolve.Resolution

	// Files is nil when the run stopped before Phase 3.
	Files []resolve.File

	// Selections is what the model chose for each record, validated. It is
	// nil when the run stopped before Phase 4, and it is the recording's
	// content: invocations are held pre-expansion, at the abstraction level
	// an author edits in.
	Selections []Selection

	// Injections is what Phase 7 wrote, or would have written under a dry
	// run. Nil when the run stopped before Phase 7.
	Injections []inject.Result

	// Variables are what the run's values resolved to, defaults filled in. They
	// are carried so that a recording can hold them: a recording that rendered
	// different text depending on invisible run state would not be one.
	Variables map[string]string

	// Unmanaged are the authorized paths a generator package declared Sedum
	// does not write. They are the run's handoff: authorized work that
	// something other than Sedum has to do.
	Unmanaged []resolve.Resolution

	// Warnings are collected from every phase. Where they go is the command's
	// decision, not this package's.
	Warnings []string

	// StoppedAfter is the phase the run halted at, or zero if it ran to the
	// end.
	StoppedAfter int
}

// Selection is one record's validated invocation list, with the files it was
// selected against.
//
// One record, one model call, one entry. The grouping is the recording's shape
// as well, which is not a coincidence: a recording is this list serialized.
type Selection struct {
	RecordID    string
	Files       []resolve.File
	Invocations []recording.Invocation

	// Calls, Rejected and Completeness are what Phase 5 spent reaching this
	// answer, and the token counts are what those calls cost as the server
	// accounted for it. They are carried rather than logged only, because "what
	// did this record cost" is a question about the run and not about one
	// phase's diagnostics (prov-2026-0811425c, prov-2026-096a4d4b). A record no
	// model was asked about carries zeroes, which is what it cost.
	Calls            int
	Rejected         int
	Completeness     int
	PromptTokens     int
	CompletionTokens int
}

// Run executes the phases in order.
func Run(ctx context.Context, cfg Config) (*Result, error) {
	log := cfg.Log
	if log == nil {
		log = runlog.Discard()
	}

	result := &Result{StoppedAfter: cfg.StopAfterPhase}

	// Phase 0 - load and validate generator packages.
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
		// Packages are wholly valid or rejected, and a run under a
		// half-loaded generators directory fails later as a missing feature
		// rather than as a broken package. Every finding is reported, not
		// just the first, so that a package is fixed in one pass.
		var lines []string
		for _, f := range findings {
			lines = append(lines, "  "+f.String())
		}
		return nil, fmt.Errorf("generator packages did not load:\n%s", strings.Join(lines, "\n"))
	}
	result.Packages = packages
	log.Info("loaded generator packages", "count", len(packages.Packages), "extensions", packages.Extensions())

	// Variables are bound before Phase 1, so a run missing one is refused while
	// nothing has been written and the diagnostic can name the package and the
	// author's description. Left to render, the same mistake surfaces in Phase 6
	// as a complaint about a template, after files are on disk.
	variables, err := packages.ResolveVariables(cfg.Variables)
	if err != nil {
		return nil, err
	}
	result.Variables = variables
	if len(variables) > 0 {
		log.Info("bound run variables", "variables", variables)
	}

	// Phase 1 - ingest provenance records.
	records, warnings, err := record.Load(cfg.Records, record.Options{Only: cfg.Only})
	result.Warnings = append(result.Warnings, warnings...)
	if err != nil {
		return nil, err
	}
	result.Records = records
	log.Info("ingested provenance records", "count", len(records.Records), "paths", len(records.Paths()))

	// Phase 2 - resolve paths to packages and to file templates.
	resolutions, warnings, err := resolve.Paths(packages, records, cfg.Lang)
	result.Warnings = append(result.Warnings, warnings...)
	if err != nil {
		return nil, err
	}
	result.Resolutions = resolutions
	for _, r := range resolutions {
		// An unmanaged path resolved to no package, by design: the check
		// runs before extension resolution so that a path no extension can
		// reach is still reportable.
		if r.Unmanaged {
			result.Unmanaged = append(result.Unmanaged, r)
			log.Info("left unmanaged", "path", r.Path,
				"declared_by", r.UnmanagedBy, "entry", r.UnmanagedAs)
			continue
		}
		log.Info("resolved path", "path", r.Path, "package", r.Package.Name,
			"template", r.Template, "default", r.Default, "captures", r.Captures)
	}

	if cfg.StopAfterPhase == PhaseResolve {
		log.Info("stopping after resolution")
		return result, nil
	}

	// Phase 3 - create files from templates.
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

	// Phases 4 and 5 - one model call per record, and deterministic
	// validation of what it returned.
	//
	// Per record rather than per run, because a record is the unit of intent:
	// its constraints govern its own paths, and one call deciding two records'
	// files would make each record's constraints apply to the other's.
	selections, err := selectAll(ctx, cfg, records, files, variables, log)
	if err != nil {
		return nil, err
	}
	result.Selections = selections

	if cfg.StopAfterPhase == PhaseValidate {
		log.Info("stopping after validated invocations")
		return result, nil
	}

	// Phase 6 - expand composites, render paths, select variants, apply
	// transforms. The model does not participate and nothing here is read
	// from it: every value comes from the validated kwargs.
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

	// Phase 7 - inject.
	applied, err := inject.Apply(resolved, inject.Options{
		Output:    cfg.Output,
		DryRun:    cfg.DryRun,
		Unwritten: unwritten(cfg, files),
		Log:       log,
	})
	if err != nil {
		return nil, err
	}
	result.Injections = applied
	return result, nil
}

// selectAll runs Phase 4 and Phase 5 once per record.
//
// A record whose paths all resolved to nothing Sedum writes is skipped without
// a model call. Its catalog would be empty, so the only valid answer is an
// empty list, and paying for a call to be told so would be a cost the run can
// see is pointless before it is incurred.
func selectAll(ctx context.Context, cfg Config, records *record.Set, files []resolve.File, variables map[string]string, log *runlog.Log) ([]Selection, error) {
	var out []Selection

	for _, rec := range records.Records {
		mine := filesOf(files, rec.ID)
		if !anyManaged(mine) {
			log.Info("no model call for record", "record", rec.ID,
				"reason", "no path it authorized is written by a generator package")
			out = append(out, Selection{RecordID: rec.ID, Files: mine})
			continue
		}

		if cfg.Client == nil {
			// A run reaching Phase 4 without a model is a caller mistake
			// rather than a user's, and it is named as one: the alternative
			// is a nil dereference at the point of the call.
			return nil, fmt.Errorf("record %s: reaching model invocation with no model configured", rec.ID)
		}

		answer, err := selection.Select(ctx, cfg.Client, selection.Request{
			RecordID:    rec.ID,
			Intent:      rec.Intent,
			Constraints: rec.Constraints,
			Files:       mine,
		}, selection.Options{Retries: cfg.Retries, Variables: variables, Log: log})
		if err != nil {
			return nil, err
		}

		out = append(out, Selection{
			RecordID:         rec.ID,
			Files:            mine,
			Invocations:      answer.Invocations,
			Calls:            answer.Calls,
			Rejected:         answer.Rejected,
			Completeness:     answer.Completeness,
			PromptTokens:     answer.PromptTokens,
			CompletionTokens: answer.CompletionTokens,
		})
	}

	return out, nil
}

// unwritten collects what Phase 3 rendered for files a dry run declined to
// create, so Phase 7 can decide against them without anything being written.
//
// It is empty for a real run, where every created file is on disk and injection
// reads it there (prov-2026-23653fdc).
func unwritten(cfg Config, files []resolve.File) map[string]string {
	if !cfg.DryRun {
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

// filesOf returns the files created for one record.
func filesOf(files []resolve.File, recordID string) []resolve.File {
	var out []resolve.File
	for _, f := range files {
		if f.RecordID == recordID {
			out = append(out, f)
		}
	}
	return out
}

// anyManaged reports whether any of a record's paths is one a generator package
// actually writes.
func anyManaged(files []resolve.File) bool {
	for _, f := range files {
		if !f.Unmanaged {
			return true
		}
	}
	return false
}
