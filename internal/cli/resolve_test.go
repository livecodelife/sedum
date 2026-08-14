package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The command layer's job is to reach the pipeline and report what it found.
// The phases themselves are tested where they live; what is asserted here is
// that the flags arrive, the output is legible, and the read-only commands stay
// read-only.

func fixtureGenerators() string {
	return filepath.Join("..", "..", "testdata", "generators")
}

// fixtureRecords writes a records directory authorizing paths across two of the
// fixture packages, so that per-file resolution is exercised end to end.
func fixtureRecords(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	body := `id: prov-2026-aaaaaaaa
intent: |
  Add a users controller and its handler.
constraints:
  - Read-only endpoints.
affected_scope:
  - app/controllers/users_controller.rb
  - internal/handlers/user.go
`
	if err := os.WriteFile(filepath.Join(dir, "r.yml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestResolveReportsPackageTemplateAndCaptures(t *testing.T) {
	out, err := exec(t, "resolve", "--generators", fixtureGenerators(), "--records", fixtureRecords(t))
	if err != nil {
		t.Fatalf("resolve: %v\n%s", err, out)
	}

	for _, want := range []string{
		"prov-2026-aaaaaaaa",
		"app/controllers/users_controller.rb",
		"rails",
		"app/controllers/{name}_controller.rb",
		"name=users",
		"internal/handlers/user.go",
		"chi",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("resolve output does not report %q:\n%s", want, out)
		}
	}
}

// resolve is the read-only inspection tool. Running it must not create the
// files it is reporting on.
func TestResolveWritesNothing(t *testing.T) {
	generators, err := filepath.Abs(fixtureGenerators())
	if err != nil {
		t.Fatal(err)
	}
	records := fixtureRecords(t)

	dir := t.TempDir()
	t.Chdir(dir)

	out, err := exec(t, "resolve", "--generators", generators, "--records", records)
	if err != nil {
		t.Fatalf("resolve: %v\n%s", err, out)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("resolve wrote %d entries into the working directory", len(entries))
	}
}

func TestResolveShowTemplateIncludesRenderedOutput(t *testing.T) {
	out, err := exec(t, "resolve", "--generators", fixtureGenerators(), "--records", fixtureRecords(t), "--show-template")
	if err != nil {
		t.Fatalf("resolve: %v\n%s", err, out)
	}

	if !strings.Contains(out, "class UserController") {
		t.Errorf("--show-template did not include the rendered template:\n%s", out)
	}
}

func TestResolveReportsAnUnclaimedExtension(t *testing.T) {
	dir := t.TempDir()
	body := "id: prov-2026-aaaaaaaa\naffected_scope:\n  - docs/readme.tex\n"
	if err := os.WriteFile(filepath.Join(dir, "r.yml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := exec(t, "resolve", "--generators", fixtureGenerators(), "--records", dir)
	wantErr(t, err, ".tex")
}

// --stop-after resolution and --stop-after files are grow's checkpoints inside
// a run that will continue, and they land with this milestone even though the
// phases after them do not.
func TestGrowStopsAfterResolutionWithoutWriting(t *testing.T) {
	dir := t.TempDir()

	out, err := exec(t, "grow",
		"--generators", fixtureGenerators(),
		"--records", fixtureRecords(t),
		"--output", dir,
		"--log", filepath.Join(t.TempDir(), "run.log"),
		"--stop-after", "resolution")
	if err != nil {
		t.Fatalf("grow --stop-after resolution: %v\n%s", err, out)
	}

	if !strings.Contains(out, "stopped after resolution") {
		t.Errorf("output does not say where the run stopped:\n%s", out)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("stopping after resolution wrote %d entries", len(entries))
	}
}

func TestGrowStopsAfterFilesHavingCreatedThem(t *testing.T) {
	dir := t.TempDir()
	records := fixtureRecords(t)

	out, err := exec(t, "grow",
		"--generators", fixtureGenerators(),
		"--records", records,
		"--output", dir,
		"--log", filepath.Join(t.TempDir(), "run.log"),
		"--stop-after", "files")
	if err != nil {
		t.Fatalf("grow --stop-after files: %v\n%s", err, out)
	}

	if !strings.Contains(out, "stopped after files") {
		t.Errorf("output does not say where the run stopped:\n%s", out)
	}
	body, err := os.ReadFile(filepath.Join(dir, "app", "controllers", "users_controller.rb"))
	if err != nil {
		t.Fatalf("the authorized file was not created: %v", err)
	}
	if !strings.Contains(string(body), "sedum:anchor:class_body") {
		t.Errorf("created file carries no markers:\n%s", body)
	}

	// Rerunning is an ordinary resume, not an error and not a rewrite.
	out, err = exec(t, "grow",
		"--generators", fixtureGenerators(),
		"--records", records,
		"--output", dir,
		"--log", filepath.Join(t.TempDir(), "run.log"),
		"--stop-after", "files")
	if err != nil {
		t.Fatalf("rerun: %v\n%s", err, out)
	}
	if !strings.Contains(out, "already present") {
		t.Errorf("a rerun did not report the files it left alone:\n%s", out)
	}
}

// Without a stop point the run would reach a phase that has not landed, and it
// says so before touching the output tree rather than after.
func TestGrowWithNoStopPointRefusesBeforeWriting(t *testing.T) {
	dir := t.TempDir()

	_, err := exec(t, "grow",
		"--generators", fixtureGenerators(),
		"--records", fixtureRecords(t),
		"--output", dir)

	wantErr(t, err, "not implemented", "M6")
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("a refused run wrote %d entries", len(entries))
	}
}

// resolve reads the generators directory and the records directory and nothing
// else, so it must answer identically from anywhere. It used to run the
// filesystem half of Phase 3 against the working directory, which meant an
// unrelated file that merely shared a relative path with an authorized one
// halted the command (prov-2026-43808a47).
func TestResolveDoesNotConsultTheWorkingDirectory(t *testing.T) {
	generators, err := filepath.Abs(fixtureGenerators())
	if err != nil {
		t.Fatal(err)
	}
	records := fixtureRecords(t)

	// A file standing exactly where an authorized path would go, carrying
	// none of the markers its template plants.
	occupied := t.TempDir()
	if err := os.MkdirAll(filepath.Join(occupied, "app", "controllers"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(occupied, "app", "controllers", "users_controller.rb"),
		[]byte("class SomethingElse\nend\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fromEmpty, err := exec(t, "resolve", "--generators", generators, "--records", records)
	if err != nil {
		t.Fatalf("resolve from an empty directory: %v\n%s", err, fromEmpty)
	}

	t.Chdir(occupied)
	fromOccupied, err := exec(t, "resolve", "--generators", generators, "--records", records)
	if err != nil {
		t.Fatalf("resolve reported on a file in the working directory: %v\n%s", err, fromOccupied)
	}

	if fromEmpty != fromOccupied {
		t.Errorf("resolve gave two answers for the same inputs:\n--- from an empty directory:\n%s\n--- from an occupied one:\n%s",
			fromEmpty, fromOccupied)
	}
}

// grow keeps the check, because it has an --output flag to aim it with.
func TestGrowStillReportsAFileMissingItsMarkers(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "app", "controllers"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "app", "controllers", "users_controller.rb"),
		[]byte("class SomethingElse\nend\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := exec(t, "grow",
		"--generators", fixtureGenerators(),
		"--records", fixtureRecords(t),
		"--output", dir,
		"--log", filepath.Join(t.TempDir(), "run.log"),
		"--stop-after", "files")

	wantErr(t, err, "users_controller.rb", "class_body")
}
