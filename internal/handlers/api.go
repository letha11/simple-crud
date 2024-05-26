package handlers

import (
	"github.com/gorilla/mux"
	controller "github.com/simple-crud-go/internal/handlers/controller"
)

func RouteHandler(r *mux.Router) {

	userPrefix := r.PathPrefix("/user").Subrouter()
	userPrefix.HandleFunc("/{username}", controller.UserByUsername).Methods("GET")
	userPrefix.HandleFunc("", controller.Users).Methods("GET")
	userPrefix.HandleFunc("", controller.CreateUser).Methods("POST")
	userPrefix.HandleFunc("/{id}", controller.UpdateUser).Methods("PUT")
}
