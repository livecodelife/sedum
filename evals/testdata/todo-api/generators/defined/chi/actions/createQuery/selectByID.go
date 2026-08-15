func Select{{resource|exported}}ByID(ctx context.Context, id int64) ({{resource|exported}}, error) {
	var {{resource|receiver}} {{resource|exported}}
	err := conn.QueryRowContext(ctx,
		`SELECT id, {{columns}}, created_at, updated_at FROM {{resource|table}} WHERE id = $1`, id).
		Scan({{scan_targets}})
	return {{resource|receiver}}, err
}
