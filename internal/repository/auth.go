package repository

type AuthRepo interface {
	Login(username string, password string) error
}
