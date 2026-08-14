package record

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Phase 1 turns a directory of provenance records into the authorized path set
// every later phase is bounded by. The cases below are organized around the two
// things that can go wrong with that: reading a record wrongly, and deciding
// wrongly what it authorizes.
//
// Records are written here as YAML strings rather than kept as fixture files,
// because most cases differ from a valid record by one line and reading that
// line next to the assertion is the point.

// validRecord is a record carrying every field Sedum reads plus a
// representative sample of the fields it must ignore.
const validRecord = `id: prov-2026-aaaaaaaa
title: Add the users controller
status: implemented
created_at: "2026-08-07"
author: someone@example.com
intent: |
  Add a users controller with an index action.
constraints:
  - Read-only endpoints only.
affected_scope:
  - app/controllers/users_controller.rb
forbidden_scope: []
type: blueprint
implements: prov-2026-bbbbbbbb
supersedes: ""
superseded_by: ""
related: []
sealed_at_sha: ""
associated_specs:
  - path: spec/controllers/users_controller_spec.rb
    type: rspec
associated_traces: []
monitors: []
tags:
  - example
`

// writeRecords materializes a records directory and returns its path.
func writeRecords(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, content := range files {
		full := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", name, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return root
}

func loadOne(t *testing.T, content string) *Record {
	t.Helper()
	set, _, err := Load(writeRecords(t, map[string]string{"record.yml": content}), Options{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(set.Records) != 1 {
		t.Fatalf("Load returned %d records, want 1", len(set.Records))
	}
	return set.Records[0]
}

func wantErr(t *testing.T, err error, fragments ...string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error naming %v, got nil", fragments)
	}
	for _, want := range fragments {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error = %q, does not mention %q", err, want)
		}
	}
}

// A provenance record is written by a governance tool, not by Sedum, and it
// carries a dozen fields Sedum has no business reading. Ingestion takes the
// four the pipeline is defined in terms of and ignores the rest - including
// fields that do not exist yet, since a record schema that grows must not
// break generation (prov-2026-dad4d877).
func TestLoadReadsFourFieldsAndIgnoresTheRest(t *testing.T) {
	r := loadOne(t, validRecord)

	if r.ID != "prov-2026-aaaaaaaa" {
		t.Errorf("ID = %q", r.ID)
	}
	if !strings.Contains(r.Intent, "users controller") {
		t.Errorf("Intent = %q", r.Intent)
	}
	if len(r.Constraints) != 1 || r.Constraints[0] != "Read-only endpoints only." {
		t.Errorf("Constraints = %v", r.Constraints)
	}
	if len(r.Paths) != 1 || r.Paths[0] != "app/controllers/users_controller.rb" {
		t.Errorf("Paths = %v", r.Paths)
	}
}

// The inverse of the rule above: a key nobody has invented yet must not fail
// the load. Strict decoding is right for a generator package, whose keys Sedum
// defines, and wrong here, whose keys it does not.
func TestLoadAcceptsUnknownFields(t *testing.T) {
	r := loadOne(t, validRecord+"\nsome_field_from_a_later_version: {a: 1, b: [2, 3]}\n")

	if r.ID != "prov-2026-aaaaaaaa" {
		t.Errorf("a record carrying an unrecognized key failed to load: ID = %q", r.ID)
	}
}

// A path names a file to create. A pattern authorizes a set and names nothing,
// so it cannot be created (prov-2026-e8671c88).
func TestAffectedScopeSplitsPathsFromPatterns(t *testing.T) {
	r := loadOne(t, `id: prov-2026-aaaaaaaa
affected_scope:
  - app/models/user.rb
  - app/controllers/**
  - config/initializers/*.rb
  - lib/tasks/
  - app/models/user.rb
`)

	wantPaths := []string{"app/models/user.rb"}
	if !equal(r.Paths, wantPaths) {
		t.Errorf("Paths = %v, want %v (deduplicated, patterns excluded)", r.Paths, wantPaths)
	}
	wantPatterns := []string{"app/controllers/**", "config/initializers/*.rb", "lib/tasks/"}
	if !equal(r.Patterns, wantPatterns) {
		t.Errorf("Patterns = %v, want %v (a trailing slash names a subtree, not a file)", r.Patterns, wantPatterns)
	}
}

// A record that authorizes only subtrees creates nothing. That is legal - the
// record may exist to bound a hand-written change - but it is worth saying so,
// because the alternative reading is that generation silently did nothing.
func TestRecordNamingNoPathWarns(t *testing.T) {
	set, warnings, err := Load(writeRecords(t, map[string]string{"r.yml": `id: prov-2026-aaaaaaaa
affected_scope:
  - internal/**
`}), Options{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(set.Records) != 1 {
		t.Fatalf("a record authorizing only patterns was dropped rather than kept")
	}
	if !containsFragment(warnings, "prov-2026-aaaaaaaa") || !containsFragment(warnings, "no path") {
		t.Errorf("warnings = %v, do not say the record names no path Sedum can create", warnings)
	}
}

// forbidden_scope means Sedum does not touch what it was not authorized to
// touch. A record that names a path and forbids it in the same breath has no
// reading that is safe to act on.
func TestPathForbiddenByItsOwnRecordIsAnError(t *testing.T) {
	_, _, err := Load(writeRecords(t, map[string]string{"r.yml": `id: prov-2026-aaaaaaaa
affected_scope:
  - app/models/user.rb
forbidden_scope:
  - app/models/**
`}), Options{})

	wantErr(t, err, "prov-2026-aaaaaaaa", "app/models/user.rb", "app/models/**")
}

// One file, two records means two model calls injecting into the same file,
// which the PRD puts out of scope. It has to fail here, where both record IDs
// are still in hand to name.
func TestTwoRecordsNamingOnePathIsAnError(t *testing.T) {
	_, _, err := Load(writeRecords(t, map[string]string{
		"a.yml": "id: prov-2026-aaaaaaaa\naffected_scope: [app/models/user.rb]\n",
		"b.yml": "id: prov-2026-bbbbbbbb\naffected_scope: [app/models/user.rb]\n",
	}), Options{})

	wantErr(t, err, "app/models/user.rb", "prov-2026-aaaaaaaa", "prov-2026-bbbbbbbb")
}

// Every authorized path is created under the output directory, so a path that
// resolves outside it is not an authorization Sedum can honor.
func TestPathEscapingTheOutputDirectoryIsAnError(t *testing.T) {
	for _, path := range []string{"/etc/passwd", "../outside.rb", "app/../../outside.rb"} {
		t.Run(path, func(t *testing.T) {
			_, _, err := Load(writeRecords(t, map[string]string{
				"r.yml": "id: prov-2026-aaaaaaaa\naffected_scope: [\"" + path + "\"]\n",
			}), Options{})
			wantErr(t, err, path)
		})
	}
}

func TestRecordWithNoIDIsAnError(t *testing.T) {
	_, _, err := Load(writeRecords(t, map[string]string{
		"nameless.yml": "affected_scope: [app/models/user.rb]\n",
	}), Options{})

	wantErr(t, err, "nameless.yml", "id")
}

func TestUnparseableRecordNamesTheFile(t *testing.T) {
	_, _, err := Load(writeRecords(t, map[string]string{
		"broken.yml": "id: prov-2026-aaaaaaaa\naffected_scope: [unclosed\n",
	}), Options{})

	wantErr(t, err, "broken.yml")
}

// Reporting one problem per run makes fixing a records directory iterative,
// which is the same reason package loading reports every finding at once.
func TestEveryProblemIsReportedAtOnce(t *testing.T) {
	_, _, err := Load(writeRecords(t, map[string]string{
		"a.yml": "affected_scope: [app/models/user.rb]\n",
		"b.yml": "id: prov-2026-bbbbbbbb\naffected_scope: [\"/etc/passwd\"]\n",
	}), Options{})

	wantErr(t, err, "a.yml", "/etc/passwd")
}

// The records directory is a governance directory, not a Sedum directory: it
// holds READMEs, indexes, and whatever else a team keeps beside its records.
func TestNonYAMLFilesAreIgnored(t *testing.T) {
	set, _, err := Load(writeRecords(t, map[string]string{
		"r.yml":     validRecord,
		"README.md": "# Provenance\n",
		"notes.txt": "not a record\n",
	}), Options{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(set.Records) != 1 {
		t.Errorf("Load returned %d records, want 1; files that are not YAML must be ignored", len(set.Records))
	}
}

func TestOnlyFiltersToNamedRecords(t *testing.T) {
	dir := writeRecords(t, map[string]string{
		"a.yml": "id: prov-2026-aaaaaaaa\naffected_scope: [a.rb]\n",
		"b.yml": "id: prov-2026-bbbbbbbb\naffected_scope: [b.rb]\n",
	})

	set, _, err := Load(dir, Options{Only: []string{"prov-2026-bbbbbbbb"}})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(set.Records) != 1 || set.Records[0].ID != "prov-2026-bbbbbbbb" {
		t.Fatalf("--only selected %v", ids(set))
	}
}

// A mistyped --only would otherwise generate nothing and report success, which
// is indistinguishable from a run that had nothing to do.
func TestOnlyNamingAnUnknownRecordIsAnError(t *testing.T) {
	dir := writeRecords(t, map[string]string{"a.yml": "id: prov-2026-aaaaaaaa\naffected_scope: [a.rb]\n"})

	_, _, err := Load(dir, Options{Only: []string{"prov-2026-cccccccc"}})
	wantErr(t, err, "prov-2026-cccccccc")
}

// Records are sorted so that a run's output does not depend on the order the
// records directory happened to be read in.
func TestRecordsAreSortedByID(t *testing.T) {
	set, _, err := Load(writeRecords(t, map[string]string{
		"z.yml": "id: prov-2026-cccccccc\naffected_scope: [c.rb]\n",
		"a.yml": "id: prov-2026-aaaaaaaa\naffected_scope: [a.rb]\n",
		"m.yml": "id: prov-2026-bbbbbbbb\naffected_scope: [b.rb]\n",
	}), Options{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	want := []string{"prov-2026-aaaaaaaa", "prov-2026-bbbbbbbb", "prov-2026-cccccccc"}
	if !equal(ids(set), want) {
		t.Errorf("record order = %v, want %v", ids(set), want)
	}
}

// A missing records directory is an I/O failure that stops ingestion from
// happening at all, not a record-level problem to collect alongside others.
func TestMissingRecordsDirectoryIsAnError(t *testing.T) {
	_, _, err := Load(filepath.Join(t.TempDir(), "absent"), Options{})
	wantErr(t, err, "absent")
}

func TestForbids(t *testing.T) {
	r := loadOne(t, `id: prov-2026-aaaaaaaa
affected_scope:
  - app/models/user.rb
forbidden_scope:
  - db/schema.rb
  - vendor/**
  - config/*.yml
`)

	tests := []struct {
		path   string
		forbid bool
	}{
		{"db/schema.rb", true},
		{"vendor/gems/thing.rb", true},
		{"vendor/thing.rb", true},
		{"vendor", true},
		{"config/database.yml", true},
		{"config/deploy/production.yml", false},
		{"app/models/user.rb", false},
		{"db/migrate/001.rb", false},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			entry, got := r.Forbids(tc.path)
			if got != tc.forbid {
				t.Errorf("Forbids(%q) = %v (%q), want %v", tc.path, got, entry, tc.forbid)
			}
		})
	}
}

func ids(s *Set) []string {
	out := make([]string, 0, len(s.Records))
	for _, r := range s.Records {
		out = append(out, r.ID)
	}
	return out
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func containsFragment(list []string, want string) bool {
	for _, s := range list {
		if strings.Contains(s, want) {
			return true
		}
	}
	return false
}
