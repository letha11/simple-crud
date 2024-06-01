package middleware

import (
	"net/http"
	"strings"
)

func TrailingSlashMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")

		h.ServeHTTP(w, r)
	})
}
