package main

import (
	"strings"
	"testing"
	"time"

	"github.com/calebcowen/sedum/evals"
	"github.com/calebcowen/sedum/internal/selection"
)

func fixture(id string, models ...string) evals.Case {
	c := evals.Case{ID: id}
	for _, m := range models {
		c.Models = append(c.Models, evals.Model{ID: "qwen", Engine: m, Quant: "q4"})
	}
	return c
}

// The timeout is the sharpest of the three things this command exists to stop
// being remembered. go test defaults to ten minutes, a fine run is forty-five,
// and a forgotten -timeout does not fail fast - it kills the run partway and no
// entry is written.
func TestTheTimeoutIsSizedFromTheRunNotDefaulted(t *testing.T) {
	fine := planFor(fixture("x", "llama.cpp", "mlx"), "", evals.Fine, 0, 0, false)
	if fine.Samples != 30 || fine.Models != 2 {
		t.Fatalf("planned %d models x %d samples, want 2 x 30", fine.Models, fine.Samples)
	}
	if want := 60 * perSample; fine.Expect != want {
		t.Errorf("expected duration is %s, want %s", fine.Expect, want)
	}
	if want := fine.Expect * headroom; fine.Timeout != want {
		t.Errorf("timeout is %s, want %s", fine.Timeout, want)
	}
	// The number that matters: comfortably past go test's own default, which is
	// what a fine run silently dies against.
	if fine.Timeout <= 10*time.Minute {
		t.Errorf("timeout of %s does not clear go test's ten-minute default", fine.Timeout)
	}

	// One row, because a 24GB machine holds one 14B resident and the arms have
	// to face the same one.
	pinned := planFor(fixture("x", "llama.cpp", "mlx"), "llama.cpp", evals.Fine, 0, 0, false)
	if pinned.Models != 1 {
		t.Errorf("the model filter left %d rows, want 1", pinned.Models)
	}
	if pinned.Timeout != fine.Timeout/2 || pinned.Expect != fine.Expect/2 {
		t.Errorf("half the rows did not halve the plan: %s/%s against %s/%s",
			pinned.Expect, pinned.Timeout, fine.Expect, fine.Timeout)
	}

	// A smoke run is two samples and would otherwise be given six minutes, which
	// a cold model load can eat on its own.
	if smoke := planFor(fixture("x", "llama.cpp"), "", evals.Smoke, 0, 0, false); smoke.Timeout != minTimeout {
		t.Errorf("a smoke run got %s, want the %s floor", smoke.Timeout, minTimeout)
	}
}

// The expectation and the ceiling answer different questions - whether this can
// be started now, and when it should be considered hung - so the plan prints
// both and says which is which. Printing only the ceiling read as the cost and
// doubled it (prov-2026-6e3c846c).
func TestThePlanSeparatesWhatARunTakesFromWhatItIsAllowed(t *testing.T) {
	arm := planFor(fixture("todo-rails-described", "llama.cpp"), "", evals.Fine, 0, 0, false)

	// 30 samples at the observed rate is the three quarters of an hour this has
	// been quoted at all along; the ceiling is twice it, and neither is the
	// other.
	if arm.Expect != 45*time.Minute {
		t.Errorf("a fine arm is expected to take %s, want 45m", arm.Expect)
	}
	if arm.Timeout != 90*time.Minute {
		t.Errorf("a fine arm is allowed %s, want 90m", arm.Timeout)
	}

	out := summarize([]plan{arm, arm}, evals.Fine)
	if !strings.Contains(out, "~45m") {
		t.Errorf("the plan does not say what the run takes:\n%s", out)
	}
	if !strings.Contains(out, "timeout 1h30m") {
		t.Errorf("the plan does not label the ceiling as a timeout:\n%s", out)
	}
	// The total is of the expectations. A total of the ceilings said three
	// hours for a pair that costs ninety minutes.
	if !strings.Contains(out, "~1h30m in total") {
		t.Errorf("the total is not of the expected durations:\n%s", out)
	}
}

// There was no way to set the timeout at all, which made the computed number a
// ceiling in the ordinary case and a wall in every other one.
func TestTheCeilingCanBeGivenRatherThanDerived(t *testing.T) {
	given := planFor(fixture("x", "llama.cpp"), "", evals.Fine, 0, 20*time.Minute, false)
	if given.Timeout != 20*time.Minute {
		t.Errorf("the override was not honoured: %s", given.Timeout)
	}
	// Downward too. A deliberately short budget is a legitimate way to find out
	// whether an endpoint is answering at all.
	if given.Timeout >= given.Expect {
		t.Errorf("a budget below the expectation was raised to %s", given.Timeout)
	}
	// The expectation is untouched by it: what the run will take does not change
	// because it was given less room.
	if given.Expect != 45*time.Minute {
		t.Errorf("the override moved the expectation to %s", given.Expect)
	}

	// A given ceiling is shown as given, so the plan never prints a number
	// nobody computed as though it had computed it.
	if !strings.Contains(summarize([]plan{given}, evals.Fine), "timeout 20m, given") {
		t.Errorf("the plan does not say the ceiling was given:\n%s", summarize([]plan{given}, evals.Fine))
	}
	if strings.Contains(summarize([]plan{planFor(fixture("x", "llama.cpp"), "", evals.Fine, 0, 0, false)}, evals.Fine), "given") {
		t.Error("a derived ceiling was reported as given")
	}
}

// An explicit count above the resolution's own is honoured, because oversampling
// is a cost decision and misleads nobody. Below it the runner refuses, so the
// budget is sized for what will actually be drawn (prov-2026-3039750e).
func TestTheBudgetFollowsTheSampleCountThatWillBeDrawn(t *testing.T) {
	over := planFor(fixture("x", "llama.cpp"), "", evals.Coarse, 40, 0, false)
	if over.Samples != 40 {
		t.Errorf("planned %d samples, want the explicit 40", over.Samples)
	}
	under := planFor(fixture("x", "llama.cpp"), "", evals.Fine, 5, 0, false)
	if under.Samples != 30 {
		t.Errorf("planned %d samples for a count the runner will refuse, want the resolution's 30", under.Samples)
	}
}

// A typo in a case id costs a second here rather than a run that quietly
// measures nothing, which is what the runner's own -eval.case does with one.
func TestAnUnknownCaseIsNamedRatherThanSilentlyMatchingNothing(t *testing.T) {
	all := []evals.Case{fixture("todo-rails-defined", "llama.cpp"), fixture("todo-rails-described", "llama.cpp")}

	_, err := selectCases(all, []string{"todo-rails-descibed"})
	if err == nil {
		t.Fatal("a misspelled case was accepted")
	}
	// The message has to carry the real ids, or the reader is left to go and
	// look them up.
	for _, want := range []string{"todo-rails-descibed", "todo-rails-defined", "todo-rails-described"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error does not name %q: %v", want, err)
		}
	}

	// The order given is the order run, because an A/B is run in the order its
	// author wrote it.
	got, err := selectCases(all, []string{"todo-rails-described", "todo-rails-defined"})
	if err != nil {
		t.Fatalf("selecting: %v", err)
	}
	if len(got) != 2 || got[0].ID != "todo-rails-described" {
		t.Errorf("selected %v, want the requested order", caseIDs(got))
	}

	// No argument is every case, which is what the bare command means.
	if got, _ := selectCases(all, nil); len(got) != 2 {
		t.Errorf("the bare command selected %d cases, want all of them", len(got))
	}
}

// Every flag this command passes is one the tagged runner already declares. It
// is a wrapper, and a wrapper that invented a flag would be measuring something
// the documented invocation does not.
func TestTheInvocationIsOneTheRunnerAlreadyAccepts(t *testing.T) {
	args := testArgs("todo-rails-described", "llama.cpp", evals.Fine, 0, 1, 0, "evals/testdata", 90*time.Minute, false)
	joined := strings.Join(args, " ")

	for _, want := range []string{
		"test -tags eval ./evals",
		"-timeout 1h30m0s",
		"-eval.case todo-rails-described",
		"-eval.resolution fine",
		"-eval.model llama.cpp",
		// The runner resolves the fixture root against the package it runs in,
		// so the repository-root path this command takes is trimmed to what go
		// test will see.
		"-eval.root testdata",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("the invocation is missing %q:\n  go %s", want, joined)
		}
	}

	// -count=1 because a cached result is not a measurement.
	if !strings.Contains(joined, "-count=1") {
		t.Errorf("the invocation would accept a cached result:\n  go %s", joined)
	}

	// An unset filter is an absent flag rather than an empty one: the runner
	// reads -eval.samples of zero as "draw what the resolution calls for", and
	// an empty -eval.model would match every label rather than none.
	bare := strings.Join(testArgs("x", "", evals.Coarse, 0, 1, 0, "evals/testdata", time.Hour, false), " ")
	for _, absent := range []string{"-eval.model", "-eval.samples"} {
		if strings.Contains(bare, absent) {
			t.Errorf("%s was passed with nothing to say:\n  go %s", absent, bare)
		}
	}
}

// The local-server defaults are applied only where the environment says nothing,
// so a run against a different endpoint is configured exactly as it was before
// this command existed.
func TestTheEndpointDefaultsWithoutOverridingTheEnvironment(t *testing.T) {
	empty := environ(func(string) string { return "" })
	if empty[selection.EnvBaseURL] != defaultBaseURL || empty[selection.EnvAPIKey] != defaultAPIKey {
		t.Errorf("an empty environment did not get the local defaults: %v", empty)
	}

	set := environ(func(k string) string {
		if k == selection.EnvBaseURL {
			return "https://api.example.com/v1"
		}
		return ""
	})
	if _, overridden := set[selection.EnvBaseURL]; overridden {
		t.Errorf("a configured endpoint was overridden: %v", set)
	}
	if set[selection.EnvAPIKey] != defaultAPIKey {
		t.Errorf("the key was not defaulted beside a configured endpoint: %v", set)
	}
}

// The between-cases refusal is the one that matters for an A/B. The first arm's
// own result file is untracked when it finishes, and Clean is read from git
// status, so the second arm would record as not re-runnable because the first
// one succeeded.
func TestTheSecondArmIsNotRunAgainstWhatTheFirstArmWrote(t *testing.T) {
	plans := []plan{{Case: "todo-rails-defined"}, {Case: "todo-rails-described"}}
	wrote := []string{"?? evals/results/todo-rails-defined.jsonl"}

	err := checkClean(wrote, false, false, plans, 1, "evals/results")
	if err == nil {
		t.Fatal("the second arm was allowed to run against the first arm's uncommitted entry")
	}
	for _, want := range []string{
		"todo-rails-defined",
		"evals/results/todo-rails-defined.jsonl",
		"todo-rails-described",
		// It stops and says what to commit rather than committing: the commit
		// needs a provenance tag this command has no way to choose.
		"git commit",
		"go run ./evals/cmd/eval todo-rails-described",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the refusal does not name %q:\n%v", want, err)
		}
	}

	// -dirty is the deliberate override, and it is what a mid-edit measurement
	// now costs.
	if err := checkClean(wrote, true, false, plans, 1, "evals/results"); err != nil {
		t.Errorf("-dirty did not permit the run: %v", err)
	}

	// A clean tree between arms is the whole point of having committed the
	// first one, and it runs without argument.
	if err := checkClean(nil, false, false, plans, 1, "evals/results"); err != nil {
		t.Errorf("a committed first arm did not let the second run: %v", err)
	}

	// Before the first case it is the ordinary check, and it says what is dirty
	// rather than only that something is.
	first := checkClean([]string{" M evals/run.go"}, false, true, plans, 0, "evals/results")
	if first == nil {
		t.Fatal("a dirty tree was allowed to produce an un-pinned entry")
	}
	if !strings.Contains(first.Error(), "evals/run.go") || !strings.Contains(first.Error(), "-dirty") {
		t.Errorf("the refusal does not name the change or the override:\n%v", first)
	}
}

// The repository this test runs in is the fixture for case loading, so the
// command and the runner agree about where cases and fixtures live.
func TestTheDefaultPathsFindTheRealCases(t *testing.T) {
	// t.Chdir rather than os.Chdir: the working directory is process-wide, and
	// a test that moved it permanently would decide where every test after it
	// ran.
	t.Chdir("../../..")

	cases, err := evals.Load("evals/cases", "evals/testdata")
	if err != nil {
		t.Fatalf("the command's default paths do not load: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("the command's default paths found no cases")
	}
}

// A plan is read in tens of minutes, and Duration.String always carries a
// seconds component it has nothing to say with.
func TestDurationsReadAsAPlanRatherThanAStopwatch(t *testing.T) {
	for _, tc := range []struct {
		d    time.Duration
		want string
	}{
		{45 * time.Minute, "45m"},
		{90 * time.Minute, "1h30m"},
		{60 * time.Minute, "1h"},
		{10 * time.Minute, "10m"},
		{2 * time.Hour, "2h"},
		{30 * time.Second, "30s"},
	} {
		if got := dur(tc.d); got != tc.want {
			t.Errorf("%s formats as %q, want %q", tc.d, got, tc.want)
		}
	}
}
