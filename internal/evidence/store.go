package evidence

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Runner is the value of an entry's runner field.
const Runner = "sedum evidence"

// Entry is one line of a record's evidence file: which record version the entry
// covers, the code hash that was seen, and when and by what it was written.
type Entry struct {
	Record   string `json:"record"`
	Version  int    `json:"version"`
	CodeHash string `json:"code_hash"`
	At       string `json:"at"`
	Runner   string `json:"runner"`
}

// Path is the evidence file of a record under the output directory.
func Path(output, record string) string {
	return filepath.Join(output, ".sedum", "evidence", record+".jsonl")
}

// Timestamp is the time in UTC, RFC 3339, with no fractional seconds.
func Timestamp(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05Z")
}

func valid(e Entry) bool {
	return e.Record != "" && e.Version >= 1 && e.CodeHash != ""
}

// Read returns the entries of a record's evidence file in file order, skipping
// blank lines; a record with no file has none. A line that is not a JSON object
// holding a record, a version of at least 1 and a code_hash, or whose record is
// not the one asked for, is an error naming the file and the line.
func Read(output, record string) ([]Entry, error) {
	path := Path(output, record)
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	defer file.Close()

	var entries []Entry
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for line := 1; scanner.Scan(); line++ {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(text), &e); err != nil {
			return nil, fmt.Errorf("%s:%d: not a JSON entry: %v", path, line, err)
		}
		if !valid(e) {
			return nil, fmt.Errorf("%s:%d: an entry needs a record, a version of at least 1 and a code_hash", path, line)
		}
		if e.Record != record {
			return nil, fmt.Errorf("%s:%d: the entry is for record %s, not %s", path, line, e.Record, record)
		}
		entries = append(entries, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return entries, nil
}

// Latest returns the last entry whose version is version, and whether there is one.
func Latest(entries []Entry, version int) (Entry, bool) {
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Version == version {
			return entries[i], true
		}
	}
	return Entry{}, false
}

// Append adds the entry as one JSON line to its record's evidence file, creating
// the .sedum/evidence directory and the file when missing, and never touching what
// the file already holds. The entry is refused, and nothing written, when its
// record is empty, holds a path separator or is "..", its version is below 1 or
// its code hash is empty.
func Append(output string, e Entry) error {
	if e.Record == "" || strings.ContainsAny(e.Record, "/"+string(os.PathSeparator)) || e.Record == ".." {
		return fmt.Errorf("cannot write evidence for the record id %q", e.Record)
	}
	if !valid(e) {
		return errors.New("an evidence entry needs a version of at least 1 and a code_hash")
	}
	path := Path(output, e.Record)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("cannot write %s: %w", path, err)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("cannot write %s: %w", path, err)
	}
	defer file.Close()

	var buf bytes.Buffer
	if info, err := file.Stat(); err == nil && info.Size() > 0 {
		last := make([]byte, 1)
		if _, err := file.ReadAt(last, info.Size()-1); err == nil && last[0] != '\n' {
			buf.WriteByte('\n')
		}
	}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("cannot write %s: %w", path, err)
	}
	if _, err := file.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("cannot write %s: %w", path, err)
	}
	return nil
}
