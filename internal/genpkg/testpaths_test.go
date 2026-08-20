package genpkg

import (
	"reflect"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/pathpat"
)

// test_paths is a declaration Sedum validates and does not act on
// (prov-2026-3f01a02d). Both halves of that sentence are load-bearing and both
// are asserted here.
//
// The validating half is ordinary: the key decodes, and an entry that cannot be
// read is an error naming it, exactly as an unmanaged entry is. The not-acting
// half is the one worth a test, because it is the kind of property that erodes
// without anyone deciding to erode it. A later change that taught some phase to
// consult the field would still pass every case below except the last, which is
// why the last one compares two resolved packages rather than inspecting one.

// withTestPaths returns the valid rails package with a test_paths declaration
// carrying both shapes the record argues for: a directory-rooted pattern, and
// one that cannot be expressed as a directory at all.
func withTestPaths(entries string) map[string]string {
	return mutated(map[string]*string{
		"rails/sedum.yaml": text(`name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
  instantize: [plural, "prefix:@"]
test_paths:
` + entries),
	})
}

func TestTestPathsLoadsAndIsCarriedInOrder(t *testing.T) {
	set, findings := loadTree(t, withTestPaths("  - \"spec/**\"\n  - \"**/*_test.go\"\n"))

	for _, f := range findings {
		t.Errorf("a package declaring test_paths reported %s: %s", f.Rule, f.Message)
	}
	if len(set.Packages) != 1 {
		t.Fatalf("loaded %d packages, want 1", len(set.Packages))
	}

	want := []string{"spec/**", "**/*_test.go"}
	if got := set.Packages[0].TestPaths; !reflect.DeepEqual(got, want) {
		t.Errorf("TestPaths = %v, want %v", got, want)
	}
}

// The Go case is the whole argument for patterns over a directory, so it gets
// its own assertion rather than riding along in the list above: a sibling
// _test.go file shares its extension, its package, and its directory with the
// thing it tests, and no directory-shaped declaration can separate them.
func TestTestPathsExpressesASiblingTestFile(t *testing.T) {
	set, _ := loadTree(t, withTestPaths("  - \"**/*_test.go\"\n"))

	pkg := set.Packages[0]
	for _, target := range []string{"internal/user/user_test.go", "user_test.go"} {
		if _, ok := pathpat.MatchAny(pkg.TestPaths, target); !ok {
			t.Errorf("%s does not match a test_paths entry declared as **/*_test.go", target)
		}
	}
	if _, ok := pathpat.MatchAny(pkg.TestPaths, "internal/user/user.go"); ok {
		t.Error("user.go matched a test_paths entry, which would make the partition useless")
	}
}

func TestTestPathsAbsentIsNotAnError(t *testing.T) {
	set, findings := loadTree(t, validPackage())

	for _, f := range findings {
		t.Errorf("a package declaring no test_paths reported %s: %s", f.Rule, f.Message)
	}
	if got := set.Packages[0].TestPaths; len(got) != 0 {
		t.Errorf("TestPaths = %v for a package that declares none, want empty", got)
	}
}

func TestTestPathsEmptyEntryIsAnError(t *testing.T) {
	_, findings := loadTree(t, withTestPaths("  - \"\"\n"))

	f := findingFor(t, findings, RuleTestPathsInvalid)
	if f.Kind != KindError {
		t.Errorf("kind = %v, want error", f.Kind)
	}
	if f.File != manifestFile {
		t.Errorf("file = %q, want %q", f.File, manifestFile)
	}
}

func TestTestPathsUnreadableEntryIsAnErrorNamingIt(t *testing.T) {
	_, findings := loadTree(t, withTestPaths("  - \"spec/[unclosed/**\"\n"))

	f := findingFor(t, findings, RuleTestPathsInvalid)
	if f.Kind != KindError {
		t.Errorf("kind = %v, want error", f.Kind)
	}
	// Naming the entry is the point. A message that said only "invalid pattern"
	// would leave an author with a list to bisect.
	if !strings.Contains(f.Message, "[unclosed") {
		t.Errorf("message does not name the entry that failed: %s", f.Message)
	}
}

// Sedum does not write an unmanaged path, so whether one also counts as a test
// is the reader's question and not Sedum's. Inventing a contradiction check
// here would be inventing a rule neither consumer asked for.
func TestTestPathsMayOverlapUnmanaged(t *testing.T) {
	files := mutated(map[string]*string{
		"rails/sedum.yaml": text(`name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
  instantize: [plural, "prefix:@"]
unmanaged:
  - "spec/spec_helper.rb"
test_paths:
  - "spec/**"
`),
	})

	_, findings := loadTree(t, files)
	for _, f := range findings {
		t.Errorf("a path declared both unmanaged and a test path reported %s: %s", f.Rule, f.Message)
	}
}

// The constraint that nothing acts on the field, expressed as the only thing
// that can actually check it: resolve the same package twice, once with the
// declaration and once without, and require that everything else about it is
// identical. A phase that started consulting test_paths would not necessarily
// break any assertion above; it would break this one as soon as its effect
// reached the resolved package.
func TestTestPathsChangesNothingElseAboutThePackage(t *testing.T) {
	with, _ := loadTree(t, withTestPaths("  - \"spec/**\"\n"))
	without, _ := loadTree(t, validPackage())

	a, b := with.Packages[0], without.Packages[0]

	if !reflect.DeepEqual(a.Extensions, b.Extensions) {
		t.Errorf("extensions differ: %v vs %v", a.Extensions, b.Extensions)
	}
	if !reflect.DeepEqual(a.Unmanaged, b.Unmanaged) {
		t.Errorf("unmanaged differs: %v vs %v", a.Unmanaged, b.Unmanaged)
	}
	if !reflect.DeepEqual(a.FileTemplates, b.FileTemplates) {
		t.Errorf("file templates differ: %v vs %v", a.FileTemplates, b.FileTemplates)
	}
	if a.CommentPrefix != b.CommentPrefix {
		t.Errorf("comment prefix differs: %q vs %q", a.CommentPrefix, b.CommentPrefix)
	}
	if len(a.Actions) != len(b.Actions) {
		t.Fatalf("action count differs: %d vs %d", len(a.Actions), len(b.Actions))
	}
	for name, action := range a.Actions {
		other, ok := b.Actions[name]
		if !ok {
			t.Errorf("action %s exists only in the package declaring test_paths", name)
			continue
		}
		if !reflect.DeepEqual(action, other) {
			t.Errorf("action %s differs between the two packages", name)
		}
	}
}
