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
}

// Result is what one invocation did.
type Result struct {
	Path    string
	Action  string
	Variant string
	// Replaced reports that an owned region with this invocation's identity
	// was already in the file and was overwritten in place.
	Replaced bool
	// Skipped reports that the region was seeded, so it was left as it is.
	Skipped bool
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

	for _, inv := range invocations {
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

		updated, result, err := applyOne(inv, content)
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

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
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
func applyOne(inv Invocation, content string) (string, Result, error) {
	result := Result{Path: inv.Path, Action: inv.Action.Name, Variant: inv.Variant}

	marker := Marker{
		Action:  inv.Action.Name,
		Variant: inv.Variant,
		Tier:    inv.Tier,
		Record:  inv.RecordID,
		Kwargs:  inv.Kwargs,
	}

	selecting := selectingKwargs(inv.Action)
	identity, err := IdentityOf(marker, selecting)
	if err != nil {
		return "", result, err
	}

	regions, err := FindRegions(inv.Package.CommentPrefix, content)
	if err != nil {
		return "", result, fmt.Errorf("file %s: %w", inv.Path, err)
	}

	// An existing region with the same identity is this invocation's region,
	// whichever record wrote it. Replacement follows from two invocations
	// being the same thing, never from a lifecycle field Sedum would have to
	// reimplement a provenance graph to read.
	for _, region := range regions {
		existing, err := IdentityOf(region.Marker, selecting)
		if err != nil {
			return "", result, fmt.Errorf("file %s: %w", inv.Path, err)
		}
		if existing != identity {
			continue
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

	offset := skipRegionsAt(regions, where.offset)
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
	open, err := marker.Open(inv.Package.CommentPrefix)
	if err != nil {
		return "", err
	}

	// The template's own trailing newline is normalized away so that a
	// template written with one and a template written without produce the
	// same region.
	body := strings.TrimRight(inv.Content, "\n")

	var b strings.Builder
	b.WriteString(indent)
	b.WriteString(open)
	b.WriteString("\n")
	if body != "" {
		b.WriteString(body)
		b.WriteString("\n")
	}
	b.WriteString(indent)
	b.WriteString(marker.Close(inv.Package.CommentPrefix))
	b.WriteString("\n")
	return b.String(), nil
}

// skipRegionsAt advances past regions already sitting at the insertion point,
// so that repeated injections at one anchor accumulate in the order they were
// applied rather than in reverse.
func skipRegionsAt(regions []Region, offset int) int {
	moved := true
	for moved {
		moved = false
		for _, region := range regions {
			if region.Start == offset {
				offset = region.End
				moved = true
			}
		}
	}
	return offset
}

// selectingKwargs returns the kwarg names that take part in a region's
// identity: the ones the action declares required.
func selectingKwargs(action *genpkg.Action) []string {
	var out []string
	for name, kwarg := range action.Kwargs {
		if kwarg.Required {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}
