package controller

import (
	"errors"
	"net/http"

	"github.com/simple-crud-go/api"
	"github.com/simple-crud-go/internal/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	Service *services.AuthService
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var (
		username = r.FormValue("username")
		password = r.FormValue("password")
	)

	if username == "" || password == "" {
		api.RequestErrorHandler(w, errors.New("username and password fields are required"), http.StatusBadRequest)
		return
	}

	token, err := c.Service.Login(username, password)
	if err != nil {
		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) && !errors.Is(err, gorm.ErrRecordNotFound) {
			api.InternalErrorHandler(w, err)
			return
		}

		api.RequestErrorHandler(w, errors.New("Username or password are wrong, please try again"), http.StatusUnauthorized)
		return
	}

	api.GenericResponseHandler(w, 200, token)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var (
		name     = r.FormValue("name")
		username = r.FormValue("username")
		password = r.FormValue("password")
	)

	if name == "" || username == "" || password == "" {
		api.RequestErrorHandler(w, errors.New("name, username and password fields are required"), http.StatusBadRequest)
		return
	}

	data, err := c.Service.Register(name, username, password)
	if err != nil {
		if errors.Is(err, services.ErrUserExist) {
			api.RequestErrorHandler(w, err, http.StatusConflict)
			return
		}

		api.InternalErrorHandler(w, err)
		return
	}

	api.GenericResponseHandler(w, 200, data)
}
