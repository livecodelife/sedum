package gitcheck

import (
	"strings"
	"testing"
)

// The marker migration proof says whether a commit only migrates marker opening
// lines and returns the records that own the changed regions.

const migRec = "prov-2026-aaaaaaaa"

const v1Line = `// sedum:addField {"tier":"owned","record":"prov-2026-aaaaaaaa","gen":"g1","pkg":"p1","kwargs":{"model":"user"}}`
const v2Line = `// sedum:addField {"v":2,"tier":"owned","record":"prov-2026-aaaaaaaa","gen":"g1","pkg":"p1","id":{"file":"a.go"},"kwargs":"sha256:1","body":"sha256:2"}`

func migFile(open string) string {
	return "head\n" + open + "\nbody\n// /sedum:addField\ntail\n"
}

type migCase struct {
	c       Classified
	base    map[string]string
	commit  map[string]string
	before  map[string]string
	after   map[string]string
	tag     Tag
	version map[string]int
}

func goodMigration() migCase {
	return migCase{
		c:       Classified{Governed: []string{"a.go"}, Recordings: []string{".sedum/recordings/r.json"}},
		base:    map[string]string{"a.go": migFile(v1Line)},
		commit:  map[string]string{"a.go": migFile(v2Line), ".sedum/recordings/r.json": "{}"},
		before:  map[string]string{migRec: "current"},
		after:   map[string]string{migRec: "current"},
		tag:     Tag{ID: migRec, Version: 3, HasVersion: true},
		version: map[string]int{migRec: 3},
	}
}

func (m migCase) run() (bool, map[string]bool) {
	return Migration(m.c, memReader(m.base), memReader(m.commit), m.before, m.after, m.tag, m.version)
}

func TestMigrationAcceptsAV1ToV2MigrationAndReturnsTheOwners(t *testing.T) {
	ok, owners := goodMigration().run()

	if !ok || len(owners) != 1 || !owners[migRec] {
		t.Errorf("ok = %v owners = %v", ok, owners)
	}
}

func TestMigrationAcceptsTheOtherDirectionAndAnAlreadyStaleRecord(t *testing.T) {
	m := goodMigration()
	m.base["a.go"], m.commit["a.go"] = migFile(v2Line), migFile(v1Line)
	m.before[migRec], m.after[migRec] = "stale", "stale"

	if ok, _ := m.run(); !ok {
		t.Error("a v2 to v1 migration of a record that stays stale was refused")
	}
}

func TestMigrationRefusals(t *testing.T) {
	cases := map[string]func(*migCase){
		"a hand edit inside a region":  func(m *migCase) { m.commit["a.go"] = strings.Replace(m.commit["a.go"], "body\n", "edited\n", 1) },
		"text outside regions changed": func(m *migCase) { m.commit["a.go"] = strings.Replace(m.commit["a.go"], "tail", "tail!", 1) },
		"the closing line changed": func(m *migCase) {
			m.commit["a.go"] = strings.Replace(m.commit["a.go"], "// /sedum:addField\n", "// /sedum:addField \n", 1)
		},
		"the opening prefix changed": func(m *migCase) {
			m.commit["a.go"] = strings.Replace(m.commit["a.go"], "// sedum:addField", "  // sedum:addField", 1)
		},
		"pkg changed":                              func(m *migCase) { m.commit["a.go"] = strings.Replace(m.commit["a.go"], `"pkg":"p1"`, `"pkg":"p2"`, 1) },
		"gen changed":                              func(m *migCase) { m.commit["a.go"] = strings.Replace(m.commit["a.go"], `"gen":"g1"`, `"gen":"g2"`, 1) },
		"a record file changed":                    func(m *migCase) { m.c.RecordFiles = []string{"provenance/x.yml"} },
		"a state file changed":                     func(m *migCase) { m.c.StateFiles = []string{".provenance/versions.jsonl"} },
		"an excluded non-recording changed":        func(m *migCase) { m.c.Excluded = []string{"README.md"} },
		"no governed file changed":                 func(m *migCase) { m.c.Governed = nil },
		"a governed file is new":                   func(m *migCase) { delete(m.base, "a.go") },
		"a recording was deleted":                  func(m *migCase) { delete(m.commit, ".sedum/recordings/r.json") },
		"the base classes could not be computed":   func(m *migCase) { m.before = nil },
		"the commit classes could not be computed": func(m *migCase) { m.after = nil },
		"a class changed":                          func(m *migCase) { m.after[migRec] = "stale" },
		"a record has no class at the base":        func(m *migCase) { delete(m.before, migRec) },
		"a record has no class at the commit":      func(m *migCase) { delete(m.after, migRec) },
		"the tag has no version":                   func(m *migCase) { m.tag = Tag{ID: migRec} },
		"the record has no current version":        func(m *migCase) { m.version = map[string]int{} },
		"the tag names another version":            func(m *migCase) { m.tag.Version = 2 },
		"the tag's record owns no changed region":  func(m *migCase) { m.tag.ID = "prov-2026-bbbbbbbb"; m.version["prov-2026-bbbbbbbb"] = 3 },
	}
	for name, mutate := range cases {
		m := goodMigration()
		mutate(&m)
		if ok, _ := m.run(); ok {
			t.Errorf("%s: Migration accepted the commit", name)
		}
	}
}

func TestMigrationIgnoresAnExcludedRecordingButStillNeedsItToBeRegular(t *testing.T) {
	m := goodMigration()
	m.c.ExcludedRecordings = []string{".sedum/recordings/r.json"}
	if ok, _ := m.run(); !ok {
		t.Error("an excluded recording blocked the migration")
	}
}

func TestMigrationOwnersCoverEveryRecordNamedByAChangedRegion(t *testing.T) {
	two := migFile(v1Line) + "// sedum:addField " + `{"tier":"owned","record":"prov-2026-bbbbbbbb","gen":"g","pkg":"p"}` + "\nx\n// /sedum:addField\n"
	two2 := strings.Replace(two, v1Line, v2Line, 1)
	m := goodMigration()
	m.base["a.go"], m.commit["a.go"] = two, two2
	m.before["prov-2026-bbbbbbbb"], m.after["prov-2026-bbbbbbbb"] = "current", "current"

	ok, owners := m.run()

	if !ok || !owners[migRec] || !owners["prov-2026-bbbbbbbb"] || len(owners) != 2 {
		t.Errorf("ok = %v owners = %v", ok, owners)
	}
}
