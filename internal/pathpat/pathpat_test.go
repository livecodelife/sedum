package pathpat

import "testing"

// The grammar two declarations now share: a record's scope entries and a
// package's unmanaged list. The cases below are the properties both rely on,
// and the reason there is one implementation rather than two.

func TestMatch(t *testing.T) {
	cases := []struct {
		entry, target string
		want          bool
	}{
		// A literal matches itself and nothing else, which is what lets one
		// function serve entries of either kind.
		{entry: "Gemfile", target: "Gemfile", want: true},
		{entry: "Gemfile", target: "Gemfile.lock", want: false},
		{entry: "Gemfile", target: "vendor/Gemfile", want: false},
		{entry: "db/schema.rb", target: "db/schema.rb", want: true},

		// * stays inside one segment. This is the whole difference between
		// config/*.yml and config/**.yml.
		{entry: "config/*.yml", target: "config/database.yml", want: true},
		{entry: "config/*.yml", target: "config/deploy/production.yml", want: false},
		{entry: "config/**.yml", target: "config/deploy/production.yml", want: false},
		{entry: "config/**/*.yml", target: "config/deploy/production.yml", want: true},

		// ** consumes any number of segments including none.
		{entry: "app/**", target: "app/models/todo.rb", want: true},
		{entry: "app/**", target: "app", want: true},
		{entry: "**/secrets.yml", target: "secrets.yml", want: true},
		{entry: "**/secrets.yml", target: "config/deploy/secrets.yml", want: true},

		// A trailing slash is a subtree, since no file can be named that.
		{entry: "config/credentials/", target: "config/credentials/production.key", want: true},
		{entry: "config/credentials/", target: "config/credentials.yml", want: false},

		{entry: "db/migrate/?.rb", target: "db/migrate/1.rb", want: true},
		{entry: "db/migrate/?.rb", target: "db/migrate/12.rb", want: false},
		{entry: "log/[ab].log", target: "log/a.log", want: true},
		{entry: "log/[ab].log", target: "log/c.log", want: false},

		// Both sides are normalized, so a declaration authored on one
		// platform reads the same on another.
		{entry: "./config/routes.rb", target: "config/routes.rb", want: true},
	}

	for _, tc := range cases {
		if got := Match(tc.entry, tc.target); got != tc.want {
			t.Errorf("Match(%q, %q) = %v, want %v", tc.entry, tc.target, got, tc.want)
		}
	}
}

// A pattern names a set; a literal names one file. Telling them apart is what
// decides whether a scope entry creates a file (prov-2026-e8671c88), so the
// distinction has to hold for entries that carry no metacharacter at all.
func TestIsPattern(t *testing.T) {
	cases := map[string]bool{
		"Gemfile":              false,
		"app/models/todo.rb":   false,
		"app/**":               true,
		"config/*.yml":         true,
		"db/migrate/?.rb":      true,
		"log/[ab].log":         true,
		"config/credentials/":  true,
		"config/master.key":    false,
		"vendor/bundle/ruby/":  true,
		"app/models/{name}.rb": false,
	}

	for entry, want := range cases {
		if got := IsPattern(entry); got != want {
			t.Errorf("IsPattern(%q) = %v, want %v", entry, got, want)
		}
	}
}

// A pattern that cannot be read would match nothing at all, so an entry meant
// to keep Sedum out of a directory would keep it out of nothing. Both readers
// reject it rather than carrying it.
func TestCheckRejectsUnreadablePatterns(t *testing.T) {
	if err := Check("config/[unclosed.yml"); err == nil {
		t.Error("an unclosed character class was accepted; left alone it matches nothing")
	}
	for _, entry := range []string{"Gemfile", "app/**", "config/*.yml", "log/[ab].log", "config/credentials/"} {
		if err := Check(entry); err != nil {
			t.Errorf("Check(%q) rejected a readable entry: %v", entry, err)
		}
	}
}

// MatchAny returns the entry that matched so a diagnostic can name the
// declaration rather than only the path it stopped.
func TestMatchAnyNamesTheEntry(t *testing.T) {
	entries := []string{"Gemfile", "config/credentials/", "db/schema.rb"}

	entry, ok := MatchAny(entries, "config/credentials/production.key")
	if !ok {
		t.Fatal("a path under a declared subtree did not match")
	}
	if entry != "config/credentials/" {
		t.Errorf("matched entry = %q, want the subtree that covered it", entry)
	}

	if _, ok := MatchAny(entries, "app/models/todo.rb"); ok {
		t.Error("a path no entry covers matched")
	}
	if _, ok := MatchAny(nil, "Gemfile"); ok {
		t.Error("an empty declaration matched something")
	}
}
