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

type PostController struct {
	Repository repository.PostRepo
}

func (c *PostController) GetPostById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	post, err := c.Repository.GetById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("Post with id = %d doesn't exist", id), 404)
			return
		}
	}

	api.GenericResponseHandler(w, http.StatusOK, post)
}

func (c *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := c.Repository.GetPosts()
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	api.GenericResponseHandler(w, http.StatusOK, posts)
}

func (c *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var (
		title = r.FormValue("title")
		body  = r.FormValue("body")
	)

	id, err := strconv.Atoi(mux.Vars(r)["authorId"])
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	if title == "" || body == "" {
		api.RequestErrorHandler(w, errors.New("username and name field are required"), http.StatusBadRequest)
		return
	}

	err = c.Repository.CreatePost(uint(id), title, body)
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	api.NoDataResponseHandler(w, http.StatusCreated, "Post successfully created")
}

func (c *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var (
		title   = r.FormValue("title")
		body    = r.FormValue("body")
		id, err = strconv.Atoi(mux.Vars(r)["id"])
	)

	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	if err = c.Repository.UpdatePost(uint(id), title, body); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("Post with id = %d doesn't exist", id), 404)
			return
		} else {
			api.InternalErrorHandler(w, err)
			return
		}
	}

	api.NoDataResponseHandler(w, http.StatusOK, fmt.Sprintf("Post with id %v successfully updated", id))
}
