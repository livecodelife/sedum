package record

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/livecodelife/provenance/migrate"
	librecord "github.com/livecodelife/provenance/record"
)

// sedum reads records through the shared library (prov-2026-5f1e4f0a), so for
// every record in the real repositories - v1 as they are, and v2 after the
// library migrates a copy - sedum's gen_hash must equal the library's. Runs only
// when PROVENANCE_REAL_REPOS names the workspace holding the repositories.
var realRepos = []string{"linespec", "sedum", "scion", "linespec-provenance-ui"}

func TestGenHashMatchesTheLibraryOnRealRecords(t *testing.T) {
	ws := os.Getenv("PROVENANCE_REAL_REPOS")
	if ws == "" {
		t.Skip("set PROVENANCE_REAL_REPOS to the workspace directory to run")
	}
	for _, repo := range realRepos {
		t.Run(repo+"/v1", func(t *testing.T) {
			compareDir(t, filepath.Join(ws, repo, "provenance"))
		})
		t.Run(repo+"/v2", func(t *testing.T) {
			src := filepath.Join(ws, repo)
			loaded, err := migrate.Load(src, migrate.LoadOptions{Workspace: ws})
			if err != nil {
				t.Fatal(err)
			}
			plan, err := migrate.MakePlan(loaded)
			if err != nil {
				t.Fatal(err)
			}
			tmp := t.TempDir()
			if err := os.CopyFS(filepath.Join(tmp, loaded.Dir), os.DirFS(filepath.Join(src, loaded.Dir))); err != nil {
				t.Fatal(err)
			}
			store := &librecord.FilesStore{Root: tmp, Dir: loaded.Dir, SpecRoot: src}
			if _, err := migrate.Apply(plan, migrate.Options{}, store); err != nil {
				t.Fatal(err)
			}
			if n := compareDir(t, filepath.Join(tmp, loaded.Dir)); n == 0 {
				t.Fatal("no v2 records compared")
			}
		})
	}
}

// compareDir checks every record file in dir and returns how many it compared.
func compareDir(t *testing.T, dir string) int {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no records in %s (%v)", dir, err)
	}
	compared := 0
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		lib, err := librecord.Decode(data)
		if err != nil {
			t.Errorf("%s: the library cannot decode it: %v", f, err)
			continue
		}
		ours, err := readRecord(f, filepath.Base(f))
		if err != nil {
			t.Errorf("%s: sedum cannot read it: %v", f, err)
			continue
		}
		if want := lib.GenHash(); ours.Gen != want {
			t.Errorf("%s: sedum gen %s, library %s", f, ours.Gen, want)
		}
		if ours.ID != lib.ID || ours.Implements != lib.Implements || ours.Deprecated != (lib.Status == librecord.StatusDeprecated) ||
			len(ours.Aliases) != len(lib.Aliases) {
			t.Errorf("%s: sedum read id %q implements %q deprecated %v aliases %v; library %q %q %q %v",
				f, ours.ID, ours.Implements, ours.Deprecated, ours.Aliases, lib.ID, lib.Implements, lib.Status, lib.Aliases)
		}
		compared++
	}
	t.Logf("%d records compared in %s", compared, dir)
	return compared
}
