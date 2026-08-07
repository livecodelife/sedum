package cli

import (
	"errors"
	"fmt"
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
