package services

import (
	"errors"

	"github.com/simple-crud-go/internal/models"
	"github.com/simple-crud-go/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostService struct {
	PostRepository repository.PostRepo
	UserRepository repository.UserRepo
}

func NewPostService(postRepo repository.PostRepo, userRepo repository.UserRepo) *PostService {
	return &PostService{
		PostRepository: postRepo,
		UserRepository: userRepo,
	}
}

func (s *PostService) GetPostById(id int) (*models.Post, error) {
	user, err := s.PostRepository.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithField("id", id).Error("Post doesn't exist")
		} else {
			logrus.Error(err)
		}
		return nil, err
	}
	return user, nil
}

func (s *PostService) GetAllPost() ([]models.Post, error) {
	users, err := s.PostRepository.GetAll()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return users, nil
}

func (s *PostService) CreatePost(authorId int, title string, body string) error {
	_, err := s.UserRepository.GetById(uint(authorId))
	if err != nil {
		return err
	}

	post := models.Post{
		UserID: uint(authorId),
		Title:  title,
		Body:   body,
	}

	err = s.PostRepository.Create(&post)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (s *PostService) UpdatePost(postId int, title string, body string) error {
	post, err := s.PostRepository.GetById(postId)
	if err != nil {
		return err
	}

	if title != "" {
		post.Title = title
	}

	if body != "" {
		post.Body = body
	}

	err = s.PostRepository.Update(post)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
