package repository

import (
	"github.com/simple-crud-go/internal/models"
	"gorm.io/gorm"
)

type UserRepo interface {
	Update(user models.User) error
	Create(user models.User) error
	GetById(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAll() ([]models.User, error)
	DeleteById(id uint) error
}

func NewUserRepository(db *gorm.DB) *gormUserRepository {
	return &gormUserRepository{
		db: db,
	}
}

type gormUserRepository struct {
	db *gorm.DB
}

func (r *gormUserRepository) Update(user models.User) error {
	return r.db.Save(&user).Error
}

func (r *gormUserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Omit("posts").Find(&users).Error

	return users, err
}

func (r *gormUserRepository) Create(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *gormUserRepository) GetByUsername(username string) (*models.User, error) {
	var user *models.User
	err := r.db.Where("username = ?", username).Preload("Posts").First(&user).Error
	return user, err
}

func (r *gormUserRepository) GetById(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Posts").First(&user, id).Error
	return &user, err
}

func (r *gormUserRepository) DeleteById(id uint) error {
	err := r.db.Delete(&models.User{}, id).Error
	return err
}
