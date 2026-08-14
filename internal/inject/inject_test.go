package inject

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/genpkg"
)

// Phase 7 is where a rendered template stops being text and becomes a change to
// a file. The cases below split along the two things that makes true: where
// content lands the first time, and what happens to it on every run after.

// generators is one package with an action per anchor kind, so that each kind's
// placement is provable against the same file rather than against a fixture
// shaped to suit it.
func generators() map[string]string {
	return map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
`,
		"rails/files/app/controllers/{name}_controller.rb": "class {{name|constantize}}Controller\n" +
			"  # sedum:anchor:class_body_top\n" +
			"\n" +
			"  # sedum:anchor:class_body\n" +
			"end\n",

		"rails/actions/actions.yaml": `actions:
  addBeforeFilter:
    kwargs:
      controller: { type: string, required: true }
      filter: { type: string, required: true }
      only: { type: list, required: false }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body_top

  createControllerMethod:
    kwargs:
      controller: { type: string, required: true }
      name: { type: string, required: true }
      collection: { type: string, required: false }
    discriminator: name
    variants: [index, show]
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body

  addHeader:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: start_of_file

  addTrailer:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: end_of_file

  addAfterClass:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: after_match
    anchor_pattern: "(?m)^class .*Controller$"

  addBeforeClass:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: before_match
    anchor_pattern: "(?m)^class .*Controller$"

  addInRegion:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: region
    anchor_start: class_body_top
    anchor_end: class_body

  addUnplanted:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body_bottom
`,
		"rails/actions/addBeforeFilter.rb": "before_action :{{filter}}\n",
		"rails/actions/createControllerMethod/index.rb": "def index\n" +
			"  {{collection|constantize}}.all\n" +
			"end\n",
		"rails/actions/createControllerMethod/show.rb":     "def show\nend\n",
		"rails/actions/createControllerMethod/_default.rb": "def {{name|snake}}\nend\n",
		"rails/actions/addHeader.rb":                       "# frozen_string_literal: true\n",
		"rails/actions/addTrailer.rb":                      "# end of file\n",
		"rails/actions/addAfterClass.rb":                   "# after\n",
		"rails/actions/addBeforeClass.rb":                  "# before\n",
		"rails/actions/addInRegion.rb":                     "# in region\n",
		"rails/actions/addUnplanted.rb":                    "# unplanted\n",
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

// loadPackage loads the fixture package, failing the test if it does not load
// cleanly: a case proving injection should never be debugging a package.
func loadPackage(t *testing.T) *genpkg.Package {
	t.Helper()

	set, findings, err := genpkg.Load(writeTree(t, generators()), genpkg.Options{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("fixture package does not load: %s", f)
		}
	}
	pkg, ok := set.Lookup("rails")
	if !ok {
		t.Fatal("fixture package rails did not load")
	}
	return pkg
}

const controllerPath = "app/controllers/users_controller.rb"

// output writes the scaffolded controller into a fresh output tree, the way
// Phase 3 would have left it.
func output(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	full := filepath.Join(root, filepath.FromSlash(controllerPath))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	scaffold := "class UsersController\n" +
		"  # sedum:anchor:class_body_top\n" +
		"\n" +
		"  # sedum:anchor:class_body\n" +
		"end\n"
	if err := os.WriteFile(full, []byte(scaffold), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return root
}

func read(t *testing.T, root string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(controllerPath)))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	return string(data)
}

// invocation builds one ready-to-write invocation against the fixture package.
func invocation(t *testing.T, pkg *genpkg.Package, name, variant, content, record string, kwargs map[string]any) Invocation {
	t.Helper()

	action, ok := pkg.Actions[name]
	if !ok {
		t.Fatalf("fixture package declares no action %q", name)
	}
	return Invocation{
		Package: pkg, Action: action, Variant: variant,
		Kwargs: kwargs, Path: controllerPath, RecordID: record, Content: content,
	}
}

func TestInjectsAtMarkerAnchor(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	inv := invocation(t, pkg, "createControllerMethod", "index", "def index\n  User.all\nend\n", "PR-014",
		map[string]any{"controller": "users", "name": "index", "collection": "users"})

	if _, err := Apply([]Invocation{inv}, Options{Output: root}); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	got := read(t, root)
	want := "class UsersController\n" +
		"  # sedum:anchor:class_body_top\n" +
		"\n" +
		"  # sedum:anchor:class_body\n" +
		`# sedum:createControllerMethod:index {"tier":"owned","record":"PR-014","kwargs":{"collection":"users","controller":"users","name":"index"}}` + "\n" +
		"def index\n  User.all\nend\n" +
		"# /sedum:createControllerMethod:index\n" +
		"end\n"

	if got != want {
		t.Errorf("injected file:\n%s\nwant:\n%s", got, want)
	}
}

// A fixture invocation list produces byte-identical output on repeat runs.
func TestRerunIsByteIdentical(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	invs := []Invocation{
		invocation(t, pkg, "createControllerMethod", "index", "def index\nend\n", "PR-014",
			map[string]any{"controller": "users", "name": "index"}),
		invocation(t, pkg, "addBeforeFilter", "", "before_action :authenticate\n", "PR-014",
			map[string]any{"controller": "users", "filter": "authenticate", "only": []string{"index"}}),
	}

	if _, err := Apply(invs, Options{Output: root}); err != nil {
		t.Fatalf("first Apply: %v", err)
	}
	first := read(t, root)

	if _, err := Apply(invs, Options{Output: root}); err != nil {
		t.Fatalf("second Apply: %v", err)
	}
	if second := read(t, root); second != first {
		t.Errorf("a rerun changed the file:\n%s\nwant:\n%s", second, first)
	}
}

// Re-running replaces the region an action owns rather than appending beside
// it. This is the parameter-space update: the invocation is the same action
// with different kwargs, rendering is deterministic, and no code is edited by
// anything.
func TestRerunReplacesOwnedRegionRatherThanDuplicating(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	kwargs := map[string]any{"controller": "users", "name": "index", "collection": "users"}
	first := invocation(t, pkg, "createControllerMethod", "index", "def index\n  User.all\nend\n", "PR-014", kwargs)

	if _, err := Apply([]Invocation{first}, Options{Output: root}); err != nil {
		t.Fatalf("first Apply: %v", err)
	}

	// A later record refines the region: same action, same variant, a
	// different non-selecting kwarg.
	refined := invocation(t, pkg, "createControllerMethod", "index",
		"def index\n  User.where(active: true)\nend\n", "PR-092",
		map[string]any{"controller": "users", "name": "index", "collection": "admins"})

	results, err := Apply([]Invocation{refined}, Options{Output: root})
	if err != nil {
		t.Fatalf("second Apply: %v", err)
	}
	if len(results) != 1 || !results[0].Replaced {
		t.Fatalf("results = %+v, want the region reported as replaced", results)
	}

	got := read(t, root)
	if n := strings.Count(got, "def index"); n != 1 {
		t.Errorf("file defines index %d times, want 1; the region was appended beside rather than replaced:\n%s", n, got)
	}
	if !strings.Contains(got, "User.where(active: true)") {
		t.Errorf("the refined body is not in the file:\n%s", got)
	}

	// The marker records who last wrote the region, so the audit trail
	// improves rather than degrading.
	if !strings.Contains(got, `"record":"PR-092"`) {
		t.Errorf("marker does not name the record that last wrote the region:\n%s", got)
	}
	if strings.Contains(got, "PR-014") {
		t.Errorf("marker still names the record that no longer parameterizes the region:\n%s", got)
	}
}

// Two invocations of one action whose selecting kwargs differ are different
// regions and coexist.
func TestDistinctRegionsCoexist(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	invs := []Invocation{
		invocation(t, pkg, "addBeforeFilter", "", "before_action :authenticate\n", "PR-014",
			map[string]any{"controller": "users", "filter": "authenticate"}),
		invocation(t, pkg, "addBeforeFilter", "", "before_action :authorize\n", "PR-014",
			map[string]any{"controller": "users", "filter": "authorize"}),
	}

	if _, err := Apply(invs, Options{Output: root}); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	got := read(t, root)
	for _, want := range []string{"before_action :authenticate", "before_action :authorize"} {
		if !strings.Contains(got, want) {
			t.Errorf("%q is not in the file:\n%s", want, got)
		}
	}

	// Applied in declaration order rather than in reverse.
	if strings.Index(got, ":authenticate") > strings.Index(got, ":authorize") {
		t.Errorf("regions accumulated in reverse order:\n%s", got)
	}
}

// A seeded region is generated once and never touched again.
func TestSeededRegionIsLeftAlone(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	kwargs := map[string]any{"controller": "users", "name": "show"}
	seed := invocation(t, pkg, "createControllerMethod", "show", "def show\n  # fill me in\nend\n", "PR-014", kwargs)
	seed.Tier = TierSeeded

	if _, err := Apply([]Invocation{seed}, Options{Output: root}); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	// A human completes the stub.
	full := filepath.Join(root, filepath.FromSlash(controllerPath))
	completed := strings.Replace(read(t, root), "  # fill me in", "  render json: @user", 1)
	if err := os.WriteFile(full, []byte(completed), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	results, err := Apply([]Invocation{seed}, Options{Output: root})
	if err != nil {
		t.Fatalf("rerun Apply: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Fatalf("results = %+v, want the region reported as skipped", results)
	}

	if got := read(t, root); got != completed {
		t.Errorf("a seeded region was overwritten on rerun:\n%s\nwant:\n%s", got, completed)
	}
}

// A missing anchor is a hard error. It means the file is not shaped the way the
// action assumed, and auto-creating the anchor would paper over exactly the
// mistake worth surfacing.
func TestMissingAnchorIsAHardError(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)

	inv := invocation(t, pkg, "addUnplanted", "", "# unplanted\n", "PR-014",
		map[string]any{"controller": "users"})

	_, err := Apply([]Invocation{inv}, Options{Output: root})
	if err == nil {
		t.Fatal("an action anchored to a marker the file does not carry was applied")
	}

	// The diagnostic names the action, the file, and the rule violated.
	for _, want := range []string{"addUnplanted", controllerPath, "class_body_bottom"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not name %q: %v", want, err)
		}
	}

	if got := read(t, root); strings.Contains(got, "unplanted") {
		t.Errorf("a failed injection still wrote to the file:\n%s", got)
	}
}

// Nothing is created that a provenance record did not authorize, so an action
// targeting a path no record named fails here rather than conjuring the file.
func TestInjectingIntoAnUncreatedFileIsAHardError(t *testing.T) {
	pkg := loadPackage(t)
	root := t.TempDir()

	inv := invocation(t, pkg, "addBeforeFilter", "", "before_action :authenticate\n", "PR-014",
		map[string]any{"controller": "users", "filter": "authenticate"})

	_, err := Apply([]Invocation{inv}, Options{Output: root})
	if err == nil {
		t.Fatal("an action injecting into a file no record authorized was applied")
	}
	if !strings.Contains(err.Error(), controllerPath) {
		t.Errorf("error does not name the file that was never created: %v", err)
	}
}

// Each anchor kind places text where the author would expect, with no parser
// and no syntax awareness of any kind.
func TestAnchorKinds(t *testing.T) {
	cases := []struct {
		action string
		body   string
		// after names lines the injected body must follow, and before
		// names lines it must precede. Placement is asserted by ordering
		// rather than by adjacency, because what matters is that the
		// content landed in the right part of the file.
		after  []string
		before []string
	}{{
		action: "addHeader", body: "# frozen_string_literal: true",
		before: []string{"class UsersController"},
	}, {
		action: "addBeforeClass", body: "# before",
		before: []string{"class UsersController"},
	}, {
		action: "addAfterClass", body: "# after",
		after:  []string{"class UsersController"},
		before: []string{"# sedum:anchor:class_body_top"},
	}, {
		// A region accumulates at its end, just inside the marker that
		// closes it.
		action: "addInRegion", body: "# in region",
		after:  []string{"# sedum:anchor:class_body_top"},
		before: []string{"# sedum:anchor:class_body"},
	}, {
		action: "addTrailer", body: "# end of file",
		after:  []string{"class UsersController", "end"},
	}}

	for _, tc := range cases {
		t.Run(tc.action, func(t *testing.T) {
			pkg := loadPackage(t)
			root := output(t)

			inv := invocation(t, pkg, tc.action, "", tc.body+"\n", "PR-014",
				map[string]any{"controller": "users"})
			if _, err := Apply([]Invocation{inv}, Options{Output: root}); err != nil {
				t.Fatalf("Apply: %v", err)
			}

			got := read(t, root)
			lines := strings.Split(got, "\n")

			at := indexOfLine(lines, tc.body)
			if at < 0 {
				t.Fatalf("injected body is not in the file:\n%s", got)
			}

			// The region's opening marker is the line before its body,
			// and its closing marker the line after.
			if at == 0 || !strings.Contains(lines[at-1], "sedum:"+tc.action) {
				t.Fatalf("body is not wrapped in its opening marker:\n%s", got)
			}
			if at+1 >= len(lines) || !strings.Contains(lines[at+1], "/sedum:"+tc.action) {
				t.Fatalf("body is not wrapped in its closing marker:\n%s", got)
			}

			for _, want := range tc.after {
				if i := indexOfLine(lines, want); i < 0 || i > at-1 {
					t.Errorf("region does not follow %q:\n%s", want, got)
				}
			}
			for _, want := range tc.before {
				if i := indexOfLine(lines, want); i < 0 || i < at+1 {
					t.Errorf("region does not precede %q:\n%s", want, got)
				}
			}
		})
	}
}

func indexOfLine(lines []string, want string) int {
	for i, line := range lines {
		if strings.TrimSpace(line) == want {
			return i
		}
	}
	return -1
}

// The anchor marker's shape is declared in genpkg, which writes it and checks
// for it, and repeated here, which locates the line it sits on. This asserts
// the two have not drifted apart: what this package searches for is what that
// one recognizes.
func TestAnchorMarkerShapesAgree(t *testing.T) {
	for _, prefix := range []string{"#", "//", ";;"} {
		planted := anchorDecl(prefix, "class_body")

		found := genpkg.MarkersIn(prefix, planted)
		if len(found) != 1 || found[0] != "class_body" {
			t.Errorf("prefix %q: genpkg does not recognize %q as planting an anchor; found %v",
				prefix, planted, found)
		}

		offset, ok := findAnchorLineStart(prefix, "class_body", "x\n"+planted+"\ny\n")
		if !ok {
			t.Errorf("prefix %q: %q was not located", prefix, planted)
			continue
		}
		if offset != 2 {
			t.Errorf("prefix %q: anchor located at %d, want 2", prefix, offset)
		}
	}
}

// A marker name is compared for equality rather than as a prefix, or an action
// anchored to class_body would land at class_body_top.
func TestAnchorNamesAreNotPrefixes(t *testing.T) {
	content := "class UsersController\n" +
		"  # sedum:anchor:class_body_top\n" +
		"\n" +
		"  # sedum:anchor:class_body\n" +
		"end\n"

	top, ok := findAnchorLineStart("#", "class_body_top", content)
	if !ok {
		t.Fatal("class_body_top was not located")
	}
	body, ok := findAnchorLineStart("#", "class_body", content)
	if !ok {
		t.Fatal("class_body was not located")
	}
	if top == body {
		t.Error("class_body matched the class_body_top line")
	}
}

// A dry run decides everything and writes nothing.
func TestDryRunWritesNothing(t *testing.T) {
	pkg := loadPackage(t)
	root := output(t)
	before := read(t, root)

	inv := invocation(t, pkg, "createControllerMethod", "index", "def index\nend\n", "PR-014",
		map[string]any{"controller": "users", "name": "index"})

	results, err := Apply([]Invocation{inv}, Options{Output: root, DryRun: true})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("results = %+v, want one", results)
	}
	if got := read(t, root); got != before {
		t.Errorf("a dry run wrote to the file:\n%s", got)
	}
}
