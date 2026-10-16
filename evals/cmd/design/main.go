package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/livecodelife/sedum/evals/design"
)

const maxK = 5

// main runs the behaviour-only-intent eval: it draws K samples of every task and
// variant through the herd's Haiku wrapper and prints the table and the total tokens.
func main() {
	tasksDir := flag.String("tasks", "evals/design/testdata/tasks", "directory of tasks")
	generators := flag.String("generators", "generators", "directory of generator packages")
	wrapper := flag.String("wrapper", filepath.Join(os.Getenv("HOME"), "workspace/.herd/sedum-haiku"), "the herd's Haiku wrapper")
	sedumBin := flag.String("sedum", "", "sedum binary (required)")
	k := flag.Int("n", maxK, "samples per cell (at most 5)")
	only := flag.String("task", "", "run only this task id")
	variant := flag.String("variant", "", "run only this variant")
	dry := flag.Bool("dry", false, "print the plan and run nothing")
	keep := flag.Bool("keep", false, "keep each sample's project on disk")
	out := flag.String("out", "", "results file (default evals/results/design-<time>.json)")
	flag.Parse()

	if *k < 1 || *k > maxK {
		fail("-n must be between 1 and %d: raising K needs approval", maxK)
	}
	tasks, err := design.LoadTasks(*tasksDir)
	if err != nil {
		fail("%v", err)
	}
	if *only != "" {
		var kept []design.Task
		for _, t := range tasks {
			if t.ID == *only {
				kept = append(kept, t)
			}
		}
		tasks = kept
	}
	variants := design.Variants()
	if *variant != "" {
		variants = []design.Variant{design.Variant(*variant)}
	}
	if len(tasks) == 0 || len(variants) == 0 {
		fail("nothing to run")
	}
	fmt.Print(design.Plan(tasks, variants, *k, 20000))
	if *dry {
		return
	}
	if *sedumBin == "" {
		fail("-sedum is required")
	}
	cfg := design.Config{Generators: *generators, Grow: design.HaikuGrow(*wrapper, *sedumBin), Keep: *keep}
	samples, err := design.Run(context.Background(), cfg, tasks, variants, *k)
	if err != nil {
		fail("%v", err)
	}
	hashes := map[string]string{}
	for _, t := range tasks {
		h, err := design.TaskHash(t)
		if err != nil {
			fail("%v", err)
		}
		hashes[t.ID] = h
	}
	report := design.NewReport(gitCommit(), "sedum", "claude-haiku-5-5", *k, hashes, samples)
	path := *out
	if path == "" {
		path = filepath.Join("evals", "results", "design-"+time.Now().UTC().Format("20060102T150405Z")+".json")
	}
	data, err := report.JSON()
	if err != nil {
		fail("%v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		fail("%v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		fail("%v", err)
	}
	fmt.Print(report.Table())
	fmt.Printf("\ntotal tokens: %d (prompt %d, completion %d)\nresults: %s\n", report.Totals.Tokens(), report.Totals.PromptTokens, report.Totals.CompletionTokens, path)
}

func gitCommit() string {
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "design: "+format+"\n", args...)
	os.Exit(1)
}
