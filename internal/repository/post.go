package repository

import (
	"github.com/simple-crud-go/internal/models"
	"gorm.io/gorm"
)

type PostRepo interface {
	Create(post *models.Post) error
	Update(post *models.Post) error
	GetById(id int) (*models.Post, error)
	GetAll() ([]models.Post, error)
	Delete(id uint) error
}

func NewPostRepository(db *gorm.DB) *gormPostRepository {
	return &gormPostRepository{
		db: db,
	}
}

type gormPostRepository struct {
	db *gorm.DB
}

func (r *gormPostRepository) GetById(id int) (*models.Post, error) {
	var post models.Post
	// err := r.db.Model(&models.Post{}).Preload("User", func(db *gorm.DB) *gorm.DB {
	// 	return db.Omit("Posts")
	// }).First(&post, id).Error
	err := r.db.Model(&models.Post{}).Preload("User").First(&post, id).Error
	return &post, err
}

func (r *gormPostRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Model(&models.Post{}).Preload("User").Find(&posts).Error

	return posts, err
}

func (r *gormPostRepository) Create(post *models.Post) error {
	return r.db.Create(&post).Error
}

func (r *gormPostRepository) Update(post *models.Post) error {
	return r.db.Save(&post).Error
}

func (r *gormPostRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}
