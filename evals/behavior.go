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
	"time"

	"github.com/calebcowen/sedum/internal/recording"
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

	Elapsed time.Duration `json:"elapsed"`

	// Err is set when the harness could not be run at all - a missing target, a
	// script that would not start. Excluded from every rate, on the same rule
	// that keeps an unreachable endpoint out of the selection denominator.
	Err error `json:"-"`
}

// Working reports whether this run produced an application that did everything
// the target asserts.
func (b BehaviorRun) Working() bool { return b.Err == nil && b.Outcome == "ok" }

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
	// samples wants none of that accumulating.
	results, err := os.MkdirTemp("", "sedum-behave-results-")
	if err != nil {
		return BehaviorRun{Err: err, Elapsed: time.Since(started)}
	}
	defer os.RemoveAll(results)

	cmd := exec.CommandContext(ctx, "bash", harnessScript, target, "--answer", answer)
	cmd.Env = append(os.Environ(), "RESULTS_DIR="+results)
	// The case's variables reach the target as environment, so the scaffold and
	// the generation read one value rather than two that agree by convention:
	// go mod init and --var are the same module name or the service does not
	// compile (prov-2026-1c33a50b).
	for _, name := range sortedNames(variables) {
		cmd.Env = append(cmd.Env, "SEDUM_VAR_"+strings.ToUpper(name)+"="+variables[name])
	}
	out, runErr := cmd.CombinedOutput()

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

// Behavior tallies what applying this measurement's selections produced. The
// second return is false when behaviour was not measured at all, which a report
// has to tell from "measured and nothing worked".
func (m Measurement) Behavior() (BehaviorTally, bool) {
	t := BehaviorTally{Phases: map[string]int{}, Failures: map[string]int{}, Details: map[string]int{}}

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
