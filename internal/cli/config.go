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

	// Force lets a record whose selection comes back empty remove every
	// region it owns (prov-2026-375f5d4a), and re-selects every record whatever
	// its hashes say (prov-2026-9621daed).
	Force bool

	// Check classifies every record from its markers and reports, with no
	// model call and nothing written (prov-2026-9621daed).
	Check bool

	// PruneOrphans removes regions whose record was removed or deprecated, in
	// files some current record authorizes (prov-2026-cc73cf61).
	PruneOrphans bool

	// ListOrphans prints every orphan region in place of the summary
	// (prov-2026-3f83335f).
	ListOrphans bool

	// Stamp writes a missing gen or pkg onto markers whose record exists, and
	// nothing else (prov-2026-cc73cf61).
	Stamp bool

	// JSON prints the run's report as the sedum.grow/1 document in place of
	// the text report (prov-2026-bf01cd6b).
	JSON bool

	Retries int

	// Repair bounds how often a selection that passed Phase 5 is sent back
	// because the project's verification failed it, and Verify adds commands
	// to the ones in .sedum/config.yml (prov-2026-b5c21813).
	Repair int
	Verify []string

	// NoVerify skips verification for the run, whatever .sedum/config.yml names.
	NoVerify bool
	Model    string
	LogPath  string
	Verbose  bool

	// LocalModel is a gguf file to serve locally instead of consulting
	// --model against a hosted or already-running endpoint. Set, it starts
	// LocalModelServer (or the goinfer-serve found on PATH) for the run's
	// lifetime rather than reading OPENAI_BASE_URL/OPENAI_API_KEY.
	LocalModel string

	// LocalModelServer is the goinfer-serve binary. Empty resolves it from
	// PATH. Only meaningful with LocalModel set.
	LocalModelServer string

	// LocalModelBackend is passed through as goinfer-serve's own -backend
	// flag (cpu, metal, cuda, webgpu). Applies to LocalModel when set, and
	// to the bundled default model when neither LocalModel nor Model is.
	LocalModelBackend string

	// FullCatalog disables the authorized-path filter on Phase 4's prompt,
	// reproducing every version before it: the model sees every exposed
	// action across the record's packages, whether or not any authorized
	// path could ever match it.
	FullCatalog bool

	// Prompt and PromptFile are caller text for Phase 4's supplementary
	// section (prov-2026-e4bab53b). PromptFile is a path, or "-" for stdin.
	// They are mutually exclusive and are read by the command, not here.
	Prompt     string
	PromptFile string

	// PriorRecording is an earlier recording whose selections for this run's
	// records are shown to the model as context (prov-2026-8602b702). It is
	// not a replay.
	PriorRecording string
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

	if c.NoVerify && len(c.Verify) > 0 {
		return errors.New("--no-verify cannot be combined with --verify: one skips verification and the other adds commands to it")
	}

	if c.Repair < 0 || c.Repair > maxRepair {
		return fmt.Errorf("--repair %d is out of range: use 0 (verify once, no repair) to %d", c.Repair, maxRepair)
	}

	if c.LocalModel == "" && c.LocalModelServer != "" {
		return errors.New("--local-model-server requires --local-model")
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
	if c.LocalModel != "" {
		ignored = append(ignored, "--local-model (replay invokes no model)")
	}
	if c.LocalModel == "" && c.LocalModelBackend != "" {
		ignored = append(ignored, "--local-model-backend (replay invokes no model)")
	}
	if c.Repair != defaultRepair {
		ignored = append(ignored, "--repair (replay invokes no model, so nothing is sent back)")
	}
	if len(c.Verify) > 0 {
		ignored = append(ignored, "--verify (replay invokes no model, so nothing is verified)")
	}
	if c.NoVerify {
		ignored = append(ignored, "--no-verify (replay verifies nothing)")
	}
	if c.Retries != defaultRetries {
		ignored = append(ignored, "--retries (replay validation is terminal, never re-prompted)")
	}
	if c.FullCatalog {
		ignored = append(ignored, "--full-catalog (replay invokes no model, so nothing builds a prompt)")
	}
	if c.Prompt != "" {
		ignored = append(ignored, "--prompt (replay invokes no model, so nothing builds a prompt)")
	}
	if c.PromptFile != "" {
		ignored = append(ignored, "--prompt-file (replay invokes no model, so nothing builds a prompt)")
	}
	if c.PriorRecording != "" {
		ignored = append(ignored, "--prior-recording (replay invokes no model, so nothing builds a prompt)")
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

// RenderConfig drives target rendering for one invocation.
type RenderConfig struct {
	Generators string
	Package    string
	Action     string
	Kwargs     string
	JSON       bool
}

// ActionsConfig drives catalog inspection.
type ActionsConfig struct {
	Generators string
	Package    string
	All        bool
	JSON       bool
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
