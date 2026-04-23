package domain

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID         int
	Username   string
	Email      string
	Password   string
	Role       string
	Active     bool
	Updated_at time.Time
	Created_at time.Time
}

func NewUser(username string, email string, password string, active bool, role string) (*User, error) {
	// Username validator
	if len(username) > 250 {
		return nil, errors.New("Name is too large.")
	} else if len(username) < 3 {
		return nil, errors.New("username is too short")
	}

	// Email validator
	if len(email) > 250 {
		return nil, errors.New("Email is too large.")
	} else if len(email) < 10 {
		return nil, errors.New("Email is too short.")
	}

	if !strings.Contains(email, "@") {
		return nil, errors.New("Invalid email.")
	} else if !strings.Contains(email, ".") {
		return nil, errors.New("Invalid email.")
	}

	// Password validator
	if len(password) > 50 {
		return nil, errors.New("Password is too large.")
	} else if len(password) < 8 {
		return nil, errors.New("Password is too short.")
	}

	if role != "admin" && role != "user" {
		return nil, errors.New("Invalid role, choose 'admin' or 'user'")
	}

	user := &User{
		Username:   username,
		Email:      email,
		Password:   password,
		Role:       role,
		Active:     active,
		Updated_at: time.Now(),
		Created_at: time.Now(),
	}

	return user, nil
}
