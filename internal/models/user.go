package models

type User struct {
	Model
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Posts    *[]Post `json:"posts"`
}
