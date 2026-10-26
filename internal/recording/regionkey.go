package recording

import (
	"encoding/json"
	"strings"
)

// RegionKey returns the string that identifies a region among a record's
// regions (prov-2026-0e23ee6a): the path, action, variant and the JSON
// encoding of the id, joined by NUL. Sorting keys as strings orders regions by
// path, then action, then variant, then id.
//
// The id is encoded as encoding/json marshals a map; a nil id, or one that
// cannot be encoded, is encoded as {}.
func RegionKey(path, action, variant string, id map[string]any) string {
	encoded, err := json.Marshal(id)
	if err != nil || id == nil {
		encoded = []byte("{}")
	}
	return strings.Join([]string{path, action, variant, string(encoded)}, "\x00")
}
