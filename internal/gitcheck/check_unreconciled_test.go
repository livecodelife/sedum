package gitcheck

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A repository can hold a hand-written file in a record's scope that the
// file template cannot be reconciled with. The class of the record is then
// unknown on both trees, and a pure marker migration of the record is still
// waived; anything more than a pure migration is not (prov-2026-9621daed).

const (
	handWritten = "class Legacy\n  def slug\n  end\nend\n"
	revise001   = "REVISE001"
)

// unreconciledRepo is a repository whose record owns app/models/user.rb, grown,
// and app/models/legacy.rb, which a person wrote. The base commit holds the
// user file with v1 markers and no recording; migration returns the v2 file
// and its recording.
type unreconciledRepo struct {
	*ckRepo
	v2        string
	recording string
	user      string
}

func newUnreconciledRepo(t *testing.T) *unreconciledRepo {
	t.Helper()
	scope := "  - app/models/user.rb\n  - app/models/legacy.rb\n"
	r := &ckRepo{repoT: newRepo(t), fake: t.TempDir()}
	writeFiles(t, r.dir, railsPackage)
	writeFiles(t, r.dir, map[string]string{
		"provenance/" + ckRecord + ".yml": strings.Replace(ckRecordText, "  - app/models/user.rb\n", scope, 1),
	})
	r.grow()

	user := filepath.Join(r.dir, "app", "models", "user.rb")
	v2, err := os.ReadFile(user)
	if err != nil {
		t.Fatal(err)
	}
	recording, err := os.ReadFile(filepath.Join(r.dir, ".sedum", "recordings", ckRecord+".json"))
	if err != nil {
		t.Fatal(err)
	}

	script := `case "$*" in
*--staged*) f="$FAKE_DIR/staged.json";;
*) sha=$(echo "$*" | sed -n 's/.*--commit \([0-9a-f]*\).*/\1/p'); f="$FAKE_DIR/$sha.json";;
esac
if [ -f "$f" ]; then cat "$f"; else echo '{"runs":[{"tool":{"driver":{"rules":[]}},"results":[]}]}'; fi
exit 1`
	t.Setenv("PROVENANCE", fakeBinary(t, t.TempDir(), "provenance", script))
	t.Setenv("FAKE_DIR", r.fake)

	// The hand-written file replaces the one grow created from the template.
	r.write("app/models/legacy.rb", handWritten)
	// The base: v1 markers, no recording.
	r.write("app/models/user.rb", asV1(t, string(v2)))
	if err := os.RemoveAll(filepath.Join(r.dir, ".sedum")); err != nil {
		t.Fatal(err)
	}
	r.commit("base [prov-2026-aaaaaaaa@3]", nil)

	return &unreconciledRepo{ckRepo: r, v2: string(v2), recording: string(recording), user: "app/models/user.rb"}
}

// asV1 rewrites each v2 marker opening line as the v1 marker with the same gen,
// pkg and body, and kwargs as an object.
func asV1(t *testing.T, text string) string {
	t.Helper()
	var out []string
	for _, line := range strings.Split(text, "\n") {
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
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// migrate commits the v2 file and its recording and returns the commit.
func (u *unreconciledRepo) migrate(userText string, extra map[string]string) string {
	u.t.Helper()
	files := map[string]string{
		u.user: userText,
		".sedum/recordings/" + ckRecord + ".json": u.recording,
	}
	for k, v := range extra {
		files[k] = v
	}
	sha := u.commit("migrate markers [prov-2026-aaaaaaaa@3]", files)
	u.findings(sha, result(revise001, 0, ckRefuse, "", ckRecord))
	return sha
}

func TestAMarkerMigrationIsWaivedWhenARecordHasAnUnreconcilableFile(t *testing.T) {
	u := newUnreconciledRepo(t)
	sha := u.migrate(u.v2, nil)

	outcome, out, err := u.check(u.options("HEAD~1..HEAD"))

	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if outcome.Waived != 1 || outcome.Unwaived != 0 {
		t.Fatalf("outcome = %+v, want the REVISE001 refusal waived", outcome)
	}
	res := decodeSARIF(t, out).Runs[0].Results
	if len(res[0].Suppressions) != 1 || res[0].Suppressions[0].Justification != "marker migration: "+sha7(sha) {
		t.Errorf("suppressions = %+v", res[0].Suppressions)
	}
}

// The class being unknown on both trees compares equal, but the other guards
// stand: a migration that changes the body of one region is not waived.
func TestAMigrationThatChangesARegionBodyIsNotWaivedWhenTheClassIsUnknown(t *testing.T) {
	u := newUnreconciledRepo(t)
	changed := strings.Replace(u.v2, "attribute :email", "attribute :email_address", 1)
	if changed == u.v2 {
		t.Fatalf("the fixture has no body to change:\n%s", u.v2)
	}
	u.migrate(changed, nil)

	outcome, _, err := u.check(u.options("HEAD~1..HEAD"))

	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if outcome.Waived != 0 || outcome.Unwaived != 1 {
		t.Errorf("outcome = %+v, want the refusal left: a region body changed", outcome)
	}
}

// A record that is unknown on only one side is not waived: here the file
// is reconcilable at the base and hand-written at the commit.
func TestAMigrationIsNotWaivedWhenTheClassIsUnknownOnOneSideOnly(t *testing.T) {
	u := newUnreconciledRepo(t)
	// Rewind the base so the legacy file reconciles there.
	u.git("checkout", "-q", "-b", "other", "HEAD")
	u.write("app/models/legacy.rb", "class Legacy\nend\n")
	u.commit("legacy reconciles [prov-2026-aaaaaaaa@3]", nil)
	u.migrate(u.v2, map[string]string{"app/models/legacy.rb": handWritten})

	outcome, _, err := u.check(u.options("HEAD~1..HEAD"))

	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if outcome.Waived != 0 || outcome.Unwaived != 1 {
		t.Errorf("outcome = %+v, want the refusal left: the class differs between the trees", outcome)
	}
}
