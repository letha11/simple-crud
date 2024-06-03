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
		title     = r.FormValue("title")
		body      = r.FormValue("body")
		ctx       = r.Context()
		authorIdS = ctx.Value(middleware.UserIdKey).(string)
	)

	authorId, err := strconv.Atoi(authorIdS)
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	if title == "" || body == "" {
		api.RequestErrorHandler(w, errors.New("title and body field are required"), http.StatusBadRequest)
		return
	}

	if err := c.Service.CreatePost(authorId, title, body); err != nil {
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
		ctx     = r.Context()
		authIdS = ctx.Value(middleware.UserIdKey).(string)
	)

	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	authId, err := strconv.Atoi(authIdS)
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	if err = c.Service.UpdatePost(authId, id, title, body); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.RequestErrorHandler(w, fmt.Errorf("Post with id = %d doesn't exist", id), 404)
			return
		} else if errors.Is(err, services.ErrMismatchAuthorID) {
			api.RequestErrorHandler(w, err, http.StatusUnauthorized)
			return
		} else {
			api.InternalErrorHandler(w, err)
			return
		}
	}

	api.NoDataResponseHandler(w, http.StatusOK, fmt.Sprintf("Post with id %v successfully updated", id))
}
