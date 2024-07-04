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

// Login Log in the user
// @summary Log in the user
// @description Log in the user
// @tags Authentication
// @id login
// @accept mpfd
// @produce json
// @param username formData string true "Username"
// @param password formData string true "Password"
// @success 200 {object} api.GenericSuccessResponse[string] "JWT Token"
// @failure 400 {object} api.ErrorResponse "Bad Request"
// @failure 401 {object} api.ErrorResponse "Unauthorized"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /login [post]
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

// Register Register a new user
// @summary Register a new user
// @description Register a new user
// @tags Authentication
// @id register
// @accept mpfd
// @produce json
// @param name formData string true "Name"
// @param username formData string true "Username"
// @param password formData string true "Password"
// @success 200 {object} api.GenericSuccessResponse[api.RegisterSuccessResponse] "User Registered"
// @failure 400 {object} api.ErrorResponse "Bad Request"
// @failure 409 {object} api.ErrorResponse "Conflict"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /register [post]
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
