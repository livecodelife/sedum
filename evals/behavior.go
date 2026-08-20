package evals

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/livecodelife/sedum/internal/recording"
)

// Applying a selection and seeing whether the result works.
//
// Everything else the harness measures is about the answer: which actions were
// selected, what their arguments were bound to. This is the only thing that is
// about the application, and it is the only question a count cannot be made to
// answer - a sample can select completely, bind correctly, and still produce a
// service that violates the record it was generated from, because what the
// generator package can express is not something the answer has any say in
// (prov-2026-83340ba0).
//
// The work is done by evals/behavior/behave.sh, which is shell because every
// step in it is a command a person would run. Nothing here knows what a target
// does; it hands over an invocation list and reads back a result file.

// harnessScript is where the runner is, relative to this package's directory.
// The eval's tests and its command both run from evals/, which is what makes a
// relative path the honest one.
const harnessScript = "behavior/behave.sh"

// BehaviorRun is what applying one sample's selection produced.
type BehaviorRun struct {
	// Outcome is one of:
	//
	//   ok             every assertion passed
	//   checks_failed  it built, booted and answered, and disagreed
	//   failed         a phase died, so the assertions never ran
	//
	// Kept apart rather than reduced to a rate, because a service that never
	// booted and one that booted and answered wrongly are different findings.
	// A single number would merge a broken generator package with a wrong one.
	Outcome string `json:"outcome"`

	// FailedPhase names the step that died, for a `failed` outcome.
	FailedPhase string `json:"failed_phase,omitempty"`

	// Detail is the tail of that phase's log.
	//
	// Without it a failure says "build" and stops, and the log that said why
	// went with the temporary project. Three times in one session the answer
	// was reached by re-running the harness by hand against a reconstructed
	// selection - and a reconstruction is not the sample that failed
	// (prov-2026-93829987).
	//
	// A tail rather than the whole log, on the rule Sample.Detail already
	// follows: enough to act on, without the entry becoming a log store. A
	// failure that needs more has --keep, which leaves the project behind.
	Detail string `json:"detail,omitempty"`

	// Checks and Passed are the assertion counts. Both are zero for a run that
	// never reached the verify phase, which is why Outcome is what a reader
	// should branch on rather than the ratio.
	Checks int `json:"checks"`
	Passed int `json:"passed"`

	// Failed names the assertions that did not hold, so a rate is never
	// reported without the ability to say which contract broke. Counted across
	// samples in the report, because one assertion failing in every sample and
	// twenty failing in one are different problems.
	Failed []string `json:"failed,omitempty"`

	// Attribution names the action whose region wrote each line the failed
	// phase's log points at.
	//
	// Absent rather than empty for a run that did not fail, for one whose log
	// named no file the project holds, and for every entry drawn before this
	// field existed - a dash is not a zero, and a reader must not take an old
	// entry for a build no action was responsible for.
	//
	// It is a list because a build naming several files and lines is several
	// findings: five identical-looking `build` deaths were three distinct
	// defects, and the difference was invisible until the lines were attributed
	// (prov-2026-27c10ac4).
	Attribution []Attribution `json:"attribution,omitempty"`

	Elapsed time.Duration `json:"elapsed"`

	// Err is set when the harness could not be run at all - a missing target, a
	// script that would not start. Excluded from every rate, on the same rule
	// that keeps an unreachable endpoint out of the selection denominator.
	Err error `json:"-"`
}

// Working reports whether this run produced an application that did everything
// the target asserts.
func (b BehaviorRun) Working() bool { return b.Err == nil && b.Outcome == "ok" }

// Attribution is one line a failed phase named, and the region that wrote it.
//
// Nothing here parses the generated language. The compiler names a file and a
// line, the marker pair enclosing that line names the action, and a target for
// a new stack inherits both without writing any lookup of its own.
type Attribution struct {
	// File is relative to the generated project, and Line is what the phase's
	// log named.
	File string `json:"file"`
	Line int    `json:"line"`

	// Action and Variant name the region enclosing that line, and are empty
	// when no marker encloses it.
	//
	// Unattributed rather than guessed at. The compiler often names the file
	// template's own text - the first of two declarations is the template's
	// line, not the injected one - and an attribution that reached for the
	// nearest marker would name an action that did not write it.
	Action  string `json:"action,omitempty"`
	Variant string `json:"variant,omitempty"`

	// Record is the provenance record that last parameterized the region, and
	// Kwargs is what it was rendered from. Both come off the marker, which is
	// why attribution needs no state the harness has to maintain.
	Record string         `json:"record,omitempty"`
	Kwargs map[string]any `json:"kwargs,omitempty"`
}

// Attributed reports whether a marker enclosed this line.
func (a Attribution) Attributed() bool { return a.Action != "" }

// Label is the action:variant pair, as the marker carries it.
func (a Attribution) Label() string {
	if a.Variant == "" {
		return a.Action
	}
	return a.Action + ":" + a.Variant
}

// harnessResult is behave.sh's results file, read at the fields this needs.
type harnessResult struct {
	Outcome      string `json:"outcome"`
	FailedPhase  string `json:"failed_phase"`
	Detail       string `json:"detail"`
	ChecksPassed int    `json:"checks_passed"`
	ChecksTotal  int    `json:"checks_total"`
	Checks       []struct {
		Check string `json:"check"`
		Pass  bool   `json:"pass"`
	} `json:"checks"`
	Attribution []Attribution `json:"attribution"`
}

// RunBehavior applies one sample's invocations to a freshly scaffolded
// application and reports whether it worked.
//
// The invocations are the sample's own, never a canned list. A checked-in
// answer would measure what the generator package can render and report it as a
// property of the model, which is the one thing this must not do.
func RunBehavior(ctx context.Context, target string, invocations []recording.Invocation, variables map[string]string) BehaviorRun {
	started := time.Now()

	if len(invocations) == 0 {
		// Nothing to apply. Callers are expected not to ask, so this is a
		// guard rather than a case: a sample with no selection is unmeasured,
		// not failed.
		return BehaviorRun{Err: fmt.Errorf("no invocations to apply"), Elapsed: time.Since(started)}
	}

	answer, err := writeAnswer(invocations)
	if err != nil {
		return BehaviorRun{Err: err, Elapsed: time.Since(started)}
	}
	defer os.Remove(answer)

	// A results directory per run, removed with it. The harness keeps its own
	// results for a person running it by hand; a measurement drawing thirty
	// samples wants none of that accumulating - which runHarness owns for both
	// arms.
	return runHarness(ctx, target, []string{"--answer", answer}, variables, started)
}

// RunBehaviorFiles is the baseline arm's half: sources the model already wrote,
// copied over the prepared scaffold instead of generated into it.
//
// Every phase but generate is the one the sedum arm runs. That is not thrift -
// it is what makes the two comparable, because a behaviour rate that came from
// two different harnesses would be two measurements rather than a comparison
// (prov-2026-a4dbe65c).
func RunBehaviorFiles(ctx context.Context, target string, files map[string]string, variables map[string]string) BehaviorRun {
	started := time.Now()

	if len(files) == 0 {
		return BehaviorRun{Err: fmt.Errorf("no files to apply"), Elapsed: time.Since(started)}
	}

	dir, err := os.MkdirTemp("", "sedum-baseline-")
	if err != nil {
		return BehaviorRun{Err: err, Elapsed: time.Since(started)}
	}
	defer os.RemoveAll(dir)

	for _, rel := range sortedNames(files) {
		full := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return BehaviorRun{Err: err, Elapsed: time.Since(started)}
		}
		if err := os.WriteFile(full, []byte(files[rel]), 0o644); err != nil {
			return BehaviorRun{Err: err, Elapsed: time.Since(started)}
		}
	}

	return runHarness(ctx, target, []string{"--files", dir}, variables, started)
}

// runHarness is everything the two arms share: the results directory, the
// variables, and the rule that the exit status is not the measurement.
func runHarness(ctx context.Context, target string, arm []string, variables map[string]string, started time.Time) BehaviorRun {
	results, err := os.MkdirTemp("", "sedum-behave-results-")
	if err != nil {
		return BehaviorRun{Err: err, Elapsed: time.Since(started)}
	}
	defer os.RemoveAll(results)

	deadlineCtx, cancel := context.WithTimeout(ctx, behaviorDeadline())
	defer cancel()

	cmd := exec.CommandContext(deadlineCtx, "bash", append([]string{harnessScript, target}, arm...)...)

	// Its own process group, so the deadline can reach what the script is
	// waiting on. bash defers a signal while a foreground command runs, and
	// behave.sh runs every phase in the foreground on purpose - so signalling
	// the script alone is queued until the command returns, which for a hung
	// bundle install or a server that never listens is never. Signalling the
	// group kills the command too, which returns the script to its own control
	// and lets its EXIT trap remove the project, the server and the database
	// (prov-2026-3957eed2).
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
	}
	// If the trap does not finish, stop waiting on it. A deadline that hung
	// waiting for its own cleanup would be the thing it exists to prevent.
	cmd.WaitDelay = 20 * time.Second

	cmd.Env = append(os.Environ(), "RESULTS_DIR="+results)
	// The case's variables reach the target as environment, so the scaffold and
	// the generation read one value rather than two that agree by convention:
	// go mod init and --var are the same module name or the service does not
	// compile (prov-2026-1c33a50b).
	for _, name := range sortedNames(variables) {
		cmd.Env = append(cmd.Env, "SEDUM_VAR_"+strings.ToUpper(name)+"="+variables[name])
	}
	out, runErr := cmd.CombinedOutput()

	// A run that passed its deadline is a failed sample naming the phase that
	// was in flight - which behave.sh's own result says, if its trap got far
	// enough to write one. If it did not, the sample is failed with the tail of
	// what it had printed.
	if deadlineCtx.Err() == context.DeadlineExceeded {
		// The outcome is ours and not the script's. Killing it lets its EXIT
		// trap run, which is what removes the project and the database - but the
		// trap writes the outcome it was holding, and a run killed mid-phase is
		// still holding "ok" with nothing asserted. A result that says ok
		// because it was interrupted before anything could go wrong is the one
		// reading a deadline must never produce.
		//
		// Failed rather than checks_failed: nothing was asserted, and a service
		// that never came up is not one that disagreed.
		phase := lastPhase(string(out))
		if parsed, _, err := readResult(results); err == nil && parsed.FailedPhase != "" {
			phase = parsed.FailedPhase
		}
		return BehaviorRun{
			Outcome:     "failed",
			FailedPhase: phase,
			Detail: fmt.Sprintf("the behaviour run passed its %s deadline during %s\n%s",
				behaviorDeadline(), phase, strings.TrimSpace(tail(string(out), 20))),
			Elapsed: time.Since(started),
		}
	}

	// The exit status is not the measurement. behave.sh exits zero on a run
	// whose assertions failed, because a failed assertion is a result; what
	// says whether the run happened at all is whether it wrote a result file.
	parsed, path, err := readResult(results)
	if err != nil {
		detail := strings.TrimSpace(tail(string(out), 20))
		if runErr != nil {
			return BehaviorRun{
				Err:     fmt.Errorf("behaviour harness failed for target %s: %w\n%s", target, runErr, detail),
				Elapsed: time.Since(started),
			}
		}
		return BehaviorRun{
			Err:     fmt.Errorf("behaviour harness wrote no result for target %s: %w\n%s", target, err, detail),
			Elapsed: time.Since(started),
		}
	}
	_ = path

	run := BehaviorRun{
		Outcome:     parsed.Outcome,
		FailedPhase: parsed.FailedPhase,
		Detail:      parsed.Detail,
		Checks:      parsed.ChecksTotal,
		Passed:      parsed.ChecksPassed,
		Attribution: parsed.Attribution,
		Elapsed:     time.Since(started),
	}
	for _, c := range parsed.Checks {
		if !c.Pass {
			run.Failed = append(run.Failed, c.Check)
		}
	}
	return run
}

// writeAnswer renders an invocation list in the envelope the stub model serves
// and Phase 4 decodes. It is the same shape a model returns, because it is
// standing in for one until grow --execute exists.
func writeAnswer(invocations []recording.Invocation) (string, error) {
	type wire struct {
		Action string         `json:"action"`
		Kwargs map[string]any `json:"kwargs"`
	}
	body := struct {
		Invocations []wire `json:"invocations"`
	}{}
	for _, inv := range invocations {
		body.Invocations = append(body.Invocations, wire{Action: inv.Action, Kwargs: inv.Kwargs})
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	f, err := os.CreateTemp("", "sedum-behave-answer-*.json")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Write(encoded); err != nil {
		return "", err
	}
	return f.Name(), nil
}

// readResult finds the one result file the run wrote. The directory is created
// per run, so exactly one is expected and more than one means two runs shared a
// directory - which is worth failing on rather than picking between.
func readResult(dir string) (harnessResult, string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return harnessResult{}, "", err
	}

	var found []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			found = append(found, filepath.Join(dir, e.Name()))
		}
	}
	switch len(found) {
	case 0:
		return harnessResult{}, "", fmt.Errorf("no result file in %s", dir)
	case 1:
	default:
		return harnessResult{}, "", fmt.Errorf("%d result files in %s; a run writes one", len(found), dir)
	}

	raw, err := os.ReadFile(found[0])
	if err != nil {
		return harnessResult{}, "", err
	}
	var parsed harnessResult
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return harnessResult{}, "", fmt.Errorf("reading %s: %w", found[0], err)
	}
	return parsed, found[0], nil
}

func tail(s string, lines int) string {
	parts := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(parts) <= lines {
		return s
	}
	return strings.Join(parts[len(parts)-lines:], "\n")
}

// BehaviorTally is what a measurement's behaviour runs added up to.
//
// Working, Disagreed and Broke are kept apart deliberately. A generator package
// that renders a service which does not compile and one that renders a service
// which compiles and violates its record are different findings with different
// fixes, and a single pass rate would report them as one number
// (prov-2026-83340ba0).
type BehaviorTally struct {
	// Measured is how many samples were applied at all. It is the denominator,
	// and it is not the sample count: a sample with no valid selection has
	// nothing to apply and is excluded rather than counted as a failure.
	Measured int

	// Working is samples where every assertion held.
	Working int
	// Disagreed is samples that built, booted and answered wrongly.
	Disagreed int
	// Broke is samples where a phase died, so the assertions never ran.
	Broke int
	// Errored is samples where the harness itself could not run. Excluded from
	// Measured, on the rule that keeps an unreachable endpoint out of the
	// selection denominator.
	Errored int

	// Checks and Passed are assertions summed across every measured sample.
	Checks int
	Passed int

	// Phases counts what broke, by phase name.
	Phases map[string]int
	// Details counts each distinct failure reason. Identical failures across
	// samples are one finding repeated, which is the thing worth knowing about
	// them - the same reason failed assertions are counted rather than listed.
	Details map[string]int
	// Failures counts each failed assertion across samples. One assertion
	// failing in every sample and twenty failing in one are different problems,
	// and only a per-assertion count tells them apart.
	Failures map[string]int

	// Actions counts the samples each action was named in, keyed by
	// action:variant, and Unattributed counts samples whose dead phase named a
	// line no marker enclosed.
	//
	// Per sample rather than per line. A build that names one action on three
	// lines is one sample dying in that action, and counting the lines would
	// report it as three - which is the confusion this exists to end, because
	// three samples dying in one action and three dying in three are the
	// different findings (prov-2026-27c10ac4).
	//
	// AttributedSamples is the denominator: samples that carried any
	// attribution at all. It is not Broke, because an entry drawn before this
	// field existed carries none and reads as absent rather than as a build no
	// action was responsible for.
	Actions           map[string]int
	Unattributed      int
	AttributedSamples int

	// Elapsed is what behaviour added to the run.
	Elapsed time.Duration
}

// Rate is the fraction of measured samples that produced a working application.
func (b BehaviorTally) Rate() float64 {
	if b.Measured == 0 {
		return 0
	}
	return float64(b.Working) / float64(b.Measured)
}

// attribute folds one sample's attributions into the tally, counting each
// action once however many of the sample's lines it wrote.
func (t *BehaviorTally) attribute(found []Attribution) {
	if len(found) == 0 {
		// Absent, not none. Nothing is incremented, so a run of entries drawn
		// before attribution existed leaves every count at zero and the report
		// says nothing rather than saying no action was responsible.
		return
	}
	t.AttributedSamples++

	// Distinct within the sample, on the rule Details already follows: one
	// finding repeated inside one observation is one finding.
	seen := map[string]bool{}
	var unattributed bool
	for _, a := range found {
		if !a.Attributed() {
			unattributed = true
			continue
		}
		if seen[a.Label()] {
			continue
		}
		seen[a.Label()] = true
		t.Actions[a.Label()]++
	}
	if unattributed {
		t.Unattributed++
	}
}

// Behavior tallies what applying this measurement's selections produced. The
// second return is false when behaviour was not measured at all, which a report
// has to tell from "measured and nothing worked".
func (m Measurement) Behavior() (BehaviorTally, bool) {
	t := BehaviorTally{
		Phases:   map[string]int{},
		Failures: map[string]int{},
		Details:  map[string]int{},
		Actions:  map[string]int{},
	}

	var any bool
	for _, s := range m.Samples {
		if s.Behavior == nil {
			continue
		}
		any = true
		b := *s.Behavior
		t.Elapsed += b.Elapsed

		if b.Err != nil {
			t.Errored++
			continue
		}

		t.Measured++
		t.Checks += b.Checks
		t.Passed += b.Passed
		for _, f := range b.Failed {
			t.Failures[f]++
		}

		switch b.Outcome {
		case "ok":
			t.Working++
		case "checks_failed":
			t.Disagreed++
		default:
			t.Broke++
			if b.FailedPhase != "" {
				t.Phases[b.FailedPhase]++
			}
			if d := strings.TrimSpace(b.Detail); d != "" {
				t.Details[d]++
			}
			t.attribute(b.Attribution)
		}
	}
	return t, any
}

// sortedNames orders variable names so a run's environment is the same every
// time, which is what keeps two runs of one case comparable.
func sortedNames(variables map[string]string) []string {
	out := make([]string, 0, len(variables))
	for name := range variables {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// behaviorDeadline is how long one behaviour run may take.
//
// Generous on purpose: the phases of one sample are tens of seconds, and a cold
// bundle install on a fresh machine is minutes and is not a hang. It exists for
// the case that is neither - a phase that will never finish - which stopped two
// runs in one session and cost a row of fifty samples, because an entry is only
// written when a row finishes (prov-2026-3957eed2).
func behaviorDeadline() time.Duration {
	if v := os.Getenv("SEDUM_BEHAVIOR_DEADLINE"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return 10 * time.Minute
}

// lastPhase reads the phase in flight out of what the script had printed, for
// the case where the deadline fired before its trap wrote a result.
//
// The script announces each phase as it starts it, so the last one announced is
// the one it was in. A guess from output rather than a field, and it is only
// reached when the field does not exist.
func lastPhase(out string) string {
	phase := ""
	for _, line := range strings.Split(out, "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "\u2192 "); ok {
			phase = strings.TrimSpace(rest)
		}
	}
	return phase
}
