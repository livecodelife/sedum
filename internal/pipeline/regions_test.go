package pipeline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// sedum regions lists every region naming a record, from markers alone, for
// provenance's evidence code hash (prov-2026-bf01cd6b).

func TestRegionsListsEveryOwnedRegionFromMarkers(t *testing.T) {
	cfg := twoOwners(t)

	report, err := Regions(RegionsConfig{Generators: cfg.Generators, Output: cfg.Output})
	if err != nil {
		t.Fatalf("Regions: %v", err)
	}
	if report.Format != "sedum.regions/1" {
		t.Errorf("format = %q", report.Format)
	}
	if len(report.Regions) != 2 {
		t.Fatalf("regions = %+v, want two", report.Regions)
	}
	byRecord := map[string]RegionEntry{}
	for _, r := range report.Regions {
		byRecord[r.Record] = r
	}
	a := byRecord["prov-2026-aaaaaaaa"]
	if a.File != "app/models/user.rb" || a.Action != "addField" || a.Tier != "owned" {
		t.Errorf("record a's region = %+v", a)
	}
	if a.Identity["field"] != "email" || a.Identity["model"] != "user" {
		t.Errorf("identity = %v, want the selecting kwargs", a.Identity)
	}
	if !strings.HasPrefix(a.Gen, "sha256:") || !strings.HasPrefix(a.Pkg, "sha256:") {
		t.Errorf("gen/pkg = %q/%q, want both hashes", a.Gen, a.Pkg)
	}

	// --record narrows the list.
	only, err := Regions(RegionsConfig{Generators: cfg.Generators, Output: cfg.Output, Only: []string{"prov-2026-bbbbbbbb"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(only.Regions) != 1 || only.Regions[0].Record != "prov-2026-bbbbbbbb" {
		t.Errorf("--record b = %+v", only.Regions)
	}
}

// A marker without hashes is reported with empty gen and pkg, never filled in.
func TestRegionsReportsMissingHashesAsEmpty(t *testing.T) {
	cfg := grown(t)
	stripHashes(t, cfg)

	report, err := Regions(RegionsConfig{Generators: cfg.Generators, Output: cfg.Output})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Regions) != 1 || report.Regions[0].Gen != "" || report.Regions[0].Pkg != "" {
		t.Errorf("regions = %+v, want one with empty gen and pkg", report.Regions)
	}
}

// With --records, a marker naming an alias is reported under the surviving id,
// and --record finds it by either.
func TestRegionsResolveAliasesWithRecords(t *testing.T) {
	cfg := grown(t)
	path := filepath.Join(cfg.Output, "app", "models", "user.rb")
	aliased := strings.Replace(userModel(t, cfg), "prov-2026-aaaaaaaa", "prov-2026-01d01d01", 1)
	if err := os.WriteFile(path, []byte(aliased), 0o644); err != nil {
		t.Fatal(err)
	}
	editRecordWith(t, cfg, "aliases: [prov-2026-01d01d01]\n")

	for _, only := range [][]string{nil, {"prov-2026-aaaaaaaa"}, {"prov-2026-01d01d01"}} {
		report, err := Regions(RegionsConfig{Generators: cfg.Generators, Output: cfg.Output, Records: cfg.Records, Only: only})
		if err != nil {
			t.Fatal(err)
		}
		if len(report.Regions) != 1 || report.Regions[0].Record != "prov-2026-aaaaaaaa" {
			t.Errorf("--record %v: regions = %+v, want one under the surviving id", only, report.Regions)
		}
	}
}

// The generators and records directories are never read.
func TestRegionsSkipsTheGeneratorsDirectory(t *testing.T) {
	cfg := config(t)
	inside := filepath.Join(cfg.Output, "generators")
	if err := os.CopyFS(inside, os.DirFS(cfg.Generators)); err != nil {
		t.Fatal(err)
	}
	fixture := filepath.Join(inside, "rails", "testdata", "user.rb")
	if err := os.MkdirAll(filepath.Dir(fixture), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fixture, []byte("# sedum:addField {\"tier\":\"owned\",\"record\":\"prov-2026-zzzzzzzz\"}\n# /sedum:addField\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := Regions(RegionsConfig{Generators: inside, Output: cfg.Output})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Regions) != 0 {
		t.Errorf("regions = %+v, want none from the generators directory", report.Regions)
	}
}

func editRecordWith(t *testing.T, cfg Config, extra string) {
	t.Helper()
	body := "id: prov-2026-aaaaaaaa\n" + extra + "intent: |\n  Something.\naffected_scope:\n  - app/models/user.rb\n"
	if err := os.WriteFile(filepath.Join(cfg.Records, "prov-2026-aaaaaaaa.yml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
