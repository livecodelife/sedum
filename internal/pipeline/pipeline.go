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

	"github.com/livecodelife/sedum/internal/expand"
	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/record"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/resolve"
	"github.com/livecodelife/sedum/internal/runlog"
	"github.com/livecodelife/sedum/internal/selection"
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

	// Force lets a record whose selection came back empty prune every region
	// it owns. Without it that is an error, because an empty reply is more
	// often a flaky model than an intent to delete a record's code
	// (prov-2026-375f5d4a). It also re-selects every record whatever its
	// hashes say (prov-2026-9621daed).
	Force bool

	// Reselect re-selects every record whatever its hashes say, without
	// Force's licence to prune on an empty selection. A recording run sets
	// it: a recording captures selections, and a skipped record makes none.
	Reselect bool

	// Check classifies every record and stops: no model call, and nothing
	// written, Phase 3's files included (prov-2026-9621daed).
	Check bool

	// PruneOrphans removes orphan regions - owned regions whose record was
	// removed or deprecated - in files some current record authorizes
	// (prov-2026-cc73cf61).
	PruneOrphans bool

	// Stamp writes the current gen and pkg onto markers missing them, for
	// regions whose record exists, and does nothing else: no model call, no
	// other byte changed (prov-2026-cc73cf61).
	Stamp bool

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

	// FullCatalog disables Phase 4's authorized-path filter, embedding every
	// exposed action in the prompt whether or not any authorized path could
	// ever match it - the prompt every version before this one sent.
	FullCatalog bool

	// Prompt is caller-supplied text added to every selected record's prompt,
	// in the supplementary section. Blank or whitespace-only adds nothing.
	Prompt string

	// Prior is an already-read earlier recording. It is context only, never
	// replayed: the invocations of the records this run selects are shown to
	// the model, and nothing else in it is used - not its sedum_version, and
	// not its own Prompt. Run filters it, so --only narrows it.
	Prior *recording.Recording
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

	// Prompt is the caller's text as given, when it was non-blank. Capture
	// writes it into the recording.
	Prompt string

	// StoppedAfter is the phase the run halted at, or zero if it ran to the
	// end.
	StoppedAfter int

	// Classifications are what grow decided for each record from its markers
	// - select, re-select, skip or re-render - and why (prov-2026-9621daed).
	// Nil when the run stopped before they were made.
	Classifications []Classification

	// Orphans are the regions whose record is no longer current, and whether
	// this run pruned them (prov-2026-cc73cf61).
	Orphans []Orphan

	// Stamped are the markers --stamp wrote a missing hash onto.
	Stamped []inject.Stamped
}

// CheckFailed reports whether --check fails: some record is stale or needs
// re-rendering. A record with no regions is unknown and does not fail.
func (r *Result) CheckFailed() bool {
	for _, c := range r.Classifications {
		if c.Failing() {
			return true
		}
	}
	// An orphan is code no current record stands behind (prov-2026-cc73cf61).
	return len(r.Orphans) > 0
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
	if strings.TrimSpace(cfg.Prompt) != "" {
		result.Prompt = cfg.Prompt
	}

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

	// Phase 3 - create files from templates. --check writes nothing, so it
	// takes Phase 3's dry-run path and reads what would have been written.
	writeless := cfg.DryRun || cfg.Check || cfg.Stamp
	files, err := resolve.Create(resolutions, resolve.Options{
		Output:    cfg.Output,
		DryRun:    writeless,
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

	owners := ownersOf(records)

	// --stamp writes a missing gen or pkg onto the markers of each record's
	// regions and nothing else (prov-2026-cc73cf61).
	if cfg.Stamp {
		var requests []inject.StampRequest
		for _, rec := range records.Records {
			requests = append(requests, inject.StampRequest{
				RecordID: rec.ID, Gen: rec.Gen, Targets: targetsOf(filesOf(files, rec.ID)),
			})
		}
		stamped, err := inject.Stamp(requests, inject.Options{
			Output:    cfg.Output,
			DryRun:    cfg.DryRun,
			Unwritten: unwrittenFiles(writeless, files),
			Owners:    owners,
		})
		if err != nil {
			return nil, err
		}
		result.Stamped = stamped
		log.Info("stamped markers", "regions", len(stamped), "dry_run", cfg.DryRun)
		return result, nil
	}

	// Classify each record from its markers: select, re-select, skip, or
	// re-render. Only the first two reach the model (prov-2026-9621daed).
	classes, err := classify(records, files, owners, cfg.Output, unwrittenFiles(writeless, files), cfg.Force || cfg.Reselect)
	if err != nil {
		return nil, err
	}
	result.Classifications = classes
	needsModel := map[string]bool{}
	for _, c := range classes {
		needsModel[c.RecordID] = c.needsModel()
		log.Info("classified record", "record", c.RecordID, "mode", c.Mode.String(),
			"stale", len(c.Stale), "rerender", len(c.Rerender), "stale_seeded", len(c.StaleSeeded))
	}
	orphans, orphanTargets, err := findOrphans(packages, records, owners, cfg.Output)
	if err != nil {
		return nil, err
	}
	result.Orphans = orphans
	for _, o := range orphans {
		log.Info("orphan region", "path", o.Path, "action", o.Action, "record", o.Record, "authorized", o.Authorized)
	}
	if cfg.Check {
		log.Info("check complete; nothing written")
		return result, nil
	}
	if !cfg.PruneOrphans {
		orphanTargets = nil
	}

	// Phases 4 and 5 - one model call per record, and deterministic
	// validation of what it returned.
	//
	// Per record rather than per run, because a record is the unit of intent:
	// its constraints govern its own paths, and one call deciding two records'
	// files would make each record's constraints apply to the other's.
	selections, err := selectAll(ctx, cfg, records, files, variables, classes, owners, log)
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
		stampGen(expanded, genOf(records, s.RecordID))
		log.Info("expanded invocations", "record", s.RecordID,
			"selected", len(s.Invocations), "resolved", len(expanded))
		resolved = append(resolved, expanded...)
	}
	// A re-rendered record's invocations were built from its markers during
	// classification; they replace its pkg-stale regions in place.
	for _, c := range classes {
		if c.Mode == ModeRerender {
			log.Info("re-rendering regions from their markers", "record", c.RecordID, "regions", len(c.rerendered))
			resolved = append(resolved, c.rerendered...)
		}
	}

	if cfg.StopAfterPhase == PhaseExpand {
		log.Info("stopping after expansion")
		return result, nil
	}

	// Phase 7 - inject.
	// Phase 7 reconciles each record the model was asked about: its owned
	// regions it did not select again are pruned, in the same all-or-nothing
	// write as everything else.
	applied, err := inject.Apply(resolved, inject.Options{
		Output:         cfg.Output,
		DryRun:         cfg.DryRun,
		Unwritten:      unwrittenFiles(cfg.DryRun, files),
		Log:            log,
		Owners:         owners,
		Reconcile:      reconcileOf(selections, needsModel),
		PruneOrphansIn: orphanTargets,
		Force:          cfg.Force,
	})
	if err != nil {
		return nil, err
	}
	result.Injections = applied
	markPruned(result.Orphans, applied)
	return result, nil
}

// selectAll runs Phase 4 and Phase 5 once per record.
//
// A record whose paths all resolved to nothing Sedum writes is skipped without
// a model call. Its catalog would be empty, so the only valid answer is an
// empty list, and paying for a call to be told so would be a cost the run can
// see is pointless before it is incurred.
func selectAll(ctx context.Context, cfg Config, records *record.Set, files []resolve.File, variables map[string]string, classes []Classification, owners map[string]inject.Owner, log *runlog.Log) ([]Selection, error) {
	byRecord := map[string]Classification{}
	for _, c := range classes {
		byRecord[c.RecordID] = c
	}

	// No duplicate-path check: records sharing a file are reconciled region
	// by region in Phase 7, where a record selecting a region another record
	// owns is the error (prov-2026-375f5d4a).

	// The prior-selection bound is per record and is checked for every selected
	// record before the first call, so an overflow costs nothing. Selections
	// from a prior recording that belong to no record in this run are never
	// looked at.
	prior := map[string][]recording.Invocation{}
	for _, rec := range records.Records {
		prior[rec.ID] = priorFor(cfg.Prior, rec.ID)
		if err := selection.CheckPrior(rec.ID, prior[rec.ID]); err != nil {
			return nil, err
		}
		// A re-selected record with no prior recording sees its own earlier
		// selections, rebuilt from its markers, so the model has less reason
		// to change kwargs that did not need changing (prov-2026-cc73cf61).
		// Rebuilt context that would overflow the bound is left out rather
		// than failing a run the caller never asked to carry it.
		if c := byRecord[rec.ID]; len(prior[rec.ID]) == 0 && c.Mode == ModeReselect {
			rebuilt := priorFromMarkers(rec.ID, c.found, owners)
			if err := selection.CheckPrior(rec.ID, rebuilt); err != nil {
				log.Info("prior selections rebuilt from markers left out", "record", rec.ID, "reason", err.Error())
			} else {
				prior[rec.ID] = rebuilt
			}
		}
	}

	var out []Selection

	for _, rec := range records.Records {
		mine := filesOf(files, rec.ID)
		if anyManaged(mine) && !byRecord[rec.ID].needsModel() {
			// Classified skip or re-render: everything it needs is on its
			// markers (prov-2026-9621daed).
			log.Info("no model call for record", "record", rec.ID,
				"reason", "its regions are current against the record")
			out = append(out, Selection{RecordID: rec.ID, Files: mine})
			continue
		}
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
			Prompt:      cfg.Prompt,
			Prior:       prior[rec.ID],
			Present:     presentFor(rec.ID, byRecord[rec.ID].found, owners),
		}, selection.Options{Retries: cfg.Retries, Variables: variables, Log: log, FullCatalog: cfg.FullCatalog, Output: cfg.Output})
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

// priorFor returns the invocations a prior recording holds for one record,
// across its phases. Only the invocations are taken: the prior recording's own
// prompt is deliberately never read.
func priorFor(prior *recording.Recording, recordID string) []recording.Invocation {
	if prior == nil {
		return nil
	}
	var out []recording.Invocation
	for _, r := range prior.Records {
		if r.RecordID != recordID {
			continue
		}
		for _, ph := range r.Phases {
			out = append(out, ph.Invocations...)
		}
	}
	return out
}

// stampGen sets the gen_hash of the record that selected them on a record's
// expanded invocations, so every region it writes says which version of the
// record wrote it (prov-2026-66bea9f3).
func stampGen(invocations []inject.Invocation, gen string) {
	for i := range invocations {
		invocations[i].Gen = gen
	}
}

// genOf returns the gen_hash of the record with the given id, or nothing when
// the run has no such record.
//
// Nothing is not a gap to fill with a guess: an empty gen means unknown, and a
// later grow reads unknown as stale.
func genOf(records *record.Set, recordID string) string {
	if records == nil {
		return ""
	}
	rec, ok := records.Lookup(recordID)
	if !ok {
		return ""
	}
	return rec.Gen
}

// ownersOf is what Phase 7 knows about every record in the directory, for
// region ownership and anchor order. Records --only left out are included,
// because the regions a run meets may belong to any of them.
func ownersOf(records *record.Set) map[string]inject.Owner {
	if records == nil {
		return nil
	}
	out := make(map[string]inject.Owner, len(records.Known))
	for _, r := range records.Known {
		// A deprecated record owns nothing: its regions are orphans
		// (prov-2026-cc73cf61).
		if r.Deprecated {
			continue
		}
		out[r.ID] = inject.Owner{ID: r.ID, Aliases: r.Aliases, CreatedAt: r.CreatedAt}
	}
	return out
}

// reconcileOf names the records whose owned regions Phase 7 reconciles: every
// record the model was asked about, over the managed files it authorizes. A
// record that made no selection - nothing to grow, skipped, or re-rendered -
// reconciles nothing.
func reconcileOf(selections []Selection, needsModel map[string]bool) []inject.Reconcile {
	var out []inject.Reconcile
	for _, s := range selections {
		if !anyManaged(s.Files) || !needsModel[s.RecordID] {
			continue
		}
		out = append(out, inject.Reconcile{RecordID: s.RecordID, Targets: targetsOf(s.Files)})
	}
	return out
}

// presentFor lists the regions other records own in a record's files, for its
// prompt (prov-2026-cc73cf61). Each names its action, the kwargs that identify
// it, and its owner - never a hash.
func presentFor(recordID string, found []inject.OwnedRegion, owners map[string]inject.Owner) []selection.Present {
	var out []selection.Present
	for _, r := range found {
		m := r.Region.Marker
		if m.Record == "" || inject.Owns(owners, recordID, m.Record) {
			continue
		}
		kwargs := m.Kwargs
		if action, ok := r.Package.Actions[m.Action]; ok {
			kwargs = map[string]any{}
			for _, name := range inject.SelectingKwargs(action) {
				if v, ok := m.Kwargs[name]; ok {
					kwargs[name] = v
				}
			}
		}
		out = append(out, selection.Present{Path: r.Path, Action: m.Action, Variant: m.Variant, Kwargs: kwargs, Owner: m.Record})
	}
	return out
}

// priorFromMarkers rebuilds a record's earlier selections from the markers of
// the regions it owns (prov-2026-cc73cf61). Only exposed actions are kept - a
// composite's regions name unexposed children the model cannot select - and
// only the kwargs each action declares, since recorded defaults are the
// package's and recorded variables are the run's, not the model's to bind.
func priorFromMarkers(recordID string, found []inject.OwnedRegion, owners map[string]inject.Owner) []recording.Invocation {
	var out []recording.Invocation
	for _, r := range found {
		m := r.Region.Marker
		if !inject.Owns(owners, recordID, m.Record) {
			continue
		}
		action, ok := r.Package.Actions[m.Action]
		if !ok || !action.Exposed {
			continue
		}
		kwargs := map[string]any{}
		for name, value := range m.Kwargs {
			if _, declared := action.Kwargs[name]; declared {
				kwargs[name] = value
			}
		}
		out = append(out, recording.Invocation{Action: m.Action, Kwargs: kwargs})
	}
	return out
}
