package domain

import "time"

type Application struct {
	ID         int64
	Url        string
	Country    string
	Created_by int64
	Updated_at time.Time
	Created_at time.Time
}

func NewApplication(url string, country string, created_by int64) (*Application, error) {
	application := Application{
		Url:        url,
		Country:    country,
		Created_by: created_by,
		Updated_at: time.Now(),
		Created_at: time.Now(),
	}

	return &application, nil
}
