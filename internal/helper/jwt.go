package helper

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/simple-crud-go/internal/configs"
)

//go:generate mockgen -destination=./mocks/jwt.go -source=./jwt.go
type JWTHelper interface {
	CreateToken(id int) (string, error)
	CheckToken(token string) error
	ExtractAudienceToken(token string) (string, error)
}

type jwtHelper struct {
	Manager JWTManager
}

func NewJWTHelper(jwtManager JWTManager) JWTHelper {
	return jwtHelper{
		Manager: jwtManager,
	}
}

func NewDefaultJWTHelper() JWTHelper {
	return jwtHelper{
		Manager: NewDefaultJWTManager(),
	}
}

func (j jwtHelper) CreateToken(id int) (string, error) {
	idString := strconv.Itoa(id)

	signedToken, err := j.Manager.SignToken(idString)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j jwtHelper) CheckToken(token string) error {
	t, err := j.Manager.ParseToken(token)

	if err != nil {
		return err
	}

	if !t.Valid {
		return jwt.ErrInvalidKey
	}

	return nil
}

func (j jwtHelper) ExtractAudienceToken(token string) (string, error) {
	t, err := j.Manager.ParseToken(token)

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

type JWTManager interface {
	SignToken(data string) (string, error)
	ParseToken(token string) (*jwt.Token, error)
}

type DefaultJWTManager struct{}

func NewDefaultJWTManager() DefaultJWTManager {
	return DefaultJWTManager{}
}

func (m DefaultJWTManager) SignToken(data string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"aud": data,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * (7 * 24)).Unix(),
	})

	return token.SignedString([]byte(configs.GetJWTSecret()))
}

func (m DefaultJWTManager) ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(configs.GetJWTSecret()), nil
	})
}
