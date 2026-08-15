package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// sedum:anchor:imports
)

// The helpers below exist so that this file uses every import it declares.
// Go rejects an unused import, so a template that imports on behalf of code
// not yet injected produces a file that will not compile. Injected handlers
// call these instead of repeating the same four lines.

// pathID reads the {id} path segment as an integer.
func pathID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

// writeJSON writes v as a JSON body with the given status.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// decodeBody reads a JSON request body into v.
func decodeBody(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// sedum:anchor:handlers
