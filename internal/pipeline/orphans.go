package pipeline

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/record"
)

// Orphan is an owned-tier region whose record is no longer current: removed
// from the records directory, or deprecated (prov-2026-cc73cf61).
type Orphan struct {
	Path    string
	Action  string
	Variant string
	Record  string
	// Authorized reports that some current record authorizes the file, so
	// --prune-orphans may remove the region. Sedum writes only authorized
	// paths, so an orphan anywhere else is reported and left alone.
	Authorized bool
	// Pruned reports that this run removed it.
	Pruned bool
}

// findOrphans walks the output tree for orphan regions.
//
// A record's regions can only be found by reading files, and an orphan's record
// is exactly the one that no longer says which files those are - so the walk
// covers every file whose extension a loaded package claims, skipping hidden
// directories. A file whose extension more than one package claims is read
// with each package's comment syntax until one finds regions.
//
// It returns the orphans and, for --prune-orphans, the authorized files they
// sit in as Phase 7 targets.
func findOrphans(packages *genpkg.Set, records *record.Set, owners map[string]inject.Owner, output string) ([]Orphan, []inject.Target, error) {
	current := currentRecords(records)

	var (
		orphans []Orphan
		targets []inject.Target
	)
	err := filepath.WalkDir(output, func(full string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if full != output && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		candidates := packages.ForExtension(path.Ext(d.Name()))
		if len(candidates) == 0 {
			return nil
		}
		rel, err := filepath.Rel(output, full)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		data, err := os.ReadFile(full)
		if err != nil {
			return err
		}

		for _, pkg := range candidates {
			regions, err := inject.FindRegions(pkg.Comment(), string(data))
			if err != nil || len(regions) == 0 {
				continue
			}
			authorized := authorizedBy(current, rel)
			found := false
			for _, region := range regions {
				m := region.Marker
				if !inject.IsOrphan(owners, m) {
					continue
				}
				found = true
				orphans = append(orphans, Orphan{
					Path: rel, Action: m.Action, Variant: m.Variant, Record: m.Record, Authorized: authorized,
				})
			}
			if found && authorized {
				targets = append(targets, inject.Target{Path: rel, Package: pkg})
			}
			break
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	sort.SliceStable(orphans, func(i, j int) bool { return orphans[i].Path < orphans[j].Path })
	return orphans, targets, nil
}

// currentRecords is every record in the directory that is not deprecated,
// whether or not this run grows it.
func currentRecords(records *record.Set) []*record.Record {
	var out []*record.Record
	for _, r := range records.Known {
		if !r.Deprecated {
			out = append(out, r)
		}
	}
	return out
}

func authorizedBy(records []*record.Record, rel string) bool {
	for _, r := range records {
		if r.Authorizes(rel) {
			return true
		}
	}
	return false
}

// markPruned records which orphans Phase 7 removed.
func markPruned(orphans []Orphan, applied []inject.Result) {
	removed := map[[4]string]int{}
	for _, r := range applied {
		if r.Orphan {
			removed[[4]string{r.Path, r.Action, r.Variant, r.Record}]++
		}
	}
	for i, o := range orphans {
		key := [4]string{o.Path, o.Action, o.Variant, o.Record}
		if removed[key] > 0 {
			orphans[i].Pruned = true
			removed[key]--
		}
	}
}
