package repository

import (
	"errors"

	"github.com/simple-crud-go/internal/models"
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
		return err
	}

	if username != "" {
		user.Username = username
	}

	if name != "" {
		user.Name = name
	}

	err = r.db.Save(&user).Error
	if err != nil {
		return err
	}

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
	err := r.db.Where("username = ?", username).Preload("Posts", func(db *gorm.DB) *gorm.DB {
		return db.Omit("title")
	}).First(&user).Error
	return &user, err
}
