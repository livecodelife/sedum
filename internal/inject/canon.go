package inject

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

// The canonical forms a marker's two hashes are taken over
// (prov-2026-8bba6ce4). Both are computable by any tool from the rules here, so
// a foreign writer can produce and check them without linking Sedum.

// maxExactInteger is 2^53, the largest integer magnitude a canonical number
// may carry. Beyond it a float64 cannot tell neighbouring integers apart, so
// two writers could disagree about which integer a spelling means.
const maxExactInteger = 9007199254740992

// CanonicalKwargs encodes a resolved kwargs map as the canonical JSON bytes a
// kwargs hash is taken over: sorted keys, no insignificant whitespace, one
// spelling per number.
func CanonicalKwargs(kwargs map[string]any) ([]byte, error) {
	if len(kwargs) == 0 {
		return []byte("{}"), nil
	}

	encoded, err := json.Marshal(kwargs)
	if err != nil {
		return nil, fmt.Errorf("kwargs cannot be hashed: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return nil, fmt.Errorf("kwargs cannot be hashed: %w", err)
	}

	var buf bytes.Buffer
	if err := writeCanonical(&buf, decoded, ""); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// KwargsHash returns the sha256: hash of a resolved kwargs map's canonical
// encoding.
func KwargsHash(kwargs map[string]any) (string, error) {
	canonical, err := CanonicalKwargs(kwargs)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// BodyHash returns the sha256: hash of a region's body text with each CRLF read
// as LF. Nothing else is changed: indentation, trailing spaces and trailing
// blank lines all take part.
func BodyHash(body string) string {
	normalized := strings.ReplaceAll(body, "\r\n", "\n")
	sum := sha256.Sum256([]byte(normalized))
	return "sha256:" + hex.EncodeToString(sum[:])
}

// writeCanonical appends the canonical JSON of one decoded value to buf.
func writeCanonical(buf *bytes.Buffer, value any, path string) error {
	switch v := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		buf.WriteByte('{')
		for i, key := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			writeCanonicalString(buf, key)
			buf.WriteByte(':')
			child := key
			if path != "" {
				child = path + "." + key
			}
			if err := writeCanonical(buf, v[key], child); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
	case []any:
		buf.WriteByte('[')
		for i, element := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := writeCanonical(buf, element, path+"["+strconv.Itoa(i)+"]"); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
	case string:
		writeCanonicalString(buf, v)
	case bool:
		if v {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case nil:
		buf.WriteString("null")
	case json.Number:
		text, err := canonicalNumber(v, path)
		if err != nil {
			return err
		}
		buf.WriteString(text)
	default:
		return fmt.Errorf("kwargs cannot be hashed: %s: unsupported value", pathText(path))
	}
	return nil
}

// pathText names a path in an error, kwargs for the root.
func pathText(path string) string {
	if path == "" {
		return "kwargs"
	}
	return path
}

// canonicalNumber returns one JSON number in canonical form, or refuses it.
func canonicalNumber(n json.Number, path string) (string, error) {
	text := string(n)
	if !isJSONNumber(text) {
		return "", fmt.Errorf("kwargs cannot be hashed: %s: not a number", pathText(path))
	}

	// Finiteness is decided on the nearest float64 before size is considered,
	// so 1e400 reports that it is not finite rather than that it is large.
	f, err := strconv.ParseFloat(text, 64)
	if err != nil && (f > 1 || f < -1) {
		return "", fmt.Errorf("kwargs cannot be hashed: %s: number is not finite", pathText(path))
	}
	if f == 0 {
		return "0", nil
	}

	if rat, ok := new(big.Rat).SetString(text); ok && rat.IsInt() {
		integer := rat.Num()
		limit := big.NewInt(maxExactInteger)
		if new(big.Int).Abs(integer).Cmp(limit) > 0 {
			return "", fmt.Errorf("kwargs cannot be hashed: %s: integer beyond 2^53", pathText(path))
		}
		return integer.String(), nil
	}
	return formatECMAScript(f), nil
}

// isJSONNumber reports whether text is a number in JSON's grammar.
func isJSONNumber(text string) bool {
	i := 0
	if i < len(text) && text[i] == '-' {
		i++
	}
	digits := func() int {
		start := i
		for i < len(text) && text[i] >= '0' && text[i] <= '9' {
			i++
		}
		return i - start
	}
	switch {
	case i < len(text) && text[i] == '0':
		i++
	case digits() == 0:
		return false
	}
	if i < len(text) && text[i] == '.' {
		i++
		if digits() == 0 {
			return false
		}
	}
	if i < len(text) && (text[i] == 'e' || text[i] == 'E') {
		i++
		if i < len(text) && (text[i] == '+' || text[i] == '-') {
			i++
		}
		if digits() == 0 {
			return false
		}
	}
	return i == len(text)
}

// formatECMAScript renders a finite, non-zero float64 the way ECMAScript's
// Number::toString does: the shortest digits that read back as the same float.
func formatECMAScript(f float64) string {
	abs := f
	if abs < 0 {
		abs = -abs
	}
	if abs >= 0.000001 && abs < 1e21 {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}

	text := strconv.FormatFloat(f, 'e', -1, 64)
	mantissa, exponent, _ := strings.Cut(text, "e")
	sign := exponent[0]
	exponent = strings.TrimLeft(exponent[1:], "0")
	return mantissa + "e" + string(sign) + exponent
}

// writeCanonicalString appends a string as a JSON string with minimal
// escaping.
func writeCanonicalString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			buf.WriteString(`\"`)
		case c == '\\':
			buf.WriteString(`\\`)
		case c == '\b':
			buf.WriteString(`\b`)
		case c == '\f':
			buf.WriteString(`\f`)
		case c == '\n':
			buf.WriteString(`\n`)
		case c == '\r':
			buf.WriteString(`\r`)
		case c == '\t':
			buf.WriteString(`\t`)
		case c < 0x20:
			buf.WriteString(`\u00`)
			buf.WriteString(strconv.FormatInt(int64(c)>>4, 16))
			buf.WriteString(strconv.FormatInt(int64(c)&0xf, 16))
		default:
			buf.WriteByte(c)
		}
	}
	buf.WriteByte('"')
}
