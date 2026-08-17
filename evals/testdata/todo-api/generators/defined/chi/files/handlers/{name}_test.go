package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// request builds a JSON request. Injected tests call it so that this file uses
// every import it declares.
func request(method, target, body string) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// run sends req through h and returns what it wrote.
func run(t *testing.T, h http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

// sedum:anchor:tests
