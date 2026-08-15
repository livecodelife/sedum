// Insert{{resource|exported}} writes in two steps - RETURNING id, then a
// separate SELECT - so each statement is a single observable operation.
func Insert{{resource|exported}}(ctx context.Context, in {{resource|exported}}) ({{resource|exported}}, error) {
	var id int64
	err := conn.QueryRowContext(ctx,
		`INSERT INTO {{resource|table}} ({{columns}}) VALUES ({{placeholders}}) RETURNING id`,
		{{insert_values}}).Scan(&id)
	if err != nil {
		return {{resource|exported}}{}, err
	}
	return Select{{resource|exported}}ByID(ctx, id)
}
