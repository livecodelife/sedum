package transform

import "testing"

// trimprefix removes a literal leader from the start of a value and unparen
// removes one matching outer pair of parentheses, so a model that writes the
// delimiters its template already supplies still renders valid code.
func TestTrimprefixRemovesALeadingLiteral(t *testing.T) {
	e := newEngine(t, Config{})

	cases := []struct{ in, want string }{
		{"// Lets the slash key jump", "Lets the slash key jump"},
		{"//Lets the slash key jump", "Lets the slash key jump"},
		{"//\tTabbed", "Tabbed"},
		{"Lets the slash key jump", "Lets the slash key jump"},
		{"a // not at the start", "a // not at the start"},
		{"// // twice", "// twice"},
		{"//", ""},
		{"", ""},
	}
	for _, c := range cases {
		if got := apply(t, e, "trimprefix://", c.in); got != c.want {
			t.Errorf("trimprefix://(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestUnparenRemovesOneMatchingOuterPair(t *testing.T) {
	e := newEngine(t, Config{})

	cases := []struct{ in, want string }{
		{"()", ""},
		{"(a: string, b: number)", "a: string, b: number"},
		{"a: string, b: number", "a: string, b: number"},
		{"", ""},
		{"(a)(b)", "(a)(b)"},
		{"(a: () => void)", "a: () => void"},
		{"((a))", "(a)"},
		{"(a", "(a"},
		{"a)", "a)"},
		{"(", "("},
	}
	for _, c := range cases {
		if got := apply(t, e, "unparen", c.in); got != c.want {
			t.Errorf("unparen(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestTrimOperationsAreBuiltinWithTheRightArguments(t *testing.T) {
	e := newEngine(t, Config{})

	if !Parameterized("trimprefix") {
		t.Error("trimprefix is not parameterized")
	}
	if err := e.Check(ParseRef("trimprefix")); err == nil {
		t.Error("trimprefix without an argument was accepted")
	}
	if err := e.Check(ParseRef("trimprefix:{{x}}")); err == nil {
		t.Error("trimprefix with a dynamic argument was accepted")
	}
	if err := e.Check(ParseRef("trimprefix://")); err != nil {
		t.Errorf("trimprefix:// rejected: %v", err)
	}

	if Parameterized("unparen") {
		t.Error("unparen is parameterized")
	}
	if err := e.Check(ParseRef("unparen")); err != nil {
		t.Errorf("unparen rejected: %v", err)
	}
	if err := e.Check(ParseRef("unparen:()")); err == nil {
		t.Error("unparen with an argument was accepted")
	}
}
