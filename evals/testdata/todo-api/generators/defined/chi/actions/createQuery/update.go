// The row's id is $1, so the assignments number from $2.
func Update{{resource|exported}}(ctx context.Context, id int64, in {{resource|exported}}) ({{resource|exported}}, error) {
	query := fmt.Sprintf(
		`UPDATE {{resource|table}} SET %s, updated_at = now() WHERE id = $1`,
		{{resource|receiver}}Assignments(2))
	args := append([]any{id}, in.insertArgs()...)
	res, err := conn.ExecContext(ctx, query, args...)
	if err != nil {
		return {{resource|exported}}{}, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return {{resource|exported}}{}, err
	}
	if affected == 0 {
		return {{resource|exported}}{}, ErrNotFound
	}
	return Select{{resource|exported}}ByID(ctx, id)
}
