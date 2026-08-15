func TestDelete{{resource|exported}}ReportsAMissingRow(t *testing.T) {
	ctx := requireConn(t)

	if err := Delete{{resource|exported}}(ctx, -1); err == nil {
		t.Error("deleting a row that does not exist returned no error")
	}
}
