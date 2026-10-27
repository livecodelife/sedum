// Package evidence keeps sedum's own code-freshness evidence for provenance
// records: the code hash of a record's regions, the per-record evidence files
// under .sedum/evidence, and the verdict a record's latest entry gives.
package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
)

// Region is one sedum region as the code hash sees it. Tier is empty when the
// marker has none; "owned" marks a region sedum regenerates and any other value,
// such as "seeded", one it does not.
type Region struct {
	Gen  string
	Pkg  string
	Tier string
}

// ErrUnstamped says that a region has no gen or pkg hash, so its record's code
// cannot be hashed; `sedum grow --stamp` supplies them.
var ErrUnstamped = errors.New("a region has no gen or pkg hash; run sedum grow --stamp to supply them")

// Owned keeps the regions whose Tier is empty or "owned", in their order, and
// returns an empty slice when none are kept.
func Owned(regions []Region) []Region {
	kept := []Region{}
	for _, r := range regions {
		if r.Tier == "" || r.Tier == "owned" {
			kept = append(kept, r)
		}
	}
	return kept
}

// CodeHash returns the code hash of the owned regions (docs/HASHING.md section
// 7): each region as the canonical JSON object {"gen":..,"pkg":..}, the
// encodings sorted by bytes with duplicates kept, joined with commas in square
// brackets, and hashed with SHA-256. An owned region without gen or pkg gives
// ErrUnstamped. It never changes its argument.
func CodeHash(regions []Region) (string, error) {
	owned := Owned(regions)
	encoded := make([]string, 0, len(owned))
	for _, r := range owned {
		if r.Gen == "" || r.Pkg == "" {
			return "", ErrUnstamped
		}
		encoded = append(encoded, canonicalRegion(r))
	}
	sort.Strings(encoded)
	sum := sha256.Sum256([]byte("[" + strings.Join(encoded, ",") + "]"))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func canonicalRegion(r Region) string {
	return `{"gen":` + canonicalString(r.Gen) + `,"pkg":` + canonicalString(r.Pkg) + `}`
}

// canonicalString writes a JSON string escaping only the double quote, the
// backslash and the code points U+0000 to U+001F; everything else is as it is.
func canonicalString(s string) string {
	const hexDigits = "0123456789abcdef"
	var b strings.Builder
	b.WriteByte('"')
	for _, c := range s {
		switch {
		case c == '"':
			b.WriteString(`\"`)
		case c == '\\':
			b.WriteString(`\\`)
		case c < 0x20:
			b.WriteString(`\u00`)
			b.WriteByte(hexDigits[c>>4])
			b.WriteByte(hexDigits[c&0xf])
		default:
			b.WriteRune(c)
		}
	}
	b.WriteByte('"')
	return b.String()
}
