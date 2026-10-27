package gitcheck

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/livecodelife/sedum/internal/pipeline"
	"github.com/livecodelife/sedum/internal/selection"
)

// The classes unit exports a tree and computes the class of each record on it
// with the grow --check pipeline.

func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

var railsPackage = map[string]string{
	"generators/rails/sedum.yaml":                 "name: rails\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\ntransforms:\n  constantize: [singular, pascal]\n",
	"generators/rails/files/app/models/{name}.rb": "class {{name|constantize}}\n  # sedum:anchor:class_body\nend\n",
	"generators/rails/actions/actions.yaml":       "actions:\n  addField:\n    kwargs:\n      model: { type: string, required: true }\n      field: { type: string, required: true }\n    injects_into: \"app/models/{{model|snake}}.rb\"\n    anchor: class_body\n",
	"generators/rails/actions/addField.rb":        "attribute :{{field}}\n",
}

func record(id, intent, scope string) map[string]string {
	return map[string]string{"provenance/" + id + ".yml": "id: " + id + "\nintent: |\n  " + intent + "\naffected_scope:\n" + scope}
}

type cannedModel struct{}

func (cannedModel) Complete(context.Context, []selection.Message) (selection.Completion, error) {
	return selection.Completion{Content: `{"invocations": [{"action": "addField", "kwargs": {"model": "user", "field": "email"}}]}`}, nil
}

// grownTree is a tree whose one record has been grown, so its region is current.
func grownTree(t *testing.T, scope string) string {
	t.Helper()
	tree := t.TempDir()
	writeFiles(t, tree, railsPackage)
	writeFiles(t, tree, record("prov-2026-aaaaaaaa", "Add email.", scope))
	result, err := pipeline.Run(context.Background(), pipeline.Config{
		Generators: filepath.Join(tree, "generators"),
		Records:    filepath.Join(tree, "provenance"),
		Output:     tree,
		Client:     cannedModel{},
	})
	if err != nil {
		t.Fatalf("grow: %v", err)
	}
	if err := pipeline.WriteRecordings(result, filepath.Join(tree, ".sedum", "recordings")); err != nil {
		t.Fatalf("WriteRecordings: %v", err)
	}
	return tree
}

func TestClassesReportsACurrentRecordAndAStaleOne(t *testing.T) {
	tree := grownTree(t, "  - app/models/user.rb\n")

	classes, skipped, err := Classes(tree, filepath.Join(tree, "generators"), filepath.Join(tree, "provenance"))

	if err != nil || classes["prov-2026-aaaaaaaa"] != "current" || len(skipped) != 0 {
		t.Fatalf("classes = %v skipped = %v err = %v, want current", classes, skipped, err)
	}

	writeFiles(t, tree, record("prov-2026-aaaaaaaa", "Something else entirely.", "  - app/models/user.rb\n"))
	classes, _, err = Classes(tree, filepath.Join(tree, "generators"), filepath.Join(tree, "provenance"))
	if err != nil || classes["prov-2026-aaaaaaaa"] != "stale" {
		t.Errorf("after the record changed: classes = %v err = %v, want stale", classes, err)
	}
}

func TestClassesReportsUnknownForARecordWithNoRegionsAndSkipsUnclaimedPaths(t *testing.T) {
	tree := t.TempDir()
	writeFiles(t, tree, railsPackage)
	writeFiles(t, tree, record("prov-2026-aaaaaaaa", "Add email.", "  - app/models/user.rb\n  - docs/NOTES.md\n  - Makefile\n"))
	writeFiles(t, tree, record("prov-2026-bbbbbbbb", "Docs.", "  - docs/NOTES.md\n"))

	classes, skipped, err := Classes(tree, filepath.Join(tree, "generators"), filepath.Join(tree, "provenance"))

	if err != nil {
		t.Fatalf("Classes: %v", err)
	}
	if classes["prov-2026-aaaaaaaa"] != "unknown" {
		t.Errorf("class = %q, want unknown (no regions yet)", classes["prov-2026-aaaaaaaa"])
	}
	if want := []string{"Makefile", "docs/NOTES.md"}; !reflect.DeepEqual(skipped, want) {
		t.Errorf("skipped = %v, want each distinct unclaimed path once: %v", skipped, want)
	}
}

func TestAMissingGeneratorsDirectoryMeansNoPackages(t *testing.T) {
	tree := t.TempDir()
	writeFiles(t, tree, record("prov-2026-aaaaaaaa", "Add email.", "  - app/models/user.rb\n"))

	classes, skipped, err := Classes(tree, filepath.Join(tree, "generators"), filepath.Join(tree, "provenance"))

	if err != nil || classes["prov-2026-aaaaaaaa"] != "none" || !reflect.DeepEqual(skipped, []string{"app/models/user.rb"}) {
		t.Errorf("classes = %v skipped = %v err = %v", classes, skipped, err)
	}
}

func TestClassesFailsWhenTheRecordsDirectoryIsMissing(t *testing.T) {
	tree := t.TempDir()
	if _, _, err := Classes(tree, filepath.Join(tree, "generators"), filepath.Join(tree, "provenance")); err == nil {
		t.Error("a missing records directory was not an error")
	}
}

func TestExportTreeExportsTheWholeTopLevelTreeOfACommit(t *testing.T) {
	r := newRepo(t)
	first := r.commit("a", map[string]string{"app/x.txt": "old", "top.txt": "t"})
	r.commit("b", map[string]string{"app/x.txt": "new"})

	dir, cleanup, err := ExportTree(filepath.Join(r.dir, "app"), first, false, r.env)
	if err != nil {
		t.Fatal(err)
	}
	for rel, want := range map[string]string{"app/x.txt": "old", "top.txt": "t"} {
		got, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
		if err != nil || string(got) != want {
			t.Errorf("%s = %q, %v; want %q", rel, got, err, want)
		}
	}
	cleanup()
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("cleanup left %s behind", dir)
	}
}

func TestExportTreeExportsTheIndexForAStagedChangeAndWritesNoObject(t *testing.T) {
	r := newRepo(t)
	r.commit("a", map[string]string{"x.txt": "head"})
	r.write("x.txt", "staged")
	r.git("add", "x.txt")
	r.write("x.txt", "working tree")
	before := r.git("count-objects", "-v")

	dir, cleanup, err := ExportTree(r.dir, "", true, r.env)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	got, err := os.ReadFile(filepath.Join(dir, "x.txt"))
	if err != nil || string(got) != "staged" {
		t.Errorf("x.txt = %q, %v; want the staged content", got, err)
	}
	if after := r.git("count-objects", "-v"); after != before {
		t.Errorf("git objects changed:\n%s\n->\n%s", before, after)
	}
}
