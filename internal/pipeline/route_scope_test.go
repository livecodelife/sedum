package pipeline

import (
	"context"
	"strings"
	"testing"
)

// A route folder written with brackets ([id], [...slug], [[...sign-in]]) is a
// path Sedum plans, so a check reads the regions in its file: a stale one makes
// the record stale (prov-2026-e8671c88).

const staleRegion = "# sedum:addField {\"tier\":\"owned\",\"record\":\"prov-2026-aaaaaaaa\"}\n# /sedum:addField\n"

func routeConfig(t *testing.T, scope string) Config {
	t.Helper()
	cfg := config(t)
	cfg.Records = recordsDir(t, map[string]string{skipRecordA: scope})
	cfg.Check = true
	return cfg
}

func TestACheckReadsTheStaleRegionInARouteFolderFile(t *testing.T) {
	for _, route := range []string{
		"app/[id]/thing.rb",
		"app/[...slug]/thing.rb",
		"app/[[...sign-in]]/thing.rb",
	} {
		t.Run(route, func(t *testing.T) {
			cfg := routeConfig(t, "  - "+route+"\n  - app/models/user.rb\n")
			writeOutput(t, cfg, route, staleRegion)

			result, err := Run(context.Background(), cfg)
			if err != nil {
				t.Fatalf("check: %v", err)
			}

			if got := modeOf(t, result, skipRecordA); got != ModeReselect {
				t.Errorf("mode = %v, want %v: the region in %s is stale", got, ModeReselect, route)
			}
		})
	}
}

// A record whose only entry is a route entry used to name no path, warn so, and
// be reported as writing nothing.
func TestARecordOfOneRouteEntryIsNotReportedAsWritingNothing(t *testing.T) {
	cfg := routeConfig(t, "  - app/[id]/thing.rb\n")
	writeOutput(t, cfg, "app/[id]/thing.rb", staleRegion)

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	for _, w := range result.Warnings {
		if strings.Contains(w, "names no path") {
			t.Errorf("warning %q: the record names a file", w)
		}
	}
	if got := modeOf(t, result, skipRecordA); got != ModeReselect {
		t.Errorf("mode = %v, want %v", got, ModeReselect)
	}
}

// Any other bracket is still a character class: it names no file, so nothing is
// planned for it.
func TestACharacterClassEntryStillNamesNoFile(t *testing.T) {
	cfg := routeConfig(t, "  - app/models/[ab].rb\n")
	writeOutput(t, cfg, "app/models/a.rb", staleRegion)

	result, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if got := modeOf(t, result, skipRecordA); got != ModeNone {
		t.Errorf("mode = %v, want %v: a character class creates and checks no file", got, ModeNone)
	}
}
