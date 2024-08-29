package middleware

import (
	"github.com/go-chi/cors"
	"net/http"
)

func CORSMiddleware() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, //TODO change in prod
		AllowedMethods:   []string{"POST", "OPTIONS", "GET", "PUT"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cache-Control", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
