// Package inject runs Phase 7: it writes a rendered action's output into the
// file and the anchored region the action names.
//
// Idempotency state lives only in the ownership markers in the generated files.
// There is no sidecar cache, no manifest, no lockfile, and nothing is read back
// from the run log. Re-running replaces the region an action owns rather than
// appending beside it, which is what makes reruns and partial regeneration safe
// and what produces the audit trail as a side effect: grepping markers yields
// file -> action -> variant -> arguments with no maintained state, and ownership
// is visible in the diff a human reviews.
package inject

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/runlog"
	"github.com/livecodelife/sedum/internal/syntax"
)

// Invocation is one action resolved to the point where only writing is left:
// Phase 6 has rendered the template, chosen the variant, and decided the file.
type Invocation struct {
	Package *genpkg.Package
	Action  *genpkg.Action
	// Variant is the discriminator value that selected the template, empty
	// for an action that declares no discriminator.
	Variant string
	// Kwargs are the bound arguments, recorded whole on the marker.
	Kwargs map[string]any
	// Path is the authorized path the action injects into, relative and
	// slash-separated.
	Path string
	// RecordID is the provenance record this invocation came from. It is
	// recorded on the marker and never matched on.
	RecordID string
	// Gen is that record's gen_hash, recorded on the marker. Empty when the
	// run had no record to hash, which a later grow reads as stale.
	Gen string
	// Content is the rendered template output.
	Content string
	// Tier declares whether Sedum may overwrite the region on a later run.
	// Empty means DefaultTier.
	Tier Tier
}

// Options controls Phase 7.
type Options struct {
	// Output is the directory the authorized paths live under.
	Output string
	// DryRun decides everything and writes nothing.
	DryRun bool

	// Unwritten is what Phase 3 rendered for the files it did not write,
	// keyed by authorized path. It exists for --dry-run: injecting into a
	// file a dry run declined to create would otherwise report a missing file
	// for every path the same run just said it would create.
	//
	// It is consulted only for a path that is not on disk, so a real run
	// never reads it and an existing file is always read from disk. That is
	// what keeps a dry run's report accurate rather than merely successful: a
	// file carrying regions from an earlier run reports a replacement, and a
	// file that would be new reports a first injection (prov-2026-23653fdc).
	Unwritten map[string]string

	// Log is the run log. A nil log discards.
	Log *runlog.Log

	// Owners are the records the run knows, keyed by id: every record in the
	// records directory, not only the ones this run grows, because ownership
	// and anchor order concern regions any record wrote (prov-2026-375f5d4a).
	// A record absent here still owns the regions naming its id; it only
	// lacks aliases and a creation date.
	Owners map[string]Owner

	// Reconcile names the records whose owned regions this run reconciles
	// against their selection: an owned region of one of these records that
	// it did not select again is pruned. Grow sets it for every record it
	// asked the model about. Replay never does, because a recording may be
	// partial (OPEN_QUESTIONS section 16).
	Reconcile []Reconcile

	// ReconcileReport turns reconciliation into a report: a region a record in
	// Reconcile owns and this run did not select is returned as Unselected and
	// left in the file. Replay of a partial recording sets it
	// (prov-2026-c637c262).
	ReconcileReport bool

	// PruneOrphansIn are the files whose orphans this run removes: owned-tier
	// regions naming a record Owners does not hold (prov-2026-cc73cf61). Owners
	// must then hold every current record. Only files some record authorizes
	// belong here, because Sedum writes nothing else.
	PruneOrphansIn []Target

	// capture, when set, receives every touched file's content after
	// application, and the run neither writes nor runs the syntax guard: it is
	// how Preview lets Phase 5 check what Phase 7 would write.
	capture *map[string]string

	// Force allows a record whose selection is empty to prune every region it
	// owns. Without it that is an error, because an empty reply is more often
	// a flaky model than an intent to delete a record's code.
	Force bool
}

// Owner is what Phase 7 knows about a record that may own regions.
type Owner struct {
	ID string
	// Aliases are ids folded into this record, whose regions it owns.
	Aliases []string
	// CreatedAt orders owners at a shared anchor, before the id does.
	CreatedAt string
}

// Reconcile is one record whose owned regions are reconciled in this run.
type Reconcile struct {
	RecordID string
	// Targets are the files the record is authorized to inject into in this
	// run. Only these are pruned.
	Targets []Target
}

// Target is one file a reconciled record may own regions in, with the package
// whose actions and comment syntax govern it.
type Target struct {
	Path    string
	Package *genpkg.Package
}

// Result is what one invocation did.
type Result struct {
	Path    string
	Action  string
	Variant string
	// Record is the record that wrote the region, or for a pruned region the
	// record that owned it.
	Record string
	// Replaced reports that an owned region with this invocation's identity
	// was already in the file and was overwritten in place.
	Replaced bool
	// Skipped reports that the region was seeded, so it was left as it is.
	Skipped bool
	// Pruned reports that the region's owner did not select it again, so it
	// was removed.
	Pruned bool
	// Unselected reports a region its owner did not select this run, left in
	// place because the run only reports (Options.ReconcileReport).
	Unselected bool
	// Held reports a prune that only --force would carry out: the owner's
	// selection was empty. Reported by a dry run; a real run fails instead.
	Held bool
	// Orphan reports that a pruned region's record is no longer current:
	// removed from the records directory, or deprecated.
	Orphan bool
	// Changed reports that the file's bytes for this region changed: an
	// insertion, a prune, or a replacement whose rendering differs from what
	// was there. A replacement rendering the same bytes is not a change
	// (prov-2026-bf01cd6b), which is what lets a caller run only the specs
	// whose code moved.
	Changed bool
}

// Apply runs Phase 7 over every invocation.
//
// Invocations targeting one file are applied in order against a single buffer
// and the file is written once, so that the order they were declared in is the
// order they appear in and a rerun produces byte-identical output.
func Apply(invocations []Invocation, opts Options) ([]Result, error) {
	log := opts.Log
	if log == nil {
		log = runlog.Discard()
	}

	var (
		results  []Result
		problems []error
		// buffers holds each touched file's evolving content, keyed by the
		// authorized path.
		buffers = map[string]string{}
		touched []string
	)

	// Two invocations resolving to one region in one run is an authoring
	// fault, and under last-write-wins one of the two silently does not
	// happen. Keyed by path as well as identity, because one identity in two
	// files is two regions.
	claimed := map[placement]string{}
	own := newOwnership(opts.Owners)
	selectedBy := map[string]int{}

	for _, inv := range invocations {
		selectedBy[inv.RecordID]++
		if identity, err := identityFor(inv.Action, Marker{
			Action: inv.Action.Name, Variant: inv.Variant, Kwargs: inv.Kwargs,
		}); err == nil {
			where := placement{path: inv.Path, identity: identity}
			if first, seen := claimed[where]; seen {
				problems = append(problems, fmt.Errorf(
					"file %s: %s claims the region already claimed by record %s in this run; both resolve to identity %s, so one would silently replace the other",
					inv.Path, inv.Action.Name, first, identity.Key))
				continue
			}
			claimed[where] = inv.RecordID
		}

		content, ok := buffers[inv.Path]
		if !ok {
			loaded, err := readTarget(opts.Output, opts.Unwritten, inv)
			if err != nil {
				problems = append(problems, err)
				continue
			}
			content = loaded
			touched = append(touched, inv.Path)
		}

		updated, result, err := applyOne(inv, content, own)
		if err != nil {
			problems = append(problems, err)
			// The buffer is still recorded so that later invocations into
			// the same file report their own problems rather than
			// re-reporting this one.
			buffers[inv.Path] = content
			continue
		}

		buffers[inv.Path] = updated
		results = append(results, result)

		switch {
		case result.Skipped:
			log.Info("region is seeded, left as it is",
				"path", inv.Path, "action", inv.Action.Name, "variant", inv.Variant)
		case result.Replaced:
			log.Info("replaced owned region",
				"path", inv.Path, "action", inv.Action.Name, "variant", inv.Variant, "record", inv.RecordID)
		default:
			log.Info("injected region",
				"path", inv.Path, "action", inv.Action.Name, "variant", inv.Variant,
				"anchor", inv.Action.Anchor, "record", inv.RecordID)
		}
	}

	// Reconciliation runs against the buffers every invocation has already
	// been applied to, and only when nothing above failed: a prune is never
	// written unless the whole run is (prov-2026-375f5d4a).
	if len(problems) == 0 {
		for _, rec := range opts.Reconcile {
			pruned, err := reconcile(rec, buffers, &touched, claimed, selectedBy[rec.RecordID], own, opts)
			if err != nil {
				problems = append(problems, err)
				continue
			}
			for _, r := range pruned {
				if r.Held {
					log.Info("region would be pruned with --force",
						"path", r.Path, "action", r.Action, "variant", r.Variant, "record", r.Record)
				} else {
					log.Info("pruned region its record no longer selects",
						"path", r.Path, "action", r.Action, "variant", r.Variant, "record", r.Record)
				}
			}
			results = append(results, pruned...)
		}
		for _, target := range opts.PruneOrphansIn {
			pruned, err := pruneOrphans(target, buffers, &touched, own, opts)
			if err != nil {
				problems = append(problems, err)
				continue
			}
			for _, r := range pruned {
				log.Info("pruned orphan region",
					"path", r.Path, "action", r.Action, "variant", r.Variant, "record", r.Record)
			}
			results = append(results, pruned...)
		}
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	if opts.capture != nil {
		*opts.capture = buffers
		return results, nil
	}

	// The syntax guard: every touched file of a package that opted into a
	// built-in checker must pass it, or nothing is written
	// (prov-2026-5403d162). Phase 5 has already shown the model any failure;
	// this is what makes the guarantee hold for every path into Phase 7.
	if err := checkSyntax(opts.Output, touched, buffers, packagesByPath(invocations, opts)); err != nil {
		return nil, err
	}
	if opts.DryRun {
		return results, nil
	}

	sort.Strings(touched)
	for _, path := range touched {
		full := filepath.Join(opts.Output, filepath.FromSlash(path))
		if err := os.WriteFile(full, []byte(buffers[path]), 0o644); err != nil {
			problems = append(problems, fmt.Errorf("write %s: %w", path, err))
		}
	}
	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	return results, nil
}

// readTarget loads the file an invocation injects into.
//
// A path that is not there is a hard error rather than something to create.
// Nothing is created that a provenance record did not authorize, so a record
// naming an implementation file but omitting its header produces exactly this
// failure, naming the file the record forgot.
//
// A dry run supplies what Phase 3 would have written for the files it declined
// to create, and that is consulted only after the filesystem has said the path
// is absent. A path in neither place is the hard error above, unchanged: a
// record that forgot a companion file fails a dry run exactly as it fails a
// real one, which is most of the value of running one.
func readTarget(output string, unwritten map[string]string, inv Invocation) (string, error) {
	full := filepath.Join(output, filepath.FromSlash(inv.Path))

	data, err := os.ReadFile(full)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		if rendered, ok := unwritten[inv.Path]; ok {
			return rendered, nil
		}
		return "", fmt.Errorf(
			"action %s injects into %s, which no provenance record authorized: the file was not created, so there is nowhere to put the injection",
			inv.Action.Name, inv.Path)
	case err != nil:
		return "", fmt.Errorf("read %s: %w", inv.Path, err)
	}
	return string(data), nil
}

// applyOne injects one invocation into one file's content.
func applyOne(inv Invocation, content string, own ownership) (string, Result, error) {
	result := Result{Path: inv.Path, Action: inv.Action.Name, Variant: inv.Variant, Record: inv.RecordID}

	marker := Marker{
		Action:  inv.Action.Name,
		Variant: inv.Variant,
		Tier:    inv.Tier,
		Record:  inv.RecordID,
		Gen:     inv.Gen,
		// A composite never owns a region: the marker names the child that
		// rendered it, so the hash is the child's, which is what re-rendering
		// would replay.
		Pkg:    inv.Action.PkgHash,
		Kwargs: inv.Kwargs,
	}

	identity, err := identityFor(inv.Action, marker)
	if err != nil {
		return "", result, err
	}

	regions, err := FindRegions(inv.Package.Comment(), content)
	if err != nil {
		return "", result, fmt.Errorf("file %s: %w", inv.Path, err)
	}

	// An existing region with the same identity is this invocation's region,
	// whichever record wrote it. Replacement follows from two invocations
	// being the same thing, never from a lifecycle field Sedum would have to
	// reimplement a provenance graph to read.
	for _, region := range regions {
		existing, err := identityFor(inv.Action, region.Marker)
		if err != nil {
			return "", result, fmt.Errorf("file %s: %w", inv.Path, err)
		}
		if existing != identity {
			continue
		}

		// A region belongs to the record its marker names, and only that
		// record rewrites it. Another record selecting it means a new record
		// was written where the owner should have been revised, which is
		// an error rather than an update (prov-2026-375f5d4a). A region
		// naming no record is unowned and is claimed here.
		if owner := region.Marker.Record; owner != "" && !own.owns(inv.RecordID, owner) {
			return "", result, fmt.Errorf(
				"file %s: record %s selected %s, which record %s owns (identity %s); revise %s, or change %s's selection",
				inv.Path, describeRecord(inv.RecordID), region.Marker.Label(), owner, identity.Key, owner, describeRecord(inv.RecordID))
		}

		if region.Marker.Tier == TierSeeded {
			// Generated once and never touched again.
			result.Skipped = true
			return content, result, nil
		}

		// The marker written back carries whatever attributes this version
		// of Sedum does not model, taken from the marker being replaced.
		// Retaining them on read is only half of the promise: this is the
		// write, and without it every annotation another tool left on the
		// region would disappear on the first rerun (prov-2026-72775ae5).
		//
		// Nothing else is taken from the old marker. Tier, record, kwargs and
		// writer describe what the region is now, and this invocation is what
		// it is now.
		refreshed := marker
		refreshed.Extra = region.Marker.Extra

		// Alignment comes from the region's own opening marker rather than
		// from the anchor, so a run refreshing a region's contents does not
		// silently re-lay-out one whose surroundings have moved.
		rendered, err := renderRegion(inv, refreshed, indentAt(content, region.Start))
		if err != nil {
			return "", result, err
		}

		result.Replaced = true
		result.Changed = rendered != content[region.Start:region.End]
		return content[:region.Start] + rendered + content[region.End:], result, nil
	}

	where, err := locate(inv.Package, inv.Action, content)
	if err != nil {
		return "", result, fmt.Errorf(
			"action %s cannot be injected into %s: %w; the file is not shaped the way the action assumes",
			inv.Action.Name, inv.Path, err)
	}

	rendered, err := renderRegion(inv, marker, where.indent)
	if err != nil {
		return "", result, err
	}

	offset := orderedOffset(regions, where.offset, own, inv.RecordID)
	result.Changed = true
	return content[:offset] + rendered + content[offset:], result, nil
}

// renderRegion wraps an invocation's rendered content in its ownership markers,
// aligned to indent.
//
// Only the markers are indented. The body is written exactly as the template
// rendered it, because a template author writes each template at the depth its
// anchor sits at - a single package legitimately mixes a fragment indented to
// sit inside a struct with a top-level declaration starting at column zero, and
// re-indenting would double-indent the first (prov-2026-df491217).
func renderRegion(inv Invocation, marker Marker, indent string) (string, error) {
	open, err := marker.Open(inv.Package.Comment())
	if err != nil {
		return "", err
	}

	// The template's own trailing newline is normalized away so that a
	// template written with one and a template written without produce the
	// same region.
	//
	// A template that ends in a blank line keeps exactly one before the
	// closing marker (prov-2026-ab24df5a). Some targets' formatters want a
	// declaration separated from a comment that follows it - gofmt rewrites a
	// top-level func followed directly by its closing marker - and only the
	// template author knows what the region holds; Sedum knows no language.
	body := strings.TrimRight(inv.Content, "\n")
	if strings.HasSuffix(inv.Content, "\n\n") && body != "" {
		body += "\n"
	}

	var b strings.Builder
	b.WriteString(indent)
	b.WriteString(open)
	b.WriteString("\n")
	if body != "" {
		b.WriteString(body)
		b.WriteString("\n")
	}
	b.WriteString(indent)
	b.WriteString(marker.Close(inv.Package.Comment()))
	b.WriteString("\n")
	return b.String(), nil
}

// orderedOffset finds where a new region goes among the regions already
// sitting at its insertion point: after every region whose owner sorts at or
// before its own, and before the first that sorts after (prov-2026-375f5d4a).
//
// Owners sort by creation date, then id, so a file regenerated from scratch
// comes out the same whatever order its records were grown in. One owner's
// regions keep the order they were applied in, because a region never sorts
// after its own owner. Nothing already in the file moves.
func orderedOffset(regions []Region, offset int, own ownership, recordID string) int {
	key := own.key(recordID)
	for {
		next := -1
		for i, region := range regions {
			if region.Start == offset {
				next = i
				break
			}
		}
		if next < 0 || own.key(regions[next].Marker.Record) > key {
			return offset
		}
		offset = regions[next].End
	}
}

// identityFor computes the identity of the region a marker describes, under the
// action's own account of what identifies one.
//
// The variant takes part only where the action declares no identity. A variant
// chooses which template renders, which is a statement about content and not
// about which region is meant — so an author who has said what identifies a
// region has already excluded it. Left in, flipping a discriminated value would
// mint a second region beside the first and both would be emitted: a rule with
// two severities, a section with two requirements (prov-2026-d71d6f76).
func identityFor(action *genpkg.Action, m Marker) (Identity, error) {
	if len(action.Identity) > 0 {
		m.Variant = ""
	}
	return IdentityOf(m, selectingKwargs(action))
}

// selectingKwargs returns the kwarg names that take part in a region's
// identity: the ones the action declares as its identity, or every required
// kwarg where it declares none.
//
// The fallback is what every package authored before identity existed was
// written against, and it stays correct for an action with no discriminator and
// no key: two imports differ because `symbol` is required. The declaration is
// for the opposite case, an action whose required kwargs are mostly content —
// a pattern, a count, a selector — where selecting on all of them means a
// changed value mints a second region rather than refining the one that is
// there (prov-2026-d71d6f76).
func selectingKwargs(action *genpkg.Action) []string {
	if len(action.Identity) > 0 {
		out := append([]string(nil), action.Identity...)
		sort.Strings(out)
		return out
	}

	var out []string
	for name, kwarg := range action.Kwargs {
		if kwarg.Required {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

// placement is one region in one file. Keyed by path as well as identity,
// because one identity in two files is two regions.
type placement struct {
	path     string
	identity Identity
}

// ownership answers which record owns a region and how owners sort.
type ownership struct {
	byID map[string]Owner
	// alias maps an alias to the id of the record that folded it in.
	alias map[string]string
}

func newOwnership(owners map[string]Owner) ownership {
	own := ownership{byID: owners, alias: map[string]string{}}
	for id, o := range owners {
		for _, a := range o.Aliases {
			own.alias[a] = id
		}
	}
	return own
}

// canonical returns the record id a marker's record names, following an alias.
func (o ownership) canonical(id string) string {
	if current, ok := o.alias[id]; ok {
		return current
	}
	return id
}

// owns reports whether the record recordID owns a region whose marker names
// owner.
func (o ownership) owns(recordID, owner string) bool {
	return owner == recordID || o.canonical(owner) == recordID
}

// orphan reports whether a region is owned-tier and names a record this
// ownership does not hold.
func (o ownership) orphan(m Marker) bool {
	if m.Tier == TierSeeded || m.Record == "" {
		return false
	}
	_, ok := o.byID[o.canonical(m.Record)]
	return !ok
}

// key is the sort key of a region's owner at a shared anchor: creation date,
// then id. A region naming no record sorts first, and an owner the run does
// not know sorts by id alone, ahead of every dated owner - deterministic
// either way, which is what reproducing a file needs.
func (o ownership) key(owner string) string {
	if owner == "" {
		return ""
	}
	id := o.canonical(owner)
	return o.byID[id].CreatedAt + "\x00" + id
}

// describeRecord names a record in a diagnostic, including the absence of one.
func describeRecord(id string) string {
	if id == "" {
		return "(none)"
	}
	return id
}

// reconcile prunes the owned regions of one record that it did not select in
// this run, in the files it is authorized for (prov-2026-375f5d4a).
//
// A region survives when its owner selected its identity in that file this
// run - which is every region the run just replaced or inserted for it. Seeded
// regions, other owners' regions and unowned regions are never touched. A
// region whose action the package no longer declares cannot have been
// selected, so it is pruned too.
//
// An empty selection over owned regions is refused unless opts.Force, because
// an empty reply is more often a flaky model than an intent to delete a
// record's code. A dry run reports those prunes as held instead of failing, so
// a person can see what --force would remove.
func reconcile(rec Reconcile, buffers map[string]string, touched *[]string, claimed map[placement]string, selected int, own ownership, opts Options) ([]Result, error) {
	type cut struct {
		path   string
		region Region
	}
	var cuts []cut
	loaded := map[string]string{}

	for _, target := range rec.Targets {
		content, ok := buffers[target.Path]
		if !ok {
			var found bool
			content, found = readExisting(opts.Output, opts.Unwritten, target.Path)
			if !found {
				// Nothing on disk means nothing to prune.
				continue
			}
		}
		loaded[target.Path] = content

		regions, err := FindRegions(target.Package.Comment(), content)
		if err != nil {
			return nil, fmt.Errorf("file %s: %w", target.Path, err)
		}
		for _, region := range regions {
			m := region.Marker
			if m.Tier == TierSeeded || m.Record == "" || !own.owns(rec.RecordID, m.Record) {
				continue
			}
			if action, ok := target.Package.Actions[m.Action]; ok {
				identity, err := identityFor(action, m)
				if err != nil {
					return nil, fmt.Errorf("file %s: %w", target.Path, err)
				}
				if claimed[placement{path: target.Path, identity: identity}] == rec.RecordID {
					continue
				}
			}
			cuts = append(cuts, cut{path: target.Path, region: region})
		}
	}

	if len(cuts) == 0 {
		return nil, nil
	}

	if opts.ReconcileReport {
		var results []Result
		for _, c := range cuts {
			results = append(results, Result{
				Path: c.path, Action: c.region.Marker.Action, Variant: c.region.Marker.Variant,
				Record: c.region.Marker.Record, Unselected: true,
			})
		}
		return results, nil
	}

	held := selected == 0 && !opts.Force
	if held && !opts.DryRun {
		var where []string
		for _, c := range cuts {
			where = append(where, c.path+" "+c.region.Marker.Label())
		}
		return nil, fmt.Errorf(
			"record %s selected nothing but owns %d region(s): %s; rerun with --force to remove them, or check the model's reply",
			rec.RecordID, len(cuts), strings.Join(where, ", "))
	}

	var results []Result
	// Cut from the end of each file backwards so earlier offsets stay valid.
	for i := len(cuts) - 1; i >= 0; i-- {
		c := cuts[i]
		results = append([]Result{{
			Path: c.path, Action: c.region.Marker.Action, Variant: c.region.Marker.Variant,
			Record: c.region.Marker.Record, Pruned: true, Held: held, Changed: !held,
		}}, results...)
		if held {
			continue
		}
		content := loaded[c.path]
		loaded[c.path] = content[:c.region.Start] + content[c.region.End:]
	}

	if !held {
		changed := map[string]bool{}
		for _, c := range cuts {
			changed[c.path] = true
		}
		for path := range changed {
			if _, seen := buffers[path]; !seen {
				*touched = append(*touched, path)
			}
			buffers[path] = loaded[path]
		}
	}
	return results, nil
}

// readExisting reads a file Phase 7 may prune in, from disk or, in a dry run,
// from what Phase 3 would have written. It reports whether there was one.
func readExisting(output string, unwritten map[string]string, path string) (string, bool) {
	data, err := os.ReadFile(filepath.Join(output, filepath.FromSlash(path)))
	if err == nil {
		return string(data), true
	}
	if content, ok := unwritten[path]; ok {
		return content, true
	}
	return "", false
}

// pruneOrphans removes the orphan regions in one file (prov-2026-cc73cf61).
func pruneOrphans(target Target, buffers map[string]string, touched *[]string, own ownership, opts Options) ([]Result, error) {
	content, ok := buffers[target.Path]
	if !ok {
		content, ok = readExisting(opts.Output, opts.Unwritten, target.Path)
		if !ok {
			return nil, nil
		}
	}
	regions, err := FindRegions(target.Package.Comment(), content)
	if err != nil {
		return nil, fmt.Errorf("file %s: %w", target.Path, err)
	}

	var results []Result
	for i := len(regions) - 1; i >= 0; i-- {
		m := regions[i].Marker
		if !own.orphan(m) {
			continue
		}
		content = content[:regions[i].Start] + content[regions[i].End:]
		results = append([]Result{{
			Path: target.Path, Action: m.Action, Variant: m.Variant, Record: m.Record,
			Pruned: true, Orphan: true, Changed: true,
		}}, results...)
	}
	if len(results) > 0 {
		if _, seen := buffers[target.Path]; !seen {
			*touched = append(*touched, target.Path)
		}
		buffers[target.Path] = content
	}
	return results, nil
}

// Preview applies invocations in memory and returns every touched file's
// content, writing nothing and running no syntax guard. Phase 5 uses it to
// check, before an answer is accepted, the files Phase 7 would write
// (prov-2026-5403d162).
func Preview(invocations []Invocation, opts Options) (map[string]string, error) {
	var buffers map[string]string
	opts.capture = &buffers
	opts.DryRun = true
	if _, err := Apply(invocations, opts); err != nil {
		return nil, err
	}
	return buffers, nil
}

// packagesByPath is the package governing each file this run may touch.
func packagesByPath(invocations []Invocation, opts Options) map[string]*genpkg.Package {
	out := map[string]*genpkg.Package{}
	for _, inv := range invocations {
		out[inv.Path] = inv.Package
	}
	for _, rec := range opts.Reconcile {
		for _, t := range rec.Targets {
			out[t.Path] = t.Package
		}
	}
	for _, t := range opts.PruneOrphansIn {
		out[t.Path] = t.Package
	}
	return out
}

// checkSyntax runs each touched file of a package with a syntax: key through
// that built-in checker.
func checkSyntax(output string, touched []string, buffers map[string]string, packages map[string]*genpkg.Package) error {
	var problems []error
	for _, path := range touched {
		pkg := packages[path]
		if pkg == nil || pkg.Syntax == "" {
			continue
		}
		checker, ok := syntax.Lookup(pkg.Syntax)
		if !ok {
			continue // rejected at load; unreachable for a loaded package
		}
		file := syntax.File{Path: path, Content: buffers[path], Siblings: syntax.Siblings(output, path, buffers)}
		if err := checker.Check(file); err != nil {
			problems = append(problems, fmt.Errorf("file %s would not pass the %s check, so nothing was written: %w", path, pkg.Syntax, err))
		}
	}
	return errors.Join(problems...)
}

// Canonical returns the record id a marker's record names: the surviving id
// when it names an alias of a record owners holds, and the name itself
// otherwise (prov-2026-bf01cd6b).
func Canonical(owners map[string]Owner, id string) string {
	return newOwnership(owners).canonical(id)
}
