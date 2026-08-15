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
//	  -eval.samples   runs per model, default 5
//	  -eval.concurrency  samples in flight at once, default 1
//	  -eval.root      where the fixture applications live, default ../..
//	                  (go test runs with cwd at evals/, so this is the workspace)
package evals

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/selection"
)

var (
	caseID      = flag.String("eval.case", "", "run only the case with this id")
	modelID     = flag.String("eval.model", "", "run only models whose label contains this; one row at a time when memory is tight")
	samples     = flag.Int("eval.samples", 5, "runs per model")
	concurrency = flag.Int("eval.concurrency", 1, "samples in flight at once; raise it against a server with continuous batching")
	root        = flag.String("eval.root", "../..", "directory the fixture applications live under")
)

func TestEval(t *testing.T) {
	if os.Getenv(selection.EnvBaseURL) == "" && os.Getenv(selection.EnvAPIKey) == "" {
		t.Skipf("no model endpoint: set %s (a local server) or %s", selection.EnvBaseURL, selection.EnvAPIKey)
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
				m, err := Run(context.Background(), c, model, *samples, *concurrency)
				if err != nil {
					t.Fatalf("running case: %v", err)
				}

				// Written to stdout rather than through t.Log so the table
				// arrives unindented and can be pasted into a record as it
				// stands.
				Report(os.Stdout, m)

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
