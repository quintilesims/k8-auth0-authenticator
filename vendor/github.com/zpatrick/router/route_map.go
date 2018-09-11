package router

import "net/http"

type MethodHandlers map[string]http.Handler

type RouteMap map[string]MethodHandlers

func (r RouteMap) GlobMatch() []RouteMatcher {
	return r.constructMatchers(NewGlobRouteMatcher)
}

func (r RouteMap) RegexMatch() []RouteMatcher {
	return r.constructMatchers(NewRegexRouteMatcher)
}

func (r RouteMap) StringMatch() []RouteMatcher {
	return r.constructMatchers(NewStringRouteMatcher)
}

func (r RouteMap) VariableMatch() []RouteMatcher {
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
