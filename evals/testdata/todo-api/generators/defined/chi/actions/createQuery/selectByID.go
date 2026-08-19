func Select{{resource|exported}}ByID(ctx context.Context, id int64) ({{resource|exported}}, error) {
	query := fmt.Sprintf(
		`SELECT id, %s, created_at, updated_at FROM {{resource|table}} WHERE id = $1`,
		{{resource|receiver}}ColumnList())
	var out {{resource|exported}}
	if err := conn.QueryRowContext(ctx, query, id).Scan(out.scanTargets()...); err != nil {
		return {{resource|exported}}{}, err
	}
	return out, nil
}
