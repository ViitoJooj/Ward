package repository

func (r *SQLite) DeleteUserByID(id int) error {
	_, err := r.db.Exec(`
		DELETE FROM users
		WHERE id = ?
	`, id)

	return err
}

func (r *SQLite) DeleteApplicationByID(id int) error {
	_, err := r.db.Exec(`
		DELETE FROM applications
		WHERE id = ?
	`, id)

	return err
}
