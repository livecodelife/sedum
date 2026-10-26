package inject

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/genpkg"
)

// The v2 marker corpus (prov-2026-54b6adc7). Like the v1 corpus, it is read by
// Sedum's own tests so that the published file is the thing that fails when
// Sedum regresses. Its cases carry no decision field.

const v2CorpusPath = "../../conformance/markers/v2cases.json"

type v2Corpus struct {
	Format      string       `json:"format"`
	Version     int          `json:"corpus_version"`
	Description string       `json:"description"`
	Cases       []v2CaseJSON `json:"cases"`
}

type v2CaseJSON struct {
	ID            string    `json:"id"`
	Kind          string    `json:"kind"`
	Why           string    `json:"why"`
	CommentPrefix string    `json:"comment_prefix"`
	CommentSuffix string    `json:"comment_suffix"`
	Input         string    `json:"input"`
	Existing      string    `json:"existing"`
	Marker        *v2Marker `json:"marker"`
	Parsed        *v2Parsed `json:"parsed"`
	Output        string    `json:"output"`
	Error         bool      `json:"error"`
	ErrorNames    string    `json:"error_names"`
}

func (c v2CaseJSON) comment() genpkg.Comment {
	return genpkg.Comment{Prefix: c.CommentPrefix, Suffix: c.CommentSuffix}
}

// v2Marker is the marker object of a case: the fields of the v1 corpus and
// version, kwargs_hash, body and foreign. For a v2 marker kwargs is the id.
type v2Marker struct {
	Action     string                     `json:"action"`
	Variant    string                     `json:"variant"`
	Tier       string                     `json:"tier"`
	Record     string                     `json:"record"`
	Gen        string                     `json:"gen"`
	Pkg        string                     `json:"pkg"`
	Writer     string                     `json:"writer"`
	Kwargs     map[string]any             `json:"kwargs"`
	Extra      map[string]json.RawMessage `json:"extra"`
	Version    int                        `json:"version"`
	KwargsHash string                     `json:"kwargs_hash"`
	Body       string                     `json:"body"`
	Foreign    bool                       `json:"foreign"`
}

type v2Parsed struct {
	Recognized bool           `json:"recognized"`
	Action     string         `json:"action"`
	Variant    string         `json:"variant"`
	Tier       string         `json:"tier"`
	Record     string         `json:"record"`
	Gen        string         `json:"gen"`
	Pkg        string         `json:"pkg"`
	Writer     string         `json:"writer"`
	Kwargs     map[string]any `json:"kwargs"`
	Extra      map[string]any `json:"extra"`
	Version    int            `json:"version"`
	KwargsHash string         `json:"kwargs_hash"`
	Body       string         `json:"body"`
	Foreign    bool           `json:"foreign"`
}

func (m v2Marker) marker() Marker {
	built := Marker{
		Action:     m.Action,
		Variant:    m.Variant,
		Tier:       Tier(m.Tier),
		Record:     m.Record,
		Gen:        m.Gen,
		Pkg:        m.Pkg,
		Writer:     m.Writer,
		Version:    m.Version,
		KwargsHash: m.KwargsHash,
		Body:       m.Body,
		Foreign:    m.Foreign,
	}
	if len(m.Kwargs) > 0 {
		built.Kwargs = m.Kwargs
	}
	if len(m.Extra) > 0 {
		built.Extra = m.Extra
	}
	return built
}

func loadV2Corpus(t *testing.T) v2Corpus {
	t.Helper()
	data, err := os.ReadFile(filepath.FromSlash(v2CorpusPath))
	if err != nil {
		t.Fatalf("reading the v2 conformance corpus: %v", err)
	}

	// Strict decoding: a field the runner does not read is a field the corpus
	// carries for no one, and a decision field in particular is not allowed.
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var c v2Corpus
	if err := decoder.Decode(&c); err != nil {
		t.Fatalf("the v2 conformance corpus is not readable: %v", err)
	}
	if c.Format != "sedum.marker.v2" {
		t.Fatalf("corpus format = %q, want sedum.marker.v2", c.Format)
	}
	if c.Version != 1 {
		t.Fatalf("corpus_version = %d, want 1", c.Version)
	}
	if c.Description == "" {
		t.Fatal("the corpus has no description")
	}
	if len(c.Cases) == 0 {
		t.Fatal("the v2 conformance corpus carries no cases")
	}

	var raw struct {
		Cases []map[string]json.RawMessage `json:"cases"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("the v2 conformance corpus is not valid JSON: %v", err)
	}
	for _, tc := range raw.Cases {
		if _, ok := tc["decision"]; ok {
			t.Errorf("case %s carries a decision field; the v2 corpus has none", tc["id"])
		}
	}
	return c
}

func TestMarkerV2Conformance(t *testing.T) {
	corpus := loadV2Corpus(t)
	seen := map[string]bool{}
	kinds := map[string]bool{}
	for _, tc := range corpus.Cases {
		t.Run(tc.ID, func(t *testing.T) {
			if tc.ID == "" || seen[tc.ID] {
				t.Fatalf("case id %q is empty or repeated", tc.ID)
			}
			seen[tc.ID] = true
			kinds[tc.Kind] = true
			if tc.Why == "" {
				t.Error("case says nothing about what goes wrong for an implementation that gets it wrong")
			}
			if tc.CommentPrefix == "" {
				t.Error("case declares no comment_prefix")
			}

			switch tc.Kind {
			case "parse":
				checkV2Parse(t, tc)
			case "emit":
				checkV2Emit(t, tc, tc.Marker.marker())
			case "round_trip":
				checkV2Parse(t, tc)
				checkV2RoundTrip(t, tc)
			case "replace":
				checkV2Replace(t, tc)
			case "error":
				checkV2Error(t, tc)
			default:
				t.Fatalf("unknown case kind %q", tc.Kind)
			}
		})
	}
	for _, kind := range []string{"parse", "emit", "round_trip", "replace", "error"} {
		if !kinds[kind] {
			t.Errorf("the v2 corpus has no %s case", kind)
		}
	}
}

func checkV2Parse(t *testing.T, tc v2CaseJSON) {
	t.Helper()
	if tc.Parsed == nil {
		t.Fatal("case has no parsed object")
	}
	got, ok, err := parseOpen(tc.comment(), tc.Input)
	if err != nil {
		t.Fatalf("parseOpen: %v", err)
	}
	if ok != tc.Parsed.Recognized {
		t.Fatalf("recognized = %v, want %v", ok, tc.Parsed.Recognized)
	}
	if !ok {
		return
	}

	want := tc.Parsed
	checks := []struct {
		name      string
		got, want string
	}{
		{"action", got.Action, want.Action},
		{"variant", got.Variant, want.Variant},
		{"tier", string(got.Tier), want.Tier},
		{"record", got.Record, want.Record},
		{"gen", got.Gen, want.Gen},
		{"pkg", got.Pkg, want.Pkg},
		{"writer", got.Writer, want.Writer},
		{"kwargs_hash", got.KwargsHash, want.KwargsHash},
		{"body", got.Body, want.Body},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
	if got.Version != want.Version {
		t.Errorf("version = %d, want %d", got.Version, want.Version)
	}
	if got.Foreign != want.Foreign {
		t.Errorf("foreign = %v, want %v", got.Foreign, want.Foreign)
	}
	if !sameValues(got.Kwargs, want.Kwargs) {
		t.Errorf("kwargs = %v, want %v", got.Kwargs, want.Kwargs)
	}

	// Extra is compared as JSON values, not as text.
	extra := map[string]any{}
	for name, raw := range got.Extra {
		var value any
		if err := json.Unmarshal(raw, &value); err != nil {
			t.Fatalf("carried key %q is not valid JSON: %v", name, err)
		}
		extra[name] = value
	}
	if !sameValues(extra, want.Extra) {
		t.Errorf("extra = %v, want %v", extra, want.Extra)
	}
}

func checkV2Emit(t *testing.T, tc v2CaseJSON, m Marker) {
	t.Helper()
	got, err := m.Open(tc.comment())
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if got != tc.Output {
		t.Errorf("Open =\n  %s\nwant\n  %s", got, tc.Output)
	}
}

func checkV2RoundTrip(t *testing.T, tc v2CaseJSON) {
	t.Helper()
	parsed, ok, err := parseOpen(tc.comment(), tc.Input)
	if err != nil || !ok {
		t.Fatalf("parseOpen = %v, %v", ok, err)
	}
	checkV2Emit(t, tc, parsed)

	again, ok, err := parseOpen(tc.comment(), tc.Output)
	if err != nil || !ok {
		t.Fatalf("parseOpen on the expected output = %v, %v", ok, err)
	}
	stable, err := again.Open(tc.comment())
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if stable != tc.Output {
		t.Errorf("emitting twice is not stable:\n  %s\n  %s", tc.Output, stable)
	}
}

// checkV2Replace is the rule applyOne implements: every modelled field comes
// from the marker being written, and the unmodelled keys from the one being
// replaced.
func checkV2Replace(t *testing.T, tc v2CaseJSON) {
	t.Helper()
	existing, ok, err := parseOpen(tc.comment(), tc.Existing)
	if err != nil || !ok {
		t.Fatalf("parseOpen on the existing marker = %v, %v", ok, err)
	}

	replacement := tc.Marker.marker()
	replacement.Extra = existing.Extra
	checkV2Emit(t, tc, replacement)

	if replacement.Version == 2 && strings.Contains(tc.Output, `"kwargs":{`) {
		t.Error("the replacing v2 marker still carries a v1 kwargs object")
	}
}

func checkV2Error(t *testing.T, tc v2CaseJSON) {
	t.Helper()
	if !tc.Error || tc.ErrorNames == "" {
		t.Fatal("an error case sets error and error_names")
	}
	_, _, err := parseOpen(tc.comment(), tc.Input)
	if err == nil {
		t.Fatal("a marker the corpus says is unreadable was accepted")
	}
	if !strings.Contains(err.Error(), tc.ErrorNames) {
		t.Errorf("error does not name %q: %v", tc.ErrorNames, err)
	}
}
