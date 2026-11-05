package selection

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/livecodelife/sedum/internal/expand"
	"github.com/livecodelife/sedum/internal/genpkg"
	"github.com/livecodelife/sedum/internal/inject"
	"github.com/livecodelife/sedum/internal/recording"
	"github.com/livecodelife/sedum/internal/resolve"
)

// prov-2026-291dbe16: the Phase 5 syntax preview reconciles the invoking
// record's owned regions the way Phase 7 would, in memory only.
//
// Every fixture is a real Go file written by the real inject.Apply with the
// repository's own go generator, so the markers on disk are the ones a real
// run leaves. The only thing under test is what CheckSyntaxReconciled reports (and, for constraint 11, what CheckSyntax still does not reconcile).

const reconcileFile = "pkg/a.go"

const goFileTemplate = "// sedum:anchor:header\n\npackage main\n\n// sedum:anchor:imports\n\n// sedum:anchor:decls\n"

func goPackage(t *testing.T) *genpkg.Package {
	t.Helper()
	set, findings, err := genpkg.Load(filepath.Join("..", "..", "generators"), genpkg.Options{Only: []string{"go"}})
	if err != nil {
		t.Fatalf("load go generator: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("go generator does not load: %s", f)
		}
	}
	pkg, ok := set.Lookup("go")
	if !ok {
		t.Fatal("go generator package not found")
	}
	return pkg
}

func writeOnDisk(t *testing.T, root, path, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readOnDisk(t *testing.T, root, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func goFiles(pkg *genpkg.Package, record string, paths ...string) []resolve.File {
	var out []resolve.File
	for _, p := range paths {
		out = append(out, resolve.File{
			Resolution: resolve.Resolution{RecordID: record, Path: p, Package: pkg},
			Rendered:   goFileTemplate,
			Existed:    true,
		})
	}
	return out
}

func addImport(path, imp string) recording.Invocation {
	return recording.Invocation{Action: "addImport", Kwargs: map[string]any{"file": path, "path": imp}}
}

func setFunction(path, name, code string) recording.Invocation {
	return recording.Invocation{Action: "setFunction", Kwargs: map[string]any{"file": path, "name": name, "code": code}}
}

// write runs Phase 7 for record over path, so the file on disk holds the
// regions record owns after that run. tier, when set, overrides the tier every
// region is written with.
func write(t *testing.T, root string, pkg *genpkg.Package, record, path string, tier inject.Tier, invs ...recording.Invocation) {
	t.Helper()
	files := goFiles(pkg, record, path)
	expanded, err := expand.Expand(record, files, invs, nil)
	if err != nil {
		t.Fatalf("expand %s: %v", record, err)
	}
	if tier != "" {
		for i := range expanded {
			expanded[i].Tier = tier
		}
	}
	if _, err := inject.Apply(expanded, inject.Options{Output: root}); err != nil {
		t.Fatalf("apply %s: %v", record, err)
	}
}

// fresh returns an output directory holding the go template at path.
func fresh(t *testing.T, path string) string {
	t.Helper()
	root := t.TempDir()
	writeOnDisk(t, root, path, goFileTemplate)
	return root
}

// dropRegion deletes the whole region of action from the file on disk, markers
// and body, leaving everything else byte for byte.
func dropRegion(t *testing.T, root, path, action string) {
	t.Helper()
	lines := strings.SplitAfter(readOnDisk(t, root, path), "\n")
	var out []string
	in := false
	for _, l := range lines {
		switch {
		case strings.Contains(l, "// sedum:"+action+" "):
			in = true
		case strings.Contains(l, "// /sedum:"+action):
			in = false
		case !in:
			out = append(out, l)
		}
	}
	writeOnDisk(t, root, path, strings.Join(out, ""))
}

// writeUnusedImport leaves path holding an import region of record that
// nothing uses. Phase 7 would refuse to write that, so the fixture is written
// with a user of the import and the user is then removed by hand.
func writeUnusedImport(t *testing.T, root string, pkg *genpkg.Package, record, path string, tier inject.Tier, imp string) {
	t.Helper()
	write(t, root, pkg, record, path, tier,
		addImport(path, imp),
		setFunction(path, "Scaffold", "func Scaffold() int {\n\treturn "+imp+".Getpid()\n}\n"))
	dropRegion(t, root, path, "setFunction")
}

func details(vs []Violation) string {
	var parts []string
	for _, v := range vs {
		parts = append(parts, v.Detail)
	}
	return strings.Join(parts, "\n---\n")
}

func wantNoViolations(t *testing.T, got []Violation) {
	t.Helper()
	if len(got) != 0 {
		t.Fatalf("check reported %d violation(s), want none:\n%s", len(got), details(got))
	}
}

func wantViolationMentioning(t *testing.T, got []Violation, substr string) {
	t.Helper()
	if len(got) == 0 {
		t.Fatalf("check reported nothing, want a violation mentioning %q", substr)
	}
	for _, v := range got {
		if v.Rule != RuleSyntax {
			t.Errorf("violation rule = %q, want %q", v.Rule, RuleSyntax)
		}
	}
	if !strings.Contains(details(got), substr) {
		t.Fatalf("no violation mentions %q:\n%s", substr, details(got))
	}
}

const usesSyscall = "func Foo() int {\n\treturn syscall.Getpid()\n}\n"
const plainFoo = "func Foo() int {\n\treturn 1\n}\n"

// The prov-2026-aabeac10 v4 regression: the record's earlier version imported
// syscall, the new one no longer uses it, and the stale import must not fail it.
func TestCheckSyntaxReconciled_StaleOwnedImportPrunedInPreview(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	write(t, root, pkg, "R", reconcileFile, "",
		addImport(reconcileFile, "syscall"), setFunction(reconcileFile, "Foo", usesSyscall))

	got := CheckSyntaxReconciled("R", goFiles(pkg, "R", reconcileFile),
		[]recording.Invocation{setFunction(reconcileFile, "Foo", plainFoo)}, nil, root)
	wantNoViolations(t, got)
}

// A stale owned declaration the new selection drops. It names a package the
// file no longer imports (the import was removed by hand), so it only passes
// if the preview prunes it.
func TestCheckSyntaxReconciled_StaleOwnedDeclarationPrunedInPreview(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	write(t, root, pkg, "R", reconcileFile, "",
		addImport(reconcileFile, "strings"),
		setFunction(reconcileFile, "Old", "func Old() string {\n\treturn strings.ToUpper(\"x\")\n}\n"),
		setFunction(reconcileFile, "Foo", plainFoo))
	dropRegion(t, root, reconcileFile, "addImport")

	got := CheckSyntaxReconciled("R", goFiles(pkg, "R", reconcileFile),
		[]recording.Invocation{setFunction(reconcileFile, "Foo", plainFoo)}, nil, root)
	wantNoViolations(t, got)
}

// A region another record owns is not R's to prune: S's import is used only by
// S's code, and both stay in the buffer, so nothing is reported.
func TestCheckSyntaxReconciled_OtherRecordRegionNotPruned(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	write(t, root, pkg, "S", reconcileFile, "",
		addImport(reconcileFile, "syscall"), setFunction(reconcileFile, "Bar", strings.ReplaceAll(usesSyscall, "Foo", "Bar")))

	got := CheckSyntaxReconciled("R", goFiles(pkg, "R", reconcileFile),
		[]recording.Invocation{setFunction(reconcileFile, "Foo", plainFoo)}, nil, root)
	wantNoViolations(t, got)
}

// The same, from the other side: S's region is itself broken (an import
// nothing uses). It is still in the buffer, so it is still reported.
func TestCheckSyntaxReconciled_OtherRecordBrokenRegionStillReported(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	writeUnusedImport(t, root, pkg, "S", reconcileFile, "", "syscall")

	got := CheckSyntaxReconciled("R", goFiles(pkg, "R", reconcileFile),
		[]recording.Invocation{setFunction(reconcileFile, "Foo", plainFoo)}, nil, root)
	wantViolationMentioning(t, got, "syscall")
}

// A seeded region of the invoking record is kept: its unused import is still in
// the buffer, so it is reported. Were it pruned, nothing would be.
func TestCheckSyntaxReconciled_SeededRegionKept(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	writeUnusedImport(t, root, pkg, "R", reconcileFile, inject.TierSeeded, "syscall")

	got := CheckSyntaxReconciled("R", goFiles(pkg, "R", reconcileFile),
		[]recording.Invocation{setFunction(reconcileFile, "Foo", plainFoo)}, nil, root)
	wantViolationMentioning(t, got, "syscall")
}

// R selects again regions it already owns, with new code: they are replaced in
// place. Pruning either would leave Foo naming an unimported package.
func TestCheckSyntaxReconciled_ReselectedRegionNotPruned(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	write(t, root, pkg, "R", reconcileFile, "",
		addImport(reconcileFile, "syscall"), setFunction(reconcileFile, "Foo", usesSyscall))

	got := CheckSyntaxReconciled("R", goFiles(pkg, "R", reconcileFile),
		[]recording.Invocation{
			addImport(reconcileFile, "syscall"),
			setFunction(reconcileFile, "Foo", "func Foo() int {\n\treturn syscall.Getppid()\n}\n"),
		}, nil, root)
	wantNoViolations(t, got)
}

type snapshotEntry struct {
	data  string
	mtime time.Time
}

func snapshot(t *testing.T, root string) map[string]snapshotEntry {
	t.Helper()
	out := map[string]snapshotEntry{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out[p] = snapshotEntry{string(data), info.ModTime()}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// The preview is in memory only: on the very case that prunes a region, every
// byte, every mtime, and the set of files are unchanged afterwards.
func TestCheckSyntaxReconciled_DoesNotWriteToDisk(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	write(t, root, pkg, "R", reconcileFile, "",
		addImport(reconcileFile, "syscall"), setFunction(reconcileFile, "Foo", usesSyscall))
	// A file on disk the selection reaches only as a new path.
	files := append(goFiles(pkg, "R", reconcileFile), goFiles(pkg, "R", "pkg/new.go")...)
	files[1].Existed = false

	before := snapshot(t, root)
	CheckSyntaxReconciled("R", files, []recording.Invocation{
		setFunction(reconcileFile, "Foo", plainFoo),
		setFunction("pkg/new.go", "Baz", "func Baz() {}\n"),
	}, nil, root)
	after := snapshot(t, root)

	if !reflect.DeepEqual(before, after) {
		t.Fatalf("the output directory changed.\nbefore: %d files\nafter:  %d files\n(first file now: %q)",
			len(before), len(after), after[filepath.Join(root, filepath.FromSlash(reconcileFile))].data)
	}
}

// Pruning what is stale must not hide what is wrong: an import the new code
// does not use is reported, and the stale one is not.
func TestCheckSyntaxReconciled_RealUnusedImportStillReported(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	write(t, root, pkg, "R", reconcileFile, "",
		addImport(reconcileFile, "syscall"), setFunction(reconcileFile, "Foo", usesSyscall))

	got := CheckSyntaxReconciled("R", goFiles(pkg, "R", reconcileFile),
		[]recording.Invocation{addImport(reconcileFile, "os"), setFunction(reconcileFile, "Foo", plainFoo)}, nil, root)
	wantViolationMentioning(t, got, "os (line")
	if strings.Contains(details(got), "syscall") {
		t.Errorf("the stale syscall import is reported although it is pruned:\n%s", details(got))
	}
}

func TestCheckSyntaxReconciled_MissingQualifierStillReported(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	write(t, root, pkg, "R", reconcileFile, "",
		addImport(reconcileFile, "syscall"), setFunction(reconcileFile, "Foo", usesSyscall))

	got := CheckSyntaxReconciled("R", goFiles(pkg, "R", reconcileFile),
		[]recording.Invocation{setFunction(reconcileFile, "Foo", "func Foo() string {\n\treturn strings.ToUpper(\"x\")\n}\n")}, nil, root)
	wantViolationMentioning(t, got, "strings")
}

// A selection Phase 7 cannot apply is reported by Phase 7, not here.
func TestCheckSyntaxReconciled_UnappliableSelectionReportsNothing(t *testing.T) {
	pkg := goPackage(t)

	t.Run("ownership conflict", func(t *testing.T) {
		root := fresh(t, reconcileFile)
		write(t, root, pkg, "S", reconcileFile, "", setFunction(reconcileFile, "Foo", plainFoo))
		got := CheckSyntaxReconciled("R", goFiles(pkg, "R", reconcileFile),
			[]recording.Invocation{setFunction(reconcileFile, "Foo", plainFoo)}, nil, root)
		wantNoViolations(t, got)
	})

	t.Run("missing anchor", func(t *testing.T) {
		root := t.TempDir()
		writeOnDisk(t, root, reconcileFile, "package main\n")
		files := goFiles(pkg, "R", reconcileFile)
		files[0].Rendered = "package main\n"
		got := CheckSyntaxReconciled("R", files,
			[]recording.Invocation{addImport(reconcileFile, "os")}, nil, root)
		wantNoViolations(t, got)
	})
}

func TestCheckSyntaxReconciled_NoSyntaxPackageOrEmptySelection(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	write(t, root, pkg, "R", reconcileFile, "", addImport(reconcileFile, "syscall"), setFunction(reconcileFile, "Foo", usesSyscall))

	t.Run("empty selection", func(t *testing.T) {
		wantNoViolations(t, CheckSyntaxReconciled("R", goFiles(pkg, "R", reconcileFile), nil, nil, root))
	})

	t.Run("package without a syntax key", func(t *testing.T) {
		bare := *pkg
		bare.Syntax = ""
		got := CheckSyntaxReconciled("R", goFiles(&bare, "R", reconcileFile),
			[]recording.Invocation{addImport(reconcileFile, "os")}, nil, root)
		wantNoViolations(t, got)
	})

	t.Run("no package at all", func(t *testing.T) {
		files := goFiles(nil, "R", reconcileFile)
		got := CheckSyntaxReconciled("R", files,
			[]recording.Invocation{addImport(reconcileFile, "os")}, nil, root)
		wantNoViolations(t, got)
	})
}

// An unmanaged file is not Sedum's to write, so the preview does not reconcile
// it. Its stale-looking R-owned import is used by hand-written code, so pruning
// it would break the file: it must be left alone.
func TestCheckSyntaxReconciled_UnmanagedFileNotReconciled(t *testing.T) {
	pkg := goPackage(t)
	const other = "pkg/b.go"
	root := fresh(t, reconcileFile)
	writeOnDisk(t, root, other, goFileTemplate)
	writeUnusedImport(t, root, pkg, "R", other, "", "syscall")
	writeOnDisk(t, root, other, readOnDisk(t, root, other)+"\nfunc ByHand() int {\n\treturn syscall.Getpid()\n}\n")

	files := append(goFiles(pkg, "R", reconcileFile), goFiles(pkg, "R", other)...)
	files[1].Unmanaged = true

	got := CheckSyntaxReconciled("R", files,
		[]recording.Invocation{setFunction(reconcileFile, "Foo", plainFoo)}, nil, root)
	wantNoViolations(t, got)
}

// Constraint 11: CheckSyntax does not reconcile. The very setup the reconciled
// preview passes (the stale syscall import of R's earlier version) is still
// reported here, because replay does not prune and the stale region is in the
// buffer. Replay's validation depends on this.
func TestCheckSyntax_NonReconcilingKeepsStaleRegion(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	write(t, root, pkg, "R", reconcileFile, "",
		addImport(reconcileFile, "syscall"), setFunction(reconcileFile, "Foo", usesSyscall))
	selection := []recording.Invocation{setFunction(reconcileFile, "Foo", plainFoo)}

	wantViolationMentioning(t, CheckSyntax("R", goFiles(pkg, "R", reconcileFile), selection, nil, root), "syscall")
}

// Constraint 12: grow's syntax check reconciles. A revised record's answer
// drops the import its earlier version declared; Select must accept it. With
// the non-reconciling check the stale import is reported on every attempt and
// the record is rejected after exhausting its retries.
func TestSelect_GrowAcceptsARevisionThatDropsAnImport(t *testing.T) {
	pkg := goPackage(t)
	root := fresh(t, reconcileFile)
	write(t, root, pkg, "R", reconcileFile, "",
		addImport(reconcileFile, "syscall"), setFunction(reconcileFile, "Foo", usesSyscall))

	response := `{"invocations":[{"action":"setFunction","kwargs":{"file":"` + reconcileFile +
		`","name":"Foo","code":"func Foo() int {\n\treturn 1\n}\n"}}]}`
	// A second copy is staged in case the completeness question is asked; it
	// does not change what is being judged.
	client := &stub{responses: []string{response, response, response, response}}
	answer, err := Select(context.Background(), client, Request{
		RecordID: "R",
		Intent:   "Foo no longer needs syscall.",
		Files:    goFiles(pkg, "R", reconcileFile),
	}, Options{Retries: 3, Output: root})
	if err != nil {
		t.Fatalf("grow rejected a revision that only drops an import: %v", err)
	}
	if answer.Rejected != 0 {
		t.Errorf("Rejected = %d, want 0: the stale import was reported to the model", answer.Rejected)
	}
}
