func TestGet{{resource|exported}}RejectsANonIntegerID(t *testing.T) {
	rec := run(t, Get{{resource|exported}}, request(http.MethodGet, "/{{resource|table}}/abc", ""))

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
