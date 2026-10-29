package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// A check run that skips scope paths no package claims also skips an existing
// file its file template cannot reconcile: the file is left alone, listed, and
// the record that owns it is unknown (prov-2026-9621daed).

const unreconcilableUser = "class User\n  def slug\n  end\n" +
	"  # sedum:addField {\"tier\":\"owned\",\"record\":\"prov-2026-aaaaaaaa\"}\n" +
	"  # /sedum:addField\nend\n"

func writeOutput(t *testing.T, cfg Config, rel, content string) {
	t.Helper()
	full := filepath.Join(cfg.Output, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func unreconciledConfig(t *testing.T, records map[string]string) Config {
	t.Helper()
	cfg := config(t)
	cfg.Records = recordsDir(t, records)
	cfg.Check = true
	cfg.SkipUnclaimedScope = true
	return cfg
}

func TestCheckSkipsAnUnreconcilableFileAndListsIt(t *testing.T) {
	cfg := unreconciledConfig(t, map[string]string{skipRecordA: "  - app/models/user.rb\n"})
	writeOutput(t, cfg, "app/models/user.rb", unreconcilableUser)

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if want := []string{"app/models/user.rb"}; !reflect.DeepEqual(result.SkippedUnreconciledPaths, want) {
		t.Errorf("SkippedUnreconciledPaths = %v, want %v", result.SkippedUnreconciledPaths, want)
	}
	got, err := os.ReadFile(filepath.Join(cfg.Output, "app/models/user.rb"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != unreconcilableUser {
		t.Errorf("the skipped file was modified:\n%s", got)
	}
}

// The markers of the file say the record is stale; the skipped file makes it
// unknown whatever they say.
func TestTheOwnerOfAnUnreconcilableFileIsUnknownWhateverItsMarkersSay(t *testing.T) {
	cfg := unreconciledConfig(t, map[string]string{skipRecordA: "  - app/models/user.rb\n"})
	writeOutput(t, cfg, "app/models/user.rb", unreconcilableUser)

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if got := modeOf(t, result, skipRecordA); got != ModeSelect {
		t.Errorf("mode = %v, want %v (unknown)", got, ModeSelect)
	}

	// The same markers in a file that reconciles are read as they are: the
	// marker has no gen, so the record is stale.
	reconcilable := "class User\n" +
		"  # sedum:addField {\"tier\":\"owned\",\"record\":\"prov-2026-aaaaaaaa\"}\n" +
		"  # /sedum:addField\n  # sedum:anchor:class_body\nend\n"
	writeOutput(t, cfg, "app/models/user.rb", reconcilable)
	result, err = Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if got := modeOf(t, result, skipRecordA); got != ModeReselect {
		t.Errorf("control: mode = %v, want %v (stale) for a file that reconciles", got, ModeReselect)
	}
	if len(result.SkippedUnreconciledPaths) != 0 {
		t.Errorf("control: SkippedUnreconciledPaths = %v, want none", result.SkippedUnreconciledPaths)
	}
}

func TestAnotherRecordIsClassifiedAsWithoutTheSkip(t *testing.T) {
	cfg := unreconciledConfig(t, map[string]string{
		skipRecordA: "  - app/models/user.rb\n",
		skipRecordB: "  - app/models/post.rb\n",
	})
	writeOutput(t, cfg, "app/models/user.rb", unreconcilableUser)
	writeOutput(t, cfg, "app/models/post.rb", "class Post\n  # sedum:anchor:class_body\nend\n")

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if want := []string{"app/models/user.rb"}; !reflect.DeepEqual(result.SkippedUnreconciledPaths, want) {
		t.Errorf("SkippedUnreconciledPaths = %v, want %v", result.SkippedUnreconciledPaths, want)
	}
	if got := modeOf(t, result, skipRecordB); got != ModeSelect {
		t.Errorf("the other record's mode = %v, want %v as without the skip", got, ModeSelect)
	}
}

func TestEachUnreconcilablePathIsListedOnceAndSorted(t *testing.T) {
	cfg := unreconciledConfig(t, map[string]string{
		skipRecordA: "  - app/models/user.rb\n  - app/models/post.rb\n",
		skipRecordB: "  - app/models/post.rb\n",
	})
	writeOutput(t, cfg, "app/models/user.rb", unreconcilableUser)
	writeOutput(t, cfg, "app/models/post.rb", "class Post\n  def slug\n  end\nend\n")

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if want := []string{"app/models/post.rb", "app/models/user.rb"}; !reflect.DeepEqual(result.SkippedUnreconciledPaths, want) {
		t.Errorf("SkippedUnreconciledPaths = %v, want each path once, sorted: %v", result.SkippedUnreconciledPaths, want)
	}
}

// A check that does not skip unclaimed scope, and a run that is not a check,
// both halt on such a file as they always did.
func TestWithoutTheToleranceAnUnreconcilableFileStillHaltsTheRun(t *testing.T) {
	cases := []struct {
		name   string
		adjust func(*Config)
	}{
		{"a check that does not skip unclaimed scope", func(c *Config) { c.SkipUnclaimedScope = false }},
		{"a run that is not a check", func(c *Config) { c.Check = false }},
	}
	for _, c := range cases {
		cfg := unreconciledConfig(t, map[string]string{skipRecordA: "  - app/models/user.rb\n"})
		writeOutput(t, cfg, "app/models/user.rb", unreconcilableUser)
		c.adjust(&cfg)

		if _, err := Run(context.Background(), cfg); err == nil {
			t.Errorf("%s: an unreconcilable file did not halt the run", c.name)
		}
	}
}
