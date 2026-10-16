package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/livecodelife/provenance/migrate"
	"github.com/livecodelife/provenance/record"
)

const (
	bpOld = "prov-2026-0000b001"
	bpNew = "prov-2026-0000b002"
	imp   = "prov-2026-0000a001"
)

// migrateFixture: blueprint bpOld superseded by bpNew, imprint imp
// implementing bpOld, a .linespec.yml whose only content is the provenance
// block, and a v1 hash manifest.
func migrateFixture(t *testing.T) string {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".linespec.yml"), "# provenance settings\nprovenance:\n  # where records live\n  dir: provenance\n  commit_on_status_change: false\n")
	writeFile(t, filepath.Join(root, ".linespec/hash_manifest.json"), "{}\n")
	d := filepath.Join(root, "provenance")
	writeFile(t, filepath.Join(d, bpOld+".yml"), strings.Replace(v1(bpOld, "blueprint", "superseded", ""), `superseded_by: ""`, "superseded_by: "+bpNew, 1))
	writeFile(t, filepath.Join(d, bpNew+".yml"), strings.Replace(v1(bpNew, "blueprint", "implemented", ""), `supersedes: ""`, "supersedes: "+bpOld, 1))
	writeFile(t, filepath.Join(d, imp+".yml"), v1(imp, "imprint", "implemented", "implements: "+bpOld+"\n"))
	gitRepo(t, root)
	return root
}

func TestMigratePlanWritesNothing(t *testing.T) {
	root := migrateFixture(t)
	cfg := filepath.Join(root, ".linespec.yml")
	out := mustRun(t, "-c", cfg, "migrate", "--plan", "--json")
	var view planJSON
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatalf("plan json: %v\n%s", err, out)
	}
	if view.RecordsBefore != 3 || view.RecordsAfter != 2 {
		t.Fatalf("records %d → %d, want 3 → 2", view.RecordsBefore, view.RecordsAfter)
	}
	if got := view.Chains[bpNew]; !reflect.DeepEqual(got, []string{bpOld, bpNew}) {
		t.Fatalf("chain %v", got)
	}
	human := mustRun(t, "-c", cfg, "migrate", "--plan")
	for _, want := range []string{"3 records → 2", bpNew + "  ← " + bpOld + ", " + bpNew, "Nothing was written"} {
		if !strings.Contains(human, want) {
			t.Errorf("human plan lacks %q:\n%s", want, human)
		}
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Fatalf("--plan wrote files:\n%s", st)
	}
}

func TestMigrateApply(t *testing.T) {
	root := migrateFixture(t)
	out := mustRun(t, "-c", filepath.Join(root, ".linespec.yml"), "migrate")
	if !strings.Contains(out, "Migrated provenance to v2: 2 records (from 3), 3 snapshots") {
		t.Fatalf("summary:\n%s", out)
	}
	store := &record.FilesStore{Root: root}
	head, _, err := store.Load(bpNew)
	if err != nil {
		t.Fatal(err)
	}
	if head.Schema != 2 || head.Version != 2 || !reflect.DeepEqual(head.Aliases, []string{bpOld}) {
		t.Fatalf("head schema %d version %d aliases %v", head.Schema, head.Version, head.Aliases)
	}
	if _, err := os.Stat(filepath.Join(root, "provenance", bpOld+".yml")); !os.IsNotExist(err) {
		t.Fatalf("folded record file still present (%v)", err)
	}
	child, _, _ := store.Load(imp)
	if child.Implements != bpNew || child.ImplementsVersion != 1 {
		t.Fatalf("child pin %s@%d", child.Implements, child.ImplementsVersion)
	}

	// Config moved, dedented, comments kept; .linespec.yml held nothing else.
	if got := readFile(t, filepath.Join(root, ".provenance.yml")); got != "# where records live\ndir: provenance\ncommit_on_status_change: false\n" {
		t.Fatalf(".provenance.yml:\n%s", got)
	}
	if _, err := os.Stat(filepath.Join(root, ".linespec.yml")); !os.IsNotExist(err) {
		t.Fatalf(".linespec.yml should be removed (%v)", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".linespec/hash_manifest.json")); !os.IsNotExist(err) {
		t.Fatalf("v1 hash manifest should be removed (%v)", err)
	}

	// The manifest has one line per record, from its latest snapshot.
	var lines []manifestLine
	for _, l := range strings.Split(strings.TrimSpace(readFile(t, filepath.Join(root, manifestPath))), "\n") {
		var m manifestLine
		if err := json.Unmarshal([]byte(l), &m); err != nil {
			t.Fatal(err)
		}
		lines = append(lines, m)
	}
	if len(lines) != 2 || lines[0].ID != imp || lines[1].ID != bpNew || lines[1].Version != 2 {
		t.Fatalf("manifest %+v", lines)
	}
	snap, err := store.Snapshot(bpNew, 2)
	if err != nil || lines[1].ContentHash != snap.ContentHash {
		t.Fatalf("manifest hash %s, snapshot %v (%v)", lines[1].ContentHash, snap, err)
	}

	// Migrating again, now from .provenance.yml, finds nothing to do.
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "migrate")
	again := mustRun(t, "-c", filepath.Join(root, ".provenance.yml"), "migrate")
	if !strings.Contains(again, "Nothing to migrate") {
		t.Fatalf("second migrate:\n%s", again)
	}
}

func TestMigrateKeepsOtherLinespecConfig(t *testing.T) {
	root := migrateFixture(t)
	cfg := filepath.Join(root, ".linespec.yml")
	writeFile(t, cfg, "service:\n  name: demo\nprovenance:\n  dir: provenance\n\n# tests\ntests:\n  dir: specs\n")
	gitIn(t, root, "commit", "-qam", "config")
	mustRun(t, "-c", cfg, "migrate")
	if got := readFile(t, cfg); got != "service:\n  name: demo\n\n# tests\ntests:\n  dir: specs\n" {
		t.Fatalf(".linespec.yml:\n%s", got)
	}
	if got := readFile(t, filepath.Join(root, ".provenance.yml")); got != "dir: provenance\n" {
		t.Fatalf(".provenance.yml:\n%s", got)
	}
}

func TestMigrateRefusals(t *testing.T) {
	root := migrateFixture(t)
	cfg := filepath.Join(root, ".linespec.yml")
	writeFile(t, filepath.Join(root, "provenance", bpNew+".yml"), v1(bpNew, "blueprint", "open", ""))
	if _, err := run(t, "-c", cfg, "migrate"); err == nil {
		t.Fatal("migrate ran over a dirty records directory")
	}
	gitIn(t, root, "checkout", "--", ".")
	if _, err := run(t, "-c", cfg, "migrate", "--json"); err == nil {
		t.Fatal("--json without --plan accepted")
	}
	if _, err := run(t, "-c", cfg, "migrate", "--plan", "--keep-separate", imp); err == nil {
		t.Fatal("--keep-separate accepted an id the plan does not fold")
	}
	if _, err := run(t, "-c", cfg, "migrate", "--plan", "--keep-separate", bpOld, "--fold", bpOld); err == nil {
		t.Fatal("--keep-separate and --fold on one id accepted")
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Fatalf("refused migrations wrote files:\n%s", st)
	}
}

// TestMigrateThisRepoPlan: --plan --json on a clone of this repository agrees
// with migrate.MakePlan.
func TestMigrateThisRepoPlan(t *testing.T) {
	src, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	clone := filepath.Join(t.TempDir(), "provenance")
	gitClone(t, src, clone)
	out := mustRun(t, "-c", filepath.Join(clone, ".linespec.yml"), "migrate", "--plan", "--json")
	var view planJSON
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatal(err)
	}
	repo, err := migrate.Load(clone, migrate.LoadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	p, err := migrate.MakePlan(repo)
	if err != nil {
		t.Fatal(err)
	}
	want := planView(p, migrate.Options{KeepSeparate: map[string]bool{}})
	if view.RecordsBefore != len(repo.Records) || view.RecordsAfter != want.RecordsAfter || len(view.Entries) != len(want.Entries) {
		t.Fatalf("plan %d → %d (%d entries), MakePlan %d → %d (%d entries)",
			view.RecordsBefore, view.RecordsAfter, len(view.Entries), len(repo.Records), want.RecordsAfter, len(want.Entries))
	}
}

// TestMigrateRealRepos: on clones of the four real repos, `provenance migrate`
// writes exactly the tree migrate.Apply writes, and the originals are only
// read. Runs when PROVENANCE_REAL_REPOS names their workspace directory.
func TestMigrateRealRepos(t *testing.T) {
	ws := os.Getenv("PROVENANCE_REAL_REPOS")
	if ws == "" {
		t.Skip("set PROVENANCE_REAL_REPOS to the workspace directory to run")
	}
	for _, name := range []string{"linespec", "sedum", "scion", "linespec-provenance-ui"} {
		t.Run(name, func(t *testing.T) {
			before := gitStatus(t, filepath.Join(ws, name), "provenance")
			tmp := t.TempDir()
			cli, lib := filepath.Join(tmp, "cli", name), filepath.Join(tmp, "lib", name)
			for _, dst := range []string{cli, lib} {
				if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
					t.Fatal(err)
				}
				gitClone(t, filepath.Join(ws, name), dst)
			}
			// Shared repos resolve against the real workspace in both runs.
			mustRun(t, "-c", filepath.Join(cli, ".linespec.yml"), "migrate", "--workspace", ws)

			repo, err := migrate.Load(lib, migrate.LoadOptions{Workspace: ws})
			if err != nil {
				t.Fatal(err)
			}
			p, err := migrate.MakePlan(repo)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := migrate.Apply(p, migrate.Options{}, &record.FilesStore{Root: lib, Dir: repo.Dir}); err != nil {
				t.Fatal(err)
			}
			dir := repo.Dir
			diffTrees(t, filepath.Join(lib, dir), filepath.Join(cli, dir))
			if after := gitStatus(t, filepath.Join(ws, name), "provenance"); after != before {
				t.Fatalf("real repo %s changed:\n%s", name, after)
			}
		})
	}
}

// diffTrees fails unless the two directory trees hold identical files.
func diffTrees(t *testing.T, a, b string) {
	t.Helper()
	files := func(root string) map[string]string {
		out := map[string]string{}
		err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			rel, _ := filepath.Rel(root, p)
			out[rel] = readFile(t, p)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	fa, fb := files(a), files(b)
	for rel, body := range fa {
		if other, ok := fb[rel]; !ok {
			t.Errorf("only in %s: %s", a, rel)
		} else if other != body {
			t.Errorf("differs: %s", rel)
		}
	}
	for rel := range fb {
		if _, ok := fa[rel]; !ok {
			t.Errorf("only in %s: %s", b, rel)
		}
	}
}

func TestMigrateRewritesExcludePaths(t *testing.T) {
	for name, tc := range map[string]struct{ block, want string }{
		"flow": {
			"provenance:\n  dir: provenance\n  exclude_paths: ['.gitignore', '.linespec.yml', '.linespec/hash_manifest.json', 'docs/DESIGN.md']\n  shared_repos: []\n",
			"dir: provenance\nexclude_paths: ['.gitignore', '.provenance.yml', 'docs/DESIGN.md']\nshared_repos: []\n",
		},
		"block": {
			"provenance:\n  dir: provenance\n  exclude_paths:\n    - .linespec.yml\n    - .provenance.yml\n    - README.md\n\n  # after\n  commit_tag_required: true\n",
			"dir: provenance\nexclude_paths:\n    - .provenance.yml\n    - README.md\n\n# after\ncommit_tag_required: true\n",
		},
		"untouched": {
			"provenance:\n  dir: provenance\n  exclude_paths: [ 'README.md' ]\n",
			"dir: provenance\nexclude_paths: [ 'README.md' ]\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := migrateFixture(t)
			cfg := filepath.Join(root, ".linespec.yml")
			writeFile(t, cfg, tc.block)
			gitIn(t, root, "commit", "-qam", "config")
			mustRun(t, "-c", cfg, "migrate")
			if got := readFile(t, filepath.Join(root, ".provenance.yml")); got != tc.want {
				t.Fatalf(".provenance.yml:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}
