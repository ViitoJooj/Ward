package repository

func (r *SQLite) DeleteUserByID(id int64) error {
	_, err := r.db.Exec(`
		DELETE FROM users
		WHERE id = ?
	`, id)

	return err
}

func (r *SQLite) DeleteApplicationByID(id int64) error {
	_, err := r.db.Exec(`
		DELETE FROM applications
		WHERE id = $1
	`, id)

	return err
}
