package services

import (
	"errors"

	"github.com/simple-crud-go/api"
	"github.com/simple-crud-go/internal/helper"
	"github.com/simple-crud-go/internal/models"
	"github.com/simple-crud-go/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService struct {
	UserRepository repository.UserRepo
}

func NewAuthService(userRepo repository.UserRepo) *AuthService {
	return &AuthService{
		UserRepository: userRepo,
	}
}

func (s *AuthService) Login(username string, password string) (string, error) {
	user, err := s.UserRepository.GetByUsername(username)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if user == nil || user.ID == 0 {
		logrus.Error("user doesn't exist")
		return "", gorm.ErrRecordNotFound
	}

	if err = helper.ComparePassword(user.Password, password); err != nil {
		logrus.Error(err)
		return "", err
	}

	// token, err := helper.CreateToken(username)
	token, err := helper.CreateToken(int(user.ID))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(name string, username string, password string) (*api.RegisterSuccessResponse, error) {
	user, err := s.UserRepository.GetByUsername(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return nil, err
	}

	if user.ID != 0 {
		logrus.Error("User already exists")
		return nil, ErrUserExist
	}

	hashedPassword, err := helper.HashPassword(password)
	if err != nil {
		logrus.Error(nil)
		return nil, err
	}

	newUser := models.User{
		Name:     name,
		Username: username,
		Password: hashedPassword,
	}

	err = s.UserRepository.Create(newUser)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	user, err = s.UserRepository.GetByUsername(username)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	token, err := helper.CreateToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	data := api.RegisterSuccessResponse{
		Token: token,
		User:  user,
	}

	return &data, nil
}
