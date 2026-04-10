package domain

import "time"

type Application struct {
	ID         int
	Url        string
	Country    string
	Created_by int
	Updated_at time.Time
	Created_at time.Time
}

func NewApplication(url string, country string, created_by int) (*Application, error) {
	application := Application{
		Url:        url,
		Country:    country,
		Created_by: created_by,
		Updated_at: time.Now(),
		Created_at: time.Now(),
	}

	return &application, nil
}
