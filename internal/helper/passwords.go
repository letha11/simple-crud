package helper

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordCrypto interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword string, password string) error
}

type BcryptPasswordCrypto struct{}

func (b BcryptPasswordCrypto) HashPassword(password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passHash), nil
}

func (b BcryptPasswordCrypto) ComparePassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
