package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/simple-crud-go/api"
	"github.com/simple-crud-go/internal/models"
	"github.com/simple-crud-go/internal/usecases"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserParams struct {
	Username *string `json:"username"`
	Name     *string `json:"name"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
		usecases.UpdateUserUsername(uint(id), username)
	} else if username == "" {
		usecases.UpdateUserName(uint(id), name)
	} else {
		usecases.UpdateUserFull(uint(id), username, name)
	}

	api.NoDataResponseHandler(w, http.StatusOK, fmt.Sprintf("User with ID=%v successfully updated", id))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var p UserParams
	err := json.NewDecoder(r.Body).Decode(&p)

	if p.Name == nil || p.Username == nil {
		logrus.Error(err)
		api.RequestErrorHandler(w, errors.New("username and name field are required"), 400)
		return
	}

	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	err = usecases.CreateUser(*p.Username, *p.Name)
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	api.NoDataResponseHandler(w, http.StatusCreated, "User successfully created")
}

func Users(w http.ResponseWriter, r *http.Request) {
	user, err := usecases.GetUsers()
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

func UserByUsername(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	if username == "" {
		w.Write([]byte("Username cannot be empty"))
		return
	}

	user, err := usecases.GetUserByUsername(username)
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
