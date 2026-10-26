package pipeline

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// A record's scope names every kind of file the change touches: docs, scripts,
// config, files with no extension. Only a file a generator package claims can
// carry markers, so a migration reads the ones it can and skips the rest rather
// than stopping on the first path no package claims.
func TestMigrateSkipsScopePathsNoGeneratorPackageClaims(t *testing.T) {
	f := newMigFixture(t)
	f.cfg.Records = recordsDir(t, map[string]string{
		v2RecordA: "  - app/models/user.rb\n  - docs/NOTES.md\n  - scripts/run.sh\n  - config/app.yml\n  - Makefile\n  - data/blob.json\n",
	})
	f.mc.Records = f.cfg.Records
	gen, pkg := f.current(t)
	f.writeUser(t, f.email(gen, pkg))

	report := f.migrate(t)

	if want := []MigrateNote{{Path: migPath, Action: "addField"}}; !reflect.DeepEqual(report.Migrated, want) {
		t.Errorf("Migrated = %+v, want only the .rb region", report.Migrated)
	}
	if !reflect.DeepEqual(report.Files, []string{migPath}) || !reflect.DeepEqual(report.Recordings, []string{v2RecordA}) {
		t.Errorf("Files = %v, Recordings = %v", report.Files, report.Recordings)
	}
	if report.Skipped != 5 {
		t.Errorf("Skipped = %d, want the five paths no package claims", report.Skipped)
	}
}

// A scope path whose extension is claimed is never skipped, and a scope with
// nothing unclaimed reports none skipped.
func TestMigrateReportsNoSkippedPathsWhenEveryScopePathIsClaimed(t *testing.T) {
	f := newMigFixture(t)
	gen, pkg := f.current(t)
	f.writeUser(t, f.email(gen, pkg))

	report := f.migrate(t)

	if report.Skipped != 0 {
		t.Errorf("Skipped = %d, want 0", report.Skipped)
	}
}

// Two packages can claim one extension. A file's regions belong to the package
// that declares their actions, so a migration must not judge them against
// another package that merely claims the same extension.
func TestMigrateJudgesRegionsAgainstThePackageThatDeclaresTheirActions(t *testing.T) {
	f := newMigFixture(t)
	gen, pkg := f.current(t)
	f.writeUser(t, f.email(gen, pkg))
	other := filepath.Join(f.cfg.Generators, "aaa-other")
	if err := os.MkdirAll(filepath.Join(other, "actions"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(other, "sedum.yaml"), []byte("name: aaa-other\nextensions: [\".rb\"]\ncomment_prefix: \"#\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(other, "actions", "actions.yaml"), []byte("actions:\n  unrelated:\n    kwargs:\n      x: { type: string, required: true }\n    injects_into: \"lib/{{x}}.rb\"\n    anchor: body\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(other, "actions", "unrelated.rb"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.cfg.Records = recordsDir(t, map[string]string{v2RecordA: "  - docs/NOTES.md\n"})
	f.mc.Records = f.cfg.Records

	report := f.migrate(t)

	if len(report.Refused) != 0 {
		t.Errorf("Refused = %+v, want none: the region's action is declared by the rails package", report.Refused)
	}
	if want := []MigrateNote{{Path: migPath, Action: "addField"}}; !reflect.DeepEqual(report.Migrated, want) {
		t.Errorf("Migrated = %+v, want the addField region", report.Migrated)
	}
}

// Whatever packages are loaded decide which scope paths can carry markers: with
// a package that claims another extension loaded, that extension's files are
// processed as normal and only the paths no loaded package claims are skipped.
func TestMigrateProcessesTheFilesOfEveryLoadedPackageAndSkipsTheRest(t *testing.T) {
	f := newMigFixture(t)
	gen, pkg := f.current(t)
	f.writeUser(t, f.email(gen, pkg))

	other := filepath.Join(f.cfg.Generators, "scripts")
	if err := os.MkdirAll(filepath.Join(other, "actions"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"sedum.yaml":           "name: scripts\nextensions: [\".ts\"]\ncomment_prefix: \"//\"\n",
		"actions/actions.yaml": "actions:\n  addConst:\n    kwargs:\n      name: { type: string, required: true }\n    injects_into: \"src/{{name}}.ts\"\n    anchor: body\n",
		"actions/addConst.ts":  "export const {{name}} = 1\n",
	}
	for name, text := range files {
		if err := os.WriteFile(filepath.Join(other, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	const tsPath = "src/limit.ts"
	line := `// sedum:addConst {"tier":"owned","record":"` + v2RecordA + `","gen":"sha256:g","pkg":"sha256:p","kwargs":{"name":"limit"}}` + "\n"
	ts := "// sedum:anchor:body\n" + line + "export const limit = 1\n// /sedum:addConst\n"
	if err := os.MkdirAll(filepath.Join(f.cfg.Output, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.cfg.Output, tsPath), []byte(ts), 0o644); err != nil {
		t.Fatal(err)
	}
	f.cfg.Records = recordsDir(t, map[string]string{v2RecordA: "  - " + migPath + "\n  - " + tsPath + "\n  - docs/NOTES.md\n  - Makefile\n"})
	f.mc.Records = f.cfg.Records

	report := f.migrate(t)

	if len(report.Refused) != 0 || len(report.Left) != 0 {
		t.Errorf("Refused = %+v, Left = %+v", report.Refused, report.Left)
	}
	var paths []string
	for _, n := range report.Migrated {
		paths = append(paths, n.Path)
	}
	if want := []string{migPath, tsPath}; !reflect.DeepEqual(paths, want) {
		t.Errorf("Migrated paths = %v, want %v", paths, want)
	}
	if report.Skipped != 2 {
		t.Errorf("Skipped = %d, want the two paths no loaded package claims", report.Skipped)
	}
}
