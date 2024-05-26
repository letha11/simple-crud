package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/simple-crud-go/internal/configs"
	"github.com/simple-crud-go/internal/database"
	"github.com/simple-crud-go/internal/handlers"
	"github.com/simple-crud-go/internal/middleware"
	"github.com/simple-crud-go/internal/models"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Error(fmt.Sprintf("Error loading .env file, error: %v", err))
	}

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(middleware.OnlyJson)
	handlers.RouteHandler(r)

	db := database.InitDB()

	err = db.AutoMigrate(&models.User{}, &models.Post{})

	if err != nil {
		panic("failed to migrate")
	}

	port := configs.GetPort()
	fmt.Printf("Listening at http://127.0.0.1:%v\n", port)
	err = http.ListenAndServe(":"+port, middleware.TrailingSlashMiddleware(r))

	if err != nil {
		logrus.Error(err)
	}
}
