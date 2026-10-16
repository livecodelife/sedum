package inject

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// StampRequest is one record whose regions --stamp brings up to date.
type StampRequest struct {
	RecordID string
	// Gen is the record's current gen_hash.
	Gen     string
	Targets []Target
}

// Stamped is one marker --stamp wrote, and which of the two hashes it added.
type Stamped struct {
	Path    string
	Action  string
	Variant string
	Record  string
	Gen     bool
	Pkg     bool
}

// Stamp writes the current gen and pkg onto the markers of each record's
// regions that are missing them, and changes nothing else (prov-2026-cc73cf61).
//
// It exists for a repository whose markers predate the two hashes: running it
// asserts that the code is current, so the next grow skips instead of
// re-selecting every record. Only a missing value is written. A value already
// present is never replaced, because replacing it would hide a real change -
// that is what grow and --check are for. The record id is not touched either,
// so a marker naming an alias keeps it until its region is next written.
//
// Only the opening marker line of a stamped region changes. Every file is
// written once, after all of them were stamped in memory, and a dry run writes
// nothing.
func Stamp(requests []StampRequest, opts Options) ([]Stamped, error) {
	own := newOwnership(opts.Owners)
	buffers := map[string]string{}
	var (
		out      []Stamped
		problems []error
	)

	for _, req := range requests {
		for _, target := range req.Targets {
			content, ok := buffers[target.Path]
			if !ok {
				content, ok = readExisting(opts.Output, opts.Unwritten, target.Path)
				if !ok {
					continue
				}
			}
			regions, err := FindRegions(target.Package.Comment(), content)
			if err != nil {
				problems = append(problems, fmt.Errorf("file %s: %w", target.Path, err))
				continue
			}

			changed := false
			// Backwards, so a rewritten line does not move the offsets of the
			// regions still to come.
			for i := len(regions) - 1; i >= 0; i-- {
				region := regions[i]
				m := region.Marker
				if m.Record == "" || !own.owns(req.RecordID, m.Record) {
					continue
				}
				stamp := Stamped{Path: target.Path, Action: m.Action, Variant: m.Variant, Record: m.Record}
				if m.Gen == "" && req.Gen != "" {
					m.Gen, stamp.Gen = req.Gen, true
				}
				if action, ok := target.Package.Actions[m.Action]; ok && m.Pkg == "" && action.PkgHash != "" {
					m.Pkg, stamp.Pkg = action.PkgHash, true
				}
				if !stamp.Gen && !stamp.Pkg {
					continue
				}

				open, err := m.Open(target.Package.Comment())
				if err != nil {
					problems = append(problems, fmt.Errorf("file %s: %w", target.Path, err))
					continue
				}
				lineEnd := strings.IndexByte(content[region.Start:], '\n')
				if lineEnd < 0 {
					lineEnd = len(content) - region.Start
				}
				content = content[:region.Start] + indentAt(content, region.Start) + open + content[region.Start+lineEnd:]
				out = append(out, stamp)
				changed = true
			}
			if changed {
				buffers[target.Path] = content
			}
		}
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	if opts.DryRun {
		return out, nil
	}

	paths := make([]string, 0, len(buffers))
	for path := range buffers {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		full := filepath.Join(opts.Output, filepath.FromSlash(path))
		if err := os.WriteFile(full, []byte(buffers[path]), 0o644); err != nil {
			problems = append(problems, fmt.Errorf("write %s: %w", path, err))
		}
	}
	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	return out, nil
}
