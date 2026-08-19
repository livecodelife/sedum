package evals

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FixtureDigest is a fingerprint of everything a case's numbers were drawn
// against: its generator package set, its records, and the behaviour target it
// names.
//
// It exists so that history can tell an entry drawn against one fixture from an
// entry drawn against another. Seven imprints under prov-2026-d318e3dd changed
// the packages materially - an API-only Rails application, a controller that
// stops skipping a filter, a Go service that reads its port from the
// environment - and every entry stored before them describes packages that no
// longer exist. Those entries are true and stay; what stops is comparing across
// the change (prov-2026-43505e1b).
//
// Computed from file contents and paths relative to each root, so the same
// fixture digests identically on any machine and in any checkout.
func FixtureDigest(c Case) (string, error) {
	sum := sha256.New()

	// Sorted, so the digest does not depend on directory iteration order.
	roots := []struct{ label, dir string }{
		{"generators", c.Generators},
		{"records", c.Records},
	}
	if c.Expect.Behavior != nil && c.Expect.Behavior.Target != "" {
		// The target is included even for a run that did not measure
		// behaviour. It errs toward incomparability, which is the safe
		// direction: a case that names a target is the same case whether or
		// not a given run spent the minutes.
		roots = append(roots, struct{ label, dir string }{
			"target", filepath.Join("behavior", "targets", c.Expect.Behavior.Target),
		})
	}

	for _, root := range roots {
		if root.dir == "" {
			// A baseline case names no generators, and its absence is the arm
			// rather than an omission. It is hashed as absent so that a
			// baseline and a sedum entry never collide.
			fmt.Fprintf(sum, "%s\x00absent\x00", root.label)
			continue
		}
		if err := hashTree(sum, root.label, root.dir); err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(sum.Sum(nil))[:12], nil
}

// hashTree folds every file under dir into sum, by path relative to dir and by
// contents.
//
// Paths are relative so that two checkouts of the same commit agree, and they
// are sorted so that the walk order cannot change the answer.
func hashTree(sum io.Writer, label, dir string) error {
	var paths []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		paths = append(paths, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return fmt.Errorf("digesting %s at %s: %w", label, dir, err)
	}
	sort.Strings(paths)

	fmt.Fprintf(sum, "%s\x00", label)
	for _, rel := range paths {
		content, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
		if err != nil {
			return fmt.Errorf("digesting %s at %s: %w", label, rel, err)
		}
		fmt.Fprintf(sum, "%s\x00%d\x00", rel, len(content))
		if _, err := sum.Write(content); err != nil {
			return err
		}
	}
	return nil
}

// sameFixture reports whether two entries were drawn against the same fixture.
//
// Two entries that both predate the field compare as they always did. That is a
// compromise and not a claim: entries written before the digest existed may well
// have been drawn against different packages from each other, and there is no
// way to find out now. Marking every one of them would add a mark to every
// historical row while telling the reader nothing they could act on, and the
// commit column already says these are different points in time.
//
// What the field can speak to is the boundary. A digest against no digest is
// different, so exactly one mark appears where the record of what was measured
// begins.
func sameFixture(a, b string) bool {
	return strings.EqualFold(a, b)
}
