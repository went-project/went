package http

func AuthTemplate() (string, error) {
	return `package middlewares

import (
	"net/http"

	"gorm.io/gorm"
)

func Auth(jwtSecret string, db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authentication logic here
			next.ServeHTTP(w, r)
		})
	}
}
`, nil
}
