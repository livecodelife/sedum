package inject

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// The kwargs hash and the body hash are computable by any tool from the same
// rules (prov-2026-8bba6ce4). The expected canonical bytes below are derived by
// hand from the rules, and the expected hashes are taken over those bytes by an
// independent SHA-256, never by calling the function under test.

const (
	hashOfEmptyText   = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	hashOfEmptyObject = "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a"
)

// sha256Of is the hash text a marker carries, taken over exactly the bytes given.
func sha256Of(text string) string {
	sum := sha256.Sum256([]byte(text))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func TestCanonicalKwargs(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]any
		want string
	}{
		{"nil map", nil, `{}`},
		{"empty map", map[string]any{}, `{}`},
		{"one string", map[string]any{"a": "x"}, `{"a":"x"}`},
		{"keys sort whatever the insertion order",
			map[string]any{"b": 1, "a": 2, "c": 3}, `{"a":2,"b":1,"c":3}`},
		{"nested objects sort at every depth",
			map[string]any{"b": map[string]any{"z": 1, "a": 2}, "a": []any{map[string]any{"y": 1, "x": 2}}},
			`{"a":[{"x":2,"y":1}],"b":{"a":2,"z":1}}`},
		{"arrays keep their order", map[string]any{"l": []string{"c", "a", "b"}}, `{"l":["c","a","b"]}`},
		{"empty array and empty object", map[string]any{"l": []any{}, "o": map[string]any{}}, `{"l":[],"o":{}}`},
		{"explicit null is present", map[string]any{"a": nil}, `{"a":null}`},
		{"bools", map[string]any{"t": true, "f": false}, `{"f":false,"t":true}`},
		{"go int", map[string]any{"n": 1}, `{"n":1}`},
		{"go float that is an integer", map[string]any{"n": 1.0}, `{"n":1}`},
		{"go float64 fraction", map[string]any{"n": 0.5}, `{"n":0.5}`},
		{"small float uses an exponent", map[string]any{"n": 0.00000015}, `{"n":1.5e-7}`},
		{"negative integer", map[string]any{"n": -5}, `{"n":-5}`},
		{"HTML characters are written as themselves", map[string]any{"s": "<a&b>"}, `{"s":"<a&b>"}`},
		{"no whitespace anywhere", map[string]any{"a": []any{1, map[string]any{"b": "c d"}}}, `{"a":[1,{"b":"c d"}]}`},
		{"keys sort by UTF-8 bytes, not UTF-16 units",
			// U+FF5E is EF BD 9E and U+10000 is F0 90 80 80: bytes put the
			// first before the second, UTF-16 code units put it after.
			map[string]any{"\U00010000": 1, "～": 2}, "{\"～\":2,\"\U00010000\":1}"},
		{"multibyte key sorts after ASCII", map[string]any{"é": 1, "z": 2, "a": 3}, `{"a":3,"z":2,"é":1}`},
		{"strings are not normalized",
			map[string]any{"composed": "é", "decomposed": "é"},
			"{\"composed\":\"é\",\"decomposed\":\"é\"}"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CanonicalKwargs(tc.in)
			if err != nil {
				t.Fatalf("CanonicalKwargs: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("CanonicalKwargs =\n  %s\nwant\n  %s", got, tc.want)
			}
		})
	}
}

// The empty map and the nil map are the two bytes {}, never the null that
// encoding/json writes for a nil map.
func TestCanonicalKwargsNilMapIsNotNull(t *testing.T) {
	got, err := CanonicalKwargs(nil)
	if err != nil {
		t.Fatalf("CanonicalKwargs: %v", err)
	}
	if !bytes.Equal(got, []byte("{}")) {
		t.Errorf("CanonicalKwargs(nil) = %q, want {}", got)
	}
}

func TestCanonicalKwargsHasNoTrailingNewline(t *testing.T) {
	got, err := CanonicalKwargs(map[string]any{"a": "x"})
	if err != nil {
		t.Fatalf("CanonicalKwargs: %v", err)
	}
	if bytes.HasSuffix(got, []byte("\n")) {
		t.Errorf("CanonicalKwargs = %q, want no trailing newline", got)
	}
}

func TestCanonicalKwargsDoesNotMutateItsArgument(t *testing.T) {
	in := map[string]any{
		"b": []any{"y", "x"},
		"a": map[string]any{"n": json.Number("1.0"), "s": "v"},
		"c": nil,
	}
	want := map[string]any{
		"b": []any{"y", "x"},
		"a": map[string]any{"n": json.Number("1.0"), "s": "v"},
		"c": nil,
	}
	if _, err := CanonicalKwargs(in); err != nil {
		t.Fatalf("CanonicalKwargs: %v", err)
	}
	if !reflect.DeepEqual(in, want) {
		t.Errorf("argument changed:\n got %#v\nwant %#v", in, want)
	}
}

// Two maps with equal keys and values encode alike whatever the Go types the
// values were built as.
func TestCanonicalKwargsIgnoresBuiltType(t *testing.T) {
	a, err := CanonicalKwargs(map[string]any{"l": []string{"a", "b"}, "n": 3, "m": map[string]any{"k": 1.0}})
	if err != nil {
		t.Fatalf("CanonicalKwargs: %v", err)
	}
	b, err := CanonicalKwargs(map[string]any{"m": map[string]any{"k": json.Number("1e0")}, "n": json.Number("3.0"), "l": []any{"a", "b"}})
	if err != nil {
		t.Fatalf("CanonicalKwargs: %v", err)
	}
	if !bytes.Equal(a, b) {
		t.Errorf("equal kwargs encoded differently:\n  %s\n  %s", a, b)
	}
}

// A number spelled 1, 1.0 or 1e0 is one integer, and every spelling of negative
// zero is 0.
func TestCanonicalKwargsNumberSpellings(t *testing.T) {
	tests := []struct {
		in   json.Number
		want string
	}{
		{"1", `{"n":1}`},
		{"1.0", `{"n":1}`},
		{"1e0", `{"n":1}`},
		{"1E0", `{"n":1}`},
		{"0.1e1", `{"n":1}`},
		{"10e-1", `{"n":1}`},
		{"-0", `{"n":0}`},
		{"-0.0", `{"n":0}`},
		{"-0e5", `{"n":0}`},
	}
	for _, tc := range tests {
		t.Run(string(tc.in), func(t *testing.T) {
			got, err := CanonicalKwargs(map[string]any{"n": tc.in})
			if err != nil {
				t.Fatalf("CanonicalKwargs: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestCanonicalKwargsErrors(t *testing.T) {
	const prefix = "kwargs cannot be hashed: "
	tests := []struct {
		name string
		in   map[string]any
		want string // required text after the prefix, as a prefix of the cause
	}{
		{"a function cannot be marshalled", map[string]any{"f": func() {}}, "json: unsupported type"},
		{"a channel cannot be marshalled", map[string]any{"c": make(chan int)}, "json: unsupported type"},
		{"NaN cannot be marshalled", map[string]any{"n": nan()}, "json: unsupported value"},
		{"integer beyond 2^53", map[string]any{"n": json.Number("9007199254740993")}, "n: integer beyond 2^53"},
		{"negative integer beyond 2^53", map[string]any{"n": json.Number("-9007199254740993")}, "n: integer beyond 2^53"},
		{"1e30 is an integer beyond 2^53", map[string]any{"n": json.Number("1e30")}, "n: integer beyond 2^53"},
		{"1e400 is not finite", map[string]any{"n": json.Number("1e400")}, "n: number is not finite"},
		{"path through objects and arrays",
			map[string]any{"a": []any{1, map[string]any{"b": json.Number("9007199254740993")}}},
			"a[1].b: integer beyond 2^53"},
		{"path through nested objects",
			map[string]any{"a": map[string]any{"b": map[string]any{"c": json.Number("1e400")}}},
			"a.b.c: number is not finite"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CanonicalKwargs(tc.in)
			if err == nil {
				t.Fatalf("CanonicalKwargs accepted the argument and returned %q", got)
			}
			if got != nil {
				t.Errorf("CanonicalKwargs returned %q with an error, want a nil slice", got)
			}
			if !strings.HasPrefix(err.Error(), prefix+tc.want) {
				t.Errorf("error = %q, want it to start with %q", err.Error(), prefix+tc.want)
			}
		})
	}
}

// The largest integer a double holds exactly is accepted, at either sign.
func TestCanonicalKwargsAcceptsTheLargestSafeInteger(t *testing.T) {
	got, err := CanonicalKwargs(map[string]any{"p": json.Number("9007199254740992"), "m": json.Number("-9007199254740992")})
	if err != nil {
		t.Fatalf("CanonicalKwargs: %v", err)
	}
	if want := `{"m":-9007199254740992,"p":9007199254740992}`; string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func nan() float64 {
	zero := 0.0
	return zero / zero
}

func TestWriteCanonical(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want string
	}{
		{"nil is null", nil, `null`},
		{"true", true, `true`},
		{"false", false, `false`},
		{"string", "a\"b", `"a\"b"`},
		{"number", json.Number("1.0"), `1`},
		{"empty map", map[string]any{}, `{}`},
		{"empty slice", []any{}, `[]`},
		{"map keys sort by bytes", map[string]any{"b": nil, "a": true, "B": false}, `{"B":false,"a":true,"b":null}`},
		{"map keys are escaped as strings", map[string]any{"a\"b": json.Number("2")}, `{"a\"b":2}`},
		{"slice keeps order", []any{"b", "a", json.Number("1")}, `["b","a",1]`},
		{"nesting", map[string]any{"a": []any{map[string]any{"c": nil, "b": []any{}}}}, `{"a":[{"b":[],"c":null}]}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := writeCanonical(&buf, tc.in, ""); err != nil {
				t.Fatalf("writeCanonical: %v", err)
			}
			if got := buf.String(); got != tc.want {
				t.Errorf("writeCanonical = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestWriteCanonicalAppendsToTheBuffer(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("prefix:")
	if err := writeCanonical(&buf, "x", ""); err != nil {
		t.Fatalf("writeCanonical: %v", err)
	}
	if got, want := buf.String(), `prefix:"x"`; got != want {
		t.Errorf("buffer = %q, want %q", got, want)
	}
}

func TestWriteCanonicalErrorPaths(t *testing.T) {
	tests := []struct {
		name string
		in   any
		path string
		want string
	}{
		{"unsupported at the root", 5, "", "kwargs cannot be hashed: kwargs: unsupported value"},
		{"unsupported float64 at the root", 1.5, "", "kwargs cannot be hashed: kwargs: unsupported value"},
		{"unsupported with a path", 5, "a.b", "kwargs cannot be hashed: a.b: unsupported value"},
		{"unsupported typed slice", []string{"x"}, "l", "kwargs cannot be hashed: l: unsupported value"},
		{"unsupported typed map", map[string]string{"k": "v"}, "m", "kwargs cannot be hashed: m: unsupported value"},
		{"unsupported inside a map", map[string]any{"a": 5}, "", "kwargs cannot be hashed: a: unsupported value"},
		{"unsupported inside nested maps", map[string]any{"a": map[string]any{"b": 5}}, "", "kwargs cannot be hashed: a.b: unsupported value"},
		{"unsupported inside a slice", []any{"x", 5}, "", "kwargs cannot be hashed: [1]: unsupported value"},
		{"unsupported inside a slice under a key", map[string]any{"a": []any{[]any{true, 5}}}, "", "kwargs cannot be hashed: a[0][1]: unsupported value"},
		{"unsupported mixed path", map[string]any{"a": []any{map[string]any{"b": 5}}}, "", "kwargs cannot be hashed: a[0].b: unsupported value"},
		{"path continues from the parent's", map[string]any{"b": 5}, "a", "kwargs cannot be hashed: a.b: unsupported value"},
		{"number error is passed through unchanged", json.Number("9007199254740993"), "a[2]", "kwargs cannot be hashed: a[2]: integer beyond 2^53"},
		{"number error at the root", json.Number("1e400"), "", "kwargs cannot be hashed: kwargs: number is not finite"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := writeCanonical(&buf, tc.in, tc.path)
			if err == nil {
				t.Fatalf("writeCanonical accepted the value and wrote %q", buf.String())
			}
			if err.Error() != tc.want {
				t.Errorf("error = %q, want %q", err.Error(), tc.want)
			}
		})
	}
}

func TestCanonicalNumber(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		// Integers, whatever the spelling.
		{"0", "0"},
		{"1", "1"},
		{"1.0", "1"},
		{"1.000", "1"},
		{"1e0", "1"},
		{"1E0", "1"},
		{"1e+0", "1"},
		{"0.1e1", "1"},
		{"10e-1", "1"},
		{"100", "100"},
		{"1e2", "100"},
		{"1.5e1", "15"},
		{"-5", "-5"},
		{"-5.0", "-5"},
		{"-1e2", "-100"},
		{"123456789012", "123456789012"},
		// Negative zero in any spelling.
		{"-0", "0"},
		{"-0.0", "0"},
		{"-0e0", "0"},
		{"-0e5", "0"},
		{"0e-5", "0"},
		// The bound.
		{"9007199254740992", "9007199254740992"},
		{"-9007199254740992", "-9007199254740992"},
		{"9007199254740992.0", "9007199254740992"},
		{"9.007199254740992e15", "9007199254740992"},
		// Not integers: the shortest digits that read back as the float64.
		{"0.5", "0.5"},
		{"-0.5", "-0.5"},
		{"0.1", "0.1"},
		{"0.25", "0.25"},
		{"2.5e-1", "0.25"},
		{"-1.25", "-1.25"},
		{"100.5", "100.5"},
		{"4.35", "4.35"},
		{"0.30000000000000004", "0.30000000000000004"},
		{"123456789.123456789", "123456789.12345679"},
		{"0.000001", "0.000001"},
		{"1e-6", "0.000001"},
		{"0.000001234", "0.000001234"},
		{"0.00000015", "1.5e-7"},
		{"1.5e-7", "1.5e-7"},
		{"1e-7", "1e-7"},
		{"0.0000001", "1e-7"},
		{"1.2345e-7", "1.2345e-7"},
		{"5e-324", "5e-324"},
		// A non-integer takes the integer of its nearest double.
		{"9007199254740993.5", "9007199254740994"},
		// A non-integer that rounds to zero.
		{"1e-400", "0"},
		{"-1e-400", "0"},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := canonicalNumber(json.Number(tc.in), "n")
			if err != nil {
				t.Fatalf("canonicalNumber(%s): %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("canonicalNumber(%s) = %s, want %s", tc.in, got, tc.want)
			}
		})
	}
}

func TestCanonicalNumberErrors(t *testing.T) {
	tests := []struct {
		in   string
		path string
		want string
	}{
		{"9007199254740993", "n", "kwargs cannot be hashed: n: integer beyond 2^53"},
		{"-9007199254740993", "n", "kwargs cannot be hashed: n: integer beyond 2^53"},
		{"9007199254740993.0", "n", "kwargs cannot be hashed: n: integer beyond 2^53"},
		{"1e30", "a.b", "kwargs cannot be hashed: a.b: integer beyond 2^53"},
		{"1e21", "a[0]", "kwargs cannot be hashed: a[0]: integer beyond 2^53"},
		{"100000000000000000000000000000", "n", "kwargs cannot be hashed: n: integer beyond 2^53"},
		{"1e400", "n", "kwargs cannot be hashed: n: number is not finite"},
		{"-1e400", "n", "kwargs cannot be hashed: n: number is not finite"},
		{"1e400", "", "kwargs cannot be hashed: kwargs: number is not finite"},
		{"abc", "n", "kwargs cannot be hashed: n: not a number"},
		{"", "n", "kwargs cannot be hashed: n: not a number"},
		{"1/2", "n", "kwargs cannot be hashed: n: not a number"},
		{"0x10", "n", "kwargs cannot be hashed: n: not a number"},
		{"1.2.3", "n", "kwargs cannot be hashed: n: not a number"},
		{"abc", "", "kwargs cannot be hashed: kwargs: not a number"},
		{"9007199254740993", "", "kwargs cannot be hashed: kwargs: integer beyond 2^53"},
	}
	for _, tc := range tests {
		t.Run(tc.in+"@"+tc.path, func(t *testing.T) {
			got, err := canonicalNumber(json.Number(tc.in), tc.path)
			if err == nil {
				t.Fatalf("canonicalNumber(%q) = %q, want an error", tc.in, got)
			}
			if got != "" {
				t.Errorf("canonicalNumber(%q) returned %q with an error, want an empty string", tc.in, got)
			}
			if err.Error() != tc.want {
				t.Errorf("error = %q, want %q", err.Error(), tc.want)
			}
		})
	}
}

func TestWriteCanonicalString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", `""`},
		{"plain", "abc", `"abc"`},
		{"double quote", `a"b`, `"a\"b"`},
		{"backslash", `a\b`, `"a\\b"`},
		{"quote after backslash", `\"`, `"\\\""`},
		{"backspace", "\b", `"\b"`},
		{"form feed", "\f", `"\f"`},
		{"line feed", "\n", `"\n"`},
		{"carriage return", "\r", `"\r"`},
		{"tab", "\t", `"\t"`},
		{"NUL", "\x00", `"\u0000"`},
		{"U+0001", "\x01", `"\u0001"`},
		{"vertical tab", "\v", `"\u000b"`},
		{"escape", "\x1b", `"\u001b"`},
		{"U+001f", "\x1f", `"\u001f"`},
		{"lower-case hex digits", "\x0e\x0f", `"\u000e\u000f"`},
		{"space is written as itself", " ", `" "`},
		{"slash is not escaped", "a/b", `"a/b"`},
		{"less than, greater than, ampersand", "<a>&", `"<a>&"`},
		{"DEL is written as itself", "\x7f", "\"\x7f\""},
		{"U+2028 is written as itself", " ", "\" \""},
		{"U+2029 is written as itself", " ", "\" \""},
		{"multibyte text", "héllo \U0001F600", "\"héllo \U0001F600\""},
		{"composed letter stays composed", "é", "\"é\""},
		{"decomposed letter stays decomposed", "é", "\"é\""},
		{"mixed", "a\tb\"c\\d\x02", `"a\tb\"c\\d\u0002"`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			writeCanonicalString(&buf, tc.in)
			if got := buf.String(); got != tc.want {
				t.Errorf("writeCanonicalString(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestWriteCanonicalStringAppends(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("k:")
	writeCanonicalString(&buf, "v")
	if got, want := buf.String(), `k:"v"`; got != want {
		t.Errorf("buffer = %q, want %q", got, want)
	}
}

func TestKwargsHash(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]any
		want string
	}{
		{"empty kwargs", map[string]any{}, hashOfEmptyObject},
		{"nil kwargs", nil, hashOfEmptyObject},
		{"one kwarg", map[string]any{"a": 1}, "sha256:015abd7f5cc57a2dd94b7590f04ad8084273905ee33ec5cebeae62276a97f862"},
		{"explicit null", map[string]any{"a": nil}, "sha256:d091f9c83c091f79652fe8786375b3fe4ce0861a56f5bfbafedbe431877ff0e8"},
		{"key order does not matter",
			map[string]any{"name": "read", "rule": "admin?", "controller": "users"},
			"sha256:498b9bece8c809c60f6d4422e9aa2559c0bfbb2409b335021638015aa3f04880"},
		{"hash is over the canonical bytes",
			map[string]any{"b": []any{1.0, "x"}, "a": map[string]any{"z": true, "y": nil}},
			sha256Of(`{"a":{"y":null,"z":true},"b":[1,"x"]}`)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := KwargsHash(tc.in)
			if err != nil {
				t.Fatalf("KwargsHash: %v", err)
			}
			if got != tc.want {
				t.Errorf("KwargsHash = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestKwargsHashAbsentAndNullDiffer(t *testing.T) {
	absent, err := KwargsHash(map[string]any{})
	if err != nil {
		t.Fatalf("KwargsHash: %v", err)
	}
	null, err := KwargsHash(map[string]any{"a": nil})
	if err != nil {
		t.Fatalf("KwargsHash: %v", err)
	}
	if absent == null {
		t.Errorf("an absent kwarg and an explicit null hash alike: %s", absent)
	}
}

func TestKwargsHashShape(t *testing.T) {
	got, err := KwargsHash(map[string]any{"a": "x"})
	if err != nil {
		t.Fatalf("KwargsHash: %v", err)
	}
	digits, ok := strings.CutPrefix(got, "sha256:")
	if !ok || len(digits) != 64 || digits != strings.ToLower(digits) {
		t.Errorf("KwargsHash = %q, want sha256: and 64 lower-case hex digits", got)
	}
	if _, err := hex.DecodeString(digits); err != nil {
		t.Errorf("KwargsHash digits are not hex: %v", err)
	}
}

// An error from CanonicalKwargs comes back unchanged, with an empty hash.
func TestKwargsHashErrorIsCanonicalKwargsError(t *testing.T) {
	in := map[string]any{"a": map[string]any{"n": json.Number("9007199254740993")}}
	_, want := CanonicalKwargs(in)
	if want == nil {
		t.Fatal("CanonicalKwargs accepted an integer beyond 2^53")
	}
	got, err := KwargsHash(in)
	if err == nil {
		t.Fatalf("KwargsHash accepted the argument and returned %q", got)
	}
	if got != "" {
		t.Errorf("KwargsHash returned %q with an error, want an empty string", got)
	}
	if err.Error() != want.Error() {
		t.Errorf("error = %q, want %q", err.Error(), want.Error())
	}

	if got, err := KwargsHash(map[string]any{"f": func() {}}); err == nil || got != "" ||
		!strings.HasPrefix(err.Error(), "kwargs cannot be hashed: ") {
		t.Errorf("KwargsHash of a function = %q, %v; want an empty hash and a kwargs cannot be hashed error", got, err)
	}
}

func TestBodyHash(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty text", "", hashOfEmptyText},
		{"LF text", "a\nb\n", sha256Of("a\nb\n")},
		{"CRLF is read as LF", "a\r\nb\r\n", sha256Of("a\nb\n")},
		{"mixed endings", "a\r\nb\nc\r\n", sha256Of("a\nb\nc\n")},
		{"lone CR is left as it is", "a\rb\n", sha256Of("a\rb\n")},
		{"trailing lone CR is left as it is", "a\r", sha256Of("a\r")},
		{"CR before CRLF is left, the CRLF is read as LF", "\r\r\n", sha256Of("\r\n")},
		{"CR then text then LF", "a\r\n\r\nb", sha256Of("a\n\nb")},
		{"indentation takes part", "  a\n", sha256Of("  a\n")},
		{"trailing spaces take part", "a \n", sha256Of("a \n")},
		{"trailing blank lines take part", "a\n\n", sha256Of("a\n\n")},
		{"no trailing newline", "a", sha256Of("a")},
		{"tabs are left alone", "\ta\t\n", sha256Of("\ta\t\n")},
		{"UTF-8 text", "héllo\n", sha256Of("héllo\n")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := BodyHash(tc.in); got != tc.want {
				t.Errorf("BodyHash(%q) = %s, want %s", tc.in, got, tc.want)
			}
		})
	}
}

func TestBodyHashLiteralValues(t *testing.T) {
	if got := BodyHash("x"); got != "sha256:2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881" {
		t.Errorf("BodyHash(x) = %s", got)
	}
	if got := BodyHash("policy :read, admin?\n"); got != "sha256:ca235cb34f9574af0f8880305104f597e91d812fca6b884d347fa0775a6f7076" {
		t.Errorf("BodyHash(policy) = %s", got)
	}
}

// Text that differs in anything the hash reads must hash differently.
func TestBodyHashDistinguishesWhatItReads(t *testing.T) {
	base := "def a\n  1\nend\n"
	seen := map[string]string{BodyHash(base): "base"}
	for name, text := range map[string]string{
		"reindented":           "def a\n    1\nend\n",
		"trailing space":       "def a \n  1\nend\n",
		"extra blank line":     base + "\n",
		"no final newline":     strings.TrimSuffix(base, "\n"),
		"lone carriage return": "def a\r  1\nend\n",
	} {
		h := BodyHash(text)
		if prev, dup := seen[h]; dup {
			t.Errorf("%s hashes like %s", name, prev)
		}
		seen[h] = name
	}
	if BodyHash("def a\r\n  1\r\nend\r\n") != BodyHash(base) {
		t.Error("the CRLF form of a body hashes differently from its LF form")
	}
}
