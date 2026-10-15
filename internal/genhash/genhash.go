// Package genhash computes a provenance record's gen_hash: the hash a marker
// stores as "gen" so a later grow can tell whether the record that wrote a
// region has changed since (prov-2026-66bea9f3).
//
// gen_hash covers exactly the fields Sedum reads from a record - intent,
// constraints, affected_scope and forbidden_scope - and nothing else, so a
// status flip, a seal or a re-pin never looks like a content change.
//
// The definition is not Sedum's. It belongs to the shared provenance record
// library (github.com/livecodelife/provenance/record), whose conformance corpus
// keeps every implementation producing the same bytes. Until Sedum can import
// that module, Of implements the same canonical encoding here; when it can, Of
// delegates to record.GenHash and the encoding below is deleted. Input mirrors
// the library's GenInput field for field so that swap is a conversion, and
// nothing outside this package changes.
package genhash

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// Input is the part of a record gen_hash covers. A nil list and an empty one
// are the same record.
type Input struct {
	Intent         string
	Constraints    []string
	AffectedScope  []string
	ForbiddenScope []string
}

// Of returns the record's gen_hash as "sha256:" followed by 64 lowercase hex
// digits.
func Of(in Input) string {
	sum := sha256.Sum256(canonical(in))
	return "sha256:" + hex.EncodeToString(sum[:])
}

// canonical encodes the four fields as one JSON object with byte-sorted keys
// and no insignificant whitespace. Every string is NFC-normalized, trimmed, and
// has each run of Unicode White_Space collapsed to one space. List order is
// kept and nothing is deduplicated or dropped.
//
// The keys are written in sorted order by hand rather than by marshalling a
// map, because encoding/json escapes <, >, & and U+2028/U+2029, and the
// canonical form escapes only what JSON requires.
func canonical(in Input) []byte {
	var b strings.Builder
	b.WriteString(`{"affected_scope":`)
	writeList(&b, in.AffectedScope)
	b.WriteString(`,"constraints":`)
	writeList(&b, in.Constraints)
	b.WriteString(`,"forbidden_scope":`)
	writeList(&b, in.ForbiddenScope)
	b.WriteString(`,"intent":`)
	writeString(&b, normalize(in.Intent))
	b.WriteString(`}`)
	return []byte(b.String())
}

func writeList(b *strings.Builder, items []string) {
	b.WriteByte('[')
	for i, item := range items {
		if i > 0 {
			b.WriteByte(',')
		}
		writeString(b, normalize(item))
	}
	b.WriteByte(']')
}

func normalize(s string) string {
	s = norm.NFC.String(s)
	var b strings.Builder
	pendingSpace := false
	for _, r := range s {
		if unicode.Is(unicode.White_Space, r) {
			pendingSpace = b.Len() > 0
			continue
		}
		if pendingSpace {
			b.WriteByte(' ')
			pendingSpace = false
		}
		b.WriteRune(r)
	}
	return b.String()
}

// writeString writes s as a JSON string, escaping only the quote, the
// backslash, and U+0000 through U+001F, which is always written \u00xx in
// lowercase hex - never the short forms like \n - as HASHING.md specifies.
func writeString(b *strings.Builder, s string) {
	const hexDigits = "0123456789abcdef"
	b.WriteByte('"')
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == '"':
			b.WriteString(`\"`)
		case r == '\\':
			b.WriteString(`\\`)
		case r < 0x20:
			b.WriteString(`\u00`)
			b.WriteByte(hexDigits[r>>4])
			b.WriteByte(hexDigits[r&0xf])
		default:
			b.WriteString(s[i : i+size])
		}
		i += size
	}
	b.WriteByte('"')
}
