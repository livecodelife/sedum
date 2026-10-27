package gitcheck

// Regeneration reports whether a commit only regenerates governed regions: no
// record file or provenance state file changed; every governed file, which is
// the Governed paths together with the Recordings that are not excluded, has the
// same regions at the commit as at the base (same tier, record and gen; the
// body, closing line, prefix, pkg and other JSON keys may change); and every
// record named by a region is `current` in classes, the class map of the
// commit's tree, which is nil when classes could not be computed.
func Regeneration(c Classified, base, commit TreeReader, classes map[string]string) bool {
	if c.Blocked() || classes == nil {
		return false
	}
	excluded := map[string]bool{}
	for _, p := range c.ExcludedRecordings {
		excluded[p] = true
	}
	set := append([]string{}, c.Governed...)
	for _, p := range c.Recordings {
		if !excluded[p] {
			set = append(set, p)
		}
	}
	files, ok := GovernedChanges(set, base, commit)
	if !ok {
		return false
	}
	for _, f := range files {
		if !SameShape(f.Base, f.Commit, false) {
			return false
		}
		_, regions, ok := ReadRegions(f.Commit)
		if !ok {
			return false
		}
		for _, r := range regions {
			if classes[r.Record] != "current" {
				return false
			}
		}
	}
	return true
}
