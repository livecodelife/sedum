func TestSelectAll{{resource|collection}}(t *testing.T) {
	ctx := requireConn(t)

	all, err := SelectAll{{resource|collection}}(ctx)
	if err != nil {
		t.Fatalf("SelectAll{{resource|collection}}: %v", err)
	}
	if all == nil {
		t.Error("select all returned nil rather than an empty slice")
	}
}
