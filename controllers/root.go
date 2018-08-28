package controllers

import (
	"github.com/zpatrick/fireball"
)

type RootController struct {
	auth0Domain   string
	auth0ClientID string
}

func NewRootController(auth0Domain, auth0ClientID string) *RootController {
	return &RootController{
		auth0Domain:   auth0Domain,
		auth0ClientID: auth0ClientID,
	}
}

func (r *RootController) Routes() []*fireball.Route {
	routes := []*fireball.Route{
		{
			Path: "/",
			Handlers: fireball.Handlers{
				"GET": r.index,
			},
		},
	}

	return routes
}

func (r *RootController) index(c *fireball.Context) (fireball.Response, error) {
	data := struct {
		Auth0Domain   string
		Auth0ClientID string
	}{
		Auth0Domain:   r.auth0Domain,
		Auth0ClientID: r.auth0ClientID,
	}

	return c.HTML(200, "index.html", data)
}
