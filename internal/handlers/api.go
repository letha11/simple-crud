package handlers

import (
	"github.com/gorilla/mux"
	"github.com/simple-crud-go/internal/handlers/controller"
	"github.com/simple-crud-go/internal/repository"
	"gorm.io/gorm"
)

func RouteHandler(r *mux.Router, db *gorm.DB) {
	var (
		postRepository = repository.NewPostRepository(db)
		postController = controller.PostController{Repository: postRepository}
		userRepository = repository.NewUserRepository(db)
		userController = controller.UserController{Repository: userRepository}
	)
	userPrefix := r.PathPrefix("/user").Subrouter()
	userPrefix.HandleFunc("/{username}", userController.UserByUsername).Methods("GET")
	userPrefix.HandleFunc("", userController.Users).Methods("GET")
	userPrefix.HandleFunc("", userController.CreateUser).Methods("POST")
	userPrefix.HandleFunc("/{id}", userController.UpdateUser).Methods("PUT")

	postPrefix := r.PathPrefix("/post").Subrouter()
	postPrefix.HandleFunc("/{id}", postController.GetPostById).Methods("GET")
	postPrefix.HandleFunc("/{authorId}", postController.CreatePost).Methods("POST")
	postPrefix.HandleFunc("/", postController.GetPosts).Methods("GET")
	postPrefix.HandleFunc("/{id}", postController.UpdatePost).Methods("PUT")
}
