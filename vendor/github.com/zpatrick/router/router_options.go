package router

import "net/http"

type RouterOption func(*Router)

func WithNotFound(f func(http.ResponseWriter, *http.Request)) RouterOption {
	return func(r *Router) {
		r.NotFound = f
	}
}
