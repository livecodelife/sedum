// Update{{resource|exported}} uses Exec with a RowsAffected check rather than
// RETURNING, then re-reads, so the write and the read stay separate operations.
// The id is $1 so the assignments can number from $2 upward.
func Update{{resource|exported}}(ctx context.Context, id int64, in {{resource|exported}}) ({{resource|exported}}, error) {
	result, err := conn.ExecContext(ctx,
		`UPDATE {{resource|table}} SET {{set_clause}}, updated_at = NOW() WHERE id = $1`,
		{{update_values}})
	if err != nil {
		return {{resource|exported}}{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return {{resource|exported}}{}, err
	}
	if affected == 0 {
		return {{resource|exported}}{}, ErrNotFound
	}

	return Select{{resource|exported}}ByID(ctx, id)
}
