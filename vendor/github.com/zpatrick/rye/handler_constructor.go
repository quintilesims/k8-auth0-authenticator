package rye

import "net/http"

type HandlerConstructor func(r *http.Request) http.Handler

func (h HandlerConstructor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(r).ServeHTTP(w, r)
}

func ToHandler(fn func(r *http.Request) http.Handler) http.Handler {
	return HandlerConstructor(fn)
}
