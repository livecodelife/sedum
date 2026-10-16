package inject

import (
	"fmt"

	"github.com/livecodelife/sedum/internal/genpkg"
)

// OwnedRegion is one region a record owns, as found in a file.
type OwnedRegion struct {
	Path    string
	Package *genpkg.Package
	Region  Region
}

// ScanOwned returns the regions recordID owns in its targets: those whose
// marker names its id or one of its aliases (prov-2026-375f5d4a). Seeded
// regions are included and say so on their marker; a caller deciding what to
// regenerate filters them.
//
// It reads the files from disk or, in a dry run, from what Phase 3 would have
// written, and writes nothing. A target that does not exist owns nothing yet.
// This is what lets grow decide per record whether to skip, re-render or
// re-select from markers alone, with no manifest (prov-2026-9621daed).
func ScanOwned(recordID string, targets []Target, owners map[string]Owner, output string, unwritten map[string]string) ([]OwnedRegion, error) {
	all, err := ScanRegions(targets, output, unwritten)
	if err != nil {
		return nil, err
	}
	own := newOwnership(owners)
	var out []OwnedRegion
	for _, r := range all {
		if r.Region.Marker.Record != "" && own.owns(recordID, r.Region.Marker.Record) {
			out = append(out, r)
		}
	}
	return out, nil
}

// ScanRegions returns every marked region in the targets, whoever owns it.
func ScanRegions(targets []Target, output string, unwritten map[string]string) ([]OwnedRegion, error) {
	var out []OwnedRegion
	for _, target := range targets {
		content, ok := readExisting(output, unwritten, target.Path)
		if !ok {
			continue
		}
		regions, err := FindRegions(target.Package.Comment(), content)
		if err != nil {
			return nil, fmt.Errorf("file %s: %w", target.Path, err)
		}
		for _, region := range regions {
			out = append(out, OwnedRegion{Path: target.Path, Package: target.Package, Region: region})
		}
	}
	return out, nil
}

// Owns reports whether the record recordID owns a region whose marker names
// owner, directly or through an alias.
func Owns(owners map[string]Owner, recordID, owner string) bool {
	return owner != "" && newOwnership(owners).owns(recordID, owner)
}

// IsOrphan reports whether a region is an orphan: owned-tier, naming a record
// that is not among owners by id or alias (prov-2026-cc73cf61). owners must
// hold every current record - the records directory less deprecated ones -
// for the answer to mean anything. Seeded and unowned regions are never
// orphans.
func IsOrphan(owners map[string]Owner, m Marker) bool {
	return newOwnership(owners).orphan(m)
}

// SelectingKwargs returns the kwarg names that take part in a region's
// identity for an action: the ones it declares as its identity, or every
// required kwarg where it declares none.
func SelectingKwargs(action *genpkg.Action) []string {
	return selectingKwargs(action)
}
