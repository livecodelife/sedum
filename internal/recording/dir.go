package recording

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// checkRecordID rejects a record id that cannot name a recording file.
func checkRecordID(id string) error {
	if id == "" || strings.ContainsAny(id, `/\`) || strings.Contains(id, "..") {
		return fmt.Errorf("record id %q cannot name a recording file", id)
	}
	return nil
}

// WriteRecords writes each record of a recording to its own file,
// <record_id>.json, in dir (prov-2026-2a2edf2b). Each file is a valid
// recording holding exactly that record, with the packages its files name and
// its regions sorted by RegionKey.
func WriteRecords(dir string, rec Recording) error {
	for _, r := range rec.Records {
		if err := checkRecordID(r.RecordID); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating recordings directory %s: %w", dir, err)
	}

	for _, r := range rec.Records {
		one := r
		if r.Regions != nil {
			keys := make(map[int]string, len(r.Regions))
			idx := make([]int, len(r.Regions))
			for i, g := range r.Regions {
				idx[i] = i
				keys[i] = RegionKey(g.Path, g.Action, g.Variant, g.ID)
			}
			sort.SliceStable(idx, func(a, b int) bool { return keys[idx[a]] < keys[idx[b]] })
			sorted := make([]Region, len(r.Regions))
			for i, j := range idx {
				sorted[i] = r.Regions[j]
			}
			one.Regions = sorted
		}

		packages := map[string]Package{}
		for _, f := range r.Files {
			if p, ok := rec.Packages[f.Package]; ok {
				packages[f.Package] = p
			}
		}

		out := Recording{
			SedumVersion: rec.SedumVersion,
			Packages:     packages,
			Records:      []Record{one},
			Variables:    rec.Variables,
		}
		path := filepath.Join(dir, r.RecordID+".json")
		if err := writeFileAtomic(dir, path, out); err != nil {
			return fmt.Errorf("writing recording %s: %w", path, err)
		}
	}
	return nil
}

// writeFileAtomic writes the recording to a temporary file in dir and renames
// it to path.
func writeFileAtomic(dir, path string, rec Recording) (err error) {
	tmp, err := os.CreateTemp(dir, ".recording-*.tmp")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = os.Remove(tmp.Name())
		}
	}()
	if err = Write(tmp, rec); err != nil {
		_ = tmp.Close()
		return err
	}
	if err = tmp.Chmod(0o644); err != nil {
		_ = tmp.Close()
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// ReadRecord reads the recording file dir holds for one record
// (prov-2026-296d6901). The returned recording holds only that record.
func ReadRecord(dir, id string) (Recording, error) {
	if err := checkRecordID(id); err != nil {
		return Recording{}, err
	}
	path := filepath.Join(dir, id+".json")
	f, err := os.Open(path)
	if err != nil {
		return Recording{}, fmt.Errorf("%w", err)
	}
	defer f.Close()

	rec, err := Read(f)
	if err != nil {
		return Recording{}, fmt.Errorf("%s: %s", path, err)
	}
	for _, r := range rec.Records {
		if r.RecordID == id {
			rec.Records = []Record{r}
			return rec, nil
		}
	}
	return Recording{}, fmt.Errorf("%s holds no record %q", path, id)
}
