package helper_test

import (
	"testing"

	"github.com/simple-crud-go/internal/helper"
)

func TestBcryptPasswordCrypto(t *testing.T) {
	passwordCrypto := helper.BcryptPasswordCrypto{}

	t.Run("HashPassword", func(t *testing.T) {
		cases := []struct {
			name     string
			password string
		}{
			{"regular password", "R3GuLarPsswd"},
			{"empty password", ""},
			{"long password", "ajdslkaaslkdjsaldj092103lkasjd02913lsadjlksajd"},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				hashedPassword, err := passwordCrypto.HashPassword(c.password)

				if err != nil {
					t.Fatal(err.Error())
				}

				if hashedPassword == c.password {
					t.Fatal("HashPassword doesn't hash the password correctly")
				}

				if hashedPassword == "" {
					t.Fatal("HashPassword return empty hashed password")
				}

			})
		}
	})

	t.Run("ComparePassword", func(t *testing.T) {
		hashedPassword, _ := passwordCrypto.HashPassword("123")

		cases := []struct {
			name       string
			password   string
			shouldFail bool
		}{
			{"correct password", "123", false},
			{"wrong password", "1234", true},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				err := passwordCrypto.ComparePassword(hashedPassword, c.password)

				if err != nil && !c.shouldFail {
					t.Fatal(err.Error())
				}
			})
		}
	})

}
