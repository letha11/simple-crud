package usecases

import (
	"github.com/simple-crud-go/internal/database"
	"github.com/simple-crud-go/internal/models"
)

func updateUser(username string, name string, user models.User) {
	db := database.GetDBGorm()

	// var user models.User
	// db.First(&user, id)
	user.Username = username
	user.Name = name

	db.Save(&user)
}
func UpdateUserFull(id uint, username string, name string) {
	db := database.GetDBGorm()

	var user models.User
	db.First(&user, id)

	updateUser(username, name, user)
}
func UpdateUserUsername(id uint, username string) {
	db := database.GetDBGorm()

	var user models.User
	db.First(&user, id)

	updateUser(username, user.Name, user)
}
func UpdateUserName(id uint, name string) {
	db := database.GetDBGorm()

	var user models.User
	db.First(&user, id)

	updateUser(user.Username, name, user)
}

func GetUsers() ([]models.User, error) {
	db := database.GetDBGorm()

	var users []models.User
	err := db.Find(&users).Error

	return users, err
}

func CreateUser(username string, name string) error {
	db := database.GetDBGorm()

	user := models.User{
		Username: username,
		Name:     name,
	}

	err := db.Create(&user).Error
	return err
}

func GetUserByUsername(username string) (models.User, error) {
	db := database.GetDBGorm()

	var user models.User
	err := db.Where("username = ?", username).First(&user).Error
	return user, err
}
