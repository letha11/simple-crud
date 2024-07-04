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

// GetPostById Get post by id
// @summary Get post by id
// @description Get post by id
// @tags Post
// @id get-post-by-id
// @produce json
// @param id path int true "Post ID"
// @success 200 {object} api.GenericSuccessResponse[models.Post] "Success"
// @failure 404 {object} api.ErrorResponse "Not Found"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /post/{id} [get]
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

// GetPosts Get all posts
// @summary Get all posts
// @description Get all posts
// @tags Post
// @id get-all-posts
// @produce json
// @success 200 {object} api.GenericSuccessResponse[[]models.Post] "Success"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /post [get]
func (c *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := c.Service.GetAllPost()
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	api.GenericResponseHandler(w, http.StatusOK, posts)
}

// CreatePost Create a post
// @summary Create a post
// @description Create a post
// @tags Post
// @id create-post
// @accept mpfd
// @produce json
// @param title formData string true "Post Title"
// @param body formData string true "Post Body"
// @success 200 {object} api.NoDataResponse "Post created"
// @failure 400 {object} api.ErrorResponse "Conflict"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /post [post]
// @security Bearer
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

// UpdatePost Update a posted post
// @summary Update a posted post
// @description Update a posted post
// @tags Post
// @id update-post
// @accept mpfd
// @produce json
// @param id path int true "Post ID"
// @param title formData string false "Post Title"
// @param body formData string false "Post Body"
// @success 200 {object} api.NoDataResponse "Post updated"
// @failure 404 {object} api.ErrorResponse "Not Found"
// @failure 401 {object} api.ErrorResponse "Unauthorized"
// @failure 400 {object} api.ErrorResponse "Conflict"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /post/{id} [put]
// @security Bearer
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

// DeletePostById Delete a post by id
// @summary Delete a post by id
// @description Delete a post by id
// @tags Post
// @id delete-post
// @produce json
// @param id path int true "Post ID"
// @success 200 {object} api.NoDataResponse "Post deleted"
// @failure 404 {object} api.ErrorResponse "Not Found"
// @failure 401 {object} api.ErrorResponse "Unauthorized"
// @failure 500 {object} api.ErrorResponse "Internal Server Error"
// @router /post/{id} [delete]
// @security Bearer
func (c *PostController) DeletePostById(w http.ResponseWriter, r *http.Request) {
	var (
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

	if err = c.Service.DeletePostById(authId, id); err != nil {
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

	api.NoDataResponseHandler(w, http.StatusOK, fmt.Sprintf("Post with id %v successfully deleted", id))
}
