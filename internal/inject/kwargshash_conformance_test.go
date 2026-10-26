package inject

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The kwargs and body hash corpus (prov-2026-76fac245). A tool that does not
// link Sedum computes both hashes from the same rules; this test reads the
// published file and checks Sedum against it.

const kwargsHashCorpusPath = "../../conformance/kwargshash/cases.json"

var hashText = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

type kwargsHashCorpus struct {
	Format      string               `json:"format"`
	Version     int                  `json:"corpus_version"`
	Description string               `json:"description"`
	Cases       []kwargsHashCaseJSON `json:"cases"`
}

type kwargsHashCaseJSON struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Why        string `json:"why"`
	Input      string `json:"input"`
	Canonical  string `json:"canonical"`
	Hash       string `json:"hash"`
	Error      bool   `json:"error"`
	ErrorNames string `json:"error_names"`
}

func TestKwargsHashConformance(t *testing.T) {
	data, err := os.ReadFile(filepath.FromSlash(kwargsHashCorpusPath))
	if err != nil {
		t.Fatalf("reading the kwargs hash corpus: %v", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var corpus kwargsHashCorpus
	if err := decoder.Decode(&corpus); err != nil {
		t.Fatalf("the kwargs hash corpus is not readable: %v", err)
	}
	if corpus.Format != "sedum.kwargshash" {
		t.Fatalf("corpus format = %q, want sedum.kwargshash", corpus.Format)
	}
	if corpus.Version != 1 {
		t.Fatalf("corpus_version = %d, want 1", corpus.Version)
	}
	if corpus.Description == "" {
		t.Fatal("the corpus has no description")
	}
	if len(corpus.Cases) == 0 {
		t.Fatal("the kwargs hash corpus carries no cases")
	}

	seen := map[string]bool{}
	kinds := map[string]int{}
	for _, tc := range corpus.Cases {
		t.Run(tc.ID, func(t *testing.T) {
			if tc.ID == "" || seen[tc.ID] {
				t.Fatalf("case id %q is empty or repeated", tc.ID)
			}
			seen[tc.ID] = true
			kinds[tc.Kind]++
			if tc.Why == "" {
				t.Error("case says nothing about what goes wrong for an implementation that gets it wrong")
			}

			switch tc.Kind {
			case "kwargs":
				checkKwargsCase(t, tc)
			case "body":
				checkBodyCase(t, tc)
			default:
				t.Fatalf("unknown case kind %q", tc.Kind)
			}
		})
	}
	if kinds["kwargs"] == 0 || kinds["body"] == 0 {
		t.Errorf("the corpus has %d kwargs cases and %d body cases, want both kinds", kinds["kwargs"], kinds["body"])
	}
}

// decodeKwargsInput reads a case's JSON text the way a number must be read to
// be hashed: as the digits it was written with, never as a double.
func decodeKwargsInput(t *testing.T, text string) map[string]any {
	t.Helper()
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()
	var kwargs map[string]any
	if err := decoder.Decode(&kwargs); err != nil {
		t.Fatalf("input is not a JSON object: %v", err)
	}
	return kwargs
}

func checkKwargsCase(t *testing.T, tc kwargsHashCaseJSON) {
	t.Helper()
	kwargs := decodeKwargsInput(t, tc.Input)

	if tc.Error {
		if tc.ErrorNames == "" {
			t.Fatal("an error case names the text the failure must contain")
		}
		if tc.Canonical != "" || tc.Hash != "" {
			t.Error("an error case carries no canonical text and no hash")
		}
		got, err := CanonicalKwargs(kwargs)
		if err == nil {
			t.Fatalf("CanonicalKwargs accepted the input and returned %q", got)
		}
		if !strings.Contains(err.Error(), tc.ErrorNames) {
			t.Errorf("CanonicalKwargs error does not name %q: %v", tc.ErrorNames, err)
		}
		hash, err := KwargsHash(kwargs)
		if err == nil {
			t.Fatalf("KwargsHash accepted the input and returned %q", hash)
		}
		if !strings.Contains(err.Error(), tc.ErrorNames) {
			t.Errorf("KwargsHash error does not name %q: %v", tc.ErrorNames, err)
		}
		return
	}

	if !hashText.MatchString(tc.Hash) {
		t.Fatalf("case hash %q is not sha256: and 64 lower-case hex digits", tc.Hash)
	}
	got, err := CanonicalKwargs(kwargs)
	if err != nil {
		t.Fatalf("CanonicalKwargs: %v", err)
	}
	if string(got) != tc.Canonical {
		t.Errorf("CanonicalKwargs =\n  %s\nwant\n  %s", got, tc.Canonical)
	}
	hash, err := KwargsHash(kwargs)
	if err != nil {
		t.Fatalf("KwargsHash: %v", err)
	}
	if hash != tc.Hash {
		t.Errorf("KwargsHash = %s, want %s", hash, tc.Hash)
	}
	if want := sha256Of(tc.Canonical); tc.Hash != want {
		t.Errorf("the case's hash %s is not the SHA-256 of its canonical text (%s)", tc.Hash, want)
	}
}

func checkBodyCase(t *testing.T, tc kwargsHashCaseJSON) {
	t.Helper()
	if !hashText.MatchString(tc.Hash) {
		t.Fatalf("case hash %q is not sha256: and 64 lower-case hex digits", tc.Hash)
	}
	if got := BodyHash(tc.Input); got != tc.Hash {
		t.Errorf("BodyHash(%q) = %s, want %s", tc.Input, got, tc.Hash)
	}
}

// The corpus pins the relationships its cases are chosen to show, so a case
// edited to agree with a regression cannot hide it.
func TestKwargsHashCorpusShowsWhatItIsFor(t *testing.T) {
	data, err := os.ReadFile(filepath.FromSlash(kwargsHashCorpusPath))
	if err != nil {
		t.Fatalf("reading the kwargs hash corpus: %v", err)
	}
	var corpus kwargsHashCorpus
	if err := json.Unmarshal(data, &corpus); err != nil {
		t.Fatalf("the kwargs hash corpus is not valid JSON: %v", err)
	}
	byID := map[string]kwargsHashCaseJSON{}
	for _, tc := range corpus.Cases {
		byID[tc.ID] = tc
	}
	get := func(id string) kwargsHashCaseJSON {
		t.Helper()
		tc, ok := byID[id]
		if !ok {
			t.Fatalf("the corpus has no case %q", id)
		}
		return tc
	}

	same := [][]string{
		{"kwargs/integer-spelled-1", "kwargs/integer-spelled-1.0", "kwargs/integer-spelled-1e0"},
		{"body/lf-text", "body/crlf-text"},
		{"kwargs/absent-kwarg", "kwargs/empty-kwargs"},
	}
	for _, ids := range same {
		for _, id := range ids[1:] {
			if get(id).Hash != get(ids[0]).Hash {
				t.Errorf("%s and %s should hash alike", ids[0], id)
			}
		}
	}
	differ := [][]string{
		{"kwargs/composed-letter", "kwargs/decomposed-letter"},
		{"kwargs/absent-kwarg", "kwargs/explicit-null"},
		{"body/lf-text", "body/lone-carriage-return"},
		{"body/lf-text", "body/re-indented"},
		{"body/lf-text", "body/trailing-space"},
		{"body/lf-text", "body/extra-trailing-blank-line"},
	}
	for _, ids := range differ {
		if get(ids[0]).Hash == get(ids[1]).Hash {
			t.Errorf("%s and %s should hash differently", ids[0], ids[1])
		}
	}

	if got := get("kwargs/empty-kwargs").Hash; got != hashOfEmptyObject {
		t.Errorf("empty kwargs hash = %s, want %s", got, hashOfEmptyObject)
	}
	if got := get("body/empty").Hash; got != hashOfEmptyText {
		t.Errorf("empty body hash = %s, want %s", got, hashOfEmptyText)
	}
	if got := get("kwargs/negative-zero").Canonical; got != `{"n":0}` {
		t.Errorf("negative zero canonical text = %s", got)
	}
	if got := get("kwargs/small-fraction-uses-an-exponent").Canonical; got != `{"n":1.5e-7}` {
		t.Errorf("0.00000015 canonical text = %s", got)
	}
	for id, names := range map[string]string{
		"kwargs/integer-beyond-2^53-refused": "integer beyond 2^53",
		"kwargs/1e30-refused":                "integer beyond 2^53",
		"kwargs/1e400-refused-as-not-finite": "number is not finite",
	} {
		if tc := get(id); !tc.Error || tc.ErrorNames != names {
			t.Errorf("%s = error %v naming %q, want an error naming %q", id, tc.Error, tc.ErrorNames, names)
		}
	}
}
