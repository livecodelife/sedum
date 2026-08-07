package filetmpl

import (
	"fmt"
	"path"
	"strings"
)

// A segment is an alternating sequence of literal text and wildcards, because
// captures are sub-segment: {name}_controller.rb is literal text around a
// capture, not a segment that happens to be a capture. Two wildcards may not be
// adjacent (prov-2026-bc163215), so a literal always separates them, which is
// what makes both matching and ranking a positional walk over pieces.
type segment []piece

type piece struct {
	kind pieceKind
	// Literal text for a literal piece, the capture name for a capture,
	// empty for a glob.
	text string
}

// The declared order is the specificity order: literal is more specific than
// capture, which is more specific than glob. The ranking in PRD.md compares
// these tiers directly, so the values are load-bearing, not incidental.
type pieceKind int

const (
	literal pieceKind = iota
	capture
	glob
)

func (k pieceKind) String() string {
	switch k {
	case capture:
		return "capture"
	case glob:
		return "glob"
	default:
		return "literal"
	}
}

// pattern is a parsed file-template path pattern. A template's own path is its
// pattern; there is no mapping table and no registry of file types.
type pattern struct {
	// raw is the pattern exactly as the caller supplied it, so diagnostics
	// name what the author wrote rather than a normalized rewrite of it.
	raw      string
	segments []segment
}

// normalize puts a pattern or a path into slash-separated form so that a
// package authored on one platform matches identically on another.
func normalize(p string) string {
	p = strings.ReplaceAll(p, `\`, "/")
	if p == "" {
		return ""
	}
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	if p == "." {
		return ""
	}
	return p
}

func parse(raw string) (*pattern, error) {
	norm := normalize(raw)
	if norm == "" {
		return nil, fmt.Errorf("file template pattern is empty")
	}

	parts := strings.Split(norm, "/")
	p := &pattern{raw: raw, segments: make([]segment, 0, len(parts))}
	names := make(map[string]bool)

	for _, part := range parts {
		seg, err := parseSegment(part, names)
		if err != nil {
			return nil, fmt.Errorf("file template pattern %q: %w", raw, err)
		}
		p.segments = append(p.segments, seg)
	}
	return p, nil
}

// parseSegment splits one path segment into its pieces. names carries the
// capture names seen across the whole pattern, so a name reused in a later
// segment is caught too.
func parseSegment(s string, names map[string]bool) (segment, error) {
	if s == "" {
		return nil, fmt.Errorf("segment is empty")
	}

	var (
		seg segment
		lit strings.Builder
	)
	flush := func() {
		if lit.Len() > 0 {
			seg = append(seg, piece{kind: literal, text: lit.String()})
			lit.Reset()
		}
	}
	// A wildcard may only follow literal text. Adjacent wildcards have no
	// non-arbitrary split point, so they are rejected rather than resolved by
	// a convention nobody wrote down.
	addWildcard := func(p piece, shown string) error {
		flush()
		if n := len(seg); n > 0 && seg[n-1].kind != literal {
			return fmt.Errorf("%s follows a %s in segment %q with no literal text between them, so where one ends and the next begins is ambiguous", shown, seg[n-1].kind, s)
		}
		seg = append(seg, p)
		return nil
	}

	for i := 0; i < len(s); {
		switch s[i] {
		case '{':
			end := strings.IndexByte(s[i:], '}')
			if end < 0 {
				return nil, fmt.Errorf("segment %q has an unclosed '{'", s)
			}
			name := s[i+1 : i+end]
			if name == "" {
				return nil, fmt.Errorf("segment %q has a capture with no name", s)
			}
			if strings.ContainsAny(name, "{*") {
				return nil, fmt.Errorf("segment %q has a capture named %q; a capture name may not contain '{' or '*'", s, name)
			}
			if names[name] {
				// Not a backreference: two pieces named the same almost
				// certainly meant two different values.
				return nil, fmt.Errorf("capture {%s} is declared twice; give the second one a different name", name)
			}
			names[name] = true
			if err := addWildcard(piece{kind: capture, text: name}, "capture {"+name+"}"); err != nil {
				return nil, err
			}
			i += end + 1
		case '*':
			if err := addWildcard(piece{kind: glob}, "glob '*'"); err != nil {
				return nil, err
			}
			i++
		case '}':
			return nil, fmt.Errorf("segment %q has a '}' with no matching '{'", s)
		default:
			lit.WriteByte(s[i])
			i++
		}
	}
	flush()
	return seg, nil
}

// match reports whether path matches, binding any captures. Because there is no
// segment-spanning wildcard, patterns and paths must agree on segment count.
func (p *pattern) match(target string) (map[string]string, bool) {
	parts := strings.Split(normalize(target), "/")
	if len(parts) != len(p.segments) {
		return nil, false
	}

	captures := make(map[string]string)
	for i, seg := range p.segments {
		if !matchSegment(seg, parts[i], captures) {
			return nil, false
		}
	}
	return captures, true
}

// matchSegment consumes pieces left to right. A wildcard takes at least one
// character, and where a literal follows it the literal binds at its leftmost
// viable occurrence, so a capture is the shortest text that still lets the rest
// of the pattern complete. Where that placement dead-ends, later occurrences
// are tried, so the result is "shortest that works" rather than "shortest".
//
// Captures written on an abandoned branch are always overwritten before a match
// is reported, because a complete match visits every piece.
func matchSegment(seg segment, s string, captures map[string]string) bool {
	if len(seg) == 0 {
		return s == ""
	}

	head, rest := seg[0], seg[1:]
	if head.kind == literal {
		if !strings.HasPrefix(s, head.text) {
			return false
		}
		return matchSegment(rest, s[len(head.text):], captures)
	}

	bind := func(text string) {
		if head.kind == capture {
			captures[head.text] = text
		}
	}

	// A trailing wildcard takes the rest of the segment, which must be at
	// least one character.
	if len(rest) == 0 {
		if s == "" {
			return false
		}
		bind(s)
		return true
	}

	// Parsing guarantees a literal follows a wildcard.
	next := rest[0].text
	for i := 1; i+len(next) <= len(s); i++ {
		if s[i:i+len(next)] != next {
			continue
		}
		bind(s[:i])
		if matchSegment(rest[1:], s[i+len(next):], captures) {
			return true
		}
	}
	return false
}

// compare orders two patterns by specificity, most specific first, the way
// sort.Slice comparators read: negative means a is more specific than b.
//
// Segments are compared leftmost first and the first difference decides, so a
// pattern that is more specific early beats one that is more specific late.
// The order is total and independent of the order patterns were supplied in.
func compare(a, b *pattern) int {
	for i := range min(len(a.segments), len(b.segments)) {
		if c := compareSegment(a.segments[i], b.segments[i]); c != 0 {
			return c
		}
	}
	// Only reachable between patterns that cannot match a common path, since
	// no wildcard spans segments. Ordered anyway to keep compare total.
	return len(a.segments) - len(b.segments)
}

func compareSegment(a, b segment) int {
	for i := range max(len(a), len(b)) {
		switch {
		case i >= len(a):
			// One ran out of pieces. An extra literal makes b more
			// specific; an extra wildcard makes it less.
			if b[i].kind == literal {
				return 1
			}
			return -1
		case i >= len(b):
			if a[i].kind == literal {
				return -1
			}
			return 1
		}

		pa, pb := a[i], b[i]
		if pa.kind != pb.kind {
			// literal beats capture beats glob.
			return int(pa.kind) - int(pb.kind)
		}
		if pa.kind == literal && len(pa.text) != len(pb.text) {
			// Longer literal beats shorter at the same position.
			return len(pb.text) - len(pa.text)
		}
	}
	return 0
}
