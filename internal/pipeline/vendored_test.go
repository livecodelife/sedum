package pipeline

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// The marker walk never enters vendored or build-output directories, and a
// project can skip more with an exclude list (prov-2026-fc05416c).

const vendoredMarker = "# sedum:addField {\"tier\":\"owned\",\"record\":\"prov-2026-aaaaaaaa\"}\n# /sedum:addField\n"

// putMarked writes a Ruby file carrying one marked region at rel under output.
func putMarked(t *testing.T, output, rel string) {
	t.Helper()
	path := filepath.Join(output, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(vendoredMarker), 0o644); err != nil {
		t.Fatal(err)
	}
}

func putConfigYAML(t *testing.T, output, body string) {
	t.Helper()
	path := filepath.Join(output, ".sedum", "config.yml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// regionFiles is the files Regions reports, sorted.
func regionFiles(t *testing.T, cfg RegionsConfig) []string {
	t.Helper()
	report, err := Regions(cfg)
	if err != nil {
		t.Fatalf("Regions: %v", err)
	}
	var files []string
	for _, r := range report.Regions {
		files = append(files, r.File)
	}
	sort.Strings(files)
	return files
}

func sameFiles(t *testing.T, got []string, want ...string) {
	t.Helper()
	sort.Strings(want)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("files = %v, want %v", got, want)
	}
}

func TestRegionsSkipsVendoredDirectoriesAtAnyDepth(t *testing.T) {
	gens := config(t).Generators
	output := t.TempDir()
	putMarked(t, output, "app/models/keep.rb")
	for _, name := range []string{"node_modules", ".next", "vendor", "dist", "target"} {
		putMarked(t, output, name+"/top.rb")
		putMarked(t, output, "app/"+name+"/nested.rb")
		putMarked(t, output, "a/b/"+name+"/c/deep.rb")
	}
	// A name that only contains a vendored name is an ordinary directory.
	putMarked(t, output, "app/vendored/kept_too.rb")

	sameFiles(t, regionFiles(t, RegionsConfig{Generators: gens, Output: output}),
		"app/models/keep.rb", "app/vendored/kept_too.rb")
}

// The output directory itself is never skipped, so an --output that is or sits
// inside a vendored directory is still read.
func TestRegionsReadsAnOutputInsideAVendoredDirectory(t *testing.T) {
	gens := config(t).Generators
	output := filepath.Join(t.TempDir(), "node_modules", "pkg")
	putMarked(t, output, "lib/a.rb")
	putMarked(t, output, "node_modules/inner.rb")

	sameFiles(t, regionFiles(t, RegionsConfig{Generators: gens, Output: output}), "lib/a.rb")
}

func TestConfigExcludeNamesAndPaths(t *testing.T) {
	gens := config(t).Generators
	output := t.TempDir()
	putMarked(t, output, "app/models/keep.rb")
	putMarked(t, output, "generated/a.rb")
	putMarked(t, output, "pkg/generated/b.rb")
	putMarked(t, output, "app/legacy/c.rb")
	putMarked(t, output, "other/app/legacy/d.rb")
	putConfigYAML(t, output, "exclude:\n  - generated\n  - app/legacy/\n")

	// A bare name matches at any depth; an entry with a slash names exactly
	// that directory under the output, and a trailing slash is ignored.
	sameFiles(t, regionFiles(t, RegionsConfig{Generators: gens, Output: output}),
		"app/models/keep.rb", "other/app/legacy/d.rb")
}

func TestExtraExcludeEntriesAreAddedToTheConfigs(t *testing.T) {
	gens := config(t).Generators
	output := t.TempDir()
	putMarked(t, output, "app/models/keep.rb")
	putMarked(t, output, "generated/a.rb")
	putMarked(t, output, "scratch/b.rb")
	putConfigYAML(t, output, "exclude: [generated]\n")

	sameFiles(t, regionFiles(t, RegionsConfig{Generators: gens, Output: output, Exclude: []string{"scratch"}}),
		"app/models/keep.rb")
}

// A config without an exclude list, or no config at all, changes nothing.
func TestNoExcludeConfigReadsEverythingButVendored(t *testing.T) {
	gens := config(t).Generators
	output := t.TempDir()
	putMarked(t, output, "generated/a.rb")
	sameFiles(t, regionFiles(t, RegionsConfig{Generators: gens, Output: output}), "generated/a.rb")

	putConfigYAML(t, output, "verify:\n  - true\n")
	sameFiles(t, regionFiles(t, RegionsConfig{Generators: gens, Output: output}), "generated/a.rb")
}

func TestInvalidExcludeEntriesAreRejectedBeforeAnyFileIsRead(t *testing.T) {
	gens := config(t).Generators
	for name, entry := range map[string]string{
		"empty":    "",
		"absolute": "/etc",
		"parent":   "../x",
		"embedded": "a/../b",
	} {
		t.Run(name, func(t *testing.T) {
			output := t.TempDir()
			putMarked(t, output, "app/a.rb")

			_, err := Regions(RegionsConfig{Generators: gens, Output: output, Exclude: []string{entry}})
			if err == nil {
				t.Fatalf("entry %q was accepted", entry)
			}
			if !strings.Contains(err.Error(), "exclude") || (entry != "" && !strings.Contains(err.Error(), entry)) {
				t.Errorf("error = %q, want it to name the exclude entry %q", err, entry)
			}

			// The same entry in the config file is the same error.
			putConfigYAML(t, output, "exclude:\n  - \""+entry+"\"\n")
			if _, err := Regions(RegionsConfig{Generators: gens, Output: output}); err == nil {
				t.Errorf("config entry %q was accepted", entry)
			}
		})
	}
}

// grow's orphan scan reads the same walk, so a marked file under a vendored
// directory is neither reported as an orphan nor pruned.
func TestGrowOrphanScanSkipsVendoredAndExcluded(t *testing.T) {
	orphanAt := func(t *testing.T, rel string, extra []string, config string) *Result {
		t.Helper()
		cfg := config0(t)
		if config != "" {
			putConfigYAML(t, cfg.Output, config)
		}
		path := filepath.Join(cfg.Output, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		gone := strings.Replace(vendoredMarker, "prov-2026-aaaaaaaa", "prov-2026-dddddddd", 1)
		if err := os.WriteFile(path, []byte(gone), 0o644); err != nil {
			t.Fatal(err)
		}
		cfg.Check = true
		cfg.Exclude = extra
		result, err := Run(context.Background(), cfg)
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		return result
	}

	if got := orphanAt(t, "lib/old.rb", nil, "").Orphans; len(got) != 1 {
		t.Fatalf("control: orphans = %+v, want the one in lib/", got)
	}
	if got := orphanAt(t, "node_modules/pkg/old.rb", nil, "").Orphans; len(got) != 0 {
		t.Errorf("orphans under node_modules = %+v, want none", got)
	}
	if got := orphanAt(t, "lib/old.rb", nil, "exclude: [lib]\n").Orphans; len(got) != 0 {
		t.Errorf("orphans under a config-excluded lib = %+v, want none", got)
	}
	if got := orphanAt(t, "lib/old.rb", []string{"lib"}, "").Orphans; len(got) != 0 {
		t.Errorf("orphans under an Exclude-field lib = %+v, want none", got)
	}
}

// config0 is a Config whose one record authorizes every .rb file, so an
// orphan in any directory counts as being in an authorized file.
func config0(t *testing.T) Config {
	t.Helper()
	cfg := config(t)
	cfg.Records = writeTree(t, map[string]string{
		"prov-2026-aaaaaaaa.yml": recordFile("prov-2026-aaaaaaaa", "  - app/models/user.rb\n  - lib/**\n  - node_modules/**\n", ""),
	})
	return cfg
}

func TestWhitespaceAroundAnEntryIsIgnored(t *testing.T) {
	gens := config(t).Generators
	output := t.TempDir()
	putMarked(t, output, "app/keep.rb")
	putMarked(t, output, "generated/a.rb")
	putConfigYAML(t, output, "exclude:\n  - \"  generated  \"\n")

	sameFiles(t, regionFiles(t, RegionsConfig{Generators: gens, Output: output}), "app/keep.rb")
}

// Entries are names and paths: a glob character is an ordinary character.
func TestExcludeEntriesAreNotGlobs(t *testing.T) {
	gens := config(t).Generators
	output := t.TempDir()
	putMarked(t, output, "generated/a.rb")
	putMarked(t, output, "gen*/b.rb")

	sameFiles(t, regionFiles(t, RegionsConfig{Generators: gens, Output: output, Exclude: []string{"gen*"}}),
		"generated/a.rb")
}

func TestAMalformedConfigFailsNamingTheFile(t *testing.T) {
	gens := config(t).Generators
	for name, body := range map[string]string{
		"not yaml":      "exclude: [unclosed\n",
		"not a list":    "exclude:\n  a: b\n",
		"not a strings": "exclude:\n  - [x]\n",
	} {
		t.Run(name, func(t *testing.T) {
			output := t.TempDir()
			putConfigYAML(t, output, body)
			_, err := Regions(RegionsConfig{Generators: gens, Output: output})
			if err == nil || !strings.Contains(err.Error(), ".sedum/config.yml") {
				t.Errorf("error = %v, want one naming .sedum/config.yml", err)
			}
		})
	}
}

// grow rejects an unusable entry before it reads a generator package or
// writes a file, so a refused run leaves the output tree as it was.
func TestGrowRejectsAnUnusableExcludeEntryBeforeWritingAnything(t *testing.T) {
	for name, entry := range map[string]string{"absolute": "/etc", "parent": "../x", "empty": ""} {
		t.Run(name, func(t *testing.T) {
			cfg := config0(t)
			putConfigYAML(t, cfg.Output, "exclude:\n  - \""+entry+"\"\n")
			putMarked(t, cfg.Output, "app/existing.rb")
			before := treeContents(t, cfg.Output)

			_, err := Run(context.Background(), cfg)
			if err == nil || !strings.Contains(err.Error(), "exclude") {
				t.Fatalf("error = %v, want an exclude entry error", err)
			}
			if after := treeContents(t, cfg.Output); !reflect.DeepEqual(before, after) {
				t.Errorf("tree changed: before %v, after %v", before, after)
			}
		})
	}
}

// The per-record scans read a file a record names wherever it sits, vendored
// or excluded directories included. One record per directory kind, so finding
// either file alone cannot satisfy both cases.
func TestRecordScansReadANamedFileInASkippedDirectory(t *testing.T) {
	for name, tc := range map[string]struct{ file, config string }{
		"vendored": {"vendor/x.rb", ""},
		"excluded": {"gen/y.rb", "exclude: [gen]\n"},
	} {
		t.Run(name, func(t *testing.T) {
			cfg := config(t)
			cfg.Records = writeTree(t, map[string]string{
				"prov-2026-aaaaaaaa.yml": recordFile("prov-2026-aaaaaaaa", "  - "+tc.file+"\n", ""),
			})
			putMarked(t, cfg.Output, tc.file)
			if tc.config != "" {
				putConfigYAML(t, cfg.Output, tc.config)
			}
			cfg.Check = true

			result, err := Run(context.Background(), cfg)
			if err != nil {
				t.Fatalf("Run: %v", err)
			}
			if len(result.Classifications) != 1 || result.Classifications[0].Mode != ModeReselect {
				t.Errorf("classifications = %+v, want the record found by its owned region (re-select, not select)", result.Classifications)
			}
		})
	}
}

// An unusable entry is reported even when the generators directory does not
// exist: the entries are validated before any package is read.
func TestExcludeEntriesAreValidatedBeforeGeneratorsLoad(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-generators")
	output := t.TempDir()

	_, err := Regions(RegionsConfig{Generators: missing, Output: output, Exclude: []string{"/etc"}})
	if err == nil || !strings.Contains(err.Error(), "exclude entry") {
		t.Errorf("Regions error = %v, want the exclude entry error", err)
	}

	cfg := config(t)
	cfg.Generators = missing
	cfg.Output = output
	putConfigYAML(t, output, "exclude: ['../x']\n")
	_, err = Run(context.Background(), cfg)
	if err == nil || !strings.Contains(err.Error(), "exclude entry") {
		t.Errorf("Run error = %v, want the exclude entry error", err)
	}
}

// treeContents is every file under root with its content.
func treeContents(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		out[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
