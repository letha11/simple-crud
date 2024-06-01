package user

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

func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{
		db: db,
	}
}

type GormUserRepository struct {
	db *gorm.DB
}

func (r *GormUserRepository) UpdateUser(id int, username string, name string) error {
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

func (r *GormUserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Omit("posts").Find(&users).Error

	return users, err
}

func (r *GormUserRepository) CreateUser(username string, name string) error {
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

func (r *GormUserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}
