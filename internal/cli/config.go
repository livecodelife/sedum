package cli

import (
	"errors"
	"fmt"
	"strings"
)

// The config types below are the whole contract between flag parsing and the
// pipeline. A command's job is to fill one in, check the interdependence rules
// the flags cannot express on their own, and hand it off. No filesystem work,
// no template handling, and no pipeline logic lives in this package.

// GrowConfig drives the full pipeline.
type GrowConfig struct {
	Generators string
	Records    string
	Output     string
	Lang       []string
	Only       []string
	RecordTo   string
	Vars       []string
	Execute    string
	DryRun     bool
	StopAfter  string
	Retries    int
	Model      string
	LogPath    string
	Verbose    bool
}

// Replaying reports whether this run replays a recording instead of invoking a
// model.
func (c *GrowConfig) Replaying() bool { return c.Execute != "" }

// Validate checks the flag interdependence rules cobra cannot express. It runs
// before any work begins, so an unusable combination costs nothing.
func (c *GrowConfig) Validate() error {
	if !c.Replaying() && c.Records == "" {
		return errors.New("--records is required unless --execute names a recording to replay")
	}

	if c.StopAfter == "" {
		return nil
	}

	sp, ok := lookupStopPoint(c.StopAfter)
	if !ok {
		return fmt.Errorf("--stop-after %q is not a phase boundary; expected one of: %s", c.StopAfter, stopPointNames())
	}

	if c.Replaying() {
		if !sp.runByReplay {
			return fmt.Errorf(
				"--stop-after %s cannot be combined with --execute: replay enters the pipeline at file creation with %s already decided, so that phase never runs",
				sp.name, sp.name)
		}
		return nil
	}

	if sp.requiresRecord && c.RecordTo == "" {
		return fmt.Errorf(
			"--stop-after %s requires --record: without it the model's output is discarded and there is nothing to resume from, having already paid for the call",
			sp.name)
	}

	return nil
}

// IgnoredFlags names the flags this configuration sets that the run will not
// consult, so the command can say so rather than appearing to honor them.
func (c *GrowConfig) IgnoredFlags() []string {
	if !c.Replaying() {
		return nil
	}

	var ignored []string
	if len(c.Lang) > 0 {
		ignored = append(ignored, "--lang (package resolution is recorded per file)")
	}
	if c.Model != "" {
		ignored = append(ignored, "--model (replay invokes no model)")
	}
	if c.Retries != defaultRetries {
		ignored = append(ignored, "--retries (replay validation is terminal, never re-prompted)")
	}
	return ignored
}

// ValidateConfig drives Phase 0 in isolation.
type ValidateConfig struct {
	Generators string
	Packages   []string
	Strict     bool
}

// ResolveConfig drives Phases 0 through 3 with no model.
type ResolveConfig struct {
	Generators   string
	Records      string
	Lang         []string
	Only         []string
	Vars         []string
	ShowTemplate bool
}

// ActionsConfig drives catalog inspection.
type ActionsConfig struct {
	Generators string
	Package    string
	All        bool
	JSON       bool
}

// notImplemented reports that a command's owning milestone has not landed. It
// names the milestone so an unimplemented phase is never mistaken for a run
// that did nothing.
func notImplemented(command, milestone, summary string) error {
	return fmt.Errorf("%s is not implemented yet (%s: %s)", command, milestone, summary)
}

// parseVars turns repeated --var name=value flags into bindings.
//
// Split on the first = only, because a value may contain one and a name may
// not: a Go module path or a .NET namespace is a plausible value, and nothing
// about a variable's name is Sedum's to interpret beyond where it ends.
//
// A repeated name is an error rather than last-one-wins. Two values for one
// variable on a single command line means the author meant one of them, and
// there is no way to tell which.
func parseVars(vars []string) (map[string]string, error) {
	if len(vars) == 0 {
		return nil, nil
	}
	out := make(map[string]string, len(vars))
	for _, v := range vars {
		name, value, ok := strings.Cut(v, "=")
		if !ok || name == "" {
			return nil, fmt.Errorf("--var %q is not name=value", v)
		}
		if existing, dup := out[name]; dup {
			return nil, fmt.Errorf("--var %s given twice, as %q and %q; one of them was meant and nothing says which",
				name, existing, value)
		}
		out[name] = value
	}
	return out, nil
}
