package runlog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewCreatesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".sedum", "run.log")

	log, err := New(path, false)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer log.Close()

	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatalf("parent directory was not created: %v", err)
	}
}

func TestLogWritesToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "run.log")

	log, err := New(path, false)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	log.Info("resolved package", "path", "app/controllers/users_controller.rb", "package", "rails")
	if err := log.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	for _, want := range []string{"resolved package", "app/controllers/users_controller.rb", "rails"} {
		if !strings.Contains(string(body), want) {
			t.Errorf("run log does not contain %q; got:\n%s", want, body)
		}
	}
}

// The run log is diagnostic output cleared per run. Nothing reads it back, so
// carrying a previous run's lines forward would only ever mislead.
func TestLogIsTruncatedPerRun(t *testing.T) {
	path := filepath.Join(t.TempDir(), "run.log")

	first, err := New(path, false)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	first.Info("from the first run")
	if err := first.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	second, err := New(path, false)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	second.Info("from the second run")
	if err := second.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if strings.Contains(string(body), "from the first run") {
		t.Errorf("run log was not truncated; got:\n%s", body)
	}
	if !strings.Contains(string(body), "from the second run") {
		t.Errorf("run log is missing the current run's output; got:\n%s", body)
	}
}

func TestVerboseMirrorsToTheGivenWriter(t *testing.T) {
	path := filepath.Join(t.TempDir(), "run.log")

	var mirror strings.Builder
	log, err := NewWithMirror(path, &mirror)
	if err != nil {
		t.Fatalf("NewWithMirror: %v", err)
	}
	log.Info("matched template", "template", "app/controllers/{name}_controller.rb")
	if err := log.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if !strings.Contains(mirror.String(), "matched template") {
		t.Errorf("verbose output was not mirrored; got: %q", mirror.String())
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(body), "matched template") {
		t.Errorf("mirroring must not replace the file; got:\n%s", body)
	}
}

func TestNewReportsAnUnwritablePath(t *testing.T) {
	dir := t.TempDir()
	// A file where a directory must be: creating the parent has to fail.
	blocker := filepath.Join(dir, "blocked")
	if err := os.WriteFile(blocker, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if _, err := New(filepath.Join(blocker, "run.log"), false); err == nil {
		t.Fatal("expected an error for an unwritable log path, got nil")
	} else if !strings.Contains(err.Error(), "run.log") {
		t.Errorf("error %q does not name the log path", err)
	}
}
