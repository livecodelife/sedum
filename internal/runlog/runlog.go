// Package runlog writes Sedum's diagnostic run log.
//
// The log records package resolution, file template matches and captures, the
// model's raw response, validation failures and retries, composite expansions,
// resolved paths, selected variants, transform resolutions, and anchor matches.
// It is diagnostic output, cleared per run. Nothing in Sedum reads it back:
// idempotency state lives in the ownership markers written into generated
// files, not here.
package runlog

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// Log is a run log open on a file, and optionally mirrored to a second writer.
type Log struct {
	*slog.Logger

	file *os.File
}

// New opens path as the run log, creating parent directories and truncating any
// previous run's output. When verbose is set the log is mirrored to stdout.
func New(path string, verbose bool) (*Log, error) {
	var mirror io.Writer
	if verbose {
		mirror = os.Stdout
	}
	return NewWithMirror(path, mirror)
}

// NewWithMirror is New with an explicit mirror destination. A nil mirror writes
// to the file alone.
func NewWithMirror(path string, mirror io.Writer) (*Log, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create run log directory for %s: %w", path, err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open run log %s: %w", path, err)
	}

	var dest io.Writer = file
	if mirror != nil {
		dest = io.MultiWriter(file, mirror)
	}

	handler := slog.NewTextHandler(dest, &slog.HandlerOptions{Level: slog.LevelDebug})
	return &Log{Logger: slog.New(handler), file: file}, nil
}

// Close releases the underlying file.
func (l *Log) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	return l.file.Close()
}

// Discard returns a run log that writes nowhere. It exists so callers that have
// no log configured can hold a non-nil logger rather than nil-checking at every
// call site.
func Discard() *Log {
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})
	return &Log{Logger: slog.New(handler)}
}
