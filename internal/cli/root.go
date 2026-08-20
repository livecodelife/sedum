// Package cli is Sedum's command surface.
//
// Commands parse flags, check the interdependence rules the flags cannot
// express on their own, build a config struct, and delegate. Everything else
// lives behind the internal packages the pipeline is built from.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/livecodelife/sedum/internal/recording"
)

const (
	defaultOutput  = "."
	defaultLogPath = ".sedum/run.log"
	defaultRetries = 3
)

// NewRootCommand builds a fresh command tree. Nothing is shared between trees,
// so callers and tests never contend over global state.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "sedum",
		Short: "Generate boilerplate from provenance records using generator packages you author",
		Long: `Sedum generates boilerplate code from provenance records using generator
packages that teams author themselves.

A provenance record declares intent, constraints, and the files a change is
authorized to touch. A generator package declares your conventions: the shape of
each file type, which code-injection actions exist, what arguments they take,
what templates they render, and where in a file the results belong.

Sedum's core contains no language-specific knowledge. Adding a language or
framework means authoring a generator package, never modifying Sedum.`,

		// A caller told to invoke a binary must be able to check it is the one
		// the instruction was written for (prov-2026-b5465dfa). The value is
		// the recording's, so the surface and the artifact cannot disagree.
		Version: recording.Version,

		// Errors are reported once, by Execute, with a consistent prefix. Usage
		// is not dumped on every failure: a specific diagnostic buried under
		// forty lines of flag help is a worse diagnostic.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// A bare semver, not cobra's "sedum version X" sentence. The consumer this
	// exists for is checking a version floor before invoking a command
	// (prov-2026-b5465dfa), and a string it has to strip a prefix off is a
	// parsing rule stated nowhere.
	root.SetVersionTemplate("{{.Version}}\n")

	root.AddCommand(
		newGrowCommand(),
		newValidateCommand(),
		newResolveCommand(),
		newActionsCommand(),
		newRenderCommand(),
	)

	return root
}

// Execute runs the CLI and returns the process exit code.
func Execute() int {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "sedum:", err)
		return 1
	}
	return 0
}
