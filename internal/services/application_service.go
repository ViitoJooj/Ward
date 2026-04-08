package services

import (
	"errors"
	"log"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/ViitoJooj/door/internal/repository"
)

type ApplicationService struct {
	ApplicationRepo repository.ApplicationRepository
	AuthRepo        repository.UserRepository
}

func NewApplicationService(applicationRepo repository.ApplicationRepository, authRepo repository.UserRepository) *ApplicationService {
	return &ApplicationService{
		ApplicationRepo: applicationRepo,
		AuthRepo:        authRepo,
	}
}

func (s ApplicationService) Create(applicationUrl string, applicationCountry string, userID int64) (*domain.Application, *domain.User, error) {
	exists, err := s.ApplicationRepo.FindApplicationByURL(applicationUrl)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	} else if exists != nil {
		log.Println(err)
		return nil, nil, errors.New("This app already exists.")
	}

	user, err := s.AuthRepo.FindUserByID(userID)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	} else if user == nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	}

	newApplication, err := domain.NewApplication(applicationUrl, applicationCountry, user.ID)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	}

	if err := s.ApplicationRepo.CreateApplication(newApplication); err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	}

	return newApplication, user, nil
}
