//go:build eval

// The runner is behind a build tag so that go test ./... never reaches it. It
// needs a model endpoint, it takes minutes rather than milliseconds, and it
// reports a sampled rate rather than a verdict - three reasons it must not sit
// in a suite whose whole value is being deterministic.
//
//	OPENAI_BASE_URL=http://127.0.0.1:1234/v1 OPENAI_API_KEY=local \
//	  go test -tags eval ./evals -v -timeout 60m
//
//	  -eval.case      run one case by id, default all
//	  -eval.model     run only models whose label contains this substring
//	  -eval.resolution  the question being asked: smoke, coarse or fine, default coarse
//	  -eval.samples   runs per model, default the resolution's own count
//	  -eval.concurrency  samples in flight at once, default 1
//	  -eval.retries   re-prompts a rejected answer may spend, default 0
//	  -eval.root      where the vendored fixtures live, default testdata
//	  -eval.results   where results are appended, default results ("" disables)
package evals

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/selection"
)

var (
	caseID      = flag.String("eval.case", "", "run only the case with this id")
	modelID     = flag.String("eval.model", "", "run only models whose label contains this; one row at a time when memory is tight")
	resolution  = flag.String("eval.resolution", string(DefaultResolution), "the question this run is asking: smoke (plumbing only), coarse (large differences) or fine (moving a rate that is already high)")
	samples     = flag.Int("eval.samples", 0, "runs per model; zero draws what the resolution calls for, and a count below it is refused rather than relabelled")
	concurrency = flag.Int("eval.concurrency", 1, "samples in flight at once; raise it against a server with continuous batching")
	root        = flag.String("eval.root", "testdata", "directory the vendored fixtures live under")
	results     = flag.String("eval.results", "results", "directory results are appended to; empty disables recording")
	retries     = flag.Int("eval.retries", 0, "re-prompts a rejected answer may spend; zero measures what one call produces, and raising it buys valid samples at the cost of comparable timing")
)

func TestEval(t *testing.T) {
	if os.Getenv(selection.EnvBaseURL) == "" && os.Getenv(selection.EnvAPIKey) == "" {
		t.Skipf("no model endpoint: set %s (a local server) or %s", selection.EnvBaseURL, selection.EnvAPIKey)
	}

	// Parsed before anything is loaded, so a misspelled resolution costs a
	// second rather than the minutes it would take to find out after the first
	// case has run.
	res, err := ParseResolution(*resolution)
	if err != nil {
		t.Fatalf("%v", err)
	}

	cases, err := Load("cases", *root)
	if err != nil {
		t.Fatalf("loading cases: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("no cases found under cases/")
	}

	for _, c := range cases {
		if *caseID != "" && c.ID != *caseID {
			continue
		}
		for _, model := range c.Models {
			if *modelID != "" && !strings.Contains(model.Label(), *modelID) {
				continue
			}
			t.Run(c.ID+"/"+model.Label(), func(t *testing.T) {
				m, err := Run(context.Background(), c, model, Options{
					Resolution:  res,
					Samples:     *samples,
					Concurrency: *concurrency,
					Retries:     *retries,
				})
				if err != nil {
					t.Fatalf("running case: %v", err)
				}

				// Written to stdout rather than through t.Log so the table
				// arrives unindented and can be pasted into a record as it
				// stands.
				Report(os.Stdout, m)

				if *results != "" {
					entry := NewEntry(m, os.Getenv(selection.EnvBaseURL))
					if !entry.Clean {
						// The commit is what makes an entry re-runnable, and a
						// dirty tree makes it a lie. The run is still recorded
						// - refusing would rule out measuring mid-edit, which
						// is most of how this gets used - but it says so here
						// and in the entry (prov-2026-eb283c56).
						fmt.Fprintf(os.Stdout,
							"  WARNING: %d uncommitted change(s); this entry is not re-runnable from %s\n",
							len(entry.Dirty), entry.Commit)
					}
					if err := Append(*results, entry); err != nil {
						t.Errorf("recording the result: %v", err)
					}
				}

				// A measurement does not fail. The one thing that does is a
				// case that produced nothing to measure, because that is a
				// broken harness rather than a poor rate.
				if _, scored, _ := m.Scored(); scored == 0 {
					t.Fatal("no sample completed; the endpoint or the case is broken, not the model")
				}
			})
		}
	}
}
