package services_test

import (
	"errors"
	"testing"

	// "github.com/simple-crud-go/internal/helper"
	mock_helper "github.com/simple-crud-go/internal/helper/mocks"
	"github.com/simple-crud-go/internal/models"
	mock_repository "github.com/simple-crud-go/internal/repository/mocks"
	"github.com/simple-crud-go/internal/services"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

var errUnexpected = errors.New("unexpected")

func userServiceWithMock(t *testing.T) (*mock_repository.MockUserRepo, *services.UserService, *mock_helper.MockPasswordCrypto) {
	ctrl := gomock.NewController(t)

	userRepoMock := mock_repository.NewMockUserRepo(ctrl)
	passwordCryptoMock := mock_helper.NewMockPasswordCrypto(ctrl)

	service := services.NewUserService(userRepoMock, passwordCryptoMock)

	return userRepoMock, service, passwordCryptoMock
}

func TestGetUserById(t *testing.T) {
	var (
		user = models.User{
			ID:       1,
			Name:     "Ibka",
			Username: "ibkaanhar",
			Password: "dummy",
		}

		userRepoMock, service, _ = userServiceWithMock(t)
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
		user     *models.User
	}{
		{
			"User Not Found",
			func() {
				userRepoMock.EXPECT().GetById(uint(1)).Return(nil, gorm.ErrRecordNotFound).Times(1)
			},
			gorm.ErrRecordNotFound,
			nil,
		},
		{
			"Unknown Error",
			func() {
				userRepoMock.EXPECT().GetById(uint(1)).Return(nil, errUnexpected).Times(1)
			},
			errUnexpected,
			nil,
		},
		{
			"Success",
			func() {
				userRepoMock.EXPECT().GetById(uint(1)).Return(&user, nil).Times(1)
			},
			nil,
			&user,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			u, err := service.GetUserById(1)
			assert.Equal(t, err, c.err)
			assert.Equal(t, u, c.user)
		})
	}

	/// This Code below are the same as above code.
	// // No user with the following id (Record Not Found)
	// userRepoMock.EXPECT().GetById(uint(1)).Return(nil, gorm.ErrRecordNotFound).Times(1)
	//
	// _, err := service.GetUserById(1)
	//
	// assert.Equal(t, err, gorm.ErrRecordNotFound)
	//
	// // Unknown Error
	// userRepoMock.EXPECT().GetById(uint(1)).Return(nil, errUnexpected).Times(1)
	//
	// _, err = service.GetUserById(1)
	//
	// assert.Equal(t, err, errUnexpected)
	//
	// // Success
	// userRepoMock.EXPECT().GetById(uint(1)).Return(&user, nil).Times(1)
	//
	// res, err := service.GetUserById(1)
	//
	// assert.Nil(t, err)
	// assert.Equal(t, res.Username, user.Username)
}

func TestGetUserByUsername(t *testing.T) {
	var (
		username = "sane"
		user     = models.User{
			ID:       1,
			Name:     "Ibka",
			Username: "ibkaanhar",
			Password: "dummy",
		}

		userRepoMock, service, _ = userServiceWithMock(t)
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
		user     *models.User
	}{
		{
			"User Not Found",
			func() {
				userRepoMock.EXPECT().GetByUsername(username).Return(nil, gorm.ErrRecordNotFound).Times(1)
			},
			gorm.ErrRecordNotFound,
			nil,
		},
		{
			"Unknown Error",
			func() {
				userRepoMock.EXPECT().GetByUsername(username).Return(nil, errUnexpected).Times(1)
			},
			errUnexpected,
			nil,
		},
		{
			"Success",
			func() {
				userRepoMock.EXPECT().GetByUsername(username).Return(&user, nil).Times(1)
			},
			nil,
			&user,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			u, err := service.GetUserByUsername(username)
			assert.Equal(t, err, c.err)
			assert.Equal(t, u, c.user)
		})
	}
}

func TestGetAllUser(t *testing.T) {
	var (
		users = []models.User{
			{
				ID:       1,
				Name:     "Ibka",
				Username: "ibkaanhar",
				Password: "dummy",
			},
			{
				ID:       2,
				Name:     "Ibka 2",
				Username: "ibkaanhar2",
				Password: "dummy",
			},
		}

		userRepoMock, service, _ = userServiceWithMock(t)
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
		users    []models.User
	}{
		{
			"Unknown Error",
			func() {
				userRepoMock.EXPECT().GetAll().Return(nil, errUnexpected).Times(1)
			},
			errUnexpected,
			nil,
		},
		{
			"Success",
			func() {
				userRepoMock.EXPECT().GetAll().Return(users, nil).Times(1)
			},
			nil,
			users,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			u, err := service.GetAllUser()
			assert.Equal(t, err, c.err)
			assert.Equal(t, u, c.users)
		})
	}
}

func TestCreateUser(t *testing.T) {
	var (
		foundUser = models.User{
			ID:       1,
			Name:     "Ibka 1",
			Username: "ibkaanhar",
			Password: "dummy",
		}

		hashedPass = "dummy"
		newUser    = models.User{
			Name:     "Ibka",
			Username: "ibkaanhar",
			Password: hashedPass,
		}

		userRepoMock, service, passwordCryptoMock = userServiceWithMock(t)
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
	}{
		{
			"Unknown Error on Getting the username",
			func() {
				userRepoMock.EXPECT().GetByUsername(newUser.Username).Return(nil, errUnexpected).Times(1)
			},
			errUnexpected,
		},
		{
			"User with the same Username exists",
			func() {
				userRepoMock.EXPECT().GetByUsername(newUser.Username).Return(&foundUser, nil).Times(1)
			},
			services.ErrUserExist,
		},
		{
			"Unknown Error when Hasing the password",
			func() {
				userRepoMock.EXPECT().GetByUsername(newUser.Username).Return(&models.User{ID: 0}, gorm.ErrRecordNotFound).Times(1)
				passwordCryptoMock.EXPECT().HashPassword(newUser.Password).Return("", errUnexpected).Times(1)
			},
			errUnexpected,
		},
		{
			"Success",
			func() {
				userRepoMock.EXPECT().GetByUsername(newUser.Username).Return(&models.User{}, gorm.ErrRecordNotFound).Times(1)
				passwordCryptoMock.EXPECT().HashPassword(newUser.Password).Return(hashedPass, nil).Times(1)
				userRepoMock.EXPECT().Create(newUser).Return(nil).Times(1)
			},
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			err := service.CreateUser(newUser.Username, newUser.Name, newUser.Password)
			assert.Equal(t, err, c.err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	var (
		hashedPass       = "dummy"
		userSameUsername = models.User{
			ID:       1,
			Name:     "Ibka",
			Username: "ibkaanhar2",
			Password: "abc",
		}
		userDiffID = models.User{
			ID:       2,
			Name:     "Ibka",
			Username: "ibkaanhar1",
			Password: "abc",
		}
		userDiffUsername = models.User{
			ID:       1,
			Name:     "Ibka",
			Username: "ibkaanhar1",
			Password: "abc",
		}
		newDataUser = models.User{
			ID:       1,
			Name:     "Ibka",
			Username: "ibkaanhar",
			Password: hashedPass,
		}

		userRepoMock, service, passwordCryptoMock = userServiceWithMock(t)
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
	}{
		{
			"User not found and unknown error on GetById",
			func() {
				userRepoMock.EXPECT().GetById(newDataUser.ID).Return(nil, gorm.ErrRecordNotFound).Times(1)
			},
			gorm.ErrRecordNotFound,
		},
		{
			"Logged in user and User to be updated ID doesn't match",
			func() {
				userRepoMock.EXPECT().GetById(newDataUser.ID).Return(&userDiffID, nil).Times(1)
			},
			services.ErrMismatchID,
		},
		{
			"Unknown Error when checking another user with the username",
			func() {
				userRepoMock.EXPECT().GetById(newDataUser.ID).Return(&userDiffUsername, nil).Times(1)
				userRepoMock.EXPECT().GetByUsername(newDataUser.Username).Return(nil, errUnexpected).Times(1)
			},
			errUnexpected,
		},
		{
			"New Username already used by another user",
			func() {
				userRepoMock.EXPECT().GetById(newDataUser.ID).Return(&userDiffUsername, nil).Times(1)
				userRepoMock.EXPECT().GetByUsername(newDataUser.Username).Return(&userSameUsername, nil).Times(1)
			},
			services.ErrUserExist,
		},
		{
			"Unknown error when hashing password",
			func() {
				userRepoMock.EXPECT().GetById(newDataUser.ID).Return(&userDiffUsername, nil).Times(1)
				userRepoMock.EXPECT().GetByUsername(newDataUser.Username).Return(&models.User{}, nil).Times(1)
				passwordCryptoMock.EXPECT().HashPassword(newDataUser.Password).Return("", errUnexpected).Times(1)
			},
			errUnexpected,
		},
		{
			"Success",
			func() {
				userDiffUsername.Username = "ibkaanhar2" // reset
				userRepoMock.EXPECT().GetById(newDataUser.ID).Return(&userDiffUsername, nil).Times(1)
				userRepoMock.EXPECT().GetByUsername(newDataUser.Username).Return(&models.User{}, gorm.ErrRecordNotFound).Times(1)
				passwordCryptoMock.EXPECT().HashPassword(newDataUser.Password).Return(hashedPass, nil).Times(1)
				userRepoMock.EXPECT().Update(newDataUser).Return(nil).Times(1)
			},
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			err := service.UpdateUser(int(newDataUser.ID), newDataUser.Username, newDataUser.Name, newDataUser.Password)
			assert.Equal(t, err, c.err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	var (
		hashedPass   = "dummy"
		existingUser = models.User{
			ID:       1,
			Name:     "Ibka",
			Username: "ibkaanhar",
			Password: hashedPass,
		}
		existingDiffUser = models.User{
			ID:       2,
			Name:     "Ibka",
			Username: "ibkaanhar",
			Password: hashedPass,
		}

		userRepoMock, service, _ = userServiceWithMock(t)
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
	}{
		{
			"Unknown error when getting the user with the passed in id",
			func() {
				userRepoMock.EXPECT().GetById(uint(1)).Return(nil, errUnexpected).Times(1)
			},
			errUnexpected,
		},
		{
			"Logged in user and User to be deleted ID doesn't match",
			func() {
				userRepoMock.EXPECT().GetById(uint(1)).Return(&existingDiffUser, nil).Times(1)
			},
			services.ErrMismatchID,
		},
		{
			"Unknown Error when deleting the user",
			func() {
				userRepoMock.EXPECT().GetById(uint(1)).Return(&existingUser, nil).Times(1)
				userRepoMock.EXPECT().DeleteById(uint(1)).Return(errUnexpected).Times(1)
			},
			errUnexpected,
		},
		{
			"Success",
			func() {
				userRepoMock.EXPECT().GetById(uint(1)).Return(&existingUser, nil).Times(1)
				userRepoMock.EXPECT().DeleteById(uint(1)).Return(nil).Times(1)
			},
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			err := service.DeleteUserById(1)

			assert.Equal(t, err, c.err)
		})
	}
}
