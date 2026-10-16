package cli

import (
	"sort"

	"github.com/livecodelife/sedum/internal/pipeline"
)

// growFormat names the JSON document `sedum grow --json` prints. The version
// changes only when a field changes meaning (conformance/formats).
const growFormat = "sedum.grow/1"

type growReport struct {
	Format      string         `json:"format"`
	DryRun      bool           `json:"dry_run"`
	Records     []recordReport `json:"records"`
	Orphans     []orphanReport `json:"orphans"`
	Stamped     []stampReport  `json:"stamped"`
	CheckFailed bool           `json:"check_failed"`
}

type recordReport struct {
	Record  string `json:"record"`
	Mode    string `json:"mode"`
	Changed bool   `json:"changed"`
	// RepairAttempts is how many times the project's verification sent the
	// record's selection back; omitted when none (prov-2026-b5c21813).
	RepairAttempts int            `json:"repair_attempts,omitempty"`
	Regions        []regionChange `json:"regions"`
}

type regionChange struct {
	File    string `json:"file"`
	Action  string `json:"action"`
	Variant string `json:"variant"`
	Change  string `json:"change"`
}

type orphanReport struct {
	File       string `json:"file"`
	Action     string `json:"action"`
	Variant    string `json:"variant"`
	Record     string `json:"record"`
	Authorized bool   `json:"authorized"`
	Pruned     bool   `json:"pruned"`
}

type stampReport struct {
	File    string `json:"file"`
	Action  string `json:"action"`
	Variant string `json:"variant"`
	Record  string `json:"record"`
	Gen     bool   `json:"gen"`
	Pkg     bool   `json:"pkg"`
}

// buildGrowReport is the sedum.grow/1 document for a run (prov-2026-bf01cd6b).
// A record's changed is true exactly when one of its regions was inserted,
// pruned, or rewritten with different bytes, so a caller runs only the specs
// whose code moved.
func buildGrowReport(result *pipeline.Result, dryRun, replay, check bool) growReport {
	report := growReport{
		Format:  growFormat,
		DryRun:  dryRun || check,
		Records: []recordReport{},
		Orphans: []orphanReport{},
		Stamped: []stampReport{},
	}

	byRecord := map[string]*recordReport{}
	add := func(id, mode string) *recordReport {
		if r, ok := byRecord[id]; ok {
			return r
		}
		r := &recordReport{Record: id, Mode: mode, Regions: []regionChange{}}
		byRecord[id] = r
		return r
	}
	for _, c := range result.Classifications {
		add(c.RecordID, c.Mode.String())
	}
	if replay {
		for _, s := range result.Selections {
			add(s.RecordID, "replay")
		}
	}

	for _, s := range result.Selections {
		if s.Repairs > 0 {
			add(s.RecordID, "").RepairAttempts = s.Repairs
		}
	}

	for _, inj := range result.Injections {
		if inj.Orphan {
			continue // reported under orphans: its record is gone
		}
		r := add(inj.Record, "")
		change := "inserted"
		switch {
		case inj.Held:
			change = "held"
		case inj.Pruned:
			change = "pruned"
		case inj.Skipped:
			change = "seeded"
		case inj.Replaced && !inj.Changed:
			change = "unchanged"
		case inj.Replaced:
			change = "replaced"
		}
		r.Regions = append(r.Regions, regionChange{File: inj.Path, Action: inj.Action, Variant: inj.Variant, Change: change})
		if inj.Changed {
			r.Changed = true
		}
	}

	for _, r := range byRecord {
		if r.Mode == "nothing to grow" {
			r.Mode = "none"
		}
		report.Records = append(report.Records, *r)
	}
	sort.Slice(report.Records, func(i, j int) bool { return report.Records[i].Record < report.Records[j].Record })

	for _, o := range result.Orphans {
		report.Orphans = append(report.Orphans, orphanReport{
			File: o.Path, Action: o.Action, Variant: o.Variant, Record: o.Record, Authorized: o.Authorized, Pruned: o.Pruned,
		})
	}
	for _, s := range result.Stamped {
		report.Stamped = append(report.Stamped, stampReport{
			File: s.Path, Action: s.Action, Variant: s.Variant, Record: s.Record, Gen: s.Gen, Pkg: s.Pkg,
		})
	}
	if check {
		report.CheckFailed = result.CheckFailed()
	}
	return report
}
