package genpkg

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// pkgCorpus is conformance/pkghash/cases.json. The test reads the published
// corpus rather than a parallel copy of it, for the reason the marker corpus
// does: what a stranger checks against has to be what fails when Sedum
// regresses.
type pkgCorpus struct {
	Cases []struct {
		ID        string            `json:"id"`
		Files     map[string]string `json:"files"`
		Action    string            `json:"action"`
		Canonical string            `json:"canonical"`
		Pkg       string            `json:"pkg"`
	} `json:"cases"`
}

func readPkgCorpus(t *testing.T) pkgCorpus {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "conformance", "pkghash", "cases.json"))
	if err != nil {
		t.Fatal(err)
	}
	var corpus pkgCorpus
	if err := json.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}
	if len(corpus.Cases) == 0 {
		t.Fatal("the pkg corpus has no cases")
	}
	return corpus
}

func TestPkgHashConformance(t *testing.T) {
	for _, c := range readPkgCorpus(t).Cases {
		t.Run(c.ID, func(t *testing.T) {
			set, findings := loadTree(t, c.Files)
			if findings.HasErrors() {
				t.Fatalf("corpus package does not load: %v", findings)
			}
			pkg, ok := set.Lookup("rails")
			if !ok {
				t.Fatal("corpus package rails did not load")
			}
			action, ok := pkg.Actions[c.Action]
			if !ok {
				t.Fatalf("corpus action %s did not load", c.Action)
			}

			canonical, err := pkgCanonical(pkg, action)
			if err != nil {
				t.Fatal(err)
			}
			if string(canonical) != c.Canonical {
				t.Errorf("canonical encoding differs\n got: %s\nwant: %s", canonical, c.Canonical)
			}
			if action.PkgHash != c.Pkg {
				t.Errorf("PkgHash = %s, want %s", action.PkgHash, c.Pkg)
			}
		})
	}
}

// pkgHashes loads a package tree and returns every action's pkg hash.
func pkgHashes(t *testing.T, files map[string]string) map[string]string {
	t.Helper()
	set, findings := loadTree(t, files)
	if findings.HasErrors() {
		t.Fatalf("package does not load: %v", findings)
	}
	pkg, ok := set.Lookup("rails")
	if !ok {
		t.Fatal("package rails did not load")
	}
	out := map[string]string{}
	for name, action := range pkg.Actions {
		if !strings.HasPrefix(action.PkgHash, "sha256:") || len(action.PkgHash) != len("sha256:")+64 {
			t.Fatalf("action %s: PkgHash %q is not sha256:<64 hex>", name, action.PkgHash)
		}
		out[name] = action.PkgHash
	}
	return out
}

func TestPkgHashIsStableAcrossLoads(t *testing.T) {
	first := pkgHashes(t, validPackage())
	for i := 0; i < 5; i++ {
		again := pkgHashes(t, validPackage())
		for name, h := range first {
			if again[name] != h {
				t.Fatalf("action %s hashed %s, then %s", name, h, again[name])
			}
		}
	}
}

// What a template renders for fixed kwargs changes pkg; what only the model
// reads does not.
func TestPkgHashCoversWhatRendersAndNothingElse(t *testing.T) {
	base := pkgHashes(t, validPackage())

	cases := []struct {
		name    string
		edits   map[string]*string
		changes []string // actions whose hash must change; every other must not
	}{
		{
			name: "an action's own template",
			edits: map[string]*string{
				"rails/actions/addBeforeFilter.rb": text("before_action :{{filter}}, prepend: true\n"),
			},
			changes: []string{"addBeforeFilter"},
		},
		{
			name: "one variant's template",
			edits: map[string]*string{
				"rails/actions/createControllerMethod/index.rb": text("def index\nend\n"),
			},
			changes: []string{"createControllerMethod"},
		},
		{
			name: "a package transform reaches every action",
			edits: map[string]*string{
				"rails/sedum.yaml": text(strings.Replace(validPackage()["rails/sedum.yaml"],
					`instantize: [plural, "prefix:@"]`, `instantize: [plural, snake, "prefix:@"]`, 1)),
			},
			changes: allActions(base),
		},
		{
			name: "a kwarg description reaches only the model",
			edits: map[string]*string{
				"rails/actions/actions.yaml": text(strings.Replace(validPackage()["rails/actions/actions.yaml"],
					"filter: { type: string, required: true }",
					"filter: { type: string, required: true, description: The filter method. }", 1)),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := pkgHashes(t, mutated(c.edits))
			for name, h := range base {
				changed := got[name] != h
				want := contains(c.changes, name)
				if changed != want {
					t.Errorf("action %s: changed=%v, want changed=%v", name, changed, want)
				}
			}
		})
	}
}

func allActions(hashes map[string]string) []string {
	var out []string
	for name := range hashes {
		out = append(out, name)
	}
	return out
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
