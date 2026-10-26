package cli

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

// Recordings, hand edits and replay at the command line (prov-2026-9621daed,
// prov-2026-c637c262, prov-2026-5cc2f7f6, prov-2026-bf01cd6b).
//
// Nothing here calls a real model. A run that needs a selection talks to a
// local server that answers a canned completion, so a test can also count the
// model calls a run made, which is what "skip" and "re-render" promise is zero.

const v2Record = "prov-2026-aaaaaaaa"

const (
	v2EmailSelection = `{"invocations":[{"action":"addField","kwargs":{"model":"user","field":"email"}}]}`
	v2NameSelection  = `{"invocations":[{"action":"addField","kwargs":{"model":"team","field":"name"}}]}`

	v2CheckFailed     = "check failed: some records are stale, edited, need re-rendering, or have left orphan regions; run sedum grow (with --prune-orphans for orphans)"
	v2CheckFailedJSON = "check failed: some records are stale, edited, need re-rendering, or have left orphan regions"
)

// v2WantError compares an error to the text a record gives, ignoring the one
// full stop the record's sentence may end with.
func v2WantError(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected the error %q, got nil", want)
	}
	if got := strings.TrimSuffix(err.Error(), "."); got != strings.TrimSuffix(want, ".") {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

// v2Model is a server speaking the chat completions protocol, answering the
// responses it was given in turn, and counting the calls it got.
type v2Model struct {
	mu        sync.Mutex
	responses []string
	calls     int
}

func (m *v2Model) Calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls
}

func newV2Model(t *testing.T, responses ...string) *v2Model {
	t.Helper()
	m := &v2Model{responses: responses}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		m.mu.Lock()
		m.calls++
		content := m.responses[(m.calls-1)%len(m.responses)]
		m.mu.Unlock()

		body, _ := json.Marshal(map[string]any{
			"id": "chatcmpl-1", "object": "chat.completion", "model": "fake",
			"choices": []any{map[string]any{
				"index": 0, "finish_reason": "stop",
				"message": map[string]any{"role": "assistant", "content": content},
			}},
			"usage": map[string]any{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)
	t.Setenv("OPENAI_BASE_URL", srv.URL)
	t.Setenv("OPENAI_API_KEY", "test")
	return m
}

// v2Generators is a throwaway generator directory: one package, .rb files with
// # comments, and four actions.
func v2Generators(t *testing.T) string {
	t.Helper()
	return writeFiles(t, map[string]string{
		"rails/sedum.yaml":                 "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n",
		"rails/files/app/models/{name}.rb": "class Model\n  # sedum:anchor:class_body\nend\n",
		"rails/actions/actions.yaml": `actions:
  addField:
    kwargs:
      model: { type: string, required: true }
      field: { type: string, required: true }
    injects_into: "app/models/{{model|snake}}.rb"
    anchor: class_body
  setWeight:
    kwargs:
      model: { type: string, required: true }
      weight: { type: int, required: true }
    injects_into: "app/models/{{model|snake}}.rb"
    anchor: class_body
  addColumn:
    kwargs:
      model: { type: string, required: true }
      name: { type: string, required: true }
      kind: { type: string, required: false, default: string }
    injects_into: "app/models/{{model|snake}}.rb"
    anchor: class_body
    identity: [model, name]
`,
		"rails/actions/addField.rb":  "attribute :{{field}}\n",
		"rails/actions/setWeight.rb": "weight {{weight}}\n",
		"rails/actions/addColumn.rb": "column :{{name}}\n",
	})
}

func v2RecordYAML(id, intent string, files ...string) string {
	scope := ""
	for _, f := range files {
		scope += "  - " + f + "\n"
	}
	return "id: " + id + "\nintent: |\n  " + intent + "\naffected_scope:\n" + scope
}

// v2Fix is a generators directory, a records directory and an empty output
// tree, with a model behind them.
type v2Fix struct {
	t               *testing.T
	Gens, Recs, Out string
	Model           *v2Model
}

func newV2Fix(t *testing.T, responses ...string) *v2Fix {
	t.Helper()
	t.Setenv("SEDUM_GROW_LOG", "off")
	t.Setenv("SEDUM_EXCLUDE", "")
	if len(responses) == 0 {
		responses = []string{v2EmailSelection}
	}
	return &v2Fix{
		t:    t,
		Gens: v2Generators(t),
		Recs: writeFiles(t, map[string]string{
			v2Record + ".yml": v2RecordYAML(v2Record, "Add the email field.", "app/models/user.rb"),
		}),
		Out:   t.TempDir(),
		Model: newV2Model(t, responses...),
	}
}

// newV2TwoFix has a second record, prov-2026-bbbbbbbb, owning a region of
// app/models/team.rb. The model answers the records in the order of their ids.
func newV2TwoFix(t *testing.T) *v2Fix {
	t.Helper()
	f := newV2Fix(t, v2EmailSelection, v2NameSelection)
	if err := os.WriteFile(filepath.Join(f.Recs, "prov-2026-bbbbbbbb.yml"),
		[]byte(v2RecordYAML("prov-2026-bbbbbbbb", "Add the name field.", "app/models/team.rb")), 0o644); err != nil {
		t.Fatal(err)
	}
	return f
}

func (f *v2Fix) path(rel string) string { return filepath.Join(f.Out, filepath.FromSlash(rel)) }

func (f *v2Fix) read(rel string) string {
	f.t.Helper()
	b, err := os.ReadFile(f.path(rel))
	if err != nil {
		f.t.Fatal(err)
	}
	return string(b)
}

func (f *v2Fix) write(rel, text string) {
	f.t.Helper()
	full := f.path(rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		f.t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(text), 0o644); err != nil {
		f.t.Fatal(err)
	}
}

func (f *v2Fix) grow(args ...string) (string, error) {
	f.t.Helper()
	base := []string{"grow", "--generators", f.Gens, "--records", f.Recs, "--output", f.Out,
		"--model", "fake", "--log", filepath.Join(f.t.TempDir(), "run.log")}
	return exec(f.t, append(base, args...)...)
}

func (f *v2Fix) check(args ...string) (string, error) {
	f.t.Helper()
	return f.grow(append([]string{"--check"}, args...)...)
}

func (f *v2Fix) recordingPath(id string) string {
	return f.path(".sedum/recordings/" + id + ".json")
}

// edit changes the text of the email region by hand.
func (f *v2Fix) edit() {
	f.t.Helper()
	text := f.read("app/models/user.rb")
	if !strings.Contains(text, "attribute :email") {
		f.t.Fatalf("no email region to edit:\n%s", text)
	}
	f.write("app/models/user.rb", strings.Replace(text, "attribute :email", "attribute :email_by_hand", 1))
}

// movePkg changes the addField template, so every region of the action is
// pkg-stale and its new text is known.
func (f *v2Fix) movePkg() {
	f.t.Helper()
	if err := os.WriteFile(filepath.Join(f.Gens, "rails", "actions", "addField.rb"), []byte("attribute :{{field}} # moved\n"), 0o644); err != nil {
		f.t.Fatal(err)
	}
}

func (f *v2Fix) appendToGeneratorManifest(text string) {
	f.t.Helper()
	path := filepath.Join(f.Gens, "rails", "sedum.yaml")
	b, err := os.ReadFile(path)
	if err != nil {
		f.t.Fatal(err)
	}
	if err := os.WriteFile(path, append(b, text...), 0o644); err != nil {
		f.t.Fatal(err)
	}
}

func (f *v2Fix) setRecordIntent(intent string) {
	f.t.Helper()
	if err := os.WriteFile(filepath.Join(f.Recs, v2Record+".yml"),
		[]byte(v2RecordYAML(v2Record, intent, "app/models/user.rb")), 0o644); err != nil {
		f.t.Fatal(err)
	}
}

func v2ReadJSON(t *testing.T, path string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatalf("%s is not a JSON object: %v", path, err)
	}
	return doc
}

func v2JSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func v2Region0(t *testing.T, path string) map[string]any {
	t.Helper()
	rec := v2ReadJSON(t, path)["records"].([]any)[0].(map[string]any)
	regions, _ := rec["regions"].([]any)
	if len(regions) == 0 {
		t.Fatalf("%s has no regions", path)
	}
	return regions[0].(map[string]any)
}

// v2Marker is one marked region of a file as the test reads it, with the
// attribute object left as raw JSON.
type v2Marker struct {
	Line   string
	Label  string
	Indent string
	Attrs  map[string]json.RawMessage
	Body   string
}

func v2Markers(t *testing.T, text string) []v2Marker {
	t.Helper()
	var out []v2Marker
	lines := strings.SplitAfter(text, "\n")
	for i := 0; i < len(lines); i++ {
		rest, ok := strings.CutPrefix(strings.TrimSpace(lines[i]), "# sedum:")
		if !ok || strings.HasPrefix(rest, "anchor:") {
			continue
		}
		label, attrs, _ := strings.Cut(rest, " ")
		m := v2Marker{Line: lines[i], Label: label, Indent: lines[i][:len(lines[i])-len(strings.TrimLeft(lines[i], " \t"))]}
		if err := json.Unmarshal([]byte(attrs), &m.Attrs); err != nil {
			t.Fatalf("marker line %q: %v", lines[i], err)
		}
		for j := i + 1; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) == "# /sedum:"+label {
				i = j
				break
			}
			m.Body += lines[j]
		}
		out = append(out, m)
	}
	return out
}

func v2Str(t *testing.T, raw json.RawMessage) string {
	t.Helper()
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("%s is not a JSON string: %v", raw, err)
	}
	return s
}

// v2ToV1 rewrites every v2 marker of a file as the v1 marker that carried its
// kwargs, for a region whose identity is all its kwargs.
func v2ToV1(t *testing.T, text string) string {
	t.Helper()
	for _, m := range v2Markers(t, text) {
		v1 := m.Indent + "# sedum:" + m.Label + ` {"tier":` + string(m.Attrs["tier"]) + `,"record":` + string(m.Attrs["record"]) +
			`,"gen":` + string(m.Attrs["gen"]) + `,"pkg":` + string(m.Attrs["pkg"]) + `,"kwargs":` + string(m.Attrs["id"]) + "}\n"
		text = strings.Replace(text, m.Line, v1, 1)
	}
	return text
}

type v2File struct {
	content string
	mtime   time.Time
}

func v2Snap(t *testing.T, root string) map[string]v2File {
	t.Helper()
	out := map[string]v2File{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		info, err := d.Info()
		if err != nil {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(rel)] = v2File{string(b), info.ModTime()}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// v2AgeTree moves every file's modification time into the past, so a write of
// identical bytes still shows.
func v2AgeTree(t *testing.T, root string) {
	t.Helper()
	past := time.Now().Add(-time.Hour)
	for rel := range v2Snap(t, root) {
		if err := os.Chtimes(filepath.Join(root, filepath.FromSlash(rel)), past, past); err != nil {
			t.Fatal(err)
		}
	}
}

func v2SameSnap(t *testing.T, before, after map[string]v2File) {
	t.Helper()
	for rel, b := range before {
		a, ok := after[rel]
		switch {
		case !ok:
			t.Errorf("%s disappeared", rel)
		case a.content != b.content:
			t.Errorf("%s changed:\n before %q\n after  %q", rel, b.content, a.content)
		case !a.mtime.Equal(b.mtime):
			t.Errorf("%s was rewritten with the same bytes (mtime %v -> %v)", rel, b.mtime, a.mtime)
		}
	}
	for rel := range after {
		if _, ok := before[rel]; !ok {
			t.Errorf("%s appeared", rel)
		}
	}
}

func v2Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// v2HandEditLine is what --check prints for the region of the email field.
var v2HandEditLine = strings.Repeat(" ", 13) + "app/models/user.rb  addField  (hand-edited: body hash differs)"

func TestRunGrowRecordings(t *testing.T) {
	t.Run("a grow writes one recording per record into the default directory", func(t *testing.T) {
		f := newV2Fix(t)
		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		path := f.recordingPath(v2Record)
		doc := v2ReadJSON(t, path)

		if _, ok := doc["complete"]; ok {
			t.Error("a per-record file carries complete")
		}
		if _, ok := doc["prompt"]; ok {
			t.Error("a per-record file carries a prompt")
		}
		records := doc["records"].([]any)
		if len(records) != 1 {
			t.Fatalf("the file holds %d records, want 1", len(records))
		}
		rec := records[0].(map[string]any)
		if rec["record_id"] != v2Record {
			t.Errorf("record_id = %v", rec["record_id"])
		}
		inv := rec["phases"].([]any)[0].(map[string]any)["invocations"].([]any)
		if len(inv) != 1 || inv[0].(map[string]any)["action"] != "addField" {
			t.Errorf("phases = %s, want the model's one addField", v2JSON(rec["phases"]))
		}

		markers := v2Markers(t, f.read("app/models/user.rb"))
		if len(markers) != 1 {
			t.Fatalf("user.rb has %d regions, want 1", len(markers))
		}
		m := markers[0]
		regions := rec["regions"].([]any)
		if len(regions) != 1 {
			t.Fatalf("the recording has %d regions, want 1", len(regions))
		}
		r := regions[0].(map[string]any)

		// The marker carries hashes, never kwargs.
		if string(m.Attrs["v"]) != "2" {
			t.Errorf("marker v = %s, want 2", m.Attrs["v"])
		}
		if string(m.Attrs["id"]) != `{"field":"email","model":"user"}` {
			t.Errorf("marker id = %s", m.Attrs["id"])
		}
		wantK := v2KwargsHash("field", "email", "model", "user")
		wantB := v2BodyHash(m.Body)
		if v2Str(t, m.Attrs["kwargs"]) != wantK || v2Str(t, m.Attrs["body"]) != wantB {
			t.Errorf("marker kwargs/body = %s / %s, want %s / %s", m.Attrs["kwargs"], m.Attrs["body"], wantK, wantB)
		}

		want := map[string]any{
			"path": "app/models/user.rb", "action": "addField", "variant": "", "tier": "owned",
			"pkg":         v2Str(t, m.Attrs["pkg"]),
			"id":          map[string]any{"field": "email", "model": "user"},
			"kwargs":      map[string]any{"field": "email", "model": "user"},
			"kwargs_hash": wantK, "body": wantB,
		}
		if v2JSON(r) != v2JSON(want) {
			t.Errorf("region = %s\nwant    %s", v2JSON(r), v2JSON(want))
		}
		if rec["gen"] != v2Str(t, m.Attrs["gen"]) {
			t.Errorf("recording gen = %v, marker gen = %s", rec["gen"], m.Attrs["gen"])
		}
	})

	t.Run("a record whose regions are all current keeps its files exactly", func(t *testing.T) {
		f := newV2Fix(t)
		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		v2AgeTree(t, f.Out)
		before := v2Snap(t, f.Out)

		if out, err := f.grow(); err != nil {
			t.Fatalf("second grow: %v\n%s", err, out)
		}
		if f.Model.Calls() != 1 {
			t.Errorf("the second grow called the model: %d calls", f.Model.Calls())
		}
		v2SameSnap(t, before, v2Snap(t, f.Out))
	})

	t.Run("dry run, check and stamp write no recording", func(t *testing.T) {
		for _, flag := range []string{"--dry-run", "--check", "--stamp"} {
			f := newV2Fix(t)
			f.grow(flag) // --check on nothing fails nowhere; the outcome is the tree
			if v2Exists(f.path(".sedum/recordings")) {
				t.Errorf("%s wrote a recordings directory", flag)
			}
		}
	})

	t.Run("record names the directory and forces nothing", func(t *testing.T) {
		f := newV2Fix(t)
		dir := filepath.Join(t.TempDir(), "recs")
		if out, err := f.grow("--record", dir); err != nil {
			t.Fatalf("grow --record: %v\n%s", err, out)
		}
		if !v2Exists(filepath.Join(dir, v2Record+".json")) {
			t.Fatalf("no recording in %s", dir)
		}
		if v2Exists(f.path(".sedum/recordings")) {
			t.Error("the default directory was written too")
		}

		if out, err := f.grow("--record", dir); err != nil {
			t.Fatalf("second grow --record: %v\n%s", err, out)
		}
		if f.Model.Calls() != 1 {
			t.Errorf("--record with a directory forced a re-select: %d model calls", f.Model.Calls())
		}

		// --check reads the recordings from that directory.
		out, err := f.check("--record", dir)
		if err != nil {
			t.Fatalf("--check --record <dir>: %v\n%s", err, out)
		}
		if !strings.Contains(out, "1 current, 0 stale, 0 to re-render, 0 edited, 0 unknown") {
			t.Errorf("--check --record <dir>:\n%s", out)
		}
		// Without it the default directory holds nothing: stale.
		out, err = f.check()
		v2WantError(t, err, v2CheckFailed)
		if !strings.Contains(out, "stale      "+v2Record+"\n") ||
			!strings.Contains(out, strings.Repeat(" ", 13)+"app/models/user.rb  addField  (no recording)\n") {
			t.Errorf("--check without the recordings directory:\n%s", out)
		}
	})

	t.Run("a path ending in json re-selects every record and writes one recording of the run", func(t *testing.T) {
		f := newV2Fix(t)
		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		all := filepath.Join(t.TempDir(), "all.json")
		out, err := f.grow("--record", all)
		if err != nil {
			t.Fatalf("grow --record %s: %v\n%s", all, err, out)
		}
		if f.Model.Calls() != 2 {
			t.Errorf("a current record was not re-selected: %d model calls, want 2", f.Model.Calls())
		}
		if !strings.Contains(out, "re-selected") {
			t.Errorf("the run does not say every record was re-selected:\n%s", out)
		}
		doc := v2ReadJSON(t, all)
		if doc["complete"] != true || len(doc["records"].([]any)) != 1 {
			t.Errorf("the recording of the run = %s, want one record and complete: true", v2JSON(doc))
		}
		// Not a dry run, a check, a stamp or a replay, and no --stop-after:
		// the per-record files are written as well.
		if !v2Exists(f.recordingPath(v2Record)) {
			t.Error("the per-record recording is gone")
		}
	})

	t.Run("stop-after writes per-record files only for a directory", func(t *testing.T) {
		f := newV2Fix(t)
		all := filepath.Join(t.TempDir(), "all.json")
		if out, err := f.grow("--stop-after", "invocations", "--record", all); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		if !v2Exists(all) {
			t.Error("the recording of the run was not written")
		}
		if v2Exists(f.path(".sedum/recordings")) {
			t.Error("--stop-after with a .json --record wrote per-record recordings")
		}

		g := newV2Fix(t)
		dir := filepath.Join(t.TempDir(), "recs")
		if out, err := g.grow("--stop-after", "invocations", "--record", dir); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		if !v2Exists(filepath.Join(dir, v2Record+".json")) {
			t.Error("--stop-after with a directory --record wrote no per-record recording")
		}
	})

	t.Run("a stale recording is a re-select, whatever is wrong with it", func(t *testing.T) {
		damage := []struct {
			name   string
			reason string
			do     func(t *testing.T, f *v2Fix)
		}{
			{"missing", "no recording", func(t *testing.T, f *v2Fix) { os.Remove(f.recordingPath(v2Record)) }},
			{"unreadable", "recording unreadable", func(t *testing.T, f *v2Fix) {
				os.WriteFile(f.recordingPath(v2Record), []byte("this is not json"), 0o644)
			}},
			{"no entry for the region", "recording has no entry for this region", func(t *testing.T, f *v2Fix) {
				v2Rewrite(t, f.recordingPath(v2Record), func(rec map[string]any) { delete(rec, "regions") })
			}},
			{"hash differs from the marker", "recording does not match the marker", func(t *testing.T, f *v2Fix) {
				v2Rewrite(t, f.recordingPath(v2Record), func(rec map[string]any) {
					rec["regions"].([]any)[0].(map[string]any)["kwargs_hash"] = v2sha([]byte("other"))
				})
			}},
			{"kwargs differ from their own hash", "recording does not match the marker", func(t *testing.T, f *v2Fix) {
				v2Rewrite(t, f.recordingPath(v2Record), func(rec map[string]any) {
					rec["regions"].([]any)[0].(map[string]any)["kwargs"].(map[string]any)["field"] = "other"
				})
			}},
		}
		for _, d := range damage {
			t.Run(d.name, func(t *testing.T) {
				f := newV2Fix(t)
				if out, err := f.grow(); err != nil {
					t.Fatalf("grow: %v\n%s", err, out)
				}
				original := f.read(".sedum/recordings/" + v2Record + ".json")
				d.do(t, f)

				out, err := f.check()
				v2WantError(t, err, v2CheckFailed)
				for _, want := range []string{
					"stale      " + v2Record + "\n",
					strings.Repeat(" ", 13) + "app/models/user.rb  addField  (" + d.reason + ")\n",
					"0 current, 1 stale, 0 to re-render, 0 edited, 0 unknown",
				} {
					if !strings.Contains(out, want) {
						t.Errorf("--check lacks %q:\n%s", want, out)
					}
				}

				// The next grow re-selects and writes the recording back.
				if out, err := f.grow(); err != nil {
					t.Fatalf("grow: %v\n%s", err, out)
				}
				if f.Model.Calls() != 2 {
					t.Errorf("%d model calls, want the re-select to make the second", f.Model.Calls())
				}
				if got := f.read(".sedum/recordings/" + v2Record + ".json"); got != original {
					t.Errorf("the recording was not rewritten as it was:\n%s\nwant\n%s", got, original)
				}
			})
		}
	})

	t.Run("a pkg-stale region is re-rendered from its recording with no model call", func(t *testing.T) {
		f := newV2Fix(t)
		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		before := v2Region0(t, f.recordingPath(v2Record))
		f.movePkg()

		out, err := f.check()
		v2WantError(t, err, v2CheckFailed)
		if !strings.Contains(out, "re-render  "+v2Record+" (package changed)\n") ||
			!strings.Contains(out, "0 current, 0 stale, 1 to re-render, 0 edited, 0 unknown") {
			t.Errorf("--check after the template moved:\n%s", out)
		}

		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		if f.Model.Calls() != 1 {
			t.Errorf("the re-render called the model: %d calls", f.Model.Calls())
		}
		text := f.read("app/models/user.rb")
		if !strings.Contains(text, "attribute :email # moved") {
			t.Fatalf("the region was not re-rendered:\n%s", text)
		}
		m := v2Markers(t, text)[0]
		after := v2Region0(t, f.recordingPath(v2Record))
		if after["body"] != v2BodyHash(m.Body) || after["body"] == before["body"] {
			t.Errorf("recording body = %v, want the new text's hash %s (was %v)", after["body"], v2BodyHash(m.Body), before["body"])
		}
		if after["pkg"] != v2Str(t, m.Attrs["pkg"]) || after["pkg"] == before["pkg"] {
			t.Errorf("recording pkg = %v, marker pkg = %s, was %v", after["pkg"], m.Attrs["pkg"], before["pkg"])
		}
		if after["kwargs_hash"] != before["kwargs_hash"] || v2Str(t, m.Attrs["kwargs"]) != before["kwargs_hash"] {
			t.Errorf("the kwargs hash moved though no kwarg did: %v -> %v", before["kwargs_hash"], after["kwargs_hash"])
		}
		if out, err := f.check(); err != nil {
			t.Errorf("--check after the re-render: %v\n%s", err, out)
		}
	})
}

func v2Rewrite(t *testing.T, path string, edit func(rec map[string]any)) {
	t.Helper()
	doc := v2ReadJSON(t, path)
	edit(doc["records"].([]any)[0].(map[string]any))
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

// A re-render never overwrites text a person edited: it is refused with the
// regions named, --force overwrites, and a run that re-selects overwrites as
// it always has.
func TestGrowRefusesToReRenderOverHandEdits(t *testing.T) {
	t.Run("refused with the regions named", func(t *testing.T) {
		f := newV2Fix(t)
		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		f.edit()
		f.movePkg()
		edited := f.read("app/models/user.rb")
		recordingBefore := f.read(".sedum/recordings/" + v2Record + ".json")

		out, err := f.grow()
		if err == nil {
			t.Fatalf("grow re-rendered over a hand edit:\n%s", out)
		}
		lines := strings.Split(strings.TrimRight(err.Error(), "\n"), "\n")
		want := []string{"refusing to re-render over hand edits; --force overwrites them", "  app/models/user.rb  addField"}
		if strings.Join(lines, "\n") != strings.Join(want, "\n") {
			t.Errorf("error = %q, want the lines %q", err.Error(), want)
		}
		if strings.TrimSpace(out) == "" {
			t.Error("the report was not printed before the error")
		}
		if f.Model.Calls() != 1 {
			t.Errorf("a refused re-render made %d model calls", f.Model.Calls()-1)
		}
		if f.read("app/models/user.rb") != edited {
			t.Error("a refused record's file was written, its pkg-stale regions included")
		}
		if f.read(".sedum/recordings/"+v2Record+".json") != recordingBefore {
			t.Error("a refused record's recording was written")
		}

		// --check shows the same decision: edited, not re-render.
		cout, cerr := f.check()
		v2WantError(t, cerr, v2CheckFailed)
		for _, w := range []string{"edited     " + v2Record + "\n", v2HandEditLine + "\n",
			"0 current, 0 stale, 0 to re-render, 1 edited, 0 unknown"} {
			if !strings.Contains(cout, w) {
				t.Errorf("--check lacks %q:\n%s", w, cout)
			}
		}
	})

	t.Run("force overwrites the edit", func(t *testing.T) {
		f := newV2Fix(t)
		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		f.edit()
		f.movePkg()
		out, err := f.grow("--force")
		if err != nil {
			t.Fatalf("grow --force: %v\n%s", err, out)
		}
		if f.Model.Calls() != 2 {
			t.Errorf("--force re-selects every record: %d model calls, want 2", f.Model.Calls())
		}
		text := f.read("app/models/user.rb")
		if strings.Contains(text, "by_hand") || !strings.Contains(text, "# moved") {
			t.Errorf("the edit survived --force, or the region was not re-rendered:\n%s", text)
		}
		if out, err := f.check(); err != nil {
			t.Errorf("--check after --force: %v\n%s", err, out)
		}
	})

	t.Run("a record that is also stale is re-selected over the edit", func(t *testing.T) {
		f := newV2Fix(t)
		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		f.edit()
		f.setRecordIntent("A different intent, so the gen moves.")

		out, err := f.check()
		v2WantError(t, err, v2CheckFailed)
		if !strings.Contains(out, "stale      "+v2Record+"\n") || strings.Contains(out, "edited") {
			t.Errorf("a stale record is never edited:\n%s", out)
		}
		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		if f.Model.Calls() != 2 || strings.Contains(f.read("app/models/user.rb"), "by_hand") {
			t.Errorf("the re-select did not overwrite the edit (%d model calls)", f.Model.Calls())
		}
	})

	t.Run("an edit with nothing to re-render is left alone and is no error", func(t *testing.T) {
		f := newV2Fix(t)
		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		f.edit()
		edited := f.read("app/models/user.rb")
		v2AgeTree(t, f.Out)
		before := v2Snap(t, f.Out)

		if out, err := f.grow(); err != nil {
			t.Fatalf("grow over an edited, otherwise current record: %v\n%s", err, out)
		}
		if f.Model.Calls() != 1 {
			t.Errorf("an edited record made a model call: %d", f.Model.Calls())
		}
		if f.read("app/models/user.rb") != edited {
			t.Error("the edit was overwritten")
		}
		v2SameSnap(t, before, v2Snap(t, f.Out))
	})

	t.Run("the other records are re-rendered and recorded around a refusal", func(t *testing.T) {
		f := newV2TwoFix(t)
		if out, err := f.grow(); err != nil {
			t.Fatalf("grow: %v\n%s", err, out)
		}
		if f.Model.Calls() != 2 {
			t.Fatalf("%d model calls for two records", f.Model.Calls())
		}
		teamBefore := f.read(".sedum/recordings/prov-2026-bbbbbbbb.json")
		userRecordingBefore := f.read(".sedum/recordings/" + v2Record + ".json")
		f.edit()
		f.movePkg()
		editedUser := f.read("app/models/user.rb")

		// aaaa is edited and pkg-stale, bbbb is only pkg-stale, and an edited
		// record counts only as edited.
		out, err := f.check()
		v2WantError(t, err, v2CheckFailed)
		for _, w := range []string{"edited     " + v2Record + "\n", "re-render  prov-2026-bbbbbbbb (package changed)\n",
			"0 current, 0 stale, 1 to re-render, 1 edited, 0 unknown"} {
			if !strings.Contains(out, w) {
				t.Errorf("--check lacks %q:\n%s", w, out)
			}
		}

		_, err = f.grow()
		if err == nil || !strings.HasPrefix(err.Error(), "refusing to re-render over hand edits; --force overwrites them\n  app/models/user.rb  addField") {
			t.Fatalf("grow error = %v", err)
		}
		if strings.Contains(err.Error(), "team.rb") {
			t.Errorf("the error names a region that was not refused: %v", err)
		}
		if f.Model.Calls() != 2 {
			t.Errorf("model called during a re-render run: %d", f.Model.Calls())
		}
		if f.read("app/models/user.rb") != editedUser || f.read(".sedum/recordings/"+v2Record+".json") != userRecordingBefore {
			t.Error("the refused record was written")
		}
		team := f.read("app/models/team.rb")
		if !strings.Contains(team, "attribute :name # moved") {
			t.Errorf("the other record was not re-rendered:\n%s", team)
		}
		teamAfter := f.read(".sedum/recordings/prov-2026-bbbbbbbb.json")
		if teamAfter == teamBefore {
			t.Error("a refused run did not write the recordings of its other records")
		}
		if got, want := v2Region0(t, f.recordingPath("prov-2026-bbbbbbbb"))["body"], v2BodyHash(v2Markers(t, team)[0].Body); got != want {
			t.Errorf("the other record's recorded body = %v, want %v", got, want)
		}
	})
}

// A recording is written or replayed, never both; a stamp never records; and
// --check only reads, so a .json path (which would re-select) is refused with
// it.
func TestGrowRecordFlagCombinations(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "recs")

	for _, rec := range []string{dir, "all.json"} {
		_, err := exec(t, growArgs(t, "--record", rec, "--execute", "in.json")...)
		wantErr(t, err, "record", "execute")

		_, err = exec(t, growArgs(t, "--stamp", "--record", rec)...)
		v2WantError(t, err, "--stamp cannot be combined with --record: --stamp only writes missing hashes onto markers, with no model call")
	}

	const tail = ": --check only reports, from markers, with no model call and nothing written"
	for _, c := range []struct {
		flag  string
		extra []string
	}{
		{"--record", []string{"--record", "all.json"}},
		{"--dry-run", []string{"--dry-run"}},
		{"--force", []string{"--force"}},
		{"--stop-after", []string{"--stop-after", "files"}},
	} {
		_, err := exec(t, growArgs(t, append([]string{"--check"}, c.extra...)...)...)
		v2WantError(t, err, "--check cannot be combined with "+c.flag+tail)
	}

	// A directory is the one --record --check accepts: it reaches the pipeline,
	// which fails on the generators directory this case does not create.
	_, err := exec(t, growArgs(t, "--check", "--record", dir)...)
	if err == nil || strings.Contains(err.Error(), "cannot be combined") {
		t.Errorf("--check --record <dir> was rejected as a combination: %v", err)
	}
}

func TestGrowFlagsV2(t *testing.T) {
	cmd := findCommand(t, "grow")
	flag := func(name string) string {
		t.Helper()
		f := cmd.Flags().Lookup(name)
		if f == nil {
			t.Fatalf("grow has no --%s", name)
		}
		return f.Usage
	}

	const record = "Write one recording per record into this directory, default .sedum/recordings under --output; " +
		"a path ending in .json instead writes one recording of the whole run to that file and re-selects every record."
	if got := flag("record"); got != record {
		t.Errorf("--record help =\n%s\nwant\n%s", got, record)
	}
	if f := cmd.Flags().Lookup("record"); f.Value.Type() != "string" || f.DefValue != "" {
		t.Errorf("--record is %s default %q, want string default empty", f.Value.Type(), f.DefValue)
	}

	for _, want := range []string{"edited", "body hash", "unverifiable", "v1"} {
		if !strings.Contains(flag("check"), want) {
			t.Errorf("--check help does not say %q: %s", want, flag("check"))
		}
	}
	// The re-render limit: a v1 marker's kwargs, or a v2 region's recording
	// entry, are replayed as they are.
	for _, name := range []string{"check", "force"} {
		for _, want := range []string{"v1", "recording", "default", "--var", "--force"} {
			if !strings.Contains(flag(name), want) {
				t.Errorf("--%s help does not mention %q: %s", name, want, flag(name))
			}
		}
	}
	for _, want := range []string{"overwrite", "edited by hand"} {
		if !strings.Contains(flag("force"), want) {
			t.Errorf("--force help does not mention %q: %s", want, flag("force"))
		}
	}
}

// A replay writes the recordings of the records it replayed, keeping what it
// replayed (prov-2026-c637c262).
func TestReplayWritesRecordings(t *testing.T) {
	f := newV2Fix(t)
	rec := filepath.Join(t.TempDir(), "rec.json")
	if out, err := f.grow("--record", rec); err != nil {
		t.Fatalf("grow --record: %v\n%s", err, out)
	}
	grown := f.read("app/models/user.rb")
	original := v2ReadJSON(t, f.recordingPath(v2Record))["records"].([]any)[0].(map[string]any)

	replay := func(output string, extra ...string) (string, error) {
		args := []string{"grow", "--generators", f.Gens, "--records", f.Recs, "--output", output,
			"--execute", rec, "--log", filepath.Join(t.TempDir(), "run.log")}
		return exec(t, append(args, extra...)...)
	}

	t.Run("writes one recording per replayed record", func(t *testing.T) {
		out := t.TempDir()
		if o, err := replay(out); err != nil {
			t.Fatalf("replay: %v\n%s", err, o)
		}
		text, err := os.ReadFile(filepath.Join(out, "app/models/user.rb"))
		if err != nil || string(text) != grown {
			t.Fatalf("replayed file = %q (%v), want what grow wrote:\n%s", text, err, grown)
		}
		path := filepath.Join(out, ".sedum", "recordings", v2Record+".json")
		doc := v2ReadJSON(t, path)
		if _, ok := doc["complete"]; ok {
			t.Error("a per-record file carries complete")
		}
		got := doc["records"].([]any)
		if len(got) != 1 {
			t.Fatalf("the file holds %d records", len(got))
		}
		r := got[0].(map[string]any)
		for _, key := range []string{"regions", "phases", "gen", "record_id"} {
			if v2JSON(r[key]) != v2JSON(original[key]) {
				t.Errorf("%s = %s, want what the grow recorded: %s", key, v2JSON(r[key]), v2JSON(original[key]))
			}
		}
	})

	t.Run("a dry run and a check write none", func(t *testing.T) {
		for _, flag := range []string{"--dry-run", "--check"} {
			out := t.TempDir()
			replay(out, flag)
			if v2Exists(filepath.Join(out, ".sedum")) {
				t.Errorf("a replay with %s wrote recordings", flag)
			}
		}
	})

	t.Run("a check over a tree the replay would leave as it is passes and writes nothing", func(t *testing.T) {
		out := t.TempDir()
		if o, err := replay(out); err != nil {
			t.Fatalf("replay: %v\n%s", err, o)
		}
		v2AgeTree(t, out)
		before := v2Snap(t, out)
		if o, err := replay(out, "--check"); err != nil {
			t.Fatalf("replay --check: %v\n%s", err, o)
		}
		v2SameSnap(t, before, v2Snap(t, out))
	})

	t.Run("the recording holds the restamped gen, and the markers carry it", func(t *testing.T) {
		f.setRecordIntent("An edited record, so its current gen is not the recorded one.")

		plain := t.TempDir()
		if o, err := replay(plain); err != nil {
			t.Fatalf("replay: %v\n%s", err, o)
		}
		if got := v2Region0Gen(t, plain); got != original["gen"] {
			t.Errorf("a replay without --restamp-gen recorded gen %v, want the recorded %v", got, original["gen"])
		}

		out := t.TempDir()
		o, err := replay(out, "--restamp-gen")
		if err != nil {
			t.Fatalf("replay --restamp-gen: %v\n%s", err, o)
		}
		if !regexp.MustCompile(`restamped ` + v2Record + `: gen [0-9a-f]{8} -> [0-9a-f]{8}`).MatchString(o) {
			t.Errorf("no restamped line:\n%s", o)
		}
		gen := v2Region0Gen(t, out)
		if gen == original["gen"] {
			t.Error("the recording kept the old gen")
		}
		text, _ := os.ReadFile(filepath.Join(out, "app/models/user.rb"))
		if m := v2Markers(t, string(text))[0]; v2Str(t, m.Attrs["gen"]) != gen {
			t.Errorf("recording gen = %s, marker gen = %s", gen, m.Attrs["gen"])
		}
		if !strings.Contains(o, "gen "+strings.TrimPrefix(original["gen"].(string), "sha256:")[:8]+" -> "+strings.TrimPrefix(gen.(string), "sha256:")[:8]) {
			t.Errorf("the restamped line does not show the first eight characters of each gen:\n%s", o)
		}
	})

	t.Run("record is rejected with execute, whatever its value", func(t *testing.T) {
		for _, value := range []string{filepath.Join(t.TempDir(), "recs"), "x.json"} {
			_, err := replay(t.TempDir(), "--record", value)
			wantErr(t, err, "record", "execute")
		}
	})
}

// v2Region0Gen is the gen the per-record recording under an output tree holds.
func v2Region0Gen(t *testing.T, out string) any {
	t.Helper()
	doc := v2ReadJSON(t, filepath.Join(out, ".sedum", "recordings", v2Record+".json"))
	return doc["records"].([]any)[0].(map[string]any)["gen"]
}

// sedum regions reports a v2 region's identity as the id object on its marker,
// read without a recording, and a v1 region's as the action's selecting
// kwargs (prov-2026-bf01cd6b).
func TestRegionsIdentityOfAV2Marker(t *testing.T) {
	const (
		rec = "prov-2026-aaaaaaaa"
	)
	gen := "sha256:" + strings.Repeat("c3", 32)
	pkg := "sha256:" + strings.Repeat("d4", 32)

	v2Line := func(action, id string, kwargs map[string]string, body string) string {
		var pairs []string
		for k, v := range kwargs {
			pairs = append(pairs, k, v)
		}
		return `# sedum:` + action + ` {"v":2,"tier":"owned","record":"` + rec + `","gen":"` + gen + `","pkg":"` + pkg +
			`","id":` + id + `,"kwargs":"` + v2KwargsHash(pairs...) + `","body":"` + v2BodyHash(body) + `"}` + "\n"
	}
	v1Line := func(action, kwargs string) string {
		return `# sedum:` + action + ` {"tier":"owned","record":"` + rec + `","gen":"` + gen + `","pkg":"` + pkg + `","kwargs":` + kwargs + "}\n"
	}
	region := func(open, action, body string) string {
		return open + body + "# /sedum:" + action + "\n"
	}

	out := writeFiles(t, map[string]string{
		// A v2 region: identity is the id object, not every kwarg (kind is not
		// in it), and there is no recording anywhere.
		"app/models/a.rb": region(v2Line("addColumn", `{"model":"user","name":"age"}`,
			map[string]string{"kind": "string", "model": "user", "name": "age"}, "column :age\n"), "addColumn", "column :age\n"),
		// A v1 region: the action's selecting kwargs (its identity), not all
		// of what its marker holds.
		"app/models/b.rb": region(v1Line("addColumn", `{"kind":"int","model":"user","name":"age2"}`), "addColumn", "column :age2\n"),
		// A v2 region whose action is gone: still its id object.
		"app/models/c.rb": region(v2Line("gone", `{"model":"c","y":"2"}`, map[string]string{"model": "c", "y": "2"}, "gone\n"), "gone", "gone\n"),
		// A v1 region whose action is gone: every kwarg its marker holds.
		"app/models/d.rb": region(v1Line("gone", `{"a":"1","b":"2"}`), "gone", "gone\n"),
	})

	t.Setenv("SEDUM_EXCLUDE", "")
	text, err := exec(t, "regions", "--generators", v2Generators(t), "--output", out, "--json")
	if err != nil {
		t.Fatalf("regions: %v\n%s", err, text)
	}
	var doc struct {
		Format  string           `json:"format"`
		Regions []map[string]any `json:"regions"`
	}
	if err := json.Unmarshal([]byte(text), &doc); err != nil {
		t.Fatalf("regions --json is not JSON: %v\n%s", err, text)
	}
	if doc.Format != "sedum.regions/1" {
		t.Errorf("format = %q", doc.Format)
	}

	wantIdentity := []string{
		`{"model":"user","name":"age"}`,
		`{"model":"user","name":"age2"}`,
		`{"model":"c","y":"2"}`,
		`{"a":"1","b":"2"}`,
	}
	wantAction := []string{"addColumn", "addColumn", "gone", "gone"}
	if len(doc.Regions) != 4 {
		t.Fatalf("%d regions, want 4:\n%s", len(doc.Regions), text)
	}
	for i, r := range doc.Regions {
		file := "app/models/" + string(rune('a'+i)) + ".rb"
		if r["file"] != file || r["action"] != wantAction[i] || r["record"] != rec || r["tier"] != "owned" || r["variant"] != "" {
			t.Errorf("region %d = %s, want %s %s of %s", i, v2JSON(r), file, wantAction[i], rec)
		}
		if v2JSON(r["identity"]) != wantIdentity[i] {
			t.Errorf("%s identity = %s, want %s", file, v2JSON(r["identity"]), wantIdentity[i])
		}
		// Verbatim from the marker.
		if r["gen"] != gen || r["pkg"] != pkg {
			t.Errorf("%s gen/pkg = %v / %v", file, r["gen"], r["pkg"])
		}
	}
	if v2Exists(filepath.Join(out, ".sedum")) {
		t.Error("regions wrote something")
	}
}
