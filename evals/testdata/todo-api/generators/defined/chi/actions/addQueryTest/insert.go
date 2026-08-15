func TestInsert{{resource|exported}}(t *testing.T) {
	ctx := requireConn(t)

	created, err := Insert{{resource|exported}}(ctx, {{resource|exported}}{})
	if err != nil {
		t.Fatalf("Insert{{resource|exported}}: %v", err)
	}
	if created.ID == 0 {
		t.Error("insert returned a row with no id")
	}
	if created.CreatedAt.IsZero() {
		t.Error("insert returned a row with no created_at")
	}
}
