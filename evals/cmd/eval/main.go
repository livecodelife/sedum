// Command eval runs the eval harness without the invocation being something to
// remember.
//
// The documented way to start a measurement is two environment variables, a
// build tag, a verbosity flag, a hand-guessed timeout and four flags behind a
// prefix. Three of those are failure modes rather than verbosity: go test's
// ten-minute default kills a fine run partway, a dirty tree produces a number
// that is not pinned to a commit, and a second arm run straight after a first
// records dirty because the first one's result file is untracked
// (prov-2026-d31cbfee).
//
//	go run ./evals/cmd/eval                                  # every case, coarse
//	go run ./evals/cmd/eval todo-rails-described             # one case
//	go run ./evals/cmd/eval -res fine -model llama.cpp \
//	  todo-rails-defined todo-rails-described                # the description A/B
//	go run ./evals/cmd/eval -dry -res fine todo-rails-defined
//
// It is a wrapper, not a second runner. Every flag it passes is one the tagged
// runner already accepts, and it prints the invocation before running it -
// nothing about how a sample is drawn, scored or recorded lives here. It shells
// out to go test rather than calling the harness directly because the runner is
// deliberately behind a build tag in a test file, and a second entry point into
// it would duplicate that judgement and eventually disagree with it.
//
// Like cmd/history it lives outside the sedum binary: the harness is not part of
// Sedum, and nothing under internal/ or cmd/ imports it (prov-2026-c0f55691).
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/calebcowen/sedum/evals"
	"github.com/calebcowen/sedum/internal/selection"
)

const (
	// defaultBaseURL is where a local server listens by default. It is applied
	// only when the environment names none, so a run against anything else is
	// configured the same way it always was.
	defaultBaseURL = "http://127.0.0.1:1234/v1"
	// defaultAPIKey is a placeholder a local server ignores and the client
	// requires.
	defaultAPIKey = "local"

	// perSample is what one sample is budgeted at when sizing the timeout. The
	// observed range on the local 14B row is about 66 to 90 seconds, so this is
	// generous by roughly a factor of two - which is the right direction to be
	// wrong in, because the cost of overestimating is a timeout that never
	// fires and the cost of underestimating is a run killed after forty
	// minutes with nothing recorded.
	perSample = 3 * time.Minute
	// minTimeout keeps a two-sample smoke run from being given a budget so
	// small that a cold model load eats it.
	minTimeout = 15 * time.Minute
)

func main() {
	var (
		res         = flag.String("res", string(evals.DefaultResolution), "the question this run is asking: smoke, coarse or fine")
		model       = flag.String("model", "", "run only models whose label contains this")
		samples     = flag.Int("n", 0, "runs per model; zero draws what the resolution calls for")
		concurrency = flag.Int("c", 1, "samples in flight at once")
		retries     = flag.Int("retries", 0, "re-prompts a rejected answer may spend")
		caseDir     = flag.String("cases", "evals/cases", "directory the case files live in")
		root        = flag.String("root", "evals/testdata", "directory the vendored fixtures live under")
		allowDirty  = flag.Bool("dirty", false, "run against a dirty tree; the entry records as not re-runnable")
		dry         = flag.Bool("dry", false, "print the invocation and the plan, run nothing")
	)
	flag.Parse()

	if err := run(*res, *model, *samples, *concurrency, *retries, *caseDir, *root, *allowDirty, *dry, flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "eval: %v\n", err)
		os.Exit(1)
	}
}

func run(res, model string, samples, concurrency, retries int, caseDir, root string, allowDirty, dry bool, want []string) error {
	// Parsed first, so a misspelled resolution costs a second rather than the
	// minutes it takes to find out after a case has run - the same reason the
	// runner parses it before loading anything.
	resolution, err := evals.ParseResolution(res)
	if err != nil {
		return err
	}

	all, err := evals.Load(caseDir, root)
	if err != nil {
		return fmt.Errorf("loading cases: %w", err)
	}
	selected, err := selectCases(all, want)
	if err != nil {
		return err
	}

	plans := make([]plan, 0, len(selected))
	for _, c := range selected {
		plans = append(plans, planFor(c, model, resolution, samples))
	}
	for _, p := range plans {
		if p.Models == 0 {
			return fmt.Errorf("case %s has no model matching %q; it declares %s",
				p.Case, model, strings.Join(labels(byID(selected, p.Case).Models), ", "))
		}
	}

	fmt.Println(summarize(plans, resolution))

	// The results directory is relative to the package go test runs in, which
	// is why the runner's default is "results" and this one's is the path from
	// the repository root.
	resultsDir := filepath.Join(filepath.Dir(caseDir), "results")

	for i, p := range plans {
		args := testArgs(p.Case, model, resolution, samples, concurrency, retries, root, p.Timeout)
		fmt.Printf("\n$ %s\n\n", strings.Join(append(envPrefix(), append([]string{"go"}, args...)...), " "))
		if dry {
			continue
		}

		dirty, err := gitDirty()
		if err != nil {
			dirty = nil
		}
		if err := checkClean(dirty, allowDirty, i == 0, plans, i, resultsDir); err != nil {
			return err
		}
		if err := goTest(args); err != nil {
			return fmt.Errorf("case %s: %w", p.Case, err)
		}
	}
	return nil
}

// plan is what one case is about to cost.
type plan struct {
	Case    string
	Models  int
	Samples int
	Timeout time.Duration
}

// planFor sizes one case: how many model rows the filter leaves, how many
// samples each draws, and how long that is allowed to take.
//
// The timeout is computed rather than defaulted because go test's own default
// is ten minutes and a fine run is forty-five. A forgotten -timeout does not
// fail fast - it kills the run partway and no entry is written.
func planFor(c evals.Case, model string, res evals.Resolution, samples int) plan {
	n := samples
	if n < res.Samples() {
		n = res.Samples()
	}

	models := 0
	for _, m := range c.Models {
		if model == "" || strings.Contains(m.Label(), model) {
			models++
		}
	}

	timeout := time.Duration(models*n) * perSample
	if timeout < minTimeout {
		timeout = minTimeout
	}
	return plan{Case: c.ID, Models: models, Samples: n, Timeout: timeout}
}

// selectCases resolves the requested ids against what is on disk. An unknown id
// is an error rather than a case that quietly matches nothing, which is what the
// runner's own -eval.case does with a typo.
func selectCases(all []evals.Case, want []string) ([]evals.Case, error) {
	if len(want) == 0 {
		return all, nil
	}

	out := make([]evals.Case, 0, len(want))
	for _, id := range want {
		c := byID(all, id)
		if c.ID == "" {
			return nil, fmt.Errorf("unknown case %q; the cases are %s", id, strings.Join(caseIDs(all), ", "))
		}
		out = append(out, c)
	}
	return out, nil
}

func byID(cases []evals.Case, id string) evals.Case {
	for _, c := range cases {
		if c.ID == id {
			return c
		}
	}
	return evals.Case{}
}

func caseIDs(cases []evals.Case) []string {
	out := make([]string, 0, len(cases))
	for _, c := range cases {
		out = append(out, c.ID)
	}
	return out
}

func labels(models []evals.Model) []string {
	out := make([]string, 0, len(models))
	for _, m := range models {
		out = append(out, m.Label())
	}
	return out
}

// summarize says what the run will cost before it starts, because the number
// worth knowing at that point is hours rather than flags.
func summarize(plans []plan, res evals.Resolution) string {
	var b strings.Builder
	var total time.Duration

	fmt.Fprintf(&b, "%s resolution, %d case(s):\n", res, len(plans))
	for _, p := range plans {
		fmt.Fprintf(&b, "  %-24s %d model(s) x %d samples, up to %s\n",
			p.Case, p.Models, p.Samples, p.Timeout)
		total += p.Timeout
	}
	if len(plans) > 1 {
		fmt.Fprintf(&b, "  %-24s up to %s in total\n", "", total)
	}
	return b.String()
}

// testArgs is the go test invocation for one case. Every flag here is one the
// tagged runner already declares.
func testArgs(caseID, model string, res evals.Resolution, samples, concurrency, retries int, root string, timeout time.Duration) []string {
	args := []string{
		"test", "-tags", "eval", "./evals", "-v", "-count=1",
		"-timeout", timeout.String(),
		"-eval.case", caseID,
		"-eval.resolution", string(res),
		"-eval.concurrency", fmt.Sprint(concurrency),
		"-eval.retries", fmt.Sprint(retries),
		// The runner resolves this against the package it runs in, so the
		// repository-root path this command takes is trimmed back to what
		// go test will see.
		"-eval.root", filepath.Base(root),
	}
	if model != "" {
		args = append(args, "-eval.model", model)
	}
	if samples > 0 {
		args = append(args, "-eval.samples", fmt.Sprint(samples))
	}
	return args
}

// environ applies the local-server defaults without overriding anything the
// environment already sets, so a run against a different endpoint is configured
// exactly as it was before this command existed.
func environ(get func(string) string) map[string]string {
	out := map[string]string{}
	if get(selection.EnvBaseURL) == "" {
		out[selection.EnvBaseURL] = defaultBaseURL
	}
	if get(selection.EnvAPIKey) == "" {
		out[selection.EnvAPIKey] = defaultAPIKey
	}
	return out
}

func envPrefix() []string {
	var out []string
	for k, v := range environ(os.Getenv) {
		out = append(out, k+"="+v)
	}
	return out
}

func goTest(args []string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	cmd.Env = os.Environ()
	for k, v := range environ(os.Getenv) {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	return cmd.Run()
}

// checkClean refuses to spend model calls on a number that will not be pinned.
//
// Before the first case it is the ordinary dirty-tree check. Between cases it is
// the one that matters: the previous case's own result file is untracked, so the
// next arm would record dirty because the previous arm succeeded. Stopping is
// the whole point - an A/B whose second arm is un-pinned is half a measurement,
// and it arrives looking like a whole one.
//
// The working tree arrives as an argument rather than being read here, so the
// refusal can be tested without the test's own repository deciding the answer.
// An empty list also covers the case where git could not be read at all: a
// measurement taken outside a checkout is still a measurement, and the entry
// already records that it carries no commit.
func checkClean(dirty []string, allowDirty, first bool, plans []plan, i int, resultsDir string) error {
	if allowDirty || len(dirty) == 0 {
		return nil
	}

	if first {
		return fmt.Errorf("the tree is dirty, so this entry would not be re-runnable from its commit:\n  %s\ncommit first, or pass -dirty to record it as un-pinned anyway",
			strings.Join(dirty, "\n  "))
	}
	return fmt.Errorf("%s wrote %s and it is uncommitted, so %s would record as not re-runnable.\ncommit it, then re-run for the remaining case(s):\n\n  git add %s && git commit -m \"...[prov-...]\"\n  go run ./evals/cmd/eval %s",
		plans[i-1].Case,
		filepath.Join(resultsDir, plans[i-1].Case+".jsonl"),
		plans[i].Case,
		filepath.Join(resultsDir, plans[i-1].Case+".jsonl"),
		strings.Join(planIDs(plans[i:]), " "))
}

func planIDs(plans []plan) []string {
	out := make([]string, 0, len(plans))
	for _, p := range plans {
		out = append(out, p.Case)
	}
	return out
}

func gitDirty() ([]string, error) {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return nil, err
	}
	var dirty []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			dirty = append(dirty, line)
		}
	}
	return dirty, nil
}
