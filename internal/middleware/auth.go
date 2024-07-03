package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/simple-crud-go/api"
	"github.com/simple-crud-go/internal/helper"
)

type CtxKey uint

var UserIdKey CtxKey = 0

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtHelper := helper.NewDefaultJWTHelper()

		token := r.Header.Get("Authorization")
		if token == "" {
			api.RequestErrorHandler(w, errors.New("Missing authentication"), http.StatusUnauthorized)
			return
		}

		token = token[len("bearer:"):]

		if err := jwtHelper.CheckToken(token); err != nil {
			if !errors.Is(err, jwt.ErrInvalidKey) {
				api.InternalErrorHandler(w, err)
				return
			}

			api.RequestErrorHandler(w, errors.New("Invalid token"), http.StatusUnauthorized)
			return
		}

		aud, err := jwtHelper.ExtractAudienceToken(token)
		if err != nil {
			api.InternalErrorHandler(w, err)
		}

		r = r.WithContext(context.WithValue(r.Context(), UserIdKey, aud))

		next.ServeHTTP(w, r)
	})
}
