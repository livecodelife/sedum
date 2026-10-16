// Package record ingests provenance records and collects the authorized path
// set every later phase is bounded by.
//
// A provenance record declares intent, constraints, and the files a change is
// authorized to touch. Sedum reads five fields - id, intent, constraints,
// affected_scope, and forbidden_scope - and nothing else. It does not interpret
// a record's status or lifecycle, and it never writes to one: the record schema
// belongs to the governance tool that produces it, so decoding is lenient and a
// key Sedum does not model is ignored rather than rejected (prov-2026-dad4d877).
//
// An affected_scope entry is either a literal path, which names a file to
// create, or a pattern, which authorizes matches without naming any. Only the
// literal paths enter the authorized path set; Sedum does not invent a member
// of the set a pattern describes (prov-2026-e8671c88).
package record

import (
	"errors"
	"fmt"
	"github.com/livecodelife/sedum/internal/genhash"
	"github.com/livecodelife/sedum/internal/pathpat"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Options controls which records a load considers.
type Options struct {
	// Only names the record IDs to ingest. Empty ingests every record in the
	// directory.
	Only []string
}

// Record is the part of a provenance record Sedum reads.
type Record struct {
	ID string
	// Source is the file the record was read from, for diagnostics.
	Source      string
	Intent      string
	Constraints []string

	// Paths are the literal paths affected_scope names: the files this record
	// authorizes Sedum to create. Slash-normalized, deduplicated, and sorted.
	Paths []string
	// Patterns are the affected_scope entries that authorize matches without
	// naming a file. They bound what a run may touch and create nothing.
	Patterns []string
	// Forbidden are the forbidden_scope entries, literal or pattern alike.
	Forbidden []string

	// Gen is the record's gen_hash, which every region it writes carries as
	// "gen" (prov-2026-66bea9f3). It is taken over the four fields exactly as
	// the file states them - before trimming, deduplication, sorting, or the
	// split into paths and patterns - because the hash is shared with every
	// other tool that reads records, and none of them see Sedum's view.
	Gen string

	// CreatedAt and Aliases are what region ownership reads
	// (prov-2026-375f5d4a): a record owns the regions whose marker names its
	// id or one of its aliases - ids folded into it by migration - and owners
	// at a shared anchor are ordered by creation date, then id.
	CreatedAt string
	Aliases   []string

	// Deprecated reports a retired record (prov-2026-cc73cf61). It is the one
	// lifecycle fact Sedum reads, and only the value "deprecated" is: such a
	// record takes no new versions, is never grown, and owns nothing, so its
	// regions are orphans.
	Deprecated bool
}

// Forbids reports whether the record's forbidden_scope covers path, and returns
// the entry that covers it.
func (r *Record) Forbids(target string) (string, bool) {
	target = normalize(target)
	for _, entry := range r.Forbidden {
		if matchScope(entry, target) {
			return entry, true
		}
	}
	return "", false
}

// Set is the ingested records directory.
type Set struct {
	// Records are the ingested records, sorted by ID.
	Records []*Record

	// Known is every record read from the directory, sorted by ID, including
	// those --only left out and deprecated ones, which Records never holds. A run grows only Records, but the regions it
	// meets may belong to any of them, and ownership and anchor order need
	// their aliases and creation dates (prov-2026-375f5d4a).
	Known []*Record
}

// Paths returns every literal path the set authorizes, sorted.
//
// A path may appear once per record naming it. Ingestion no longer rejects a
// path named by two records - CheckDuplicatePaths does, at Phase 4's entry - so
// a caller that needs each path once deduplicates.
func (s *Set) Paths() []string {
	var out []string
	for _, r := range s.Records {
		out = append(out, r.Paths...)
	}
	sort.Strings(out)
	return out
}

// Lookup returns the record with the given ID.
func (s *Set) Lookup(id string) (*Record, bool) {
	for _, r := range s.Records {
		if r.ID == id {
			return r, true
		}
	}
	return nil, false
}

// Load reads every provenance record under dir.
//
// The returned error is a joined list of everything wrong with the directory's
// contents, because a diagnostic that stops at the first bad record makes
// fixing a records directory iterative. A separate error is returned for an I/O
// problem that stops ingestion from happening at all.
//
// Warnings are conditions worth reporting that do not stop a run.
func Load(dir string, opts Options) (*Set, []string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("read records directory %s: %w", dir, err)
	}

	var (
		records  []*Record
		warnings []string
		problems []error
	)

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && isRecordFile(e.Name()) {
			names = append(names, e.Name())
		}
	}
	// Sorted so that the problems reported do not depend on directory
	// iteration order.
	sort.Strings(names)

	for _, name := range names {
		rec, err := readRecord(filepath.Join(dir, name), name)
		if err != nil {
			problems = append(problems, err)
			continue
		}
		records = append(records, rec)
	}

	known := append([]*Record(nil), records...)
	sort.Slice(known, func(i, j int) bool { return known[i].ID < known[j].ID })

	records, err = filterOnly(records, opts.Only)
	records = slices.DeleteFunc(records, func(r *Record) bool { return r.Deprecated })
	if err != nil {
		problems = append(problems, err)
	}

	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })

	problems = append(problems, checkForbidden(records)...)
	warnings = append(warnings, warnEmpty(records)...)

	if len(problems) > 0 {
		return nil, warnings, errors.Join(problems...)
	}
	return &Set{Records: records, Known: known}, warnings, nil
}

// isRecordFile reports whether a directory entry is a record. The records
// directory is a governance directory rather than a Sedum one: it holds
// READMEs, indexes, and whatever else a team keeps beside its records, and
// those are not ours to complain about.
func isRecordFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".yml" || ext == ".yaml"
}

// document is the shape Sedum decodes a record into. Decoding is deliberately
// not strict: every other key a record carries belongs to the governance tool
// that wrote it (prov-2026-dad4d877).
type document struct {
	ID             string   `yaml:"id"`
	Intent         string   `yaml:"intent"`
	Constraints    []string `yaml:"constraints"`
	AffectedScope  []string `yaml:"affected_scope"`
	ForbiddenScope []string `yaml:"forbidden_scope"`
	CreatedAt      string   `yaml:"created_at"`
	Aliases        []string `yaml:"aliases"`
	Status         string   `yaml:"status"`
}

func readRecord(full, name string) (*Record, error) {
	data, err := os.ReadFile(full)
	if err != nil {
		return nil, fmt.Errorf("read record %s: %w", name, err)
	}

	var doc document
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("record %s could not be read: %w", name, err)
	}
	if doc.ID == "" {
		return nil, fmt.Errorf("record %s declares no id; the id is what --only names and what every diagnostic identifies a record by", name)
	}

	rec := &Record{
		ID:          doc.ID,
		Source:      name,
		Intent:      doc.Intent,
		Constraints: doc.Constraints,
		CreatedAt:   doc.CreatedAt,
		Aliases:     doc.Aliases,
		Deprecated:  strings.TrimSpace(doc.Status) == "deprecated",
		Gen: genhash.Of(genhash.Input{
			Intent:         doc.Intent,
			Constraints:    doc.Constraints,
			AffectedScope:  doc.AffectedScope,
			ForbiddenScope: doc.ForbiddenScope,
		}),
	}

	var problems []error
	seen := map[string]bool{}
	for _, entry := range doc.AffectedScope {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if isPattern(entry) {
			if err := checkPattern(entry); err != nil {
				problems = append(problems, fmt.Errorf("record %s: affected_scope entry %s: %w", name, entry, err))
				continue
			}
			rec.Patterns = append(rec.Patterns, entry)
			continue
		}
		clean := normalize(entry)
		if escapes(clean) {
			problems = append(problems, fmt.Errorf(
				"record %s authorizes %s, which resolves outside the output directory; every authorized path is created under it",
				name, entry))
			continue
		}
		if !seen[clean] {
			seen[clean] = true
			rec.Paths = append(rec.Paths, clean)
		}
	}
	for _, entry := range doc.ForbiddenScope {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if err := checkPattern(entry); err != nil {
			problems = append(problems, fmt.Errorf("record %s: forbidden_scope entry %s: %w", name, entry, err))
			continue
		}
		rec.Forbidden = append(rec.Forbidden, entry)
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}

	sort.Strings(rec.Paths)
	sort.Strings(rec.Patterns)
	return rec, nil
}

func filterOnly(records []*Record, only []string) ([]*Record, error) {
	if len(only) == 0 {
		return records, nil
	}

	kept := make([]*Record, 0, len(only))
	var retired []string
	for _, r := range records {
		if slices.Contains(only, r.ID) {
			if r.Deprecated {
				retired = append(retired, r.ID)
				continue
			}
			kept = append(kept, r)
		}
	}
	if len(retired) > 0 {
		return kept, fmt.Errorf("--only names %s, which is deprecated: a deprecated record is never grown, and its regions are orphans for --prune-orphans",
			strings.Join(retired, ", "))
	}

	// A mistyped --only would otherwise generate nothing and report success,
	// which is indistinguishable from a run that had nothing to do.
	var missing []string
	for _, id := range only {
		if !slices.ContainsFunc(kept, func(r *Record) bool { return r.ID == id }) {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		return kept, fmt.Errorf("--only names %s, which no record in the directory declares",
			strings.Join(missing, ", "))
	}
	return kept, nil
}

// checkForbidden rejects a record that names a path and forbids it in the same
// breath. Neither half can be honored silently: acting on the authorization
// would make forbidden_scope advisory, and skipping the path would make
// affected_scope advisory.
func checkForbidden(records []*Record) []error {
	var problems []error
	for _, r := range records {
		for _, p := range r.Paths {
			if entry, forbidden := r.Forbids(p); forbidden {
				problems = append(problems, fmt.Errorf(
					"record %s authorizes %s and forbids it through forbidden_scope entry %s; the record has no reading Sedum can act on",
					r.ID, p, entry))
			}
		}
	}
	return problems
}

// warnEmpty reports a record that authorizes only patterns. Generating nothing
// is legal - the record may exist to bound a hand-written change - but it must
// be distinguishable from failing to generate.
func warnEmpty(records []*Record) []string {
	var out []string
	for _, r := range records {
		if len(r.Paths) == 0 {
			out = append(out, fmt.Sprintf(
				"record %s names no path Sedum can create; its affected_scope authorizes matches but names no file", r.ID))
		}
	}
	return out
}

// normalize puts a path in the slash-separated, cleaned form every phase
// compares in, so that a record authored on one platform reads identically on
// another.
func normalize(p string) string {
	return path.Clean(filepath.ToSlash(p))
}

// escapes reports whether a cleaned path resolves outside the output directory.
func escapes(clean string) bool {
	return path.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, "../")
}

// Authorizes reports whether this record's affected_scope covers a path, by
// naming it literally or by matching a pattern entry.
//
// Both kinds count. A pattern authorizes matches without naming a file, which
// is what makes it usable for a path a record bounds but does not enumerate.
func (r *Record) Authorizes(target string) bool {
	target = normalize(target)
	for _, p := range r.Paths {
		if p == target {
			return true
		}
	}
	_, ok := pathpat.MatchAny(r.Patterns, target)
	return ok
}

// Authorize checks a set of paths against every record in the set.
//
// A path is authorized when some record's affected_scope covers it, and refused
// when a record covering it also forbids it. Replay uses this when records are
// supplied: it is asking whether the paths a recording names were authorized,
// which is a question about scope alone and involves no model and no phase that
// consults one (prov-2026-dc227be7).
//
// Every unauthorized path is reported rather than only the first, so a
// synthesised recording is corrected in one pass rather than one path at a
// time.
func (s *Set) Authorize(paths []string) error {
	var problems []error

	for _, target := range paths {
		var authorized bool
		for _, r := range s.Records {
			if !r.Authorizes(target) {
				continue
			}
			// A record that authorizes and forbids the same path has no
			// reading Sedum can act on, and ingestion already refuses that
			// combination for literal paths. A pattern reaching it here is the
			// same fault and is named the same way.
			if entry, forbidden := r.Forbids(target); forbidden {
				problems = append(problems, fmt.Errorf(
					"record %s authorizes %s and forbids it through forbidden_scope entry %s; the record has no reading Sedum can act on",
					r.ID, target, entry))
				continue
			}
			authorized = true
			break
		}
		if authorized {
			continue
		}
		problems = append(problems, fmt.Errorf(
			"path %s is not authorized by any supplied record; no affected_scope entry names or matches it", target))
	}

	return errors.Join(problems...)
}
