package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/simple-crud-go/api"
	"github.com/simple-crud-go/internal/models"
	"github.com/simple-crud-go/internal/repository/user"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserController struct {
	Repository user.Repository
}

type UserParams struct {
	Username *string `json:"username"`
	Name     *string `json:"name"`
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}
	username := r.FormValue("username")
	name := r.FormValue("name")

	if err := r.ParseForm(); err != nil {
		logrus.Error(err)
		api.InternalErrorHandler(w, err)
		return
	}

	if name == "" {
		c.Repository.UpdateUserUsername(uint(id), username)
	} else if username == "" {
		c.Repository.UpdateUserName(uint(id), name)
	} else {
		c.Repository.UpdateUserFull(uint(id), username, name)
	}

	api.NoDataResponseHandler(w, http.StatusOK, fmt.Sprintf("User with ID=%v successfully updated", id))
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var p UserParams
	err := json.NewDecoder(r.Body).Decode(&p)

	if p.Name == nil || p.Username == nil {
		logrus.Error(err)
		api.RequestErrorHandler(w, errors.New("username and name field are required"), http.StatusBadRequest)
		return
	}

	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	err = c.Repository.CreateUser(*p.Username, *p.Name)
	if err != nil && errors.Is(err, user.UserExistErr) {
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
		logrus.Debug(err)
		return
	}

	resp := api.GenericSuccessReponse[[]models.User]{
		StatusCode: http.StatusOK,
		Data:       user,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (c *UserController) UserByUsername(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	if username == "" {
		w.Write([]byte("Username cannot be empty"))
		return
	}

	user, err := c.Repository.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.Write([]byte(fmt.Sprintf("User with %s not found", username)))

			return
		}
	}

	resp := api.GenericSuccessReponse[models.User]{
		StatusCode: http.StatusOK,
		Data:       user,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
