package repository

import (
	"errors"

	"github.com/simple-crud-go/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var ErrUserExist = errors.New("User with the same username already exist")

type UserRepo interface {
	UpdateUser(id int, username string, name string) error
	CreateUser(username string, name string) error
	GetByUsername(username string) (*models.User, error)
	GetUsers() ([]models.User, error)
}

func NewUserRepository(db *gorm.DB) *gormUserRepository {
	return &gormUserRepository{
		db: db,
	}
}

type gormUserRepository struct {
	db *gorm.DB
}

func (r *gormUserRepository) UpdateUser(id int, username string, name string) error {
	var user models.User
	err := r.db.First(&user, id).Error

	if err != nil {
		logrus.Error(err)
		return err
	}

	if username != "" {
		user.Username = username
	}

	if name != "" {
		user.Name = name
	}

	r.db.Save(&user)

	return nil
}

func (r *gormUserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Omit("posts").Find(&users).Error

	return users, err
}

func (r *gormUserRepository) CreateUser(username string, name string) error {
	var err error
	var userUsername models.User
	err = r.db.Where("username = ?", username).First(&userUsername).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if userUsername.ID != 0 {
		logrus.Println("user dengan username sudah ada")
		return ErrUserExist
	}

	user := models.User{
		Username: username,
		Name:     name,
	}

	err = r.db.Create(&user).Error

	return err
}

func (r *gormUserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
