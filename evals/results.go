package evals

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/calebcowen/sedum/internal/recording"
)

// An Entry is one run, recorded as it happened.
//
// A recording (internal/recording) is a decision and carries no volatile
// fields, so that two recordings of equivalent runs diff cleanly. This is the
// opposite and has to be: it says what happened rather than what to do, so it
// carries when, against what, and how - and it is never edited. A measurement
// that turns out to be wrong is superseded by a later entry rather than
// corrected in place (prov-2026-eb283c56).
type Entry struct {
	// Schema is the entry format's version, so a reader a year from now can
	// tell what it is looking at. Fields are added rather than repurposed, so
	// this moves only if something existing changes meaning.
	Schema int `json:"schema"`

	// Commit and Clean are what make an entry re-runnable. Because the fixtures
	// are vendored, the commit pins the generator package, the record and
	// Sedum's own code - everything except the model's sampling.
	//
	// It pins none of that if the tree was dirty, so Clean is carried and the
	// harness warns. A dirty run is still recorded: refusing to measure
	// mid-edit would rule out most of how this is used. It is recorded as
	// visibly not re-runnable, which is the fact that was missing both times a
	// number went stale.
	Commit string `json:"commit"`
	Clean  bool   `json:"clean"`
	// Dirty lists the paths that made it dirty, so the caveat is specific.
	Dirty []string `json:"dirty,omitempty"`

	At time.Time `json:"at"`

	Case      string `json:"case"`
	Model     Model  `json:"model"`
	Language  string `json:"language,omitempty"`
	Framework string `json:"framework,omitempty"`
	Tightness string `json:"tightness,omitempty"`
	Arm       string `json:"arm,omitempty"`

	// Endpoint is where the model was served from, as configured. Recorded
	// because a hosted row and a local row are not evidence about the same
	// question.
	Endpoint    string `json:"endpoint,omitempty"`
	Samples     int    `json:"samples"`
	Concurrency int    `json:"concurrency"`
	Retries     int    `json:"retries"`

	// Fixture is a digest of the generator packages, the records and the
	// behaviour target this entry was drawn against.
	//
	// It is what lets history tell an entry drawn against one fixture from an
	// entry drawn against another, and stop comparing across the difference.
	// Omitted on entries written before the field, which read as unknown and
	// are compared against nothing - there is no way to recover what those runs
	// used beyond their commit (prov-2026-43505e1b).
	Fixture string `json:"fixture,omitempty"`

	// Resolution is the question the sample size was drawn for: smoke, coarse
	// or fine.
	//
	// Omitted when there is none, and the entries already in results/ have none
	// - they were drawn at five because five was the default, and stamping them
	// "coarse" now would invent a decision nobody made. A reader sees "unstated"
	// for those, which is the true thing to say about them (prov-2026-3039750e).
	Resolution string `json:"resolution,omitempty"`

	Valid   int `json:"valid"`
	Invalid int `json:"invalid"`
	Failed  int `json:"failed"`

	WallMS int64 `json:"wall_ms"`

	// Runs keeps each sample rather than only their aggregate.
	//
	// This is what an aggregate cannot answer. The first run to establish
	// expectations reported observed means of 0.33, 0.67 and 3.33 invocations -
	// numbers no single answer can have produced - so it could say selection was
	// inconsistent and could not say whether any one sample was complete, which
	// is the question an expectation has to be set from.
	Runs []SampleRun `json:"runs"`

	// Details are the distinct reasons samples were invalid or failed, so a
	// stored rate can always say what its misses looked like.
	Details []string `json:"details,omitempty"`
}

// SampleRun is one sample's outcome as stored.
type SampleRun struct {
	Outcome string         `json:"outcome"` // valid | invalid | failed
	Counts  map[string]int `json:"counts,omitempty"`
	First   string         `json:"first,omitempty"`
	MS      int64          `json:"ms"`

	// Calls, Rejected and Completeness are what the sample cost.
	//
	// Added rather than replacing anything, so the entries written before them
	// stay readable and keep meaning what they meant - a reader ignores what it
	// does not recognise, and an older entry simply reports no call counts
	// rather than reporting zero of them. Cost in calls is what makes two
	// entries at different retry budgets comparable, and Rejected is what makes
	// first-call validity survive a raised budget (prov-2026-0811425c).
	Calls        int `json:"calls,omitempty"`
	Rejected     int `json:"rejected,omitempty"`
	Completeness int `json:"completeness,omitempty"`

	// Behavior is what applying this sample's selection produced, or nil when
	// behaviour was not measured. Added rather than replacing anything, on the
	// same rule as the counts above: a reader ignores what it does not
	// recognise, and every entry written before this one simply reports no
	// behaviour rather than reporting that nothing worked.
	//
	// Stored per sample rather than only as a rate, because the interesting
	// question later is which contract broke on which draw - and a rate cannot
	// be re-scored against a question sharpened after the run
	// (prov-2026-83340ba0).
	Behavior *BehaviorRun `json:"behavior,omitempty"`

	// PromptTokens and CompletionTokens are the server's accounting of what
	// those calls cost. Omitted when the endpoint reported none, so an entry
	// says "not measured" rather than "cost nothing".
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`

	// Rules are the rule slugs an invalid sample was rejected under, and
	// Invocations is what a valid one bound.
	//
	// Both are additive and both are stored rather than scored: nothing here
	// derives a rate from them, because a scoring rule is a decision and the
	// change that first makes data visible should not also decide what to
	// conclude from it (prov-2026-2256e6fa).
	//
	// An entry written before these carries neither, which reads as "not
	// recorded" exactly as the earlier additions do - the schema is additive
	// and a reader ignores what it does not recognise (prov-2026-eb283c56).
	Rules []string `json:"rules,omitempty"`
	// Violations is the rendered text of each rule Rules names, in the same
	// order. A slug can be tallied and a sentence cannot, which is why both are
	// kept: the tally says how often, the sentence says what the answer was
	// (prov-2026-986ac4ca).
	Violations  []string               `json:"violations,omitempty"`
	Invocations []recording.Invocation `json:"invocations,omitempty"`

	// Fill, Syntax and Idempotent are the three derived signals for this
	// sample, stored per sample for the same reason Behavior is: a rate cannot
	// be re-scored against a question sharpened after the run, and the
	// interesting question later is which file failed on which draw.
	//
	// Pointers, so that an entry written before these reads as "not recorded"
	// rather than as a sample that filled no anchors and parsed nothing. The
	// schema is additive and a reader ignores what it does not recognise
	// (prov-2026-eb283c56, prov-2026-d61010a4).
	Fill       *AnchorFill  `json:"fill,omitempty"`
	Syntax     *SyntaxCheck `json:"syntax,omitempty"`
	Idempotent *Idempotency `json:"idempotent,omitempty"`

	// Files is what a baseline sample wrote, keyed by path, and is absent on
	// the sedum arm - where Invocations above is the answer and the files are
	// derivable from it and the package.
	//
	// The contents, not the paths. Storing a rate without the source that
	// produced it makes an entry single-use: the first baseline run failed the
	// same fifteen assertions in all five samples, and answering *why* meant
	// re-deriving it from a separate hand-run probe because the samples
	// themselves were gone. That is exactly the re-scoring prov-2026-2256e6fa
	// stores invocations to allow, and the baseline arm's equivalent is this
	// (prov-2026-a4dbe65c).
	Files map[string]string `json:"files,omitempty"`

	// Missing is the authorized paths a baseline answer left out, and
	// Unexpected the paths it wrote that nothing authorized - the second never
	// reached disk. Both are empty on the sedum arm.
	Missing    []string `json:"missing,omitempty"`
	Unexpected []string `json:"unexpected,omitempty"`
}

// EntrySchema is the current entry format. Additive changes do not move it.
const EntrySchema = 1

// NewEntry builds an entry from a finished measurement.
//
// The retry budget is read off the measurement rather than passed alongside it,
// because the two disagreeing would put a number in the file that describes a
// run nobody made (prov-2026-b4555efc).
func NewEntry(m Measurement, endpoint string) Entry {
	commit, clean, dirty := gitState()
	t := m.Tally()

	e := Entry{
		Schema:      EntrySchema,
		Commit:      commit,
		Clean:       clean,
		Dirty:       dirty,
		At:          time.Now().UTC().Truncate(time.Second),
		Case:        m.Case.ID,
		Model:       m.Model,
		Language:    m.Case.Language,
		Framework:   m.Case.Framework,
		Tightness:   m.Case.Tightness,
		Arm:         m.Case.Arm,
		Endpoint:    endpoint,
		Samples:     len(m.Samples),
		Concurrency: m.Concurrency,
		Retries:     m.Retries,
		Fixture:     fixtureOf(m.Case),
		Resolution:  string(m.Resolution),
		Valid:       t.Valid,
		Invalid:     t.Invalid,
		Failed:      t.Failed,
		WallMS:      m.Wall.Milliseconds(),
		Details:     m.Details(),
	}

	for _, s := range m.Samples {
		r := SampleRun{
			Outcome:          "valid",
			Counts:           s.Counts,
			First:            s.First,
			MS:               s.Elapsed.Milliseconds(),
			Calls:            s.Calls,
			Rejected:         s.Rejected,
			Completeness:     s.Completeness,
			PromptTokens:     s.PromptTokens,
			CompletionTokens: s.CompletionTokens,
			Rules:            s.Rules,
			Violations:       s.Violations,
			Invocations:      s.Invocations,
			Behavior:         s.Behavior,
		}
		switch {
		case s.Err != nil:
			r.Outcome = "failed"
			r.Counts = nil
		case s.Invalid:
			r.Outcome = "invalid"
			r.Counts = nil
		default:
			// Only a sample that reached Phase 7 has these to record. A
			// rejected answer never wrote a file, and storing its zeroes
			// would put them in a denominator they do not belong in.
			syntax := s.Syntax
			r.Syntax = &syntax

			// Anchor fill and idempotency need a package and an invocation
			// list. Storing them zero-valued for a baseline would put a
			// measurement nobody made on disk, where the report's care about
			// not printing them cannot help a later reader
			// (prov-2026-a4dbe65c).
			if m.Case.Arm == "baseline" {
				r.Files = s.Files
				r.Missing = s.Missing
				r.Unexpected = s.Unexpected
				break
			}
			fill, idem := s.Fill, s.Idempotent
			r.Fill, r.Idempotent = &fill, &idem
		}
		e.Runs = append(e.Runs, r)
	}
	return e
}

// Append writes an entry to the case's results file, creating it if needed.
//
// Append-only, one entry per line: a run adds lines and rewrites nothing, so
// two runs never conflict and the history cannot be edited by accident.
func Append(dir string, e Entry) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, e.Case+".jsonl")

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Compact rather than indented: one entry is one line, which is what makes
	// the file append-only in the way that matters.
	raw, err := json.Marshal(e)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(raw, '\n')); err != nil {
		return err
	}
	return nil
}

// ReadEntries loads a case's history, oldest first.
func ReadEntries(dir, caseID string) ([]Entry, error) {
	f, err := os.Open(filepath.Join(dir, caseID+".jsonl"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var out []Entry
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("%s: %w", caseID, err)
		}
		out = append(out, e)
	}
	return out, sc.Err()
}

// Complete reports the samples whose counts match every expectation exactly.
//
// This is what the aggregate could not answer, and it is the whole reason runs
// are kept per sample: an expectation is set from an answer that was actually
// complete, and a mean cannot say whether any one answer was.
func (e Entry) Complete(expected map[string]int) []int {
	if len(expected) == 0 {
		return nil
	}
	var out []int
	for i, r := range e.Runs {
		if r.Outcome != "valid" {
			continue
		}
		ok := true
		for name, want := range expected {
			if r.Counts[name] != want {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, i)
		}
	}
	return out
}

// Vocabulary returns every action any valid sample selected, so a case with no
// expectations can be read for what a complete answer might contain.
func (e Entry) Vocabulary() []string {
	seen := map[string]bool{}
	for _, r := range e.Runs {
		for name := range r.Counts {
			seen[name] = true
		}
	}
	out := make([]string, 0, len(seen))
	for name := range seen {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// gitState reports the commit an entry ran against and whether the tree was
// clean. A failure to read either is not fatal - a result measured outside a
// checkout is still a result, and says so by carrying no commit.
func gitState() (commit string, clean bool, dirty []string) {
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return "", false, nil
	}
	commit = strings.TrimSpace(string(out))

	status, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return commit, false, nil
	}
	for _, line := range strings.Split(strings.TrimSpace(string(status)), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			dirty = append(dirty, line)
		}
	}
	return commit, len(dirty) == 0, dirty
}

// fixtureOf digests the case, or returns empty when it cannot.
//
// A digest that could not be computed is recorded as absent rather than as a
// failed run. It is a reading aid for history, and a measurement that was
// otherwise fine should not be lost because a fixture directory moved between
// the run and the write.
func fixtureOf(c Case) string {
	digest, err := FixtureDigest(c)
	if err != nil {
		return ""
	}
	return digest
}
