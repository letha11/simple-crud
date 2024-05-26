package middleware

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func TrailingSlashMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		logrus.Println(r.RequestURI)

		h.ServeHTTP(w, r)
	})
}
