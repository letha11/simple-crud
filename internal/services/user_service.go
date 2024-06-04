package services

import (
	"errors"

	"github.com/simple-crud-go/internal/helper"
	"github.com/simple-crud-go/internal/models"
	"github.com/simple-crud-go/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var ErrUserExist = errors.New("User with the same username already exist")

type UserService struct {
	UserRepository repository.UserRepo
	PasswordCrypto helper.PasswordCrypto
}

func NewUserService(userRepo repository.UserRepo, passwordCrypto helper.PasswordCrypto) *UserService {
	return &UserService{
		UserRepository: userRepo,
		PasswordCrypto: passwordCrypto,
	}
}

func (s *UserService) GetUserById(id int) (*models.User, error) {
	user, err := s.UserRepository.GetById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithField("id", id).Error("User doesn't exist")
		} else {
			logrus.Error(err)
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.UserRepository.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithField("username", username).Error("User doesn't exist")
		} else {
			logrus.Error(err)
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAllUser() ([]models.User, error) {
	user, err := s.UserRepository.GetAll()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(username string, name string, password string) error {
	var err error
	user, err := s.UserRepository.GetByUsername(username)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return err
	}

	if user.ID != 0 {
		logrus.WithField("username", username).Error("User already exists")
		return ErrUserExist
	}

	hashedPass, err := s.PasswordCrypto.HashPassword(password)
	if err != nil {
		return err
	}

	newUser := models.User{
		Username: username,
		Name:     name,
		Password: hashedPass,
	}

	return s.UserRepository.Create(newUser)
}

func (s *UserService) UpdateUser(id int, username string, name string, password string) error {
	user, err := s.UserRepository.GetById(uint(id))
	if err != nil {
		logrus.Error(err)
		return err
	}

	if username != "" {
		if user.Username != username {
			userWithUsername, err := s.UserRepository.GetByUsername(username)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.Error(err)
				return err
			}

			if userWithUsername.ID != 0 {
				return ErrUserExist
			}
		}

		user.Username = username
	}

	if name != "" {
		user.Name = name
	}

	if password != "" {
		hashed, err := s.PasswordCrypto.HashPassword(password)
		if err != nil {
			logrus.Error(err)
			return err
		}

		user.Password = hashed
	}

	return s.UserRepository.Update(*user)
}
