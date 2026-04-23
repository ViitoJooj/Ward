package services

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
	"github.com/ViitoJooj/ward/pkg/cryptography"
	"github.com/ViitoJooj/ward/pkg/logger"
)

var (
	ErrUserForbidden = errors.New("forbidden")
	ErrUserNotFound  = errors.New("user not found")
)

type UserService struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

func NewUserService(userRepo repository.UserRepository, log *logger.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   log,
	}
}

func (s *UserService) CreateByAdmin(adminID int, username string, email string, role string, active bool) (*domain.User, string, error) {
	admin, err := s.userRepo.FindUserByID(adminID)
	if err != nil {
		s.logger.Error("failed to find admin user / error: " + err.Error())
		return nil, "", errors.New("internal error")
	}
	if admin == nil || admin.Role != "admin" {
		return nil, "", ErrUserForbidden
	}

	existingUser, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		s.logger.Error("failed to find user by email / error: " + err.Error())
		return nil, "", errors.New("internal error")
	}
	if existingUser != nil {
		return nil, "", errors.New("user already exists")
	}

	if role == "" {
		role = "user"
	}

	temporaryPassword, err := generateTemporaryPassword(12)
	if err != nil {
		s.logger.Error("failed to generate temporary password / error: " + err.Error())
		return nil, "", errors.New("internal error")
	}

	newUser, err := domain.NewUser(username, email, temporaryPassword, active, role)
	if err != nil {
		return nil, "", err
	}

	hashedPassword, err := cryptography.HashPassword(newUser.Password)
	if err != nil {
		s.logger.Error("failed to hash temporary password / error: " + err.Error())
		return nil, "", errors.New("internal error")
	}

	newUser.Password = hashedPassword

	if err := s.userRepo.CreateUser(newUser); err != nil {
		s.logger.Error("failed to create user / error: " + err.Error())
		return nil, "", errors.New("internal error")
	}

	return newUser, temporaryPassword, nil
}

func (s *UserService) GetAll(adminID int) ([]*domain.User, error) {
	admin, err := s.userRepo.FindUserByID(adminID)
	if err != nil {
		s.logger.Error("failed to find admin user / error: " + err.Error())
		return nil, errors.New("internal error")
	}
	if admin == nil || admin.Role != "admin" {
		return nil, ErrUserForbidden
	}

	users, err := s.userRepo.ListUsers()
	if err != nil {
		s.logger.Error("failed to list users / error: " + err.Error())
		return nil, errors.New("internal error")
	}

	return users, nil
}

func (s *UserService) GetByID(adminID int, userID int) (*domain.User, error) {
	admin, err := s.userRepo.FindUserByID(adminID)
	if err != nil {
		s.logger.Error("failed to find admin user / error: " + err.Error())
		return nil, errors.New("internal error")
	}
	if admin == nil || admin.Role != "admin" {
		return nil, ErrUserForbidden
	}

	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		s.logger.Error("failed to find user by id / error: " + err.Error())
		return nil, errors.New("internal error")
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) UpdateByAdmin(adminID int, userID int, username string, email string, password string, role string, active bool) (*domain.User, error) {
	admin, err := s.userRepo.FindUserByID(adminID)
	if err != nil {
		s.logger.Error("failed to find admin user / error: " + err.Error())
		return nil, errors.New("internal error")
	}
	if admin == nil || admin.Role != "admin" {
		return nil, ErrUserForbidden
	}

	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		s.logger.Error("failed to find user by id / error: " + err.Error())
		return nil, errors.New("internal error")
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	validatedUser, err := domain.NewUser(username, email, password, active, role)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := cryptography.HashPassword(validatedUser.Password)
	if err != nil {
		s.logger.Error("failed to hash password / error: " + err.Error())
		return nil, errors.New("internal error")
	}

	user.Username = validatedUser.Username
	user.Email = validatedUser.Email
	user.Password = hashedPassword
	user.Role = validatedUser.Role
	user.Active = validatedUser.Active
	user.Updated_at = time.Now()

	if err := s.userRepo.UpdateUser(user); err != nil {
		s.logger.Error("failed to update user / error: " + err.Error())
		return nil, errors.New("internal error")
	}

	return user, nil
}

func (s *UserService) DeleteByID(adminID int, userID int) error {
	admin, err := s.userRepo.FindUserByID(adminID)
	if err != nil {
		s.logger.Error("failed to find admin user / error: " + err.Error())
		return errors.New("internal error")
	}
	if admin == nil || admin.Role != "admin" {
		return ErrUserForbidden
	}

	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		s.logger.Error("failed to find user by id / error: " + err.Error())
		return errors.New("internal error")
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := s.userRepo.DeleteUserByID(userID); err != nil {
		s.logger.Error("failed to delete user / error: " + err.Error())
		return errors.New("internal error")
	}

	return nil
}

func (s *UserService) UpdateOwnData(userID int, username string, email string, password string) (*domain.User, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		s.logger.Error("failed to find user / error: " + err.Error())
		return nil, errors.New("internal error")
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	validatedUser, err := domain.NewUser(username, email, password, user.Active, user.Role)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := cryptography.HashPassword(validatedUser.Password)
	if err != nil {
		s.logger.Error("failed to hash password / error: " + err.Error())
		return nil, errors.New("internal error")
	}

	user.Username = validatedUser.Username
	user.Email = validatedUser.Email
	user.Password = hashedPassword
	user.Updated_at = time.Now()

	if err := s.userRepo.UpdateUser(user); err != nil {
		s.logger.Error("failed to update own data / error: " + err.Error())
		return nil, errors.New("internal error")
	}

	return user, nil
}

func generateTemporaryPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*"
	max := big.NewInt(int64(len(charset)))
	password := make([]byte, length)

	for i := range password {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		password[i] = charset[n.Int64()]
	}

	return string(password), nil
}
