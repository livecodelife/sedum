package resolve

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/record"
	"github.com/calebcowen/sedum/internal/runlog"
)

// Phases 2 and 3 are where a path stops being a string and becomes a file on
// disk. The cases below split along that seam: which package a path resolves to
// and why, then what gets written and what deliberately does not.

// generators is a two-package directory: one claiming .rb, one claiming .go,
// with file templates shaped like the targets they came from. Cases that need a
// third package or a contested extension mutate a copy.
func generators() map[string]string {
	return map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
`,
		"rails/files/app/controllers/{name}_controller.rb": "class {{name|constantize}}Controller\n" +
			"  # sedum:anchor:class_body\nend\n",
		"rails/files/app/models/{name}.rb": "class {{name|constantize}}\n  # sedum:anchor:class_body\nend\n",
		"rails/files/_default.rb":          "# sedum:anchor:top\n",
		"rails/actions/actions.yaml": `actions:
  addBeforeFilter:
    kwargs:
      filter: { type: string, required: true }
    injects_into: "app/controllers/x_controller.rb"
    anchor: class_body
`,
		"rails/actions/addBeforeFilter.rb": "before_action :{{filter}}\n",

		"chi/sedum.yaml": `name: chi
extensions: [".go"]
comment_prefix: "//"
transforms:
  exported: [pascal]
`,
		"chi/files/internal/handlers/{name}.go": "package handlers\n\n// sedum:anchor:handlers\n",
		"chi/actions/actions.yaml": `actions:
  addHandlerFunc:
    kwargs:
      resource: { type: string, required: true }
    injects_into: "internal/handlers/{{resource|snake}}.go"
    anchor: handlers
`,
		"chi/actions/addHandlerFunc.go": "func {{resource|exported}}() {}\n",
	}
}

func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
	return root
}

// loadPackages loads a generators directory and fails the test if any package
// was rejected: a resolution case must never be debugging a broken fixture.
func loadPackages(t *testing.T, files map[string]string) *genpkg.Set {
	t.Helper()
	set, findings, err := genpkg.Load(writeTree(t, files), genpkg.Options{})
	if err != nil {
		t.Fatalf("genpkg.Load: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("fixture package %s is invalid: %s: %s", f.Package, f.Rule, f.Message)
		}
	}
	return set
}

// records builds a one-record set authorizing the given paths.
func records(t *testing.T, paths ...string) *record.Set {
	t.Helper()
	body := "id: prov-2026-aaaaaaaa\naffected_scope:\n"
	for _, p := range paths {
		body += "  - \"" + p + "\"\n"
	}
	set, _, err := record.Load(writeTree(t, map[string]string{"r.yml": body}), record.Options{})
	if err != nil {
		t.Fatalf("record.Load: %v", err)
	}
	return set
}

func wantErr(t *testing.T, err error, fragments ...string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error naming %v, got nil", fragments)
	}
	for _, want := range fragments {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error = %q, does not mention %q", err, want)
		}
	}
}

func byPath(resolutions []Resolution) map[string]Resolution {
	out := map[string]Resolution{}
	for _, r := range resolutions {
		out[r.Path] = r
	}
	return out
}

// Resolution is per file, not per run: one record may legitimately touch a .rb
// and a .go path, and each is generated under its own package's conventions.
func TestResolutionIsPerFile(t *testing.T) {
	set := loadPackages(t, generators())
	recs := records(t, "app/models/user.rb", "internal/handlers/user.go")

	resolutions, _, err := Paths(set, recs, nil)
	if err != nil {
		t.Fatalf("Packages: %v", err)
	}

	got := byPath(resolutions)
	if p := got["app/models/user.rb"].Package; p == nil || p.Name != "rails" {
		t.Errorf("app/models/user.rb resolved to %v, want rails", p)
	}
	if p := got["internal/handlers/user.go"].Package; p == nil || p.Name != "chi" {
		t.Errorf("internal/handlers/user.go resolved to %v, want chi", p)
	}
	if got["app/models/user.rb"].RecordID != "prov-2026-aaaaaaaa" {
		t.Errorf("resolution lost the record it came from")
	}
}

// Nothing infers a target from a path's shape or a directory's name. The only
// thing that resolves a path is an extension a package declared.
func TestUnclaimedExtensionIsAnError(t *testing.T) {
	set := loadPackages(t, generators())

	_, _, err := Paths(set, records(t, "app/views/users/index.erb"), nil)
	wantErr(t, err, "app/views/users/index.erb", ".erb")
}

// A path with no extension resolves through no rule at all. There is no default
// package and no fallback: guessing here is exactly the language knowledge the
// core is not allowed to hold.
func TestPathWithNoExtensionIsAnError(t *testing.T) {
	set := loadPackages(t, generators())

	_, _, err := Paths(set, records(t, "Makefile"), nil)
	wantErr(t, err, "Makefile")
}

// contested adds a second package claiming .rb, which is legal at load time and
// only becomes a problem when a .rb path appears.
func contested() map[string]string {
	files := generators()
	files["sinatra/sedum.yaml"] = "name: sinatra\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n"
	files["sinatra/files/app/{name}.rb"] = "# sedum:anchor:top\n"
	files["sinatra/actions/actions.yaml"] = "actions: {}\n"
	return files
}

func TestContestedExtensionWithoutALangFlagNamesBothPackages(t *testing.T) {
	set := loadPackages(t, contested())

	_, _, err := Paths(set, records(t, "app/models/user.rb"), nil)
	wantErr(t, err, "app/models/user.rb", ".rb", "rails", "sinatra", "--lang")
}

func TestLangFlagResolvesAContestedExtension(t *testing.T) {
	set := loadPackages(t, contested())

	resolutions, _, err := Paths(set, records(t, "app/models/user.rb"), []string{"sinatra"})
	if err != nil {
		t.Fatalf("Packages: %v", err)
	}
	if p := resolutions[0].Package; p == nil || p.Name != "sinatra" {
		t.Fatalf("--lang sinatra resolved to %v", p)
	}
}

// Preferring both candidates is not a disambiguation.
func TestLangFlagNamingBothCandidatesStillErrors(t *testing.T) {
	set := loadPackages(t, contested())

	_, _, err := Paths(set, records(t, "app/models/user.rb"), []string{"rails", "sinatra"})
	wantErr(t, err, "rails", "sinatra")
}

// A preference that cannot be honored says so and does not refuse the work
// (prov-2026-fe1e68b8).
func TestUnusableLangFlagsWarn(t *testing.T) {
	set := loadPackages(t, generators())

	tests := []struct {
		name     string
		prefer   string
		fragment string
	}{
		{"absent package", "react", "generators directory"},
		{"claims no extension in the record set", "chi", "no extension"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, warnings, err := Paths(set, records(t, "app/models/user.rb"), []string{tc.prefer})
			if err != nil {
				t.Fatalf("an unusable --lang refused the run: %v", err)
			}
			if !containsFragment(warnings, tc.prefer) || !containsFragment(warnings, tc.fragment) {
				t.Errorf("warnings = %v, want one naming %q and %q", warnings, tc.prefer, tc.fragment)
			}
		})
	}
}

// created runs both phases and returns the files that ended up on disk.
func created(t *testing.T, set *genpkg.Set, recs *record.Set, out string) []File {
	t.Helper()
	resolutions, _, err := Paths(set, recs, nil)
	if err != nil {
		t.Fatalf("Packages: %v", err)
	}
	files, err := Create(resolutions, Options{Output: out})
	if err != nil {
		t.Fatalf("Files: %v", err)
	}
	return files
}

func TestFilesAreCreatedFromTheMatchedTemplate(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()

	files := created(t, set, records(t, "app/controllers/users_controller.rb"), out)

	if len(files) != 1 {
		t.Fatalf("created %d files, want 1", len(files))
	}
	if want := "app/controllers/{name}_controller.rb"; files[0].Template != want {
		t.Errorf("template = %q, want %q", files[0].Template, want)
	}
	if got := files[0].Captures["name"]; got != "users" {
		t.Errorf("captures = %v, want name=users", files[0].Captures)
	}

	// The capture is bound and then run through the package's own pipeline -
	// constantize is [singular, pascal] here - so the rendered name proves
	// both that the capture reached the template and that the package's
	// transforms, not Sedum's opinions, shaped it.
	body := readFile(t, filepath.Join(out, "app/controllers/users_controller.rb"))
	if !strings.Contains(body, "class UserController") {
		t.Errorf("file body = %q, template was not rendered with its captures", body)
	}
	if !strings.Contains(body, "# sedum:anchor:class_body") {
		t.Errorf("file body = %q, the template's marker was not planted", body)
	}
}

// Sedum touches only what a record authorized. No sibling expansion, no
// companion files, nothing inferred from a convention.
func TestNothingUnauthorizedIsCreated(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()

	created(t, set, records(t, "app/models/user.rb"), out)

	want := []string{"app/models/user.rb"}
	if got := tree(t, out); !equal(got, want) {
		t.Errorf("output tree = %v, want exactly %v", got, want)
	}
}

// The matcher and its ranking are M1a's; that Phase 3 defers to them rather
// than to directory order is this milestone's.
func TestTheMostSpecificTemplateWins(t *testing.T) {
	files := generators()
	files["rails/files/app/models/admin/{name}.rb"] = "class Admin::{{name|constantize}}\n  # sedum:anchor:class_body\nend\n"
	set := loadPackages(t, files)
	out := t.TempDir()

	got := created(t, set, records(t, "app/models/admin/user.rb"), out)

	if want := "app/models/admin/{name}.rb"; got[0].Template != want {
		t.Errorf("template = %q, want %q", got[0].Template, want)
	}
}

// A path nothing matches falls to the package's _default for that extension
// (prov-2026-326598ac).
func TestNoMatchFallsToDefault(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()

	files := created(t, set, records(t, "lib/tasks/import.rb"), out)

	if !files[0].Default {
		t.Errorf("lib/tasks/import.rb did not fall to _default: template = %q", files[0].Template)
	}
	if body := readFile(t, filepath.Join(out, "lib/tasks/import.rb")); !strings.Contains(body, "sedum:anchor:top") {
		t.Errorf("file body = %q, _default was not rendered", body)
	}
}

// The default is selected by the target's extension and nothing else, so a
// package's _default for one extension never stands in for another
// (prov-2026-326598ac).
func TestDefaultIsSelectedByExtension(t *testing.T) {
	files := generators()
	files["rails/sedum.yaml"] = strings.Replace(files["rails/sedum.yaml"], `[".rb"]`, `[".rb", ".erb"]`, 1)
	set := loadPackages(t, files)
	out := t.TempDir()

	got := created(t, set, records(t, "app/views/users/index.erb"), out)

	if got[0].Template != "" || got[0].Default {
		t.Fatalf("an .erb path used the .rb default: template = %q", got[0].Template)
	}
	if body := readFile(t, filepath.Join(out, "app/views/users/index.erb")); body != "" {
		t.Errorf("file body = %q, want empty", body)
	}
}

// No match and no default is not an error: a migration or a plain config file
// may legitimately start blank. It is a log line, because a file appearing with
// no boilerplate should be explicable.
func TestNoMatchAndNoDefaultCreatesAnEmptyFileAndLogsIt(t *testing.T) {
	files := generators()
	delete(files, "rails/files/_default.rb")
	set := loadPackages(t, files)
	out := t.TempDir()

	var log bytes.Buffer
	resolutions, _, err := Paths(set, records(t, "db/migrate/001_init.rb"), nil)
	if err != nil {
		t.Fatalf("Packages: %v", err)
	}
	if _, err := Create(resolutions, Options{Output: out, Log: mirrorLog(t, &log)}); err != nil {
		t.Fatalf("Files: %v", err)
	}

	if body := readFile(t, filepath.Join(out, "db/migrate/001_init.rb")); body != "" {
		t.Errorf("file body = %q, want empty", body)
	}
	if !strings.Contains(log.String(), "db/migrate/001_init.rb") {
		t.Errorf("run log = %q, does not mention the file created empty", log.String())
	}
}

// Phase 3 is create-if-absent. Re-rendering would destroy the injected regions
// a file already carries, which is what makes stopping and resuming a run an
// ordinary workflow rather than an edge case.
func TestExistingFilesAreNotReRendered(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()
	recs := records(t, "app/models/user.rb")

	created(t, set, recs, out)

	full := filepath.Join(out, "app/models/user.rb")
	edited := readFile(t, full) + "\n# hand-written afterwards\n"
	if err := os.WriteFile(full, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	files := created(t, set, recs, out)

	if !files[0].Existed {
		t.Errorf("a rerun reported the existing file as newly created")
	}
	if got := readFile(t, full); got != edited {
		t.Errorf("a rerun rewrote the file:\n got %q\nwant %q", got, edited)
	}
}

// A file that exists but lacks its template's markers was written by something
// other than Sedum, or its template changed shape after it was generated.
// Either way the injections aimed at it will not land.
func TestExistingFileMissingItsMarkersHalts(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()

	full := filepath.Join(out, "app/models/user.rb")
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte("class User\nend\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	resolutions, _, err := Paths(set, records(t, "app/models/user.rb"), nil)
	if err != nil {
		t.Fatalf("Packages: %v", err)
	}
	_, err = Create(resolutions, Options{Output: out})

	wantErr(t, err, "app/models/user.rb", "app/models/{name}.rb", "class_body")
}

// A dry run reports every decision and writes nothing, which is what makes
// `sedum resolve` a read-only inspection tool over the same code path a real
// run takes.
func TestDryRunWritesNothing(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()

	resolutions, _, err := Paths(set, records(t, "app/models/user.rb"), nil)
	if err != nil {
		t.Fatalf("Packages: %v", err)
	}
	files, err := Create(resolutions, Options{Output: out, DryRun: true})
	if err != nil {
		t.Fatalf("Files: %v", err)
	}

	if !strings.Contains(files[0].Rendered, "class User") {
		t.Errorf("a dry run did not report what it would render: %q", files[0].Rendered)
	}
	if got := tree(t, out); len(got) != 0 {
		t.Errorf("a dry run wrote %v", got)
	}
}

// An authorized path standing where a directory already is cannot be created,
// and silently skipping it would leave a later injection failing for a reason
// nothing explained.
func TestPathOccupiedByADirectoryHalts(t *testing.T) {
	set := loadPackages(t, generators())
	out := t.TempDir()
	if err := os.MkdirAll(filepath.Join(out, "app/models/user.rb"), 0o755); err != nil {
		t.Fatal(err)
	}

	resolutions, _, err := Paths(set, records(t, "app/models/user.rb"), nil)
	if err != nil {
		t.Fatalf("Packages: %v", err)
	}
	_, err = Create(resolutions, Options{Output: out})

	wantErr(t, err, "app/models/user.rb", "directory")
}

// The fixture packages on disk are the worked example of what a generator
// package looks like, and one of them is for a target that does not exist:
// a different extension, comment prefix, pipeline vocabulary, and directory
// shape from every other fixture. If it resolves and generates on the same
// terms as the others, no target-specific knowledge is in the path.
func TestEveryFixturePackageGeneratesOnTheSameTerms(t *testing.T) {
	set, findings, err := genpkg.Load(filepath.Join("..", "..", "testdata", "generators"), genpkg.Options{})
	if err != nil {
		t.Fatalf("genpkg.Load(testdata/generators): %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("fixture package %s is invalid: %s", f.Package, f.Message)
		}
	}
	if len(set.Packages) < 3 {
		t.Fatalf("testdata/generators holds %d packages; the no-target-knowledge case needs a third", len(set.Packages))
	}

	out := t.TempDir()
	recs := records(t,
		"app/controllers/users_controller.rb",
		"internal/handlers/user.go",
		"Units/ingest/Manifest.crn",
	)
	files := created(t, set, recs, out)

	if len(files) != 3 {
		t.Fatalf("created %d files, want 3", len(files))
	}
	for _, f := range files {
		if f.Template == "" {
			t.Errorf("%s matched no file template", f.Path)
		}
		if f.Rendered == "" {
			t.Errorf("%s was created empty", f.Path)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// tree returns every file under root, slash-separated and sorted.
func tree(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		out = append(out, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	sort.Strings(out)
	return out
}

func mirrorLog(t *testing.T, buf *bytes.Buffer) *runlog.Log {
	t.Helper()
	log, err := runlog.NewWithMirror(filepath.Join(t.TempDir(), "run.log"), buf)
	if err != nil {
		t.Fatalf("runlog: %v", err)
	}
	t.Cleanup(func() { log.Close() })
	return log
}

func containsFragment(list []string, want string) bool {
	for _, s := range list {
		if strings.Contains(s, want) {
			return true
		}
	}
	return false
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
