// Insert{{resource|exported}} writes in two steps - RETURNING id, then a
// separate SELECT - so each statement is a single observable operation.
//
// The statement is built from {{resource|receiver}}Columns rather than from
// column names written here, so the columns, the placeholders and the arguments
// cannot disagree about order: all three come from the same invocations.
func Insert{{resource|exported}}(ctx context.Context, in {{resource|exported}}) ({{resource|exported}}, error) {
	var id int64
	query := fmt.Sprintf(
		`INSERT INTO {{resource|table}} (%s) VALUES (%s) RETURNING id`,
		{{resource|receiver}}ColumnList(), {{resource|receiver}}Placeholders(1))
	if err := conn.QueryRowContext(ctx, query, in.insertArgs()...).Scan(&id); err != nil {
		return {{resource|exported}}{}, err
	}
	return Select{{resource|exported}}ByID(ctx, id)
}
