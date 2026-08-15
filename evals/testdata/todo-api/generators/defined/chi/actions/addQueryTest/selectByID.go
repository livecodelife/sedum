func TestSelect{{resource|exported}}ByIDReportsAMissingRow(t *testing.T) {
	ctx := requireConn(t)

	if _, err := Select{{resource|exported}}ByID(ctx, -1); err == nil {
		t.Error("selecting a row that does not exist returned no error")
	}
}
