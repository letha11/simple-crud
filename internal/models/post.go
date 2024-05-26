package models

type Post struct {
	Model
	UserID uint
	User   User
}
