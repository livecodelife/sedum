package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"regexp"

	"github.com/livecodelife/sedum/internal/recording"
)

// The flag surface is contract, not convenience: PRD.md's CLI tables are
// reproduced here so that adding or dropping a flag fails a test rather than
// silently drifting from the document.

type flagSpec struct {
	name      string
	shorthand string
	typ       string // pflag's Value.Type(): string, bool, int, stringSlice
	def       string
}

func findCommand(t *testing.T, name string) *cobra.Command {
	t.Helper()
	cmd, _, err := NewRootCommand().Find([]string{name})
	if err != nil {
		t.Fatalf("subcommand %q not found: %v", name, err)
	}
	if cmd.Name() != name {
		t.Fatalf("Find(%q) resolved to %q", name, cmd.Name())
	}
	return cmd
}

func assertFlags(t *testing.T, cmdName string, want []flagSpec) {
	t.Helper()
	cmd := findCommand(t, cmdName)

	got := map[string]flagSpec{}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == "help" {
			return
		}
		got[f.Name] = flagSpec{f.Name, f.Shorthand, f.Value.Type(), f.DefValue}
	})

	for _, w := range want {
		g, ok := got[w.name]
		if !ok {
			t.Errorf("%s: missing flag --%s", cmdName, w.name)
			continue
		}
		if g != w {
			t.Errorf("%s: flag --%s = %+v, want %+v", cmdName, w.name, g, w)
		}
		delete(got, w.name)
	}
	for name := range got {
		t.Errorf("%s: undocumented flag --%s (PRD.md lists no such flag)", cmdName, name)
	}
}

func TestGrowFlagSurface(t *testing.T) {
	assertFlags(t, "grow", []flagSpec{
		{"generators", "", "string", ""},
		{"records", "", "string", ""},
		{"output", "", "string", "."},
		{"lang", "", "stringSlice", "[]"},
		{"only", "", "stringSlice", "[]"},
		{"var", "", "stringArray", "[]"},
		{"record", "", "string", ""},
		{"execute", "", "string", ""},
		{"dry-run", "", "bool", "false"},
		{"force", "", "bool", "false"},
		{"check", "", "bool", "false"},
		{"prune-orphans", "", "bool", "false"},
		{"restamp-gen", "", "bool", "false"},
		{"stamp", "", "bool", "false"},
		{"json", "", "bool", "false"},
		{"stop-after", "", "string", ""},
		{"retries", "", "int", "3"},
		{"model", "", "string", ""},
		{"local-model", "", "string", ""},
		{"local-model-server", "", "string", ""},
		{"local-model-backend", "", "string", ""},
		{"log", "", "string", ".sedum/run.log"},
		{"verbose", "v", "bool", "false"},
		{"full-catalog", "", "bool", "false"},
		// prov-2026-91c54941: opt-in, default off. Off, the request sedum
		// sends is byte-for-byte what it is today (internal/selection's
		// TestCompleteWithoutSchemaMatchesTodaysRequestByteForByte).
		{"response-schema", "", "bool", "false"},
		// prov-2026-e4bab53b and prov-2026-8602b702: caller-supplied context
		// for Phase 4's prompt. All three default to absent, and absent, the
		// prompt is byte-for-byte what it was.
		{"prompt", "", "string", ""},
		{"prompt-file", "", "string", ""},
		{"prior-recording", "", "string", ""},
	})
}

func TestValidateFlagSurface(t *testing.T) {
	assertFlags(t, "validate", []flagSpec{
		{"generators", "", "string", ""},
		{"package", "", "stringSlice", "[]"},
		{"strict", "", "bool", "false"},
	})
}

// sedum regions is read-only and needs no records (prov-2026-bf01cd6b).
func TestRegionsFlagSurface(t *testing.T) {
	assertFlags(t, "regions", []flagSpec{
		{"generators", "", "string", ""},
		{"output", "", "string", "."},
		{"records", "", "string", ""},
		{"record", "", "stringArray", "[]"},
		{"json", "", "bool", "false"},
	})
}

func TestResolveFlagSurface(t *testing.T) {
	assertFlags(t, "resolve", []flagSpec{
		{"generators", "", "string", ""},
		{"records", "", "string", ""},
		{"lang", "", "stringSlice", "[]"},
		{"only", "", "stringSlice", "[]"},
		{"var", "", "stringArray", "[]"},
		{"show-template", "", "bool", "false"},
	})
}

func TestActionsFlagSurface(t *testing.T) {
	assertFlags(t, "actions", []flagSpec{
		{"generators", "", "string", ""},
		{"package", "", "string", ""},
		{"all", "", "bool", "false"},
		{"json", "", "bool", "false"},
	})
}

func exec(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := NewRootCommand()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(args)
	// Execute first: as return operands these evaluate left to right, so
	// reading the buffer inline would always capture it empty.
	err := root.Execute()
	return out.String(), err
}

// wantErr asserts the command failed and that the message names every needle.
// A diagnostic that does not name the flags involved is a defect under M0's
// "every error names the artifact it concerns" constraint.
func wantErr(t *testing.T, err error, needles ...string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	msg := err.Error()
	for _, n := range needles {
		if !strings.Contains(msg, n) {
			t.Errorf("error %q does not mention %q", msg, n)
		}
	}
}

func TestFlagInterdependence(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		needles []string
	}{
		{
			name:    "record and execute are mutually exclusive",
			args:    []string{"grow", "--generators", "g", "--records", "r", "--record", "out.json", "--execute", "in.json"},
			needles: []string{"record", "execute"},
		},
		{
			name:    "stop-after invocations requires record",
			args:    []string{"grow", "--generators", "g", "--records", "r", "--stop-after", "invocations"},
			needles: []string{"invocations", "--record"},
		},
		{
			name:    "stop-after expansion requires record",
			args:    []string{"grow", "--generators", "g", "--records", "r", "--stop-after", "expansion"},
			needles: []string{"expansion", "--record"},
		},
		{
			name:    "stop-after rejects an unknown phase",
			args:    []string{"grow", "--generators", "g", "--records", "r", "--stop-after", "injection"},
			needles: []string{"injection", "resolution", "files", "invocations", "expansion"},
		},
		{
			name:    "replay does not run resolution so it cannot stop after it",
			args:    []string{"grow", "--generators", "g", "--execute", "in.json", "--stop-after", "resolution"},
			needles: []string{"resolution", "--execute"},
		},
		{
			name:    "records is required without execute",
			args:    []string{"grow", "--generators", "g"},
			needles: []string{"--records"},
		},
		{
			name:    "generators is always required",
			args:    []string{"grow", "--records", "r"},
			needles: []string{"generators"},
		},
		{
			name:    "validate requires generators",
			args:    []string{"validate"},
			needles: []string{"generators"},
		},
		{
			name:    "resolve requires records",
			args:    []string{"resolve", "--generators", "g"},
			needles: []string{"records"},
		},
		{
			name:    "actions requires package",
			args:    []string{"actions", "--generators", "g"},
			needles: []string{"package"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := exec(t, tc.args...)
			wantErr(t, err, tc.needles...)
		})
	}
}

// Combinations the PRD explicitly blesses must survive validation. What they
// fail on afterwards says how far they got: a combination whose phases have all
// landed reaches Phase 0 and fails on the generators directory these cases name
// but do not create, and every other combination names the milestone it is
// waiting on.
func TestLegalFlagCombinations(t *testing.T) {
	const (
		reachesThePipeline = "read generators directory"
		// A run that will reach Phase 4 builds its model client before any
		// phase runs, so this fragment proves the flags were accepted and the
		// run got as far as needing a model - falling back to a bundled
		// default this test binary was never built to carry.
		reachesTheModel = "no bundled model found"
		// A replay builds no model client, so it gets as far as reading the
		// recording. Failing to open a file these cases name but do not create
		// is what proves the flags were accepted.
		reachesTheRecording = "opening the recording"
	)

	tests := []struct {
		name string
		args []string
		// fails is the fragment that proves validation passed.
		fails string
	}{
		{"record composes with dry-run", []string{"grow", "--generators", "g", "--records", "r", "--record", "out.json", "--dry-run"}, reachesTheModel},
		{"execute composes with dry-run", []string{"grow", "--generators", "g", "--execute", "in.json", "--dry-run"}, reachesTheRecording},
		{"execute makes records optional", []string{"grow", "--generators", "g", "--execute", "in.json"}, reachesTheRecording},
		{"execute accepts records for scope validation", []string{"grow", "--generators", "g", "--records", "r", "--execute", "in.json"}, reachesTheRecording},
		{"replay may stop after expansion without record", []string{"grow", "--generators", "g", "--execute", "in.json", "--stop-after", "expansion"}, reachesTheRecording},
		{"stop-after files needs no record", []string{"grow", "--generators", "g", "--records", "r", "--stop-after", "files"}, reachesThePipeline},
		{"stop-after resolution needs no record", []string{"grow", "--generators", "g", "--records", "r", "--stop-after", "resolution"}, reachesThePipeline},
		{"repeated lang and only", []string{"grow", "--generators", "g", "--records", "r", "--lang", "rails", "--lang", "chi", "--only", "PR-1", "--only", "PR-2"}, reachesTheModel},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// A run that gets past flag validation opens its run log, so it
			// is pointed somewhere disposable rather than at the working
			// directory these tests run in.
			args := append(tc.args, "--log", filepath.Join(t.TempDir(), "run.log"))

			_, err := exec(t, args...)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tc.fails) {
				t.Fatalf("flag validation rejected a legal combination: %v", err)
			}
		})
	}
}

// Every milestone has landed, so no command names one it is waiting on. The
// test that asserted otherwise was retired with M7 rather than left asserting
// over an empty table (prov-2026-dc227be7).

// validate is the one command a milestone has landed, so it is exercised end
// to end here: the checks themselves are genpkg's to test, but that the command
// reaches them and reports what they found is this package's.
func TestValidateReportsOnTheGeneratorsDirectory(t *testing.T) {
	out, err := exec(t, "validate", "--generators", filepath.Join("..", "..", "testdata", "generators"))
	if err != nil {
		t.Fatalf("validate on the fixture packages failed: %v\n%s", err, out)
	}
	for _, want := range []string{"4 package(s) loaded", "0 error(s)"} {
		if !strings.Contains(out, want) {
			t.Errorf("validate output = %q, does not report %q", out, want)
		}
	}
}

func TestValidateFailsOnARejectedPackage(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "broken"), 0o755); err != nil {
		t.Fatal(err)
	}

	out, err := exec(t, "validate", "--generators", dir)
	wantErr(t, err, "broken")
	if !strings.Contains(out, "manifest_missing") {
		t.Errorf("validate output = %q, does not name the rule violated", out)
	}
}

func TestRootListsEverySubcommand(t *testing.T) {
	want := map[string]bool{"grow": false, "validate": false, "resolve": false, "actions": false, "regions": false}
	for _, c := range NewRootCommand().Commands() {
		if _, ok := want[c.Name()]; ok {
			want[c.Name()] = true
		}
		if c.Short == "" {
			t.Errorf("subcommand %q has no Short description", c.Name())
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("root is missing subcommand %q", name)
		}
	}
}

// --var is name=value, split on the first = only: a Go module path or a .NET
// namespace is a plausible value and may contain one, while a name may not.
func TestVarParsing(t *testing.T) {
	got, err := parseVars([]string{"module=github.com/acme/todo", "registry=ghcr.io"})
	if err != nil {
		t.Fatalf("parseVars: %v", err)
	}
	if got["module"] != "github.com/acme/todo" || got["registry"] != "ghcr.io" {
		t.Errorf("parsed %v", got)
	}

	if _, err := parseVars([]string{"module"}); err == nil {
		t.Error("a --var with no = was accepted")
	}
	if _, err := parseVars([]string{"=todo"}); err == nil {
		t.Error("a --var with no name was accepted")
	}

	// Two values for one variable means one of them was meant and nothing says
	// which, so it is refused rather than resolved by order.
	if _, err := parseVars([]string{"module=a", "module=b"}); err == nil {
		t.Error("a repeated --var silently took one of the two values")
	}
}

// A caller told to invoke a binary must be able to check it is the one the
// instruction was written for, and the value it reads must be the one a
// recording carries - otherwise the two drift and the check is theatre
// (prov-2026-b5465dfa).

func TestVersionEmitsABareSemver(t *testing.T) {
	out, err := exec(t, "--version")
	if err != nil {
		t.Fatalf("--version: %v\n%s", err, out)
	}

	got := strings.TrimSpace(out)
	if got != recording.Version {
		t.Errorf("--version emitted %q, which is not the value a recording carries (%q)", got, recording.Version)
	}
	// A prefix would be a parsing rule stated nowhere.
	if !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(got) {
		t.Errorf("--version emitted %q, which is not a bare semver", got)
	}
}

// prov-2026-e4bab53b and prov-2026-8602b702.
//
// Every case here points --generators and --records at directories that do not
// exist. A flag problem is reported before any phase runs, so it fails with its
// own message; a run whose flags were all accepted fails further on, at
// reachesTheModel, which is how a case proves its flags were legal.

// execIn is exec with a stdin, for --prompt-file -. It assumes the command reads
// the input through cobra (cmd.InOrStdin) rather than os.Stdin directly.
func execIn(t *testing.T, stdin string, args ...string) (string, error) {
	t.Helper()
	root := NewRootCommand()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetIn(io.NopCloser(strings.NewReader(stdin)))
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), err
}

func growArgs(t *testing.T, extra ...string) []string {
	t.Helper()
	args := []string{"grow", "--generators", "g", "--records", "r", "--log", filepath.Join(t.TempDir(), "run.log")}
	return append(args, extra...)
}

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// wantReachesTheModel asserts a run's flags were all accepted.
func wantReachesTheModel(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "no bundled model found") {
		t.Fatalf("flag validation rejected a legal combination: %v", err)
	}
}

func TestPromptAndPromptFileAreMutuallyExclusive(t *testing.T) {
	file := writeTemp(t, "p.txt", "x")
	_, err := exec(t, growArgs(t, "--prompt", "x", "--prompt-file", file)...)
	wantErr(t, err, "prompt", "prompt-file")
	// Not the flag-parsing failure of a flag that does not exist.
	if strings.Contains(err.Error(), "unknown flag") {
		t.Errorf("--prompt or --prompt-file does not exist yet: %v", err)
	}
}

func TestPromptTextIsAcceptedAndReachesTheModel(t *testing.T) {
	for name, args := range map[string][]string{
		"inline text":          {"--prompt", "FAIL: undefined method"},
		"whitespace only":      {"--prompt", " \n\t "},
		"exactly 16384 B":      {"--prompt", strings.Repeat("a", 16384)},
		"combined with --only": {"--prompt", "x", "--only", "PR-1"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := exec(t, growArgs(t, args...)...)
			wantReachesTheModel(t, err)
		})
	}
}

func TestPromptOverTheCapIsAnErrorNamingTheLimitAndTheSize(t *testing.T) {
	_, err := exec(t, growArgs(t, "--prompt", strings.Repeat("a", 20000))...)
	wantErr(t, err, "16384", "20000")
	if strings.Contains(err.Error(), "unknown flag") {
		t.Errorf("--prompt does not exist yet: %v", err)
	}

	// The limit is bytes, not characters: 8200 two-byte runes is 16400 bytes.
	_, err = exec(t, growArgs(t, "--prompt", strings.Repeat("\u00e9", 8200))...)
	wantErr(t, err, "16384", "16400")

	_, err = exec(t, growArgs(t, "--prompt", strings.Repeat("a", 16385))...)
	wantErr(t, err, "16384", "16385")
}

func TestPromptFileIsReadFromAPath(t *testing.T) {
	ok := writeTemp(t, "ok.txt", "FAIL: validator output\n")
	_, err := exec(t, growArgs(t, "--prompt-file", ok)...)
	wantReachesTheModel(t, err)

	big := writeTemp(t, "big.txt", strings.Repeat("a", 20000))
	_, err = exec(t, growArgs(t, "--prompt-file", big)...)
	wantErr(t, err, "16384", "20000")

	missing := filepath.Join(t.TempDir(), "does-not-exist.txt")
	_, err = exec(t, growArgs(t, "--prompt-file", missing)...)
	wantErr(t, err, missing)

	blank := writeTemp(t, "blank.txt", "\n  \n")
	_, err = exec(t, growArgs(t, "--prompt-file", blank)...)
	wantReachesTheModel(t, err)
}

func TestPromptFileDashReadsStdin(t *testing.T) {
	_, err := execIn(t, "FAIL from stdin\n", growArgs(t, "--prompt-file", "-")...)
	wantReachesTheModel(t, err)

	_, err = execIn(t, strings.Repeat("a", 20000), growArgs(t, "--prompt-file", "-")...)
	wantErr(t, err, "16384", "20000")
}

func TestPriorRecordingIsAcceptedWhenItHoldsRecords(t *testing.T) {
	rec := writeTemp(t, "prior.json", `{"sedum_version":"0.7.0","packages":{},"records":[`+
		`{"record_id":"PR-1","files":[],"phases":[{"name":"default","invocations":[{"action":"a","kwargs":{}}]}]}]}`)
	_, err := exec(t, growArgs(t, "--prior-recording", rec)...)
	wantReachesTheModel(t, err)

	// A sedum_version nobody has heard of is not consulted for compatibility.
	odd := writeTemp(t, "odd.json", `{"sedum_version":"99.0.0","packages":{},"records":[{"record_id":"PR-1","files":[],"phases":[]}]}`)
	_, err = exec(t, growArgs(t, "--prior-recording", odd)...)
	wantReachesTheModel(t, err)

	// Composes with --prompt.
	_, err = exec(t, growArgs(t, "--prior-recording", rec, "--prompt", "x")...)
	wantReachesTheModel(t, err)
}

// A retry run that silently lacked its context would look like a normal run, so
// each of these is an error naming the file, raised before any phase runs: the
// generators directory these cases name does not exist, and its own error would
// not name the recording.
func TestPriorRecordingThatCannotBeUsedIsAnEarlyErrorNamingThePath(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.json")
	tests := map[string]string{
		"missing":        missing,
		"not json":       writeTemp(t, "prior.json", "this is not json"),
		"empty file":     writeTemp(t, "empty.json", ""),
		"no records":     writeTemp(t, "norecords.json", `{"sedum_version":"0.7.0","packages":{},"records":[]}`),
		"records absent": writeTemp(t, "absent.json", `{"sedum_version":"0.7.0"}`),
		"empty object":   writeTemp(t, "obj.json", `{}`),
	}
	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := exec(t, growArgs(t, "--prior-recording", path)...)
			wantErr(t, err, path)
			if strings.Contains(err.Error(), "unknown flag") {
				t.Errorf("--prior-recording does not exist yet: %v", err)
			}
			if strings.Contains(err.Error(), "read generators directory") {
				t.Errorf("the run got as far as Phase 0 before rejecting the recording: %v", err)
			}
		})
	}
}

// Replay invokes no model, so neither flag has anything to do and neither is
// consulted: an unreadable file or an oversized text is not an error there.
func TestPromptAndPriorRecordingAreIgnoredUnderExecute(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")

	tests := []struct {
		name   string
		args   []string
		notice string
	}{
		{"prompt", []string{"--prompt", "x"}, "--prompt (replay invokes no model"},
		{"oversized prompt", []string{"--prompt", strings.Repeat("a", 20000)}, "--prompt (replay invokes no model"},
		{"prompt-file", []string{"--prompt-file", missing}, "--prompt-file (replay invokes no model"},
		{"prior-recording", []string{"--prior-recording", missing}, "--prior-recording (replay invokes no model"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"grow", "--generators", "g", "--execute", "in.json",
				"--log", filepath.Join(t.TempDir(), "run.log")}, tc.args...)
			out, err := exec(t, args...)
			if err == nil || !strings.Contains(err.Error(), "opening the recording") {
				t.Fatalf("replay did not get as far as the recording it was told to read: %v", err)
			}
			if !strings.Contains(out, "ignoring "+tc.notice) {
				t.Errorf("output does not list the flag as ignored with its reason (want %q):\n%s", tc.notice, out)
			}
		})
	}
}

// prov-2026-d2f1bdf2: --model and --local-model take precedence over the
// bundled default, whether or not this build carries one. BundledLocalConfig
// resolves paths from os.Executable(), so a bundled default is staged next to
// this test binary the same way selection's own test does it.
//
// What proves precedence is whose name appears in the failure: the bundled
// paths are stand-in files that could never serve a model, so a run that chose
// them would fail naming them.
func TestModelFlagsOverrideTheBundledDefault(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Dir(exe)
	bundledServer := filepath.Join(dir, "goinfer-serve")
	bundledModel := filepath.Join(dir, "sedum-default-model.gguf")

	stage := func(t *testing.T, present bool) {
		t.Helper()
		for _, p := range []string{bundledServer, bundledModel} {
			os.Remove(p)
			if present {
				if err := os.WriteFile(p, []byte("stand-in for a test"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
		}
		t.Cleanup(func() { os.Remove(bundledServer); os.Remove(bundledModel) })
	}

	for _, present := range []bool{true, false} {
		name := "bundled absent"
		if present {
			name = "bundled present"
		}
		t.Run(name+"/--model", func(t *testing.T) {
			stage(t, present)
			t.Setenv("OPENAI_BASE_URL", "")
			t.Setenv("OPENAI_API_KEY", "")

			_, err := exec(t, growArgs(t, "--model", "some-model")...)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			// It tried to build a hosted client for the named model.
			if !strings.Contains(err.Error(), "OPENAI_BASE_URL") {
				t.Errorf("--model did not take the hosted path: %v", err)
			}
			if strings.Contains(err.Error(), "no bundled model found") || strings.Contains(err.Error(), bundledModel) {
				t.Errorf("the bundled default was consulted despite --model: %v", err)
			}
		})
		t.Run(name+"/--local-model", func(t *testing.T) {
			stage(t, present)
			userServer := filepath.Join(t.TempDir(), "my-serve")
			userModel := filepath.Join(t.TempDir(), "mine.gguf")

			_, err := exec(t, growArgs(t, "--local-model", userModel, "--local-model-server", userServer)...)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), userServer) {
				t.Errorf("--local-model did not use the server it was given: %v", err)
			}
			if strings.Contains(err.Error(), bundledServer) || strings.Contains(err.Error(), "no bundled model found") {
				t.Errorf("the bundled default was consulted despite --local-model: %v", err)
			}
		})
	}
}

// --check only reports, so every flag that writes, records or forces a
// selection is refused with it, naming both flags (prov-2026-9621daed).
// --execute is the exception: with it --check is a dry replay that fails when
// replaying would change the tree (prov-2026-c637c262).
func TestCheckRefusesWhatContradictsIt(t *testing.T) {
	for _, extra := range [][]string{
		{"--record", "rec.json"},
		{"--dry-run"},
		{"--force"},
		{"--stop-after", "files"},
	} {
		_, err := exec(t, growArgs(t, append([]string{"--check"}, extra...)...)...)
		wantErr(t, err, "--check", extra[0])
	}
}

// The re-render limit is stated where a user deciding whether to trust
// --check or reach for --force will read it.
func TestCheckAndForceHelpNameTheReRenderLimit(t *testing.T) {
	cmd := findCommand(t, "grow")
	for _, name := range []string{"check", "force"} {
		usage := cmd.Flags().Lookup(name).Usage
		for _, want := range []string{"verbatim", "default", "--var", "--force"} {
			if !strings.Contains(usage, want) {
				t.Errorf("--%s help does not mention %q: %s", name, want, usage)
			}
		}
	}
}

// --stamp writes only missing hashes, so it refuses everything that selects,
// replays, records or checks (prov-2026-cc73cf61).
func TestStampRefusesWhatContradictsIt(t *testing.T) {
	for _, extra := range [][]string{
		{"--check"},
		{"--execute", "rec.json"},
		{"--record", "rec.json"},
		{"--force"},
	} {
		_, err := exec(t, growArgs(t, append([]string{"--stamp"}, extra...)...)...)
		wantErr(t, err, "--stamp", extra[0])
	}
}
