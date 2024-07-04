package services_test

import (
	"testing"

	"github.com/simple-crud-go/internal/models"
	mock_repository "github.com/simple-crud-go/internal/repository/mocks"
	"github.com/simple-crud-go/internal/services"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func postServiceWithMock(t *testing.T) (*mock_repository.MockPostRepo, *mock_repository.MockUserRepo, *services.PostService) {
	ctrl := gomock.NewController(t)

	userRepoMock := mock_repository.NewMockUserRepo(ctrl)
	postRepoMock := mock_repository.NewMockPostRepo(ctrl)

	service := services.NewPostService(postRepoMock, userRepoMock)

	return postRepoMock, userRepoMock, service
}

func TestGetPostById(t *testing.T) {
	var (
		id                   = 1
		postRepo, _, service = postServiceWithMock(t)
		foundPost            = models.Post{
			ID:     1,
			Title:  "dummy title",
			Body:   "dummy body",
			UserID: 2,
		}
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
		post     *models.Post
	}{
		{
			"Post not found",
			func() {
				postRepo.EXPECT().GetById(id).Return(nil, gorm.ErrRecordNotFound).Times(1)
			},
			gorm.ErrRecordNotFound,
			nil,
		},
		{
			"Unexpected Error",
			func() {
				postRepo.EXPECT().GetById(id).Return(nil, errUnexpected).Times(1)
			},
			errUnexpected,
			nil,
		},
		{
			"Success",
			func() {
				postRepo.EXPECT().GetById(id).Return(&foundPost, nil).Times(1)
			},
			nil,
			&foundPost,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			p, err := service.GetPostById(id)

			assert.Equal(t, err, c.err)
			assert.Equal(t, p, c.post)
		})
	}
}

func TestGetAllPost(t *testing.T) {
	var (
		postRepo, _, service = postServiceWithMock(t)
		posts                = []models.Post{
			{
				ID:     1,
				Title:  "dummy title",
				Body:   "dummy body",
				UserID: 2,
			},
			{
				ID:     2,
				Title:  "dummy title ",
				Body:   "dummy body 2",
				UserID: 3,
			},
		}
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
		post     []models.Post
	}{
		{
			"Unexpected Error",
			func() {
				postRepo.EXPECT().GetAll().Return(nil, errUnexpected).Times(1)
			},
			errUnexpected,
			nil,
		},
		{
			"Success",
			func() {
				postRepo.EXPECT().GetAll().Return(posts, nil).Times(1)
			},
			nil,
			posts,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			p, err := service.GetAllPost()

			assert.Equal(t, err, c.err)
			assert.Equal(t, p, c.post)
		})
	}
}

func TestCreatePost(t *testing.T) {
	var (
		postRepo, userRepo, service = postServiceWithMock(t)
		author                      = models.User{
			ID:       2,
			Name:     "Ibka",
			Username: "ibkaanhar",
			Password: "123",
		}
		newPost = models.Post{
			Title:  "dummy title",
			Body:   "dummy body",
			UserID: 2,
		}
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
		post     *models.Post
	}{
		{
			"Unknown Error on finding associated user",
			func() {
				userRepo.EXPECT().GetById(author.ID).Return(nil, errUnexpected)
			},
			errUnexpected,
			nil,
		},
		{
			"Unknown Error on creating the Post",
			func() {
				userRepo.EXPECT().GetById(author.ID).Return(&author, nil)
				postRepo.EXPECT().Create(&newPost).Return(errUnexpected)
			},
			errUnexpected,
			nil,
		},
		{
			"Success",
			func() {
				userRepo.EXPECT().GetById(author.ID).Return(&author, nil)
				postRepo.EXPECT().Create(&newPost).Return(nil)
			},
			nil,
			&newPost,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			err := service.CreatePost(int(newPost.UserID), newPost.Title, newPost.Body)

			assert.Equal(t, err, c.err)
		})
	}
}

func TestUpdatePost(t *testing.T) {
	var (
		postRepo, _, service = postServiceWithMock(t)
		loggedInUser         = models.User{
			ID:       2,
			Name:     "Ibka",
			Username: "ibkaanhar",
			Password: "123",
		}
		// diffUser = models.User{
		// 	ID:       3,
		// 	Name:     "Ibka",
		// 	Username: "ibkaanhar",
		// 	Password: "123",
		// }
		newPost = models.Post{
			ID:     2,
			Title:  "dummy title",
			Body:   "dummy body",
			UserID: 2,
		}
		postDiffUser = models.Post{
			ID:     3,
			Title:  "dummy title 2",
			Body:   "dummy body 2",
			UserID: 3,
		}
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
		post     *models.Post
	}{
		{
			"User Not Found and Unexpected Error on finding associated user",
			func() {
				postRepo.EXPECT().GetById(int(newPost.ID)).Return(nil, errUnexpected)
			},
			errUnexpected,
			nil,
		},
		{
			"Current Logged in User ID mismatch with Post.UserID",
			func() {
				postRepo.EXPECT().GetById(int(newPost.ID)).Return(&postDiffUser, nil)
			},
			services.ErrMismatchAuthorID,
			nil,
		},
		{
			"Unexpected Error when updating post",
			func() {
				postRepo.EXPECT().GetById(int(newPost.ID)).Return(&newPost, nil)
				postRepo.EXPECT().Update(&newPost).Return(errUnexpected)
			},
			errUnexpected,
			nil,
		},
		{
			"Success",
			func() {
				postRepo.EXPECT().GetById(int(newPost.ID)).Return(&newPost, nil)
				postRepo.EXPECT().Update(&newPost).Return(nil)
			},
			nil,
			&newPost,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			err := service.UpdatePost(int(loggedInUser.ID), int(newPost.ID), newPost.Title, newPost.Body)

			assert.Equal(t, err, c.err)
		})
	}
}

func TestDeletePost(t *testing.T) {
	var (
		postRepo, _, service = postServiceWithMock(t)
		loggedInUser         = models.User{
			ID:       2,
			Name:     "Ibka",
			Username: "ibkaanhar",
			Password: "123",
		}
		toBeDeletedPost = models.Post{
			ID:     4,
			Title:  "dummy title 2",
			Body:   "dummy body 2",
			UserID: 2,
		}
		diffUserAndPost = models.Post{
			ID:     3,
			Title:  "dummy title 2",
			Body:   "dummy body 2",
			UserID: 3,
		}
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
	}{
		{
			"Unknown Error when getting post",
			func() {
				postRepo.EXPECT().GetById(int(toBeDeletedPost.ID)).Return(nil, errUnexpected)
			},
			errUnexpected,
		},
		{
			"Logged in user ID mismatch with the post Author/UserID to be deleted",
			func() {
				postRepo.EXPECT().GetById(int(toBeDeletedPost.ID)).Return(&diffUserAndPost, nil)
			},
			services.ErrMismatchAuthorID,
		},
		{
			"Unknown Error when deleting the post",
			func() {
				postRepo.EXPECT().GetById(int(toBeDeletedPost.ID)).Return(&toBeDeletedPost, nil)
				postRepo.EXPECT().Delete(toBeDeletedPost.ID).Return(errUnexpected)
			},
			errUnexpected,
		},
		{
			"Success",
			func() {
				postRepo.EXPECT().GetById(int(toBeDeletedPost.ID)).Return(&toBeDeletedPost, nil)
				postRepo.EXPECT().Delete(toBeDeletedPost.ID).Return(nil)
			},
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			err := service.DeletePostById(int(loggedInUser.ID), int(toBeDeletedPost.ID))

			assert.Equal(t, err, c.err)
		})
	}
}
