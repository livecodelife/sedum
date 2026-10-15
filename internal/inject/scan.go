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
	own := newOwnership(owners)

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
			if region.Marker.Record == "" || !own.owns(recordID, region.Marker.Record) {
				continue
			}
			out = append(out, OwnedRegion{Path: target.Path, Package: target.Package, Region: region})
		}
	}
	return out, nil
}
