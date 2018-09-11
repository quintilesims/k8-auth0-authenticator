package router

import "net/http"

type MethodHandlers map[string]http.Handler

type RouteMap map[string]MethodHandlers

func (r RouteMap) Glob() []RouteMatcher {
	return r.constructMatchers(NewGlobRouteMatcher)
}

func (r RouteMap) Regex() []RouteMatcher {
	return r.constructMatchers(NewRegexRouteMatcher)
}

func (r RouteMap) Variable() []RouteMatcher {
	return r.constructMatchers(NewVariableRouteMatcher)
}

func (r RouteMap) constructMatchers(constructor func(string, string, http.Handler) RouteMatcher) []RouteMatcher {
	matchers := []RouteMatcher{}
	for pattern, methodHandlers := range r {
		for method, handler := range methodHandlers {
			matchers = append(matchers, constructor(method, pattern, handler))
		}
	}

	return matchers
}
