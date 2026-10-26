package pipeline

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/record"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/resolve"
)

// MigrateConfig holds what a marker migration reads (prov-2026-0a398351). It
// has no model client and no field for excluded directories: a migration calls
// no model, and passes none.
type MigrateConfig struct {
	Generators string
	Records    string
	Output     string
	Variables  map[string]string
	DryRun     bool
}

// MigrateNote names one region a migration reported. Why is empty for a
// migrated region and says why for a refused or a left one.
type MigrateNote struct {
	Path    string
	Action  string
	Variant string
	Why     string
}

// MigrateReport holds what a marker migration did. Migrated holds the regions
// whose marker was rewritten, Refused the regions whose kwargs cannot be hashed
// or whose action is not a key of their package's Actions map, and Left the v1
// regions left as they are. Files holds the paths rewritten, relative to the
// output directory and slash-separated, and Recordings the record ids whose
// recordings were written.
type MigrateReport struct {
	Migrated   []MigrateNote
	Refused    []MigrateNote
	Left       []MigrateNote
	Files      []string
	Recordings []string
	Skipped    int
}

// migratedRegion is one region a migration rewrote, with what its recording
// needs. Region holds the region as it was read, with its v1 marker, Entry the
// entry its recording gets, and Package the generator package whose comment
// syntax found it.
type migratedRegion struct {
	Path    string
	Package *genpkg.Package
	Region  inject.Region
	Entry   recording.Region
}

// splitOpeningLine splits the line starting at start - a region's opening
// marker line - into its indentation, the offset of the end of its text, and
// its line ending: a line feed, a carriage return and line feed, or nothing
// when the content ends first.
func splitOpeningLine(content string, start int) (indent string, end int, ending string) {
	line := content[start:]
	trimmed := strings.TrimLeft(line, " \t")
	indent = line[:len(line)-len(trimmed)]

	lf := strings.IndexByte(line, '\n')
	if lf < 0 {
		return indent, len(content), ""
	}
	if lf > 0 && line[lf-1] == '\r' {
		return indent, start + lf - 1, "\r\n"
	}
	return indent, start + lf, "\n"
}

// migrateContent rewrites the v1 markers in one file's text to v2 and returns
// the recording entries for them. A region with a foreign or v2 marker is not
// touched and not reported. A v1 region whose marker names no known record is
// left, one whose action the package does not declare or whose kwargs cannot be
// hashed is refused; both keep their v1 marker.
func migrateContent(path, content string, pkg *genpkg.Package, known func(string) bool) (string, []migratedRegion, []MigrateNote, []MigrateNote, error) {
	regions, err := inject.FindRegions(pkg.Comment(), content)
	if err != nil {
		return content, nil, nil, nil, fmt.Errorf("file %s: %w", path, err)
	}

	type rewrite struct {
		start, end int
		text       string
	}
	var (
		migrated []migratedRegion
		refused  []MigrateNote
		left     []MigrateNote
		edits    []rewrite
	)
	for _, region := range regions {
		m := region.Marker
		if m.Foreign || m.Version == 2 {
			continue
		}
		note := MigrateNote{Path: path, Action: m.Action, Variant: m.Variant}
		if m.Record == "" {
			note.Why = "marker names no record"
			left = append(left, note)
			continue
		}
		if !known(m.Record) {
			note.Why = "record not in the records directory"
			left = append(left, note)
			continue
		}
		action, declared := pkg.Actions[m.Action]
		if !declared {
			note.Why = "action no longer declared by its package"
			refused = append(refused, note)
			continue
		}
		hash, err := inject.KwargsHash(m.Kwargs)
		if err != nil {
			note.Why = err.Error()
			refused = append(refused, note)
			continue
		}

		id := map[string]any{}
		for _, name := range inject.SelectingKwargs(action) {
			if v, ok := m.Kwargs[name]; ok {
				id[name] = v
			}
		}
		next := m
		next.Version = 2
		next.KwargsHash = hash
		next.Body = region.BodyHash
		next.Kwargs = id
		open, err := next.Open(pkg.Comment())
		if err != nil {
			note.Why = err.Error()
			refused = append(refused, note)
			continue
		}

		indent, end, ending := splitOpeningLine(content, region.Start)
		edits = append(edits, rewrite{start: region.Start, end: end + len(ending), text: indent + open + ending})
		migrated = append(migrated, migratedRegion{
			Path:    path,
			Package: pkg,
			Region:  region,
			Entry: recording.Region{
				Path:       path,
				Action:     m.Action,
				Variant:    m.Variant,
				Tier:       string(m.Tier),
				Pkg:        m.Pkg,
				ID:         id,
				Kwargs:     m.Kwargs,
				KwargsHash: hash,
				Body:       region.BodyHash,
			},
		})
	}

	// From the last region to the first, so the offsets of the earlier ones
	// stay valid.
	out := content
	for i := len(edits) - 1; i >= 0; i-- {
		e := edits[i]
		out = out[:e.start] + e.text + out[e.end:]
	}
	return out, migrated, refused, left, nil
}

// buildMigrationRecording builds the recording of one record from the regions
// migrated for it. When the record already has a recording its files, phases,
// gen and created_at are kept and the migrated entries replace the entries
// with their keys; otherwise the recording is built from the record, the
// resolved files and the markers as they were.
func buildMigrationRecording(recordID string, regions []migratedRegion, records *record.Set, files []resolve.File, owners map[string]inject.Owner, variables map[string]string, existing *recording.Recording) recording.Recording {
	ordered := append([]migratedRegion(nil), regions...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Path != ordered[j].Path {
			return ordered[i].Path < ordered[j].Path
		}
		return ordered[i].Region.Start < ordered[j].Region.Start
	})

	var entries []recording.Region
	migratedKeys := map[string]bool{}
	for _, m := range ordered {
		entries = append(entries, m.Entry)
		migratedKeys[recording.RegionKey(m.Entry.Path, m.Entry.Action, m.Entry.Variant, m.Entry.ID)] = true
	}

	if existing != nil {
		for i, r := range existing.Records {
			if r.RecordID != recordID {
				continue
			}
			out := *existing
			out.Records = append([]recording.Record(nil), existing.Records...)
			var kept []recording.Region
			for _, e := range r.Regions {
				if !migratedKeys[recording.RegionKey(e.Path, e.Action, e.Variant, e.ID)] {
					kept = append(kept, e)
				}
			}
			r.Regions = append(kept, entries...)
			out.Records[i] = r
			return out
		}
	}

	out := recording.Recording{
		SedumVersion: recording.Version,
		Packages:     map[string]recording.Package{},
		Variables:    variables,
	}
	rec := recording.Record{
		RecordID:  recordID,
		CreatedAt: createdAtOf(records, recordID),
		Regions:   entries,
	}
	if len(ordered) > 0 {
		rec.Gen = ordered[0].Region.Marker.Gen
	}
	for _, f := range filesOf(files, recordID) {
		if f.Unmanaged || f.Package == nil {
			continue
		}
		rec.Files = append(rec.Files, recording.File{
			Path:     f.Path,
			Package:  f.Package.Name,
			Template: f.Template,
			Captures: capturesOf(f),
		})
		out.Packages[f.Package.Name] = recording.Package{Extensions: f.Package.Extensions}
	}

	var owned []inject.OwnedRegion
	for _, m := range ordered {
		owned = append(owned, inject.OwnedRegion{Path: m.Path, Package: m.Package, Region: m.Region})
	}
	rec.Phases = []recording.Phase{{
		Name:        recording.DefaultPhase,
		Invocations: priorFromMarkers(recordID, owned, owners),
	}}
	out.Records = []recording.Record{rec}
	return out
}

// Migrate rewrites every v1 marker under the output directory to v2 and writes
// the recordings of the records they name (prov-2026-0a398351). It calls no
// model and changes no region's text. A region it cannot migrate is reported
// and keeps its v1 marker. With DryRun nothing is written, and the report says
// what would have been.
func Migrate(cfg MigrateConfig) (*MigrateReport, error) {
	rules, err := projectSkipRules(cfg.Output, nil)
	if err != nil {
		return nil, err
	}

	packages, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{})
	if err != nil {
		return nil, err
	}
	var lines []string
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			lines = append(lines, "  "+f.String())
		}
	}
	if len(lines) > 0 {
		return nil, fmt.Errorf("generator packages did not load:\n%s", strings.Join(lines, "\n"))
	}

	variables, err := packages.ResolveVariables(cfg.Variables)
	if err != nil {
		return nil, err
	}
	records, _, err := record.Load(cfg.Records, record.Options{})
	if err != nil {
		return nil, err
	}
	owners := ownersOf(records)
	// A record's scope names every kind of file its change touches, and only
	// a file a package claims can carry markers, so the paths no package
	// claims are left out of the migration and counted.
	skipped := 0
	for i := range records.Records {
		var kept []string
		for _, p := range records.Records[i].Paths {
			if _, _, unmanaged := packages.Unmanaged(p); !unmanaged {
				if ext := strings.ToLower(path.Ext(p)); ext == "" || len(packages.ForExtension(ext)) == 0 {
					skipped++
					continue
				}
			}
			kept = append(kept, p)
		}
		records.Records[i].Paths = kept
	}
	resolutions, _, err := resolve.Paths(packages, records, nil)
	if err != nil {
		return nil, err
	}
	// Each path is created on its own. Migration reads only a file's path,
	// package, template and captures, so a file that cannot be reconciled with
	// its template - it is not rewritten for that - is carried by its
	// resolution alone rather than stopping the migration.
	var files []resolve.File
	for _, res := range resolutions {
		created, err := resolve.Create([]resolve.Resolution{res}, resolve.Options{Output: cfg.Output, DryRun: true, Variables: variables})
		if err != nil {
			created = []resolve.File{{Resolution: res}}
		}
		files = append(files, created...)
	}

	known := func(id string) bool {
		_, ok := records.Lookup(inject.Canonical(owners, id))
		return ok
	}

	type rewritten struct {
		path    string
		content string
		mode    fs.FileMode
	}
	var (
		report    = &MigrateReport{Skipped: skipped}
		rewrites  []rewritten
		byRecord  = map[string][]migratedRegion{}
		walkError error
	)
	err = walkMarked(packages, cfg.Output, []string{cfg.Generators, cfg.Records}, rules, func(rel string, pkg *genpkg.Package, _ []inject.Region) {
		if walkError != nil {
			return
		}
		full := filepath.Join(cfg.Output, filepath.FromSlash(rel))
		data, err := os.ReadFile(full)
		if err != nil {
			walkError = err
			return
		}
		info, err := os.Stat(full)
		if err != nil {
			walkError = err
			return
		}
		pkg = packageDeclaringActions(packages, rel, string(data), pkg)
		next, migrated, refused, left, err := migrateContent(rel, string(data), pkg, known)
		if err != nil {
			walkError = err
			return
		}
		report.Refused = append(report.Refused, refused...)
		report.Left = append(report.Left, left...)
		if len(migrated) == 0 {
			return
		}
		for _, m := range migrated {
			report.Migrated = append(report.Migrated, MigrateNote{Path: m.Path, Action: m.Entry.Action, Variant: m.Entry.Variant})
			id := inject.Canonical(owners, m.Region.Marker.Record)
			byRecord[id] = append(byRecord[id], m)
		}
		report.Files = append(report.Files, rel)
		rewrites = append(rewrites, rewritten{path: full, content: next, mode: info.Mode().Perm()})
	})
	if err != nil {
		return nil, err
	}
	if walkError != nil {
		return nil, walkError
	}

	dir := filepath.Join(cfg.Output, ".sedum", "recordings")
	ids := make([]string, 0, len(byRecord))
	for id := range byRecord {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var built []recording.Recording
	for _, id := range ids {
		var existing *recording.Recording
		got, err := recording.ReadRecord(dir, id)
		switch {
		case err == nil:
			existing = &got
		case errors.Is(err, fs.ErrNotExist):
		default:
			return nil, err
		}
		built = append(built, buildMigrationRecording(id, byRecord[id], records, files, owners, variables, existing))
	}
	report.Recordings = ids

	if !cfg.DryRun {
		// Every replacement file is staged beside the file it replaces before
		// anything is changed, so a file that cannot be replaced stops the
		// migration with the tree as it was. A rename needs the directory to
		// be writable, not the file, so a locked file is replaced and stays
		// locked.
		staged := make([]string, 0, len(rewrites))
		discard := func() {
			for _, tmp := range staged {
				_ = os.Remove(tmp)
			}
		}
		for _, w := range rewrites {
			tmp, err := stageFile(w.path, w.content, w.mode)
			if err != nil {
				discard()
				return nil, err
			}
			staged = append(staged, tmp)
		}
		for _, rec := range built {
			if err := recording.WriteRecords(dir, rec); err != nil {
				discard()
				return nil, err
			}
		}
		for i, w := range rewrites {
			if err := os.Rename(staged[i], w.path); err != nil {
				discard()
				return nil, err
			}
		}
	}

	byPath := func(notes []MigrateNote) {
		sort.SliceStable(notes, func(i, j int) bool {
			a, b := notes[i], notes[j]
			if a.Path != b.Path {
				return a.Path < b.Path
			}
			if a.Action != b.Action {
				return a.Action < b.Action
			}
			return a.Variant < b.Variant
		})
	}
	byPath(report.Migrated)
	byPath(report.Refused)
	byPath(report.Left)
	sort.Strings(report.Files)
	return report, nil
}

// packageDeclaringActions is the package, among those that claim a file's
// extension, whose actions the file's regions name: the first one that declares
// the most of them. Packages can share an extension, and a region is judged by
// the package that wrote it, not by one that merely claims the same files.
// found, the package the walk reached the file through, is the answer when no
// candidate does better.
func packageDeclaringActions(packages *genpkg.Set, rel, text string, found *genpkg.Package) *genpkg.Package {
	score := func(pkg *genpkg.Package) int {
		regions, err := inject.FindRegions(pkg.Comment(), text)
		if err != nil {
			return -1
		}
		declared := 0
		for _, r := range regions {
			if _, ok := pkg.Actions[r.Marker.Action]; ok {
				declared++
			}
		}
		return declared
	}

	best, bestScore := found, score(found)
	for _, candidate := range packages.ForExtension(strings.ToLower(path.Ext(rel))) {
		if s := score(candidate); s > bestScore {
			best, bestScore = candidate, s
		}
	}
	return best
}

// stageFile writes content to a new file in the directory of target, with the
// given mode, and returns its path. The file is removed again on failure.
func stageFile(target, content string, mode fs.FileMode) (string, error) {
	tmp, err := os.CreateTemp(filepath.Dir(target), ".sedum-migrate-*")
	if err != nil {
		return "", err
	}
	name := tmp.Name()
	fail := func(err error) (string, error) {
		_ = tmp.Close()
		_ = os.Remove(name)
		return "", err
	}
	if _, err := tmp.WriteString(content); err != nil {
		return fail(err)
	}
	if err := tmp.Chmod(mode); err != nil {
		return fail(err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return "", err
	}
	return name, nil
}
