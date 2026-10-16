package main

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/livecodelife/provenance/record"
)

// created returns the id of the record a create run reported.
func created(t *testing.T, out string) string {
	t.Helper()
	m := regexp.MustCompile(`✓ Created (\S+)\.yml`).FindStringSubmatch(out)
	if m == nil {
		t.Fatalf("no ✓ Created line:\n%s", out)
	}
	return filepath.Base(m[1])
}

func TestCreateDraft(t *testing.T) {
	parent := "prov-2026-0000e001"
	root, cfg := v2Repo(t, false, map[string]string{parent: v2Draft(parent, "open")})
	out := mustRun(t, "-c", cfg, "create", "--no-edit", "--title", "New thing", "--intent", "Make it",
		"--type", "imprint", "--tag", "cli,dogfood", "--tag", "extra", "--scope", "a.go", "--scope", "b/**",
		"--implements", parent)
	id := created(t, out)
	if !regexp.MustCompile(`^prov-\d{4}-[0-9a-f]{8}$`).MatchString(id) || !strings.HasPrefix(id, "prov-"+time.Now().Format("2006")+"-") {
		t.Fatalf("id %q", id)
	}
	if !strings.Contains(out, filepath.Join(root, "provenance", id+".yml")) {
		t.Fatalf("output does not name the file's full path:\n%s", out)
	}
	store := &record.FilesStore{Root: root}
	rec, _, err := store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Schema != 2 || rec.Version != 1 || rec.Status != record.StatusDraft || rec.Type != record.TypeImprint {
		t.Fatalf("schema %d version %d status %s type %s", rec.Schema, rec.Version, rec.Status, rec.Type)
	}
	if rec.Title != "New thing" || rec.Intent != "Make it" || rec.Author != "test@example.com" || rec.CreatedAt != time.Now().Format("2006-01-02") {
		t.Fatalf("fields %+v", rec)
	}
	if !reflect.DeepEqual(rec.Tags, []string{"cli", "dogfood", "extra"}) || !reflect.DeepEqual(rec.AffectedScope, []string{"a.go", "b/**"}) {
		t.Fatalf("tags %v scope %v", rec.Tags, rec.AffectedScope)
	}
	if rec.Implements != parent || rec.ImplementsVersion != 1 {
		t.Fatalf("implements %s@%d", rec.Implements, rec.ImplementsVersion)
	}
	if err := rec.Validate(); err != nil {
		t.Fatal(err)
	}
	if vs, _ := store.Versions(id); len(vs) != 0 {
		t.Fatalf("create wrote snapshots %v", vs)
	}
}

func TestCreateMinimalAndSuffix(t *testing.T) {
	root, cfg := v2Repo(t, false, nil)
	id := created(t, mustRun(t, "-c", cfg, "create", "--no-edit", "--id-suffix", "user-service"))
	if !strings.HasSuffix(id, "-user-service") {
		t.Fatalf("id %q", id)
	}
	rec, _, err := (&record.FilesStore{Root: root}).Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Type != "" || rec.Implements != "" || len(rec.Tags) != 0 || rec.Locked {
		t.Fatalf("defaults %+v", rec)
	}
	// -i is the id-suffix short form, as in linespec.
	if id := created(t, mustRun(t, "-c", cfg, "create", "--no-edit", "-i", "svc")); !strings.HasSuffix(id, "-svc") {
		t.Fatalf("-i id %q", id)
	}
}

func TestCreateOpensEditor(t *testing.T) {
	root, cfg := v2Repo(t, false, nil)
	marker := filepath.Join(root, "edited")
	editor := filepath.Join(root, "editor.sh")
	writeFile(t, editor, "#!/bin/sh\necho \"$1\" > "+marker+"\n")
	if err := os.Chmod(editor, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("EDITOR", editor)
	id := created(t, mustRun(t, "-c", cfg, "create", "--title", "Edit me"))
	if got := strings.TrimSpace(readFile(t, marker)); got != filepath.Join(root, "provenance", id+".yml") {
		t.Fatalf("editor got %q", got)
	}
}

func TestCreateRefusals(t *testing.T) {
	root, cfg := v2Repo(t, false, nil)
	for name, args := range map[string][]string{
		"bad type":           {"create", "--no-edit", "--type", "bug"},
		"unknown implements": {"create", "--no-edit", "--implements", "prov-2026-0000ffff"},
		"stray arg":          {"create", "--no-edit", "extra"},
	} {
		if out, err := run(t, append([]string{"-c", cfg}, args...)...); err == nil {
			t.Errorf("%s: accepted\n%s", name, out)
		}
	}
	if files, _ := filepath.Glob(filepath.Join(root, "provenance", "*.yml")); len(files) != 0 {
		t.Fatalf("refused creates wrote %v", files)
	}
}
