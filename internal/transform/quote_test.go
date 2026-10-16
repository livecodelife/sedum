package transform

import (
	"encoding/json"
	"strconv"
	"testing"
)

// quote renders a value as a complete double-quoted string literal, so model
// output bound into a string cannot break the generated file
// (prov-2026-f63c340a).
func TestQuoteRendersAStringLiteral(t *testing.T) {
	e := newEngine(t, Config{})

	cases := []struct{ in, want string }{
		{"Create a record", `"Create a record"`},
		{`say "hi"`, `"say \"hi\""`},
		{`C:\path`, `"C:\\path"`},
		{"two\nlines\tand a tab", `"two\nlines\tand a tab"`},
		{"<b>&</b>", `"<b>&</b>"`},
		{"café 日本", `"café 日本"`},
		{"", `""`},
	}
	for _, c := range cases {
		if got := apply(t, e, "quote", c.in); got != c.want {
			t.Errorf("quote(%q) = %s, want %s", c.in, got, c.want)
		}
	}
}

// Whatever the model binds, the result reads back as exactly that value - as
// a JSON string and as a Go string literal, the two readings the targets that
// use it share. The check is in a test; the operation knows no language.
func TestQuoteRoundTripsAnyInput(t *testing.T) {
	e := newEngine(t, Config{})
	inputs := []string{
		"plain", `"`, `\`, `\"`, "\x00\x01\x1f\x7f", "\u2028\u2029", "tab\there", "\r\n",
		"emoji 🌱", "back`tick", "it's", "{{not a template}}", "%s %d", "</script>",
	}
	for _, in := range inputs {
		out := apply(t, e, "quote", in)

		var fromJSON string
		if err := json.Unmarshal([]byte(out), &fromJSON); err != nil || fromJSON != in {
			t.Errorf("quote(%q) = %s does not read back as JSON: %q, %v", in, out, fromJSON, err)
		}
		fromGo, err := strconv.Unquote(out)
		if err != nil || fromGo != in {
			t.Errorf("quote(%q) = %s does not read back as a Go literal: %q, %v", in, out, fromGo, err)
		}
	}
}

// A number or a bool is quoted as its text; a list is refused, as it is for
// every operation but join.
func TestQuoteTakesScalarsOnly(t *testing.T) {
	e := newEngine(t, Config{})
	if out, err := e.Apply(ParseRef("quote"), 42); err != nil || out != `"42"` {
		t.Errorf("quote(42) = %s, %v; want \"42\"", out, err)
	}
	if _, err := e.Apply(ParseRef("quote"), []any{"a", "b"}); err == nil {
		t.Error("quote accepted a list")
	}
}
