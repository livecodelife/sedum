package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
		{"record", "", "string", ""},
		{"execute", "", "string", ""},
		{"dry-run", "", "bool", "false"},
		{"stop-after", "", "string", ""},
		{"retries", "", "int", "3"},
		{"model", "", "string", ""},
		{"log", "", "string", ".sedum/run.log"},
		{"verbose", "v", "bool", "false"},
	})
}

func TestValidateFlagSurface(t *testing.T) {
	assertFlags(t, "validate", []flagSpec{
		{"generators", "", "string", ""},
		{"package", "", "stringSlice", "[]"},
		{"strict", "", "bool", "false"},
	})
}

func TestResolveFlagSurface(t *testing.T) {
	assertFlags(t, "resolve", []flagSpec{
		{"generators", "", "string", ""},
		{"records", "", "string", ""},
		{"lang", "", "stringSlice", "[]"},
		{"only", "", "stringSlice", "[]"},
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
	return out.String(), root.Execute()
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

// Combinations the PRD explicitly blesses must survive validation and fail only
// on the milestone that has not landed yet.
func TestLegalFlagCombinations(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"record composes with dry-run", []string{"grow", "--generators", "g", "--records", "r", "--record", "out.json", "--dry-run"}},
		{"execute composes with dry-run", []string{"grow", "--generators", "g", "--execute", "in.json", "--dry-run"}},
		{"execute makes records optional", []string{"grow", "--generators", "g", "--execute", "in.json"}},
		{"execute accepts records for scope validation", []string{"grow", "--generators", "g", "--records", "r", "--execute", "in.json"}},
		{"replay may stop after expansion without record", []string{"grow", "--generators", "g", "--execute", "in.json", "--stop-after", "expansion"}},
		{"stop-after files needs no record", []string{"grow", "--generators", "g", "--records", "r", "--stop-after", "files"}},
		{"repeated lang and only", []string{"grow", "--generators", "g", "--records", "r", "--lang", "rails", "--lang", "chi", "--only", "PR-1", "--only", "PR-2"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := exec(t, tc.args...)
			if err == nil {
				t.Fatalf("expected the not-implemented error, got nil")
			}
			if !strings.Contains(err.Error(), "not implemented") {
				t.Fatalf("flag validation rejected a legal combination: %v", err)
			}
		})
	}
}

// Until a milestone lands its command must fail loudly and name the milestone,
// so a user never mistakes an unimplemented phase for a no-op run.
func TestUnimplementedCommandsNameTheirMilestone(t *testing.T) {
	tests := []struct {
		args      []string
		milestone string
	}{
		{[]string{"validate", "--generators", "g"}, "M1"},
		{[]string{"resolve", "--generators", "g", "--records", "r"}, "M3"},
		{[]string{"actions", "--generators", "g", "--package", "rails"}, "M6"},
		{[]string{"grow", "--generators", "g", "--records", "r"}, "M6"},
	}

	for _, tc := range tests {
		t.Run(tc.args[0], func(t *testing.T) {
			_, err := exec(t, tc.args...)
			wantErr(t, err, "not implemented", tc.milestone)
		})
	}
}

func TestRootListsEverySubcommand(t *testing.T) {
	want := map[string]bool{"grow": false, "validate": false, "resolve": false, "actions": false}
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
