func TestCreate{{resource|exported}}RejectsAMalformedBody(t *testing.T) {
	rec := run(t, Create{{resource|exported}}, request(http.MethodPost, "/{{resource|table}}", `not json`))

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
