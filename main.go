package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/simple-crud-go/docs"
	"github.com/simple-crud-go/internal/configs"
	"github.com/simple-crud-go/internal/database"
	"github.com/simple-crud-go/internal/handlers"
	"github.com/simple-crud-go/internal/middleware"
	"github.com/simple-crud-go/internal/models"
	"github.com/sirupsen/logrus"
)

// @title Simple CRUD & Authentication
// @version 1.0
// @description This is a learning project, the purpose of this API are just me getting familiar with the language and learn about how to build an REST API in Golang

// @contact.name Ibka Anhar Fatcha
// @contact.url https://github.com/letha11
// @contact.email ibkaanhar1@gmail.com

// @host localhost:5000
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT Token
func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Error(fmt.Sprintf("Error loading .env file, error: %v", err))
	}

	db := database.InitDB()
	err = db.AutoMigrate(&models.User{}, &models.Post{})
	if err != nil {
		panic("failed to migrate")
	}

	r := mux.NewRouter()
	// r.Use(middleware.OnlyJson)
	handlers.RouteHandler(r, db)

	port := configs.GetPort()
	fmt.Printf("Listening at http://127.0.0.1:%v\n", port)
	err = http.ListenAndServe(":"+port, middleware.TrailingSlashMiddleware(r))

	if err != nil {
		logrus.Error(err)
	}
}
