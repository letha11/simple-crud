package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/simple-crud-go/internal/handlers/controller"
	"github.com/simple-crud-go/internal/helper"
	"github.com/simple-crud-go/internal/middleware"
	"github.com/simple-crud-go/internal/repository"
	"github.com/simple-crud-go/internal/services"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func RouteHandler(r *mux.Router, db *gorm.DB) {
	var (
		bcryptPassCrypto = helper.BcryptPasswordCrypto{}
		jwtHelper        = helper.NewDefaultJWTHelper()

		userRepository = repository.NewUserRepository(db)
		postRepository = repository.NewPostRepository(db)

		userService = services.NewUserService(userRepository, bcryptPassCrypto)
		postService = services.NewPostService(postRepository, userRepository)
		authService = services.NewAuthService(userRepository, bcryptPassCrypto, jwtHelper)

		userController = controller.UserController{Service: userService}
		postController = controller.PostController{Service: postService}
		authController = controller.AuthController{Service: authService}
	)

	r.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	r = r.PathPrefix("/api").Subrouter()

	r.HandleFunc("/login", authController.Login).Methods("POST")
	r.HandleFunc("/register", authController.Register).Methods("POST")

	userPrefix := r.PathPrefix("/user").Subrouter()
	userPrefix.HandleFunc("/{username}", userController.UserByUsername).Methods("GET")
	userPrefix.HandleFunc("", userController.Users).Methods("GET")
	// userPrefix.HandleFunc("", userController.CreateUser).Methods("POST")
	userPrefix.HandleFunc("/{id}", middleware.AuthMiddleware(http.HandlerFunc(userController.UpdateUser)).ServeHTTP).Methods("PUT")
	userPrefix.HandleFunc("", middleware.AuthMiddleware(http.HandlerFunc(userController.DeleteUserById)).ServeHTTP).Methods("DELETE")

	postPrefix := r.PathPrefix("/post").Subrouter()
	postPrefix.HandleFunc("", postController.GetPosts).Methods("GET")
	postPrefix.HandleFunc("/{id}", postController.GetPostById).Methods("GET")
	postPrefix.HandleFunc("", middleware.AuthMiddleware(http.HandlerFunc(postController.CreatePost)).ServeHTTP).Methods("POST")
	postPrefix.HandleFunc("/{id}", middleware.AuthMiddleware(http.HandlerFunc(postController.UpdatePost)).ServeHTTP).Methods("PUT")
	postPrefix.HandleFunc("/{id}", middleware.AuthMiddleware(http.HandlerFunc(postController.DeletePostById)).ServeHTTP).Methods("DELETE")

}
