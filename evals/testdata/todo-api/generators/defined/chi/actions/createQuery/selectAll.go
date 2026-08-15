// SelectAll{{resource|collection}} prepares before querying, which forces the
// Extended Query Protocol rather than a simple query.
func SelectAll{{resource|collection}}(ctx context.Context) ([]{{resource|exported}}, error) {
	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, {{columns}}, created_at, updated_at FROM {{resource|table}} ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []{{resource|exported}}{}
	for rows.Next() {
		var {{resource|receiver}} {{resource|exported}}
		if err := rows.Scan({{scan_targets}}); err != nil {
			return nil, err
		}
		out = append(out, {{resource|receiver}})
	}
	return out, rows.Err()
}
