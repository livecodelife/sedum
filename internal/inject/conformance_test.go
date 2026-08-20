package inject

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Sedum's own tests read the published conformance corpus rather than carrying
// a parallel copy of the same expectations (prov-2026-c6580f1e).
//
// The direction matters. A corpus maintained beside the tests it duplicates
// goes stale within a release and then actively misleads, because the stranger
// checking against it has no way to know. A corpus the tests depend on cannot:
// the published artifact is the thing that fails when Sedum regresses. So the
// cases that used to live here as Go literals live in conformance/markers now,
// and what remains is the driver.
//
// It is deliberately not a runner anyone else is meant to use. A runner in one
// language is a runner for one implementation, and the audience for the corpus
// is a writer who is not using Go at all.

const corpusPath = "../../conformance/markers/cases.json"

type corpus struct {
	Format  string       `json:"format"`
	Version int          `json:"corpus_version"`
	Cases   []corpusCase `json:"cases"`
}

type corpusCase struct {
	ID            string        `json:"id"`
	Kind          string        `json:"kind"`
	Decision      string        `json:"decision"`
	Why           string        `json:"why"`
	CommentPrefix string        `json:"comment_prefix"`
	Input         string        `json:"input"`
	Existing      string        `json:"existing"`
	Marker        *corpusMarker `json:"marker"`
	Parsed        *corpusParse  `json:"parsed"`
	Output        string        `json:"output"`
	Error         bool          `json:"error"`
	ErrorNames    string        `json:"error_names"`
}

type corpusMarker struct {
	Action  string                     `json:"action"`
	Variant string                     `json:"variant"`
	Tier    string                     `json:"tier"`
	Record  string                     `json:"record"`
	Writer  string                     `json:"writer"`
	Kwargs  map[string]any             `json:"kwargs"`
	Extra   map[string]json.RawMessage `json:"extra"`
}

type corpusParse struct {
	Recognized bool           `json:"recognized"`
	Action     string         `json:"action"`
	Variant    string         `json:"variant"`
	Tier       string         `json:"tier"`
	Record     string         `json:"record"`
	Writer     string         `json:"writer"`
	Kwargs     map[string]any `json:"kwargs"`
	Extra      map[string]any `json:"extra"`
}

func (m corpusMarker) marker() Marker {
	built := Marker{
		Action:  m.Action,
		Variant: m.Variant,
		Tier:    Tier(m.Tier),
		Record:  m.Record,
		Writer:  m.Writer,
	}
	if len(m.Kwargs) > 0 {
		built.Kwargs = m.Kwargs
	}
	if len(m.Extra) > 0 {
		built.Extra = m.Extra
	}
	return built
}

func loadCorpus(t *testing.T) corpus {
	t.Helper()
	data, err := os.ReadFile(filepath.FromSlash(corpusPath))
	if err != nil {
		t.Fatalf("reading the conformance corpus: %v", err)
	}
	var c corpus
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatalf("the conformance corpus is not valid JSON: %v", err)
	}
	if c.Format != "sedum.marker" {
		t.Fatalf("corpus format = %q, want sedum.marker", c.Format)
	}
	if len(c.Cases) == 0 {
		t.Fatal("the conformance corpus carries no cases")
	}
	return c
}

func TestMarkerConformanceCorpus(t *testing.T) {
	for _, tc := range loadCorpus(t).Cases {
		t.Run(tc.ID, func(t *testing.T) {
			// Every case says which decision it pins. A case that cites none
			// is coverage for its own sake, which this corpus does not accept.
			if tc.Decision == "" {
				t.Fatal("case cites no decision")
			}

			switch tc.Kind {
			case "parse":
				checkParse(t, tc)
			case "emit":
				checkEmit(t, tc, tc.Marker.marker())
			case "round_trip":
				checkParse(t, tc)
				checkRoundTrip(t, tc)
			case "replace":
				checkReplace(t, tc)
			case "close":
				if got := tc.Marker.marker().Close(tc.CommentPrefix); got != tc.Output {
					t.Errorf("Close =\n  %s\nwant\n  %s", got, tc.Output)
				}
			case "error":
				checkError(t, tc)
			default:
				t.Fatalf("unknown case kind %q", tc.Kind)
			}
		})
	}
}

func checkParse(t *testing.T, tc corpusCase) {
	t.Helper()
	got, ok, err := parseOpen(tc.CommentPrefix, tc.Input)
	if err != nil {
		t.Fatalf("parseOpen: %v", err)
	}
	if ok != tc.Parsed.Recognized {
		t.Fatalf("recognized = %v, want %v", ok, tc.Parsed.Recognized)
	}
	if !ok {
		return
	}

	if got.Action != tc.Parsed.Action {
		t.Errorf("action = %q, want %q", got.Action, tc.Parsed.Action)
	}
	if got.Variant != tc.Parsed.Variant {
		t.Errorf("variant = %q, want %q", got.Variant, tc.Parsed.Variant)
	}
	if string(got.Tier) != tc.Parsed.Tier {
		t.Errorf("tier = %q, want %q", got.Tier, tc.Parsed.Tier)
	}
	if got.Record != tc.Parsed.Record {
		t.Errorf("record = %q, want %q", got.Record, tc.Parsed.Record)
	}
	if got.Writer != tc.Parsed.Writer {
		t.Errorf("writer = %q, want %q", got.Writer, tc.Parsed.Writer)
	}
	if !sameValues(got.Kwargs, tc.Parsed.Kwargs) {
		t.Errorf("kwargs = %v, want %v", got.Kwargs, tc.Parsed.Kwargs)
	}

	// Extra is compared as JSON values rather than as text: the corpus says
	// what was retained, and the bytes it is retained as are the emit cases'
	// business.
	extra := map[string]any{}
	for name, raw := range got.Extra {
		var value any
		if err := json.Unmarshal(raw, &value); err != nil {
			t.Fatalf("carried key %q is not valid JSON: %v", name, err)
		}
		extra[name] = value
	}
	if !sameValues(extra, tc.Parsed.Extra) {
		t.Errorf("extra = %v, want %v", extra, tc.Parsed.Extra)
	}
}

func checkEmit(t *testing.T, tc corpusCase, m Marker) {
	t.Helper()
	got, err := m.Open(tc.CommentPrefix)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if got != tc.Output {
		t.Errorf("Open =\n  %s\nwant\n  %s", got, tc.Output)
	}
}

func checkRoundTrip(t *testing.T, tc corpusCase) {
	t.Helper()
	parsed, ok, err := parseOpen(tc.CommentPrefix, tc.Input)
	if err != nil || !ok {
		t.Fatalf("parseOpen = %v, %v", ok, err)
	}
	checkEmit(t, tc, parsed)

	// Re-emission is idempotent, or a rerun churns the file without changing
	// what the marker says.
	again, ok, err := parseOpen(tc.CommentPrefix, tc.Output)
	if err != nil || !ok {
		t.Fatalf("parseOpen on the expected output = %v, %v", ok, err)
	}
	stable, err := again.Open(tc.CommentPrefix)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if stable != tc.Output {
		t.Errorf("emitting twice is not stable:\n  %s\n  %s", tc.Output, stable)
	}
}

// checkReplace is the rule applyOne implements, stated as the corpus states it:
// every modelled field comes from the marker being written, and the unmodelled
// keys come from the marker already in the file.
func checkReplace(t *testing.T, tc corpusCase) {
	t.Helper()
	existing, ok, err := parseOpen(tc.CommentPrefix, tc.Existing)
	if err != nil || !ok {
		t.Fatalf("parseOpen on the existing marker = %v, %v", ok, err)
	}

	replacement := tc.Marker.marker()
	replacement.Extra = existing.Extra
	checkEmit(t, tc, replacement)
}

func checkError(t *testing.T, tc corpusCase) {
	t.Helper()
	_, _, err := parseOpen(tc.CommentPrefix, tc.Input)
	if err == nil {
		t.Fatal("a marker the corpus says is unreadable was accepted")
	}
	if !strings.Contains(err.Error(), tc.ErrorNames) {
		t.Errorf("error does not name %q: %v", tc.ErrorNames, err)
	}
}

func sameValues(got, want map[string]any) bool {
	if len(got) == 0 && len(want) == 0 {
		return true
	}
	return reflect.DeepEqual(got, want)
}
