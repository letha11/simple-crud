package user

import (
	"errors"

	"github.com/simple-crud-go/internal/database"
	"github.com/simple-crud-go/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var UserExistErr = errors.New("User with the same username already exist")

// type UserExistErr struct{}
//
// func (r *UserExistErr) Error() string {
// 	return "User with the same username already exist"
// }

type Repository struct {
	DB *gorm.DB
}

func updateUser(username string, name string, user models.User) {
	db := database.GetDBGorm()

	// var user models.User
	// db.First(&user, id)
	user.Username = username
	user.Name = name

	db.Save(&user)
}
func (r *Repository) UpdateUserFull(id uint, username string, name string) {
	var user models.User
	r.DB.First(&user, id)

	updateUser(username, name, user)
}
func (r *Repository) UpdateUserUsername(id uint, username string) {

	var user models.User
	r.DB.First(&user, id)

	updateUser(username, user.Name, user)
}
func (r *Repository) UpdateUserName(id uint, name string) {
	var user models.User
	r.DB.First(&user, id)

	updateUser(user.Username, name, user)
}

func (r *Repository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.DB.Omit("posts").Find(&users).Error

	return users, err
}

func (r *Repository) CreateUser(username string, name string) error {
	var err error
	var userUsername models.User
	err = r.DB.Where("username = ?", username).First(&userUsername).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if userUsername.ID != 0 {
		logrus.Println("user dengan username sudah ada")
		return UserExistErr
	}

	user := models.User{
		Username: username,
		Name:     name,
	}

	err = r.DB.Create(&user).Error

	return err
}

func (r *Repository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return user, err
}
