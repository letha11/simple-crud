package helper_test

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/simple-crud-go/internal/helper"
	mock_helper "github.com/simple-crud-go/internal/helper/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func jwtHelperWithMock(t *testing.T) (helper.JWTHelper, *mock_helper.MockJWTManager) {
	ctrl := gomock.NewController(t)

	jwtManager := mock_helper.NewMockJWTManager(ctrl)

	return helper.NewJWTHelper(jwtManager), jwtManager
}

func TestCheckToken(t *testing.T) {
	jwtHelper, jwtManager := jwtHelperWithMock(t)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
	}{
		{
			"Valid token",
			func() {
				jwtManager.EXPECT().ParseToken(gomock.Any()).Return(&jwt.Token{Valid: true}, nil).Times(1)
			},
			nil,
		},
		{
			"Invalid Token",
			func() {
				jwtManager.EXPECT().ParseToken(gomock.Any()).Return(&jwt.Token{Valid: false}, nil).Times(1)
			},
			jwt.ErrInvalidKey,
		},
		{
			"Parse error",
			func() {
				jwtManager.EXPECT().ParseToken(gomock.Any()).Return(nil, jwt.ErrSignatureInvalid).Times(1)
			},
			jwt.ErrSignatureInvalid,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			err := jwtHelper.CheckToken("asd")

			assert.Equal(t, err, c.err)
		})
	}
}

func TestCreateToken(t *testing.T) {
	jwtHelper, jwtManager := jwtHelperWithMock(t)

	cases := []struct {
		name        string
		mockFunc    func()
		err         error
		signedToken string
	}{
		{
			"Success",
			func() {
				jwtManager.EXPECT().SignToken(gomock.Any()).Return("token", nil).Times(1)
			},
			nil,
			"token",
		},
		{
			"Error",
			func() {
				jwtManager.EXPECT().SignToken(gomock.Any()).Return("", jwt.ErrInvalidKey).Times(1)
			},
			jwt.ErrInvalidKey,
			"",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			signedToken, err := jwtHelper.CreateToken(1)

			assert.Equal(t, signedToken, c.signedToken)
			assert.Equal(t, err, c.err)
		})
	}
}

func TestExtractAudienceToken(t *testing.T) {
	var (
		aud                   = "1"
		jwtHelper, jwtManager = jwtHelperWithMock(t)
	)

	cases := []struct {
		name     string
		mockFunc func()
		err      error
		aud      string
	}{
		{
			"Success",
			func() {
				jwtManager.EXPECT().ParseToken(gomock.Any()).Return(&jwt.Token{
					Claims: jwt.RegisteredClaims{
						Audience: jwt.ClaimStrings{aud},
					},
					Valid: true,
				},
					nil,
				).Times(1)
			},
			nil,
			aud,
		},
		{
			"Invalid Token",
			func() {
				jwtManager.EXPECT().ParseToken(gomock.Any()).Return(&jwt.Token{Valid: false}, nil).Times(1)
			},
			jwt.ErrInvalidKey,
			"",
		},
		{
			"Parse error",
			func() {
				jwtManager.EXPECT().ParseToken(gomock.Any()).Return(nil, jwt.ErrSignatureInvalid).Times(1)
			},
			jwt.ErrSignatureInvalid,
			"",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc()
			aud, err := jwtHelper.ExtractAudienceToken("completelynormaltoken")

			assert.Equal(t, aud, c.aud)
			assert.Equal(t, err, c.err)
		})
	}
}
