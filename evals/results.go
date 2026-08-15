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

	// PromptTokens and CompletionTokens are the server's accounting of what
	// those calls cost. Omitted when the endpoint reported none, so an entry
	// says "not measured" rather than "cost nothing".
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
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
		}
		switch {
		case s.Err != nil:
			r.Outcome = "failed"
			r.Counts = nil
		case s.Invalid:
			r.Outcome = "invalid"
			r.Counts = nil
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
