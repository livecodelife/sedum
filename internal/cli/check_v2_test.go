package cli

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/pipeline"
)

// grow --check after marker v2 (prov-2026-9621daed): a fourth class, edited,
// foreign and unverifiable regions listed beside the records, a stale region
// that says why, and a summary line with the edited count.

func v2Print(result *pipeline.Result) []string {
	var out bytes.Buffer
	printClassifications(&out, result)
	return strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
}

func v2Index(lines []string, want string) int {
	for i, l := range lines {
		if l == want {
			return i
		}
	}
	return -1
}

func v2WantLine(t *testing.T, lines []string, want string) int {
	t.Helper()
	i := v2Index(lines, want)
	if i < 0 {
		t.Errorf("no line %q in:\n%s", want, strings.Join(lines, "\n"))
	}
	return i
}

var v2Summary = regexp.MustCompile(`^\d+ current, \d+ stale, \d+ to re-render, \d+ edited, \d+ unknown$`)

func TestModeEditedFailsACheck(t *testing.T) {
	if pipeline.ModeEdited.String() != "edited" {
		t.Errorf("ModeEdited.String() = %q, want edited", pipeline.ModeEdited.String())
	}
	if pipeline.ModeEdited <= pipeline.ModeRerender {
		t.Error("ModeEdited must be declared after ModeRerender")
	}
	if !(pipeline.Classification{Mode: pipeline.ModeEdited}).Failing() {
		t.Error("an edited record does not fail --check")
	}
	result := &pipeline.Result{Classifications: []pipeline.Classification{{RecordID: "prov-2026-aaaaaaaa", Mode: pipeline.ModeEdited}}}
	if !result.CheckFailed() {
		t.Error("a result with an edited record does not fail --check")
	}
}

func TestPrintClassificationsV2(t *testing.T) {
	ref := func(path, action, variant, reason string) pipeline.RegionRef {
		return pipeline.RegionRef{Path: path, Action: action, Variant: variant, Reason: reason}
	}
	pad := strings.Repeat(" ", 13)

	t.Run("an edited record and its regions", func(t *testing.T) {
		lines := v2Print(&pipeline.Result{Classifications: []pipeline.Classification{{
			RecordID: "prov-2026-aaaaaaaa", Mode: pipeline.ModeEdited,
			Edited: []pipeline.RegionRef{
				ref("app/models/user.rb", "addField", "", ""),
				ref("app/controllers/users_controller.rb", "createControllerMethod", "index", ""),
			},
		}}})

		at := v2WantLine(t, lines, "edited     prov-2026-aaaaaaaa")
		v2WantLine(t, lines, pad+"app/models/user.rb  addField  (hand-edited: body hash differs)")
		v2WantLine(t, lines, pad+"app/controllers/users_controller.rb  createControllerMethod:index  (hand-edited: body hash differs)")
		if at != 0 {
			t.Errorf("the edited record is not the first line:\n%s", strings.Join(lines, "\n"))
		}
		if last := lines[len(lines)-1]; last != "0 current, 0 stale, 0 to re-render, 1 edited, 0 unknown" {
			t.Errorf("summary = %q", last)
		}
	})

	t.Run("every class is counted once, and an edited record only as edited", func(t *testing.T) {
		lines := v2Print(&pipeline.Result{Classifications: []pipeline.Classification{
			{RecordID: "prov-2026-aaaaaaaa", Mode: pipeline.ModeSkip},
			{RecordID: "prov-2026-bbbbbbbb", Mode: pipeline.ModeReselect, Stale: []pipeline.RegionRef{
				ref("a.rb", "addField", "", "no recording"),
				ref("b.rb", "addField", "", ""),
			}},
			{RecordID: "prov-2026-cccccccc", Mode: pipeline.ModeRerender, Rerender: []pipeline.RegionRef{ref("c.rb", "addField", "", "")}},
			{RecordID: "prov-2026-dddddddd", Mode: pipeline.ModeSelect},
			// An edited record keeps the regions a re-render would have
			// rewritten; it is still one edited record.
			{RecordID: "prov-2026-eeeeeeee", Mode: pipeline.ModeEdited,
				Edited:   []pipeline.RegionRef{ref("e.rb", "addField", "", "")},
				Rerender: []pipeline.RegionRef{ref("e.rb", "addField", "", "")}},
		}})

		v2WantLine(t, lines, "current    prov-2026-aaaaaaaa")
		v2WantLine(t, lines, "stale      prov-2026-bbbbbbbb")
		v2WantLine(t, lines, "re-render  prov-2026-cccccccc (package changed)")
		v2WantLine(t, lines, "edited     prov-2026-eeeeeeee")
		v2WantLine(t, lines, pad+"e.rb  addField  (hand-edited: body hash differs)")
		if last := lines[len(lines)-1]; last != "1 current, 1 stale, 1 to re-render, 1 edited, 1 unknown" {
			t.Errorf("summary = %q", last)
		}
		if !v2Summary.MatchString(lines[len(lines)-1]) {
			t.Error("the summary does not have the five counts in their order")
		}
	})

	t.Run("a stale region says why, when it has a reason", func(t *testing.T) {
		lines := v2Print(&pipeline.Result{Classifications: []pipeline.Classification{{
			RecordID: "prov-2026-bbbbbbbb", Mode: pipeline.ModeReselect,
			Stale: []pipeline.RegionRef{
				ref("a.rb", "addField", "", "recording does not match the marker"),
				ref("b.rb", "createControllerMethod", "show", "recording has no entry for this region"),
				ref("c.rb", "addField", "", ""),
			},
		}}})
		v2WantLine(t, lines, pad+"a.rb  addField  (recording does not match the marker)")
		v2WantLine(t, lines, pad+"b.rb  createControllerMethod:show  (recording has no entry for this region)")
		// No reason, no parentheses.
		v2WantLine(t, lines, pad+"c.rb  addField")
	})

	t.Run("records with unverifiable regions are listed after the record lines", func(t *testing.T) {
		lines := v2Print(&pipeline.Result{Classifications: []pipeline.Classification{
			{RecordID: "prov-2026-aaaaaaaa", Mode: pipeline.ModeSkip, Unverifiable: []pipeline.RegionRef{
				ref("a.rb", "addField", "", ""), ref("b.rb", "createControllerMethod", "index", ""),
			}},
			{RecordID: "prov-2026-bbbbbbbb", Mode: pipeline.ModeSkip},
		}})

		current := v2WantLine(t, lines, "current    prov-2026-aaaaaaaa")
		other := v2WantLine(t, lines, "current    prov-2026-bbbbbbbb")
		head := v2WantLine(t, lines, "unverifiable prov-2026-aaaaaaaa")
		v2WantLine(t, lines, pad+"a.rb  addField")
		v2WantLine(t, lines, pad+"b.rb  createControllerMethod:index")
		if head < current || head < other {
			t.Errorf("the unverifiable block comes before the record lines:\n%s", strings.Join(lines, "\n"))
		}
		if v2Index(lines, "unverifiable prov-2026-bbbbbbbb") >= 0 {
			t.Error("a record with nothing unverifiable is listed as unverifiable")
		}
		// Unverifiable never fails and is not a count of its own.
		if last := lines[len(lines)-1]; last != "2 current, 0 stale, 0 to re-render, 0 edited, 0 unknown" {
			t.Errorf("summary = %q", last)
		}
	})

	t.Run("a foreign region is named and left alone", func(t *testing.T) {
		lines := v2Print(&pipeline.Result{Classifications: []pipeline.Classification{{
			RecordID: "prov-2026-aaaaaaaa", Mode: pipeline.ModeSkip,
			Foreign: []pipeline.RegionRef{ref("app/models/x.rb", "addField", "", ""), ref("y.rb", "createControllerMethod", "index", "")},
		}}})
		v2WantLine(t, lines, "foreign app/models/x.rb  addField  (marker version not known; left alone)")
		v2WantLine(t, lines, "foreign y.rb  createControllerMethod:index  (marker version not known; left alone)")
		if last := lines[len(lines)-1]; last != "1 current, 0 stale, 0 to re-render, 0 edited, 0 unknown" {
			t.Errorf("summary = %q: a foreign region takes no part in a decision", last)
		}
	})
}

// The check on a real tree, end to end: what a person sees after editing a
// generated region by hand.
func TestGrowCheckReportsAnEditedRegion(t *testing.T) {
	f := newV2Fix(t)
	if out, err := f.grow(); err != nil {
		t.Fatalf("grow: %v\n%s", err, out)
	}

	// Current first: five counts, and exit zero.
	out, err := f.check()
	if err != nil {
		t.Fatalf("--check on a current tree: %v\n%s", err, out)
	}
	if !strings.Contains(out, "current    "+v2Record+"\n") ||
		!strings.Contains(out, "\n1 current, 0 stale, 0 to re-render, 0 edited, 0 unknown\n") {
		t.Errorf("--check on a current tree:\n%s", out)
	}

	f.edit()
	out, err = f.check()
	v2WantError(t, err, v2CheckFailed)
	for _, want := range []string{
		"edited     " + v2Record + "\n" + v2HandEditLine + "\n",
		"\n0 current, 0 stale, 0 to re-render, 1 edited, 0 unknown\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("--check lacks %q:\n%s", want, out)
		}
	}
	if f.Model.Calls() != 1 {
		t.Errorf("--check called the model: %d calls", f.Model.Calls())
	}

	// A region with no record yet is unknown, not edited.
	g := newV2Fix(t)
	out, err = g.check()
	if err != nil || !strings.Contains(out, "0 current, 0 stale, 0 to re-render, 0 edited, 1 unknown") {
		t.Errorf("--check on an empty tree: %v\n%s", err, out)
	}
}

// --check writes nothing: not the files, not a recording.
func TestGrowCheckWritesNothingOverAnEdit(t *testing.T) {
	f := newV2Fix(t)
	if out, err := f.grow(); err != nil {
		t.Fatalf("grow: %v\n%s", err, out)
	}
	f.edit()
	f.movePkg()
	v2AgeTree(t, f.Out)
	before := v2Snap(t, f.Out)

	f.check()
	v2SameSnap(t, before, v2Snap(t, f.Out))
}

// A v1 region has no body hash: the check can neither pass nor fail it.
func TestGrowCheckCallsAV1RegionUnverifiable(t *testing.T) {
	f := newV2Fix(t)
	if out, err := f.grow(); err != nil {
		t.Fatalf("grow: %v\n%s", err, out)
	}
	f.write("app/models/user.rb", v2ToV1(t, f.read("app/models/user.rb")))

	// An edit cannot be seen in a v1 region.
	f.write("app/models/user.rb", strings.Replace(f.read("app/models/user.rb"), "attribute :email", "attribute :email_by_hand", 1))

	out, err := f.check()
	if err != nil {
		t.Fatalf("a v1 region failed --check: %v\n%s", err, out)
	}
	for _, want := range []string{
		"current    " + v2Record + "\n",
		"unverifiable " + v2Record + "\n" + strings.Repeat(" ", 13) + "app/models/user.rb  addField\n",
		"1 current, 0 stale, 0 to re-render, 0 edited, 0 unknown",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("--check lacks %q:\n%s", want, out)
		}
	}
}

// The grow report names an edited record and, under --check, its edited
// regions, without calling them changed (prov-2026-bf01cd6b).
func TestGrowReportEdited(t *testing.T) {
	f := newV2Fix(t)
	if out, err := f.grow(); err != nil {
		t.Fatalf("grow: %v\n%s", err, out)
	}
	f.edit()

	out, err := f.check("--json")
	v2WantError(t, err, v2CheckFailedJSON)

	var report struct {
		Format      string `json:"format"`
		DryRun      bool   `json:"dry_run"`
		CheckFailed bool   `json:"check_failed"`
		Records     []struct {
			Record  string `json:"record"`
			Mode    string `json:"mode"`
			Changed bool   `json:"changed"`
			Regions []struct {
				File    string `json:"file"`
				Action  string `json:"action"`
				Variant string `json:"variant"`
				Change  string `json:"change"`
			} `json:"regions"`
		} `json:"records"`
	}
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("--check --json did not print JSON: %v\n%s", err, out)
	}
	if report.Format != "sedum.grow/1" || !report.CheckFailed || !report.DryRun {
		t.Errorf("format %q, check_failed %v, dry_run %v", report.Format, report.CheckFailed, report.DryRun)
	}
	if len(report.Records) != 1 {
		t.Fatalf("%d records, want 1:\n%s", len(report.Records), out)
	}
	r := report.Records[0]
	if r.Record != v2Record || r.Mode != "edited" {
		t.Errorf("record %q has mode %q, want %s edited", r.Record, r.Mode, v2Record)
	}
	if r.Changed {
		t.Error("an edited region set its record's changed flag")
	}
	if len(r.Regions) != 1 || r.Regions[0].File != "app/models/user.rb" || r.Regions[0].Action != "addField" ||
		r.Regions[0].Variant != "" || r.Regions[0].Change != "edited" {
		t.Errorf("regions = %+v, want the one edited region", r.Regions)
	}

	// A current tree reports skip, with nothing edited.
	g := newV2Fix(t)
	if o, err := g.grow(); err != nil {
		t.Fatalf("grow: %v\n%s", err, o)
	}
	out, err = g.check("--json")
	if err != nil {
		t.Fatalf("--check --json on a current tree: %v\n%s", err, out)
	}
	if strings.Contains(out, `"edited"`) {
		t.Errorf("a current tree reports something edited:\n%s", out)
	}
}
