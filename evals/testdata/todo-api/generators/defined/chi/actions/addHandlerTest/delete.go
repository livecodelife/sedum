func TestDelete{{resource|exported}}RejectsANonIntegerID(t *testing.T) {
	rec := run(t, Delete{{resource|exported}}, request(http.MethodDelete, "/{{resource|table}}/abc", ""))

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
