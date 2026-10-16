package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Orphans, deprecated records, stamping, and prompt context
// (prov-2026-cc73cf61).

func recordFile(id, scope, extra string) string {
	return "id: " + id + "\n" + extra + "intent: |\n  Something.\naffected_scope:\n" + scope
}

const userScope = "  - app/models/user.rb\n"

// twoOwners grows records a and b into one file, each owning one region.
func twoOwners(t *testing.T) Config {
	t.Helper()
	cfg := config(t)
	cfg.Records = writeTree(t, map[string]string{
		"prov-2026-aaaaaaaa.yml": recordFile("prov-2026-aaaaaaaa", userScope, ""),
		"prov-2026-bbbbbbbb.yml": recordFile("prov-2026-bbbbbbbb", userScope, ""),
	})
	cfg.Client = &stub{sequence: []string{fields("email"), fields("name")}}
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("first grow: %v", err)
	}
	return cfg
}

func writeRecord(t *testing.T, cfg Config, id, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(cfg.Records, id+".yml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func removeRecord(t *testing.T, cfg Config, id string) {
	t.Helper()
	if err := os.Remove(filepath.Join(cfg.Records, id+".yml")); err != nil {
		t.Fatal(err)
	}
}

func TestADeprecatedRecordIsNeverGrown(t *testing.T) {
	cfg := config(t)
	cfg.Records = writeTree(t, map[string]string{
		"prov-2026-aaaaaaaa.yml": recordFile("prov-2026-aaaaaaaa", userScope, "status: deprecated\n"),
	})
	client := &stub{response: fields("email")}
	cfg.Client = client
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("grow: %v", err)
	}
	if client.calls != 0 {
		t.Errorf("a deprecated record made %d model call(s)", client.calls)
	}

	cfg.Only = []string{"prov-2026-aaaaaaaa"}
	_, err := Run(context.Background(), cfg)
	if err == nil || !strings.Contains(err.Error(), "deprecated") {
		t.Errorf("--only naming a deprecated record = %v, want an error saying it is deprecated", err)
	}
}

func TestPruneOrphansRemovesRegionsOfMissingAndDeprecatedRecords(t *testing.T) {
	for _, retire := range []struct {
		name string
		do   func(t *testing.T, cfg Config)
	}{
		{"removed", func(t *testing.T, cfg Config) { removeRecord(t, cfg, "prov-2026-bbbbbbbb") }},
		{"deprecated", func(t *testing.T, cfg Config) {
			writeRecord(t, cfg, "prov-2026-bbbbbbbb", recordFile("prov-2026-bbbbbbbb", userScope, "status: deprecated\n"))
		}},
	} {
		t.Run(retire.name, func(t *testing.T) {
			cfg := twoOwners(t)
			retire.do(t, cfg)

			// Without the flag the orphan stays, and is reported.
			result, _ := rerun(t, cfg)
			if !strings.Contains(userModel(t, cfg), "attribute :name") {
				t.Fatalf("an orphan was pruned without --prune-orphans")
			}
			if len(result.Orphans) != 1 || result.Orphans[0].Record != "prov-2026-bbbbbbbb" || !result.Orphans[0].Authorized {
				t.Errorf("orphans = %+v, want b's region, in an authorized file", result.Orphans)
			}

			// --check fails on it.
			check := cfg
			check.Check = true
			result, _ = rerun(t, check)
			if len(result.Orphans) != 1 || !result.CheckFailed() {
				t.Errorf("--check: orphans = %+v, failed = %v; want the orphan and a failure", result.Orphans, result.CheckFailed())
			}

			prune := cfg
			prune.PruneOrphans = true
			result, _ = rerun(t, prune)
			got := userModel(t, cfg)
			if strings.Contains(got, "attribute :name") || !strings.Contains(got, "attribute :email") {
				t.Errorf("--prune-orphans should remove only b's region:\n%s", got)
			}
			if len(result.Orphans) != 1 || !result.Orphans[0].Pruned {
				t.Errorf("orphans = %+v, want one pruned", result.Orphans)
			}
		})
	}
}

// Sedum writes only authorized paths, so an orphan in a file no record
// authorizes any more is reported and left alone.
func TestAnOrphanInAnUnauthorizedFileIsOnlyReported(t *testing.T) {
	cfg := config(t)
	cfg.Records = writeTree(t, map[string]string{
		"prov-2026-aaaaaaaa.yml": recordFile("prov-2026-aaaaaaaa", userScope, ""),
		"prov-2026-bbbbbbbb.yml": recordFile("prov-2026-bbbbbbbb", "  - app/models/order.rb\n", ""),
	})
	orderField := `{"invocations": [{"action": "addField", "kwargs": {"model": "order", "field": "total"}}]}`
	cfg.Client = &stub{sequence: []string{fields("email"), orderField}}
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("first grow: %v", err)
	}
	orderPath := filepath.Join(cfg.Output, "app", "models", "order.rb")
	before := readFile(t, orderPath)

	removeRecord(t, cfg, "prov-2026-bbbbbbbb")
	prune := cfg
	prune.PruneOrphans = true
	result, _ := rerun(t, prune)

	if got := readFile(t, orderPath); got != before {
		t.Errorf("an unauthorized file was changed:\n%s", got)
	}
	if len(result.Orphans) != 1 || result.Orphans[0].Authorized || result.Orphans[0].Pruned {
		t.Errorf("orphans = %+v, want one reported in an unauthorized file and not pruned", result.Orphans)
	}
}

// stripHashes removes gen and pkg from every marker, as markers written before
// they existed look.
func stripHashes(t *testing.T, cfg Config) {
	t.Helper()
	path := filepath.Join(cfg.Output, "app", "models", "user.rb")
	content := userModel(t, cfg)
	for _, key := range []string{`"gen":"`, `"pkg":"`} {
		for {
			i := strings.Index(content, key)
			if i < 0 {
				break
			}
			end := strings.Index(content[i+len(key):], `"`) + i + len(key) + 1
			content = content[:i] + content[end+1:] // the value and its trailing comma
		}
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestStampWritesOnlyMissingHashes(t *testing.T) {
	cfg := grown(t)
	stamped := userModel(t, cfg)
	stripHashes(t, cfg)
	stripped := userModel(t, cfg)
	if strings.Contains(stripped, `"gen"`) {
		t.Fatalf("fixture still carries gen:\n%s", stripped)
	}

	// A dry run reports and writes nothing.
	dry := cfg
	dry.Stamp, dry.DryRun = true, true
	result, client := rerun(t, dry)
	if client.calls != 0 || len(result.Stamped) != 1 {
		t.Errorf("dry --stamp: calls=%d stamped=%+v; want 0 calls and one region reported", client.calls, result.Stamped)
	}
	if userModel(t, cfg) != stripped {
		t.Error("dry --stamp wrote")
	}

	stamp := cfg
	stamp.Stamp = true
	if _, client := rerun(t, stamp); client.calls != 0 {
		t.Errorf("--stamp made %d model call(s)", client.calls)
	}
	if got := userModel(t, cfg); got != stamped {
		t.Errorf("--stamp should restore exactly the hashes and nothing else:\n got %q\nwant %q", got, stamped)
	}

	// The record is now current.
	again, client := rerun(t, cfg)
	if client.calls != 0 || modeOf(t, again, "prov-2026-aaaaaaaa") != ModeSkip {
		t.Error("after --stamp the record is not current")
	}

	// A value already present is never overwritten, so stamping cannot hide a
	// real change.
	editRecord(t, cfg, "Something else.")
	stamp.Stamp = true
	if result, _ := rerun(t, stamp); len(result.Stamped) != 0 {
		t.Errorf("--stamp overwrote a gen that was present: %+v", result.Stamped)
	}
	check := cfg
	check.Check = true
	if result, _ := rerun(t, check); !result.CheckFailed() {
		t.Error("after an edit and --stamp, --check should still report the record stale")
	}
}

// The model is told which regions other records own in its files, and a
// re-selected record sees its own earlier selections, rebuilt from markers.
func TestTheSelectionPromptCarriesPresentRegionsAndRebuiltPriors(t *testing.T) {
	cfg := twoOwners(t)
	editRecord(t, cfg, "Something else.")

	client := &stub{response: fields("email")}
	only := cfg
	only.Only = []string{"prov-2026-aaaaaaaa"}
	only.Client = client
	if _, err := Run(context.Background(), only); err != nil {
		t.Fatalf("grow: %v", err)
	}
	if client.calls != 1 {
		t.Fatalf("calls = %d, want 1", client.calls)
	}
	var user string
	for _, m := range client.prompts[0] {
		user += m.Content
	}
	for _, want := range []string{"Already present", "prov-2026-bbbbbbbb", `"field":"name"`, `"field":"email"`} {
		if !strings.Contains(user, want) {
			t.Errorf("the prompt does not carry %q:\n%s", want, user)
		}
	}
	if strings.Contains(user, "sha256:") {
		t.Errorf("a hash reached the prompt:\n%s", user)
	}
}
