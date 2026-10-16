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
