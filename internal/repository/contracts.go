package repository

import (
	"sync"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/cockroachdb/pebble"
)

type Store struct {
	db *pebble.DB
	mu sync.Mutex
}

func NewRocksDBRepository(db *pebble.DB) (UserRepository, ApplicationRepository, RequestLogRepository) {
	s := &Store{db: db}
	return s, s, s
}

type UserRepository interface {
	CreateUser(user *domain.User) error
	FindUserByID(id int) (*domain.User, error)
	FindUserByEmail(email string) (*domain.User, error)
	FindUserByUsername(username string) (*domain.User, error)
	ListUsers() ([]*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUserByID(id int) error
}

type ApplicationRepository interface {
	CreateApplication(application *domain.Application) error
	FindApplicationByID(id int) (*domain.Application, error)
	FindApplicationByURL(url string) (*domain.Application, error)
	FindApplicationByCountry(country string) (*domain.Application, error)
	ListApplications() ([]*domain.Application, error)
	UpdateApplication(application *domain.Application) error
	DeleteApplicationByID(id int) error
}

type RequestLogRepository interface {
	InsertRequestLog(log *domain.RequestLog) error
	ListRequestLogs() ([]*domain.RequestLog, error)
}
