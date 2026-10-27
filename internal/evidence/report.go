package evidence

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Verdict is the freshness of a record's code.
type Verdict string

// The verdicts, in the order Judge decides them: no regions, then unknown, then
// none, then fresh or stale.
const (
	Fresh     Verdict = "fresh"
	Stale     Verdict = "stale"
	None      Verdict = "none"
	Unknown   Verdict = "unknown"
	NoRegions Verdict = "no regions"
)

// Result is one record's verdict. CodeHash is the code hash of the owned regions
// now (empty when unknown or no regions), Recorded the hash of the latest entry
// for the version (empty unless fresh or stale) and Regions the number of owned
// regions.
type Result struct {
	Record   string  `json:"record"`
	Version  int     `json:"version"`
	Verdict  Verdict `json:"verdict"`
	CodeHash string  `json:"code_hash"`
	Recorded string  `json:"recorded"`
	Regions  int     `json:"regions"`
}

// Judge gives a record's verdict from its regions and its evidence entries: no
// regions when none is owned, unknown when an owned region has no gen or pkg,
// none when no entry covers the version, otherwise fresh when the latest entry
// for the version has the current hash and stale when it does not.
func Judge(record string, version int, regions []Region, entries []Entry) Result {
	res := Result{Record: record, Version: version, Regions: len(Owned(regions))}
	if res.Regions == 0 {
		res.Verdict = NoRegions
		return res
	}
	hash, err := CodeHash(regions)
	if err == ErrUnstamped {
		res.Verdict = Unknown
		return res
	}
	res.CodeHash = hash
	latest, ok := Latest(entries, version)
	if !ok {
		res.Verdict = None
		return res
	}
	res.Recorded = latest.CodeHash
	if latest.CodeHash == hash {
		res.Verdict = Fresh
	} else {
		res.Verdict = Stale
	}
	return res
}

// Failing reports whether any result is stale or unknown.
func Failing(results []Result) bool {
	for _, r := range results {
		if r.Verdict == Stale || r.Verdict == Unknown {
			return true
		}
	}
	return false
}

// Format is the name and version of the JSON document.
const Format = "sedum.evidence/1"

// Document is the JSON report.
type Document struct {
	Format  string   `json:"format"`
	Records []Result `json:"records"`
}

func sortedCopy(results []Result) []Result {
	out := append([]Result{}, results...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Record < out[j].Record })
	return out
}

func short(hash string) string {
	hex := strings.TrimPrefix(hash, "sha256:")
	if hex == "" {
		return "-"
	}
	if len(hex) > 12 {
		hex = hex[:12]
	}
	return hex
}

// WriteText writes one line per result in record id order, then an empty line
// and the count of each verdict.
func WriteText(w io.Writer, results []Result) error {
	counts := map[Verdict]int{}
	for _, r := range sortedCopy(results) {
		counts[r.Verdict]++
		if _, err := fmt.Fprintf(w, "%s  %s  now=%s recorded=%s\n", r.Record, r.Verdict, short(r.CodeHash), short(r.Recorded)); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "\n%d %s, %d %s, %d %s, %d %s, %d %s\n",
		counts[Fresh], Fresh, counts[Stale], Stale, counts[None], None, counts[Unknown], Unknown, counts[NoRegions], NoRegions)
	return err
}

// WriteJSON writes the Document of the results in record id order, two-space
// indented, without HTML escaping, with a final newline.
func WriteJSON(w io.Writer, results []Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(Document{Format: Format, Records: sortedCopy(results)})
}
