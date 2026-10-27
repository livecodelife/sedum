package gitcheck

// Migration reports whether a commit only migrates marker opening lines, and
// returns the records that own the changed regions. The commit must change
// governed files and recordings only: no record file, state file or excluded
// non-recording; each recording a regular file at the commit; each governed file
// with the same regions, bodies, closing lines, prefixes, pkg values and tier,
// record and gen values as at the base (other opening-line JSON keys may
// differ); and every record named by a region in the same class at the base tree
// (before) and the commit tree (after). The tag must be [id@N] with N the record's
// current version in versions, and its record must own a changed region.
func Migration(c Classified, base, commit TreeReader, before, after map[string]string, tag Tag, versions map[string]int) (bool, map[string]bool) {
	if before == nil || after == nil || c.Blocked() || len(c.Excluded) > 0 {
		return false, nil
	}
	for _, p := range c.Recordings {
		if !commit.Regular(p) {
			return false, nil
		}
	}
	files, ok := GovernedChanges(c.Governed, base, commit)
	if !ok {
		return false, nil
	}
	owners := map[string]bool{}
	for _, f := range files {
		if !SameShape(f.Base, f.Commit, true) {
			return false, nil
		}
		_, regions, ok := ReadRegions(f.Commit)
		if !ok {
			return false, nil
		}
		for _, r := range regions {
			owners[r.Record] = true
		}
	}
	for record := range owners {
		was, ok1 := before[record]
		now, ok2 := after[record]
		if !ok1 || !ok2 || was != now {
			return false, nil
		}
	}
	if !tag.HasVersion {
		return false, nil
	}
	if current, ok := versions[tag.ID]; !ok || current != tag.Version {
		return false, nil
	}
	if !owners[tag.ID] {
		return false, nil
	}
	return true, owners
}
