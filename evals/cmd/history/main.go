// Command history prints what the eval results file already knows.
//
// It lives outside the sedum binary on purpose. The harness is not part of
// Sedum - nothing under internal/ or cmd/ imports it - and a reporting command
// that read eval results from inside the tool would make that separation untrue
// by adjacency (prov-2026-c0f55691).
//
// It needs no build tag of its own, unlike the runner: reading committed JSONL
// takes milliseconds, calls no model, needs no endpoint, and is deterministic,
// so none of the reasons the runner is tagged apply here.
//
//	go run ./evals/cmd/history               # every case
//	go run ./evals/cmd/history todo-chi-defined
//	go run ./evals/cmd/history -dir evals/results
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/livecodelife/sedum/evals"
)

func main() {
	dir := flag.String("dir", "evals/results", "directory the results files live in")
	flag.Parse()

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	cases := flag.Args()
	if len(cases) == 0 {
		found, err := evals.Cases(*dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reading %s: %v\n", *dir, err)
			os.Exit(1)
		}
		if len(found) == 0 {
			fmt.Fprintf(os.Stderr, "no results under %s\n", *dir)
			os.Exit(1)
		}
		cases = found
	}

	for i, id := range cases {
		entries, err := evals.ReadEntries(*dir, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reading %s: %v\n", id, err)
			os.Exit(1)
		}
		if i > 0 {
			fmt.Fprintln(out)
		}
		evals.History(out, id, entries)
	}
}
