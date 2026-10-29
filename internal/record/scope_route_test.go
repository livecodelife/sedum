package record

import "testing"

// A Next.js route folder is written with brackets: [id], [...slug],
// [[...sign-in]]. Such an entry names a file, so it is a path Sedum plans, not
// a pattern; any other bracket is still a character class
// (prov-2026-e8671c88).

func TestRouteEntriesAreLiteralPaths(t *testing.T) {
	r := loadOne(t, `id: prov-2026-aaaaaaaa
affected_scope:
  - src/app/api/apps/[id]/route.ts
  - src/app/api/[...slug]/route.ts
  - src/app/[[...sign-in]]/page.ts
  - src/app/[user_id]/[post-1]/page.ts
  - src/plain/route.ts
`)

	want := []string{
		"src/app/[[...sign-in]]/page.ts",
		"src/app/[user_id]/[post-1]/page.ts",
		"src/app/api/[...slug]/route.ts",
		"src/app/api/apps/[id]/route.ts",
		"src/plain/route.ts",
	}
	got := append([]string(nil), r.Paths...)
	sortStrings(got)
	if !equal(got, want) {
		t.Errorf("Paths = %v, want %v", got, want)
	}
	if len(r.Patterns) != 0 {
		t.Errorf("Patterns = %v, want none: every entry names a file", r.Patterns)
	}
	for _, p := range want {
		if !r.Authorizes(p) {
			t.Errorf("the record does not authorize %s", p)
		}
	}
}

func TestOtherBracketsStayPatterns(t *testing.T) {
	r := loadOne(t, `id: prov-2026-aaaaaaaa
affected_scope:
  - src/[ab].ts
  - src/file[0-9].ts
  - src/app/[id!]/route.ts
  - src/app/[id]x/route.ts
  - src/app/x[id]/route.ts
`)

	if len(r.Paths) != 0 {
		t.Errorf("Paths = %v, want none: each of these is a character class", r.Paths)
	}
	if len(r.Patterns) != 5 {
		t.Errorf("Patterns = %v, want all five entries", r.Patterns)
	}
}

func TestARouteSegmentWithAWildcardOrTrailingSlashIsAPattern(t *testing.T) {
	r := loadOne(t, `id: prov-2026-aaaaaaaa
affected_scope:
  - src/app/[id]/*.ts
  - src/app/[id]/route?.ts
  - src/app/[id]/**
  - src/app/[id]/
`)

	if len(r.Paths) != 0 {
		t.Errorf("Paths = %v, want none", r.Paths)
	}
	if len(r.Patterns) != 4 {
		t.Errorf("Patterns = %v, want all four entries", r.Patterns)
	}
}

// A route entry still matches only itself: it is a literal path, not a class
// that also reads as src/app/i/route.ts.
func TestARouteEntryAuthorizesItsOwnPathOnly(t *testing.T) {
	r := loadOne(t, `id: prov-2026-aaaaaaaa
affected_scope:
  - src/app/[id]/route.ts
`)

	if !r.Authorizes("src/app/[id]/route.ts") {
		t.Error("the route entry does not authorize the path it names")
	}
	for _, other := range []string{"src/app/i/route.ts", "src/app/d/route.ts", "src/app/[slug]/route.ts"} {
		if r.Authorizes(other) {
			t.Errorf("the route entry also authorizes %s", other)
		}
	}
}

func TestARecordOfRouteEntriesDoesNotWarnItNamesNoPath(t *testing.T) {
	_, warnings, err := Load(writeRecords(t, map[string]string{"r.yml": `id: prov-2026-aaaaaaaa
affected_scope:
  - src/app/api/apps/[id]/route.ts
`}), Options{})

	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if containsFragment(warnings, "no path") {
		t.Errorf("warnings = %v: the record names a file", warnings)
	}
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
