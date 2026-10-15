package genhash

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// The shared hashing corpus lives in the provenance repository
// (conformance/hashing/cases.json), which owns gen_hash. This test runs its
// gen_hash cases against the interim implementation while Sedum cannot yet
// import the library; once Of delegates to record.GenHash, the library's own
// tests cover this and the test goes with the interim code.
//
// PROVENANCE_HASHING_CORPUS names the file; by default it is looked for in a
// provenance checkout beside this repository, and the test skips when there
// is none rather than carrying a copy that would go stale.
func corpusPath() string {
	if p := os.Getenv("PROVENANCE_HASHING_CORPUS"); p != "" {
		return p
	}
	return filepath.Join("..", "..", "..", "provenance", "conformance", "hashing", "cases.json")
}

func TestOfPassesTheSharedHashingCorpus(t *testing.T) {
	data, err := os.ReadFile(corpusPath())
	if os.IsNotExist(err) {
		t.Skipf("no shared hashing corpus at %s", corpusPath())
	}
	if err != nil {
		t.Fatal(err)
	}

	var corpus struct {
		Cases []struct {
			ID     string `json:"id"`
			Kind   string `json:"kind"`
			Record struct {
				Intent         string   `json:"intent"`
				Constraints    []string `json:"constraints"`
				AffectedScope  []string `json:"affected_scope"`
				ForbiddenScope []string `json:"forbidden_scope"`
			} `json:"record"`
			Canonical string `json:"canonical"`
			Hash      string `json:"hash"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}

	ran := 0
	for _, c := range corpus.Cases {
		if c.Kind != "gen_hash" {
			continue
		}
		ran++
		t.Run(c.ID, func(t *testing.T) {
			in := Input(c.Record)
			if got := string(canonical(in)); got != c.Canonical {
				t.Errorf("canonical\n got: %s\nwant: %s", got, c.Canonical)
			}
			if got := Of(in); got != c.Hash {
				t.Errorf("Of() = %s, want %s", got, c.Hash)
			}
		})
	}
	if ran == 0 {
		t.Fatal("the shared corpus has no gen_hash cases")
	}
}
