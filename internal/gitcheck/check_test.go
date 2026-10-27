package gitcheck

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/pipeline"
)

// The orchestration runs provenance for each commit, decides which findings the
// proofs waive, and writes the merged output.

const (
	ckRecord = "prov-2026-aaaaaaaa"
	ckRefuse = ckRecord + " is implemented: run 'provenance revise " + ckRecord + "' to start a new version, then tag the commit with it."
)

var ckRecordText = "id: " + ckRecord + "\ntitle: Add email\nstatus: implemented\nversion: 3\ncreated_at: \"2026-10-27\"\nauthor: a@example.com\nintent: |\n  Add email.\nconstraints:\n  - Add an email field\naffected_scope:\n  - app/models/user.rb\n"

// ckRepo is a git repository holding a generator package, one implemented
// record and its grown region, committed once.
type ckRepo struct {
	*repoT
	fake string
}

func newCkRepo(t *testing.T) *ckRepo {
	t.Helper()
	r := &ckRepo{repoT: newRepo(t), fake: t.TempDir()}
	writeFiles(t, r.dir, railsPackage)
	writeFiles(t, r.dir, map[string]string{"provenance/" + ckRecord + ".yml": ckRecordText})
	r.grow()
	r.commit("grow [prov-2026-aaaaaaaa@3]", nil)
	// The fake provenance prints, per commit, the log saved under the commit's
	// sha, and exits 1 as provenance does when it has findings.
	script := `case "$*" in
*--staged*) f="$FAKE_DIR/staged.json";;
*) sha=$(echo "$*" | sed -n 's/.*--commit \([0-9a-f]*\).*/\1/p'); f="$FAKE_DIR/$sha.json";;
esac
if [ -f "$f" ]; then cat "$f"; else echo '{"runs":[{"tool":{"driver":{"rules":[]}},"results":[]}]}'; fi
exit 1`
	bin := fakeBinary(t, t.TempDir(), "provenance", script)
	t.Setenv("PROVENANCE", bin)
	t.Setenv("FAKE_DIR", r.fake)
	return r
}

func (r *ckRepo) grow() {
	r.t.Helper()
	result, err := pipeline.Run(context.Background(), pipeline.Config{
		Generators: filepath.Join(r.dir, "generators"),
		Records:    filepath.Join(r.dir, "provenance"),
		Output:     r.dir,
		Client:     cannedModel{},
		Force:      true,
	})
	if err != nil {
		r.t.Fatalf("grow: %v", err)
	}
	if err := pipeline.WriteRecordings(result, filepath.Join(r.dir, ".sedum", "recordings")); err != nil {
		r.t.Fatalf("WriteRecordings: %v", err)
	}
}

func (r *ckRepo) findings(sha string, results ...string) {
	r.t.Helper()
	log := `{"runs":[{"tool":{"driver":{"rules":[{"id":"STATE001"},{"id":"SCOPE001"},{"id":"TAG001"}]}},"results":[` + strings.Join(results, ",") + `]}]}`
	if err := os.WriteFile(filepath.Join(r.fake, sha+".json"), []byte(log), 0o644); err != nil {
		r.t.Fatal(err)
	}
}

func (r *ckRepo) options(rng string) Options {
	top, _ := TopLevel(r.dir, r.env)
	return Options{
		Range:      rng,
		Output:     r.dir,
		TopLevel:   top,
		Generators: "generators",
		Records:    "provenance",
		Format:     "sarif",
		Paths:      PathConfig{RecordsDir: "provenance", ExcludePaths: []string{"generators/**", ".sedum/**"}},
	}
}

func (r *ckRepo) check(opts Options) (Outcome, string, error) {
	var out bytes.Buffer
	t := r.t
	t.Helper()
	// Check reads the environment: keep it free of the developer's git state.
	for _, kv := range r.env {
		if k, v, ok := strings.Cut(kv, "="); ok {
			t.Setenv(k, v)
		}
	}
	outcome, err := Check(opts, &out)
	return outcome, out.String(), err
}

type sarifOut struct {
	Runs []struct {
		Results []struct {
			RuleID       string `json:"ruleId"`
			Suppressions []struct {
				Kind          string `json:"kind"`
				Justification string `json:"justification"`
			} `json:"suppressions"`
		} `json:"results"`
		Properties map[string]any `json:"properties"`
	} `json:"runs"`
}

func decodeSARIF(t *testing.T, text string) sarifOut {
	t.Helper()
	var s sarifOut
	if err := json.Unmarshal([]byte(text), &s); err != nil || len(s.Runs) != 1 {
		t.Fatalf("output is not one SARIF run: %v\n%s", err, text)
	}
	return s
}

func refusalResult() string {
	return result("STATE001", 0, ckRefuse, "", ckRecord)
}

func sha7(sha string) string { return sha[:7] }

func TestCheckWaivesARegenerationWithItsJustification(t *testing.T) {
	r := newCkRepo(t)
	r.write("generators/rails/actions/addField.rb", "attribute :{{field}}, :string\n")
	r.grow()
	sha := r.commit("regenerate [prov-2026-aaaaaaaa@3]", nil)
	r.findings(sha, refusalResult())

	outcome, out, err := r.check(r.options("HEAD~1..HEAD"))

	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if outcome.Waived != 1 || outcome.Unwaived != 0 {
		t.Errorf("outcome = %+v, want one waived and none unwaived", outcome)
	}
	res := decodeSARIF(t, out).Runs[0].Results
	if len(res) != 1 || len(res[0].Suppressions) != 1 || res[0].Suppressions[0].Justification != "regeneration: "+sha7(sha) {
		t.Errorf("results = %+v", res)
	}
}

func TestCheckLeavesTheFindingWhenTheCommitIsMoreThanARegeneration(t *testing.T) {
	r := newCkRepo(t)
	r.write("generators/rails/actions/addField.rb", "attribute :{{field}}, :string\n")
	r.grow()
	r.write("notes.txt", "a hand-written file\n")
	sha := r.commit("regenerate and more [prov-2026-aaaaaaaa@3]", nil)
	r.findings(sha, refusalResult())

	outcome, _, err := r.check(r.options("HEAD~1..HEAD"))

	if err != nil || outcome.Unwaived != 1 || outcome.Waived != 0 {
		t.Errorf("outcome = %+v err = %v, want the finding left (a new governed file changed)", outcome, err)
	}
}

func TestCheckLeavesTheFindingWhenThereIsNoTag(t *testing.T) {
	r := newCkRepo(t)
	r.write("generators/rails/actions/addField.rb", "attribute :{{field}}, :string\n")
	r.grow()
	sha := r.commit("regenerate without a tag", nil)
	r.findings(sha, refusalResult())

	outcome, _, err := r.check(r.options("HEAD~1..HEAD"))

	if err != nil || outcome.Unwaived != 1 {
		t.Errorf("outcome = %+v err = %v", outcome, err)
	}
}

func TestCheckGivesARootCommitNoWaiver(t *testing.T) {
	r := newCkRepo(t)
	root := r.git("rev-list", "--max-parents=0", "HEAD")
	r.findings(root, refusalResult())

	outcome, _, err := r.check(r.options(root + "^!"))

	if err != nil || outcome.Unwaived != 1 || outcome.Waived != 0 {
		t.Errorf("outcome = %+v err = %v", outcome, err)
	}
}

func TestCheckWaivesAMarkerMigrationAndTheScopeFindingsOfTheTag(t *testing.T) {
	r := newCkRepo(t)
	model := filepath.Join(r.dir, "app", "models", "user.rb")
	v2, err := os.ReadFile(model)
	if err != nil {
		t.Fatal(err)
	}
	// Step back to v1 markers: the same gen, pkg and body, kwargs as an object,
	// and no recording; commit that as the base.
	var v1 []string
	for _, line := range strings.Split(string(v2), "\n") {
		if strings.Contains(line, "sedum:addField") && !strings.Contains(line, "/sedum:") {
			i := strings.Index(line, "{")
			var m map[string]any
			if err := json.Unmarshal([]byte(line[i:]), &m); err != nil {
				t.Fatal(err)
			}
			old := map[string]any{"tier": m["tier"], "record": m["record"], "gen": m["gen"], "pkg": m["pkg"], "kwargs": map[string]string{"field": "email", "model": "user"}}
			data, _ := json.Marshal(old)
			line = line[:i] + string(data)
		}
		v1 = append(v1, line)
	}
	if err := os.WriteFile(model, []byte(strings.Join(v1, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(r.dir, ".sedum")); err != nil {
		t.Fatal(err)
	}
	r.commit("back to v1 markers", nil)
	// The migration: the v2 file and its recording.
	r.grow()
	if err := os.WriteFile(model, v2, 0o644); err != nil {
		t.Fatal(err)
	}
	sha := r.commit("migrate markers [prov-2026-aaaaaaaa@3]", nil)
	r.findings(sha, refusalResult(), result("SCOPE001", 1, "forbids changes to a.go", "a.go", ckRecord), result("TAG001", 2, "tag", "", ckRecord))

	outcome, out, err := r.check(r.options("HEAD~1..HEAD"))

	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if outcome.Waived != 2 || outcome.Unwaived != 1 {
		t.Errorf("outcome = %+v, want the refusal and the scope finding waived, the tag finding left", outcome)
	}
	res := decodeSARIF(t, out).Runs[0].Results
	for _, i := range []int{0, 1} {
		if len(res[i].Suppressions) != 1 || res[i].Suppressions[0].Justification != "marker migration: "+sha7(sha) {
			t.Errorf("result %d suppressions = %+v", i, res[i].Suppressions)
		}
	}
	if len(res[2].Suppressions) != 0 {
		t.Errorf("the tag finding was suppressed: %+v", res[2].Suppressions)
	}
}

func TestCheckDoesNotWaiveAMigrationForAnUnversionedTag(t *testing.T) {
	r := newCkRepo(t)
	r.write("generators/rails/actions/addField.rb", "attribute :{{field}}, :string\n")
	r.grow()
	sha := r.commit("regenerate [prov-2026-aaaaaaaa]", nil)
	r.findings(sha, result("SCOPE001", 0, "forbids changes to a.go", "a.go", ckRecord))

	outcome, _, err := r.check(r.options("HEAD~1..HEAD"))

	if err != nil || outcome.Unwaived != 1 {
		t.Errorf("outcome = %+v err = %v: a scope finding is waived by migration only", outcome, err)
	}
}

func TestCheckWritesAnEmptyLogForAnEmptyRange(t *testing.T) {
	r := newCkRepo(t)

	outcome, out, err := r.check(r.options("HEAD..HEAD"))

	if err != nil || outcome != (Outcome{}) {
		t.Fatalf("outcome = %+v err = %v", outcome, err)
	}
	if res := decodeSARIF(t, out).Runs[0].Results; len(res) != 0 {
		t.Errorf("results = %v", res)
	}
}

func TestCheckWritesTheHumanFormat(t *testing.T) {
	r := newCkRepo(t)
	sha := r.commit("another", map[string]string{"x.txt": "1"})
	r.findings(sha, result("TAG001", 0, "Commit tag required", "", ""))
	opts := r.options("HEAD~1..HEAD")
	opts.Format = "human"

	outcome, out, err := r.check(opts)

	if err != nil || outcome.Unwaived != 1 {
		t.Fatalf("outcome = %+v err = %v", outcome, err)
	}
	if !strings.Contains(out, "TAG001: Commit tag required\n") || !strings.HasSuffix(out, "skipped 0 scope paths that no generator package claims\nwaived: 0\n") {
		t.Errorf("output = %q", out)
	}
}

func TestCheckFailsWhenProvenanceCannotRun(t *testing.T) {
	r := newCkRepo(t)
	t.Setenv("PROVENANCE", filepath.Join(t.TempDir(), "missing"))

	_, out, err := r.check(r.options("HEAD~1..HEAD"))

	if err == nil || !strings.HasPrefix(err.Error(), "provenance:") || out != "" {
		t.Errorf("err = %v out = %q, want a provenance: error and no output", err, out)
	}
}

func TestCheckCountsEachDistinctSkippedScopePathOnce(t *testing.T) {
	r := newCkRepo(t)
	r.write("generators/rails/actions/addField.rb", "attribute :{{field}}, :string\n")
	r.grow()
	writeFiles(t, r.dir, map[string]string{"provenance/prov-2026-bbbbbbbb.yml": strings.Replace(strings.Replace(ckRecordText, ckRecord, "prov-2026-bbbbbbbb", -1), "app/models/user.rb", "docs/NOTES.md", 1)})
	sha := r.commit("regenerate [prov-2026-aaaaaaaa@3]", nil)
	r.findings(sha, refusalResult())

	_, out, err := r.check(r.options("HEAD~1..HEAD"))

	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if got := decodeSARIF(t, out).Runs[0].Properties["skippedScopePaths"]; got != float64(1) {
		t.Errorf("skippedScopePaths = %v, want 1 (docs/NOTES.md, once)", got)
	}
}

func TestCheckWaivesAStagedRegenerationNamingStaged(t *testing.T) {
	r := newCkRepo(t)
	r.write("generators/rails/actions/addField.rb", "attribute :{{field}}, :string\n")
	r.grow()
	r.git("add", "-A")
	msg := filepath.Join(t.TempDir(), "MSG")
	if err := os.WriteFile(msg, []byte("regenerate [prov-2026-aaaaaaaa@3]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(r.fake, "staged.json"), []byte(`{"runs":[{"tool":{"driver":{"rules":[{"id":"STATE001"}]}},"results":[`+refusalResult()+`]}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	opts := r.options("")
	opts.Staged = true
	opts.MessageFile = msg

	outcome, out, err := r.check(opts)

	if err != nil || outcome.Waived != 1 || outcome.Unwaived != 0 {
		t.Fatalf("outcome = %+v err = %v", outcome, err)
	}
	res := decodeSARIF(t, out).Runs[0].Results
	if len(res) != 1 || len(res[0].Suppressions) != 1 || res[0].Suppressions[0].Justification != "regeneration: staged" {
		t.Errorf("results = %+v", res)
	}
}
