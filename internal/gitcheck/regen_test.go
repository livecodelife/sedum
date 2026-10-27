package gitcheck

import (
	"strings"
	"testing"
)

// The regeneration proof says whether a commit only regenerates governed
// regions.

// memReader is a TreeReader over an in-memory set of files.
func memReader(files map[string]string) TreeReader {
	return TreeReader{
		read: func(p string) ([]byte, bool, error) {
			content, ok := files[p]
			return []byte(content), ok, nil
		},
		regular: func(p string) bool {
			_, ok := files[p]
			return ok
		},
	}
}

const (
	regenFile = "app/models/user.rb"
	regenRec  = "prov-2026-aaaaaaaa"
)

func rubyRegion(record, gen, pkg, body string) string {
	return "# header\n" + open("#", "addField", "owned", record, gen, pkg) + "\n" + body + closeLine("#", "addField") + "\n# footer\n"
}

func currentClasses(records ...string) map[string]string {
	m := map[string]string{}
	for _, r := range records {
		m[r] = "current"
	}
	return m
}

func governed(paths ...string) Classified { return Classified{Governed: paths} }

func TestRegenerationAcceptsARegeneratedRegion(t *testing.T) {
	base := memReader(map[string]string{regenFile: rubyRegion(regenRec, "g1", "p1", "old\n")})
	commit := memReader(map[string]string{regenFile: rubyRegion(regenRec, "g1", "p2", "new\nbody\n")})

	if !Regeneration(governed(regenFile), base, commit, currentClasses(regenRec)) {
		t.Error("a region whose body and pkg changed was refused")
	}
}

func TestRegenerationRefusals(t *testing.T) {
	baseText := rubyRegion(regenRec, "g1", "p1", "old\n")
	cases := []struct {
		name    string
		c       Classified
		base    map[string]string
		commit  map[string]string
		classes map[string]string
	}{
		{"a record file changed", Classified{Governed: []string{regenFile}, RecordFiles: []string{"provenance/x.yml"}}, map[string]string{regenFile: baseText}, map[string]string{regenFile: baseText}, currentClasses(regenRec)},
		{"a state file changed", Classified{Governed: []string{regenFile}, StateFiles: []string{".provenance/versions.jsonl"}}, map[string]string{regenFile: baseText}, map[string]string{regenFile: baseText}, currentClasses(regenRec)},
		{"classes could not be computed", governed(regenFile), map[string]string{regenFile: baseText}, map[string]string{regenFile: baseText}, nil},
		{"the record is not current", governed(regenFile), map[string]string{regenFile: baseText}, map[string]string{regenFile: baseText}, map[string]string{regenRec: "stale"}},
		{"the record has no class", governed(regenFile), map[string]string{regenFile: baseText}, map[string]string{regenFile: baseText}, map[string]string{}},
		{"no governed file changed", Classified{Excluded: []string{"README.md"}}, nil, nil, currentClasses(regenRec)},
		{"a governed file was deleted", governed(regenFile), map[string]string{regenFile: baseText}, map[string]string{}, currentClasses(regenRec)},
		{"a governed file is new", governed(regenFile), map[string]string{}, map[string]string{regenFile: baseText}, currentClasses(regenRec)},
		{"the base has no region", governed("a.rb"), map[string]string{"a.rb": "plain\n"}, map[string]string{"a.rb": "plain\n"}, currentClasses(regenRec)},
		{"gen changed", governed(regenFile), map[string]string{regenFile: baseText}, map[string]string{regenFile: rubyRegion(regenRec, "g9", "p1", "old\n")}, currentClasses(regenRec)},
		{"record changed", governed(regenFile), map[string]string{regenFile: baseText}, map[string]string{regenFile: rubyRegion("prov-2026-bbbbbbbb", "g1", "p1", "old\n")}, currentClasses(regenRec, "prov-2026-bbbbbbbb")},
		{"text outside regions changed", governed(regenFile), map[string]string{regenFile: baseText}, map[string]string{regenFile: strings.Replace(baseText, "# footer", "# footer!", 1)}, currentClasses(regenRec)},
		{"a region was added", governed(regenFile), map[string]string{regenFile: baseText}, map[string]string{regenFile: baseText + rubyRegion(regenRec, "g1", "p1", "x\n")}, currentClasses(regenRec)},
		{"a changed recording is a governed file with no regions", Classified{Governed: []string{regenFile}, Recordings: []string{".sedum/recordings/r.json"}}, map[string]string{regenFile: baseText, ".sedum/recordings/r.json": "{}"}, map[string]string{regenFile: baseText, ".sedum/recordings/r.json": "{\"a\":1}"}, currentClasses(regenRec)},
	}
	for _, c := range cases {
		if Regeneration(c.c, memReader(c.base), memReader(c.commit), c.classes) {
			t.Errorf("%s: Regeneration accepted the commit", c.name)
		}
	}
}

func TestRegenerationRequiresEveryRecordNamedByAChangedRegionToBeCurrent(t *testing.T) {
	two := rubyRegion(regenRec, "g1", "p1", "a\n") + rubyRegion("prov-2026-bbbbbbbb", "g2", "p2", "b\n")
	base := memReader(map[string]string{regenFile: two})

	if Regeneration(governed(regenFile), base, base, map[string]string{regenRec: "current", "prov-2026-bbbbbbbb": "edited"}) {
		t.Error("a second record that is not current was ignored")
	}
	if !Regeneration(governed(regenFile), base, base, currentClasses(regenRec, "prov-2026-bbbbbbbb")) {
		t.Error("two current records were refused")
	}
}

func TestRegenerationIgnoresExcludedPathsAndExcludedRecordings(t *testing.T) {
	baseText := rubyRegion(regenRec, "g1", "p1", "old\n")
	base := memReader(map[string]string{regenFile: baseText})
	commit := memReader(map[string]string{regenFile: rubyRegion(regenRec, "g1", "p2", "new\n")})
	c := Classified{
		Governed:           []string{regenFile},
		Excluded:           []string{"README.md"},
		Recordings:         []string{".sedum/recordings/r.json"},
		ExcludedRecordings: []string{".sedum/recordings/r.json"},
	}

	if !Regeneration(c, base, commit, currentClasses(regenRec)) {
		t.Error("an excluded path or an excluded recording blocked the regeneration")
	}
}
