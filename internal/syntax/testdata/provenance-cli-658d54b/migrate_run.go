package main

// runMigrate converts a repository's v1 records to v2 through migrate.Plan
// and migrate.Apply. Temporary bridge: moves into `go` generator regions once
// that package lands (prov-2026-60e72d0e).

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/livecodelife/provenance/migrate"
)

func runMigrate(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}
	planOnly, _ := cmd.Flags().GetBool("plan")
	asJSON, _ := cmd.Flags().GetBool("json")
	accept, _ := cmd.Flags().GetBool("accept-suggestions")
	keep, _ := cmd.Flags().GetStringArray("keep-separate")
	fold, _ := cmd.Flags().GetStringArray("fold")
	workspace, _ := cmd.Flags().GetString("workspace")
	if asJSON && !planOnly {
		return errors.New("--json needs --plan")
	}

	if !r.Legacy {
		// migrate.Load reads the v1 config, which a migrated repo no longer has.
		v1, err := r.schema1Records()
		if err != nil {
			return err
		}
		if len(v1) == 0 {
			fmt.Fprintf(out, "\n✓ Nothing to migrate: every record in %s is already v2.\n\n", r.Config.Dir)
			return nil
		}
		return fmt.Errorf("%d v1 records (%s, …) but no .linespec.yml provenance block to migrate from", len(v1), v1[0])
	}
	src, err := migrate.Load(r.Root, migrate.LoadOptions{Workspace: workspace})
	if err != nil {
		return err
	}
	if src.Dir != r.Config.Dir {
		return fmt.Errorf("records directory mismatch: %s config says %q, migration reads %q", r.ConfigPath, r.Config.Dir, src.Dir)
	}
	plan, err := migrate.MakePlan(src)
	if err != nil {
		return err
	}
	opts, err := migrateOptions(plan, accept, keep, fold)
	if err != nil {
		return err
	}

	if planOnly {
		view := planView(plan, opts)
		if asJSON {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(view)
		}
		printPlan(out, view)
		return nil
	}

	if !needsMigration(src) {
		fmt.Fprintf(out, "\n✓ Nothing to migrate: every record in %s is already v2.\n\n", src.Dir)
		return nil
	}
	if dirty, err := r.git("status", "--porcelain", "--", r.Config.Dir, legacyConfigFile, configFile); err == nil && dirty != "" {
		r.fail("Cannot migrate: %s has uncommitted changes.\n\n%s\n\n  Commit or stash them first, so the migration is one reviewable change.", r.Config.Dir, dirty)
		return errors.New("dirty records directory")
	}

	res, err := migrate.Apply(plan, opts, r.Store)
	if err != nil {
		return fmt.Errorf("migrate: %w (records written so far: %d; inspect with git status and restore with git checkout)", err, res.Records)
	}
	moved := ""
	if r.Legacy {
		if moved, err = moveConfig(r.Root); err != nil {
			return err
		}
	}
	if err := r.writeManifest(); err != nil {
		return err
	}
	removedManifest := false
	if err := os.Remove(filepath.Join(r.Root, legacyManifest)); err == nil {
		removedManifest = true
	} else if !os.IsNotExist(err) {
		return err
	}

	fmt.Fprintf(out, "\n✓ Migrated %s to v2: %d records (from %d), %d snapshots\n\n", src.Dir, res.Records, len(src.Records), res.Snapshots)
	if len(res.Removed) > 0 {
		fmt.Fprintf(out, "  Folded into other records (kept as aliases): %d\n", len(res.Removed))
	}
	if moved != "" {
		fmt.Fprintf(out, "  Config: %s\n", moved)
	}
	fmt.Fprintf(out, "  Version manifest: %s\n", manifestPath)
	if removedManifest {
		fmt.Fprintf(out, "  Removed the v1 hash manifest: %s\n", legacyManifest)
	}
	fmt.Fprintln(out, "  Hooks: reinstall them with `provenance hooks install` once it is available (CLI part 2).")
	fmt.Fprintln(out)
	return nil
}

// migrateOptions turns the keep-separate flags into migrate.Options. Every
// id they name must be an absorbed record in the plan.
func migrateOptions(p *migrate.Plan, accept bool, keep, fold []string) (migrate.Options, error) {
	opts := migrate.Options{AcceptSuggestions: accept, KeepSeparate: map[string]bool{}}
	for _, set := range []struct {
		ids  []string
		keep bool
	}{{keep, true}, {fold, false}} {
		for _, id := range set.ids {
			if _, ok := p.Entries[id]; !ok {
				return opts, fmt.Errorf("%s is not folded into another record by this plan", id)
			}
			if prev, dup := opts.KeepSeparate[id]; dup && prev != set.keep {
				return opts, fmt.Errorf("%s is named by both --keep-separate and --fold", id)
			}
			opts.KeepSeparate[id] = set.keep
		}
	}
	return opts, nil
}

// needsMigration reports whether any local record is still schema 1.
func needsMigration(src *migrate.Repo) bool {
	for _, s := range src.Records {
		if s.Rec.Schema != 2 {
			return true
		}
	}
	return false
}

// schema1Records lists the ids of records still in the v1 schema.
func (r *repo) schema1Records() ([]string, error) {
	ids, err := r.Store.List()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, id := range ids {
		rec, _, err := r.Store.Load(id)
		if err != nil {
			return nil, err
		}
		if rec.Schema != 2 {
			out = append(out, id)
		}
	}
	return out, nil
}

// planJSON is what `migrate --plan --json` prints.
type planJSON struct {
	Repo           string              `json:"repo"`
	RecordsBefore  int                 `json:"records_before"`
	RecordsAfter   int                 `json:"records_after"`
	Chains         map[string][]string `json:"chains"`
	Entries        []entryJSON         `json:"entries"`
	Suggestions    []suggestionJSON    `json:"suggestions"`
	KeptDeprecated []keptJSON          `json:"kept_deprecated"`
	StatusMappings []map[string]any    `json:"status_mappings"`
	Warnings       []map[string]any    `json:"warnings"`
}

type entryJSON struct {
	ID      string `json:"id"`
	Target  string `json:"target"`
	Rule    string `json:"rule"`
	Suggest string `json:"suggest,omitempty"`
	Kept    bool   `json:"kept"`
}

type suggestionJSON struct {
	ID     string `json:"id"`
	Target string `json:"target"`
	Rule   string `json:"rule"`
	Why    string `json:"why"`
}

type keptJSON struct {
	ID         string `json:"id"`
	Extends    string `json:"extends"`
	Implements string `json:"implements,omitempty"`
}

// planView is the plan under opts: a kept-separate id leaves its chain and
// stands on its own.
func planView(p *migrate.Plan, opts migrate.Options) planJSON {
	v := planJSON{
		Repo: p.Repo.Name, RecordsBefore: len(p.Repo.Records), RecordsAfter: p.RecordsAfter(opts),
		Chains: map[string][]string{}, Entries: []entryJSON{}, Suggestions: []suggestionJSON{},
		KeptDeprecated: []keptJSON{}, StatusMappings: []map[string]any{}, Warnings: []map[string]any{},
	}
	for s, ids := range p.Versions {
		var chain []string
		for _, id := range ids {
			if p.SurvivorWith(id, opts) == s {
				chain = append(chain, id)
			}
		}
		if len(chain) > 1 {
			v.Chains[s] = chain
		}
	}
	ids := make([]string, 0, len(p.Entries))
	for id := range p.Entries {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		for _, e := range p.Entries[id] {
			v.Entries = append(v.Entries, entryJSON{ID: id, Target: e.Target, Rule: e.Rule, Suggest: e.Suggest, Kept: p.SurvivorWith(id, opts) == id})
		}
	}
	for _, s := range p.Suggestions {
		v.Suggestions = append(v.Suggestions, suggestionJSON{ID: s.ID, Target: s.Target, Rule: s.Rule, Why: s.Why})
	}
	for _, k := range p.KeptDeprecated {
		v.KeptDeprecated = append(v.KeptDeprecated, keptJSON{ID: k.ID, Extends: k.Extends, Implements: k.AsImprint.Implements})
	}
	for _, w := range p.StatusMappings(opts) {
		v.StatusMappings = append(v.StatusMappings, warningMap(w))
	}
	for _, w := range p.Warnings {
		v.Warnings = append(v.Warnings, warningMap(w))
	}
	return v
}

func warningMap(w migrate.Warning) map[string]any {
	m := map[string]any{"kind": w.Kind}
	for k, val := range w.Fields {
		m[k] = val
	}
	return m
}

func printPlan(out io.Writer, v planJSON) {
	fmt.Fprintf(out, "\nMigration plan for %s: %d records → %d\n", v.Repo, v.RecordsBefore, v.RecordsAfter)

	survivors := make([]string, 0, len(v.Chains))
	for s := range v.Chains {
		survivors = append(survivors, s)
	}
	sort.Strings(survivors)
	fmt.Fprintf(out, "\nVersion chains (%d), oldest version first:\n", len(survivors))
	for _, s := range survivors {
		fmt.Fprintf(out, "  %s  ← %s\n", s, strings.Join(v.Chains[s], ", "))
	}
	if len(v.Suggestions) > 0 {
		fmt.Fprintf(out, "\nSuggested to keep separate (%d); accept with --accept-suggestions or --keep-separate <id>:\n", len(v.Suggestions))
		for _, s := range v.Suggestions {
			fmt.Fprintf(out, "  %s → %s (%s): %s\n", s.ID, s.Target, s.Rule, s.Why)
		}
	}
	if len(v.KeptDeprecated) > 0 {
		fmt.Fprintf(out, "\nDeprecated bugs kept as deprecated imprints (%d):\n", len(v.KeptDeprecated))
		for _, k := range v.KeptDeprecated {
			fmt.Fprintf(out, "  %s (extends %s)\n", k.ID, k.Extends)
		}
	}
	if len(v.StatusMappings) > 0 {
		fmt.Fprintf(out, "\nStatus mappings to review (%d):\n", len(v.StatusMappings))
		for _, m := range v.StatusMappings {
			fmt.Fprintf(out, "  %v: %v → %v\n", m["id"], m["from"], m["to"])
		}
	}
	if len(v.Warnings) > 0 {
		fmt.Fprintf(out, "\nWarnings (%d):\n", len(v.Warnings))
		for _, w := range v.Warnings {
			fmt.Fprintf(out, "  %s\n", formatWarning(w))
		}
	}
	fmt.Fprintln(out, "\nNothing was written. Run `provenance migrate` with the same flags to apply.")
	fmt.Fprintln(out)
}

func formatWarning(w map[string]any) string {
	keys := make([]string, 0, len(w))
	for k := range w {
		if k != "kind" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	parts := []string{fmt.Sprint(w["kind"])}
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, w[k]))
	}
	return strings.Join(parts, " ")
}

// moveConfig moves the provenance block of .linespec.yml to .provenance.yml,
// dedented to the top level and otherwise byte for byte, and removes it from
// .linespec.yml, deleting that file when nothing else is left in it.
func moveConfig(root string) (string, error) {
	legacy := filepath.Join(root, legacyConfigFile)
	data, err := os.ReadFile(legacy)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(filepath.Join(root, configFile)); err == nil {
		return "", fmt.Errorf("%s already exists; move the provenance block from %s by hand", configFile, legacyConfigFile)
	}
	lines := strings.SplitAfter(string(data), "\n")
	start := -1
	for i, l := range lines {
		if strings.TrimRight(l, " \r\n") == "provenance:" {
			start = i
			break
		}
	}
	if start < 0 {
		return "", fmt.Errorf("%s has no block-style provenance: key to move", legacyConfigFile)
	}
	end := start + 1
	for end < len(lines) {
		l := lines[end]
		trimmed := strings.TrimSpace(l)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") && l[0] != ' ' && l[0] != '\t' {
			break
		}
		end++
	}
	// Trailing blank and comment lines belong to whatever follows the block.
	for end > start+1 {
		trimmed := strings.TrimSpace(lines[end-1])
		if trimmed != "" && !(strings.HasPrefix(trimmed, "#") && lines[end-1][0] == '#') {
			break
		}
		end--
	}
	block := lines[start+1 : end]
	indent := ""
	for _, l := range block {
		if strings.TrimSpace(l) != "" {
			indent = l[:len(l)-len(strings.TrimLeft(l, " \t"))]
			break
		}
	}
	var moved strings.Builder
	for _, l := range block {
		moved.WriteString(strings.TrimPrefix(l, indent))
	}
	text, err := rewriteExcludePaths(moved.String())
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(root, configFile), []byte(text), 0o644); err != nil {
		return "", err
	}
	rest := strings.Join(append(append([]string{}, lines[:start]...), lines[end:]...), "")
	if strings.TrimSpace(stripComments(rest)) == "" {
		if err := os.Remove(legacy); err != nil {
			return "", err
		}
		return fmt.Sprintf("moved to %s; %s removed (it held only the provenance block)", configFile, legacyConfigFile), nil
	}
	if err := os.WriteFile(legacy, []byte(rest), 0o644); err != nil {
		return "", err
	}
	return fmt.Sprintf("provenance block moved from %s to %s", legacyConfigFile, configFile), nil
}

// rewriteExcludePaths fixes the moved config's exclude_paths for the
// migrated repo: .linespec.yml becomes .provenance.yml (once), and the v1
// hash manifest, which migrate deletes, is dropped. Only that key's lines are
// rewritten, in the list's own style; everything else stays byte for byte.
func rewriteExcludePaths(text string) (string, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(text), &doc); err != nil {
		return "", fmt.Errorf("%s: %w", configFile, err)
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return text, nil
	}
	root := doc.Content[0]
	for i := 0; i+1 < len(root.Content); i += 2 {
		key, seq := root.Content[i], root.Content[i+1]
		if key.Value != "exclude_paths" || seq.Kind != yaml.SequenceNode {
			continue
		}
		var items []*yaml.Node
		changed, seen := false, map[string]bool{}
		for _, it := range seq.Content {
			v := it.Value
			switch v {
			case legacyConfigFile:
				v, changed = configFile, true
			case legacyManifest:
				changed = true
				continue
			}
			if seen[v] {
				changed = true
				continue
			}
			seen[v] = true
			n := *it
			n.Value = v
			items = append(items, &n)
		}
		if !changed {
			return text, nil
		}
		next := len(text)
		lines := strings.SplitAfter(text, "\n")
		if i+2 < len(root.Content) {
			next = offsetOfLine(lines, root.Content[i+2].Line)
		}
		start := offsetOfLine(lines, key.Line)
		newSeq := *seq
		newSeq.Content = items
		out, err := yaml.Marshal(&yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{key, &newSeq}})
		if err != nil {
			return "", err
		}
		// Keep the blank and comment lines that end the old span: they belong
		// to the key that follows.
		span := strings.SplitAfter(text[start:next], "\n")
		cut := len(span)
		for cut > 1 {
			l := strings.TrimSpace(span[cut-1])
			if l != "" && !strings.HasPrefix(l, "#") {
				break
			}
			cut--
		}
		keep := strings.Join(span[cut:], "")
		return text[:start] + string(out) + keep + text[next:], nil
	}
	return text, nil
}

// offsetOfLine is the byte offset where 1-based line n starts.
func offsetOfLine(lines []string, n int) int {
	off := 0
	for i := 0; i < n-1 && i < len(lines); i++ {
		off += len(lines[i])
	}
	return off
}

func stripComments(s string) string {
	var b strings.Builder
	for _, l := range strings.Split(s, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(l), "#") {
			b.WriteString(l)
			b.WriteByte('\n')
		}
	}
	return b.String()
}
