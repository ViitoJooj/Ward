package repository

import "github.com/ViitoJooj/door/internal/domain"

func (r *SQLite) CreateUser(user *domain.User) error {
	_, err := r.db.Exec(`INSERT INTO users (username, email, password, updated_at, created_at) VALUES ($1, $2, $3, $4, $5)`,
		user.Username,
		user.Email,
		user.Password,
		user.Updated_at,
		user.Created_at,
	)
	return err
}

func (r *SQLite) CreateApplication(application *domain.Application) error {
	_, err := r.db.Exec(`INSERT INTO applications (url, country, created_by, updated_at, created_at) VALUES ($1, $2, $3, $4, $5)`,
		application.Url,
		application.Country,
		application.Created_by,
		application.Updated_at,
		application.Created_at,
	)
	return err
}
