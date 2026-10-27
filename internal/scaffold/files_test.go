package scaffold

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// The scaffold files unit creates or updates the context file, the .gitignore
// and the provenance config exclude list of a directory and reports each as
// created, updated, unchanged or skipped.

func write(t *testing.T, dir, rel, content string, mode os.FileMode) {
	t.Helper()
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(full, mode); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, dir, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, rel))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestContextTemplateIsTheElevenLinesProvenanceInitWrites(t *testing.T) {
	want := "# Context for sedum grow\n" +
		"\n" +
		"Pass with: sedum grow --generators <generators dir> --records provenance --prompt-file .sedum/context.md --only <record id>\n" +
		"Advisory only: behaviour belongs in the record, and editing this file regenerates nothing.\n" +
		"Write here what the model needs but the records should not carry.\n" +
		"\n" +
		"## Conventions\n" +
		"- (package conventions, output rules, error handling)\n" +
		"\n" +
		"## Helpers that exist (call them; do not redefine)\n" +
		"- (signatures the generated code may call)\n"
	if ContextTemplate != want {
		t.Errorf("ContextTemplate =\n%q\nwant\n%q", ContextTemplate, want)
	}
}

func TestEnsureContextCreatesTheFileOnceAndNeverOverwritesIt(t *testing.T) {
	dir := t.TempDir()

	item, err := ensureContext(dir)
	if err != nil || item.Name != ".sedum/context.md" || item.Status != StatusCreated {
		t.Fatalf("first: %+v %v", item, err)
	}
	if read(t, dir, ".sedum/context.md") != ContextTemplate {
		t.Error("the new context file is not the template")
	}
	info, _ := os.Stat(filepath.Join(dir, ".sedum", "context.md"))
	if info.Mode().Perm() != 0o644 {
		t.Errorf("mode = %v, want 0644", info.Mode().Perm())
	}

	write(t, dir, ".sedum/context.md", "my own text", 0o644)
	item, err = ensureContext(dir)
	if err != nil || item.Status != StatusUnchanged || read(t, dir, ".sedum/context.md") != "my own text" {
		t.Errorf("second: %+v %v, content %q", item, err, read(t, dir, ".sedum/context.md"))
	}
}

func TestEnsureGitignore(t *testing.T) {
	cases := []struct {
		name   string
		before *string
		status string
		after  string
	}{
		{"missing", nil, StatusCreated, ".sedum/run.log\n"},
		{"empty file", ptr(""), StatusUpdated, ".sedum/run.log\n"},
		{"ends with a newline", ptr("node_modules\n"), StatusUpdated, "node_modules\n.sedum/run.log\n"},
		{"no final newline", ptr("node_modules"), StatusUpdated, "node_modules\n.sedum/run.log\n"},
		{"already there", ptr("a\n.sedum/run.log\nb\n"), StatusUnchanged, "a\n.sedum/run.log\nb\n"},
		{"there with whitespace", ptr("  .sedum/run.log \t\n"), StatusUnchanged, "  .sedum/run.log \t\n"},
		{"a wider pattern does not count", ptr(".sedum/*\n"), StatusUpdated, ".sedum/*\n.sedum/run.log\n"},
	}
	for _, c := range cases {
		dir := t.TempDir()
		if c.before != nil {
			write(t, dir, ".gitignore", *c.before, 0o644)
		}
		item, err := ensureGitignore(dir)
		if err != nil || item.Name != ".gitignore" || item.Status != c.status || read(t, dir, ".gitignore") != c.after {
			t.Errorf("%s: item %+v err %v content %q", c.name, item, err, read(t, dir, ".gitignore"))
		}
	}
}

func ptr(s string) *string { return &s }

func TestEnsureExcludesWithoutAConfigIsSkipped(t *testing.T) {
	dir := t.TempDir()

	item, err := ensureExcludes(dir)

	if err != nil || item.Name != ".provenance.yml" || item.Status != StatusSkipped || item.Detail != "no .provenance.yml in this directory" {
		t.Errorf("item = %+v err = %v", item, err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".provenance.yml")); !os.IsNotExist(err) {
		t.Error("a config was created")
	}
}

func TestEnsureExcludes(t *testing.T) {
	cases := []struct {
		name   string
		before string
		status string
		after  string
	}{
		{"both present in a flow list", "exclude_paths: ['a', \"generators/**\", .sedum/**]\n", StatusUnchanged, "exclude_paths: ['a', \"generators/**\", .sedum/**]\n"},
		{"both missing from a flow list", "dir: provenance\nexclude_paths: ['README.md']\nother: 1\n", StatusUpdated, "dir: provenance\nexclude_paths: ['README.md', 'generators/**', '.sedum/**']\nother: 1\n"},
		{"one missing from a flow list", "exclude_paths: [generators/**]\n", StatusUpdated, "exclude_paths: [generators/**, '.sedum/**']\n"},
		{"an empty flow list", "exclude_paths: []\n", StatusUpdated, "exclude_paths: ['generators/**', '.sedum/**']\n"},
		{"a flow list with a trailing comment", "exclude_paths: ['a'] # keep this\n", StatusUpdated, "exclude_paths: ['a', 'generators/**', '.sedum/**'] # keep this\n"},
		{"a regex entry holding commas and a bracket", "exclude_paths: ['/a,b]c/', 'x']\n", StatusUpdated, "exclude_paths: ['/a,b]c/', 'x', 'generators/**', '.sedum/**']\n"},
		{"a block list", "exclude_paths:\n  - a\n  - 'generators/**'\nnext: 1\n", StatusUpdated, "exclude_paths:\n  - a\n  - 'generators/**'\n  - '.sedum/**'\nnext: 1\n"},
		{"a block list with comments between items and a deeper indent", "exclude_paths:\n    - a\n    # note\n\n    - b\nnext: 1\n", StatusUpdated, "exclude_paths:\n    - a\n    # note\n\n    - b\n    - 'generators/**'\n    - '.sedum/**'\nnext: 1\n"},
		{"a block list at the end of the file", "exclude_paths:\n  - a", StatusUpdated, "exclude_paths:\n  - a\n  - 'generators/**'\n  - '.sedum/**'\n"},
		{"an empty block list", "exclude_paths:\nnext: 1\n", StatusUpdated, "exclude_paths:\n  - 'generators/**'\n  - '.sedum/**'\nnext: 1\n"},
		{"no key", "dir: provenance\n", StatusUpdated, "dir: provenance\nexclude_paths: ['generators/**', '.sedum/**']\n"},
		{"no key and no final newline", "dir: provenance", StatusUpdated, "dir: provenance\nexclude_paths: ['generators/**', '.sedum/**']\n"},
	}
	for _, c := range cases {
		dir := t.TempDir()
		write(t, dir, ".provenance.yml", c.before, 0o640)

		item, err := ensureExcludes(dir)

		if err != nil || item.Name != ".provenance.yml" || item.Status != c.status || read(t, dir, ".provenance.yml") != c.after {
			t.Errorf("%s: item %+v err %v content %q", c.name, item, err, read(t, dir, ".provenance.yml"))
		}
		info, _ := os.Stat(filepath.Join(dir, ".provenance.yml"))
		if info.Mode().Perm() != 0o640 {
			t.Errorf("%s: mode = %v, want the original 0640", c.name, info.Mode().Perm())
		}
	}
}

func TestEnsureExcludesRefusesOtherFormsWithoutWriting(t *testing.T) {
	for name, before := range map[string]string{
		"a flow list over several lines": "exclude_paths: [\n  'a',\n]\n",
		"a scalar":                       "exclude_paths: nope\n",
	} {
		dir := t.TempDir()
		write(t, dir, ".provenance.yml", before, 0o644)

		_, err := ensureExcludes(dir)

		if err == nil || err.Error() != "exclude_paths is not a one-line list or a block list; add generators/** and .sedum/** to it by hand" {
			t.Errorf("%s: err = %v", name, err)
		}
		if read(t, dir, ".provenance.yml") != before {
			t.Errorf("%s: the file was written", name)
		}
	}
}

func TestInitReturnsTheThreeItemsInOrderAndIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, ".provenance.yml", "exclude_paths: ['a']\n", 0o644)

	first, err := Init(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names, statuses []string
	for _, it := range first {
		names, statuses = append(names, it.Name), append(statuses, it.Status)
	}
	if !reflect.DeepEqual(names, []string{".sedum/context.md", ".gitignore", ".provenance.yml"}) ||
		!reflect.DeepEqual(statuses, []string{StatusCreated, StatusCreated, StatusUpdated}) {
		t.Errorf("first run: %v %v", names, statuses)
	}

	second, err := Init(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, it := range second {
		if it.Status != StatusUnchanged {
			t.Errorf("second run: %+v, want everything unchanged", it)
		}
	}
}

func TestInitWithoutAConfigReportsTheConfigSkipped(t *testing.T) {
	items, err := Init(t.TempDir())
	if err != nil || len(items) != 3 || items[2].Status != StatusSkipped {
		t.Errorf("items = %+v err = %v", items, err)
	}
}

func TestInitWrapsAStepErrorWithTheItemNameAndStopsThere(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, ".provenance.yml", "exclude_paths: nope\n", 0o644)

	items, err := Init(dir)

	if items != nil || err == nil || !strings.HasPrefix(err.Error(), ".provenance.yml: ") {
		t.Errorf("items = %+v err = %v", items, err)
	}
	// The steps before the failing one already ran.
	if read(t, dir, ".gitignore") != ".sedum/run.log\n" {
		t.Error("the earlier steps did not run")
	}
}

func TestInitStopsAtTheFirstFailingStepWithoutRunningTheLaterOnes(t *testing.T) {
	dir := t.TempDir()
	// A file named .sedum makes the context step fail.
	write(t, dir, ".sedum", "not a directory", 0o644)

	items, err := Init(dir)

	if items != nil || err == nil || !strings.HasPrefix(err.Error(), ".sedum/context.md: ") {
		t.Fatalf("items = %+v err = %v", items, err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, ".gitignore")); !os.IsNotExist(statErr) {
		t.Error("a later step ran")
	}
}
