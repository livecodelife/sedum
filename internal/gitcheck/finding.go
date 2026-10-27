package gitcheck

import (
	"regexp"
	"strconv"
	"strings"
)

// Finding is one provenance result read as a finding.
type Finding struct {
	RuleID   string
	Message  string
	File     string
	RecordID string
}

// FindingOf reads a result as a finding. File is empty when the result has no
// location, and RecordID when it has no recordId property.
func FindingOf(r Result) Finding {
	return Finding{
		RuleID:   r.RuleID(),
		Message:  r.MessageText(),
		File:     r.FirstLocationURI(),
		RecordID: r.Property("recordId"),
	}
}

const refusalText = "is implemented: run 'provenance revise"

// IsRefusal reports whether f is the refusal that a tagged record is implemented
// and needs a revision: rule STATE001 with that message text.
func IsRefusal(f Finding) bool {
	return f.RuleID == "STATE001" && strings.Contains(f.Message, refusalText)
}

// Allowed reports whether the proof may waive f for the tagged record. The
// proof is "regeneration", which waives the refusal of the tagged record, or
// "migration", which waives that refusal and the SCOPE001 findings of the tagged
// record. A finding without a record, of another record, or of any other rule is
// never waived.
func Allowed(proof string, f Finding, tagRecord string) bool {
	if f.RecordID == "" || f.RecordID != tagRecord {
		return false
	}
	switch proof {
	case "regeneration":
		return IsRefusal(f)
	case "migration":
		return IsRefusal(f) || f.RuleID == "SCOPE001"
	}
	return false
}

// Tag is a record tag of a commit message: [id] or [id@N].
type Tag struct {
	ID         string
	Version    int
	HasVersion bool
}

var tagPattern = regexp.MustCompile(`\[(prov-[0-9]{4}-(?:[0-9a-f]{8}|[0-9]{3})(?:-[a-z0-9-]+)?)(?:@([0-9]+))?\]`)

// ParseTags returns the tags of a commit message in message order. A tag is the
// bracket form [id] or [id@N] anywhere in the message, N a positive integer.
func ParseTags(message string) []Tag {
	var tags []Tag
	for _, m := range tagPattern.FindAllStringSubmatch(message, -1) {
		tag := Tag{ID: m[1]}
		if m[2] != "" {
			n, err := strconv.Atoi(m[2])
			if err != nil || n < 1 {
				continue
			}
			tag.Version, tag.HasVersion = n, true
		}
		tags = append(tags, tag)
	}
	return tags
}
