package helper

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/simple-crud-go/internal/configs"
)

// func CreateToken(username string) (string, error) {
func CreateToken(id int) (string, error) {
	idString := strconv.Itoa(id)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		// "aud": username,
		"aud": idString,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * (7 * 24)).Unix(),
	})

	signedToken, err := token.SignedString([]byte(configs.GetJWTSecret()))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func CheckToken(token string) error {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(configs.GetJWTSecret()), nil
	})

	if err != nil {
		return err
	}

	if !t.Valid {
		return jwt.ErrInvalidKey
	}

	return nil
}

func ExtractAudienceToken(token string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(configs.GetJWTSecret()), nil
	})

	if err != nil {
		return "", err
	}

	if !t.Valid {
		return "", jwt.ErrInvalidKey
	}

	claims, err := t.Claims.GetAudience()

	if err != nil {
		return "", err
	}

	return claims[0], nil
}
