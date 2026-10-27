// Package gitcheck holds the logic behind `sedum check`: it reads the findings
// provenance reports for a commit and decides which of them sedum can prove to
// be its own regenerations or marker migrations.
package gitcheck

import (
	"bytes"
	"strings"
)

// Region is one sedum region of a file: the values of its opening line and the
// lines it holds.
type Region struct {
	Tier   string
	Record string
	Gen    string
	Pkg    string
	// Prefix is the text of the opening line before its JSON object.
	Prefix string
	// Body is the lines between the opening and closing lines, each followed
	// by a newline.
	Body string
	// Close is the closing line as it is in the file.
	Close string
}

// regionPlaceholder stands in for a whole region in the text outside regions,
// so that two files compare by it.
const regionPlaceholder = "// sedum-region"

func isSpace(c byte) bool  { return c == ' ' || c == '\t' || c == '\r' }
func isIndent(c byte) bool { return c == ' ' || c == '\t' }

func trimRightSpace(s string) string {
	end := len(s)
	for end > 0 && isSpace(s[end-1]) {
		end--
	}
	return s[:end]
}

func skipSpace(s string, pos int) int {
	for pos < len(s) && isSpace(s[pos]) {
		pos++
	}
	return pos
}

// ReadRegions splits content into the text outside regions, with each whole
// region replaced by one fixed line, and the regions in file order. The bool is
// false when the file cannot be read: a region is never closed, or opens inside
// another. A line that looks like an opening line but does not parse is text.
func ReadRegions(content []byte) (string, []Region, bool) {
	var lines []string
	text := string(content)
	start := 0
	for i := 0; i <= len(text); i++ {
		if i == len(text) || text[i] == '\n' {
			lines = append(lines, text[start:i])
			start = i + 1
		}
	}

	var out []string
	regions := []Region{}
	var open *Region
	var leader, action, body string
	for _, line := range lines {
		if open == nil {
			if l, a, r, ok := parseOpening(line); ok {
				leader, action, body = l, a, ""
				open = &r
				out = append(out, regionPlaceholder)
				continue
			}
			out = append(out, line)
			continue
		}
		if closesRegion(line, leader, action) {
			open.Body = body
			open.Close = line
			regions = append(regions, *open)
			open = nil
			continue
		}
		if _, _, _, ok := parseOpening(line); ok {
			return "", nil, false
		}
		body += line + "\n"
	}
	if open != nil {
		return "", nil, false
	}
	return strings.Join(out, "\n"), regions, true
}

// parseOpening reads an opening line: indentation, a comment leader, " sedum:",
// an action name, one space and a JSON object carrying string values for tier,
// record, gen and pkg.
func parseOpening(line string) (leader, action string, r Region, ok bool) {
	i := 0
	for i < len(line) && isIndent(line[i]) {
		i++
	}
	rest := line[i:]
	switch {
	case strings.HasPrefix(rest, "//"):
		leader = "//"
	case strings.HasPrefix(rest, "#"):
		leader = "#"
	default:
		return "", "", Region{}, false
	}
	pos := i + len(leader)
	const marker = " sedum:"
	if !strings.HasPrefix(line[pos:], marker) {
		return "", "", Region{}, false
	}
	pos += len(marker)
	sp := strings.IndexByte(line[pos:], ' ')
	if sp <= 0 {
		return "", "", Region{}, false
	}
	action = line[pos : pos+sp]
	jsonStart := pos + sp + 1
	if jsonStart >= len(line) || line[jsonStart] != '{' {
		return "", "", Region{}, false
	}
	values, ok := parseObject(trimRightSpace(line[jsonStart:]))
	if !ok {
		return "", "", Region{}, false
	}
	tier, ok1 := values["tier"]
	record, ok2 := values["record"]
	gen, ok3 := values["gen"]
	pkg, ok4 := values["pkg"]
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return "", "", Region{}, false
	}
	return leader, action, Region{Tier: tier, Record: record, Gen: gen, Pkg: pkg, Prefix: line[:jsonStart]}, true
}

// closesRegion reports whether line closes a region opened with leader and action.
func closesRegion(line, leader, action string) bool {
	i := 0
	for i < len(line) && isIndent(line[i]) {
		i++
	}
	rest := line[i:]
	if len(rest) < len(leader)+1 || rest[:len(leader)] != leader || rest[len(leader)] != ' ' {
		return false
	}
	rest = rest[len(leader)+1:]
	want := "/sedum:" + action
	if !strings.HasPrefix(rest, want) {
		return false
	}
	return trimRightSpace(rest[len(want):]) == ""
}

// parseObject reads one JSON object and returns its string-valued members;
// members of any other type are skipped.
func parseObject(s string) (map[string]string, bool) {
	values := map[string]string{}
	i := skipSpace(s, 0)
	if i >= len(s) || s[i] != '{' {
		return nil, false
	}
	i = skipSpace(s, i+1)
	if i < len(s) && s[i] == '}' {
		i++
	} else {
		for {
			key, next, ok := parseString(s, i)
			if !ok {
				return nil, false
			}
			i = skipSpace(s, next)
			if i >= len(s) || s[i] != ':' {
				return nil, false
			}
			i = skipSpace(s, i+1)
			if i < len(s) && s[i] == '"' {
				v, next, ok := parseString(s, i)
				if !ok {
					return nil, false
				}
				values[key] = v
				i = next
			} else {
				next, ok := skipValue(s, i)
				if !ok {
					return nil, false
				}
				i = next
			}
			i = skipSpace(s, i)
			if i < len(s) && s[i] == ',' {
				i = skipSpace(s, i+1)
				continue
			}
			if i < len(s) && s[i] == '}' {
				i++
				break
			}
			return nil, false
		}
	}
	if skipSpace(s, i) != len(s) {
		return nil, false
	}
	return values, true
}

// parseString reads a JSON string starting at pos and returns its value and the
// position after the closing quote.
func parseString(s string, pos int) (string, int, bool) {
	if pos >= len(s) || s[pos] != '"' {
		return "", pos, false
	}
	var buf []byte
	i := pos + 1
	for i < len(s) {
		c := s[i]
		if c == '"' {
			return string(buf), i + 1, true
		}
		if c == '\\' {
			if i+1 >= len(s) {
				return "", pos, false
			}
			switch e := s[i+1]; e {
			case 'n':
				buf = append(buf, '\n')
			case 't':
				buf = append(buf, '\t')
			case 'r':
				buf = append(buf, '\r')
			case 'b', 'f':
			case 'u':
				if i+5 >= len(s) {
					return "", pos, false
				}
				i += 4
			default:
				buf = append(buf, e)
			}
			i += 2
			continue
		}
		buf = append(buf, c)
		i++
	}
	return "", pos, false
}

// skipValue skips one JSON value that is not a string and returns the position
// after it.
func skipValue(s string, pos int) (int, bool) {
	depth := 0
	i := pos
	for i < len(s) {
		c := s[i]
		if c == '"' {
			_, next, ok := parseString(s, i)
			if !ok {
				return pos, false
			}
			i = next
			continue
		}
		if c == '{' || c == '[' {
			depth++
		} else if c == '}' || c == ']' {
			if depth == 0 {
				break
			}
			depth--
		} else if c == ',' && depth == 0 {
			break
		}
		i++
	}
	if i == pos || depth != 0 {
		return pos, false
	}
	return i, true
}

// SameShape reports whether the commit's content has the same regions as the
// base's, in the same order, with the same text outside them. The base must be
// readable and hold at least one region. Each region keeps its tier, record and
// gen; with strict, its pkg, body, closing line and opening-line prefix as well.
func SameShape(base, commit []byte, strict bool) bool {
	baseOutside, baseRegions, ok := ReadRegions(base)
	if !ok || len(baseRegions) == 0 {
		return false
	}
	commitOutside, commitRegions, ok := ReadRegions(commit)
	if !ok {
		return false
	}
	if !bytes.Equal([]byte(baseOutside), []byte(commitOutside)) || len(baseRegions) != len(commitRegions) {
		return false
	}
	for i, a := range baseRegions {
		b := commitRegions[i]
		if a.Tier != b.Tier || a.Record != b.Record || a.Gen != b.Gen {
			return false
		}
		if strict && (a.Pkg != b.Pkg || a.Body != b.Body || a.Close != b.Close || a.Prefix != b.Prefix) {
			return false
		}
	}
	return true
}
