func TestUpdate{{resource|exported}}ReportsAMissingRow(t *testing.T) {
	ctx := requireConn(t)

	if _, err := Update{{resource|exported}}(ctx, -1, {{resource|exported}}{}); err == nil {
		t.Error("updating a row that does not exist returned no error")
	}
}
