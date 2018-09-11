package router

import (
	"net/http"
)

type Router struct {
	NotFound func(http.ResponseWriter, *http.Request)
	matchers []RouteMatcher
}

func NewRouter(matchers []RouteMatcher, options ...RouterOption) *Router {
	r := &Router{
		NotFound: http.NotFound,
		matchers: matchers,
	}

	for _, option := range options {
		option(r)
	}

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, match := range r.matchers {
		handler, ok := match(req)
		if ok {
			handler.ServeHTTP(w, req)
			return
		}
	}

	r.NotFound(w, req)
}
