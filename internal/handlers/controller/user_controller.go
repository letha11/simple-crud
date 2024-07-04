package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/simple-crud-go/api"
	"github.com/simple-crud-go/internal/middleware"
	"github.com/simple-crud-go/internal/services"
	"gorm.io/gorm"
)

type UserController struct {
	Service *services.UserService
}

// UpdateUser Update authenticated user
// @summary Update authenticated user
// @description Update authenticated user
// @tags User
// @id update-user
// @accept mpfd
// @produce json
// @param username formData string false "Username"
// @param name formData string false "Name"
// @param password formData string false "Password"
// @success 200 {object} api.NoDataResponse "Success"
// @failure 404 {object} api.ErrorResponse "Not Found"
// @failure 409 {object} api.ErrorResponse "Conflict"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /user/{id} [put]
// @security Bearer
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		username = r.FormValue("username")
		name     = r.FormValue("name")
		password = r.FormValue("password")
		ctx      = r.Context()
		authIdS  = ctx.Value(middleware.UserIdKey).(string)
	)

	if err := r.ParseForm(); err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	authId, err := strconv.Atoi(authIdS)
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	if err := c.Service.UpdateUser(authId, username, name, password); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("User with id %d doesn't exist", authId), 404)
			return
		} else if errors.Is(err, services.ErrUserExist) {
			api.RequestErrorHandler(w, err, http.StatusConflict)
			return
		} else if errors.Is(err, services.ErrMismatchID) {
			api.RequestErrorHandler(w, err, http.StatusUnauthorized)
			return
		} else {
			api.InternalErrorHandler(w, err)
			return
		}
	}

	api.NoDataResponseHandler(w, http.StatusOK, fmt.Sprintf("User with ID=%v successfully updated", authId))
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		username = r.FormValue("username")
		name     = r.FormValue("name")
		password = r.FormValue("password")
	)

	if name == "" || username == "" || password == "" {
		api.RequestErrorHandler(w, errors.New("username, name and password fields are required"), http.StatusBadRequest)
		return
	}

	err := c.Service.CreateUser(username, name, password)
	if err != nil && errors.Is(err, services.ErrUserExist) {
		api.RequestErrorHandler(w, err, http.StatusConflict)
		return
	} else if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	api.NoDataResponseHandler(w, http.StatusCreated, "User successfully created")
}

// Users Get all users
// @summary Get all users
// @description Get all users
// @tags User
// @id get-users
// @produce json
// @success 200 {object} []models.User "Success"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /user [get]
func (c *UserController) Users(w http.ResponseWriter, r *http.Request) {
	user, err := c.Service.GetAllUser()
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	api.GenericResponseHandler(w, http.StatusOK, user)
}

// UserByUsername Get user by username
// @summary Get user by username
// @description Get user by username
// @tags User
// @id get-user-by-username
// @param username path string true "Username"
// @produce json
// @success 200 {object} models.User "Success"
// @failure 404 {object} api.ErrorResponse "Not Found"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /user/{username} [get]
func (c *UserController) UserByUsername(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	if username == "" {
		api.RequestErrorHandler(w, fmt.Errorf("Username cannot be empty"), 400)
		return
	}

	user, err := c.Service.GetUserByUsername(username)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("User with %s not found", username), 404)
			return
		}

		api.InternalErrorHandler(w, err)
	}

	api.GenericResponseHandler(w, http.StatusOK, user)
}

// DeleteUserById Delete authenticated user
// @summary Delete authenticated user
// @description Delete authenticated/logged in user
// @tags User
// @id delete-user-by-id
// @produce json
// @success 200 {object} api.NoDataResponse "Success"
// @failure 404 {object} api.ErrorResponse "Not Found"
// @failure 401 {object} api.ErrorResponse "Unauthorized"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /user [delete]
// @security Bearer
func (c *UserController) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		authIdS = ctx.Value(middleware.UserIdKey).(string)
	)

	authId, err := strconv.Atoi(authIdS)
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	err = c.Service.DeleteUserById(authId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("user with id = %d not found", authId), 404)
			return
		}

		if errors.Is(err, services.ErrMismatchID) {
			api.RequestErrorHandler(w, err, http.StatusUnauthorized)
			return
		}

		api.InternalErrorHandler(w, err)
	}

	api.NoDataResponseHandler(w, http.StatusOK, "Successfully deleted user")
}
