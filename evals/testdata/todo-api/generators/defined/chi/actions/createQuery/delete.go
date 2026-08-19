func Delete{{resource|exported}}(ctx context.Context, id int64) error {
	res, err := conn.ExecContext(ctx, `DELETE FROM {{resource|table}} WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
