// Package genhash computes a provenance record's gen_hash: the hash a marker
// stores as "gen" so a later grow can tell whether the record that wrote a
// region has changed since (prov-2026-66bea9f3).
//
// gen_hash covers exactly the fields Sedum reads from a record - intent,
// constraints, affected_scope and forbidden_scope - and nothing else, so a
// status flip, a seal or a re-pin never looks like a content change.
//
// The definition is not Sedum's. It belongs to the shared provenance record
// library, whose conformance corpus keeps every implementation producing the
// same bytes, and Of delegates to it (prov-2026-dac49f78). This package stays
// as the seam: Sedum's callers depend on Input and Of, not on the library.
package genhash

import "github.com/livecodelife/provenance/record"

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
	return record.GenHash(record.GenInput(in))
}
