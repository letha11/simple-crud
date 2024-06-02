package handlers

import (
	"github.com/gorilla/mux"
	"github.com/simple-crud-go/internal/handlers/controller"
	"github.com/simple-crud-go/internal/repository"
	"github.com/simple-crud-go/internal/services"
	"gorm.io/gorm"
)

func RouteHandler(r *mux.Router, db *gorm.DB) {
	var (
		userRepository = repository.NewUserRepository(db)
		userService    = services.NewUserService(userRepository)
		userController = controller.UserController{Service: userService}
		postRepository = repository.NewPostRepository(db)
		postService    = services.NewPostService(postRepository, userRepository)
		postController = controller.PostController{Service: postService}
	)
	userPrefix := r.PathPrefix("/user").Subrouter()
	userPrefix.HandleFunc("/{username}", userController.UserByUsername).Methods("GET")
	userPrefix.HandleFunc("", userController.Users).Methods("GET")
	userPrefix.HandleFunc("", userController.CreateUser).Methods("POST")
	userPrefix.HandleFunc("/{id}", userController.UpdateUser).Methods("PUT")

	postPrefix := r.PathPrefix("/post").Subrouter()
	postPrefix.HandleFunc("/{id}", postController.GetPostById).Methods("GET")
	// FIXME update when proper authentication works
	postPrefix.HandleFunc("/{authorId}", postController.CreatePost).Methods("POST")
	postPrefix.HandleFunc("", postController.GetPosts).Methods("GET")
	postPrefix.HandleFunc("/{id}", postController.UpdatePost).Methods("PUT")
}
