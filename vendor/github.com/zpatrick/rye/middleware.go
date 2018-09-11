package rye

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("TODO")
		handler.ServeHTTP(w, r)
		log.Printf("TODO")
	})
}
