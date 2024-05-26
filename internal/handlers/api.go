package handlers

import (
	"github.com/gorilla/mux"
	controller "github.com/simple-crud-go/internal/handlers/controller"
	"github.com/simple-crud-go/internal/repository/user"
	"gorm.io/gorm"
)

func RouteHandler(r *mux.Router, db *gorm.DB) {
	userRepository := user.Repository{DB: db}
	userController := controller.UserController{Repository: userRepository}

	userPrefix := r.PathPrefix("/user").Subrouter()
	userPrefix.HandleFunc("/{username}", userController.UserByUsername).Methods("GET")
	userPrefix.HandleFunc("", userController.Users).Methods("GET")
	userPrefix.HandleFunc("", userController.CreateUser).Methods("POST")
	userPrefix.HandleFunc("/{id}", userController.UpdateUser).Methods("PUT")
}
