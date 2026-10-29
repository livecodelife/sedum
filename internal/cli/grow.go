package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/pipeline"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/runlog"
	"github.com/livecodelife/sedum/internal/selection"
	"github.com/livecodelife/sedum/internal/verify"
)

func newGrowCommand() *cobra.Command {
	var cfg GrowConfig

	// responseSchema is kept local rather than added to GrowConfig, since
	// GrowConfig's file is outside this record's affected_scope: the flag is
	// threaded through runGrow's own signature instead of the config struct
	// every other flag here goes through.
	var responseSchema bool

	cmd := &cobra.Command{
		Use:   "grow",
		Short: "Run the full pipeline: load, ingest, resolve, create, invoke, validate, expand, inject",
		Long: `Runs every phase in order: load and validate generator packages, ingest
provenance records, resolve paths to packages, create files from templates,
invoke the model, validate its output, expand and resolve, and inject.

Nothing is created that a provenance record did not authorize.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cfg.Validate(); err != nil {
				return err
			}
			for _, ignored := range cfg.IgnoredFlags() {
				fmt.Fprintf(cmd.ErrOrStderr(), "sedum: ignoring %s\n", ignored)
			}
			// Replay invokes no model, so nothing here would ever build a
			// schema or a request to carry it on - the same reason --model
			// and --full-catalog are ignored under --execute.
			if cfg.Replaying() && responseSchema {
				fmt.Fprintf(cmd.ErrOrStderr(), "sedum: ignoring --response-schema (replay invokes no model)\n")
			}
			// The caller's context is read and checked here, before any phase
			// and before a model client is built, so a bad file or an
			// oversized text costs nothing. Replay consults neither, so under
			// --execute they are not read and cannot fail the run.
			var ctxIn promptContext
			if !cfg.Replaying() {
				var err error
				if ctxIn, err = readPromptContext(cmd.InOrStdin(), cfg); err != nil {
					return err
				}
			}
			return runGrow(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), cfg, responseSchema, ctxIn)
		},
	}

	f := cmd.Flags()
	f.StringVar(&cfg.Generators, "generators", "", "Generators directory. Required.")
	f.StringVar(&cfg.Records, "records", "", "Provenance records directory. Required unless --execute is given.")
	f.StringVar(&cfg.Output, "output", defaultOutput, "Output directory.")
	f.StringSliceVar(&cfg.Lang, "lang", nil, "Prefer the named package where an extension is contested. Repeatable.")
	f.StringSliceVar(&cfg.Only, "only", nil, "Generate only the named provenance record. Repeatable.")
	f.StringVar(&cfg.RecordTo, "record", "", "Write one recording per record into this directory, default .sedum/recordings under --output; "+
		"a path ending in .json instead writes one recording of the whole run to that file and re-selects every record.")
	f.StringArrayVar(&cfg.Vars, "var", nil, "Bind a variable a generator package declares, as name=value. Repeatable.")
	f.StringVar(&cfg.Execute, "execute", "", "Replay a recording. Skips model invocation.")
	f.BoolVar(&cfg.DryRun, "dry-run", false, "Run every phase, write nothing.")
	f.BoolVar(&cfg.Force, "force", false,
		"Re-select every record even when its markers say it is current, and let a record whose selection comes back empty "+
			"remove every region it owns (without --force that is an error). It overwrites a region a re-render refused because its text was edited by hand. "+
			"Re-rendering replays the kwargs of a v1 marker, or of the recording entry of a v2 region, verbatim, so a changed kwarg default or --var value needs --force. Ignored with --execute.")
	f.BoolVar(&cfg.PruneOrphans, "prune-orphans", false,
		"Remove regions whose record was removed from the records directory or deprecated, in files some current record "+
			"authorizes. Orphans in other files are reported and left alone. Ignored with --execute.")
	f.BoolVar(&cfg.RestampGen, "restamp-gen", false,
		"With --execute and --records: stamp each recorded record's current gen onto its regions in place of the gen the recording "+
			"carries, so a replay after a record edit reads as current. The recorded selection must still pass replay's own checks. "+
			"Asserts the selection still satisfies the edited record.")
	f.BoolVar(&cfg.ListOrphans, "orphans", false,
		"Print every orphan region, one line each, in place of the per-fate summary (counts and the busiest directories). "+
			"The --json report always lists them all.")
	f.BoolVar(&cfg.Stamp, "stamp", false,
		"Write the current gen and pkg onto markers missing them, where the owning record exists, and change nothing else. "+
			"No model call. Asserts that the code is current, so a migrated repository's next grow skips instead of re-selecting. "+
			"Cannot be combined with --check, --execute, --record or --force; with --dry-run it reports what it would stamp.")
	f.BoolVar(&cfg.JSON, "json", false,
		"Print the run's report as the sedum.grow/1 JSON document (conformance/formats) instead of text. "+
			"Each record says whether its regions changed, so a caller runs only the specs whose code moved.")
	f.BoolVar(&cfg.Check, "check", false,
		"Report which records are current, stale, need re-rendering, edited, or have no regions yet (unknown), from their markers and recordings. "+
			"A region whose text no longer matches its body hash is reported as edited and fails the check; v1 regions have no body hash and are reported unverifiable without failing. "+
			"No model call, nothing written. Exits non-zero when any record is stale, edited or needs re-rendering; unknown does not fail. "+
			"Re-rendering replays the kwargs of a v1 marker, or of the recording entry of a v2 region, verbatim, so a changed kwarg default or --var value needs --force. "+
			"Cannot be combined with --execute, --dry-run, --force, --stop-after or --record with a path ending in .json.")
	f.StringVar(&cfg.StopAfter, "stop-after", "", "Halt after the named phase. One of: "+stopPointNames()+".")
	f.IntVar(&cfg.Retries, "retries", defaultRetries, "Model output validation retry limit. Ignored with --execute.")
	f.IntVar(&cfg.Repair, "repair", defaultRepair,
		"How many times a selection that passed validation is sent back, with the output, because a verification command failed it "+
			"(0 verifies once and fails the run, at most 3). Commands come from .sedum/config.yml under the output directory and --verify, run in a scratch copy of the tree. Ignored with --execute.")
	f.BoolVar(&cfg.NoVerify, "no-verify", false,
		"Skip verification for this run, whatever .sedum/config.yml names. Cannot be combined with --verify. Ignored with --execute.")
	f.StringArrayVar(&cfg.Verify, "verify", nil,
		"A shell command run in a scratch copy of the tree with a record's selection applied; repeatable, added after .sedum/config.yml's. "+
			"It must exit 0. Ignored with --execute.")
	f.StringVar(&cfg.Model, "model", "", "Model identifier. Endpoint and credentials come from environment. Ignored with --execute.")
	f.StringVar(&cfg.LocalModel, "local-model", "", "Path to a gguf file to serve locally instead of --model. Starts and stops its own server for the run. Ignored with --execute.")
	f.StringVar(&cfg.LocalModelServer, "local-model-server", "", "Path to the goinfer-serve binary. Defaults to one found on PATH. Requires --local-model.")
	f.StringVar(&cfg.LocalModelBackend, "local-model-backend", "", "Compute backend passed to goinfer-serve's own -backend flag (cpu, metal, cuda, webgpu). Requires --local-model.")
	f.StringVar(&cfg.LogPath, "log", defaultLogPath, "Run log location.")
	f.BoolVarP(&cfg.Verbose, "verbose", "v", false, "Mirror the run log to stdout.")
	f.BoolVar(&cfg.FullCatalog, "full-catalog", false,
		"Show the model every exposed action, including ones no authorized path can reach. Ignored with --execute.")
	f.BoolVar(&responseSchema, "response-schema", false,
		"Compile the record's filtered catalog into a JSON Schema and send it as response_format, "+
			"constraining the model's completion. Off by default, and off, the request is byte-for-byte "+
			"what it is today. Fails clearly if the endpoint rejects response_format; never falls back "+
			"unconstrained. Ignored with --execute.")

	f.StringVar(&cfg.Prompt, "prompt", "",
		"Supplementary text added to every selected record's prompt, after its intent and constraints. "+
			"Advisory: it cannot widen what a record authorizes. At most 16384 bytes. Ignored with --execute.")
	f.StringVar(&cfg.PromptFile, "prompt-file", "",
		"Read the supplementary text from a file, or from stdin with -. Same rules as --prompt. Ignored with --execute.")
	f.StringVar(&cfg.PriorRecording, "prior-recording", "",
		"A recording whose selections for this run's records are shown to the model as context. "+
			"Not a replay: the model is still called and its output validated as ever. Ignored with --execute.")

	mustMarkRequired(cmd, "generators")

	// A recording is either being written or being replayed, never both.
	cmd.MarkFlagsMutuallyExclusive("record", "execute")

	// A run consults exactly one model source.
	cmd.MarkFlagsMutuallyExclusive("model", "local-model")

	// Text comes from one place.
	cmd.MarkFlagsMutuallyExclusive("prompt", "prompt-file")

	return cmd
}

// runGrow runs the pipeline.
//
// The model client is built before any phase runs, so a run that cannot reach a
// model fails while nothing has been written. Creating files and then failing
// at the model call would leave the output tree half built for a reason the
// user could have been told first.
func runGrow(ctx context.Context, out, errOut io.Writer, cfg GrowConfig, responseSchema bool, pc promptContext) error {
	// No --stop-after means run every phase. Zero is not a stop point, so it
	// reads as "stop after nothing".
	stopAfter := 0
	name := ""
	if cfg.StopAfter != "" {
		sp, ok := lookupStopPoint(cfg.StopAfter)
		if !ok {
			return fmt.Errorf("--stop-after %q is not a phase boundary; expected one of: %s", cfg.StopAfter, stopPointNames())
		}
		stopAfter, name = sp.afterPhase, sp.name
	}

	if err := checkCombinations(cfg); err != nil {
		return err
	}

	log, err := runlog.New(cfg.LogPath, cfg.Verbose)
	if err != nil {
		return err
	}
	defer log.Close()

	if cfg.Replaying() {
		return runReplay(out, errOut, cfg, stopAfter, name, log)
	}

	// A run halting before Phase 4 never consults a model, so it never needs
	// one configured. That is what makes --stop-after resolution and
	// --stop-after files usable with no endpoint at all.
	//
	// The client is built before any phase runs. A run that cannot reach a
	// model fails while nothing has been written, rather than creating files
	// and then failing at the call for a reason the user could have been told
	// first.
	var client selection.Client
	if !cfg.Check && !cfg.Stamp && (stopAfter == 0 || stopAfter >= pipeline.PhaseSelect) {
		var localCfg selection.LocalConfig
		var useLocal bool

		switch {
		case cfg.LocalModel != "":
			localCfg = selection.LocalConfig{
				ModelPath:  cfg.LocalModel,
				ServerPath: cfg.LocalModelServer,
				Backend:    cfg.LocalModelBackend,
			}
			useLocal = true
		case cfg.Model == "":
			// Neither --model nor --local-model: the bundled default, if
			// this build carries one, rather than an immediate error.
			bundled, ok := selection.BundledLocalConfig()
			if !ok {
				return errors.New(
					"no --model or --local-model given, and no bundled model found next to this binary; " +
						"pass one of them, or install the full archive with: " +
						"curl -fsSL https://raw.githubusercontent.com/livecodelife/sedum/main/install.sh | sh")
			}
			bundled.Backend = cfg.LocalModelBackend
			localCfg = bundled
			useLocal = true
		}

		if useLocal {
			local, err := selection.NewLocal(ctx, localCfg)
			if err != nil {
				return err
			}
			// The subprocess belongs to this run; nothing past runGrow needs
			// it, so it stops here rather than outliving the pipeline result
			// that already consumed its answers.
			defer local.Close()
			client = local
		} else {
			built, err := selection.NewOpenAI(cfg.Model)
			if err != nil {
				return err
			}
			client = built
		}

		// --response-schema opts this run's client into structured output.
		// The wrapper is recognized inside selection.Select, which is the one
		// place that knows each record's own (authorized-path-filtered)
		// catalog to compile a schema from; wrapping here rather than passing
		// a bool through pipeline.Config means this record's flag reaches
		// Select without pipeline.go - which only ever forwards a
		// selection.Client unchanged - needing a change of its own.
		//
		// *Local (the bundled goinfer-serve backend) implements no
		// structured-output method today, so wrapping it changes nothing
		// about whether it can honor --response-schema: Select's own type
		// assertion against selection.SchemaCapable still fails for it, and
		// the run still fails clearly rather than silently running
		// unconstrained. Whether goinfer-serve could honor response_format,
		// or would reproduce prov-2026-4bcabb2f's collapse if it did, is
		// unverified and this record does not resolve it.
		if responseSchema {
			client = selection.WithResponseSchema(client)
		}
	}

	variables, err := parseVars(cfg.Vars)
	if err != nil {
		return err
	}

	verifyFn, repair, err := verifier(errOut, cfg)
	if err != nil {
		return err
	}

	started := time.Now()
	result, err := pipeline.Run(ctx, pipeline.Config{
		Generators: cfg.Generators,
		Records:    cfg.Records,
		Output:     cfg.Output,
		Lang:       cfg.Lang,
		Only:       cfg.Only,
		DryRun:     cfg.DryRun,
		Force:      cfg.Force,
		// Only a --record path ending in .json re-selects every record: one
		// recording of the whole run holds a selection for each, and a record
		// skipped or re-rendered from its markers makes none
		// (prov-2026-9621daed). A directory just says where recordings go.
		Reselect:       recordsWholeRun(cfg.RecordTo),
		Recordings:     recordingsDir(cfg.RecordTo),
		Check:          cfg.Check,
		PruneOrphans:   cfg.PruneOrphans,
		Stamp:          cfg.Stamp,
		StopAfterPhase: stopAfter,
		Client:         client,
		Retries:        cfg.Retries,
		Repair:         repair,
		Verify:         verifyFn,
		Variables:      variables,
		Log:            log,
		FullCatalog:    cfg.FullCatalog,
		Prompt:         pc.text,
		Prior:          pc.prior,
	})
	writeGrowLog(errOut, cfg, result, err, started)
	if err != nil {
		return err
	}

	printWarnings(errOut, result.Warnings)
	if cfg.JSON {
		if err := writeJSON(out, buildGrowReport(result, cfg.DryRun, false, cfg.Check)); err != nil {
			return err
		}
		if cfg.Check && result.CheckFailed() {
			return errors.New("check failed: some records are stale, edited, need re-rendering, or have left orphan regions")
		}
	} else {
		if cfg.Check {
			printClassifications(out, result)
			printOrphans(out, result, cfg.ListOrphans)
			if result.CheckFailed() {
				return errors.New("check failed: some records are stale, edited, need re-rendering, or have left orphan regions; run sedum grow (with --prune-orphans for orphans)")
			}
			return nil
		}
		if cfg.Stamp {
			printStamped(out, result, cfg.DryRun)
			return nil
		}
		if recordsWholeRun(cfg.RecordTo) {
			fmt.Fprintln(out, "--record: every record re-selected, whatever its markers say, so the recording holds a selection for each")
		}
		report(out, result, stopAfter, cfg.DryRun, cfg.ListOrphans)
		if name != "" {
			fmt.Fprintf(out, "\nstopped after %s\n", name)
		}
	}
	if cfg.Check || cfg.Stamp {
		return nil
	}

	// Written after the report, so a run whose recording cannot be written
	// still tells the user what it decided. The model call is already paid
	// for and its output is the thing at risk of being lost. The per-record
	// recordings are Sedum's, not Run's: a second grow skips a record whose
	// regions are current only when this one wrote its recording.
	if !cfg.DryRun && (name == "" || (cfg.RecordTo != "" && !recordsWholeRun(cfg.RecordTo))) {
		dir := recordingsDir(cfg.RecordTo)
		if dir == "" {
			dir = filepath.Join(cfg.Output, ".sedum", "recordings")
		}
		if err := pipeline.WriteRecordings(result, dir); err != nil {
			return err
		}
	}
	if recordsWholeRun(cfg.RecordTo) {
		if err := writeRecording(cfg.RecordTo, capture(result, cfg)); err != nil {
			return err
		}
		if !cfg.JSON {
			fmt.Fprintf(out, "\nrecording written to %s\n", cfg.RecordTo)
		}
	}

	return refusedError(result)
}

// recordsWholeRun reports whether --record names a file for one recording of
// the whole run, which is a path ending in .json, rather than a directory.
func recordsWholeRun(record string) bool {
	return strings.HasSuffix(record, ".json")
}

// recordingsDir is the directory --record names for the per-record recordings:
// the value when it is set and does not end in .json, and empty otherwise.
func recordingsDir(record string) string {
	if record == "" || recordsWholeRun(record) {
		return ""
	}
	return record
}

// refusedError is the error a run returns when a re-render was refused over a
// hand edit: the first line, then each refused region on a line of its own.
func refusedError(result *pipeline.Result) error {
	var lines []string
	for _, c := range result.Classifications {
		if !c.Refused() {
			continue
		}
		for _, e := range c.Edited {
			for _, r := range c.Rerender {
				if e.Path == r.Path && e.Action == r.Action && e.Variant == r.Variant {
					lines = append(lines, fmt.Sprintf("  %s  %s", e.Path, describeRef(e)))
					break
				}
			}
		}
	}
	if len(lines) == 0 {
		return nil
	}
	return errors.New("refusing to re-render over hand edits; --force overwrites them\n" + strings.Join(lines, "\n"))
}

// maxPromptBytes bounds the caller's text. Over it is an error, not a
// truncation: the tail of a failure log is usually the part that says what
// went wrong.
const maxPromptBytes = 16384

// promptContext is the caller's context for Phase 4, already read and checked.
type promptContext struct {
	text  string
	prior *recording.Recording
}

// readPromptContext resolves --prompt / --prompt-file and --prior-recording.
//
// Blank text is no text. --prior-recording must hold at least one record: a
// retry run that silently lacked its context would look like a normal run. No
// version check is made on the prior recording: the recording carries no schema
// version and sedum_version is not consulted for compatibility anywhere.
func readPromptContext(stdin io.Reader, cfg GrowConfig) (promptContext, error) {
	var pc promptContext

	text := cfg.Prompt
	if cfg.PromptFile != "" {
		var data []byte
		var err error
		if cfg.PromptFile == "-" {
			data, err = io.ReadAll(stdin)
			if err != nil {
				return pc, fmt.Errorf("reading --prompt-file from stdin: %w", err)
			}
		} else {
			data, err = os.ReadFile(cfg.PromptFile)
			if err != nil {
				return pc, fmt.Errorf("reading --prompt-file %s: %w", cfg.PromptFile, err)
			}
		}
		text = string(data)
	}
	if len(text) > maxPromptBytes {
		return pc, fmt.Errorf("the supplementary prompt is %d bytes, over the %d-byte limit; nothing is truncated, so shorten it",
			len(text), maxPromptBytes)
	}
	if strings.TrimSpace(text) != "" {
		pc.text = text
	}

	if cfg.PriorRecording != "" {
		rec, err := readRecording(cfg.PriorRecording)
		if err != nil {
			return pc, fmt.Errorf("--prior-recording %s: %w", cfg.PriorRecording, err)
		}
		if len(rec.Records) == 0 {
			return pc, fmt.Errorf("--prior-recording %s holds no records, so it has no selections to show", cfg.PriorRecording)
		}
		pc.prior = &rec
	}
	return pc, nil
}

// runReplay executes a recording instead of consulting a model.
//
// Replay enters the pipeline at Phase 3 with resolution already decided and
// skips Phase 4, so nothing here builds a model client: the flags that
// configure one are reported as ignored rather than honored.
func runReplay(out, errOut io.Writer, cfg GrowConfig, stopAfter int, name string, log *runlog.Log) error {
	rec, err := readRecording(cfg.Execute)
	if err != nil {
		return err
	}

	variables, err := parseVars(cfg.Vars)
	if err != nil {
		return err
	}

	result, err := pipeline.Replay(rec, pipeline.ReplayConfig{
		Generators: cfg.Generators,
		Output:     cfg.Output,
		// Optional under --execute. Supplied, every recorded path is checked
		// against affected_scope and forbidden_scope; omitted, the recording
		// executes as written.
		Records:    cfg.Records,
		Variables:  variables,
		RestampGen: cfg.RestampGen,
		// --check is a dry replay that fails when it would change anything.
		DryRun:         cfg.DryRun || cfg.Check,
		StopAfterPhase: stopAfter,
		Log:            log,
	})
	if err != nil {
		return err
	}

	printWarnings(errOut, result.Warnings)
	if cfg.JSON {
		if err := writeJSON(out, buildGrowReport(result, cfg.DryRun, true, false)); err != nil {
			return err
		}
	} else {
		for _, r := range result.Restamps {
			fmt.Fprintln(out, r)
		}
		report(out, result, stopAfter, cfg.DryRun || cfg.Check, cfg.ListOrphans)
		if name != "" {
			fmt.Fprintf(out, "\nstopped after %s\n", name)
		}
	}
	if cfg.Check && pipeline.ReplayWouldChange(result.Injections) {
		return fmt.Errorf("check failed: replaying %s would change the tree (stale gen, a changed region, or a region the recording no longer selects)", cfg.Execute)
	}
	// The recordings of the records replayed are written after the report, so
	// a later grow finds what this replay wrote. A replay stopped at a phase
	// boundary injected nothing, so it leaves the recordings as they are.
	if !cfg.DryRun && !cfg.Check && stopAfter == 0 {
		return pipeline.WriteReplayRecordings(rec, result, filepath.Join(cfg.Output, ".sedum", "recordings"))
	}
	return nil
}

// capture is the recording a run writes. A run that re-selected every record
// (no --only) writes a complete one, which replay may reconcile against
// (prov-2026-c637c262).
func capture(result *pipeline.Result, cfg GrowConfig) recording.Recording {
	rec := pipeline.Capture(result)
	rec.Complete = len(cfg.Only) == 0
	return rec
}

// readRecording loads a recording, naming the file in anything that goes wrong
// with it. The recording is an input format, so a caller authoring one against
// these diagnostics has to be able to tell which file was rejected.
func readRecording(path string) (recording.Recording, error) {
	f, err := os.Open(path)
	if err != nil {
		return recording.Recording{}, fmt.Errorf("opening the recording: %w", err)
	}
	defer f.Close()

	rec, err := recording.Read(f)
	if err != nil {
		return recording.Recording{}, fmt.Errorf("%s: %w", path, err)
	}
	return rec, nil
}

// writeRecording writes a recording, creating the directories above it.
func writeRecording(path string, rec recording.Recording) error {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating the directory for the recording: %w", err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating the recording: %w", err)
	}
	defer f.Close()

	if err := recording.Write(f, rec); err != nil {
		return err
	}
	return f.Close()
}

// report prints what the run got as far as, at the detail the stop point makes
// useful. A run halted at resolution has decided nothing to show but
// resolutions; one that finished has files, selections, and injections, and the
// interesting part is the last of those.
func report(out io.Writer, result *pipeline.Result, stopAfter int, dryRun, listOrphans bool) {
	if stopAfter == pipeline.PhaseResolve {
		printResolutions(out, result, nil, false)
		return
	}

	printFiles(out, result, dryRun)
	if stopAfter == pipeline.PhaseCreate {
		return
	}

	printStaleSeeded(out, result)
	printSelections(out, result)
	if stopAfter == pipeline.PhaseValidate || stopAfter == pipeline.PhaseExpand {
		return
	}
	printInjections(out, result, dryRun)
	printOrphans(out, result, listOrphans)
}

// printSelections reports what the model chose, grouped by record because a
// record is one model call and the unit a reader is reasoning about.
func printSelections(out io.Writer, result *pipeline.Result) {
	fmt.Fprintln(out)
	var total int
	modes := map[string]pipeline.Mode{}
	for _, c := range result.Classifications {
		modes[c.RecordID] = c.Mode
	}
	for _, s := range result.Selections {
		fmt.Fprintf(out, "%s\n", s.RecordID)
		switch modes[s.RecordID] {
		case pipeline.ModeSkip:
			fmt.Fprintf(out, "  (current; no model call)\n")
			continue
		case pipeline.ModeRerender:
			fmt.Fprintf(out, "  (package changed; regions re-rendered from their markers, no model call)\n")
			continue
		case pipeline.ModeEdited:
			fmt.Fprintf(out, "  (edited by hand; no model call)\n")
			continue
		}
		if len(s.Invocations) == 0 {
			fmt.Fprintf(out, "  (no action selected)\n")
			continue
		}
		for _, inv := range s.Invocations {
			fmt.Fprintf(out, "  %s %s\n", inv.Action, describeKwargs(inv.Kwargs))
			total++
		}
		if s.Repairs > 0 {
			fmt.Fprintf(out, "  repaired after %d attempt(s)\n", s.Repairs)
		}
	}
	fmt.Fprintf(out, "\n%d invocation(s) selected across %d record(s)\n", total, len(result.Selections))
}

func printInjections(out io.Writer, result *pipeline.Result, dryRun bool) {
	verb := "injected"
	if dryRun {
		verb = "would inject"
	}

	fmt.Fprintln(out)
	var replaced, skipped, pruned, held, unselected int
	for _, r := range result.Injections {
		switch {
		case r.Held:
			held++
			fmt.Fprintf(out, "held     %s  %s (record %s selected nothing; --force would remove it)\n", r.Path, describeRegion(r), r.Record)
		case r.Unselected:
			unselected++
			fmt.Fprintf(out, "unselected %s  %s (record %s; the recording is partial, so replay leaves it)\n", r.Path, describeRegion(r), r.Record)
		case r.Orphan:
			// Reported once, by printOrphans, which says why: its record is
			// gone or deprecated, not that it stopped selecting the region.
			continue
		case r.Pruned:
			pruned++
			if dryRun {
				fmt.Fprintf(out, "would prune %s  %s (record %s no longer selects it)\n", r.Path, describeRegion(r), r.Record)
			} else {
				fmt.Fprintf(out, "pruned   %s  %s (record %s no longer selects it)\n", r.Path, describeRegion(r), r.Record)
			}
		case r.Skipped:
			skipped++
			fmt.Fprintf(out, "seeded   %s  %s (left as it is)\n", r.Path, describeRegion(r))
		case r.Replaced:
			replaced++
			fmt.Fprintf(out, "replaced %s  %s\n", r.Path, describeRegion(r))
		default:
			fmt.Fprintf(out, "%s %s  %s\n", verb, r.Path, describeRegion(r))
		}
	}
	orphans := 0
	for _, r := range result.Injections {
		if r.Orphan {
			orphans++
		}
	}
	fmt.Fprintf(out, "\n%d region(s) %s, %d replaced, %d left as seeded\n",
		len(result.Injections)-pruned-held-orphans-unselected, verb, replaced, skipped)
	if unselected > 0 {
		fmt.Fprintf(out, "%d region(s) unselected: owned by a record this partial recording names, and left in place\n", unselected)
	}
	if pruned > 0 || held > 0 {
		fmt.Fprintf(out, "%d region(s) pruned, %d held for --force\n", pruned, held)
	}
}

func describeRegion(r inject.Result) string { return label(r.Action, r.Variant) }

// label is action, or action:variant for a discriminated one.
func label(action, variant string) string {
	if variant == "" {
		return action
	}
	return action + ":" + variant
}

// describeKwargs renders bound arguments in name order, so that one run's
// report reads the same way as another's.
func describeKwargs(kwargs map[string]any) string {
	names := make([]string, 0, len(kwargs))
	for name := range kwargs {
		names = append(names, name)
	}
	sort.Strings(names)

	pairs := make([]string, 0, len(names))
	for _, name := range names {
		pairs = append(pairs, fmt.Sprintf("%s=%v", name, kwargs[name]))
	}
	return strings.Join(pairs, " ")
}

// printFiles reports what Phase 3 did, in three groups.
//
// Files it created, files it found already there, and files a package declares
// it does not write. Create-if-absent makes the second group normal rather than
// exceptional, so it is reported rather than hidden. The third is the run's
// handoff - authorized work something other than Sedum has to do - and is
// gathered at the end rather than scattered through the listing, because a
// person acting on it wants it as a list.
//
// An unmanaged path resolved to no package, by design: the check runs before
// extension resolution so a path no extension can reach is still reportable.
// That is why it is not merely a different line here but a different branch -
// there is no package name to print.
func printFiles(out io.Writer, result *pipeline.Result, dryRun bool) {
	verb := "created"
	if dryRun {
		verb = "would create"
	}

	var created, existing, unmanaged int
	for _, f := range result.Files {
		switch {
		case f.Unmanaged:
			unmanaged++
		case f.Existed:
			existing++
		default:
			created++
			fmt.Fprintf(out, "%s  %s  (%s %s)\n", verb, f.Path, f.Package.Name, describeTemplate(f.Resolution))
		}
	}
	for _, f := range result.Files {
		if f.Existed && !f.Unmanaged {
			fmt.Fprintf(out, "exists   %s  (left as it is)\n", f.Path)
		}
	}

	fmt.Fprintf(out, "\n%d path(s) authorized, %d %s, %d already present, %d left unmanaged\n",
		len(result.Files), created, verb, existing, unmanaged)

	if unmanaged == 0 {
		return
	}
	fmt.Fprintf(out, "\n%d path(s) left unmanaged, for a person or another tool:\n", unmanaged)
	for _, f := range result.Files {
		if f.Unmanaged {
			fmt.Fprintf(out, "  %s (%s declares %q)\n", f.Path, f.UnmanagedBy, f.UnmanagedAs)
		}
	}
}

// checkCombinations rejects flag combinations --check cannot honor. --check
// writes nothing and calls no model, so anything that writes, records, replays
// or forces a selection contradicts it (prov-2026-9621daed).
func checkCombinations(cfg GrowConfig) error {
	if cfg.Stamp {
		// --stamp asserts code is current and writes only hashes, so nothing
		// that selects, replays, records or checks belongs with it
		// (prov-2026-cc73cf61).
		for _, c := range []struct {
			set  bool
			flag string
		}{
			{cfg.Check, "--check"},
			{cfg.Execute != "", "--execute"},
			{cfg.RecordTo != "", "--record"},
			{cfg.Force, "--force"},
		} {
			if c.set {
				return fmt.Errorf("--stamp cannot be combined with %s: --stamp only writes missing hashes onto markers, with no model call", c.flag)
			}
		}
	}
	if !cfg.Check {
		return nil
	}
	for _, c := range []struct {
		set  bool
		flag string
	}{
		{recordsWholeRun(cfg.RecordTo), "--record"},
		{cfg.DryRun, "--dry-run"},
		{cfg.Force, "--force"},
		{cfg.StopAfter != "", "--stop-after"},
	} {
		if c.set {
			return fmt.Errorf("--check cannot be combined with %s: --check only reports, from markers, with no model call and nothing written", c.flag)
		}
	}
	return nil
}

// printClassifications is --check's report: each record's state and the
// regions behind it, then the records with regions that cannot be verified,
// then the regions whose marker version is not known.
func printClassifications(out io.Writer, result *pipeline.Result) {
	const pad = "             "
	var current, stale, rerender, edited, unknown int
	for _, c := range result.Classifications {
		switch c.Mode {
		case pipeline.ModeNone:
			fmt.Fprintf(out, "%-10s %s (nothing Sedum writes)\n", "none", c.RecordID)
		case pipeline.ModeSelect:
			unknown++
			fmt.Fprintf(out, "%-10s %s (no regions yet)\n", "unknown", c.RecordID)
		case pipeline.ModeSkip:
			current++
			fmt.Fprintf(out, "%-10s %s\n", "current", c.RecordID)
		case pipeline.ModeReselect:
			stale++
			fmt.Fprintf(out, "%-10s %s\n", "stale", c.RecordID)
			for _, r := range c.Stale {
				if r.Reason != "" {
					fmt.Fprintf(out, "%s%s  %s  (%s)\n", pad, r.Path, describeRef(r), r.Reason)
					continue
				}
				fmt.Fprintf(out, "%s%s  %s\n", pad, r.Path, describeRef(r))
			}
		case pipeline.ModeRerender:
			rerender++
			fmt.Fprintf(out, "%-10s %s (package changed)\n", "re-render", c.RecordID)
			for _, r := range c.Rerender {
				fmt.Fprintf(out, "%s%s  %s\n", pad, r.Path, describeRef(r))
			}
		case pipeline.ModeEdited:
			edited++
			fmt.Fprintf(out, "%-10s %s\n", "edited", c.RecordID)
			for _, r := range c.Edited {
				fmt.Fprintf(out, "%s%s  %s  (hand-edited: body hash differs)\n", pad, r.Path, describeRef(r))
			}
		}
	}
	for _, c := range result.Classifications {
		if len(c.Unverifiable) == 0 {
			continue
		}
		fmt.Fprintf(out, "%-10s %s\n", "unverifiable", c.RecordID)
		for _, r := range c.Unverifiable {
			fmt.Fprintf(out, "%s%s  %s\n", pad, r.Path, describeRef(r))
		}
	}
	for _, c := range result.Classifications {
		for _, r := range c.Foreign {
			fmt.Fprintf(out, "foreign %s  %s  (marker version not known; left alone)\n", r.Path, describeRef(r))
		}
	}
	printStaleSeeded(out, result)
	fmt.Fprintf(out, "\n%d current, %d stale, %d to re-render, %d edited, %d unknown\n", current, stale, rerender, edited, unknown)
}

// printStaleSeeded lists seeded regions whose record has changed since they
// were written. Sedum never rewrites them; a person should look.
func printStaleSeeded(out io.Writer, result *pipeline.Result) {
	for _, c := range result.Classifications {
		for _, r := range c.StaleSeeded {
			fmt.Fprintf(out, "seeded region %s  %s was written from an older version of %s; Sedum leaves it alone, so look at it\n",
				r.Path, describeRef(r), c.RecordID)
		}
	}
}

func describeRef(r pipeline.RegionRef) string { return label(r.Action, r.Variant) }

// printStamped is --stamp's report.
func printStamped(out io.Writer, result *pipeline.Result, dryRun bool) {
	verb := "stamped"
	if dryRun {
		verb = "would stamp"
	}
	for _, s := range result.Stamped {
		var added []string
		if s.Gen {
			added = append(added, "gen")
		}
		if s.Pkg {
			added = append(added, "pkg")
		}
		ref := describeRef(pipeline.RegionRef{Path: s.Path, Action: s.Action, Variant: s.Variant})
		fmt.Fprintf(out, "%s %s  %s (%s, record %s)\n", verb, s.Path, ref, strings.Join(added, " and "), s.Record)
	}
	fmt.Fprintf(out, "\n%d marker(s) %s\n", len(result.Stamped), verb)
}

// verifier is the project's verification for this run: the commands in
// .sedum/config.yml under the output directory, then --verify, run by
// verify.Run in a scratch copy of the tree. It returns a nil function when
// there is nothing to run or --no-verify is set, which is how a run behaves as
// it always did (prov-2026-b5c21813).
func verifier(out io.Writer, cfg GrowConfig) (func(context.Context, map[string]string) (string, error), int, error) {
	if cfg.NoVerify || cfg.Check || cfg.Stamp {
		return nil, 0, nil
	}
	loaded, err := verify.LoadConfig(cfg.Output)
	if err != nil {
		return nil, 0, err
	}
	loaded = loaded.With(cfg.Verify)
	if !loaded.Enabled() {
		return nil, 0, nil
	}
	loaded.Skip = []string{filepath.Base(cfg.Generators), filepath.Base(cfg.Records)}
	return func(ctx context.Context, candidate map[string]string) (string, error) {
		failures, err := verify.Run(ctx, loaded, cfg.Output, candidate, out)
		if err != nil {
			return "", err
		}
		if len(failures) == 0 {
			return "", nil
		}
		return verify.Message(failures), nil
	}, cfg.Repair, nil
}
