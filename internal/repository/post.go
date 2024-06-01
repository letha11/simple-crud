package repository

import (
	"github.com/simple-crud-go/internal/models"
	"gorm.io/gorm"
)

type PostRepo interface {
	CreatePost(authorId uint, title string, body string) error
	UpdatePost(postId uint, title string, body string) error
	GetById(id int) (*models.Post, error)
	GetPosts() ([]models.Post, error)
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
	err := r.db.Model(&models.Post{}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("Posts")
	}).First(&post, id).Error

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *gormPostRepository) GetPosts() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// CreatePost(authorId uint, title string, body string) error
func (r *gormPostRepository) CreatePost(authorId uint, title string, body string) error {
	var user models.User
	err := r.db.Find(&user, authorId).Error

	if err != nil {
		return err
	}

	post := models.Post{
		UserID: authorId,
		Title:  title,
		Body:   body,
	}

	err = r.db.Create(&post).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *gormPostRepository) UpdatePost(postId uint, title string, body string) error {
	var post models.Post
	err := r.db.First(&post, postId).Error

	if err != nil {
		return err
	}

	if title != "" {
		post.Title = title
	}

	if body != "" {
		post.Body = body
	}

	err = r.db.Save(&post).Error
	if err != nil {
		return err
	}

	return nil
}
