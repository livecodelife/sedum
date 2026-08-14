// Package pipeline runs Sedum's phases in order.
//
// There is no plan artifact and nothing here decides anything. The execution
// sequence is a consequence of the record and the configuration, so this
// package's whole job is to run each phase, hand its output to the next as that
// phase's only input, and stop where it was told to stop.
//
// Phases 0 through 3 are implemented. A phase that fails halts the run, which
// is what makes a stop point mean "everything before this is complete and
// nothing after it started" rather than "some of it happened".
package pipeline

import (
	"fmt"
	"strings"

	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/record"
	"github.com/calebcowen/sedum/internal/resolve"
	"github.com/calebcowen/sedum/internal/runlog"
)

// The phase boundaries a run can be halted at. The values are the phase numbers
// the PRD uses, so that a stop point named on the command line maps to one
// without a translation table.
const (
	PhaseLoad = iota
	PhaseIngest
	PhaseResolve
	PhaseCreate
)

// Config is everything a run needs. It is filled in by the command and read
// nowhere else.
type Config struct {
	Generators string
	Records    string
	Output     string
	Lang       []string
	Only       []string

	// DryRun runs every phase and writes nothing.
	DryRun bool

	// StopAfterPhase halts the run after the named phase. Zero runs every
	// phase this milestone implements, since Phase 0 is not a stop point.
	StopAfterPhase int

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

	// Unmanaged are the authorized paths a generator package declared Sedum
	// does not write. They are the run's handoff: authorized work that
	// something other than Sedum has to do.
	Unmanaged []resolve.Resolution

	// Warnings are collected from every phase. Where they go is the command's
	// decision, not this package's.
	Warnings []string

	// StoppedAfter is the phase the run halted at, or zero if it ran to the
	// end of what is implemented.
	StoppedAfter int
}

// Run executes the phases in order.
func Run(cfg Config) (*Result, error) {
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
		Output: cfg.Output,
		DryRun: cfg.DryRun,
		Log:    log,
	})
	if err != nil {
		return nil, err
	}
	result.Files = files

	if cfg.StopAfterPhase == PhaseCreate {
		log.Info("stopping after file creation")
	}
	return result, nil
}
