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

type PostController struct {
	Service *services.PostService
}

func (c *PostController) GetPostById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	post, err := c.Service.GetPostById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("Post with id = %d doesn't exist", id), 404)
			return
		}
	}

	api.GenericResponseHandler(w, http.StatusOK, post)
}

func (c *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := c.Service.GetAllPost()
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

	// FIXME it should not been hardcoded authorId (update when proper authentication works)
	id, err := strconv.Atoi(mux.Vars(r)["authorId"])
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	if title == "" || body == "" {
		api.RequestErrorHandler(w, errors.New("title and body field are required"), http.StatusBadRequest)
		return
	}

	if err = c.Service.CreatePost(id, title, body); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("User with id %v doesn't exists", id), 404)
		}
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

	if err = c.Service.UpdatePost(id, title, body); err != nil {
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
