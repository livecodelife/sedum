package pipeline

import (
	"errors"
	"fmt"

	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/record"
)

// RegionsFormat names the JSON document `sedum regions --json` prints. The
// version changes only when a field changes meaning (conformance/formats).
const RegionsFormat = "sedum.regions/1"

// RegionsConfig is what `sedum regions` reads (prov-2026-bf01cd6b).
type RegionsConfig struct {
	Generators string
	Output     string
	// Records is optional. Given, a marker naming an alias is reported under
	// the surviving id, and Only matches either.
	Records string
	// Only names the records whose regions to list. Empty lists every record.
	Only []string
	// Exclude adds directories to the ones the walk skips, besides the
	// vendored names and the exclude list in .sedum/config.yml
	// (prov-2026-fc05416c).
	Exclude []string
}

// RegionEntry is one region whose marker names a record.
type RegionEntry struct {
	Record   string         `json:"record"`
	File     string         `json:"file"`
	Action   string         `json:"action"`
	Variant  string         `json:"variant"`
	Tier     string         `json:"tier"`
	Identity map[string]any `json:"identity"`
	Gen      string         `json:"gen"`
	Pkg      string         `json:"pkg"`
}

// RegionsReport is the sedum.regions/1 document.
type RegionsReport struct {
	Format  string        `json:"format"`
	Regions []RegionEntry `json:"regions"`
}

// Regions lists every region whose marker names a record, read from markers
// under the output directory (prov-2026-bf01cd6b). It is read-only: no model,
// no record required, nothing written. It is what the provenance library's
// evidence code hash is computed from, because that library never parses a
// marker.
//
// gen and pkg are reported exactly as the marker carries them, "" when it does
// not: a consumer reads that as unknown, never as current, and grow --stamp is
// how a marker gains them.
func Regions(cfg RegionsConfig) (*RegionsReport, error) {
	rules, err := projectSkipRules(cfg.Output, cfg.Exclude)
	if err != nil {
		return nil, err
	}

	packages, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{})
	if err != nil {
		return nil, err
	}
	if findings.HasErrors() {
		var problems []error
		for _, f := range findings {
			if f.Kind == genpkg.KindError {
				problems = append(problems, errors.New(f.String()))
			}
		}
		return nil, fmt.Errorf("generator packages did not load:\n%w", errors.Join(problems...))
	}

	var owners map[string]inject.Owner
	if cfg.Records != "" {
		records, _, err := record.Load(cfg.Records, record.Options{})
		if err != nil {
			return nil, err
		}
		owners = ownersOf(records)
	}

	only := map[string]bool{}
	for _, id := range cfg.Only {
		only[inject.Canonical(owners, id)] = true
	}

	report := &RegionsReport{Format: RegionsFormat, Regions: []RegionEntry{}}
	err = walkMarked(packages, cfg.Output, []string{cfg.Generators, cfg.Records}, rules, func(rel string, pkg *genpkg.Package, regions []inject.Region) {
		for _, region := range regions {
			m := region.Marker
			if m.Record == "" {
				continue
			}
			owner := inject.Canonical(owners, m.Record)
			if len(only) > 0 && !only[owner] {
				continue
			}
			report.Regions = append(report.Regions, RegionEntry{
				Record:   owner,
				File:     rel,
				Action:   m.Action,
				Variant:  m.Variant,
				Tier:     string(m.Tier),
				Identity: identityKwargs(pkg, m),
				Gen:      m.Gen,
				Pkg:      m.Pkg,
			})
		}
	})
	if err != nil {
		return nil, err
	}
	return report, nil
}

// identityKwargs is the kwargs that identify a region: the action's selecting
// kwargs, or every recorded kwarg when the package no longer declares it.
func identityKwargs(pkg *genpkg.Package, m inject.Marker) map[string]any {
	out := map[string]any{}
	action, ok := pkg.Actions[m.Action]
	if !ok {
		for k, v := range m.Kwargs {
			out[k] = v
		}
		return out
	}
	for _, name := range inject.SelectingKwargs(action) {
		if v, ok := m.Kwargs[name]; ok {
			out[name] = v
		}
	}
	return out
}
