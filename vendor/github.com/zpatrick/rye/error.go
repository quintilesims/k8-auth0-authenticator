package rye

import (
	"net/http"
)

func Error(status int, err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), status)
	})
}
