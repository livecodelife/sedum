package inject

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// The ownership marker: the only place idempotency state lives.
//
// A region is written between a pair of markers naming the action that produced
// it, the ownership tier, the kwargs it was rendered from, and the record that
// last parameterized it:
//
//	# sedum:createControllerMethod:index {"tier":"owned","record":"PR-014","kwargs":{"controller":"users"}}
//	def index
//	  @users = User.all
//	end
//	# /sedum:createControllerMethod:index
//
// The action and variant stay literal on the line because they are the audit
// trail: grepping markers yields file -> action -> variant with no maintained
// state. Everything else is one JSON object, because the parser has to tolerate
// fields it does not recognize and default the ones that are absent
// (prov-2026-36c8a99c). A marker sits in a generated codebase long after the
// version that wrote it is gone, so the reader and the writer of any given
// marker are routinely different versions of Sedum. Positional fields would
// make every field added later a migration across every repository that already
// carries markers; a JSON object makes it an addition.
//
// There is a concrete field waiting: cross-region ordering wants markers to
// record what a region exposes and what it consumes, so a dependency graph can
// be built by scanning markers rather than by parsing the target language.

// Tier declares whether Sedum may overwrite a region.
type Tier string

const (
	// TierOwned means Sedum generated the region and replaces it on every
	// run. Hand edits are lost.
	TierOwned Tier = "owned"
	// TierSeeded means Sedum generated the region once and never touches it
	// again. Present in the file, skipped on rerun.
	TierSeeded Tier = "seeded"
)

// DefaultTier is what an absent tier field means. Only owned is exercised by
// this milestone; seeded is honored on read so that a region which has stopped
// being Sedum's to overwrite can say so.
const DefaultTier = TierOwned

const (
	openKeyword  = "sedum:"
	closeKeyword = "/sedum:"

	// anchorKeyword is the label a file template plants to declare an anchor
	// point: "sedum:anchor:class_body". Anchor declarations and ownership
	// markers share the "sedum:" namespace, so the reader has to tell them
	// apart - otherwise an anchor a template planted reads as an ownership
	// marker for an action named "anchor", and the first region found after
	// it appears to be closing that one.
	//
	// The consequence is that "anchor" is not available as an action name.
	anchorKeyword = "anchor"
)

// Marker is one region's opening marker.
type Marker struct {
	Action string
	// Variant is the discriminator value that selected the template, or
	// empty for an action that declares no discriminator.
	Variant string
	Tier    Tier
	// Record is the provenance record that last wrote the region. It is
	// recorded and never matched on: a region is identified by what it is,
	// not by who wrote it, so that a later record refining an earlier region
	// replaces it in place instead of minting a duplicate beside it.
	Record string
	// Kwargs is every kwarg the region was rendered from. All of them are
	// recorded rather than the subset that currently selects a region,
	// because deciding today which kwargs select and which parameterize
	// would be guessing, and recording all of them keeps the question open.
	Kwargs map[string]any
}

// attrs is the JSON object carried on an opening marker.
//
// Decoding is deliberately lenient in both directions: encoding/json ignores a
// key this struct does not declare, and a key this struct declares but the
// marker omits keeps its zero value, which normalize turns into the documented
// default.
type attrs struct {
	Tier   Tier           `json:"tier"`
	Record string         `json:"record,omitempty"`
	Kwargs map[string]any `json:"kwargs,omitempty"`
}

// Label is the action:variant pair as it appears on the marker line.
func (m Marker) Label() string {
	if m.Variant == "" {
		return m.Action
	}
	return m.Action + ":" + m.Variant
}

// Open renders the opening marker line, without a trailing newline.
//
// HTML escaping is turned off. encoding/json escapes &, < and > by default,
// which is meant for JSON embedded in a web page and is wrong here: a kwarg
// holding "&t.ID" would be recorded as a & escape, and a Go generic bound
// or a C++ template argument would fare worse. The escaped form round-trips
// correctly, so nothing breaks - but a marker is read by people and by grep,
// and it should say what the region was rendered from.
func (m Marker) Open(commentPrefix string) (string, error) {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(attrs{
		Tier:   m.tierOrDefault(),
		Record: m.Record,
		Kwargs: m.Kwargs,
	}); err != nil {
		return "", fmt.Errorf("action %s: kwargs cannot be recorded on a marker: %w", m.Label(), err)
	}

	// Encode terminates the value with a newline, which would split the
	// marker across two lines.
	encoded := strings.TrimRight(buf.String(), "\n")
	return commentPrefix + " " + openKeyword + m.Label() + " " + encoded, nil
}

// Close renders the closing marker line, without a trailing newline.
func (m Marker) Close(commentPrefix string) string {
	return commentPrefix + " " + closeKeyword + m.Label()
}

func (m Marker) tierOrDefault() Tier {
	if m.Tier == "" {
		return DefaultTier
	}
	return m.Tier
}

// parseOpen reads an opening marker from one line, reporting whether the line
// is one at all.
//
// A line that is not a marker is not an error - most lines in a generated file
// are not markers. A line that is a marker but carries an unreadable attribute
// object is an error, because that is corruption rather than version skew: the
// object's encoding is the committed shape, and a newer version adding a field
// still produces JSON an older one can read.
func parseOpen(commentPrefix, line string) (Marker, bool, error) {
	rest, ok := trimMarkerPrefix(commentPrefix, line, openKeyword)
	if !ok {
		return Marker{}, false, nil
	}

	label, encoded, _ := strings.Cut(rest, " ")
	action, variant := splitLabel(strings.TrimSpace(label))
	if action == "" || action == anchorKeyword {
		return Marker{}, false, nil
	}

	marker := Marker{Action: action, Variant: variant, Tier: DefaultTier}

	encoded = strings.TrimSpace(encoded)
	if encoded == "" {
		// Every field is absent, so every field takes its default. A marker
		// is never rejected for carrying less than the reading version
		// expects.
		return marker, true, nil
	}

	var a attrs
	if err := json.Unmarshal([]byte(encoded), &a); err != nil {
		return Marker{}, false, fmt.Errorf(
			"marker for %s carries attributes that are not readable as JSON: %w", marker.Label(), err)
	}
	if a.Tier != "" {
		marker.Tier = a.Tier
	}
	marker.Record = a.Record
	marker.Kwargs = a.Kwargs
	return marker, true, nil
}

// parseClose reads a closing marker's label from one line.
func parseClose(commentPrefix, line string) (string, bool) {
	rest, ok := trimMarkerPrefix(commentPrefix, line, closeKeyword)
	if !ok {
		return "", false
	}
	return strings.TrimSpace(rest), true
}

// trimMarkerPrefix matches "<comment_prefix> <keyword>" at the start of a
// line's trimmed text and returns what follows.
//
// The comment prefix is the package's declared one and is never hardcoded and
// never inferred from a file extension: #, // and -- all appear across targets.
func trimMarkerPrefix(commentPrefix, line, keyword string) (string, bool) {
	trimmed := strings.TrimSpace(line)

	after, ok := strings.CutPrefix(trimmed, commentPrefix)
	if !ok {
		return "", false
	}
	after, ok = strings.CutPrefix(strings.TrimSpace(after), keyword)
	if !ok {
		return "", false
	}
	return after, true
}

func splitLabel(label string) (action, variant string) {
	action, variant, _ = strings.Cut(label, ":")
	return action, variant
}

// Identity is what makes two invocations the same region.
//
// The record ID takes no part in it. Under record-scoped identity a later
// record could never claim a region an earlier one wrote - it would always mint
// a second region beside the first, and a record whose intent is "PUT should
// support partial updates" would produce two definitions of one function rather
// than a replacement.
//
// Which kwargs select a region and which merely parameterize it is a rule in
// code, not a fact on disk: the marker records all of them so the rule can be
// redefined without migrating generated codebases. The rule this milestone
// implements is that an action's required kwargs select and its optional kwargs
// parameterize. A required kwarg is one the action cannot be invoked without,
// which makes it definitional; an optional one refines what is already
// identified. That is what makes two before_action filters in one controller
// distinct regions while re-invoking one of them with a different "only" list
// replaces it in place.
type Identity struct {
	Action  string
	Variant string
	// Key is a canonical encoding of the selecting kwargs.
	Key string
}

// IdentityOf computes a marker's identity under the given selecting kwarg
// names.
func IdentityOf(m Marker, selecting []string) (Identity, error) {
	selected := map[string]any{}
	for _, name := range selecting {
		if value, ok := m.Kwargs[name]; ok {
			selected[name] = value
		}
	}

	// Round-tripping normalizes what a value was built as into what it would
	// have been read back as, so a fixture's []string and a parsed marker's
	// []any compare equal. json.Marshal sorts map keys, which is what makes
	// the result canonical.
	encoded, err := json.Marshal(selected)
	if err != nil {
		return Identity{}, fmt.Errorf("action %s: kwargs cannot be compared: %w", m.Label(), err)
	}
	var normalized map[string]any
	if err := json.Unmarshal(encoded, &normalized); err != nil {
		return Identity{}, fmt.Errorf("action %s: kwargs cannot be compared: %w", m.Label(), err)
	}
	canonical, err := json.Marshal(normalized)
	if err != nil {
		return Identity{}, fmt.Errorf("action %s: kwargs cannot be compared: %w", m.Label(), err)
	}

	return Identity{Action: m.Action, Variant: m.Variant, Key: string(canonical)}, nil
}

// Region is one marked block found in a file.
type Region struct {
	Marker Marker
	// Start is the byte offset of the first byte of the opening marker line.
	Start int
	// End is the byte offset one past the newline that ends the closing
	// marker line, so content[Start:End] is the whole region.
	End int
}

// FindRegions returns every marked region in content, in the order they appear.
//
// A region whose opening marker is never closed is an error: its extent is
// unknown, so replacing it would either destroy the rest of the file or append
// beside it, and both are worse than halting.
func FindRegions(commentPrefix, content string) ([]Region, error) {
	var (
		out  []Region
		open *Region
	)

	for offset := 0; offset < len(content); {
		end := bytes.IndexByte([]byte(content[offset:]), '\n')
		lineEnd := len(content)
		next := len(content)
		if end >= 0 {
			lineEnd = offset + end
			next = lineEnd + 1
		}
		line := content[offset:lineEnd]

		if open == nil {
			marker, ok, err := parseOpen(commentPrefix, line)
			if err != nil {
				return nil, err
			}
			if ok {
				open = &Region{Marker: marker, Start: offset}
			}
		} else if label, ok := parseClose(commentPrefix, line); ok {
			if label != open.Marker.Label() {
				return nil, fmt.Errorf(
					"region opened by marker %q is closed by marker %q; a region's markers must name the same action and variant",
					open.Marker.Label(), label)
			}
			open.End = next
			out = append(out, *open)
			open = nil
		}

		offset = next
	}

	if open != nil {
		return nil, fmt.Errorf(
			"region opened by marker %q is never closed; its extent is unknown, so it cannot be replaced",
			open.Marker.Label())
	}
	return out, nil
}
