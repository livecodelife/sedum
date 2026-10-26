package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

// `sedum markers migrate` rewrites v1 markers to v2 and writes the recordings,
// with no model and without changing any region's text (prov-2026-0a398351,
// prov-2026-fa3497af, prov-2026-f6643ae2, prov-2026-8210461f).
//
// The hashes in these fixtures are computed here, from the rules of
// prov-2026-8bba6ce4 and the standard library alone, never with the package
// under test: a test that asked the implementation for its own expected value
// would pass whatever it computed.

// v2sha is the text sha256: followed by 64 lower-case hex digits of the bytes.
func v2sha(b []byte) string {
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// v2KwargsHash is the kwargs hash of a map of string kwargs, given as
// name, value pairs. It covers only what these fixtures need of the canonical
// form: sorted keys, no whitespace, plain ASCII strings.
func v2KwargsHash(pairs ...string) string {
	kv := map[string]string{}
	for i := 0; i+1 < len(pairs); i += 2 {
		kv[pairs[i]] = pairs[i+1]
	}
	names := make([]string, 0, len(kv))
	for name := range kv {
		names = append(names, name)
	}
	sort.Strings(names)

	var canon bytes.Buffer
	canon.WriteByte('{')
	for i, name := range names {
		if i > 0 {
			canon.WriteByte(',')
		}
		k, _ := json.Marshal(name)
		v, _ := json.Marshal(kv[name])
		canon.Write(k)
		canon.WriteByte(':')
		canon.Write(v)
	}
	canon.WriteByte('}')
	return v2sha(canon.Bytes())
}

// v2BodyHash is the body hash of region text: each CRLF read as LF, nothing
// else changed.
func v2BodyHash(text string) string {
	return v2sha([]byte(strings.ReplaceAll(text, "\r\n", "\n")))
}

// The helpers must agree with the two hashes prov-2026-8bba6ce4 and its
// implementation records state outright, or every fixture below is wrong.
func TestV2FixtureHashesAgreeWithTheRecords(t *testing.T) {
	if got, want := v2BodyHash(""), "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"; got != want {
		t.Errorf("empty body hash = %s, want %s", got, want)
	}
	if got, want := v2KwargsHash(), "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a"; got != want {
		t.Errorf("empty kwargs hash = %s, want %s", got, want)
	}
	if v2BodyHash("a\r\nb\r\n") != v2BodyHash("a\nb\n") {
		t.Error("CRLF must hash as LF")
	}
}

func TestRootHasMarkers(t *testing.T) {
	var found bool
	for _, c := range NewRootCommand().Commands() {
		if c.Name() == "markers" {
			found = true
		}
	}
	if !found {
		t.Fatal("root has no markers command")
	}

	cmd, _, err := NewRootCommand().Find([]string{"markers", "migrate"})
	if err != nil || cmd.Name() != "migrate" {
		t.Fatalf("markers migrate is not reachable from the root: %v (%v)", cmd, err)
	}
}

func TestMarkersCommand(t *testing.T) {
	cmd := findCommand(t, "markers")
	if cmd.Use != "markers" {
		t.Errorf("Use = %q, want markers", cmd.Use)
	}
	if cmd.Short != "Work on ownership markers" {
		t.Errorf("Short = %q, want %q", cmd.Short, "Work on ownership markers")
	}

	subs := cmd.Commands()
	if len(subs) != 1 {
		t.Fatalf("markers has %d subcommands, want exactly migrate", len(subs))
	}
	if subs[0].Name() != "migrate" || subs[0].Use != "migrate" {
		t.Errorf("subcommand = %q (Use %q), want migrate", subs[0].Name(), subs[0].Use)
	}
}

func TestMarkersMigrateFlags(t *testing.T) {
	cmd, _, err := NewRootCommand().Find([]string{"markers", "migrate"})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][2]string{ // name -> type, default
		"generators": {"string", ""},
		"records":    {"string", ""},
		"output":     {"string", "."},
		"dry-run":    {"bool", "false"},
		"var":        {"stringArray", "[]"},
	}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == "help" {
			return
		}
		w, ok := want[f.Name]
		if !ok {
			t.Errorf("markers migrate has an undocumented flag --%s", f.Name)
			return
		}
		if f.Value.Type() != w[0] || f.DefValue != w[1] {
			t.Errorf("--%s is %s default %q, want %s default %q", f.Name, f.Value.Type(), f.DefValue, w[0], w[1])
		}
		delete(want, f.Name)
	})
	for name := range want {
		t.Errorf("markers migrate has no flag --%s", name)
	}

	// --generators and --records are required.
	_, err = exec(t, "markers", "migrate", "--records", t.TempDir())
	wantErr(t, err, "generators")
	_, err = exec(t, "markers", "migrate", "--generators", t.TempDir())
	wantErr(t, err, "records")

	// No positional argument.
	_, err = exec(t, "markers", "migrate", "--generators", t.TempDir(), "--records", t.TempDir(), "extra")
	wantErr(t, err, "extra")

	// --var is name=value, parsed with parseVars, so a malformed one fails
	// before anything is read.
	_, err = exec(t, "markers", "migrate", "--generators", t.TempDir(), "--records", t.TempDir(), "--var", "bad")
	wantErr(t, err, "--var", "bad")
}

var (
	migGenA = "sha256:" + strings.Repeat("a1", 32)
	migPkgA = "sha256:" + strings.Repeat("b2", 32)
)

func migV1Line(indent, action, record string, kwargs string) string {
	attrs := `{"tier":"owned"`
	if record != "" {
		attrs += `,"record":"` + record + `","gen":"` + migGenA + `","pkg":"` + migPkgA + `"`
	}
	return indent + "# sedum:" + action + " " + attrs + `,"kwargs":` + kwargs + "}\n"
}

// migFixture is an output tree of v1 markers with the generators and records
// that name them. With refusals it also holds two regions that cannot be
// migrated.
type migFixture struct {
	gens, recs, out string
}

func newMigFixture(t *testing.T, refusals bool) migFixture {
	t.Helper()
	const rec = "prov-2026-aaaaaaaa"
	scope := "  - app/models/user.rb\n"
	if refusals {
		scope += "  - app/models/bad.rb\n  - app/models/heavy.rb\n"
	}
	recs := writeFiles(t, map[string]string{
		rec + ".yml": "id: " + rec + "\nintent: |\n  Add fields.\naffected_scope:\n" + scope,
	})
	files := map[string]string{
		"app/models/user.rb": "class User\n" +
			migV1Line("  ", "addField", rec, `{"field":"email","model":"user"}`) +
			"  attribute :email\n  # /sedum:addField\n" +
			migV1Line("  ", "addField", rec, `{"field":"name","model":"user"}`) +
			"  attribute :name\n  # /sedum:addField\n" +
			"  # sedum:anchor:class_body\nend\n",
		// A record that is not in the records directory.
		"app/models/team.rb": "class Team\n" +
			migV1Line("  ", "addField", "prov-2026-bbbbbbbb", `{"field":"title","model":"team"}`) +
			"  attribute :title\n  # /sedum:addField\nend\n",
		// A marker that names no record.
		"app/models/misc.rb": "class Misc\n" +
			migV1Line("  ", "addField", "", `{"field":"misc","model":"misc"}`) +
			"  attribute :misc\n  # /sedum:addField\nend\n",
	}
	if refusals {
		// An action its package no longer declares.
		files["app/models/bad.rb"] = "class Bad\n" +
			migV1Line("  ", "gone", rec, `{"model":"bad","x":"1"}`) +
			"  gone\n  # /sedum:gone\nend\n"
		// kwargs that cannot be hashed: an integer beyond 2^53.
		files["app/models/heavy.rb"] = "class Heavy\n" +
			migV1Line("  ", "setWeight", rec, `{"model":"heavy","weight":9007199254740993}`) +
			"  weight 9007199254740993\n  # /sedum:setWeight\nend\n"
	}
	return migFixture{gens: v2Generators(t), recs: recs, out: writeFiles(t, files)}
}

func (f migFixture) run(t *testing.T, extra ...string) (string, error) {
	t.Helper()
	args := []string{"markers", "migrate", "--generators", f.gens, "--records", f.recs, "--output", f.out}
	return exec(t, append(args, extra...)...)
}

func (f migFixture) read(t *testing.T, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(f.out, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// migV2Line is the v2 marker line the migration must write for a v1 region
// with the fixture's gen and pkg: the identity object, the kwargs hash and the
// body hash of the text present, in the key order of prov-2026-a3771fbf.
func migV2Line(action, record, field, model, body string) string {
	return `  # sedum:` + action + ` {"v":2,"tier":"owned","record":"` + record + `","gen":"` + migGenA + `","pkg":"` + migPkgA +
		`","id":{"field":"` + field + `","model":"` + model + `"},"kwargs":"` + v2KwargsHash("field", field, "model", model) +
		`","body":"` + v2BodyHash(body) + `"}` + "\n"
}

func TestMarkersMigrate(t *testing.T) {
	const rec = "prov-2026-aaaaaaaa"
	wantUser := "class User\n" +
		migV2Line("addField", rec, "email", "user", "  attribute :email\n") +
		"  attribute :email\n  # /sedum:addField\n" +
		migV2Line("addField", rec, "name", "user", "  attribute :name\n") +
		"  attribute :name\n  # /sedum:addField\n" +
		"  # sedum:anchor:class_body\nend\n"

	t.Run("rewrites what it can and says what it left", func(t *testing.T) {
		f := newMigFixture(t, false)
		before := map[string]string{"app/models/team.rb": f.read(t, "app/models/team.rb"), "app/models/misc.rb": f.read(t, "app/models/misc.rb")}

		out, err := f.run(t)
		if err != nil {
			t.Fatalf("migrate: %v\n%s", err, out)
		}

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		var migrated, left []string
		for _, l := range lines {
			switch {
			case strings.HasPrefix(l, "migrated "):
				migrated = append(migrated, l)
			case strings.HasPrefix(l, "left "):
				left = append(left, l)
			}
		}
		if len(migrated) != 2 || migrated[0] != "migrated app/models/user.rb  addField" || migrated[1] != migrated[0] {
			t.Errorf("migrated lines = %q, want two of %q", migrated, "migrated app/models/user.rb  addField")
		}
		if len(left) != 2 ||
			!regexp.MustCompile(`^left app/models/misc\.rb  addField: .+`).MatchString(left[0]) ||
			!regexp.MustCompile(`^left app/models/team\.rb  addField: .+`).MatchString(left[1]) {
			t.Errorf("left lines = %q, want one for misc.rb and one for team.rb, each ending with the Why text", left)
		}
		if strings.Contains(out, "dry run") || strings.Contains(out, "would") {
			t.Errorf("a real run says it is a dry run:\n%s", out)
		}
		if last := lines[len(lines)-1]; last != "2 regions migrated, 1 files rewritten, 1 recordings written" {
			t.Errorf("final line = %q", last)
		}
		// Migrated lines come before left ones, as the contract lists them.
		if strings.Index(out, "migrated app/") > strings.Index(out, "left app/") {
			t.Errorf("a left line precedes a migrated line:\n%s", out)
		}

		// The text of every region is unchanged and the markers are v2.
		if got := f.read(t, "app/models/user.rb"); got != wantUser {
			t.Errorf("user.rb =\n%s\nwant\n%s", got, wantUser)
		}
		// A left region keeps its v1 marker, byte for byte.
		for rel, text := range before {
			if got := f.read(t, rel); got != text {
				t.Errorf("%s changed though its regions were left:\n%s", rel, got)
			}
		}
	})

	t.Run("writes the recording with the gen found on the markers", func(t *testing.T) {
		f := newMigFixture(t, false)
		if out, err := f.run(t); err != nil {
			t.Fatalf("migrate: %v\n%s", err, out)
		}

		entries, err := os.ReadDir(filepath.Join(f.out, ".sedum", "recordings"))
		if err != nil {
			t.Fatalf("no recordings directory: %v", err)
		}
		if len(entries) != 1 || entries[0].Name() != rec+".json" {
			t.Fatalf("recordings = %v, want only %s.json (a left region has no record to write)", entries, rec)
		}

		doc := v2ReadJSON(t, filepath.Join(f.out, ".sedum", "recordings", rec+".json"))
		records := doc["records"].([]any)
		if len(records) != 1 {
			t.Fatalf("recording holds %d records, want 1", len(records))
		}
		r := records[0].(map[string]any)
		if r["record_id"] != rec {
			t.Errorf("record_id = %v", r["record_id"])
		}
		// Not the record's current gen: a stale marker stays stale.
		if r["gen"] != migGenA {
			t.Errorf("recording gen = %v, want the marker's %s", r["gen"], migGenA)
		}
		regions := r["regions"].([]any)
		if len(regions) != 2 {
			t.Fatalf("recording has %d regions, want the 2 migrated", len(regions))
		}
		for i, field := range []string{"email", "name"} { // sorted by region key
			g := regions[i].(map[string]any)
			body := "  attribute :" + field + "\n"
			want := map[string]any{
				"path": "app/models/user.rb", "action": "addField", "variant": "", "tier": "owned", "pkg": migPkgA,
				"id":          map[string]any{"field": field, "model": "user"},
				"kwargs":      map[string]any{"field": field, "model": "user"},
				"kwargs_hash": v2KwargsHash("field", field, "model", "user"),
				"body":        v2BodyHash(body),
			}
			wb, _ := json.Marshal(want)
			gb, _ := json.Marshal(g)
			if string(wb) != string(gb) {
				t.Errorf("region %d = %s, want %s", i, gb, wb)
			}
		}
	})

	t.Run("a second migration changes nothing", func(t *testing.T) {
		f := newMigFixture(t, false)
		if out, err := f.run(t); err != nil {
			t.Fatalf("migrate: %v\n%s", err, out)
		}
		v2AgeTree(t, f.out)
		before := v2Snap(t, f.out)

		out, err := f.run(t)
		if err != nil {
			t.Fatalf("second migrate: %v\n%s", err, out)
		}
		v2SameSnap(t, before, v2Snap(t, f.out))
		if strings.Contains(out, "migrated app/") {
			t.Errorf("the second run migrated something:\n%s", out)
		}
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if last := lines[len(lines)-1]; last != "0 regions migrated, 0 files rewritten, 0 recordings written" {
			t.Errorf("final line = %q", last)
		}
	})

	t.Run("dry run writes nothing and says would", func(t *testing.T) {
		f := newMigFixture(t, false)
		before := v2Snap(t, f.out)

		out, err := f.run(t, "--dry-run")
		if err != nil {
			t.Fatalf("migrate --dry-run: %v\n%s", err, out)
		}
		v2SameSnap(t, before, v2Snap(t, f.out))
		if _, err := os.Stat(filepath.Join(f.out, ".sedum")); err == nil {
			t.Error("--dry-run wrote a recording")
		}

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if lines[0] != "dry run: nothing is written" {
			t.Errorf("first line = %q, want %q", lines[0], "dry run: nothing is written")
		}
		var would int
		for _, l := range lines {
			if l == "would migrate app/models/user.rb  addField" {
				would++
			}
			if strings.HasPrefix(l, "migrated ") {
				t.Errorf("a dry run says migrated: %q", l)
			}
		}
		if would != 2 {
			t.Errorf("%d would-migrate lines, want 2:\n%s", would, out)
		}
		if last := lines[len(lines)-1]; last != "2 regions would be migrated, 1 files would be rewritten, 1 recordings would be written" {
			t.Errorf("final line = %q", last)
		}
	})

	t.Run("a refused region exits non-zero and keeps its v1 marker", func(t *testing.T) {
		f := newMigFixture(t, true)
		bad, heavy := f.read(t, "app/models/bad.rb"), f.read(t, "app/models/heavy.rb")

		out, err := f.run(t)
		if err == nil || err.Error() != "2 regions could not be migrated" {
			t.Fatalf("error = %v, want %q", err, "2 regions could not be migrated")
		}

		var refused []string
		for _, l := range strings.Split(out, "\n") {
			if strings.HasPrefix(l, "refused ") {
				refused = append(refused, l)
			}
		}
		if len(refused) != 2 ||
			!regexp.MustCompile(`^refused app/models/bad\.rb  gone: .+`).MatchString(refused[0]) ||
			!regexp.MustCompile(`^refused app/models/heavy\.rb  setWeight: .+`).MatchString(refused[1]) {
			t.Errorf("refused lines = %q", refused)
		}
		// The report is printed before the error is returned, and the other
		// regions are migrated regardless.
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if last := lines[len(lines)-1]; last != "2 regions migrated, 1 files rewritten, 1 recordings written" {
			t.Errorf("final line = %q", last)
		}
		if f.read(t, "app/models/bad.rb") != bad || f.read(t, "app/models/heavy.rb") != heavy {
			t.Error("a refused region lost its v1 marker")
		}
		if f.read(t, "app/models/user.rb") != wantUser {
			t.Error("the regions beside a refusal were not migrated")
		}
		// Migrated, then refused, then left, as the contract lists them.
		if !(strings.Index(out, "migrated app/") < strings.Index(out, "refused app/") &&
			strings.Index(out, "refused app/") < strings.Index(out, "left app/")) {
			t.Errorf("lines are not grouped migrated, refused, left:\n%s", out)
		}
	})

	t.Run("a refusal under dry run still exits non-zero and writes nothing", func(t *testing.T) {
		f := newMigFixture(t, true)
		before := v2Snap(t, f.out)

		out, err := f.run(t, "--dry-run")
		if err == nil || err.Error() != "2 regions could not be migrated" {
			t.Fatalf("error = %v, want %q", err, "2 regions could not be migrated")
		}
		v2SameSnap(t, before, v2Snap(t, f.out))
		if !strings.Contains(out, "refused app/models/bad.rb  gone: ") {
			t.Errorf("dry run does not report the refusal:\n%s", out)
		}
	})

	t.Run("nothing to migrate is no error", func(t *testing.T) {
		gens := v2Generators(t)
		recs := writeFiles(t, map[string]string{"prov-2026-aaaaaaaa.yml": "id: prov-2026-aaaaaaaa\nintent: |\n  x\naffected_scope:\n  - app/models/user.rb\n"})
		out, err := exec(t, "markers", "migrate", "--generators", gens, "--records", recs, "--output", t.TempDir())
		if err != nil {
			t.Fatalf("migrate over an empty tree: %v\n%s", err, out)
		}
		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if last := lines[len(lines)-1]; last != "0 regions migrated, 0 files rewritten, 0 recordings written" {
			t.Errorf("final line = %q", last)
		}
	})
}

// Migrating a tree grow wrote must give back exactly what grow wrote, and
// grow --check must say the same before and after (prov-2026-0a398351).
func TestMarkersMigrateMatchesWhatGrowWrites(t *testing.T) {
	f := newV2Fix(t)
	if out, err := f.grow(); err != nil {
		t.Fatalf("grow: %v\n%s", err, out)
	}
	grown := f.read("app/models/user.rb")
	original := v2ReadJSON(t, f.path(".sedum/recordings/"+v2Record+".json"))

	// Turn the region back into the v1 marker that carried its kwargs.
	if err := os.WriteFile(f.path("app/models/user.rb"), []byte(v2ToV1(t, grown)), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(f.path(".sedum")); err != nil {
		t.Fatal(err)
	}

	// A v1 region is current and cannot be verified.
	out, err := f.check()
	if err != nil {
		t.Fatalf("--check over v1 markers: %v\n%s", err, out)
	}
	for _, want := range []string{"current    " + v2Record + "\n", "unverifiable " + v2Record + "\n",
		strings.Repeat(" ", 13) + "app/models/user.rb  addField\n",
		"1 current, 0 stale, 0 to re-render, 0 edited, 0 unknown\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("--check before migration lacks %q:\n%s", want, out)
		}
	}

	mout, err := exec(t, "markers", "migrate", "--generators", f.Gens, "--records", f.Recs, "--output", f.Out)
	if err != nil {
		t.Fatalf("migrate: %v\n%s", err, mout)
	}
	if got := f.read("app/models/user.rb"); got != grown {
		t.Errorf("migrated file differs from the one grow wrote:\n%s\nwant\n%s", got, grown)
	}
	migrated := v2ReadJSON(t, f.path(".sedum/recordings/"+v2Record+".json"))
	mr := migrated["records"].([]any)[0].(map[string]any)
	or := original["records"].([]any)[0].(map[string]any)
	if mj, oj := v2JSON(mr["regions"]), v2JSON(or["regions"]); mj != oj {
		t.Errorf("migrated regions = %s, want what grow recorded: %s", mj, oj)
	}
	if mr["gen"] != or["gen"] {
		t.Errorf("migrated recording gen = %v, want %v", mr["gen"], or["gen"])
	}

	// The same current class, and nothing left to verify.
	out, err = f.check()
	if err != nil {
		t.Fatalf("--check after migration: %v\n%s", err, out)
	}
	if !strings.Contains(out, "1 current, 0 stale, 0 to re-render, 0 edited, 0 unknown") || strings.Contains(out, "unverifiable") {
		t.Errorf("--check after migration:\n%s", out)
	}
	if f.Model.Calls() != 1 {
		t.Errorf("migration or check called the model: %d calls, want the 1 of the first grow", f.Model.Calls())
	}
}

// A stale marker stays stale through migration: the recording carries the gen
// the marker had, not the record's current one.
func TestMarkersMigrateKeepsAStaleClassStale(t *testing.T) {
	f := newV2Fix(t)
	if out, err := f.grow(); err != nil {
		t.Fatalf("grow: %v\n%s", err, out)
	}
	v1 := v2ToV1(t, f.read("app/models/user.rb"))
	if err := os.WriteFile(f.path("app/models/user.rb"), []byte(v1), 0o644); err != nil {
		t.Fatal(err)
	}
	f.setRecordIntent("A different intent, so the gen moves.")

	if out, err := f.check(); err == nil || !strings.Contains(out, "1 stale") {
		t.Fatalf("control: v1 region should be stale before migration: %v\n%s", err, out)
	}
	if out, err := exec(t, "markers", "migrate", "--generators", f.Gens, "--records", f.Recs, "--output", f.Out); err != nil {
		t.Fatalf("migrate: %v\n%s", err, out)
	}
	out, err := f.check()
	if err == nil || !strings.Contains(out, "0 current, 1 stale, 0 to re-render, 0 edited, 0 unknown") {
		t.Errorf("--check after migration should be stale as before: %v\n%s", err, out)
	}

	marker := v2Markers(t, f.read("app/models/user.rb"))[0]
	rec := v2ReadJSON(t, f.path(".sedum/recordings/"+v2Record+".json"))["records"].([]any)[0].(map[string]any)
	if rec["gen"] != v2Str(t, marker.Attrs["gen"]) {
		t.Errorf("recording gen = %v, want the marker's gen %s", rec["gen"], v2Str(t, marker.Attrs["gen"]))
	}
}

// Re-rendering a migrated region from its recording reproduces its body: here
// the package moves without any template changing, so the text must not.
func TestMarkersMigrateThenReRenderReproducesTheBody(t *testing.T) {
	f := newV2Fix(t)
	if out, err := f.grow(); err != nil {
		t.Fatalf("grow: %v\n%s", err, out)
	}
	if err := os.WriteFile(f.path("app/models/user.rb"), []byte(v2ToV1(t, f.read("app/models/user.rb"))), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := exec(t, "markers", "migrate", "--generators", f.Gens, "--records", f.Recs, "--output", f.Out); err != nil {
		t.Fatalf("migrate: %v\n%s", err, out)
	}
	before := v2Markers(t, f.read("app/models/user.rb"))[0]

	// A transform nothing uses moves the package hash and no text.
	f.appendToGeneratorManifest("transforms:\n  unused: [snake]\n")
	if out, err := f.check(); err == nil || !strings.Contains(out, "re-render") {
		t.Fatalf("control: the moved package should need a re-render: %v\n%s", err, out)
	}
	out, err := f.grow()
	if err != nil {
		t.Fatalf("grow: %v\n%s", err, out)
	}
	if f.Model.Calls() != 1 {
		t.Errorf("re-render called the model: %d calls", f.Model.Calls())
	}
	after := v2Markers(t, f.read("app/models/user.rb"))[0]
	if after.Body != before.Body || v2Str(t, after.Attrs["body"]) != v2Str(t, before.Attrs["body"]) {
		t.Errorf("re-render changed the body:\n%q\n%q", before.Body, after.Body)
	}
	if v2Str(t, after.Attrs["pkg"]) == v2Str(t, before.Attrs["pkg"]) {
		t.Error("the pkg hash did not move, so nothing was re-rendered")
	}
}
