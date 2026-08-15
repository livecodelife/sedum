func TestUpdate{{resource|exported}}RejectsANonIntegerID(t *testing.T) {
	rec := run(t, Update{{resource|exported}}, request(http.MethodPut, "/{{resource|table}}/abc", `{}`))

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
