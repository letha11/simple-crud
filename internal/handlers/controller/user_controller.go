package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/simple-crud-go/api"
	"github.com/simple-crud-go/internal/repository"
	"gorm.io/gorm"
)

type UserController struct {
	Repository repository.UserRepo
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		username = r.FormValue("username")
		name     = r.FormValue("name")
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

	if err := c.Repository.UpdateUser(id, username, name); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("User with id %d doesn't exist", id), 404)
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
	)

	if name == "" || username == "" {
		api.RequestErrorHandler(w, errors.New("username and name field are required"), http.StatusBadRequest)
		return
	}

	err := c.Repository.CreateUser(username, name)
	if err != nil && errors.Is(err, repository.ErrUserExist) {
		api.RequestErrorHandler(w, err, http.StatusConflict)
		return
	} else if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	api.NoDataResponseHandler(w, http.StatusCreated, "User successfully created")
}

func (c *UserController) Users(w http.ResponseWriter, r *http.Request) {
	user, err := c.Repository.GetUsers()
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

	user, err := c.Repository.GetByUsername(username)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("User with %s not found", username), 404)
			return
		}
	}

	api.GenericResponseHandler(w, http.StatusOK, user)
}
