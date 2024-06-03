package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/simple-crud-go/api"
	"github.com/simple-crud-go/internal/services"
	"gorm.io/gorm"
)

type UserController struct {
	Service *services.UserService
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		username = r.FormValue("username")
		name     = r.FormValue("name")
		password = r.FormValue("password")
		id, err  = strconv.Atoi(mux.Vars(r)["id"])
	)

	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	if err := c.Service.UpdateUser(id, username, name, password); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("User with id %d doesn't exist", id), 404)
			return
		} else if errors.Is(err, services.ErrUserExist) {
			api.RequestErrorHandler(w, err, http.StatusConflict)
			return
		} else {
			api.InternalErrorHandler(w, err)
			return
		}
	}

	api.NoDataResponseHandler(w, http.StatusOK, fmt.Sprintf("User with ID=%v successfully updated", id))
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

func (c *UserController) Users(w http.ResponseWriter, r *http.Request) {
	user, err := c.Service.GetAllUser()
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	api.GenericResponseHandler(w, http.StatusOK, user)
}

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
	}

	api.GenericResponseHandler(w, http.StatusOK, user)
}
