package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Hash-aware grow (prov-2026-9621daed): each record is skipped, re-rendered or
// re-selected by comparing its markers' gen and pkg with the record and the
// package as they are now.

// snapshot is every file under root with its bytes and modification time.
type snapshot map[string]struct {
	content string
	mtime   time.Time
}

func snap(t *testing.T, root string) snapshot {
	t.Helper()
	out := snapshot{}
	for _, rel := range tree(t, root) {
		full := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Stat(full)
		if err != nil {
			t.Fatal(err)
		}
		out[rel] = struct {
			content string
			mtime   time.Time
		}{readFile(t, full), info.ModTime()}
	}
	return out
}

// ageTree moves every file's modification time into the past, so that any
// write - even of identical bytes - shows up as a changed mtime.
func ageTree(t *testing.T, root string) {
	t.Helper()
	past := time.Now().Add(-time.Hour)
	for _, rel := range tree(t, root) {
		if err := os.Chtimes(filepath.Join(root, filepath.FromSlash(rel)), past, past); err != nil {
			t.Fatal(err)
		}
	}
}

func sameSnapshot(t *testing.T, before, after snapshot) {
	t.Helper()
	if len(before) != len(after) {
		t.Errorf("file count changed: %d before, %d after", len(before), len(after))
	}
	for rel, b := range before {
		a, ok := after[rel]
		if !ok {
			t.Errorf("%s disappeared", rel)
			continue
		}
		if a.content != b.content {
			t.Errorf("%s changed:\n before %q\n after  %q", rel, b.content, a.content)
		}
		if !a.mtime.Equal(b.mtime) {
			t.Errorf("%s was rewritten (mtime %v -> %v)", rel, b.mtime, a.mtime)
		}
	}
}

// grown runs a first grow with one record selecting the email field.
func grown(t *testing.T) Config {
	t.Helper()
	cfg := config(t)
	cfg.Client = &stub{response: fields("email")}
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatalf("first grow: %v", err)
	}
	return cfg
}

func rerun(t *testing.T, cfg Config) (*Result, *stub) {
	t.Helper()
	client := &stub{response: fields("email")}
	cfg.Client = client
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("grow: %v", err)
	}
	return result, client
}

func modeOf(t *testing.T, result *Result, recordID string) Mode {
	t.Helper()
	for _, c := range result.Classifications {
		if c.RecordID == recordID {
			return c.Mode
		}
	}
	t.Fatalf("no classification for %s in %+v", recordID, result.Classifications)
	return 0
}

func editRecord(t *testing.T, cfg Config, intent string) {
	t.Helper()
	body := "id: prov-2026-aaaaaaaa\nintent: |\n  " + intent + "\naffected_scope:\n  - app/models/user.rb\n"
	if err := os.WriteFile(filepath.Join(cfg.Records, "prov-2026-aaaaaaaa.yml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func editTemplate(t *testing.T, cfg Config, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(cfg.Generators, "rails", "actions", "addField.rb"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Phase 2's done-check: a second grow over unchanged records makes no model
// calls and edits no files.
func TestASecondGrowOverUnchangedRecordsDoesNothing(t *testing.T) {
	cfg := grown(t)
	ageTree(t, cfg.Output)
	before := snap(t, cfg.Output)

	result, client := rerun(t, cfg)
	if client.calls != 0 {
		t.Errorf("an unchanged record made %d model call(s)", client.calls)
	}
	if got := modeOf(t, result, "prov-2026-aaaaaaaa"); got != ModeSkip {
		t.Errorf("mode = %v, want skip", got)
	}
	sameSnapshot(t, before, snap(t, cfg.Output))
}

func TestAChangedRecordIsReselected(t *testing.T) {
	cfg := grown(t)
	editRecord(t, cfg, "Something else.")

	result, client := rerun(t, cfg)
	if client.calls != 1 {
		t.Errorf("a changed record made %d model call(s), want 1", client.calls)
	}
	if got := modeOf(t, result, "prov-2026-aaaaaaaa"); got != ModeReselect {
		t.Errorf("mode = %v, want re-select", got)
	}
}

// A package change alone re-renders from the marker: no model call, the new
// template's output in place, the record's gen untouched.
func TestAChangedTemplateReRendersWithoutAModelCall(t *testing.T) {
	cfg := grown(t)
	genBefore := markerGenIn(t, userModel(t, cfg))
	editTemplate(t, cfg, "attribute :{{field}}, :string\n")

	result, client := rerun(t, cfg)
	if client.calls != 0 {
		t.Errorf("a package change made %d model call(s), want 0", client.calls)
	}
	if got := modeOf(t, result, "prov-2026-aaaaaaaa"); got != ModeRerender {
		t.Errorf("mode = %v, want re-render", got)
	}
	got := userModel(t, cfg)
	if !strings.Contains(got, "attribute :email, :string") {
		t.Errorf("the region was not re-rendered from the new template:\n%s", got)
	}
	if n := strings.Count(got, "attribute :email"); n != 1 {
		t.Errorf("the region appears %d times, want 1:\n%s", n, got)
	}
	if markerGenIn(t, got) != genBefore {
		t.Errorf("re-rendering changed the region's gen")
	}

	// And it is current afterwards.
	again, client := rerun(t, cfg)
	if client.calls != 0 || modeOf(t, again, "prov-2026-aaaaaaaa") != ModeSkip {
		t.Errorf("after re-rendering, the record is not current")
	}
}

func TestForceAndReselectIgnoreTheHashes(t *testing.T) {
	cfg := grown(t)

	forced := cfg
	forced.Force = true
	if _, client := rerun(t, forced); client.calls != 1 {
		t.Errorf("--force made %d model call(s), want 1", client.calls)
	}

	recording := cfg
	recording.Reselect = true
	if _, client := rerun(t, recording); client.calls != 1 {
		t.Errorf("a recording run made %d model call(s), want 1", client.calls)
	}
}

// --check reports without a model call and without writing, Phase 3 files
// included.
func TestCheckReportsAndWritesNothing(t *testing.T) {
	fresh := config(t)
	fresh.Check = true
	client := &stub{}
	fresh.Client = client
	result, err := Run(context.Background(), fresh)
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if client.calls != 0 {
		t.Errorf("--check made %d model call(s)", client.calls)
	}
	if got := tree(t, fresh.Output); len(got) != 0 {
		t.Errorf("--check created files: %v", got)
	}
	if got := modeOf(t, result, "prov-2026-aaaaaaaa"); got != ModeSelect {
		t.Errorf("a record with no regions: mode = %v, want select (reported as unknown)", got)
	}
	if result.CheckFailed() {
		t.Error("an unknown record alone failed --check")
	}

	cfg := grown(t)
	check := cfg
	check.Check = true

	editTemplate(t, cfg, "attribute :{{field}}, :string\n")
	ageTree(t, cfg.Output)
	before := snap(t, cfg.Output)
	result, client = rerun(t, check)
	if client.calls != 0 || modeOf(t, result, "prov-2026-aaaaaaaa") != ModeRerender || !result.CheckFailed() {
		t.Errorf("a package change: calls=%d, classifications=%+v; want re-render and a failing check", client.calls, result.Classifications)
	}
	sameSnapshot(t, before, snap(t, cfg.Output))

	editRecord(t, cfg, "Something else.")
	result, _ = rerun(t, check)
	if modeOf(t, result, "prov-2026-aaaaaaaa") != ModeReselect || !result.CheckFailed() {
		t.Errorf("a changed record: classifications=%+v; want re-select and a failing check", result.Classifications)
	}
	sameSnapshot(t, before, snap(t, cfg.Output))
}

// Seeded regions decide nothing, but a seeded region whose record moved is
// reported.
func TestAStaleSeededRegionIsReported(t *testing.T) {
	cfg := grown(t)
	path := filepath.Join(cfg.Output, "app", "models", "user.rb")
	seeded := strings.Replace(userModel(t, cfg), `"tier":"owned"`, `"tier":"seeded"`, 1)
	if err := os.WriteFile(path, []byte(seeded), 0o644); err != nil {
		t.Fatal(err)
	}
	editRecord(t, cfg, "Something else.")

	check := cfg
	check.Check = true
	result, _ := rerun(t, check)
	for _, c := range result.Classifications {
		if c.RecordID == "prov-2026-aaaaaaaa" {
			if c.Mode != ModeSelect {
				t.Errorf("a record owning only seeded regions: mode = %v, want select", c.Mode)
			}
			if len(c.StaleSeeded) != 1 {
				t.Errorf("stale seeded regions = %+v, want the one seeded region", c.StaleSeeded)
			}
		}
	}
}

func markerGenIn(t *testing.T, content string) string {
	t.Helper()
	_, after, ok := strings.Cut(content, `"gen":"`)
	if !ok {
		t.Fatalf("no gen on the marker:\n%s", content)
	}
	gen, _, _ := strings.Cut(after, `"`)
	return gen
}
