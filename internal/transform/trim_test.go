package transform

import "testing"

// trimprefix removes a literal leader from the start of a value and unwrap
// removes one matching outer pair, so a model that writes the delimiters its
// template already supplies still renders valid code.
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

func TestUnwrapRemovesOneMatchingOuterPair(t *testing.T) {
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
		if got := apply(t, e, "unwrap:()", c.in); got != c.want {
			t.Errorf("unwrap:()(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestUnwrapAcceptsOtherDelimiterPairs(t *testing.T) {
	e := newEngine(t, Config{})

	if got := apply(t, e, "unwrap:<>", "<T, U>"); got != "T, U" {
		t.Errorf("unwrap:<>(%q) = %q, want %q", "<T, U>", got, "T, U")
	}
	if got := apply(t, e, "unwrap:{}", "{ a, b }"); got != " a, b " {
		t.Errorf("unwrap:{}(%q) = %q, want %q", "{ a, b }", got, " a, b ")
	}
}

func TestTrimOperationsAreBuiltinAndTakeLiteralArguments(t *testing.T) {
	e := newEngine(t, Config{})

	for _, op := range []string{"trimprefix", "unwrap"} {
		if !Parameterized(op) {
			t.Errorf("%s is not parameterized", op)
		}
		if err := e.Check(ParseRef(op)); err == nil {
			t.Errorf("%s without an argument was accepted", op)
		}
		if err := e.Check(ParseRef(op + ":{{x}}")); err == nil {
			t.Errorf("%s with a dynamic argument was accepted", op)
		}
	}
	if err := e.Check(ParseRef("trimprefix://")); err != nil {
		t.Errorf("trimprefix:// rejected: %v", err)
	}
	if err := e.Check(ParseRef("unwrap:()")); err != nil {
		t.Errorf("unwrap:() rejected: %v", err)
	}
}

func TestUnwrapRefusesAnArgumentThatIsNotAPair(t *testing.T) {
	e := newEngine(t, Config{})

	for _, arg := range []string{"(", "(()", "abc"} {
		if _, err := e.Apply(ParseRef("unwrap:"+arg), "(x)"); err == nil {
			t.Errorf("unwrap:%s was accepted, want an error naming the two characters it needs", arg)
		}
	}
}
