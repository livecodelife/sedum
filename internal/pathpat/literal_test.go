package pathpat

import "testing"

// A segment that is written exactly as the path's segment matches it, so a
// path that carries bracket characters (a Next.js route folder such as [id] or
// [...slug]) is covered by the entry that names it, while a bracket segment in
// a pattern still reads as a character class against any other segment.

func TestMatchSegmentEqualToThePathSegmentMatchesLiterally(t *testing.T) {
	cases := []struct {
		entry, target string
		want          bool
	}{
		{"src/app/api/apps/[id]/build/route.ts", "src/app/api/apps/[id]/build/route.ts", true},
		{"src/app/[...notfound]/page.tsx", "src/app/[...notfound]/page.tsx", true},
		{"src/app/[id]/page.tsx", "src/app/[id]/page.tsx", true},
		{"src/app/[[...slug]]/page.tsx", "src/app/[[...slug]]/page.tsx", true},
		{"src/app/**/page.tsx", "src/app/[...notfound]/page.tsx", true},
		{"src/app/api/apps/[id]/**", "src/app/api/apps/[id]/build/route.ts", true},
		{"src/app/[id]/*.ts", "src/app/[id]/route.ts", true},

		// A different segment is still a different path.
		{"src/app/[id]/page.tsx", "src/app/[slug]/page.tsx", false},
		{"src/app/[id]/page.tsx", "src/app/x/page.tsx", false},
		{"src/app/[...notfound]/page.tsx", "src/app/[...other]/page.tsx", false},
		{"src/app/[id]/page.tsx", "src/app/[id]/layout.tsx", false},

		// A bracket segment in an entry is still a character class against a
		// segment it does not equal.
		{"src/app/[id]/page.tsx", "src/app/i/page.tsx", true},
		{"src/app/[id]/page.tsx", "src/app/d/page.tsx", true},
		{"src/app/[a-c]/page.tsx", "src/app/b/page.tsx", true},
		{"src/app/[a-c]/page.tsx", "src/app/z/page.tsx", false},
	}
	for _, c := range cases {
		if got := Match(c.entry, c.target); got != c.want {
			t.Errorf("Match(%q, %q) = %v, want %v", c.entry, c.target, got, c.want)
		}
	}
}

func TestMatchAnyNamesTheEntryThatCoversABracketPath(t *testing.T) {
	entries := []string{"docs/**", "src/app/api/apps/[id]/build/route.ts"}
	got, ok := MatchAny(entries, "src/app/api/apps/[id]/build/route.ts")
	if !ok || got != "src/app/api/apps/[id]/build/route.ts" {
		t.Errorf("MatchAny = %q, %v; want the literal entry", got, ok)
	}
}

func TestCheckStillRejectsAnUnclosedClassAndAcceptsBracketRouteFolders(t *testing.T) {
	for _, entry := range []string{"src/app/[id]/page.tsx", "src/app/[...notfound]/page.tsx", "src/app/[[...slug]]/page.tsx"} {
		if err := Check(entry); err != nil {
			t.Errorf("Check(%q) = %v, want nil", entry, err)
		}
	}
	if err := Check("src/app/[id/page.tsx"); err == nil {
		t.Error("Check accepted an unclosed character class")
	}
}
