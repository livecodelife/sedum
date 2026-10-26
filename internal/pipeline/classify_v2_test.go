package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/record"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/resolve"
)

// Marker v2 classification (prov-2026-9621daed, revised): a v2 owned region
// whose gen is current is checked against its record's recording, a region
// whose text no longer matches its body hash is edited, a v1 region is
// re-rendered from its own kwargs, and a foreign marker takes no part.

const (
	v2UserRB    = "app/models/user.rb"
	handEdited  = "attribute :email, :hand\n"
	originalRow = "attribute :email\n"
)

// v2Edit replaces the text of user.rb once.
func v2Edit(t *testing.T, cfg Config, old, replacement string) {
	t.Helper()
	content := readFile(t, filepath.Join(cfg.Output, v2UserRB))
	if !strings.Contains(content, old) {
		t.Fatalf("user.rb has no %q:\n%s", old, content)
	}
	v2WriteFile(t, cfg, v2UserRB, strings.Replace(content, old, replacement, 1))
}

// v2Env is what classify and rerenderOrSkip are handed, loaded the way Run
// loads it.
type v2Env struct {
	packages *genpkg.Set
	records  *record.Set
	files    []resolve.File
	owners   map[string]inject.Owner
}

func v2Load(t *testing.T, cfg Config) v2Env {
	t.Helper()
	packages, findings, err := genpkg.Load(cfg.Generators, genpkg.Options{})
	if err != nil || findings.HasErrors() {
		t.Fatalf("genpkg.Load: %v %v", err, findings)
	}
	variables, err := packages.ResolveVariables(cfg.Variables)
	if err != nil {
		t.Fatal(err)
	}
	records, _, err := record.Load(cfg.Records, record.Options{})
	if err != nil {
		t.Fatal(err)
	}
	resolutions, _, err := resolve.Paths(packages, records, cfg.Lang)
	if err != nil {
		t.Fatal(err)
	}
	files, err := resolve.Create(resolutions, resolve.Options{Output: cfg.Output, DryRun: true, Variables: variables})
	if err != nil {
		t.Fatal(err)
	}
	return v2Env{packages: packages, records: records, files: files, owners: ownersOf(records)}
}

func (e v2Env) current(t *testing.T, cfg Config) []inject.OwnedRegion {
	t.Helper()
	regions, err := inject.ScanRegions(targetsOf(filesOf(e.files, v2RecordA)), cfg.Output, nil)
	if err != nil {
		t.Fatal(err)
	}
	return regions
}

func TestClassifyV2ACurrentRegionWithItsRecordingIsSkipped(t *testing.T) {
	cfg, _ := v2Grow(t)

	c := v2Class(t, v2Check(t, cfg), v2RecordA)
	if c.Mode != ModeSkip {
		t.Fatalf("mode = %v, want skip", c.Mode)
	}
	if len(c.Stale)+len(c.Rerender)+len(c.Edited)+len(c.Unverifiable)+len(c.Foreign) != 0 {
		t.Errorf("a current region is reported: %+v", c)
	}
	if c.Failing() {
		t.Error("a skipped record fails --check")
	}
}

func TestClassifyV2StaleReasons(t *testing.T) {
	zeros := "sha256:" + strings.Repeat("0", 64)

	for _, tc := range []struct {
		name   string
		mutate func(t *testing.T, cfg Config, dir string)
		reason string
	}{
		{
			name: "no recording",
			mutate: func(t *testing.T, cfg Config, dir string) {
				if err := os.Remove(filepath.Join(dir, v2RecordA+".json")); err != nil {
					t.Fatal(err)
				}
			},
			reason: "no recording",
		},
		{
			name: "no recordings directory at all",
			mutate: func(t *testing.T, cfg Config, dir string) {
				if err := os.RemoveAll(dir); err != nil {
					t.Fatal(err)
				}
			},
			reason: "no recording",
		},
		{
			name: "recording unreadable",
			mutate: func(t *testing.T, cfg Config, dir string) {
				if err := os.WriteFile(filepath.Join(dir, v2RecordA+".json"), []byte("not json"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			reason: "recording unreadable",
		},
		{
			name: "recording holds no record with the id",
			mutate: func(t *testing.T, cfg Config, dir string) {
				v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) { r.RecordID = "prov-2026-bbbbbbbb" })
				// Moved to the file of the record it now names; the file this
				// record's id names is gone.
				if err := os.Remove(filepath.Join(dir, v2RecordA+".json")); err != nil {
					t.Fatal(err)
				}
				if err := os.Rename(filepath.Join(dir, "prov-2026-bbbbbbbb.json"), filepath.Join(dir, v2RecordA+".json")); err != nil {
					t.Fatal(err)
				}
			},
			reason: "recording unreadable",
		},
		{
			name: "recording has no entry for this region: no regions",
			mutate: func(t *testing.T, cfg Config, dir string) {
				v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) { r.Regions = nil })
			},
			reason: "recording has no entry for this region",
		},
		{
			name: "recording has no entry for this region: other identity",
			mutate: func(t *testing.T, cfg Config, dir string) {
				v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) {
					r.Regions[0].ID = map[string]any{"field": "other", "model": "user"}
				})
			},
			reason: "recording has no entry for this region",
		},
		{
			name: "recording has no entry for this region: other path",
			mutate: func(t *testing.T, cfg Config, dir string) {
				v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) { r.Regions[0].Path = "app/models/other.rb" })
			},
			reason: "recording has no entry for this region",
		},
		{
			name: "recording has no entry for this region: other action",
			mutate: func(t *testing.T, cfg Config, dir string) {
				v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) { r.Regions[0].Action = "addOther" })
			},
			reason: "recording has no entry for this region",
		},
		{
			name: "entry kwargs hash differs from the marker's",
			mutate: func(t *testing.T, cfg Config, dir string) {
				v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) { r.Regions[0].KwargsHash = zeros })
			},
			reason: "recording does not match the marker",
		},
		{
			name: "entry kwargs cannot be hashed",
			mutate: func(t *testing.T, cfg Config, dir string) {
				v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) {
					r.Regions[0].Kwargs = map[string]any{"field": "email", "model": "user", "n": 1 << 60}
				})
			},
			reason: "recording does not match the marker",
		},
		{
			name: "entry kwargs do not hash to the entry's kwargs hash",
			mutate: func(t *testing.T, cfg Config, dir string) {
				v2EditRecording(t, dir, v2RecordA, func(r *recording.Record) {
					r.Regions[0].Kwargs = map[string]any{"field": "email", "model": "order"}
				})
			},
			reason: "recording does not match the marker",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg, dir := v2Grow(t)
			tc.mutate(t, cfg, dir)

			result := v2Check(t, cfg)
			c := v2Class(t, result, v2RecordA)
			if c.Mode != ModeReselect {
				t.Fatalf("mode = %v, want re-select", c.Mode)
			}
			if len(c.Stale) != 1 {
				t.Fatalf("stale = %+v, want the one v2 region", c.Stale)
			}
			ref := c.Stale[0]
			if ref.Reason != tc.reason {
				t.Errorf("Reason = %q, want %q", ref.Reason, tc.reason)
			}
			if ref.Path != v2UserRB || ref.Action != "addField" || ref.Variant != "" {
				t.Errorf("stale ref = %+v", ref)
			}
			if !c.Failing() || !result.CheckFailed() {
				t.Error("a record stale for want of its recording passes --check")
			}
		})
	}
}

// Only a v2 region stale against its recording has a Reason; a region stale by
// gen has none.
func TestClassifyV2AStaleGenHasNoReason(t *testing.T) {
	cfg, _ := v2Grow(t)
	editRecord(t, cfg, "Something else.")

	c := v2Class(t, v2Check(t, cfg), v2RecordA)
	if c.Mode != ModeReselect || len(c.Stale) != 1 {
		t.Fatalf("classification = %+v, want a stale region", c)
	}
	if c.Stale[0].Reason != "" {
		t.Errorf("Reason = %q for a region stale by gen, want empty", c.Stale[0].Reason)
	}
}

// A region stale by gen is never looked up in the recording, so a missing
// recording adds no reason to it.
func TestClassifyV2TheRecordingIsOnlyCheckedForRegionsWhoseGenIsCurrent(t *testing.T) {
	cfg, dir := v2Grow(t)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	editRecord(t, cfg, "Something else.")

	c := v2Class(t, v2Check(t, cfg), v2RecordA)
	if len(c.Stale) != 1 || c.Stale[0].Reason != "" {
		t.Errorf("stale = %+v, want one region with no reason", c.Stale)
	}
}

func TestClassifyV2AV1RegionTakesNoPartInTheRecordingCheck(t *testing.T) {
	cfg, dir := v2Grow(t)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	v2ToV1(t, cfg)

	result := v2Check(t, cfg)
	c := v2Class(t, result, v2RecordA)
	if c.Mode != ModeSkip {
		t.Fatalf("mode = %v, want skip: a v1 region needs no recording", c.Mode)
	}
	if len(c.Stale) != 0 {
		t.Errorf("stale = %+v", c.Stale)
	}
	if len(c.Unverifiable) != 1 || c.Unverifiable[0].Path != v2UserRB || c.Unverifiable[0].Action != "addField" {
		t.Errorf("unverifiable = %+v, want the v1 region", c.Unverifiable)
	}
	if len(c.Edited) != 0 {
		t.Errorf("edited = %+v, want none: a v1 region cannot be shown edited", c.Edited)
	}
	if c.Failing() || result.CheckFailed() {
		t.Error("an unverifiable region fails --check")
	}
}

func TestClassifyV2AV2RegionWithNoBodyHashIsUnverifiable(t *testing.T) {
	cfg, _ := v2Grow(t)
	content := readFile(t, filepath.Join(cfg.Output, v2UserRB))
	stripped := regexp.MustCompile(`,"body":"[^"]*"`).ReplaceAllString(content, "")
	if stripped == content {
		t.Fatalf("the v2 marker carries no body hash:\n%s", content)
	}
	v2WriteFile(t, cfg, v2UserRB, stripped)

	result := v2Check(t, cfg)
	c := v2Class(t, result, v2RecordA)
	if c.Mode != ModeSkip || len(c.Unverifiable) != 1 || len(c.Edited) != 0 {
		t.Errorf("classification = %+v, want skip with one unverifiable region", c)
	}
	if result.CheckFailed() {
		t.Error("an unverifiable region fails --check")
	}
}

func TestClassifyV2ASeededRegionIsNeverUnverifiableOrEdited(t *testing.T) {
	cfg, _ := v2Grow(t)
	v2Edit(t, cfg, `"tier":"owned"`, `"tier":"seeded"`)
	v2Edit(t, cfg, originalRow, handEdited)

	c := v2Class(t, v2Check(t, cfg), v2RecordA)
	if c.Mode != ModeSelect {
		t.Errorf("mode = %v, want select: a record owning only seeded regions owns none", c.Mode)
	}
	if len(c.Edited)+len(c.Unverifiable) != 0 {
		t.Errorf("edited = %+v, unverifiable = %+v; seeded regions are in neither", c.Edited, c.Unverifiable)
	}
}

// Edited.

func TestClassifyV2AHandEditedRegionMakesTheRecordEdited(t *testing.T) {
	cfg, _ := v2Grow(t)
	v2Edit(t, cfg, originalRow, handEdited)

	result := v2Check(t, cfg)
	c := v2Class(t, result, v2RecordA)
	if c.Mode != ModeEdited {
		t.Fatalf("mode = %v, want edited", c.Mode)
	}
	if len(c.Edited) != 1 || c.Edited[0].Path != v2UserRB || c.Edited[0].Action != "addField" {
		t.Errorf("edited = %+v", c.Edited)
	}
	if len(c.Stale) != 0 || len(c.Rerender) != 0 {
		t.Errorf("stale = %+v, rerender = %+v; the region is current and its pkg too", c.Stale, c.Rerender)
	}
	if !c.Failing() || !result.CheckFailed() {
		t.Error("an edited record passes --check")
	}
	if c.Refused() {
		t.Error("an edited record with nothing to re-render is refused")
	}
}

func TestClassifyV2ModeEditedIsDeclaredAfterRerenderAndPrintsAsEdited(t *testing.T) {
	if !(ModeEdited > ModeRerender) {
		t.Errorf("ModeEdited = %d is not declared after ModeRerender = %d", ModeEdited, ModeRerender)
	}
	for _, earlier := range []Mode{ModeNone, ModeSelect, ModeReselect, ModeSkip, ModeRerender} {
		if earlier == ModeEdited {
			t.Errorf("ModeEdited equals %v", earlier)
		}
	}
	if got := ModeEdited.String(); got != "edited" {
		t.Errorf("ModeEdited.String() = %q, want edited", got)
	}
}

func TestClassifyV2FailingAndRefused(t *testing.T) {
	for mode, want := range map[Mode]bool{
		ModeNone: false, ModeSelect: false, ModeSkip: false,
		ModeReselect: true, ModeRerender: true, ModeEdited: true,
	} {
		if got := (Classification{Mode: mode}).Failing(); got != want {
			t.Errorf("Classification{Mode: %v}.Failing() = %v, want %v", mode, got, want)
		}
	}
	// A forced re-select does not fail --check.
	if (Classification{Mode: ModeReselect, Forced: true}).Failing() {
		t.Error("a forced re-select fails --check")
	}

	ref := RegionRef{Path: "a.rb", Action: "addField"}
	other := RegionRef{Path: "b.rb", Action: "addField"}
	for _, tc := range []struct {
		name string
		c    Classification
		want bool
	}{
		{"edited with the region in both lists", Classification{Mode: ModeEdited, Edited: []RegionRef{ref}, Rerender: []RegionRef{ref}}, true},
		{"edited, lists disjoint", Classification{Mode: ModeEdited, Edited: []RegionRef{ref}, Rerender: []RegionRef{other}}, false},
		{"edited, nothing to re-render", Classification{Mode: ModeEdited, Edited: []RegionRef{ref}}, false},
		{"edited, overlap among several", Classification{Mode: ModeEdited, Edited: []RegionRef{other, ref}, Rerender: []RegionRef{ref, other}}, true},
		{"not edited, overlap", Classification{Mode: ModeRerender, Edited: []RegionRef{ref}, Rerender: []RegionRef{ref}}, false},
		{"skip", Classification{Mode: ModeSkip}, false},
	} {
		if got := tc.c.Refused(); got != tc.want {
			t.Errorf("%s: Refused() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestClassifyV2AnEditedRegionWithAStalePkgIsRefusedAndNothingIsWritten(t *testing.T) {
	cfg, dir := v2Grow(t)
	v2Edit(t, cfg, originalRow, handEdited)
	editTemplate(t, cfg, "attribute :{{field}}, :string\n")
	ageTree(t, cfg.Output)
	before := snap(t, cfg.Output)

	result, client := v2Run(t, cfg, fields("email"))
	if client.calls != 0 {
		t.Errorf("an edited record made %d model call(s)", client.calls)
	}
	c := v2Class(t, result, v2RecordA)
	if c.Mode != ModeEdited {
		t.Fatalf("mode = %v, want edited", c.Mode)
	}
	if len(c.Edited) != 1 || len(c.Rerender) != 1 {
		t.Errorf("edited = %+v, rerender = %+v; the regions a re-render would have rewritten stay listed", c.Edited, c.Rerender)
	}
	if !c.Refused() {
		t.Error("an edited region a re-render would have rewritten is not refused")
	}
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}
	sameSnapshot(t, before, snap(t, cfg.Output))
}

func TestClassifyV2AnEditedRecordWithACurrentPkgWritesNothingEither(t *testing.T) {
	cfg, dir := v2Grow(t)
	v2Edit(t, cfg, originalRow, handEdited)
	ageTree(t, cfg.Output)
	before := snap(t, cfg.Output)

	result, client := v2Run(t, cfg, fields("email"))
	if client.calls != 0 || v2Class(t, result, v2RecordA).Mode != ModeEdited {
		t.Fatalf("calls = %d, classifications = %+v", client.calls, result.Classifications)
	}
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatal(err)
	}
	sameSnapshot(t, before, snap(t, cfg.Output))
}

func TestClassifyV2ForceOverwritesAnEditedRegion(t *testing.T) {
	cfg, _ := v2Grow(t)
	v2Edit(t, cfg, originalRow, handEdited)

	forced := cfg
	forced.Force = true
	result, client := v2Run(t, forced, fields("email"))
	if client.calls != 1 {
		t.Errorf("--force made %d model call(s), want 1", client.calls)
	}
	if got := v2Class(t, result, v2RecordA); got.Mode != ModeReselect {
		t.Errorf("classification = %+v, want a forced re-select", got)
	}
	if got := userModel(t, cfg); strings.Contains(got, handEdited) || !strings.Contains(got, originalRow) {
		t.Errorf("--force did not overwrite the hand edit:\n%s", got)
	}
}

func TestClassifyV2AStaleOrNewRecordIsNeverEdited(t *testing.T) {
	t.Run("stale gen", func(t *testing.T) {
		cfg, _ := v2Grow(t)
		v2Edit(t, cfg, originalRow, handEdited)
		editRecord(t, cfg, "Something else.")

		c := v2Class(t, v2Check(t, cfg), v2RecordA)
		if c.Mode != ModeReselect {
			t.Errorf("mode = %v, want re-select", c.Mode)
		}
		if len(c.Edited) != 0 {
			t.Errorf("edited = %+v: only current regions are", c.Edited)
		}
	})
	t.Run("forced", func(t *testing.T) {
		cfg, _ := v2Grow(t)
		v2Edit(t, cfg, originalRow, handEdited)
		cfg.Force = true

		if c := v2Class(t, v2Check(t, cfg), v2RecordA); c.Mode != ModeReselect {
			t.Errorf("mode = %v, want re-select", c.Mode)
		}
	})
	t.Run("stale for want of a recording", func(t *testing.T) {
		cfg, dir := v2Grow(t)
		v2Edit(t, cfg, originalRow, handEdited)
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}

		c := v2Class(t, v2Check(t, cfg), v2RecordA)
		if c.Mode != ModeReselect || len(c.Stale) != 1 {
			t.Errorf("classification = %+v, want re-select", c)
		}
	})
	t.Run("new", func(t *testing.T) {
		cfg := config(t)
		if c := v2Class(t, v2Check(t, cfg), v2RecordA); c.Mode != ModeSelect {
			t.Errorf("mode = %v, want select", c.Mode)
		}
	})
}

// Foreign.

const foreignRegion = "  # sedum:addField {\"v\":3,\"tier\":\"owned\",\"record\":\"prov-2026-aaaaaaaa\"}\n" +
	"attribute :foreign\n" +
	"  # /sedum:addField\n"

func v2AddForeign(t *testing.T, cfg Config) {
	t.Helper()
	content := readFile(t, filepath.Join(cfg.Output, v2UserRB))
	v2WriteFile(t, cfg, v2UserRB, strings.TrimSuffix(content, "end\n")+foreignRegion+"end\n")
}

func TestClassifyV2AForeignRegionIsReportedAndTakesNoPartInTheDecision(t *testing.T) {
	cfg, _ := v2Grow(t)
	v2AddForeign(t, cfg)

	result := v2Check(t, cfg)
	c := v2Class(t, result, v2RecordA)
	if c.Mode != ModeSkip {
		t.Errorf("mode = %v, want skip: a foreign region decides nothing", c.Mode)
	}
	if len(c.Foreign) != 1 || c.Foreign[0].Path != v2UserRB || c.Foreign[0].Action != "addField" || c.Foreign[0].Variant != "" {
		t.Errorf("foreign = %+v, want the one foreign region", c.Foreign)
	}
	if len(c.Stale)+len(c.Edited)+len(c.Unverifiable)+len(c.Rerender) != 0 {
		t.Errorf("a foreign region was counted as something else: %+v", c)
	}
	if result.CheckFailed() {
		t.Error("a foreign region fails --check")
	}
}

func TestClassifyV2AForeignRegionAloneDoesNotMakeARecordOwnAnything(t *testing.T) {
	cfg := config(t)
	cfg.Client = &stub{}
	if _, err := Run(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	v2AddForeign(t, cfg)

	c := v2Class(t, v2Check(t, cfg), v2RecordA)
	if c.Mode != ModeSelect {
		t.Errorf("mode = %v, want select: the foreign region is owned by no one", c.Mode)
	}
	if len(c.Foreign) != 1 {
		t.Errorf("foreign = %+v", c.Foreign)
	}
}

func TestClassifyV2AGrowLeavesAForeignRegionExactlyAsItIs(t *testing.T) {
	cfg, _ := v2Grow(t)
	v2AddForeign(t, cfg)
	editRecord(t, cfg, "Something else.")

	_, client := v2Run(t, cfg, fields("phone"))
	if client.calls != 1 {
		t.Errorf("model calls = %d, want 1", client.calls)
	}
	got := userModel(t, cfg)
	if !strings.Contains(got, foreignRegion) {
		t.Errorf("the foreign region was changed or pruned:\n%s", got)
	}
	if !strings.Contains(got, "attribute :phone") || strings.Contains(got, originalRow) {
		t.Errorf("the record's own regions were not re-selected:\n%s", got)
	}
}

// The recordings directory.

func TestClassifyV2ReadsTheRecordingFromTheRecordingsDirectoryOfTheRun(t *testing.T) {
	cfg := config(t)
	cfg.Recordings = filepath.Join(t.TempDir(), "elsewhere")
	result, _ := v2Run(t, cfg, fields("email"))
	if err := WriteRecordings(result, cfg.Recordings); err != nil {
		t.Fatal(err)
	}

	if c := v2Class(t, v2Check(t, cfg), v2RecordA); c.Mode != ModeSkip {
		t.Errorf("with the recordings directory set: mode = %v, want skip", c.Mode)
	}

	// The default directory under the output holds nothing, and is not read.
	if got := tree(t, cfg.Output); len(got) != 1 {
		t.Errorf("output tree = %v, want only the generated file", got)
	}
	cfg.Recordings = ""
	c := v2Class(t, v2Check(t, cfg), v2RecordA)
	if c.Mode != ModeReselect || len(c.Stale) != 1 || c.Stale[0].Reason != "no recording" {
		t.Errorf("with the default directory: %+v, want re-select for want of a recording", c)
	}
}

func TestClassifyV2AnEmptyRecordingsDirectoryMeansTheOneUnderTheOutput(t *testing.T) {
	cfg, dir := v2Grow(t)
	if want := filepath.Join(cfg.Output, ".sedum", "recordings"); dir != want {
		t.Fatalf("test setup: dir = %q", dir)
	}
	cfg.Recordings = ""
	if c := v2Class(t, v2Check(t, cfg), v2RecordA); c.Mode != ModeSkip {
		t.Errorf("mode = %v, want skip from .sedum/recordings under the output", c.Mode)
	}
}

func TestClassifyV2TheFunctionTakesTheRecordingsDirectoryLast(t *testing.T) {
	cfg, dir := v2Grow(t)
	env := v2Load(t, cfg)
	skip := []string{cfg.Generators, cfg.Records}

	got, err := classify(env.records, env.files, env.owners, cfg.Output, nil, false, skip, dir)
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	if len(got) != 1 || got[0].Mode != ModeSkip {
		t.Fatalf("classify = %+v, want skip", got)
	}

	// Another directory holds no recording.
	got, err = classify(env.records, env.files, env.owners, cfg.Output, nil, false, skip, t.TempDir())
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	if got[0].Mode != ModeReselect || len(got[0].Stale) != 1 || got[0].Stale[0].Reason != "no recording" {
		t.Errorf("classify with an empty directory = %+v, want re-select, no recording", got[0])
	}

	// Asking for a re-select ignores the hashes.
	got, err = classify(env.records, env.files, env.owners, cfg.Output, nil, true, skip, dir)
	if err != nil {
		t.Fatalf("classify: %v", err)
	}
	if got[0].Mode != ModeReselect || !got[0].Forced {
		t.Errorf("classify with reselect = %+v, want a forced re-select", got[0])
	}
}

func TestClassifyV2ARecordingVouchesForEveryRegionOfItsRecord(t *testing.T) {
	// Two regions of one record, each vouched for by the one recording.
	cfg, dir := v2Grow(t)
	cfg.Client = &stub{response: fields("email", "phone")}
	editRecord(t, cfg, "Something else.")
	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatal(err)
	}

	c := v2Class(t, v2Check(t, cfg), v2RecordA)
	if c.Mode != ModeSkip || len(c.Stale) != 0 {
		t.Errorf("classification = %+v, want skip", c)
	}
}

// Re-render.

func TestRerenderFromRecordingAV2RegionIsReRenderedWithNoModelCall(t *testing.T) {
	cfg, dir := v2Grow(t)
	editTemplate(t, cfg, "attribute :{{field}}, :string\n")

	check := v2Class(t, v2Check(t, cfg), v2RecordA)
	if check.Mode != ModeRerender || len(check.Rerender) != 1 || check.Rerender[0].Reason != "" {
		t.Fatalf("classification = %+v, want one region to re-render", check)
	}

	result, client := v2Run(t, cfg, fields("email"))
	if client.calls != 0 {
		t.Errorf("a re-render made %d model call(s)", client.calls)
	}
	got := userModel(t, cfg)
	if !strings.Contains(got, "attribute :email, :string") || strings.Count(got, "attribute :email") != 1 {
		t.Errorf("the region was not re-rendered in place:\n%s", got)
	}
	region := v2Regions(t, cfg, v2UserRB)[0].Marker
	if region.Version != 2 || region.Body != hashBodyEmailT || region.KwargsHash != hashKwargsEmail {
		t.Errorf("marker after re-render = %+v", region)
	}

	if err := WriteRecordings(result, dir); err != nil {
		t.Fatal(err)
	}
	// And it is current afterwards, against the recording just rewritten.
	again, client := v2Run(t, cfg, fields("email"))
	if client.calls != 0 || v2Class(t, again, v2RecordA).Mode != ModeSkip {
		t.Errorf("after re-rendering: calls = %d, %+v; want skip", client.calls, again.Classifications)
	}
}

// v2TypedGenerators is the rails fixture with an action whose identity is two
// of its three kwargs, so a v2 marker's id holds less than the region was
// rendered from.
func v2TypedGenerators(t *testing.T) string {
	t.Helper()
	return writeTree(t, map[string]string{
		"rails/sedum.yaml":                 "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\ntransforms:\n  constantize: [singular, pascal]\n",
		"rails/files/app/models/{name}.rb": "class {{name|constantize}}\n  # sedum:anchor:class_body\nend\n",
		"rails/actions/actions.yaml": `actions:
  addTyped:
    kwargs:
      model: { type: string, required: true }
      field: { type: string, required: true }
      kind: { type: string, required: false, default: "string" }
    identity: [model, field]
    injects_into: "app/models/{{model|snake}}.rb"
    anchor: class_body
`,
		"rails/actions/addTyped.rb": "attribute :{{field}}, :{{kind}}\n",
	})
}

// A v2 marker carries only the identity kwargs. The kwargs that parameterize
// the region, and so its text, are in the recording entry, and a re-render
// replays them.
func TestRerenderFromRecordingReplaysTheEntrysKwargsNotTheMarkers(t *testing.T) {
	cfg := config(t)
	cfg.Generators = v2TypedGenerators(t)
	selection := `{"invocations": [{"action": "addTyped", "kwargs": {"model": "user", "field": "age", "kind": "integer"}}]}`
	result, _ := v2Run(t, cfg, selection)
	dir := v2Dir(cfg)
	if err := WriteRecordings(result, dir); err != nil {
		t.Fatal(err)
	}

	marker := v2Regions(t, cfg, v2UserRB)[0].Marker
	if jsonOf(t, marker.Kwargs) != `{"field":"age","model":"user"}` {
		t.Fatalf("marker id = %v, want only the identity kwargs", marker.Kwargs)
	}
	entry := v2ReadRecord(t, dir, v2RecordA).Regions[0]
	if jsonOf(t, entry.Kwargs) != `{"field":"age","kind":"integer","model":"user"}` || jsonOf(t, entry.ID) != `{"field":"age","model":"user"}` {
		t.Fatalf("entry = %+v", entry)
	}
	if !strings.Contains(userModel(t, cfg), "attribute :age, :integer") {
		t.Fatalf("test setup: %s", userModel(t, cfg))
	}

	if err := os.WriteFile(filepath.Join(cfg.Generators, "rails", "actions", "addTyped.rb"), []byte("attribute :{{field}}, :{{kind}}, null: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, client := v2Run(t, cfg, selection)
	if client.calls != 0 {
		t.Errorf("a re-render made %d model call(s)", client.calls)
	}
	if got := userModel(t, cfg); !strings.Contains(got, "attribute :age, :integer, null: false") {
		t.Errorf("the re-render did not replay the recorded kind:\n%s", got)
	}
}

func TestRerenderFromRecordingAV1RegionIsReRenderedFromItsOwnKwargs(t *testing.T) {
	cfg, dir := v2Grow(t)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	v2ToV1(t, cfg)
	genBefore := markerGenIn(t, userModel(t, cfg))
	editTemplate(t, cfg, "attribute :{{field}}, :string\n")

	check := v2Class(t, v2Check(t, cfg), v2RecordA)
	if check.Mode != ModeRerender || len(check.Rerender) != 1 {
		t.Fatalf("classification = %+v, want re-render from the v1 marker with no recording", check)
	}

	_, client := v2Run(t, cfg, fields("email"))
	if client.calls != 0 {
		t.Errorf("a re-render made %d model call(s)", client.calls)
	}
	got := userModel(t, cfg)
	if !strings.Contains(got, "attribute :email, :string") || strings.Count(got, "attribute :email") != 1 {
		t.Errorf("the v1 region was not re-rendered:\n%s", got)
	}
	marker := v2Regions(t, cfg, v2UserRB)[0].Marker
	if marker.Version != 2 || marker.KwargsHash != hashKwargsEmail || marker.Body != hashBodyEmailT {
		t.Errorf("a re-rendered v1 region is not written with a v2 marker: %+v", marker)
	}
	if marker.Gen != genBefore {
		t.Errorf("re-rendering changed the gen: %q -> %q", genBefore, marker.Gen)
	}
}

// rerenderOrSkip, called directly.

func TestRerenderFromRecordingRerenderOrSkip(t *testing.T) {
	cfg, dir := v2Grow(t)
	editTemplate(t, cfg, "attribute :{{field}}, :string\n")
	env := v2Load(t, cfg)
	rec := env.records.Records[0]
	mine := filesOf(env.files, rec.ID)
	entries := v2ReadRecord(t, dir, rec.ID).Regions
	current := env.current(t, cfg)
	if len(current) != 1 || len(entries) != 1 {
		t.Fatalf("test setup: current = %d, entries = %d", len(current), len(entries))
	}

	t.Run("a v2 region needing a re-render is re-rendered from its entry", func(t *testing.T) {
		c := Classification{RecordID: rec.ID}
		if got := rerenderOrSkip(&c, rec, mine, current, entries); got != ModeRerender {
			t.Errorf("mode = %v, want re-render", got)
		}
		if len(c.Rerender) != 1 || c.Rerender[0].Path != v2UserRB {
			t.Errorf("Rerender = %+v", c.Rerender)
		}
	})

	t.Run("a v2 region with no entry sends its record to re-select", func(t *testing.T) {
		for name, regions := range map[string][]recording.Region{"nil": nil, "empty": {}} {
			c := Classification{RecordID: rec.ID}
			if got := rerenderOrSkip(&c, rec, mine, current, regions); got != ModeReselect {
				t.Errorf("%s entries: mode = %v, want re-select", name, got)
			}
		}
		other := []recording.Region{entries[0]}
		other[0].ID = map[string]any{"field": "other", "model": "user"}
		c := Classification{RecordID: rec.ID}
		if got := rerenderOrSkip(&c, rec, mine, current, other); got != ModeReselect {
			t.Errorf("an entry with another identity: mode = %v, want re-select", got)
		}
	})

	t.Run("the entry's kwargs decide what is rendered, not the marker's", func(t *testing.T) {
		moved := []recording.Region{entries[0]}
		moved[0].Kwargs = map[string]any{"field": "email", "model": "order"}
		c := Classification{RecordID: rec.ID}
		// Rendered from the entry the region would land in order.rb, not the
		// file it sits in.
		if got := rerenderOrSkip(&c, rec, mine, current, moved); got != ModeReselect {
			t.Errorf("mode = %v, want re-select: the entry's kwargs render into another file", got)
		}
	})

	t.Run("every pkg current is a skip", func(t *testing.T) {
		pkg := current[0].Package
		fresh := current[0]
		fresh.Region.Marker.Pkg = pkg.Actions["addField"].PkgHash
		c := Classification{RecordID: rec.ID}
		if got := rerenderOrSkip(&c, rec, mine, []inject.OwnedRegion{fresh}, entries); got != ModeSkip {
			t.Errorf("mode = %v, want skip", got)
		}
		if len(c.Rerender) != 0 {
			t.Errorf("Rerender = %+v", c.Rerender)
		}
	})

	t.Run("a v1 region is re-rendered from its own marker with no recording", func(t *testing.T) {
		v1 := v1Region(rec, current[0], map[string]any{"model": "user", "field": "email"})
		c := Classification{RecordID: rec.ID}
		if got := rerenderOrSkip(&c, rec, mine, []inject.OwnedRegion{v1}, nil); got != ModeRerender {
			t.Errorf("mode = %v, want re-render", got)
		}
		if len(c.Rerender) != 1 {
			t.Errorf("Rerender = %+v", c.Rerender)
		}
	})

	t.Run("a v1 region ignores a recording entry", func(t *testing.T) {
		v1 := v1Region(rec, current[0], map[string]any{"model": "user", "field": "email"})
		moved := []recording.Region{entries[0]}
		moved[0].Kwargs = map[string]any{"field": "email", "model": "order"}
		c := Classification{RecordID: rec.ID}
		if got := rerenderOrSkip(&c, rec, mine, []inject.OwnedRegion{v1}, moved); got != ModeRerender {
			t.Errorf("mode = %v, want re-render from the marker's own kwargs", got)
		}
	})

	t.Run("a v1 region that would land in another file sends its record to re-select", func(t *testing.T) {
		v1 := v1Region(rec, current[0], map[string]any{"model": "order", "field": "email"})
		c := Classification{RecordID: rec.ID}
		if got := rerenderOrSkip(&c, rec, mine, []inject.OwnedRegion{v1}, nil); got != ModeReselect {
			t.Errorf("mode = %v, want re-select", got)
		}
	})

	t.Run("a region whose action no longer exists sends its record to re-select", func(t *testing.T) {
		v1 := v1Region(rec, current[0], map[string]any{"model": "user", "field": "email"})
		v1.Region.Marker.Action = "removedAction"
		c := Classification{RecordID: rec.ID}
		if got := rerenderOrSkip(&c, rec, mine, []inject.OwnedRegion{v1}, nil); got != ModeReselect {
			t.Errorf("mode = %v, want re-select", got)
		}

		v2 := current[0]
		v2.Region.Marker.Action = "removedAction"
		gone := []recording.Region{entries[0]}
		gone[0].Action = "removedAction"
		if got := rerenderOrSkip(&c, rec, mine, []inject.OwnedRegion{v2}, gone); got != ModeReselect {
			t.Errorf("a v2 region: mode = %v, want re-select", got)
		}
	})
}

// v1Region is a current owned v1 region of the record's file whose pkg is
// stale, as grow found before marker v2.
func v1Region(rec *record.Record, from inject.OwnedRegion, kwargs map[string]any) inject.OwnedRegion {
	return inject.OwnedRegion{
		Path:    from.Path,
		Package: from.Package,
		Region: inject.Region{Marker: inject.Marker{
			Action: "addField", Tier: inject.TierOwned, Record: rec.ID, Gen: rec.Gen,
			Pkg: "sha256:stale-pkg", Kwargs: kwargs,
		}},
	}
}
